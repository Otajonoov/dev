
Когда Go-программа запускается, ядро отображает исполняемый файл в **виртуальное адресное пространство** процесса. Не каждый байт файла становится частью этого in-memory образа, и даже отображённые части обычно подгружаются **постранично по требованию** (page-by-page on demand). 
Чтобы понять почему, рассмотрим структуру самого бинарника.

Когда мы собираем исполняемый файл, результирующий бинарник разделён на **секции**. Эти секции можно инспектировать простой командой:

```bash
# Linux arm64/amd64
$ readelf -S main

# Darwin arm64
$ otool -l main
```

Layout незначительно различается между архитектурами, но ключевые секции **остаются теми же**. Таблица соответствия Linux и Darwin:

|Linux|Darwin|
|---|---|
|`.text`|`__text`|
|`.rodata`|`__rodata`|
|`.typelink`|`__typelink`|
|`.itablink`|`__itablink`|
|`.noptrdata`|`__noptrdata`|
|`.data`|`__data`|
|`.bss`|`__bss`|
|`.noptrbss`|`__noptrbss`|

Это называется **section header**. Section header описывают логический layout файла — какие секции существуют (код, данные, таблицы символов, отладочная информация) и где они расположены в бинарнике. Section header полезны для инструментов вроде linker, loader или debugger, которым нужно детально понимать структуру файла.

Каждая секция имеет **флаги**, указывающие на использование памяти и разрешения:

```
  Sections                Flags
  [ 1] .text              AX
  [ 2] .rodata            A
  ...
  [ 9] .noptrdata         WA
  [10] .data              WA
  [11] .bss               WA
  [12] .noptrbss          WA
```

Значения этих флагов:

- **A (allocatable):** секция должна существовать в виртуальном адресном пространстве процесса.
- **X (executable):** секция содержит инструкции, которые CPU может выполнять (например, `.text`).
- **W (writable):** секция содержит данные, которые можно изменять (например, глобальные переменные в `.data` и `.bss`).

Однако флаг allocatable (A) лишь говорит, что linker **может** поместить секцию в загружаемый segment. Что фактически загружается в память, определяется **program header**, а не section header. Ядро читает program header для нахождения адресов segment, размеров памяти и разрешений (read, write, execute).

---
### Program Headers

Предполагая, что наш бинарник нацелен на Linux/amd64, вот выдержка из program header:

```bash
$ readelf -l main

Elf file type is EXEC (Executable file)
Entry point 0x46dcc0
There are 6 program headers, starting at offset 64

Program Headers:
  Type      Offset             VirtAddr           PhysAddr
            FileSiz            MemSiz             Flags  Align
  PHDR      0x0000000000000040 0x0000000000400040 0x0000000000400040
            0x0000000000000150 0x0000000000000150 R      0x1000
  NOTE      0x0000000000000f78 0x0000000000400f78 0x0000000000400f78
            0x0000000000000064 0x0000000000000064 R      0x4
  LOAD      0x0000000000000000 0x0000000000400000 0x0000000000400000
            0x000000000006fb31 0x000000000006fb31 R E    0x1000
  LOAD      0x0000000000070000 0x0000000000470000 0x0000000000470000
            0x0000000000090d38 0x0000000000090d38 R      0x1000
  LOAD      0x0000000000101000 0x0000000000501000 0x0000000000501000
            0x0000000000003de0 0x0000000000038ac0 RW     0x1000
  GNU_STACK 0x0000000000000000 0x0000000000000000 0x0000000000000000
            0x0000000000000000 0x0000000000000000 RW     0x8

Section to Segment mapping:
  Segment Sections...
   00
   01     .note.go.buildid
   02     .text .note.gnu.build-id .note.go.buildid
   03     .rodata .typelink .itablink .gosymtab .gopclntab
   04     .go.buildinfo .go.fipsinfo .noptrdata .data .bss .noptrbss
   05
```

Program header сообщают операционной системе, **как загрузить** исполняемый файл в память и создать работающий процесс. Без них loader не знал бы, как отобразить файл в память, и выполнение **завершилось бы неудачей**.

