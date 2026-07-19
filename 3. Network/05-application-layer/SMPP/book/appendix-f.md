# Appendix F. Glossary — atamalar lug'ati

> English termin → bir qatorlik o'zbekcha izoh. Termin birinchi chuqur yoritilgan bob qavsda.

## Ekotizim va rollar

| Termin | Izoh |
|---|---|
| **SMSC** (Short Message Service Centre) | Operatorning SMS markazi: xabarlarni qabul qiladi, saqlaydi, yetkazadi (1) |
| **ESME** (External Short Message Entity) | Mobil tarmoq TASHQARISIDAGI SMPP mijozi — sizning gateway (1) |
| **SME** (Short Message Entity) | Xabar yubora/ola biladigan har qanday uchastka (telefon ham, server ham) (1) |
| **MC** (Message Centre) | v5.0'da SMSC'ning umumlashgan nomi (16) |
| **Aggregator** | Operatorlar bilan sizning orangizdagi vositachi: bir tomonga SMSC, ikkinchisiga ESME (1) |
| **MO / MT** (Mobile Originated/Terminated) | Telefon YUBORGAN / telefonGA boradigan xabar (1, 5) |
| **A2P / P2P** | Application-to-Person / Person-to-Person traffic turlari (5) |
| **MSISDN** | Abonentning to'liq telefon raqami (xalqaro formatda) (6) |
| **SS7 / MAP** | Operator ichki signalizatsiya tarmog'i va protokoli — SMPP undan izolyatsiya qiladi (1) |
| **HLR / MSC** | Abonent registri / kommutatsiya markazi — MT yetkazish zanjiri qatnashchilari (1) |

## Protokol asoslari

