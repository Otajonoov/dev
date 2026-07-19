package session

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"smpp/pdu"
)

// Request — window orqali yuboriladigan har qanday request PDU: o'z
// command_id'sini biladi va berilgan seq bilan to'liq frame yasay oladi.
// pdu package'idagi barcha request turlari (Bind, SubmitSM, QuerySM...)
// bu interfeysga avtomatik mos.
type Request interface {
	Encode(seq uint32) ([]byte, error)
	Cmd() pdu.CommandID
}

// Config — Session sozlamalari. Zero qiymatlar oqilona default'larga
// to'ldiriladi. DIQQAT: timer qiymatlari spec'da YO'Q (§7.2 "outside the
// scope") — bular industriya konventsiyasi, operator TZ'si bilan moslanadi.
type Config struct {
	WindowSize      int           // outstanding limit; default 10 (§2.5.2 Note guideline)
	ResponseTimeout time.Duration // request → resp max kutish; default 10s
	EnquireLink     time.Duration // yuborish intervali; default 30s; manfiy = o'chiq
	InboundQueue    int           // deliver_sm dispatch bufferi; default 64
	MaxPDUSize      uint32        // ReadFrame himoyasi; default 64KB

	// OnInbound — SMSC'dan kelgan request'lar (deliver_sm, data_sm,
	// alert_notification) uchun handler. BITTA dispatcher goroutine'dan
	// ketma-ket chaqiriladi; resp ALLAQACHON yuborilgan bo'ladi (ack sync,
	// processing async — 9-bob qoidasi). nil bo'lsa inbound'lar tashlanadi.
	OnInbound func(Resp)

	// Logf — ixtiyoriy diagnostika (notanish seq, tashlangan PDU...).
	Logf func(format string, args ...any)
}

func (c Config) withDefaults() Config {
	if c.WindowSize <= 0 {
		c.WindowSize = 10
	}
	if c.ResponseTimeout <= 0 {
		c.ResponseTimeout = 10 * time.Second
	}
	if c.EnquireLink == 0 {
		c.EnquireLink = 30 * time.Second
	}
	if c.InboundQueue <= 0 {
		c.InboundQueue = 64
	}
	if c.MaxPDUSize == 0 {
		c.MaxPDUSize = 64 * 1024
	}
	return c
}

// Session — bitta TCP ulanish ustidagi SMPP runtime: reader goroutine,
// pending window, dispatcher, enquire_link va expire timer'lari.
// Goroutine modeli (12-bob):
//
//	reader   — FAQAT ReadFrame + route (hech qachon bloklanmaydi!)
//	dispatch — OnInbound'ni ketma-ket chaqiradi (bounded queue ortida)
//	enquire  — davriy enquire_link (javobsizlik = sessiya o'limi)
//	expire   — window deadline scanner (map leak'ka qarshi)
//
// Yozish alohida goroutine EMAS — writeMu bilan himoyalangan to'liq-frame
// Write: bitta PDU = bitta Write chaqiruvi, baytlar interleave bo'lmaydi.
type Session struct {
	conn net.Conn
	cfg  Config
	seq  Sequencer
	win  *window

	writeMu sync.Mutex

	stateMu sync.Mutex
	state   State

	closing atomic.Bool // graceful Close boshlandi — yangi Send'lar rad
	inbound chan Resp
	done    chan struct{}

	termOnce sync.Once
	termErr  error
}

// New sessiyani net.Conn ustida ishga tushiradi (state=OPEN — bind hali
// yo'q; Bind chaqiriladi). conn odatda *net.TCPConn: Go'da TCP_NODELAY
// DEFAULT yoqiq (net paketi SetNoDelay(true)'ni o'zi qiladi) — SMPP'ning
// kichik PDU'lari uchun ayni kerakli rejim; uni O'CHIRMANG.
func New(conn net.Conn, cfg Config) *Session {
	cfg = cfg.withDefaults()
	s := &Session{
		conn:    conn,
		cfg:     cfg,
		win:     newWindow(cfg.WindowSize),
		state:   Open,
		inbound: make(chan Resp, cfg.InboundQueue),
		done:    make(chan struct{}),
	}
	go s.dispatchLoop()
	go s.readLoop()
	go s.expireLoop()
	go s.enquireLoop()
	return s
}

