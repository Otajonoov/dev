package pdu

import (
	"bytes"
	"errors"
	"testing"
)

func TestUintHelpersBigEndian(t *testing.T) {
	var b bytes.Buffer
	writeUint32(&b, 0x0000002F)
	writeUint16(&b, 0x0210)
	writeUint8(&b, 0x34)

	want := []byte{0x00, 0x00, 0x00, 0x2F, 0x02, 0x10, 0x34}
	if !bytes.Equal(b.Bytes(), want) {
		t.Fatalf("big-endian yozuv: % X, kutilgan % X", b.Bytes(), want)
	}

	r := bytes.NewReader(b.Bytes())
	if v, err := readUint32(r, "u32"); err != nil || v != 0x2F {
		t.Errorf("readUint32 = 0x%X, %v", v, err)
	}
	if v, err := readUint16(r, "u16"); err != nil || v != 0x0210 {
		t.Errorf("readUint16 = 0x%X, %v", v, err)
	}
	if v, err := readUint8(r, "u8"); err != nil || v != 0x34 {
		t.Errorf("readUint8 = 0x%X, %v", v, err)
	}
}

func TestCStringRoundTrip(t *testing.T) {
	tests := []struct {
		s       string
		max     int
		wantLen int // yozilgan oktetlar soni (NULL bilan)
	}{
		{"", 16, 1}, // bo'sh string = yagona 0x00 (§3.1)
		{"SMPP3TEST", 16, 10},
		{"secret08", 9, 9}, // aynan max: 8 belgi + NULL
	}
	for _, tt := range tests {
		var b bytes.Buffer
		if err := writeCString(&b, tt.s, tt.max, "f"); err != nil {
			t.Errorf("writeCString(%q): %v", tt.s, err)
			continue
		}
		if b.Len() != tt.wantLen {
			t.Errorf("writeCString(%q) %d oktet yozdi, kutilgan %d", tt.s, b.Len(), tt.wantLen)
		}
		got, err := readCString(bytes.NewReader(b.Bytes()), tt.max, "f")
		if err != nil || got != tt.s {
			t.Errorf("readCString = %q, %v; kutilgan %q", got, err, tt.s)
		}
	}
}

func TestWriteCStringTooLong(t *testing.T) {
	// "Max 9" NULL'ni o'z ichiga oladi (§3.1 note iii): 9 belgili string sig'maydi.
	var b bytes.Buffer
	if err := writeCString(&b, "secret089", 9, "password"); err == nil {
		t.Error("9 belgili string max 9 field'ga sig'masligi kerak edi")
	}
}

func TestWriteCStringEmbeddedNull(t *testing.T) {
	var b bytes.Buffer
	if err := writeCString(&b, "ab\x00cd", 16, "f"); err == nil {
		t.Error("ichki NULL baytli string rad etilishi kerak edi")
	}
}

func TestReadCStringNoTerminator(t *testing.T) {
	// max chegaragacha NULL yo'q — ErrNoTerminator.
	data := []byte("0123456789ABCDEF") // 16 oktet, terminatorsiz
	_, err := readCString(bytes.NewReader(data), 16, "system_id")
	if !errors.Is(err, ErrNoTerminator) {
		t.Errorf("ErrNoTerminator kutilgan edi, keldi: %v", err)
	}
}

func TestReadCStringStreamEnds(t *testing.T) {
	// Stream terminatordan oldin tugadi.
	_, err := readCString(bytes.NewReader([]byte("abc")), 16, "f")
	if err == nil {
		t.Error("terminatorsiz tugagan stream'da xato kutilgan edi")
	}
}
