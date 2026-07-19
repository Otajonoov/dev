package tlv

import (
	"bytes"
	"encoding/hex"
	"errors"
	"strings"
	"testing"
)

func mustHex(t *testing.T, s string) []byte {
	t.Helper()
	clean := strings.NewReplacer(" ", "", "\n", "", "\t", "").Replace(s)
	b, err := hex.DecodeString(clean)
	if err != nil {
		t.Fatalf("hex decode xatosi: %v", err)
	}
	return b
}

// goldenDLRTailHex — tipik DLR deliver_sm'ning TLV tail'i (9-bobda to'liq
// kontekstda ko'ramiz): receipted_message_id "7F3A9B" + message_state=2
// (DELIVERED) + network_error_code GSM/0.
const goldenDLRTailHex = `
00 1E 00 07 37 46 33 41 39 42 00
04 27 00 01 02
04 23 00 03 03 00 00`

func TestDecodeGoldenTail(t *testing.T) {
	tlvs, err := Decode(mustHex(t, goldenDLRTailHex))
	if err != nil {
		t.Fatalf("Decode xatosi: %v", err)
	}
	if len(tlvs) != 3 {
		t.Fatalf("%d TLV, kutilgan 3", len(tlvs))
	}

	if tlvs[0].Tag != ReceiptedMessageID {
		t.Errorf("tlvs[0].Tag = %s, kutilgan receipted_message_id", tlvs[0].Tag)
	}
	if id, err := tlvs[0].CStringValue(); err != nil || id != "7F3A9B" {
		t.Errorf("CStringValue = %q, %v; kutilgan \"7F3A9B\"", id, err)
	}

	if tlvs[1].Tag != MessageState {
		t.Errorf("tlvs[1].Tag = %s, kutilgan message_state", tlvs[1].Tag)
	}
	if st, err := tlvs[1].Uint8Value(); err != nil || st != 2 {
		t.Errorf("Uint8Value = %d, %v; kutilgan 2 (DELIVERED)", st, err)
	}

	if tlvs[2].Tag != NetworkErrorCode {
		t.Errorf("tlvs[2].Tag = %s, kutilgan network_error_code", tlvs[2].Tag)
	}
	if ne, err := tlvs[2].NetworkError(); err != nil || ne != (NetworkError{Type: 3, Code: 0}) {
		t.Errorf("NetworkError = %+v, %v; kutilgan {Type:3 Code:0}", ne, err)
	}
}

func TestEncodeGoldenTail(t *testing.T) {
	var b bytes.Buffer
	err := Encode(&b, []TLV{
		CString(ReceiptedMessageID, "7F3A9B"),
		U8(MessageState, 2),
		{Tag: NetworkErrorCode, Value: []byte{0x03, 0x00, 0x00}},
	})
	if err != nil {
		t.Fatalf("Encode xatosi: %v", err)
	}
	want := mustHex(t, goldenDLRTailHex)
	if !bytes.Equal(b.Bytes(), want) {
		t.Errorf("Encode = % X,\nkutilgan % X", b.Bytes(), want)
	}
}

func TestZeroLengthTLV(t *testing.T) {
	// alert_on_message_delivery (§5.3.2.41) — value uzunligi 0.
	wire := mustHex(t, `13 0C 00 00`)

	tlvs, err := Decode(wire)
	if err != nil {
		t.Fatalf("Decode xatosi: %v", err)
	}
	if len(tlvs) != 1 || tlvs[0].Tag != AlertOnMessageDelivery || len(tlvs[0].Value) != 0 {
		t.Fatalf("kutilgan zero-length alert_on_message_delivery, keldi: %+v", tlvs)
	}

	var b bytes.Buffer
	if err := Encode(&b, []TLV{Empty(AlertOnMessageDelivery)}); err != nil {
		t.Fatalf("Encode xatosi: %v", err)
	}
	if !bytes.Equal(b.Bytes(), wire) {
		t.Errorf("Encode = % X, kutilgan % X", b.Bytes(), wire)
	}
}

