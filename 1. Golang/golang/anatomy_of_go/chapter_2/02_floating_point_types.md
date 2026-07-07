# 2. Floating-Point Types: precision, representation va IEEE-754

Go ikkita floating-point type beradi:

- `float32` - single precision
- `float64` - double precision

Ikkalasi IEEE-754 standardiga amal qiladi.

```go
var a float32 = 0.1
var b float64 = 0.0001
```

Bit soni integer type'larga o'xshasa ham, float'lar ancha katta range'ni ifodalay oladi. Buning narxi - precision.

```text
float32: 1.4e-45 to 3.4e+38
float64: 4.9e-324 to 1.8e+308
```

Go 1'dan oldin architecture-dependent `float` type bo'lgan, lekin u type system'ni ortiqcha murakkab qilgani uchun 2011-yilda olib tashlangan. Go aniq `float32` va `float64` bilan qolgan.

Type explicit berilmasa, floating literal default `float64` bo'ladi:

```go
a := 0.1
fmt.Printf("Type of a is %T\n", a) // "Type of a is float64"
```

`float64` default sifatida tanlangan, chunki Go standard library'da floating calculation uchun aniqroq va kengroq precision beradi.

## 2.1 Floating-point literal formatlari

Float literal'lar uch asosiy formatda yoziladi:

- Decimal form
- Scientific notation
- Hexadecimal form

```go
// Decimal form
d1 := 0.1
d2 := 1.
d3 := .1

// Scientific notation
s1 := 1.2e3 // 1.2 * 10^3 = 1200
s2 := 1.E-3 // 1 * 10^-3 = 0.001
s3 := .2e3  // 0.2 * 10^3 = 200

// Hexadecimal form
h1 := 0x1.2p3 // (1 + 2/16) * 2^3 = 1.125 * 8 = 9
h2 := 0x1P3   // 1 * 2^3 = 8
h3 := 0X1.p3  // 1 * 2^3 = 8
```

Hexadecimal float exponent uchun `p`/`P` ishlatadi, chunki exponent base-2 bilan bog'liq. Hex literal ichida `e`/`E` scientific exponent sifatida ishlatilmaydi.

"Floating-point literal" nomi biroz chalg'itishi mumkin. Ba'zi formatlar integer assignment'da ham ishlaydi, agar qiymat integer type'ga exact sig'sa:

```go
var a int = 1e3     // 1000
var b float64 = 0x1 // 1
```

Lekin fractional qiymat integer'ga bevosita sig'maydi:

```go
// Error: cannot use 1e-1 (untyped float constant 0.1)
// as int value in variable declaration (truncated)
var a int = 1e-1

var b int = 1.2e3 // This works fine, equals 1200
```

Go literal'larni avval untyped constant sifatida ko'radi. Keyin context'ga qarab type beradi, faqat qiymat maqsad type uchun ma'noli bo'lsa.

## 2.2 Type conversion va precision trade-off

Float type'lar orasida conversion ham explicit:

```go
// Precision loss
var a float64 = 1.123456789
var b float32 = float32(a) // 1.1234568 (some precision is lost)
```

`float32 -> float64` odatda range va precision jihatdan xavfsizroq. `float64 -> float32` esa precision yo'qotadi.

Float'ni integer'ga convert qilish rounding qilmaydi, fractional qismni kesib tashlaydi:

```go
var a float64 = 1.999999
var b int = int(a) // 1 (decimal is dropped, not rounded)
```

Integer'dan float'ga o'tishda ham katta sonlarda precision yo'qolishi mumkin:

```go
var c int = 123456789
var d float32 = float32(c) // 123456792 (some rounding happens)
```

> **Caution:** pul hisob-kitoblarida float ishlatish ko'p uchraydigan xato. Dollar/sent ko'rinishida float qulay tuyuladi, lekin rounding error yig'ilib ketadi. Pul uchun integer (masalan, cent sifatida), fixed-point yoki decimal library ishlatish yaxshi.

## 2.3 IEEE-754 behavior: +Inf, -Inf va NaN

Go float'larni IEEE-754 bo'yicha saqlaydi. Shuning uchun float division by zero integer'dan farq qiladi:

