## DDD (Domain-Driven Design) nima?

**Qisqacha:** Murakkab dasturiy ta'minotni loyihalash va amalga oshirish uchun asboblar to'plami.

**Maqsadi:**

- Biznes ehtiyojlariga mos dasturiy ta'minot yaratish
- Texnik murakkablik emas, **domain murakkabligiga** e'tibor berish

**Asosiy g'oya:** Biznes-ekspertlar va dasturchilar birgalikda umumiy til (Ubiquitous Language) yaratib, aniq chegaralangan kontekst (Bounded Context) ichida model quradilar.

---

## Big Ball of Mud nima?

**Ta'rif:** Chalkash, tartibsiz, chegarasiz tizim — dasturlashdagi eng yomon anti-pattern.

**Belgilari:**

- Juda ko'p kontseptsiyalar bitta modelda
- Aniq chegaralar yo'q
- Bir nechta til aralashgan
- Bir nechta jamoa ishlashi kerak (muammoli)
- Mustaqil kontseptsiyalar modullar bo'ylab tarqalgan
- Test qilish juda qiyin
- Qo'llab-quvvatlash mumkin emas

**Sabablari:**

- Loyihalashdan voz kechish
- Biznes-ekspertlarni tinglamaslik
- Texnologiyaga ortiqcha e'tibor
- Qachon to'xtatishni bilmaslik

**Oldini olish:** Bounded Context va Ubiquitous Language ishlatish!

---

## Bounded Context nima?

**Ta'rif:** Semantik kontekstual chegara — har bir komponent aniq ma'no va vazifaga ega bo'lgan chegara.

**Asosiy xususiyatlari:**

- Bir **Bounded Context** = bir **jamoa**
- Bir **Bounded Context** = bir **repository**
- Bir **Bounded Context** = bir **Ubiquitous Language**
- Ichida barcha komponentlar aniq va semantik jihatdan asoslangan

**Problem Space vs Solution Space:**

- **Problem Space:** Strategik tahlil, CORE DOMAIN aniqlash
- **Solution Space:** Kod yozish, amalga oshirish

```go
// Har bir Bounded Context - alohida Go module
myapp/
├── go.mod                          // Root module
├── cmd/                            // Entry points
├── internal/                       // Private code
│   ├── catalog/                    // Catalog Bounded Context
│   │   ├── domain/                 // Domain Model
│   │   │   ├── product.go          // Aggregate
│   │   │   ├── category.go         // Entity
│   │   │   └── price.go            // Value Object
│   │   ├── application/            // Application Services
│   │   └── infrastructure/         // Adapters
│   ├── ordering/                   // Ordering Bounded Context
│   └── inventory/                  // Inventory Bounded Context
└── pkg/                            // Public shared code
```

**Qoida:** `internal/` ichidagi har bir package - alohida Bounded Context bo'lishi mumkin.

Entity vs Value Object Go da

```Go
// Entity - ID bilan aniqlanadi
type Product struct {
    id          ProductID        // Identity
    name        string
    price       Money            // Value Object
    category    CategoryID
    createdAt   time.Time
}

// ID bilan taqqoslash
func (p *Product) Equals(other *Product) bool {
    return p.id == other.id
}

// Identity
type ProductID string

func NewProductID() ProductID {
    return ProductID(uuid.New().String())
}
```

#### Value Object (qiymat bilan aniqlanadi)

```go
// Value Object - immutable, ID yo'q
type Money struct {
    amount   decimal.Decimal
    currency string
}

// Qiymat bilan taqqoslash
func (m Money) Equals(other Money) bool {
    return m.amount.Equal(other.amount) && 
           m.currency == other.currency
}

// Yangi qiymat qaytaradi, o'zgarmas
func (m Money) Add(other Money) (Money, error) {
    if m.currency != other.currency {
        return Money{}, errors.New("currency mismatch")
    }
    return Money{
        amount:   m.amount.Add(other.amount),
        currency: m.currency,
    }, nil
}
```

**Farq:**

- **Entity:** O'zgaruvchan (mutable), ID bor, hayot tsikli bor
- **Value Object:** O'zgarmas (immutable), ID yo'q, faqat qiymat


---

## Core Domain nima?

**Ta'rif:** Tashkilotning asosiy strategik tashabbusi sifatida ishlab chiqilayotgan Bounded Context.

**Ahamiyati:**

- Eng muhim dasturiy model
- Eng yaxshi resurslar ajratilishi kerak
- Biznesning asosiy muammolarini hal qilishi kerak

---

## Ubiquitous Language nima?

