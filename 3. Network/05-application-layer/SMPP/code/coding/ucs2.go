package coding

import (
	"fmt"
	"unicode/utf16"
)

// UCS2 (data_coding=8): 2 oktetlik unit'lar, BIG-ENDIAN, BOM'siz.
// 140 oktet = 70 unit. Rasmiy UCS2 faqat BMP (U+0000–U+FFFF), lekin amalda
// zanjir UTF-16 sifatida ishlaydi: BMP'dan tashqari belgilar (emoji) surrogate
// pair bo'lib ketadi — 1 belgi = 2 unit = 4 oktet.

// EncodeUCS2 matnni UTF-16BE baytlarga o'tkazadi (surrogate-safe).
func EncodeUCS2(s string) []byte {
	units := utf16.Encode([]rune(s))
	out := make([]byte, len(units)*2)
	for i, u := range units {
		out[i*2] = byte(u >> 8)
		out[i*2+1] = byte(u)
	}
	return out
}

// DecodeUCS2 UTF-16BE baytlarni matnga qaytaradi. Toq uzunlik — xato
// (UCS2 oqimi doim juft oktet; toq kelishi kesilgan/buzilgan matn belgisi).
func DecodeUCS2(b []byte) (string, error) {
	if len(b)%2 != 0 {
		return "", fmt.Errorf("coding: UCS2 oqimi %d oktet — juft bo'lishi kerak", len(b))
	}
	units := make([]uint16, len(b)/2)
	for i := range units {
		units[i] = uint16(b[i*2])<<8 | uint16(b[i*2+1])
	}
	return string(utf16.Decode(units)), nil
}
