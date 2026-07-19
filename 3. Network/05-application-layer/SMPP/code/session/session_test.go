package session

import (
	"context"
	"errors"
	"fmt"
	"net"
	"testing"
	"time"

	"smpp/pdu"
)

// testPeer — net.Pipe'ning "SMSC tomoni": testlar undan frame o'qib,
// stsenariy bo'yicha javob yozadi. net.Pipe sinxron va buffersiz —
// deadlock'lar darhol ko'rinadi (aynan shuning uchun testda qimmatli).
type testPeer struct {
	t    *testing.T
	conn net.Conn
}

func (p *testPeer) read() (pdu.PDU, pdu.Header) {
	p.t.Helper()
	p.conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	frame, err := pdu.ReadFrame(p.conn, 64*1024)
	if err != nil {
		p.t.Fatalf("peer read: %v", err)
	}
	pp, h, err := pdu.Decode(frame)
	if err != nil {
		p.t.Fatalf("peer decode: %v", err)
	}
	return pp, h
}

func (p *testPeer) write(frame []byte) {
	p.t.Helper()
	p.conn.SetWriteDeadline(time.Now().Add(2 * time.Second))
	if _, err := p.conn.Write(frame); err != nil {
		p.t.Fatalf("peer write: %v", err)
	}
}

func newPair(t *testing.T, cfg Config) (*Session, *testPeer) {
	t.Helper()
	c, srv := net.Pipe()
	s := New(c, cfg)
	t.Cleanup(func() { s.terminate(nil); srv.Close() })
	return s, &testPeer{t: t, conn: srv}
}

func testSubmit(dest string) pdu.SubmitSM {
	return pdu.SubmitSM{SMFields: pdu.SMFields{
		Dest:         pdu.Address{TON: pdu.TONInternational, NPI: pdu.NPIISDN, Addr: dest},
		ShortMessage: []byte("test"),
	}}
}

// bindTRX sessiyani test uchun bound holatga keltiradi.
func bindTRX(t *testing.T, s *Session, peer *testPeer) {
	t.Helper()
	done := make(chan error, 1)
	go func() {
		_, err := s.Bind(context.Background(), pdu.Bind{Mode: pdu.CmdBindTransceiver, SystemID: "test"})
		done <- err
	}()
	_, h := peer.read()
	frame, err := pdu.BindResp{Mode: pdu.CmdBindTransceiverResp, SystemID: "PEER"}.Encode(h.Sequence)
	if err != nil {
		t.Fatal(err)
	}
	peer.write(frame)
	if err := <-done; err != nil {
		t.Fatalf("bind: %v", err)
	}
	if s.State() != BoundTRX {
		t.Fatalf("state = %s, BOUND_TRX kutilgan", s.State())
	}
}

// Out-of-order javoblar: server ikkinchi submit'ga OLDIN javob beradi —
// korrelyatsiya seq bo'yicha, tartib muhim emas (§2.5.2).
func TestSendOutOfOrderResponses(t *testing.T) {
	s, peer := newPair(t, Config{EnquireLink: -1})
	bindTRX(t, s, peer)

	type result struct {
		resp Resp
		err  error
	}
	res := make(chan result, 2)
	send := func(dest string) {
		r, err := s.Send(context.Background(), testSubmit(dest))
		res <- result{r, err}
	}
	go send("998900000001")
	_, h1 := peer.read()
	go send("998900000002")
	_, h2 := peer.read()

	// Javoblar TESKARI tartibda; message_id = "MSG-<seq>" — kim kimniki
	// ekani shu orqali isbotlanadi.
	for _, h := range []pdu.Header{h2, h1} {
		frame, err := pdu.SubmitSMResp{MessageID: fmt.Sprintf("MSG-%d", h.Sequence)}.Encode(h.Sequence)
		if err != nil {
			t.Fatal(err)
		}
		peer.write(frame)
	}
	for i := 0; i < 2; i++ {
		r := <-res
		if r.err != nil {
			t.Fatalf("Send xato: %v", r.err)
		}
		resp := r.resp.PDU.(pdu.SubmitSMResp)
		want := fmt.Sprintf("MSG-%d", r.resp.Header.Sequence)
		if resp.MessageID != want {
			t.Errorf("korrelyatsiya buzildi: seq=%d uchun %q keldi", r.resp.Header.Sequence, resp.MessageID)
		}
	}
	if s.WindowDepth() != 0 {
		t.Errorf("window bo'shamadi: depth=%d", s.WindowDepth())
	}
}

