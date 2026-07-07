
Может показаться, что компилятор Go просто перемещает deferred-вызов в конец функции, перед каждым `return`. 
Для простейших случаев эта ментальная модель работает:

|Go Source Code|Compiler Transformation|
|---|---|
|`defer println(1)`|`println(4)`|
|`defer println(2)`|`println(3)`|
|`defer println(3)`|`println(2)`|
|`println(4)`|`println(1)`|

Однако **это не то, как defer работает** на самом деле, особенно в сложных случаях. 
Go требует **дополнительной runtime-поддержки**.

Одна из проблем — компилятор часто **не может знать** заранее, какие deferred-вызовы нужно выполнить или сколько раз `defer` будет вызван:

```go
func do(b1 bool) {
    if b1 {
        defer fmt.Println("b1")
    }
    defer fmt.Println("b2")
}
```

Здесь `defer fmt.Println("b1")` выполняется **только если** `b1` равно `true` в runtime. 
Компилятор не может знать это заранее.

Кроме того, аргументы deferred-функции **всегда вычисляются сразу** при достижении `defer`
Они не откладываются до момента выполнения.

Вместо использования только компилятора Go использует **runtime-систему** для отслеживания deferred-вызовов — через stack goroutine и специальные runtime-функции.

---
### Структура _defer

Когда программа достигает `defer`, она создаёт запись `_defer` для запоминания, что нужно сделать позже:

```go
type _defer struct {
    heap      bool                    // Флаг: stack или heap allocation.
    rangefunc bool                    // Специальный флаг range-over-func.
    sp        uintptr                 // Stack pointer при defer.
    pc        uintptr                 // Program counter при defer.
    fn        func()                  // Deferred-функция для выполнения.
    link      *_defer                 // Следующий defer в LIFO-списке.
    head      *atomic.Pointer[_defer] // Head atomic rangefunc списка.
}
// The Defer Record (src/runtime/runtime2.go)_
```

Эта запись хранит всю информацию для последующего выполнения deferred-функции и отслеживает порядок. Три важных поля: **link**, **fn** и **heap**.

Поле `link` указывает на **следующую** запись `_defer`, создавая простой linked list:

![](../../assets/obsidian-images/Pasted%20image%2020260219112545.png)

Каждая goroutine хранит самую последнюю запись `_defer` для текущего function frame:

```go
type g struct {
    ...
    _defer *_defer // самый внутренний defer
}
```

Может показаться странным, что deferred-функции группируются по **goroutine**, а не по функции. Однако `return` из функции — не единственный способ запуска defer. 
При panic **все** defer в текущем goroutine stack также выполняются.

Эта структура образует **singly linked list** записей `_defer` для всего call stack. 
Поле `_defer` goroutine — **head** этого списка:

![](../../assets/obsidian-images/Pasted%20image%2020260219113014.png)
### Нормализация deferred-функции через closure

Поле `fn func()` **не имеет** аргументов и return-значений. 
Как это работает, если вы defer-ите функцию с аргументами?

```go
func main() {
    defer func(a, b int) int {
        return a + b
    }(1, 2)
}
```

Компилятор Go использует **closure** и процесс «нормализации» функции:

1. Компилятор **немедленно вычисляет** аргументы `arg1` и `arg2` при достижении `defer`.
2. Создаёт **новую, скрытую** анонимную функцию (closure) без аргументов.
3. Closure **захватывает** значения `arg1` и `arg2`.
4. Closure просто вызывает оригинальную функцию с захваченными значениями.

| Go                  | Pseudo Go        |
| ------------------- | ---------------- |
| `defer Hello(a, b)` | `a1 := a`        |
|                     | `b1 := b`        |
|                     | `defer func() {` |
|                     | `Hello(a1, b1)`  |
|                     | `}()`            |

Эти closure видны в stack trace с именами вроде `<OuterFunctionName>.deferwrap<N>`.

---
### Три вида defer

В Go существует **три вида** defer-вызовов: 
- **open-coded defer**
- **stack defer**
- **heap defer**

```go
func (s *state) stmt(n ir.Node) {
    ...
    case ir.ODEFER:
        n := n.(*ir.GoDeferStmt)
        if s.hasOpenDefers {
            defertype = "open-coded"
        } else if n.Esc() == ir.EscNever {
            defertype = "stack-allocated"
        } else {
            defertype = "heap-allocated"
        }
        ...
}
// Calling stmt defer (src/cmd/compile/internal/ssagen/ssa.go)_
```

