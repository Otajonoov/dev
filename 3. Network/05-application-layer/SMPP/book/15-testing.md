# 15-bob. Testing va tooling: piramida, fuzzer va Wireshark

Kitob davomida testlar har bobda yozildi — bu bob ularni TIZIMGA aylantiradi: nima qayerda testlanadi (piramida), qo'lda o'ylab topib bo'lmaydigan xatolarni kim topadi (fuzzing), tezlik qayerda o'lchanadi (benchmark), sim ustidagi haqiqatni kim ko'rsatadi (Wireshark) va bularning bari qanday qilib har commit'da avtomatik yuguradi (ci.sh). Va bobning eng qimmatli qismi rejalashtirilmagan edi: **fuzzer bu bobni yozish jarayonida kodimizdan haqiqiy bug topdi** — o'sha voqeani to'liq, yashirmasdan ko'rsatamiz, chunki u fuzzing'ning qiymatini istalgan nazariy tushuntirishdan yaxshi isbotlaydi.

## 15.1 Test piramidasi: nima qayerda yashaydi

Loyihada to'planganini qatlamlarga yoysak:

| Qatlam | Transport | Misollar | Nimani ushlaydi |
|---|---|---|---|
| Codec (pure) | yo'q — baytlar | golden hex testlar (2–11-boblar) | Har baytning to'g'riligi, spec mosligi |
| State/logic | yo'q | Table 2-1, Classify, NormalizeID | Qoidalar va tasniflar |
| Session | net.Pipe | 12-bob testlari | Concurrency, deadlock, backpressure |
| Integratsiya | real listener | client↔smsc, quirk testlari (13–14) | Qatlamlar birga ishlashi, lifecycle |
| Fuzz | baytlar | FuzzDecode, FuzzParse... | O'ylab TOPILMAGAN kirishlar |
| Qo'lda/interop | tashqi | Melrose, SMPPSim, Wireshark | Boshqa implementatsiyalar bilan til topish |

**Golden hex intizomi** — poydevor qatlamning uslubi: kutilgan baytlar testda HEX-STRING konstanta sifatida turadi (binary fayl emas — diff o'qiladi, git xursand) va python bilan MUSTAQIL yasalgan (encoder o'z natijasini o'zi "golden" deb e'lon qilsa test aylanma isbotga aylanadi). Happy-path yetarli emasligi ham intizomning qismi: truncated header, length<16, NULL-terminatorsiz C-string, TLV length > qolgan baytlar — bular alohida case'lar sifatida 2–3-boblardan beri turibdi. Katta loyihalarda golden'lar `testdata/` fayllarga ko'chadi — bu pattern'ni ham namoyish qildik (`pdu/goldenfile_test.go`): `go test ./pdu -update` golden fayllarni qayta yozadi, oddiy yugurish esa solishtiradi. Kitobda asosiy golden'lar ataylab konstanta bo'lib qoldi — matn bilan yonma-yon turishi o'quv qiymati.

**net.Pipe vs real listener** tanlovining qoidasi (12-bobda boshlangan): pipe — sinxron, buffersiz, "qattiq" muhit: reader'ning har yashirin bloklanishi testni darhol qotiradi; real listener — realistic buffering, dial/close semantikasi, port bandligi. Amaliy taqsimot: codec — hech qaysi (pure function), session mexanikasi — pipe, to'liq stack — listener. Pipe'ning tuzog'i esda tursin: har ikki uchini ALOHIDA goroutine yuritishi shart, aks holda birinchi Write'da test o'zini deadlock qiladi.

## 15.2 Fuzzing: o'ylamagan kirishlar fabrikasi

Go'ning native fuzzer'i (1.18+) coverage-guided: seed corpus'dan boshlab kirishlarni mutatsiya qiladi va YANGI kod yo'lini ochgan har variantni "qiziq" deb saqlab qoladi. Bizning targetlar — tashqaridan bayt oladigan har chegara:

