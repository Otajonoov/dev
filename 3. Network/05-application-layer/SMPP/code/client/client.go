package client

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"smpp/coding"
	"smpp/pdu"
	"smpp/session"
)

// StatusError — SMSC resp'ida command_status != 0. Xatoning O'ZI tasnifga
// tayyor: errors.As bilan olib Classify(e.Status) qilinadi.
type StatusError struct {
	Status pdu.CommandStatus
}

func (e StatusError) Error() string {
	return fmt.Sprintf("smpp: %s (%s)", e.Status, Classify(e.Status))
}

// ErrNotBound — hozir bound sessiya yo'q (reconnect ketmoqda). FAIL-FAST
// dizayni: Submit navbatda KUTMAYDI — xato darhol qaytadi, retry qarori
// chaqiruvchining queue'sida (aks holda client ichida ko'rinmas buffer
// o'sib boradi va uzilish paytidagi xabarlar taqdiri noaniq bo'lib qoladi).
var ErrNotBound = errors.New("client: bound sessiya yo'q (reconnect ketmoqda)")

// Config — Client sozlamalari.
type Config struct {
	Addr string // SMSC manzili ("host:port")

	SystemID   string
	Password   string
	SystemType string
	Mode       pdu.CommandID // bind turi; default CmdBindTransceiver

	Session session.Config // window/timer'lar (12-bob) — OnInbound'ni Client o'zi egallaydi

	// Reconnect backoff: Base*2^n, Max plato, har kutishga ±50% jitter
	// (thundering herd davosi — 12-bob). Default 1s/60s.
	ReconnectBase time.Duration
	ReconnectMax  time.Duration

	// RateLimiter — submit'dan OLDIN chaqiriladi (operator TPS'i).
	// nil = cheklovsiz. x/time/rate.Limiter ham to'g'ri keladi.
	RateLimiter RateLimiter

	// TLS — berilsa ulanish SMPP over TLS bo'ladi (16-bob).
	// DefaultTLSConfig bilan boshlang; InsecureSkipVerify — anti-pattern.
	TLS *tls.Config

	// Metrics — monitoring hook'lari (16-bob); nil = o'chiq.
	Metrics Metrics

	// OnDeliver — kelgan deliver_sm'lar (MO ham, DLR ham — ajratish
	// EsmClass'da, 9-bob). Session dispatcher goroutine'idan ketma-ket
	// chaqiriladi; deliver_sm_resp ALLAQACHON yuborilgan.
	OnDeliver func(pdu.DeliverSM, pdu.Header)

	Logf func(format string, args ...any)
}

func (c Config) withDefaults() Config {
	if c.Mode == 0 {
		c.Mode = pdu.CmdBindTransceiver
	}
	if c.ReconnectBase <= 0 {
		c.ReconnectBase = time.Second
	}
	if c.ReconnectMax <= 0 {
		c.ReconnectMax = 60 * time.Second
	}
	return c
}

// Client — session engine (12-bob) ustidagi foydalanuvchi qatlami:
// auto-reconnect, rate limit, SubmitLong va DLR dispatch. Uch kutubxona
// tajribasidan olingan dizayn: context-first API (ajankovic) + sync Submit
// (fiorix) + to'liq lifecycle (gosmpp) — ularning tuzoqlari testlarda
// regression sifatida qotirilgan.
type Client struct {
	cfg  Config
	refs *coding.RefCounter

	mu   sync.Mutex
	sess *session.Session // nil = hozir ulanish yo'q

	closed  chan struct{}
	closeMu sync.Once
	wg      sync.WaitGroup
}

