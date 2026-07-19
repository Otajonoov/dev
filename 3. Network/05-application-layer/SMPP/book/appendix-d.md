# Appendix D. GSM 03.38 default alphabet — to'liq jadval

> Mustaqil lookup-hujjat: GSM 7-bit default alphabet (TS 23.038 §6.2.1) barcha 128 o'rni, extension table va "GSM7'da YO'Q mashhur belgilar" ro'yxati. Kontekst va packing algoritmi: [7-bob](07-encoding.md). Kodda: `code/coding/gsm7.go` (`gsm7Basic`, `gsm7Ext`).

## D.1 Basic jadval (128 o'rin)

O'qish tartibi: belgining GSM kodi = ustun sarlavhasi + qator raqami. Masalan `A` = 0x40 ustuni, 0x01 qatori → **0x41**; `Δ` = 0x10 + 0x00 → **0x10**.

| Qator | +0x00 | +0x10 | +0x20 | +0x30 | +0x40 | +0x50 | +0x60 | +0x70 |
|---|---|---|---|---|---|---|---|---|
| 0x00 | @ | Δ | (space) | 0 | ¡ | P | ¿ | p |
| 0x01 | £ | _ | ! | 1 | A | Q | a | q |
| 0x02 | $ | Φ | " | 2 | B | R | b | r |
| 0x03 | ¥ | Γ | # | 3 | C | S | c | s |
| 0x04 | è | Λ | ¤ | 4 | D | T | d | t |
| 0x05 | é | Ω | % | 5 | E | U | e | u |
| 0x06 | ù | Π | & | 6 | F | V | f | v |
| 0x07 | ì | Ψ | ' | 7 | G | W | g | w |
| 0x08 | ò | Σ | ( | 8 | H | X | h | x |
| 0x09 | Ç | Θ | ) | 9 | I | Y | i | y |
| 0x0A | LF | Ξ | * | : | J | Z | j | z |
| 0x0B | Ø | **ESC** | + | ; | K | Ä | k | ä |
| 0x0C | ø | Æ | , | < | L | Ö | l | ö |
| 0x0D | CR | æ | - | = | M | Ñ | m | ñ |
| 0x0E | Å | ß | . | > | N | Ü | n | ü |
| 0x0F | å | É | / | ? | O | § | o | à |

Diqqatga loyiq o'rinlar (ASCII bilan adashtirmaslik uchun):

| Belgi | GSM kodi | ASCII kodi | Izoh |
|---|---|---|---|
| `@` | **0x00** | 0x40 | GSM'ning "noli" — dc=0 talqin xatolarining klassik indikatori |
| `_` | **0x11** | 0x5F | |
| `$` | 0x02 | 0x24 | GSM 0x24 esa `¤` (currency sign) |
| `'` (apostrof) | 0x27 | 0x27 | Mos! — o'zbek normalizatsiyasi shunga tayanadi |
| A–Z, a–z, 0–9 | ASCII bilan mos | | Harf-raqamlar bir xil — faqat maxsus belgilar farq qiladi |

## D.2 Extension table (ESC = 0x1B orqali)

Har belgi simda **2 septet** (ESC + kod):

| Belgi | Kod | | Belgi | Kod |
|---|---|---|---|---|
| € | 0x65 | | \[ | 0x3C |
| { | 0x28 | | \] | 0x3E |
| } | 0x29 | | ~ | 0x3D |
| \\ | 0x2F | | \| | 0x40 |
| ^ | 0x14 | | FF (form feed) | 0x0A |

Qabul qiluvchi extension'ni tushunmasa: ESC → bo'sh joy sifatida ko'rsatiladi (TS 23.038). Notanish extension kodi kelsa ham xuddi shunday tolerant o'qiladi.

## D.3 GSM7'da YO'Q mashhur belgilar

Quyidagilarning HAR BIRI xabarni UCS2'ga tushiradi (limit 160 → 70), agar normalizatsiya qilinmasa:

| Belgi(lar) | Unicode | Izoh / davo |
|---|---|---|
| ʻ (oʻ, gʻ tarkibida) | U+02BB | O'zbek rasmiy imlosi; normalize → ASCII `'` |
| ʼ (tutuq belgisi) | U+02BC | normalize → `'` |
| ' ' (aqlli bir qo'shtirnoqlar) | U+2018/U+2019 | Klaviatura/CMS'lardan keladi; normalize → `'` |
| " " (aqlli qo'sh qo'shtirnoqlar) | U+201C/U+201D | GSM7'dagi `"` (0x22) bilan almashtirsa bo'ladi |
| — – (tire'lar) | U+2014/U+2013 | GSM7'dagi `-` bilan almashtirsa bo'ladi |
| … (ellipsis) | U+2026 | `...` bilan almashtirsa bo'ladi |
| Barcha kirill harflar | U+0400–U+04FF | Davo yo'q — UCS2 (yoki lotinga transliteratsiya) |
| Emoji | U+1F300+ va h.k. | UCS2'da ham 2 unit = 4 oktet (surrogate pair) |
| č š ž ō ā kabi diakritikalar | har xil | GSM7'da faqat jadvaldagi diakritikalar bor |

Eslatma: `è é ù ì ò Ç Ø å Æ ß É Ä Ö Ñ Ü à ä ö ñ ü` va grek harflari (Δ Φ Γ Λ Ω Π Ψ Σ Θ Ξ) — GSM7'da BOR, bular limitni buzmaydi.

## D.4 Tezkor tekshiruv

Kod bilan: `coding.IsGSM7(r)` — bitta belgi; `coding.SeptetLen(s)` — matnning septet uzunligi (extension = 2); `coding.Choose(text)` — normalizatsiya + dc tanlash. To'liq kontekst: [7-bob](07-encoding.md).
