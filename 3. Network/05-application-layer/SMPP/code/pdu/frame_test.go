package pdu

import (
	"bytes"
	"errors"
	"io"
	"testing"
	"testing/iotest"
)

// enquireLinkHex — header-only PDU: enquire_link, seq=2 (v3.4 §4.11.1).
const enquireLinkHex = `00 00 00 10 00 00 00 15 00 00 00 00 00 00 00 02`

const testMaxSize = 4096

func TestReadFrameGolden(t *testing.T) {
	golden := mustHex(t, specBindTransmitterHex)
	r := bytes.NewReader(golden)

	frame, err := ReadFrame(r, testMaxSize)
	if err != nil {
		t.Fatalf("ReadFrame xatosi: %v", err)
	}
	if !bytes.Equal(frame, golden) {
		t.Errorf("frame = % X,\nkutilgan % X", frame, golden)
	}
	// Stream frame'lar orasida toza tugadi — io.EOF.
	if _, err := ReadFrame(r, testMaxSize); !errors.Is(err, io.EOF) {
		t.Errorf("bo'sh stream'da io.EOF kutilgan edi, keldi: %v", err)
	}
}

func TestReadFrameByteByByte(t *testing.T) {
	// "1 PDU = 1 TCP segment" degan faraz yo'q: stream'ni ataylab
	// bir baytdan tomizamiz — framing baribir to'g'ri ishlashi kerak.
	bind := mustHex(t, specBindTransmitterHex)
	enq := mustHex(t, enquireLinkHex)
	stream := iotest.OneByteReader(bytes.NewReader(append(append([]byte{}, bind...), enq...)))

	first, err := ReadFrame(stream, testMaxSize)
	if err != nil {
		t.Fatalf("birinchi frame: %v", err)
	}
	if !bytes.Equal(first, bind) {
		t.Errorf("birinchi frame bind_transmitter emas: % X", first)
	}
	second, err := ReadFrame(stream, testMaxSize)
	if err != nil {
		t.Fatalf("ikkinchi frame: %v", err)
	}
	if !bytes.Equal(second, enq) {
		t.Errorf("ikkinchi frame enquire_link emas: % X", second)
	}
}

func TestReadFrameLengthTooShort(t *testing.T) {
	// command_length=0x0C < 16 — header ham sig'maydi.
	r := bytes.NewReader(mustHex(t, `00 00 00 0C 00 00 00 15 00 00 00 00`))
	if _, err := ReadFrame(r, testMaxSize); !errors.Is(err, ErrFrameTooShort) {
		t.Errorf("ErrFrameTooShort kutilgan edi, keldi: %v", err)
	}
}

func TestReadFrameLengthTooLarge(t *testing.T) {
	// command_length=0x7FFFFFFF — 2 GB'lik "PDU". maxSize'siz bu make([]byte, 2<<30)
	// bo'lar edi; validatsiya allocation'dan OLDIN kesishi kerak.
	r := bytes.NewReader(mustHex(t, `7F FF FF FF 00 00 00 04 00 00 00 00 00 00 00 01`))
	if _, err := ReadFrame(r, testMaxSize); !errors.Is(err, ErrFrameTooLarge) {
		t.Errorf("ErrFrameTooLarge kutilgan edi, keldi: %v", err)
	}
}

func TestReadFrameTruncated(t *testing.T) {
	t.Run("length field o'rtasida", func(t *testing.T) {
		r := bytes.NewReader([]byte{0x00, 0x00})
		_, err := ReadFrame(r, testMaxSize)
		if err == nil || errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Errorf("qisman length'da ErrUnexpectedEOF'li xato kutilgan edi, keldi: %v", err)
		}
	})
	t.Run("body o'rtasida", func(t *testing.T) {
		full := mustHex(t, specBindTransmitterHex)
		r := bytes.NewReader(full[:20]) // header + 4 oktet body, qolgani "kelmagan"
		_, err := ReadFrame(r, testMaxSize)
		if !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Errorf("io.ErrUnexpectedEOF kutilgan edi, keldi: %v", err)
		}
	})
}
