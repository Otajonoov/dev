// External test package: smsc'ni import qiladi (smsc → session import
// qilgani uchun in-package testda cycle bo'lardi — tashqi package'da mumkin).
package session_test

import (
	"context"
	"net"
	"testing"
	"time"

	"smpp/pdu"
	"smpp/session"
	"smpp/smsc"
)

// Real listener bilan integratsiya: smsc.TestServer'ga qarshi to'liq
// lifecycle (bind → avto enquire_link'lar → graceful close). net.Pipe emas —
// realistic buffering bilan.
func TestSessionAgainstTestServer(t *testing.T) {
	srv, err := smsc.StartTestServer()
	if err != nil {
		t.Fatal(err)
	}
	defer srv.Close()

	conn, err := net.Dial("tcp", srv.Addr())
	if err != nil {
		t.Fatal(err)
	}
	s := session.New(conn, session.Config{
		EnquireLink:     20 * time.Millisecond, // tez keepalive — server javob beradi
		ResponseTimeout: 500 * time.Millisecond,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	br, err := s.Bind(ctx, pdu.Bind{Mode: pdu.CmdBindTransceiver, SystemID: "esme1", InterfaceVersion: pdu.InterfaceVersion34})
	if err != nil {
		t.Fatalf("bind: %v", err)
	}
	if br.SystemID != "TESTSMSC" || !br.HasSCVersion {
		t.Fatalf("bind_resp: %+v", br)
	}
	// Bir necha enquire_link davri o'tsin — javoblar kelmasa sessiya
	// o'zini o'ldirardi (TestEnquireLinkDeath'ning teskarisi).
	select {
	case <-s.Done():
		t.Fatalf("sessiya nobud bo'ldi: %v", s.Err())
	case <-time.After(100 * time.Millisecond):
	}
	if err := s.Close(ctx); err != nil {
		t.Fatalf("graceful close: %v", err)
	}
}