```go
// FuzzDecode — dispatcher va barcha PDU decoder'lari uchun fuzz target.
// Tamoyil: "Never trust incoming data" — istalgan bayt to'plami PANIC EMAS,
// yo (PDU, nil) yo xato qaytarishi kerak.
func FuzzDecode(f *testing.F) {
	seeds := []string{specBindTransmitterHex, goldenDataSMHex, ...}
	...
	f.Fuzz(func(t *testing.T, data []byte) {
		p, h, err := Decode(data)
		if err != nil {
			return // xato — normal natija
		}
		if p.Cmd() != h.ID {
			t.Fatalf("Cmd()=%s, header=%s", p.Cmd(), h.ID)
		}
	})
}
```

Uch qoida targetlarni foydali qiladi: (1) **seed'siz fuzzer ko'r** — `f.Add` bilan har PDU turidan valid namuna beriladi (bizda golden'larning o'zi!), aks holda mutatsiya 16-baytlik header strukturasigacha "yetib bormaydi"; (2) **invariant tekshiring, faqat panic emas** — FuzzGSM7RoundTrip encode muvaffaqiyatida decode AYNAN matnni qaytarishini talab qiladi (round-trip — codec'ning eng kuchli umumiy invarianti); (3) **topilgan crash regression'ga aylanadi** — fuzzer yiqilgan kirishni `testdata/fuzz/<Target>/` ga yozadi va u ENDI ODDIY `go test`da doim yuguradi.

### Fuzzer bizdan nima topdi: ikkita haqiqiy hikoya

**Hikoya 1 — haqiqiy bug (dlr.Parse, 3 soniyada).** Kirish: `"\xEAerr:"` — bitta buzuq UTF-8 bayt + valid kalit. Panic: `slice bounds out of range [7:5]`. Ildiz: parser kalitlarni registrsiz topish uchun `strings.ToLower(s)` nusxasida qidirar, qiymatlarni esa topilgan POZITSIYALAR bo'yicha ORIGINAL satrdan kesib olardi. `strings.ToLower` esa buzuq UTF-8 baytni U+FFFD belgisiga almashtiradi — u 3 BAYT: nusxaning uzunligi o'zgaradi, pozitsiyalar suriladi, kesish chegaradan chiqadi. Operator DLR matnida buzuq bayt yuborishi mumkinmi? Bemalol (9-bob: format vendor-specific, encoding kafolati yo'q) — ya'ni bu production'da OTILADIGAN panic edi. Davo — uzunlikni saqlaydigan ASCII-only kichraytirish:

```go
// lowerASCII faqat ASCII harflarni kichraytiradi — UZUNLIK SAQLANADI.
// strings.ToLower ishlatib BO'LMAYDI: buzuq UTF-8 baytlarni (operator
// yuborishi mumkin!) U+FFFD (3 bayt) bilan almashtirib pozitsiyalarni
// suradi — fuzzer topgan real panic (testdata/fuzz corpus'ida regression).
func lowerASCII(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'A' && c <= 'Z' {
			b[i] = c + ('a' - 'A')
		}
	}
	return string(b)
}
```

O'ylab ko'ring: 9-bobda parser uchun 6 xil operator formati, TLV ziddiyatlari, token-chegara case'lari — jami o'nlab test yozganmiz. BITTASI ham buzuq-UTF-8 ni o'ylamagan. Fuzzer uch soniyada o'yladi.

**Hikoya 2 — xato invariant (FuzzSplit, bir soniyada).** Fuzz target'ga "har segment ≤ 140 oktet" invariantini yozdik — mantiqiy tuyuladi: SMS payload 140 bayt-ku (7-bob). Fuzzer 145 belgilik matn bilan darhol yiqitdi: yakka GSM7 segment 145 OKTET chiqdi. Bug emas — 7-bobning o'zagi esdan chiqqan: **short_message'da GSM7 UNPACKED** (1 belgi = 1 oktet), 160 belgili xabar = 160 oktet, 140 limiti PACKED havo-interfeysiga tegishli; SMPP'dagi mutlaq chegara — sm_length 254. Invariant tuzatildi (254 mutlaq; UCS2 ≤ 140; multi-segment GSM7 ≤ 159). Saboq ikki tomonlama: fuzzer nafaqat kodni, TESTNING O'ZINI ham tekshiradi — va "o'z invariantingni himoya qilolmasang, uni noto'g'ri tushungansan".

`-race` haqidagi eslatma shu kontekstda o'z o'rnini topadi: race detector ham fuzzer'ga o'xshab faqat HAQIQATDA YUZ BERGAN muammoni ko'radi — shuning uchun concurrency testlari ko'p iteratsiya bilan (`-count=5`) va CI'da DOIM `-race` bilan yuguradi.

## 15.3 Benchmark: raqamlar bilan gaplashish

Hot path — har xabarda yuguradigan kod: encode, decode, framing. O'lchov (`go test ./pdu -bench=. -benchmem`, Apple M1 Pro):

```
BenchmarkSubmitSMEncode-10       8368974    127.8 ns/op    288 B/op    3 allocs/op
BenchmarkDecodeSubmitSM-10       6438294    186.2 ns/op    128 B/op    7 allocs/op
BenchmarkDecodeDispatcher-10     4864201    245.0 ns/op    288 B/op    8 allocs/op
BenchmarkReadFrame-10           29867276     40.0 ns/op    100 B/op    2 allocs/op
```

Bu raqamlarni KONTEKSTGA qo'yish muhim: to'liq encode+decode+framing ≈ 400ns — soniyasiga MILLIONLAB PDU degani; operator limiti esa yuzlab TPS (12-bob). Ya'ni bizning bottleneck HECH QACHON codec bo'lmaydi — u tarmoq RTT va operator TPS'i. Shu xulosaning o'zi benchmark'ning qiymati: **optimallashtirMASlik qarori ham o'lchov bilan qabul qilinadi.** sync.Pool bilan allocation'larni 3→0 ga tushirish mumkin (va buffer-reuse pattern'i sifatida mashqda bor), lekin 288 baytlik 3 allocation soniyasiga 200 marta — GC uchun sezilmas yuk; pool esa kod murakkabligi va "buffer'ni kim egallab turibdi" degan yangi savollar. Qoida: avval o'lcha, keyin (kerak bo'lsa!) optimallashtir — `-cpuprofile`/`-memprofile` + `go tool pprof` chuqurlashish yo'li tayyor turibdi.

