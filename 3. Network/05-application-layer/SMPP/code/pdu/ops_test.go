package pdu

import (
	"errors"
	"reflect"
	"testing"

	"smpp/tlv"
)

// 10-bob golden'lari — barchasi python struct bilan mustaqil yasalgan
// (kitobdagi hex dump'lar bilan bir xil baytlar).

// data_sm: Bank(5/0) → 998901234567(1/1), matn message_payload TLV'da,
// DLR so'ralgan; seq=4. 51 (0x33) oktet.
const goldenDataSMHex = "0000003300000103000000000000000400050042616E6B000101393938393031323334353637000001000424000553616C6F6D"

// data_sm_resp: message_id "7F3A9C"; seq=4.
const goldenDataSMRespHex = "0000001780000103000000000000000437463341394300"

// submit_multi: Bank(5/0) → 2 SME + 1 DL ("vip_mijozlar"), "Salom",
// registered_delivery=0x01; seq=7. 86 (0x56) oktet.
const goldenSubmitMultiHex = "0000005600000021000000000000000700050042616E6B00030101013939383930313233343536370001010139393839303736353433323100027669705F6D696A6F7A6C6172000000000000010000000553616C6F6D"

// submit_multi_resp: message_id "7F3AA0", no_unsuccess=1
// (998907654321 → 0x0B RINVDSTADR); seq=7. 43 (0x2B) oktet.
const goldenSubmitMultiRespHex = "0000002B80000021000000000000000737463341413000010101393938393037363534333231000000000B"

// query_sm: message_id "7F3A9B" (5-bob submit'iniki!), source Bank(5/0); seq=8.
const goldenQuerySMHex = "0000001E00000003000000000000000837463341394200050042616E6B00"

// query_sm_resp: DELIVERED (2), final_date "260717120600004+", err=0; seq=8.
const goldenQuerySMRespHex = "0000002A800000030000000000000008374633413942003236303731373132303630303030342B000200"

func TestDataSMGolden(t *testing.T) {
	want := mustHex(t, goldenDataSMHex)
	d := DataSM{
		Source:             Address{TON: TONAlphanumeric, NPI: NPIUnknown, Addr: "Bank"},
		Dest:               Address{TON: TONInternational, NPI: NPIISDN, Addr: "998901234567"},
		RegisteredDelivery: DLRFinal,
		TLVs:               []tlv.TLV{{Tag: tlv.MessagePayload, Value: []byte("Salom")}},
	}
	got, err := d.Encode(4)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("encode:\n got  %X\n want %X", got, want)
	}
	back, h, err := DecodeDataSM(want)
	if err != nil {
		t.Fatal(err)
	}
	if h.Sequence != 4 || !reflect.DeepEqual(back, d) {
		t.Fatalf("decode: %+v (seq=%d)", back, h.Sequence)
	}
	// Matn message_payload'da ekanini tipli helper bilan olish.
	mp, ok := tlv.Find(back.TLVs, tlv.MessagePayload)
	if !ok || string(mp.Value) != "Salom" {
		t.Fatalf("message_payload: %+v", mp)
	}
}

func TestDataSMRespGolden(t *testing.T) {
	want := mustHex(t, goldenDataSMRespHex)
	got, err := DataSMResp{MessageID: "7F3A9C"}.Encode(4)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("encode:\n got  %X\n want %X", got, want)
	}
	back, _, err := DecodeDataSMResp(want)
	if err != nil || back.MessageID != "7F3A9C" {
		t.Fatalf("decode: %+v err=%v", back, err)
	}
	// Status != 0 → body yo'q (submit_sm_resp qoidasi data_sm_resp'da ham).
	failFrame, err := DataSMResp{Status: 0x58, MessageID: "IGNORED"}.Encode(5)
	if err != nil {
		t.Fatal(err)
	}
	if len(failFrame) != HeaderSize {
		t.Fatalf("status!=0 frame %d oktet, faqat header kutilgan", len(failFrame))
	}
}

func TestSubmitMultiGolden(t *testing.T) {
	want := mustHex(t, goldenSubmitMultiHex)
	s := SubmitMulti{
		Source: Address{TON: TONAlphanumeric, NPI: NPIUnknown, Addr: "Bank"},
		Dests: []DestAddress{
			{SME: Address{TON: TONInternational, NPI: NPIISDN, Addr: "998901234567"}},
			{SME: Address{TON: TONInternational, NPI: NPIISDN, Addr: "998907654321"}},
			{DLName: "vip_mijozlar"},
		},
		RegisteredDelivery: DLRFinal,
		ShortMessage:       []byte("Salom"),
	}
	got, err := s.Encode(7)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("encode:\n got  %X\n want %X", got, want)
	}
	back, h, err := DecodeSubmitMulti(want)
	if err != nil {
		t.Fatal(err)
	}
	if h.Sequence != 7 || !reflect.DeepEqual(back, s) {
		t.Fatalf("decode: %+v", back)
	}
	if !back.Dests[2].IsDistList() || back.Dests[0].IsDistList() {
		t.Fatal("dest_flag union noto'g'ri tiklandi")
	}
}

