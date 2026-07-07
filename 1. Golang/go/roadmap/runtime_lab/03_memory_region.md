# 03. Memory Region

Allocatorning birinchi qatlami OS dan raw memory olishdir.

## Maqsad

Quyidagi API ni yozing:

```go
type Region struct {
    base unsafe.Pointer
    size uintptr
    data []byte
}

func MmapRegion(size uintptr) (*Region, error)
func (r *Region) Close() error
func (r *Region) Ptr(offset uintptr) unsafe.Pointer
func (r *Region) Bytes() []byte
```

`data []byte` ni struct ichida saqlash muhim. `syscall.Mmap` yoki `unix.Mmap` qaytargan slice `Munmap` uchun kerak bo'ladi.

## Backendlar

Bosqichma-bosqich:

1. `[]byte` backed region: oson, test uchun.
2. `unix.Mmap` backed region: Linux/macOS uchun.
3. `mprotect` guard page: debug uchun.

## Nima Bilan Yozamiz

Tanlov:

```text
avval: make([]byte, size)
keyin: golang.org/x/sys/unix.Mmap
emas:  cgo / C.malloc
```

Sabab:

- `[]byte` backend unit test va fuzz test uchun qulay.
- `unix.Mmap` haqiqiy off-heap memory beradi.
- `cgo` build, cross-compile va pointer qoidalarini murakkablashtiradi.
- Allocator algoritmlari baribir Pure Go + `unsafe` bilan yoziladi.

## Alignment

Har helper power-of-two alignment talab qilsin:

```go
func AlignUp(n, align uintptr) uintptr {
    if align == 0 || align&(align-1) != 0 {
        panic("alignment must be power of two")
    }
    return (n + align - 1) &^ (align - 1)
}
```

## Guard Page

Debug mode:

```text
[ guard ][ usable memory ][ guard ]
```

Overrun bo'lsa process tez crash qiladi. Bu yaxshi: xato jim qolmaydi.

## Testlar

- zero size rad qilinadi
- size page boundaryga align qilinadi
- `Ptr(0)` basega teng
- `Ptr(size)` panic yoki nil
- `Close` dan keyin region ishlatilmaydi
- double close panic emas, error qaytaradi

## Eslatma

Windows uchun alohida backend keyin yoziladi. Birinchi versiyani Linux bilan cheklash normal.
