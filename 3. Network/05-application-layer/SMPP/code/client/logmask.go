package client

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"unicode/utf8"

	"smpp/pdu"
)

// PII masking (16-bob): MSISDN va xabar matni — personal data. Log'ga
// tushgan PII'ni "keyin tozalab" bo'lmaydi (GDPR o'chirish so'rovi log
// pipeline'iga yetmaydi) — masking APPLICATION darajasida, log yozilishidan
// OLDIN bo'lishi shart. Uch strategiya: mask (qisman ko'rinish — kundalik
// debug), hash (deterministik korrelyatsiya — incident tekshiruvi),
// redact (to'liq yashirish — matn uchun default).

// MaskMSISDN raqamning prefiksini (davlat/operator kodi — routing debug
// uchun kerak) va oxirgi 2 raqamini qoldiradi: "998901234567" →
// "9989******67". Qisqa/bo'sh qiymatlar butunlay yulduzlanadi.
func MaskMSISDN(msisdn string) string {
	if len(msisdn) < 7 {
		return strings.Repeat("*", len(msisdn))
	}
	return msisdn[:4] + strings.Repeat("*", len(msisdn)-6) + msisdn[len(msisdn)-2:]
}

// RedactText matnni to'liq yashiradi, faqat uzunlik qoladi (u ham signal:
// segment hisobini tekshirishga yetadi).
func RedactText(sm []byte) string {
	return fmt.Sprintf("[%d belgi yashirildi]", utf8.RuneCount(sm))
}

// HashPII deterministik qisqa hash qaytaradi — qiymatni oshkor qilmasdan
// korrelyatsiya: bir xil raqam har doim bir xil hash, log'lar bo'ylab
// incident'ni kuzatish mumkin, raqamni tiklab BO'LMAYDI.
func HashPII(v string) string {
	sum := sha256.Sum256([]byte(v))
	return hex.EncodeToString(sum[:6]) // 48 bit — log korrelyatsiyasiga yetarli
}

// MaskedSubmit submit_sm'ning log-safe bir qatorli tavsifi: manzillar
// masked+hashed, matn redacted. To'liq hex dump FAQAT debug rejimida va
// baribir shu masking'dan keyin bo'lishi kerak.
func MaskedSubmit(sm pdu.SubmitSM) string {
	return fmt.Sprintf("submit_sm src=%s dst=%s(#%s) dc=0x%02X reg=0x%02X sm=%s",
		MaskMSISDN(sm.Source.Addr),
		MaskMSISDN(sm.Dest.Addr), HashPII(sm.Dest.Addr),
		sm.DataCoding, uint8(sm.RegisteredDelivery),
		RedactText(sm.ShortMessage))
}
