# 5-bob mashqlari: submit_sm va deliver_sm

> Avval mashqlarni mustaqil bajaring, keyin yechimlarga qarang. Bob matni: [book/05-submit-deliver.md](../book/05-submit-deliver.md)

---

## Mashq 1. Uch xato yashiringan submit_sm

Quyidagi frame'da UCHTA xato yashiringan (biri header'da, biri vaqt field'ida, biri uzunlikda):

```
00 00 00 40 00 00 00 04 00 00 00 05 00 00 00 02
00 05 00 42 61 6E 6B 00 01 01 39 39 38 39 30 31
32 33 34 35 36 37 00 00 00 00 00 30 30 30 30 30
31 30 30 30 30 00 01 00 00 00 07 53 61 6C 6F 6D
```

1. Frame'ni qo'lda, field-ma-field parse qilib uchala xatoni toping.
2. Bizning `DecodeSubmitSM` bu frame'da qaysi xatoni USHLAYDI va qaysilarini o'tkazib yuboradi? Nega?
3. Har xato uchun: to'g'ri server bunga qanday munosabatda bo'lishi kerak (qaysi command_status)?
4. Har xato qanday client-bug'dan kelib chiqqan bo'lishi mumkin?

## Mashq 2. registered_delivery = 0x11

submit_sm'da `registered_delivery = 0x11` yuborildi.

1. Baytni bit-guruhlarga ajratib yozing: qaysi bitlar o'rnatilgan, har guruh nimani so'rayapti?
2. SMSC bunga javoban qanday deliver_sm'lar yuborishi mumkin — har birining esm_class bit 5–2 qiymati bilan sanang.
3. `pdu` package helper'lari bilan tekshiring: `RegisteredDelivery(0x11)` uchun `WantsDLR()`, `WantsIntermediate()`, `DLRRequest()` nima qaytaradi?
4. Xuddi shu savol 0x21 uchun — nima o'zgaradi va bu qiymat qayerda adashuvdan kelib chiqishi mumkin?

## Mashq 3. "SMS yetkazildimi?"

Sizning gateway submit_sm yubordi va submit_sm_resp status=0, message_id="A81F03" oldi. Mahsulot menejeri so'rayapti: "Demak mijozga SMS bordi-da?"

To'liq, texnik asoslangan javob yozing: (1) status=0 aynan nimani KAFOLATLAYDI va nimani KAFOLATLAMAYDI; (2) xabar hali qaysi taqdirlarga uchrashi mumkin (kamida 4 stsenariy, message_state qiymatlari bilan); (3) "yetkazildi" deyish uchun qanday dalil kerak va u qachon/qanday keladi; (4) DLR umuman kelmasa-chi?

---

# Yechimlar

## Yechim 1

**Parse va uch xato:**

