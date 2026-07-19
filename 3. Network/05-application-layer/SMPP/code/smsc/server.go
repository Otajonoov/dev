package smsc

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"smpp/dlr"
	"smpp/pdu"
	"smpp/session"
	"smpp/tlv"
)

// Bu fayl — 14-bob: TO'LIQ mock SMSC. TestServer skeletidan (4/13-boblar)
// farqi: auth, server-tomon state machine (Table 2-1), message store,
// DLR generatori, MO injection va — testlar uchun oltin koni — operator
// QUIRK simulyatsiya rejimlari.

// Account — ruxsat etilgan ESME (system_id + parol).
type Account struct {
	SystemID string
	Password string
}

// Quirks — real operatorlarning "g'alati"liklarini ataylab qaytarish
// rejimlari. Har biri kitobda hujjatlangan real og'ish (9/11/12-boblar);
// mock'ning bosh qiymati ideal-spec server emas, AYNAN shu rejimlar.
type Quirks struct {
	// HexIDDecimalDLR: submit_sm_resp'da message_id HEX ko'rinishda,
	// DLR matnining id: field'ida esa DECIMAL (Kannel #334 dunyosi, 9-bob).
	// receipted_message_id TLV'sida esa resp'dagi (hex) qiymat ketadi —
	// real SMSC'larda hujjatlangan kombinatsiya.
	HexIDDecimalDLR bool

	// DLRTextOnly: DLR'da TLV'lar YO'Q — faqat Appendix B matni
	// ("spec bo'yicha SHART" TLV'larni yubormaydigan ko'pchilik, 9-bob).
	DLRTextOnly bool

	// ThrottleEveryN: har N-chi submit_sm RTHROTTLED bilan rad etiladi
	// (rate limit simulyatsiyasi, 11-bob). 0 = o'chiq.
	ThrottleEveryN int

	// SlowResp: submit_sm_resp shu muddatga kechiktiriladi
	// (response timeout trigger'i, 12-bob).
	SlowResp time.Duration

	// IgnoreEnquireLink: enquire_link'ka javob BERILMAYDI
	// (half-open simulyatsiyasi — client o'zini o'ldirishi kerak, 12-bob).
	IgnoreEnquireLink bool

	// OutOfOrderResp: submit javoblari JUFT-JUFT teskari tartibda yuboriladi
	// (§2.5.2 out-of-order ruxsatining stress-testi, 12-bob).
	OutOfOrderResp bool
}

// Config — server sozlamalari.
type Config struct {
	SystemID string    // bind_resp'dagi server nomi; default "MOCKSMSC"
	Accounts []Account // bo'sh = autentifikatsiyasiz (har kim kiradi)

	// DLRDelay — submit qabul qilingandan DLR yuborilguncha kutish
	// (default 10ms; real hayotda soniyalar-soatlar).
	DLRDelay time.Duration

	// SessionInitTimeout — connect → bind orasidagi max vaqt (§7.2:
	// session_init_timer SMSC'da AKTIV bo'lishi kerak). Default 5s.
	SessionInitTimeout time.Duration

	// InactivityTimeout — bound sessiyada shuncha vaqt PDU kelmasa ulanish
	// reset qilinadi (§7.2 inactivity_timer). 0 = o'chiq (test default'i).
	InactivityTimeout time.Duration

	Quirks Quirks
	Logf   func(format string, args ...any)
}

// storedMsg — message store yozuvi: query_sm/cancel_sm va DLR uchun.
type storedMsg struct {
	respID   string // resp'da berilgan id (quirk'da hex)
	dlrID    string // DLR matnida ishlatiladigan id (quirk'da decimal)
	systemID string
	sm       pdu.SMFields
	state    dlr.MessageState
	submitAt time.Time
}

// Server — to'liq mock SMSC.
type Server struct {
	cfg Config
	ln  net.Listener
	wg  sync.WaitGroup

	idCounter atomic.Uint64

	mu    sync.Mutex
	conns map[*serverConn]struct{}
	msgs  map[string]*storedMsg // kalit: respID
}

