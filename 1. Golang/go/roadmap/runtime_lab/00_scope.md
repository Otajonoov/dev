# 00. Scope

## Nima Quramiz

Bu roadmapda Go runtime konseptlarining kichik nusxasini yozamiz:

- OS dan memory olish: `mmap` yoki `x/sys/unix`
- off-heap arena
- allocator: bump, free list, size class, tiny, small, large
- runtime-style slice/string/map/channel
- type descriptor, struct layout, interface layout
- mini stack/frame modeli
- defer/panic/recover simulator
- toy GC: mark-sweep, tri-color, write barrier
- GMP scheduler simulator: G, M, P, run queues, work stealing
- profiling va benchmark

## Nima Qurmaymiz

- Go compiler yozmaymiz
- Go runtime'ni almashtirmaymiz
- Go GC'ni o'chirmaymiz
- haqiqiy goroutine stackni boshqarmaymiz
- haqiqiy Go scheduler/GMP'ni almashtirmaymiz
- production allocator yozmaymiz
- `runtime` internal APIlariga tayanmaymiz
- `cgo` orqali `malloc/free`ga tayanmaymiz

## Asosiy Chegara

Go ichida yozilgan kodni Go runtime boshqaradi. Shuning uchun haqiqiy Go heap/stack sizning qo'lingizda emas.

Biz boshqaradigan narsa:

```text
OS memory -> mmap region -> custom allocator -> custom data structures
```

Go boshqaradigan narsa:

```text
goroutine stack, Go heap, GC, scheduler, compiler escape analysis
```

## Dizayn Tamoyillari

- Avval to'g'ri model, keyin tezlik.
- Har bir data structure allocator qabul qilsin.
- Off-heap ichida Go pointer saqlashdan qoching.
- Har bir unsafe operatsiyaga test yozing.
- Har bosqichda standard Go behavior bilan solishtiring.
- Debug mode production modedan muhimroq: canary, poison, bounds check, stats.

## Manual-Memory Yo'ldan Olib Qolingan Fikrlar

Alternativ til yo'li alohida kerak emas, lekin undan quyidagi fikrlar foydali:

- allocator explicit bo'lishi kerak
- memory lifetime aniq ko'rinishi kerak
- `defer` cleanup uchun intizom beradi
- arena vaqtinchalik objectlar uchun kuchli model
- manual memory bilan ishlaganda ASan/Valgrind mental modeli foydali

Bu fikrlarni Go'da `unsafe`, `mmap`, test va profiling orqali ishlatamiz.

## Texnik Yo'l

```text
allocator logic: Pure Go + unsafe
test memory:      []byte
real memory:      golang.org/x/sys/unix.Mmap
cgo:              ishlatilmaydi
```
