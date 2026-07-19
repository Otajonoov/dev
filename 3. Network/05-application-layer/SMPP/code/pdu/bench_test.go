package pdu

import (
	"bytes"
	"testing"
)

// Hot path benchmark'lari: gateway'da har xabar shu yo'ldan o'tadi.
// Yugurish: go test ./pdu -bench=. -benchmem

func benchSubmit() SubmitSM {
	return SubmitSM{SMFields: SMFields{
		Source:             Address{TON: TONAlphanumeric, Addr: "Bank"},
		Dest:               Address{TON: TONInternational, NPI: NPIISDN, Addr: "998901234567"},
		RegisteredDelivery: DLRFinal,
		ShortMessage:       []byte("Assalomu alaykum! Kodingiz: 5521"),
	}}
}

func BenchmarkSubmitSMEncode(b *testing.B) {
	sm := benchSubmit()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := sm.Encode(uint32(i%maxInt31 + 1)); err != nil {
			b.Fatal(err)
		}
	}
}

const maxInt31 = 0x7FFFFFFE

func BenchmarkDecodeSubmitSM(b *testing.B) {
	frame, err := benchSubmit().Encode(42)
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, _, err := DecodeSubmitSM(frame); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodeDispatcher(b *testing.B) {
	frame, err := benchSubmit().Encode(42)
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, _, err := Decode(frame); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkReadFrame(b *testing.B) {
	frame, err := benchSubmit().Encode(42)
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	r := bytes.NewReader(frame)
	for i := 0; i < b.N; i++ {
		r.Reset(frame)
		if _, err := ReadFrame(r, 64*1024); err != nil {
			b.Fatal(err)
		}
	}
}
