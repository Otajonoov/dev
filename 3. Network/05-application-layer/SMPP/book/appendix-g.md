# Appendix G. Tezkor shpargalka: DLR, TON/NPI, data_coding

> Kundalik ishda eng ko'p ochiladigan jadvallar bir joyda. Batafsil: DLR — [9-bob](09-dlr.md), manzillar — [6-bob](06-addressing.md), encoding — [7-bob](07-encoding.md).

## DLR (Appendix B formati)

Tipik matn (NORMATIV EMAS — "vendor specific... typical example"!):

```
id:7F3A9B sub:001 dlvrd:001 submit date:2607171205 done date:2607171206 stat:DELIVRD err:000 text:Salom
```

| Field | Ma'no | Ogohlantirish |
|---|---|---|
| `id:` | Original message_id | Resp'dagi bilan hex/dec FARQ qilishi mumkin — NormalizeID |
| `sub:` / `dlvrd:` | Submit/yetkazilgan qismlar soni | Yo'q bo'lishi mumkin (bizda -1 sentinel) |
| `submit date:` / `done date:` | YYMMDDhhmm (10) yoki sekundli (12) | Timezone — operator lokali, spec'da YO'Q |
| `stat:` | Table B-2 qisqartmasi | "DELIVERED" kabi uzun variantlar keladi — prefiks bo'yicha o'qing |
| `err:` | Tarmoq/SMSC kodi | Operator-specific — jadvalini so'rang; command_status EMAS |
| `text:` | Original matn boshi (≤20) | Ko'pincha yo'q; hech qachon to'liq matn kafolati emas |

## message_state (§5.2.28 ↔ Table B-2)

| Raqam | Nom | stat: | Final? |
|---|---|---|---|
| 1 | ENROUTE | — | yo'q |
| 2 | DELIVERED | DELIVRD | **HA** |
| 3 | EXPIRED | EXPIRED | **HA** |
| 4 | DELETED | DELETED | **HA** |
| 5 | UNDELIVERABLE | UNDELIV | **HA** |
| 6 | ACCEPTED | ACCEPTD | yo'q* |
| 7 | UNKNOWN | UNKNOWN | yo'q* |
| 8 | REJECTED | REJECTD | **HA** |

\* DELIVRD deb SANAMANG; alohida metrikada kuzating ([9-bob](09-dlr.md)).

DLR TLV'lari: `receipted_message_id` 0x001E (matndan USTUVOR), `message_state` 0x0427, `network_error_code` 0x0423 (3 oktet: type 1=ANSI-136/2=IS-95/3=GSM + 2 oktet kod).

DLR'ni tanish: `(esm_class >> 2) & 0x0F == 0x01` (intermediate = 0x08). Manzillar ALMASHGAN bo'ladi.

## TON (Table 5-3)

| Qiymat | Nom | Qachon |
|---|---|---|
| 0 | Unknown | Noma'lum/default; bo'sh source |
| 1 | International | To'liq xalqaro raqam ('+'SIZ!): 998901234567 |
| 2 | National | Milliy formatdagi raqam |
| 3 | Network Specific | Short code (konventsiya) |
| 4 | Subscriber Number | Kam ishlatiladi |
| 5 | Alphanumeric | "Bank" — max 11 GSM7 belgi |
| 6 | Abbreviated | Qisqartirilgan |

## NPI (Table 5-4) — qiymatlar KETMA-KET EMAS!

| Qiymat | Nom |
|---|---|
| 0 | Unknown |
| 1 | ISDN (E.163/E.164) — oddiy telefon |
| 3 | Data (X.121) |
| 4 | Telex (F.69) |
| 6 | Land Mobile (E.212) |
| 8 | National |
| 9 | Private |
| 10 (0x0A) | ERMES |
| 14 (0x0E) | Internet (IP) |
| 18 (0x12) | WAP Client Id |

Amaliy juftliklar (KONVENTSIYA, spec emas): international **1/1**; national 2/1 yoki 2/8; alphanumeric **5/0**; short code 3/0 (yoki 6/0, 0/0 — operator aytadi); bo'sh source 0/0.

## data_coding (§5.2.19)

| Qiymat | Ma'no | Limit (yakka/segment) |
|---|---|---|
| 0x00 | SMSC default — ANIQLANMAGAN, operator bilan kelishiladi! | 160/153 (GSM7 bo'lsa) |
| 0x01 | IA5/ASCII | — |
| 0x02, 0x04 | 8-bit binary | 140/134 |
| 0x03 | Latin-1 (ISO-8859-1) | 140/134 |
| 0x05 | JIS | — |
| 0x06 | Kirill ISO-8859-5 (nazariy — ISHLATILMAYDI) | — |
| 0x07 | Latin/Hebrew | — |
| **0x08** | **UCS2 (UTF-16BE)** — kirill/emoji yo'li | **70/67** |
| 0x0D / 0x0E | Extended Kanji / KSC 5601 | — |
| 0xC0–0xDF | GSM MWI control | — |
| 0xF0–0xFF | GSM message class | — |

Segment limitlari (UDH 8-bit bilan): GSM7 **153**, 8-bit **134**, UCS2 **67**; 16-bit ref bilan: 152/133/66. SMPP'da short_message odatda **UNPACKED** (belgi=oktet), sm_length OKTETLARDA (UCS2 70 belgi = 140 oktet).

O'zbek matni esdaligi: kirill → UCS2 (70/67); lotin — U+02BB/02BC/2018/2019 belgilarini ASCII apostrofga NORMALIZE qiling, aks holda bitta belgi xabarni 160→70 ga tushiradi ([7-bob](07-encoding.md)).
