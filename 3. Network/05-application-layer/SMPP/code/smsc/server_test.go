package smsc

import (
	"context"
	"errors"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"smpp/client"
	"smpp/dlr"
	"smpp/pdu"
	"smpp/session"
)

// dlrCollector — OnDeliver'dan kelgan deliver_sm'larni yig'uvchi.
type dlrCollector struct {
	mu    sync.Mutex
	items []pdu.DeliverSM
	ch    chan pdu.DeliverSM
}

func newCollector() *dlrCollector {
	return &dlrCollector{ch: make(chan pdu.DeliverSM, 16)}
}

func (dc *dlrCollector) onDeliver(d pdu.DeliverSM, h pdu.Header) {
	dc.mu.Lock()
	dc.items = append(dc.items, d)
	dc.mu.Unlock()
	dc.ch <- d
}

func (dc *dlrCollector) wait(t *testing.T) pdu.DeliverSM {
	t.Helper()
	select {
	case d := <-dc.ch:
		return d
	case <-time.After(3 * time.Second):
		t.Fatal("deliver_sm kutish timeout")
		return pdu.DeliverSM{}
	}
}

func startServer(t *testing.T, cfg Config) *Server {
	t.Helper()
	srv, err := Start(cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(srv.Close)
	return srv
}

func dialClient(t *testing.T, srv *Server, dc *dlrCollector) *client.Client {
	t.Helper()
	cfg := client.Config{
		Addr:     srv.Addr(),
		SystemID: "esme1",
		Password: "secret",
		Session:  session.Config{EnquireLink: -1, ResponseTimeout: 2 * time.Second},
	}
	if dc != nil {
		cfg.OnDeliver = dc.onDeliver
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	c, err := client.Dial(ctx, cfg)
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		c.Close(ctx)
	})
	return c
}

func regSubmit(dest string) pdu.SubmitSM {
	return pdu.SubmitSM{SMFields: pdu.SMFields{
		Source:             pdu.Address{TON: pdu.TONAlphanumeric, Addr: "Bank"},
		Dest:               pdu.Address{TON: pdu.TONInternational, NPI: pdu.NPIISDN, Addr: dest},
		RegisteredDelivery: pdu.DLRFinal,
		ShortMessage:       []byte("Salom"),
	}}
}

// To'liq oqim: submit → resp → DLR → parse → korrelyatsiya (normal rejim).
func TestFullFlowWithDLR(t *testing.T) {
	srv := startServer(t, Config{})
	dc := newCollector()
	c := dialClient(t, srv, dc)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	table := dlr.NewTable()
	id, err := c.Submit(ctx, regSubmit("998901234567"))
	if err != nil {
		t.Fatalf("Submit: %v", err)
	}
	table.Register(id)

	d := dc.wait(t)
	if !d.EsmClass.IsDeliveryReceipt() {
		t.Fatalf("esm_class=0x%02X — DLR emas", uint8(d.EsmClass))
	}
	// Manzillar ALMASHGAN bo'lishi kerak (§2.11).
	if d.Source.Addr != "998901234567" || d.Dest.Addr != "Bank" {
		t.Fatalf("manzillar almashmagan: %s -> %s", d.Source.Addr, d.Dest.Addr)
	}
	r, err := dlr.Parse(d.ShortMessage, d.TLVs)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if r.State != dlr.StateDelivered {
		t.Fatalf("State=%v", r.State)
	}
	canon, ok := table.Resolve(r.ID)
	if !ok || canon != id {
		t.Fatalf("korrelyatsiya: %q %v (kutilgan %q)", canon, ok, id)
	}
	// TLV'lar bor rejimda receipted_message_id resp id bilan bir xil.
	if r.Text != "Salom" || r.Sub != 1 || r.Dlvrd != 1 {
		t.Fatalf("receipt: %+v", r)
	}
}

// QUIRK: resp'da hex, DLR matnida decimal — va TLV'siz. Korrelyatsiya faqat
// NormalizeID tufayli ishlaydi (9-bob katta dardi).
func TestQuirkHexIDDecimalDLR(t *testing.T) {
	srv := startServer(t, Config{Quirks: Quirks{HexIDDecimalDLR: true, DLRTextOnly: true}})
	dc := newCollector()
	c := dialClient(t, srv, dc)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	table := dlr.NewTable()
	id, err := c.Submit(ctx, regSubmit("998901234567"))
	if err != nil {
		t.Fatal(err)
	}
	table.Register(id)

	d := dc.wait(t)
	if len(d.TLVs) != 0 {
		t.Fatalf("DLRTextOnly'da TLV bo'lmasligi kerak: %v", d.TLVs)
	}
	r, err := dlr.Parse(d.ShortMessage, d.TLVs)
	if err != nil {
		t.Fatal(err)
	}
	if r.ID == id {
		t.Fatalf("quirk ishlamadi: DLR id resp id bilan bir xil (%q)", id)
	}
	// String tengligi YO'Q, lekin NormalizeID orqali topiladi.
	canon, ok := table.Resolve(r.ID)
	if !ok || canon != id {
		t.Fatalf("hex/dec korrelyatsiya ishlamadi: dlrID=%q respID=%q", r.ID, id)
	}
}

