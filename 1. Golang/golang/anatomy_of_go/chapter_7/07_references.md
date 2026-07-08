# 07 — Manbalar (References)

> **The Anatomy of Go** (Phuong Le) — 7-bob "Memory" uchun manbalar.
> Bu ro'yxat kitobdagi asl havolalarni va Go xotira modeli hamda GC'ni chuqurroq o'rganish uchun qo'shimcha rasmiy manbalarni o'z ichiga oladi.

[← 06 Xulosa](06_summary.md)

---

## Kitobdagi asl havolalar

Kitob 7-bobda quyidagi manbalarga havola qiladi:

| Belgi | Nomi | Havola |
|---|---|---|
| `[go1.22-darwin]` | Go 1.22 Release Notes / Ports / Darwin | https://tip.golang.org/doc/go1.22#darwin |
| `[ctgst]` | Contiguous Stacks (dizayn hujjati) | https://docs.google.com/document/d/1wAaf1rYoM4S4gtnPh0zOlGzWtrZFQ5suE8qr2sD8uWQ/pub |
| `[fptr]` | Reducing Go Execution Tracer Overhead With Frame Pointer Unwinding | https://blog.felixge.de/reducing-gos-execution-tracer-overhead-with-frame-pointer-unwinding/ |
| `[tiny]` | runtime: mark tiny blocks at GC start (Gerrit CL) | https://go-review.googlesource.com/c/go/+/31456 |

---

## Rasmiy Go hujjatlari

### GC va xotira boshqaruvi asoslari

- **A Guide to the Go Garbage Collector** — GC'ning eng to'liq rasmiy qo'llanmasi (GOGC, GOMEMLIMIT, pacing, latency): https://tip.golang.org/doc/gc-guide
- **Go Memory Model** — xotira modeli, happens-before, sinxronizatsiya: https://go.dev/ref/mem
- **`runtime` paketi hujjati** — `runtime.GC`, `runtime.SetFinalizer`, `GOGC`, `GOMEMLIMIT`, `MemStats`: https://pkg.go.dev/runtime
- **`runtime/debug` paketi** — `debug.SetGCPercent`, `debug.SetMemoryLimit`, `debug.ReadGCStats`, `debug.FreeOSMemory`: https://pkg.go.dev/runtime/debug

### GOMEMLIMIT

- **Go 1.19 Release Notes** — `GOMEMLIMIT` tanishtirilgan versiya: https://go.dev/doc/go1.19
- **`GOMEMLIMIT` uchun soft memory limit proposal** (Michael Knyszek): https://github.com/golang/proposal/blob/master/design/48409-soft-memory-limit.md

---

## Go blog maqolalari

- **Getting to Go: The Journey of Go's Garbage Collector** (Rick Hudson) — Go GC'ning evolyutsiyasi, STW'dan concurrent'gacha, latency maqsadlari: https://go.dev/blog/ismmkeynote
- **Go GC: Prioritizing low latency and simplicity** (2015) — Go 1.5 concurrent GC e'loni: https://go.dev/blog/go15gc
- **Go's Memory Allocator** (background) — `tcmalloc`ga asoslangan allokator g'oyalari.

---

## Runtime manba kodi (src/runtime)

GC va xotira boshqaruvini kodda o'rganish uchun eng muhim fayllar (Go rasmiy repozitoriyasi):

| Fayl | Mavzu |
|---|---|
| `src/runtime/mgc.go` | GC'ning asosiy tsikli, fazalar, `gcStart`, `gcMarkTermination` |
| `src/runtime/mgcpacer.go` | Pacer: `heapGoal`, `trigger`, `consMark`, `gcBackgroundUtilization` |
| `src/runtime/mgcmark.go` | Marking: root skanerlash, `scanobject`, oblet, assist |
| `src/runtime/mgcsweep.go` | Sweeping: `sweepone`, `bgsweep`, `sweepgen` |
| `src/runtime/mgcscavenge.go` | Scavenger: bo'sh sahifalarni OS'ga qaytarish |
| `src/runtime/mbarrier.go` | Write barrier: `wbMove`, `wbZero`, hybrid barrier |
| `src/runtime/malloc.go` | `mallocgc` — asosiy allokator, `gcmarknewobject` |
| `src/runtime/mheap.go` | `mheap`, `mspan`, `mcentral`, `mheap.reclaim` (span reclaimer) |
| `src/runtime/mcache.go` | Per-P kesh, tiny allocator |
| `src/runtime/mcentral.go` | `mcentral` struktura, `partialSwept`/`fullUnswept` |
| `src/runtime/stack.go` | Stek o'sishi: `morestack`, `newstack`, `copystack` |
| `src/runtime/mgcwork.go` | GC work buferlari: `gcWork`, `wbuf1`/`wbuf2`, global full/empty |

> Manba: https://github.com/golang/go/tree/master/src/runtime

---

## Kompilyator va escape analysis

- **`go build -gcflags="-m"`** — escape analysis qarorlarini ko'rish (`-m=2` batafsilroq).
- **`src/cmd/compile/internal/escape/`** — escape analysis implementatsiyasi.
- **`src/cmd/compile/internal/ssa/`** — SSA pass'lar, jumladan `writebarrier` pass (`compile.go`).

---

## Diagnostika va profiling asboblari

- **`GODEBUG=gctrace=1`** — har GC siklidan keyin trace chop etadi (heap goal, marking vaqti, STW pauzalar).
- **`GODEBUG=scavtrace=1`** — scavenger faoliyatini kuzatish.
- **`runtime/pprof`** va **`go tool pprof`** — heap profil (allokatsiya, in-use xotira). 6-bobda batafsil.
- **`go tool trace`** — GC worker'lar (dedicated/fractional/idle) va STW pauzalarni timeline'da ko'rish.
- **`runtime.ReadMemStats`** — dasturiy tarzda `HeapAlloc`, `HeapSys`, `NextGC`, `PauseNs` kabi statistikalarni o'qish.

---

## Nazariy asoslar (write barrier va tri-color)

- **Dijkstra et al. (1978)** — "On-the-Fly Garbage Collection: An Exercise in Cooperation" — tri-color abstraksiyasi va insertion barrier'ning ilk manbasi.
- **Yuasa (1990)** — "Real-time garbage collection on general-purpose machines" — deletion (snapshot-at-the-beginning) barrier.
- **The Garbage Collection Handbook** (Jones, Hosking, Moss) — GC nazariyasining eng to'liq akademik manbasi.

---

## Bog'liq boblar (shu kitob ichida)

- **5-bob** — Compiler: SSA pass'lar, `writebarrier` pass qayerda ishlashi.
- **6-bob** — Functionality: MPG modeli, stek freymlari, profiling va trace (GC'ni kuzatish uchun).
- **8-bob** — Scheduler: preemption (async preemption GC uchun muhim), `sudog`, run queue'lar.

---

[← 06 Xulosa](06_summary.md)
