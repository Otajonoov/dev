package pdu

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

// Bu fayl v3.4 §3.1 data type'lari uchun past darajali encode/decode
// helper'larini beradi. Barcha Integer'lar big-endian (v3.4 §3.1: "MSB first").
//
// C-Octet String qoidalari (v3.4 §3.1, Table 3-1):
//   - ASCII belgilar ketma-ketligi + NULL (0x00) terminator;
//   - bo'sh string = yagona 0x00 okteti;
//   - spec'dagi "Max N" o'lchamlar NULL terminator'ni O'Z ICHIGA OLADI
//     (§3.1 note iii: 8 belgili string 9 oktetda kodlanadi).

// ErrNoTerminator — C-Octet String maksimal uzunlikkacha o'qildi,
// lekin NULL terminator topilmadi.
var ErrNoTerminator = errors.New("pdu: C-Octet String NULL terminator'siz")

// writeUint8/16/32 — big-endian Integer'lar (v3.4 §3.1).
// bytes.Buffer'ga yozish hech qachon xato qaytarmaydi.

func writeUint8(b *bytes.Buffer, v uint8) {
	b.WriteByte(v)
}

func writeUint16(b *bytes.Buffer, v uint16) {
	b.WriteByte(byte(v >> 8))
	b.WriteByte(byte(v))
}

func writeUint32(b *bytes.Buffer, v uint32) {
	b.WriteByte(byte(v >> 24))
	b.WriteByte(byte(v >> 16))
	b.WriteByte(byte(v >> 8))
	b.WriteByte(byte(v))
}

func readUint8(r *bytes.Reader, field string) (uint8, error) {
	v, err := r.ReadByte()
	if err != nil {
		return 0, fmt.Errorf("pdu: %s field'ini o'qishda stream tugadi", field)
	}
	return v, nil
}

func readUint16(r *bytes.Reader, field string) (uint16, error) {
	var buf [2]byte
	if _, err := io.ReadFull(r, buf[:]); err != nil {
		return 0, fmt.Errorf("pdu: %s field'ini o'qishda stream tugadi", field)
	}
	return uint16(buf[0])<<8 | uint16(buf[1]), nil
}

func readUint32(r *bytes.Reader, field string) (uint32, error) {
	var buf [4]byte
	if _, err := io.ReadFull(r, buf[:]); err != nil {
		return 0, fmt.Errorf("pdu: %s field'ini o'qishda stream tugadi", field)
	}
	return uint32(buf[0])<<24 | uint32(buf[1])<<16 | uint32(buf[2])<<8 | uint32(buf[3]), nil
}

// writeCString s'ni NULL terminator bilan yozadi. max — spec'dagi field o'lchami
// (NULL terminator'ni O'Z ICHIGA olgan holda, §3.1 note iii). Ichki 0x00 baytga
// ruxsat yo'q — u terminator bilan aralashib butun PDU'ni buzadi.
func writeCString(b *bytes.Buffer, s string, max int, field string) error {
	if len(s)+1 > max {
		return fmt.Errorf("pdu: %s uzunligi %d oktet, max %d (NULL bilan)", field, len(s)+1, max)
	}
	for i := 0; i < len(s); i++ {
		if s[i] == 0x00 {
			return fmt.Errorf("pdu: %s ichida NULL bayt (indeks %d)", field, i)
		}
	}
	b.WriteString(s)
	b.WriteByte(0x00)
	return nil
}

// readCString NULL terminator'gacha o'qiydi (terminator iste'mol qilinadi,
// natijaga kirmaydi). max — spec'dagi field o'lchami (NULL bilan); shu
// chegaragacha terminator topilmasa ErrNoTerminator qaytadi.
func readCString(r *bytes.Reader, max int, field string) (string, error) {
	var sb []byte
	for i := 0; i < max; i++ {
		c, err := r.ReadByte()
		if err != nil {
			return "", fmt.Errorf("pdu: %s field'ini o'qishda stream tugadi", field)
		}
		if c == 0x00 {
			return string(sb), nil
		}
		sb = append(sb, c)
	}
	return "", fmt.Errorf("%w: %s (max %d oktet)", ErrNoTerminator, field, max)
}
