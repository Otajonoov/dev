package pdu

import (
	"bytes"
	"fmt"

	"smpp/tlv"
)

// DataSM — data_sm PDU (§4.7): submit_sm/deliver_sm'ga muqobil, IKKI
// yo'nalishli (ESME↔SMSC), WAP kabi interaktiv application'lar uchun
// kiritilgan. Farqlari:
//   - mandatory body QISQA: short_message/sm_length YO'Q — matn faqat
//     message_payload TLV'da (0x0424);
//   - manzillar max 65 oktet (submit_sm'dagi 21 emas);
//   - protocol_id, priority, schedule/validity, replace_if_present,
//     sm_default_msg_id yo'q — kerak bo'lsa TLV ekvivalentlari ishlatiladi;
//   - transaction messaging mode'ni QO'LLAYDI (§2.10.3) — data_sm_resp
//     end-to-end natijani qaytarishi mumkin.
type DataSM struct {
	ServiceType        string
	Source             Address // max 65 (§4.7.1)
	Dest               Address // max 65
	EsmClass           EsmClass
	RegisteredDelivery RegisteredDelivery
	DataCoding         uint8
	TLVs               []tlv.TLV // matn message_payload TLV'da keladi
}

// Encode to'liq wire frame yasaydi (request → command_status=0).
func (d DataSM) Encode(seq uint32) ([]byte, error) {
	var b bytes.Buffer
	if err := writeCString(&b, d.ServiceType, maxServiceType, "service_type"); err != nil {
		return nil, err
	}
	if err := writeAddressLong(&b, d.Source, "source_addr"); err != nil {
		return nil, err
	}
	if err := writeAddressLong(&b, d.Dest, "destination_addr"); err != nil {
		return nil, err
	}
	writeUint8(&b, uint8(d.EsmClass))
	writeUint8(&b, uint8(d.RegisteredDelivery))
	writeUint8(&b, d.DataCoding)
	if err := tlv.Encode(&b, d.TLVs); err != nil {
		return nil, err
	}
	return encodePDU(CmdDataSM, 0, seq, b.Bytes()), nil
}

// DecodeDataSM to'liq frame'dan DataSM o'qiydi.
func DecodeDataSM(frame []byte) (DataSM, Header, error) {
	h, err := DecodeHeader(frame)
	if err != nil {
		return DataSM{}, Header{}, err
	}
	if h.ID != CmdDataSM {
		return DataSM{}, h, fmt.Errorf("pdu: %s data_sm emas", h.ID)
	}
	r := bytes.NewReader(frame[HeaderSize:])
	var d DataSM
	if d.ServiceType, err = readCString(r, maxServiceType, "service_type"); err != nil {
		return d, h, err
	}
	if d.Source, err = readAddressLong(r, "source_addr"); err != nil {
		return d, h, err
	}
	if d.Dest, err = readAddressLong(r, "destination_addr"); err != nil {
		return d, h, err
	}
	var esm, reg uint8
	if esm, err = readUint8(r, "esm_class"); err != nil {
		return d, h, err
	}
	d.EsmClass = EsmClass(esm)
	if reg, err = readUint8(r, "registered_delivery"); err != nil {
		return d, h, err
	}
	d.RegisteredDelivery = RegisteredDelivery(reg)
	if d.DataCoding, err = readUint8(r, "data_coding"); err != nil {
		return d, h, err
	}
	if d.TLVs, err = readTLVTail(r); err != nil {
		return d, h, err
	}
	return d, h, nil
}

// DataSMResp — data_sm_resp (§4.7.2): message_id + ixtiyoriy TLV'lar.
// delivery_failure_reason / network_error_code / dpf_result TLV'lari faqat
// transaction mode'da ma'noli (§4.7.2 izohlari); additional_status_info_text —
// ASCII izoh. Bular submit_sm_resp'da YO'Q — data_sm'ning "end-to-end natija"
// imkoniyati aynan shu TLV'larda.
type DataSMResp struct {
	Status    uint32
	MessageID string
	TLVs      []tlv.TLV
}

// Encode to'liq wire frame yasaydi. Status != 0 → body YO'Q (§4.7.2'da ham
// submit_sm_resp'dagi qoida amal qiladi).
func (r DataSMResp) Encode(seq uint32) ([]byte, error) {
	if r.Status != 0 {
		return encodePDU(CmdDataSMResp, r.Status, seq, nil), nil
	}
	var b bytes.Buffer
	if err := writeCString(&b, r.MessageID, maxMessageID, "message_id"); err != nil {
		return nil, err
	}
	if err := tlv.Encode(&b, r.TLVs); err != nil {
		return nil, err
	}
	return encodePDU(CmdDataSMResp, 0, seq, b.Bytes()), nil
}

// DecodeDataSMResp to'liq frame'dan resp o'qiydi.
func DecodeDataSMResp(frame []byte) (DataSMResp, Header, error) {
	h, err := DecodeHeader(frame)
	if err != nil {
		return DataSMResp{}, Header{}, err
	}
	if h.ID != CmdDataSMResp {
		return DataSMResp{}, h, fmt.Errorf("pdu: %s data_sm_resp emas", h.ID)
	}
	resp := DataSMResp{Status: h.Status}
	if h.Status != 0 {
		return resp, h, nil
	}
	r := bytes.NewReader(frame[HeaderSize:])
	if resp.MessageID, err = readCString(r, maxMessageID, "message_id"); err != nil {
		return DataSMResp{}, h, err
	}
	if resp.TLVs, err = readTLVTail(r); err != nil {
		return DataSMResp{}, h, err
	}
	return resp, h, nil
}
