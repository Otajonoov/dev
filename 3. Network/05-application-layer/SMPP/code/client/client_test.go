package client

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"smpp/pdu"
	"smpp/session"
	"smpp/smsc"
)

func dialTest(t *testing.T, cfg Config) (*Client, *smsc.TestServer) {
	t.Helper()
	srv, err := smsc.StartTestServer()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(srv.Close)
	cfg.Addr = srv.Addr()
	if cfg.SystemID == "" {
		cfg.SystemID = "esme1"
	}
	if cfg.Session.EnquireLink == 0 {
		cfg.Session.EnquireLink = -1 // testda shovqin bo'lmasin
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	c, err := Dial(ctx, cfg)
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		c.Close(ctx)
	})
	return c, srv
}

func testSM(dest string) pdu.SubmitSM {
	return pdu.SubmitSM{SMFields: pdu.SMFields{
		Source:       pdu.Address{TON: pdu.TONAlphanumeric, Addr: "Bank"},
		Dest:         pdu.Address{TON: pdu.TONInternational, NPI: pdu.NPIISDN, Addr: dest},
		ShortMessage: []byte("Salom"),
	}}
}

// To'liq oqim: Dial (sinxron bind) → Submit → message_id.
func TestDialAndSubmit(t *testing.T) {
	c, _ := dialTest(t, Config{})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	id, err := c.Submit(ctx, testSM("998901234567"))
	if err != nil {
		t.Fatalf("Submit: %v", err)
	}
	if !strings.HasPrefix(id, "TST") {
		t.Fatalf("message_id = %q", id)
	}
}

// Server ulanishni QO'POL uzadi → client backoff bilan qayta bind qiladi
// va Submit yana ishlaydi (gosmpp #151 regression'i: rebind tiklanishi).
func TestReconnectAfterDrop(t *testing.T) {
	c, srv := dialTest(t, Config{ReconnectBase: 10 * time.Millisecond, ReconnectMax: 50 * time.Millisecond})
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if _, err := c.Submit(ctx, testSM("998901234567")); err != nil {
		t.Fatalf("birinchi Submit: %v", err)
	}
	srv.DropConnections()

	// Uzilish darhol sezilmaydi — fail-fast xatolardan biri keladi,
	// keyin reconnect tugagach Submit yana o'tadi.
	deadline := time.Now().Add(2 * time.Second)
	for {
		id, err := c.Submit(ctx, testSM("998901234567"))
		if err == nil {
			if !strings.HasPrefix(id, "TST") {
				t.Fatalf("reconnect'dan keyin id=%q", id)
			}
			return
		}
		if !errors.Is(err, ErrNotBound) && !errors.Is(err, session.ErrSessionClosed) {
			t.Fatalf("kutilmagan xato turi: %v", err)
		}
		if time.Now().After(deadline) {
			t.Fatalf("reconnect ulgurmadi: %v", err)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// SubmitLong: 200 belgilik GSM7 matn → 2 segment, HAR BIRI alohida seq va
// UNIQUE message_id (gosmpp #178 regression'i: server id'ni seq'dan yasaydi —
// seq'lar bir xil bo'lsa id'lar ham bir xil bo'lib qolardi).
func TestSubmitLongUniqueSequences(t *testing.T) {
	c, _ := dialTest(t, Config{})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	text := strings.Repeat("A", 200) // GSM7: 160 sig'maydi → 153+47, 2 segment
	ids, err := c.SubmitLong(ctx, LongMessage{
		Source: pdu.Address{TON: pdu.TONAlphanumeric, Addr: "Bank"},
		Dest:   pdu.Address{TON: pdu.TONInternational, NPI: pdu.NPIISDN, Addr: "998901234567"},
		Text:   text,
	})
	if err != nil {
		t.Fatalf("SubmitLong: %v", err)
	}
	if len(ids) != 2 {
		t.Fatalf("%d id keldi, 2 kutilgan", len(ids))
	}
	if ids[0] == ids[1] {
		t.Fatalf("segmentlar BIR XIL seq olgan (#178!): %v", ids)
	}
}

// Qisqa matn — segment YO'Q, bitta oddiy submit.
func TestSubmitLongSingle(t *testing.T) {
	c, _ := dialTest(t, Config{})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	ids, err := c.SubmitLong(ctx, LongMessage{
		Source: pdu.Address{TON: pdu.TONAlphanumeric, Addr: "Bank"},
		Dest:   pdu.Address{TON: pdu.TONInternational, NPI: pdu.NPIISDN, Addr: "998901234567"},
		Text:   "Assalomu alaykum! Kodingiz: 5521",
	})
	if err != nil || len(ids) != 1 {
		t.Fatalf("ids=%v err=%v", ids, err)
	}
}

// Rate limiter Submit'dan OLDIN ishlaydi: 100/s bilan 5 submit >= 40ms.
func TestRateLimiterThrottlesSubmit(t *testing.T) {
	c, _ := dialTest(t, Config{RateLimiter: PerSecond(100)})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	start := time.Now()
	for i := 0; i < 5; i++ {
		if _, err := c.Submit(ctx, testSM("998901234567")); err != nil {
			t.Fatalf("Submit #%d: %v", i, err)
		}
	}
	if el := time.Since(start); el < 40*time.Millisecond {
		t.Fatalf("5 submit %v da o'tdi — limiter ishlamayapti", el)
	}
	// Va ctx bekor qilinsa Wait ham to'xtaydi.
	slow := PerSecond(1)
	slow.Wait(context.Background()) // birinchisi tekin
	ctx2, cancel2 := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel2()
	if err := slow.Wait(ctx2); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("limiter ctx'ni hurmat qilmadi: %v", err)
	}
}

// Close'dan keyin reconnect ham to'xtaydi, Submit xato qaytaradi.
func TestCloseStopsEverything(t *testing.T) {
	c, _ := dialTest(t, Config{})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := c.Close(ctx); err != nil {
		t.Fatalf("Close: %v", err)
	}
	if _, err := c.Submit(ctx, testSM("998901234567")); err == nil {
		t.Fatal("yopiq client'da Submit o'tmasligi kerak")
	}
}
