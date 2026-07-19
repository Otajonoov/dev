# 13-bob mashqlari: ESME client API

> Avval mashqlarni mustaqil bajaring, keyin yechimlarga qarang. Bob matni: [book/13-client-api.md](../book/13-client-api.md)

---

## Mashq 1. RespTimeout=1s stsenariysi

fiorix kutubxonasining default RespTimeout'i 1 soniya. Operator SMSC'si esa peak soatlarda submit_sm_resp'ni 1.5–3 soniyada qaytaradi.

1. Bitta OTP xabari uchun hodisalar zanjirini yozing: client nima ko'radi, SMSC nima qiladi, abonent nima oladi?
2. Bu stsenariyda naive retry ("timeout — qayta yubor") kunlik 100 000 xabarda taxminan nechta duplicate beradi (peak'da traffic'ning 20%'i 1s+ javob oladi desak)?
3. Bizning client bu xatoni qaysi ikki qaror bilan oldini oladi?
4. RespTimeout'ni qancha qilish kerak va uning yuqori chegarasini nima belgilaydi?

## Mashq 2. ErrBusy siyosati

Bizning fail-fast client'ga "submit navbati" qo'shmoqchisiz: Client ichida sig'imi 500 lik queue; Submit navbatga soladi, alohida worker'lar yuboradi; navbat to'lsa `ErrBusy` qaytadi.

1. Bu dizayn fail-fast falsafamizning qaysi muammosini QAYTARIB olib keladi va qaysi birini hal qiladi?
2. Queue'dagi xabar "eskirishi" muammosini qanday hal qilasiz (OTP 5 daqiqa keyin yuborilmasligi kerak)?
3. Shu siyosatni implement qiling: `BufferedClient` wrapper (Client'ni o'zgartirmasdan!) — Submit(ctx, sm) ErrBusy yoki async natija.
4. Qaysi metrikalar shart bo'ladi (kamida 3)?

## Mashq 3. #178'ni qayta yaratish

gosmpp #178: `SubmitSM.Split()` multipart qismlariga bir xil sequence_number bergan.

1. Bu bug'da SMSC tomonda nima sodir bo'ladi va client qaysi javoblarni "yo'qotadi"?
2. Bizning kodda xuddi shu bug'ni SUN'IY yaratish uchun qaysi qatorni qanday buzish kerak bo'lardi? (Nazariy — buzmang!)
3. `TestSubmitLongUniqueSequences` bu regression'ni qanday ushlaydi — mexanizmini tushuntiring.
4. Testni haqiqatan buzib ko'ring (vaqtincha!): sequencer o'rniga fixed seq qo'yilsa test qanday xato bilan yiqiladi? Tekshirib, qaytaring.

---

# Yechimlar

## Yechim 1

**1.** T+0: submit_sm ketdi. T+1s: client timeout deb topadi (lekin SMSC xabarni OLGAN, qayta ishlayapti). T+1.5s: submit_sm_resp keladi — client uchun "notanish seq" (window'dan expire bo'lgan), log'ga tushib yo'qoladi. Client retry qiladi: T+1s+ε ikkinchi submit_sm. Natija: SMSC IKKITA xabar oldi → abonentga IKKITA OTP. Client hisobotida esa: 1 timeout, 1 muvaffaqiyat.

**2.** 100 000 × 20% = 20 000 xabar timeout oladi; har biri kamida bitta duplicate — **~20 000 duplicate/kun** (retry ham peak'ga tushsa kaskadlanib undan ko'p). Har duplicate — abonent ishonchsizligi va pul.

**3.** (a) ResponseTimeout default 10s — SMSC'ning real p99'idan katta; (b) timeout'ni StatusError'dan AJRATIB qaytaramiz (`session.ErrResponseTimeout`) — chaqiruvchi "taqdiri noma'lum" rejimini bilib, ko'r-ko'rona retry qilmaydi (11-bob).

**4.** Operator resp-latency p99'idan kamida 2–3 barobar katta (o'lchab olinadi — 16-bob metrikasi); yuqori chegara — window drain vaqti va xabar "eskirish"i: 60s timeout bilan window'dagi xabarlar uzilishda 60s "osilib" turadi. Amaliy oraliq 10–30s.

