package coding

import (
	"bytes"
	"strings"
	"testing"

	"smpp/tlv"
)

func TestUDHEncodeGolden(t *testing.T) {
	// 8-bit: 05 00 03 RR TT SS; 16-bit: 06 08 04 RR RR TT SS.
	u8 := UDH{RefNum: 0x5A, Total: 2, Seq: 1}
	if got := u8.Encode(); !bytes.Equal(got, mustHex(t, "05 00 03 5A 02 01")) {
		t.Errorf("UDH8 = % X", got)
	}
	u16 := UDH{RefNum: 0x1234, Total: 3, Seq: 2, Is16bit: true}
	if got := u16.Encode(); !bytes.Equal(got, mustHex(t, "06 08 04 12 34 03 02")) {
		t.Errorf("UDH16 = % X", got)
	}
}

func TestParseUDHRoundTrip(t *testing.T) {
	for _, u := range []UDH{
		{RefNum: 0x5A, Total: 2, Seq: 2},
		{RefNum: 0xABCD, Total: 5, Seq: 3, Is16bit: true},
	} {
		sm := append(u.Encode(), []byte("matn")...)
		got, payload, found, err := ParseUDH(sm)
		if err != nil || !found {
			t.Fatalf("ParseUDH(%+v): %v, found=%v", u, err, found)
		}
		if got != u || string(payload) != "matn" {
			t.Errorf("round-trip: %+v, payload=%q", got, payload)
		}
	}
}

func TestParseUDHUnknownIESkipped(t *testing.T) {
	// Port addressing IE (0x05) + concat IE birga: concat topiladi, port skip.
	sm := mustHex(t, "0B 05 04 0B 84 0B 84 00 03 5A 02 01")
	sm = append(sm, 'x')
	u, payload, found, err := ParseUDH(sm)
	if err != nil || !found {
		t.Fatalf("ParseUDH: %v, found=%v", err, found)
	}
	if u.RefNum != 0x5A || u.Total != 2 || u.Seq != 1 || string(payload) != "x" {
		t.Errorf("u=%+v payload=%q", u, payload)
	}
	// Faqat notanish IE: found=false, lekin xato emas.
	only := mustHex(t, "06 05 04 0B 84 0B 84")
	if _, _, found, err := ParseUDH(only); err != nil || found {
		t.Errorf("faqat port IE: found=%v, err=%v", found, err)
	}
}

func TestParseUDHErrors(t *testing.T) {
	if _, _, _, err := ParseUDH(nil); err == nil {
		t.Error("bo'sh sm xato berishi kerak edi")
	}
	if _, _, _, err := ParseUDH(mustHex(t, "05 00 03 5A")); err == nil {
		t.Error("UDHL'dan qisqa sm xato berishi kerak edi")
	}
}

func TestSplitGSM7TwoSegments(t *testing.T) {
	// 161 ta 'a' — 160'dan bitta oshiq: UDH8 bilan 153 + 8.
	text := strings.Repeat("a", 161)
	segs, err := Split(text, DCDefault, MethodUDH8, 0x5A)
	if err != nil {
		t.Fatal(err)
	}
	if len(segs) != 2 {
		t.Fatalf("%d segment, kutilgan 2", len(segs))
	}
	// 1-segment: UDH golden + 153 'a'.
	wantPrefix := mustHex(t, "05 00 03 5A 02 01")
	if !bytes.Equal(segs[0].Data[:6], wantPrefix) {
		t.Errorf("seg1 UDH = % X", segs[0].Data[:6])
	}
	if len(segs[0].Data) != 6+153 || len(segs[1].Data) != 6+8 {
		t.Errorf("uzunliklar: %d, %d; kutilgan 159, 14", len(segs[0].Data), len(segs[1].Data))
	}
	if segs[1].Data[5] != 2 {
		t.Errorf("seg2 seqnum = %d, kutilgan 2", segs[1].Data[5])
	}
	// Payload'lar birlashtirilsa original matn chiqadi.
	p1, _ := DecodeGSM7(segs[0].Data[6:])
	p2, _ := DecodeGSM7(segs[1].Data[6:])
	if p1+p2 != text {
		t.Error("segmentlar birlashmasi original bilan mos emas")
	}
}

func TestSplitUCS2(t *testing.T) {
	// 68 kirill belgi — 70 ga sig'adi-yu? Yo'q: 68 ≤ 70, bitta segment!
	text68 := strings.Repeat("д", 68)
	segs, err := Split(text68, DCUCS2, MethodUDH8, 1)
	if err != nil || len(segs) != 1 {
		t.Fatalf("68 belgi: %d segment, %v; kutilgan 1", len(segs), err)
	}
	// 71 belgi — endi 2 segment: 67 + 4.
	text71 := strings.Repeat("д", 71)
	segs, err = Split(text71, DCUCS2, MethodUDH8, 1)
	if err != nil || len(segs) != 2 {
		t.Fatalf("71 belgi: %d segment, %v; kutilgan 2", len(segs), err)
	}
	if len(segs[0].Data) != 6+67*2 || len(segs[1].Data) != 6+4*2 {
		t.Errorf("uzunliklar: %d, %d", len(segs[0].Data), len(segs[1].Data))
	}
}

