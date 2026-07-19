# 9-bob mashqlari: Delivery receipt (DLR)

> Avval mashqlarni mustaqil bajaring, keyin yechimlarga qarang. Bob matni: [book/09-dlr.md](../book/09-dlr.md)

---

## Mashq 1. Uch operator, uch DLR

Uchta turli SMSC'dan deliver_sm keldi (esm_class hammasida 0x04, TLV yo'q). short_message matnlari:

**(a)** `id:9F44C1 sub:001 dlvrd:001 submit date:2607191030 done date:2607191031 stat:DELIVRD err:000`

**(b)** `id:5077244932 submit date:260719103544 done date:260719104012 stat:UNDELIV err:034 text:Assalomu alayku`

**(c)** `stat:EXPIRED err:255 id:0000AB12 done date:2607190900`

Har biri uchun ayting:

1. Qaysi field'lar bor, qaysilari yo'q — va `dlr.Parse` natijasidagi Receipt qanday to'ladi (Sub/Dlvrd, sanalar, State qiymatlari bilan)?
2. (b)'dagi sana formati (a)'dan nimasi bilan farq qiladi va parser buni qanday biladi?
3. (c)'da `id:` leading zero bilan — korrelyatsiyada bu qanday muammo tug'dirishi mumkin?
4. Uchchalasini `dlr.Parse` bilan dasturda parse qilib, javobingizni testda tekshiring.

## Mashq 2. Hex/dec bog'lash

Bazangizda submit_sm_resp'dan yozilgan message_id: `04000000086ECD50`. DLR keldi: `id:288230376293190992 stat:DELIVRD ...`

