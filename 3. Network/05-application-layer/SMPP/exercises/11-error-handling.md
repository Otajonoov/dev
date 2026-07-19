# 11-bob mashqlari: Error handling

> Avval mashqlarni mustaqil bajaring, keyin yechimlarga qarang. Bob matni: [book/11-error-handling.md](../book/11-error-handling.md)

---

## Mashq 1. Log'dagi raqam

Uch xil gateway log'idan uch qator:

```
(a) submit failed: command_status=20
(b) submit failed: command_status=0x0A
(c) DLR: stat:UNDELIV err:020
```

1. (a) va (b) qaysi kodlar? Har birining toifasi (Classify) va to'g'ri reaksiyasi?
2. (a)'ni "0x20" deb o'qigan operator qanday xato tashxis qo'yadi?
3. (c)'dagi 020 bilan (a)'dagi 20 bir xil narsami? Nega?
4. `pdu.CommandStatus(20).String()` va `client.Classify(20)` nima qaytaradi — kod bilan tekshiring.

## Mashq 2. RTHROTTLED zanjiri

Yuk testida gateway'ingiz 500 mps'ga chiqdi; operator limiti 100 mps. SMSC RTHROTTLED qaytara boshladi. Kutubxonangiz esa xatoli submit'ni DARHOL qayta yuboradi.

1. Keyingi 60 soniyada nima bo'lishini bosqichma-bosqich yozing (SMSC'ning "javobsiz tashlash" xulqini ham hisobga oling).
2. Zanjirning qaysi nuqtasida sessiya "o'ladi" va nega reconnect vaziyatni YOMONLASHTIRADI?
3. To'g'ri arxitektura bu zanjirni qaysi UCH nuqtada uzadi?
4. Window=300 bo'lsa empirik qoida bo'yicha uni nechaga tushirish kerak?

## Mashq 3. Buzuq PDU'ga javob

SMSC sifatida ishlayapsiz (14-bob tayyorgarligi). Stream'dan 4 bayt keldi: `00 00 00 0C` — va undan keyin yana 8 bayt: `00 00 00 04 00 00 00 2A`.

1. Bu frame nimasi bilan buzuq? (2-bob qoidasi.)
2. Javob PDU'sini BAYT-MA-BAYT hex'da yig'ing: qaysi command_id, qaysi command_status, qaysi sequence_number va NEGA?
3. Javobni yuborib bo'lgach sessiya bilan nima qilasiz?
4. Xuddi shu savol, lekin endi to'liq 16-baytlik header keldi: `00 00 00 10 00 00 00 77 00 00 00 00 00 00 00 05` — javob nima bo'ladi?

---

# Yechimlar

## Yechim 1

**1.** (a) decimal 20 = 0x14 = **ESME_RMSGQFUL** — Transient: backoff bilan retry, max-age oynasi. (b) 0x0A = **ESME_RINVSRCADR** — Permanent: retry taqiqlangan, source konfiguratsiyasini tekshirish kerak.

**2.** 0x20 Table 5-2'da yo'q (reserved) — operator "notanish/vendor xato" deb adashadi va haqiqiy muammoni (SMSC navbati to'la — ehtimol o'zi juda tez yuboryapti) ko'rmaydi.

**3.** YO'Q — ikki fazo: (a) SMPP command_status (Table 5-2, universal), (c) DLR `err:` (operator-specific, universal jadvali yo'q). `err:020` ning ma'nosi faqat shu operatorning hujjatida.

**4.** `ESME_RMSGQFUL` va `ClassTransient`:

```go
fmt.Println(pdu.CommandStatus(20))          // ESME_RMSGQFUL
fmt.Println(client.Classify(pdu.CommandStatus(20))) // transient
```

## Yechim 2

**1.** Zanjir: (i) 100 mps'dan oshgan har submit RTHROTTLED oladi; (ii) darhol-retry bu xabarlarni oqimga qaytaradi — endi SMSC'ga 500 original + 400 retry = yanada ko'proq PDU boradi; (iii) SMSC himoyaga o'tadi: PDU'larni JAVOBSIZ tashlay boshlaydi; (iv) javobsiz submit'lar response timeout'gacha pending window'da osilib turadi — window to'ladi, yangi submit bloklanadi; (v) enquire_link ham javobsiz qolishi mumkin → kutubxona "sessiya o'ldi" deb topadi.

**2.** (v)-nuqtada — enquire_link/response timeout'lar ketma-ket otganda. Reconnect yomonlashtiradi, chunki: queue endi ham original, ham retry xabarlar bilan to'la; yangi sessiya ochilishi bilan hammasi birdan otiladi — spiral yangi kuch bilan aylanadi (+ tez-tez rebind'ning o'zi operator limitlariga uriladi).

**3.** (1) Rate limiter submit'dan OLDIN — 100 mps'dan oshirmaydi (token bucket, 13-bob); (2) RTHROTTLED'ga Classify → Transient → queue'ga qaytarish + backoff (darhol emas); (3) window'ni operator mps'iga moslash — batch oxirida throttling to'planmaydi.

**4.** Empirik qoida window ≤ 2–3 × mps limiti: 100 mps uchun **200–300** — ya'ni 300 chegarada; lekin muammo window'da emas, RATE'da edi — limiter qo'yilgach window 300 qolsa ham bo'ladi (u endi hech qachon to'lmaydi).

## Yechim 3

**1.** command_length = 0x0C = 12 < 16 — header'ning o'zidan kichik (2-bob: minimal PDU = 16 bayt header). Frame yaroqsiz, undan keyingi baytlar chegarasi ishonchsiz.

**2.** generic_nack, RINVCMDLEN, **seq=0** — chunki "12 baytlik PDU" degan da'vo bilan kelgan oqimda haqiqiy sequence qayerdaligiga ishonch yo'q (keyingi 8 baytni "shu PDU'ning davomi" deb o'qish ham taxmin):

```
00 00 00 10   <- command_length = 16 (faqat header)
80 00 00 00   <- command_id = generic_nack
00 00 00 02   <- command_status = ESME_RINVCMDLEN
00 00 00 00   <- sequence_number = 0 (decode ishonchsiz - §4.3.1)
```

**3.** Yopish tarafga o'tasiz: framing buzilgan oqimda davom etish har keyingi "PDU"ni ham axlat qiladi. Log + TCP close (server sifatida unbind kutish shart emas — stream baribir ishonchsiz).

**4.** Endi header VALID (length=16, seq=5 aniq o'qildi) — muammo faqat notanish command_id 0x77. Javob: generic_nack, status=RINVCMDID (0x03), **seq=5** (endi ishonchli!):

```
00 00 00 10  80 00 00 00  00 00 00 03  00 00 00 05
```

Va bu safar sessiyani UZMAYMIZ — framing sog'lom, shunchaki bitta notanish PDU keldi (ehtimol versiya farqi); forward-compatibility ruhida davom etamiz.
