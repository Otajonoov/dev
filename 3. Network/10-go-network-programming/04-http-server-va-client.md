# 04. HTTP server va client — production darajasidagi `net/http`

## Muammo / Hook

TCP echo server yozdik — lekin real dunyoda hech kim "xom" TCP ustiga o'z protokolini yozmaydi. Web API, mikroservis, brauzer — hammasi **HTTP** ustida ishlaydi. Va Go'ning eng kuchli tomonlaridan biri — `net/http` paketi. Bir necha qatorda production'ga chidamli web server ko'tarasan.

Ammo "ishlaydi" bilan "production'ga tayyor" orasida katta farq bor. Standart `http.ListenAndServe(":8080", nil)` — o'quv uchun yaxshi, lekin **timeout'lari yo'q**: bitta yovuz client ulanib, hech narsa yubormasa, server resursini abadiy ushlab turadi (Slowloris hujumi). Bu darsda ikkalasini ham — **to'g'ri sozlangan server** va **to'g'ri client** — o'rganamiz.

> `http.ListenAndServe(":8080", nil)` — demo uchun. Production'da har doim `&http.Server{...}` timeout'lar bilan.

## Analogiya — restoran

HTTP server'ni **restoran** deb tasavvur qil:

- **`http.Server`** — restoranning o'zi (bino, ish vaqti, qoidalar).
- **`ServeMux`** (router) — administrator: qaysi buyurtma qaysi oshpazga borishini hal qiladi (URL -> handler).
- **`Handler`** — oshpaz: bitta turdagi taomni tayyorlaydi (bitta endpoint'ga javob beradi).
- **`Middleware`** — kirishdagi qorovul + kassir: har mijoz oshpazga yetguncha ular orqali o'tadi (logging, auth, ...).
- **Timeout'lar** — "bir mijoz stolda qancha o'tira oladi" qoidasi: cheksiz o'tirsa, boshqalarga joy qolmaydi.

Analogiya chegarasi: restoranda bitta oshpaz bir vaqtda bitta taom qiladi; Go'da har so'rov **alohida goroutine**da, shuning uchun minglab mijoz parallel xizmatlanadi.

## Sodda ta'rif

> **`net/http`** — Go standart kutubxonasining HTTP server va client uchun paketi: `http.Server` so'rovlarni qabul qiladi, `ServeMux` ularni yo'naltiradi, `http.Handler` javob beradi, `http.Client` esa so'rov yuboradi.

## Diagramma — so'rovning yo'li

```mermaid
flowchart LR
    Req["HTTP so'rov<br/>GET /users/42"] --> S["http.Server<br/>(timeout'lar)"]
    S --> M1["Middleware:<br/>logging"]
    M1 --> M2["Middleware:<br/>auth"]
    M2 --> Mux["ServeMux<br/>(routing)"]
    Mux --> H["Handler<br/>getUser(42)"]
    H --> Resp["HTTP javob<br/>200 OK + JSON"]
```

## Worked example 1 — production server (mux + handler)

Go 1.22+ da `ServeMux` **method va path parametr** bilan routing qila oladi — endi tashqi router kutubxona kerak emas.

```go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	// --- 1-qadam: o'z ServeMux'imizni yaratamiz (global default emas) ---
	mux := http.NewServeMux()

	// --- 2-qadam: Go 1.22+ uslubida marshrutlar (method + path param) ---
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	mux.HandleFunc("GET /users/{id}", getUser)

	// --- 3-qadam: serverni TIMEOUT'lar bilan sozlaymiz ---
	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
	}

	log.Println("Server 8080-portda ishlamoqda")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server xatosi: %v", err)
	}
}

func getUser(w http.ResponseWriter, r *http.Request) {
	// --- path parametrni Go 1.22+ uslubida olamiz ---
	id := r.PathValue("id")
	user := User{ID: id, Name: "Ali"}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "kodlash xatosi", http.StatusInternalServerError)
	}
}
```

Bloklarni tushuntiramiz:

- **1-qadam** — `http.NewServeMux()`. Nega global `http.DefaultServeMux` emas? Chunki u **umumiy** (global) — boshqa import qilingan paket ham unga handler qo'shishi mumkin, bu esa xavfsizlik va debugging muammosi. O'zingniki bo'lsin.
- **2-qadam** — `"GET /users/{id}"` — Go 1.22+ sintaksisi. Method (`GET`) va path parametr (`{id}`) to'g'ridan-to'g'ri qo'llab-quvvatlanadi. Ilgari buning uchun gorilla/mux kabi tashqi kutubxona kerak edi.
- **3-qadam** — **eng muhim qism**: timeout'lar. Ularsiz server DoS'ga ochiq.

### Timeout'lar — nega har biri kerak

| Timeout | Nimani cheklaydi | Bo'lmasa nima bo'ladi |
| --- | --- | --- |
| `ReadHeaderTimeout` | Header'ni o'qish vaqti | Slowloris: sekin header yuborib ulanish ushlaydi |
| `ReadTimeout` | Butun so'rovni o'qish vaqti | Sekin body yuborib goroutine ushlaydi |
| `WriteTimeout` | Javob yozish vaqti | Sekin o'qiydigan client goroutine'ni ushlaydi |
| `IdleTimeout` | Keep-alive'da bo'sh turish | Ochiq, ishlamaydigan ulanishlar to'planadi |

**Notional machine:** har HTTP so'rov Go'da **alohida goroutine**da ishlanadi. Timeout'siz yovuz client minglab ulanish ochib, hech narsa yubormay tursa — minglab goroutine "abadiy kutish"da qoladi, xotira va file descriptor tugaydi. Timeout'lar bu goroutine'larni majburan tugatadi.

## Worked example 2 — middleware pattern

**Middleware** — bu `http.Handler`ni "o'raydigan" funksiya: asosiy handler'dan **oldin** yoki **keyin** kod ishga tushiradi (logging, auth, CORS, ...).

```mermaid
flowchart LR
    In["So'rov"] --> LM["Logging<br/>(oldin: vaqt boshla)"]
    LM --> AM["Auth<br/>(token tekshir)"]
    AM --> H["Asosiy handler"]
    H --> AM2["Auth (keyin)"]
    AM2 --> LM2["Logging<br/>(keyin: davomiylik log)"]
    LM2 --> Out["Javob"]
```

```go
// Middleware — Handler'ni qabul qilib, o'ralgan Handler qaytaradigan tip
type Middleware func(http.Handler) http.Handler

// --- Logging middleware: har so'rov davomiyligini yozadi ---
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r) // asosiy handler'ni chaqiramiz
		log.Printf("%s %s -> %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// --- Auth middleware: token bo'lmasa 401 qaytaradi ---
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			http.Error(w, "ruxsat yo'q", http.StatusUnauthorized)
			return // asosiy handler'ga BORMAYMIZ
		}
		next.ServeHTTP(w, r)
	})
}

// --- Chain: bir nechta middleware'ni ketma-ket ulaydi ---
func Chain(h http.Handler, mws ...Middleware) http.Handler {
	// Teskari tartibda o'raymiz: birinchi yozilgan tashqarida bo'ladi
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}
```

Ishlatish:

```go
protected := Chain(mux, Logging, Auth)
srv := &http.Server{Addr: ":8080", Handler: protected, /* timeout'lar */}
```

Bloklarni tushuntiramiz:

- **`Middleware`** tipi — `func(http.Handler) http.Handler`. Handler oladi, "o'ralgan" handler qaytaradi.
- **`next.ServeHTTP(w, r)`** — bu asosiy handler'ni (yoki keyingi middleware'ni) chaqirish. Undan **oldingi** kod so'rovdan oldin, **keyingi** kod javobdan keyin ishlaydi.
- **`Auth`da `return`** — token yo'q bo'lsa `next`ni **chaqirmaymiz**, so'rov shu yerda to'xtaydi (401). Bu middleware'ning "qorovul" roli.
- **`Chain`** — teskari tartibda o'raydi, shuning uchun `Chain(mux, Logging, Auth)` da so'rov avval Logging'dan, keyin Auth'dan o'tadi.

