
Предположим, запись `_panic` уже recovered. 
Как это **фактически меняет** control flow? 
Ключ — вызов `mcall(recovery)`:

> **Примечание.** `mcall()` — специальная assembly-функция, переключающая с stack текущей goroutine на **system stack**. Это переключение важно, потому что позволяет Go безопасно изменить состояние goroutine — слишком рискованно делать это, оставаясь на том же goroutine stack.

```go
func recovery(gp *g) {
    ...
    // Заставляем deferproc для этого d вернуться снова,
    // на этот раз возвращая 1. Вызывающая функция
    // перейдёт к стандартному return epilogue.
    gp.sched.sp = sp
    gp.sched.pc = pc
    gp.sched.lr = 0

    switch {
    case goarch.IsAmd64 != 0:
        gp.sched.bp = fp - 2*goarch.PtrSize
    case goarch.IsArm64 != 0:
        gp.sched.bp = sp - goarch.PtrSize
    }

    gp.sched.ret = 1 // <---
    gogo(&gp.sched)
}
// recovery (src/runtime/panic.go)_
```

Go runtime берёт на себя **управление control flow**. 

Нормально `deferproc` и `deferprocStack` возвращают 0, сигнализируя обычный control flow. При recovery они **возвращают 1** и переходят к функции `deferreturn()` для выполнения оставшихся deferred-функций в этом frame.

Recovery выполняет следующее:

- **Устанавливает** stack pointer, program counter и frame pointer для caller, восстанавливая состояние от момента первого вызова `deferprocStack`.
- **Устанавливает return-значение в 1** (`gp.sched.ret = 1`) — как если бы `deferproc`/`deferprocStack` возвращал 1.
- **Вызывает `gogo(&gp.sched)`** для фактического возобновления goroutine в модифицированной точке.

Значение в `gp.sched.ret` **не находится** на goroutine stack; оно копируется непосредственно в **return-value регистр** архитектуры (AX на amd64, R0 на arm64 и т.д.) при выполнении `gogo`.

---
### Return-значение при восстановлении из panic

Что происходит, если мы установили return-значение **до** panic?

```go
func doSomething() (res uint) {
    defer func() {
        fmt.Println("This remaining deferred function will execute")
    }()

    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered:", r)
        }
    }()

    res = 100
    panic("oops")
    return 42
}

func main() {
    result := doSomething()
    fmt.Println("After doSomething, result is", result)
}
// Output:
// Recovered: oops
// This remaining deferred function will execute
// After doSomething 100
```

Поскольку panic восстановлен, функция **возвращается нормально**. 
Значение `res` всё ещё 100 благодаря предыдущему присваиванию — это и возвращается.

---

## Recovery и Argument Pointer

Рассмотрим случай, когда deferred-функция вызывает **другую функцию**, и та helper-функция вызывает `recover()`. Должно ли это работать?

```go
func recoverHelper() {
    if r := recover(); r != nil {
        fmt.Println("Recovered:", r)
    }
}

func doSomething() (res int, err error) {
    defer func() {
        recoverHelper()
    }()
    panic("oops")
}
```

Чтобы понять, что допустимо, рассмотрим, как `recover()` работает — точнее, как работает `gorecover()`:

```go
func gorecover(argp uintptr) any {
    gp := getg()
    p := gp._panic
    if p != nil && !p.goexit && !p.recovered && argp == uintptr(p.argp) {
        p.recovered = true
        return p.arg
    }
    return nil
}
// gorecover (src/runtime/panic.go)_
```

> **Примечание: Флаг p.goexit**
> 
> Существует специальная runtime-функция `runtime.Goexit()`. Она завершает текущую goroutine, подобно panic, но **без вызова** panic и без влияния на другие goroutine. При вызове все deferred-функции выполняются, но программа **не падает**. Можно думать о ней как о «return statement» для всей goroutine.

---
### Проверка argument pointer

Самая важная часть — **проверка argument pointer** (`argp == uintptr(p.argp)`). Эта проверка гарантирует, что только **правильная** deferred-функция может восстановить panic.

Когда происходит panic, Go runtime (через `runtime.gopanic()`) сохраняет позицию **outgoing argument area** в `_panic.argp`

![](../../assets/obsidian-images/Pasted%20image%2020260219173458.png)

