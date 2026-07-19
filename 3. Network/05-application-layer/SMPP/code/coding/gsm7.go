// Package coding SMS matn encoding'larini beradi: GSM 03.38 7-bit alphabet
// (unpacked va packed), Latin-1, UCS2 hamda data_coding tanlash/normalizatsiya.
// Normativ manba: 3GPP TS 23.038 (GSM 03.38 vorisi).
package coding

import "fmt"

// esc — extension table'ga o'tish belgisi (TS 23.038 §6.2.1.1).
const esc = 0x1B

// gsm7Basic — TS 23.038 §6.2.1 default alphabet: indeks = GSM kodi (0x00–0x7F).
// 0x1B pozitsiyasi ESC — belgi emas, extension prefiksi (qiymati 0 qoldirilgan).
var gsm7Basic = [128]rune{
	'@', '£', '$', '¥', 'è', 'é', 'ù', 'ì', 'ò', 'Ç', '\n', 'Ø', 'ø', '\r', 'Å', 'å',
	'Δ', '_', 'Φ', 'Γ', 'Λ', 'Ω', 'Π', 'Ψ', 'Σ', 'Θ', 'Ξ', 0, 'Æ', 'æ', 'ß', 'É',
	' ', '!', '"', '#', '¤', '%', '&', '\'', '(', ')', '*', '+', ',', '-', '.', '/',
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', ':', ';', '<', '=', '>', '?',
	'¡', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O',
	'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', 'Ä', 'Ö', 'Ñ', 'Ü', '§',
	'¿', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o',
	'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'ä', 'ö', 'ñ', 'ü', 'à',
}

// gsm7Ext — extension table (TS 23.038 §6.2.1.1): ESC + kod → belgi.
// 10 ta belgi; har biri simda 2 septet egallaydi.
var gsm7Ext = map[byte]rune{
	0x0A: '\f', // Form Feed (sahifa uzilishi)
	0x14: '^',
	0x28: '{',
	0x29: '}',
	0x2F: '\\',
	0x3C: '[',
	0x3D: '~',
	0x3E: ']',
	0x40: '|',
	0x65: '€',
}

// Teskari jadvallar (rune → GSM kodi), init'da quriladi.
var (
	gsm7BasicRev = make(map[rune]byte, 127)
	gsm7ExtRev   = make(map[rune]byte, len(gsm7Ext))
)

func init() {
	for code, r := range gsm7Basic {
		if code == esc {
			continue
		}
		gsm7BasicRev[r] = byte(code)
	}
	for code, r := range gsm7Ext {
		gsm7ExtRev[r] = code
	}
}

// IsGSM7 r GSM 7-bit alphabet'da (basic yoki extension) bor-yo'qligini aytadi.
func IsGSM7(r rune) bool {
	if _, ok := gsm7BasicRev[r]; ok {
		return true
	}
	_, ok := gsm7ExtRev[r]
	return ok
}

// SeptetLen matnning septet hisobidagi uzunligi: basic belgi = 1 septet,
// extension belgi (€, {, [, ...) = 2 septet. Segment matematikasi (8-bob)
// aynan shu son ustida ishlaydi.
func SeptetLen(s string) (int, error) {
	n := 0
	for _, r := range s {
		switch {
		case gsm7ExtRev[r] != 0:
			n += 2
		default:
			if _, ok := gsm7BasicRev[r]; !ok {
				return 0, fmt.Errorf("coding: %q GSM7 alphabet'da yo'q", r)
			}
			n++
		}
	}
	return n, nil
}

// EncodeGSM7 matnni UNPACKED GSM7 baytlarga o'tkazadi: 1 septet = 1 oktet
// (yuqori bit 0), extension belgi = ESC + kod (2 oktet). SMPP short_message
// amalda deyarli hamisha shu ko'rinishda yuriladi; packing'ni SMSC o'zi qiladi.
func EncodeGSM7(s string) ([]byte, error) {
	out := make([]byte, 0, len(s))
	for _, r := range s {
		if code, ok := gsm7BasicRev[r]; ok {
			out = append(out, code)
			continue
		}
		if code, ok := gsm7ExtRev[r]; ok {
			out = append(out, esc, code)
			continue
		}
		return nil, fmt.Errorf("coding: %q GSM7 alphabet'da yo'q", r)
	}
	return out, nil
}

// DecodeGSM7 unpacked GSM7 baytlarni matnga qaytaradi. Bayt >0x7F — xato
// (unpacked GSM7 7-bitli). Notanish extension kodi TS 23.038 tavsiyasi
// bo'yicha bo'sh joy sifatida o'qiladi (xato emas — eski telefon xulqi).
func DecodeGSM7(b []byte) (string, error) {
	out := make([]rune, 0, len(b))
	for i := 0; i < len(b); i++ {
		c := b[i]
		if c > 0x7F {
			return "", fmt.Errorf("coding: 0x%02X unpacked GSM7'da bo'lishi mumkin emas (bit 7 o'rnatilgan)", c)
		}
		if c == esc {
			i++
			if i >= len(b) {
				return "", fmt.Errorf("coding: ESC bilan tugagan GSM7 oqimi")
			}
			if r, ok := gsm7Ext[b[i]]; ok {
				out = append(out, r)
			} else {
				out = append(out, ' ')
			}
			continue
		}
		out = append(out, gsm7Basic[c])
	}
	return string(out), nil
}

// Pack septet'larni TS 23.038 §6.1.2.1 bo'yicha zichlaydi: har septetning
// 7 biti ketma-ket teriladi, 8 septet 7 oktetga sig'adi. SMPP'da odatda
// KERAK EMAS (unpacked yuboriladi) — havo interfeysi formati va TP-OA
// (6-bob) ni tushunish uchun.
func Pack(septets []byte) []byte {
	if len(septets) == 0 {
		return nil
	}
	out := make([]byte, 0, (len(septets)*7+7)/8)
	var acc uint16
	accBits := 0
	for _, s := range septets {
		acc |= uint16(s&0x7F) << accBits
		accBits += 7
		for accBits >= 8 {
			out = append(out, byte(acc))
			acc >>= 8
			accBits -= 8
		}
	}
	if accBits > 0 {
		out = append(out, byte(acc))
	}
	return out
}

// Unpack zichlangan oqimni septet'larga yoyadi; n — kutilayotgan septet soni
// (packed formatda uzunlik alohida tashiladi — masalan TP-UDL).
func Unpack(data []byte, n int) ([]byte, error) {
	if n < 0 || n > len(data)*8/7 {
		return nil, fmt.Errorf("coding: %d oktetdan %d septet chiqmaydi", len(data), n)
	}
	out := make([]byte, 0, n)
	var acc uint16
	accBits := 0
	for _, b := range data {
		acc |= uint16(b) << accBits
		accBits += 8
		for accBits >= 7 && len(out) < n {
			out = append(out, byte(acc&0x7F))
			acc >>= 7
			accBits -= 7
		}
	}
	if len(out) < n {
		return nil, fmt.Errorf("coding: data tugadi — %d septetdan %d o'qildi", n, len(out))
	}
	return out, nil
}
