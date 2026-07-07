# 12. GMP Scheduler

Bu bosqichda haqiqiy Go scheduler'ni almashtirmaymiz. GMP modelini simulator sifatida yozamiz.

## Chegara

Go runtime ichidagi real G, M, P sizga to'g'ridan-to'g'ri boshqarilmaydi. Siz yozadigan model o'z tasklarini boshqaradi:

```text
real goroutine -> Go scheduler
runtime-lab task -> siz yozgan GMP scheduler simulator
```

## Model

```go
type G struct {
    id     int64
    status GStatus
    fn     func(*Context)
    stack  *Stack
}

type M struct {
    id int64
    p  *P
}

type P struct {
    id       int64
    runq     RingQueue[*G]
    runqNext *G
}

type Scheduler struct {
    globalRunq Queue[*G]
    ps         []*P
    ms         []*M
}
```

## G Status

```go
type GStatus uint8

const (
    GIdle GStatus = iota
    GRunnable
    GRunning
    GWaiting
    GSyscall
    GDead
)
```

## Bosqich 1: Cooperative Scheduler

Avval tasklar o'zi yield qiladi:

```go
func (s *Scheduler) Go(fn func(*Context))
func (ctx *Context) Yield()
func (s *Scheduler) Run()
```

Bu bosqichda preemption yo'q.

## Bosqich 2: Local Va Global Run Queue

Qoida:

- yangi G avval local P queuega tushadi
- local queue to'lsa global queuega o'tadi
- M avval o'z P local queueidan oladi
- local bo'sh bo'lsa global queuega qaraydi

## Bosqich 3: Work Stealing

P bo'sh qolsa boshqa P queueidan yarmini o'g'irlaydi:

```text
P0 empty -> P1 runq half steal -> P0 run
```

Test:

- tasklar P'lar bo'yicha tarqaladi
- starvation bo'lmaydi
- global queue cheksiz o'smaydi

## Bosqich 4: Parking Va Unparking

Channel, mutex, timer kabi holatlarda G waiting bo'ladi:

```go
func (ctx *Context) Park(reason string)
func (s *Scheduler) Ready(g *G)
```

Bu `channel` implementatsiyasi bilan bog'lanadi:

```text
Recv empty channel -> Park
Send kelganda -> Ready(receiver)
```

## Bosqich 5: Syscall Blocking Simulator

Syscallga kirgan G M'ni band qilib qo'yadi. Real Go runtime P'ni ajratib olib boshqa M'ga beradi.

Lab modeli:

```text
G enters syscall
M blocks
P detached
another M picks P
syscall returns
G runnable
```

## Bosqich 6: Preemption Simulator

Real async preemption emas, lekin instruction budget bilan model qiling:

```go
func (ctx *Context) Step()
```

Budget tugasa:

```text
GRunning -> GRunnable
queuega qaytadi
```

## Netpoll/Timer Mini Model

Advanced:

- timer heap
- sleep
- network event simulator
- parked G timer tugaganda ready bo'ladi

## Testlar

- cooperative yield order
- local/global queue
- work stealing
- park/unpark
- syscall P handoff
- preemption budget
- channel bilan integration

## Manba O'qish

Go source:

- `runtime/proc.go`
- `runtime/runtime2.go`
- `runtime/time.go`
- `runtime/netpoll.go`
- `runtime/chan.go`