| Field | Qiymat | Baho |
|---|---|---|
| command_length | 0x40 = 64 | to'g'ri (frame haqiqatan 64 oktet) |
| command_id | 0x00000004 submit_sm | to'g'ri |
| **command_status** | **0x00000005** | **XATO 1: request'da command_status NULL bo'lishi SHART (§5.1.3)** — 0x05 (RALYBND qiymatiga to'g'ri kelib qoladi) bu yerda ma'nosiz |
| sequence | 2 | to'g'ri |
| service_type.."Bank"..dest | — | to'g'ri (goldendagi bilan bir xil) |
| esm/protocol/priority | 00 00 00 | to'g'ri |
| schedule | 00 (bo'sh) | to'g'ri |
| **validity_period** | **"0000010000" (10 belgi)** | **XATO 2: "1 or 17" qoidasi (§7.1.1)** — yo bo'sh, yo AYNAN 16 belgi + NULL; 10 belgili qiymat format buzilishi |
| reg/replace/dc/sm_def | 01 00 00 00 | to'g'ri |
| **sm_length** | **0x07** | **XATO 3: frame'da short_message uchun faqat 5 oktet ("Salom") qolgan** — sm_length yolg'on gapiryapti |

**2. Bizning decoder:** birinchi ikkita xatoni parse darajasida USHLAMAYDI: command_status'ni `DecodeSubmitSM` tekshirmaydi (header sintaktik to'g'ri — semantik tekshiruv dispatcher/server qatlamining ishi, 10/14-boblar); validity ham shunchaki C-Octet String sifatida o'qiladi (10 belgili string sintaktik valid — u SEMANTIK xato, `ParseTime` chaqirilganda yoki encode'da `writeTimeField`'da ushlanadi). Uchinchi xato esa aniq parse xatosi: `sm_length=7, lekin frame'da 5 oktet qoldi` — TLV tail'ga "kirib ketishdan" himoya ishlaydi. Bu qatlamlanish ataylab: **sintaksis codec'da, semantika yuqorida**.

**3. Server munosabati:** request'dagi nolsiz status — spec bo'yicha "should be NULL"; qattiq server ESME_RINVCMDID emas... to'g'risi: PDU'ni qabul qilib, lekin qoida buzilishiga ko'ra rad etsa ESME_RSYSERR (0x08) yoki oddiy ignore (ko'p SMSC'lar shunchaki e'tibor bermaydi — eng keng tarqalgan xulq). Validity xatosi → **ESME_RINVEXPIRY (0x62)**. sm_length nomuvofiqligi → **ESME_RINVMSGLEN (0x01)**. Uchala javob ham submit_sm_resp ichida keladi (generic_nack EMAS — header butun!).

**4. Tipik ildizlar:** status≠0 — header struct'ini nusxalashda resp'dan qolgan qiymatni tozalamaslik; validity — "YYMMDDhhmm" (10 belgili DLR-sana formati!) bilan "YYMMDDhhmmsstnnp" (16 belgili §7.1.1) ni adashtirish; sm_length — matnni keyin qisqartirib, uzunlik field'ini yangilashni unutish (yoki belgi soni bilan bayt sonini adashtirish — UCS2'da klassika, 7-bob).

## Yechim 2

**1.** 0x11 = `0001 0001`:

| Bitlar | Qiymat | So'rov |
|---|---|---|
| bit 1–0 (SMSC DLR) | `01` | final holat uchun DLR — muvaffaqiyat HAM, xato HAM |
| bit 3–2 (SME ack) | `00` | SME ack so'ralmagan |
| bit 4 (Intermediate) | `1` | oraliq notification'lar HAM so'ralgan (erratum esda: bit 4 = 0x10, "bit 5" emas!) |

**2.** Kelishi mumkin: oraliq holatlarda **Intermediate Delivery Notification** — esm_class bit 5–2 = `1000` (0x20) (masalan birinchi urinish muvaffaqiyatsiz, xabar hali SMSC'da); yakunda **SMSC Delivery Receipt** — esm_class bit 5–2 = `0001` (0x04). SME ack'lar kelmaydi (so'ralmagan). Diqqat: intermediate support SMSC ixtiyorida — kelmasligi ham normal.

**3.**

```go
rd := pdu.RegisteredDelivery(0x11)
rd.WantsDLR()          // true  (bit 1-0 = 01)
rd.WantsIntermediate() // true  (bit 4)
rd.DLRRequest()        // 0x01 = pdu.DLRFinal
```

**4.** 0x21 = `0010 0001`: bit 1–0 = 01 (DLR bor), bit 4 = 0 (intermediate YO'Q), bit 5 = 1 — v3.4'da bit 5 registered_delivery'da **reserved**. `WantsIntermediate()` false qaytaradi. Bu qiymat qayerdan chiqadi? Aynan bit 4/5 erratum'idan: spec matnidagi "bit 5"ni so'zma-so'z olgan kod intermediate uchun 0x20 qo'yadi (cloudhopper'ning eski konstantasi shu xatoni qilgan — issue #54). Natija: intermediate so'ralmagan bo'lib chiqadi, va hech qanday xato ham ko'rinmaydi — jim degradatsiya.

## Yechim 3

**(1) status=0 kafolatlari:** SMSC PDU'ni sintaktik qabul qildi, xabarni O'Z NAVBATIGA oldi va unga message_id ajratdi ("A81F03" — opaque, saqlab qo'yamiz). KAFOLATLAMAYDI: telefonga yetganini, hatto yetkazish urinishi BOSHLANGANINI ham. Xabar hozir SMSC store'ida ENROUTE (message_state=1) holatida.

**(2) Mumkin taqdirlar:** DELIVERED (2) — yetdi; EXPIRED (3) — telefon o'chiq bo'lib validity tugadi; UNDELIVERABLE (5) — raqam mavjud emas/bloklangan; REJECTED (8) — SMSC/operator filtri keyinroq rad etdi (masalan spam-filtr, sender ro'yxatdan o'tmagan); DELETED (4) — operator tozaladi. Aggregator zanjirida bularga "keyingi hop rad etdi" varianti ham qo'shiladi — u ham oxir-oqibat UNDELIV/REJECTD DLR bo'lib keladi.

**(3) "Yetkazildi" dalili:** esm_class'ida DLR biti (bit 5–2 = 0001) o'rnatilgan deliver_sm, ichida `stat:DELIVRD` (yoki message_state TLV = 2) va `id:` (yoki receipted_message_id TLV) bizning "A81F03" bilan bog'lanadigan qiymat. Kelish vaqti — sekundlardan (telefon yoniq, tarmoq band emas) soatlargacha (validity oxirigacha qayta urinishlar). Buning uchun submit'da registered_delivery=0x01 so'ralgan bo'lishi shart edi!

**(4) DLR kelmasa:** bu "yetmadi" degani EMAS — DLR yo'qolishi mumkin (zanjir uzun), operator DLR route'i buzilgan bo'lishi mumkin, yoki xabar haqiqatan ENROUTE'da osilib turibdi. Amaliy siyosat: timeout'dan keyin (masalan validity + zaxira) xabarni "holati noma'lum" deb belgilash; ba'zi hollarda query_sm bilan so'rab ko'rish mumkin (10-bob), lekin polling asosiy mexanizm emas; monitoring darajasida esa "DLR kelish ulushi" metrikasi route sog'lig'ining asosiy indikatori (16-bob). Mahsulot menejeriga qisqa javob: **"SMSC qabul qildi; yetkazilganini DLR tasdiqlaydi — u kelguncha 'yo'lda' deb hisoblaymiz."**
