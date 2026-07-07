
`panic` и `recover` — две встроенные функции в Go для обработки **исключительных ситуаций**, которые не могут быть обработаны через обычную обработку ошибок.

Используйте `panic`, когда что-то идёт **неожиданно неправильно** — проблемы, которые обычно не обрабатываются обычными проверками ошибок: деление на ноль, разыменование nil pointer, доступ к индексу вне диапазона или баг, который, по вашему мнению, **никогда** не должен произойти:

Когда вы вызываете `panic()`, нормальное выполнение **останавливается**, и Go начинает «раскрутку» (unwinding) stack. Это означает, что он возвращается по call stack и выполняет все deferred-функции **в той же goroutine**. 
Мы видели это с `deferreturn` — создаётся запись `_panic` и обходится defer chain.

Panicking продолжает выполнять deferred-функции в LIFO-порядке, пока не найдёт deferred-функцию, вызывающую `recover()`:

```go
func A() {
    println("A called")
}

func main() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered:", r)
        }
    }()

    defer A()

    fmt.Println("Starting the program")
    panic("A severe error occurred")
    fmt.Println("This will not be printed")
}
// Output:
// Starting the program
// A called
// Recovered: A severe error occurred
```

В этом примере deferred-функция настроена для вызова `recover`. 
После вывода «Starting the program» вызывается panic. 
Go немедленно **останавливает** обычное выполнение и начинает раскрутку stack, выполняя все defer. Deferred-функция с `recover` вызывается и получает значение panic. 
Поскольку значение **не nil**, программа выводит «Recovered: A severe error occurred». Последний `fmt.Println` **пропускается**.

**Ключевой момент:** после восстановления из panic функция завершается выполнением оставшихся deferred-функций, затем **завершается нормально**. 
Оригинальный flow не продолжается от точки panic.

---
### Return-значение после recover

Что насчёт значения, которое функция возвращает после восстановления из panic?

```go
func add(a, b int) int {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered:", r)
        }
    }()

    panic("oops")
    return a + b
}
```

Логически, когда panic восстановлен, функция не может продолжить выполнение. Функция останавливается и возвращает **значение по умолчанию** своего типа. 
`a + b` никогда не достигается, поэтому функция возвращает `0` (default для `int`).

Это может быть рискованно. Чтобы использовать panic как ошибку для возврата вызывающему, можно использовать **named return values**:

```go
func add(a, b int) (res int, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic: %v", r)
        }
    }()

    res = a + b
    panic("oops")
    return
}
```

Теперь panic **и восстановлен, и конвертирован** в ошибку, позволяя caller обработать failure вместо тихого отказа программы.

---
### Значение panic может быть чем угодно

Значение, передаваемое в `panic(...)` и получаемое из `recover()`, может быть **чем угодно**. Аргумент — `interface{}` (empty interface), поэтому вы не ограничены ошибками. Можно использовать строки, числа, struct или что-либо ещё.

Внутренне значение вроде `panic("A severe error occurred")` выглядит как `eface{ type: string, value: &"A severe error occurred" }`.

---
### panic(nil) — поведение по версиям

Что произойдёт при `panic(nil)`?

```go
defer func() {
    if r := recover(); r != nil {
        fmt.Println("Recovered:", r)
    }
}()
panic(nil)
```

Ответ зависит от версии Go:

**До Go 1.21:** при вызове `panic(nil)` и recover значение, возвращаемое `recover()`, тоже было `nil`. Это создавало путаницу — **невозможно** было отличить:

- Panic не произошёл (нормальный return из recover)
- Panic вызван с nil значением

**Начиная с Go 1.21:** если вызвать `panic(nil)`, runtime автоматически оборачивает nil в новый тип `PanicNilError` с сообщением: «panic called with nil argument».

---

## Что Runtime делает при Panic

Рассмотрим пример:

```go
func main() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered:", r)
        }
    }()

    defer func() {
        fmt.Println("defer but not recovered")
    }()

    panic("Unexpected wrong")
}
```

Вызов `panic(any)` транслируется в функцию `runtime.gopanic(any)`.
Эта функция настраивает контекст panic и запускает процесс раскрутки stack:

```go
// Реализация предопределённой функции panic.
func gopanic(e any) {
    ... // if e == nil {

    // Настройка контекста panic
    var p _panic
    p.arg = e
    ...
    p.start(getcallerpc(), unsafe.Pointer(getcallersp()))

    for {
        fn, ok := p.nextDefer()
        if !ok {
            break
        }
        fn()
    }

    preprintpanics(&p)
    fatalpanic(&p)
    *(*int)(nil) = 0 // не достижимо
}
// panic(any) или gopanic(any) (src/runtime/panic.go)_
```

Go runtime создаёт запись `_panic` (которую мы уже видели в `deferreturn`) и захватывает program counter и stack pointer. Запись `_panic` добавляется в **список всех активных panic** в текущей goroutine — как записи `_defer`. 
Это обрабатывается в `_panic.start()` и позволяет **вложенные panic**.

Процесс очень похож на то, что мы обсуждали с `deferreturn`
Функция итерирует по defer chain через `nextDefer`, забирая и выполняя каждую deferred-функцию. После обхода chain вызывается `preprintpanics` для подготовки всех значений panic к печати, затем `fatalpanic(&p)` для **аварийного завершения** программы.

---
### Как recover меняет control flow

Если смотреть только на control flow, кажется, что `fatalpanic(&p)` **всегда** выполняется и программа всегда падает. Но мы знаем, что если deferred-функция корректно вызовет `recover()`, panic может быть **остановлен**.

Runtime нужен способ перехватить и **изменить flow**. 
Этот механизм находится в функции `p.nextDefer()`:

```go
func (p *_panic) nextDefer() (func(), bool) {
    gp := getg()

    if !p.deferreturn {
        ...
        if p.recovered {
            mcall(recovery) // <- Это процесс восстановления
            throw("recovery failed")
        }
    }

    p.argp = add(p.startSP, sys.MinFrameSize)

    for {
        ...
    Recheck:
        if d := gp._defer; d != nil && d.sp == uintptr(p.sp) {
            ...
            fn := d.fn
            d.fn = nil
            p.retpc = d.pc

            // Отвязываем и освобождаем _defer.
            gp._defer = d.link
            freedefer(d)
            return fn, true
        }

        // Если текущий _defer не в текущем frame,
        // переходим к предыдущему frame.
        if !p.nextFrame() {
            return nil, false
        }
    }
}
// nextDefer (src/runtime/panic.go)_
```

Ключевая часть — **флаг `p.recovered`**. Если этот флаг `true`, panic был восстановлен в deferred-функции.

![](../../assets/obsidian-images/Pasted%20image%2020260219170507.png)

Функция обходит defer chain и возвращает следующую deferred-функцию для выполнения при раскрутке stack.

Когда вы вызываете `recover()` внутри deferred-функции, это транслируется в `runtime.gorecover(...)`, который устанавливает флаг `recovered` записи `_panic` в `true`.

---