## PRIMM — bashorat qil

> 🤔 **O'ylab ko'r:** `Auth` middleware'da token yo'q bo'lganda `http.Error(...)`dan keyin `return` yozishni **unutdik**. Nima bo'ladi?

<details>
<summary>💡 Javobni ko'rish</summary>

Ikkita muammo yuzaga keladi:
1. `http.Error` 401 status va xato matnini yozadi, **lekin** `return` yo'qligi uchun kod davom etib `next.ServeHTTP(w, r)`ni ham chaqiradi — ya'ni himoyalanmagan handler ham ishga tushadi. **Auth butunlay buziladi**: ruxsatsiz so'rov ham asosiy handler'ga o'tadi.
2. Ikki marta yozish urinishi: `http.Error` allaqachon status yozgan, keyin handler yana yozmoqchi bo'ladi -> "superfluous WriteHeader call" ogohlantirishi.

Xulosa: middleware'da so'rovni to'xtatmoqchi bo'lsang, javob yozgandan keyin **doim `return`** qil.
</details>

## Worked example 3 — to'g'ri HTTP client

Client tomonda ham `http.Get(url)` (default client) production uchun **yaramaydi** — timeout'i yo'q. To'g'ri client — transport bilan sozlangan:

```go
// --- 1-qadam: transport'ni connection pooling bilan sozlaymiz ---
transport := &http.Transport{
	MaxIdleConns:        100,              // umumiy bo'sh ulanishlar
	MaxIdleConnsPerHost: 10,               // har host uchun
	IdleConnTimeout:     90 * time.Second, // bo'sh ulanish qancha yashaydi
	DialContext: (&net.Dialer{
		Timeout:   5 * time.Second, // ulanish o'rnatish timeout'i
		KeepAlive: 30 * time.Second,
	}).DialContext,
}

// --- 2-qadam: client'ga umumiy timeout beramiz ---
client := &http.Client{
	Transport: transport,
	Timeout:   15 * time.Second, // butun so'rov+javob uchun
}

// --- 3-qadam: context bilan so'rov (bekor qilish mumkin) ---
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

req, err := http.NewRequestWithContext(ctx, http.MethodGet,
	"https://api.example.com/users/42", nil)
if err != nil {
	log.Fatalf("so'rov yaratish xatosi: %v", err)
}

resp, err := client.Do(req)
if err != nil {
	log.Fatalf("so'rov xatosi: %v", err)
}
// --- 4-qadam: MUHIM — body'ni yop va TO'LIQ o'qi (connection reuse uchun) ---
defer resp.Body.Close()

body, err := io.ReadAll(resp.Body)
if err != nil {
	log.Fatalf("body o'qish xatosi: %v", err)
}
fmt.Printf("Status: %d, Body: %s\n", resp.StatusCode, body)
```

