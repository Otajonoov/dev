# 02. Unsafe Va GC Qoidalari

Bu roadmapdagi eng muhim xavfsizlik bo'limi. Allocator yozishdan oldin shu qoidalar yod bo'lishi kerak.

## Asosiy Qoidalar

1. `unsafe.Pointer` pointer, `uintptr` esa oddiy integer.
2. `uintptr` GC root emas.
3. Go pointerini `uintptr`ga aylantirib uzoq saqlamang.
4. Pointer arithmetic uchun `unsafe.Add` ishlating.
5. Off-heap memory ichida Go heap pointer saqlamang.
6. Go objectning yashab turishini kafolatlash kerak bo'lsa `runtime.KeepAlive` ishlating.
7. `reflect.SliceHeader` va `reflect.StringHeader` bilan yangi value yasamang. Zamonaviy helperlar: `unsafe.Slice`, `unsafe.String`, `unsafe.SliceData`, `unsafe.StringData`.

## Off-Heap Qoida

Custom allocator memorysi Go GC tomonidan scan qilinmaydi.

Shuning uchun off-heap ichida quyidagilar xavfsizroq:

- `int`, `uint64`, `float64`
- fixed-size pointer-free struct
- raw bytes
- offset yoki handle

Quyidagilar xavfli:

- `*T`
- `string`
- `[]T`
- `map`
- `chan`
- `interface`

Chunki ularning ichida Go heap pointer bo'lishi mumkin.

## Handle Model

Pointer o'rniga offset saqlang:

```go
type Handle uint64

type Region struct {
    base unsafe.Pointer
    size uintptr
}

func (r *Region) Ptr(h Handle) unsafe.Pointer {
    return unsafe.Add(r.base, uintptr(h))
}
```

Bu model map, slice va channel uchun keyin juda qulay bo'ladi.

## Debug Qoidalar

Har allocator debug modega ega bo'lsin:

- allocation header
- requested size
- actual size
- alignment
- magic number
- freed flag
- canary
- poison on free
- bounds check

Minimal header:

```go
type Header struct {
    size  uintptr
    magic uint32
    state uint8
}
```

## Tekshiruv

Har unsafe package uchun quyidagilarni ishlating:

```bash
go test -race ./...
go test -gcflags=all=-d=checkptr=2 ./...
go test -run TestAllocator -count=1000 ./...
```
