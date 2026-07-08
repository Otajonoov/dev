
==Context== - bu Go dasturlash tilida 1.7 versiyasidan boshlab qo'shilgan standart paket bo'lib, u orqali:

- Operatsiyalarning maksimal ishlash muddati (deadlines)
- Bekor qilish signallari (cancellation)
- So'rovlar bo'yicha qiymatlarni uzatish

##### Asosiy tamoyil

Go tilidagi kontekstlarda quyidagi qoida mavjud:
> **"Context bekor qilinganda, undan olingan BARCHA contextlar ham bekor qilinadi.  
> Lekin bu context olingan ASL kontekstlar bekor qilinmaydi."

Buni quyidagi misol orqali tushunish mumkin:
```
// 1. Asosiy context yaratamiz (hech qachon avtomatik bekor bo'lmaydi)
rootCtx := context.Background()

// 2. rootCtx dan yangi context olamiz (cancel funksiyasi bilan)
childCtx, cancel := context.WithCancel(rootCtx)

// 3. childCtx dan yana bir context olamiz (timeout bilan)
grandchildCtx, _ := context.WithTimeout(childCtx, 10*time.Second)

// 4. Endi cancel() ni chaqiramiz
cancel()
```

1. **childCtx** bekor qilinadi (chunki biz `cancel()` ni chaqirdik)
2. **grandchildCtx** ham avtomatik bekor qilinadi (chunki u childCtx dan olingan)
3. **rootCtx** o'zgarmaydi (u hech qachon bekor qilinmaydi)


##### Contextning asosiy afzalliklari
1. **Bekor qilishni boshqarish**:
    - Agar foydalanuvchi so'rovni bekor qilsa, barcha bog'liq jarayonlar ham to'xtatiladi
    - Resurslarni tejash (CPU, xotira, tarmoq)
    
2. **Vaqt chegaralari**:
    - Operatsiyalar uchun maksimal ishlash vaqtini belgilash
    - Uzoq davom etadigan operatsiyalarni avtomatik to'xtatish
    
3. **Xavfsizlik**:
    - Context bir vaqtning o'zida bir nechta gorutinalarda ishlatilishi mumkin (thread-safe)


##### Context turlari va yaratish usullari: 
```
// Asosiy kontekstlar
ctx := context.Background()  // Boshlang'ich kontekst
ctx := context.TODO()        // Vaqtincha kontekst

// O'zgartirilgan kontekstlar
ctx, cancel := context.WithCancel(parentCtx)               // Bekor qilish bilan
ctx, cancel := context.WithTimeout(parentCtx, 10*time.Second) // time bilan
ctx, cancel := context.WithDeadline(parentCtx, time.Now().Add(10*time.Second)) // Aniq vaqt bilan
ctx := context.WithValue(parentCtx, key, value)      // Qiymat qo'shish bilan
```


==WithTimeout vs WithDeadline==

|Xususiyat|WithTimeout|WithDeadline|
|---|---|---|
|**Parametr**|`time.Duration` (10 soniya)|`time.Time` (aniq vaqt)|
|**Foydalanish**|Nisbiy vaqt (hozirdan boshlab)|Absolyut vaqt (aniq sana/vaqt)|
|**Ishlatish**|HTTP so'rovlari, umumiy timeout|Cron vazifalar, rejalashtirish|
|**Ichki logika**|`WithDeadline` ni chaqiradi|To'g'ridan-to'g'ri ishlaydi|

==WithTimeout==: 
```
// HTTP so'rov uchun 5 soniya timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// Bu so'rov 5 soniyadan ortiq davom etsa, avtomatik bekor bo'ladi
resp, err := http.NewRequestWithContext(ctx, "GET", "https://example.com", nil)
```

==WithDeadline:== 
```
// Aniq vaqtgacha ishlaydigan vazifa
deadline := time.Date(2023, 6, 20, 14, 30, 0, 0, time.UTC)
ctx, cancel := context.WithDeadline(context.Background(), deadline)
defer cancel()

// Bu vazifa 2023-06-20 14:30:00 da avtomatik to'xtaydi
doSomeWork(ctx)
```

