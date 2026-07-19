package pdu

import "fmt"

// CommandStatus — header'dagi command_status field'ining tipli ko'rinishi
// (v3.4 §5.1.3, Table 5-2). Request PDU'larda doim 0 (NULL); ma'no faqat
// response'larda. Header struct'ida xom uint32 saqlanadi — bu tip nomlash,
// log va tasnif uchun: CommandStatus(h.Status).String().
//
// DIQQAT — bu jadval eski 09-smpp.md darsidagi XATO jadvalning to'g'rilangan
// versiyasi: 0x0E = RINVPASWD (0x06 EMAS — u RINVPRTFLG), 0x14 = RMSGQFUL
// (0x0A EMAS — u RINVSRCADR), RINVDCS degan kod v3.4'da YO'Q (u v5.0'ning
// 0x104 oilasi). Har qiymat spec Table 5-2'dan aynan olingan.
type CommandStatus uint32

// v3.4 Table 5-2 (§5.1.3) — to'liq rasmiy ro'yxat.
const (
	StatusROK              CommandStatus = 0x00000000 // No Error
	StatusRInvMsgLen       CommandStatus = 0x00000001 // Message Length is invalid
	StatusRInvCmdLen       CommandStatus = 0x00000002 // Command Length is invalid
	StatusRInvCmdID        CommandStatus = 0x00000003 // Invalid Command ID
	StatusRInvBndSts       CommandStatus = 0x00000004 // Incorrect BIND Status for given command
	StatusRAlyBnd          CommandStatus = 0x00000005 // ESME Already in Bound State
	StatusRInvPrtFlg       CommandStatus = 0x00000006 // Invalid Priority Flag
	StatusRInvRegDlvFlg    CommandStatus = 0x00000007 // Invalid Registered Delivery Flag
	StatusRSysErr          CommandStatus = 0x00000008 // System Error
	StatusRInvSrcAdr       CommandStatus = 0x0000000A // Invalid Source Address
	StatusRInvDstAdr       CommandStatus = 0x0000000B // Invalid Dest Addr
	StatusRInvMsgID        CommandStatus = 0x0000000C // Message ID is invalid
	StatusRBindFail        CommandStatus = 0x0000000D // Bind Failed
	StatusRInvPaswd        CommandStatus = 0x0000000E // Invalid Password
	StatusRInvSysID        CommandStatus = 0x0000000F // Invalid System ID
	StatusRCancelFail      CommandStatus = 0x00000011 // Cancel SM Failed
	StatusRReplaceFail     CommandStatus = 0x00000013 // Replace SM Failed
	StatusRMsgQFul         CommandStatus = 0x00000014 // Message Queue Full
	StatusRInvSerTyp       CommandStatus = 0x00000015 // Invalid Service Type
	StatusRInvNumDests     CommandStatus = 0x00000033 // Invalid number of destinations
	StatusRInvDLName       CommandStatus = 0x00000034 // Invalid Distribution List name
	StatusRInvDestFlag     CommandStatus = 0x00000040 // Destination flag is invalid (submit_multi)
	StatusRInvSubRep       CommandStatus = 0x00000042 // Invalid "submit with replace" request
	StatusRInvEsmClass     CommandStatus = 0x00000043 // Invalid esm_class field data
	StatusRCntSubDL        CommandStatus = 0x00000044 // Cannot Submit to Distribution List
	StatusRSubmitFail      CommandStatus = 0x00000045 // submit_sm or submit_multi failed
	StatusRInvSrcTON       CommandStatus = 0x00000048 // Invalid Source address TON
	StatusRInvSrcNPI       CommandStatus = 0x00000049 // Invalid Source address NPI
	StatusRInvDstTON       CommandStatus = 0x00000050 // Invalid Destination address TON
	StatusRInvDstNPI       CommandStatus = 0x00000051 // Invalid Destination address NPI
	StatusRInvSysTyp       CommandStatus = 0x00000053 // Invalid system_type field
	StatusRInvRepFlag      CommandStatus = 0x00000054 // Invalid replace_if_present flag
	StatusRInvNumMsgs      CommandStatus = 0x00000055 // Invalid number of messages
	StatusRThrottled       CommandStatus = 0x00000058 // Throttling error (rate limit oshildi)
	StatusRInvSched        CommandStatus = 0x00000061 // Invalid Scheduled Delivery Time
	StatusRInvExpiry       CommandStatus = 0x00000062 // Invalid message validity period
	StatusRInvDftMsgID     CommandStatus = 0x00000063 // Predefined Message Invalid or Not Found
	StatusRxTAppn          CommandStatus = 0x00000064 // ESME Receiver Temporary App Error
	StatusRxPAppn          CommandStatus = 0x00000065 // ESME Receiver Permanent App Error
	StatusRxRAppn          CommandStatus = 0x00000066 // ESME Receiver Reject Message Error
	StatusRQueryFail       CommandStatus = 0x00000067 // query_sm request failed
	StatusRInvOptParStream CommandStatus = 0x000000C0 // Error in the optional part of the PDU Body
	StatusROptParNotAllwd  CommandStatus = 0x000000C1 // Optional Parameter not allowed
	StatusRInvParLen       CommandStatus = 0x000000C2 // Invalid Parameter Length
	StatusRMissingOptParam CommandStatus = 0x000000C3 // Expected Optional Parameter missing
	StatusRInvOptParamVal  CommandStatus = 0x000000C4 // Invalid Optional Parameter Value
	StatusRDeliveryFailure CommandStatus = 0x000000FE // Delivery Failure (data_sm_resp uchun)
	StatusRUnknownErr      CommandStatus = 0x000000FF // Unknown Error
)

