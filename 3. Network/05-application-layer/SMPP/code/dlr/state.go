// Package dlr SMSC Delivery Receipt (DLR) bilan ishlaydi: deliver_sm ichida
// kelgan receipt'ni tolerant parse qilish (Appendix B matni + TLV'lar) va
// message_id'ni original submit bilan korrelyatsiya qilish (hex↔dec quirk
// bilan birga).
//
// Bu PDU codec EMAS — deliver_sm'ni pdu package decode qiladi; dlr esa uning
// short_message va TLV tail'i USTIDAGI biznes-mantiq.
package dlr

import "strings"

// MessageState — xabarning SMSC ichidagi holati (v3.4 §5.2.28, Table 5-6).
// Qiymatlar 1–8; zero qiymat "holat yo'q/aniqlanmadi" degani (spec'da 0 yo'q).
type MessageState uint8

const (
	StateEnroute       MessageState = 1 // yo'lda — yagona "tirik" no-final holat
	StateDelivered     MessageState = 2
	StateExpired       MessageState = 3 // validity period tugadi
	StateDeleted       MessageState = 4
	StateUndeliverable MessageState = 5
	StateAccepted      MessageState = 6 // "operator nomidan qabul qilindi"
	StateUnknown       MessageState = 7
	StateRejected      MessageState = 8
)

// Final xabar taqdiri hal bo'lganini bildiradi: SMSC endi qayta urinmaydi.
// ENROUTE — hali yo'lda; ACCEPTED va UNKNOWN — spec bo'yicha g'alati oraliq
// holatlar (Appendix B'da qisqartmasi bor, lekin retry siyosati operator'ga
// bog'liq) — ularni final deb hisoblamaymiz.
func (s MessageState) Final() bool {
	switch s {
	case StateDelivered, StateExpired, StateDeleted, StateUndeliverable, StateRejected:
		return true
	}
	return false
}

func (s MessageState) String() string {
	switch s {
	case StateEnroute:
		return "ENROUTE"
	case StateDelivered:
		return "DELIVERED"
	case StateExpired:
		return "EXPIRED"
	case StateDeleted:
		return "DELETED"
	case StateUndeliverable:
		return "UNDELIVERABLE"
	case StateAccepted:
		return "ACCEPTED"
	case StateUnknown:
		return "UNKNOWN"
	case StateRejected:
		return "REJECTED"
	}
	return "NONE"
}

// Abbrev Appendix B Table B-2'dagi 7 belgilik stat qisqartmasini qaytaradi.
// ENROUTE B-2'da yo'q (final emas), lekin ayrim SMSC'lar intermediate
// notification'da baribir yuboradi — shuning uchun u ham qaytariladi.
func (s MessageState) Abbrev() string {
	switch s {
	case StateEnroute:
		return "ENROUTE"
	case StateDelivered:
		return "DELIVRD"
	case StateExpired:
		return "EXPIRED"
	case StateDeleted:
		return "DELETED"
	case StateUndeliverable:
		return "UNDELIV"
	case StateAccepted:
		return "ACCEPTD"
	case StateUnknown:
		return "UNKNOWN"
	case StateRejected:
		return "REJECTD"
	}
	return ""
}

// StateFromStat DLR matnidagi stat: qiymatini MessageState'ga aylantiradi.
// Tolerant: katta-kichik farqsiz va PREFIKS bo'yicha — chunki vendor'lar
// "DELIVRD" o'rniga "DELIVERED", "REJECTD" o'rniga "REJECTED" yozib yuboradi
// (format Appendix B bo'yicha normativ emas). Tanilmasa 0 qaytadi.
func StateFromStat(stat string) MessageState {
	u := strings.ToUpper(strings.TrimSpace(stat))
	switch {
	case strings.HasPrefix(u, "DELIV"):
		return StateDelivered
	case strings.HasPrefix(u, "EXPIR"):
		return StateExpired
	case strings.HasPrefix(u, "DELET"):
		return StateDeleted
	case strings.HasPrefix(u, "UNDELIV"):
		return StateUndeliverable
	case strings.HasPrefix(u, "ACCEPT"):
		return StateAccepted
	case strings.HasPrefix(u, "REJECT"):
		return StateRejected
	case strings.HasPrefix(u, "ENROUTE"):
		return StateEnroute
	case strings.HasPrefix(u, "UNKNOWN"):
		return StateUnknown
	}
	return 0
}
