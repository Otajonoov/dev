package pdu

import (
	"bytes"
	"fmt"
)

// QuerySM — query_sm PDU (§4.8): ilgari yuborilgan xabar holatini so'rash.
// Matching qoidasi: message_id + source manzil ORIGINAL submit'dagi bilan
// AYNAN mos kelishi shart — original'da source NULL bo'lgan bo'lsa bu yerda
// ham NULL (§4.8). Aks holda SMSC xabarni "topmaydi" (ESME_RINVMSGID yoki
// ESME_RQUERYFAIL).
type QuerySM struct {
	MessageID string
	Source    Address
}

// Encode to'liq wire frame yasaydi.
func (q QuerySM) Encode(seq uint32) ([]byte, error) {
	var b bytes.Buffer
	if err := writeCString(&b, q.MessageID, maxMessageID, "message_id"); err != nil {
		return nil, err
	}
	if err := writeAddress(&b, q.Source, "source_addr"); err != nil {
		return nil, err
	}
	return encodePDU(CmdQuerySM, 0, seq, b.Bytes()), nil
}

// DecodeQuerySM to'liq frame'dan QuerySM o'qiydi.
func DecodeQuerySM(frame []byte) (QuerySM, Header, error) {
	h, err := DecodeHeader(frame)
	if err != nil {
		return QuerySM{}, Header{}, err
	}
	if h.ID != CmdQuerySM {
		return QuerySM{}, h, fmt.Errorf("pdu: %s query_sm emas", h.ID)
	}
	r := bytes.NewReader(frame[HeaderSize:])
	var q QuerySM
	if q.MessageID, err = readCString(r, maxMessageID, "message_id"); err != nil {
		return q, h, err
	}
	if q.Source, err = readAddress(r, "source_addr"); err != nil {
		return q, h, err
	}
	return q, h, nil
}

// QuerySMResp — query_sm_resp (§4.8.2): xabarning joriy holati.
//
// FinalDate — "1 or 17" vaqt field'i (§7.1.1): xabar final holatga
// yetmagan bo'lsa NULL (bo'sh string). MessageState — §5.2.28 qiymatlari
// (1=ENROUTE..8=REJECTED; dlr package'dagi MessageState bilan bir fazo).
// ErrorCode — network-specific kod (DLR'dagi err: bilan bir fazo,
// command_status EMAS).
type QuerySMResp struct {
	Status       uint32
	MessageID    string
	FinalDate    string
	MessageState uint8
	ErrorCode    uint8
}

// Encode to'liq wire frame yasaydi. Status != 0 → body YO'Q.
func (r QuerySMResp) Encode(seq uint32) ([]byte, error) {
	if r.Status != 0 {
		return encodePDU(CmdQuerySMResp, r.Status, seq, nil), nil
	}
	var b bytes.Buffer
	if err := writeCString(&b, r.MessageID, maxMessageID, "message_id"); err != nil {
		return nil, err
	}
	if err := writeTimeField(&b, r.FinalDate, "final_date"); err != nil {
		return nil, err
	}
	writeUint8(&b, r.MessageState)
	writeUint8(&b, r.ErrorCode)
	return encodePDU(CmdQuerySMResp, 0, seq, b.Bytes()), nil
}

// DecodeQuerySMResp to'liq frame'dan resp o'qiydi.
func DecodeQuerySMResp(frame []byte) (QuerySMResp, Header, error) {
	h, err := DecodeHeader(frame)
	if err != nil {
		return QuerySMResp{}, Header{}, err
	}
	if h.ID != CmdQuerySMResp {
		return QuerySMResp{}, h, fmt.Errorf("pdu: %s query_sm_resp emas", h.ID)
	}
	resp := QuerySMResp{Status: h.Status}
	if h.Status != 0 {
		return resp, h, nil
	}
	r := bytes.NewReader(frame[HeaderSize:])
	if resp.MessageID, err = readCString(r, maxMessageID, "message_id"); err != nil {
		return QuerySMResp{}, h, err
	}
	if resp.FinalDate, err = readCString(r, maxTimeField, "final_date"); err != nil {
		return QuerySMResp{}, h, err
	}
	if resp.MessageState, err = readUint8(r, "message_state"); err != nil {
		return QuerySMResp{}, h, err
	}
	if resp.ErrorCode, err = readUint8(r, "error_code"); err != nil {
		return QuerySMResp{}, h, err
	}
	return resp, h, nil
}
