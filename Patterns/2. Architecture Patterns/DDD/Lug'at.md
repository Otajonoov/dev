
## A

### Aggregate (Agregat)

Bir xil transaksiya chegarasida boshqariladigan Entity va Value Object'larning klasteridir. Har bir Aggregate bitta Aggregate Root'ga ega. Consistency boundary (izchillik chegarasi)ni belgilaydi.

### Aggregate Root (Agregat Ildizi)

Aggregate ichidagi asosiy Entity. Tashqi dunyo aggregate bilan faqat shu root orqali muloqot qiladi. Barcha o'zgarishlar va biznes qoidalari shu yerda amalga oshiriladi.

### Anemic Domain Model (Anemik Domen Model)

**Anti-pattern.** Faqat getter/setter'larga ega bo'lgan va biznes mantiqi bo'lmagan domain model. Biznes mantiqi service'larda yozilgan bo'ladi, bu yomon dizayn hisoblanadi.

### Anticorruption Layer (ACL) (Himoya Qatlami)

O'z modelingizni tashqi tizim (legacy yoki 3rd-party) modelidan izolyatsiya qiluvchi translation qatlami. Ikki turli Ubiquitous Language o'rtasida tarjimon vazifasini bajaradi.

### Application Service (Ilova Xizmati)

Use case'larni orchestrate qiluvchi service. Aggregate'larni yuklaydi, command'larni bajaradi, tranzaksiyalarni boshqaradi, domain event'larni nashr etadi. Domain mantiqni o'z ichiga olmaydi.

---

## B

### Behavior (Xatti-harakat)

Aggregate yoki Entity'ning metodlari orqali ifodalangan biznes mantiqi. Anemic model'dan qochish uchun barcha biznes mantiqi behavior method'lar orqali amalga oshirilishi kerak.

### Big Ball of Mud (Iflos Loy To'pi)

**Anti-pattern.** Hech qanday aniq arxitektura yoki modelga ega bo'lmagan, chigal va murakkab tizim. Ko'plab legacy tizimlarning holati. Bunday tizimlardan qochish yoki ACL orqali izolyatsiya qilish kerak.

### Bounded Context (Chegaralangan Kontekst)

Muayyan model va Ubiquitous Language qo'llaniladigan aniq chegaralangan hudud. Bir tizimda bir nechta Bounded Context bo'lishi mumkin. Bu DDD ning strategik dizaynining asosiy tushunchasi.

### Business Logic (Biznes Mantiqi)

Biznes qoidalari, jarayonlar va hisob-kitoblar. DDD da bu mantiq Domain Layer'da, asosan Aggregate'larda joylashgan bo'lishi kerak.

### Business Rules (Biznes Qoidalari)

Domain ekspertlari tomonidan belgilangan qoidalar va cheklovlar. Masalan: "faqat aktiv foydalanuvchiga rol tayinlash mumkin" yoki "buyurtma summasi 0 dan katta bo'lishi kerak".

---

## C

### Command (Buyruq)

Tizimda biror o'zgarishni amalga oshirishga niyatni ifodalovchi obyekt. CQRS pattern'ida ishlatiladi. Nomi buyruq shaklida (imperativ) bo'ladi: CreateUser, AssignRole.

### Command Handler (Buyruq Boshqaruvchisi)

Command'ni qabul qilib, uni bajaradigan component. Application Service'ning bir turi. Use case'ni amalga oshiradi.

### Conformist (Moslashuvchi)

Context Mapping pattern'i. Downstream jamoasi upstream modeliga to'liq moslashib, o'z translation layer'ini yaratmaydi. Odatda katta, mashhur tizimlar bilan integratsiyada ishlatiladi (masalan, Salesforce, AWS).

### Consistency Boundary (Izchillik Chegarasi)

Aggregate chegarasi. Bir tranzaksiya ichida izchil bo'lishi kafolatlangan ma'lumotlar to'plami. Aggregate chegarasidan tashqaridagi ma'lumotlar Eventual Consistency orqali yangilanadi.

### Context Map (Kontekst Xaritasi)

Tizimda mavjud barcha Bounded Context'lar va ular orasidagi munosabatlarni vizual ko'rsatuvchi xarita. Strategic dizaynning muhim vositasi.

### Context Mapping (Kontekst Xaritalash)

Turli Bounded Context'lar orasidagi integratsiya munosabatlarini aniqlash va boshqarish jarayoni. Partnership, Shared Kernel, Customer-Supplier kabi pattern'larni o'z ichiga oladi.

