# Go'da Channel'lar - Batafsil Tushuntirish

## Channel nima?

**Channel** - goroutine'lar orasidagi "quvur" yoki "kanal". Bir goroutine ma'lumot yuboradi, ikkinchisi qabul qiladi.

**Yaratish:**

```go
ch := make(chan int)           // Unbuffered channel
ch := make(chan int, 3)        // Buffered channel (3 ta element sig'adi)
```

**Asosiy operatsiyalar:**

```go
ch <- 42        // Yuborish (send)
x := <-ch       // Qabul qilish (receive)
<-ch            // Qabul qilish, lekin natijani ishlatmaslik
close(ch)       // Yopish
```

---

## 1. Unbuffered Channel - To'g'ridan-to'g'ri aloqa

### Qanday ishlaydi?

**Asosiy qoida:** Yuboruvchi va qabul qiluvchi bir vaqtda uchrashishi kerak.

go

```go
ch := make(chan int)

// Goroutine 1
ch <- 42        // KUTADI qabul qiluvchini

// Goroutine 2  
x := <-ch       // KUTADI yuboruvchini
```

**Hayotiy misol:** Ikkita odam qo'l berib gap berish - birovning qo'lini ushlamasangiz, ikkinchisi kutadi. Hech kim o'rtada kutib turmaydi.

### Happens-before garantiyasi

**Muhim tushuncha:** Unbuffered channel'da qabul qilish **avtomatik ravishda** yuborishdan keyin sodir bo'ladi.

go

```go
var done = make(chan struct{})

// Background goroutine
go func() {
    io.Copy(os.Stdout, conn)
    log.Println("done")
    done <- struct{}{}  // Signal yuborish
}()

// Main goroutine
mustCopy(conn, os.Stdin)
conn.Close()
<-done  // Signal kutish - bu yerga yetguncha yuqoridagi ishlar tugagan
```

**Natija:** `<-done` dan keyin biz 100% ishonchimiz komil - goroutine tugadi va barcha o'zgarishlar ko'rinadi.

---

## 2. Pipeline Pattern - Zanjir ishlov berish

**Pipeline** - bir nechta goroutine'larni ketma-ket ulash, har biri ma'lumotni qayta ishlab keyingisiga uzatadi.

### Oddiy misol: Sonlarni kvadratga ko'tarish

go

```go
func main() {
    naturals := make(chan int)
    squares := make(chan int)
    
    // 1-bosqich: Sonlar generatsiya qilish
    go func() {
        for x := 0; x < 100; x++ {
            naturals <- x
        }
        close(naturals)  // Tugadim deb signal
    }()
    
    // 2-bosqich: Kvadratga ko'tarish
    go func() {
        for x := range naturals {  // Yopilguncha o'qiydi
            squares <- x * x
        }
        close(squares)
    }()
    
    // 3-bosqich: Chiqarish
    for x := range squares {
        fmt.Println(x)
    }
}
```

### Channel yopish (close)

**Muhim qoidalar:**

1. **Close faqat yuboruvchi tomondan** - qabul qiluvchi yopmaydi
2. **Yopilgan channel'ga yuborish** - panic
3. **Yopilgan channel'dan o'qish:**
    - Buffer'dagi qolgan ma'lumotlar qaytadi
    - Buffer bo'shagach - zero value qaytadi (0, "", nil, false)

**Yopilganligini tekshirish:**

go

```go
x, ok := <-ch
if !ok {
    // Channel yopilgan va bo'sh
}
```

**Range bilan:**

go

```go
for x := range ch {  // Yopilguncha davom etadi
    fmt.Println(x)
}
```

### Qachon close qilish kerak?

- ✅ **Kerak:** Qabul qiluvchiga "boshqa ma'lumot yo'q" deb bildirish uchun
- ❌ **Kerak emas:** Har doim - garbage collector o'zi tozalaydi
- ⚠️ **Muhim:** Fayl'lardan farqli - channel'ni yopmaslik xatolik emas

---

## 3. Unidirectional Channels - Bir yo'nalishli kanallar

**Muammo:** Funksiya parametri sifatida channel o'tkazganda, uni noto'g'ri ishlatish mumkin.

**Yechim:** Type system orqali cheklash.

### Turlar:

go

```go
chan<- int      // Faqat yuborish (send-only)
<-chan int      // Faqat qabul qilish (receive-only)
chan int        // Ikkalasi ham (bi-directional)
```

### Pipeline misoli:

go

```go
func counter(out chan<- int) {  // Faqat yuboradi
    for x := 0; x < 100; x++ {
        out <- x
    }
    close(out)
}

func squarer(out chan<- int, in <-chan int) {  // in'dan o'qiydi, out'ga yozadi
    for v := range in {
        out <- v * v
    }
    close(out)
}

func printer(in <-chan int) {  // Faqat o'qiydi
    for v := range in {
        fmt.Println(v)
    }
}

func main() {
    naturals := make(chan int)
    squares := make(chan int)
    
    go counter(naturals)           // chan int -> chan<- int (avtomatik)
    go squarer(squares, naturals)  // chan int -> ikki yo'nalishli
    printer(squares)               // chan int -> <-chan int (avtomatik)
}
```

**Foydasi:**

- Compile-time xavfsizlik
- Funksiya niyati aniq
- Noto'g'ri ishlatish - kompilator xatosi

**Muhim:** Bi-directional → Unidirectional o'tish mumkin, lekin aksi yo'q!

---

## 4. Buffered Channels - Buffer bilan kanallar

### Farqi nima?

**Unbuffered:**

- Yuboruvchi va qabul qiluvchi bir vaqtda uchrashadi
- Har bir operatsiya sinxronlashadi

**Buffered:**

