package pdu

import (
	"bytes"
	"testing"

	"smpp/tlv"
)

// goldenSubmitSMHex — to'liq submit_sm (54 = 0x36 oktet), seq=2:
// "Bank" (alphanumeric 5/0) → 998901234567 (international 1/1),
// matn "Salom" (ASCII, dc=0), DLR so'ralgan (registered_delivery=0x01).
const goldenSubmitSMHex = `
00 00 00 36 00 00 00 04 00 00 00 00 00 00 00 02
00
05 00 42 61 6E 6B 00
01 01 39 39 38 39 30 31 32 33 34 35 36 37 00
00 00 00
00 00
01 00 00 00
05 53 61 6C 6F 6D`

var goldenSubmitSM = SubmitSM{SMFields: SMFields{
	Source:             Address{TON: 5, NPI: 0, Addr: "Bank"},
	Dest:               Address{TON: 1, NPI: 1, Addr: "998901234567"},
	RegisteredDelivery: DLRFinal,
	ShortMessage:       []byte("Salom"),
}}

func TestSubmitSMEncodeGolden(t *testing.T) {
	frame, err := goldenSubmitSM.Encode(2)
	if err != nil {
		t.Fatalf("Encode xatosi: %v", err)
	}
	want := mustHex(t, goldenSubmitSMHex)
	if !bytes.Equal(frame, want) {
		t.Errorf("Encode = % X,\nkutilgan % X", frame, want)
	}
}

func TestSubmitSMDecodeGolden(t *testing.T) {
	sm, h, err := DecodeSubmitSM(mustHex(t, goldenSubmitSMHex))
	if err != nil {
		t.Fatalf("DecodeSubmitSM xatosi: %v", err)
	}
	if h.Sequence != 2 {
		t.Errorf("seq = %d, kutilgan 2", h.Sequence)
	}
	if sm.Source != goldenSubmitSM.Source || sm.Dest != goldenSubmitSM.Dest {
		t.Errorf("manzillar: %+v / %+v", sm.Source, sm.Dest)
	}
	if !bytes.Equal(sm.ShortMessage, []byte("Salom")) {
		t.Errorf("short_message = % X", sm.ShortMessage)
	}
	if !sm.RegisteredDelivery.WantsDLR() {
		t.Error("DLR so'ralgan bo'lishi kerak (0x01)")
	}
	if len(sm.TLVs) != 0 {
		t.Errorf("TLV kutilmagan edi: %+v", sm.TLVs)
	}
}

// goldenSubmitSMRespHex — submit_sm_resp (23 = 0x17 oktet), seq=2,
// message_id "7F3A9B" (2-bob mashqidagi frame bilan solishtiring).
const goldenSubmitSMRespHex = `
00 00 00 17 80 00 00 04 00 00 00 00 00 00 00 02
37 46 33 41 39 42 00`

func TestSubmitSMRespGolden(t *testing.T) {
	in := SubmitSMResp{MessageID: "7F3A9B"}
	frame, err := in.Encode(2)
	if err != nil {
		t.Fatalf("Encode xatosi: %v", err)
	}
	want := mustHex(t, goldenSubmitSMRespHex)
	if !bytes.Equal(frame, want) {
		t.Fatalf("Encode = % X,\nkutilgan % X", frame, want)
	}
	out, h, err := DecodeSubmitSMResp(frame)
	if err != nil {
		t.Fatalf("DecodeSubmitSMResp xatosi: %v", err)
	}
	if out.MessageID != "7F3A9B" || h.Sequence != 2 {
		t.Errorf("round-trip: %+v seq=%d", out, h.Sequence)
	}
}

func TestSubmitSMRespErrorNoBody(t *testing.T) {
	// §4.4.2: status != 0 → message_id qaytarilmaydi.
	in := SubmitSMResp{Status: 0x58 /* ESME_RTHROTTLED */, MessageID: "IGNORED"}
	frame, err := in.Encode(7)
	if err != nil {
		t.Fatalf("Encode xatosi: %v", err)
	}
	if len(frame) != HeaderSize {
		t.Fatalf("xatoli resp %d oktet, 16 bo'lishi kerak", len(frame))
	}
	out, _, err := DecodeSubmitSMResp(frame)
	if err != nil {
		t.Fatalf("DecodeSubmitSMResp xatosi: %v", err)
	}
	if out.Status != 0x58 || out.MessageID != "" {
		t.Errorf("xatoli resp: %+v", out)
	}
}