Eng ko'p qilinadigan xato — **4-qadam**. Go HTTP client ulanishni **faqat** body to'liq o'qilib, yopilgach qayta ishlatadi (connection reuse). Agar body'ni o'qimay tashlab ketsang, Go har so'rovga **yangi** TCP ulanish ochadi — bu file descriptor'larni tugatadi va serverni sekinlashtiradi.

> Qoida: `defer resp.Body.Close()` **va** body'ni to'liq o'qi (`io.ReadAll` yoki `io.Copy(io.Discard, resp.Body)`). Ikkalasi ham kerak.

## Worked example 4 — HTTP server'da graceful shutdown

2-darsdagi TCP graceful shutdown'ni `http.Server` o'zining `Shutdown` metodi bilan ancha soddalashtiradi:

```go
func main() {
	srv := &http.Server{Addr: ":8080", Handler: mux /* + timeout'lar */}

	// --- 1-qadam: serverni alohida goroutine'da ishga tushiramiz ---
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()
	log.Println("Server ishga tushdi")

	// --- 2-qadam: SIGINT/SIGTERM'ni kutamiz ---
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	// --- 3-qadam: graceful shutdown (30 soniya berib) ---
	log.Println("Shutdown boshlandi...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Shutdown xatosi: %v", err)
	}
	log.Println("Server toza yopildi")
}
```

`srv.Shutdown(ctx)` — yangi ulanishlarni to'xtatadi, mavjud so'rovlarni tugatishga imkon beradi, `ctx` muddati tugasa majburan yopadi. `http.ErrServerClosed` — bu **normal** xato (`Shutdown` chaqirilgani sababli), shuning uchun uni `Fatal` qilmaymiz.

## Ko'p uchraydigan xatolar