// Start serverni 127.0.0.1'dagi bo'sh portda ishga tushiradi.
func Start(cfg Config) (*Server, error) {
	if cfg.SystemID == "" {
		cfg.SystemID = "MOCKSMSC"
	}
	if cfg.DLRDelay <= 0 {
		cfg.DLRDelay = 10 * time.Millisecond
	}
	if cfg.SessionInitTimeout <= 0 {
		cfg.SessionInitTimeout = 5 * time.Second
	}
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}
	s := &Server{
		cfg:   cfg,
		ln:    ln,
		conns: make(map[*serverConn]struct{}),
		msgs:  make(map[string]*storedMsg),
	}
	s.wg.Add(1)
	go s.acceptLoop()
	return s, nil
}

// Addr server manzilini qaytaradi.
func (s *Server) Addr() string { return s.ln.Addr().String() }

// Close listener va barcha ulanishlarni yopadi.
func (s *Server) Close() {
	s.ln.Close()
	s.mu.Lock()
	for c := range s.conns {
		c.conn.Close()
	}
	s.mu.Unlock()
	s.wg.Wait()
}

func (s *Server) logf(format string, args ...any) {
	if s.cfg.Logf != nil {
		s.cfg.Logf(format, args...)
	}
}

func (s *Server) acceptLoop() {
	defer s.wg.Done()
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			return
		}
		c := &serverConn{srv: s, conn: conn, state: session.Open}
		s.mu.Lock()
		s.conns[c] = struct{}{}
		s.mu.Unlock()
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			c.serve()
			s.mu.Lock()
			delete(s.conns, c)
			s.mu.Unlock()
		}()
	}
}

// deliverableConn systemID'ga deliver_sm yuborsa bo'ladigan sessiyani topadi
// (BOUND_RX yoki BOUND_TRX — SMPPSim qoidasi: "same system ID"ning RX/TRX
// ulanishlariga yuboriladi).
func (s *Server) deliverableConn(systemID string) *serverConn {
	s.mu.Lock()
	defer s.mu.Unlock()
	for c := range s.conns {
		st, id := c.status()
		if id == systemID && (st == session.BoundRX || st == session.BoundTRX) {
			return c
		}
	}
	return nil
}

// InjectMO systemID'ning RX/TRX sessiyasiga MO xabar (esm_class=0) yuboradi —
// abonentdan kelgan SMS simulyatsiyasi.
func (s *Server) InjectMO(systemID string, from, to pdu.Address, text string) error {
	c := s.deliverableConn(systemID)
	if c == nil {
		return fmt.Errorf("smsc: %q uchun RX/TRX sessiya yo'q", systemID)
	}
	d := pdu.DeliverSM{SMFields: pdu.SMFields{
		Source:       from,
		Dest:         to,
		ShortMessage: []byte(text),
	}}
	frame, err := d.Encode(c.seq.Next())
	if err != nil {
		return err
	}
	return c.write(frame)
}

// newMessageID id juftligini yasaydi: resp'dagi va DLR matnidagi.
// Quirk yoqiq bo'lsa ular BIR SONNING IKKI YOZUVI bo'ladi.
func (s *Server) newMessageID() (respID, dlrID string) {
	n := s.idCounter.Add(1) + 0x4000000000 // "real"ga o'xshagan katta qiymat
	if s.cfg.Quirks.HexIDDecimalDLR {
		return fmt.Sprintf("%X", n), strconv.FormatUint(n, 10)
	}
	id := fmt.Sprintf("%X", n)
	return id, id
}

// serverConn — bitta ESME ulanishi.
type serverConn struct {
	srv  *Server
	conn net.Conn

	writeMu sync.Mutex
	seq     session.Sequencer

	mu       sync.Mutex
	state    session.State
	systemID string

	submitCount int
	heldResp    []byte // OutOfOrderResp quirk'i uchun ushlab turilgan javob
}

func (c *serverConn) status() (session.State, string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.state, c.systemID
}

func (c *serverConn) setBound(st session.State, systemID string) {
	c.mu.Lock()
	c.state = st
	c.systemID = systemID
	c.mu.Unlock()
}

// write to'liq frame'ni bitta Write bilan yozadi (12-bob qoidasi server
// tomonda ham: DLR goroutine'lari bilan resp'lar interleave bo'lmasin).
func (c *serverConn) write(frame []byte) error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	_, err := c.conn.Write(frame)
	return err
}

