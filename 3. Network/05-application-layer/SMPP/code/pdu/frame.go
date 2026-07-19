package pdu

import (
	"errors"
	"fmt"
	"io"
)

// TCP — stream protokol: unda "xabar chegarasi" tushunchasi yo'q. Bitta Read
// chaqiruvi yarim PDU ham, uch yarim PDU ham qaytarishi mumkin. SMPP framing'i
// command_length'ga tayanadi (v3.4 §3.2.2): avval 4 oktet length o'qiladi,
// keyin qolgan (length − 4) oktet TO'LIQ yig'ilguncha o'qiladi — io.ReadFull.

var (
	// ErrFrameTooShort — command_length < 16: header ham sig'maydigan PDU
	// bo'lmaydi; stream buzilgan (§2.8 bo'yicha bunga generic_nack qaytariladi).
	ErrFrameTooShort = errors.New("pdu: command_length header'dan ham kichik")
	// ErrFrameTooLarge — command_length maxSize'dan katta: buzilgan stream yoki
	// ataylab yuborilgan gigant qiymat; ko'r-ko'rona make([]byte, length)
	// qilishdan oldin kesib tashlanadi (OOM himoyasi).
	ErrFrameTooLarge = errors.New("pdu: command_length ruxsat etilgan max'dan katta")
)

// ReadFrame r'dan bitta to'liq PDU frame o'qib qaytaradi (header bilan birga,
// ya'ni natija uzunligi = command_length). maxSize — qabul qilinadigan eng
// katta PDU (odatda bir necha KB; message_payload max 64K ekanini yodda tuting).
//
// Xato semantikasi: frame'lar ORASIDA toza uzilgan stream io.EOF qaytaradi;
// frame O'RTASIDA uzilgani io.ErrUnexpectedEOF.
func ReadFrame(r io.Reader, maxSize uint32) ([]byte, error) {
	var lenBuf [4]byte
	if _, err := io.ReadFull(r, lenBuf[:]); err != nil {
		if errors.Is(err, io.ErrUnexpectedEOF) {
			// length field'ining o'zi ham to'liq kelmagan.
			return nil, fmt.Errorf("pdu: command_length o'qishda stream uzildi: %w", err)
		}
		return nil, err
	}
	length := getUint32(lenBuf[:])
	if length < HeaderSize {
		return nil, fmt.Errorf("%w: command_length=%d, kamida %d bo'lishi kerak", ErrFrameTooShort, length, HeaderSize)
	}
	if length > maxSize {
		return nil, fmt.Errorf("%w: command_length=%d, max %d", ErrFrameTooLarge, length, maxSize)
	}
	frame := make([]byte, length)
	copy(frame, lenBuf[:])
	if _, err := io.ReadFull(r, frame[4:]); err != nil {
		if errors.Is(err, io.EOF) {
			err = io.ErrUnexpectedEOF
		}
		return nil, fmt.Errorf("pdu: frame body o'qishda stream uzildi (%d oktetdan %d kutilgan edi): %w", 4, length, err)
	}
	return frame, nil
}
