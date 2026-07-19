# 16-bob mashqlari: Production

> Avval mashqlarni mustaqil bajaring, keyin yechimlarga qarang. Bob matni: [book/16-production.md](../book/16-production.md)

---

## Mashq 1. Integratsiya so'rovnomasi

Yangi operator bilan SMPP integratsiyasi boshlanmoqda. Kickoff'da so'raladigan KAMIDA 10 savolli checklist tuzing — har savolga: nega muhim (qaysi bob dardi) va javob kodning qaysi konfiguratsiyasiga tushadi.

## Mashq 2. Sabablar daraxti

Alert: "DLR latency o'rtachasi 8s → 90s ga oshdi; submit success rate o'zgarmagan (99.7%)". Mumkin sabablarning daraxtini tuzing — kamida 3 shox, har shoxda tekshirish usuli va tasdiqlovchi/inkor etuvchi signal.

## Mashq 3. Masked log bilan incident

Support'ga shikoyat: "+998 90 123 45 67 raqamiga kecha 14:30 atrofida OTP kelmagan". Log'laringiz to'liq masked (MaskedSubmit formatida). Incident'ni qanday tekshirasiz — qadam-baqadam, va oxirida: masking tekshiruvga XALAQIT berdimi?

---

# Yechimlar

## Yechim 1

1. **data_coding=0 qanday talqin qilinadi?** (7-bob: GSM7/Latin-1/UTF-8 farqi — `@` buzilishi) → test xabar + provisioning hujjat.
2. **short_message packed yoki unpacked GSM7?** (7-bob quirk) → deyarli hamisha unpacked, lekin tasdiq shart.
3. **DLR id: formati — resp bilan bir xilmi (hex/dec)?** (9-bob) → NormalizeID baribir himoya qiladi, lekin bilish monitoring uchun kerak.
4. **DLR'da TLV'lar bormi yoki faqat matn?** (9-bob) → parser tayyor; TLV bo'lsa ustuvor.
5. **DLR `err:` kodlari jadvali?** (9/11-bob) → operator-specific mapping konfiguratsiyaga.
6. **TPS limiti qancha va throttling xulqi qanday (RTHROTTLED'mi, jim tashlashmi)?** (11-bob) → `RateLimiter = PerSecond(N)`, window ≤ 2–3×N.
7. **TRX qo'llanadimi yoki TX+RX juftlik kerakmi? Bind limiti nechta?** (4/13-bob) → Mode konfiguratsiyasi, RALYBND siyosati.
8. **enquire_link intervali va inactivity timeout qiymatlari?** (4/12-bob) → `Session.EnquireLink` < ularning inactivity'si.
9. **TLS talab qilinadimi, qaysi port, mTLS/sertifikat jarayoni? IP whitelist'ga qaysi IP'lar beriladi?** (16-bob) → `Config.TLS`, deploy IP rejasi.
10. **message_id qancha saqlanadi (query_sm oynasi)? Scheduled delivery/cancel/replace qo'llanadimi?** (10-bob) → istisno-yo'llar dizayni.
11. (bonus) **MO route kerakmi va qanday sozlanadi?** (5-bob) → alohida kelishuv.
12. (bonus) **Vendor xato kodlari (0x400–0x4FF) jadvali?** (11-bob) → konfiguratsiya mapping.

## Yechim 2

**Shox A — Route/operator degradatsiyasi (eng ehtimoliy).** SMSC qabul qilyapti (submit OK), lekin telefonlarga yetkazish sekinlashgan: grey route almashgan, ichki navbat o'sgan, ma'lum yo'nalish (bitta mobil operator) yotib qolgan. Tekshirish: DLR latency'ni DEST prefiks bo'yicha kesish (hamma yo'nalishdami yoki bittasidami?); stat taqsimoti (EXPIRED/UNDELIV o'sganmi?); operator status sahifasi/support. Tasdiq: bitta prefiksda lokalizatsiya — route muammosi.

**Shox B — Bizning qabul zanjirimiz sekinlashgan.** DLR'lar KELYAPTI-yu, biz sekin qayta ishlayapmiz: inbound queue to'lib RX_T_APPN oqyapti (12-bob) → SMSC retry bilan keyinroq yetkazyapti — "latency" aslida bizning kechikishimiz. Tekshirish: RX_T_APPN counter'i, inbound queue metrikasi, OnDeliver handler'ining ish vaqti (DB sekinlashganmi?). Tasdiq: queue-full loglari alert vaqti bilan mos.

**Shox C — O'lchov/korrelyatsiya artefakti.** Xabarlar aslida tez yetyapti, lekin korrelyatsiya kechikyapti: out-of-order race'da "kutish xonasi"ga tushib qolgan DLR'lar (9-bob) retry bilan kech Resolve bo'lyapti; yoki Register yozuvi (DB) sekinlashgan. Tekshirish: kutish-xonasi hajmi/age metrikasi; Resolve-topilmadi loglari. Tasdiq: deliver_sm'ning KELISH vaqti (session log) bilan Resolve vaqti orasida farq katta.

(Bonus shox D — vaqt manbai: server soati siljigan/timezone chalkashligi — done date operator lokalida (9-bob), agar latency'ni DLR matni sanalaridan hisoblayotgan bo'lsangiz. Tekshirish: o'z timestamp'laringiz bilan solishtirish.)

## Yechim 3

1. Raqamni kanonlashtiring: `998901234567` → `HashPII` = masalan `a1b2c3d4e5f6`.
2. Log'dan qidiruv: `grep 'a1b2c3d4e5f6' gateway.log` — kecha 14:25–14:40 oralig'ini filtrlaysiz.
3. Topilgan zanjirni o'qiysiz: `submit_sm src=Bank dst=9989******67(#a1b2c3d4e5f6) ...` → resp qatori (message_id KO'RINADI — u PII emas!) → DLR qatori: `stat:UNDELIV err:034` deylik.
4. Xulosa chiqarish: submit ketgan, SMSC olgan, tarmoq UNDELIV degan — muammo bizda emas; `err:034`ni operator jadvalidan ochasiz (masalan "absent subscriber" — telefon o'chiq bo'lgan). Yoki DLR qatori umuman yo'q — 2-mashq daraxtiga o'tasiz.
5. Javob: masking XALAQIT BERMADI — hash to'liq zanjirni bog'ladi, message_id ochiq handle bo'lib xizmat qildi, raqam va matn esa log'da hech qayerda ochilmadi. Aynan shu — "incident tekshiruvi mumkin, PII oshkor emas" balansining isboti (agar hash BO'LMAGANIDA — masked raqam bo'yicha qidirib bo'lmasdi: oxirgi 2 raqam + prefiks yuzlab abonentga mos).