- Ma'lum miqdorda element saqlay oladi
- Buffer to'lmaguncha yuboruvchi kutmaydi
- Buffer bo'shmaguncha qabul qiluvchi kutmaydi

### Yaratish va ishlatish:

go

```go
ch := make(chan string, 3)  // 3 ta element sig'adi

// To'ldiramiz
ch <- "A"  // Kutmaydi
ch <- "B"  // Kutmaydi  
ch <- "C"  // Kutmaydi
ch <- "D"  // KUTADI - buffer to'lgan!

// O'qiymiz
fmt.Println(<-ch)  // "A" - endi joy ochildi
ch <- "D"          // Endi yuborish mumkin
```

### Hajmni bilish:

go

```go
fmt.Println(cap(ch))  // 3 - umumiy hajm
fmt.Println(len(ch))  // 2 - hozirgi element soni
```

⚠️ **Ogohlantirish:** `len()` parallel dasturlashda kam foydali - qiymat darhol eskiradi.

### Konveyer analogiyasi

**Konditerlar misoli:**

**Buffer yo'q (unbuffered):**

- 3 ta konditer: birinchisi tort pishiradi, ikkinchisi glazur quyadi, uchinchisi bezaydi
- Har biri keyingisi tayyor bo'lguncha kutadi
- Sekin - har bir qadam sinxronlashgan

**Buffer bor (buffered):**

- Konditerlar orasida stol bor (buffer)
- Birinchisi tortni stolga qo'yib, keyingisiga o'tadi
- Ikkinchisi tayyor bo'lganda oladi
- Tezroq - vaqt farqlari silliqlashadi

**Asosiy g'oya:** Buffer goroutine'larning tezlik farqlarini muvozanatlaydi.

### Qachon buffered channel kerak?

✅ **Kerak bo'lganda:**

1. Yuboriluvchi element sonini oldindan bilsangiz
2. Goroutine'lar tezligi o'zgaruvchan
3. Burst load'larni qabul qilish kerak

❌ **Kerak bo'lmaganda:**

1. Sinxronizatsiya muhim bo'lsa (unbuffered yaxshiroq)
2. Buffer hajmini noto'g'ri tanlasangiz (deadlock xavfi)

### mirroredQuery misoli - Eng tez javobni olish

go

```go
func mirroredQuery() string {
    responses := make(chan string, 3)  // 3 ta server
    
    // Har biriga parallel so'rov
    go func() { responses <- request("asia.gopl.io") }()
    go func() { responses <- request("europe.gopl.io") }()
    go func() { responses <- request("americas.gopl.io") }()
    
    return <-responses  // Birinchisini qaytarish
}
```

**Nega buffer kerak?**

- Unbuffered bo'lsa: sekin 2 ta goroutine bloklanadi - **goroutine leak**
- Buffered bilan: hammalari yozib tugatadi va tugatadi

### Buffer hajmini tanlash

**Optimal hajm topish:**

go

```go
// Juda kichik - goroutine'lar ko'p kutadi
ch := make(chan int, 1)

// To'g'ri - maksimal parallel vazifalar soni
ch := make(chan int, maxConcurrent)

// Juda katta - xotira isrofi, foydasiz
ch := make(chan int, 1000000)
```

**Qoida:** Buffer hajmi = Maksimal parallel bajariluvchi vazifalar soni

---

## Xulosa va Best Practices

### Unbuffered vs Buffered - qaysi birini tanlash?

|Vazifa|Tanlov|Sabab|
|---|---|---|
|Sinxronizatsiya kerak|Unbuffered|Kafolatlangan happens-before|
|Performance muhim|Buffered|Goroutine'lar kutmaydi|
|Element soni ma'lum|Buffered|Optimal hajm tanlanadi|
|Pipeline|Ikkalasi ham|Vazifaga bog'liq|

### Channel yopish qoidalari

1. ✅ Faqat yuboruvchi yopadi
2. ✅ Boshqa ma'lumot yo'qligini bildirish uchun yopiladi
3. ❌ Har bir channel'ni yopish shart emas
4. ❌ Yopilgan channel'ni qayta yopish - panic
5. ✅ `for range` bilan yopilish avtomatik aniqlanadi

### Xavfsizlik

**To'g'ri:**

go

```go
// Pipeline tugashini kutish
for x := range ch {
    // Process
}
```

**Noto'g'ri:**

go

```go
// Bitta goroutine'da buffer'ni queue sifatida ishlatish
for i := 0; i < 10; i++ {
    ch <- i  // Deadlock xavfi!
}
```

### Unidirectional channels afzalliklari

- **Type safety** - noto'g'ri ishlatish compile error
- **Documentation** - funksiya niyati aniq
- **Refactoring xavfsizligi** - o'zgarishlar kontrollangan

---

**Eng muhim prinsip:**

- Unbuffered channel - **sinxronizatsiya** uchun
- Buffered channel - **performance va decoupling** uchun
- Unidirectional channels - **xavfsizlik va aniqlik** uchun

Channel yopish - qabul qiluvchiga **signal yuborish** mexanizmi, majburiyat emas!


---
---
Буферизованный канал имеет очередь элементов. Максимальный размер очере­ ди определяется при создании канала с помощью аргумента емкости функции make.

ch := make(chan int) // Channel yaratish ch <- x // Yuborish (send) x = <-ch // Qabul qilish (receive) close(ch) // Yopish

## Unidirectional Channel - Bir yo'nalishli kanal

Xavfsizlik va aniqlik uchun channel'ni faqat yuborish yoki faqat qabul qilish uchun belgilash mumkin:

go

`func counter(out chan<- int) // Faqat yuborish func squarer(out chan<- int, in <-chan int) // in faqat qabul, out faqat yuborish func printer(in <-chan int) // Faqat qabul qilish`