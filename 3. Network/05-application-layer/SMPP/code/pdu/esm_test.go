package pdu

import "testing"

func TestEsmClassBits(t *testing.T) {
	tests := []struct {
		name    string
		esm     EsmClass
		msgType MessageType
		isDLR   bool
		udhi    bool
	}{
		{"oddiy MO xabar", 0x00, TypeNormal, false, false},
		{"DLR", 0x04, TypeDeliveryReceipt, true, false},
		{"DLR + UDHI (0x44) — klassik tuzoq", 0x44, TypeDeliveryReceipt, true, true},
		{"MO + UDHI (0x40)", 0x40, TypeNormal, false, true},
		{"SME Delivery Ack", 0x08, TypeSMEDeliveryAck, false, false},
		{"SME Manual Ack", 0x10, TypeSMEManualAck, false, false},
		{"Conversation Abort", 0x18, TypeConversationAbort, false, false},
		{"Intermediate Notification", 0x20, TypeIntermediate, false, false},
		{"Reply Path + DLR", 0x84, TypeDeliveryReceipt, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.esm.MessageType(); got != tt.msgType {
				t.Errorf("MessageType() = 0x%02X, kutilgan 0x%02X", uint8(got), uint8(tt.msgType))
			}
			if got := tt.esm.IsDeliveryReceipt(); got != tt.isDLR {
				t.Errorf("IsDeliveryReceipt() = %v, kutilgan %v", got, tt.isDLR)
			}
			if got := tt.esm.HasUDHI(); got != tt.udhi {
				t.Errorf("HasUDHI() = %v, kutilgan %v", got, tt.udhi)
			}
		})
	}
}

func TestEsmClassNaiveComparisonTrap(t *testing.T) {
	// Hujjatlashtirilgan tuzoq: esm_class == 0x04 deb solishtirish UDHI'li
	// DLR'ni (0x44) o'tkazib yuboradi. Bizning helper ikkalasini ham ushlaydi.
	var udhiDLR EsmClass = 0x44
	if udhiDLR == 0x04 {
		t.Fatal("bu tenglik hech qachon true bo'lmasligi kerak edi")
	}
	if !udhiDLR.IsDeliveryReceipt() {
		t.Error("0x44 DLR sifatida tanilishi kerak")
	}
}

func TestEsmClassMode(t *testing.T) {
	if got := (ModeDatagram | FlagUDHI).Mode(); got != ModeDatagram {
		t.Errorf("Mode() = 0x%02X, kutilgan datagram (0x01)", uint8(got))
	}
	if got := EsmClass(0x44).Mode(); got != ModeDefault {
		t.Errorf("Mode() = 0x%02X, kutilgan default (0x00)", uint8(got))
	}
}

func TestEsmClassWithUDHI(t *testing.T) {
	if got := EsmClass(0x04).WithUDHI(); got != 0x44 {
		t.Errorf("WithUDHI() = 0x%02X, kutilgan 0x44", uint8(got))
	}
}

func TestRegisteredDeliveryBits(t *testing.T) {
	tests := []struct {
		name         string
		rd           RegisteredDelivery
		wantsDLR     bool
		wantsInterim bool
	}{
		{"default — hech narsa", 0x00, false, false},
		{"final DLR", 0x01, true, false},
		{"faqat failure DLR", 0x02, true, false},
		{"DLR + intermediate (0x11)", 0x11, true, true},
		{"faqat intermediate", 0x10, false, true},
		{"SME ack DLR'siz", 0x04, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rd.WantsDLR(); got != tt.wantsDLR {
				t.Errorf("WantsDLR() = %v, kutilgan %v", got, tt.wantsDLR)
			}
			if got := tt.rd.WantsIntermediate(); got != tt.wantsInterim {
				t.Errorf("WantsIntermediate() = %v, kutilgan %v", got, tt.wantsInterim)
			}
		})
	}
}

func TestIntermediateIsBit4(t *testing.T) {
	// Erratum qotirildi: intermediate notification biti 0x10 (bit 4),
	// 0x20 (bit 5) EMAS — v5.0 va cloudhopper issue #54 bo'yicha.
	if Intermediate != 0x10 {
		t.Errorf("Intermediate = 0x%02X, 0x10 bo'lishi kerak", uint8(Intermediate))
	}
}
