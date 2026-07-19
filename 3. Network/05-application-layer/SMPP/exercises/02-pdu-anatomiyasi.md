# 2-bob mashqlari: PDU anatomiyasi

> Avval mashqlarni mustaqil bajaring, keyin yechimlarga qarang. Bob matni: [book/02-pdu-anatomiyasi.md](../book/02-pdu-anatomiyasi.md)

---

## Mashq 1. Aralash stream'ni qo'lda decode qilish

TCP stream'dan quyidagi baytlar ketma-ket keldi (bo'shliqlar faqat o'qish qulayligi uchun):

```
00 00 00 10 00 00 00 15 00 00 00 00 00 00 00 1C
00 00 00 17 80 00 00 04 00 00 00 00 00 00 00 1B
37 46 33 41 39 42 00
```

Hech qanday tool'siz, faqat qog'oz-qalam bilan:

1. Stream'da nechta PDU bor? Chegaralarini qanday topdingiz?
2. Har PDU uchun: turi (nomi bilan), request'mi/response'mi, command_status, sequence_number.
3. Body bor PDU'larning body'sini field'ma-field o'qing.
4. Bu stream'ga qarab peer (qarshi tomon)dan qanday harakat kutiladi?
5. Sequence raqamlariga qarang: 0x1C keyin 0x1B keldi — bu xatomi?

## Mashq 2. Buzilgan PDU

Stream'dan `command_length = 0x0000000C` bilan boshlanadigan frame keldi.

1. Bu frame nimasi bilan xato? Aniq qoidani ayting.
2. Bizning `ReadFrame` bunga qanday javob beradi (qaysi error)?
3. Spec bo'yicha protokol darajasida qanday javob PDU yuborilishi kerak va undan keyin session bilan nima qilgan ma'qul?
4. Xuddi shu savol `command_length = 0x00FFFFFF` (16 MB) uchun — bu qiymat "qonuniy" ko'rinadi-ku?

## Mashq 3. "1 or 17" encoder

v3.4 §3.1'dagi "Fixed 1 or 17" qoidasini implement qiling: vaqt field'i (masalan `schedule_delivery_time`) yo bitta 0x00 okteti (qiymat berilmagan), yo aynan 17 oktet (16 belgi + NULL) bo'ladi — oraliq yo'q.

`writeTimeField(b *bytes.Buffer, s string) error` funksiyasini yozing:

- `s == ""` → bitta 0x00 yoziladi;
- `len(s) == 16` → 16 belgi + NULL yoziladi;
- boshqa har qanday uzunlik → xato.

Va unga table-driven test yozing: bo'sh string (1 oktet), to'g'ri 16 belgili string (17 oktet), 12 belgili string (xato), 17 belgili string (xato). (Belgilar MAZMUNINI — YYMMDDhhmmsstnnp format qoidalarini — tekshirish shart emas, u 5-bobda `pdu/time.go`'da qilinadi.)

---

# Yechimlar

## Yechim 1

**1. Ikkita PDU.** Chegara faqat command_length orqali topiladi: birinchi 4 oktet `00 00 00 10` = 16 → birinchi PDU 16 oktet (header-only), stream'ning 16-oktetidan keyingi 4 oktet `00 00 00 17` = 23 → ikkinchi PDU 23 oktet. 16 + 23 = 39 — stream'dagi baytlar soni bilan aynan mos.

**2–3. PDU'lar:**

| | 1-PDU | 2-PDU |
|---|---|---|
| command_length | 0x10 = 16 (body yo'q) | 0x17 = 23 (body 7 oktet) |
| command_id | 0x00000015 = **enquire_link** | 0x80000004 = **submit_sm_resp** (bit 31 → response) |
| command_status | 0 (request'da majburiy NULL) | 0 = ESME_ROK — muvaffaqiyat |
| sequence_number | 0x1C = 28 | 0x1B = 27 |
| body | — | `37 46 33 41 39 42 00` = C-Octet String **"7F3A9B"** + NULL = message_id |

message_id "7F3A9B" hex ko'rinishli — SMSC message_id'ni qanday formatda berishi operatorga bog'liq (opaque qiymat); bu 9-bobdagi hex/decimal korrelyatsiya muammosining darak beruvchisi.

**4. Kutiladigan harakat:** enquire_link — request, unga **darhol** enquire_link_resp (command_status=0, sequence=28!) qaytarish shart, aks holda peer sessiyani o'lik deb topadi. submit_sm_resp esa o'zi javob — unga javob berilmaydi, faqat sequence=27 bo'yicha kutayotgan submit_sm bilan korrelyatsiya qilinadi.

**5. Xato emas.** Bular ikki MUSTAQIL narsa: 27 — bizning avvalgi request'imizga javob (peer'ning javobi qachon kelishi uning ixtiyorida), 28 — peer'ning O'Z request'i (peer o'z sequence fazosini yuritadi — aslida bu misolda 28 bizniki bo'lishi ham mumkin emas edi, chunki enquire_link peer'dan keldi; peer'ning sequence'lari bizning raqamlarimizdan butunlay mustaqil). Async protokolda javoblar tartibi ham kafolatlanmagan (v3.4 §2.5–2.7): har ikki tomon out-of-order'ni qabul qila olishi SHART.

