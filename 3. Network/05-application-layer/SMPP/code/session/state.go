// Package session SMPP session state machine'ini va (12-bobdan boshlab)
// session runtime'ini beradi.
package session

import (
	"fmt"

	"smpp/pdu"
)

// State — v3.4 §2.2 bo'yicha session holati. Spec'da AYNAN 5 state bor;
// UNBOUND va OUTBOUND alohida state sifatida v5.0'da (§2.3) kiritilgan —
// v3.4 modelida YO'Q (library'lardagi 7-state enum'lar v5.0'dan olingan).
type State uint8

const (
	Open     State = iota // TCP ulangan, bind hali qilinmagan
	BoundTX               // bind_transmitter muvaffaqiyatli
	BoundRX               // bind_receiver muvaffaqiyatli
	BoundTRX              // bind_transceiver muvaffaqiyatli
	Closed                // unbind qilingan va/yoki ulanish uzilgan
)

// String v3.4 §2.2'dagi rasmiy state nomlarini qaytaradi.
func (s State) String() string {
	switch s {
	case Open:
		return "OPEN"
	case BoundTX:
		return "BOUND_TX"
	case BoundRX:
		return "BOUND_RX"
	case BoundTRX:
		return "BOUND_TRX"
	case Closed:
		return "CLOSED"
	}
	return fmt.Sprintf("State(%d)", uint8(s))
}

// stateMask — allowedStates jadvalining ichki bit to'plami.
type stateMask uint8

const (
	inOpen stateMask = 1 << iota
	inTX
	inRX
	inTRX
)

const anyBound = inTX | inRX | inTRX

// allowedStates — v3.4 Table 2-1 (§2.3): har PDU qaysi session state'da
// yuborilishi mumkin. Yo'nalish (ESME'danmi, SMSC'danmi) bu jadvalga kirmaydi —
// uni server tomon alohida tekshiradi (14-bob): masalan submit_sm faqat
// ESME'dan keladi, lekin state jihatdan BOUND_TX/TRX talab qilinadi.
var allowedStates = map[pdu.CommandID]stateMask{
	pdu.CmdBindTransmitter:     inOpen,
	pdu.CmdBindTransmitterResp: inOpen,
	pdu.CmdBindReceiver:        inOpen,
	pdu.CmdBindReceiverResp:    inOpen,
	pdu.CmdBindTransceiver:     inOpen,
	pdu.CmdBindTransceiverResp: inOpen,
	pdu.CmdOutbind:             inOpen,
	pdu.CmdUnbind:              anyBound,
	pdu.CmdUnbindResp:          anyBound,
	pdu.CmdSubmitSM:            inTX | inTRX,
	pdu.CmdSubmitSMResp:        inTX | inTRX,
	pdu.CmdSubmitMulti:         inTX | inTRX,
	pdu.CmdSubmitMultiResp:     inTX | inTRX,
	pdu.CmdDataSM:              anyBound, // yagona message PDU'si: uchala bound state'da ham
	pdu.CmdDataSMResp:          anyBound,
	pdu.CmdDeliverSM:           inRX | inTRX,
	pdu.CmdDeliverSMResp:       inRX | inTRX,
	pdu.CmdQuerySM:             inTX | inTRX,
	pdu.CmdQuerySMResp:         inTX | inTRX,
	pdu.CmdCancelSM:            inTX | inTRX,
	pdu.CmdCancelSMResp:        inTX | inTRX,
	pdu.CmdReplaceSM:           inTX, // Table 2-1: FAQAT BOUND_TX — TRX'da YO'Q!
	pdu.CmdReplaceSMResp:       inTX,
	pdu.CmdEnquireLink:         anyBound,
	pdu.CmdEnquireLinkResp:     anyBound,
	pdu.CmdAlertNotification:   inRX | inTRX,
	pdu.CmdGenericNack:         anyBound,
}

func (s State) mask() stateMask {
	switch s {
	case Open:
		return inOpen
	case BoundTX:
		return inTX
	case BoundRX:
		return inRX
	case BoundTRX:
		return inTRX
	}
	return 0 // Closed: hech qanday PDU yuborilmaydi
}

// CanSend id'ni s state'da yuborish Table 2-1 bo'yicha joizligini aytadi.
// Notanish command_id uchun har doim false.
func CanSend(id pdu.CommandID, s State) bool {
	return allowedStates[id]&s.mask() != 0
}
