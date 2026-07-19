package pdu

import (
	"bytes"
	"testing"
)

// goldenBindTRXHex — qo'lda yig'ilgan bind_transceiver (34 = 0x22 oktet):
// system_id "uzsms", password "s3cr3t", system_type bo'sh,
// interface_version 0x34, addr_ton/npi 0, address_range bo'sh, seq=1.
const goldenBindTRXHex = `
00 00 00 22 00 00 00 09 00 00 00 00 00 00 00 01
75 7A 73 6D 73 00
73 33 63 72 33 74 00
00
34 00 00
00`

// goldenBindTRXRespHex — muvaffaqiyatli bind_transceiver_resp (30 = 0x1E oktet):
// system_id "TESTSMSC" + sc_interface_version TLV (0x34), seq=1.
const goldenBindTRXRespHex = `
00 00 00 1E 80 00 00 09 00 00 00 00 00 00 00 01
54 45 53 54 53 4D 53 43 00
02 10 00 01 34`

var goldenBind = Bind{
	Mode:             CmdBindTransceiver,
	SystemID:         "uzsms",
	Password:         "s3cr3t",
	InterfaceVersion: InterfaceVersion34,
}

func TestBindEncodeGolden(t *testing.T) {
	frame, err := goldenBind.Encode(1)
	if err != nil {
		t.Fatalf("Encode xatosi: %v", err)
	}
	want := mustHex(t, goldenBindTRXHex)
	if !bytes.Equal(frame, want) {
		t.Errorf("Encode = % X,\nkutilgan % X", frame, want)
	}
}

func TestBindDecodeGolden(t *testing.T) {
	b, h, err := DecodeBind(mustHex(t, goldenBindTRXHex))
	if err != nil {
		t.Fatalf("DecodeBind xatosi: %v", err)
	}
	if h.Sequence != 1 || h.Status != 0 {
		t.Errorf("header: seq=%d status=%d, kutilgan 1/0", h.Sequence, h.Status)
	}
	if b != goldenBind {
		t.Errorf("DecodeBind = %+v,\nkutilgan %+v", b, goldenBind)
	}
}

func TestBindRoundTripAllModes(t *testing.T) {
	for _, mode := range []CommandID{CmdBindTransmitter, CmdBindReceiver, CmdBindTransceiver} {
		t.Run(mode.String(), func(t *testing.T) {
			in := Bind{
				Mode:             mode,
				SystemID:         "InternetGW",
				Password:         "pwd",
				SystemType:       "OTA",
				InterfaceVersion: InterfaceVersion34,
				AddrTON:          1,
				AddrNPI:          1,
				AddressRange:     "^9989",
			}
			frame, err := in.Encode(42)
			if err != nil {
				t.Fatalf("Encode xatosi: %v", err)
			}
			out, h, err := DecodeBind(frame)
			if err != nil {
				t.Fatalf("DecodeBind xatosi: %v", err)
			}
			if out != in || h.Sequence != 42 {
				t.Errorf("round-trip: %+v (seq=%d), kutilgan %+v (seq=42)", out, h.Sequence, in)
			}
		})
	}
}

func TestBindFieldLimits(t *testing.T) {
	tests := []struct {
		name string
		bind Bind
	}{
		{"system_id 16 belgi (max 15)", Bind{Mode: CmdBindTransceiver, SystemID: "0123456789ABCDEF"}},
		{"password 9 belgi (max 8)", Bind{Mode: CmdBindTransceiver, SystemID: "x", Password: "012345678"}},
		{"mode bind emas", Bind{Mode: CmdSubmitSM, SystemID: "x"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.bind.Encode(1); err == nil {
				t.Error("xato kutilgan edi, nil keldi")
			}
		})
	}
}

func TestDecodeBindRejectsOtherPDU(t *testing.T) {
	if _, _, err := DecodeBind(mustHex(t, enquireLinkHex)); err == nil {
		t.Error("enquire_link frame'ida DecodeBind xato qaytarishi kerak edi")
	}
}

func TestBindRespGolden(t *testing.T) {
	in := BindResp{
		Mode:         CmdBindTransceiverResp,
		SystemID:     "TESTSMSC",
		SCVersion:    InterfaceVersion34,
		HasSCVersion: true,
	}
	frame, err := in.Encode(1)
	if err != nil {
		t.Fatalf("Encode xatosi: %v", err)
	}
	want := mustHex(t, goldenBindTRXRespHex)
	if !bytes.Equal(frame, want) {
		t.Fatalf("Encode = % X,\nkutilgan % X", frame, want)
	}

	out, h, err := DecodeBindResp(frame)
	if err != nil {
		t.Fatalf("DecodeBindResp xatosi: %v", err)
	}
	if out != in || h.Sequence != 1 {
		t.Errorf("round-trip: %+v, kutilgan %+v", out, in)
	}
}

func TestBindRespErrorHasNoBody(t *testing.T) {
	// §4.1.2 Note: status != 0 bo'lsa body qaytarilmaydi — frame 16 oktet.
	in := BindResp{Mode: CmdBindTransceiverResp, Status: 0x0E /* ESME_RINVPASWD */, SystemID: "IGNORED"}
	frame, err := in.Encode(7)
	if err != nil {
		t.Fatalf("Encode xatosi: %v", err)
	}
	if len(frame) != HeaderSize {
		t.Fatalf("xatoli resp %d oktet, faqat header (16) bo'lishi kerak", len(frame))
	}

	out, _, err := DecodeBindResp(frame)
	if err != nil {
		t.Fatalf("DecodeBindResp xatosi: %v", err)
	}
	if out.Status != 0x0E || out.SystemID != "" || out.HasSCVersion {
		t.Errorf("xatoli resp: %+v — SystemID bo'sh, TLV yo'q bo'lishi kerak", out)
	}
}

func TestBindRespWithoutTLV(t *testing.T) {
	// sc_interface_version YO'QLIGI = "SMSC TLV qo'llamaydi" (§3.4) —
	// HasSCVersion=false bilan ifodalanadi.
	in := BindResp{Mode: CmdBindReceiverResp, SystemID: "OLDSMSC"}
	frame, err := in.Encode(3)
	if err != nil {
		t.Fatalf("Encode xatosi: %v", err)
	}
	out, _, err := DecodeBindResp(frame)
	if err != nil {
		t.Fatalf("DecodeBindResp xatosi: %v", err)
	}
	if out.HasSCVersion {
		t.Error("TLV'siz resp'da HasSCVersion=false bo'lishi kerak")
	}
	if out.SystemID != "OLDSMSC" {
		t.Errorf("SystemID = %q", out.SystemID)
	}
}

func TestOutbindRoundTrip(t *testing.T) {
	in := Outbind{SystemID: "SMSC01", Password: "secret"}
	frame, err := in.Encode(1)
	if err != nil {
		t.Fatalf("Encode xatosi: %v", err)
	}
	out, h, err := DecodeOutbind(frame)
	if err != nil {
		t.Fatalf("DecodeOutbind xatosi: %v", err)
	}
	if out != in || h.ID != CmdOutbind {
		t.Errorf("round-trip: %+v, kutilgan %+v", out, in)
	}
}