Ядро создаёт memory mapping для **трёх LOAD segment** (segment 02, 03 и 04):

- **Segment 02** загружает исполняемый код начиная с виртуального адреса `0x400000` с разрешениями **read и execute**, содержащий секцию `.text` и build ID notes.
- **Segment 03** загружает read-only данные начиная с виртуального адреса `0x470000`, содержащий `.rodata`, `.typelink`, `.itablink`, `.gosymtab`, `.gopclntab`.
- **Segment 04** загружает read-write данные начиная с виртуального адреса `0x501000`, содержащий `.go.buildinfo`, `.go.fipsinfo`, `.noptrdata`, `.data`, `.bss` и `.noptrbss`.

Остальные записи тоже важны, но работают **иначе**:

- **PHDR** указывает на саму таблицу program header. Ядро читает эту таблицу, чтобы знать, как отобразить файл. Многие linker размещают таблицу в начале файла, поэтому её байты часто попадают **внутрь первого LOAD segment** и отображаются как его часть.
- **NOTE** описывает область notes — build ID и другие metadata. Байты файла для notes часто размещаются **внутри** некоторого LOAD segment и отображаются таким образом, хотя тип program header entry — NOTE, а не LOAD.
- **GNU_STACK** — это вообще не файловые данные. Это просто **флаг**, указывающий ядру, какие разрешения должен иметь stack процесса (например, non-executable).

---
### Зачем нужны не-allocatable секции

Если program header решают, что отображается в runtime, зачем файл всё ещё включает non-allocatable секции?

```
[ 0]                      -
[ 1] .text                AX
[ 2] .rodata              A
[ 3] .typelink            A
[ 4] .itablink            A
[ 5] .gosymtab            A
[ 6] .gopclntab           A
[ 7] .go.buildinfo        WA
[ 8] .go.fipsinfo         WA
[ 9] .noptrdata           WA
[10] .data                WA
[11] .bss                 WA
[12] .noptrbss            WA
[13] .debug_abbrev        -
[14] .debug_line          -
[15] .debug_frame         -
[16] .debug_gdb_scripts   -
[17] .debug_info          -
[18] .debug_loclists      -
[19] .debug_rnglists      -
[20] .debug_addr          -
[21] .note.gnu.build-id   A
[22] .note.go.buildid     A
[23] .shstrtab            -
[24] .symtab              -
[25] .strtab              -
```

Конкретный пример: когда **pprof** собирает sample, он записывает адрес program counter (адрес инструкции, которая выполнялась). Чтобы сделать это полезным, он ищет этот адрес в отладочной информации бинарника и транслирует его в соответствующее местоположение в исходном коде, показываемое как `file:line` (например, `theanatomyofgo/main.go:158`).

Эта информация поступает из **двух источников**: собственная PC/line таблица Go и стандартные DWARF debug-секции вроде `.debug_info` и `.debug_line`.

Go хранит свою PC/line таблицу в `.gopclntab` — **Go-специфичная** таблица, отображающая program counter на строку исходного кода:

```
$ readelf -S main
Section headers:
  ...
  [ 6] .gopclntab    PROGBITS  00000000004a49a0 000a49a0
       000000000005c398 0000000000000000   A   0   0  32
```

Эта секция создаётся Go linker и **загружается в память**. Благодаря этой runtime-таблице вы видите информацию `file:line` в stack trace при panic; Go также использует её для `runtime.Caller()`. При символизации профилей `cmd/pprof` тоже читает эти таблицы (вместе с `.gosymtab` или эквивалентными данными символов) из исполняемого файла.

Однако эта таблица очень **компактная** и в основном предназначена для работы с Go-функциями. Для более детальной информации, особенно для кода, написанного **не на Go**, существуют DWARF debug-секции, которые **не загружаются** в память. Они служат мостом между Go-специфичной таблицей и стандартными инструментами отладки, ожидающими DWARF.

Два практических напоминания:

