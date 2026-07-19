# 15-bob mashqlari: Testing va tooling

> Avval mashqlarni mustaqil bajaring, keyin yechimlarga qarang. Bob matni: [book/15-testing.md](../book/15-testing.md)

---

## Mashq 1. Fuzzer'ga bug topdirish

Kodga ATAYLAB bug qo'yib, fuzzer uni topishini kuzating (tajribadan keyin qaytaring!).

1. `pdu/frame.go`dagi `ReadFrame`dan `length < HeaderSize` tekshiruvini vaqtincha olib tashlang. Qaysi fuzz target, qancha vaqtda, qanday xato bilan topadi? Yugurtirib ko'ring: `go test ./pdu -run='^$' -fuzz=FuzzReadFrame -fuzztime=30s`
2. `tlv.Decode`dagi `len(data)-off < length` tekshiruvini buzing (masalan `<=` qilib) — FuzzDecode buni sezadimi? Nega (yoki nega yo'q)?
3. Fuzzer topgan crash fayli qayerga yoziladi va u keyingi `go test`larda qanday rol o'ynaydi?
4. Xulosa: qaysi TURDAGI buglar fuzzer uchun oson, qaysilari qiyin?

## Mashq 2. O'z traffic'ingizni Wireshark'da o'qing

1. `go run ./examples/localsmsc -dlr 2s` ishga tushiring; portini yozib oling.
2. `tcpdump -i lo0 -w smpp.pcap port <PORT>` yozuvni boshlang va boshqa terminalda bitta integration test yuguring (masalan `go test ./smsc -run TestFullFlowWithDLR -count=1` — server portini testniki bilan almashtirmang, test o'z serverini ochadi: tcpdump'ni testdan keyin port'ga emas, `-i lo0 tcp` bilan yozish osonroq).
3. Wireshark'da oching, Decode As → SMPP, reassembly'ni yoqing. Kamida UCH PDU'ni toping va har biri uchun yozing: command_id (nomi), sequence_number, command_status, body'dagi 2 muhim field.
4. `smpp.command_status != 0` filtri nima ko'rsatadi? DLR deliver_sm'ini toping — short_message ichidagi receipt matnini dissector qanday ko'rsatadi?

## Mashq 3. 0-allocation encode

`BenchmarkSubmitSMEncode` 3 allocation ko'rsatadi. Ularni 0 ga tushiring — lekin PUBLIC API'ni buzmasdan.

1. Avval allocation'lar QAYERDAN kelayotganini toping: `go test ./pdu -bench=BenchmarkSubmitSMEncode -benchmem -memprofile=mem.out` + `go tool pprof -alloc_objects mem.out` (top, list).
2. `EncodeTo(buf *bytes.Buffer, seq uint32) error` uslubidagi qo'shimcha API + sync.Pool bilan yechim skeletini yozing.
3. Benchmark bilan isbotlang: yangi yo'l 0 allocs/op.
4. Va halol savol: bu optimallashtirish bizning loyihada QACHON o'zini oqlaydi? (15.3-bo'limdagi mulohazaga qarshi/yon argument keltiring.)

---

# Yechimlar

## Yechim 1

**1.** `FuzzReadFrame` odatda SONIYALAR ichida topadi: length<16 bo'lgan frame'da `length-HeaderSize` underflow bo'lib katta songa aylanadi (uint32) yoki manfiy slice size panic beradi (implementatsiyaga qarab: bizda `make([]byte, length-16)` — length=4 bo'lsa `length-16` uint32'da ~4 milliard → OOM-urinish/allocation panic). Xato: `runtime: makeslice: len out of range` yoki xotira portlashi.

**2.** Sezadi — lekin panic sifatida EMAS: `<=` bo'lganda oxirgi-baytgacha-to'liq TLV'lar xato "truncated" deb rad etiladi, ya'ni VALID kirishlar yiqila boshlaydi. Buni fuzzer emas, MAVJUD testlar (golden TLV round-trip'lar) ushlaydi. Saboq: fuzzer "crash topuvchi", regression testlar "to'g'rilik qo'riqchisi" — ikkalasi birga ishlaydi.

**3.** `testdata/fuzz/<FuzzTarget>/<hash>` fayliga (`go test fuzz v1` formatida). U seed corpus'ning qismiga aylanadi: ENDI har oddiy `go test` (fuzz'siz ham!) shu kirishni yuguradi — bug qaytsa darhol yiqiladi. Commit qilinadi — jamoaviy regression.

**4.** Oson: panic/index-range/OOM/infinite-loop — "kod yiqildi" sinfi (oracle bepul: crash o'zi signal). Qiyin: SEMANTIK xatolar ("decode noto'g'ri qiymat qaytardi, lekin yiqilmadi") — ular uchun invariant/oracle YOZISH kerak (round-trip kabi), va invariant qanchalik kuchli bo'lsa fuzzer shunchalik "aqlli".

## Yechim 2

**3.** (Namuna — sizning capture'ingizda raqamlar farq qiladi.) bind_transceiver: command_id=0x00000009, seq=1, status=0 (request'da doim 0!), system_id="esme1", interface_version=0x34. submit_sm: 0x00000004, seq=2, dest_addr="998901234567", registered_delivery=0x01. deliver_sm (DLR): 0x00000005, server seq=1, esm_class=0x04, short_message="id:... stat:DELIVRD...".

**4.** `smpp.command_status != 0` faqat XATO resp'larni ko'rsatadi (toza oqimda bo'sh; ThrottleEveryN quirk bilan yugursangiz RTHROTTLED'lar chiqadi). DLR'ning short_message'ini dissector "Message" sifatida xom ko'rsatadi — receipt MATNI SMPP strukturasi emas (Appendix B normativ emas — 9-bob!), shuning uchun dissector uni field'larga ochmaydi; TLV'lar (receipted_message_id, message_state) esa alohida, nomlari bilan ochiladi.

## Yechim 3

**1.** pprof uchtasini ko'rsatadi: `bytes.Buffer` ichki o'sishi (body yig'ishda), `encodePDU`dagi natija slice, va header yozuvidagi kichik buffer. **2.** Skelet:

```go
var bufPool = sync.Pool{New: func() any { return new(bytes.Buffer) }}

// EncodeTo sm'ni buf ustiga yozadi — allocation'siz yo'l (buf Pool'dan).
func (s SubmitSM) EncodeTo(buf *bytes.Buffer, seq uint32) error {
	buf.Reset()
	// encodePDU mantig'i buf ustida: header uchun joy band qilinadi,
	// body yoziladi, length orqaga patch qilinadi.
	...
}

// Chaqiruvchi tomonda:
buf := bufPool.Get().(*bytes.Buffer)
defer bufPool.Put(buf)
sm.EncodeTo(buf, seq)
conn.Write(buf.Bytes()) // Write'dan KEYIN Put - baytlar egasi aniq!
```

**3.** `BenchmarkSubmitSMEncodeTo` bilan: 0 allocs/op (birinchi iteratsiyalardan keyin — pool isishi). **4.** Oqlanadi: bitta protsess o'nlab-yuzlab MING PDU/s qayta ishlasa (masalan SMPP proxy/router yozsangiz) — GC pause'lar latency'da seziladi. Bizning gateway profilida (yuzlab TPS) — oqlanmaydi: murakkablik (+"buffer'ni Put'dan keyin ishlatish" bug sinfi) foydadan katta. Ikkala javob ham RAQAMDAN keladi — bu mashqning asl maqsadi.
