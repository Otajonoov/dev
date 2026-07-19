package dlr

import (
	"slices"
	"testing"
)

func TestNormalizeID(t *testing.T) {
	tests := []struct {
		in   string
		want []string // kamida shu variantlar bo'lishi kerak
	}{
		// Decimal id → hex va padded-hex variantlar (Kannel #334 case'i,
		// python bilan tekshirilgan: 0x04000000086ECD50 = 288230376293190992).
		{"288230376293190992", []string{"288230376293190992", "4000000086ECD50", "04000000086ECD50"}},
		// Hex id → decimal variant.
		{"04000000086ECD50", []string{"04000000086ECD50", "288230376293190992", "4000000086ECD50"}},
		// Faqat-raqamli qisqa id — ham decimal, ham hex talqin:
		// dec("100")=100 → hex "64"; hex("100")=256 → dec "256".
		{"100", []string{"100", "64", "256"}},
		// Kichik harfli hex katta harfga keltiriladi.
		{"6ecd50", []string{"6ECD50", "7261520"}},
		// Opaque id (raqamga o'xshamaydi) — faqat o'zi.
		{"a4b8-c1d2-uuid", []string{"A4B8-C1D2-UUID"}},
		{"", nil},
	}
	for _, tt := range tests {
		got := NormalizeID(tt.in)
		for _, w := range tt.want {
			if !slices.Contains(got, w) {
				t.Errorf("NormalizeID(%q) = %v — %q varianti yo'q", tt.in, got, w)
			}
		}
		if tt.in == "" && got != nil {
			t.Errorf("bo'sh id uchun nil kutilgan, keldi %v", got)
		}
	}
}

// 9-bob mashqidagi stsenariy: bazada hex, DLR'da decimal.
func TestTableHexDecQuirk(t *testing.T) {
	tab := NewTable()
	tab.Register("04000000086ECD50") // submit_sm_resp shu ko'rinishda berdi

	canon, ok := tab.Resolve("288230376293190992") // DLR decimal keldi
	if !ok || canon != "04000000086ECD50" {
		t.Fatalf("Resolve(dec) = %q, %v — kanonik hex id kutilgan", canon, ok)
	}
	// Teskari yo'nalish ham: bazada decimal, DLR'da hex.
	tab2 := NewTable()
	tab2.Register("288230376293190992")
	canon, ok = tab2.Resolve("04000000086ECD50")
	if !ok || canon != "288230376293190992" {
		t.Fatalf("Resolve(hex) = %q, %v", canon, ok)
	}
}

func TestTableExactAndCase(t *testing.T) {
	tab := NewTable()
	tab.Register("7F3A9B")
	if canon, ok := tab.Resolve("7f3a9b"); !ok || canon != "7F3A9B" {
		t.Fatalf("registr farqi to'sqinlik qildi: %q %v", canon, ok)
	}
	if _, ok := tab.Resolve("DEADBEEF"); ok {
		t.Fatal("notanish id topilmasligi kerak edi")
	}
}

func TestTableForget(t *testing.T) {
	tab := NewTable()
	tab.Register("7F3A9B")
	tab.Forget("7F3A9B")
	if _, ok := tab.Resolve("7F3A9B"); ok {
		t.Fatal("Forget'dan keyin id topilmasligi kerak")
	}
	if len(tab.m) != 0 {
		t.Fatalf("jadval bo'sh emas: %v", tab.m)
	}
}