Приоритет и производительность: **open-coded** > **stack-allocated** > **heap-allocated**. Основное различие между stack и heap — где хранится запись `_defer`:

- **Stack-allocated defer:** Запись `_defer` хранится **на stack функции**. Быстрее, так как избегает heap allocation и garbage collection.
- **Heap-allocated defer:** Запись `_defer` хранится **на heap**. Используется, когда компилятор не может определить заранее, сколько `defer` будет создано (цикл, conditional).

---

## Heap-Allocated Defers

Escape analysis определяет, попадает ли `defer` на heap:

```go
func (e *escape) goDeferStmt(n *ir.GoDeferStmt) {
    k := e.heapHole()
    if n.Op() == ir.ODEFER && e.loopDepth == 1 && n.DeferAt == nil {
        // Top-level defer аргументы не escape на heap,
        // но должны сохраняться до вызова.
        k = e.later(e.discardHole())
        n.SetEsc(ir.EscNever)
    }
    // ...
}
// goDeferStmt (src/cmd/compile/internal/escape/call.go)_
```

Ключевое условие — `e.loopDepth == 1`. 
Компилятор отслеживает `loopDepth`, начинающийся с 1 (тело основной функции) и увеличивающийся при входе в цикл.

Если `defer` на **верхнем уровне** функции (`loopDepth == 1`), escape analysis устанавливает статус `EscNever` — запись `_defer` будет **stack-allocated**. Если `defer` внутри цикла — **heap-allocated**.

Пример heap-allocated defer:

```go
package main

func main() {
    for i := 0; i < 3; i++ {
        defer println("Defer in loop", i)
    }
}
```

Для проверки:

```bash
go build -gcflags=-d=defer=1
# theanatomyofgo
./main.go:5:3: heap-allocated defer
```

---
### Assembly heap-allocated defer

Каждый heap-allocated defer приводит к затратному runtime-вызову `runtime.deferproc`:

```asm
0x0034 00052 MOVD $type:noalg.struct { F uintptr; X0 int }(SB), R0
0x003c 00060 CALL runtime.newobject(SB)
0x0040 00064 MOVD $main.main.deferwrap1(SB), R1
0x0048 00072 MOVD R1, (R0)
0x004c 00076 MOVD main.i-8(SP), R2
0x0050 00080 MOVD R2, 8(R0)
0x0054 00084 CALL runtime.deferproc(SB)
0x0058 00088 CMP $0, R0
0x005c 00092 BNE 100
0x0060 00096 JMP 32
0x0064 00100 CALL runtime.deferreturn(SB)
...
```

Первый шаг — создание **closure** для deferred-функции. 
Это closure-объект (называемый `funcval`) содержит function pointer `F`, указывающий на wrapper-функцию, и поле аргумента `X0`, хранящее захваченное значение переменной цикла `i`:

> **Иллюстрация 194.** Go closure создан как heap-allocated объект _(Рисунок не загружен)_

Два ключевых runtime-вызова: `CALL runtime.deferproc(SB)` и `CALL runtime.deferreturn(SB)`:

> **Иллюстрация 195.** runtime.deferproc планирует deferred-вызов _(Рисунок не загружен)_

---
### runtime.deferproc

Функция `runtime.deferproc` создаёт запись `_defer` и связывает её с defer chain текущей goroutine:

```go
func deferproc(fn func()) {
    gp := getg()
    if gp.m.curg != gp {
        throw("defer on system stack")
    }
    d := newdefer()
    d.link = gp._defer
    gp._defer = d
    d.fn = fn
    d.pc = getcallerpc()
    d.sp = getcallersp()
    return0()
}
```

_Figure 117. deferproc (src/runtime/panic.go)_

Сначала runtime получает текущую goroutine (`getg()`)
Затем создаёт новую запись `_defer` (`newdefer()`). 
Новая запись **связывается** с defer chain goroutine через поле `link`:

> **Иллюстрация 196.** Все активные defer связаны в текущей goroutine _(Рисунок не загружен)_

Запись `_defer` также сохраняет **program counter** (`pc`) и **stack pointer** (`sp`) текущего function frame — эти поля помогают runtime корректно обрабатывать deferred-вызовы и panic.

