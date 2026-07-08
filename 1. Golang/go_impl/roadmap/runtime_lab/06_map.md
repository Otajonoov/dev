# 06. Map

Mapni ikki bosqichda yozing: avval oddiy model, keyin Go runtimega yaqin model.

## Bosqich 1: Chained Hash Map

Maqsad: hash, collision, resize tushunish.

```go
type Entry[K comparable, V any] struct {
    key   K
    value V
    next  *Entry[K, V]
}
```

Bu bosqich oson, lekin Go runtime mapiga o'xshamaydi.

## Bosqich 2: Old Bucket Model

Go'ning eski map modeliga yaqin:

```go
const bucketSize = 8

type bucket[K comparable, V any] struct {
    tophash [bucketSize]uint8
    keys    [bucketSize]K
    values  [bucketSize]V
    overflow *bucket[K, V]
}
```

O'rganiladigan mavzular:

- top hash
- overflow bucket
- load factor
- grow
- evacuation
- iteration randomization

## Bosqich 3: Swiss Table Model

Go 1.24+ map implementatsiyasi Swiss Table dizayniga asoslangan. Muhim tuzatma: Go source ichida group 8 slotli (`abi.MapGroupSlots`) model sifatida keladi. 16 slot SIMD g'oyasi umumiy Swiss Table optimizatsiyasi bo'lishi mumkin, lekin labda 8 slotdan boshlang.

Minimal group:

```go
const groupSlots = 8

type group[K comparable, V any] struct {
    ctrl  [groupSlots]uint8
    keys  [groupSlots]K
    vals  [groupSlots]V
}
```

Control byte:

```text
empty
deleted
used + H2 lower 7 bits
```

Hash qismlari:

```text
H1: upper bits, group/probe uchun
H2: lower 7 bits, ctrl byte uchun
```

## Type-Specific Hash

`K` ni `unsafe` bilan bytega aylantirib hash qilish umumiy holatda noto'g'ri.

Shuning uchun map constructor hasher va equal qabul qilsin:

```go
type Hasher[K any] func(K, uint64) uint64
type Equal[K any] func(a, b K) bool

type Map[K any, V any] struct {
    hash Hasher[K]
    eq   Equal[K]
    seed uint64
}
```

Boshlanish uchun:

- `string`
- `int`
- `uint64`

Keyin struct keylar uchun type descriptor qo'shiladi.

## Nil Map Semantikasi

Go semantikasini test qiling:

- nil mapdan read zero value qaytaradi
- nil mapga write panic
- delete nil mapda no-op
- len nil mapda 0

## Testlar

- put/get/delete
- overwrite
- collision
- resize
- tombstone
- random operation: standard map bilan solishtirish
- iteration paytida delete/add semantikasi alohida hujjatlansin
