package pdu

import (
	"bytes"
	"fmt"
)

// DeliverSM — deliver_sm PDU (§4.6): SMSC→ESME. Format submit_sm bilan
// BIR XIL (§4.6.1), semantikasi farq qiladi: MO xabar yoki DLR/ack —
// turi EsmClass.MessageType()'dan aniqlanadi.
//
// Ishlatilmaydigan field'lar NULL bo'lishi SHART (§4.6.1):
// schedule_delivery_time, validity_period, replace_if_present_flag,
// sm_default_msg_id — encoder shuni tekshiradi; decoder esa tolerant
// (spec buzgan SMSC'dan kelgan PDU baribir o'qiladi).
type DeliverSM struct {
	SMFields
}

// Encode to'liq wire frame yasaydi (server tomon — 14-bob mock SMSC ishlatadi).
func (d DeliverSM) Encode(seq uint32) ([]byte, error) {
	if d.ScheduleDeliveryTime != "" || d.ValidityPeriod != "" {
		return nil, fmt.Errorf("pdu: deliver_sm'da schedule/validity NULL bo'lishi shart (§4.6.1)")
	}
	if d.ReplaceIfPresent != 0 || d.SMDefaultMsgID != 0 {
		return nil, fmt.Errorf("pdu: deliver_sm'da replace_if_present/sm_default_msg_id NULL bo'lishi shart (§4.6.1)")
	}
	body, err := d.SMFields.encode()
	if err != nil {
		return nil, err
	}
	return encodePDU(CmdDeliverSM, 0, seq, body), nil
}

// DecodeDeliverSM to'liq frame'dan DeliverSM o'qiydi.
func DecodeDeliverSM(frame []byte) (DeliverSM, Header, error) {
	h, err := DecodeHeader(frame)
	if err != nil {
		return DeliverSM{}, Header{}, err
	}
	if h.ID != CmdDeliverSM {
		return DeliverSM{}, h, fmt.Errorf("pdu: %s deliver_sm emas", h.ID)
	}
	m, err := decodeSMBody(frame[HeaderSize:])
	if err != nil {
		return DeliverSM{}, h, err
	}
	return DeliverSM{SMFields: m}, h, nil
}

// DeliverSMResp — deliver_sm_resp (§4.6.2): body'dagi message_id "unused,
// set to NULL" — ya'ni doim bitta 0x00 okteti. MessageID field'i yo'q —
// tashish ham, o'qish ham ma'nosiz.
type DeliverSMResp struct {
	Status uint32
}

// Encode to'liq wire frame yasaydi. Status qandayligidan qat'i nazar
// NULL message_id yoziladi — field mandatory, qiymati doim bo'sh.
func (r DeliverSMResp) Encode(seq uint32) []byte {
	return encodePDU(CmdDeliverSMResp, r.Status, seq, []byte{0x00})
}

// DecodeDeliverSMResp to'liq frame'dan resp o'qiydi. Body'siz (16 oktet)
// variantga ham toqat qilinadi — ayrim stack'lar shunday yuboradi.
func DecodeDeliverSMResp(frame []byte) (DeliverSMResp, Header, error) {
	h, err := DecodeHeader(frame)
	if err != nil {
		return DeliverSMResp{}, Header{}, err
	}
	if h.ID != CmdDeliverSMResp {
		return DeliverSMResp{}, h, fmt.Errorf("pdu: %s deliver_sm_resp emas", h.ID)
	}
	if len(frame) > HeaderSize {
		r := bytes.NewReader(frame[HeaderSize:])
		if _, err := readCString(r, maxMessageID, "message_id"); err != nil {
			return DeliverSMResp{}, h, err
		}
	}
	return DeliverSMResp{Status: h.Status}, h, nil
}
