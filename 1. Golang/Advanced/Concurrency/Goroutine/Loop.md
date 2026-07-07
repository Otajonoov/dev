# Go'da Parallel Loop'lar va Advanced Channel Patternlar

## 1. Parallel Loop'lar - Mustaqil vazifalarni parallel bajarish

**Asosiy g'oya:** Agar loop'ning har bir iteratsiyasi boshqasiga bog'liq bo'lmasa, ularni parallel bajarsak tezroq bo'ladi.

**Misol: Rasm thumbnail'larini yaratish**

**Noto'g'ri yondashuv #1:**

go

```go
func makeThumbnails2(filenames []string) {
    for _, f := range filenames {
        go thumbnail.ImageFile(f)  // Parallel ishga tushadi
    }
    // XATO: Funksiya darhol qaytadi, goroutine'lar tugamaydi
}
```

**Muammo:** Funksiya goroutine'lar tugashini kutmaydi.

**To'g'ri yondashuv - Channel bilan:**

go

```go
func makeThumbnails3(filenames []string) {
    ch := make(chan struct{})
    for _, f := range filenames {
        go func(f string) {
            thumbnail.ImageFile(f)
            ch <- struct{}{}  // Signal yuborish
        }(f)
    }
    // Barcha goroutine'larni kutish
    for range filenames {
        <-ch
    }
}
```

**Muhim: Loop variable capture muammosi**

go

```go
// NOTO'G'RI:
for _, f := range filenames {
    go func() {
        thumbnail.ImageFile(f)  // f o'zgaradi!
    }()
}

// TO'G'RI:
for _, f := range filenames {
    go func(f string) {  // f parametr sifatida
        thumbnail.ImageFile(f)
    }(f)
}
```

## 2. WaitGroup - Goroutine'larni sanash

**Muammo:** Goroutine'lar soni noma'lum bo'lganda nima qilish kerak?

**Yechim: sync.WaitGroup**

go

```go
var wg sync.WaitGroup
for f := range filenames {
    wg.Add(1)  // Goroutine qo'shamiz deb signal
    go func(f string) {
        defer wg.Done()  // Tugagach kamaytirish
        thumbnail.ImageFile(f)
    }(f)
}
wg.Wait()  // Hammasini kutish
```

**Muhim qoidalar:**

- `Add()` goroutine ishga tushishidan **OLDIN** chaqiriladi
- `Done()` goroutine **ICHIDA** chaqiriladi (odatda defer bilan)
- `Wait()` alohida goroutine'da yoki asosiy codeda

## 3. Select - Bir nechta channel'dan kutish

**Select nima?** Bir nechta channel operatsiyasidan birinchi tayyor bo'lganini tanlaydi.

**Asosiy sintaksis:**

go

```go
select {
case msg := <-ch1:
    // ch1'dan xabar kelsa
case ch2 <- value:
    // ch2'ga yuborish mumkin bo'lsa
case <-time.After(1 * time.Second):
    // 1 sekunddan keyin
default:
    // Hech biri tayyor bo'lmasa (non-blocking)
}
```

**Real misol: Raketa ishga tushirish countdown'i**

go

```go
select {
case <-time.After(10 * time.Second):
    launch()  // 10 sekund o'tdi
case <-abort:
    fmt.Println("Bekor qilindi!")
    return  // Foydalanuvchi Enter bosdi
}
```

**Select xususiyatlari:**

- Bir nechta case tayyor bo'lsa - **tasodifiy** birini tanlaydi
- `default` bor bo'lsa - hech narsa tayyor bo'lmasa darhol ishlaydi (non-blocking)
- `nil` channel - hech qachon tayyor bo'lmaydi (case'ni "o'chirish" uchun ishlatiladi)

## 4. Cancellation Pattern - Bekor qilish mexanizmi

**Muammo:** Ko'plab goroutine'larni bir vaqtning o'zida to'xtatish kerak.

**Yechim: Channel yopish orqali broadcast**

go

