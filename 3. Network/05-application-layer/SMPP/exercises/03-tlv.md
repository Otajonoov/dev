# 3-bob mashqlari: TLV

> Avval mashqlarni mustaqil bajaring, keyin yechimlarga qarang. Bob matni: [book/03-tlv.md](../book/03-tlv.md)

---

## Mashq 1. TLV tail'ni qo'lda parse qilish

PDU'ning mandatory qismidan keyin quyidagi baytlar qoldi:

```
02 04 00 02 00 2A 14 50 00 04 DE AD BE EF 04 26 00 01 01
```

Qog'oz-qalam bilan:

1. Nechta TLV bor? Har birining Tag / Length / Value'sini ajrating.
2. Har tag'ni nomlang (Table 5-7 yoki bob jadvalidan). Nomsiz chiqqani qaysi blokka tegishli va unga qanday munosabatda bo'lish kerak?
3. Tanish TLV'larning value'larini TALQIN qiling (nima degani?).
4. Agar oxirgi TLV'ning Length'i `00 01` emas, `00 02` bo'lganida `Decode` nima qilardi?

## Mashq 2. message_payload va sm_length=0 qoidasi

v3.4 §3.2.3: matn short_message field'iga (max 254 oktet) YOKI message_payload TLV'siga solinadi — **hech qachon ikkalasiga birga emas**; message_payload ishlatilsa sm_length=0 bo'lishi shart.

`tlv` package'idan foydalanib `validatePayload(smLength int, tlvs []tlv.TLV) error` funksiyasini yozing: qoida buzilsa (sm_length ≠ 0 va tail'da message_payload bor) xato qaytarsin. Unga uch case'li test yozing: faqat short_message (ok), faqat message_payload (ok), ikkalasi birga (xato).

Bonus savol: bu tekshiruvni qaysi tomonda qilish to'g'riroq — encoder'da (yuborishdan oldin) yoki decoder'da (qabul qilganda)? Ikkala javobni asoslang.

## Mashq 3. Vendor TLV design

Operator sizdan har submit_sm'da "kampaniya identifikatori"ni 0x1401 tag'ida yuborishni so'radi, formatini esa sizga qo'yib berdi. Qanday value format tanlaysiz va nega? Kamida uchta variantni solishtiring (masalan: 4-oktet Integer, C-Octet String, xom UTF-8 baytlar) va tanlovingizning: (a) qarshi tomon parseri, (b) Wireshark'da debugging, (c) kelajakda format o'zgarishi nuqtai nazaridan oqibatlarini yozing.

---

# Yechimlar

## Yechim 1

**1–2. Uchta TLV:**

| Baytlar | Tag | Length | Value | Nomi |
|---|---|---|---|---|
| `02 04 00 02 00 2A` | 0x0204 | 2 | `00 2A` | user_message_reference |
| `14 50 00 04 DE AD BE EF` | 0x1450 | 4 | `DE AD BE EF` | nomsiz — **vendor bloki** (0x1400–0x3FFF) |
| `04 26 00 01 01` | 0x0426 | 1 | `01` | more_messages_to_send |

Jami 6 + 8 + 5 = 19 oktet — tail to'liq yopildi, ortiqcha bayt yo'q ✓.

0x1450 — SMSC vendor-specific TLV. Munosabat (§3.3): agar bu operator hujjatida tanish bo'lsa — o'sha hujjat bo'yicha talqin; bo'lmasa — **jim ignore**, lekin codec darajasida saqlab qo'yiladi (bizning `Decode` shunday qiladi; `Tag.String()` uni `vendor(0x1450)` deb ko'rsatadi). Error-level log YOZILMAYDI.

