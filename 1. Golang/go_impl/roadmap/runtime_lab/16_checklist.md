# 16. Checklist

## Scope

- [ ] Go runtime'ni almashtirmasligimiz aniq yozilgan
- [ ] off-heap memory modeli tushunilgan
- [ ] Go pointerlarni off-heapda saqlash xavfi tushunilgan

## Unsafe

- [ ] `unsafe.Pointer`
- [ ] `uintptr` xavfi
- [ ] `unsafe.Add`
- [ ] `unsafe.Slice`
- [ ] `unsafe.String`
- [ ] `runtime.KeepAlive`
- [ ] `checkptr`

## Memory Region

- [ ] byte-backed region
- [ ] mmap-backed region
- [ ] munmap
- [ ] page alignment
- [ ] guard page
- [ ] double close handling

## Allocator

- [ ] bump allocator
- [ ] free list allocator
- [ ] pool allocator
- [ ] size classes
- [ ] tiny allocator
- [ ] small allocator
- [ ] large allocator
- [ ] debug canary
- [ ] poison on free
- [ ] stats
- [ ] benchmarks

## Array, Slice, String

- [ ] custom array
- [ ] custom slice
- [ ] append/grow
- [ ] slicing
- [ ] copy
- [ ] immutable string
- [ ] Go string copy
- [ ] fuzz with std slice

## Map

- [ ] hasher/equal abstraction
- [ ] chained map
- [ ] bucket map
- [ ] overflow bucket
- [ ] resize
- [ ] evacuation
- [ ] tombstone
- [ ] Swiss 8-slot group
- [ ] fuzz with std map

## Channel

- [ ] buffered channel
- [ ] unbuffered channel
- [ ] close semantics
- [ ] send to closed panic
- [ ] receive from closed zero value
- [ ] wait queue
- [ ] race detector clean

## Type System

- [ ] struct layout calculator
- [ ] field offsets
- [ ] padding
- [ ] type descriptor
- [ ] eface model
- [ ] iface/itab model
- [ ] generic vector via type descriptor

## Stack, Defer, Panic

- [ ] frame stack
- [ ] stack allocator
- [ ] heap allocator integration
- [ ] defer LIFO
- [ ] panic unwind
- [ ] recover rule
- [ ] nested panic test

## Garbage Collector

- [ ] object header
- [ ] root set
- [ ] pointer mask
- [ ] mark-sweep
- [ ] cycle graph traversal
- [ ] tri-color model
- [ ] write barrier simulator
- [ ] incremental `Step`
- [ ] mark assist simulator

## GMP Scheduler

- [ ] G struct
- [ ] M struct
- [ ] P struct
- [ ] G statuses
- [ ] cooperative scheduler
- [ ] local run queue
- [ ] global run queue
- [ ] work stealing
- [ ] park/unpark
- [ ] syscall blocking simulator
- [ ] preemption budget
- [ ] timer/netpoll mini model
- [ ] channel integration

## Verification

- [ ] `go test ./...`
- [ ] `go test -race ./...`
- [ ] `go test -gcflags=all=-d=checkptr=2 ./...`
- [ ] fuzz tests
- [ ] benchmarks
- [ ] pprof
- [ ] trace
