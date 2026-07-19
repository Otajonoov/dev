package dlr

import (
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"

	"smpp/tlv"
)

// Receipt — bitta DLR'dan ajratib olingan ma'lumot. Matn (Appendix B uslubi)
// va TLV'lar birga hisobga olinadi; TLV bor bo'lsa — u haqiqat manbai.
type Receipt struct {
	ID    string       // original xabar message_id'si (TLV ustun, keyin id:)
	State MessageState // yakuniy holat (TLV ustun, keyin stat:); 0 = topilmadi

	// Matn field'lari — topilmagani uchun Sub/Dlvrd = -1 (0 ma'noli qiymat:
	// "dlvrd:000" = hech biri yetmadi), sanalar zero time.Time.
	Sub        int       // sub: — nechta qism submit qilingan
	Dlvrd      int       // dlvrd: — nechta qism yetkazilgan
	SubmitDate time.Time // submit date: — YYMMDDhhmm yoki YYMMDDhhmmss
	DoneDate   time.Time // done date:
	Stat       string    // stat: — xom matn (State unga tolerant mapping)
	Err        string    // err: — operator-specific kod; universal jadvali YO'Q
	Text       string    // text: — original matnning boshi (faqat ma'lumot)

	NetErr *tlv.NetworkError // network_error_code TLV'si kelgan bo'lsa
}

// ErrNotReceipt — kirishda DLR'ga xos hech narsa topilmadi: birorta
// key:value ham, DLR TLV'lari ham yo'q. Bunday deliver_sm katta ehtimol
// MO xabar — dispatch'da esm_class tekshirilmagan bo'lishi mumkin.
var ErrNotReceipt = errors.New("dlr: matnda ham, TLV'larda ham receipt belgisi yo'q")

// Matnda qidiriladigan kalitlar. Bir kalitning bir nechta imlosi bor —
// format vendor-specific, "submit_date" yoki "sub date" yozadiganlar uchraydi.
// Uzunroq variant OLDIN turishi kerak: "submit date" topilgan joyda "date"
// qidirilmaydi.
var receiptKeys = []struct {
	name  string // kanonik nom
	forms []string
}{
	{"id", []string{"id:"}},
	{"sub", []string{"sub:"}},
	{"dlvrd", []string{"dlvrd:"}},
	{"submitdate", []string{"submit date:", "submit_date:", "sub date:"}},
	{"donedate", []string{"done date:", "done_date:", "donedate:"}},
	{"stat", []string{"stat:", "status:"}},
	{"err", []string{"err:"}},
	{"text", []string{"text:"}},
}

// Parse deliver_sm'ning short_message'i va TLV tail'idan Receipt yig'adi.
//
// Tolerantlik qoidalari (Appendix B "vendor specific... typical example"
// bo'lgani uchun):
//   - kalitlar istalgan TARTIBda va istalgan REGISTRda bo'lishi mumkin;
//   - yetishmagan field xato emas — Receipt qisman to'ladi;
//   - sana 10 (YYMMDDhhmm) yoki 12 (YYMMDDhhmmss) raqamli;
//   - buzuq qiymat (raqam bo'lmagan sub: kabi) shunchaki tashlanadi.
//
// Xato faqat kirishda umuman receipt belgisi bo'lmaganda (ErrNotReceipt).
func Parse(shortMessage []byte, tlvs []tlv.TLV) (Receipt, error) {
	r := Receipt{Sub: -1, Dlvrd: -1}
	found := r.parseText(string(shortMessage))

	// TLV'lar — ustuvor manba (§5.3.2.12, §5.3.2.35): matn bilan zid kelsa
	// TLV yutadi, chunki u strukturali va encoding'ga bog'liq emas.
	if t, ok := tlv.Find(tlvs, tlv.ReceiptedMessageID); ok {
		if id, err := t.CStringValue(); err == nil && id != "" {
			r.ID = id
			found = true
		}
	}
	if t, ok := tlv.Find(tlvs, tlv.MessageState); ok {
		if v, err := t.Uint8Value(); err == nil {
			r.State = MessageState(v)
			found = true
		}
	}
	if t, ok := tlv.Find(tlvs, tlv.NetworkErrorCode); ok {
		if ne, err := t.NetworkError(); err == nil {
			r.NetErr = &ne
		}
	}

	// TLV holat bermagan bo'lsa — stat: matnidan tolerant mapping.
	if r.State == 0 && r.Stat != "" {
		r.State = StateFromStat(r.Stat)
	}
	if !found {
		return r, ErrNotReceipt
	}
	return r, nil
}

