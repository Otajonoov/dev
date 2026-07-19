package dlr

import (
	"math/big"
	"strings"
	"sync"
)

// Bu fayl DLR korrelyatsiyasining eng mashhur og'rig'ini davolaydi:
// submit_sm_resp'dagi message_id bilan DLR id:'si BIR XIL SONNING IKKI
// YOZUVI bo'lishi mumkin — biri hex ("04000000086ECD50"), ikkinchisi
// decimal ("288230376293190992"). To'g'ridan-to'g'ri string solishtirish
// DLR'ni "yo'qotadi" (Kannel bug #334; Vonage'da buning uchun alohida
// account setting bor). Yechim — Kannel patch'idagi yondashuv: id'ning
// barcha mumkin talqinlarini chiqarib, lookup'ni shu variantlar bo'yicha
// qilish.

// NormalizeID id'ning lookup uchun kanonik variantlarini qaytaradi:
//
//   - id'ning o'zi (katta harfga keltirilgan — hex registri farq qilmasin);
//   - faqat raqamlardan iborat bo'lsa: decimal deb o'qib hex yozuvi
//     (va juft uzunlikka yetkazish uchun "0" prefiksli varianti — SMSC'lar
//     hex id'ni bayt chegarasiga to'ldirib yozadi);
//   - valid hex bo'lsa: hex deb o'qib decimal yozuvi.
//
// Faqat-raqamli id ikkala talqinga ham tushadi — "12345" ham decimal, ham
// hex bo'lishi mumkin, ikkala variant ham qaytariladi. message_id opaque
// (§5.2.23): raqamga o'xshamagan id (UUID, harfli prefiks) faqat o'z
// ko'rinishida qaytadi.
func NormalizeID(id string) []string {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil
	}
	up := strings.ToUpper(id)
	out := []string{up}
	add := func(s string) {
		for _, have := range out {
			if have == s {
				return
			}
		}
		out = append(out, s)
	}

	// Decimal talqin → hex variantlar.
	if isDigits(up) {
		if n, ok := new(big.Int).SetString(up, 10); ok {
			h := strings.ToUpper(n.Text(16))
			add(h)
			if len(h)%2 == 1 {
				add("0" + h) // "4000000086ECD50" → "04000000086ECD50"
			}
		}
	}
	// Hex talqin → decimal variant (faqat-raqamli id ham valid hex).
	if isHexDigits(up) {
		if n, ok := new(big.Int).SetString(up, 16); ok {
			add(n.Text(10))
		}
	}
	// Leading-zero'li hex/decimal uchun kesilgan varianti ham foydali:
	// "007F3A9B" saqlangan-u DLR'da "7F3A9B" kelishi mumkin.
	if trimmed := strings.TrimLeft(up, "0"); trimmed != "" && trimmed != up {
		add(trimmed)
	}
	return out
}

func isDigits(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

func isHexDigits(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if (c < '0' || c > '9') && (c < 'A' || c > 'F') {
			return false
		}
	}
	return true
}

// Table — thread-safe korrelyatsiya jadvali: har register qilingan
// message_id'ning BARCHA NormalizeID variantlari kanonik id'ga ko'rsatadi.
// Resolve DLR'dan kelgan id'ni (qaysi yozuvda bo'lsa ham) topadi.
//
// Bu ataylab minimal in-memory struktura — real gateway'da bu jadval
// Redis/PG'da yashaydi, lekin kalitlash printsipi aynan shu.
type Table struct {
	mu sync.Mutex
	m  map[string]string // variant → kanonik (register paytidagi) id
}

func NewTable() *Table {
	return &Table{m: make(map[string]string)}
}

// Register submit_sm_resp'dan kelgan message_id'ni jadvalga kiritadi.
func (t *Table) Register(messageID string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	for _, v := range NormalizeID(messageID) {
		t.m[v] = messageID
	}
}

// Resolve DLR id'sini kanonik message_id'ga qaytaradi. Avval to'g'ridan-
// to'g'ri variantlar qidiriladi; topilmasa dlrID'ning O'Z variantlari
// bo'yicha ham urinamiz (ikkala tomon ham normalize qilinadi).
func (t *Table) Resolve(dlrID string) (string, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	for _, v := range NormalizeID(dlrID) {
		if canon, ok := t.m[v]; ok {
			return canon, true
		}
	}
	return "", false
}

// Forget yakuniy DLR qabul qilingandan keyin jadvalni tozalaydi —
// aks holda uzoq yashaydigan protsessda jadval cheksiz o'sadi.
func (t *Table) Forget(messageID string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	for _, v := range NormalizeID(messageID) {
		if t.m[v] == messageID {
			delete(t.m, v)
		}
	}
}