⚠️ **Xato 1 — timeout'siz `http.ListenAndServe`.**
Noto'g'ri tasavvur: "ishlayapti, demak tayyor." To'g'risi: doim `&http.Server{...}` da `ReadHeaderTimeout`, `ReadTimeout`, `WriteTimeout`, `IdleTimeout` sozla.

⚠️ **Xato 2 — client body'ni yopmaslik yoki o'qmaslik.**
`resp.Body`ni yopmasang goroutine/ulanish leak; o'qmasang connection reuse ishlamaydi. To'g'risi: `defer resp.Body.Close()` + body'ni to'liq o'qi.

⚠️ **Xato 3 — global `http.DefaultServeMux` va `http.DefaultClient`.**
Ular umumiy — boshqa paket ularga qo'shilishi mumkin va timeout'lari yo'q. To'g'risi: `http.NewServeMux()` va o'z `&http.Client{}`ing.

⚠️ **Xato 4 — middleware'da javob yozib, `return` unutish.**
So'rovni to'xtatmoqchi bo'lganda `return`siz kod davom etadi va asosiy handler ham ishlaydi. To'g'risi: javob yozgach doim `return`.

## Xulosa

- Production HTTP server = `&http.Server{}` + **to'rt xil timeout** (Slowloris va resurs sizishidan himoya).
- Go 1.22+ `ServeMux` method va path parametr routing'ni tashqi kutubxonasiz beradi (`"GET /users/{id}"`, `r.PathValue("id")`).
- **Middleware** = `func(http.Handler) http.Handler`; `next.ServeHTTP` bilan zanjir; to'xtatishda `return` shart.
- To'g'ri client = sozlangan `http.Transport` (connection pool) + umumiy `Timeout` + context; body'ni **doim yop va o'qi**.
- Graceful shutdown = `srv.Shutdown(ctx)`; `http.ErrServerClosed` — normal holat.
- Har so'rov alohida goroutine'da ishlaydi — shuning uchun timeout'lar hayotiy muhim.

## 🧠 Eslab qol

- `ListenAndServe(":8080", nil)` = demo; production'da `&http.Server{}` + timeout'lar.
- Go 1.22+ da router o'z ichida: `"GET /path/{id}"` + `r.PathValue`.
- Middleware = handler'ni o'raydi; to'xtatishda `return` unutma.
- Client body'ni doim yop VA to'liq o'qi (connection reuse).
- Graceful shutdown = `srv.Shutdown(ctx)`.

## ✅ O'z-o'zini tekshir (retrieval practice)

**1.** `http.ListenAndServe(":8080", nil)` production'da nega xavfli? Aniq bir hujum turini ayt.

<details>
<summary>Javob</summary>

Uning **timeout'lari yo'q**. **Slowloris** hujumida yovuz client ulanib, header'ni juda sekin (bir baytdan) yuboradi va hech qachon tugatmaydi. Timeout bo'lmagani uchun server goroutine'i abadiy kutadi. Minglab shunday ulanish bilan server goroutine va file descriptor'lari tugaydi -> DoS. `ReadHeaderTimeout` buni oldini oladi.
</details>

**2.** HTTP client'da `resp.Body`ni yopdik, lekin o'qimay tashlab ketdik. Qanday muammo chiqadi?

<details>
<summary>Javob</summary>

Go HTTP client ulanishni **faqat** body to'liq o'qilgach qayta ishlatadi. O'qimay yopsang, ulanish qayta ishlatilmaydi va har so'rovga **yangi** TCP ulanish ochiladi. Yuk ostida bu ephemeral port'lar va file descriptor'larni tugatadi, serverni sekinlashtiradi. Yechim: `io.Copy(io.Discard, resp.Body)` bilan qolganini o'qib tashla.
</details>

**3.** `Chain(mux, Logging, Auth)` da so'rov qaysi tartibda o'tadi va nega `Chain` teskari siklda o'raydi?

<details>
<summary>Javob</summary>

So'rov avval **Logging**, keyin **Auth**, keyin asosiy handler orqali o'tadi. `Chain` teskari (oxiridan boshiga) o'raydi, chunki oxirgi o'ralgan middleware **tashqi** qavat bo'ladi. `Logging` birinchi berilgani uchun eng tashqarida bo'lishi kerak — shuning uchun u eng oxirida o'raladi.
</details>

