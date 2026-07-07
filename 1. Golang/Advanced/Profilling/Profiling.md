# 5. Profiling

**pprof** — Go'da performance tahlili uchun asosiy vosita. U ikki qismdan iborat: 
1. `runtime/pprof` (ma'lumot to'plash) 
2. `go tool pprof` (tahlil va vizualizatsiya). 
Profiler dastur ishlayotgan vaqtda stack trace'lar yig'adi va har bir funksiya qancha resurs (CPU, memory, kutish vaqti) sarflayotganini ko'rsatadi.

---

### O'rnatilgan profiler'lar

To'rt guruhga bo'linadi:

- **Memory-sampling** (heap, allocs) — standart holatda yoqilgan. `MemProfileRate` orqali boshqariladi (standart 512 KB). Profiler har taxminan 512 KB ajratilgan memory uchun bitta namuna yozadi. Sampling intervali eksponensial taqsimot bilan tasodifiy tanlanadi — ba'zan 120 KB, ba'zan 1.3 MB, lekin o'rtacha 512 KB atrofida. Namuna olinganda, Go natijani scaling qiladi: 2 KB allocation namuna olinsa, profiler uni ~256.5 ga ko'paytirib haqiqiy hajmga yaqinlashtiradi. heap va allocs bir xil ma'lumotdan foydalanadi, faqat standart ko'rinish farq qiladi — heap `inuse_space`, allocs esa `alloc_space` ko'rsatadi.
- **Concurrency** (goroutine, threadcreate) — doimo faol, runtime tomonidan yangilanib turadi, deyarli nol overhead. Goroutine profiler stop-the-world pauza bilan barcha faol goroutine'larning stack trace'ini oladi. Threadcreate profiler esa har bir OS thread yaratilgan nuqtadagi stack trace'ni `m.createstack` maydonidan yig'adi, lekin amalda kamdan-kam foydali, chunki ko'pchilik thread'lar scheduler tomonidan yaratiladi.
- **CPU** — faqat `StartCPUProfile`/`StopCPUProfile` orasida ishlaydi. Unix'da `ITIMER_PROF` taymeri 100 Hz da SIGPROF signal yuboradi. Signal handler to'xtatilgan thread'ning stack trace'ini P buffer'iga yozadi. Alohida `profileWriter` goroutine har 100 ms da buffer'larni o'qib, bir xil stack trace'larni bitta yozuvga birlashtiradi. Linux'da thread-ga xos taymerlar ham qo'llab-quvvatlanadi — har bir thread o'z SIGPROF'ini oladi, bu aniqroq natija beradi. Call stack chuqurligi 64 frame bilan cheklangan.
- **Contention** (block, mutex) — standart holatda o'chirilgan. Mutex profiler `SetMutexProfileFraction(n)` bilan yoqiladi — har n ta contention hodisasidan bittasi yoziladi. Block profiler `SetBlockProfileRate` bilan yoqiladi — vaqt-asoslangan threshold ishlaydi: threshold'dan uzun block doimo yoziladi, qisqasi ehtimollik bilan.

---

### Profile to'plash usullari

**`go test` flag'lari** — `-cpuprofile`, `-memprofile`, `-blockprofile`, `-mutexprofile` flag'lari orqali testlar va benchmark'lar davomida profile to'planadi. `testing` paketi `M.before()` da profiler'larni sozlaydi, `M.after()` da to'xtatib diskka yozadi.

**`runtime/pprof` API** — `pprof.Lookup("heap")` bilan profile'ga kirish, `WriteTo(w, debug)` bilan yozish. 
- `debug=0` binary protobuf
- `debug=1` inson o'qiy oladigan matn
- `debug=2` (faqat goroutine) to'liq stack dump beradi
CPU profile uchun `pprof.StartCPUProfile(w)` va `pprof.StopCPUProfile()` ishlatiladi.

**`net/http/pprof`** — `import _ "net/http/pprof"` qo'shish bilan `/debug/pprof/` endpoint'lar avtomatik ro'yxatdan o'tadi. CPU profile uchun handler `StartCPUProfile` chaqirib 30 soniya (yoki `seconds=N`) kutadi. Production'da bu endpoint'larni himoya qilish kerak. Custom `ServeMux` ishlatilsa, handler'larni qo'lda ro'yxatdan o'tkazish kerak.

---

### `go tool pprof`

Interaktiv shell buyruqlari:

