package coding

import "strings"

// DataCoding — data_coding baytining biz ishlatadigan qiymatlari (v3.4 §5.2.19).
// MUHIM: 0x00 "default" qiymat EMAS — "SMSC'ning default alphabet'i" degani;
// qaysi alphabet ekani spec'da YO'Q (note c), operator bilan kelishiladi.
type DataCoding uint8

const (
	DCDefault DataCoding = 0x00 // SMSC default alphabet — amalda ko'pincha GSM7
	DCLatin1  DataCoding = 0x03 // ISO-8859-1
	DCBinary  DataCoding = 0x04 // 8-bit binary (UDH'li binary payload'lar, 8-bob)
	DCUCS2    DataCoding = 0x08 // UCS2 / UTF-16BE
)

// normalizeReplacer — o'zbek lotin matnini SMS-safe qilish jadvali:
//
//	U+02BB (oʻ/gʻ dagi rasmiy belgi — MODIFIER LETTER TURNED COMMA)
//	U+02BC (tutuq belgisi — MODIFIER LETTER APOSTROPHE)
//	U+2018/U+2019 ("aqlli" bir qo'shtirnoqlar — klaviatura/CMS'lardan keladi)
//
// hammasi → ASCII ' (0x27, GSM7'da bor). Bitta shunday belgi butun xabarni
// UCS2'ga tushirib, 160 → 70 limit qiladi — normalizatsiya shu narxni to'laydi.
// Imloviy jihatdan ASCII apostrof — rasmiy U+02BB'ning keng qabul qilingan
// surrogati (ko'p o'zbek saytlari, jumladan davlat saytlari ham shuni ishlatadi).
var normalizeReplacer = strings.NewReplacer(
	"ʻ", "'",
	"ʼ", "'",
	"‘", "'",
	"’", "'",
)

// Normalize o'zbek lotin matnidagi GSM7'da yo'q apostrof-belgilarni ASCII '
// ga almashtiradi. Kirill matniga TA'SIR QILMAYDI (u baribir UCS2'ga ketadi).
func Normalize(s string) string {
	return normalizeReplacer.Replace(s)
}

// Choose matn uchun data_coding tanlaydi va matnni shu encoding'da baytlaydi.
// Algoritm: AVVAL Normalize (tartib muhim — aks holda oʻ'li matn behuda
// UCS2'ga ketadi), KEYIN skan: hamma belgi GSM7'da bo'lsa → DCDefault +
// unpacked GSM7 baytlar; aks holda → DCUCS2 + UTF-16BE baytlar.
//
// DCDefault tanlovi "SMSC default = GSM7" kelishuviga tayanadi — bu eng keng
// tarqalgan holat, lekin provisioning'da TASDIQLANISHI shart (7-bob, dc=0 tuzog'i).
func Choose(text string) (DataCoding, []byte) {
	text = Normalize(text)
	if b, err := EncodeGSM7(text); err == nil {
		return DCDefault, b
	}
	return DCUCS2, EncodeUCS2(text)
}
