# 16-bob. Production: security, monitoring, quirk katalogi va yakuniy demo

Kutubxona tayyor — endi uni operator bilan real shartnoma, real trafik va real tungi alert'lar dunyosida YASHATISH haqida gaplashamiz. Bu bob to'rt himoya qatlamini qo'shadi (TLS, rate, monitoring, PII masking), kitob davomida to'plangan operator quirk'larini yagona katalogga yig'adi, v5.0 savoliga nuqta qo'yadi va — eng muhimi — butun kitob kodini BITTA stsenariyda birlashtiruvchi yakuniy e2e demo bilan yopiladi.

## 16.1 Security: ochiq protokolni himoyalash

Qattiq haqiqatdan boshlaymiz: **SMPP v3.4'da hech qanday kriptografiya YO'Q.** Autentifikatsiya — bind ichidagi system_id + parol (max 9 bayt! — 1999-yil merosi), va u simda OCHIQ ketadi; xabar matnlari ham. Protokol himoyani transport qatlamiga qoldirgan, sanoat esa uch qatlamli konventsiyaga kelgan:

1. **SMPP over TLS.** Portlar haqida aniqlik (chalkash mavzu): IANA'da faqat **2775 (`smpp`)** ro'yxatdan o'tgan; **3550 ("ssmpp")** — keng tarqalgan **de-facto konventsiya** (node-smpp kabi kutubxonalarning default'i, Ozeki hujjatlari), lekin IANA birlamchi manbasidan tasdiqlanmagan va smpp.org uni tilga olmaydi; amalda esa operator qaysi portni bersa — o'sha. Bizning `client/tls.go` buni bir field qiladi: `Config.TLS = client.DefaultTLSConfig("smsc.operator.uz")` — MinVersion TLS 1.2 (Melrose simulyatori ham TLS 1.1+dan pastini rad etadi; 1.2+ bugungi minimal odob).
2. **Mutual TLS** — operator client sertifikat talab qilsa: `tls.Config.Certificates`ga o'z juftligingiz. Ikkala tomon ham tasdiqlangan — MITM'ga qarshi to'liq javob.
3. **IP whitelisting** — TLS bor-yo'qligidan qat'i nazar deyarli har operatorda majburiy qatlam: bind faqat kelishilgan IP'lardan. Bu SIZNING deploy'ingizga talab qo'yadi: gateway chiqish IP'si barqaror bo'lishi kerak (NAT ortidagi suzuvchi IP'lar bilan integratsiya boshlanmasdan tugaydi).

> **⚠ Amaliyotda — `InsecureSkipVerify: true` haqida.** Bu flag test muhitida (self-signed sertifikatli staging SMSC) paydo bo'ladi va production config'iga KO'CHIB QOLADI — TLS'ning MITM himoyasini butunlay o'chirib, "shifrlangan, lekin KIMGA shifrlanganini bilmaymiz" holatiga tushiradi. Bizning `DefaultTLSConfig`da bu field yo'q va bo'lmaydi; staging uchun to'g'ri yo'l — o'z CA'ingizni `RootCAs` pool'iga qo'shish. Parol gigienasi ham shu ro'yxatda: 9-baytlik parol baribir zaif — uni "himoya" emas, "identifikatsiya" deb biling; haqiqiy himoya TLS+mTLS+whitelist uchligida.

## 16.2 Monitoring: nimani o'lchash va qachon uyg'onish

12–13-boblarda qurilgan `Metrics` interfeysi to'rt hodisani beradi (SubmitObserved, DeliverReceived, SessionState, ReconnectAttempt) + pull-gauge sifatida `Session.WindowDepth()`. Ularning ustiga quriladigan minimal production dashboard (nomlar Prometheus konventsiyasida — `examples/prometheus` adapteri shu formatda chiqaradi, stdlib-only!):