// Window to'lganda Send BLOKLANADI (ctx'gacha); slot bo'shagach davom etadi.
func TestWindowFullBlocks(t *testing.T) {
	s, peer := newPair(t, Config{WindowSize: 1, EnquireLink: -1, ResponseTimeout: time.Second})
	bindTRX(t, s, peer)

	first := make(chan error, 1)
	go func() {
		_, err := s.Send(context.Background(), testSubmit("998900000001"))
		first <- err
	}()
	_, h1 := peer.read() // server oldi, lekin javob bermay turadi

	// Window=1 to'la: ikkinchi Send ctx timeout'gacha kutib xato oladi.
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	_, err := s.Send(ctx, testSubmit("998900000002"))
	cancel()
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("to'la window'da DeadlineExceeded kutilgan, keldi: %v", err)
	}

	// Birinchisiga javob → slot bo'shaydi → yangi Send o'tadi.
	frame, _ := pdu.SubmitSMResp{MessageID: "A"}.Encode(h1.Sequence)
	peer.write(frame)
	if err := <-first; err != nil {
		t.Fatalf("birinchi Send: %v", err)
	}
	done := make(chan error, 1)
	go func() {
		_, err := s.Send(context.Background(), testSubmit("998900000003"))
		done <- err
	}()
	_, h3 := peer.read()
	frame, _ = pdu.SubmitSMResp{MessageID: "B"}.Encode(h3.Sequence)
	peer.write(frame)
	if err := <-done; err != nil {
		t.Fatalf("uchinchi Send: %v", err)
	}
}