func TestUnknownTagRoundTrip(t *testing.T) {
	// Vendor tag 0x1401 — notanish, lekin Decode uni SAQLASHI kerak (§3.3
	// forward compatibility: ignore ≠ yo'qotish) va qayta Encode'da baytlar
	// aynan tiklanishi kerak.
	wire := mustHex(t, `14 01 00 02 AB CD 02 10 00 01 34`)

	tlvs, err := Decode(wire)
	if err != nil {
		t.Fatalf("Decode xatosi: %v", err)
	}
	if len(tlvs) != 2 {
		t.Fatalf("%d TLV, kutilgan 2", len(tlvs))
	}
	if tlvs[0].Tag != Tag(0x1401) || !tlvs[0].Tag.IsVendor() {
		t.Errorf("tlvs[0].Tag = %s, vendor 0x1401 kutilgan edi", tlvs[0].Tag)
	}
	if v, err := tlvs[1].Uint8Value(); err != nil || v != 0x34 {
		t.Errorf("sc_interface_version = 0x%02X, %v; kutilgan 0x34", v, err)
	}

	var b bytes.Buffer
	if err := Encode(&b, tlvs); err != nil {
		t.Fatalf("Encode xatosi: %v", err)
	}
	if !bytes.Equal(b.Bytes(), wire) {
		t.Errorf("round-trip buzildi: % X != % X", b.Bytes(), wire)
	}
}

func TestDecodeTruncated(t *testing.T) {
	tests := []struct {
		name string
		hex  string
	}{
		{"Tag+Length'ning o'zi chala", `00 1E 00`},
		{"Length va'da qilgan Value kelmagan", `04 24 00 05 41 42 43`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := Decode(mustHex(t, tt.hex)); !errors.Is(err, ErrTruncated) {
				t.Errorf("ErrTruncated kutilgan edi, keldi: %v", err)
			}
		})
	}
}

func TestDecodeEmptyTail(t *testing.T) {
	// TLV'siz PDU normal holat — bo'sh tail xato emas.
	tlvs, err := Decode(nil)
	if err != nil || len(tlvs) != 0 {
		t.Errorf("bo'sh tail: %v, %v; kutilgan nil, nil", tlvs, err)
	}
}

func TestEncodeValueTooLarge(t *testing.T) {
	var b bytes.Buffer
	err := Encode(&b, []TLV{{Tag: MessagePayload, Value: make([]byte, 0x10000)}})
	if err == nil {
		t.Error("65536 oktetlik value 2 oktetlik Length'ga sig'masligi kerak edi")
	}
}

func TestFindAndDuplicates(t *testing.T) {
	tlvs := []TLV{
		CString(CallbackNum, "998901111111"),
		U8(MessageState, 2),
		CString(CallbackNum, "998902222222"),
	}
	if _, ok := Find(tlvs, MessagePayload); ok {
		t.Error("yo'q tag topildi")
	}
	got, ok := Find(tlvs, CallbackNum)
	if !ok {
		t.Fatal("callback_num topilmadi")
	}
	// Find birinchisini qaytaradi; takrorlar to'plamda saqlanadi.
	if s, _ := got.CStringValue(); s != "998901111111" {
		t.Errorf("Find birinchi callback_num'ni qaytarishi kerak, keldi %q", s)
	}
	var count int
	for _, tl := range tlvs {
		if tl.Tag == CallbackNum {
			count++
		}
	}
	if count != 2 {
		t.Errorf("takror callback_num'lar yo'qoldi: %d ta, kutilgan 2", count)
	}
}

func TestTagString(t *testing.T) {
	tests := []struct {
		tag  Tag
		want string
	}{
		{MessagePayload, "message_payload"},
		{ScInterfaceVersion, "sc_interface_version"},
		{Tag(0x1401), "vendor(0x1401)"},
		{Tag(0x0099), "unknown(0x0099)"},
	}
	for _, tt := range tests {
		if got := tt.tag.String(); got != tt.want {
			t.Errorf("Tag(0x%04X).String() = %q, kutilgan %q", uint16(tt.tag), got, tt.want)
		}
	}
}

func TestCStringValueTolerant(t *testing.T) {
	// Spec shakli: NULL bilan.
	withNull := TLV{Tag: ReceiptedMessageID, Value: []byte("ABC\x00")}
	if s, err := withNull.CStringValue(); err != nil || s != "ABC" {
		t.Errorf("NULL'li value: %q, %v", s, err)
	}
	// Amaliyot: ayrim SMSC'lar terminatorsiz yuboradi — xato emas.
	noNull := TLV{Tag: ReceiptedMessageID, Value: []byte("ABC")}
	if s, err := noNull.CStringValue(); err != nil || s != "ABC" {
		t.Errorf("terminatorsiz value: %q, %v", s, err)
	}
	// Ichki NULL — har doim xato.
	embedded := TLV{Tag: ReceiptedMessageID, Value: []byte("A\x00C")}
	if _, err := embedded.CStringValue(); err == nil {
		t.Error("ichki NULL'li value xato bo'lishi kerak edi")
	}
}