---
### Скрытое return-значение через return0()

Тонкая деталь: хотя `deferproc` определён **без return-значения**, assembly всё равно проверяет return-значение в регистре R0:

Если return-значение **не 0** — функция возвращается из-за panic, и код переходит к `runtime.deferreturn`.

Как это работает, если Go-сигнатура `deferproc` не имеет return-значения? 
Здесь вступает **специальная функция** `return0()`:

```go
func deferproc(fn func()) {
    ...
    return0()
}
```

Это **assembly stub**, который возвращает 0 в машинный регистр напрямую, без использования обычной Go return-системы:

```asm
TEXT runtime·return0(SB), NOSPLIT, $0
    MOVW $0, R0
    RET
```


Этот трюк делает функцию «возвращающей значение» на уровне assembly, сохраняя **void-сигнатуру** на уровне Go. Если `deferproc` возвращает 0 — выполнение продолжается нормально. Если 1 (при panic recovery) — переход к концу функции, где вызывается `runtime.deferreturn`, **пропуская** остаток тела функции.

> **Примечание: Почему deferproc не возвращает значение обычным Go-способом?**
> 
> На уровне языка `defer` должен быть «невидимым» — обычный statement, не возвращающий значение. Внутренне runtime нужен способ передать один бит информации: был ли panic recovered? Для этого runtime устанавливает значение напрямую в первый return-регистр через assembly stub. Если бы `deferproc` имел `int` return в Go-сигнатуре, callers должны были бы обрабатывать return-значение, раскрывая implementation detail.

---
### Per-Processor Pooling (Defer Pool)

Может показаться, что каждый defer — это **два** allocation: одно для closure, другое для `_defer`. Но Go runtime использует **per-processor pool** для кеширования часто используемых `_defer` структур:

```go
type p struct {
    ...
    deferpool    []*_defer
    deferpoolbuf [32]*_defer
    ...
}
```

Каждый логический processor (P) имеет собственный **локальный кеш** `_defer` структур (`p.deferpool`).

Когда нужен новый `_defer`, runtime сначала проверяет **локальный pool**. Это per-processor кеширование **избегает** synchronization overhead при allocation и освобождении `_defer` записей. Каждый P может обращаться к своему pool **без lock** — нет contention, так как только один thread обращается к каждому pool в момент времени.

Когда локальный pool пуст, он заполняется из **global pool** (`sched.deferpool`):

![](../../assets/obsidian-images/Pasted%20image%2020260219123424.png)

Извлечение из global pool **требует lock**, поскольку несколько thread могут обращаться одновременно. Для минимизации lock contention локальный pool забирает **половину** своей capacity (16 элементов) за один раз:

```go
func newdefer() *_defer {
    var d *_defer
    mp := acquirem()
    pp := mp.p.ptr()

    // Проверяем, пуст ли локальный pool, пытаемся заполнить из global
    if len(pp.deferpool) == 0 && sched.deferpool != nil {
        ...
        // Перемещаем половину defer из global pool в локальный
        for len(pp.deferpool) < cap(pp.deferpool)/2 && sched.deferpool != nil {
            d := sched.deferpool
            sched.deferpool = d.link
            d.link = nil
            pp.deferpool = append(pp.deferpool, d)
        }
        ...
    }

    // Pop defer из локального pool
    if n := len(pp.deferpool); n > 0 {
        d = pp.deferpool[n-1]
        pp.deferpool[n-1] = nil
        pp.deferpool = pp.deferpool[:n-1]
    }

    releasem(mp)
    mp, pp = nil, nil

    // Если pool всё ещё пуст — allocate с heap
    if d == nil {
        d = new(_defer)
    }

    d.heap = true
    return d
}
```

Если и локальный, и global pool пусты — runtime выделяет **свежий** `_defer` с heap (`new(_defer)`). 
Heap allocation дороже и создаёт давление на GC.

---
### Обратная связь: локальный → global pool

Если локальный pool P становится **слишком полным** (более 32 элементов), около половины `_defer` записей перемещается **обратно** в global pool:

![](../../assets/obsidian-images/Pasted%20image%2020260219123734.png)

Как Go предотвращает **бесконечный рост** global pool? Garbage collector. 
В начале каждого GC-цикла он **отвязывает** весь global list, очищает `link` поля и сбрасывает head pointer в `nil`. 
Вся партия становится обычным мусором и освобождается в следующей sweep-фазе.

