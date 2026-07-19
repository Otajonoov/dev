package pdu

import (
	"bytes"
	"fmt"

	"smpp/tlv"
)

// Field o'lchamlari (NULL bilan).
const (
	maxServiceType = 6   // §5.2.11
	maxMessageID   = 65  // §5.2.23
	maxShortMsg    = 254 // §5.2.21: sm_length 0–254, 255 taqiqlangan
	maxTimeField   = 17  // §7.1.1: "1 or 17"
)

// SMFields — submit_sm va deliver_sm'ning UMUMIY body'si: format bir xil
// (§4.6.1: "The deliver_sm PDU has the same format as the submit_sm PDU"),
// farq faqat semantikada. Bitta codec — ikki PDU.
type SMFields struct {
	ServiceType          string
	Source               Address
	Dest                 Address
	EsmClass             EsmClass
	ProtocolID           uint8
	PriorityFlag         uint8
	ScheduleDeliveryTime string // "" yoki 16 belgi (§7.1.1)
	ValidityPeriod       string
	RegisteredDelivery   RegisteredDelivery
	ReplaceIfPresent     uint8
	DataCoding           uint8
	SMDefaultMsgID       uint8
	ShortMessage         []byte
	TLVs                 []tlv.TLV
}

// writeTimeField "1 or 17" qoidasi (§7.1.1): bo'sh = yagona NULL,
// aks holda aynan 16 belgi + NULL.
func writeTimeField(b *bytes.Buffer, s, field string) error {
	if s == "" {
		b.WriteByte(0x00)
		return nil
	}
	if len(s) != 16 {
		return fmt.Errorf("pdu: %s vaqt field'i aynan 16 belgi bo'lishi kerak, keldi %d", field, len(s))
	}
	return writeCString(b, s, maxTimeField, field)
}

