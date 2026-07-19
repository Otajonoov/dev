package pdu

import (
	"bytes"
	"fmt"
)

// CancelSM — cancel_sm PDU (§4.9): hali yetkazilmagan xabar(lar)ni bekor
// qilish. IKKI rejim (§4.9.1):
//
//  1. MessageID berilgan — o'sha BITTA xabar bekor qilinadi (source mos
//     kelsa); Dest NULL bo'lishi mumkin.
//  2. MessageID bo'sh (NULL) — Source + Dest (+ ServiceType berilgan bo'lsa)
//     mos kelgan BARCHA kutayotgan xabarlar bekor qilinadi; guruh rejimida
//     Dest majburiy.
//
// Encoder shu ikki rejimdan biriga tushishni tekshiradi.
type CancelSM struct {
	ServiceType string
	MessageID   string
	Source      Address
	Dest        Address
}

// Encode to'liq wire frame yasaydi.
func (c CancelSM) Encode(seq uint32) ([]byte, error) {
	if c.MessageID == "" && c.Dest.Addr == "" {
		return nil, fmt.Errorf("pdu: cancel_sm — message_id ham, destination ham bo'sh: hech qaysi rejimga tushmaydi (§4.9.1)")
	}
	var b bytes.Buffer
	if err := writeCString(&b, c.ServiceType, maxServiceType, "service_type"); err != nil {
		return nil, err
	}
	if err := writeCString(&b, c.MessageID, maxMessageID, "message_id"); err != nil {
		return nil, err
	}
	if err := writeAddress(&b, c.Source, "source_addr"); err != nil {
		return nil, err
	}
	if err := writeAddress(&b, c.Dest, "destination_addr"); err != nil {
		return nil, err
	}
	return encodePDU(CmdCancelSM, 0, seq, b.Bytes()), nil
}

// DecodeCancelSM to'liq frame'dan CancelSM o'qiydi. Decoder tolerant —
// rejim tekshiruvi server biznes-mantig'ining ishi (ESME_RCANCELFAIL).
func DecodeCancelSM(frame []byte) (CancelSM, Header, error) {
	h, err := DecodeHeader(frame)
	if err != nil {
		return CancelSM{}, Header{}, err
	}
	if h.ID != CmdCancelSM {
		return CancelSM{}, h, fmt.Errorf("pdu: %s cancel_sm emas", h.ID)
	}
	r := bytes.NewReader(frame[HeaderSize:])
	var c CancelSM
	if c.ServiceType, err = readCString(r, maxServiceType, "service_type"); err != nil {
		return c, h, err
	}
	if c.MessageID, err = readCString(r, maxMessageID, "message_id"); err != nil {
		return c, h, err
	}
	if c.Source, err = readAddress(r, "source_addr"); err != nil {
		return c, h, err
	}
	if c.Dest, err = readAddress(r, "destination_addr"); err != nil {
		return c, h, err
	}
	return c, h, nil
}

// CancelSMResp — cancel_sm_resp (§4.9.2): body YO'Q, faqat header.
// Muvaffaqiyatsizlik command_status'da (ESME_RCANCELFAIL 0x11).
type CancelSMResp struct {
	Status uint32
}

// Encode to'liq wire frame yasaydi.
func (r CancelSMResp) Encode(seq uint32) []byte {
	return encodePDU(CmdCancelSMResp, r.Status, seq, nil)
}

// DecodeCancelSMResp to'liq frame'dan resp o'qiydi.
func DecodeCancelSMResp(frame []byte) (CancelSMResp, Header, error) {
	h, err := DecodeHeader(frame)
	if err != nil {
		return CancelSMResp{}, Header{}, err
	}
	if h.ID != CmdCancelSMResp {
		return CancelSMResp{}, h, fmt.Errorf("pdu: %s cancel_sm_resp emas", h.ID)
	}
	return CancelSMResp{Status: h.Status}, h, nil
}
