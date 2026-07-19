package pdu

// Bu fayl header-only PDU'larni (body'siz, command_length=16) va PDU yig'ish
// helper'ini beradi: unbind/unbind_resp (§4.2), generic_nack (§4.3),
// enquire_link/enquire_link_resp (§4.11).

// encodePDU header + body'ni yaxlit wire frame'ga yig'adi. Barcha PDU
// encoder'lari shu yo'ldan o'tadi — command_length hisoblash bitta joyda.
func encodePDU(id CommandID, status, seq uint32, body []byte) []byte {
	h := EncodeHeader(Header{
		Length:   uint32(HeaderSize + len(body)),
		ID:       id,
		Status:   status,
		Sequence: seq,
	})
	frame := make([]byte, 0, HeaderSize+len(body))
	frame = append(frame, h[:]...)
	frame = append(frame, body...)
	return frame
}

// EncodeEnquireLink — "tirikmisan?" so'rovi (§4.11). Ikkala tomon ham yuboradi.
func EncodeEnquireLink(seq uint32) []byte { return encodePDU(CmdEnquireLink, 0, seq, nil) }

// EncodeEnquireLinkResp — enquire_link javobi; darhol qaytarilishi shart.
func EncodeEnquireLinkResp(seq uint32) []byte { return encodePDU(CmdEnquireLinkResp, 0, seq, nil) }

// EncodeUnbind — sessiyani yopish so'rovi (§4.2). Ikkala tomon ham yuboradi.
func EncodeUnbind(seq uint32) []byte { return encodePDU(CmdUnbind, 0, seq, nil) }

// EncodeUnbindResp — unbind javobi; bundan keyingina TCP yopiladi.
func EncodeUnbindResp(status, seq uint32) []byte { return encodePDU(CmdUnbindResp, status, seq, nil) }

// EncodeGenericNack — buzilgan header'ga universal rad javobi (§4.3).
// seq=0 "originalni decode qilib bo'lmadi" degani (§4.3.1; 11-bobda batafsil).
func EncodeGenericNack(status, seq uint32) []byte { return encodePDU(CmdGenericNack, status, seq, nil) }
