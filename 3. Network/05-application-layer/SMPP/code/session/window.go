package session

import (
	"context"
	"errors"
	"sync"
	"time"

	"smpp/pdu"
)

// Pending window — javob kutilayotgan (outstanding) request'lar jadvali.
// "Window" so'zi spec'da YO'Q: §2.9/§2.5.2 faqat asinxronlikka ruxsat beradi
// ("10 outstanding" — Note darajasidagi guideline); limit — implementatsiya
// va operator kelishuvi.
//
// Struktura tanlovi: sync.Mutex + map (sync.Map EMAS). Sabab: har entry
// bir marta yoziladi va bir marta o'chiriladi (write-heavy, doim yangi
// key'lar) — bu sync.Map'ning yomon case'i; ustiga len() va deadline scan
// iteratsiyasi kerak (victoriametrics tahlili).

// ErrResponseTimeout — request'ga ResponseTimeout ichida javob kelmadi.
// MUHIM (11-bob "uchinchi rejim"): bu "SMSC rad etdi" degani EMAS —
// "taqdiri noma'lum" degani; retry qarori duplicate xavfi bilan o'lchanadi.
var ErrResponseTimeout = errors.New("session: response timeout — javob kelmadi, taqdiri noma'lum")

// ErrSessionClosed — sessiya yopilgan; pending request'lar shu xato bilan
// yakunlanadi (rebind'da eski window yangi ulanishga KO'CHMAYDI).
var ErrSessionClosed = errors.New("session: sessiya yopilgan")

// Resp — window orqali kelgan javob: konkret PDU + header (status/seq).
type Resp struct {
	PDU    pdu.PDU
	Header pdu.Header
}

type pending struct {
	cmd      pdu.CommandID // request command_id — resp mosligini tekshirish uchun
	deadline time.Time
	ch       chan windowResult // sig'imi 1 — resolve hech qachon bloklanmaydi
}

type windowResult struct {
	resp Resp
	err  error
}

type window struct {
	slots chan struct{} // bo'sh joylar semafori — to'lganda Send bloklanadi

	mu sync.Mutex
	m  map[uint32]*pending
}

func newWindow(size int) *window {
	return &window{
		slots: make(chan struct{}, size),
		m:     make(map[uint32]*pending, size),
	}
}

// acquire bo'sh slot kutadi: window to'la bo'lsa ctx bekor bo'lguncha
// bloklanadi — backpressure Submit chaqiruvchisiga OQIB CHIQADI.
func (w *window) acquire(ctx context.Context) error {
	select {
	case w.slots <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// add seq uchun pending entry ochadi (slot allaqachon olingan bo'lishi kerak).
func (w *window) add(seq uint32, cmd pdu.CommandID, deadline time.Time) *pending {
	p := &pending{cmd: cmd, deadline: deadline, ch: make(chan windowResult, 1)}
	w.mu.Lock()
	w.m[seq] = p
	w.mu.Unlock()
	return p
}

// resolve kelgan javobni seq bo'yicha egasiga yetkazadi. Topilmasa false —
// notanish seq (kechikkan/duplicate javob) — chaqiruvchi log qiladi,
// generic_nack YUBORMAYDI (11-bob: resp'ga nack yo'q).
//
// command_id mosligi ham tekshiriladi (§2.5.2 korrelyatsiya faqat seq
// orqali, lekin buggy SMSC noto'g'ri resp turi qaytarishi mumkin):
// generic_nack har qanday request'ga javob sanaladi.
func (w *window) resolve(seq uint32, r Resp) bool {
	w.mu.Lock()
	p, ok := w.m[seq]
	if ok {
		if _, isNack := r.PDU.(pdu.GenericNack); !isNack && r.Header.ID != p.cmd.Resp() {
			w.mu.Unlock()
			return false // seq mos-u, turi boshqa — bu bizning javob emas
		}
		delete(w.m, seq)
	}
	w.mu.Unlock()
	if !ok {
		return false
	}
	p.ch <- windowResult{resp: r}
	<-w.slots // slot bo'shadi
	return true
}

// fail bitta pending'ni xato bilan yakunlaydi (masalan encode'dan keyin
// yozish muvaffaqiyatsiz bo'lsa).
func (w *window) fail(seq uint32, err error) {
	w.mu.Lock()
	p, ok := w.m[seq]
	if ok {
		delete(w.m, seq)
	}
	w.mu.Unlock()
	if ok {
		p.ch <- windowResult{err: err}
		<-w.slots
	}
}

// expire deadline'i o'tgan barcha entry'larni ErrResponseTimeout bilan
// yakunlaydi. Scanner goroutine davriy chaqiradi — busiz map LEAK bo'ladi:
// javobsiz qolgan har submit abadiy joy egallab window'ni "toraytiradi".
func (w *window) expire(now time.Time) int {
	w.mu.Lock()
	var expired []*pending
	for seq, p := range w.m {
		if now.After(p.deadline) {
			delete(w.m, seq)
			expired = append(expired, p)
		}
	}
	w.mu.Unlock()
	for _, p := range expired {
		p.ch <- windowResult{err: ErrResponseTimeout}
		<-w.slots
	}
	return len(expired)
}

// failAll sessiya yopilganda BARCHA pending'larni xato bilan yakunlaydi.
func (w *window) failAll(err error) {
	w.mu.Lock()
	all := make([]*pending, 0, len(w.m))
	for seq, p := range w.m {
		delete(w.m, seq)
		all = append(all, p)
	}
	w.mu.Unlock()
	for _, p := range all {
		p.ch <- windowResult{err: err}
		<-w.slots
	}
}

// depth — hozirgi outstanding soni (monitoring metrikasi, 16-bob).
func (w *window) depth() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return len(w.m)
}