// lowerASCII faqat ASCII harflarni kichraytiradi — UZUNLIK SAQLANADI.
// strings.ToLower ishlatib BO'LMAYDI: buzuq UTF-8 baytlarni (operator
// yuborishi mumkin!) U+FFFD (3 bayt) bilan almashtirib pozitsiyalarni
// suradi — fuzzer topgan real panic (testdata/fuzz corpus'ida regression).
func lowerASCII(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'A' && c <= 'Z' {
			b[i] = c + ('a' - 'A')
		}
	}
	return string(b)
}

// parseText matnni skan qiladi; kamida bitta kalit topilsa true.
func (r *Receipt) parseText(s string) bool {
	lower := lowerASCII(s)

	// Har kanonik kalitning matndagi BIRINCHI uchrashuvi: [boshlanish,
	// qiymat boshlanishi]. Keyin pozitsiya bo'yicha saralab, qiymatni
	// "keyingi kalitgacha" deb olamiz — tartibga bog'lanmaymiz.
	type hit struct {
		name       string
		start, val int
	}
	var hits []hit
	for _, k := range receiptKeys {
		for _, form := range k.forms {
			if i := indexKey(lower, form); i >= 0 {
				hits = append(hits, hit{k.name, i, i + len(form)})
				break
			}
		}
	}
	if len(hits) == 0 {
		return false
	}
	sort.Slice(hits, func(i, j int) bool { return hits[i].start < hits[j].start })

	for i, h := range hits {
		end := len(s)
		if i+1 < len(hits) {
			end = hits[i+1].start
		}
		v := strings.TrimSpace(s[h.val:end])
		switch h.name {
		case "id":
			r.ID = v
		case "sub":
			if n, err := strconv.Atoi(v); err == nil {
				r.Sub = n
			}
		case "dlvrd":
			if n, err := strconv.Atoi(v); err == nil {
				r.Dlvrd = n
			}
		case "submitdate":
			r.SubmitDate = parseReceiptDate(v)
		case "donedate":
			r.DoneDate = parseReceiptDate(v)
		case "stat":
			r.Stat = v
		case "err":
			r.Err = v
		case "text":
			r.Text = v
		}
	}
	return true
}

// indexKey kalitni faqat token BOSHIda qidiradi: yo satr boshi, yo oldingi
// belgi harf-raqam emas. Aks holda "sub:" kaliti "resubmit:" ichidan,
// "id:" esa "receipted_id:" ichidan topilib ketadi.
func indexKey(lower, key string) int {
	from := 0
	for {
		i := strings.Index(lower[from:], key)
		if i < 0 {
			return -1
		}
		i += from
		if i == 0 || !isAlnum(lower[i-1]) {
			return i
		}
		from = i + 1
	}
}

func isAlnum(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= '0' && c <= '9' || c == '_'
}

// parseReceiptDate YYMMDDhhmm (10) yoki YYMMDDhhmmss (12) formatni o'qiydi.
// DIQQAT: timezone spec'da YO'Q — operator o'z lokal vaqtini yozadi;
// biz UTC deb qaytaramiz, real integratsiyada offset operator bilan
// kelishiladi. Buzuq sana — zero time.Time, xato emas.
func parseReceiptDate(v string) time.Time {
	var layout string
	switch len(v) {
	case 10:
		layout = "0601021504"
	case 12:
		layout = "060102150405"
	default:
		return time.Time{}
	}
	t, err := time.ParseInLocation(layout, v, time.UTC)
	if err != nil {
		return time.Time{}
	}
	return t
}
