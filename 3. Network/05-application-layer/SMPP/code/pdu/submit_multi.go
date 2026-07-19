package pdu

import (
	"bytes"
	"fmt"

	"smpp/tlv"
)

// submit_multi chegaralari (§5.2.24–5.2.27).
const (
	maxDests  = 254 // number_of_dests 1 oktet, 0 va 255 ma'nosiz (§5.2.24)
	maxDLName = 21  // dl_name — C-Octet, max 21 (§5.2.27)
)

// dest_flag qiymatlari (§5.2.25, Table 5-4a).
const (
	destFlagSME      uint8 = 1 // SME Address: ton + npi + addr
	destFlagDistList uint8 = 2 // Distribution List: dl_name
)

// DestAddress — submit_multi'dagi bitta qabul qiluvchi: UNION tipi
// (§4.5.1, Table 4-13/14/15): yo SME manzili, yo SMSC'da oldindan
// ro'yxatlangan Distribution List nomi. DLName bo'sh bo'lmasa — DL rejimi,
// aks holda SME rejimi.
type DestAddress struct {
	SME    Address // dest_flag=1 bo'lganda
	DLName string  // dest_flag=2 bo'lganda (bo'sh emas = DL rejimi)
}

// IsDistList qabul qiluvchi Distribution List ekanini bildiradi.
func (d DestAddress) IsDistList() bool { return d.DLName != "" }

// SubmitMulti — submit_multi PDU (§4.5): bitta xabarni 254 tagacha qabul
// qiluvchiga (yoki DL'larga) yuborish. Body submit_sm bilan deyarli bir xil,
// faqat yakka dest bloki o'rniga number_of_dests + dest_address ro'yxati.
//
// MUHIM farq (§4.5.1 izohi, 81-b): replace_if_present_flag bu yerda
// RESERVED — NULL bo'lishi shart; shuning uchun struct'da bunday field yo'q
// (encoder doim 0 yozadi). Transaction mode ham qo'llanmaydi (§4.5).
type SubmitMulti struct {
	ServiceType          string
	Source               Address
	Dests                []DestAddress
	EsmClass             EsmClass
	ProtocolID           uint8
	PriorityFlag         uint8
	ScheduleDeliveryTime string
	ValidityPeriod       string
	RegisteredDelivery   RegisteredDelivery
	DataCoding           uint8
	SMDefaultMsgID       uint8
	ShortMessage         []byte
	TLVs                 []tlv.TLV
}

// Encode to'liq wire frame yasaydi.
func (s SubmitMulti) Encode(seq uint32) ([]byte, error) {
	if len(s.Dests) == 0 || len(s.Dests) > maxDests {
		return nil, fmt.Errorf("pdu: submit_multi dest soni %d — 1..%d bo'lishi kerak (§5.2.24)", len(s.Dests), maxDests)
	}
	if len(s.ShortMessage) > maxShortMsg {
		return nil, fmt.Errorf("pdu: short_message %d oktet, max %d", len(s.ShortMessage), maxShortMsg)
	}
	var b bytes.Buffer
	if err := writeCString(&b, s.ServiceType, maxServiceType, "service_type"); err != nil {
		return nil, err
	}
	if err := writeAddress(&b, s.Source, "source_addr"); err != nil {
		return nil, err
	}
	writeUint8(&b, uint8(len(s.Dests)))
	for i, d := range s.Dests {
		if d.IsDistList() {
			writeUint8(&b, destFlagDistList)
			if err := writeCString(&b, d.DLName, maxDLName, "dl_name"); err != nil {
				return nil, err
			}
		} else {
			writeUint8(&b, destFlagSME)
			if err := writeAddress(&b, d.SME, fmt.Sprintf("dest_address[%d]", i)); err != nil {
				return nil, err
			}
		}
	}
	writeUint8(&b, uint8(s.EsmClass))
	writeUint8(&b, s.ProtocolID)
	writeUint8(&b, s.PriorityFlag)
	if err := writeTimeField(&b, s.ScheduleDeliveryTime, "schedule_delivery_time"); err != nil {
		return nil, err
	}
	if err := writeTimeField(&b, s.ValidityPeriod, "validity_period"); err != nil {
		return nil, err
	}
	writeUint8(&b, uint8(s.RegisteredDelivery))
	writeUint8(&b, 0x00) // replace_if_present_flag — Reserved (§4.5.1)
	writeUint8(&b, s.DataCoding)
	writeUint8(&b, s.SMDefaultMsgID)
	writeUint8(&b, uint8(len(s.ShortMessage)))
	b.Write(s.ShortMessage)
	if err := tlv.Encode(&b, s.TLVs); err != nil {
		return nil, err
	}
	return encodePDU(CmdSubmitMulti, 0, seq, b.Bytes()), nil
}