### Core Domain (Asosiy Domen)

Biznesingizga raqobat ustunligi beradigan, eng muhim va strategik subdomain. Bu yerga eng ko'p investitsiya va eng yaxshi dasturchilar jalb qilinishi kerak.

### CQRS (Command Query Responsibility Segregation)

Command (yozish) va Query (o'qish) operatsiyalarini ajratish pattern'i. Har biri o'z model va data store'iga ega bo'lishi mumkin.

### Customer-Supplier (Mijoz-Yetkazib beruvchi)

Context Mapping pattern'i. Upstream (Supplier) downstream (Customer) ehtiyojlarini qondirishi kerak, lekin yakuniy qarorlarni o'zi qabul qiladi. Odatda bir tashkilot ichidagi jamoalar o'rtasida ishlatiladi.

---

## D

### Domain (Domen)

Biznes sohaси, muammo maydoni. Masalan: e-commerce, banking, healthcare. DDD shu domenda dasturiy model yaratishga qaratilgan.

### Domain Event (Domen Hodisasi)

Domenda sodir bo'lgan muhim biznes hodisasini ifodalovchi obyekt. O'tmish zamonda nomlanadi: UserRegistered, OrderPlaced. Aggregate'lar tomonidan nashr etiladi va Eventual Consistency uchun ishlatiladi.

### Domain Expert (Domen Eksperti)

Biznes sohasida chuqur bilimga ega bo'lgan odam. Product owner, biznes tahlilchi, foydalanuvchi yoki har qanday domain haqida batafsil biladigan shaxs. Ular bilan yaqin hamkorlik DDD ning kaliti.

### Domain Layer (Domen Qatlami)

Biznes mantiqining joylashgan joyи. Aggregate'lar, Entity'lar, Value Object'lar, Domain Service'lar shu qatlamda. Hech qanday texnik dependencylar bo'lmasligi kerak (database, framework).

### Domain Logic (Domen Mantiqi)

Biznes qoidalarini dasturiy shaklda ifodalash. Masalan: narxlarni hisoblash, validatsiya, workflow mantiqi.

### Domain Model (Domen Modeli)

Real biznes jarayonlarini dasturiy abstraksiyalar orqali ifodalash. Entity'lar, Value Object'lar, Aggregate'lar yig'indisi.

### Domain Service (Domen Xizmati)

Bitta Aggregate'ga tegishli bo'lmagan biznes mantiqni o'z ichiga olgan service. Bir nechta Aggregate bilan ishlaydi yoki stateless operatsiyalarni bajaradi. Masalan: AuthorizationService, PricingService.

### DTO (Data Transfer Object)

Qatlamlar o'rtasida ma'lumot uzatish uchun ishlatiluvchi oddiy data strukturasi. Domain modelini tashqi dunyodan izolyatsiya qilish uchun kerak.

---

## E

### Entity (Mavjudlik)

O'zining noyob identifikatsiyasiga ega bo'lgan domain obyekti. Ikki entity bir xil ma'lumotlarga ega bo'lsa ham, agar ID'lari farq qilsa, ular turli obyektlar. Odatda o'zgaruvchan (mutable).

### Event Sourcing (Hodisalar Manbai)

Aggregate holatini to'g'ridan-to'g'ri saqlash o'rniga, uni o'zgartirgan barcha Domain Event'larni saqlash pattern'i. Har qanday vaqtdagi holatni event'larni qayta o'ynash orqali tiklash mumkin.

### Event Storming (Hodisa Shturmi)

Domain'ni tez o'rganish va loyihalash uchun workshop metodikasi. Domain ekspertlari va dasturchilar birgalikda stikerlar yordamida biznes jarayonlarini Domain Event'lar orqali modellashtiadi.

### Eventual Consistency (Oxir-oqibat Izchillik)

Bir Aggregate o'zgarganida, boshqa Aggregate'lar darhol emas, balki ma'lum vaqt ichida (odatda soniyalar yoki daqiqalar) yangilanadi. Domain Event'lar orqali amalga oshiriladi.

---

## F

### Factory (Zavod)

Murakkab Aggregate yoki Entity yaratish mantiqini inkapsulyatsiya qiluvchi pattern. Constructor yetarli bo'lmagan hollarda ishlatiladi.

### Factory Method (Zavod Metodi)

Aggregate ichida static yoki instance metod sifatida amalga oshirilgan factory. Masalan: User.Register(), Order.CreateFrom().

---

## G

### Generic Subdomain (Universal Quyi Domen)

Raqobat ustunligi bermaydigan, tayyor yechimlar mavjud bo'lgan subdomain. Masalan: email yuborish, PDF generatsiya. Kutubxonalar yoki 3rd-party servislardan foydalanish mumkin.

---

## I

### Idempotent (Idempotent)

Bir xil operatsiyani bir necha marta bajarish bitta bajarishga teng natija bersa. Messaging tizimlarida muhim, chunki xabar bir necha marta yetkazilishi mumkin.

### Infrastructure Layer (Infratuzilma Qatlami)

Texnik implementatsiyalar qatlami: database, messaging, external API'lar, file system. Domain Layer interface'larini implement qiladi.

### Intention Revealing Interface (Niyatni Ko'rsatuvchi Interface)

Metod nomlari nima qilishini emas, balki nima uchun qilinishini ko'rsatishi kerak. Masalan: CalculateTotal() o'rniga ApplyDiscountPolicy().

### Invariant (O'zgarmas Qoida)

Har doim bajarilishi kerak bo'lgan biznes qoidasi. Masalan: Order total must be positive, User must have unique email. Aggregate bu invariantlarni himoya qiladi.

---

## L

### Layered Architecture (Qatlamli Arxitektura)

Dasturni mantiqiy qatlamlarga bo'lish: Interface, Application, Domain, Infrastructure. Har bir qatlam faqat o'zidan past qatlamlarga bog'liq.

### Legacy System (Meros Tizim)

Eski, ko'pincha yomon dizaynlangan tizim. DDD da bunday tizimlar bilan ACL orqali integratsiya qilinadi yoki Strangler Fig pattern orqali asta-sekin almashtiriladi.

---

## M

### Model (Model)

Real dunyo jarayonlarining dasturiy abstraksiyasi. DDD da model faqat ma'lumotlar strukturasi emas, balki biznes mantiq va xatti-harakatni ham o'z ichiga oladi.

### Module (Modul)

Bog'liq domain tushunchalarini guruhlash uchun ishlatiluvchi logical package. Go da package, Java da package, C# da namespace. Modullar Ubiquitous Language'ni aks ettirishi kerak.

---

## O

### Open Host Service (Ochiq Xost Xizmati)

Context Mapping pattern'i. Bir Bounded Context boshqalarga integratsiya uchun yaxshi hujjatlashtirilgan, barqaror API taqdim etadi. Ko'pincha Published Language bilan birgalikda ishlatiladi.

---

## P

### Partnership (Hamkorlik)

Context Mapping pattern'i. Ikki jamoa bir-biriga bog'liq bo'lib, birgalikda muvaffaqiyatga erishadi yoki mag'lub bo'ladi. Doimiy kommunikatsiya va koordinatsiya talab etiladi.

### Persistence Ignorance (Persistentsiyadan Bexabarlik)

Domain model database strukturasi haqida bilmasligi kerak printsipi. Repository pattern bu printsipni amalga oshiradi.

### Problem Space (Muammo Maydoni)

Hal qilish kerak bo'lgan biznes muammosi. Strategic dizaynda Subdomain'lar Problem Space'ni tashkil qiladi.

### Published Language (Nashr Qilingan Til)

Turli sistemalar o'rtasida ma'lumot almashish uchun yaxshi hujjatlashtirilgan, standartlashtirilgan format. Masalan: JSON Schema, XML Schema, Protobuf.

---

## Q

### Query (So'rov)

Tizimdan ma'lumot olish operatsiyasi, o'zgarishlar kiritмайdi. CQRS da Command'dan ajratiladi.

### Query Handler (So'rov Boshqaruvchisi)

Query'ni bajarib, DTO'lar shaklida ma'lumot qaytaruvchi component. Read model bilan ishlaydi.

---

## R

### Repository (Repositoriya)

Aggregate'larni saqlash va olish uchun abstraksiya. Domain Layer'da interface, Infrastructure Layer'da implementation. Collection kabi ishlaydi.

### Rich Domain Model (Boy Domen Model)

Biznes mantiq va xatti-harakatga ega bo'lgan domain model. Anemic model'ning teskarisi. DDD da maqsad - Rich Domain Model yaratish.

---

## S

### Separate Ways (Alohida Yo'llar)

Context Mapping pattern'i. Ikki Bounded Context o'rtasida integratsiya juda qimmat yoki keraksiz. Har biri o'z yo'li bilan ketadi.

### Shared Kernel (Umumiy Yadro)

Context Mapping pattern'i. Ikki yoki undan ortiq jamoa kichik bir qism modelni baham ko'radi. Juda ehtiyotkorlik bilan boshqarish kerak, chunki bir jamoa o'zgarishi ikkinchisiga ta'sir qiladi.

### Side Effect Free Function (Yon Ta'sirsiz Funksiya)

O'z parametrlarini o'zgartirmaydigan, faqat yangi qiymat qaytaradigan funksiya. Funktsional dasturlashda va Value Object'larda ishlatiladi.

### Solution Space (Yechim Maydoni)

Muammoni qanday hal qilish. Strategic dizaynda Bounded Context'lar Solution Space'ni tashkil qiladi.

### Specification (Spetsifikatsiya)

Biznes qoidasini testable predicate shaklida ifodalash pattern'i. Masalan: eligibleForDiscount(customer), canShipToCountry(order, country).

### Subdomain (Quyi Domen)

Domainning mantiqiy qismi. Uch turi: Core Domain, Supporting Subdomain, Generic Subdomain. Problem Space'da yashaydi.

### Supporting Subdomain (Yordamchi Quyi Domen)

Zarur, ammo raqobat ustunligi bermaydigan subdomain. Tayyor yechim mavjud emas, shuning uchun o'zimiz yozamiz, lekin eng yaxshi resurslarni boshqa joylarga sarflaymiz.

---

## T

### Tactical Design (Taktik Dizayn)

Code level dizayn: Aggregate, Entity, Value Object, Repository, Domain Service va boshqalar. "How" (qanday) ni javob beradi.

### Transaction (Tranzaksiya)

Ma'lumotlar izchilligini ta'minlovchi atomik operatsiya. DDD da bir tranzaksiyada faqat bitta Aggregate o'zgaradi.

---

## U

### Ubiquitous Language (Hamma Joyda Ishlatiladigan Til)

Domen ekspertlari va dasturchilar o'rtasida umumiy til. Kod, hujjatlar, suhbatlarda bir xil terminlar ishlatiladi. DDD ning asosiy printsipi.

### Use Case (Foydalanish Holati)

Foydalanuvchi amalga oshirmoqchi bo'lgan maqsad. Application Service'lar use case'larni amalga oshiradi.

---

## V

### Value Object (Qiymat Obyekti)

Identifikatsiyasi yo'q, faqat qiymati bilan farqlanadigan obyekt. Immutable (o'zgarmas) bo'lishi kerak. Masalan: Money, Email, Address, DateRange.

