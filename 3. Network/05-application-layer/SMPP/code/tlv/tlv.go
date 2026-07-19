// Package tlv SMPP v3.4 optional parameter'lari — TLV (Tag-Length-Value) —
// codec'ini beradi (v3.4 §3.2.4).
//
// Format (§3.2.4.1, Table 3-4): Tag — 2 oktet Integer; Length — 2 oktet
// Integer, FAQAT Value uzunligi (Tag+Length kirmaydi); Value — o'zgaruvchan.
// TLV'lar doim PDU'ning mandatory qismidan KEYIN, o'zaro ixtiyoriy tartibda
// keladi (§3.2.4).
//
// Dizayn: TLV to'plami map emas, SLICE sifatida yuritiladi — simdagi tartib
// saqlanadi va bitta tag'ning takrorlanishiga ruxsat qoladi (callback_num
// PDU'da bir necha marta kelishi mumkin — v3.4 §4.4.1 izohi).
package tlv

import (
	"bytes"
	"errors"
	"fmt"
)

// Tag — TLV'ning 2 oktetlik identifikatori (v3.4 §3.2.4.1).
// Bloklari (§5.3.2, Table 5-7): 0x0001–0x00FF, 0x0200–0x05FF va
// 0x1200–0x13FF — SMPP-defined; 0x1400–0x3FFF — SMSC vendor-specific;
// qolgani reserved.
type Tag uint16

// v3.4 Table 5-7 (§5.3.2) bo'yicha SMPP-defined tag'lar.
const (
	DestAddrSubunit          Tag = 0x0005
	DestNetworkType          Tag = 0x0006
	DestBearerType           Tag = 0x0007
	DestTelematicsID         Tag = 0x0008
	SourceAddrSubunit        Tag = 0x000D
	SourceNetworkType        Tag = 0x000E
	SourceBearerType         Tag = 0x000F
	SourceTelematicsID       Tag = 0x0010
	QosTimeToLive            Tag = 0x0017
	PayloadType              Tag = 0x0019
	AdditionalStatusInfoText Tag = 0x001D
	ReceiptedMessageID       Tag = 0x001E
	MsMsgWaitFacilities      Tag = 0x0030
	PrivacyIndicator         Tag = 0x0201
	SourceSubaddress         Tag = 0x0202
	DestSubaddress           Tag = 0x0203
	UserMessageReference     Tag = 0x0204
	UserResponseCode         Tag = 0x0205
	SourcePort               Tag = 0x020A
	DestinationPort          Tag = 0x020B
	SarMsgRefNum             Tag = 0x020C
	LanguageIndicator        Tag = 0x020D
	SarTotalSegments         Tag = 0x020E
	SarSegmentSeqnum         Tag = 0x020F
	ScInterfaceVersion       Tag = 0x0210
	CallbackNumPresInd       Tag = 0x0302
	CallbackNumAtag          Tag = 0x0303
	NumberOfMessages         Tag = 0x0304
	CallbackNum              Tag = 0x0381
	DpfResult                Tag = 0x0420
	SetDpf                   Tag = 0x0421
	MsAvailabilityStatus     Tag = 0x0422
	NetworkErrorCode         Tag = 0x0423
	MessagePayload           Tag = 0x0424
	DeliveryFailureReason    Tag = 0x0425
	MoreMessagesToSend       Tag = 0x0426
	MessageState             Tag = 0x0427
	UssdServiceOp            Tag = 0x0501
	DisplayTime              Tag = 0x1201
	SmsSignal                Tag = 0x1203
	MsValidity               Tag = 0x1204
	AlertOnMessageDelivery   Tag = 0x130C
	ItsReplyType             Tag = 0x1380
	ItsSessionInfo           Tag = 0x1383
)

// IsVendor tag SMSC vendor-specific blokida (0x1400–0x3FFF) ekanini bildiradi.
func (t Tag) IsVendor() bool { return t >= 0x1400 && t <= 0x3FFF }

