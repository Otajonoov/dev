package coding

import "fmt"

// Latin-1 (ISO-8859-1, data_coding=3): 1 belgi = 1 oktet, kod nuqtasi = bayt
// qiymati (U+0000–U+00FF). 140 oktet = 140 belgi. G'arbiy Yevropa tillari
// uchun; kirill va o'zbek maxsus harflari SIG'MAYDI.

// EncodeLatin1 matnni ISO-8859-1 baytlarga o'tkazadi.
func EncodeLatin1(s string) ([]byte, error) {
	out := make([]byte, 0, len(s))
	for _, r := range s {
		if r > 0xFF {
			return nil, fmt.Errorf("coding: %q Latin-1'da yo'q (U+%04X > U+00FF)", r, r)
		}
		out = append(out, byte(r))
	}
	return out, nil
}

// DecodeLatin1 ISO-8859-1 baytlarni matnga qaytaradi (har bayt = kod nuqtasi;
// xato holati yo'q — barcha 256 qiymat aniqlangan).
func DecodeLatin1(b []byte) string {
	out := make([]rune, len(b))
	for i, c := range b {
		out[i] = rune(c)
	}
	return string(out)
}
