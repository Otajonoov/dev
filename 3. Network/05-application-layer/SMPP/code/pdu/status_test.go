package pdu

import (
	"strings"
	"testing"
)

// Eski 09-smpp.md darsidagi xato jadval tuzatilganining isboti — eng ko'p
// adashtiriladigan qiymatlar aynan spec Table 5-2 bo'yicha.
func TestStatusCriticalValues(t *testing.T) {
	tests := []struct {
		s    CommandStatus
		code uint32
		name string
	}{
		{StatusROK, 0x00, "ESME_ROK"},
		{StatusRInvPrtFlg, 0x06, "ESME_RINVPRTFLG"}, // 0x06 RINVPASWD EMAS!
		{StatusRInvSrcAdr, 0x0A, "ESME_RINVSRCADR"}, // 0x0A RMSGQFUL EMAS!
		{StatusRInvPaswd, 0x0E, "ESME_RINVPASWD"},   // to'g'ri joyi
		{StatusRInvSysID, 0x0F, "ESME_RINVSYSID"},
		{StatusRMsgQFul, 0x14, "ESME_RMSGQFUL"}, // to'g'ri joyi (dec 20)
		{StatusRThrottled, 0x58, "ESME_RTHROTTLED"},
		{StatusRxTAppn, 0x64, "ESME_RX_T_APPN"},
		{StatusRDeliveryFailure, 0xFE, "ESME_RDELIVERYFAILURE"},
		{StatusRUnknownErr, 0xFF, "ESME_RUNKNOWNERR"},
	}
	for _, tt := range tests {
		if uint32(tt.s) != tt.code {
			t.Errorf("%s = 0x%02X, kutilgan 0x%02X", tt.name, uint32(tt.s), tt.code)
		}
		if tt.s.String() != tt.name {
			t.Errorf("String(0x%02X) = %q, kutilgan %q", tt.code, tt.s.String(), tt.name)
		}
	}
}

// String() to'liq qamrov: har konstanta rasmiy ESME_* nom qaytaradi.
func TestStatusStringFullCoverage(t *testing.T) {
	if len(statusNames) != 48 {
		t.Fatalf("statusNames %d ta — Table 5-2'da 48 nomlangan kod bor", len(statusNames))
	}
	for s, name := range statusNames {
		if got := s.String(); got != name {
			t.Errorf("0x%08X: String() = %q, kutilgan %q", uint32(s), got, name)
		}
		if !strings.HasPrefix(name, "ESME_R") {
			t.Errorf("%q — rasmiy nom ESME_R bilan boshlanadi", name)
		}
	}
}

func TestStatusVendorAndUnknown(t *testing.T) {
	v := CommandStatus(0x00000410)
	if !v.IsVendor() || v.String() != "vendor(0x00000410)" {
		t.Errorf("vendor kod: IsVendor=%v String=%q", v.IsVendor(), v.String())
	}
	u := CommandStatus(0x00000030) // Reserved oraliq (0x16–0x32 dan keyingi bo'shliq)
	if u.IsVendor() || !strings.HasPrefix(u.String(), "unknown(") {
		t.Errorf("notanish kod: %q", u.String())
	}
	// RINVDCS v3.4'da YO'Q — v5.0'ning 0x104 kodi bizda notanish bo'lishi shart.
	if CommandStatus(0x104).String() != "unknown(0x00000104)" {
		t.Errorf("0x104 v3.4'da nomlanmasligi kerak: %q", CommandStatus(0x104).String())
	}
}
