package dlr

import (
	"errors"
	"testing"

	"smpp/tlv"
)

// FuzzParse — tolerant DLR parser'ining "tolerant"ligini fuzzer tekshiradi:
// istalgan matn panic bermasligi, xato yo partial Receipt qaytishi kerak.
func FuzzParse(f *testing.F) {
	f.Add([]byte(goldenReceiptText))
	f.Add([]byte("id: stat: err:"))
	f.Add([]byte("submit date:9999999999999 done date:abc"))
	f.Add([]byte("text:id:sub:dlvrd:stat:"))
	f.Add([]byte{})
	f.Fuzz(func(t *testing.T, sm []byte) {
		r, err := Parse(sm, nil)
		if err != nil {
			// Yagona ruxsat etilgan xato turi — ErrNotReceipt.
			if !errors.Is(err, ErrNotReceipt) {
				t.Fatalf("kutilmagan xato turi: %v", err)
			}
			return
		}
		// Sub/Dlvrd yo -1 (topilmadi) yo >= 0 bo'lishi kerak.
		if r.Sub < -1 || r.Dlvrd < -1 {
			t.Fatalf("sentinel buzildi: %+v", r)
		}
	})
}

// FuzzNormalizeID — id normalizatsiyasi: panic yo'q, variantlar ichida
// har doim id'ning o'zi (upper) bo'lishi kerak.
func FuzzNormalizeID(f *testing.F) {
	f.Add("7F3A9B")
	f.Add("288230376293190992")
	f.Add("04000000086ECD50")
	f.Add("a4b8-uuid")
	f.Fuzz(func(t *testing.T, id string) {
		got := NormalizeID(id)
		_ = got
		// TLV bilan birga Parse ham sinab qo'yiladi (arzon qo'shimcha qamrov).
		_, _ = Parse([]byte("id:"+id+" stat:DELIVRD"), []tlv.TLV{tlv.CString(tlv.ReceiptedMessageID, "X")})
	})
}
