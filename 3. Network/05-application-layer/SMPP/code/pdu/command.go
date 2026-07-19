// Package pdu SMPP v3.4 PDU (Protocol Data Unit) turlari va codec'larini
// o'z ichiga oladi. Har bir PDU 16 oktetlik header bilan boshlanadi; header'ning
// ikkinchi field'i — command_id — PDU turini aniqlaydi (v3.4 §3.2, §5.1.2).
package pdu

import "fmt"

// CommandID — PDU header'idagi command_id field'i (v3.4 §5.1.2).
// Request PDU'lar 0x00000000–0x000001FF diapazonida yotadi; response PDU'da
// bit 31 o'rnatiladi, ya'ni diapazon 0x80000000–0x800001FF bo'ladi (v3.4 §3.2).
type CommandID uint32

// v3.4 Table 5-1 (§5.1.2.1) bo'yicha command_id qiymatlari.
// outbind va alert_notification'ning response'i yo'q (v3.4 §2.8, §4.1.7.1, §4.12.1);
// generic_nack faqat response ko'rinishida mavjud (v3.4 §4.3.1).
const (
	CmdGenericNack         CommandID = 0x80000000
	CmdBindReceiver        CommandID = 0x00000001
	CmdBindReceiverResp    CommandID = 0x80000001
	CmdBindTransmitter     CommandID = 0x00000002
	CmdBindTransmitterResp CommandID = 0x80000002
	CmdQuerySM             CommandID = 0x00000003
	CmdQuerySMResp         CommandID = 0x80000003
	CmdSubmitSM            CommandID = 0x00000004
	CmdSubmitSMResp        CommandID = 0x80000004
	CmdDeliverSM           CommandID = 0x00000005
	CmdDeliverSMResp       CommandID = 0x80000005
	CmdUnbind              CommandID = 0x00000006
	CmdUnbindResp          CommandID = 0x80000006
	CmdReplaceSM           CommandID = 0x00000007
	CmdReplaceSMResp       CommandID = 0x80000007
	CmdCancelSM            CommandID = 0x00000008
	CmdCancelSMResp        CommandID = 0x80000008
	CmdBindTransceiver     CommandID = 0x00000009
	CmdBindTransceiverResp CommandID = 0x80000009
	CmdOutbind             CommandID = 0x0000000B
	CmdEnquireLink         CommandID = 0x00000015
	CmdEnquireLinkResp     CommandID = 0x80000015
	CmdSubmitMulti         CommandID = 0x00000021
	CmdSubmitMultiResp     CommandID = 0x80000021
	CmdAlertNotification   CommandID = 0x00000102
	CmdDataSM              CommandID = 0x00000103
	CmdDataSMResp          CommandID = 0x80000103
)

// respBit — command_id'ning 31-biti: response PDU belgisi (v3.4 §5.1.2).
const respBit CommandID = 0x80000000

// IsResponse id response PDU ekanini bildiradi (bit 31 o'rnatilgan bo'lsa true).
func (id CommandID) IsResponse() bool { return id&respBit != 0 }

// Resp request command_id'ga mos response command_id'ni qaytaradi
// (bit 31 o'rnatiladi). Javobi bo'lmagan PDU'lar (CmdOutbind, CmdAlertNotification)
// uchun chaqirish ma'nosiz — Table 5-1'da bunday resp qiymatlar Reserved.
func (id CommandID) Resp() CommandID { return id | respBit }

// commandNames — Table 5-1'dagi rasmiy protokol nomlari.
var commandNames = map[CommandID]string{
	CmdGenericNack:         "generic_nack",
	CmdBindReceiver:        "bind_receiver",
	CmdBindReceiverResp:    "bind_receiver_resp",
	CmdBindTransmitter:     "bind_transmitter",
	CmdBindTransmitterResp: "bind_transmitter_resp",
	CmdQuerySM:             "query_sm",
	CmdQuerySMResp:         "query_sm_resp",
	CmdSubmitSM:            "submit_sm",
	CmdSubmitSMResp:        "submit_sm_resp",
	CmdDeliverSM:           "deliver_sm",
	CmdDeliverSMResp:       "deliver_sm_resp",
	CmdUnbind:              "unbind",
	CmdUnbindResp:          "unbind_resp",
	CmdReplaceSM:           "replace_sm",
	CmdReplaceSMResp:       "replace_sm_resp",
	CmdCancelSM:            "cancel_sm",
	CmdCancelSMResp:        "cancel_sm_resp",
	CmdBindTransceiver:     "bind_transceiver",
	CmdBindTransceiverResp: "bind_transceiver_resp",
	CmdOutbind:             "outbind",
	CmdEnquireLink:         "enquire_link",
	CmdEnquireLinkResp:     "enquire_link_resp",
	CmdSubmitMulti:         "submit_multi",
	CmdSubmitMultiResp:     "submit_multi_resp",
	CmdAlertNotification:   "alert_notification",
	CmdDataSM:              "data_sm",
	CmdDataSMResp:          "data_sm_resp",
}

// String Table 5-1'dagi rasmiy nomni qaytaradi; notanish id uchun hex ko'rinish.
// Notanish id'lar log'da aynan spec'dagi raqam bilan ko'rinishi debugging'da muhim.
func (id CommandID) String() string {
	if name, ok := commandNames[id]; ok {
		return name
	}
	return fmt.Sprintf("unknown(0x%08X)", uint32(id))
}
