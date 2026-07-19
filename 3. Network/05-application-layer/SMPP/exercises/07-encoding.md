# 7-bob mashqlari: Text encoding

> Avval mashqlarni mustaqil bajaring, keyin yechimlarga qarang. Bob matni: [book/07-encoding.md](../book/07-encoding.md)

---

## Mashq 1. dc va segment hisobi

Ikki xabar uchun data_coding, sm_length (oktetlarda!) va segment sonini aniqlang:

1. `Assalomu alaykum! Kodingiz: 5521`
2. `Ассалому алайкум!`

Har biri uchun hisobni ko'rsating. Keyin `coding.Choose` bilan tekshiring.

## Mashq 2. Bu UCS2'mi?

deliver_sm keldi: data_coding=8, short_message boshlanishi:

```
D0 90 D0 BB D0 BE
```

Abonent esa "telefonda koreyscha belgilar ko'rinyapti" deb shikoyat qilmoqda.

1. dc=8 bo'yicha bu baytlarni UCS2 sifatida o'qing — qaysi belgilar chiqadi?
2. Baytlar naqshiga qarab ASL muammoni tashxislang: yuboruvchi aslida nima qilgan?
3. Xabar aslida qanday matn bo'lishi kerak edi va TO'G'RI baytlari qanday?
4. Bunday bug qaysi kod qatlamida paydo bo'ladi va qanday test uni ushlaydi?

## Mashq 3. Qo'lda packing

"hello" so'zini (GSM kodlari: h=0x68, e=0x65, l=0x6C, l=0x6C, o=0x6F) qo'lda septet-packing qiling: har belgining 7 bitini yozib, oktetlarga taxlang. Natijani `coding.Pack` bilan solishtiring (kutilayotgan javob bobda bor — unga qaramasdan bajaring, keyin tekshiring). Bonus: nega 5 belgi packing'dan keyin ham 5 oktet? Nechanchi belgidan boshlab "tejash" boshlanadi?

## Mashq 4. 160 + 1 belgi

Marketing bo'limi aynan 160 belgilik (barchasi GSM7) matn tayyorladi — 1 segment. Keyin matnga "soʻm" so'zini qo'shishdi (matn 161 belgi bo'ldi, ichida bitta U+02BB).

1. Normalizatsiyasiz yuborilsa nechta segment ketadi? Hisob bilan.
2. Normalize qilinsa-chi?
3. Farqni pulda ifodalang (segment narxi N deb).
4. Bu stsenariy uchun qanday himoya mexanizmi (kod/jarayon darajasida) taklif qilasiz?

---

# Yechimlar

## Yechim 1

**1. `Assalomu alaykum! Kodingiz: 5521`** — 32 belgi, hammasi ASCII harf/raqam/tinish — GSM7'da bor (apostrof ham yo'q, extension ham yo'q). Natija: **dc=0 (GSM7), sm_length=32** (unpacked: 1 belgi = 1 oktet, extension yo'q), 32 ≤ 160 → **1 segment**.

**2. `Ассалому алайкум!`** — 17 belgi, kirill → GSM7'da yo'q, dc=6 amalda ishlamaydi → **dc=8 (UCS2)**. sm_length = 17 × 2 = **34 oktet** (belgi emas, OKTET!). 17 ≤ 70 → **1 segment**. E'tibor: matn 1-xabardan qisqa bo'lsa ham, "sig'im"ning kichikligi tufayli uzunroq kirill xabarlar tezroq segmentlanadi.

`Choose` bilan: birinchisi `(DCDefault, 32 bayt)`, ikkinchisi `(DCUCS2, 34 bayt)` qaytaradi.

## Yechim 2

**1.** UCS2 (big-endian) o'qish: `D090 D0BB D0BE` — uch unit. U+D090, U+D0BB, U+D0BE — bular Hangul Syllables blokidan (U+AC00–U+D7AF): koreys bo'g'inlari chiqadi. Shikoyat aynan mos!

