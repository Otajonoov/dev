# Appendix E. SMPP timer'lari

> Mustaqil lookup-hujjat: v3.4 §2.9/§7.2 timer'lari. Batafsil: [4-bob](04-bind-session.md) (tushunchalar), [12-bob](12-session-engine.md) (client implementatsiyasi), [14-bob](14-mock-smsc.md) (server implementatsiyasi).

## Bosh qoida

**Spec'da BIRORTA timer'ning qiymati YO'Q** — §7.2: "outside the scope of this specification", hammasi konfiguratsiyalanadigan bo'lishi shart. Internetda ko'rgan har qanday raqam (30s, 60s...) — operator konventsiyasi yoki kutubxona default'i, spec fakti EMAS.

## To'rt timer

| Timer | Kimda faol | Nimani o'lchaydi | Muddati o'tganda | Bizning kodda |
|---|---|---|---|---|
| **session_init** | SMSC'da SHART | TCP connect → bind orasi | Ulanish UZILADI | server: bind'gacha read-deadline (`smsc/server.go`) |
| **enquire_link** | Ikkala tomonda | Oxirgi yuborilgan ping'dan beri | enquire_link YUBORILADI | client: `Session.Config.EnquireLink` ticker |
| **inactivity** | Odatda SMSC'da | Har qanday PDU'siz jimlik | Sessiya RESET qilinadi | server: har PDU'da yangilanadigan read-deadline |
| **response** | Request yuborganda | request → resp orasi | Request "javobsiz" — operatsiyaga mos chora | client: window deadline + expire scanner |

## Munosabat qoidasi

```
response_timeout  <  enquire_link_interval  <  inactivity_timeout
```

- response < enquire: ping'ning o'z javob oynasi bo'lishi kerak;
- enquire < inactivity (QARSHI tomonning inactivity'si!): jim daqiqalarda ping'lar sessiyani "to'ydirib" turadi — aks holda SMSC sizni uzadi.

## Amaliy diapazonlar (konventsiya!)

| Timer | Keng tarqalgan qiymatlar | Izoh |
|---|---|---|
| session_init | 5–30s | Server resurs himoyasi; bizning default 5s |
| enquire_link | **30–60s** (konsensus) | Manbalar 15s–15min oralig'ida tarqoq; operator TZ'si hal qiladi |
| inactivity | 60–120s | Ko'p provayderlarda shu tartibda (~80s ham uchraydi) |
| response | 10–60s | SMSC RTT p99'idan 2–3 barobar katta; transaction-mode data_sm'da SEKUNDLAR kerak ([10-bob](10-boshqa-operatsiyalar.md)) |

## Bog'liq tushunchalar

- **Half-open connection** — enquire_link'ning asl dushmani: L4 tirik, L7 o'lik ([12-bob](12-session-engine.md)). TCP keepalive (default ~2 soat, faqat stack'ni tekshiradi) uni ALMASHTIRA OLMAYDI.
- **Javobsiz enquire_link** = sessiya o'ldi: log + close + reconnect (backoff-jitter bilan).
- Eski darsdagi "enquire_link 30s, reconnect 60/120s" qiymatlari — BIR operator shartnomasidan misol, standart emas.
