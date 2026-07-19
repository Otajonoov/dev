package pdu

import (
	"bytes"
	"fmt"
	"strings"
	"unicode/utf8"
)

// Address — SMPP manzil uchligi: TON + NPI + manzil matni (§5.2.5–5.2.9).
// 5-bobda codec darajasida kiritilgan edi; bu bobda (6-bob) konstruktorlar
// va Validate qo'shildi. Codec semantikani TEKSHIRMAYDI (har baytni o'tkazadi);
// Validate esa yuborishdan oldin chaqiriladigan ixtiyoriy semantik filtr.
type Address struct {
	TON  uint8
	NPI  uint8
	Addr string
}

// TON qiymatlari (§5.2.5, Table 5-3).
const (
	TONUnknown         uint8 = 0x00
	TONInternational   uint8 = 0x01 // to'liq xalqaro format (country code bilan)
	TONNational        uint8 = 0x02 // milliy format (country code'siz)
	TONNetworkSpecific uint8 = 0x03 // operator ichki raqami / short code konventsiyasi
	TONSubscriber      uint8 = 0x04
	TONAlphanumeric    uint8 = 0x05 // harf-raqamli sender ("Bank")
	TONAbbreviated     uint8 = 0x06 // qisqartirilgan raqam
)

// NPI qiymatlari (§5.2.6, Table 5-4). DIQQAT: qiymatlar KETMA-KET EMAS —
// 1'dan keyin to'g'ridan-to'g'ri 3 keladi (2 yo'q!). "Tartib bilan sanalgan"
// deb taxmin qilib enum yozish — klassik xato.
const (
	NPIUnknown    uint8 = 0x00
	NPIISDN       uint8 = 0x01 // E.163/E.164 — oddiy telefon raqamlari
	NPIData       uint8 = 0x03 // X.121
	NPITelex      uint8 = 0x04 // F.69
	NPILandMobile uint8 = 0x06 // E.212
	NPINational   uint8 = 0x08
	NPIPrivate    uint8 = 0x09
	NPIERMES      uint8 = 0x0A
	NPIInternet   uint8 = 0x0E // IP manzil, "aaa.bbb.ccc.ddd" ko'rinishda; IPv6 v3.4'da YO'Q
	NPIWAP        uint8 = 0x12
)

// maxAddrShort — submit_sm/deliver_sm'dagi manzil o'lchami (§5.2.8–5.2.9:
// max 21, NULL bilan). data_sm'da limit 65 — u 10-bobda alohida.
const maxAddrShort = 21

// maxAlphanumeric — alphanumeric sender limiti: 11 belgi. Bu SMPP limiti EMAS
// (SMPP field'i 20 belgigacha sig'diradi) — havo interfeysi limiti:
// TS 23.040 §9.1.2.5 TP-OA'da alphanumeric manzil max 10 oktet, 7-bit
// packing bilan 11 ta GSM7 belgi.
const maxAlphanumeric = 11

// gsm7SenderPunct — GSM7 default alphabet'ning ASCII tinish belgilari qismi.
// To'liq GSM7 jadvali 7-bobda (coding package); sender tekshiruvi uchun ASCII
// qismi yetarli — GSM7'ning ASCII bo'lmagan belgilari (è, Δ...) sender'larda
// amalda ishlatilmaydi, ko'p operatorlar esa bundan ham torroq ro'yxat talab qiladi.
const gsm7SenderPunct = "@$ !\"#%&'()*+,-./:;<=>?_"

func isGSM7SenderRune(r rune) bool {
	switch {
	case r >= 'A' && r <= 'Z', r >= 'a' && r <= 'z', r >= '0' && r <= '9':
		return true
	}
	return strings.ContainsRune(gsm7SenderPunct, r)
}

