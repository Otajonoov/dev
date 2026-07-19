package coding

import (
	"bytes"
	"encoding/hex"
	"strings"
	"testing"
)

func mustHex(t *testing.T, s string) []byte {
	t.Helper()
	clean := strings.NewReplacer(" ", "", "\n", "").Replace(s)
	b, err := hex.DecodeString(clean)
	if err != nil {
		t.Fatalf("hex decode xatosi: %v", err)
	}
	return b
}

func TestGSM7BasicAlphabetRoundTrip(t *testing.T) {
	// Barcha 127 basic belgi (ESC'dan tashqari) yakka holda round-trip.
	for code := 0; code < 128; code++ {
		if code == esc {
			continue
		}
		r := gsm7Basic[code]
		s, err := DecodeGSM7([]byte{byte(code)})
		if err != nil {
			t.Fatalf("DecodeGSM7(0x%02X): %v", code, err)
		}
		if s != string(r) {
			t.Errorf("DecodeGSM7(0x%02X) = %q, kutilgan %q", code, s, string(r))
		}
		enc, err := EncodeGSM7(string(r))
		if err != nil {
			t.Fatalf("EncodeGSM7(%q): %v", string(r), err)
		}
		if len(enc) != 1 || enc[0] != byte(code) {
			t.Errorf("EncodeGSM7(%q) = % X, kutilgan %02X", string(r), enc, code)
		}
	}
}

func TestGSM7ExtensionChars(t *testing.T) {
	// 10 extension belgi: har biri 2 oktet (ESC + kod), 2 septet.
	extChars := "\f^{}\\[~]|€"
	for _, r := range extChars {
		enc, err := EncodeGSM7(string(r))
		if err != nil {
			t.Fatalf("EncodeGSM7(%q): %v", string(r), err)
		}
		if len(enc) != 2 || enc[0] != esc {
			t.Errorf("EncodeGSM7(%q) = % X, ESC+kod kutilgan", string(r), enc)
		}
		dec, err := DecodeGSM7(enc)
		if err != nil || dec != string(r) {
			t.Errorf("round-trip(%q) = %q, %v", string(r), dec, err)
		}
		if n, err := SeptetLen(string(r)); err != nil || n != 2 {
			t.Errorf("SeptetLen(%q) = %d, %v; kutilgan 2", string(r), n, err)
		}
	}
	// € ning aniq kodi — spec misoli (TS 23.038: 0x1B 0x65).
	enc, _ := EncodeGSM7("€")
	if !bytes.Equal(enc, []byte{0x1B, 0x65}) {
		t.Errorf("€ = % X, kutilgan 1B 65", enc)
	}
}

func TestGSM7AsciiDiffGolden(t *testing.T) {
	// GSM7 ≠ ASCII isboti: "a_b@c" — harflar mos, _ va @ butunlay boshqa kodda.
	enc, err := EncodeGSM7("a_b@c")
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{0x61, 0x11, 0x62, 0x00, 0x63}
	if !bytes.Equal(enc, want) {
		t.Fatalf("EncodeGSM7(a_b@c) = % X, kutilgan % X", enc, want)
	}
	dec, err := DecodeGSM7(want)
	if err != nil || dec != "a_b@c" {
		t.Errorf("DecodeGSM7 = %q, %v", dec, err)
	}
}

func TestGSM7NotInAlphabet(t *testing.T) {
	for _, s := range []string{"Салом", "ʻo'zbek", "😀"} {
		if _, err := EncodeGSM7(s); err == nil {
			t.Errorf("EncodeGSM7(%q): xato kutilgan edi", s)
		}
	}
	if IsGSM7('ʻ') || IsGSM7('С') {
		t.Error("U+02BB va kirill GSM7'da bo'lmasligi kerak")
	}
	if !IsGSM7('@') || !IsGSM7('€') || !IsGSM7('\'') {
		t.Error("@, €, ' GSM7'da bo'lishi kerak")
	}
}

func TestDecodeGSM7Errors(t *testing.T) {
	if _, err := DecodeGSM7([]byte{0x80}); err == nil {
		t.Error("bit 7 o'rnatilgan bayt xato berishi kerak edi")
	}
	if _, err := DecodeGSM7([]byte{esc}); err == nil {
		t.Error("ESC bilan tugagan oqim xato berishi kerak edi")
	}
	// Notanish extension kodi → bo'sh joy (TS 23.038 tolerant xulqi).
	s, err := DecodeGSM7([]byte{esc, 0x41})
	if err != nil || s != " " {
		t.Errorf("notanish ext kod: %q, %v; kutilgan bo'sh joy", s, err)
	}
}

func TestPackHelloGolden(t *testing.T) {
	// Klassik reference misol: "hello" → E8 32 9B FD 06.
	septets, err := EncodeGSM7("hello")
	if err != nil {
		t.Fatal(err)
	}
	packed := Pack(septets)
	want := mustHex(t, "E8 32 9B FD 06")
	if !bytes.Equal(packed, want) {
		t.Fatalf("Pack(hello) = % X, kutilgan % X", packed, want)
	}
	// 8 belgi → 7 oktet qoidasi.
	eight, _ := EncodeGSM7("hihihihi")
	if p := Pack(eight); len(p) != 7 {
		t.Errorf("8 septet %d oktetga zichlandi, kutilgan 7", len(p))
	}
}

func TestPackUnpackRoundTrip(t *testing.T) {
	texts := []string{"", "a", "hello", "Salom dunyo!", "12345678", "€uro {test}"}
	for _, s := range texts {
		septets, err := EncodeGSM7(s)
		if err != nil {
			t.Fatal(err)
		}
		unpacked, err := Unpack(Pack(septets), len(septets))
		if err != nil {
			t.Fatalf("Unpack(%q): %v", s, err)
		}
		if !bytes.Equal(unpacked, septets) {
			t.Errorf("pack/unpack(%q): % X != % X", s, unpacked, septets)
		}
	}
}

func TestUnpackErrors(t *testing.T) {
	if _, err := Unpack([]byte{0xE8}, 3); err == nil {
		t.Error("1 oktetdan 3 septet chiqmasligi kerak edi")
	}
}

func TestSeptetLen(t *testing.T) {
	tests := []struct {
		s    string
		want int
	}{
		{"", 0},
		{"hello", 5},
		{"€", 2}, // extension = 2 septet
		{"5€", 3},
		{"Salom { }", 11}, // 7 basic + 2×2 ext
	}
	for _, tt := range tests {
		if n, err := SeptetLen(tt.s); err != nil || n != tt.want {
			t.Errorf("SeptetLen(%q) = %d, %v; kutilgan %d", tt.s, n, err, tt.want)
		}
	}
	if _, err := SeptetLen("кирилл"); err == nil {
		t.Error("kirill matnda SeptetLen xato berishi kerak edi")
	}
}