// goldenDeliverDLRHex — to'liq deliver_sm DLR (175 = 0xAF oktet), seq=9:
// manzillar ALMASHGAN (source=998..., dest=Bank), esm_class=0x04 (DLR),
// short_message = Appendix B uslubidagi matn, TLV tail = 3-bob goldeni
// (receipted_message_id "7F3A9B" + message_state=2 + network_error_code GSM/0).
const goldenDeliverDLRHex = `
00 00 00 AF 00 00 00 05 00 00 00 00 00 00 00 09
00
01 01 39 39 38 39 30 31 32 33 34 35 36 37 00
05 00 42 61 6E 6B 00
04 00 00
00 00
00 00 00 00
67
69 64 3A 37 46 33 41 39 42 20 73 75 62 3A 30 30
31 20 64 6C 76 72 64 3A 30 30 31 20 73 75 62 6D
69 74 20 64 61 74 65 3A 32 36 30 37 31 37 31 32
30 35 20 64 6F 6E 65 20 64 61 74 65 3A 32 36 30
37 31 37 31 32 30 36 20 73 74 61 74 3A 44 45 4C
49 56 52 44 20 65 72 72 3A 30 30 30 20 74 65 78
74 3A 53 61 6C 6F 6D
00 1E 00 07 37 46 33 41 39 42 00
04 27 00 01 02
04 23 00 03 03 00 00`

const dlrText = "id:7F3A9B sub:001 dlvrd:001 submit date:2607171205 done date:2607171206 stat:DELIVRD err:000 text:Salom"

func TestDeliverSMDLRGolden(t *testing.T) {
	in := DeliverSM{SMFields: SMFields{
		Source:       Address{TON: 1, NPI: 1, Addr: "998901234567"},
		Dest:         Address{TON: 5, NPI: 0, Addr: "Bank"},
		EsmClass:     0x04,
		ShortMessage: []byte(dlrText),
		TLVs: []tlv.TLV{
			tlv.CString(tlv.ReceiptedMessageID, "7F3A9B"),
			tlv.U8(tlv.MessageState, 2),
			{Tag: tlv.NetworkErrorCode, Value: []byte{0x03, 0x00, 0x00}},
		},
	}}
	frame, err := in.Encode(9)
	if err != nil {
		t.Fatalf("Encode xatosi: %v", err)
	}
	want := mustHex(t, goldenDeliverDLRHex)
	if !bytes.Equal(frame, want) {
		t.Fatalf("Encode = % X,\nkutilgan % X", frame, want)
	}

	out, h, err := DecodeDeliverSM(frame)
	if err != nil {
		t.Fatalf("DecodeDeliverSM xatosi: %v", err)
	}
	if h.Sequence != 9 {
		t.Errorf("seq = %d", h.Sequence)
	}
	if !out.EsmClass.IsDeliveryReceipt() {
		t.Error("esm_class=0x04 DLR sifatida tanilishi kerak")
	}
	if string(out.ShortMessage) != dlrText {
		t.Errorf("DLR matni buzildi: %q", out.ShortMessage)
	}
	if st, ok := tlv.Find(out.TLVs, tlv.MessageState); !ok {
		t.Error("message_state TLV topilmadi")
	} else if v, _ := st.Uint8Value(); v != 2 {
		t.Errorf("message_state = %d, kutilgan 2 (DELIVERED)", v)
	}
	if id, ok := tlv.Find(out.TLVs, tlv.ReceiptedMessageID); !ok {
		t.Error("receipted_message_id TLV topilmadi")
	} else if s, _ := id.CStringValue(); s != "7F3A9B" {
		t.Errorf("receipted_message_id = %q", s)
	}
}

func TestDeliverSMUnusedFieldsEnforced(t *testing.T) {
	// §4.6.1: deliver_sm'da schedule/validity/replace/sm_default NULL SHART.
	bad := DeliverSM{SMFields: SMFields{
		Dest:           Address{TON: 1, NPI: 1, Addr: "998901234567"},
		ValidityPeriod: "000001000000000R",
	}}
	if _, err := bad.Encode(1); err == nil {
		t.Error("validity_period'li deliver_sm rad etilishi kerak edi")
	}
}

func TestDeliverSMRespNullMessageID(t *testing.T) {
	frame := DeliverSMResp{}.Encode(9)
	// Header + bitta NULL oktet (§4.6.2).
	if len(frame) != HeaderSize+1 || frame[HeaderSize] != 0x00 {
		t.Fatalf("deliver_sm_resp = % X, header+0x00 kutilgan", frame)
	}
	out, h, err := DecodeDeliverSMResp(frame)
	if err != nil {
		t.Fatalf("DecodeDeliverSMResp xatosi: %v", err)
	}
	if out.Status != 0 || h.Sequence != 9 {
		t.Errorf("resp: %+v seq=%d", out, h.Sequence)
	}
	// Body'siz variantga toqat (ba'zi stack'lar shunday yuboradi).
	hdrOnly := EncodeHeader(Header{Length: 16, ID: CmdDeliverSMResp, Sequence: 3})
	if _, _, err := DecodeDeliverSMResp(hdrOnly[:]); err != nil {
		t.Errorf("header-only deliver_sm_resp qabul qilinishi kerak edi: %v", err)
	}
}