## Yechim 2

**1.** Qaytadi: xabarlar yana client ichida "ko'rinmas" holatda (application kuzata olmaydi). Hal bo'ladi: qisqa uzilishlar chaqiruvchiga sezilmaydi (retry-siz silliq o'tadi) va backpressure ErrBusy sifatida aniq chegarada.

**2.** Har entry'ga deadline (enqueue vaqti + TTL yoki ctx deadline'ining o'zi); worker olganda avval tekshiradi — muddati o'tgan bo'lsa yubormasdan xato callback/channel'ga qaytaradi. OTP uchun TTL=validity bilan tenglashtiriladi.

**3.** Skelet:

```go
type BufferedClient struct {
	c     *client.Client
	queue chan job
}

type job struct {
	ctx  context.Context
	sm   pdu.SubmitSM
	resp chan submitResult // {id string; err error}
}

func (b *BufferedClient) Submit(ctx context.Context, sm pdu.SubmitSM) (string, error) {
	j := job{ctx: ctx, sm: sm, resp: make(chan submitResult, 1)}
	select {
	case b.queue <- j:
	default:
		return "", ErrBusy // navbat to'la - aniq chegara
	}
	select {
	case r := <-j.resp:
		return r.id, r.err
	case <-ctx.Done():
		return "", ctx.Err() // worker baribir yuborishi mumkin - hujjatlang!
	}
}

func (b *BufferedClient) worker() {
	for j := range b.queue {
		if j.ctx.Err() != nil {
			j.resp <- submitResult{err: j.ctx.Err()} // eskirgan - yubormaymiz
			continue
		}
		id, err := b.c.Submit(j.ctx, j.sm)
		if errors.Is(err, client.ErrNotBound) {
			// reconnect ketmoqda: qisqa kutish bilan qayta navbatga -
			// yoki darhol xato; siyosat shu yerda ko'rinadi.
		}
		j.resp <- submitResult{id: id, err: err}
	}
}
```

**4.** Metrikalar: queue depth gauge (to'lish trendi), ErrBusy counter (rad darajasi), queue-age histogram (eskirish), worker throughput. Busiz buffer "ko'r quti" bo'lib qoladi.

## Yechim 3

**1.** SMSC uchta submit_sm oladi (uchchalasi seq=N!), uchta submit_sm_resp qaytaradi — hammasi seq=N bilan. Client window'ida seq=N ostida BITTA entry bor: birinchi resp uni yopadi, qolgan ikkitasi "notanish seq" bo'lib tashlanadi. Natija: 2-3-segmentlarning message_id'lari YO'QOLADI → ularning DLR'larini bog'lab bo'lmaydi; yomonroq variantlarda birinchi resp'ning message_id'si "butun xabarniki" deb saqlanib, hisobot butunlay chalkashadi.

**2.** `session.send` ichidagi `seq := s.seq.Next()` ni sikldan tashqariga chiqarish — SubmitLong barcha segmentlariga bitta seq berardi. Bizning dizaynda buni QILIB BO'LMAYDI: seq Send ichida olinadi, SubmitLong esa har segmentga alohida Submit chaqiradi — bug uchun arxitekturani buzish kerak.

**3.** TestServer message_id'ni `"TST%08X" % seq` qilib yasaydi → seq'lar unique bo'lsa id'lar unique; bug bo'lsa ikkala segment BIR XIL id oladi va `ids[0] == ids[1]` tekshiruvi yiqiladi. Regression testning go'zalligi: bug'ning ICHKI sababini (seq) TASHQI kuzatiladigan natija (id) orqali ushlaydi.

**4.** (Tajriba) `Sequencer.Next`ni vaqtincha `return 42` qilsangiz: `TestSubmitLongUniqueSequences` "segmentlar BIR XIL seq olgan (#178!)" bilan yiqiladi; qo'shimcha, `TestSendOutOfOrderResponses` ham sinadi (ikkala Send bitta seq'ni window'ga qo'yishga urinadi). Qaytarishni unutmang — `go test ./... -race` toza bo'lsin.