- `cmd/pprof` использует Go-таблицы для Go-функций; если адрес не может быть разрешён таким образом (обычно C/cgo frame), он пробует DWARF, если тот доступен.
- Если вы собираете с `-w`, linker **опускает DWARF-секции**; Go stack trace всё ещё работают (они полагаются на `.gopclntab`), но DWARF-based символизация **не будет работать**.

---

## Виртуальная память процесса

Это обзор Go бинарника; теперь рассмотрим, что фактически находится в виртуальной памяти. Ниже — упрощённая диаграмма layout виртуальной памяти:

Эти segment отображаются в виртуальную память ядром. Исполняемая область TEXT обычно начинается с базы вроде `0x400000` для non-PIE бинарников. 
Выбор `0x400000` (4 MiB) как базового адреса — **общепринятая конвенция** для x86_64 Linux исполняемых файлов. Она обеспечивает разумный промежуток от null pointer (`0x0`), не потребляя слишком много нижнего адресного пространства.

Разобравшись с layout файла, пройдёмся по тому, как эти segment выглядят после запуска процесса.

---

### TEXT

Рассмотрим пример:

```go
func main() {
    fmt.Println(add)
}

//go:noinline
func add(x, y int) int {
    return x + y
}
```

Когда мы пишем `add` без скобок, мы **не вызываем** её. Мы ссылаемся на саму функцию — **function value**.

На практике function value содержит адрес, по которому CPU может перейти для начала выполнения этой функции. При печати Go показывает его как hex-число вроде `0x491b00`. Думайте об этом как о «местоположении» внутри процесса.

На Linux/amd64, если исполняемый файл не собран как PIE, программа отображается по **фиксированному базовому адресу** при каждом запуске. Поскольку база фиксирована, адрес `main.add` **стабилен** между запусками:

```bash
$ ./main
0x491b00
$ ./main
0x491b00
$ ./main
0x491b00
```

Чтобы понять, откуда берётся этот адрес, инспектируем символы исполняемого файла с помощью `go tool nm`:

> **Примечание.** `go tool nm` — утилита командной строки, перечисляющая символы, определённые или используемые объектными файлами, archive или исполняемыми файлами. Это Go-эквивалент Unix-команды `nm`, разработанный для работы с Go-бинарниками. Мы использовали этот инструмент в Chapter 4 при обсуждении interface table, хранимой в секции `.itablink`.

```bash
$ go tool nm ./main | grep main.add
  491b00 T main.add
```

`go tool nm` выводит три поля: **адрес**, **буква типа** и **имя символа**. Здесь `T` означает, что символ находится в области исполняемого кода.

Важная идея: linker выбирает виртуальный адрес для каждой функции, `go tool nm` показывает этот адрес, и на non-PIE бинарнике вы наблюдаете **тот же адрес** в runtime, потому что исполняемый файл отображается по фиксированной базе.

---

### PIE и ASLR

**PIE** (Position-Independent Executable) — бинарник, собранный так, чтобы корректно работать, даже если ОС отображает весь исполняемый файл по **другому базовому адресу** при каждом запуске. Функция всё ещё по тому же offset внутри бинарника, но сам бинарник может быть сдвинут в памяти, поэтому финальный адрес при печати **меняется**.

Напротив, **non-PIE** бинарник слинкован для фиксированного базового адреса. Если ОС отображает его по этой фиксированной базе, абсолютный адрес функции **стабилен** между запусками.

На macOS/arm64 Go исполняемые файлы являются **PIE**. При каждом запуске приложение выводит другой адрес:

bash

```bash
$ ./main
0x1003f4700
$ ./main
0x10230c700
$ ./main
0x100814700
$ ./main
0x104a0c700
```

Причина в том, что на macOS/arm64 Go-исполняемые файлы — PIE. PIE executable разработан для работы **независимо от того**, куда ОС помещает его в память.

**ASLR** (Address Space Layout Randomization) — функция ОС, намеренно **меняющая** расположение в памяти между запусками. С PIE, ASLR может рандомизировать и базовый адрес основного исполняемого файла: функция остаётся по тому же offset внутри бинарника, но весь бинарник **сдвигается** на новую базу в runtime, поэтому финальный адрес функции при печати меняется при каждом запуске.