---
## Выполнение Deferred-вызовов (Defer Return)

При возврате из функции с `defer` вызывается `runtime.deferreturn`, который удаляет `_defer` записи из defer chain goroutine и **выполняет** их по одной, начиная с самой последней (head linked list):

```go
func deferreturn() {
    var p _panic
    p.deferreturn = true
    p.start(getcallerpc(), unsafe.Pointer(getcallersp()))
    for {
        fn, ok := p.nextDefer()
        if !ok {
            break
        }
        fn()
    }
}
// deferreturn (src/runtime/panic.go)_
```

### Определение принадлежности _defer к функции

Как runtime узнаёт, какие `_defer` записи принадлежат **текущей функции**, если все записи управляются per-goroutine?

При каждом создании defer runtime сохраняет **stack pointer caller**: `d.sp = getcallersp()`. Stack pointer (sp) действует как **уникальный идентификатор** для callee frame:

![](../../assets/obsidian-images/Pasted%20image%2020260219124043.png)

При возврате из функции runtime проверяет, какие `_defer` записи имеют `sp`, **совпадающий** с текущим frame. Только эти defer выполняются.

Метод `p.nextDefer()` также **утилизирует** запись `_defer` обратно в processor pool после завершения deferred-вызова. Он возвращает `false`, когда больше нет записей для текущего stack frame.

---

## Stack-Allocated Defers

Heap-allocated defer имеют **стоимость**: heap allocation и overhead кеш-менеджмента. 
Но `defer` встроен в Go и используется повсеместно. 
В высоконагруженных системах с тысячами или миллионами goroutine дополнительная стоимость heap allocation **быстро накапливается**. Для решения этой проблемы, начиная с Go 1.13, команда Go ввела **«stack-allocated defers»** как оптимизацию.

Как создать stack-allocated defer на практике? 
Если функция имеет хотя бы один heap-allocated defer, то другие простые `defer` в той же функции **не используют** open-coded defer, а становятся **stack-allocated**:

```go
//go:noinline
func printsmt() {
    println("Oh hey! I'm stack-allocated deferred!")
}

func main() {
    // Stack-allocated defer
    defer printsmt()

    // Heap-allocated defers
    for i := 0; i < 3; i++ {
        defer println(i)
    }
}
```

Если deferred-функция **не принимает** аргументов и не возвращает значений (как `printsmt()`), компилятору не нужно создавать wrapper или closure. 
Runtime может вызвать `deferproc` напрямую с function pointer.

---

## Open-Coded Defers

Open-coded defer — оптимизация, введённая в **Go 1.14**, которая может **устранить runtime overhead** defer во многих случаях. Вместо создания `_defer` записей в runtime, компилятор помещает deferred-вызовы **непосредственно** в каждую точку выхода функции.

### Условия для open-coded defer

Для применения open-coded оптимизации должны выполняться **несколько правил**:

1. **Флаг `-N` не указан.** При `-N` большинство оптимизаций отключены, включая open-coded defer.
2. **Функция не имеет heap-allocated result parameters.** Если функция возвращает heap-allocated результаты, результат должен копироваться обратно в stack slot при каждом выходе — это добавляет шаги:

```go
if s.hasOpenDefers {
    for _, f := range s.curfn.Type().Results() {
        if !f.Nname.(*ir.Name).OnStack() {
            s.hasOpenDefers = false
            break
        }
    }
}
```

3. **Произведение return point × defer ≤ 15.** Компилятор должен inline defer-вызов при каждом выходе. Слишком много комбинаций создают **раздутый код**.
4. **Число defer в функции ≤ 8.** Go использует **один uint8 bitmap**, называемый `deferBits`, для отслеживания всех `defer` в функции:
    
    - Первый defer использует bit 0 (least significant bit)
    - Второй — bit 1
    - Третий — bit 2, и так далее до bit 7
    
    `uint8` имеет **только 8 бит**, поэтому только 8 defer можно отслеживать. Если больше — компилятор отключает open-coded defer и fallback на обычный runtime defer list.
3. **Нет heap-allocated defer в той же функции.** Если есть хотя бы один heap-allocated defer — open-coded defer **не используется** для остальных:

