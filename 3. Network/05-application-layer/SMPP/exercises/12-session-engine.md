# 12-bob mashqlari: Session engine

> Avval mashqlarni mustaqil bajaring, keyin yechimlarga qarang. Bob matni: [book/12-session-engine.md](../book/12-session-engine.md)

---

## Mashq 1. Kaskad anatomiyasi

12-bobdagi production kaskadini Mermaid sequence diagram qilib chizing: DLR burst → jobQueue → OnPDU → read loop → enquire_link → ulanish o'limi zanjiri ko'rinsin. So'ng bizning `session` package kodidan kaskadni uzadigan ANIQ joy(lar)ni ko'rsating (fayl + funksiya) va har biri zanjirning qaysi bo'g'inini uzishini ayting.

## Mashq 2. Window matematikasi

RTT = 80ms, operator limiti 200 TPS.

1. Window=1 bilan 1000 SMS qancha vaqt oladi? Window=32 bilan-chi?
2. 200 TPS'ga yetish uchun minimal window nechchi? Empirik "2–3×TPS" qoidasi bo'yicha maksimal oqilona window-chi?
3. Window=500 qo'ysangiz nima yaxshilanadi va nima YOMONLASHADI (kamida 3 xavf)?
4. `TestWindowFullBlocks` testidagi stsenariyni window=2 bilan qayta yozing (qog'ozda): qaysi Send qachon bloklanadi?

## Mashq 3. Read loop ichida sync Submit — o'z-o'zini deadlock

Faraz qiling, kimdir `OnInbound` handler'ida (masalan kelgan MO xabarga avtomatik javob yubormoqchi bo'lib) `s.Send(ctx, submit)` chaqirdi — va session'imizda dispatcher YO'Q bo'lib, OnInbound to'g'ridan-to'g'ri reader'dan chaqirilardi (gosmpp modeli).

1. Deadlock zanjirini qadam-baqadam yozing: Send nimani kutadi, reader nima bilan band?
2. Bizning arxitekturada (dispatcher goroutine bor) xuddi shu kod deadlock qiladimi? Qaysi RESURS baribir xavf ostida?
3. Buni isbotlaydigan test yozing: `OnInbound` ichida Send chaqiring va uning window to'la+handler band holatida nima bo'lishini kuzating.
4. Xulosa qoidasini bitta jumla qilib ayting.

---

# Yechimlar

## Yechim 1

```mermaid
sequenceDiagram
    participant S as SMSC
    participant R as read loop (gosmpp)
    participant Q as jobQueue (500)
    participant W as 5 worker
    S->>R: DLR burst (yuzlab deliver_sm)
    R->>Q: OnPDU -> enqueue... enqueue...
    Note over W: worker'lar DB'da band<br/>(ba'zilari getMessageWithRetry'da UXLAB yotibdi)
    Note over Q: jobQueue TO'LDI
    R--xQ: enqueue BLOKLANDI (OnPDU qaytmayapti)
    Note over R: read loop keyingi baytni O'QIMAYDI
    S->>R: enquire_link_resp (socket buferida yotib qoladi)
    Note over R: 30s "jimlik" (aslida javob kelgan!)
    R->>R: ulanish "o'lik" deb topildi -> close
    Note over S,W: barcha Submit -> ErrConnectionClosing;<br/>reconnect -> queue hali ham to'la -> sikl boshiga
```

Uzuvchi joylar (`code/session/session.go`):

