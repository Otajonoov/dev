# 06. Buffer Manager

Buffer manager disk page'larni memoryda cache qiladi.

## Frame

```go
type Frame struct {
    PageID PageID
    Data   storage.Page
    Dirty  bool
    Pin    int
    Usage  uint8
    LSN    uint64
}
```

## Manager

```go
type Manager struct {
    frames []Frame
    table  map[PageID]*Frame
}

func (m *Manager) FetchPage(id PageID) (*Frame, error)
func (m *Manager) NewPage() (*Frame, error)
func (m *Manager) UnpinPage(id PageID, dirty bool) error
func (m *Manager) FlushPage(id PageID) error
func (m *Manager) FlushAll() error
```

## Replacement

Avval LRU yozish mumkin. Postgresga yaqinroq model uchun CLOCK:

```text
usage_count > 0 -> decrement
usage_count == 0 and pin == 0 -> victim
```

## WAL Rule

Dirty page flushdan oldin:

```text
page.LSN <= wal.flushedLSN
```

Ya'ni WAL oldin diskka tushadi, keyin data page.

## Testlar

- fetch same page returns same frame
- pin bo'lgan page evict qilinmaydi
- dirty page flush bo'ladi
- CLOCK victim tanlaydi
- WAL rule buzilsa flush error
- buffer size kichik bo'lganda stress test
