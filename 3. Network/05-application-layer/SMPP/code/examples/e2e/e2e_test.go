package e2e

import (
	"context"
	"sync"
	"testing"
	"time"

	"smpp/client"
	"smpp/coding"
	"smpp/dlr"
	"smpp/pdu"
	"smpp/session"
	"smpp/smsc"
)

// TestEndToEnd — butun kitob kodini birlashtiruvchi stsenariy. Har bosqich
// qaysi bobning mehnati ekani izohlarda.
func TestEndToEnd(t *testing.T) {
	// ── 14-bob: mock SMSC — eng yovuz kombinatsiyada: resp'da HEX id,
	// DLR matnida DECIMAL, va TLV'lar umuman YO'Q.
	srv, err := smsc.Start(smsc.Config{
		Accounts: []smsc.Account{{SystemID: "e2e", Password: "sirli"}},
		DLRDelay: 50 * time.Millisecond,
		Quirks:   smsc.Quirks{HexIDDecimalDLR: true, DLRTextOnly: true},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer srv.Close()

	// ── 9-bob: korrelyatsiya jadvali va biznes-holat.
	table := dlr.NewTable()
	var mu sync.Mutex
	delivered := make(map[string]dlr.MessageState) // kanonik id → final holat
	done := make(chan string, 16)

	// ── 13-bob: client; OnDeliver ichida 9-bob zanjiri.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	c, err := client.Dial(ctx, client.Config{
		Addr:     srv.Addr(),
		SystemID: "e2e",
		Password: "sirli",
		Session:  session.Config{EnquireLink: -1, WindowSize: 10},
		OnDeliver: func(d pdu.DeliverSM, h pdu.Header) {
			if !d.EsmClass.IsDeliveryReceipt() {
				return // MO oqimi bu demoda yo'q
			}
			r, err := dlr.Parse(d.ShortMessage, d.TLVs) // tolerant parser
			if err != nil {
				t.Errorf("DLR parse: %v", err)
				return
			}
			canon, ok := table.Resolve(r.ID) // hex↔dec shu yerda yechiladi!
			if !ok {
				t.Errorf("korrelyatsiya topilmadi: dlr id=%q", r.ID)
				return
			}
			mu.Lock()
			delivered[canon] = r.State
			mu.Unlock()
			if r.State.Final() {
				table.Forget(canon)
			}
			done <- canon
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close(ctx)

	var expected []string // kutilayotgan kanonik id'lar
	submitLong := func(text string) []string {
		t.Helper()
		ids, err := c.SubmitLong(ctx, client.LongMessage{
			Source:             pdu.Address{TON: pdu.TONAlphanumeric, Addr: "Kitob"},
			Dest:               pdu.Address{TON: pdu.TONInternational, NPI: pdu.NPIISDN, Addr: "998901234567"},
			Text:               text,
			RegisteredDelivery: pdu.DLRFinal,
		})
		if err != nil {
			t.Fatalf("SubmitLong(%q...): %v", text[:20], err)
		}
		for _, id := range ids {
			table.Register(id) // 9-bob: resp kelishi bilan DARHOL
			expected = append(expected, id)
		}
		return ids
	}

	// ── 7-bob: o'zbek lotin, U+02BB bilan — Normalize uni ' ga aylantiradi,
	// matn GSM7'da qoladi (aks holda UCS2 70-limitga tushardi).
	uzLatin := "Assalomu alaykum! Yangi oʻzbek xabari: kodingiz 5521"
	if ids := submitLong(uzLatin); len(ids) != 1 {
		t.Fatalf("lotin matn %d segment bo'ldi, 1 kutilgan", len(ids))
	}

	// ── 7-bob: kirill matn — UCS2 (dc=8) bilan ketishi shart.
	uzKirill := "Ассалому алайкум! Кодингиз: 5521"
	if dc, n, err := coding.CountSegments(uzKirill); err != nil || dc != coding.DCUCS2 || n != 1 {
		t.Fatalf("kirill: dc=%v n=%d err=%v", dc, n, err)
	}
	if ids := submitLong(uzKirill); len(ids) != 1 {
		t.Fatalf("kirill matn %d segment", len(ids))
	}

	// ── 8-bob: concatenation — 200+ belgili lotin matn 2 segmentga
	// bo'linadi (UDH 8-bit, har segment alohida seq va alohida message_id).
	long := "Hurmatli mijoz! Sizning soʻrovingiz qabul qilindi. " +
		"Ushbu xabar ataylab uzun yozilmoqda, chunki maqsad concatenation " +
		"mexanizmini oxirigacha tekshirish: UDH, reference raqami, segment " +
		"tartibi va har bir qism uchun alohida delivery receipt."
	ids := submitLong(long)
	if len(ids) != 2 {
		t.Fatalf("uzun matn %d segment bo'ldi, 2 kutilgan", len(ids))
	}
	if ids[0] == ids[1] {
		t.Fatal("segment id'lari bir xil (#178 sinfi regression)")
	}

	// ── 9/12-boblar: barcha DLR'larni yig'amiz (jami 4: 1+1+2).
	for i := 0; i < len(expected); i++ {
		select {
		case <-done:
		case <-ctx.Done():
			t.Fatalf("%d/%d DLR keldi, qolgani yo'q", i, len(expected))
		}
	}
	mu.Lock()
	defer mu.Unlock()
	for _, id := range expected {
		st, ok := delivered[id]
		if !ok {
			t.Errorf("id=%s uchun DLR yo'q", id)
			continue
		}
		if st != dlr.StateDelivered {
			t.Errorf("id=%s holati %v, DELIVERED kutilgan", id, st)
		}
	}
	t.Logf("e2e OK: %d xabar (lotin+kirill+2 segment), hammasi hex/dec quirk ostida DELIVRD deb korrelyatsiya qilindi", len(expected))
}
