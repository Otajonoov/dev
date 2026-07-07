# 01. Prerequisites

## Go

Bilish kerak:

- `io.Reader`, `io.Writer`
- `os.File`, `ReadAt`, `WriteAt`, `Sync`
- `encoding/binary`
- `context`
- `sync.Mutex`, `sync.RWMutex`, `sync.Cond`
- `testing`, fuzz, benchmark
- errors wrapping
- binary protocol parsing

## OS

Kerakli tushunchalar:

- file descriptor
- page cache
- fsync
- partial write
- torn page
- crash consistency
- directory fsync
- endian

## Database Internals

Tushunish kerak:

- page
- slotted page
- tuple
- heap table
- buffer pool
- WAL
- LSN
- transaction ID
- MVCC
- index
- query plan

## Minimal Mashqlar

Boshlashdan oldin:

1. 8KB byte slice ichiga binary header yozing.
2. `os.File.WriteAt` bilan page 0 va page 1 yozing.
3. Programni qayta ishga tushirib page headerni o'qing.
4. 1000 row insert qilib file size va page countni tekshiring.