**3. Talqinlar:** user_message_reference = 0x002A = 42 — ESME'ning O'Z ichki reference raqami; SMSC uni ack oqimlarida aynan qaytaradi (masalan SME ack'larda original xabarni ko'rsatish uchun). more_messages_to_send = 1 — "shu manzilga ketma-ket yana xabar(lar) yuboraman" (GSM'da SMSC radio kanalni ochiq ushlab turishi mumkin — samaradorlik optimizatsiyasi, §5.3.2.34).

**4.** Length=2 va'da qilingan, lekin tail'da value uchun bitta `01` okteti qolgan — `Decode` `ErrTruncated` qaytaradi (`more_messages_to_send Length=2, qolgan baytlar 1`). Protokol darajasida bunga ESME_RINVOPTPARSTREAM (0xC0) mos keladi.

## Yechim 2

```go
package main

import (
	"errors"
	"fmt"

	"smpp/tlv"
)

// validatePayload §3.2.3 qoidasi: short_message (sm_length>0) va
// message_payload TLV birga kelishi taqiqlangan.
func validatePayload(smLength int, tlvs []tlv.TLV) error {
	_, hasPayload := tlv.Find(tlvs, tlv.MessagePayload)
	if hasPayload && smLength != 0 {
		return fmt.Errorf("short_message (sm_length=%d) va message_payload birga taqiqlangan (v3.4 §3.2.3)", smLength)
	}
	return nil
}
```

Test:

```go
func TestValidatePayload(t *testing.T) {
	payload := []tlv.TLV{tlv.CString(tlv.MessagePayload, "juda uzun matn...")}

	if err := validatePayload(10, nil); err != nil {
		t.Errorf("faqat short_message: kutilmagan xato %v", err)
	}
	if err := validatePayload(0, payload); err != nil {
		t.Errorf("faqat message_payload: kutilmagan xato %v", err)
	}
	if err := validatePayload(10, payload); err == nil {
		t.Error("ikkalasi birga: xato kutilgan edi")
	}
}
```

(Eslatma: message_payload value'si aslida C-Octet String emas — xom Octet String; misolda qulaylik uchun CString ishlatildi, real kodda `tlv.TLV{Tag: tlv.MessagePayload, Value: encodedText}` bo'ladi — 8-bobda to'g'ri ishlatamiz.)

Bonus: **ikkala tomonda ham**, lekin har xil natija bilan. Encoder'da — o'z xatomizni SMSC'ga yetkazmasdan lokal ushlash uchun (fail fast, ESME_RINVPARLEN/RSYSERR o'rniga aniq Go error). Decoder'da (server/mock SMSC yozganda, 14-bob) — spec'ni enforce qilish uchun: qoidabuzar PDU'ga xato status qaytariladi. Faqat bir tomonga ishonish bo'lmaydi: simning narigi tomoni har doim "begona kod".

## Yechim 3

Uch variant taqqosi:

| Variant | Parser (a) | Wireshark (b) | Evolyutsiya (c) |
|---|---|---|---|
| **4-oktet Integer (big-endian)** | Eng sodda, xato qilish qiyin; qat'iy 4 oktet — Length tekshiruvi trivial | Hex'da ko'rinadi, lekin qiymatni "o'qish" uchun hisoblash kerak | Diapazon tugasa (4 mlrd) format sindiriladi; "keyin string qilamiz" degan yo'l YO'Q |
| **C-Octet String (ASCII, NULL bilan)** | Sal murakkabroq (NULL, max uzunlik kelishish kerak), lekin bizning CStringValue kabi helper'lar bor | Dump'da ODAMga o'qiladi ("CAMP-2026-07" ko'rinib turadi) — debugging'da katta plus | Istalgan yangi format (prefiks, versiya belgisi) matn ichida hal bo'ladi — eng moslashuvchan |
| **Xom UTF-8 baytlar** | Uzunlik Length'dan ma'lum, lekin ASCII bo'lmagan belgilar telekom tizimlarida kutilmagan joyda sinadi | O'qiladi-yu, encoding taxmin qilinadi | Moslashuvchan, lekin "bu UTF-8" degan bilim faqat og'zaki kelishuvda yashaydi |

**Tavsiya etiladigan tanlov: C-Octet String.** Sabablari: kampaniya id'si tabiatan identifikator (raqam emas — arifmetika qilinmaydi), string har ikki tomonda intuitiv parse bo'ladi, log/dump'larda ko'rinadi va format evolyutsiyasiga chidamli. Qo'shimcha design qoidalari (qaysi formatni tanlamang): (1) hujjatlashtiring — tag, format, max uzunlik, misol dump; (2) bitta tag = bitta format, "ba'zan int, ba'zan string" MUMKIN EMAS; (3) qarshi tomon bu TLV'ni umuman yubormasligi/ignore qilishi mumkinligini nazarda tuting (§3.3 — u optional bo'lib qoladi); (4) 0x1400–0x3FFF blokidan tashqariga chiqmang — SMPP-defined diapazonda "o'z" tag e'lon qilish kelajakdagi spec/kutubxona to'qnashuviga yo'l ochadi.