// State joriy session holatini qaytaradi.
func (s *Session) State() State {
	s.stateMu.Lock()
	defer s.stateMu.Unlock()
	return s.state
}

func (s *Session) setState(st State) {
	s.stateMu.Lock()
	s.state = st
	s.stateMu.Unlock()
}

// Done sessiya tugaganda yopiladigan kanal (reconnect signali — 13-bob).
func (s *Session) Done() <-chan struct{} { return s.done }

// Err sessiya tugash sababini qaytaradi (Done yopilgach ma'noli).
func (s *Session) Err() error {
	<-s.done
	return s.termErr
}

// WindowDepth — hozirgi outstanding soni (monitoring, 16-bob).
func (s *Session) WindowDepth() int { return s.win.depth() }

func (s *Session) logf(format string, args ...any) {
	if s.cfg.Logf != nil {
		s.cfg.Logf(format, args...)
	}
}

// Send request'ni yuborib javobini kutadi (sync-over-async: ichkarida
// window). Har chaqiruv alohida goroutine'dan bo'lishi mumkin — parallellik
// window bilan cheklanadi: window to'la bo'lsa ctx bekor bo'lguncha kutadi
// (backpressure chaqiruvchiga oqib chiqadi).
func (s *Session) Send(ctx context.Context, req Request) (Resp, error) {
	if s.closing.Load() {
		return Resp{}, ErrSessionClosed
	}
	if st := s.State(); !CanSend(req.Cmd(), st) {
		return Resp{}, fmt.Errorf("session: %s holatida %s yuborib bo'lmaydi (Table 2-1)", st, req.Cmd())
	}
	return s.send(ctx, req.Cmd(), req.Encode)
}

// send — ichki yo'l: state/closing tekshiruvisiz (unbind va enquire_link
// shu yo'ldan yuradi).
func (s *Session) send(ctx context.Context, cmd pdu.CommandID, encode func(uint32) ([]byte, error)) (Resp, error) {
	seq := s.seq.Next()
	frame, err := encode(seq)
	if err != nil {
		return Resp{}, err
	}
	if err := s.win.acquire(ctx); err != nil {
		return Resp{}, err
	}
	p := s.win.add(seq, cmd, time.Now().Add(s.cfg.ResponseTimeout))
	if err := s.writeFrame(frame); err != nil {
		s.win.fail(seq, err)
	}
	select {
	case r := <-p.ch:
		return r.resp, r.err
	case <-ctx.Done():
		// Poyga: resolve allaqachon yozgan bo'lishi mumkin. fail entry
		// mavjud bo'lsagina yozadi — har ikki holda ch'da AYNAN bitta
		// natija bor.
		s.win.fail(seq, ctx.Err())
		r := <-p.ch
		if r.err == nil {
			return r.resp, nil // javob ctx bilan bir vaqtda yetib keldi
		}
		return Resp{}, r.err
	}
}

// writeFrame to'liq frame'ni BITTA Write bilan yozadi (interleave taqiqlangan).
// Yozish xatosi ulanish o'limi — sessiya terminate qilinadi.
func (s *Session) writeFrame(frame []byte) error {
	s.writeMu.Lock()
	_, err := s.conn.Write(frame)
	s.writeMu.Unlock()
	if err != nil {
		s.terminate(fmt.Errorf("session: yozish xatosi: %w", err))
	}
	return err
}

// Bind bind PDU'sini yuborib, muvaffaqiyatda state'ni yangilaydi.
func (s *Session) Bind(ctx context.Context, b pdu.Bind) (pdu.BindResp, error) {
	r, err := s.Send(ctx, b)
	if err != nil {
		return pdu.BindResp{}, err
	}
	br, ok := r.PDU.(pdu.BindResp)
	if !ok {
		return pdu.BindResp{}, fmt.Errorf("session: bind'ga %s keldi", r.Header.ID)
	}
	if br.Status != 0 {
		return br, fmt.Errorf("session: bind rad etildi: %s", pdu.CommandStatus(br.Status))
	}
	switch b.Mode {
	case pdu.CmdBindTransmitter:
		s.setState(BoundTX)
	case pdu.CmdBindReceiver:
		s.setState(BoundRX)
	case pdu.CmdBindTransceiver:
		s.setState(BoundTRX)
	}
	return br, nil
}

