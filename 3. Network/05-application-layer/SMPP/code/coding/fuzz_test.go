package coding

import (
	"testing"
	"unicode/utf8"
)

// FuzzGSM7RoundTrip — GSM7 codec'ining asosiy invarianti: encode
// muvaffaqiyatli bo'lsa, decode AYNAN o'sha matnni qaytarishi kerak.
func FuzzGSM7RoundTrip(f *testing.F) {
	f.Add("hello")
	f.Add("Assalomu alaykum! Kodingiz: 5521")
	f.Add("a_b@c")
	f.Add("chegara {[|]} ~ \\ ^ belgi")
	f.Fuzz(func(t *testing.T, s string) {
		enc, err := EncodeGSM7(s)
		if err != nil {
			return // GSM7'da yo'q belgi — normal rad
		}
		dec, err := DecodeGSM7(enc)
		if err != nil {
			t.Fatalf("encode o'tdi-yu decode yiqildi: %v", err)
		}
		if dec != s {
			t.Fatalf("round-trip buzildi: %q -> %q", s, dec)
		}
		// Pack/Unpack ham xuddi shu invariantda (septet darajasi).
		packed := Pack(enc)
		unpacked, err := Unpack(packed, len(enc))
		if err != nil || string(unpacked) != string(enc) {
			t.Fatalf("pack round-trip buzildi (len=%d, err=%v)", len(enc), err)
		}
	})
}

// FuzzUCS2RoundTrip — UTF-16BE codec: valid UTF-8 kirish uchun round-trip.
func FuzzUCS2RoundTrip(f *testing.F) {
	f.Add("Salom")
	f.Add("emoji \U0001F600 aralash")
	f.Fuzz(func(t *testing.T, s string) {
		if !utf8.ValidString(s) {
			return
		}
		dec, err := DecodeUCS2(EncodeUCS2(s))
		if err != nil {
			t.Fatalf("decode: %v", err)
		}
		if dec != s {
			t.Fatalf("round-trip: %q -> %q", s, dec)
		}
	})
}

// FuzzSplit — segmentlash hech qachon panic bermasligi va budjetlarni
// buzmasligi kerak. DIQQAT — birinchi urinishda invariant XATO yozilgan edi
// ("hamma segment ≤140 oktet"): unpacked GSM7'da (7-bob) 1 belgi = 1 oktet,
// yakka 160-belgili segment = 160 oktet — QONUNIY (140 limiti PACKED havo
// interfeysiga tegishli, sm_length chegarasi esa 254). Fuzzer 145 belgilik
// matn bilan buni bir soniyada eslatdi — corpus'da regression sifatida turibdi.
func FuzzSplit(f *testing.F) {
	f.Add("qisqa")
	f.Add(string(make([]byte, 500)))
	f.Fuzz(func(t *testing.T, s string) {
		if !utf8.ValidString(s) {
			return
		}
		dc, _ := Choose(s)
		segs, err := Split(Normalize(s), dc, MethodUDH8, 0x42)
		if err != nil {
			return
		}
		for i, seg := range segs {
			// sm_length chegarasi — mutlaq (§5.2.21).
			if len(seg.Data) > 254 {
				t.Fatalf("segment %d: %d oktet (sm_length max 254)", i, len(seg.Data))
			}
			// UCS2 baytlari havo limitiga to'g'ridan-to'g'ri boradi.
			if dc == DCUCS2 && len(seg.Data) > 140 {
				t.Fatalf("UCS2 segment %d: %d oktet (max 140)", i, len(seg.Data))
			}
			// Multi-segment GSM7 (unpacked): UDH 6 + 153 belgi = max 159.
			if dc == DCDefault && len(segs) > 1 && len(seg.Data) > 6+udh8GSM7 {
				t.Fatalf("GSM7 segment %d: %d oktet (max %d)", i, len(seg.Data), 6+udh8GSM7)
			}
		}
	})
}