## 15.4 Wireshark: simdagi haqiqat

Barcha testlar O'ZIMIZ yozgan kod ichida — Wireshark esa mustaqil hakam: uning SMPP dissector'i bizning frame'larni BOSHQA implementatsiya ko'zi bilan o'qiydi. Amaliy retsept:

1. localsmsc'ni ishga tushiring (14-bob), client testini yuguring, trafikni yozing: `tcpdump -i lo0 -w smpp.pcap port <PORT>` (macOS'da loopback `lo0`).
2. Wireshark'da oching. Port 2775 bo'lmagani uchun dissector avtomatik ulanmaydi: paketga o'ng tugma → **Decode As... → SMPP**.
3. Ikki preference MUHIM: SMPP'da "Reassemble SMPP over TCP messages spanning multiple TCP segments" va TCP'da "Allow subdissectors to reassemble TCP streams" — busiz katta PDU'lar (yoki bitta segmentga sig'gan bir nechta PDU) buzuq ko'rinadi. Bu 2-bobdagi framing darsimizning Wireshark'dagi aksi: u ham xuddi bizdek 16-baytlik header'dan length o'qib yig'adi.
4. Filter'lar: `smpp` (hammasi), `smpp.command_id == 0x00000004` (submit'lar), `smpp.sequence_number == 42` (bitta tranzaksiya juftligi), `smpp.command_status != 0` (xatolar).

Dissector har field'ni nomi bilan ochib beradi — o'z hex dump'laringizni u bilan solishtirish (2-bobdagi qo'lda o'qish mashqidan keyin) juda qoniqarli tajriba. Tashqi simulatorlar bilan interop ham shu bosqichning ishi: Melrose Labs onlayn simulatori (bepul, v3.4, TLS; credentials ~90 kun) yoki Docker'dagi SMPPSim'ga localsmsc o'rniga ulanib, bizning client boshqa implementatsiya bilan gaplasha olishini tekshirasiz. Qoida esa o'zgarmaydi (14-bob): avtomatlashtirilgan testlar FAQAT o'z mock'imizga tayanadi — tashqi xizmat CI'da flaky dependency; simulatorlar qo'lda interop-tekshiruv uchun.

## 15.5 ci.sh: yagona haqiqat manbai

