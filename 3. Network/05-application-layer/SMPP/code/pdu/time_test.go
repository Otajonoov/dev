package pdu

import (
	"testing"
	"time"
)

func TestParseSpecRelativeExample(t *testing.T) {
	// Spec'ning o'z misoli (§7.1.1.2): hozirdan 2 yil 6 oy 10 kun 23:34:29 keyin.
	v, err := ParseTime("020610233429000R")
	if err != nil {
		t.Fatalf("ParseTime xatosi: %v", err)
	}
	want := TimeValue{Relative: true, Years: 2, Months: 6, Days: 10, Hours: 23, Minutes: 34, Seconds: 29}
	if v != want {
		t.Errorf("ParseTime = %+v,\nkutilgan %+v", v, want)
	}
	// Encode qaytib aynan spec matnini berishi kerak.
	s, err := EncodeRelativeValue(v)
	if err != nil {
		t.Fatalf("EncodeRelativeValue xatosi: %v", err)
	}
	if s != "020610233429000R" {
		t.Errorf("EncodeRelativeValue = %q, kutilgan spec misoli", s)
	}
}

func TestEncodeRelativeDuration(t *testing.T) {
	// 1 kunlik validity — eng tipik amaliy qiymat.
	s, err := EncodeRelative(24 * time.Hour)
	if err != nil {
		t.Fatalf("EncodeRelative xatosi: %v", err)
	}
	if s != "000001000000000R" {
		t.Errorf("1 kun = %q, kutilgan 000001000000000R", s)
	}
	// 2 soat 30 minut.
	s, err = EncodeRelative(2*time.Hour + 30*time.Minute)
	if err != nil {
		t.Fatalf("EncodeRelative xatosi: %v", err)
	}
	if s != "000000023000000R" {
		t.Errorf("2h30m = %q, kutilgan 000000023000000R", s)
	}
	// 100+ kun — DD'ga sig'maydi.
	if _, err := EncodeRelative(100 * 24 * time.Hour); err == nil {
		t.Error("100 kunlik Duration xato qaytarishi kerak edi")
	}
}

func TestAbsoluteRoundTrip(t *testing.T) {
	// Toshkent vaqti: UTC+5 = 20 chorak soat.
	tz := time.FixedZone("UZT", 5*3600)
	orig := time.Date(2026, 7, 17, 12, 30, 45, 300_000_000, tz)

	s, err := EncodeAbsolute(orig, tz)
	if err != nil {
		t.Fatalf("EncodeAbsolute xatosi: %v", err)
	}
	if s != "260717123045320+" {
		t.Fatalf("EncodeAbsolute = %q, kutilgan 260717123045320+", s)
	}

	v, err := ParseTime(s)
	if err != nil {
		t.Fatalf("ParseTime xatosi: %v", err)
	}
	if v.Relative || !v.HasOffset {
		t.Fatalf("absolute + offset kutilgan edi: %+v", v)
	}
	if !v.At.Equal(orig) {
		t.Errorf("round-trip: %v != %v", v.At, orig)
	}
}

func TestParseAbsoluteNegativeOffset(t *testing.T) {
	// UTC-4 (16 chorak soat, '-').
	v, err := ParseTime("990825120000016-")
	if err != nil {
		t.Fatalf("ParseTime xatosi: %v", err)
	}
	// Appendix C sliding window: 99 → 1999.
	want := time.Date(1999, 8, 25, 12, 0, 0, 0, time.FixedZone("", -4*3600))
	if !v.At.Equal(want) {
		t.Errorf("ParseTime = %v, kutilgan %v", v.At, want)
	}
}

func TestParseSMSCLocalVariant(t *testing.T) {
	// 12 belgili SMSC lokal vaqt (§7.1.1 Note) — masalan query_sm_resp final_date.
	v, err := ParseTime("260717120000")
	if err != nil {
		t.Fatalf("ParseTime xatosi: %v", err)
	}
	if v.Relative || v.HasOffset {
		t.Errorf("12 belgili variant: Relative=false, HasOffset=false kutilgan: %+v", v)
	}
	if v.At.Year() != 2026 || v.At.Month() != 7 || v.At.Day() != 17 {
		t.Errorf("sana noto'g'ri: %v", v.At)
	}
}

func TestParseTimeErrors(t *testing.T) {
	bad := []string{
		"02061023342900R",  // 15 belgi
		"020610233429000X", // noma'lum p
		"0206102334290A0R", // raqam emas
		"260717120000049+", // nn=49 > 48
		"26071712000",      // 11 belgi
	}
	for _, s := range bad {
		if _, err := ParseTime(s); err == nil {
			t.Errorf("ParseTime(%q) xato qaytarishi kerak edi", s)
		}
	}
}

func TestEncodeAbsoluteBadZone(t *testing.T) {
	// Chorak soatga karrali bo'lmagan zona (masalan UTC+5:37 degan uydirma).
	tz := time.FixedZone("ODD", 5*3600+37*60)
	if _, err := EncodeAbsolute(time.Now(), tz); err == nil {
		t.Error("chorak soatga karrali bo'lmagan zona xato berishi kerak edi")
	}
}
