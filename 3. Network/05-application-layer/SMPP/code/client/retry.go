// Package client ESME client API'sini beradi (13-bobda to'liq quriladi).
// Bu fayl — 11-bob milestone'i: command_status tasnifi va retry siyosati
// skeleti. 12–13-boblarda session engine va Client shunga ulanadi.
package client

import (
	"time"

	"smpp/pdu"
)

// Class — command_status'ga TO'G'RI reaksiya toifasi.
//
// MUHIM: bu tasnif spec'da RASMAN YO'Q (faqat RX_T_APPN/RX_P_APPN nomlarida
// ishora bor) — bu NowSMS/Kannel/aggregator amaliyotidan kelgan industriya
// konsensusi. Spec faqat kodlarni beradi; ularga qanday munosabatda bo'lish —
// implementatsiya qarori.
type Class int

const (
	// ClassTransient — vaqtinchalik: backoff bilan retry ma'noli
	// (RTHROTTLED, RMSGQFUL, RSYSERR, RX_T_APPN).
	ClassTransient Class = iota
	// ClassPermanent — retry MA'NOSIZ: xabar/so'rov o'zi noto'g'ri
	// (butun "Invalid *" oilasi, RX_P_APPN, RX_R_APPN). Qayta yuborish
	// natijani o'zgartirmaydi, faqat spam hisoblanadi.
	ClassPermanent
	// ClassSessionLevel — muammo xabarda emas, SESSIYADA: rebind/credential
	// yo'li (RINVBNDSTS, RALYBND, RINVPASWD, RINVSYSID, RBINDFAIL).
	// Xabarni retry qilish foydasiz — avval sessiya tuzatiladi.
	ClassSessionLevel
	// ClassUnknown — tasnifga tushmagan: vendor kodlari (0x400–0x4FF),
	// RSUBMITFAIL, RUNKNOWNERR va notanish qiymatlar. Default siyosat —
	// CHEKLANGAN retry + operator hujjatiga qarash.
	ClassUnknown
)

func (c Class) String() string {
	switch c {
	case ClassTransient:
		return "transient"
	case ClassPermanent:
		return "permanent"
	case ClassSessionLevel:
		return "session-level"
	}
	return "unknown"
}

// Classify command_status'ni reaksiya toifasiga ajratadi.
func Classify(status pdu.CommandStatus) Class {
	switch status {
	case pdu.StatusRThrottled, pdu.StatusRMsgQFul, pdu.StatusRSysErr, pdu.StatusRxTAppn:
		return ClassTransient

	case pdu.StatusRInvBndSts, pdu.StatusRAlyBnd, pdu.StatusRInvPaswd,
		pdu.StatusRInvSysID, pdu.StatusRBindFail:
		return ClassSessionLevel

	case pdu.StatusRInvMsgLen, pdu.StatusRInvCmdLen, pdu.StatusRInvCmdID,
		pdu.StatusRInvPrtFlg, pdu.StatusRInvRegDlvFlg,
		pdu.StatusRInvSrcAdr, pdu.StatusRInvDstAdr, pdu.StatusRInvMsgID,
		pdu.StatusRCancelFail, pdu.StatusRReplaceFail, pdu.StatusRInvSerTyp,
		pdu.StatusRInvNumDests, pdu.StatusRInvDLName, pdu.StatusRInvDestFlag,
		pdu.StatusRInvSubRep, pdu.StatusRInvEsmClass, pdu.StatusRCntSubDL,
		pdu.StatusRInvSrcTON, pdu.StatusRInvSrcNPI,
		pdu.StatusRInvDstTON, pdu.StatusRInvDstNPI, pdu.StatusRInvSysTyp,
		pdu.StatusRInvRepFlag, pdu.StatusRInvNumMsgs,
		pdu.StatusRInvSched, pdu.StatusRInvExpiry, pdu.StatusRInvDftMsgID,
		pdu.StatusRxPAppn, pdu.StatusRxRAppn, pdu.StatusRQueryFail,
		pdu.StatusRInvOptParStream, pdu.StatusROptParNotAllwd,
		pdu.StatusRInvParLen, pdu.StatusRMissingOptParam,
		pdu.StatusRInvOptParamVal:
		return ClassPermanent
	}
	// RSUBMITFAIL, RDELIVERYFAILURE, RUNKNOWNERR, vendor va notanish kodlar.
	return ClassUnknown
}

// RetryPolicy — transient xatolar uchun exponential backoff siyosati.
// 12-bobda reconnect'ga, 13-bobda Client.Submit retry'siga ulanadi.
//
// MaxAge — "necha URINISH" emas, "qancha VAQT" chegarasi: RTHROTTLED va
// RMSGQFUL "keyinroq albatta o'tadi" degan signal bo'lgani uchun ularga
// attempt-limit emas, muddat-limit to'g'ri (NowSMS amaliyoti): xabar
// validity'si tugagunicha urinishda davom etish mumkin.
type RetryPolicy struct {
	BaseDelay time.Duration // 1-urinishdan keyingi kutish (masalan 1s)
	MaxDelay  time.Duration // backoff shifti (masalan 60s)
	MaxAge    time.Duration // birinchi urinishdan beri jami muddat chegarasi
}

// DefaultRetryPolicy — boshlang'ich nuqta; qiymatlar operator TZ'siga
// moslanadi (spec'da bunday raqamlar YO'Q).
var DefaultRetryPolicy = RetryPolicy{
	BaseDelay: 1 * time.Second,
	MaxDelay:  60 * time.Second,
	MaxAge:    1 * time.Hour,
}

// NextDelay n-urinishdan keyingi kutish vaqti (n 1'dan boshlanadi):
// BaseDelay×2^(n-1), MaxDelay bilan kesilgan. Jitter ATAYIN yo'q — u
// reconnect'da qo'shiladi (12-bob): bitta sessiya ichidagi retry'da
// determinizm testlarni soddalashtiradi.
func (p RetryPolicy) NextDelay(attempt int) time.Duration {
	if attempt < 1 {
		attempt = 1
	}
	d := p.BaseDelay
	for i := 1; i < attempt; i++ {
		d *= 2
		if d >= p.MaxDelay {
			return p.MaxDelay
		}
	}
	if d > p.MaxDelay {
		return p.MaxDelay
	}
	return d
}

// ShouldRetry status va birinchi urinish vaqtiga qarab yana urinish
// kerakligini aytadi. Qoida: faqat Transient — cheksiz (MaxAge ichida);
// Unknown — ehtiyotkor (MaxAge/4); Permanent va SessionLevel — hech qachon
// (session-level'da retry EMAS, rebind kerak — u boshqa yo'l).
func (p RetryPolicy) ShouldRetry(status pdu.CommandStatus, firstAttempt, now time.Time) bool {
	age := now.Sub(firstAttempt)
	switch Classify(status) {
	case ClassTransient:
		return age < p.MaxAge
	case ClassUnknown:
		return age < p.MaxAge/4
	}
	return false
}