**Ta'rif:** Jamoa a'zolari (dasturchilar va domain ekspertlar) o'rtasida umumiy, rasmiy, aniq til.

**Xususiyatlari:**

- **Rasmiy:** Qat'iy, aniq, konkret, ifodali
- **Hamma joyda:** Muloqot va kodda bir xil
- **Jonli:** Vaqt o'tishi bilan rivojlanadi
- **Umumiy:** Barcha jamoa a'zolari tushunadilar

**Tarkibi:**

- Ot so'zlar (noun): `Product`, `Sprint`, `BacklogItem`
- Fe'llar (verb): `commit`, `schedule`, `approve`
- Ravishlar va boshqa grammatik konstruksiyalar
- **Stsenariylar:** Qanday ishlashi kerakligi

---

## Subdomain nima?

**Ta'rif:** Biznesning alohida qismi — butun domenni kichikroq qismlarga bo'lish.

**Turlari:**

1. **Core Subdomain:** Raqobatdagi ustunlik
2. **Supporting Subdomain:** Core'ni qo'llab-quvvatlash
3. **Generic Subdomain:** Umumiy yechimlar (sotib olish mumkin)

**Maqsadi:**

- Legacy tizimlar murakkabligini kamaytirish
- Yangi loyihalarda aniq yo'nalish berish

---

## Context Mapping nima?

**Ta'rif:** Bir nechta Bounded Context'larni birlashtirish metodikasi.

**Context Map:** Bu metodikaning natijasi — kontekstlar o'rtasidagi munosabatlarni ko'rsatuvchi xarita/diagramma.

**Belgilanadi:**

- Jamoalar o'rtasidagi munosabatlar
- Texnik mexanizmlar
- Integratsiya strategiyalari

---

## Strategic Design (Strategik loyihalash) nima?

**Ta'rif:** Katta cho'tkalar bilan chizilgan rasm — umumiy yo'nalish.

**Asosiy vositalar:**

- Bounded Context
- Ubiquitous Language
- Context Mapping
- Subdomain

**Maqsadi:**

- Strategik muhim jihatlarni ajratish
- Ishni muhimlik darajasiga qarab bo'lish
- Integratsiya strategiyasini belgilash

---

## Tactical Design (Taktik loyihalash) nima?

**Ta'rif:** Nozik cho'tkacha — tafsilotlarni chizish.

**Asosiy vositalar:**

- **Aggregate:** Entity va Value Object'larni to'g'ri o'lchamda guruhlash
- **Domain Event:** Domenda sodir bo'layotgan voqealarni modellashtirish
- **Entity:** ID ga ega obyekt
- **Value Object:** ID siz, qiymat bilan aniqlanadi

**Maqsadi:** Aniq, to'g'ri dasturiy model yaratish

---

## Domain Expert kim?

**Ta'rif:** Domain bilan chuqur tanish, biznes bilimiga ega odam.

**Roli:**

- Ubiquitous Language'ni shakllantirish
- Biznes qoidalarini tushuntirish
- Stsenariylarni tasdiqlanish
- Dasturchilar bilan hamkorlik

**Muhim:** Bu lavozim emas, xarakteristika! Product Owner har doim Domain Expert bo'lavermaydi.

---

## Event Storming nima?

**Ta'rif:** Jamoaviy sessiya — domain voqealarini va jarayonlarini tez aniqlash usuli.

**Maqsadi:**

- Bilimni tez orttirish
- Stsenariylarni topish
- Ustuvorliklarni belgilash
- Bounded Context'larni aniqlash

---

## Ports and Adapters arxitekturasi nima?

**Ta'rif:** Hexagonal Architecture — Domain Model'ni texnologiyadan ajratish arxitekturasi.

**Qatlamlar:**

1. **Input Adapter** (tashqi dunyo → dastur)
2. **Application Service** (use case orkestratsiyasi)
3. **Domain Model** (biznes logika, texnologiyasiz!)
4. **Output Adapter** (dastur → tashqi dunyo)

**Maqsadi:** Domain Model texnologiyadan mustaqil bo'lsin!

---


## Aggregate Pattern Go da