func allDigits(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

// International to'liq xalqaro (E.164) raqamdan manzil yasaydi: TON=1/NPI=1.
// Bosh '+' OLIB TASHLANADI — ko'p SMSC'lar TON=1 bilan faqat raqam kutadi,
// '+' bilan kelganini ESME_RINVDSTADR (0x0B) bilan rad etadi.
func International(msisdn string) (Address, error) {
	s := strings.TrimPrefix(msisdn, "+")
	if s == "" {
		return Address{}, fmt.Errorf("pdu: bo'sh msisdn")
	}
	if !allDigits(s) {
		return Address{}, fmt.Errorf("pdu: international manzilda raqam bo'lmagan belgi: %q", msisdn)
	}
	if len(s) > 15 {
		return Address{}, fmt.Errorf("pdu: %q — E.164 max 15 raqam, keldi %d", s, len(s))
	}
	return Address{TON: TONInternational, NPI: NPIISDN, Addr: s}, nil
}

// Alphanumeric brend-sender yasaydi: TON=5/NPI=0 (industriya konventsiyasi).
// Max 11 belgi, faqat GSM7 belgilar (TS 23.040 TP-OA limiti). Eslatma:
// bunday sender'ga abonent JAVOB YOZA OLMAYDI va ko'p mamlakatlarda u
// oldindan ro'yxatdan o'tkazilishi shart.
func Alphanumeric(name string) (Address, error) {
	a := Address{TON: TONAlphanumeric, NPI: NPIUnknown, Addr: name}
	if err := a.Validate(); err != nil {
		return Address{}, err
	}
	return a, nil
}

// ShortCode qisqa raqamdan manzil yasaydi: TON=3/NPI=0 — ENG KENG TARQALGAN
// konventsiya, lekin standartlashmagan: ayrim operatorlar 6/0 yoki 0/1 kutadi —
// integratsiya hujjatidan tekshiring.
func ShortCode(s string) (Address, error) {
	if s == "" || !allDigits(s) {
		return Address{}, fmt.Errorf("pdu: short code faqat raqamlardan iborat bo'lishi kerak: %q", s)
	}
	if len(s) > 8 {
		return Address{}, fmt.Errorf("pdu: short code %d raqam — 8'dan oshiq qiymat short code emas", len(s))
	}
	return Address{TON: TONNetworkSpecific, NPI: NPIUnknown, Addr: s}, nil
}

// NullSource — bo'sh source manzil (TON=0/NPI=0/addr=""): SMSC hisobga
// bog'langan default sender'ni o'zi qo'yadi (§5.2.8).
func NullSource() Address { return Address{} }

// Validate manzilni TON'iga mos semantik qoidalar bilan tekshiradi.
// Codec bu tekshiruvni CHAQIRMAYDI (kelgan har qanday baytni o'qish kerak);
// yuborish yo'lida esa (13-bob client) xatoni SMSC'dan emas, lokal olish afzal.
func (a Address) Validate() error {
	if len(a.Addr)+1 > maxAddrShort {
		return fmt.Errorf("pdu: manzil %d belgi — max %d (NULL bilan %d oktet)", len(a.Addr), maxAddrShort-1, maxAddrShort)
	}
	switch a.TON {
	case TONAlphanumeric:
		if a.Addr == "" {
			return fmt.Errorf("pdu: alphanumeric sender bo'sh bo'lmaydi")
		}
		if n := utf8.RuneCountInString(a.Addr); n > maxAlphanumeric {
			return fmt.Errorf("pdu: alphanumeric sender %d belgi — max %d (TS 23.040 TP-OA)", n, maxAlphanumeric)
		}
		for _, r := range a.Addr {
			if !isGSM7SenderRune(r) {
				return fmt.Errorf("pdu: alphanumeric sender'da GSM7'da bo'lmagan belgi: %q", r)
			}
		}
	case TONInternational:
		if a.Addr == "" {
			return fmt.Errorf("pdu: international manzil bo'sh bo'lmaydi")
		}
		if strings.HasPrefix(a.Addr, "+") {
			return fmt.Errorf("pdu: TON=International bilan '+' yuborilmaydi (ko'p SMSC RINVDSTADR qaytaradi) — International() konstruktori '+'ni o'zi olib tashlaydi")
		}
		if !allDigits(a.Addr) {
			return fmt.Errorf("pdu: international manzilda raqam bo'lmagan belgi: %q", a.Addr)
		}
	default:
		// Bo'sh manzil faqat to'liq NULL uchlikda normal (NULL source, §5.2.8).
		if a.Addr == "" && (a.TON != TONUnknown || a.NPI != NPIUnknown) {
			return fmt.Errorf("pdu: bo'sh manzil faqat TON=0/NPI=0 bilan yuboriladi (keldi %d/%d)", a.TON, a.NPI)
		}
	}
	return nil
}

func writeAddress(b *bytes.Buffer, a Address, field string) error {
	writeUint8(b, a.TON)
	writeUint8(b, a.NPI)
	return writeCString(b, a.Addr, maxAddrShort, field)
}

// maxAddrLong — data_sm va alert_notification'dagi manzil o'lchami (§4.7.1,
// §4.12.1): Var. max 65 — submit_sm'dagi 21 emas. Bir xil Address struct,
// faqat simdagi limit boshqa.
const maxAddrLong = 65

func writeAddressLong(b *bytes.Buffer, a Address, field string) error {
	writeUint8(b, a.TON)
	writeUint8(b, a.NPI)
	return writeCString(b, a.Addr, maxAddrLong, field)
}

func readAddressLong(r *bytes.Reader, field string) (Address, error) {
	var a Address
	var err error
	if a.TON, err = readUint8(r, field+"_ton"); err != nil {
		return a, err
	}
	if a.NPI, err = readUint8(r, field+"_npi"); err != nil {
		return a, err
	}
	if a.Addr, err = readCString(r, maxAddrLong, field); err != nil {
		return a, err
	}
	return a, nil
}

func readAddress(r *bytes.Reader, field string) (Address, error) {
	var a Address
	var err error
	if a.TON, err = readUint8(r, field+"_ton"); err != nil {
		return a, err
	}
	if a.NPI, err = readUint8(r, field+"_npi"); err != nil {
		return a, err
	}
	if a.Addr, err = readCString(r, maxAddrShort, field); err != nil {
		return a, err
	}
	return a, nil
}
