# 09. Stack, Heap, Escape

Go'ning haqiqiy goroutine stackini Go ichidan boshqarmaysiz. Bu bosqichda mini VM-style stack yozasiz.

## Go Haqiqati

Go compiler escape analysis qiladi:

```bash
go build -gcflags="-m=2" ./...
```

Natija:

- escape qilmagan object stackda qolishi mumkin
- escape qilgan object heapga ketadi
- goroutine stack o'suvchi stack
- stack copy bo'lishi mumkin

## Lab Model

Mini frame stack:

```go
type Stack struct {
    base unsafe.Pointer
    size uintptr
    sp   uintptr
}

type Frame struct {
    base   uintptr
    size   uintptr
    defers []Defer
}
```

API:

```go
func (s *Stack) PushFrame(size, align uintptr) Frame
func (s *Stack) PopFrame(f Frame)
func (s *Stack) Alloc(size, align uintptr) unsafe.Pointer
```

## Heap Model

Heap allocator oldingi `04_allocator.md` dan keladi.

Escape simulator:

```go
type Lifetime uint8

const (
    Local Lifetime = iota
    Escapes
)

func AllocObject(l Lifetime, size, align uintptr) unsafe.Pointer
```

## Maqsad

Quyidagini modellashtirish:

```text
local temporary -> frame stack
returned/shared -> heap allocator
```

Bu haqiqiy Go compiler emas, lekin stack/heap qarorini tushunish uchun juda yaxshi mashq.

## Testlar

- frame LIFO pop
- frame ichidagi objectlar popdan keyin invalid
- heap object frame popdan keyin yashaydi
- stack overflow guard
- recursive frame stress
