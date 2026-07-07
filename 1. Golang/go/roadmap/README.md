# Go Runtime Lab Roadmap

Bu papkadagi roadmapning maqsadi: Go runtime mexanizmlarini o'rganish uchun ularning kichik, tajribaviy nusxalarini Go'da yozish.

Bu compiler roadmap emas. Bu yerda asosiy yo'nalish:

- off-heap memory region
- allocator: tiny, small, large
- array, slice, string
- map
- channel
- struct layout, interface layout, generics modeli
- stack/heap, escape analysis
- defer, panic, recover
- toy GC: mark-sweep, tri-color, write barrier
- GMP scheduler simulator
- test, fuzz, benchmark, profiling

## Asosiy Yo'l

Yangi, yagona roadmap:

- [runtime_lab/README.md](runtime_lab/README.md)

Oldingi ikki yo'l bitta roadmapga birlashtirildi. Endi alohida alternativ til yo'li yo'q, chunki bu loyiha Go ichida Go runtime mexanizmlarini laboratoriya sifatida yozishga qaratilgan.

Oldingi manual-memory yo'ldan faqat foydali fikrlar olib qolindi:

- har bir data structure allocator qabul qilishi
- ownership va lifetime haqida aniq o'ylash
- arena/pool orqali manual memory discipline
- cleanupni `defer` bilan tartibli qilish
- OS/C darajasidagi memory tushunchalarini Go tajribalariga bog'lash

## Qisqa Start

1. [00_scope.md](runtime_lab/00_scope.md) ni o'qing.
2. [02_unsafe_gc_rules.md](runtime_lab/02_unsafe_gc_rules.md) ni yaxshi tushuning.
3. [03_memory_region.md](runtime_lab/03_memory_region.md) bo'yicha `mmap` region yozing.
4. [04_allocator.md](runtime_lab/04_allocator.md) bo'yicha bump allocator va free listdan boshlang.
5. Keyin slice/string, map, channelga o'ting.

## Texnik Tanlov

- allocator logikasi: Pure Go + `unsafe`
- test backend: `[]byte`
- real off-heap backend: `golang.org/x/sys/unix.Mmap`
- `cgo`: bu roadmapda ishlatilmaydi
- `syscall`: faqat tushuncha uchun, yangi kodda `x/sys/unix` afzal