func (m SMFields) encode() ([]byte, error) {
	if len(m.ShortMessage) > maxShortMsg {
		return nil, fmt.Errorf("pdu: short_message %d oktet, max %d — uzun matn uchun message_payload TLV (§3.2.3)", len(m.ShortMessage), maxShortMsg)
	}
	if _, ok := tlv.Find(m.TLVs, tlv.MessagePayload); ok && len(m.ShortMessage) > 0 {
		return nil, fmt.Errorf("pdu: short_message va message_payload birga taqiqlangan (§3.2.3)")
	}
	var b bytes.Buffer
	if err := writeCString(&b, m.ServiceType, maxServiceType, "service_type"); err != nil {
		return nil, err
	}
	if err := writeAddress(&b, m.Source, "source_addr"); err != nil {
		return nil, err
	}
	if err := writeAddress(&b, m.Dest, "destination_addr"); err != nil {
		return nil, err
	}
	writeUint8(&b, uint8(m.EsmClass))
	writeUint8(&b, m.ProtocolID)
	writeUint8(&b, m.PriorityFlag)
	if err := writeTimeField(&b, m.ScheduleDeliveryTime, "schedule_delivery_time"); err != nil {
		return nil, err
	}
	if err := writeTimeField(&b, m.ValidityPeriod, "validity_period"); err != nil {
		return nil, err
	}
	writeUint8(&b, uint8(m.RegisteredDelivery))
	writeUint8(&b, m.ReplaceIfPresent)
	writeUint8(&b, m.DataCoding)
	writeUint8(&b, m.SMDefaultMsgID)
	writeUint8(&b, uint8(len(m.ShortMessage)))
	b.Write(m.ShortMessage)
	if err := tlv.Encode(&b, m.TLVs); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func decodeSMBody(data []byte) (SMFields, error) {
	var m SMFields
	r := bytes.NewReader(data)
	var err error
	if m.ServiceType, err = readCString(r, maxServiceType, "service_type"); err != nil {
		return m, err
	}
	if m.Source, err = readAddress(r, "source_addr"); err != nil {
		return m, err
	}
	if m.Dest, err = readAddress(r, "destination_addr"); err != nil {
		return m, err
	}
	var esm, reg uint8
	if esm, err = readUint8(r, "esm_class"); err != nil {
		return m, err
	}
	m.EsmClass = EsmClass(esm)
	if m.ProtocolID, err = readUint8(r, "protocol_id"); err != nil {
		return m, err
	}
	if m.PriorityFlag, err = readUint8(r, "priority_flag"); err != nil {
		return m, err
	}
	if m.ScheduleDeliveryTime, err = readCString(r, maxTimeField, "schedule_delivery_time"); err != nil {
		return m, err
	}
	if m.ValidityPeriod, err = readCString(r, maxTimeField, "validity_period"); err != nil {
		return m, err
	}
	if reg, err = readUint8(r, "registered_delivery"); err != nil {
		return m, err
	}
	m.RegisteredDelivery = RegisteredDelivery(reg)
	if m.ReplaceIfPresent, err = readUint8(r, "replace_if_present_flag"); err != nil {
		return m, err
	}
	if m.DataCoding, err = readUint8(r, "data_coding"); err != nil {
		return m, err
	}
	if m.SMDefaultMsgID, err = readUint8(r, "sm_default_msg_id"); err != nil {
		return m, err
	}
	smLen, err := readUint8(r, "sm_length")
	if err != nil {
		return m, err
	}
	if int(smLen) > r.Len() {
		return m, fmt.Errorf("pdu: sm_length=%d, lekin frame'da %d oktet qoldi", smLen, r.Len())
	}
	if smLen > 0 {
		m.ShortMessage = make([]byte, smLen)
		if _, err := r.Read(m.ShortMessage); err != nil {
			return m, fmt.Errorf("pdu: short_message o'qishda: %w", err)
		}
	}
	// Qolgan hamma bayt — TLV tail (§3.2.4).
	rest := make([]byte, r.Len())
	if len(rest) > 0 {
		if _, err := r.Read(rest); err != nil {
			return m, fmt.Errorf("pdu: TLV tail o'qishda: %w", err)
		}
	}
	if m.TLVs, err = tlv.Decode(rest); err != nil {
		return m, err
	}
	return m, nil
}

// SubmitSM — submit_sm PDU (§4.4.1, Table 4-10): ESME→SMSC xabar yuborish.
// Transaction messaging mode'ni QO'LLAMAYDI (§4.4).
type SubmitSM struct {
	SMFields
}

// Encode to'liq wire frame yasaydi (request → command_status=0).
func (s SubmitSM) Encode(seq uint32) ([]byte, error) {
	body, err := s.SMFields.encode()
	if err != nil {
		return nil, err
	}
	return encodePDU(CmdSubmitSM, 0, seq, body), nil
}

// DecodeSubmitSM to'liq frame'dan SubmitSM o'qiydi.
func DecodeSubmitSM(frame []byte) (SubmitSM, Header, error) {
	h, err := DecodeHeader(frame)
	if err != nil {
		return SubmitSM{}, Header{}, err
	}
	if h.ID != CmdSubmitSM {
		return SubmitSM{}, h, fmt.Errorf("pdu: %s submit_sm emas", h.ID)
	}
	m, err := decodeSMBody(frame[HeaderSize:])
	if err != nil {
		return SubmitSM{}, h, err
	}
	return SubmitSM{SMFields: m}, h, nil
}

// SubmitSMResp — submit_sm_resp (§4.4.2): body faqat message_id.
// MessageID — SMSC bergan OPAQUE qiymat (§5.2.23): formati noma'lum,
// songa aylantirilmaydi, faqat string sifatida saqlanadi (9-bob: hex/dec tuzog'i).
type SubmitSMResp struct {
	Status    uint32
	MessageID string
}

// Encode to'liq wire frame yasaydi. Status != 0 → body YO'Q (§4.4.2).
func (r SubmitSMResp) Encode(seq uint32) ([]byte, error) {
	if r.Status != 0 {
		return encodePDU(CmdSubmitSMResp, r.Status, seq, nil), nil
	}
	var b bytes.Buffer
	if err := writeCString(&b, r.MessageID, maxMessageID, "message_id"); err != nil {
		return nil, err
	}
	return encodePDU(CmdSubmitSMResp, 0, seq, b.Bytes()), nil
}

// DecodeSubmitSMResp to'liq frame'dan resp o'qiydi. Status != 0 → body parse yo'q.
func DecodeSubmitSMResp(frame []byte) (SubmitSMResp, Header, error) {
	h, err := DecodeHeader(frame)
	if err != nil {
		return SubmitSMResp{}, Header{}, err
	}
	if h.ID != CmdSubmitSMResp {
		return SubmitSMResp{}, h, fmt.Errorf("pdu: %s submit_sm_resp emas", h.ID)
	}
	resp := SubmitSMResp{Status: h.Status}
	if h.Status != 0 {
		return resp, h, nil
	}
	r := bytes.NewReader(frame[HeaderSize:])
	if resp.MessageID, err = readCString(r, maxMessageID, "message_id"); err != nil {
		return SubmitSMResp{}, h, err
	}
	return resp, h, nil
}
