# 01. Prerequisites

## Go Bilimlari

Bilish shart:

- pointer: `*T`, `&v`
- slice: len, cap, backing array
- map zero value, nil map, iteration
- channel: buffered, unbuffered, close
- interface va type assertion
- generics: `T any`, `comparable`
- defer, panic, recover
- goroutine va scheduler haqida umumiy tasavvur
- tests, benchmarks, race detector

## Memory Bilimlari

Kerakli tushunchalar:

- virtual memory
- page va page alignment
- stack vs heap
- alignment va padding
- endianness
- cache line
- false sharing
- fragmentation
- use-after-free
- double-free
- dangling pointer

## OS Bilimlari

Boshlang'ich darajada bilish kerak:

- `mmap`, `munmap`, `mprotect`
- anonymous mapping
- file mapping
- page size
- syscall nima
- Linux memory overcommit

## Compiler Bilimlari

Kerakli komandalar:

```bash
go build -gcflags="-m=2" ./...
go test -race ./...
go test -gcflags=all=-d=checkptr=2 ./...
go test -bench=. -benchmem ./...
go tool pprof cpu.prof
go tool trace trace.out
```

## Boshlashdan Oldin Sinov

Quyidagilarni tushuntira olsangiz tayyorsiz:

- Nega `uintptr` pointer emas?
- Nega off-heap ichidagi Go pointer xavfli?
- `len` va `cap` farqi nima?
- Nil mapga yozish nega panic?
- Closed channeldan receive nima qaytaradi?
- `recover` qachon ishlaydi, qachon ishlamaydi?