## Yechim 2

1. `command_length = 12 < 16` — PDU o'z header'idan ham kichik bo'la olmaydi (header majburiy 16 oktet, v3.4 §3.2). Bunday qiymat "kichik PDU" emas, **buzilgan stream** belgisi.
2. `ReadFrame` allocation qilmasdan `ErrFrameTooShort` qaytaradi (frame_test.go'dagi `TestReadFrameLengthTooShort` aynan shu holat).
3. Spec bo'yicha: header'i buzuq PDU'ga **generic_nack** qaytariladi (v3.4 §2.8), command_status = ESME_RINVCMDLEN (0x02 — "Invalid Command Length"). Lekin amaliy davomi muhim: length ishonchsiz bo'lsa, KEYINGI frame chegarasini ham bilmaymiz — stream'dagi qolgan hamma narsa shubhali. Eng xavfsiz strategiya: generic_nack yuborib, sessiyani yopish va qayta ulanish (11-bobda batafsil).
4. 0x00FFFFFF formal jihatdan "katta PDU" xolos — lekin real PDU'lar bir necha yuz oktet, eng kattasi ham (message_payload bilan) 64 KB atrofida (v3.4 §3.2.3). 16 MB so'ragan frame yo buzilgan stream, yo hujum. `ReadFrame`'da bu maxSize bilan kesiladi → `ErrFrameTooLarge`, xotira ajratilMAYdi. Himoyasiz implementatsiya esa `make([]byte, 16<<20)` qilib, bunday frame'lar seriyasidan OOM bo'ladi. Javob PDU darajasida yana generic_nack (ESME_RINVCMDLEN) + reconnect.

## Yechim 3

```go
package main

import (
	"bytes"
	"fmt"
)

// writeTimeField "Fixed 1 or 17" field'ini yozadi (v3.4 §3.1):
// bo'sh qiymat = yagona NULL okteti; aks holda aynan 16 belgi + NULL terminator.
func writeTimeField(b *bytes.Buffer, s string) error {
	if s == "" {
		b.WriteByte(0x00)
		return nil
	}
	if len(s) != 16 {
		return fmt.Errorf("vaqt field'i aynan 16 belgi bo'lishi kerak, keldi: %d", len(s))
	}
	b.WriteString(s)
	b.WriteByte(0x00)
	return nil
}
```

Table-driven test:

```go
func TestWriteTimeField(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		wantLen int  // yozilgan oktetlar; -1 = xato kutiladi
	}{
		{"bo'sh qiymat", "", 1},
		{"to'liq absolute vaqt", "020610233429000R", 17},
		{"juda qisqa", "0206102334", -1},
		{"juda uzun", "020610233429000R0", -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b bytes.Buffer
			err := writeTimeField(&b, tt.in)
			if tt.wantLen == -1 {
				if err == nil {
					t.Fatalf("xato kutilgan edi, nil keldi (yozildi: % X)", b.Bytes())
				}
				return
			}
			if err != nil {
				t.Fatalf("kutilmagan xato: %v", err)
			}
			if b.Len() != tt.wantLen {
				t.Errorf("%d oktet yozildi, kutilgan %d", b.Len(), tt.wantLen)
			}
		})
	}
}
```

Ikki nozik nuqta: (1) "1 or 17" dagi 17 — NULL BILAN (§3.1 note iii bu yerda ham amal qiladi): 16 belgi + terminator; (2) bo'sh holatda `writeCString(b, "", 17)` chaqirsangiz ham xuddi shu natija (bitta 0x00) chiqadi — ya'ni "1 or 17" aslida C-Octet String'ning maxsus holati, alohida funksiya faqat "oraliq uzunlik taqiqlangan" qoidasini qo'shadi.