// Close — graceful shutdown (§4.2 tartibi): yangi submit STOP → window
// drain (ctx chegarasida) → unbind → unbind_resp → TCP close. ctx tugasa
// ham sessiya YOPILADI (majburan), faqat xato qaytadi.
func (s *Session) Close(ctx context.Context) error {
	if !s.closing.CompareAndSwap(false, true) {
		<-s.done
		return nil // allaqachon yopilmoqda/yopilgan
	}
	var retErr error
	// 1. Drain: outstanding'lar javobini kutamiz.
	for s.win.depth() > 0 {
		select {
		case <-ctx.Done():
			retErr = fmt.Errorf("session: drain tugallanmadi: %w", ctx.Err())
		case <-s.done:
			return s.termErr
		case <-time.After(5 * time.Millisecond):
			continue
		}
		break
	}
	// 2. Unbind (faqat bound holatda ma'noli).
	if retErr == nil && s.State() != Open && s.State() != Closed {
		if _, err := s.send(ctx, pdu.CmdUnbind, func(seq uint32) ([]byte, error) {
			return pdu.EncodeUnbind(seq), nil
		}); err != nil && !errors.Is(err, ErrSessionClosed) {
			retErr = fmt.Errorf("session: unbind: %w", err)
		}
	}
	// 3. Yakuniy yopish.
	s.terminate(nil)
	return retErr
}

// terminate sessiyani qat'iy yakunlaydi: state=CLOSED, conn yopiladi,
// barcha pending'lar xato bilan qaytadi. Rebind'da eski window YANGI
// ulanishga ko'chmaydi — resend avtomatik EMAS (duplicate xavfi, 5-bob
// dilemmasi caller'ga qoldiriladi).
func (s *Session) terminate(err error) {
	s.termOnce.Do(func() {
		s.termErr = err
		s.setState(Closed)
		s.conn.Close()
		s.win.failAll(ErrSessionClosed)
		close(s.done)
	})
}

// readLoop — sessiyaning yuragi. QOIDA: bu loop HECH QACHON bloklanmaydi
// (g.md kaskadi: read to'xtasa enquire_link_resp o'qilmaydi → 30s jimlik →
// ulanish o'limi). Shuning uchun inbound enqueue non-blocking, handler'lar
// alohida goroutine'da.
func (s *Session) readLoop() {
	for {
		frame, err := pdu.ReadFrame(s.conn, s.cfg.MaxPDUSize)
		if err != nil {
			s.terminate(fmt.Errorf("session: o'qish xatosi: %w", err))
			return
		}
		p, h, err := pdu.Decode(frame)
		if err != nil {
			if errors.Is(err, pdu.ErrUnknownCommandID) {
				// §4.3: notanish command_id → generic_nack (seq bor).
				s.writeFrame(pdu.EncodeGenericNack(uint32(pdu.StatusRInvCmdID), h.Sequence))
				continue
			}
			// Body decode xatosi: request bo'lsa to'g'ri javob server
			// ishi (14-bob); client sifatida log + davom.
			s.logf("session: decode xatosi: %v", err)
			continue
		}
		switch p.(type) {
		case pdu.EnquireLink:
			// Read path'da DARHOL — queue'siz (half-open davosi).
			s.writeFrame(pdu.EncodeEnquireLinkResp(h.Sequence))
		case pdu.Unbind:
			// Peer-initiated unbind: resp → yopish (§4.2).
			s.writeFrame(pdu.EncodeUnbindResp(0, h.Sequence))
			s.terminate(fmt.Errorf("session: peer unbind qildi"))
			return
		case pdu.GenericNack:
			if h.Sequence != 0 && s.win.resolve(h.Sequence, Resp{PDU: p, Header: h}) {
				continue
			}
			// seq=0 yoki korrelyatsiyasiz nack — framing shubhasi (11-bob).
			s.logf("session: korrelyatsiyasiz generic_nack status=%s seq=%d", pdu.CommandStatus(h.Status), h.Sequence)
		default:
			if h.ID.IsResponse() {
				if !s.win.resolve(h.Sequence, Resp{PDU: p, Header: h}) {
					// Kechikkan (expire bo'lgan) yoki duplicate javob.
					// Resp'ga nack YUBORILMAYDI (11-bob).
					s.logf("session: notanish seq=%d (%s) — e'tiborsiz", h.Sequence, h.ID)
				}
				continue
			}
			s.handleInboundRequest(p, h)
		}
	}
}

