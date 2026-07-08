# 05. Array, Slice, String

Bu bosqichda allocator ustida Go composite typelarining mini modelini yozasiz.

## Array

Array fixed-size memory:

```go
type Array[T any] struct {
    data unsafe.Pointer
    len  int
}
```

API:

```go
func NewArray[T any](a Allocator, n int) Array[T]
func (a Array[T]) Get(i int) T
func (a Array[T]) Set(i int, v T)
func (a Array[T]) Len() int
```

Test:

- bounds check
- element size
- alignment
- zero value

## Slice

Slice header:

```go
type Slice[T any] struct {
    data unsafe.Pointer
    len  int
    cap  int
    alloc Allocator
}
```

API:

```go
func Make[T any](a Allocator, len, cap int) Slice[T]
func (s Slice[T]) Get(i int) T
func (s Slice[T]) Set(i int, v T)
func (s Slice[T]) Append(v T) Slice[T]
func (s Slice[T]) Slice(lo, hi int) Slice[T]
func (s Slice[T]) Free()
```

## Growth

Go growth modeliga yaqinlashtiring:

- kichik caplarda 2x
- 256 dan keyin asta-sekin 1.25x tomonga
- element sizega qarab allocation roundup

Soddalashtirilgan formula:

```go
func nextCap(newLen, oldCap int) int {
    newCap := oldCap
    double := newCap + newCap
    if newLen > double {
        return newLen
    }
    if oldCap < 256 {
        return double
    }
    for newCap < newLen {
        newCap += (newCap + 3*256) >> 2
    }
    return newCap
}
```

## String

String immutable bo'lsin:

```go
type String struct {
    data unsafe.Pointer
    len  int
    alloc Allocator
}
```

API:

```go
func NewString(a Allocator, b []byte) String
func (s String) Len() int
func (s String) Byte(i int) byte
func (s String) GoString() string
func (s String) Free()
```

`GoString()` copy qilsin. Off-heap memoryni zero-copy Go stringga aylantirish xavfli, chunki lifetime boshqaruvi sizda bo'ladi.

## Muhim Cheklov

Birinchi versiyada `T` pointer-free bo'lsin. Buni dokumentatsiyada aniq ayting.

Keyingi versiyada:

- `TypeDesc`
- pointer mask
- custom mark/sweep simulator

## Testlar

- append capacity o'sishi
- appenddan keyin eski slice o'zgarmasligi kerak bo'lgan holatlar
- slicing bounds
- string immutability
- allocator stats
- fuzz: append/get standard slice bilan solishtiriladi