func TestSubmitMultiRespGolden(t *testing.T) {
	want := mustHex(t, goldenSubmitMultiRespHex)
	r := SubmitMultiResp{
		MessageID: "7F3AA0",
		Unsuccess: []UnsuccessSME{{
			Addr:            Address{TON: TONInternational, NPI: NPIISDN, Addr: "998907654321"},
			ErrorStatusCode: 0x0B, // ESME_RINVDSTADR
		}},
	}
	got, err := r.Encode(7)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("encode:\n got  %X\n want %X", got, want)
	}
	back, _, err := DecodeSubmitMultiResp(want)
	if err != nil || !reflect.DeepEqual(back, r) {
		t.Fatalf("decode: %+v err=%v", back, err)
	}
	// Umumiy status=0, lekin 1 manzil muvaffaqiyatsiz — "qisman muvaffaqiyat".
	if back.Status != 0 || len(back.Unsuccess) != 1 {
		t.Fatalf("qisman muvaffaqiyat semantikasi: %+v", back)
	}
}

func TestQuerySMGolden(t *testing.T) {
	want := mustHex(t, goldenQuerySMHex)
	q := QuerySM{
		MessageID: "7F3A9B",
		Source:    Address{TON: TONAlphanumeric, NPI: NPIUnknown, Addr: "Bank"},
	}
	got, err := q.Encode(8)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("encode:\n got  %X\n want %X", got, want)
	}
	back, _, err := DecodeQuerySM(want)
	if err != nil || !reflect.DeepEqual(back, q) {
		t.Fatalf("decode: %+v err=%v", back, err)
	}
}

func TestQuerySMRespGolden(t *testing.T) {
	want := mustHex(t, goldenQuerySMRespHex)
	r := QuerySMResp{
		MessageID:    "7F3A9B",
		FinalDate:    "260717120600004+",
		MessageState: 2, // DELIVERED
	}
	got, err := r.Encode(8)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("encode:\n got  %X\n want %X", got, want)
	}
	back, _, err := DecodeQuerySMResp(want)
	if err != nil || !reflect.DeepEqual(back, r) {
		t.Fatalf("decode: %+v err=%v", back, err)
	}
}

func TestCancelSMRoundTrip(t *testing.T) {
	// Rejim 1: bitta xabar message_id bo'yicha.
	one := CancelSM{
		MessageID: "7F3A9B",
		Source:    Address{TON: TONAlphanumeric, NPI: NPIUnknown, Addr: "Bank"},
	}
	frame, err := one.Encode(9)
	if err != nil {
		t.Fatal(err)
	}
	back, h, err := DecodeCancelSM(frame)
	if err != nil || h.Sequence != 9 || !reflect.DeepEqual(back, one) {
		t.Fatalf("rejim 1: %+v err=%v", back, err)
	}
	// Rejim 2: guruh — message_id NULL, source+dest bo'yicha.
	group := CancelSM{
		ServiceType: "CMT",
		Source:      Address{TON: TONAlphanumeric, NPI: NPIUnknown, Addr: "Bank"},
		Dest:        Address{TON: TONInternational, NPI: NPIISDN, Addr: "998901234567"},
	}
	frame, err = group.Encode(10)
	if err != nil {
		t.Fatal(err)
	}
	back, _, err = DecodeCancelSM(frame)
	if err != nil || !reflect.DeepEqual(back, group) {
		t.Fatalf("rejim 2: %+v err=%v", back, err)
	}
	// Hech qaysi rejimga tushmaydi: ikkalasi ham bo'sh.
	if _, err := (CancelSM{Source: one.Source}).Encode(11); err == nil {
		t.Fatal("message_id ham dest ham bo'sh — xato kutilgan edi")
	}
	// Resp round-trip.
	rframe := CancelSMResp{Status: 0x11}.Encode(9) // RCANCELFAIL
	rback, _, err := DecodeCancelSMResp(rframe)
	if err != nil || rback.Status != 0x11 {
		t.Fatalf("resp: %+v err=%v", rback, err)
	}
}

