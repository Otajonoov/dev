package pdu

import (
	"bytes"
	"fmt"
)

// ReplaceSM — replace_sm PDU (§4.10): hali yetkazilmagan xabarni YANGISI
// bilan almashtirish. Matching: message_id + source (original NULL bo'lsa
// bunda ham NULL). submit_sm+replace_if_present'dan tub farqi: mos xabar
// TOPILMASA yangi xabar yaratilmaydi — XATO qaytadi (ESME_RREPLACEFAIL).
//
// Body'da destination_addr YO'Q (manzil almashtirilmaydi) va data_coding
// YO'Q — yangi short_message ORIGINAL xabarning kodlashida talqin qilinadi:
// UCS2 xabarni GSM7 matn bilan almashtirsangiz telefon axlat ko'radi.
// Bu cheklov API'da ham atayin ko'rinib turadi: struct'da bunday field yo'q.
type ReplaceSM struct {
	MessageID            string
	Source               Address
	ScheduleDeliveryTime string // "" yoki 16 belgi (§7.1.1)
	ValidityPeriod       string
	RegisteredDelivery   RegisteredDelivery
	SMDefaultMsgID       uint8
	ShortMessage         []byte
}

// Encode to'liq wire frame yasaydi.
func (p ReplaceSM) Encode(seq uint32) ([]byte, error) {
	if len(p.ShortMessage) > maxShortMsg {
		return nil, fmt.Errorf("pdu: short_message %d oktet, max %d", len(p.ShortMessage), maxShortMsg)
	}
	var b bytes.Buffer
	if err := writeCString(&b, p.MessageID, maxMessageID, "message_id"); err != nil {
		return nil, err
	}
	if err := writeAddress(&b, p.Source, "source_addr"); err != nil {
		return nil, err
	}
	if err := writeTimeField(&b, p.ScheduleDeliveryTime, "schedule_delivery_time"); err != nil {
		return nil, err
	}
	if err := writeTimeField(&b, p.ValidityPeriod, "validity_period"); err != nil {
		return nil, err
	}
	writeUint8(&b, uint8(p.RegisteredDelivery))
	writeUint8(&b, p.SMDefaultMsgID)
	writeUint8(&b, uint8(len(p.ShortMessage)))
	b.Write(p.ShortMessage)
	return encodePDU(CmdReplaceSM, 0, seq, b.Bytes()), nil
}

// DecodeReplaceSM to'liq frame'dan ReplaceSM o'qiydi.
func DecodeReplaceSM(frame []byte) (ReplaceSM, Header, error) {
	h, err := DecodeHeader(frame)
	if err != nil {
		return ReplaceSM{}, Header{}, err
	}
	if h.ID != CmdReplaceSM {
		return ReplaceSM{}, h, fmt.Errorf("pdu: %s replace_sm emas", h.ID)
	}
	r := bytes.NewReader(frame[HeaderSize:])
	var p ReplaceSM
	if p.MessageID, err = readCString(r, maxMessageID, "message_id"); err != nil {
		return p, h, err
	}
	if p.Source, err = readAddress(r, "source_addr"); err != nil {
		return p, h, err
	}
	if p.ScheduleDeliveryTime, err = readCString(r, maxTimeField, "schedule_delivery_time"); err != nil {
		return p, h, err
	}
	if p.ValidityPeriod, err = readCString(r, maxTimeField, "validity_period"); err != nil {
		return p, h, err
	}
	var reg uint8
	if reg, err = readUint8(r, "registered_delivery"); err != nil {
		return p, h, err
	}
	p.RegisteredDelivery = RegisteredDelivery(reg)
	if p.SMDefaultMsgID, err = readUint8(r, "sm_default_msg_id"); err != nil {
		return p, h, err
	}
	smLen, err := readUint8(r, "sm_length")
	if err != nil {
		return p, h, err
	}
	if int(smLen) > r.Len() {
		return p, h, fmt.Errorf("pdu: sm_length=%d, lekin frame'da %d oktet qoldi", smLen, r.Len())
	}
	if smLen > 0 {
		p.ShortMessage = make([]byte, smLen)
		if _, err := r.Read(p.ShortMessage); err != nil {
			return p, h, fmt.Errorf("pdu: short_message o'qishda: %w", err)
		}
	}
	return p, h, nil
}

// ReplaceSMResp — replace_sm_resp (§4.10.2): body YO'Q, faqat header.
type ReplaceSMResp struct {
	Status uint32
}

// Encode to'liq wire frame yasaydi.
func (r ReplaceSMResp) Encode(seq uint32) []byte {
	return encodePDU(CmdReplaceSMResp, r.Status, seq, nil)
}

// DecodeReplaceSMResp to'liq frame'dan resp o'qiydi.
func DecodeReplaceSMResp(frame []byte) (ReplaceSMResp, Header, error) {
	h, err := DecodeHeader(frame)
	if err != nil {
		return ReplaceSMResp{}, Header{}, err
	}
	if h.ID != CmdReplaceSMResp {
		return ReplaceSMResp{}, h, fmt.Errorf("pdu: %s replace_sm_resp emas", h.ID)
	}
	return ReplaceSMResp{Status: h.Status}, h, nil
}
