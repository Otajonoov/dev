# 17. Resources

## Go Source

- `runtime/malloc.go` - allocator, tiny/small/large, mcache/mcentral/mheap
- `runtime/slice.go` - growslice, nextslicecap
- `runtime/string.go` - string allocation and conversion helpers
- `internal/runtime/maps/map.go` - Swiss Table based map
- `runtime/chan.go` - channel implementation
- `runtime/panic.go` - defer, panic, recover
- `runtime/mgc.go` - GC orchestration
- `runtime/mgcmark.go` - marking
- `runtime/mgcsweep.go` - sweeping
- `runtime/mbarrier.go` - write barriers
- `runtime/proc.go` - scheduler and GMP
- `runtime/runtime2.go` - core runtime structs
- `runtime/time.go` - timers
- `runtime/netpoll.go` - network poller
- `cmd/compile/internal/escape` - escape analysis

## Go Docs

- Go specification
- Go memory model
- `unsafe` package docs
- `sync/atomic` package docs
- `runtime` package docs

## Books

- The Go Programming Language
- 100 Go Mistakes and How to Avoid Them
- Concurrency in Go
- Operating Systems: Three Easy Pieces
- Computer Systems: A Programmer's Perspective
- Database Internals

## Tools

- `go test`
- `go test -race`
- `go test -fuzz`
- `go test -bench`
- `benchstat`
- `pprof`
- `go tool trace`
- `go build -gcflags="-m=2"`
- `go tool objdump`

## Source Reading Order

1. `runtime/slice.go`
2. `runtime/string.go`
3. `runtime/malloc.go` top comments
4. `runtime/chan.go` top invariants
5. `internal/runtime/maps/map.go` top comments
6. `runtime/panic.go`
7. `runtime/mgc.go`
8. `runtime/mgcmark.go`
9. `runtime/mgcsweep.go`
10. `runtime/mbarrier.go`
11. `runtime/proc.go`
12. `cmd/compile/internal/escape`