var tagNames = map[Tag]string{
	DestAddrSubunit:          "dest_addr_subunit",
	DestNetworkType:          "dest_network_type",
	DestBearerType:           "dest_bearer_type",
	DestTelematicsID:         "dest_telematics_id",
	SourceAddrSubunit:        "source_addr_subunit",
	SourceNetworkType:        "source_network_type",
	SourceBearerType:         "source_bearer_type",
	SourceTelematicsID:       "source_telematics_id",
	QosTimeToLive:            "qos_time_to_live",
	PayloadType:              "payload_type",
	AdditionalStatusInfoText: "additional_status_info_text",
	ReceiptedMessageID:       "receipted_message_id",
	MsMsgWaitFacilities:      "ms_msg_wait_facilities",
	PrivacyIndicator:         "privacy_indicator",
	SourceSubaddress:         "source_subaddress",
	DestSubaddress:           "dest_subaddress",
	UserMessageReference:     "user_message_reference",
	UserResponseCode:         "user_response_code",
	SourcePort:               "source_port",
	DestinationPort:          "destination_port",
	SarMsgRefNum:             "sar_msg_ref_num",
	LanguageIndicator:        "language_indicator",
	SarTotalSegments:         "sar_total_segments",
	SarSegmentSeqnum:         "sar_segment_seqnum",
	ScInterfaceVersion:       "sc_interface_version",
	CallbackNumPresInd:       "callback_num_pres_ind",
	CallbackNumAtag:          "callback_num_atag",
	NumberOfMessages:         "number_of_messages",
	CallbackNum:              "callback_num",
	DpfResult:                "dpf_result",
	SetDpf:                   "set_dpf",
	MsAvailabilityStatus:     "ms_availability_status",
	NetworkErrorCode:         "network_error_code",
	MessagePayload:           "message_payload",
	DeliveryFailureReason:    "delivery_failure_reason",
	MoreMessagesToSend:       "more_messages_to_send",
	MessageState:             "message_state",
	UssdServiceOp:            "ussd_service_op",
	DisplayTime:              "display_time",
	SmsSignal:                "sms_signal",
	MsValidity:               "ms_validity",
	AlertOnMessageDelivery:   "alert_on_message_delivery",
	ItsReplyType:             "its_reply_type",
	ItsSessionInfo:           "its_session_info",
}

// String Table 5-7'dagi rasmiy nomni qaytaradi; vendor blok uchun "vendor(...)",
// notanish tag uchun hex ko'rinish.
func (t Tag) String() string {
	if name, ok := tagNames[t]; ok {
		return name
	}
	if t.IsVendor() {
		return fmt.Sprintf("vendor(0x%04X)", uint16(t))
	}
	return fmt.Sprintf("unknown(0x%04X)", uint16(t))
}

// TLV — bitta optional parameter. Value nil yoki bo'sh bo'lishi mumkin
// (zero-length TLV — masalan alert_on_message_delivery, §5.3.2.41).
type TLV struct {
	Tag   Tag
	Value []byte
}

// ErrTruncated — TLV tail o'rtasida uzilgan: Tag/Length'ning o'zi to'liq emas
// yoki Length va'da qilgan Value oxirigacha kelmagan.
var ErrTruncated = errors.New("tlv: TLV tail o'rtasida uzilgan")

// Encode tlv'larni b'ga simdagi tartibda yozadi. Length — faqat Value uzunligi.
func Encode(b *bytes.Buffer, tlvs []TLV) error {
	for _, t := range tlvs {
		if len(t.Value) > 0xFFFF {
			return fmt.Errorf("tlv: %s value %d oktet — 2 oktetlik Length'ga sig'maydi", t.Tag, len(t.Value))
		}
		b.WriteByte(byte(t.Tag >> 8))
		b.WriteByte(byte(t.Tag))
		b.WriteByte(byte(len(t.Value) >> 8))
		b.WriteByte(byte(len(t.Value)))
		b.Write(t.Value)
	}
	return nil
}

// Decode PDU'ning TLV tail'ini (mandatory qismdan keyingi BARCHA baytlar)
// parse qiladi. Notanish tag'lar ham SAQLANADI (forward compatibility, §3.3):
// tashlab yuborish caller ixtiyorida, codec ma'lumot yo'qotmaydi.
func Decode(data []byte) ([]TLV, error) {
	var tlvs []TLV
	for off := 0; off < len(data); {
		if len(data)-off < 4 {
			return nil, fmt.Errorf("%w: %d oktet qoldi, Tag+Length uchun 4 kerak", ErrTruncated, len(data)-off)
		}
		tag := Tag(uint16(data[off])<<8 | uint16(data[off+1]))
		length := int(uint16(data[off+2])<<8 | uint16(data[off+3]))
		off += 4
		if len(data)-off < length {
			return nil, fmt.Errorf("%w: %s Length=%d, qolgan baytlar %d", ErrTruncated, tag, length, len(data)-off)
		}
		// Value nusxalanadi: frame buffer'i qayta ishlatilsa TLV yashab qolsin.
		value := make([]byte, length)
		copy(value, data[off:off+length])
		off += length
		tlvs = append(tlvs, TLV{Tag: tag, Value: value})
	}
	return tlvs, nil
}

// Find birinchi mos tag'ni qaytaradi. Takrorlanadigan tag'lar (callback_num)
// uchun to'plamni to'g'ridan-to'g'ri aylanib chiqish kerak.
func Find(tlvs []TLV, tag Tag) (TLV, bool) {
	for _, t := range tlvs {
		if t.Tag == tag {
			return t, true
		}
	}
	return TLV{}, false
}
