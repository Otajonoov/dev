# 07. WAL And Recovery

WAL database crashdan keyin tiklanishi uchun kerak.

## LSN

```go
type LSN uint64
```

LSN WAL ichidagi byte offset yoki logical record id bo'lishi mumkin.

## WAL Record

```go
type RecordType uint8

const (
    RecordBegin RecordType = iota
    RecordCommit
    RecordAbort
    RecordHeapInsert
    RecordHeapDelete
    RecordHeapUpdate
    RecordPageImage
    RecordCheckpoint
)
```

Record header:

```go
type RecordHeader struct {
    LSN  LSN
    Type RecordType
    XID  TxID
    Len  uint32
    CRC  uint32
}
```

## WAL Manager

```go
func (w *Manager) Append(rec Record) (LSN, error)
func (w *Manager) Flush(lsn LSN) error
func (w *Manager) Replay(from LSN, apply func(Record) error) error
```

## Write-Ahead Rule

```text
data page flush qilinishidan oldin unga tegishli WAL record fsync qilingan bo'lishi kerak
```

## Recovery V1

Soddalashtirilgan:

```text
1. WAL boshidan o'qiladi
2. committed transaction recordlari topiladi
3. committed heap changes qayta apply qilinadi
4. uncommitted changes ignore qilinadi
```

Keyin ARIESga yaqinlashasiz:

- analysis
- redo
- undo
- checkpoint
- page LSN

## Testlar

- insert WAL record yoziladi
- commit flush qiladi
- crash simulation: data page flushsiz, WAL bor -> recovery rowni tiklaydi
- uncommitted insert recoverydan keyin ko'rinmaydi
- corrupted WAL CRC bilan rad qilinadi
