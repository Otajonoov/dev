package client

import (
	"context"
	"sync"
	"time"
)

// RateLimiter — submit oqimini operator TPS limitiga moslash nuqtasi.
// Interfeys ATAYLAB golang.org/x/time/rate.Limiter'ning Wait metodiga mos:
// x/time/rate ishlatmoqchi bo'lsangiz uni to'g'ridan-to'g'ri Config'ga
// qo'yasiz (drop-in). Core repo esa stdlib-only bo'lib qoladi — quyidagi
// PerSecond ko'p holat uchun yetarli.
//
// Window bilan adashtirmang (12-bob): window PARALLELIZMNI cheklaydi
// ("bir vaqtda nechta javobsiz"), rate TEZLIKNI ("soniyasiga nechta").
type RateLimiter interface {
	// Wait navbatdagi yuborishga ruxsat berilguncha bloklanadi
	// (yoki ctx bekor bo'lsa xato qaytaradi).
	Wait(ctx context.Context) error
}

// intervalLimiter — minimal token-bucket'ning yanada soddasi: yuborishlar
// orasida kamida `interval` bo'lishini ta'minlaydi (burst'siz, tekis oqim).
// Operator TPS shartnomalari uchun aynan shu xulq kerak: "sekundiga 100"
// deganda ko'p SMSC'lar 10ms'lik tekis oqimni kutadi, sekund boshida
// 100 talik portlashni emas.
type intervalLimiter struct {
	interval time.Duration

	mu   sync.Mutex
	next time.Time // navbatdagi yuborishga ruxsat vaqti
}

// PerSecond soniyasiga n ta yuborishga ruxsat beruvchi tekis limiter yasaydi.
func PerSecond(n int) RateLimiter {
	return &intervalLimiter{interval: time.Second / time.Duration(n)}
}

func (l *intervalLimiter) Wait(ctx context.Context) error {
	l.mu.Lock()
	now := time.Now()
	if l.next.Before(now) {
		l.next = now
	}
	wait := l.next.Sub(now)
	l.next = l.next.Add(l.interval)
	l.mu.Unlock()

	if wait <= 0 {
		return nil
	}
	t := time.NewTimer(wait)
	defer t.Stop()
	select {
	case <-t.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
