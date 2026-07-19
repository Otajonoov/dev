package pdu

import (
	"bytes"
	"testing"
)

func TestEncodeEnquireLinkGolden(t *testing.T) {
	// 2-bobdagi golden bilan aynan mos kelishi kerak (frame_test.go'dagi const).
	got := EncodeEnquireLink(2)
	want := mustHex(t, enquireLinkHex)
	if !bytes.Equal(got, want) {
		t.Errorf("EncodeEnquireLink(2) = % X, kutilgan % X", got, want)
	}
}

func TestHeaderOnlyPDUs(t *testing.T) {
	tests := []struct {
		name   string
		frame  []byte
		id     CommandID
		status uint32
		seq    uint32
	}{
		{"enquire_link_resp", EncodeEnquireLinkResp(28), CmdEnquireLinkResp, 0, 28},
		{"unbind", EncodeUnbind(99), CmdUnbind, 0, 99},
		{"unbind_resp", EncodeUnbindResp(0, 99), CmdUnbindResp, 0, 99},
		{"generic_nack seq'li", EncodeGenericNack(0x02, 5), CmdGenericNack, 0x02, 5},
		{"generic_nack seq=0 (decode-fail)", EncodeGenericNack(0x02, 0), CmdGenericNack, 0x02, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.frame) != HeaderSize {
				t.Fatalf("header-only PDU %d oktet, 16 bo'lishi kerak", len(tt.frame))
			}
			h, err := DecodeHeader(tt.frame)
			if err != nil {
				t.Fatalf("DecodeHeader xatosi: %v", err)
			}
			want := Header{Length: HeaderSize, ID: tt.id, Status: tt.status, Sequence: tt.seq}
			if h != want {
				t.Errorf("header = %+v, kutilgan %+v", h, want)
			}
		})
	}
}