// handleInboundRequest SMSC'dan kelgan request'ni qayta ishlaydi:
// AVVAL bounded queue'ga joylash urinishi, KEYIN ack — teskari tartibda
// "ack qildik-u tashladik" (jimgina yo'qolgan DLR) chiqadi.
func (s *Session) handleInboundRequest(p pdu.PDU, h pdu.Header) {
	enqueued := false
	select {
	case s.inbound <- Resp{PDU: p, Header: h}:
		enqueued = true
	default:
		// Queue to'la. Read loop BLOKLANMAYDI — buning o'rniga SMSC'ga
		// halol signal: RX_T_APPN = "vaqtincha ololmayman, qayta urin"
		// (11-bob). Bu g.md kaskadining kodlashtirilgan davosi.
		s.logf("session: inbound queue to'la — %s tashlab yuborildi", h.ID)
	}
	switch p.(type) {
	case pdu.DeliverSM:
		status := uint32(0)
		if !enqueued {
			status = uint32(pdu.StatusRxTAppn)
		}
		s.writeFrame(pdu.DeliverSMResp{Status: status}.Encode(h.Sequence))
	case pdu.DataSM:
		var frame []byte
		var err error
		if enqueued {
			frame, err = pdu.DataSMResp{}.Encode(h.Sequence)
		} else {
			frame, err = pdu.DataSMResp{Status: uint32(pdu.StatusRxTAppn)}.Encode(h.Sequence)
		}
		if err == nil {
			s.writeFrame(frame)
		}
	case pdu.AlertNotification:
		// Resp YO'Q (§4.12) — enqueue bo'lmasa shunchaki yo'qoladi.
	default:
		// ESME sessiyasida kutilmagan request (masalan submit_sm) —
		// to'g'ri enforcement server ishi (14-bob); bu yerda faqat log.
		s.logf("session: kutilmagan request %s — e'tiborsiz", h.ID)
	}
}

// dispatchLoop — OnInbound'ni KETMA-KET chaqiradigan yagona goroutine.
// Handler sekin bo'lsa queue to'ladi va handleInboundRequest RX_T_APPN
// qaytara boshlaydi — backpressure protokol darajasiga halol chiqadi.
func (s *Session) dispatchLoop() {
	for {
		select {
		case r := <-s.inbound:
			if s.cfg.OnInbound != nil {
				s.cfg.OnInbound(r)
			}
		case <-s.done:
			return
		}
	}
}

// expireLoop — javobsiz qolgan pending'larni davriy tozalaydi.
func (s *Session) expireLoop() {
	interval := s.cfg.ResponseTimeout / 4
	if interval < 5*time.Millisecond {
		interval = 5 * time.Millisecond
	}
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case now := <-t.C:
			if n := s.win.expire(now); n > 0 {
				s.logf("session: %d ta request expire bo'ldi", n)
			}
		case <-s.done:
			return
		}
	}
}

// enquireLoop — davriy enquire_link. Javob kelmasa sessiya O'LIK deb
// topiladi (half-open connection davosi): TCP keepalive L4'da stack
// tirikligini tekshiradi, enquire_link esa L7'da SMPP application'ning
// o'zini — server process hang bo'lsa farqi aynan shu yerda ko'rinadi.
func (s *Session) enquireLoop() {
	if s.cfg.EnquireLink < 0 {
		return
	}
	t := time.NewTicker(s.cfg.EnquireLink)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			st := s.State()
			if st != BoundTX && st != BoundRX && st != BoundTRX {
				continue // bind'dan oldin/keyin yuborilmaydi (Table 2-1)
			}
			ctx, cancel := context.WithTimeout(context.Background(), s.cfg.ResponseTimeout)
			_, err := s.send(ctx, pdu.CmdEnquireLink, func(seq uint32) ([]byte, error) {
				return pdu.EncodeEnquireLink(seq), nil
			})
			cancel()
			if err != nil && !errors.Is(err, ErrSessionClosed) {
				s.terminate(fmt.Errorf("session: enquire_link javobsiz: %w", err))
				return
			}
		case <-s.done:
			return
		}
	}
}
