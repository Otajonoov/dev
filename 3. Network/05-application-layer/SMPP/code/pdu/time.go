package pdu

import (
	"fmt"
	"time"
)

// Vaqt formati (§7.1.1): 16 belgi "YYMMDDhhmmsstnnp" (+ NULL = 17 oktet).
//   t  — sekundning o'ndan biri (0–9)
//   nn — UTC'dan farq CHORAK SOATLARDA (00–48)
//   p  — '+' (UTC'dan oldinda) / '-' (orqada) / 'R' (relative)
// Relative'da (p='R') qiymat SMSC joriy vaqtidan boshlab davr; t='0', nn='00'.
// SMSC javoblarida (masalan query_sm_resp final_date) 12 belgili
// "YYMMDDhhmmss" varianti ham uchraydi — SMSC lokal vaqti (§7.1.1 Note).

// TimeValue — parse qilingan SMPP vaqt qiymati.
type TimeValue struct {
	Relative bool

	// Absolute holat uchun (Relative=false):
	At        time.Time
	HasOffset bool // false = 12-belgili variant (zona noma'lum, At UTC deb qurilgan)

	// Relative holat uchun (Relative=true) — davr komponentlari.
	// time.Duration EMAS: oy/yil kalendarga bog'liq, aniq davr SMSC vaqtida hal bo'ladi.
	Years, Months, Days, Hours, Minutes, Seconds int
}

// EncodeAbsolute t'ni absolute formatga o'tkazadi (§7.1.1.1); vaqt loc
// zonasida ifodalanadi, nn = shu zonaning UTC'dan farqi chorak soatlarda.
func EncodeAbsolute(t time.Time, loc *time.Location) (string, error) {
	t = t.In(loc)
	_, offSec := t.Zone()
	sign := byte('+')
	if offSec < 0 {
		sign = '-'
		offSec = -offSec
	}
	if offSec%900 != 0 {
		return "", fmt.Errorf("pdu: zona offseti %ds chorak soatga karrali emas", offSec)
	}
	quarters := offSec / 900
	if quarters > 48 {
		return "", fmt.Errorf("pdu: zona offseti %d chorak soat, max 48", quarters)
	}
	tenth := t.Nanosecond() / 100_000_000
	return fmt.Sprintf("%02d%02d%02d%02d%02d%02d%d%02d%c",
		t.Year()%100, int(t.Month()), t.Day(),
		t.Hour(), t.Minute(), t.Second(), tenth, quarters, sign), nil
}

// EncodeRelative davomiylikni relative formatga o'tkazadi (§7.1.1.2):
// d kun/soat/minut/sekundga bo'linadi (yil/oy komponentlari 0 — Duration'da
// kalendar oy tushunchasi yo'q). 100 kundan uzun davr uchun TimeValue'ni
// qo'lda to'ldirib EncodeRelativeValue ishlatiladi.
func EncodeRelative(d time.Duration) (string, error) {
	if d < 0 {
		return "", fmt.Errorf("pdu: manfiy relative davr: %v", d)
	}
	total := int(d / time.Second)
	days := total / 86400
	if days > 99 {
		return "", fmt.Errorf("pdu: %d kun — DD field'iga (max 99) sig'maydi; oy/yil komponentlari bilan EncodeRelativeValue ishlating", days)
	}
	return EncodeRelativeValue(TimeValue{
		Relative: true,
		Days:     days,
		Hours:    total % 86400 / 3600,
		Minutes:  total % 3600 / 60,
		Seconds:  total % 60,
	})
}

// EncodeRelativeValue relative TimeValue'ni matnga o'tkazadi.
func EncodeRelativeValue(v TimeValue) (string, error) {
	if !v.Relative {
		return "", fmt.Errorf("pdu: EncodeRelativeValue faqat relative qiymat uchun")
	}
	for _, c := range []struct {
		name string
		val  int
	}{{"yil", v.Years}, {"oy", v.Months}, {"kun", v.Days}, {"soat", v.Hours}, {"minut", v.Minutes}, {"sekund", v.Seconds}} {
		if c.val < 0 || c.val > 99 {
			return "", fmt.Errorf("pdu: relative %s=%d — 00–99 oralig'idan tashqarida", c.name, c.val)
		}
	}
	return fmt.Sprintf("%02d%02d%02d%02d%02d%02d000R",
		v.Years, v.Months, v.Days, v.Hours, v.Minutes, v.Seconds), nil
}

// ParseTime 16 belgili (absolute/relative) yoki 12 belgili (SMSC lokal)
// vaqt matnini o'qiydi. Bo'sh string uchun chaqirilmaydi — "1 or 17"
// qoidasidagi NULL holatini codec o'zi hal qiladi.
func ParseTime(s string) (TimeValue, error) {
	switch len(s) {
	case 16:
		// davomida
	case 12:
		var v TimeValue
		t, err := parseDigits(s, "YYMMDDhhmmss")
		if err != nil {
			return v, err
		}
		v.At = time.Date(slideY2K(t[0]), time.Month(t[1]), t[2], t[3], t[4], t[5], 0, time.UTC)
		return v, nil
	default:
		return TimeValue{}, fmt.Errorf("pdu: vaqt matni %d belgi — 16 yoki 12 kutiladi: %q", len(s), s)
	}

	p := s[15]
	t, err := parseDigits(s[:13], "YYMMDDhhmmsst")
	if err != nil {
		return TimeValue{}, err
	}
	nn, err := parseDigits(s[13:15], "nn")
	if err != nil {
		return TimeValue{}, err
	}

	switch p {
	case 'R':
		return TimeValue{
			Relative: true,
			Years:    t[0], Months: t[1], Days: t[2],
			Hours: t[3], Minutes: t[4], Seconds: t[5],
		}, nil
	case '+', '-':
		if nn[0] > 48 {
			return TimeValue{}, fmt.Errorf("pdu: nn=%d — chorak soatlar max 48", nn[0])
		}
		offSec := nn[0] * 900
		if p == '-' {
			offSec = -offSec
		}
		loc := time.FixedZone("", offSec)
		return TimeValue{
			At: time.Date(slideY2K(t[0]), time.Month(t[1]), t[2],
				t[3], t[4], t[5], t[6]*100_000_000, loc),
			HasOffset: true,
		}, nil
	default:
		return TimeValue{}, fmt.Errorf("pdu: p belgisi %q — '+', '-' yoki 'R' kutiladi", p)
	}
}

// slideY2K — Appendix C sliding window: 38–99 → 19xx, 00–37 → 20xx.
func slideY2K(yy int) int {
	if yy >= 38 {
		return 1900 + yy
	}
	return 2000 + yy
}

// parseDigits s'ni 2 belgilik raqam guruhlariga bo'lib o'qiydi (oxirgi toq
// belgi — bitta raqam, masalan 't'). layout faqat xato matni uchun.
func parseDigits(s, layout string) ([]int, error) {
	var out []int
	for i := 0; i < len(s); i += 2 {
		if i+1 < len(s) {
			d1, d2 := s[i]-'0', s[i+1]-'0'
			if d1 > 9 || d2 > 9 {
				return nil, fmt.Errorf("pdu: %q vaqt matnida raqam emas belgi (%s)", s, layout)
			}
			out = append(out, int(d1)*10+int(d2))
		} else {
			d := s[i] - '0'
			if d > 9 {
				return nil, fmt.Errorf("pdu: %q vaqt matnida raqam emas belgi (%s)", s, layout)
			}
			out = append(out, int(d))
		}
	}
	return out, nil
}