- positive float / zero -> `+Inf`
- negative float / zero -> `-Inf`
- zero / zero -> `NaN`

```go
func main() {
    zero := 0.0

    println(1 / zero)    // +Inf
    println(-1 / zero)   // -Inf
    println(zero / zero) // NaN
}
```

`NaN` bilan har qanday operation yana `NaN` beradi:

```go
func main() {
    NaN := math.NaN()
    Inf := math.Inf(1)
    NegInf := math.Inf(-1)

    println(NaN - 100)    // NaN
    println(NaN + Inf)    // NaN
    println(NaN + NegInf) // NaN
}
```

Infinity bilan natija operation'ga bog'liq:

```go
println(Inf + Inf)       // +Inf
println(NegInf + NegInf) // -Inf
println(NegInf + Inf)    // NaN
```

`NaN` comparison'da ham alohida:

```go
println(NaN == NaN) // false
println(NaN != NaN) // true
```

`+Inf` finite sonlarning hammasidan katta, `-Inf` esa hammasidan kichik:

```go
println(Inf > 10)         // true
println(Inf == Inf)       // true
println(NegInf < 10)      // true
println(NegInf == NegInf) // true
```

`NaN` tekshirishning eng sodda formasi `n != n`, lekin readable code uchun `math` package ishlatgan yaxshi:

```go
n := math.NaN()
if n != n {
    fmt.Println("n is NaN")
}
```

```go
import "math"

math.IsNaN(n)
math.IsInf(n, 0) // 0: any infinity, -1: negative infinity, 1: positive infinity
```

## Rounding error va precision loss

Float'larda eng mashhur holat:

```go
a := 0.1
b := 0.2
c := 0.3

a + b == c // false
a + b      // 0.30000000000000004
c          // 0.3
```

Sabab: `0.1` va `0.3` kabi decimal sonlarni binary floating-point formatda exact ifodalash mumkin emas. Base-10 da `1/3` cheksiz `0.3333...` bo'lgani kabi, base-2 da ham ba'zi sonlar cheksiz representation talab qiladi.

Katta float32 sonlarda precision yanada yaqqol ko'rinadi:

```go
a := float32(16777216)
b := a + 1

a == b // true
b      // 16777216
```

`float32` butun sonlarni 16,777,216 gacha aniq ajrata oladi. Undan keyin har bir qo'shni integer alohida ko'rinmaydi.

| Number | Float32 |
|--------|---------|
| 16777216 | 16777216 |
| 16777217 | 16777216 |
| 16777218 | 16777218 |
| 16777219 | 16777220 |
| 16777220 | 16777220 |
| 16777221 | 16777220 |
| 16777222 | 16777222 |
| 66777310 | 66777312 |
| 66777311 | 66777312 |
| 66777312 | 66777312 |
| 66777313 | 66777312 |
| 66777314 | 66777312 |
| 66777315 | 66777316 |

`math.Nextafter32` keyingi representable float32 qiymatni ko'rsatadi:

```go
math.Nextafter32(66777310, float32(math.Inf(+1))) // 66777316
math.Nextafter32(66777312, float32(math.Inf(+1))) // 66777316
```

Float magnitude juda katta bo'lsa, overflow `+Inf` yoki `-Inf` ga olib boradi. Juda kichik magnitude esa zero tomonga underflow qiladi.

```go
max := float32(math.MaxFloat32) // 3.4028235e+38
max + 1                         // 3.4028235e+38
max + max                       // +Inf
```

`max + 1` o'zgarmaydi, chunki shu scale'da `1` juda kichik. `max + max` esa range'dan chiqib, `+Inf` bo'ladi.

## Eslab qol

- Go'da `float` yo'q; `float32` va `float64` bor.
- Default floating literal type - `float64`.
- Float division by zero panic emas, IEEE-754 special value beradi.
- `NaN` hech narsaga, hatto o'ziga ham teng emas.
- Float equality comparison'da ehtiyot bo'l; tolerance/epsilon yondashuvi ko'pincha kerak bo'ladi.
- Pul uchun float ishlatma.