// DecodeSubmitMulti to'liq frame'dan SubmitMulti o'qiydi.
func DecodeSubmitMulti(frame []byte) (SubmitMulti, Header, error) {
	h, err := DecodeHeader(frame)
	if err != nil {
		return SubmitMulti{}, Header{}, err
	}
	if h.ID != CmdSubmitMulti {
		return SubmitMulti{}, h, fmt.Errorf("pdu: %s submit_multi emas", h.ID)
	}
	r := bytes.NewReader(frame[HeaderSize:])
	var s SubmitMulti
	if s.ServiceType, err = readCString(r, maxServiceType, "service_type"); err != nil {
		return s, h, err
	}
	if s.Source, err = readAddress(r, "source_addr"); err != nil {
		return s, h, err
	}
	n, err := readUint8(r, "number_of_dests")
	if err != nil {
		return s, h, err
	}
	for i := 0; i < int(n); i++ {
		flag, err := readUint8(r, "dest_flag")
		if err != nil {
			return s, h, err
		}
		var d DestAddress
		switch flag {
		case destFlagSME:
			if d.SME, err = readAddress(r, fmt.Sprintf("dest_address[%d]", i)); err != nil {
				return s, h, err
			}
		case destFlagDistList:
			if d.DLName, err = readCString(r, maxDLName, "dl_name"); err != nil {
				return s, h, err
			}
		default:
			return s, h, fmt.Errorf("pdu: dest_flag=%d — 1 (SME) yoki 2 (DL) kutilgan (§5.2.25)", flag)
		}
		s.Dests = append(s.Dests, d)
	}
	var esm, reg uint8
	if esm, err = readUint8(r, "esm_class"); err != nil {
		return s, h, err
	}
	s.EsmClass = EsmClass(esm)
	if s.ProtocolID, err = readUint8(r, "protocol_id"); err != nil {
		return s, h, err
	}
	if s.PriorityFlag, err = readUint8(r, "priority_flag"); err != nil {
		return s, h, err
	}
	if s.ScheduleDeliveryTime, err = readCString(r, maxTimeField, "schedule_delivery_time"); err != nil {
		return s, h, err
	}
	if s.ValidityPeriod, err = readCString(r, maxTimeField, "validity_period"); err != nil {
		return s, h, err
	}
	if reg, err = readUint8(r, "registered_delivery"); err != nil {
		return s, h, err
	}
	s.RegisteredDelivery = RegisteredDelivery(reg)
	if _, err = readUint8(r, "replace_if_present_flag"); err != nil { // Reserved — o'qib tashlanadi
		return s, h, err
	}
	if s.DataCoding, err = readUint8(r, "data_coding"); err != nil {
		return s, h, err
	}
	if s.SMDefaultMsgID, err = readUint8(r, "sm_default_msg_id"); err != nil {
		return s, h, err
	}
	smLen, err := readUint8(r, "sm_length")
	if err != nil {
		return s, h, err
	}
	if int(smLen) > r.Len() {
		return s, h, fmt.Errorf("pdu: sm_length=%d, lekin frame'da %d oktet qoldi", smLen, r.Len())
	}
	if smLen > 0 {
		s.ShortMessage = make([]byte, smLen)
		if _, err := r.Read(s.ShortMessage); err != nil {
			return s, h, fmt.Errorf("pdu: short_message o'qishda: %w", err)
		}
	}
	if s.TLVs, err = readTLVTail(r); err != nil {
		return s, h, err
	}
	return s, h, nil
}

