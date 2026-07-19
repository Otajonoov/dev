package pdu

import "fmt"

// HeaderSize — har PDU boshidagi majburiy header hajmi oktetlarda (v3.4 §3.2).
const HeaderSize = 16

// Sequence_number diapazoni (v3.4 §5.1.4). 0 qiymat faqat generic_nack'ning
// "originalni decode qilib bo'lmadi" holatida uchraydi (§4.3.1; qarang 11-bob).
const (
	MinSequence uint32 = 0x00000001
	MaxSequence uint32 = 0x7FFFFFFF
)

// Header — 16 oktetlik PDU header'i: to'rtta 4-oktetlik big-endian Integer
// (v3.4 §3.2, Table 3-2).
type Header struct {
	Length   uint32    // command_length: BUTUN PDU, shu field'ning o'zi ham kiradi (§5.1.1)
	ID       CommandID // command_id: PDU turi (§5.1.2)
	Status   uint32    // command_status: faqat response'da ma'noli; request'da 0 SHART (§5.1.3)
	Sequence uint32    // sequence_number: request↔response korrelyatsiyasi (§5.1.4)
}

// EncodeHeader header'ni 16 oktetlik big-endian ko'rinishga o'tkazadi.
func EncodeHeader(h Header) [HeaderSize]byte {
	var b [HeaderSize]byte
	putUint32(b[0:4], h.Length)
	putUint32(b[4:8], uint32(h.ID))
	putUint32(b[8:12], h.Status)
	putUint32(b[12:16], h.Sequence)
	return b
}

// DecodeHeader data'ning birinchi 16 oktetidan header o'qiydi.
// Length'ning o'zi bu yerda TEKSHIRILMAYDI — framing darajasidagi
// validatsiya ReadFrame'da (frame.go).
func DecodeHeader(data []byte) (Header, error) {
	if len(data) < HeaderSize {
		return Header{}, fmt.Errorf("pdu: header uchun %d oktet yetarli emas (kamida %d kerak)", len(data), HeaderSize)
	}
	return Header{
		Length:   getUint32(data[0:4]),
		ID:       CommandID(getUint32(data[4:8])),
		Status:   getUint32(data[8:12]),
		Sequence: getUint32(data[12:16]),
	}, nil
}

func putUint32(b []byte, v uint32) {
	b[0] = byte(v >> 24)
	b[1] = byte(v >> 16)
	b[2] = byte(v >> 8)
	b[3] = byte(v)
}

func getUint32(b []byte) uint32 {
	return uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])
}
