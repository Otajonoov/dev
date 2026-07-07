# 10. Defer, Panic, Recover

Bu bosqichda haqiqiy Go `panic/recover`ni qayta yozmaymiz. O'z frame stack modelimiz ustida defer va panic unwind yozamiz.

## Defer Model

```go
type Defer struct {
    fn   func(*Context)
    args unsafe.Pointer
    next *Defer
}

type Frame struct {
    defers *Defer
}
```

Qoida:

- defer LIFO tartibda ishlaydi
- normal return paytida ishlaydi
- panic unwind paytida ishlaydi

## Panic Model

```go
type Panic struct {
    value     any
    recovered bool
}

type Context struct {
    stack  *Stack
    panic  *Panic
}
```

Algorithm:

```text
panic(value)
  current panic = value
  frame by frame unwind
  har frame defers LIFO chaqiriladi
  defer ichida recover bo'lsa panic recovered
  recovered bo'lsa normal control flowga qaytadi
  recovered bo'lmasa top-level fatal
```

## Recover Qoidasi

Go semantikasiga yaqin:

- `recover` faqat deferred function ichida ishlaydi
- normal function ichida `recover()` nil
- panic bo'lmasa `recover()` nil
- bir panic faqat bir marta recover qilinadi

Labda buni context flag bilan modellashtiring:

```go
type Context struct {
    inDeferredCall bool
    panic *Panic
}
```

## Open-Coded Defer

Keyingi bosqich:

- har defer heap/listga ketmasin
- frame ichida bitmask bilan open-coded defer modelini yozing

Bu Go compiler/runtime optimizationini tushunish uchun foydali.

## Testlar

- normal return defer order
- panic defer order
- recover deferred function ichida ishlaydi
- nested panic
- panic while panicking
- recoverdan keyin execution model dokumentatsiya qilinadi