На macOS динамический loader (**dyld**) — компонент, применяющий этот сдвиг адреса при отображении программы в память, что объясняет, почему `fmt.Println(add)` выводит **разное число** каждый раз.

---

### RODATA

Следующая область — **read-only data** (RODATA (R)). После запуска процесса эта область трактуется ОС как read-only, и её содержимое **не меняется** во время работы программы.

> **Иллюстрация 270.** Layout памяти Go бинарника с подсвеченной RODATA _(Рисунок не загружен)_

Read-only data хранит **неизменяемую информацию**, необходимую программе в runtime: type metadata и различные constant blob, используемые runtime и вашим кодом. Поскольку memory mapping read-only, программа может безопасно **указывать** на эти данные напрямую. Когда ей нужна writable копия (например, при конвертации `string` в `[]byte`), Go **копирует** байты в writable память.

Если запустить `go tool nm`, вы обычно **не увидите** каждый string literal как отдельный символ. Это потому, что Go linker не создаёт отдельный экспортированный символ для каждой строковой константы. Вместо этого он группирует строковые данные и выставляет их под синтетическим символом вроде `go:string.*`, представляющим область строковых данных **целиком**.

---

### Global Data

Область global data идёт следующей. Здесь живут **package-level переменные**. Она отображена как read-write, потому что глобальные переменные могут обновляться во время работы программы.

На высоком уровне она имеет две знакомые части: **DATA (RW)** для глобальных переменных с явным начальным значением и **BSS (RW)** для глобальных переменных, начинающихся с нуля. Область BSS **не нуждается** в хранении байтов в файле; ОС инициализирует её нулями при создании процесса:

> **Иллюстрация 271.** Layout global data и BSS памяти Go бинарника _(Рисунок не загружен)_

Go добавляет ещё одно отличие, которое большинству языков не нужно: **содержит ли** часть global data pointer. Garbage collector Go должен находить pointer, чтобы сохранять объекты, на которые они ссылаются, но сканирование больших блоков памяти слово за словом было бы **расточительным**, когда данные заведомо не содержат pointer.

Чтобы избежать этой работы, linker разделяет глобальные переменные на **pointer-free** и **pointer-carrying** области:

- Pointer-free глобальные переменные идут в `.noptrdata` и `.noptrbss` — garbage collector может **игнорировать** эти области.
- Глобальные переменные, которые могут содержать pointer, остаются в обычных `.data` и `.bss` — garbage collector **сканирует** их, используя metadata, сгенерированные linker.

---

### Heap & Main Thread Stack

На Linux Go-программа имеет **две области**, ведущие себя как heap: традиционный **brk-based heap** и **memory-mapped region**.

> **Примечание.** Memory-mapped region не является частью heap в традиционном смысле, но всё равно используется для динамических allocation. Мы упростим и будем трактовать его как часть Go heap.

**brk-based heap** располагается над `.bss` и растёт непрерывно при использовании C-библиотек вроде `malloc`. Сам Go **не выделяет** из этой области; вместо этого runtime растит heap, используя **mmap**, резервируя большие куски виртуального адресного пространства и коммитя физическую память по требованию. Во многих современных системах user-space allocator тоже предпочитают mmap для больших allocation, поэтому brk-based heap может оставаться маленьким.

Далее, говоря «heap» обобщённо, мы имеем в виду **mmap-backed region**, который Go использует для heap allocation.

Go heap — это часть этой области, управляемая garbage collector. Здесь живут **heap-allocated Go объекты**, такие как значения, сбежавшие со stack. Она **не включает** вещи вроде OS thread stack или памяти, выделенной C-библиотеками через `malloc`:

> **Иллюстрация 272.** Go runtime heap и layout main thread stack _(Рисунок не загружен)_

