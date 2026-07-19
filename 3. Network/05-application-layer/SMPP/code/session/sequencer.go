package session

import "sync/atomic"

// maxSequence — sequence_number diapazonining yuqori chegarasi (§5.1.4:
// 0x00000001–0x7FFFFFFF). Undan keyin 1'ga qaytish spec'da YO'Q — bu
// industriya konventsiyasi (cloudhopper va boshqalar shunday qiladi);
// wrap'ni unutgan implementatsiya uzoq yashaydigan sessiyada sinadi.
const maxSequence = 0x7FFFFFFF

// Sequencer — thread-safe sequence_number generatori. Zero qiymati tayyor:
// birinchi Next() 1 qaytaradi. BITTA sessiyaning barcha PDU turlari shu
// bitta fazodan oladi (submit ham, enquire_link ham).
type Sequencer struct {
	n atomic.Uint32
}

// Next keyingi sequence_number'ni qaytaradi (1..0x7FFFFFFF, keyin yana 1).
// CAS loop — wrap nuqtasida ham race'siz: ikki goroutine bir xil qiymat
// olmaydi.
func (s *Sequencer) Next() uint32 {
	for {
		cur := s.n.Load()
		next := cur + 1
		if next > maxSequence {
			next = 1
		}
		if s.n.CompareAndSwap(cur, next) {
			return next
		}
	}
}
