# 04. Allocator

Allocator bu roadmapning yuragi. Boshqa hamma narsa shu ustida quriladi.

Allocator logikasi Pure Go + `unsafe` bilan yoziladi. Memory manbasi esa `Region` bo'ladi: boshlanishda `[]byte`, keyin `unix.Mmap`. Bu bosqichda `cgo` ishlatilmaydi.

## Umumiy Interfeys

```go
type Allocator interface {
    Alloc(size, align uintptr) unsafe.Pointer
    Free(ptr unsafe.Pointer, size uintptr)
    Reset()
    Stats() Stats
}

type Stats struct {
    TotalAllocated uintptr
    InUse          uintptr
    NumAllocs      uint64
    NumFrees       uint64
}
```

## Bosqich 1: Bump Allocator

Maqsad: eng sodda allocator.

Xususiyat:

- faqat oldinga yuradi
- individual `Free` yo'q
- `Reset` hammasini bo'shatadi
- alignment to'g'ri bo'lishi shart

Test:

- allocation pointerlari overlap qilmaydi
- alignment to'g'ri
- out of memory nil qaytaradi yoki panic qiladi
- resetdan keyin memory qayta ishlatiladi

## Bosqich 2: Free List Allocator

Maqsad: variable-size allocation.

Blok header:

```go
type block struct {
    size uintptr
    free bool
    next *block
    prev *block
}
```

O'rganiladigan mavzular:

- split
- coalesce
- first fit
- best fit
- fragmentation
- header overhead

## Bosqich 3: Pool Allocator

Maqsad: bir xil size blocklar.

```go
type Pool struct {
    slotSize uintptr
    free     unsafe.Pointer
}
```

Free blockning birinchi wordi keyingi free block pointeri yoki offseti bo'ladi.

## Bosqich 4: Size Classes

Go runtime small objectlarni size classlarga ajratadi. Siz mini variant yozing:

```text
8, 16, 24, 32, 48, 64, 80, 96, 112, 128,
160, 192, 224, 256, 320, 384, 448, 512,
640, 768, 896, 1024, ...
```

Maqsad:

- requested size -> class
- class -> pool/slab
- internal fragmentation o'lchash

## Bosqich 5: Tiny, Small, Large

Runtime-lab allocator:

```text
tiny  <= 16B, pointer-free, explicit free yo'q
small <= 32KB, size class orqali
large >  32KB, alohida mmap yoki region chunk
```

Tiny allocator:

- bitta 16B block ichida kichik objectlar
- faqat pointer-free data
- debugda tiny stats ko'rsin

Small allocator:

- size class
- slab/span
- bitmap yoki free list

Large allocator:

- page aligned
- alohida region
- `munmap` yoki large free list

## Go Runtime Bilan Bog'lash

Go runtime allocatorida asosiy model:

```text
mcache -> mcentral -> mheap -> OS
```

Labda soddalashtirilgan model:

```text
local cache -> central size classes -> region/page allocator
```

## Keyingi Bosqichga O'tish Sharti

Quyidagilar tayyor bo'lmaguncha slice/stringga o'tmang:

- `Alloc`
- `Free`
- `Reset`
- `Stats`
- debug poison
- alignment tests
- fragmentation test
- benchmark: tiny/small/large
