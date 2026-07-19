package session

import (
	"testing"

	"smpp/pdu"
)

var allStates = []State{Open, BoundTX, BoundRX, BoundTRX, Closed}

// TestCanSendTable — v3.4 Table 2-1'ning to'liq table-driven aksi: har PDU
// uchun ruxsat etilgan state'lar ro'yxati; qolgan barcha state'larda false.
func TestCanSendTable(t *testing.T) {
	tests := []struct {
		id      pdu.CommandID
		allowed []State
	}{
		{pdu.CmdBindTransmitter, []State{Open}},
		{pdu.CmdBindTransmitterResp, []State{Open}},
		{pdu.CmdBindReceiver, []State{Open}},
		{pdu.CmdBindTransceiver, []State{Open}},
		{pdu.CmdOutbind, []State{Open}},
		{pdu.CmdUnbind, []State{BoundTX, BoundRX, BoundTRX}},
		{pdu.CmdUnbindResp, []State{BoundTX, BoundRX, BoundTRX}},
		{pdu.CmdSubmitSM, []State{BoundTX, BoundTRX}},
		{pdu.CmdSubmitSMResp, []State{BoundTX, BoundTRX}},
		{pdu.CmdSubmitMulti, []State{BoundTX, BoundTRX}},
		{pdu.CmdDataSM, []State{BoundTX, BoundRX, BoundTRX}},
		{pdu.CmdDeliverSM, []State{BoundRX, BoundTRX}},
		{pdu.CmdDeliverSMResp, []State{BoundRX, BoundTRX}},
		{pdu.CmdQuerySM, []State{BoundTX, BoundTRX}},
		{pdu.CmdCancelSM, []State{BoundTX, BoundTRX}},
		{pdu.CmdReplaceSM, []State{BoundTX}}, // TRX'da YO'Q — Table 2-1 nozikligi
		{pdu.CmdReplaceSMResp, []State{BoundTX}},
		{pdu.CmdEnquireLink, []State{BoundTX, BoundRX, BoundTRX}},
		{pdu.CmdEnquireLinkResp, []State{BoundTX, BoundRX, BoundTRX}},
		{pdu.CmdAlertNotification, []State{BoundRX, BoundTRX}},
		{pdu.CmdGenericNack, []State{BoundTX, BoundRX, BoundTRX}},
	}
	for _, tt := range tests {
		t.Run(tt.id.String(), func(t *testing.T) {
			allowed := make(map[State]bool, len(tt.allowed))
			for _, s := range tt.allowed {
				allowed[s] = true
			}
			for _, s := range allStates {
				if got := CanSend(tt.id, s); got != allowed[s] {
					t.Errorf("CanSend(%s, %s) = %v, kutilgan %v", tt.id, s, got, allowed[s])
				}
			}
		})
	}
}

func TestCanSendUnknownCommand(t *testing.T) {
	for _, s := range allStates {
		if CanSend(pdu.CommandID(0xAA), s) {
			t.Errorf("notanish command_id %s state'da ruxsat olmasligi kerak", s)
		}
	}
}

func TestClosedSendsNothing(t *testing.T) {
	ids := []pdu.CommandID{
		pdu.CmdBindTransceiver, pdu.CmdSubmitSM, pdu.CmdDeliverSM,
		pdu.CmdEnquireLink, pdu.CmdUnbind, pdu.CmdGenericNack,
	}
	for _, id := range ids {
		if CanSend(id, Closed) {
			t.Errorf("CLOSED state'da %s yuborilmasligi kerak", id)
		}
	}
}

func TestStateString(t *testing.T) {
	tests := []struct {
		s    State
		want string
	}{
		{Open, "OPEN"},
		{BoundTX, "BOUND_TX"},
		{BoundRX, "BOUND_RX"},
		{BoundTRX, "BOUND_TRX"},
		{Closed, "CLOSED"},
	}
	for _, tt := range tests {
		if got := tt.s.String(); got != tt.want {
			t.Errorf("String() = %q, kutilgan %q", got, tt.want)
		}
	}
}
