# 13. Testing, Benchmark, Debug

Unsafe loyiha testsiz yurmaydi. Har bosqich test va benchmark bilan yopiladi.

## Unit Tests

Har package:

```bash
go test ./...
```

Allocator uchun minimal testlar:

- alignment
- overlap yo'q
- OOM
- free reuse
- double free
- use-after-free debug detection
- stats

## Fuzz Tests

Standard Go type bilan solishtiring:

```bash
go test -fuzz=FuzzMap ./...
```

Misollar:

- custom slice vs `[]int`
- custom map vs `map[string]int`
- custom channel semantics uchun stress test

## Race Detector

Concurrent package uchun:

```bash
go test -race ./...
```

`mychan`, concurrent allocator cache, stats counterlarda shart.

## Checkptr

Unsafe code uchun:

```bash
go test -gcflags=all=-d=checkptr=2 ./...
```

## Benchmarks

```bash
go test -bench=. -benchmem -count=10 ./... > bench.txt
benchstat bench.txt
```

Benchmark yo'nalishlari:

- tiny alloc
- small alloc
- large alloc
- slice append
- map put/get/delete
- channel send/recv

## Profiling

```bash
go test -bench=. -cpuprofile=cpu.prof ./...
go test -bench=. -memprofile=mem.prof ./...
go test -bench=. -trace=trace.out ./...

go tool pprof -http=:8080 cpu.prof
go tool trace trace.out
```

## Debug Allocator

Debug mode:

- magic number
- canary
- red zone
- poison on free
- allocation site
- stack trace optional

Poison values:

```text
0xAA allocated
0xDD freed
0xCC red zone
```

## CI Gate

Har milestone oxirida:

```bash
go test ./...
go test -race ./...
go test -gcflags=all=-d=checkptr=2 ./...
go test -bench=. -benchmem ./...
```