Здесь `runtime.gopanic()` подготавливает аргументы для deferred-функции.

В этот момент runtime **не знает**, какая deferred-функция (если вообще) вызовет `recover()`
Ключ в том, что только **непосредственный callee** имеет правильный argument pointer, совпадающий с подготовленным `gopanic()`.

Когда deferred-функция вызывает `recover()`, она передаёт свой argument pointer в `runtime.gorecover(argp)`. Если эта deferred-функция — та, что вызвана **непосредственно** frame `gopanic()`, её argument pointer **совпадает** с `_panic.argp`, проверка проходит, и panic восстанавливается:

---
### Типичная ошибка с recover

```go
func customRecover() error {
    if r := recover(); r != nil {
        recoverTime.Increase() // metric
        log.Printf("Recovered from panic: %v", r) // log
        return fmt.Errorf("panic: %v", r)
    }
    return nil
}

func doSomething() (err error) {
    // func 1
    defer func() {
        err = errors.Join(err, customRecover())
    }()
    panic("oops")
}
```

Здесь разработчик пытается объединить ошибку с ошибкой от panic. 
Идея — собирать метрики и логировать в пользовательской функции. 

Проблема: argument pointer записи `_panic` (`_panic.argp`) **не совпадает** с argument pointer при вызове `recover()` в `customRecover`. Из-за этого проверка `argp == uintptr(p.argp)` **не проходит**, и panic не может быть восстановлен.

```go
func main() {
    defer func(i int) {
         recoverHelper()
    }(1)
    panic("oops")
}
```

---
### Другие ошибки

Вызов `recover()` **напрямую** в defer statement (не внутри deferred-функции) тоже не работает:

```go
func main() {
    defer recover()
    panic("Unexpected wrong")
}
```

---
### Два panic и один recover

Рассмотрим ситуацию с **двумя panic** и одним recover. Ловит ли один recover оба panic, только первый или только последний?

```go
func main() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered:", r)
        }
    }()

    defer func() {
        panic("Second panic")
    }()

    panic("First panic")
}
```

Ответ: один `recover()` может поймать **только один panic** за раз. В данном случае **второй panic** будет пойман, и программа вернётся нормально без аварийного завершения. Panic **не накапливаются**; `recover()` получает только самый **последний** panic:

```
Recovered: Second panic
```

---

## Заключительные замечания

Текущий механизм panic-recover не идеален и использует **множество трюков** для надёжной работы. Восстановление из panic означает, что runtime должен **очень осторожно** раскрутить stack, отслеживать, какие frame были созданы реальными вызовами функций, а какие добавлены runtime, и координироваться с deferred-функциями, которые могут сами вызвать дополнительные panic.

Многое зависит от конвенций: как хранятся return program counter, как разные архитектуры управляют link register и где компилятор вставляет `deferreturn`. Поскольку эти конвенции **различаются** между архитектурами и меняются с новыми оптимизациями компилятора, Go runtime использует **эвристики** и таблицы special-case для сохранения корректной работы.

---

## Xulosa

**1. panic — bu error handling emas, bu «dastur buzildi» signali.**
Go'da error handling `error` interface orqali amalga oshiriladi: `return err`
`panic` faqat **haqiqatan kutilmagan** holatlar uchun: nil pointer dereference, index out of range, yoki «bu hech qachon bo'lmasligi kerak» degan invariant buzilganda. 
Agar siz `panic` ni odatiy error handling o'rniga ishlatsangiz — bu Go idiom'lariga **zid**.

**2. recover faqat to'g'ridan-to'g'ri deferred-funksiya ichida ishlaydi.** Bu eng muhim qoida va eng ko'p xato manbayi. `recover()` muvaffaqiyatli bo'lishi uchun **argument pointer** tekshiruvi o'tishi kerak: `argp == uintptr(p.argp)`. Bu shuni anglatadi:

- `defer func() { recover() }()` — ✅ **ishlaydi**
- `defer func() { helperThatCallsRecover() }()` — ❌ **ishlamaydi** (argp mos kelmaydi)
- `defer recover()` — ❌ **ishlamaydi**

Bu cheklov ataylab qilingan — faqat **«birinchi darajali»** deferred-funksiya panic'ni to'xtata oladi.

**3. `gopanic` → `nextDefer` → `gorecover` → `recovery` → `gogo` — to'liq cycle.** 
Panic jarayoni quyidagicha ishlaydi:

