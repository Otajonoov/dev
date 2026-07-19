package pdu

import (
	"bytes"
	"errors"
	"fmt"

	"smpp/tlv"
)

// Bu fayl 10-bob yakuni: barcha 15 PDU turini bitta kirish nuqtasidan
// taniydigan umumiy dispatcher. Session engine (12-bob) simdan frame o'qib,
// shu Decode'ga beradi va konkret turga type switch qiladi.

// PDU — Decode qaytaradigan umumiy interfeys: har konkret PDU turi o'zining
// command_id'sini biladi. Header ma'lumotlari (status, sequence) Decode'ning
// ikkinchi qaytish qiymatida.
type PDU interface {
	Cmd() CommandID
}

// ErrUnknownCommandID — frame'dagi command_id Table 5-1'da yo'q. To'g'ri
// reaksiya (§4.3): generic_nack + ESME_RINVCMDID (0x03). Sentinel xato —
// caller errors.Is bilan ushlab nack yuboradi.
var ErrUnknownCommandID = errors.New("pdu: notanish command_id")

// Header-only PDU'larning dispatcher uchun turlari. Encode tomoni simple.go'da
// (EncodeEnquireLink va h.k.) — bu struct'lar decode natijasini type switch'da
// farqlash uchun.

// EnquireLink — enquire_link (§4.11): body'siz "tirikmisan" so'rovi.
type EnquireLink struct{}

// EnquireLinkResp — enquire_link_resp (§4.11.2).
type EnquireLinkResp struct{}

// Unbind — unbind (§4.2): sessiyani yopish so'rovi.
type Unbind struct{}

// UnbindResp — unbind_resp (§4.2.2).
type UnbindResp struct{ Status uint32 }

// GenericNack — generic_nack (§4.3): faqat header'li salbiy javob.
// Status'da sabab (RINVCMDLEN/RINVCMDID), Sequence — original PDU'niki
// yoki decode bo'lmagan bo'lsa 0 (§4.3.1).
type GenericNack struct{ Status uint32 }

// Cmd metodlari — har konkret tur o'z command_id'sini qaytaradi.

func (b Bind) Cmd() CommandID            { return b.Mode }
func (b BindResp) Cmd() CommandID        { return b.Mode }
func (Outbind) Cmd() CommandID           { return CmdOutbind }
func (SubmitSM) Cmd() CommandID          { return CmdSubmitSM }
func (SubmitSMResp) Cmd() CommandID      { return CmdSubmitSMResp }
func (DeliverSM) Cmd() CommandID         { return CmdDeliverSM }
func (DeliverSMResp) Cmd() CommandID     { return CmdDeliverSMResp }
func (DataSM) Cmd() CommandID            { return CmdDataSM }
func (DataSMResp) Cmd() CommandID        { return CmdDataSMResp }
func (QuerySM) Cmd() CommandID           { return CmdQuerySM }
func (QuerySMResp) Cmd() CommandID       { return CmdQuerySMResp }
func (CancelSM) Cmd() CommandID          { return CmdCancelSM }
func (CancelSMResp) Cmd() CommandID      { return CmdCancelSMResp }
func (ReplaceSM) Cmd() CommandID         { return CmdReplaceSM }
func (ReplaceSMResp) Cmd() CommandID     { return CmdReplaceSMResp }
func (SubmitMulti) Cmd() CommandID       { return CmdSubmitMulti }
func (SubmitMultiResp) Cmd() CommandID   { return CmdSubmitMultiResp }
func (AlertNotification) Cmd() CommandID { return CmdAlertNotification }
func (EnquireLink) Cmd() CommandID       { return CmdEnquireLink }
func (EnquireLinkResp) Cmd() CommandID   { return CmdEnquireLinkResp }
func (Unbind) Cmd() CommandID            { return CmdUnbind }
func (UnbindResp) Cmd() CommandID        { return CmdUnbindResp }
func (GenericNack) Cmd() CommandID       { return CmdGenericNack }

// readTLVTail reader'da qolgan BARCHA baytlarni TLV tail sifatida o'qiydi
// (§3.2.4: TLV'lar doim mandatory qismdan keyin, PDU oxirigacha).
func readTLVTail(r *bytes.Reader) ([]tlv.TLV, error) {
	if r.Len() == 0 {
		return nil, nil
	}
	rest := make([]byte, r.Len())
	if _, err := r.Read(rest); err != nil {
		return nil, fmt.Errorf("pdu: TLV tail o'qishda: %w", err)
	}
	return tlv.Decode(rest)
}

// Decode to'liq frame'ni (header + body) konkret PDU turiga aylantiradi —
// command_id bo'yicha dispatch. Notanish command_id → ErrUnknownCommandID
// (caller generic_nack yuboradi); framing xatolari DecodeHeader'dan keladi.
func Decode(frame []byte) (PDU, Header, error) {
	h, err := DecodeHeader(frame)
	if err != nil {
		return nil, Header{}, err
	}
	switch h.ID {
	case CmdBindTransmitter, CmdBindReceiver, CmdBindTransceiver:
		p, h, err := DecodeBind(frame)
		return p, h, err
	case CmdBindTransmitterResp, CmdBindReceiverResp, CmdBindTransceiverResp:
		p, h, err := DecodeBindResp(frame)
		return p, h, err
	case CmdOutbind:
		p, h, err := DecodeOutbind(frame)
		return p, h, err
	case CmdSubmitSM:
		p, h, err := DecodeSubmitSM(frame)
		return p, h, err
	case CmdSubmitSMResp:
		p, h, err := DecodeSubmitSMResp(frame)
		return p, h, err
	case CmdDeliverSM:
		p, h, err := DecodeDeliverSM(frame)
		return p, h, err
	case CmdDeliverSMResp:
		p, h, err := DecodeDeliverSMResp(frame)
		return p, h, err
	case CmdDataSM:
		p, h, err := DecodeDataSM(frame)
		return p, h, err
	case CmdDataSMResp:
		p, h, err := DecodeDataSMResp(frame)
		return p, h, err
	case CmdQuerySM:
		p, h, err := DecodeQuerySM(frame)
		return p, h, err
	case CmdQuerySMResp:
		p, h, err := DecodeQuerySMResp(frame)
		return p, h, err
	case CmdCancelSM:
		p, h, err := DecodeCancelSM(frame)
		return p, h, err
	case CmdCancelSMResp:
		p, h, err := DecodeCancelSMResp(frame)
		return p, h, err
	case CmdReplaceSM:
		p, h, err := DecodeReplaceSM(frame)
		return p, h, err
	case CmdReplaceSMResp:
		p, h, err := DecodeReplaceSMResp(frame)
		return p, h, err
	case CmdSubmitMulti:
		p, h, err := DecodeSubmitMulti(frame)
		return p, h, err
	case CmdSubmitMultiResp:
		p, h, err := DecodeSubmitMultiResp(frame)
		return p, h, err
	case CmdAlertNotification:
		p, h, err := DecodeAlertNotification(frame)
		return p, h, err
	case CmdEnquireLink:
		return EnquireLink{}, h, nil
	case CmdEnquireLinkResp:
		return EnquireLinkResp{}, h, nil
	case CmdUnbind:
		return Unbind{}, h, nil
	case CmdUnbindResp:
		return UnbindResp{Status: h.Status}, h, nil
	case CmdGenericNack:
		return GenericNack{Status: h.Status}, h, nil
	}
	return nil, h, fmt.Errorf("%w: 0x%08X", ErrUnknownCommandID, uint32(h.ID))
}