```
$ ./ci.sh
== gofmt tekshiruvi
== go vet
== go build (examples bilan)
== go test -race
== qisqa fuzz (regression rejimi)
== HAMMASI TOZA
```

Zanjir qat'iy tartibda: format (arzon, tez yiqiladi) → vet → build (examples ham!) → to'liq test `-race` bilan → har fuzz target 3 soniyadan. Qisqa fuzz YANGI bug qidirmaydi — corpus'dagi regression'lar va seed'lar hali ham o'tishini tasdiqlaydi; chuqur qidiruv (daqiqalab) — qo'lda, vaqti-vaqti bilan. Skript lokal ham, CI serverda ham AYNAN bir xil yuguradi — "mashinamda ishlagan edi" sinfi muammolar uchun yagona davo.

## Xulosa

Test piramidasi qatlamlari aniq vazifalarga bo'lindi: golden hex (mustaqil yasalgan!) baytlarni, net.Pipe concurrency'ni, real listener integratsiyani, fuzzer esa tasavvur chegarasini qo'riqlaydi. Fuzzer ikki darsni amalda berdi: dlr.Parse'dagi haqiqiy panic (strings.ToLower'ning UTF-8 uzunlik-siljitishi — endi testdata/fuzz'da abadiy regression) va FuzzSplit invariantining o'zi xato bo'lgani (unpacked GSM7 — 140 emas!). Benchmark'lar "optimallashtirmaslik" qarorini raqam bilan asosladi: codec millionlab PDU/s beradi, bottleneck operatorda. Wireshark mustaqil hakam sifatida framing'imizni tasdiqlaydi, tashqi simulatorlar interop uchun (CI'da emas!). Va ci.sh hammasini bitta buyruqqa yig'di. Kitobning texnik qurilishi tugadi — oxirgi bob qolgan yagona savolga javob beradi: bularni PRODUCTION'da qanday yashatish kerak.

**Takrorlash savollari** (javoblar matnda bor — o'zingizni tekshiring):

1. Golden hex nega python bilan mustaqil yasaladi va nega hex-string binary fayldan yaxshi?
2. Fuzz target'ga seed bermaslik nimaga olib keladi?
3. dlr.Parse bug'ining ildizi nima edi va nega oddiy testlar uni topa olmadi?
4. FuzzSplit hikoyasi test yozuvchiga qanday saboq beradi?
5. sync.Pool'ni codec'ga qo'shmaslik qarori qanday asoslandi?
6. Wireshark'da SMPP'ni 2775 bo'lmagan portda o'qish uchun nima qilinadi va qaysi ikki preference shart?
7. Nega tashqi simulatorlar CI testlariga kiritilmaydi?
8. ci.sh'dagi qisqa fuzz bilan qo'lda chuqur fuzz'ning vazifalari qanday farqlanadi?

**Mashqlar:** [exercises/15-testing.md](../exercises/15-testing.md) — ataylab bug qo'yib fuzzer'ga topdirish, pcap tahlili va 0-allocation encode.

---

**Oldingi bob:** [14-bob. Mock SMSC](14-mock-smsc.md) · **Keyingi bob:** [16-bob. Production](16-production.md) — TLS, monitoring, PII masking, quirk katalogi va yakuniy e2e demo.

## Manbalar

- [Go Fuzzing (rasmiy hujjat)](https://go.dev/doc/security/fuzz/) — seed corpus, testdata/fuzz regression mexanizmi
- [Eli Bendersky — File-driven testing in Go](https://eli.thegreenplace.net/2022/file-driven-testing-in-go/) — testdata + -update pattern'ining kanonik izohi
- [Wireshark Wiki — SMPP](https://wiki.wireshark.org/SMPP) va [SMPP display filter reference](https://www.wireshark.org/docs/dfref/s/smpp.html) — dissector, reassembly preference, filter maydonlari
- [Melrose Labs SMSC Simulator](https://melroselabs.com/services/smsc-simulator/) va [komuw/smpp_server_docker](https://github.com/komuw/smpp_server_docker) — tashqi simulatorlar (qo'lda interop uchun)
- [sync.Pool: when it helps and when it hurts](https://harrisonsec.com/blog/go-sync-pool-buffer-reuse-when-it-helps/) — pool trade-off'lari (bizning "qo'shmaslik" qarorimizning asosi)