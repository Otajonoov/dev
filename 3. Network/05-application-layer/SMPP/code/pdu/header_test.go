package pdu

import (
	"bytes"
	"encoding/hex"
	"strings"
	"testing"
)

// mustHex bo'shliq/yangi-qatorli hex matnni baytlarga aylantiradi.
func mustHex(t *testing.T, s string) []byte {
	t.Helper()
	clean := strings.NewReplacer(" ", "", "\n", "", "\t", "").Replace(s)
	b, err := hex.DecodeString(clean)
	if err != nil {
		t.Fatalf("hex decode xatosi: %v", err)
	}
	return b
}

// specBindTransmitterHex — v3.4 §3.2.2'dagi RASMIY bind_transmitter misoli,
// 47 (0x2F) oktet. Diqqat: spec'ning o'z misolida interface_version=0x00
// (§5.2.4 "v3.4 uchun 0x34" qoidasiga zid — spec'ning ichki g'alatiligi).
const specBindTransmitterHex = `
00 00 00 2F 00 00 00 02 00 00 00 00 00 00 00 01
53 4D 50 50 33 54 45 53 54 00
73 65 63 72 65 74 30 38 00
53 55 42 4D 49 54 31 00
00 01 01 00`

func TestSpecExampleHeader(t *testing.T) {
	frame := mustHex(t, specBindTransmitterHex)
	if len(frame) != 47 {
		t.Fatalf("spec misoli %d oktet bo'ldi, 47 bo'lishi kerak", len(frame))
	}

	h, err := DecodeHeader(frame)
	if err != nil {
		t.Fatalf("DecodeHeader xatosi: %v", err)
	}
	want := Header{Length: 0x2F, ID: CmdBindTransmitter, Status: 0, Sequence: 1}
	if h != want {
		t.Errorf("DecodeHeader = %+v, kutilgan %+v", h, want)
	}

	// Encode qaytib aynan spec'dagi 16 oktetni berishi kerak.
	enc := EncodeHeader(h)
	if !bytes.Equal(enc[:], frame[:HeaderSize]) {
		t.Errorf("EncodeHeader = % X, kutilgan % X", enc[:], frame[:HeaderSize])
	}
}

func TestSpecExampleBody(t *testing.T) {
	frame := mustHex(t, specBindTransmitterHex)
	r := bytes.NewReader(frame[HeaderSize:])

	// bind_transmitter body tartibi (v3.4 §4.1.1, Table 4-1).
	systemID, err := readCString(r, 16, "system_id")
	if err != nil {
		t.Fatal(err)
	}
	password, err := readCString(r, 9, "password")
	if err != nil {
		t.Fatal(err)
	}
	systemType, err := readCString(r, 13, "system_type")
	if err != nil {
		t.Fatal(err)
	}
	ifVersion, err := readUint8(r, "interface_version")
	if err != nil {
		t.Fatal(err)
	}
	addrTON, err := readUint8(r, "addr_ton")
	if err != nil {
		t.Fatal(err)
	}
	addrNPI, err := readUint8(r, "addr_npi")
	if err != nil {
		t.Fatal(err)
	}
	addrRange, err := readCString(r, 41, "address_range")
	if err != nil {
		t.Fatal(err)
	}

	if systemID != "SMPP3TEST" || password != "secret08" || systemType != "SUBMIT1" {
		t.Errorf("string field'lar: %q/%q/%q, kutilgan SMPP3TEST/secret08/SUBMIT1",
			systemID, password, systemType)
	}
	// Spec misolining g'alatiligi: interface_version=0x00 (0x34 emas!).
	if ifVersion != 0x00 {
		t.Errorf("interface_version = 0x%02X, spec misolida 0x00", ifVersion)
	}
	if addrTON != 0x01 || addrNPI != 0x01 || addrRange != "" {
		t.Errorf("addr uchligi: %d/%d/%q, kutilgan 1/1/\"\"", addrTON, addrNPI, addrRange)
	}
	if r.Len() != 0 {
		t.Errorf("body'da %d oktet ortib qoldi", r.Len())
	}
}

func TestDecodeHeaderTooShort(t *testing.T) {
	if _, err := DecodeHeader(make([]byte, 10)); err == nil {
		t.Error("10 oktetlik data'da DecodeHeader xato qaytarishi kerak edi")
	}
}

func TestHeaderRoundTrip(t *testing.T) {
	h := Header{Length: 0x100, ID: CmdDeliverSM, Status: 0, Sequence: MaxSequence}
	enc := EncodeHeader(h)
	got, err := DecodeHeader(enc[:])
	if err != nil {
		t.Fatalf("DecodeHeader xatosi: %v", err)
	}
	if got != h {
		t.Errorf("round-trip: %+v != %+v", got, h)
	}
}
