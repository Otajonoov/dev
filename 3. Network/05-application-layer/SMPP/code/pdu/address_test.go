package pdu

import "testing"

func TestInternationalConstructor(t *testing.T) {
	tests := []struct {
		in   string
		want string // kutilgan Addr ("" = xato kutiladi)
	}{
		{"998901234567", "998901234567"},
		{"+998901234567", "998901234567"}, // '+' olib tashlanadi
		{"+15551234567", "15551234567"},
		{"", ""},
		{"+", ""},
		{"998-90-123", ""},       // raqam bo'lmagan belgi
		{"9989012345678901", ""}, // 16 raqam — E.164 max 15
	}
	for _, tt := range tests {
		a, err := International(tt.in)
		if tt.want == "" {
			if err == nil {
				t.Errorf("International(%q): xato kutilgan edi, keldi %+v", tt.in, a)
			}
			continue
		}
		if err != nil {
			t.Errorf("International(%q): kutilmagan xato %v", tt.in, err)
			continue
		}
		if a.TON != TONInternational || a.NPI != NPIISDN || a.Addr != tt.want {
			t.Errorf("International(%q) = %+v, kutilgan 1/1/%q", tt.in, a, tt.want)
		}
	}
}

func TestAlphanumericConstructor(t *testing.T) {
	tests := []struct {
		in string
		ok bool
	}{
		{"Bank", true},
		{"MyBrand 24", true},    // bo'sh joy GSM7'da bor
		{"O'zBank", true},       // ASCII apostrof GSM7'da bor
		{"ELEVENCHARS", true},   // aynan 11 belgi
		{"TWELVECHARSX", false}, // 12 belgi — TP-OA limiti
		{"", false},
		{"Uz~Card", false}, // '~' GSM7 extension'da (2 septet), basic'da yo'q
		{"Ozbankʻ", false}, // U+02BB — GSM7'da yo'q (7-bob mavzusi)
	}
	for _, tt := range tests {
		a, err := Alphanumeric(tt.in)
		if tt.ok && err != nil {
			t.Errorf("Alphanumeric(%q): kutilmagan xato %v", tt.in, err)
		}
		if !tt.ok && err == nil {
			t.Errorf("Alphanumeric(%q): xato kutilgan edi, keldi %+v", tt.in, a)
		}
		if tt.ok && (a.TON != TONAlphanumeric || a.NPI != NPIUnknown) {
			t.Errorf("Alphanumeric(%q) TON/NPI = %d/%d, kutilgan 5/0", tt.in, a.TON, a.NPI)
		}
	}
}

func TestShortCodeConstructor(t *testing.T) {
	if a, err := ShortCode("1234"); err != nil || a.TON != TONNetworkSpecific || a.NPI != NPIUnknown {
		t.Errorf("ShortCode(1234) = %+v, %v; kutilgan 3/0", a, err)
	}
	for _, bad := range []string{"", "12a4", "123456789"} {
		if _, err := ShortCode(bad); err == nil {
			t.Errorf("ShortCode(%q): xato kutilgan edi", bad)
		}
	}
}

func TestValidateTable(t *testing.T) {
	tests := []struct {
		name string
		a    Address
		ok   bool
	}{
		{"NULL source", NullSource(), true},
		{"international to'g'ri", Address{TON: 1, NPI: 1, Addr: "998901234567"}, true},
		{"international '+' bilan", Address{TON: 1, NPI: 1, Addr: "+998901234567"}, false},
		{"international bo'sh", Address{TON: 1, NPI: 1, Addr: ""}, false},
		{"international harf bilan", Address{TON: 1, NPI: 1, Addr: "99890ABC"}, false},
		{"bo'sh addr, lekin TON!=0", Address{TON: 2, NPI: 1, Addr: ""}, false},
		{"20 belgili addr (chegara)", Address{TON: 0, NPI: 1, Addr: "12345678901234567890"}, true},
		{"21 belgili addr", Address{TON: 0, NPI: 1, Addr: "123456789012345678901"}, false},
		{"IP manzil NPI=14", Address{TON: 0, NPI: NPIInternet, Addr: "10.0.0.1"}, true},
		{"alphanumeric 12 belgi", Address{TON: 5, NPI: 0, Addr: "TWELVECHARSX"}, false},
		{"national raqam", Address{TON: 2, NPI: 1, Addr: "901234567"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.a.Validate()
			if tt.ok && err != nil {
				t.Errorf("kutilmagan xato: %v", err)
			}
			if !tt.ok && err == nil {
				t.Error("xato kutilgan edi")
			}
		})
	}
}

func TestNPIValuesNotSequential(t *testing.T) {
	// Table 5-4 tuzog'i hujjat sifatida testda: 2 qiymati NPI jadvalida YO'Q.
	vals := []uint8{NPIUnknown, NPIISDN, NPIData, NPITelex, NPILandMobile,
		NPINational, NPIPrivate, NPIERMES, NPIInternet, NPIWAP}
	want := []uint8{0, 1, 3, 4, 6, 8, 9, 10, 14, 18}
	for i, v := range vals {
		if v != want[i] {
			t.Errorf("NPI konstanta #%d = %d, Table 5-4 bo'yicha %d", i, v, want[i])
		}
	}
}
