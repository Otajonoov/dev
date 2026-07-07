# 9. Manbalar va qo'shimcha o'qish

> Ushbu material — Anatomy of Go kitobining 6-bobi mavzulari asosida o'zbek tilida tayyorlangan o'quv qo'llanma.

## Asl manba

- **The Anatomy of Go** — Phuong Le. Bu material shu kitobning 6-bobi ("Functionality") asosida tayyorlangan.

## Asl havolalar (kitobdan)

| Havola | Mavzu |
|--------|-------|
| [Graceful Shutdown in Go: Practical Patterns](https://victoriametrics.com/blog/go-graceful-shutdown/) | Graceful shutdown shabloni |
| [GitHub: threadcreate profile is broken](https://github.com/golang/go/issues/6104) | threadcreate profile xatoligi |
| [Go commit: PGO inlining](https://github.com/golang/go/commit/99862cd57d) | PGO ning birinchi commit'i |
| [Go.dev blog: More powerful Go execution traces](https://go.dev/blog/execution-traces-2024) | Execution trace yangiliklari (Knyszek) |
| [Go Execution Tracer (Vyukov)](https://docs.google.com/document/u/1/d/1FP5apqzBgr7ahCCgFO-yoVhk4YZrNIDNf9RybngBc14/pub) | Execution tracer arxitekturasi |
| [Profile-guided optimization](https://go.dev/doc/pgo) | PGO rasmiy hujjatlari |

## Go rasmiy hujjatlari

### Funksiyalar va closure'lar

- [Go Spec: Function declarations](https://go.dev/ref/spec#Function_declarations)
- [Go Spec: Function literals](https://go.dev/ref/spec#Function_literals)
- [Go FAQ: Why does Go not have function overloading?](https://go.dev/doc/faq#overloading)
- [Effective Go: Functions](https://go.dev/doc/effective_go#functions)

### Defer, Panic, Recover

- [Go Blog: Defer, Panic, and Recover](https://go.dev/blog/defer-panic-and-recover)
- [Go Spec: Defer statements](https://go.dev/ref/spec#Defer_statements)
- [Go Spec: Run-time panics](https://go.dev/ref/spec#Run_time_panics)

### Profiling va Tracing

- [Profiling Go Programs](https://go.dev/blog/pprof) — Russ Cox'ning klassik blog'i
- [runtime/pprof package](https://pkg.go.dev/runtime/pprof)
- [net/http/pprof package](https://pkg.go.dev/net/http/pprof)
- [runtime/trace package](https://pkg.go.dev/runtime/trace)
- [Diagnostics](https://go.dev/doc/diagnostics) — barcha tahlil asboblari

### PGO

- [PGO User Guide](https://go.dev/doc/pgo) — to'liq qo'llanma
- [PGO Design Doc](https://github.com/golang/proposal/blob/master/design/55022-pgo.md)

## Yaxshi blog'lar va maqolalar

### Dave Cheney (Go ekspert)

- [High Performance Go Workshop](https://dave.cheney.net/high-performance-go-workshop/dotgo-paris.html) — performance bo'yicha eng yaxshi material
- [Practical Go: Real World Advice for writing Maintainable Go Programs](https://dave.cheney.net/practical-go) — yaxshi amaliyotlar

### Bill Kennedy (Ardan Labs)

- [Going Deeper](https://www.ardanlabs.com/blog/) — Go internals haqida ko'p maqolalar
- Keep an eye on YouTube — Go scheduler haqida ajoyib videolari

### Michael Knyszek (Go team)

- [Memory profiling in Go](https://go.dev/blog/pprof) — profile ichki tuzilishi
- Trace bo'yicha ko'p material

## Real Go runtime kod o'qish

Agar siz juda chuqur tushunmoqchi bo'lsangiz, Go'ning o'z runtime kodini o'qishni tavsiya qilaman:

### Defer

- `src/runtime/panic.go` — `deferproc`, `deferprocStack`, `deferreturn`, `gopanic`, `gorecover`
- `src/runtime/runtime2.go` — `_defer`, `_panic`, `g`, `p` strukturalari
- `src/cmd/compile/internal/ssagen/ssa.go` — defer SSA generation
- `src/cmd/compile/internal/escape/call.go` — `goDeferStmt` (escape analysis)

### Panic & Recover

- `src/runtime/panic.go` — `gopanic`, `gorecover`, `recovery`, `nextDefer`
- `src/runtime/asm_arm64.s`, `src/runtime/asm_amd64.s` — `gogo`, `mcall`, `return0`

### Profiling

- `src/runtime/pprof/pprof.go` — Profile API
- `src/runtime/mprof.go` — runtime memory profiling
- `src/runtime/cpuprof.go` — CPU profile
- `src/runtime/trace.go` — execution tracer

### Funksiya qiymatlari

- `src/runtime/runtime2.go` — `funcval` strukturasi
- `src/cmd/compile/internal/ir/func.go` — `FuncSymName`

## Kichik amaliy qo'llanmalar

### "Profile your code" qadam-baqadam

```bash
# 1. Server'ni pprof bilan ishga tushiring
import _ "net/http/pprof"
go http.ListenAndServe("localhost:6060", nil)

# 2. CPU profile
$ go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# 3. Heap profile
$ go tool pprof http://localhost:6060/debug/pprof/heap

# 4. Goroutine profile
$ curl http://localhost:6060/debug/pprof/goroutine?debug=2 > goroutines.txt

# 5. Trace
$ curl http://localhost:6060/debug/pprof/trace?seconds=5 > trace.out
$ go tool trace trace.out
```

### "Use PGO" qadam-baqadam

```bash
# 1. Production'dan profile yig'ing (yoki test'da)
$ curl https://prod.app/debug/pprof/profile?seconds=120 > default.pgo

# 2. Repository'ga commit qiling
$ git add default.pgo
$ git commit -m "Add PGO profile"

# 3. Build qiling (avtomatik aniqlanadi)
$ go build .

# 4. Vaqti-vaqti bilan yangilang (3-6 oyda bir)
```

## Kitoblar (qo'shimcha o'qish)

### Go bo'yicha

- **The Go Programming Language** — Donovan, Kernighan. Klassik
- **100 Go Mistakes and How to Avoid Them** — Teiva Harsanyi. Real xatolar
- **Concurrency in Go** — Katherine Cox-Buday. Goroutine, channels chuqurroq

### Performance bo'yicha

- **Systems Performance** — Brendan Gregg. Linux performance umumiy
- **High Performance Browser Networking** — Ilya Grigorik. Network performance
- **Computer Systems: A Programmer's Perspective** — Bryant, O'Hallaron. Stack frame, registr — past darajada

## Video resurslar

- **GopherCon talks** — YouTube'da har xil yil
- **dotGo conference** — Yevropa Go konferensiya
- **Bill Kennedy YouTube** — Go runtime haqida
- **Damian Gryski (Twitter @damianm)** — performance va profiling

## Kamayubchilarga keng tarqalgan asboblar

| Asbob | Vazifa |
|-------|--------|
| [`pprof`](https://github.com/google/pprof) | Profile tahlil (Go'da o'rnatilgan) |
| [`go-perfbook`](https://github.com/dgryski/go-perfbook) | Performance tutoriallar to'plami |
| [`benchstat`](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat) | Benchmark natijalarini taqqoslash |
| [`flamegraph`](https://github.com/brendangregg/FlameGraph) | Flamegraph yaratish |
| [`go-torch`](https://github.com/uber-archive/go-torch) | Flamegraph (eski, hozir pprof'ga integrated) |
| [`gops`](https://github.com/google/gops) | Ishlamoqda Go jarayonlarni ko'rish |

## Kichik so'z

Bu bob Go'ning past darajadagi mexanizmlarini chuqur o'rganib chiqdi. Shu narsalarni bilish sizga **professional Go developer** bo'lishga yordam beradi:

- Sizning kodingiz **nega** ishlaydi (yoki ishlamaydi)
- **Qaerda** muammo bor
- **Qanday qilib** tezlashtirish mumkin

Keyingi kitoblarda yana chuqurroq materiallar — memory, GC, scheduler, concurrency primitivlari — bularning hammasi shu bo'limda o'rgangan tushunchalarga tayanadi.

**Quvonchbek, omadlar tilayman! Go o'rganishda davom et.**

---

**Avvalgi mavzu:** [08_summary.md](08_summary.md) — Bobning umumiy xulosasi
**Bo'lim asosiga qaytish:** [README.md](README.md)