```go
var done = make(chan struct{})

// Bekor qilish funksiyasi
func cancelled() bool {
    select {
    case <-done:
        return true
    default:
        return false
    }
}

// Bekor qilish
go func() {
    os.Stdin.Read(make([]byte, 1))  // Enter kutish
    close(done)  // Hammaga signal
}()

// Goroutine'da tekshirish
func worker() {
    for {
        if cancelled() {
            return
        }
        // Ish bajarish
    }
}
```

**Nega channel yopish?**

- Yopilgan channel'dan o'qish darhol zero value qaytaradi
- Bir marta yopsak, **barcha** goroutine'lar ko'radi (broadcast)
- Value yuborishdan farqli - value faqat bitta goroutine oladi

## 5. Semaphore Pattern - Parallel vazifalarni cheklash

**Muammo:** Juda ko'p goroutine ochilsa, resurslar tugaydi (masalan, fayllar).

**Yechim: Buffered channel semaphore sifatida**

go

```go
// Maksimum 20 ta parallel so'rov
var tokens = make(chan struct{}, 20)

func crawl(url string) []string {
    tokens <- struct{}{}  // Token olish (kutish mumkin)
    defer func() { <-tokens }()  // Token qaytarish
    
    return links.Extract(url)  // Aslida ish bajarish
}
```

**Qanday ishlaydi:**

- Channel buffer'i to'lsa - keyingi yuborish kutadi
- Biror goroutine tugasa - token ozod bo'ladi
- Natija: Maksimum N ta goroutine bir vaqtda ishlaydi

**Hayotiy misol:** Restoran eshigidagi qorovul - ichkarida joy bo'lsa kiritadi, to'lsa kutkazadi.

## 6. Chat Server - Real loyiha misoli

**Arxitektura:** 4 xil goroutine:

1. **main** - yangi client'larni qabul qiladi
2. **broadcaster** - xabarlarni tarqatadi
3. **handleConn** - har bir client uchun (o'qiydi)
4. **clientWriter** - har bir client uchun (yozadi)

**Asosiy channel'lar:**

go

```go
var (
    entering = make(chan client)  // Yangi client
    leaving = make(chan client)   // Ketayotgan client
    messages = make(chan string)  // Barcha xabarlar
)
```

**Broadcaster pattern:**

go

```go
func broadcaster() {
    clients := make(map[client]bool)
    for {
        select {
        case msg := <-messages:
            // Hammaga yuborish
            for cli := range clients {
                cli <- msg
            }
        case cli := <-entering:
            clients[cli] = true
        case cli := <-leaving:
            delete(clients, cli)
            close(cli)
        }
    }
}
```

**Nima uchun yaxshi?**

- `clients` map faqat broadcaster goroutine'da - **race condition yo'q**
- Channel'lar orqali aloqa - **lock'siz parallellık**
- Har bir client mustaqil - bitta buzilsa, boshqalari ishlaydi

---

## Asosiy Xulosalar

### 1. **Parallel loop'lar uchun:**

- Loop variable'ni parameter qilib o'tkaz
- WaitGroup yoki channel bilan tugashini kut
- Goroutine leak'dan ehtiyot bo'l

### 2. **Select operatori:**

- Bir nechta channel'dan birinchisi tayyor bo'lganini ol
- Timeout, cancellation uchun juda qulay
- `default` bilan non-blocking qilish mumkin

### 3. **Cancellation pattern:**

- Channel yopish - broadcast mexanizmi
- `cancelled()` funksiya bilan doimiy tekshirish
- Goroutine'lar tez javob berishi uchun muhim joylarni tekshir

### 4. **Semaphore pattern:**

- Buffered channel resurslarni cheklash uchun
- Token olish/qaytarish mexanizmi
- defer bilan token qaytarishni kafolatlash

### 5. **Arxitektura:**

- Channel orqali aloqa - lock'dan yaxshiroq
- Ma'lumotni bitta goroutine'da saqlash (confinement)
- Broadcaster pattern - markazlashtirilgan boshqaruv

**Eng muhim tamoyil:** Channel'lar orqali ma'lumot ulashing, xotirani emas. Go'ning mottosi: "Don't communicate by sharing memory, share memory by communicating."