### Versioning (Versiyalash)

API yoki Published Language'ning turli versiyalarini boshqarish. Backward compatibility uchun muhim.

---

## Qo'shimcha Muhim Tushunchalar

### Consistency (Izchillik)

Ma'lumotlarning to'g'ri va qarama-qarshi bo'lmagan holati. Immediate Consistency (bir tranzaksiyada) va Eventual Consistency (vaqt o'tishi bilan) turlari mavjud.

### Coupling (Bog'lanish)

Komponentlar o'rtasidagi bog'liqlik darajasi. DDD maqsadi - loose coupling (kuchsiz bog'lanish) yaratish.

### Encapsulation (Inkapsulyatsiya)

Ichki holatni va implementatsiyani yashirish, faqat aniq interface orqali kirish berish. DDD da Aggregate'lar strong encapsulation yaratadi.

### Orchestration (Orkestrlash)

Bir nechta komponentlarni birgalikda muvofiqlashtirish. Application Service'lar use case'larni orchestrate qiladi.

### Side Effect (Yon Ta'sir)

Funksiya yoki metodning o'z qiymatidan tashqari ta'siri. Masalan: database'ga yozish, xabar yuborish.

---

# Xulosa

DDD terminlari:

- **Strategic** - yuqori darajadagi dizayn (Bounded Context, Subdomain, Context Mapping)
- **Tactical** - kod darajasidagi pattern'lar (Aggregate, Entity, Value Object)
- **Arxitektura** - tizim strukturasi (Layered Architecture, CQRS, Event Sourcing)
- **Jarayon** - loyihalash va hamkorlik (Event Storming, Ubiquitous Language)

Barcha bu tushunchalar birgalikda murakkab biznes tizimlarini boshqariladigan, tushunarli va o'zgaruvchan qilishga yordam beradi.