| Termin | Izoh |
|---|---|
| **PDU** (Protocol Data Unit) | SMPP'ning bitta xabar birligi: 16-baytlik header + body (2) |
| **command_id** | PDU turi (submit_sm, deliver_sm...); bit 31 = response belgisi (2, A) |
| **command_status** | Resp'dagi natija kodi; request'da doim 0 (2, 11, B) |
| **sequence_number** | Request↔resp bog'lovchi raqam, 1..0x7FFFFFFF (2, 12) |
| **frame / framing** | TCP oqimidan bitta to'liq PDU'ni ajratish (length-prefix + io.ReadFull) (2) |
| **big-endian / network byte order** | Katta bayt oldin — SMPP'ning barcha integer'lari shunday (2) |
| **C-Octet String** | NULL (0x00) bilan tugaydigan satr; "max N" NULL'ni O'Z ICHIGA oladi (2) |
| **TLV** (Tag-Length-Value) | Optional parameter formati; Length = faqat Value (3, C) |
| **body / mandatory field** | PDU'ning majburiy, tartibi qat'iy qismi (TLV'lardan oldin) (2, 5) |
| **vendor-specific** | Spec ajratgan maxsus diapazonlar (tag 0x1400+, status 0x400+, cmd 0x10200+) (3, 11) |

## Sessiya

| Termin | Izoh |
|---|---|
| **bind / unbind** | SMPP darajasidagi login/logout (4) |
| **TX / RX / TRX** | Transmitter (yuboradi) / Receiver (oladi) / Transceiver (ikkalasi) bind turlari (4) |
| **outbind** | SMSC→ESME teskari ulanish taklifi (kam ishlatiladi) (4) |
| **session state** | OPEN → BOUND_TX/RX/TRX → CLOSED — v3.4'da AYNAN 5 holat (4) |
| **enquire_link** | L7 "tiriksanmi" ping'i; javobi majburiy (4, 12, E) |
| **half-open connection** | Bir tomoni o'lgan, ikkinchisi bilmaydigan TCP — enquire_link davolaydi (12) |
| **window** | Javob kutilayotgan request'lar soni chegarasi (so'z spec'da YO'Q) (12) |
| **backpressure** | "Ulgurmayapman" signalining yuqoriga oqishi (window to'lishi, RX_T_APPN) (12) |
| **graceful shutdown** | To'g'ri yopilish: submit stop → drain → unbind → resp → close (12) |
| **reconnect / rebind** | Uzilgandan keyin qayta ulanish; backoff+jitter SHART (12, 13) |
| **thundering herd** | Hamma bir vaqtda qayta urinishi — jitter davolaydi (12) |

## Xabar va manzillar

| Termin | Izoh |
|---|---|
| **submit_sm / deliver_sm** | ESME→SMSC yuborish / SMSC→ESME yetkazish PDU'lari (5) |
| **message_id** | SMSC bergan OPAQUE identifikator — string sifatida saqlanadi (5, 9) |
| **esm_class** | Rejim/tur/GSM-flag bitlari bayti (UDHI shu yerda) (5) |
| **registered_delivery** | DLR/ack so'rovi bayti (5, 9) |
| **validity_period** | Xabarning "yashash muddati" — tugasa EXPIRED (5) |
| **TON / NPI** | Type of Number / Numbering Plan Indicator — manzil turi juftligi (6) |
| **alphanumeric sender** | Harfli sender ("Bank"): max 11 GSM7 belgi, javob berib bo'lmaydi (6) |
| **short code** | Qisqa raqam (1234) — TON=3 konventsiyasi (6) |

## Encoding va segmentlash

| Termin | Izoh |
|---|---|
| **data_coding (DCS)** | Matn kodlash sxemasi bayti; 0 = "SMSC default" (ANIQLANMAGAN!) (7) |
| **GSM7 / GSM 03.38** | 7-bitlik SMS alifbosi, 160 belgi; extension jadvali ESC orqali (7, D) |
| **septet / packed / unpacked** | 7-bit birlik; SMPP'da odatda UNPACKED (1 belgi = 1 oktet) (7) |
| **UCS2** | 16-bitlik universal kodlash (UTF-16BE), 70 belgi; kirill/emoji yo'li (7) |
| **surrogate pair** | Emoji kabi belgilar UCS2'da 2 unit — chegarada bo'linmasligi kerak (7, 8) |
| **U+02BB** | Oʻzbek oʻ/gʻ dagi rasmiy belgi — GSM7'da YO'Q, normalize qilinadi (7) |
| **concatenation** | Uzun xabarni segmentlarga bo'lish (UDH / sar_* / payload) (8) |
| **UDH / UDHI** | User Data Header (segment metadata) va uning esm_class'dagi bayrog'i (8) |
| **segment** | Concat'ning bitta qismi = alohida submit_sm = alohida narx (8) |
| **message_payload** | Matnni TLV'da tashish usuli (sm_length=0 bilan) (8, 10) |

## DLR va xatolar

| Termin | Izoh |
|---|---|
| **DLR** (Delivery Receipt) | "Xabar yetdi/yetmadi" hisoboti — deliver_sm ichida keladi (9, G) |
| **message_state** | Xabar holati: ENROUTE(1)...REJECTED(8) (9, G) |
| **receipted_message_id** | DLR'dagi original id TLV'si — matndan USTUVOR (9) |
| **korrelyatsiya** | DLR'ni original submit bilan bog'lash (hex/dec tuzog'i bilan) (9) |
| **out-of-order** | Javoblar/DLR'lar yuborish tartibida kelmasligi — NORMAL (9, 12) |
| **generic_nack** | Faqat buzuq length/notanish id'ga beriladigan javob (11) |
| **transient / permanent** | Retry ma'noli / ma'nosiz xato toifalari (industriya tasnifi) (11) |
| **backoff (exponential)** | Har urinishda kutishni 2x oshirish (11, 12) |
| **RTHROTTLED** | "Sekinla" kodi (0x58) — darhol retry TAQIQLANGAN (11) |
| **at-least-once / at-most-once** | Resend siyosati: duplicate xavfi vs yo'qotish xavfi (5, 12) |
| **idempotent** | Takror bajarilganda natijani o'zgartirmaydigan amal (duplicate DLR'ga chidam) (9, 14) |

## Muhandislik

| Termin | Izoh |
|---|---|
| **golden test / golden hex** | Kutilgan baytlarning aniq etaloni bilan solishtirish (2, 15) |
| **round-trip** | encode→decode→tenglik invarianti (15) |
| **fuzzing** | Tasodifiy-mutatsion kirishlar bilan avtomatik test (15) |
| **seed corpus** | Fuzzer'ning boshlang'ich valid namunalari (15) |
| **race detector** | Go'ning -race bayrog'i: data race'larni ushlaydi (12, 15) |
| **net.Pipe** | Sinxron, buffersiz in-memory ulanish — qattiq test muhiti (12, 15) |
| **mock SMSC** | Testlar uchun o'zimiz yozgan server (quirk'lari bilan) (14) |
| **quirk** | Operatorning spec'dan og'ishi — kod chidashi kerak bo'lgan xulq (9, 14, 16) |
| **PII / masking** | Personal data va uni log'da yashirish (mask/redact/hash) (16) |
| **token bucket / rate limiter** | Tezlik chegaralash mexanizmi (soniyasiga N) (13, 16) |
| **mTLS** | Ikki tomonlama sertifikat tekshiruvli TLS (16) |
