package pdu

import (
	"bytes"
	"fmt"

	"smpp/tlv"
)

// AlertNotification — alert_notification PDU (§4.12): SMSC→ESME signal —
// ilgari data_sm'da set_dpf (delivery pending flag) bilan belgilangan
// abonent (MS) yana AVAILABLE bo'ldi. Tipik ishlatilishi: WAP push retry
// trigger'i. Xabar tanasi yo'q — bu faqat "endi urinsang bo'ladi" signali.
//
// Ikki o'ziga xoslik: (1) RESPONSE YO'Q (§4.12.1 — Table 5-1'da resp qiymati
// Reserved) — dispatcher'da unga javob kutilmaydi va yuborilmaydi;
// (2) manzillar max 65 (data_sm bilan bir xil): Source — available bo'lgan
// MS, ESMEAddr — dpf o'rnatgan ESME'ning manzili.
type AlertNotification struct {
	Source   Address
	ESMEAddr Address
	TLVs     []tlv.TLV // ms_availability_status (0x0422) kelishi mumkin
}

// Encode to'liq wire frame yasaydi (server tomon — mock SMSC ishlatadi).
func (a AlertNotification) Encode(seq uint32) ([]byte, error) {
	var b bytes.Buffer
	if err := writeAddressLong(&b, a.Source, "source_addr"); err != nil {
		return nil, err
	}
	if err := writeAddressLong(&b, a.ESMEAddr, "esme_addr"); err != nil {
		return nil, err
	}
	if err := tlv.Encode(&b, a.TLVs); err != nil {
		return nil, err
	}
	return encodePDU(CmdAlertNotification, 0, seq, b.Bytes()), nil
}

// DecodeAlertNotification to'liq frame'dan AlertNotification o'qiydi.
func DecodeAlertNotification(frame []byte) (AlertNotification, Header, error) {
	h, err := DecodeHeader(frame)
	if err != nil {
		return AlertNotification{}, Header{}, err
	}
	if h.ID != CmdAlertNotification {
		return AlertNotification{}, h, fmt.Errorf("pdu: %s alert_notification emas", h.ID)
	}
	r := bytes.NewReader(frame[HeaderSize:])
	var a AlertNotification
	if a.Source, err = readAddressLong(r, "source_addr"); err != nil {
		return a, h, err
	}
	if a.ESMEAddr, err = readAddressLong(r, "esme_addr"); err != nil {
		return a, h, err
	}
	if a.TLVs, err = readTLVTail(r); err != nil {
		return a, h, err
	}
	return a, h, nil
}