- **`top`** — eng ko'p resurs ishlatayotgan funksiyalarni ko'rsatadi. `flat` — funksiyaning o'zi sarflagan resurs, `cum` — o'zi + chaqirgan funksiyalari bilan birga. `-cum` flag cumulative bo'yicha saralaydi.
- **`list`** — funksiya manba kodini performance ma'lumotlari bilan qator-qator ko'rsatadi. Aniq qaysi qator qimmat ekanligini topish mumkin.
- **`peek`** — funksiyaning to'g'ridan-to'g'ri chaqiruvchilari va chaqiriluvchilarini ko'rsatadi, inline belgilarini ham ko'rsatadi.
- **`disasm`** — assembly instruksiyalar bilan performance ma'lumotlarini ko'rsatadi. Binary executable fayl kerak. Inlined funksiyalar alohida disassemble qilinmaydi.
- **`web`** — Graphviz yordamida SVG call graph yaratadi va brauzerda ochadi.

**HTTP UI rejimi** (`-http=:8080`) — flamegraph, call graph, top jadvallar, manba kod va disassembly ko'rinishlarini interaktiv brauzer interfeysida taqdim etadi.

---

### Call Graph va Flamegraph o'qish

**Call Graph:** Har bir node funksiya, hajmi flat qiymatga, rangi resurs ishlatishga (yashildan qizilga) asoslangan. To'liq edge'lar to'g'ridan-to'g'ri chaqiruvni, nuqtali edge'lar bilvosita chaqiruvni ko'rsatadi. Edge qalinligi resurs oqimiga mutanosib. Nodelet'lar — node'larga biriktirilgan kichik qutichalar bo'lib, allocation hajmlari kabi qo'shimcha metadata ko'rsatadi, Go'ning size class'lariga mos keladi.

**Flamegraph:** Har bir funksiya gorizontal chiziq, kengligi resurs ishlatishiga mutanosib. Vertikal ustun call stack'ni ifodalaydi — tepa ildiz, pastki qism leaf funksiya. Agar bola quti ota qutiga teng kenglikda bo'lsa, ko'pchilik ish bolada sodir bo'layapti. Ranglar funksiyalarni vizual ajratish uchun, performance ko'rsatmaydi.

---

### Mutex vs Block profiling

**Mutex profiler** — lock'ni **ushlovchi** (unlock qiluvchi) nuqtai nazaridan ishlaydi. "Kim lock'ni boshqalarni blokirovka qilish uchun uzoq ushlab turayapti?" savoliga javob beradi. Contention vaqti — boshqa goroutine'larning kutish yig'indisi sifatida aybni lock ushlovchiga yuklaydi. `defer` ishlatilganda turli mutex'lar bir xil stack trace ostida guruhlanishi mumkin.

**Block profiler** — **kutuvchi** nuqtai nazaridan ishlaydi. Channel, select, sync.Mutex, WaitGroup, I/O operatsiyalari kabi barcha synchronization primitive'larida goroutine qancha vaqt bloklangan bo'lishini yozadi. Sampling vaqt-asoslangan threshold bilan ishlaydi — uzoq block'lar doimo yoziladi, qisqalari ehtimollik bilan.

---

### Memory profiler'ning uch-tsikl tizimi

Allocation'lar darhol hisobga olinadi, lekin bo'shatishlar faqat GC paytida sodir bo'ladi. Bu noto'g'ri raqamlarga olib kelishi mumkin. Buni hal qilish uchun profiler uch cycle'li halqa buffer ishlatadi (C, C+1, C+2). Allocation'lar C+2 ga, deallocation'lar C+1 ga yoziladi. GC tugagandan keyin C+1 ma'lumotlari active profile'ga ko'chiriladi. Tashqi vositalar faqat active profile'ni ko'radi — bu barqaror va izchil ko'rinish beradi.

---

### Delta Profiling

Ikki suratni solishtirish orqali faqat farqni ko'rsatadi. Har bir stack trace uchun ikkinchi suratdan birinchisi ayiriladi. Musbat raqamlar o'sishni, salbiy raqamlar kamayishni bildiradi. HTTP endpoint'larga `seconds=N` qo'shib delta profile'ni avtomatik olish mumkin — handler boshlang'ich profile yozadi, N soniya kutadi, ikkinchi profile yozadi va farqni qaytaradi. Natija oddiy pprof fayl, barcha tahlil buyruqlari ishlaydi.

---





![](../../../assets/obsidian-images/Pasted%20image%2020260408212924.png)

![](../../../assets/obsidian-images/Pasted%20image%2020260408213056.png)

---