Stack в верхней части диаграммы — это **OS thread stack**, а не goroutine stack. Когда Go-программа запускается, ОС предоставляет начальному потоку (main thread) его собственный stack. Go **не выделяет** и не перемещает этот начальный OS stack. Goroutine работают на **отдельных stack**, которые runtime выделяет и наращивает по мере необходимости, используя собственные memory mapping.

---

### Эксперимент: адреса переменных

Проведём небольшой эксперимент:

```go
var global1, global2 = 1, 2
var global3, global4 int

func main() {
    stack1, stack2 := 3, 4
    heap1, heap2 := escaped(), escaped()

    println("initialized global1:", &global1, uintptr(unsafe.Pointer(&global1)))
    println("initialized global2:", &global2, uintptr(unsafe.Pointer(&global2)))
    println("uninitialized global3:", &global3, uintptr(unsafe.Pointer(&global3)))
    println("uninitialized global4:", &global4, uintptr(unsafe.Pointer(&global4)))
    println("stack 1:", &stack1, uintptr(unsafe.Pointer(&stack1)))
    println("stack 2:", &stack2, uintptr(unsafe.Pointer(&stack2)))
    println("heap 1:", heap1, uintptr(unsafe.Pointer(heap1)))
    println("heap 2:", heap2, uintptr(unsafe.Pointer(heap2)))
}

//go:noinline
func escaped() *int {
    c := 100
    return &c
}
```

Мы объявляем четыре глобальные переменные: первые две (`global1`, `global2`) инициализированы значениями, последние две (`global3`, `global4`) — нет. Пример одного запуска:

```
initialized global1:    0x102b0c3d0
initialized global2:    0x102b0c3d8
uninitialized global3:  0x102b35b28
uninitialized global4:  0x102b35b30
stack 1:                0x14000060720
stack 2:                0x14000060718
heap 1:                 0x1400000e090
heap 2:                 0x1400000e098
```

Глобальные переменные живут в **нижней части** виртуальной памяти — их адреса относительно малы по сравнению со stack и heap объектами. Заметьте, что `global2` на **8 байт** дальше `global1`, и `global4` на 8 байт дальше `global3`. Каждая пара хранится **непрерывно**, что ожидаемо. Но есть заметный промежуток **около 166 KiB** между `global2` и `global3`:

> **Иллюстрация 273.** Размещение инициализированных и неинициализированных глобальных переменных _(Рисунок не загружен)_

Этот промежуток обусловлен **разделением** между DATA segment и BSS segment — переменные хранятся в разных секциях в зависимости от того, инициализированы они или нет.

Теперь рассмотрим stack и heap переменные:

```
stack 1: 0x14000060720
stack 2: 0x14000060718
heap 1:  0x1400000e090
heap 2:  0x1400000e098
```

Адреса этих переменных **значительно выше** предыдущих. `stack1` и `stack2` — локальные переменные на stack основной goroutine, `heap1` и `heap2` указывают на heap-allocated integer, возвращённые `escaped()`. Базовые идеи сохраняются: `stack2` всего на 8 байт от `stack1`, и `heap2` на 8 байт от `heap1`.

Каждая пара довольно **близко друг к другу**, хотя это не гарантировано — зависит от того, как компилятор размещает stack frame и как runtime размещает heap-объекты:

> **Иллюстрация 274.** Сравнение расположения global, heap и stack переменных _(Рисунок не загружен)_

При многократном запуске адреса goroutine stack могут оказаться **ниже** адресов heap:

```
stack 1: 0x14000060720
stack 2: 0x14000060718
heap 1:  0x14000182000
heap 2:  0x14000182008
```

Важный вывод: **goroutine stack — это не то же самое, что OS thread stack**. Go-код выполняется на goroutine stack, которые runtime выделяет в собственных memory mapping, и эти stack могут появляться в **другой части** адресного пространства, чем OS-предоставленный stack.

Размещение user-level stack в runtime-managed памяти не уникально для Go; runtime вроде **Erlang**, Kotlin coroutines и различные fiber-библиотеки делают то же самое. Stack goroutine всё ещё ведёт себя как типичный thread stack во время выполнения кода, включая **рост вниз**.

