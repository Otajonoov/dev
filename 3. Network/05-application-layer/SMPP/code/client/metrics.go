package client

import (
	"time"

	"smpp/pdu"
)

// Metrics — monitoring hook'lari (16-bob). Interfeys ATAYLAB kichik va
// backend-agnostik: Prometheus, OpenTelemetry, statsd — adapter yozish
// chaqiruvchining ishi (examples/prometheus'da stdlib-only namuna bor).
// nil Config.Metrics = nopMetrics (hech narsa qilmaydi, hech narsa
// sekinlashtirmaydi).
//
// Nima O'LCHANADI va nega (research konsensusi):
//   - submit natijasi status bo'yicha — success rate va xato taqsimoti;
//   - submit RTT — SMSC javob tezligi (response timer'ni sozlash asosi);
//   - deliver oqimi (DLR/MO ajratilgan) — DLR latency'ning xom materiali;
//   - sessiya holati va reconnect urinishlari — bind flapping indikatori.
//
// Window depth esa pull-usulda: Session.WindowDepth()'ni gauge sifatida
// davriy so'rab olish adapter zimmasida (push emas — chunki u holat, hodisa emas).
type Metrics interface {
	// SubmitObserved har submit_sm yakunida: status (0=ok) va request↔resp RTT.
	// Timeout'da chaqirilmaydi — u alohida hodisa emas, xato yo'lida ko'rinadi.
	SubmitObserved(status pdu.CommandStatus, rtt time.Duration)
	// DeliverReceived har kelgan deliver_sm'da (resp allaqachon ketgan).
	DeliverReceived(isDLR bool)
	// SessionState bound holat o'zgarganda: true=bound, false=uzildi.
	SessionState(bound bool)
	// ReconnectAttempt har qayta ulanish urinishida.
	ReconnectAttempt(success bool)
}

type nopMetrics struct{}

func (nopMetrics) SubmitObserved(pdu.CommandStatus, time.Duration) {}
func (nopMetrics) DeliverReceived(bool)                            {}
func (nopMetrics) SessionState(bool)                               {}
func (nopMetrics) ReconnectAttempt(bool)                           {}

func (c *Client) metrics() Metrics {
	if c.cfg.Metrics != nil {
		return c.cfg.Metrics
	}
	return nopMetrics{}
}
