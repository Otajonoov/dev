package dlr

import (
	"errors"
	"testing"
	"time"

	"smpp/tlv"
)

// 5-bob golden DLR'ining short_message matni — kitob baytlari bilan bir xil.
const goldenReceiptText = "id:7F3A9B sub:001 dlvrd:001 submit date:2607171205 done date:2607171206 stat:DELIVRD err:000 text:Salom"

func date(s string) time.Time {
	layout := "0601021504"
	if len(s) == 12 {
		layout = "060102150405"
	}
	t, err := time.ParseInLocation(layout, s, time.UTC)
	if err != nil {
		panic(err)
	}
	return t
}

// Turli operator formatlari — research'dagi real og'ishlar asosida.
func TestParseOperatorVariants(t *testing.T) {
	tests := []struct {
		name string
		text string
		want Receipt
	}{
		{
			// Appendix B'ning "tipik misoli" — 5-bob golden DLR'i.
			name: "appendix-b-canonical",
			text: goldenReceiptText,
			want: Receipt{
				ID: "7F3A9B", State: StateDelivered, Sub: 1, Dlvrd: 1,
				SubmitDate: date("2607171205"), DoneDate: date("2607171206"),
				Stat: "DELIVRD", Err: "000", Text: "Salom",
			},
		},
		{
			// Tartib boshqacha (stat oldinda), text: umuman yo'q.
			name: "reordered-no-text",
			text: "stat:DELIVRD err:000 id:ab021099504969 sub:001 dlvrd:001 submit date:1704181518 done date:1704181519",
			want: Receipt{
				ID: "ab021099504969", State: StateDelivered, Sub: 1, Dlvrd: 1,
				SubmitDate: date("1704181518"), DoneDate: date("1704181519"),
				Stat: "DELIVRD", Err: "000",
			},
		},
		{
			// Sekundli sana + sub/dlvrd tashlab ketilgan (real vendor case).
			name: "seconds-date-missing-sub",
			text: "id:1526758174 submit date:170124090433 done date:170124090455 stat:DELIVRD err:000 text:hellow",
			want: Receipt{
				ID: "1526758174", State: StateDelivered, Sub: -1, Dlvrd: -1,
				SubmitDate: date("170124090433"), DoneDate: date("170124090455"),
				Stat: "DELIVRD", Err: "000", Text: "hellow",
			},
		},
		{
			// Katta-kichik harf aralash, to'liq nom "DELIVERED", err harfli.
			name: "case-insensitive-long-stat",
			text: "ID:0D9A2F Sub:001 Dlvrd:001 Submit Date:2607181200 Done Date:2607181201 Stat:DELIVERED Err:00A Text:OK",
			want: Receipt{
				ID: "0D9A2F", State: StateDelivered, Sub: 1, Dlvrd: 1,
				SubmitDate: date("2607181200"), DoneDate: date("2607181201"),
				Stat: "DELIVERED", Err: "00A", Text: "OK",
			},
		},
		{
			// UNDELIV + text ichida bo'sh joylar va ':' — text oxirgi kalit,
			// qolgan hammasi unga tegishli.
			name: "undeliv-text-with-colon",
			text: "id:88 sub:001 dlvrd:000 submit date:2607190800 done date:2607190805 stat:UNDELIV err:034 text:Kod: 5521 kirish",
			want: Receipt{
				ID: "88", State: StateUndeliverable, Sub: 1, Dlvrd: 0,
				SubmitDate: date("2607190800"), DoneDate: date("2607190805"),
				Stat: "UNDELIV", Err: "034", Text: "Kod: 5521 kirish",
			},
		},
		{
			// Minimal vendor varianti: faqat id va stat.
			name: "minimal-id-stat",
			text: "id:XYZ-123 stat:EXPIRED",
			want: Receipt{ID: "XYZ-123", State: StateExpired, Sub: -1, Dlvrd: -1, Stat: "EXPIRED"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse([]byte(tt.text), nil)
			if err != nil {
				t.Fatalf("Parse xato: %v", err)
			}
			if got != tt.want {
				t.Errorf("Parse:\n  got  %+v\n  want %+v", got, tt.want)
			}
		})
	}
}

// TLV'lar matndan USTUN: zid qiymatlarda TLV yutadi (§5.3.2.12/.35 —
// strukturali manba haqiqat manbai).
func TestParseTLVPriority(t *testing.T) {
	tlvs := []tlv.TLV{
		tlv.CString(tlv.ReceiptedMessageID, "TLVID99"),
		tlv.U8(tlv.MessageState, uint8(StateExpired)),
		{Tag: tlv.NetworkErrorCode, Value: []byte{0x03, 0x01, 0x2C}}, // GSM, 300
	}
	// Matn ataylab boshqa id va boshqa stat aytadi.
	r, err := Parse([]byte(goldenReceiptText), tlvs)
	if err != nil {
		t.Fatalf("Parse xato: %v", err)
	}
	if r.ID != "TLVID99" {
		t.Errorf("ID = %q, TLV ustun bo'lishi kerak edi (TLVID99)", r.ID)
	}
	if r.State != StateExpired {
		t.Errorf("State = %v, TLV'dagi EXPIRED kutilgan", r.State)
	}
	if r.NetErr == nil || r.NetErr.Type != 3 || r.NetErr.Code != 300 {
		t.Errorf("NetErr = %+v, {Type:3 Code:300} kutilgan", r.NetErr)
	}
	// Matn field'lari baribir saqlanadi (Stat xom holida).
	if r.Stat != "DELIVRD" || r.Text != "Salom" {
		t.Errorf("matn field'lari yo'qoldi: Stat=%q Text=%q", r.Stat, r.Text)
	}
}