---

## Итоговая таблица регионов памяти

| Регион                     | Разрешения      | Описание                                                                                                                                                                                     |
| -------------------------- | --------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **TEXT**                   | RX              | Скомпилированный машинный код: ваши функции + Go runtime. На non-PIE Linux/amd64 часто по фиксированной базе ~`0x400000`; с PIE + ASLR база варьируется.                                     |
| **RODATA**                 | R               | Read-only данные: константы, string literal, type metadata и т.д.                                                                                                                            |
| **DATA**                   | RW              | Инициализированные глобальные переменные программы и runtime.                                                                                                                                |
| **BSS**                    | RW              | Zero-initialized глобальные переменные. Файл не хранит их байты; loader резервирует место, ОС предоставляет обнулённую память.                                                               |
| **Traditional heap (brk)** | RW              | Непрерывная область, растущая через `brk(2)`/`sbrk(2)`, используемая C-библиотеками. Go runtime не использует для Go heap (использует mmap).                                                 |
| **mmap region**            | mixed           | Всё, отображённое через `mmap(2)`: shared library, GC-managed heap arena Go, goroutine stack, file-backed mapping.                                                                           |
| **OS thread stack**        | RW, растёт вниз | Stack, предоставляемый ОС потоку. Main thread начинает с OS-предоставленного stack; Go запускает большинство кода на goroutine stack, runtime и cgo используют per-thread system stack (g0). |

---

## Xulosa (Mening xulosalarim)

Bu bob Go binarnikining **ichki tuzilishi** va virtual memory layout'ini — ELF section'laridan tortib goroutine stack'ning heap bilan munosabatlarigacha ochib berdi. Mana asosiy xulosalarim:

**1. Section header vs Program header — ikki xil «ko'rish», ikki xil maqsad.** Bu farq ko'p dasturchilarni chalg'itadi, lekin tushunish muhim. 
**Section header** — toollar uchun (linker, debugger, `go tool nm`, `readelf -S`): «bu faylda qanday bo'limlar bor va ular qayerda». 
**Program header** — faqat yadro (kernel) uchun (`readelf -l`): «shu baytlarni virtual memory'ning shu adresiga yukla, shu permission'lar bilan». Runtime'da **faqat** `PT_LOAD` segment'lar ahamiyatli. `PT_PHDR` va `PT_NOTE` metadata'si biror `PT_LOAD` diapazoni ichiga tushsa **bilvosita** yuklanadi, aks holda — yo'q.

**2. Go binarnikida boshqa tillardan ko'proq section bor — va bu sababi bor.**

- `.typelink` — runtime'da interface satisfaction tekshiruvi uchun (Chapter 4 dagi itab bilan bog'liq)
- `.itablink` — itab entry'larni saqlaydi
- `.gosymtab` + `.gopclntab` — **Go-spetsifik** PC→source line mapping
- `.gopclntab` **memory'ga yuklanadi** (A flag) — shuning uchun panic stack trace'da `file:line` ko'rasiz
- `.debug_*` sektsiyalar **yuklanmaydi** (flag yo'q) — faqat pprof, delve, GDB kabi tashqi toollar uchun

Amaliy qoida: `-w` flag bilan build qilsangiz DWARF yo'qoladi (binary kichikroq), lekin Go stack trace **ishlaydi** (`.gopclntab` ga bog'liq). `-s` flag bilan symbol table ham yo'qoladi — debug qilish **butunlay** qiyinlashadi.

**3. `.noptrdata` / `.noptrbss` — GC uchun «ko'rinmas» zonalar.** 
Go linker global o'zgaruvchilarni **ikki guruhga** ajratadi:

- Pointer-free (`int`, `float64`, `[10]byte` va h.k.) → `.noptrdata` / `.noptrbss` — GC bu zonalarni **butunlay o'tkazib yuboradi**
- Pointer-carrying (`*int`, `[]byte`, `map`, `string` va h.k.) → `.data` / `.bss` — GC skanerlaydi

