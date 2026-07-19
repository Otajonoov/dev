package pdu

import (
	"bytes"
	"encoding/hex"
	"strings"
	"testing"
)

// FuzzDecode — dispatcher va barcha PDU decoder'lari uchun fuzz target.
// Tamoyil: "Never trust incoming data" — istalgan bayt to'plami PANIC EMAS,
// yo (PDU, nil) yo xato qaytarishi kerak. command_length manipulyatsiyasi,
// kesik body'lar, buzuq TLV tail — hammasi shu yerdan o'tadi.
func FuzzDecode(f *testing.F) {
	// Seed corpus: har PDU turidan kamida bitta valid frame — busiz fuzzer
	// PDU strukturasigacha "yetib bormaydi".
	seeds := []string{
		specBindTransmitterHex,
		goldenDataSMHex,
		goldenDataSMRespHex,
		goldenSubmitMultiHex,
		goldenSubmitMultiRespHex,
		goldenQuerySMHex,
		goldenQuerySMRespHex,
	}
	for _, s := range seeds {
		// specBindTransmitterHex ko'p qatorli — whitespace olib tashlanadi.
		b, err := hex.DecodeString(strings.ToLower(strings.Join(strings.Fields(s), "")))
		if err != nil {
			f.Fatal(err)
		}
		f.Add(b)
	}
	// Header-only'lar va ataylab buzuqlar.
	f.Add(EncodeEnquireLink(1))
	f.Add(EncodeGenericNack(3, 0))
	f.Add([]byte{0x00, 0x00, 0x00, 0x0C}) // length < 16
	f.Add(bytes.Repeat([]byte{0xFF}, 32))

	f.Fuzz(func(t *testing.T, data []byte) {
		p, h, err := Decode(data)
		if err != nil {
			return // xato — normal natija
		}
		// Muvaffaqiyatli decode bo'lsa, invariantlar tekshiriladi.
		if p.Cmd() != h.ID {
			t.Fatalf("Cmd()=%s, header=%s", p.Cmd(), h.ID)
		}
	})
}

// FuzzReadFrame — framing qatlami: buzuq length'lar OOM/panic bermasligi.
func FuzzReadFrame(f *testing.F) {
	f.Add(EncodeEnquireLink(7))
	f.Add([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0x00}) // absurd katta length
	f.Add([]byte{0x00, 0x00, 0x00, 0x08})       // length < 16
	f.Fuzz(func(t *testing.T, data []byte) {
		frame, err := ReadFrame(bytes.NewReader(data), 64*1024)
		if err != nil {
			return
		}
		if len(frame) < HeaderSize || len(frame) > 64*1024 {
			t.Fatalf("ReadFrame chegara buzdi: %d bayt", len(frame))
		}
	})
}
