// localsmsc — mock SMSC'ni lokal ishga tushiruvchi demo (14-bob).
//
//	go run ./examples/localsmsc
//
// So'ng istalgan SMPP client (masalan bizning client package yoki tashqi
// tool) bilan ulanish mumkin: system_id=esme1, password=secret.
// Quirk rejimlarini flag'lar bilan yoqib real operator "g'alati"liklarini
// his qilib ko'ring.
package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"smpp/smsc"
)

func main() {
	hexdec := flag.Bool("hexdec", false, "resp'da hex, DLR'da decimal message_id")
	textonly := flag.Bool("textonly", false, "DLR'da TLV yo'q - faqat Appendix B matni")
	throttle := flag.Int("throttle", 0, "har N-chi submit RTHROTTLED (0=o'chiq)")
	slow := flag.Duration("slow", 0, "submit_sm_resp kechikishi (masalan 2s)")
	dlrDelay := flag.Duration("dlr", 500*time.Millisecond, "submit -> DLR kechikishi")
	flag.Parse()

	srv, err := smsc.Start(smsc.Config{
		Accounts: []smsc.Account{{SystemID: "esme1", Password: "secret"}},
		DLRDelay: *dlrDelay,
		Quirks: smsc.Quirks{
			HexIDDecimalDLR: *hexdec,
			DLRTextOnly:     *textonly,
			ThrottleEveryN:  *throttle,
			SlowResp:        *slow,
		},
		Logf: log.Printf,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("mock SMSC tinglayapti: %s (system_id=esme1, password=secret)", srv.Addr())
	log.Printf("quirk'lar: hexdec=%v textonly=%v throttle=%d slow=%v", *hexdec, *textonly, *throttle, *slow)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
	log.Println("yopilmoqda...")
	srv.Close()
}