// QUIRK: har 2-chi submit RTHROTTLED (transient — 11-bob tasnifi).
func TestQuirkThrottleEveryN(t *testing.T) {
	srv := startServer(t, Config{Quirks: Quirks{ThrottleEveryN: 2}})
	c := dialClient(t, srv, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if _, err := c.Submit(ctx, regSubmit("998901234567")); err != nil {
		t.Fatalf("1-submit: %v", err)
	}
	_, err := c.Submit(ctx, regSubmit("998901234567"))
	var se client.StatusError
	if !errors.As(err, &se) || se.Status != pdu.StatusRThrottled {
		t.Fatalf("2-submit RTHROTTLED kutilgan edi: %v", err)
	}
	if client.Classify(se.Status) != client.ClassTransient {
		t.Fatal("RTHROTTLED transient bo'lishi kerak")
	}
	if _, err := c.Submit(ctx, regSubmit("998901234567")); err != nil {
		t.Fatalf("3-submit: %v", err)
	}
}

// Auth: noto'g'ri parol → RINVPASWD, noma'lum system_id → RINVSYSID.
func TestAuth(t *testing.T) {
	srv := startServer(t, Config{Accounts: []Account{{SystemID: "esme1", Password: "secret"}}})

	dial := func(sysID, pass string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		c, err := client.Dial(ctx, client.Config{
			Addr: srv.Addr(), SystemID: sysID, Password: pass,
			Session: session.Config{EnquireLink: -1},
		})
		if err == nil {
			c.Close(ctx)
		}
		return err
	}
	if err := dial("esme1", "xato"); err == nil || !strings.Contains(err.Error(), "ESME_RINVPASWD") {
		t.Fatalf("RINVPASWD kutilgan: %v", err)
	}
	if err := dial("hacker", "x"); err == nil || !strings.Contains(err.Error(), "ESME_RINVSYSID") {
		t.Fatalf("RINVSYSID kutilgan: %v", err)
	}
	if err := dial("esme1", "secret"); err != nil {
		t.Fatalf("to'g'ri credential: %v", err)
	}
}

// QUIRK: enquire_link'ka jimlik → client sessiyasi o'zini o'ldiradi (12-bob).
func TestQuirkIgnoreEnquireLink(t *testing.T) {
	srv := startServer(t, Config{Quirks: Quirks{IgnoreEnquireLink: true}})

	conn, err := net.Dial("tcp", srv.Addr())
	if err != nil {
		t.Fatal(err)
	}
	s := session.New(conn, session.Config{
		EnquireLink:     20 * time.Millisecond,
		ResponseTimeout: 30 * time.Millisecond,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if _, err := s.Bind(ctx, pdu.Bind{Mode: pdu.CmdBindTransceiver, SystemID: "esme1"}); err != nil {
		t.Fatal(err)
	}
	select {
	case <-s.Done():
		// sessiya half-open'ni sezib o'zini yopdi — kutilgan xulq
	case <-time.After(2 * time.Second):
		t.Fatal("jim server bilan sessiya o'lishi kerak edi")
	}
}

// Server-tomon state enforcement: BOUND_RX'dan submit_sm → RINVBNDSTS.
// Xom frame'lar bilan — client/session qatlami buni lokal to'sib qo'yardi!
func TestServerStateEnforcement(t *testing.T) {
	srv := startServer(t, Config{})
	conn, err := net.Dial("tcp", srv.Addr())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	write := func(frame []byte, e error) {
		t.Helper()
		if e != nil {
			t.Fatal(e)
		}
		if _, err := conn.Write(frame); err != nil {
			t.Fatal(err)
		}
	}
	read := func() pdu.Header {
		t.Helper()
		conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		frame, err := pdu.ReadFrame(conn, maxPDUSize)
		if err != nil {
			t.Fatal(err)
		}
		h, err := pdu.DecodeHeader(frame)
		if err != nil {
			t.Fatal(err)
		}
		return h
	}

	write(pdu.Bind{Mode: pdu.CmdBindReceiver, SystemID: "esme1"}.Encode(1))
	if h := read(); h.Status != 0 {
		t.Fatalf("bind_receiver: status=%s", pdu.CommandStatus(h.Status))
	}
	write(regSubmit("998901234567").Encode(2))
	h := read()
	if h.ID != pdu.CmdSubmitSMResp || pdu.CommandStatus(h.Status) != pdu.StatusRInvBndSts {
		t.Fatalf("RINVBNDSTS'li submit_sm_resp kutilgan: %s status=%s", h.ID, pdu.CommandStatus(h.Status))
	}
}

// MO injection: abonent xabari (esm_class=0) OnDeliver'ga yetib boradi.
func TestMOInjection(t *testing.T) {
	srv := startServer(t, Config{})
	dc := newCollector()
	dialClient(t, srv, dc)

	err := srv.InjectMO("esme1",
		pdu.Address{TON: pdu.TONInternational, NPI: pdu.NPIISDN, Addr: "998901234567"},
		pdu.Address{TON: pdu.TONUnknown, NPI: pdu.NPIUnknown, Addr: "1234"},
		"STOP")
	if err != nil {
		t.Fatal(err)
	}
	d := dc.wait(t)
	if d.EsmClass.IsDeliveryReceipt() {
		t.Fatal("MO xabar DLR deb belgilanib qolgan")
	}
	if string(d.ShortMessage) != "STOP" || d.Source.Addr != "998901234567" {
		t.Fatalf("MO: %q from %s", d.ShortMessage, d.Source.Addr)
	}
}

// QUIRK: out-of-order javoblar — client window'i baribir to'g'ri bog'laydi.
func TestQuirkOutOfOrderResp(t *testing.T) {
	srv := startServer(t, Config{Quirks: Quirks{OutOfOrderResp: true}})
	c := dialClient(t, srv, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Juftlab almashtiriladi: ikkala Submit PARALLEL ketishi shart
	// (ketma-ket yuborilsa 1-javob 2-submit'gacha ushlab turiladi).
	var wg sync.WaitGroup
	ids := make([]string, 2)
	errs := make([]error, 2)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ids[i], errs[i] = c.Submit(ctx, regSubmit("998901234567"))
		}(i)
	}
	wg.Wait()
	for i, err := range errs {
		if err != nil {
			t.Fatalf("submit %d: %v", i, err)
		}
	}
	if ids[0] == ids[1] || ids[0] == "" {
		t.Fatalf("id'lar: %v", ids)
	}
}

// QUIRK: sekin resp → session ResponseTimeout (11-bob "uchinchi rejim").
func TestQuirkSlowResp(t *testing.T) {
	srv := startServer(t, Config{Quirks: Quirks{SlowResp: 200 * time.Millisecond}})
	cfg := client.Config{
		Addr: srv.Addr(), SystemID: "esme1",
		Session: session.Config{EnquireLink: -1, ResponseTimeout: 50 * time.Millisecond},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	c, err := client.Dial(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		c.Close(ctx)
	}()
	_, err = c.Submit(ctx, regSubmit("998901234567"))
	if !errors.Is(err, session.ErrResponseTimeout) {
		t.Fatalf("ErrResponseTimeout kutilgan: %v", err)
	}
}

// cancel_sm ENROUTE xabarni DELETED qiladi — DLR stat:DELETED keladi;
// query_sm holatni ko'rsatadi. session.Send bilan (client'da bu API'lar yo'q).
func TestCancelAndQueryFlow(t *testing.T) {
	srv := startServer(t, Config{DLRDelay: 300 * time.Millisecond})

	conn, err := net.Dial("tcp", srv.Addr())
	if err != nil {
		t.Fatal(err)
	}
	got := make(chan pdu.DeliverSM, 1)
	s := session.New(conn, session.Config{
		EnquireLink: -1,
		OnInbound: func(r session.Resp) {
			if d, ok := r.PDU.(pdu.DeliverSM); ok {
				got <- d
			}
		},
	})
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if _, err := s.Bind(ctx, pdu.Bind{Mode: pdu.CmdBindTransceiver, SystemID: "esme1"}); err != nil {
		t.Fatal(err)
	}
	defer s.Close(ctx)

	r, err := s.Send(ctx, regSubmit("998901234567"))
	if err != nil {
		t.Fatal(err)
	}
	id := r.PDU.(pdu.SubmitSMResp).MessageID

	// DLR kelmasidan OLDIN cancel.
	r, err = s.Send(ctx, pdu.CancelSM{MessageID: id, Source: pdu.Address{TON: pdu.TONAlphanumeric, Addr: "Bank"}})
	if err != nil || r.Header.Status != 0 {
		t.Fatalf("cancel: err=%v status=%d", err, r.Header.Status)
	}
	// query endi DELETED ko'rsatadi.
	r, err = s.Send(ctx, pdu.QuerySM{MessageID: id, Source: pdu.Address{TON: pdu.TONAlphanumeric, Addr: "Bank"}})
	if err != nil {
		t.Fatal(err)
	}
	if qr := r.PDU.(pdu.QuerySMResp); qr.MessageState != uint8(dlr.StateDeleted) {
		t.Fatalf("query state=%d, DELETED kutilgan", qr.MessageState)
	}
	// DLR baribir keladi — endi stat:DELETED bilan.
	select {
	case d := <-got:
		rc, err := dlr.Parse(d.ShortMessage, d.TLVs)
		if err != nil || rc.State != dlr.StateDeleted {
			t.Fatalf("DLR: %+v err=%v", rc, err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("DELETED DLR kelmadi")
	}
}
