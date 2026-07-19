// prometheus — client.Metrics interfeysi uchun STDLIB-ONLY adapter namunasi
// (16-bob): metrikalar Prometheus text exposition formatida /metrics'da
// beriladi. Real loyihada prometheus/client_golang olinadi — bu namuna
// (a) interfeys qanday ulanishini, (b) exposition format naqadar oddiy
// ekanini ko'rsatadi (dependency'siz).
//
//	go run ./examples/prometheus
//	curl localhost:9090/metrics
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sort"
	"sync"
	"time"

	"smpp/client"
	"smpp/pdu"
	"smpp/session"
	"smpp/smsc"
)

// promMetrics — client.Metrics implementatsiyasi: counter/gauge'lar mutex
// ostida, /metrics handler'i ularni text formatda chiqaradi.
type promMetrics struct {
	mu             sync.Mutex
	submitTotal    map[string]uint64 // label: status nomi
	submitRTTSum   time.Duration
	submitRTTCount uint64
	deliverTotal   map[bool]uint64 // label: dlr yoki mo
	bound          int
	reconnects     map[bool]uint64

	// window depth — pull-gauge: sample() davriy yangilaydi.
	windowDepth int
}

func newPromMetrics() *promMetrics {
	return &promMetrics{
		submitTotal:  make(map[string]uint64),
		deliverTotal: make(map[bool]uint64),
		reconnects:   make(map[bool]uint64),
	}
}

func (p *promMetrics) SubmitObserved(status pdu.CommandStatus, rtt time.Duration) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.submitTotal[status.String()]++
	p.submitRTTSum += rtt
	p.submitRTTCount++
}

func (p *promMetrics) DeliverReceived(isDLR bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.deliverTotal[isDLR]++
}

func (p *promMetrics) SessionState(bound bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if bound {
		p.bound = 1
	} else {
		p.bound = 0
	}
}

func (p *promMetrics) ReconnectAttempt(success bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.reconnects[success]++
}

// ServeHTTP — Prometheus text exposition format (schema juda oddiy:
// "# TYPE nom tur" + "nom{label=\"qiymat\"} son").
func (p *promMetrics) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	p.mu.Lock()
	defer p.mu.Unlock()
	fmt.Fprintln(w, "# TYPE smpp_submit_total counter")
	statuses := make([]string, 0, len(p.submitTotal))
	for st := range p.submitTotal {
		statuses = append(statuses, st)
	}
	sort.Strings(statuses)
	for _, st := range statuses {
		fmt.Fprintf(w, "smpp_submit_total{command_status=%q} %d\n", st, p.submitTotal[st])
	}
	fmt.Fprintln(w, "# TYPE smpp_submit_rtt_seconds_sum counter")
	fmt.Fprintf(w, "smpp_submit_rtt_seconds_sum %f\n", p.submitRTTSum.Seconds())
	fmt.Fprintf(w, "smpp_submit_rtt_seconds_count %d\n", p.submitRTTCount)
	fmt.Fprintln(w, "# TYPE smpp_deliver_total counter")
	fmt.Fprintf(w, "smpp_deliver_total{kind=\"dlr\"} %d\n", p.deliverTotal[true])
	fmt.Fprintf(w, "smpp_deliver_total{kind=\"mo\"} %d\n", p.deliverTotal[false])
	fmt.Fprintln(w, "# TYPE smpp_bind_status gauge")
	fmt.Fprintf(w, "smpp_bind_status %d\n", p.bound)
	fmt.Fprintln(w, "# TYPE smpp_reconnect_total counter")
	fmt.Fprintf(w, "smpp_reconnect_total{success=\"true\"} %d\n", p.reconnects[true])
	fmt.Fprintf(w, "smpp_reconnect_total{success=\"false\"} %d\n", p.reconnects[false])
	fmt.Fprintln(w, "# TYPE smpp_window_depth gauge")
	fmt.Fprintf(w, "smpp_window_depth %d\n", p.windowDepth)
}

func main() {
	// Demo: lokal mock SMSC + metrics'li client + har soniyada trafik.
	srv, err := smsc.Start(smsc.Config{DLRDelay: 300 * time.Millisecond})
	if err != nil {
		log.Fatal(err)
	}
	defer srv.Close()

	m := newPromMetrics()
	ctx := context.Background()
	c, err := client.Dial(ctx, client.Config{
		Addr:     srv.Addr(),
		SystemID: "metrics-demo",
		Metrics:  m,
		Session:  session.Config{EnquireLink: 10 * time.Second},
		OnDeliver: func(d pdu.DeliverSM, h pdu.Header) {
			// DLR'lar DeliverReceived orqali hisoblanadi — bu yerda biznes yo'q.
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close(ctx)

	go func() {
		for {
			time.Sleep(time.Second)
			_, err := c.Submit(ctx, pdu.SubmitSM{SMFields: pdu.SMFields{
				Source:             pdu.Address{TON: pdu.TONAlphanumeric, Addr: "Demo"},
				Dest:               pdu.Address{TON: pdu.TONInternational, NPI: pdu.NPIISDN, Addr: "998901234567"},
				RegisteredDelivery: pdu.DLRFinal,
				ShortMessage:       []byte("metrics demo"),
			}})
			if err != nil {
				log.Printf("submit: %v", err)
			}
		}
	}()

	http.Handle("/metrics", m)
	log.Println("metrics: http://localhost:9090/metrics")
	log.Fatal(http.ListenAndServe("localhost:9090", nil))
}
