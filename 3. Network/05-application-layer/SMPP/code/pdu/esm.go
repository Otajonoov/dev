package pdu

// esm_class va registered_delivery baytlarining bit helper'lari (§5.2.12, §5.2.17).
// Ikkala bayt ham "bit to'plami" — ularni son sifatida solishtirish (masalan
// esm_class == 0x04) UDHI kabi flag qo'shilganda sinadi; faqat mask/shift ishlatiladi.

// EsmClass — esm_class bayti (§5.2.12).
//
// ESME→SMSC (submit_sm, submit_multi, data_sm):
//
//	bit 1-0: Messaging Mode (00 default, 01 datagram, 10 forward/transaction, 11 store-and-forward)
//	bit 5-2: Message Type (0000 oddiy, 0010 ESME Delivery Ack, 0100 ESME Manual Ack)
//	bit 7-6: GSM (01 UDHI, 10 Reply Path)
//
// SMSC→ESME (deliver_sm, data_sm):
//
//	bit 1-0: e'tiborga olinmaydi
//	bit 5-2: 0000 MO xabar, 0001 DLR, 0010 SME Delivery Ack, 0100 SME Manual Ack,
//	         0110 Conversation Abort, 1000 Intermediate Notification
//	bit 7-6: xuddi shu GSM flag'lar
type EsmClass uint8

// Messaging Mode qiymatlari (bit 1-0).
const (
	ModeDefault         EsmClass = 0x00
	ModeDatagram        EsmClass = 0x01
	ModeForward         EsmClass = 0x02 // transaction; submit_sm buni QO'LLAMAYDI (§4.4)
	ModeStoreAndForward EsmClass = 0x03
)

// GSM feature flag'lari (bit 7-6).
const (
	FlagUDHI      EsmClass = 0x40 // short_message UDH bilan boshlanadi (8-bob)
	FlagReplyPath EsmClass = 0x80
)

// MessageType — esm_class'ning bit 5-2 qiymati (SMSC→ESME talqini).
type MessageType uint8

const (
	TypeNormal            MessageType = 0x00 // oddiy xabar (deliver'da: MO)
	TypeDeliveryReceipt   MessageType = 0x01 // SMSC Delivery Receipt (DLR)
	TypeSMEDeliveryAck    MessageType = 0x02
	TypeSMEManualAck      MessageType = 0x04
	TypeConversationAbort MessageType = 0x06 // Korean CDMA
	TypeIntermediate      MessageType = 0x08 // Intermediate Delivery Notification
)

// Mode messaging mode bitlarini qaytaradi (bit 1-0).
func (e EsmClass) Mode() EsmClass { return e & 0x03 }

// MessageType bit 5-2'ni ajratadi. DIQQAT: aynan shu shift+mask to'g'ri usul;
// esm_class'ni butun bayt sifatida 0x04 bilan solishtirish UDHI (0x40)
// qo'shilib kelganda (0x44) DLR'ni o'tkazib yuboradi.
func (e EsmClass) MessageType() MessageType { return MessageType((e >> 2) & 0x0F) }

// IsDeliveryReceipt — deliver_sm DLR'mi (bit 5-2 == 0001).
func (e EsmClass) IsDeliveryReceipt() bool { return e.MessageType() == TypeDeliveryReceipt }

// IsIntermediate — intermediate delivery notification'mi (bit 5-2 == 1000).
func (e EsmClass) IsIntermediate() bool { return e.MessageType() == TypeIntermediate }

// HasUDHI — short_message UDH bilan boshlanadimi (bit 6).
func (e EsmClass) HasUDHI() bool { return e&FlagUDHI != 0 }

// WithUDHI UDHI flag o'rnatilgan nusxa qaytaradi.
func (e EsmClass) WithUDHI() EsmClass { return e | FlagUDHI }

// RegisteredDelivery — registered_delivery bayti (§5.2.17).
//
//	bit 1-0: SMSC Delivery Receipt (00 yo'q, 01 final holatda — muvaffaqiyat ham,
//	         xato ham, 10 faqat xato final holatda, 11 reserved)
//	bit 3-2: SME acknowledgement (00 yo'q, 01 Delivery Ack, 10 Manual/User Ack, 11 ikkalasi)
//	bit 4:   Intermediate Notification. DIQQAT — spec matni "bit 5" deydi, lekin
//	         jadval patterni xxx1xxxx = bit 4 (0x10): v3.4'ning O'Z ichki erratum'i.
//	         v5.0 va implementatsiyalar (cloudhopper issue #54) bo'yicha to'g'risi 0x10.
type RegisteredDelivery uint8

const (
	DLRNone        RegisteredDelivery = 0x00
	DLRFinal       RegisteredDelivery = 0x01 // muvaffaqiyat YOKI xato yakunida DLR
	DLRFailureOnly RegisteredDelivery = 0x02
	SMEDeliveryAck RegisteredDelivery = 0x04
	SMEManualAck   RegisteredDelivery = 0x08
	Intermediate   RegisteredDelivery = 0x10 // bit 4 — erratum izohi yuqorida
)

// DLRRequest bit 1-0'ni qaytaradi (so'ralgan DLR rejimi).
func (r RegisteredDelivery) DLRRequest() RegisteredDelivery { return r & 0x03 }

// WantsDLR biror final-holat DLR so'ralganmi (bit 1-0 != 00; 11 reserved
// bo'lgani uchun ham "so'ralgan" deb qaraladi — liberal qabul).
func (r RegisteredDelivery) WantsDLR() bool { return r&0x03 != 0 }

// WantsIntermediate — intermediate notification so'ralganmi (bit 4).
func (r RegisteredDelivery) WantsIntermediate() bool { return r&0x10 != 0 }