// Dial birinchi ulanish+bind'ni SINXRON qiladi — noto'g'ri credential
// darhol xato bilan qaytadi (aks holda konfiguratsiya xatosi "background'da
// abadiy reconnect" bo'lib yashirinadi). Muvaffaqiyatdan keyin uzilishlarni
// background reconnect goroutine boshqaradi.
func Dial(ctx context.Context, cfg Config) (*Client, error) {
	cfg = cfg.withDefaults()
	c := &Client{cfg: cfg, refs: coding.NewRefCounter(), closed: make(chan struct{})}
	sess, err := c.connect(ctx)
	if err != nil {
		return nil, err
	}
	c.setSession(sess)
	c.metrics().SessionState(true)
	c.wg.Add(1)
	go c.reconnectLoop(sess)
	return c, nil
}

func (c *Client) logf(format string, args ...any) {
	if c.cfg.Logf != nil {
		c.cfg.Logf(format, args...)
	}
}

// connect — TCP/TLS dial + session + bind.
func (c *Client) connect(ctx context.Context) (*session.Session, error) {
	conn, err := c.dial(ctx)
	if err != nil {
		return nil, err
	}
	scfg := c.cfg.Session
	scfg.OnInbound = c.onInbound
	if scfg.Logf == nil {
		scfg.Logf = c.cfg.Logf
	}
	sess := session.New(conn, scfg)
	_, err = sess.Bind(ctx, pdu.Bind{
		Mode:             c.cfg.Mode,
		SystemID:         c.cfg.SystemID,
		Password:         c.cfg.Password,
		SystemType:       c.cfg.SystemType,
		InterfaceVersion: pdu.InterfaceVersion34,
	})
	if err != nil {
		sess.Close(ctx)
		return nil, err
	}
	return sess, nil
}

func (c *Client) setSession(s *session.Session) {
	c.mu.Lock()
	c.sess = s
	c.mu.Unlock()
}

func (c *Client) current() *session.Session {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.sess
}

// reconnectLoop sessiya o'limini kutadi va backoff+jitter bilan qayta
// ulanadi. Eski window'dagi xabarlar YANGI sessiyaga ko'chmaydi (12-bob:
// duplicate qarori chaqiruvchida).
func (c *Client) reconnectLoop(sess *session.Session) {
	defer c.wg.Done()
	for {
		select {
		case <-sess.Done():
			c.setSession(nil)
			c.metrics().SessionState(false)
			c.logf("client: sessiya uzildi: %v — reconnect boshlandi", sess.Err())
		case <-c.closed:
			return
		}

		attempt := 0
		for {
			select {
			case <-c.closed:
				return
			default:
			}
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			newSess, err := c.connect(ctx)
			cancel()
			c.metrics().ReconnectAttempt(err == nil)
			if err == nil {
				c.setSession(newSess)
				c.metrics().SessionState(true)
				c.logf("client: qayta bind muvaffaqiyatli (%d urinishdan keyin)", attempt+1)
				sess = newSess
				break // tashqi for: yangi sessiya o'limini kutishga qaytamiz
			}
			attempt++
			delay := c.backoff(attempt)
			c.logf("client: reconnect #%d muvaffaqiyatsiz: %v — %v kutamiz", attempt, err, delay)
			// Bind xatosiga HAM backoff (11-bob: RINVPASWD tight-loop = IP ban).
			select {
			case <-time.After(delay):
			case <-c.closed:
				return
			}
		}
	}
}

// backoff — exponential + to'liq bo'lmagan jitter: delay/2 + rand(delay/2).
// Jitter reconnect'da MAJBURIY (thundering herd — 12-bob); 11-bob
// RetryPolicy'sida yo'qligi bilan solishtiring.
func (c *Client) backoff(attempt int) time.Duration {
	d := c.cfg.ReconnectBase
	for i := 1; i < attempt && d < c.cfg.ReconnectMax; i++ {
		d *= 2
	}
	if d > c.cfg.ReconnectMax {
		d = c.cfg.ReconnectMax
	}
	half := d / 2
	return half + time.Duration(rand.Int63n(int64(half)+1))
}

