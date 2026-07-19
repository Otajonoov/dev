package coding

import (
	"bytes"
	"testing"
)

func TestUCS2CyrillicGolden(t *testing.T) {
	// "Салом" — kirill, UCS2 big-endian golden hex.
	got := EncodeUCS2("Салом")
	want := mustHex(t, "04 21 04 30 04 3B 04 3E 04 3C")
	if !bytes.Equal(got, want) {
		t.Fatalf("EncodeUCS2 = % X,\nkutilgan % X", got, want)
	}
	back, err := DecodeUCS2(want)
	if err != nil || back != "Салом" {
		t.Errorf("DecodeUCS2 = %q, %v", back, err)
	}
}

func TestUCS2SurrogatePair(t *testing.T) {
	// Emoji (U+1F600) — BMP tashqarisida: surrogate pair = 4 oktet.
	got := EncodeUCS2("😀")
	want := mustHex(t, "D8 3D DE 00")
	if !bytes.Equal(got, want) {
		t.Fatalf("EncodeUCS2(emoji) = % X, kutilgan % X", got, want)
	}
	back, err := DecodeUCS2(got)
	if err != nil || back != "😀" {
		t.Errorf("DecodeUCS2(emoji) = %q, %v", back, err)
	}
}

func TestUCS2OddLength(t *testing.T) {
	if _, err := DecodeUCS2([]byte{0x04, 0x21, 0x04}); err == nil {
		t.Error("toq uzunlikdagi UCS2 xato berishi kerak edi")
	}
}

func TestLatin1RoundTrip(t *testing.T) {
	s := "Grüße für 99¢?" // Latin-1 diapazonidagi belgilar
	b, err := EncodeLatin1(s)
	if err != nil {
		t.Fatalf("EncodeLatin1: %v", err)
	}
	if len(b) != 14 {
		t.Errorf("Latin-1'da 1 belgi = 1 oktet: %d oktet, kutilgan 14", len(b))
	}
	if back := DecodeLatin1(b); back != s {
		t.Errorf("round-trip: %q != %q", back, s)
	}
	if _, err := EncodeLatin1("Салом"); err == nil {
		t.Error("kirill Latin-1 ga sig'masligi kerak edi")
	}
}

func TestNormalize(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"U+02BB (oʻ/gʻ rasmiy)", "oʻzbekcha gʻoya", "o'zbekcha g'oya"},
		{"U+02BC (tutuq)", "maʼno", "ma'no"},
		{"aqlli qo'shtirnoqlar", "s‘z ’z", "s'z 'z"},
		{"ASCII o'zgarmaydi", "o'zbekcha", "o'zbekcha"},
		{"kirillga tegilmaydi", "Салом", "Салом"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Normalize(tt.in); got != tt.want {
				t.Errorf("Normalize(%q) = %q, kutilgan %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestChooseUzbekLatin(t *testing.T) {
	// Rasmiy imlodagi lotin matn: normalize'dan KEYIN to'liq GSM7 —
	// DCDefault tanlanadi, apostrof ASCII 0x27 bo'lib ketadi.
	dc, b := Choose("oʻzbekcha matn")
	if dc != DCDefault {
		t.Fatalf("dc = 0x%02X, kutilgan DCDefault", uint8(dc))
	}
	want, _ := EncodeGSM7("o'zbekcha matn")
	if !bytes.Equal(b, want) {
		t.Errorf("baytlar: % X != % X", b, want)
	}
}

func TestChooseCyrillic(t *testing.T) {
	dc, b := Choose("Ассалому алайкум")
	if dc != DCUCS2 {
		t.Fatalf("kirill uchun dc = 0x%02X, kutilgan DCUCS2", uint8(dc))
	}
	if len(b) != len([]rune("Ассалому алайкум"))*2 {
		t.Errorf("UCS2 uzunligi: %d oktet", len(b))
	}
}

func TestChooseNormalizeOrderMatters(t *testing.T) {
	// Tartib isboti: normalize'siz bu matn UCS2'ga ketardi (U+02BB GSM7'da yo'q),
	// Choose esa avval normalize qilib GSM7'da qoldiradi.
	dc, _ := Choose("gʻisht")
	if dc != DCDefault {
		t.Errorf("normalize'dan keyin GSM7 kutilgan edi, dc=0x%02X", uint8(dc))
	}
}