```go
func walkStmt(n ir.Node) ir.Node {
    ...
    case ir.ODEFER:
        n := n.(*ir.GoDeferStmt)
        ir.CurFunc.SetHasDefer(true)
        ir.CurFunc.NumDefers++
        if ir.CurFunc.NumDefers > maxOpenDefers || n.DeferAt != nil {
            ir.CurFunc.SetOpenCodedDeferDisallowed(true)
        }
        if n.Esc() != ir.EscNever {
            ir.CurFunc.SetOpenCodedDeferDisallowed(true)
        }
    ...
}
```

### Как работает deferBits

При выполнении функции каждый раз, когда control достигает `defer`, компилятор генерирует **два вида инструкций**:

- Одна сохраняет **function pointer** в фиксированный stack slot, зарезервированный для этого defer.
- Одна устанавливает **соответствующий бит** в `deferBits` операцией OR.

Пример:

```go
func main() {
    defer func() { println("defer 1") }()
    defer func() { println("defer 2") }()
    defer func() { println("defer 3") }()
    println("Hello, World!")
}
```

После первого defer байт содержит `0000 0001₂` (1). После второго — `0000 0011₂` (3). После третьего — `0000 0111₂` (7):

| Go          | Assembly                          | Pseudo Go           |
| ----------- | --------------------------------- | ------------------- |
| defer func1 | `MOVD $main.main.func1·f(SB), R0` | `R0 = &func1`       |
|             | `MOVD R0, main..autotmp_1-24(SP)` | `stack_slot_1 = R0` |
|             | `MOVD $1, R0`                     | `R0 = 1`            |
|             | `MOVB R0, main..autotmp_0-25(SP)` | `deferBits = 1`     |
| defer func2 | `MOVD $main.main.func2·f(SB), R1` | `R1 = &func2`       |
|             | `MOVD R1, main..autotmp_2-16(SP)` | `stack_slot_2 = R1` |
|             | `MOVD $3, R1`                     | `R1 = 3`            |
|             | `MOVB R1, main..autotmp_0-25(SP)` | `deferBits = 3`     |
| defer func3 | `MOVD $main.main.func3·f(SB), R2` | `R2 = &func3`       |
|             | `MOVD R2, main..autotmp_3-8(SP)`  | `stack_slot_3 = R2` |
|             | `MOVD $7, R2`                     | `R2 = 7`            |
|             | `MOVB R2, main..autotmp_0-25(SP)` | `deferBits = 7`     |

### Epilogue: выполнение в обратном порядке

При каждом не-panic выходе компилятор генерирует **inline-код**, обходящий биты от **старшего к младшему**. Для каждого установленного бита он сбрасывает бит и вызывает сохранённый function pointer:

```asm
# Вызов func3
00124 MOVD $3, R0            # clear bit 2: 7 &^ (1<<2) == 3
00128 MOVB R0, main..autotmp_0-25(SP)
00132 CALL main.main.func3(SB)

# Вызов func2
00136 MOVD $1, R0            # clear bit 1: 3 &^ (1<<1) == 1
00140 MOVB R0, main..autotmp_0-25(SP)
00144 CALL main.main.func2(SB)

# Вызов func1
00148 MOVB ZR, main..autotmp_0-25(SP)
00152 CALL main.main.func1(SB)

# Завершение
00156 LDP -8(RSP), (R29, R30)
00160 ADD $64, RSP
00164 RET (R30)
```

Обход от старшего к младшему обеспечивает обычный **LIFO-порядок**: самый последний зарегистрированный defer выполняется первым.

В данном примере компилятор может **hardcode** биты 1, 3 и 7 и вызывать напрямую. 
В более сложных случаях, когда defer скрыт за branch, компилятор **не может знать** на compile-time, должен ли бит быть установлен, и использует условные инструкции.

---

## Xulosa

**1. Defer — bu Go'ning eng elegant, lekin eng noto'g'ri tushunilgan feature'laridan biri.** 
Ko'pchilik dasturchilar `defer` ni «funksiya oxirida chaqiriladigan cleanup» deb biladi. 
Lekin **argumentlar darhol baholanadi**, receiver value **copy** qilinadi, va `defer` **funksiyaga** bog'langan, blokka emas. 
Bu uchta xususiyat eng ko'p bug manbayi. Pointer receiver yoki closure ishlatish — yagona to'g'ri yechim agar latest value kerak bo'lsa.