func (c *serverConn) serve() {
	defer c.conn.Close()
	// session_init_timer (§7.2): bind kelgunicha deadline.
	c.conn.SetReadDeadline(time.Now().Add(c.srv.cfg.SessionInitTimeout))
	for {
		frame, err := pdu.ReadFrame(c.conn, maxPDUSize)
		if err != nil {
			var ne net.Error
			if errors.As(err, &ne) && ne.Timeout() {
				st, _ := c.status()
				c.srv.logf("smsc: %s timer otdi (state=%s) — ulanish yopildi", timerName(st), st)
			}
			return
		}
		p, h, err := pdu.Decode(frame)
		if err != nil {
			// §4.3: notanish id yoki buzuq length → generic_nack.
			status := pdu.StatusRInvCmdID
			if !errors.Is(err, pdu.ErrUnknownCommandID) {
				status = pdu.StatusRInvCmdLen
			}
			c.write(pdu.EncodeGenericNack(uint32(status), h.Sequence))
			continue
		}
		if closeAfter := c.handle(p, h); closeAfter {
			return
		}
		// Inactivity timer: har PDU'dan keyin yangilanadi.
		st, _ := c.status()
		if st != session.Open && c.srv.cfg.InactivityTimeout > 0 {
			c.conn.SetReadDeadline(time.Now().Add(c.srv.cfg.InactivityTimeout))
		} else if st != session.Open {
			c.conn.SetReadDeadline(time.Time{})
		}
	}
}

func timerName(st session.State) string {
	if st == session.Open {
		return "session_init"
	}
	return "inactivity"
}

// handle bitta PDU'ni qayta ishlaydi; true = ulanish yopilsin.
func (c *serverConn) handle(p pdu.PDU, h pdu.Header) bool {
	st, _ := c.status()

	// Server-tomon state enforcement (Table 2-1): noto'g'ri holatdagi
	// PDU'ga RINVBNDSTS. Bind PDU'lari bundan mustasno ko'rib chiqiladi
	// (ular uchun RALYBND aniqroq kod).
	switch p.(type) {
	case pdu.Bind:
	default:
		if !session.CanSend(h.ID, st) {
			c.write(errorResp(h, pdu.StatusRInvBndSts))
			return false
		}
	}

	switch v := p.(type) {
	case pdu.Bind:
		return c.handleBind(v, h)
	case pdu.SubmitSM:
		return c.handleSubmit(v, h)
	case pdu.EnquireLink:
		if c.srv.cfg.Quirks.IgnoreEnquireLink {
			return false // jim — half-open simulyatsiyasi
		}
		c.write(pdu.EncodeEnquireLinkResp(h.Sequence))
	case pdu.Unbind:
		c.write(pdu.EncodeUnbindResp(0, h.Sequence))
		return true
	case pdu.QuerySM:
		c.handleQuery(v, h)
	case pdu.CancelSM:
		c.handleCancel(v, h)
	case pdu.ReplaceSM:
		// Minimal server: replace saqlanmagan — halol RREPLACEFAIL.
		c.write(pdu.ReplaceSMResp{Status: uint32(pdu.StatusRReplaceFail)}.Encode(h.Sequence))
	case pdu.DataSM:
		respID, _ := c.srv.newMessageID()
		if frame, err := (pdu.DataSMResp{MessageID: respID}).Encode(h.Sequence); err == nil {
			c.write(frame)
		}
	case pdu.SubmitMulti:
		respID, _ := c.srv.newMessageID()
		if frame, err := (pdu.SubmitMultiResp{MessageID: respID}).Encode(h.Sequence); err == nil {
			c.write(frame)
		}
	case pdu.DeliverSMResp, pdu.EnquireLinkResp, pdu.GenericNack:
		// Client javoblari — qabul qilinadi, korrelyatsiya talab emas
		// (mock o'z DLR'lariga resp kutmaydi).
	default:
		c.write(pdu.EncodeGenericNack(uint32(pdu.StatusRInvCmdID), h.Sequence))
	}
	return false
}