func TestReplaceSMRoundTrip(t *testing.T) {
	p := ReplaceSM{
		MessageID:          "7F3A9B",
		Source:             Address{TON: TONAlphanumeric, NPI: NPIUnknown, Addr: "Bank"},
		ValidityPeriod:     "000001000000000R",
		RegisteredDelivery: DLRFinal,
		ShortMessage:       []byte("Yangi matn"),
	}
	frame, err := p.Encode(12)
	if err != nil {
		t.Fatal(err)
	}
	back, h, err := DecodeReplaceSM(frame)
	if err != nil || h.Sequence != 12 || !reflect.DeepEqual(back, p) {
		t.Fatalf("decode: %+v err=%v", back, err)
	}
	rframe := ReplaceSMResp{Status: 0x13}.Encode(12) // RREPLACEFAIL
	rback, _, err := DecodeReplaceSMResp(rframe)
	if err != nil || rback.Status != 0x13 {
		t.Fatalf("resp: %+v err=%v", rback, err)
	}
}

func TestAlertNotificationRoundTrip(t *testing.T) {
	a := AlertNotification{
		Source:   Address{TON: TONInternational, NPI: NPIISDN, Addr: "998901234567"},
		ESMEAddr: Address{TON: TONInternational, NPI: NPIISDN, Addr: "170"},
		TLVs:     []tlv.TLV{tlv.U8(tlv.MsAvailabilityStatus, 0)}, // 0 = available
	}
	frame, err := a.Encode(13)
	if err != nil {
		t.Fatal(err)
	}
	back, h, err := DecodeAlertNotification(frame)
	if err != nil || h.Sequence != 13 || !reflect.DeepEqual(back, a) {
		t.Fatalf("decode: %+v err=%v", back, err)
	}
}

// Dispatcher: barcha turlarni taniydi, notanish id'ga sentinel xato.
func TestDecodeDispatcher(t *testing.T) {
	frames := map[CommandID][]byte{}
	add := func(f []byte, err error) {
		if err != nil {
			t.Fatal(err)
		}
		h, _ := DecodeHeader(f)
		frames[h.ID] = f
	}
	add(Bind{Mode: CmdBindTransceiver, SystemID: "esme1"}.Encode(1))
	add(BindResp{Mode: CmdBindTransceiverResp, SystemID: "SMSC"}.Encode(1))
	add(Outbind{SystemID: "smsc1"}.Encode(1))
	add(SubmitSM{SMFields: SMFields{Dest: Address{TON: 1, NPI: 1, Addr: "998901234567"}, ShortMessage: []byte("x")}}.Encode(2))
	add(SubmitSMResp{MessageID: "7F3A9B"}.Encode(2))
	add(DeliverSM{SMFields: SMFields{Source: Address{TON: 1, NPI: 1, Addr: "998901234567"}}}.Encode(3))
	add(DeliverSMResp{}.Encode(3), nil)
	add(mustHex(t, goldenDataSMHex), nil)
	add(mustHex(t, goldenDataSMRespHex), nil)
	add(mustHex(t, goldenQuerySMHex), nil)
	add(mustHex(t, goldenQuerySMRespHex), nil)
	add(CancelSM{MessageID: "1"}.Encode(5))
	add(CancelSMResp{}.Encode(5), nil)
	add(ReplaceSM{MessageID: "1"}.Encode(6))
	add(ReplaceSMResp{}.Encode(6), nil)
	add(mustHex(t, goldenSubmitMultiHex), nil)
	add(mustHex(t, goldenSubmitMultiRespHex), nil)
	add(AlertNotification{}.Encode(7))
	add(EncodeEnquireLink(8), nil)
	add(EncodeEnquireLinkResp(8), nil)
	add(EncodeUnbind(9), nil)
	add(EncodeUnbindResp(0, 9), nil)
	add(EncodeGenericNack(0x03, 0), nil)

	if len(frames) != 23 {
		t.Fatalf("frame to'plami %d — 23 kutilgan", len(frames))
	}
	for id, frame := range frames {
		p, h, err := Decode(frame)
		if err != nil {
			t.Errorf("%s: Decode xato: %v", id, err)
			continue
		}
		if p.Cmd() != id || h.ID != id {
			t.Errorf("%s: Cmd()=%s, header=%s", id, p.Cmd(), h.ID)
		}
	}

	// Notanish command_id → ErrUnknownCommandID (caller generic_nack yuboradi).
	bogus := encodePDU(CommandID(0x000000AA), 0, 99, nil)
	_, h, err := Decode(bogus)
	if !errors.Is(err, ErrUnknownCommandID) {
		t.Fatalf("notanish id: %v", err)
	}
	if h.Sequence != 99 {
		t.Fatalf("header baribir o'qilishi kerak (nack'ka seq uchun): %+v", h)
	}
}