Bu **katta** optimizatsiya. Agar global o'zgaruvchilaringiz ko'p bo'lsa va ularda pointer yo'q bo'lsa — GC pressure **kamayadi**. Shuning uchun `[1000000]int` global array GC ga hech qanday yuk qo'shmaydi.

**4. PIE vs non-PIE — xavfsizlik va debugging o'rtasidagi muvozanat.**

- **Non-PIE** (Linux/amd64 default): dastur har safar **bir xil** adresga yuklanadi → `go tool nm` ko'rsatgan adres = runtime adres → debug qilish oson, lekin **hujumchilarga** funksiya adresini topish oson
- **PIE** (macOS/arm64 default): ASLR har safar boshqa bazaviy adres tanlaydi → funksiya **bir xil offset'da**, lekin absolyut adres har safar farq qiladi → **xavfsizroq**, lekin debugging biroz qiyinroq

Go'da `go build -buildmode=pie` bilan PIE'ni majburlash mumkin, yoki `-buildmode=exe` bilan non-PIE.

**5. Goroutine stack ≠ OS thread stack — bu eng ko'p chalkashtiradigan nuance.** Eksperiment natijasi buni aniq ko'rsatdi:

- Global o'zgaruvchilar: `0x102b0c...` (past, binary segment ichida)
- Stack o'zgaruvchilar: `0x140000607...` (yuqori, mmap region)
- Heap o'zgaruvchilar: `0x1400000e...` (yuqori, mmap region, stack bilan **bir xil zona**)

Stack va heap **ikkalasi ham** mmap region'ida, shuning uchun ba'zan stack adresi heap adresidan **past** bo'ladi. Bu Go'ga xos xususiyat — goroutine stack'lar runtime tomonidan `mmap` bilan ajratiladi, OS thread stack bilan **bog'liq emas**.

OS thread stack faqat **main thread** va runtime/cgo uchun ishlatiladi. Barcha Go kodingiz goroutine stack'da ishlaydi — bu stack kerak bo'lganda **o'sadi** va hatto **siljishi** mumkin (stack copying mechanism).

**6. DATA vs BSS — fayl hajmi optimizatsiyasi.**

```go
var x = 42    // → .data: faylda 8 bayt saqlaydi
var y int     // → .bss: faylda HECH NARSA saqlamaydi
```

BSS sektsiyasi faqat «shu qadar joy ajrat» deydi, yadro 0 bilan to'ldiradi. Shuning uchun:

- `var bigArray [1000000]int` — faylga **0 bayt** qo'shadi (BSS'da)
- `var bigArray = [1000000]int{1}` — faylga **~8 MB** qo'shadi (DATA'da)

Bu binary hajmi optimizatsiyasi uchun muhim — imkon qadar zero-initialized qoldiring.

**7. `go tool nm` — binary'ni tushunish uchun birinchi qadam.**

```bash
go tool nm ./main | grep main.add    # funksiya adresi
go tool nm ./main | grep ' B '       # BSS section'dagi symbollar
go tool nm ./main | grep ' D '       # DATA section'dagi symbollar
go tool nm ./main | grep ' T '       # TEXT section'dagi symbollar
```

**Amaliy maslahatlar:**

- `readelf -S` (Linux) yoki `otool -l` (macOS) bilan binarnikingiz section'larini tekshiring
- `readelf -l` bilan program header'larni ko'ring — qaysi section'lar memory'ga yuklanishini tushuning
- `-w` flag bilan DWARF'ni olib tashlang — production binary hajmini **sezilarli** kamaytiradi
- `-ldflags="-s -w"` — symbol table ham, DWARF ham yo'q → eng kichik binary
- Global o'zgaruvchilaringizda pointer bo'lmasa — `.noptrdata`/`.noptrbss`'da joylashadi, GC uchun **ko'rinmas**
- Goroutine stack va heap adreslari o'rtasidagi farqqa **tayanmang** — ikkalasi ham mmap region'ida
- Binary hajmi katta bo'lsa — `go tool nm` bilan qaysi symbollar joy egallashini tekshiring