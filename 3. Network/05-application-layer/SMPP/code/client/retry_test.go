package client

import (
	"testing"
	"time"

	"smpp/pdu"
)

func TestClassify(t *testing.T) {
	tests := []struct {
		status pdu.CommandStatus
		want   Class
	}{
		// Transient — backoff bilan retry.
		{pdu.StatusRThrottled, ClassTransient},
		{pdu.StatusRMsgQFul, ClassTransient},
		{pdu.StatusRSysErr, ClassTransient},
		{pdu.StatusRxTAppn, ClassTransient},
		// Session-level — rebind yo'li.
		{pdu.StatusRInvBndSts, ClassSessionLevel},
		{pdu.StatusRAlyBnd, ClassSessionLevel},
		{pdu.StatusRInvPaswd, ClassSessionLevel},
		{pdu.StatusRInvSysID, ClassSessionLevel},
		{pdu.StatusRBindFail, ClassSessionLevel},
		// Permanent — retry ma'nosiz.
		{pdu.StatusRInvSrcAdr, ClassPermanent},
		{pdu.StatusRInvDstAdr, ClassPermanent},
		{pdu.StatusRInvEsmClass, ClassPermanent},
		{pdu.StatusRxPAppn, ClassPermanent},
		{pdu.StatusRxRAppn, ClassPermanent},
		{pdu.StatusRInvParLen, ClassPermanent},
		{pdu.StatusRMissingOptParam, ClassPermanent},
		// Noaniq — vendor va umumiy kodlar.
		{pdu.StatusRSubmitFail, ClassUnknown},
		{pdu.StatusRUnknownErr, ClassUnknown},
		{pdu.StatusRDeliveryFailure, ClassUnknown},
		{pdu.CommandStatus(0x00000410), ClassUnknown}, // vendor
		{pdu.CommandStatus(0x00000030), ClassUnknown}, // reserved oraliq
		// 0x14/0x0A tuzog'i: dec 20 = RMSGQFUL transient, 0x0A esa
		// RINVSRCADR — permanent. Ikkalasini almashtirish qimmatga tushadi.
		{pdu.CommandStatus(20), ClassTransient},
		{pdu.CommandStatus(0x0A), ClassPermanent},
	}
	for _, tt := range tests {
		if got := Classify(tt.status); got != tt.want {
			t.Errorf("Classify(%s) = %s, kutilgan %s", tt.status, got, tt.want)
		}
	}
}

func TestNextDelay(t *testing.T) {
	p := RetryPolicy{BaseDelay: time.Second, MaxDelay: 8 * time.Second}
	want := []time.Duration{
		1 * time.Second, // 1-urinish
		2 * time.Second,
		4 * time.Second,
		8 * time.Second,
		8 * time.Second, // shift: MaxDelay'da to'xtaydi (overflow yo'q)
		8 * time.Second,
	}
	for i, w := range want {
		if got := p.NextDelay(i + 1); got != w {
			t.Errorf("NextDelay(%d) = %v, kutilgan %v", i+1, got, w)
		}
	}
	// attempt<1 himoyasi.
	if p.NextDelay(0) != time.Second {
		t.Error("NextDelay(0) BaseDelay bo'lishi kerak")
	}
	// Juda katta attempt overflow qilmasligi (2^63 tuzog'i).
	if got := p.NextDelay(200); got != 8*time.Second {
		t.Errorf("NextDelay(200) = %v", got)
	}
}

func TestShouldRetry(t *testing.T) {
	p := RetryPolicy{BaseDelay: time.Second, MaxDelay: time.Minute, MaxAge: time.Hour}
	start := time.Now()

	// Transient: MaxAge ichida ha, tashqarisida yo'q.
	if !p.ShouldRetry(pdu.StatusRThrottled, start, start.Add(30*time.Minute)) {
		t.Error("transient 30min: retry kutilgan")
	}
	if p.ShouldRetry(pdu.StatusRThrottled, start, start.Add(2*time.Hour)) {
		t.Error("transient 2h: MaxAge oshdi, retry emas")
	}
	// Unknown: qisqartirilgan oyna (MaxAge/4 = 15min).
	if !p.ShouldRetry(pdu.StatusRUnknownErr, start, start.Add(10*time.Minute)) {
		t.Error("unknown 10min: retry kutilgan")
	}
	if p.ShouldRetry(pdu.StatusRUnknownErr, start, start.Add(20*time.Minute)) {
		t.Error("unknown 20min: oyna yopilgan")
	}
	// Permanent va session-level: hech qachon.
	if p.ShouldRetry(pdu.StatusRInvDstAdr, start, start.Add(time.Second)) {
		t.Error("permanent: retry taqiqlangan")
	}
	if p.ShouldRetry(pdu.StatusRInvPaswd, start, start.Add(time.Second)) {
		t.Error("session-level: retry emas, rebind yo'li")
	}
}