1. `gopanic(e)` — `_panic` record yaratadi, `p.arg = e`, `p.start()` bilan stack unwinding boshlaydi
2. `nextDefer()` — defer chain'ni iterate qiladi, har bir deferred-funksiyani execute qiladi
3. Agar deferred-funksiya `recover()` chaqirsa → `gorecover()` `p.recovered = true` qiladi
4. `nextDefer()` qayta tekshirganda `p.recovered == true` → `mcall(recovery)` chaqiradi
5. `recovery()` — `gp.sched.ret = 1` qiladi, `gogo(&gp.sched)` bilan goroutine'ni **qayta boshlaydi**
6. Execution `deferprocStack` dan «qaytganday» ko'rinadi, lekin ret=1 — shuning uchun `deferreturn` ga o'tadi
7. `deferreturn` qolgan defer'larni execute qiladi, funksiya **normal** qaytadi

**4. `gp.sched.ret = 1` — deferproc va deferprocStack bilan bog'langan trick.** Bu oldingi defer bobidagi `return0()` assembly trick'ning **davomi**. Normal holatda `deferproc` 0 qaytaradi → execution davom etadi. Recovery holatda runtime `gp.sched.ret` ni 1 ga o'zgartiradi → `gogo` orqali goroutine qayta boshlanadi → compiler `CMP $0, R0` / `BNE` bilan buni aniqlaydi → `deferreturn` ga o'tadi. **Ikki bob birgalikda ishlaydi** — defer va panic mexanizmlari bir-biriga bog'langan.

**5. Named return value + recover = idiomatic error conversion.**

```go
func safe() (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic: %v", r)
        }
    }()
    riskyOperation()
    return nil
}
```

Bu **eng to'g'ri** pattern. Named return `err` closure orqali defer'da accessible, panic recover bo'lgandan keyin `err` set bo'ladi, funksiya normal qaytadi. Agar named return **ishlatmasangiz** — funksiya zero value qaytaradi (`0`, `""`, `nil`), bu **silent bug** manbayi.

**6. panic(nil) — Go 1.21 dan oldin xavfli edi.**

```go
defer func() {
    if r := recover(); r != nil { ... } // r == nil — panic bo'lmadi yoki nil panic?
}()
panic(nil)
```

Go 1.21 dan boshlab `panic(nil)` avtomatik `PanicNilError` ga wrap qilinadi — endi **farqlash mumkin**. Bu backward-incompatible o'zgarish, lekin xavfsizlik uchun to'g'ri qaror.

**7. Bir recover() faqat bitta panic ni ushlaydi.** Agar deferred-funksiya o'zi panic qilsa — **yangi** panic eskisini «yopadi». `recover()` faqat **eng oxirgi** panic ni qaytaradi. Panic'lar stack'lanmaydi:

```go
defer func() {
    recover() // "Second panic" ni ushlaydi, "First panic" yo'qoladi
}()
defer func() { panic("Second panic") }()
panic("First panic")
```

**8. Go 1.22 dagi recover bug — inlining sabab.** Agar deferred-funksiya argument'li bo'lsa va ichida `recoverHelper()` chaqirsa, compiler anonym funksiyani **inline** qilishi mumkin — `argp` tekshiruvi noto'g'ri o'tadi va recover **ishlamasligi kerak bo'lgan joyda ishlaydi**. Bu Go 1.26 da fix bo'lishi kutilmoqda.

**Amaliy maslahatlar:**

- `panic` ni **faqat** haqiqatan exceptional holatlar uchun ishlating — library API'da error return qiling, panic qilmang
- `recover()` ni doim **to'g'ridan-to'g'ri** deferred anonymous function ichida chaqiring — helper funksiyaga o'tkazmang
- Named return value ishlating agar recover'dan keyin **ma'noli** qiymat qaytarish kerak bo'lsa
- `defer recover()` deb yozmang — bu **ishlamaydi**, chunki recover deferred-funksiya **ichida** bo'lishi kerak
- Production kodda `recover` ishlatganingizda: logging + metric qo'shing, lekin **recover ichida** qiling, helper'da emas
- Panic value `any` turi — `string`, `error`, `int`, `struct` bo'lishi mumkin, shuning uchun type assertion kerak bo'lishi mumkin