// Javob kelmasa — ErrResponseTimeout (ctx emas!) va window bo'shaydi.
func TestResponseTimeout(t *testing.T) {
	s, peer := newPair(t, Config{ResponseTimeout: 60 * time.Millisecond, EnquireLink: -1})
	bindTRX(t, s, peer)

	done := make(chan error, 1)
	go func() {
		_, err := s.Send(context.Background(), testSubmit("998900000001"))
		done <- err
	}()
	peer.read() // server oldi va JIM

	select {
	case err := <-done:
		if !errors.Is(err, ErrResponseTimeout) {
			t.Fatalf("ErrResponseTimeout kutilgan, keldi: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Send qaytmadi")
	}
	if s.WindowDepth() != 0 {
		t.Errorf("expire'dan keyin depth=%d", s.WindowDepth())
	}
}

// Kelgan enquire_link'ka javob READ PATH'da, handler'siz, darhol.
func TestEnquireLinkAutoResp(t *testing.T) {
	s, peer := newPair(t, Config{EnquireLink: -1})
	bindTRX(t, s, peer)

	peer.write(pdu.EncodeEnquireLink(777))
	p, h := peer.read()
	if _, ok := p.(pdu.EnquireLinkResp); !ok || h.Sequence != 777 {
		t.Fatalf("enquire_link_resp seq=777 kutilgan, keldi %s seq=%d", h.ID, h.Sequence)
	}
}

// Sessiya o'z enquire_link'lariga javob olmasa — o'zini O'LIK deb topadi.
func TestEnquireLinkDeath(t *testing.T) {
	s, peer := newPair(t, Config{EnquireLink: 30 * time.Millisecond, ResponseTimeout: 40 * time.Millisecond})
	bindTRX(t, s, peer)

	// Server enquire_link'ni O'QIYDI, lekin javob bermaydi (half-open
	// simulyatsiyasi: L4 tirik, L7 o'lik).
	peer.read()
	select {
	case <-s.Done():
		if s.Err() == nil {
			t.Fatal("terminate sababi yo'q")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("sessiya o'limi aniqlanmadi")
	}
}

// Inbound queue to'lganda: read loop BLOKLANMAYDI, deliver_sm'ga RX_T_APPN
// qaytadi, reader keyingi PDU'larga xizmat qilishda davom etadi (g.md
// kaskadining teskarisi).
func TestInboundQueueFullNeverBlocksReader(t *testing.T) {
	started := make(chan uint32, 1)
	release := make(chan struct{})
	s, peer := newPair(t, Config{
		EnquireLink:  -1,
		InboundQueue: 1,
		OnInbound: func(r Resp) {
			started <- r.Header.Sequence
			<-release // handler ataylab qotib turadi
		},
	})
	defer close(release)
	bindTRX(t, s, peer)

	dlr := pdu.DeliverSM{SMFields: pdu.SMFields{
		Source:       pdu.Address{TON: 1, NPI: 1, Addr: "998900000001"},
		EsmClass:     0x04,
		ShortMessage: []byte("id:1 stat:DELIVRD"),
	}}
	writeDeliver := func(seq uint32) {
		frame, err := dlr.Encode(seq)
		if err != nil {
			t.Fatal(err)
		}
		peer.write(frame)
	}

	writeDeliver(100)
	// Dispatcher 100'ni olib handler'da qotdi — endi queue bo'sh.
	if seq := <-started; seq != 100 {
		t.Fatalf("handler seq=%d oldi", seq)
	}
	p, h := peer.read() // 100 uchun ack
	if h.Sequence != 100 || h.Status != 0 {
		t.Fatalf("100: %s status=%d", h.ID, h.Status)
	}
	writeDeliver(101) // queue'ga sig'adi (sig'im 1)
	p, h = peer.read()
	if h.Sequence != 101 || h.Status != 0 {
		t.Fatalf("101: status=%d", h.Status)
	}
	writeDeliver(102) // queue TO'LA → RX_T_APPN
	p, h = peer.read()
	if _, ok := p.(pdu.DeliverSMResp); !ok || h.Sequence != 102 {
		t.Fatalf("102: %T seq=%d", p, h.Sequence)
	}
	if h.Status != uint32(pdu.StatusRxTAppn) {
		t.Fatalf("102: status=%s, RX_T_APPN kutilgan", pdu.CommandStatus(h.Status))
	}
	// Reader tirikligining isboti: enquire_link hali ham javob oladi.
	peer.write(pdu.EncodeEnquireLink(103))
	if _, h = peer.read(); h.Sequence != 103 {
		t.Fatalf("reader qotib qoldi: %d", h.Sequence)
	}
}

// Graceful Close: drain → unbind → unbind_resp → yopish; yangi Send'lar rad.
func TestCloseGracefulDrain(t *testing.T) {
	s, peer := newPair(t, Config{EnquireLink: -1})
	bindTRX(t, s, peer)

	sent := make(chan error, 1)
	go func() {
		_, err := s.Send(context.Background(), testSubmit("998900000001"))
		sent <- err
	}()
	_, h := peer.read() // submit yetib keldi, javob hali yo'q

	closed := make(chan error, 1)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		closed <- s.Close(ctx)
	}()
	// Close boshlandi — yangi Send darhol rad etiladi.
	time.Sleep(10 * time.Millisecond)
	if _, err := s.Send(context.Background(), testSubmit("998900000002")); !errors.Is(err, ErrSessionClosed) {
		t.Fatalf("closing paytida ErrSessionClosed kutilgan: %v", err)
	}
	// Server pending submit'ga javob beradi — drain tugaydi.
	frame, _ := pdu.SubmitSMResp{MessageID: "OK"}.Encode(h.Sequence)
	peer.write(frame)
	if err := <-sent; err != nil {
		t.Fatalf("pending Send drain'da xato: %v", err)
	}
	// Endi unbind kelishi kerak.
	p, uh := peer.read()
	if _, ok := p.(pdu.Unbind); !ok {
		t.Fatalf("unbind kutilgan, keldi %s", uh.ID)
	}
	peer.write(pdu.EncodeUnbindResp(0, uh.Sequence))
	if err := <-closed; err != nil {
		t.Fatalf("Close: %v", err)
	}
	if s.State() != Closed {
		t.Fatalf("state=%s", s.State())
	}
}

// Peer unbind yuborsa: resp qaytariladi va sessiya yopiladi (§4.2).
func TestPeerUnbind(t *testing.T) {
	s, peer := newPair(t, Config{EnquireLink: -1})
	bindTRX(t, s, peer)

	peer.write(pdu.EncodeUnbind(55))
	p, h := peer.read()
	if _, ok := p.(pdu.UnbindResp); !ok || h.Sequence != 55 {
		t.Fatalf("unbind_resp seq=55 kutilgan: %s seq=%d", h.ID, h.Sequence)
	}
	select {
	case <-s.Done():
	case <-time.After(time.Second):
		t.Fatal("peer unbind'dan keyin sessiya yopilmadi")
	}
}

// Table 2-1 enforcement: OPEN holatda submit taqiqlangan.
func TestSendStateEnforcement(t *testing.T) {
	s, _ := newPair(t, Config{EnquireLink: -1})
	_, err := s.Send(context.Background(), testSubmit("998900000001"))
	if err == nil {
		t.Fatal("OPEN'da submit_sm o'tmasligi kerak edi")
	}
}

// Sequencer: wrap 0x7FFFFFFF → 1 (0 hech qachon chiqmaydi).
func TestSequencerWrap(t *testing.T) {
	var s Sequencer
	if n := s.Next(); n != 1 {
		t.Fatalf("birinchi Next()=%d", n)
	}
	s.n.Store(maxSequence - 1)
	if n := s.Next(); n != maxSequence {
		t.Fatalf("maxSequence kutilgan: %d", n)
	}
	if n := s.Next(); n != 1 {
		t.Fatalf("wrap'dan keyin 1 kutilgan: %d", n)
	}
}

// Real listener bilan integratsiya testi integration_test.go'da (external
// package session_test) — smsc'ni import qiladi; smsc esa session'ni import
// qilgani uchun in-package testda bu cycle bo'lardi.
