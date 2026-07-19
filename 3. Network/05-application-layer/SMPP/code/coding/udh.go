package coding

import "fmt"

// UDH — concatenated SMS uchun User Data Header'ning concat IE'si
// (TS 23.040 §9.2.3.24): xabar qismlarini telefonda qayta yig'ish uchun
// reference + total + seqnum. UDH short_message'ning ENG BOSHIGA yoziladi
// va esm_class'da UDHI (0x40) bayrog'i majburiy (v3.4 §5.2.12 Notes).
type UDH struct {
	RefNum  uint16 // barcha segmentlarda BIR XIL; 8-bit rejimda faqat quyi bayt
	Total   uint8  // jami segmentlar (1–255)
	Seq     uint8  // shu segment raqami, 1'dan boshlanadi
	Is16bit bool   // true → IEI 0x08 (2 baytlik ref), false → IEI 0x00
}

// UDH IE identifikatorlari (TS 23.040 §9.2.3.24).
const (
	ieiConcat8  = 0x00 // 8-bit reference — de-fakto standart
	ieiConcat16 = 0x08 // 16-bit reference — kam ishlatiladi
)

// Len — to'liq UDH uzunligi oktetlarda, UDHL baytini ham qo'shib: 6 yoki 7.
func (u UDH) Len() int {
	if u.Is16bit {
		return 7
	}
	return 6
}

// Encode UDH baytlarini qaytaradi:
//
//	8-bit:  05 00 03 RR TT SS
//	16-bit: 06 08 04 RR RR TT SS
func (u UDH) Encode() []byte {
	if u.Is16bit {
		return []byte{0x06, ieiConcat16, 0x04, byte(u.RefNum >> 8), byte(u.RefNum), u.Total, u.Seq}
	}
	return []byte{0x05, ieiConcat8, 0x03, byte(u.RefNum), u.Total, u.Seq}
}

// ParseUDH short_message boshidan UDH'ni o'qiydi (esm_class'da UDHI
// o'rnatilgan bo'lsa chaqiriladi). Qaytaradi: topilgan concat ma'lumoti
// (found=false — UDH bor-u, concat IE'si yo'q: masalan faqat port addressing),
// header'dan KEYINGI payload baytlari va xato. Notanish IE'lar jimgina
// o'tkazib yuboriladi (UDH'ning o'z forward-compatibility qoidasi).
func ParseUDH(sm []byte) (udh UDH, payload []byte, found bool, err error) {
	if len(sm) == 0 {
		return UDH{}, nil, false, fmt.Errorf("coding: bo'sh short_message'da UDH yo'q")
	}
	udhl := int(sm[0])
	if 1+udhl > len(sm) {
		return UDH{}, nil, false, fmt.Errorf("coding: UDHL=%d, lekin short_message %d oktet", udhl, len(sm))
	}
	payload = sm[1+udhl:]
	ie := sm[1 : 1+udhl]
	for off := 0; off < len(ie); {
		if len(ie)-off < 2 {
			return UDH{}, nil, false, fmt.Errorf("coding: UDH IE header'i chala (%d oktet qoldi)", len(ie)-off)
		}
		iei, iedl := ie[off], int(ie[off+1])
		off += 2
		if len(ie)-off < iedl {
			return UDH{}, nil, false, fmt.Errorf("coding: IE 0x%02X uzunligi %d, qolgan %d", iei, iedl, len(ie)-off)
		}
		data := ie[off : off+iedl]
		off += iedl
		switch {
		case iei == ieiConcat8 && iedl == 3:
			udh = UDH{RefNum: uint16(data[0]), Total: data[1], Seq: data[2]}
			found = true
		case iei == ieiConcat16 && iedl == 4:
			udh = UDH{RefNum: uint16(data[0])<<8 | uint16(data[1]), Total: data[2], Seq: data[3], Is16bit: true}
			found = true
		}
		// boshqa IE'lar (port addressing 0x05 va h.k.) — e'tiborsiz
	}
	return udh, payload, found, nil
}