// Faqat TLV, matn bo'sh — spec'ga to'liq amal qilgan SMSC.
func TestParseTLVOnly(t *testing.T) {
	tlvs := []tlv.TLV{
		tlv.CString(tlv.ReceiptedMessageID, "7F3A9B"),
		tlv.U8(tlv.MessageState, uint8(StateDelivered)),
	}
	r, err := Parse(nil, tlvs)
	if err != nil {
		t.Fatalf("Parse xato: %v", err)
	}
	if r.ID != "7F3A9B" || r.State != StateDelivered {
		t.Errorf("TLV-only parse: %+v", r)
	}
}

// Buzuq/yarim DLR — XATO EMAS, partial Receipt.
func TestParsePartial(t *testing.T) {
	r, err := Parse([]byte("id:ABC stat:BROKEN_9 err:???"), nil)
	if err != nil {
		t.Fatalf("partial DLR xato bermasligi kerak: %v", err)
	}
	if r.ID != "ABC" {
		t.Errorf("ID = %q", r.ID)
	}
	if r.State != 0 {
		t.Errorf("notanish stat uchun State=0 kutilgan, keldi %v", r.State)
	}
	if r.Stat != "BROKEN_9" || r.Err != "???" {
		t.Errorf("xom qiymatlar saqlanishi kerak: %+v", r)
	}
	// Buzuq sub: raqam emas — tashlanadi, -1 qoladi.
	r2, err := Parse([]byte("id:1 sub:abc dlvrd:001"), nil)
	if err != nil {
		t.Fatal(err)
	}
	if r2.Sub != -1 || r2.Dlvrd != 1 {
		t.Errorf("Sub=%d (kutilgan -1), Dlvrd=%d (kutilgan 1)", r2.Sub, r2.Dlvrd)
	}
}

// MO xabar (oddiy matn) receipt emas — ErrNotReceipt.
func TestParseNotReceipt(t *testing.T) {
	_, err := Parse([]byte("Salom, bu oddiy MO xabar"), nil)
	if !errors.Is(err, ErrNotReceipt) {
		t.Fatalf("ErrNotReceipt kutilgan edi, keldi: %v", err)
	}
	// Kalit token o'rtasida bo'lsa sanalmaydi: "resubmit:" ichidagi "sub:".
	_, err = Parse([]byte("please resubmit:now"), nil)
	if !errors.Is(err, ErrNotReceipt) {
		t.Fatalf("token-boshi tekshiruvi ishlamadi: %v", err)
	}
}

func TestStateFromStat(t *testing.T) {
	tests := []struct {
		in   string
		want MessageState
	}{
		{"DELIVRD", StateDelivered},
		{"DELIVERED", StateDelivered},
		{"delivrd", StateDelivered},
		{"EXPIRED", StateExpired},
		{"DELETED", StateDeleted},
		{"UNDELIV", StateUndeliverable},
		{"UNDELIVERABLE", StateUndeliverable},
		{"ACCEPTD", StateAccepted},
		{"REJECTD", StateRejected},
		{"REJECTED", StateRejected},
		{"ENROUTE", StateEnroute},
		{"UNKNOWN", StateUnknown},
		{"BLABLA", 0},
		{"", 0},
	}
	for _, tt := range tests {
		if got := StateFromStat(tt.in); got != tt.want {
			t.Errorf("StateFromStat(%q) = %v, kutilgan %v", tt.in, got, tt.want)
		}
	}
}

func TestStateFinalAndAbbrev(t *testing.T) {
	finals := map[MessageState]bool{
		StateEnroute: false, StateDelivered: true, StateExpired: true,
		StateDeleted: true, StateUndeliverable: true, StateAccepted: false,
		StateUnknown: false, StateRejected: true,
	}
	for s, want := range finals {
		if s.Final() != want {
			t.Errorf("%v.Final() = %v, kutilgan %v", s, s.Final(), want)
		}
	}
	// Abbrev ↔ StateFromStat round-trip: har qisqartma o'z holatiga qaytadi.
	for s := StateEnroute; s <= StateRejected; s++ {
		if got := StateFromStat(s.Abbrev()); got != s {
			t.Errorf("StateFromStat(%q) = %v, kutilgan %v", s.Abbrev(), got, s)
		}
	}
}