```go
// Aggregate Root - faqat shu orqali kirish
type Order struct {
    // Aggregate identity
    id          OrderID
    customerID  CustomerID
    
    // Aggregate ichidagi Entity'lar
    items       []OrderItem      // Entity
    
    // Value Object'lar
    shippingAddress Address
    totalAmount     Money
    
    // Aggregate state
    status      OrderStatus
    placedAt    time.Time
}

// Factory method
func NewOrder(customerID CustomerID, address Address) *Order {
    return &Order{
        id:              NewOrderID(),
        customerID:      customerID,
        items:           make([]OrderItem, 0),
        shippingAddress: address,
        status:          OrderStatusDraft,
        placedAt:        time.Now(),
    }
}

// Biznes qoidalari - faqat Aggregate orqali
func (o *Order) AddItem(productID ProductID, quantity int, price Money) error {
    // Biznes qoidasini tekshirish
    if o.status != OrderStatusDraft {
        return errors.New("cannot modify placed order")
    }
    
    if quantity <= 0 {
        return errors.New("quantity must be positive")
    }
    
    // Qo'shish
    item := OrderItem{
        productID: productID,
        quantity:  quantity,
        price:     price,
    }
    o.items = append(o.items, item)
    
    // Total'ni yangilash
    o.recalculateTotal()
    
    return nil
}

// Private method - tashqaridan kirish yo'q
func (o *Order) recalculateTotal() {
    total := Money{amount: decimal.Zero, currency: "USD"}
    for _, item := range o.items {
        itemTotal := item.price.Multiply(item.quantity)
        total = total.Add(itemTotal)
    }
    o.totalAmount = total
}

// OrderItem - Aggregate ichidagi Entity
type OrderItem struct {
    productID ProductID
    quantity  int
    price     Money
}
```

**Aggregate qoidalari:**

1. Tashqaridan faqat **root** ga kirish mumkin
2. Root ichidagi boshqa obyektlarga to'g'ridan-to'g'ri kirish yo'q
3. Bitta tranzaksiya - bitta Aggregate
4. Aggregate'lar o'rtasida faqat ID orqali reference

---

## Repository Pattern Go da

```go
// Repository interface - Domain qatlamida
type ProductRepository interface {
    Save(ctx context.Context, product *Product) error
    FindByID(ctx context.Context, id ProductID) (*Product, error)
    FindByCategory(ctx context.Context, categoryID CategoryID) ([]*Product, error)
    Delete(ctx context.Context, id ProductID) error
}

// Konkret implementation - Infrastructure qatlamida
type PostgresProductRepository struct {
    db *sql.DB
}

func NewPostgresProductRepository(db *sql.DB) *PostgresProductRepository {
    return &PostgresProductRepository{db: db}
}

func (r *PostgresProductRepository) Save(ctx context.Context, product *Product) error {
    query := `
        INSERT INTO products (id, name, price_amount, price_currency, category_id)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (id) DO UPDATE SET
            name = EXCLUDED.name,
            price_amount = EXCLUDED.price_amount,
            price_currency = EXCLUDED.price_currency,
            category_id = EXCLUDED.category_id
    `
    _, err := r.db.ExecContext(ctx, query,
        product.id,
        product.name,
        product.price.amount,
        product.price.currency,
        product.category,
    )
    return err
}
```

**Muhim:**

- Interface - Domain qatlamida
- Implementation - Infrastructure qatlamida
- Dependency Inversion Principle

---

## Asosiy o'xshashliklar: Go ↔ DDD

| Go tushunchasi          | DDD tushunchasi               | Misol                                |
| ----------------------- | ----------------------------- | ------------------------------------ |
| `struct`                | Entity / Value Object         | `type Product struct`                |
| `struct` pointer method | Aggregate behavior            | `func (p *Product) AddBacklogItem()` |
| `interface`             | Repository, Service contracts | `type ProductRepository interface`   |
| Package                 | Bounded Context               | `package scrum`                      |
| Lowercase field         | Private/encapsulated          | `id ProductID`                       |
| Uppercase field         | Public                        | `Name string`                        |
| Factory function        | Constructor                   | `func NewProduct()`                  |
| `slice`                 | Collection                    | `backlogItems []BacklogItem`         |
| `map`                   | Lookup table                  | `handlers map[string]Handler`        |
| `channel`               | Event bus                     | `events chan DomainEvent`            |
| Goroutine               | Async event handler           | `go handler.Handle(event)`           |
| Context                 | Transaction scope             | `ctx context.Context`                |
|                         |                               |                                      |

## Xulosa

Scrum misolida ko'rdik:

1. **Bounded Context** = Go package (`internal/scrum`)
2. **Aggregate** = Go struct + method'lar (`Product`, `Sprint`, `Team`)
3. **Entity** = Identity bilan struct (`BacklogItem`, `Task`)
4. **Value Object** = Immutable struct (`Email`, `TeamQuorum`, `Volunteer`)
5. **Repository** = Go interface (domain) + implementation (infrastructure)
6. **Domain Event** = Go struct + interface
7. **Application Service** = Use case orchestrator
8. **Ubiquitous Language** = Method nomlari (`CommitToSprint`, `FormQuorum`)