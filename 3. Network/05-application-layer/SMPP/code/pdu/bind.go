package pdu

import (
	"bytes"
	"fmt"

	"smpp/tlv"
)

// Bind field o'lchamlari — NULL terminator BILAN (v3.4 §5.2.1–5.2.7, §3.1 note iii).
const (
	maxSystemID     = 16
	maxPassword     = 9
	maxSystemType   = 13
	maxAddressRange = 41
)

// InterfaceVersion34 — v3.4 uchun interface_version qiymati (§5.2.4).
// 0x00–0x33 diapazoni "v3.3 yoki undan eski" deb talqin qilinadi.
const InterfaceVersion34 uint8 = 0x34

// Bind — bind_transmitter/bind_receiver/bind_transceiver'ning umumiy tanasi:
// uchala PDU'ning body'si bir xil 7 field (v3.4 §4.1, Table 4-1/4-3/4-5),
// faqat command_id farq qiladi — u Mode'da.
type Bind struct {
	Mode             CommandID // CmdBindTransmitter, CmdBindReceiver yoki CmdBindTransceiver
	SystemID         string    // ESME identifikatori (§5.2.1)
	Password         string    // autentifikatsiya; bo'sh bo'lishi mumkin (§5.2.2)
	SystemType       string    // ESME kategoriyasi ("VMS"...); odatda bo'sh (§5.2.3)
	InterfaceVersion uint8     // ESME qo'llagan SMPP versiyasi (§5.2.4)
	AddrTON          uint8     // §5.2.5; noma'lum bo'lsa 0
	AddrNPI          uint8     // §5.2.6; noma'lum bo'lsa 0
	AddressRange     string    // UNIX regex, RX/TRX routing uchun; odatda bo'sh (§5.2.7)
}

func isBindReq(id CommandID) bool {
	return id == CmdBindTransmitter || id == CmdBindReceiver || id == CmdBindTransceiver
}

// Encode to'liq wire frame yasaydi. Request bo'lgani uchun command_status=0 (§5.1.3).
func (b Bind) Encode(seq uint32) ([]byte, error) {
	if !isBindReq(b.Mode) {
		return nil, fmt.Errorf("pdu: %s bind request emas", b.Mode)
	}
	var body bytes.Buffer
	if err := writeCString(&body, b.SystemID, maxSystemID, "system_id"); err != nil {
		return nil, err
	}
	if err := writeCString(&body, b.Password, maxPassword, "password"); err != nil {
		return nil, err
	}
	if err := writeCString(&body, b.SystemType, maxSystemType, "system_type"); err != nil {
		return nil, err
	}
	writeUint8(&body, b.InterfaceVersion)
	writeUint8(&body, b.AddrTON)
	writeUint8(&body, b.AddrNPI)
	if err := writeCString(&body, b.AddressRange, maxAddressRange, "address_range"); err != nil {
		return nil, err
	}
	return encodePDU(b.Mode, 0, seq, body.Bytes()), nil
}

// DecodeBind to'liq frame'dan (header bilan) Bind'ni o'qiydi.
func DecodeBind(frame []byte) (Bind, Header, error) {
	h, err := DecodeHeader(frame)
	if err != nil {
		return Bind{}, Header{}, err
	}
	if !isBindReq(h.ID) {
		return Bind{}, h, fmt.Errorf("pdu: %s bind request emas", h.ID)
	}
	r := bytes.NewReader(frame[HeaderSize:])
	b := Bind{Mode: h.ID}
	if b.SystemID, err = readCString(r, maxSystemID, "system_id"); err != nil {
		return Bind{}, h, err
	}
	if b.Password, err = readCString(r, maxPassword, "password"); err != nil {
		return Bind{}, h, err
	}
	if b.SystemType, err = readCString(r, maxSystemType, "system_type"); err != nil {
		return Bind{}, h, err
	}
	if b.InterfaceVersion, err = readUint8(r, "interface_version"); err != nil {
		return Bind{}, h, err
	}
	if b.AddrTON, err = readUint8(r, "addr_ton"); err != nil {
		return Bind{}, h, err
	}
	if b.AddrNPI, err = readUint8(r, "addr_npi"); err != nil {
		return Bind{}, h, err
	}
	if b.AddressRange, err = readCString(r, maxAddressRange, "address_range"); err != nil {
		return Bind{}, h, err
	}
	return b, h, nil
}