1. Bu ikki id bir xabarniki ekanini isbotlang (hisob-kitob bilan).
2. `dlr.NormalizeID("288230376293190992")` qaysi variantlarni qaytaradi? Nega oddiy hex konversiya yetarli emas (16 belgili bazadagi qiymatga e'tibor bering)?
3. `dlr.Table` bilan bog'lanishni kod yozib ko'rsating: Register → Resolve → Forget zanjiri.
4. Teskari stsenariyni ham tekshiring: bazada decimal, DLR'da hex kelsa ishlaydimi?

## Mashq 3. stat:ACCEPTD

Kampaniya yubordingiz, 100 000 xabar. DLR'lar oqib kelmoqda; 3% xabarda `stat:ACCEPTD` keldi va boshqa DLR kelmayapti.

1. ACCEPTD spec bo'yicha nimani anglatadi (§ raqami bilan) va u qaysi message_state raqamiga mos?
2. Bu xabarlarni "yetkazildi" deb hisoblash mumkinmi? "Yetkazilmadi" deb-chi? Retry yuborish kerakmi?
3. `dlr.MessageState.Final()` ACCEPTD uchun nima qaytaradi va nega ayni shu qaror qilingan?
4. Hisobot/monitoring nuqtai nazaridan bu 3% bilan nima qilish to'g'ri?

---

# Yechimlar

## Yechim 1

**1.**

- **(a)** To'liq Appendix B to'plami, faqat `text:` yo'q. Receipt: ID="9F44C1", Sub=1, Dlvrd=1, SubmitDate=2026-07-19 10:30 UTC, DoneDate=10:31, Stat="DELIVRD" → State=StateDelivered, Err="000", Text="".
- **(b)** `sub:`/`dlvrd:` yo'q → Sub=-1, Dlvrd=-1 ("topilmadi" — 0 emas!); sanalar sekundli (12 raqam); Stat="UNDELIV" → State=StateUndeliverable; Err="034" (operator fazosi — jadvalsiz ma'nosi noma'lum); Text="Assalomu alayku" (original matnning kesilgan boshi — receipt'dagi text hech qachon to'liq matn kafolati emas).
- **(c)** Tartib teskari, `sub:`/`dlvrd:`/`submit date:`/`text:` yo'q: ID="0000AB12", State=StateExpired, Err="255", DoneDate bor, SubmitDate zero. Parser tartibga bog'lanmagani uchun bemalol o'qiydi.

**2.** (a) 10 raqam — YYMMDDhhmm; (b) 12 raqam — YYMMDDhhmmss. `parseReceiptDate` format tanlovini qiymat UZUNLIGIdan qiladi (10 → daqiqagacha, 12 → soniyagacha); boshqa uzunlik — zero time, xato emas.

**3.** SMSC keyingi safar (masalan query_sm javobida yoki boshqa DLR'da) o'sha id'ni `AB12` deb, leading zero'siz yozishi mumkin — string tengligi buziladi. `NormalizeID` shu holat uchun leading-zero'si kesilgan variantni ham indekslaydi ("0000AB12" → "AB12" ham lookup kalitiga aylanadi); hex talqin varianti esa ("0000AB12" hex → decimal 43794) hex/dec quirk'ni yopadi.

**4.** Test skeleti:

```go
func TestUchOperator(t *testing.T) {
	r, err := dlr.Parse([]byte("id:9F44C1 sub:001 dlvrd:001 submit date:2607191030 done date:2607191031 stat:DELIVRD err:000"), nil)
	// err == nil; r.ID=="9F44C1", r.Sub==1, r.Dlvrd==1, r.State==dlr.StateDelivered
	r, err = dlr.Parse([]byte("id:5077244932 submit date:260719103544 done date:260719104012 stat:UNDELIV err:034 text:Assalomu alayku"), nil)
	// r.Sub==-1, r.Dlvrd==-1, r.State==dlr.StateUndeliverable, sekundli sana o'qilgan
	r, err = dlr.Parse([]byte("stat:EXPIRED err:255 id:0000AB12 done date:2607190900"), nil)
	// r.ID=="0000AB12", r.State==dlr.StateExpired, r.SubmitDate.IsZero()==true
	_ = r
	_ = err
}
```

## Yechim 2

**1.** `0x04000000086ECD50` ni decimal'ga o'girsak: 0x0400000000000000 = 2⁵⁸ = 288230376151711744; 0x086ECD50 = 141479248; yig'indi = **288230376293190992** — DLR'dagi son aynan shu. (Tezroq isbot: `python3 -c "print(int('04000000086ECD50',16))"`.)

**2.** Kamida: `288230376293190992` (o'zi), `4000000086ECD50` (hex yozuvi — 15 belgi) va `04000000086ECD50` ("0" bilan juft uzunlikka to'ldirilgani). Oddiy konversiya "4000000086ECD50" beradi — bazadagi qiymat esa 16 belgili, SMSC hex id'ni BAYT chegarasiga to'ldirib yozgan. Padded variant bo'lmasa lookup baribir o'tmaydi — shu bitta "0" ko'p integratsiyalarda haftalab qidirilgan bug.

**3.**

```go
tab := dlr.NewTable()
tab.Register("04000000086ECD50")            // submit_sm_resp'dan

canon, ok := tab.Resolve("288230376293190992") // DLR'dan
// ok == true, canon == "04000000086ECD50" — bazadagi kalit topildi

// final DLR qayta ishlangach:
tab.Forget(canon) // aks holda jadval cheksiz o'sadi
```

**4.** Ishlaydi — `Resolve` kelgan id'ning O'Z variantlarini ham chiqaradi: bazada `288230376293190992` bo'lsa, DLR'dagi `04000000086ECD50`ning decimal talqin varianti o'sha kalitga uriladi (`TestTableHexDecQuirk`ning ikkinchi yarmi aynan shu stsenariy).

## Yechim 3

**1.** §5.2.28: ACCEPTED (6) — "The message is in accepted state (i.e. has been manually read on behalf of the subscriber by customer service)" — operator xizmati abonent nomidan qabul qilgan; Appendix B qisqartmasi ACCEPTD. Amalda ko'proq "hamkor tarmoq/qabul markazi oldi, telefonga yetganini tasdiqlay olmayman" ma'nosida keladi.

**2.** Ikkalasi ham emas: telefon yetkazuvi TASDIQLANMAGAN (DELIVRD emas), lekin xabar rad ham etilmagan (UNDELIV/REJECTD emas). **Retry yubormang**: xabar SMSC/tarmoqda "qabul qilingan" — qayta yuborish duplicate xavfi. Bu "noaniq oraliq" holat, yakuniy hukm uchun ma'lumot yetarli emas.

**3.** `false` — final to'plamga faqat DELIVERED, EXPIRED, DELETED, UNDELIVERABLE, REJECTED kiradi. ACCEPTD'dan keyin nazariy jihatdan yana DLR kelishi mumkin (kelmasligi ham); uni final desak korrelyatsiya jadvalidan (Forget) erta o'chirib yuborardik va kechikkan haqiqiy final DLR "notanish" bo'lib qolardi.

**4.** Alohida metrikada ko'rsatish (masalan `dlr_state{state="accepted"}`), DELIVRD foiziga QO'SHMASLIK, va route sifati signali sifatida kuzatish: ACCEPTD ulushi keskin o'ssa — route degradatsiyasi/grey route belgisi bo'lishi mumkin. TTL bilan "kutish": N soatdan keyin ham final kelmasa hisobotda "unconfirmed" toifasiga o'tkazish — biznes qarori sifatida hujjatlashtiriladi (16-bob monitoring bo'limi).