func TestSubmitSMValidation(t *testing.T) {
	t.Run("short_message 255 oktet", func(t *testing.T) {
		sm := SubmitSM{SMFields: SMFields{
			Dest:         Address{TON: 1, NPI: 1, Addr: "998901234567"},
			ShortMessage: make([]byte, 255),
		}}
		if _, err := sm.Encode(1); err == nil {
			t.Error("255 oktetlik short_message rad etilishi kerak edi")
		}
	})
	t.Run("short_message + message_payload birga", func(t *testing.T) {
		sm := SubmitSM{SMFields: SMFields{
			Dest:         Address{TON: 1, NPI: 1, Addr: "998901234567"},
			ShortMessage: []byte("Salom"),
			TLVs:         []tlv.TLV{{Tag: tlv.MessagePayload, Value: []byte("uzun matn")}},
		}}
		if _, err := sm.Encode(1); err == nil {
			t.Error("§3.2.3 qoidasi buzilishi rad etilishi kerak edi")
		}
	})
	t.Run("vaqt field'i 10 belgi", func(t *testing.T) {
		sm := SubmitSM{SMFields: SMFields{
			Dest:           Address{TON: 1, NPI: 1, Addr: "998901234567"},
			ValidityPeriod: "0000010000",
		}}
		if _, err := sm.Encode(1); err == nil {
			t.Error("16 belgidan farqli vaqt field'i rad etilishi kerak edi")
		}
	})
}

func TestSubmitSMWithTimeAndTLVRoundTrip(t *testing.T) {
	in := SubmitSM{SMFields: SMFields{
		ServiceType:          "CMT",
		Source:               Address{TON: 1, NPI: 1, Addr: "998901111111"},
		Dest:                 Address{TON: 1, NPI: 1, Addr: "998902222222"},
		EsmClass:             EsmClass(0).WithUDHI(),
		ScheduleDeliveryTime: "",
		ValidityPeriod:       "000001000000000R",
		RegisteredDelivery:   DLRFinal | Intermediate, // 0x11
		DataCoding:           8,
		ShortMessage:         []byte{0x05, 0x00, 0x03, 0x2A, 0x02, 0x01, 0x00, 0x41},
		TLVs:                 []tlv.TLV{tlv.U16(tlv.UserMessageReference, 42)},
	}}
	frame, err := in.Encode(100)
	if err != nil {
		t.Fatalf("Encode xatosi: %v", err)
	}
	out, h, err := DecodeSubmitSM(frame)
	if err != nil {
		t.Fatalf("DecodeSubmitSM xatosi: %v", err)
	}
	if h.Sequence != 100 {
		t.Errorf("seq = %d", h.Sequence)
	}
	if out.ValidityPeriod != in.ValidityPeriod || out.ServiceType != "CMT" {
		t.Errorf("field'lar: %+v", out.SMFields)
	}
	if !out.EsmClass.HasUDHI() {
		t.Error("UDHI saqlanishi kerak edi")
	}
	if out.RegisteredDelivery != 0x11 {
		t.Errorf("registered_delivery = 0x%02X", uint8(out.RegisteredDelivery))
	}
	if !bytes.Equal(out.ShortMessage, in.ShortMessage) {
		t.Errorf("short_message: % X", out.ShortMessage)
	}
	if ref, ok := tlv.Find(out.TLVs, tlv.UserMessageReference); !ok {
		t.Error("user_message_reference yo'qoldi")
	} else if v, _ := ref.Uint16Value(); v != 42 {
		t.Errorf("user_message_reference = %d", v)
	}
}

func TestDecodeSubmitSMTruncatedSmLength(t *testing.T) {
	// sm_length frame'da qolgan baytlardan katta — xato, panic emas.
	frame, err := goldenSubmitSM.Encode(2)
	if err != nil {
		t.Fatal(err)
	}
	frame[len(frame)-6] = 200 // sm_length baytini buzamiz
	if _, _, err := DecodeSubmitSM(frame); err == nil {
		t.Error("buzuq sm_length xato berishi kerak edi")
	}
}