Texnik jihatdan `WithTimeout` ichida `WithDeadline` ishlatadi, shuning uchun ikkalasi ham bir xil natijani beradi. Farq faqat foydalanish qulayligida.

##### Xulosa:
1. Har bir context o'zining hayot tsiklini mustaqil boshqaradi
2. Memory leaklardan oldini oladi (child contextlar bekor qilinishi kafolatlanadi)


Ядром пакета `context`является `Context`тип:
```go
// Контекст переносит крайний срок, сигнал отмены и значения в области запроса 
// через границы API. Его методы безопасны для одновременного использования несколькими 
// горутинами. 
type Context interface { 

    // Done возвращает канал, который закрывается при отмене этого контекста         // или истечении времени ожидания. 
    Done() <-chan struct{} 
    
    // Err указывает причину отмены этого контекста после закрытия канала Done. 
    Err() error 
    
    // Deadline возвращает время отмены этого контекста, если таковое имеется. 
    Deadline() (deadline time.Time, ok bool) 
    
    // Value возвращает значение, связанное с ключом, или nil, если его нет. 
    Value(key interface{}) interface{} 
}
```

Метод `Done`возвращает канал, который служит сигналом отмены для функций, запущенных от имени `Context`: при закрытии канала функции должны прекратить свою работу и вернуться. `Err`Метод возвращает ошибку, указывающую причину `Context`отмены.

_У_ A `Context`нет метода по той же причине, по которой канал предназначен только для приёма: функция, получающая сигнал отмены, обычно не является функцией, которая его отправляет. В частности, когда родительская операция запускает горутины для подопераций, эти подоперации не должны иметь возможности отменить родительскую. Вместо этого функция (описанная ниже) предоставляет способ отменить новое значение.`Cancel``Done``WithCancel``Context`

A `Context`безопасен для одновременного использования несколькими горутинами. Код может передать один объект `Context`любому количеству горутин и отменить его, `Context`чтобы подать сигнал всем горутинам.

Этот `Deadline`метод позволяет функциям определять, стоит ли им вообще начинать работу; если времени осталось слишком мало, работа может оказаться бесполезной. Код также может использовать крайний срок для установки тайм-аутов для операций ввода-вывода.

`Value`Позволяет `Context`переносить данные, относящиеся к области запроса. Эти данные должны быть безопасны для одновременного использования несколькими горутинами.

Пакет `context`предоставляет функции для _получения_ новых `Context`значений из существующих. Эти значения образуют дерево: при `Context`удалении a все `Contexts`производные от него значения также удаляются.

`Background`является корнем любого `Context`дерева; он никогда не отменяется:

```go
// Фон возвращает пустой контекст. Он никогда не отменяется, не имеет крайнего срока и не имеет значений. Фон обычно используется в main, init и тестах, 
// а также как контекст верхнего уровня для входящих запросов. 
func Background() Context
```


`WithCancel`и `WithTimeout`возвращают производные `Context`значения, которые могут быть отменены раньше, чем родительский `Context`. `Context`Связанный с входящим запросом объект обычно отменяется при возврате результата обработчика запросов. `WithCancel`Также полезен для отмены избыточных запросов при использовании нескольких реплик. `WithTimeout`Полезно для установки крайнего срока для запросов к внутренним серверам:

```go
// WithCancel возвращает копию родителя, канал Done которого закрывается, как только закрывается parent.Done или вызывается cancel. 
func WithCancel(parent Context) (ctx Context, cancel CancelFunc) 

// Функция CancelFunc отменяет Context. 
type CancelFunc func() 

// WithTimeout возвращает копию родителя, канал Done которого закрывается, как только закрывается parent.Done, вызывается cancel или истечет время ожидания. Крайний срок нового Context — это ближайшее из двух значений: now+timeout и крайний срок родителя, если он есть. Если таймер все еще работает, функция cancel освобождает свои ресурсы. 
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)
```


`WithValue`предоставляет способ связать значения области запроса с `Context`:
```go
// WithValue возвращает копию родителя, метод Value которого возвращает значение для ключа. 
func WithValue(parent Context, key interface{}, val interface{}) Context
```

---

