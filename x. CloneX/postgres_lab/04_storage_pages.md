# 04. Storage Pages

Postgres-like DB diskni page'larga bo'lib ishlaydi. Birinchi asosiy component: 8KB page.

## Page Size

```go
const PageSize = 8192

type PageID uint32
type Page []byte
```

## Page Header

Minimal header:

```go
type PageHeader struct {
    LSN       uint64
    Lower     uint16
    Upper     uint16
    Special   uint16
    PageFlags uint16
}
```

Slotted page layout:

```text
[ header ][ slot array ->          free space          <- tuple data ][ special ]
```

## Slot

```go
type Slot struct {
    Offset uint16
    Length uint16
    Flags  uint16
}
```

## API

```go
func NewPage() Page
func (p Page) Header() PageHeader
func (p Page) FreeSpace() int
func (p Page) InsertCell(data []byte) (slot uint16, error)
func (p Page) GetCell(slot uint16) ([]byte, error)
func (p Page) DeleteCell(slot uint16) error
```

## Relation File

```go
type Relation struct {
    file *os.File
}

func (r *Relation) ReadPage(id PageID) (Page, error)
func (r *Relation) WritePage(id PageID, p Page) error
func (r *Relation) NumPages() (uint32, error)
func (r *Relation) ExtendPage() (PageID, Page, error)
```

## Testlar

- empty page free space to'g'ri
- cell insert/get
- page full bo'lganda error
- delete slotni dead qiladi
- page diskka yozilib qayta o'qiladi
- 1000 tuple bir nechta pagega tarqaladi

## Muhim Qaror

Birinchi versiyada checksums yo'q. Keyin page checksum va torn-page detection qo'shiladi.