| Metrika | Turi | Nimani aytadi | Alert chegarasi (boshlang'ich) |
|---|---|---|---|
| `smpp_submit_total{command_status}` | counter | Success rate + xato taqsimoti | non-OK ulushi > 5% (5 min) |
| submit RTT (sum/count → o'rtacha) | histogram | SMSC javob tezligi | p99 > ResponseTimeout/2 |
| **DLR latency** (submit → DLR) | histogram | ROUTE SIFATI — bosh indikator | o'rtacha > 30s (domestic ~10s normal; OTP <5s ideal) |
| DLR rate (delivered/submitted, 5 min oyna) | ratio | Yetkazish darajasi | < 95% |
| `smpp_bind_status` | gauge | Sessiya tirikligi | 0 holati > 1 min |
| `smpp_reconnect_total{success}` | counter | **Bind flapping** | > 3 urinish/min |
| `smpp_window_depth` | gauge | Throughput yetarliligi | doim limitda = window torlik |
| inbound queue to'lishi (RX_T_APPN soni) | counter | Handler sekinligi (12-bob) | > 0 barqaror |

Ikki nuqta alohida urg'u talab qiladi. **DLR latency — route sifatining birinchi ko'zgusi**: submit success o'zgarmasdan DLR latency o'sishi "SMSC qabul qilyapti-yu, yetkaza olmayapti" degani (grey route degradatsiyasi, operator ichki navbati, hatto hex/dec korrelyatsiya buzilishi — mashqda bu sabablar daraxti bor). **Bind flapping** — reconnect counter'ining tez o'sishi: tarmoq muammosi, operator tomonidagi limit, yoki 11-bobdagi credential-loop — har biri alohida yo'l, lekin signal bitta counter'da ko'rinadi.

DLR latency'ni KIM o'lchaydi degan savol arxitekturaviy: client uni bilmaydi (submit va DLR ikki alohida hodisa) — o'lchov korrelyatsiya nuqtasida, ya'ni sizning OnDeliver zanjiringizda: Register paytida timestamp, Resolve paytida farq. Bu atayin client'ga qo'shilmagan — korrelyatsiya biznes qatlamda (9-bob), o'lchov ham o'sha yerda.

## 16.3 PII masking: log — GDPR hujjati

MSISDN va xabar matni — personal data. Qoida qattiq: **log'ga tushgan PII'ni "keyin tozalab" bo'lmaydi** — masking log yozilishidan OLDIN, application formatter darajasida. `client/logmask.go` uch strategiyani beradi:

```go
MaskMSISDN("998901234567")  // "9989******67" — prefiks routing debug'iga yetadi
RedactText(sm)              // "[52 belgi yashirildi]" — uzunlik segment-hisobga signal
HashPII("998901234567")     // "a1b2c3d4e5f6" — deterministik 48-bit korrelyatsiya
```

Hash'ning roli incident'larda ochiladi: abonent "SMS kelmadi" deb shikoyat qildi — support raqamni oladi, `HashPII` qiymatini hisoblaydi va masked log'lardan o'sha hash bo'yicha BUTUN zanjirni (submit → resp → DLR) topadi: raqam log'da hech qayerda ochiq ko'rinmagan, tekshiruv esa to'liq o'tdi. `MaskedSubmit` — tayyor bir qatorlik formatter: `submit_sm src=Bank dst=9989******67(#a1b2c3d4e5f6) dc=0x00 reg=0x01 sm=[52 belgi yashirildi]`. To'liq PDU hex dump'lar — faqat debug rejimi va U HAM masked (matn baytlarini yashirib) — "bir kunlik debug flag" yillab yoqiq qolishini hammamiz ko'rganmiz.

## 16.4 Operator quirk katalogi — jamlanma

Kitob bo'ylab uchragan "spec'da yo'q, hayotda bor"larning yagona jadvali — integratsiya oldidan o'qib chiqiladigan ro'yxat (har biri qaysi bobda batafsil):

| Quirk | Namoyon bo'lishi | Davo | Bob |
|---|---|---|---|
| message_id hex/dec | DLR "yo'qoladi" | NormalizeID + Table | 9 |
| DLR format og'ishlari | Parser sinadi | Tolerant parser + TLV ustuvorligi | 9 |
| data_coding=0 talqini | `@` buzilishi, krakozyabra | Operator bilan kelishish; testda tekshirish | 7 |
| Packed GSM7 kutish | Matn siljigan axlat | Unpacked default; so'rab tasdiqlash | 7 |
| Session/bind limitlari | RALYBND, uzilishlar | Kelishuv + reconnect'da ehtiyot | 4, 14 |
| enquire_link majburiy intervali | Jimlikda uzish | Konfiguratsiya (30–60s) < inactivity | 4, 12 |
| Concat qismlari bitta system_id'dan | Telefon yig'olmaydi | Bir client orqali yuborish | 8 |
| Sender qayta yozilishi | Alphanumeric o'zgargan | Ro'yxatdan o'tkazish (operator jarayoni) | 6 |
| TRX yo'q, faqat TX+RX | Bind "Invalid Command ID" | Ikki client (13-bob) | 4, 13 |
| Throttling xulqlari | RTHROTTLED yoki JIM tashlash | Rate limiter + Classify + backoff | 11, 12 |
| Vendor xato kodlari | 0x400–0x4FF sirli kodlar | Jadvalni so'rab konfiguratsiyada saqlash | 11 |

Bu jadvalning amaliy shakli — **integratsiya so'rovnomasi**: operator bilan kickoff'da so'raladigan savollar ro'yxati (mashq 1'da to'liq tuzasiz): data_coding=0 nima deb talqin qilinadi? DLR id formati hex'mi/dec'mi? TPS limiti va throttling xulqi? TRX bormi? enquire_link intervali? DLR `err:` jadvali? TLS/port? IP whitelist jarayoni? message_id'ning saqlanish muddati? MO route kerakmi?

## 16.5 SMPP v5.0: bilish kerak, shoshilish shart emas

v5.0 (2003) nimalar qo'shgani va nega baribir v3.4 o'rganganimiz:

- **congestion_state TLV** — eng qiziq g'oya: har response'ga 0–100 "bandlik" ko'rsatkichi qo'shilishi mumkin; ESME tezlikni shunga moslaydi. Bu fixed window'dan nazariy jihatdan yaxshiroq flow control (spec'ning o'zi: qo'llansa window'ni tashlab yuborsa ham bo'ladi) — TCP congestion control g'oyasining SMPP'ga ko'chishi.
- **broadcast_sm oilasi** — Cell Broadcast (bir hududdagi hamma telefonga) operatsiyalari.
- **7-state session model** (4-bobda aytilgan UNBOUND/OUTBOUND qo'shimchalari) va yangi xato kodlari (jumladan mashhur **RINVDCS 0x104** — 11-bob).
- Yana: bind_resp'siz ham ishlaydigan session addressing, qo'shimcha TLV'lar.

Adoption haqiqati esa: so'rovlarda v3.4 ~54% bilan yetakchi, v5.0 ~8% atrofida (bitta so'rov — tartib sifatida o'qing); Infobip kabi yirik aggregator'lar SMPP interfeysida FAQAT v3.4 qo'llaydi. Xulosa o'zgarmadi: **v3.4 — de-fakto sanoat tili**; v5.0'ni tushunish uchun v3.4 poydevori baribir shart edi — endi sizda u bor, v5.0 spec'ini o'qish bir oqshomlik ish.

## 16.6 Yakuniy demo: examples/e2e

Va nihoyat — 1-bobda va'da qilingan Definition of Done: butun kitob kodi bitta stsenariyda, `go test ./examples/e2e -v` bilan:

```mermaid
sequenceDiagram
    participant T as e2e test
    participant C as client (13-bob)
    participant S as mock SMSC (14-bob)<br/>quirk: hex-resp/dec-DLR + TLV'siz
    T->>C: Dial (TLS'siz lokal; auth: e2e/sirli)
    C->>S: bind_transceiver -> resp (4-bob)
    T->>C: SubmitLong("...oʻzbek..." U+02BB!)
    Note over C: Normalize -> GSM7, 1 segment (7-bob)
    C->>S: submit_sm -> resp id=HEX (5-bob)
    T->>C: SubmitLong(kirill matn)
    Note over C: UCS2, 1 segment (7-bob)
    T->>C: SubmitLong(250 belgili matn)
    Note over C: UDH 8-bit, 2 segment, 2 seq (8-bob)
    S-->>C: 4 ta deliver_sm DLR: id DECIMAL, TLV YO'Q
    Note over C: Parse (tolerant!) -> NormalizeID -><br/>Table.Resolve hex↔dec TOPADI (9-bob)
    T->>T: 4/4 DELIVRD tasdiqlandi
```

```
$ go test ./examples/e2e -v
=== RUN   TestEndToEnd
    e2e_test.go:143: e2e OK: 4 xabar (lotin+kirill+2 segment), hammasi
                     hex/dec quirk ostida DELIVRD deb korrelyatsiya qilindi
--- PASS: TestEndToEnd
```

Bu 60 millisekundlik testning ichida: framing (2-bob), TLV (3), bind/state (4), submit/deliver (5), addressing (6), encoding+normalizatsiya (7), concatenation (8), DLR parse+korrelyatsiya (9), status'lar (11), session engine (12), client (13) va mock server quirk'lari (14) — hammasi birga, `-race` ostida, tashqi dunyoga bog'liqliksiz. Loyihaning bosh DoD'si bajarildi.

## Xulosa — va kitob yakuni

Production to'rt qatlam qo'shdi: TLS (2775 IANA, 3550 konventsiya, mTLS+whitelist uchligi, InsecureSkipVerify taqiqlangan), monitoring (DLR latency — route sifatining birinchi ko'zgusi; bind flapping; window depth; hammasi kichik Metrics interfeysi orqali, adapter — sizning tanlov), PII masking (mask/redact/hash — log yozilishidan OLDIN) va quirk katalogi (integratsiya so'rovnomasiga aylanadigan jadval). v5.0 — bilib qo'yiladigan kelajak, v3.4 — ishlanadigan bugun.

Kitob boshidagi va'daga qaytaylik: "materialni o'qigan Go dasturchi SMPP'ni tashqi manbasiz, to'liq va to'g'ri qayta yoza olishi kerak" edi. Yo'lda nimalar yig'ildi: 16 bob, o'nlab qo'lda tekshirilgan hex dump, 8 package (~150 test, 7 fuzz target, e2e demo), va — muhimi — har qarorning NEGAsi: nega unpacked, nega mutex+map, nega fail-fast, nega at-most-once. Endi sizda nafaqat ishlaydigan kod, balki uni HIMOYA QILA OLADIGAN tushuncha bor. Ilovalar (A–G) — kundalik lookup; mashqlar — bilimni qotirish; mock SMSC — keyingi tajribalar maydoni. Omad, va DLR'laringiz doim DELIVRD bo'lsin!

**Takrorlash savollari** (javoblar matnda bor — o'zingizni tekshiring):

1. SMPP'ning o'zida qanday himoya bor va sanoat uchligi nimalardan iborat?
2. 3550 portining maqomi qanday va bu nima uchun muhim nuance?
3. DLR latency oshgan-u submit success o'zgarmagan — bu qanday sinf muammo?
4. Nega DLR latency o'lchovi client'ga emas, korrelyatsiya qatlamiga qo'yiladi?
5. HashPII incident tekshiruvida qanday ishlaydi?
6. congestion_state fixed window'dan nimasi bilan nazariy yaxshiroq?
7. e2e demo qaysi ikki quirk'ni yoqib o'tadi va nega aynan ular?

**Mashqlar:** [exercises/16-production.md](../exercises/16-production.md) — integratsiya so'rovnomasi, DLR-latency sabablar daraxti va masked-log incident mashqi.

---

**Oldingi bob:** [15-bob. Testing va tooling](15-testing.md) · **Mundarija:** [1-bob](01-sms-ekotizimi.md)dan boshlanadi · Ilovalar: [A](appendix-a.md) · [B](appendix-b.md) · [C](appendix-c.md) · [D](appendix-d.md)

## Manbalar

- [Ozeki — Secure SMPP over SSL/TLS](https://ozeki-sms-gateway.com/p_7618-secure-smpp-connection-over-ssl-tls.html) va [EMnify — SMPP with TLS](https://www.emnify.com/developer-blog/enhanced-security-for-a2p-sms-using-smpp-with-tls) — TLS amaliyoti, 2775/3550 konventsiyasi, mTLS
- [OneUptime — Monitor SMS gateway delivery latency](https://oneuptime.com/blog/post/2026-02-06-monitor-sms-gateway-delivery-latency-opentelemetry/view) — DLR latency chegaralari (10s/30s/5s OTP), delivery rate 95% alert
- [Plivo — SMS data redaction](https://www.plivo.com/docs/messaging/concepts/sms-data-redaction) — PII masking sanoat namunasi
- [smpp.org — SMPP v5](https://smpp.org/smpp-v5.html) — congestion_state va v5 yangiliklari; [Infobip SMPP spec](https://www.infobip.com/docs/essentials/api-essentials/smpp-specification) — yirik aggregator faqat v3.4 qo'llashi
- [Ozeki — SMPP version comparison](https://ozeki-sms-gateway.com/p_6393-smpp-protocol-version-comparison.html) — 54%/8% adoption so'rovi (tartib sifatida)