**4.** `srv.Shutdown(ctx)` chaqirilganda hozir ishlayotgan so'rovlarga nima bo'ladi? Va `http.ErrServerClosed`ni nega `Fatal` qilmaymiz?

<details>
<summary>Javob</summary>

`Shutdown` yangi ulanishlarni rad qiladi, lekin **mavjud** so'rovlar `ctx` muddati ichida tugashiga imkon beradi. `ListenAndServe` esa `Shutdown` chaqirilganda `http.ErrServerClosed` qaytaradi — bu **normal** signal (xato emas), shuning uchun uni `Fatal` qilsak, dastur o'rinsiz "xato" bilan chiqib ketardi.
</details>

## 🛠 Amaliyot

**1. Oson (Modify).** `GET /users/{id}` handler'ga yana bitta path parametr qo'sh: `GET /users/{id}/posts/{postID}`. Ikkala parametrni ham `r.PathValue`bilan olib, JSON'da qaytar.

<details>
<summary>Hint</summary>

`mux.HandleFunc("GET /users/{id}/posts/{postID}", ...)`, ichida `id := r.PathValue("id"); postID := r.PathValue("postID")`.
</details>

**2. O'rta (faded example — TODO to'ldirish).** Quyidagi "rate limit" middleware skeletini to'ldir: har so'rovni sanaydi, limitdan oshsa 429 qaytaradi.

```go
func RateLimit(limit int) Middleware {
	var mu sync.Mutex
	count := 0
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: mu bilan qulflab count'ni oshir
			// TODO: agar count > limit bo'lsa 429 (StatusTooManyRequests) qaytar va RETURN
			// TODO: aks holda next.ServeHTTP(w, r)
		})
	}
}
```

<details>
<summary>Hint</summary>

`mu.Lock(); count++; current := count; mu.Unlock()`. Keyin `if current > limit { http.Error(w, "juda ko'p so'rov", http.StatusTooManyRequests); return }`. Oxirida `next.ServeHTTP(w, r)`.
</details>

**3. Qiyin (Make — noldan).** To'liq CRUD JSON API yoz: xotirada `map[string]User` saqla, `GET /users`, `GET /users/{id}`, `POST /users`, `DELETE /users/{id}` marshrutlarini qo'llab-quvvatla. `sync.RWMutex` bilan concurrency'ni himoyala, server'ga timeout'lar va graceful shutdown qo'sh.

<details>
<summary>Hint</summary>

`type Store struct { mu sync.RWMutex; users map[string]User }`. O'qishda `RLock`, yozishda `Lock`. `POST`da `json.NewDecoder(r.Body).Decode(&u)`. Server qismini Worked example 1 va 4'dan ol.
</details>

## 🔁 Takrorlash

- **Bog'liq darslar:** [01-net-package-asoslari.md](01-net-package-asoslari.md) (context va deadline shu yerda), [02-tcp-client-server.md](02-tcp-client-server.md) (graceful shutdown g'oyasi). Keyingi [05-websocket-chat.md](05-websocket-chat.md) HTTP ustida WebSocket'ni quradi.
- **Takrorlash jadvali:** "timeout'lar nega kerak", "middleware'da return", "body'ni o'qish" nuqtalariga **ertaga**, **3 kundan so'ng**, **1 haftadan so'ng** qaytib javob ber.
- **Feynman testi:** "Nega `http.Get(url)` production'da yaramaydi?" degan savolga do'stingga 3 jumlada javob ber. (Kalit: timeout yo'q + default client umumiy.)

## 📚 Manbalar

- [How to Build a Production-Ready HTTP Server in Go — OneUptime](https://oneuptime.com/blog/post/2026-02-20-go-http-server-production/view)
- [Mastering Go's net/http Package — Edgar Montano](https://www.edgarmontano.com/posts/go/go-net-http-guide/)
- [Mastering HTTP Clients in Go — DEV Community](https://dev.to/jones_charles_ad50858dbc0/mastering-http-clients-in-go-your-guide-to-the-nethttp-package-2d8b)
- [How to Implement Graceful Shutdown in Go — OneUptime](https://oneuptime.com/blog/post/2026-01-23-go-graceful-shutdown/view)
- [Efficient Use of net/http, net.Conn, and UDP — Go Optimization Guide](https://goperf.dev/02-networking/efficient-net-use/)