// errorResp request turiga mos status'li (body'siz) resp yasaydi —
// resp'i bor PDU'larga generic_nack YUBORILMAYDI (11-bob).
func errorResp(h pdu.Header, status pdu.CommandStatus) []byte {
	return encodeHeaderOnly(h.ID.Resp(), uint32(status), h.Sequence)
}

func encodeHeaderOnly(id pdu.CommandID, status, seq uint32) []byte {
	hdr := pdu.EncodeHeader(pdu.Header{Length: pdu.HeaderSize, ID: id, Status: status, Sequence: seq})
	return hdr[:]
}

func (c *serverConn) handleBind(b pdu.Bind, h pdu.Header) bool {
	st, _ := c.status()
	if st != session.Open {
		c.write(errorResp(h, pdu.StatusRAlyBnd))
		return false
	}
	if status := c.srv.auth(b); status != pdu.StatusROK {
		// Status != 0 → body'siz bind_resp (§4.1.2 Note).
		frame, _ := pdu.BindResp{Mode: h.ID.Resp(), Status: uint32(status)}.Encode(h.Sequence)
		c.write(frame)
		return true // muvaffaqiyatsiz auth — ulanish yopiladi
	}
	var newState session.State
	switch b.Mode {
	case pdu.CmdBindTransmitter:
		newState = session.BoundTX
	case pdu.CmdBindReceiver:
		newState = session.BoundRX
	default:
		newState = session.BoundTRX
	}
	c.setBound(newState, b.SystemID)
	frame, err := pdu.BindResp{
		Mode:         h.ID.Resp(),
		SystemID:     c.srv.cfg.SystemID,
		SCVersion:    pdu.InterfaceVersion34,
		HasSCVersion: true,
	}.Encode(h.Sequence)
	if err == nil {
		c.write(frame)
	}
	return false
}

func (s *Server) auth(b pdu.Bind) pdu.CommandStatus {
	if len(s.cfg.Accounts) == 0 {
		return pdu.StatusROK
	}
	for _, a := range s.cfg.Accounts {
		if a.SystemID == b.SystemID {
			if a.Password == b.Password {
				return pdu.StatusROK
			}
			return pdu.StatusRInvPaswd
		}
	}
	return pdu.StatusRInvSysID
}

func (c *serverConn) handleSubmit(sm pdu.SubmitSM, h pdu.Header) bool {
	q := c.srv.cfg.Quirks

	// Throttling quirk'i: har N-chi submit rad (11-bob RTHROTTLED).
	c.mu.Lock()
	c.submitCount++
	throttled := q.ThrottleEveryN > 0 && c.submitCount%q.ThrottleEveryN == 0
	c.mu.Unlock()
	if throttled {
		c.sendSubmitResp(errorResp(h, pdu.StatusRThrottled))
		return false
	}

	respID, dlrID := c.srv.newMessageID()
	_, systemID := c.status()
	msg := &storedMsg{
		respID:   respID,
		dlrID:    dlrID,
		systemID: systemID,
		sm:       sm.SMFields,
		state:    dlr.StateEnroute,
		submitAt: time.Now(),
	}
	c.srv.mu.Lock()
	c.srv.msgs[respID] = msg
	c.srv.mu.Unlock()

	if q.SlowResp > 0 {
		time.Sleep(q.SlowResp) // sekin SMSC simulyatsiyasi (12-bob timeout'i)
	}
	frame, err := pdu.SubmitSMResp{MessageID: respID}.Encode(h.Sequence)
	if err != nil {
		return false
	}
	c.sendSubmitResp(frame)

	// DLR so'ralgan bo'lsa — kechikish bilan generatsiya.
	if sm.RegisteredDelivery.WantsDLR() {
		time.AfterFunc(c.srv.cfg.DLRDelay, func() { c.srv.sendDLR(msg) })
	}
	return false
}