Context'lar immutable — har bir With... funksiya yangi child yaratadi:

 context.Context — chuqur tushuntirish

Nima uchun yaratilgan?

Muammo: HTTP request keldi → siz 3 ta goroutine ochib DB, cache va tashqi API'ga murojaat qildingiz → client ulanishni uzdi. Bu goroutine'lar buni qayerdan biladi? Ular ishlashda davom etib, resurs isrof qiladi (goroutine leak, keraksiz DB load). Context — shu cancellation signal'ni daraxt bo'ylab tarqatish mexanizmi.

Interface

gotype Context interface {
    Deadline() (deadline time.Time, ok bool) // deadline bormi?
    Done() <-chan struct{}                   // bekor qilinganda YOPILADIGAN channel
    Err() error                              // nima uchun bekor: Canceled yoki DeadlineExceeded
    Value(key any) any                       // request-scoped qiymat
}

Asosiy mexanizm — Done() channel. Yopilgan channel'dan o'qish darhol qaytadi, shuning uchun selectda ishlatiladi:

goselect {
case <-ctx.Done():
    return ctx.Err() // context.Canceled yoki context.DeadlineExceeded
case result := <-resultCh:
    return result
}

Context daraxti (tree)

Context'lar immutable — har bir With... funksiya yangi child yaratadi:

goctx := context.Background()                          // root
ctx, cancel := context.WithTimeout(ctx, 5*time.Second) // child: 5s timeout
defer cancel()                                       // HAR DOIM defer cancel — resource leak'dan saqlaydi
ctx = context.WithValue(ctx, traceIDKey, "abc-123")  // child: qiymat bilan

Muhim qoida: parent bekor qilinsa, BARCHA child'lar ham bekor bo'ladi. Child bekor qilinsa — parent'ga ta'sir qilmaydi. Signal faqat pastga oqadi.

Turlari

FunksiyaVazifasicontext.Background()Root context — main, init, testlardacontext.TODO()"Hali qaysi context kerakligini bilmayman" — refactoring paytidaWithCancel(parent)Qo'lda bekor qilish uchun cancel funksiya qaytaradiWithTimeout(parent, d)d vaqtdan keyin avtomatik bekorWithDeadline(parent, t)Aniq vaqtda bekorWithValue(parent, k, v)Request-scoped data (trace ID, auth user)WithCancelCause(parent) (1.20+)Bekor sababini ko'rsatish: cancel(errors.New("...")); context.Cause(ctx)WithoutCancel(parent) (1.21+)Parent bekorligidan ajratish — masalan, response qaytgach ham davom etishi kerak bo'lgan audit log uchunAfterFunc(ctx, f) (1.21+)Context bekor bo'lganda f'ni chaqirish

Best practices (intervyuda aytish kerak)


Context — har doim birinchi parametr: func Fetch(ctx context.Context, id string). Struct ichida saqlamang (yagona istisno — http.Request).
Har doim defer cancel() — hatto timeout o'zi tugasa ham, cancel chaqirilmasa parent'da child reference qolib, xotira oqadi.
WithValueni suiiste'mol qilmang — faqat request-scoped metadata (trace ID, user ID) uchun. Funksiya parametrlarini context orqali "yashirincha" uzatish — anti-pattern, kod o'qib bo'lmas holga keladi.
Value key uchun unexported custom type ishlating — collision'dan saqlaydi:


go   type ctxKey string
   const traceIDKey ctxKey = "traceID"


nil context bermang — bilmasangiz context.TODO().
DB query'larga har doim context'li variantni ishlating: db.QueryContext(ctx, ...) — client ketsa, query PostgreSQL tomonda ham bekor qilinadi (pgx buni pg_cancel_backend orqali qiladi).


Follow-up savollar


"Context bekor qilinsa, goroutine avtomatik to'xtaydimi?" — YO'Q! Context faqat signal. Goroutine o'zi ctx.Done()ni tekshirishi kerak (for-select-done pattern). Buni bilmaslik — middle uchun qizil bayroq.
"Timeout DB query'ni to'xtatadimi?" — driver context'ni qo'llab-quvvatlasa (pgx, database/sql Context metodlari) — ha, server tomonda cancel yuboriladi.