# 14. Locks And Isolation

MVCC hamma narsani hal qilmaydi. Schema, pages, index, write conflicts uchun lock kerak.

## Lock Manager

```go
type LockMode uint8

const (
    AccessShare LockMode = iota
    RowShare
    RowExclusive
    Exclusive
)

type LockTag struct {
    Kind string
    ID   uint64
}
```

API:

```go
func (lm *Manager) Lock(tx TxID, tag LockTag, mode LockMode) error
func (lm *Manager) Unlock(tx TxID, tag LockTag) error
func (lm *Manager) UnlockAll(tx TxID)
```

## Row Locks

Update/delete conflict:

```text
two tx update same tuple -> one waits or gets serialization error
```

## Deadlock Detection

Wait-for graph:

```text
T1 waits T2
T2 waits T1
```

Cycle bo'lsa victim abort.

## Isolation

Bosqichlar:

1. Read Committed
2. Repeatable Read
3. Serializable simulator

## Testlar

- incompatible locks wait
- compatible locks pass
- update conflict
- deadlock detection
- read committed sees new committed rows per statement
- repeatable read snapshot stable
