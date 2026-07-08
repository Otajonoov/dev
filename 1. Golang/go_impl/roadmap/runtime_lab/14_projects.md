# 14. Projects

## Project 1: Memory Region + Bump Allocator

Maqsad:

- `mmap`
- align
- OOM
- reset
- stats

Natija:

```go
region, _ := memory.MmapRegion(64 << 20)
alloc := allocator.NewBump(region)
ptr := alloc.Alloc(128, 8)
```

## Project 2: Unified Allocator

Maqsad:

- tiny
- small
- large
- size classes
- debug mode

API:

```go
alloc := allocator.New(runtimeRegion)
p := alloc.Alloc(24, 8)
alloc.Free(p, 24)
```

## Project 3: Runtime Slice And String

Maqsad:

- custom array
- custom slice
- append/grow
- string copy
- allocator stats

## Project 4: Runtime Map

Maqsad:

- type-specific hash/equal
- old bucket model
- Swiss group model
- resize
- tombstone

## Project 5: Runtime Channel

Maqsad:

- buffered channel
- unbuffered channel
- close semantics
- blocking send/recv
- stress test

## Project 6: Type System Lab

Maqsad:

- struct layout calculator
- type descriptor
- eface/iface model
- generic vector via `TypeDesc`

## Project 7: Mini Stack And Panic Runtime

Maqsad:

- frame stack
- defer list
- panic unwind
- recover

## Project 8: Toy GC

Maqsad:

- object header
- root set
- pointer mask
- mark-sweep
- tri-color model
- write barrier simulator
- incremental `Step`

## Project 9: GMP Scheduler Simulator

Maqsad:

- G/M/P structs
- local/global run queue
- cooperative yield
- work stealing
- park/unpark
- syscall blocking simulator
- channel integration

## Final Project: Mini Runtime Lab

Hammasini bitta modulga ulang:

```text
memory.Region
allocator.Allocator
rt.Slice
rt.String
rt.Map
rt.Chan
rt.TypeDesc
rt.Stack
rt.Panic
rt.GC
rt.Scheduler
```

Final demo:

- custom allocator stats
- custom slice append
- custom map put/get
- custom channel send/recv
- panic/defer/recover simulator
- toy GC collect cycle
- GMP scheduler task trace
- benchmark report
