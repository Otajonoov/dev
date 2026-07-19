# 10-bob mashqlari: Qolgan operatsiyalar

> Avval mashqlarni mustaqil bajaring, keyin yechimlarga qarang. Bob matni: [book/10-boshqa-operatsiyalar.md](../book/10-boshqa-operatsiyalar.md)

---

## Mashq 1. query_sm vs DLR

Ikkala mexanizm ham "xabar taqdiri"ni aytadi. Jadval tuzing — kamida 6 mezon bo'yicha taqqoslang (kim boshlaydi, qachon keladi, yuk/narx, qamrov, ishonchlilik, qo'llanish o'rni). So'ng: 50 000 xabarlik kampaniyada 30 ta DLR kelmay qoldi — qaysi mexanizm bilan, qanday tartibda aniqlaysiz?

## Mashq 2. cancel_sm shartlari

Har holat uchun ayting: nima bekor bo'ladi, yoki nima uchun xato/ma'nosiz?

1. `message_id="7F3A9B"`, source=Bank(5/0), dest bo'sh
2. `message_id=""`, source=Bank(5/0), dest=998901234567(1/1), service_type=""
3. `message_id=""`, source=Bank(5/0), dest bo'sh
4. `message_id="7F3A9B"`, source=Reklama(5/0) — lekin original submit Bank'dan ketgan
5. Xabar 5 daqiqa oldin DELIVRD bo'lgan, `message_id` to'g'ri

## Mashq 3. submit_multi_resp bilan ishlash

SMSC'dan submit_multi_resp keldi: Status=0, message_id="7F3AA0", no_unsuccess=2 (998907654321 → 0x0B, 998901111111 → 0x14).

1. Bu muvaffaqiyatmi? 254 manzilning nechtasiga xabar ketadi?
2. Ikkala error_status_code'ni oching (Table 5-2). Qaysi biriga retry ma'noli, qaysi biriga yo'q?
3. Go kodi yozing: `SubmitMultiResp`ni qabul qilib, (a) muvaffaqiyatli manzillar sonini hisoblaydigan, (b) retry-ga loyiq manzillarni alohida ro'yxatga ajratadigan funksiya.
4. registered_delivery=0x01 bo'lsa bu submit'dan nechta DLR kutasiz?

---

# Yechimlar

## Yechim 1

| Mezon | query_sm | DLR |
|---|---|---|
| Kim boshlaydi | ESME (pull) | SMSC (push) |
| Qachon | So'ralganda — holat o'sha ondagi | Holat O'ZGARGANDA — final/intermediate |
| Yuk | Har tekshiruv = 1 PDU juftligi; N xabar × M davr = N×M so'rov | Har xabar uchun maksimum 1-2 deliver_sm |
| Qamrov | Bitta message_id + source aniq mos bo'lishi shart | registered_delivery so'ragan hamma xabar |
| Ishonchlilik | Sinxron javob kafolatli (SMSC qo'llasa) | DLR yo'qolishi/kechikishi mumkin (session uzilishi, route) |
| Support | Ko'p aggregator'da YO'Q (Sinch — rasman) | Universal — SMPP'ning asosiy mexanizmi |
| O'rni | Istisno vositasi: yakka tekshiruv, incident | Asosiy oqim: holat boshqaruvi shunga quriladi |

30 yo'qolgan DLR: (1) avval lokal tekshiruv — "kutish xonasi"/korrelyatsiya jadvalida osilib qolganmi (9-bob out-of-order case'i); (2) qolganlarini query_sm bilan YAKKA-YAKKA so'rash (30 ta so'rov — pinset rejimi, polling emas); (3) query ham qo'llanmasa — operator support'iga message_id ro'yxati bilan murojaat. 50 000 xabarni qayta-query qilish — anti-pattern.

## Yechim 2

1. **Rejim 1**: faqat 7F3A9B bekor qilinadi (agar hali final bo'lmagan bo'lsa) — dest shart emas.
2. **Rejim 2 (guruh)**: Bank→998901234567 yo'nalishidagi BARCHA kutayotgan xabarlar bekor qilinadi; service_type bo'sh — filtr sifatida qatnashmaydi.
3. **Xato**: message_id ham yo'q, dest ham yo'q — hech qaysi rejimga tushmaydi. Bizning `Encode` buni lokal ushlaydi; simga chiqsa SMSC ESME_RCANCELFAIL qaytargan bo'lardi.
4. **Topilmaydi**: matching message_id + SOURCE bo'yicha — source mos emas → ESME_RCANCELFAIL (0x11). Bu himoya mexanizmi ham: boshqa sender'ning xabarini bekor qilib bo'lmaydi.
5. **Kech**: xabar final holatda — bekor qilinadigan narsa yo'q → ESME_RCANCELFAIL. SMS'ni qaytarib olish texnologiyasi mavjud emas.

## Yechim 3

**1.** Qisman muvaffaqiyat: Status=0 — PDU qabul qilindi, lekin 254−2 = **252 manzilga** ketadi; 2 manzil hech qachon olmaydi.

**2.** 0x0B = **ESME_RINVDSTADR** ("Invalid Dest Addr") — manzil formati/routing'i noto'g'ri: PERMANENT, retry ma'nosiz (raqamni tuzatmaguncha natija o'zgarmaydi). 0x14 = **ESME_RMSGQFUL** ("Message Queue Full") — o'sha abonent uchun SMSC navbati to'la: TRANSIENT, keyinroq retry ma'noli (11-bob tasnifi).

**3.**

```go
// retryable — Table 5-2 bo'yicha transient kodlar (11-bobda to'liq tasnif).
func retryable(code uint32) bool {
	switch code {
	case 0x14, 0x58, 0x08: // RMSGQFUL, RTHROTTLED, RSYSERR
		return true
	}
	return false
}

// TriageMulti resp'ni tahlil qiladi: muvaffaqiyatli manzillar soni va
// retry-ga loyiq manzillar ro'yxati.
func TriageMulti(sent int, resp pdu.SubmitMultiResp) (okCount int, retry []pdu.Address) {
	okCount = sent - len(resp.Unsuccess)
	for _, u := range resp.Unsuccess {
		if retryable(u.ErrorStatusCode) {
			retry = append(retry, u.Addr)
		}
	}
	return okCount, retry
}
```

998907654321 (0x0B) ro'yxatga kirmaydi, 998901111111 (0x14) esa retry'ga tushadi — keyinroq YAKKA submit_sm bilan.

**4.** **252 ta** — muvaffaqiyatli qabul qilingan har manzil uchun bittadan (message_id bir xil, manzil har xil — korrelyatsiya kaliti id+dest juftligi). Unsuccess'dagi 2 manzil uchun DLR bo'lmaydi: xabar ular uchun umuman navbatga kirmagan.