// UnsuccessSME — submit_multi_resp'dagi bitta muvaffaqiyatsiz qabul qiluvchi
// (§4.5.2, Table 4-17): manzil + unga tegishli command_status qiymati.
type UnsuccessSME struct {
	Addr            Address
	ErrorStatusCode uint32 // command_status fazosidagi kod (Table 5-2)
}

// SubmitMultiResp — submit_multi_resp (§4.5.2): message_id + muvaffaqiyatsiz
// manzillar ro'yxati. E'tibor: PDU'ning UMUMIY Status'i 0 bo'lishi mumkin,
// lekin ayrim manzillar baribir muvaffaqiyatsiz — "qisman muvaffaqiyat"
// submit_multi'da NORMAL holat, faqat Unsuccess bo'sh bo'lsa hammasi qabul
// qilingan.
type SubmitMultiResp struct {
	Status    uint32
	MessageID string
	Unsuccess []UnsuccessSME
}

// Encode to'liq wire frame yasaydi. Status != 0 → body YO'Q.
func (r SubmitMultiResp) Encode(seq uint32) ([]byte, error) {
	if r.Status != 0 {
		return encodePDU(CmdSubmitMultiResp, r.Status, seq, nil), nil
	}
	if len(r.Unsuccess) > maxDests {
		return nil, fmt.Errorf("pdu: no_unsuccess %d — max %d", len(r.Unsuccess), maxDests)
	}
	var b bytes.Buffer
	if err := writeCString(&b, r.MessageID, maxMessageID, "message_id"); err != nil {
		return nil, err
	}
	writeUint8(&b, uint8(len(r.Unsuccess)))
	for i, u := range r.Unsuccess {
		if err := writeAddress(&b, u.Addr, fmt.Sprintf("unsuccess_sme[%d]", i)); err != nil {
			return nil, err
		}
		writeUint32(&b, u.ErrorStatusCode)
	}
	return encodePDU(CmdSubmitMultiResp, 0, seq, b.Bytes()), nil
}

// DecodeSubmitMultiResp to'liq frame'dan resp o'qiydi.
func DecodeSubmitMultiResp(frame []byte) (SubmitMultiResp, Header, error) {
	h, err := DecodeHeader(frame)
	if err != nil {
		return SubmitMultiResp{}, Header{}, err
	}
	if h.ID != CmdSubmitMultiResp {
		return SubmitMultiResp{}, h, fmt.Errorf("pdu: %s submit_multi_resp emas", h.ID)
	}
	resp := SubmitMultiResp{Status: h.Status}
	if h.Status != 0 {
		return resp, h, nil
	}
	r := bytes.NewReader(frame[HeaderSize:])
	if resp.MessageID, err = readCString(r, maxMessageID, "message_id"); err != nil {
		return SubmitMultiResp{}, h, err
	}
	n, err := readUint8(r, "no_unsuccess")
	if err != nil {
		return SubmitMultiResp{}, h, err
	}
	for i := 0; i < int(n); i++ {
		var u UnsuccessSME
		if u.Addr, err = readAddress(r, fmt.Sprintf("unsuccess_sme[%d]", i)); err != nil {
			return SubmitMultiResp{}, h, err
		}
		if u.ErrorStatusCode, err = readUint32(r, "error_status_code"); err != nil {
			return SubmitMultiResp{}, h, err
		}
		resp.Unsuccess = append(resp.Unsuccess, u)
	}
	return resp, h, nil
}
