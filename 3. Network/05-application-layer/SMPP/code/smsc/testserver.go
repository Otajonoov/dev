// Package smsc mock SMSC server'ni beradi. Bu bosqichda (4-bob) faqat
// minimal test skeleti bor: TCP accept, bind'ga resp, enquire_link'ka resp,
// unbind'ga resp, notanish PDU'ga generic_nack. To'liq server (auth, state
// enforcement, DLR generatsiya, quirk rejimlari) 14-bobda quriladi.
package smsc

import (
	"fmt"
	"net"
	"sync"

	"smpp/pdu"
)

// esmeRInvCmdID — ESME_RINVCMDID (0x03, "Invalid Command ID"). To'liq
// command_status jadvali 11-bobda pdu/status.go'ga tushadi; skelet uchun
// bitta kod yetarli.
const esmeRInvCmdID = 0x00000003

const maxPDUSize = 64 * 1024

// TestServer — testlar uchun minimal SMSC skeleti: har ulanishga alohida
// goroutine, sodda "frame o'qi → javob yoz" sikli.
type TestServer struct {
	// SystemID bind resp'larda qaytariladigan server identifikatori.
	SystemID string

	ln net.Listener
	wg sync.WaitGroup

	mu    sync.Mutex
	conns map[net.Conn]struct{}
}

// StartTestServer 127.0.0.1'dagi bo'sh portda skelet serverni ishga tushiradi.
func StartTestServer() (*TestServer, error) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}
	s := &TestServer{SystemID: "TESTSMSC", ln: ln, conns: make(map[net.Conn]struct{})}
	s.wg.Add(1)
	go s.acceptLoop()
	return s, nil
}

// DropConnections barcha aktiv ulanishlarni QO'POL uzadi (unbind'siz) —
// reconnect testlari uchun "tarmoq uzildi" simulyatsiyasi. Listener ochiq
// qoladi: client qayta ulana oladi.
func (s *TestServer) DropConnections() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for c := range s.conns {
		c.Close()
	}
}

func (s *TestServer) track(c net.Conn) {
	s.mu.Lock()
	s.conns[c] = struct{}{}
	s.mu.Unlock()
}

func (s *TestServer) untrack(c net.Conn) {
	s.mu.Lock()
	delete(s.conns, c)
	s.mu.Unlock()
}

// Addr server manzilini qaytaradi ("127.0.0.1:PORT").
func (s *TestServer) Addr() string { return s.ln.Addr().String() }

// Close listener'ni yopadi va barcha connection handler'lar tugashini kutadi.
func (s *TestServer) Close() {
	s.ln.Close()
	s.wg.Wait()
}

func (s *TestServer) acceptLoop() {
	defer s.wg.Done()
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			return // listener yopildi
		}
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.handle(conn)
		}()
	}
}

// handle — bitta ulanishning sinxron sikli. Skelet uchun shu yetarli;
// timer'lar, state enforcement va parallel yozish 14-bobda.
func (s *TestServer) handle(conn net.Conn) {
	s.track(conn)
	defer s.untrack(conn)
	defer conn.Close()
	for {
		frame, err := pdu.ReadFrame(conn, maxPDUSize)
		if err != nil {
			return // EOF yoki buzilgan stream — skelet farqlamaydi
		}
		h, err := pdu.DecodeHeader(frame)
		if err != nil {
			return
		}
		resp, closeAfter := s.respond(h)
		if resp != nil {
			if _, err := conn.Write(resp); err != nil {
				return
			}
		}
		if closeAfter {
			return
		}
	}
}

// respond kelgan PDU header'iga qarab javob frame tuzadi.
func (s *TestServer) respond(h pdu.Header) (resp []byte, closeAfter bool) {
	switch h.ID {
	case pdu.CmdBindTransmitter, pdu.CmdBindReceiver, pdu.CmdBindTransceiver:
		br := pdu.BindResp{
			Mode:         h.ID.Resp(),
			SystemID:     s.SystemID,
			SCVersion:    pdu.InterfaceVersion34,
			HasSCVersion: true,
		}
		out, err := br.Encode(h.Sequence)
		if err != nil {
			return nil, true
		}
		return out, false
	case pdu.CmdEnquireLink:
		return pdu.EncodeEnquireLinkResp(h.Sequence), false
	case pdu.CmdSubmitSM:
		// 13-bob client testlari uchun minimal submit: message_id seq'dan
		// yasaladi (har submit'ga UNIQUE id — SubmitLong segment testi shu
		// bilan #178-regression'ni ushlaydi). To'liq server (DLR, quirk'lar,
		// auth) — 14-bobda.
		out, err := pdu.SubmitSMResp{MessageID: fmt.Sprintf("TST%08X", h.Sequence)}.Encode(h.Sequence)
		if err != nil {
			return nil, true
		}
		return out, false
	case pdu.CmdUnbind:
		// unbind_resp'dan KEYIN yopiladi — §4.2 tartibi.
		return pdu.EncodeUnbindResp(0, h.Sequence), true
	default:
		// Skelet tanimagan har qanday PDU → generic_nack (§3.3, §4.3).
		return pdu.EncodeGenericNack(esmeRInvCmdID, h.Sequence), false
	}
}