// BindResp — bind_*_resp: muvaffaqiyatda body = SMSC'ning system_id'si +
// ixtiyoriy sc_interface_version TLV (§4.1.2). MUHIM: command_status != 0
// bo'lsa body UMUMAN qaytarilmaydi (§4.1.2 Note) — encoder ham, decoder ham
// shu qoidaga amal qiladi.
type BindResp struct {
	Mode         CommandID // CmdBindTransmitterResp, CmdBindReceiverResp yoki CmdBindTransceiverResp
	Status       uint32    // command_status; 0 = muvaffaqiyat
	SystemID     string    // SMSC identifikatori (faqat Status=0'da)
	SCVersion    uint8     // sc_interface_version TLV qiymati (§5.3.2.25)
	HasSCVersion bool      // TLV kelgan-kelmagani — YO'QLIGI ham ma'lumot (§3.4)
}

func isBindResp(id CommandID) bool {
	return id == CmdBindTransmitterResp || id == CmdBindReceiverResp || id == CmdBindTransceiverResp
}

// Encode to'liq wire frame yasaydi. Status != 0 bo'lsa faqat header ketadi.
func (br BindResp) Encode(seq uint32) ([]byte, error) {
	if !isBindResp(br.Mode) {
		return nil, fmt.Errorf("pdu: %s bind response emas", br.Mode)
	}
	if br.Status != 0 {
		return encodePDU(br.Mode, br.Status, seq, nil), nil
	}
	var body bytes.Buffer
	if err := writeCString(&body, br.SystemID, maxSystemID, "system_id"); err != nil {
		return nil, err
	}
	if br.HasSCVersion {
		if err := tlv.Encode(&body, []tlv.TLV{tlv.U8(tlv.ScInterfaceVersion, br.SCVersion)}); err != nil {
			return nil, err
		}
	}
	return encodePDU(br.Mode, 0, seq, body.Bytes()), nil
}

// DecodeBindResp to'liq frame'dan BindResp'ni o'qiydi. Status != 0 bo'lsa
// body parse QILINMAYDI (spec bo'yicha u yo'q; kelgan taqdirda ham ishonchsiz).
func DecodeBindResp(frame []byte) (BindResp, Header, error) {
	h, err := DecodeHeader(frame)
	if err != nil {
		return BindResp{}, Header{}, err
	}
	if !isBindResp(h.ID) {
		return BindResp{}, h, fmt.Errorf("pdu: %s bind response emas", h.ID)
	}
	br := BindResp{Mode: h.ID, Status: h.Status}
	if h.Status != 0 {
		return br, h, nil
	}
	r := bytes.NewReader(frame[HeaderSize:])
	if br.SystemID, err = readCString(r, maxSystemID, "system_id"); err != nil {
		return BindResp{}, h, err
	}
	rest := make([]byte, r.Len())
	if _, err := r.Read(rest); err != nil && r.Len() > 0 {
		return BindResp{}, h, err
	}
	tlvs, err := tlv.Decode(rest)
	if err != nil {
		return BindResp{}, h, err
	}
	if v, ok := tlv.Find(tlvs, tlv.ScInterfaceVersion); ok {
		if sc, scErr := v.Uint8Value(); scErr == nil {
			br.SCVersion, br.HasSCVersion = sc, true
		}
	}
	return br, h, nil
}

// Outbind (§2.2.1, §4.1.7): SMSC o'zi ESME'ga ulanib "bind_receiver qil" deb
// signal beradi; password bu yerda SMSC'ni ESME'ga tanitadi. Resp PDU'si YO'Q —
// javob ESME'ning bind_receiver'i, rad etish esa TCP'ni uzish.
type Outbind struct {
	SystemID string
	Password string
}

// Encode to'liq wire frame yasaydi.
func (o Outbind) Encode(seq uint32) ([]byte, error) {
	var body bytes.Buffer
	if err := writeCString(&body, o.SystemID, maxSystemID, "system_id"); err != nil {
		return nil, err
	}
	if err := writeCString(&body, o.Password, maxPassword, "password"); err != nil {
		return nil, err
	}
	return encodePDU(CmdOutbind, 0, seq, body.Bytes()), nil
}

// DecodeOutbind to'liq frame'dan Outbind'ni o'qiydi.
func DecodeOutbind(frame []byte) (Outbind, Header, error) {
	h, err := DecodeHeader(frame)
	if err != nil {
		return Outbind{}, Header{}, err
	}
	if h.ID != CmdOutbind {
		return Outbind{}, h, fmt.Errorf("pdu: %s outbind emas", h.ID)
	}
	r := bytes.NewReader(frame[HeaderSize:])
	var o Outbind
	if o.SystemID, err = readCString(r, maxSystemID, "system_id"); err != nil {
		return Outbind{}, h, err
	}
	if o.Password, err = readCString(r, maxPassword, "password"); err != nil {
		return Outbind{}, h, err
	}
	return o, h, nil
}
