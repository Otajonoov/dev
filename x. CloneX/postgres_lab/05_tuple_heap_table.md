# 05. Tuple And Heap Table

Heap table - rowlar page ichida tartibsiz joylashgan relation file.

## TID

Postgresdagi CTIDga o'xshash:

```go
type TID struct {
    PageID PageID
    SlotID uint16
}
```

## Tuple Header

Minimal MVCC uchun:

```go
type TupleHeader struct {
    Xmin uint64
    Xmax uint64
    Flags uint16
    NullBitmapLen uint16
}
```

V0 uchun MVCCsiz ham boshlash mumkin:

```go
type TupleHeader struct {
    Flags uint16
}
```

## Row Encoding

Avval fixed + simple varlen:

```text
tuple header
column count
null bitmap
column values
```

Types:

- `INT`
- `BIGINT`
- `BOOL`
- `TEXT`

## Heap API

```go
type Heap struct {
    rel *storage.Relation
    buf *buffer.Manager
}

func (h *Heap) Insert(tx TxID, row Row) (TID, error)
func (h *Heap) Get(tid TID) (Row, error)
func (h *Heap) Delete(tx TxID, tid TID) error
func (h *Heap) SeqScan(snapshot Snapshot) Iterator
```

## Insert Algorithm

```text
for page in relation:
  if page has enough free space:
    insert tuple
    return tid
extend relation
insert tuple into new page
```

Keyin Free Space Map qo'shiladi.

## Testlar

- insert/get
- multiple pages
- seq scan
- delete marks tuple
- restartdan keyin rowlar saqlanadi
- random rows encode/decode