// sendSubmitResp OutOfOrderResp quirk'ini qo'llaydi: javoblar juft-juft
// almashtiriladi — 1-javob ushlab turiladi, 2-si kelganda avval 2, keyin 1
// yuboriladi (§2.5.2: bunga client TAYYOR bo'lishi shart).
func (c *serverConn) sendSubmitResp(frame []byte) {
	if !c.srv.cfg.Quirks.OutOfOrderResp {
		c.write(frame)
		return
	}
	c.mu.Lock()
	held := c.heldResp
	if held == nil {
		c.heldResp = frame
		c.mu.Unlock()
		return
	}
	c.heldResp = nil
	c.mu.Unlock()
	c.write(frame) // avval keyingisi
	c.write(held)  // keyin oldingisi
}

// sendDLR xabar uchun delivery receipt yasab tegishli RX/TRX sessiyaga
// yuboradi. Manzillar ALMASHADI, esm_class=DLR, matn — Appendix B uslubi.
func (s *Server) sendDLR(msg *storedMsg) {
	s.mu.Lock()
	if msg.state == dlr.StateEnroute {
		msg.state = dlr.StateDelivered
	}
	state := msg.state
	s.mu.Unlock()

	c := s.deliverableConn(msg.systemID)
	if c == nil {
		s.logf("smsc: %q uchun DLR yo'naltirib bo'lmadi (RX/TRX yo'q)", msg.systemID)
		return
	}

	now := time.Now().UTC()
	text := msg.sm.ShortMessage
	if len(text) > 20 {
		text = text[:20]
	}
	receipt := fmt.Sprintf(
		"id:%s sub:001 dlvrd:%s submit date:%s done date:%s stat:%s err:000 text:%s",
		msg.dlrID,
		map[bool]string{true: "001", false: "000"}[state == dlr.StateDelivered],
		msg.submitAt.UTC().Format("0601021504"),
		now.Format("0601021504"),
		state.Abbrev(),
		text,
	)
	d := pdu.DeliverSM{SMFields: pdu.SMFields{
		Source:       msg.sm.Dest,   // ALMASHGAN (§2.11)
		Dest:         msg.sm.Source, // ALMASHGAN
		EsmClass:     0x04,          // SMSC Delivery Receipt
		ShortMessage: []byte(receipt),
	}}
	if !s.cfg.Quirks.DLRTextOnly {
		d.TLVs = []tlv.TLV{
			tlv.CString(tlv.ReceiptedMessageID, msg.respID), // resp'dagi ko'rinish
			tlv.U8(tlv.MessageState, uint8(state)),
		}
	}
	frame, err := d.Encode(c.seq.Next())
	if err != nil {
		s.logf("smsc: DLR encode: %v", err)
		return
	}
	if err := c.write(frame); err != nil {
		s.logf("smsc: DLR yozish: %v", err)
	}
}

func (c *serverConn) handleQuery(q pdu.QuerySM, h pdu.Header) {
	c.srv.mu.Lock()
	msg, ok := c.srv.msgs[q.MessageID]
	var state dlr.MessageState
	if ok {
		state = msg.state
	}
	c.srv.mu.Unlock()
	// Matching: message_id + source (§4.8) — soddalashtirilgan: id bo'yicha.
	if !ok {
		c.write(errorResp(h, pdu.StatusRInvMsgID))
		return
	}
	resp := pdu.QuerySMResp{MessageID: q.MessageID, MessageState: uint8(state)}
	if state.Final() {
		// 16 belgili absolute vaqt (§7.1.1): YYMMDDhhmmss + t + nn + p.
		resp.FinalDate = time.Now().UTC().Format("060102150405") + "000+"
	}
	frame, err := resp.Encode(h.Sequence)
	if err != nil {
		c.srv.logf("smsc: query_sm_resp encode: %v", err)
		return
	}
	c.write(frame)
}

func (c *serverConn) handleCancel(cn pdu.CancelSM, h pdu.Header) {
	c.srv.mu.Lock()
	msg, ok := c.srv.msgs[cn.MessageID]
	canceled := false
	if ok && msg.state == dlr.StateEnroute {
		msg.state = dlr.StateDeleted // DLR timer'i endi DELETED yuboradi
		canceled = true
	}
	c.srv.mu.Unlock()
	if !canceled {
		c.write(pdu.CancelSMResp{Status: uint32(pdu.StatusRCancelFail)}.Encode(h.Sequence))
		return
	}
	c.write(pdu.CancelSMResp{}.Encode(h.Sequence))
}