**2.** Naqsh ko'zga tashlanadi: har ikkinchi bayt `D0` — bu **UTF-8'ning kirill uchun 2-baytlik prefiksi** (kirill U+0400–U+047F diapazoni UTF-8'da `D0/D1 xx` bo'ladi). Yuboruvchi matnni **UTF-8'da baytlab, dc=8 deb e'lon qilgan** — ya'ni "UCS2 deganda UTF-8 yuborish" bug'i. Go'da tipik ildiz: `[]byte(s)` (bu UTF-8!) ni to'g'ridan-to'g'ri short_message'ga qo'yish.

**3.** `D0 90 D0 BB D0 BE` UTF-8 sifatida = "Ало" boshlanishi (U+0410 А, U+043B л, U+043E о). To'g'ri UCS2: **`04 10 04 3B 04 3E`**. Kirill uchun UCS2 baytlarining "imzosi" — har juftlikning birinchisi `04` bo'lishi.

**4.** Qatlam: matn → baytlar konvertatsiyasi (bizda `coding.EncodeUCS2` — `utf16.Encode` orqali; bug esa uni chetlab `[]byte(s)` ishlatganda). Ushlaydigan test — golden hex: `TestUCS2CyrillicGolden` aynan "Салом" ning baytlarini `04 21 ...` deb qotiradi; UTF-8 yuborilsa test yiqiladi. Umumiy saboq: encoding funksiyalarini "string kirdi — bayt chiqdi" deb emas, ANIQ baytlargacha testlash.

## Yechim 3

Bit taxlash (har septet 7 bit, quyi bitdan boshlab teriladi):

```
h = 1101000   e = 1100101   l = 1101100   l = 1101100   o = 1101111

oktet 0 = e[0] + h[6..0]   = 1|1101000 = 0xE8
oktet 1 = l[1..0] + e[6..1] = 00|110010 = 0x32
oktet 2 = l[2..0] + l[6..2] = 100|11011 = 0x9B
oktet 3 = o[3..0] + l[6..3] = 1111|1101 = 0xFD
oktet 4 = 00000 + o[6..4]   = 00000|110 = 0x06
```

Natija: **E8 32 9B FD 06** — `coding.Pack(coding.EncodeGSM7("hello"))` xuddi shuni beradi (`TestPackHelloGolden`).

Bonus: 5 belgi × 7 bit = 35 bit → 5 oktetga (40 bit) sig'adi, 4 oktetga (32 bit) sig'MAYdi — tejash hali ko'rinmaydi. Formula: n belgi → ceil(7n/8) oktet. n=8'da birinchi marta 7 oktet chiqadi — **8-belgidan boshlab har 8 belgi 1 oktet tejaydi**.

## Yechim 4

**1. Normalizatsiyasiz:** bitta U+02BB butun xabarni UCS2'ga tushiradi. 161 belgi > 70 → concatenation, segment sig'imi 67 (UDH bilan): ceil(161/67) = **3 segment**.

**2. Normalize bilan:** ʻ → ' bo'lib, 161 belgining hammasi GSM7'da qoladi. 161 > 160 → baribir concat, lekin segment sig'imi 153: ceil(161/153) = **2 segment**. (Matnni 160'ga qisqartirish esa 1 segmentga qaytaradi — "soʻm" o'rniga qisqaroq so'z topish yana ham arzon!)

**3.** Narxda: normalizatsiyasiz 3N, normalize bilan 2N, matn 160'da ushlab qolinsa 1N. Ya'ni bitta belgi va bitta funksiya chaqiruvi orasidagi farq — har xabarda **2N gacha**. Millionlab notification'da bu jiddiy byudjet.

**4.** Himoya qatlamlari: (a) yuborish yo'lida MAJBURIY `Normalize` (bizning `Choose` buni o'zi qiladi); (b) segment soni oshganda ogohlantirish — masalan API javobida "3 segment ketadi" ko'rsatish yoki limitdan oshganda alert; (c) marketing matnlari uchun oldindan tekshiruvchi tool (matn kiritilganda dc/segment jonli ko'rsatiladi); (d) monitoring: o'rtacha segment/xabar metrikasi keskin oshsa — kimdir "chiroyli belgi" qo'shgan (16-bob).
