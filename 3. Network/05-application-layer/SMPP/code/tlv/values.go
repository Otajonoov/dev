package tlv

import (
	"bytes"
	"fmt"
)

// Bu fayl Value baytlarini tipli qiymatlarga aylantiruvchi helper'lar va
// encode tomonning qulay konstruktorlarini beradi.

// Uint8Value 1 oktetlik Integer value (masalan sc_interface_version, §5.3.2.25).
func (t TLV) Uint8Value() (uint8, error) {
	if len(t.Value) != 1 {
		return 0, fmt.Errorf("tlv: %s value %d oktet, 1 kutilgan", t.Tag, len(t.Value))
	}
	return t.Value[0], nil
}

// Uint16Value 2 oktetlik big-endian Integer value (masalan sar_msg_ref_num, §5.3.2.22).
func (t TLV) Uint16Value() (uint16, error) {
	if len(t.Value) != 2 {
		return 0, fmt.Errorf("tlv: %s value %d oktet, 2 kutilgan", t.Tag, len(t.Value))
	}
	return uint16(t.Value[0])<<8 | uint16(t.Value[1]), nil
}

// Uint32Value 4 oktetlik big-endian Integer value (masalan qos_time_to_live, §5.3.2.9).
func (t TLV) Uint32Value() (uint32, error) {
	if len(t.Value) != 4 {
		return 0, fmt.Errorf("tlv: %s value %d oktet, 4 kutilgan", t.Tag, len(t.Value))
	}
	return uint32(t.Value[0])<<24 | uint32(t.Value[1])<<16 | uint32(t.Value[2])<<8 | uint32(t.Value[3]), nil
}

// CStringValue C-Octet String value (masalan receipted_message_id, §5.3.2.12).
// Tolerant o'qish: spec NULL terminator talab qiladi, lekin ayrim SMSC'lar
// TLV value'da terminatorsiz yuboradi — bor bo'lsa olib tashlanadi, yo'q
// bo'lsa ham xato emas. Ichki NULL esa har doim xato.
func (t TLV) CStringValue() (string, error) {
	v := t.Value
	if i := bytes.IndexByte(v, 0x00); i >= 0 {
		if i != len(v)-1 {
			return "", fmt.Errorf("tlv: %s value ichida NULL bayt (indeks %d)", t.Tag, i)
		}
		v = v[:i]
	}
	return string(v), nil
}

// NetworkError — network_error_code (0x0423) TLV'sining 3 oktetlik strukturali
// value'si (§5.3.2.31): [network type: 1 oktet] + [error code: 2 oktet].
type NetworkError struct {
	Type uint8  // 1=ANSI-136, 2=IS-95, 3=GSM
	Code uint16 // tarmoqqa xos xato kodi
}

// NetworkError network_error_code value'sini parse qiladi.
func (t TLV) NetworkError() (NetworkError, error) {
	if len(t.Value) != 3 {
		return NetworkError{}, fmt.Errorf("tlv: %s value %d oktet, 3 kutilgan", t.Tag, len(t.Value))
	}
	return NetworkError{
		Type: t.Value[0],
		Code: uint16(t.Value[1])<<8 | uint16(t.Value[2]),
	}, nil
}

// Encode tomon konstruktorlari.

// U8 1 oktetlik Integer value'li TLV yasaydi.
func U8(tag Tag, v uint8) TLV { return TLV{Tag: tag, Value: []byte{v}} }

// U16 2 oktetlik big-endian Integer value'li TLV yasaydi.
func U16(tag Tag, v uint16) TLV { return TLV{Tag: tag, Value: []byte{byte(v >> 8), byte(v)}} }

// U32 4 oktetlik big-endian Integer value'li TLV yasaydi.
func U32(tag Tag, v uint32) TLV {
	return TLV{Tag: tag, Value: []byte{byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}}
}

// CString NULL-terminatorli C-Octet String value'li TLV yasaydi (spec shakli).
func CString(tag Tag, s string) TLV {
	v := make([]byte, 0, len(s)+1)
	v = append(v, s...)
	v = append(v, 0x00)
	return TLV{Tag: tag, Value: v}
}

// Empty zero-length value'li TLV yasaydi (alert_on_message_delivery, §5.3.2.41).
func Empty(tag Tag) TLV { return TLV{Tag: tag, Value: nil} }
