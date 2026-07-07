# 08. Transactions And MVCC

Postgresning yuragi: MVCC.

## Transaction ID

```go
type TxID uint64

type TxStatus uint8

const (
    TxInProgress TxStatus = iota
    TxCommitted
    TxAborted
)
```

## Tuple MVCC Header

```go
type TupleHeader struct {
    Xmin TxID
    Xmax TxID
    Flags uint16
}
```

Meaning:

- `xmin`: tuple yaratgan transaction
- `xmax`: tuple delete/update qilgan transaction

## Snapshot

```go
type Snapshot struct {
    Xmin TxID
    Xmax TxID
    Active map[TxID]struct{}
}
```

## Visibility

Read Committed uchun:

```text
tuple visible if:
  xmin committed
  xmin not active
  xmax is zero OR xmax aborted OR xmax active
```

Repeatable Read uchun transaction boshida bitta snapshot olinadi.

## Transaction Manager

```go
func (m *Manager) Begin() *Tx
func (m *Manager) Commit(tx *Tx) error
func (m *Manager) Abort(tx *Tx) error
func (m *Manager) Snapshot() Snapshot
func (m *Manager) Status(xid TxID) TxStatus
```

## Update Model

Postgres update in-place emas:

```text
old tuple xmax = current xid
new tuple xmin = current xid
```

Keyin HOT update optional.

## Testlar

- uncommitted insert boshqa txga ko'rinmaydi
- committed insert ko'rinadi
- delete qilingan tuple yangi txga ko'rinmaydi
- aborted insert ko'rinmaydi
- repeatable read bir xil snapshotni saqlaydi
- update eski va yangi version yaratadi
