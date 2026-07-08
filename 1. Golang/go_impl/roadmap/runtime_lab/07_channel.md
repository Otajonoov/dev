# 07. Channel

Channel Go runtime ichidagi eng qiziq mexanizmlardan biri. Bu bosqichda scheduler yozmaymiz, lekin channel semantikasini `sync.Mutex`, `sync.Cond` va wait queue bilan modellashtiramiz.

## Maqsad

Quyidagi API:

```go
type Chan[T any] struct {}

func Make[T any](cap int) *Chan[T]
func (c *Chan[T]) Send(v T)
func (c *Chan[T]) Recv() (T, bool)
func (c *Chan[T]) Close()
func (c *Chan[T]) Len() int
func (c *Chan[T]) Cap() int
```

## Buffered Channel

Ichki ring buffer:

```go
type Chan[T any] struct {
    mu     sync.Mutex
    sendq  *waitq[T]
    recvq  *waitq[T]
    buf    []T
    qcount int
    sendx  int
    recvx  int
    closed bool
}
```

## Semantika

Test bilan tekshiring:

- unbuffered send receiver kelguncha block qiladi
- unbuffered receive sender kelguncha block qiladi
- buffered send buffer to'lmaguncha block qilmaydi
- full buffered channelga send block qiladi
- empty buffered channeldan receive block qiladi
- close qilingan channeldan receive zero value, `ok=false`
- close qilingan channelga send panic
- channelni ikkinchi marta close qilish panic

## Wait Queue

Runtime ichida `sudog` kabi wait node ishlatiladi. Labda soddalashtiring:

```go
type waiter[T any] struct {
    value T
    ready chan struct{}
    ok    bool
    next  *waiter[T]
}
```

Birinchi versiyada har waiter uchun Go channel ishlatish mumkin. Keyingi versiyada `sync.Cond` yoki custom parking modelga o'ting.

## Select

`select` ni darhol yozmang. Avval channel to'g'ri bo'lsin.

Keyingi bosqich:

- non-blocking send
- non-blocking receive
- timeout
- select fairness

## Testlar

- goroutine leak yo'q
- race detector toza
- close semantics
- stress: 100 sender, 100 receiver
- benchmark: Go builtin channel bilan solishtirish