// Vendor-specific diapazon (Table 5-2): "Reserved for SMSC vendor specific".
const (
	vendorStatusLo CommandStatus = 0x00000400
	vendorStatusHi CommandStatus = 0x000004FF
)

// IsVendor kod SMSC vendor-specific blokida ekanini bildiradi — ma'nosi
// faqat operator hujjatida (balans tugashi, route yo'qligi va h.k.).
func (s CommandStatus) IsVendor() bool { return s >= vendorStatusLo && s <= vendorStatusHi }

var statusNames = map[CommandStatus]string{
	StatusROK:              "ESME_ROK",
	StatusRInvMsgLen:       "ESME_RINVMSGLEN",
	StatusRInvCmdLen:       "ESME_RINVCMDLEN",
	StatusRInvCmdID:        "ESME_RINVCMDID",
	StatusRInvBndSts:       "ESME_RINVBNDSTS",
	StatusRAlyBnd:          "ESME_RALYBND",
	StatusRInvPrtFlg:       "ESME_RINVPRTFLG",
	StatusRInvRegDlvFlg:    "ESME_RINVREGDLVFLG",
	StatusRSysErr:          "ESME_RSYSERR",
	StatusRInvSrcAdr:       "ESME_RINVSRCADR",
	StatusRInvDstAdr:       "ESME_RINVDSTADR",
	StatusRInvMsgID:        "ESME_RINVMSGID",
	StatusRBindFail:        "ESME_RBINDFAIL",
	StatusRInvPaswd:        "ESME_RINVPASWD",
	StatusRInvSysID:        "ESME_RINVSYSID",
	StatusRCancelFail:      "ESME_RCANCELFAIL",
	StatusRReplaceFail:     "ESME_RREPLACEFAIL",
	StatusRMsgQFul:         "ESME_RMSGQFUL",
	StatusRInvSerTyp:       "ESME_RINVSERTYP",
	StatusRInvNumDests:     "ESME_RINVNUMDESTS",
	StatusRInvDLName:       "ESME_RINVDLNAME",
	StatusRInvDestFlag:     "ESME_RINVDESTFLAG",
	StatusRInvSubRep:       "ESME_RINVSUBREP",
	StatusRInvEsmClass:     "ESME_RINVESMCLASS",
	StatusRCntSubDL:        "ESME_RCNTSUBDL",
	StatusRSubmitFail:      "ESME_RSUBMITFAIL",
	StatusRInvSrcTON:       "ESME_RINVSRCTON",
	StatusRInvSrcNPI:       "ESME_RINVSRCNPI",
	StatusRInvDstTON:       "ESME_RINVDSTTON",
	StatusRInvDstNPI:       "ESME_RINVDSTNPI",
	StatusRInvSysTyp:       "ESME_RINVSYSTYP",
	StatusRInvRepFlag:      "ESME_RINVREPFLAG",
	StatusRInvNumMsgs:      "ESME_RINVNUMMSGS",
	StatusRThrottled:       "ESME_RTHROTTLED",
	StatusRInvSched:        "ESME_RINVSCHED",
	StatusRInvExpiry:       "ESME_RINVEXPIRY",
	StatusRInvDftMsgID:     "ESME_RINVDFTMSGID",
	StatusRxTAppn:          "ESME_RX_T_APPN",
	StatusRxPAppn:          "ESME_RX_P_APPN",
	StatusRxRAppn:          "ESME_RX_R_APPN",
	StatusRQueryFail:       "ESME_RQUERYFAIL",
	StatusRInvOptParStream: "ESME_RINVOPTPARSTREAM",
	StatusROptParNotAllwd:  "ESME_ROPTPARNOTALLWD",
	StatusRInvParLen:       "ESME_RINVPARLEN",
	StatusRMissingOptParam: "ESME_RMISSINGOPTPARAM",
	StatusRInvOptParamVal:  "ESME_RINVOPTPARAMVAL",
	StatusRDeliveryFailure: "ESME_RDELIVERYFAILURE",
	StatusRUnknownErr:      "ESME_RUNKNOWNERR",
}

// String Table 5-2'dagi rasmiy nomni qaytaradi. Notanish kod hex ko'rinishda —
// log'da AYNAN spec yozuvida ko'rinishi tashxis uchun muhim (hex/dec
// aralashuvi — 0x14 ni "20" deb o'qib RMSGQFUL o'rniga boshqa kod deb
// o'ylash — klassik xato).
func (s CommandStatus) String() string {
	if name, ok := statusNames[s]; ok {
		return name
	}
	if s.IsVendor() {
		return fmt.Sprintf("vendor(0x%08X)", uint32(s))
	}
	return fmt.Sprintf("unknown(0x%08X)", uint32(s))
}