**2. Loop ichida defer — klassik xato.** `defer` loop ichida bo'lsa, barcha deferred-chaqiruvlar **funksiya oxirigacha** kutadi. 1000 ta faylni loop'da ochib `defer f.Close()` qilsangiz — **barcha 1000 fayl** ochiq qoladi. Yechim: anonymous function wrapper ishlatish, har bir iteratsiyada defer **o'sha wrapper ichida** execute bo'ladi.

**3. Uchta defer turi — performance gradient.**

|Tur|Qachon|Tezlik|Mexanizm|
|---|---|---|---|
|**Open-coded**|Oddiy, ≤8 defer, loop yo'q|**Eng tez**|Function pointer stack slot'da, `deferBits` bitmap, inline epilogue|
|**Stack-allocated**|Top-level, `EscNever`|**O'rtacha**|`_defer` record stack'da, `deferprocStack()`|
|**Heap-allocated**|Loop ichida, conditional|**Eng sekin**|`_defer` record heap'da, `deferproc()` + per-P pool|

**4. `_defer` record — linked list elementlari.** 
Har bir `defer` `_defer` struct yaratadi: `{fn, sp, pc, link, heap}`.
Barcha `_defer` lar goroutine'ning `g._defer` field'idan boshlanuvchi **singly linked list**. 
`link` field keyingi `_defer` ga ko'rsatadi. 
`sp` field qaysi funksiyaga tegishli ekanligini aniqlaydi — `deferreturn` faqat **joriy frame'ning sp'siga mos** keladigan `_defer` larni execute qiladi.

**5. Per-Processor Pool — allocation overhead'ni kamaytirish.** Har bir P (logical processor) o'z `deferpool` (32 element capacity) ga ega. Yangi `_defer` kerak bo'lganda: avval **lokal pool** tekshiriladi (lock-free), keyin **global pool** (lock bilan, 16 ta birdan oladi), eng oxirida **heap allocation**. Lokal pool to'lsa — yarmisi global pool'ga qaytariladi. GC har cycle'da global pool'ni **to'liq tozalaydi** — unbounded retention yo'q.

**6. `return0()` — Go runtime'ning eng qiziqarli assembly trick'laridan biri.** `deferproc` Go darajasida void return, lekin assembly darajasida R0 registrda 0 qaytaradi. Agar panic recover bo'lsa — runtime R0 ni 1 ga o'zgartiradi va execution `deferproc` dan keyin davom etadi. Compiler `CMP $0, R0` va `BNE` bilan buni tekshiradi. Bu Go tili **toza** qoladi (defer statement value qaytarmaydi), lekin runtime kerakli control'ni oladi.

**7. Open-coded defer — deferBits bitmap orqali zero-cost defer.** `uint8` bitmap har bir defer uchun 1 bit: bit 0 = 1-defer, bit 1 = 2-defer, ... bit 7 = 8-defer. Funksiya oxirida compiler **inline epilogue** generatsiya qiladi — yuqoridan pastga bit'larni tekshiradi, har bir set bit uchun `CALL` qiladi. **Hech qanday runtime overhead yo'q** — na `deferproc`, na `deferreturn`, na linked list, na pool. Faqat bitta `MOVB` + `CALL`. Lekin shartlar qat'iy: ≤8 defer, ≤15 (returns × defers), heap-allocated defer yo'q, `-N` flag yo'q.

**8. Named return value bilan defer — kuchli pattern.** `defer func() { res = 100 }()` bilan named return value'ni **o'zgartirish mumkin**. Bu panic recovery'da eng ko'p ishlatiladi: `defer func() { if r := recover(); r != nil { err = fmt.Errorf(...) } }()`. Bu pattern Go idiomatic error handling'ning asosi.

**Amaliy maslahatlar:**

- `defer` argumentlarini **darhol** baholanishini doim esda tuting — pointer yoki closure ishlating
- Loop ichida `defer` — **anonymous function wrapper** ishlatmasangiz resource leak bo'ladi
- `go build -gcflags=-d=defer=1` bilan defer turini tekshiring — performance-critical kodda muhim
- Oddiy funksiyalarda (≤8 defer, loop yo'q) `defer` **bepul** — open-coded optimize qilinadi
- Named return value + defer = **idiomatic error recovery** pattern
- `defer f.Close()` yozganingizda error'ni ham handle qiling: `defer func() { if err := f.Close(); err != nil { ... } }()`