1. **`handleInboundRequest`** — enqueue `select ... default` bilan NON-BLOCKING: queue to'la bo'lsa reader bloklanmaydi (zanjirning "OnPDU qaytmayapti" bo'g'ini uziladi), deliver_sm'ga RX_T_APPN qaytadi — yo'qotish protokol darajasida halol.
2. **`readLoop`dagi `case pdu.EnquireLink`** — ping javobi queue'dan MUTLAQO o'tmaydi, read path'da yoziladi ("ping javobi navbatda qotib qoldi" bo'g'ini strukturaviy yo'q).
3. **`dispatchLoop`** — foydalanuvchi handler'i reader'dan ajratilgan: handler qancha sekin bo'lsa ham o'qish davom etadi (kaskadning boshlanish sharti yo'qoladi).

## Yechim 2

**1.** Window=1: throughput = 1/0.08 = 12.5 TPS → 1000/12.5 = **80 soniya**. Window=32: nazariy 32/0.08 = 400 TPS, lekin operator limiti 200 TPS ushlab qoladi → **~5 soniya** (1000/200).

**2.** Minimal: window ≥ TPS × RTT = 200 × 0.08 = **16**. Empirik yuqori chegara: 2–3 × 200 = **400–600** — undan katta window foyda bermaydi, faqat xavf oshiradi.

**3.** Yaxshilanadi: hech narsa (200 TPS limiti baribir shift) — window 16'dan oshgach throughput o'zgarmaydi. Yomonlashadi: (a) ulanish uzilsa 500 tagacha xabar "taqdiri noma'lum" holatga tushadi (duplicate/yo'qotish dilemmasi ulgurji miqyosda); (b) SMSC tomonda 500 outstanding — RTHROTTLED yog'ilishi va "misbehaving client" tamg'asi; (c) xotira/scan narxi: expire scanner har intervalda 500 entry'ni aylanadi; (d) batch to'xtaganda oxirgi 500 ta birdan expire bo'lib alert bo'roni.

**4.** Window=2: Send#1 → slot 1 (o'tdi), Send#2 → slot 2 (o'tdi), Send#3 → **bloklanadi** (slots to'la); server Send#1'ga javob berishi bilan slot bo'shaydi → Send#3 davom etadi; Send#4 endi Send#2'ning javobini kutadi.

## Yechim 3

**1.** gosmpp modelida: OnInbound reader goroutine'ida ishlayapti → Send resp kutadi → resp'ni o'qiydigan YAGONA joy — o'sha reader → reader esa OnInbound'dan qaytmagan → resp hech qachon o'qilmaydi → Send abadiy kutadi → reader abadiy band. Klassik o'z-o'zini deadlock: "javobni kutayotgan kod javobni o'qiydigan kodni band qilib turibdi".

**2.** To'liq deadlock YO'Q: Send dispatcher goroutine'ida kutadi, reader esa alohida — resp o'qilib window'ga yetib boradi, Send qaytadi. Lekin xavf ostidagi resurs — **dispatcher'ning o'zi va inbound queue**: Send kutgan vaqt davomida dispatcher boshqa inbound'larni qayta ishlamaydi → queue to'ladi → yangi deliver_sm'lar RX_T_APPN olib tashlana boshlaydi. Sessiya tirik, lekin DLR oqimi degradatsiyada. (Va window to'la bo'lsa Send avval slot uchun ham kutadi — degradatsiya chuqurlashadi.)

**3.** Test skeleti (package session):

```go
func TestSendInsideHandlerDegradesQueue(t *testing.T) {
	s, peer := newPair(t, Config{EnquireLink: -1, InboundQueue: 1,
		OnInbound: func(r Resp) {
			// Handler ichida sync Send — javobni server ATAYIN bermaydi.
			ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
			defer cancel()
			s.Send(ctx, testSubmit("998900000009")) // ctx bilan tugaydi
		}})
	bindTRX(t, s, peer)
	// deliver #1 handler'ni band qiladi (u ichida Send kutmoqda)...
	// deliver #2 queue'ga sig'adi, deliver #3 esa RX_T_APPN oladi —
	// sessiya TIRIK (enquire_link'ka javob keladi), lekin DLR'lar rad
	// etilayotgani ko'rinadi. Reader esa hech qachon qotmaydi.
}
```

Kuzatuv: sessiya o'lmaydi (bizning arxitektura yutuq shu yerda), lekin queue darhol to'lib RX_T_APPN oqadi — ya'ni handler'da kutish METRIKADA ko'rinadigan degradatsiya.

**4.** Qoida: **handler va worker'larda hech qachon in-place kutmang (Send, Sleep, sekin DB) — ish navbatga, kutish esa o'z goroutine'siga; read/dispatch yo'llari faqat marshrutlaydi.**