// onInbound — session dispatcher'idan keladi: deliver_sm'larni handler'ga
// uzatadi, qolganini log qiladi.
func (c *Client) onInbound(r session.Resp) {
	switch v := r.PDU.(type) {
	case pdu.DeliverSM:
		c.metrics().DeliverReceived(v.EsmClass.IsDeliveryReceipt())
		if c.cfg.OnDeliver != nil {
			c.cfg.OnDeliver(v, r.Header)
		}
	default:
		c.logf("client: kutilmagan inbound %s — e'tiborsiz", r.Header.ID)
	}
}

// Submit bitta submit_sm yuborib message_id qaytaradi. Xato turlari:
// ErrNotBound (reconnect ketmoqda — fail fast), StatusError (SMSC rad
// etdi — Classify bilan tasniflanadi), session.ErrResponseTimeout
// (taqdiri NOMA'LUM — 11-bob "uchinchi rejim": ko'r-ko'rona retry duplicate
// xavfi!), ctx xatolari.
func (c *Client) Submit(ctx context.Context, sm pdu.SubmitSM) (string, error) {
	if c.cfg.RateLimiter != nil {
		if err := c.cfg.RateLimiter.Wait(ctx); err != nil {
			return "", err
		}
	}
	sess := c.current()
	if sess == nil {
		return "", ErrNotBound
	}
	start := time.Now()
	r, err := sess.Send(ctx, sm)
	if err != nil {
		return "", err
	}
	resp, ok := r.PDU.(pdu.SubmitSMResp)
	if !ok {
		return "", fmt.Errorf("client: submit_sm'ga %s keldi", r.Header.ID)
	}
	c.metrics().SubmitObserved(pdu.CommandStatus(resp.Status), time.Since(start))
	if resp.Status != 0 {
		return "", StatusError{Status: pdu.CommandStatus(resp.Status)}
	}
	return resp.MessageID, nil
}

// LongMessage — SubmitLong parametrlari.
type LongMessage struct {
	Source             pdu.Address
	Dest               pdu.Address
	Text               string // xom matn — encoding avtomatik (7-bob Choose)
	RegisteredDelivery pdu.RegisteredDelivery
}

// SubmitLong matnni kerak bo'lsa segmentlarga bo'lib yuboradi (8-bob:
// UDH 8-bit default) va BARCHA segment message_id'larini qaytaradi —
// "xabar yetkazildi" = hamma id'ning DLR'i kelganda (9-bob korrelyatsiya).
//
// Har segment ALOHIDA sequence oladi (gosmpp #178 regression'i testda) va
// alohida submit sifatida rate limit'dan o'tadi. Qisman muvaffaqiyat
// mumkin: xato qaytganda o'sha paytgacha olingan id'lar ham qaytariladi —
// chaqiruvchi qolganini qayta yuborish/bekor qilishni o'zi hal qiladi.
func (c *Client) SubmitLong(ctx context.Context, m LongMessage) ([]string, error) {
	dcRaw, _ := coding.Choose(m.Text)
	segs, err := coding.Split(coding.Normalize(m.Text), dcRaw, coding.MethodUDH8, c.refs.Next(m.Dest.Addr))
	if err != nil {
		return nil, err
	}
	var ids []string
	for i, seg := range segs {
		sm := pdu.SubmitSM{SMFields: pdu.SMFields{
			Source:             m.Source,
			Dest:               m.Dest,
			RegisteredDelivery: m.RegisteredDelivery,
			DataCoding:         uint8(dcRaw),
			ShortMessage:       seg.Data,
		}}
		if seg.UDH != nil {
			sm.EsmClass = sm.EsmClass.WithUDHI() // UDH bor — UDHI SHART (§5.2.12)
		}
		id, err := c.Submit(ctx, sm)
		if err != nil {
			return ids, fmt.Errorf("client: segment %d/%d: %w", i+1, len(segs), err)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// Close reconnect'ni to'xtatadi va joriy sessiyani graceful yopadi.
func (c *Client) Close(ctx context.Context) error {
	var err error
	c.closeMu.Do(func() {
		close(c.closed)
		if sess := c.current(); sess != nil {
			err = sess.Close(ctx)
		}
		c.wg.Wait()
	})
	return err
}
