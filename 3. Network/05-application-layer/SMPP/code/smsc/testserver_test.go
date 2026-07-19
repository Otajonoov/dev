package smsc

import (
	"errors"
	"io"
	"net"
	"testing"

	"smpp/pdu"
)

const testMaxPDU = 4096

func dialServer(t *testing.T) (net.Conn, *TestServer) {
	t.Helper()
	srv, err := StartTestServer()
	if err != nil {
		t.Fatalf("server ishga tushmadi: %v", err)
	}
	t.Cleanup(srv.Close)
	conn, err := net.Dial("tcp", srv.Addr())
	if err != nil {
		t.Fatalf("dial xatosi: %v", err)
	}
	t.Cleanup(func() { conn.Close() })
	return conn, srv
}

func TestBindFlow(t *testing.T) {
	conn, _ := dialServer(t)

	bind := pdu.Bind{
		Mode:             pdu.CmdBindTransceiver,
		SystemID:         "uzsms",
		Password:         "s3cr3t",
		InterfaceVersion: pdu.InterfaceVersion34,
	}
	frame, err := bind.Encode(1)
	if err != nil {
		t.Fatalf("bind encode: %v", err)
	}
	if _, err := conn.Write(frame); err != nil {
		t.Fatalf("write: %v", err)
	}

	respFrame, err := pdu.ReadFrame(conn, testMaxPDU)
	if err != nil {
		t.Fatalf("resp o'qilmadi: %v", err)
	}
	resp, h, err := pdu.DecodeBindResp(respFrame)
	if err != nil {
		t.Fatalf("DecodeBindResp: %v", err)
	}
	if h.Sequence != 1 {
		t.Errorf("resp seq = %d, request'niki (1) aynan qaytishi kerak", h.Sequence)
	}
	if resp.Status != 0 || resp.SystemID != "TESTSMSC" {
		t.Errorf("resp: status=%d system_id=%q", resp.Status, resp.SystemID)
	}
	if !resp.HasSCVersion || resp.SCVersion != pdu.InterfaceVersion34 {
		t.Errorf("sc_interface_version TLV: has=%v val=0x%02X, kutilgan 0x34", resp.HasSCVersion, resp.SCVersion)
	}
}

func TestEnquireLinkResp(t *testing.T) {
	conn, _ := dialServer(t)

	if _, err := conn.Write(pdu.EncodeEnquireLink(28)); err != nil {
		t.Fatalf("write: %v", err)
	}
	frame, err := pdu.ReadFrame(conn, testMaxPDU)
	if err != nil {
		t.Fatalf("resp o'qilmadi: %v", err)
	}
	h, err := pdu.DecodeHeader(frame)
	if err != nil {
		t.Fatalf("DecodeHeader: %v", err)
	}
	if h.ID != pdu.CmdEnquireLinkResp || h.Sequence != 28 || h.Status != 0 {
		t.Errorf("kutilgan enquire_link_resp seq=28 status=0, keldi: %+v", h)
	}
}

func TestUnknownPDUGetsGenericNack(t *testing.T) {
	conn, _ := dialServer(t)

	// Notanish command_id'li header-only frame.
	raw := pdu.EncodeHeader(pdu.Header{Length: pdu.HeaderSize, ID: pdu.CommandID(0xAA), Sequence: 7})
	if _, err := conn.Write(raw[:]); err != nil {
		t.Fatalf("write: %v", err)
	}
	frame, err := pdu.ReadFrame(conn, testMaxPDU)
	if err != nil {
		t.Fatalf("resp o'qilmadi: %v", err)
	}
	h, err := pdu.DecodeHeader(frame)
	if err != nil {
		t.Fatalf("DecodeHeader: %v", err)
	}
	if h.ID != pdu.CmdGenericNack || h.Status != esmeRInvCmdID || h.Sequence != 7 {
		t.Errorf("kutilgan generic_nack status=0x03 seq=7, keldi: %+v", h)
	}
}

func TestUnbindClosesConnection(t *testing.T) {
	conn, _ := dialServer(t)

	if _, err := conn.Write(pdu.EncodeUnbind(5)); err != nil {
		t.Fatalf("write: %v", err)
	}
	frame, err := pdu.ReadFrame(conn, testMaxPDU)
	if err != nil {
		t.Fatalf("unbind_resp o'qilmadi: %v", err)
	}
	h, err := pdu.DecodeHeader(frame)
	if err != nil {
		t.Fatalf("DecodeHeader: %v", err)
	}
	if h.ID != pdu.CmdUnbindResp || h.Sequence != 5 {
		t.Errorf("kutilgan unbind_resp seq=5, keldi: %+v", h)
	}
	// unbind_resp'dan keyin server TCP'ni yopadi (§4.2 tartibi) — keyingi
	// o'qish toza EOF bilan tugashi kerak.
	if _, err := pdu.ReadFrame(conn, testMaxPDU); !errors.Is(err, io.EOF) {
		t.Errorf("ulanish yopilishi kutilgan edi (io.EOF), keldi: %v", err)
	}
}
