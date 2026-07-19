package pdu

import "testing"

// requestPairs — Table 5-1'dagi request↔response juftliklari (resp'i bor 13 operatsiya).
var requestPairs = []struct {
	req      CommandID
	resp     CommandID
	reqName  string
	respName string
}{
	{CmdBindReceiver, CmdBindReceiverResp, "bind_receiver", "bind_receiver_resp"},
	{CmdBindTransmitter, CmdBindTransmitterResp, "bind_transmitter", "bind_transmitter_resp"},
	{CmdQuerySM, CmdQuerySMResp, "query_sm", "query_sm_resp"},
	{CmdSubmitSM, CmdSubmitSMResp, "submit_sm", "submit_sm_resp"},
	{CmdDeliverSM, CmdDeliverSMResp, "deliver_sm", "deliver_sm_resp"},
	{CmdUnbind, CmdUnbindResp, "unbind", "unbind_resp"},
	{CmdReplaceSM, CmdReplaceSMResp, "replace_sm", "replace_sm_resp"},
	{CmdCancelSM, CmdCancelSMResp, "cancel_sm", "cancel_sm_resp"},
	{CmdBindTransceiver, CmdBindTransceiverResp, "bind_transceiver", "bind_transceiver_resp"},
	{CmdEnquireLink, CmdEnquireLinkResp, "enquire_link", "enquire_link_resp"},
	{CmdSubmitMulti, CmdSubmitMultiResp, "submit_multi", "submit_multi_resp"},
	{CmdDataSM, CmdDataSMResp, "data_sm", "data_sm_resp"},
}

func TestRespRoundTrip(t *testing.T) {
	for _, p := range requestPairs {
		t.Run(p.reqName, func(t *testing.T) {
			if got := p.req.Resp(); got != p.resp {
				t.Errorf("%s.Resp() = 0x%08X, kutilgan 0x%08X", p.reqName, uint32(got), uint32(p.resp))
			}
			if p.req.IsResponse() {
				t.Errorf("%s request bo'lishi kerak, IsResponse()=true qaytdi", p.reqName)
			}
			if !p.resp.IsResponse() {
				t.Errorf("%s response bo'lishi kerak, IsResponse()=false qaytdi", p.respName)
			}
			// Resp() nomi ham "_resp" suffiksli rasmiy nom bilan mos kelishi kerak.
			if got := p.req.Resp().String(); got != p.respName {
				t.Errorf("%s.Resp().String() = %q, kutilgan %q", p.reqName, got, p.respName)
			}
		})
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		id   CommandID
		want string
	}{
		{CmdGenericNack, "generic_nack"},
		{CmdBindReceiver, "bind_receiver"},
		{CmdBindTransmitter, "bind_transmitter"},
		{CmdQuerySM, "query_sm"},
		{CmdSubmitSM, "submit_sm"},
		{CmdDeliverSM, "deliver_sm"},
		{CmdUnbind, "unbind"},
		{CmdReplaceSM, "replace_sm"},
		{CmdCancelSM, "cancel_sm"},
		{CmdBindTransceiver, "bind_transceiver"},
		{CmdOutbind, "outbind"},
		{CmdEnquireLink, "enquire_link"},
		{CmdSubmitMulti, "submit_multi"},
		{CmdAlertNotification, "alert_notification"},
		{CmdDataSM, "data_sm"},
		// Notanish id'lar — hex ko'rinishda.
		{CommandID(0x000000AA), "unknown(0x000000AA)"},
		{CommandID(0x80000063), "unknown(0x80000063)"},
	}
	for _, tt := range tests {
		if got := tt.id.String(); got != tt.want {
			t.Errorf("CommandID(0x%08X).String() = %q, kutilgan %q", uint32(tt.id), got, tt.want)
		}
	}
}

func TestTable51Values(t *testing.T) {
	// Table 5-1'dagi raqamlar kodda aynan spec'dagidek ekanini qotirib qo'yamiz:
	// konstanta qiymati tasodifan o'zgarsa shu test sinadi.
	values := []struct {
		id   CommandID
		want uint32
	}{
		{CmdGenericNack, 0x80000000},
		{CmdBindReceiver, 0x00000001},
		{CmdBindTransmitter, 0x00000002},
		{CmdQuerySM, 0x00000003},
		{CmdSubmitSM, 0x00000004},
		{CmdDeliverSM, 0x00000005},
		{CmdUnbind, 0x00000006},
		{CmdReplaceSM, 0x00000007},
		{CmdCancelSM, 0x00000008},
		{CmdBindTransceiver, 0x00000009},
		{CmdOutbind, 0x0000000B},
		{CmdEnquireLink, 0x00000015},
		{CmdSubmitMulti, 0x00000021},
		{CmdAlertNotification, 0x00000102},
		{CmdDataSM, 0x00000103},
	}
	for _, v := range values {
		if uint32(v.id) != v.want {
			t.Errorf("%s = 0x%08X, Table 5-1 bo'yicha 0x%08X bo'lishi kerak", v.id, uint32(v.id), v.want)
		}
	}
}

func TestGenericNackIsResponseOnly(t *testing.T) {
	// generic_nack faqat response ko'rinishida mavjud (v3.4 §4.3.1):
	// bit 31 doim o'rnatilgan, Resp() uni o'zgartirmaydi.
	if !CmdGenericNack.IsResponse() {
		t.Error("CmdGenericNack.IsResponse() = false, true bo'lishi kerak")
	}
	if got := CmdGenericNack.Resp(); got != CmdGenericNack {
		t.Errorf("CmdGenericNack.Resp() = 0x%08X, o'zi qolishi kerak edi", uint32(got))
	}
}