func TestSplitExtensionBoundary(t *testing.T) {
	// 152 'a' + '€' + 10 'b': € 2 septet — 153-byudjetga sig'maydi,
	// BUTUN holda 2-segmentga o'tishi kerak (ESC juftligi bo'linmaydi).
	text := strings.Repeat("a", 152) + "€" + strings.Repeat("b", 10)
	segs, err := Split(text, DCDefault, MethodUDH8, 7)
	if err != nil {
		t.Fatal(err)
	}
	if len(segs) != 2 {
		t.Fatalf("%d segment, kutilgan 2", len(segs))
	}
	if len(segs[0].Data) != 6+152 {
		t.Errorf("seg1 %d oktet — € kirib qolganmi?", len(segs[0].Data))
	}
	p2, err := DecodeGSM7(segs[1].Data[6:])
	if err != nil {
		t.Fatalf("seg2 decode: %v — ESC juftligi bo'lingan!", err)
	}
	if !strings.HasPrefix(p2, "€") {
		t.Errorf("seg2 € bilan boshlanishi kerak: %q", p2)
	}
}

func TestSplitSurrogateBoundary(t *testing.T) {
	// Jami 78 unit (>70 — concat kerak); 1-segment byudjeti 67: 66 'д'dan keyin
	// emoji (2 unit) sig'maydi — butun holda keyingi segmentga o'tishi kerak.
	text := strings.Repeat("д", 66) + "😀" + strings.Repeat("д", 10)
	segs, err := Split(text, DCUCS2, MethodUDH8, 7)
	if err != nil {
		t.Fatal(err)
	}
	if len(segs) != 2 {
		t.Fatalf("%d segment, kutilgan 2", len(segs))
	}
	for i, s := range segs {
		if _, err := DecodeUCS2(s.Data[6:]); err != nil {
			t.Errorf("seg%d decode xatosi: %v — surrogate bo'lingan!", i+1, err)
		}
	}
}

func TestSplitSarMethod(t *testing.T) {
	text := strings.Repeat("a", 200)
	segs, err := Split(text, DCDefault, MethodSarTLV, 0x0102)
	if err != nil || len(segs) != 2 {
		t.Fatalf("%d segment, %v", len(segs), err)
	}
	for i, s := range segs {
		if s.UDH != nil || s.Sar == nil {
			t.Fatalf("sar usulida UDH bo'lmasligi, Sar bo'lishi kerak: %+v", s)
		}
		if s.Sar.RefNum != 0x0102 || s.Sar.Total != 2 || s.Sar.Seq != uint8(i+1) {
			t.Errorf("seg%d sar: %+v", i+1, s.Sar)
		}
	}
	// TLV round-trip: uchlik → SarFromTLVs.
	got, found, err := SarFromTLVs(segs[0].Sar.TLVs())
	if err != nil || !found || got != *segs[0].Sar {
		t.Errorf("SarFromTLVs = %+v, found=%v, %v", got, found, err)
	}
}

func TestSarFromTLVsIncompleteIgnored(t *testing.T) {
	// Spec (§5.3.2.22): uchlik chala bo'lsa — ignore (xato emas).
	partial := []tlv.TLV{tlv.U16(tlv.SarMsgRefNum, 9)}
	if _, found, err := SarFromTLVs(partial); found || err != nil {
		t.Errorf("chala uchlik: found=%v, err=%v; kutilgan false, nil", found, err)
	}
	// Buzuq value — endi xato.
	bad := []tlv.TLV{
		{Tag: tlv.SarMsgRefNum, Value: []byte{1}}, // 1 oktet — 2 kerak
		tlv.U8(tlv.SarTotalSegments, 2),
		tlv.U8(tlv.SarSegmentSeqnum, 1),
	}
	if _, _, err := SarFromTLVs(bad); err == nil {
		t.Error("buzuq sar value xato berishi kerak edi")
	}
}

func TestSplitPayloadMethod(t *testing.T) {
	text := strings.Repeat("a", 500)
	segs, err := Split(text, DCDefault, MethodPayload, 0)
	if err != nil || len(segs) != 1 {
		t.Fatalf("payload: %d segment, %v; kutilgan 1", len(segs), err)
	}
	if len(segs[0].Data) != 500 || segs[0].UDH != nil || segs[0].Sar != nil {
		t.Errorf("payload segment: len=%d, %+v", len(segs[0].Data), segs[0])
	}
}

func TestSplitSingleFits(t *testing.T) {
	segs, err := Split("Salom", DCDefault, MethodUDH8, 1)
	if err != nil || len(segs) != 1 || segs[0].UDH != nil {
		t.Errorf("qisqa matn UDH'siz bitta segment bo'lishi kerak: %+v, %v", segs, err)
	}
}

func TestCountSegments(t *testing.T) {
	tests := []struct {
		text string
		dc   DataCoding
		n    int
	}{
		{"Assalomu alaykum! Kodingiz: 5521", DCDefault, 1},
		{strings.Repeat("a", 161), DCDefault, 2},
		{strings.Repeat("д", 300), DCUCS2, 5},              // ceil(300/67)
		{"soʻm " + strings.Repeat("a", 156), DCDefault, 2}, // normalize'dan keyin GSM7 161 belgi — 160 dan oshadi
	}
	for _, tt := range tests {
		dc, n, err := CountSegments(tt.text)
		if err != nil || dc != tt.dc || n != tt.n {
			t.Errorf("CountSegments(len=%d): dc=0x%02X n=%d %v; kutilgan 0x%02X/%d",
				len(tt.text), uint8(dc), n, err, uint8(tt.dc), tt.n)
		}
	}
}

func TestRefCounter(t *testing.T) {
	rc := NewRefCounter()
	if rc.Next("998901111111") != 1 || rc.Next("998901111111") != 2 {
		t.Error("bir dest uchun ketma-ket o'sishi kerak")
	}
	if rc.Next("998902222222") != 1 {
		t.Error("boshqa dest mustaqil hisoblanishi kerak")
	}
}
