# Loyihalash Patternlari — Asoslar

---

## 1. Pattern nima?

**Design Pattern (loyihalash patterni)** — dastur arxitekturasini loyihalashda tez-tez uchraydigan muammolarga **sinovdan o'tgan, umumlashtirilgan yechim**.

Pattern — bu tayyor funksiya yoki kutubxona emas. Uni shunchaki nusxalab kodga qo'yib bo'lmaydi. Pattern — bu **yechim kontseptsiyasi**, uni har bir dasturning o'ziga xos ehtiyojlariga moslashtirib qo'llash kerak.

### Pattern nimalardan iborat?

| Bo'lim | Mazmuni |
|--------|---------|
| **Muammo** | Pattern qaysi muammoni hal qiladi |
| **Motivatsiya** | Nima uchun aynan shu yechim taklif qilinadi |
| **Struktura** | Klasslar va ular orasidagi munosabatlar (UML) |
| **Kod misoli** | Biror dasturlash tilida amalga oshirish |
| **Qo'llash holatlari** | Qachon va qaerda ishlatish kerak |
| **Boshqa patternlar bilan aloqasi** | Qaysi patternlar bilan birgalikda ishlatiladi |

---

## 2. Yaxshi dizayn nima?

> - **S1:** Yaxshi kodni yomonidan qanday **ajratish** mumkin?
> - **S2:** Qanday **mezonlar** bo'yicha baholash kerak?
> - **S3:** Moslashuvchanlik, bog'liqlik, boshqaruvchanlik, barqarorlik va tushunarlilikni qanday **ta'minlash** mumkin?

---

### 2.1 → S1: Yaxshi kodni yomonidan qanday ajratish mumkin?

Yaxshi kod — bu **o'zgarishga bardoshli** kod. Agar bitta talabni o'zgartirganda butun tizim "titrasа" — bu yomon dizayn. Agar o'zgarish faqat bir joyda izolyatsiyalansa — bu yaxshi dizayn.

Bundan tashqari, **kodni qayta ishlatish imkoniyati** ham ajratuvchi mezon hisoblanadi. Erix Gamma (GoF muallifi) kodni qayta ishlatishning uch darajasini ko'rsatadi:

```
Daraja 3 — Fremvorklar    → butun tizim skeleti (Django, Spring, Echo)
Daraja 2 — Patternlar     → g'oyalar va munosabatlar qayta ishlatiladi  ← eng muvozanatli
Daraja 1 — Kutubxonalar   → tayyor funksiyalar (math, json, sort)
```

Patternlar — fremvorklarga qaraganda **kamroq xavfli** va **arzonroq** qayta ishlatish yo'li, chunki ular konkret kod emas, g'oyalarni qayta ishlatadi.

---

### 2.2 → S2: Qanday mezonlar bo'yicha baholash kerak?

| Mezon | Ta'rifi | Yomon bo'lsa... |
|-------|---------|-----------------|
| **Moslashuvchanlik** | Yangi talablarga oson moslashadi | Kichik o'zgarish uchun ko'p joy tahrirlash kerak |
| **Kengaytiruvchanlik** | Mavjud kodni o'zgartirmasdan yangi funksiya qo'shiladi | Yangi funksiya eski kodni buzadi |
| **Boshqaruvchanlik** | Kodni tushunish va o'zgartirish oson | Yangi dasturchi kodni o'qib tushunmaydi |
| **Barqarorlik** | Bir qism o'zgarganda boshqalar buzilmaydi | Bitta xato boshqa modullarga tarqaladi |
| **Qayta ishlatish** | Komponentlar boshqa loyihalarda ham ishlaydi | Har safar noldan yoziladi |

> **Amaliy test:** Agar yangi talabni amalga oshirishda "bu yerdan ham o'zgartirish kerak, u yerdan ham..." degan fikr kelsa — dizayn yomon.

---

### 2.3 → S3: Bu sifatlarni qanday ta'minlash mumkin?

Ikki asosiy tushuncha orqali:

#### Coupling (Bog'liqlik darajasi) — mumkin qadar past bo'lsin

**Coupling** — modullar orasidagi o'zaro bog'liqlik. Biri o'zgarganda ikkinchisi ham o'zgarishga majbur bo'lsa — coupling yuqori.

```
Yuqori Coupling (yomon):         Past Coupling (yaxshi):
A ═══════ B                      A - - - - B
A o'zgarsa, B ham o'zgaradi      A o'zgarsa, B ta'sirlanmaydi
```

Pasaytirish usullari: interface orqali muloqot, Dependency Injection, event/message orqali aloqa.

#### Cohesion (Yaxlitlik darajasi) — mumkin qadar yuqori bo'lsin

**Cohesion** — bir modul ichidagi elementlarning bir-biriga tegishliligi. Bir klassda "email yuborish", "fayl o'qish", "parol shifrlash" bo'lsa — cohesion past.

```
Past Cohesion (yomon):           Yuqori Cohesion (yaxshi):
┌──────────────────┐             ┌──────────┐  ┌──────────┐
│ "Utils" klassi   │      →      │FileReader│  │EmailSender│
│ + faylOqi()      │             │ + oqi()  │  │+ yuborish│
│ + emailYubor()   │             └──────────┘  └──────────┘
│ + parolShifr()   │
└──────────────────┘
```

**Formula:**

```
Past Coupling  +  Yuqori Cohesion  =  Yaxshi dizayn
```

---

## 3. Dizayn tamoyillari

Bu tamoyillar yuqoridagi sifatlarni amalda ta'minlash uchun ishlatiladi. Ko'pchilik patternlar aynan shu tamoyillarga asoslanadi.

### Tamoyil 1: O'zgaruvchan narsani kapsulalash

> *Ko'p o'zgaradigan qismlarni topib, o'zgarmaydigan qismdan ajrating.*

**Analogiya:** Kema seksiyalarga bo'linsa, bitta mina faqat bitta seksiyani yo'q qiladi — qolganlar sog'lom qoladi. Xuddi shunday, o'zgaruvchan mantiqni izolyatsiya qilsangiz, o'zgarish faqat shu joyda qoladi.

**Amalda:** Soliq hisoblash mantiqini asosiy metod ichidan alohida metodga, keyin alohida klassga chiqarish — shu tamoyilning misoli. Soliq qonunlari o'zgarganda faqat shu klass o'zgaradi.

---

### Tamoyil 2: Interfeysga dasturlang, amalga oshirishga emas

> *Kod konkret klasslarga emas, abstraktsiyalarga (interfeys) bog'liq bo'lsin.*

**Nima uchun:** Konkret klassga bog'liq kod o'sha klass o'zgarganda buziladi. Interfeyga bog'liq kod esa — interfeys o'zgarmasa ishlayveradi, ichki amalga oshirish qanday o'zgarsayam.

**4 qadam:**
1. Ob'ekt boshqa ob'ektdan nima kerakligini aniqla (qaysi metodlar)
2. Shu metodlarni alohida interfeys sifatida e'lon qil
3. Bog'liq klassni shu interfeysni amalga oshirishga o'tkazdir
4. Asosiy kod endi interfeysga bog'liq — konkretga emas

**Natija:** Keyinchalik amalga oshirishni almashtirsangiz, asosiy kodga tegilmaydi.

> Bu tamoyilni amalga oshirishning bir ko'rinishi — **Factory Method** pattern.

---

### Tamoyil 3: Meros olishdan ko'ra kompozitsiyani afzal ko'ring

> *"Has-a" ("tarkibida bor") munosabatini "Is-a" ("hisoblanadi") munosabatiga afzal ko'ring.*

**Merosning 5 muammosi:**

| Muammo | Izoh |
|--------|------|
| Interfeydan voz kechib bo'lmaydi | Keraksiz metodlarni ham amalga oshirish kerak |
| Inkapsulatsiya buziladi | Avlod ota-klass ichki tafsilotlarini ko'radi |
| Qattiq bog'liqlik | Ota-klass o'zgarganda avlodlar buzilishi mumkin |
| Klass portlashi | Har yangi kombinatsiya yangi klass talab qiladi |
| Ko'p merosdan foydalanib bo'lmaydi | Ko'pchilik tillar buni qo'llab-quvvatlamaydi |

**Klass portlashi misoli:** Avtomobil (elektr/benzin) × (yengil/yuk) × (qo'l/avtopilot) = 2×2×2 = **8 klass**. Yana bir parametr qo'shilsa — 16 ta. Kompozitsiyada esa 3 ta interfeys va bir nechta klass yetarli, kombinatsiyalar ob'ektlar orqali hal qilinadi.

**Qo'shimcha afzallik:** Kompozitsiyada xatti-harakatni **dastur ishlayotgan paytda** (runtime) almashtirsa bo'ladi. Merosda esa bu mumkin emas — bog'lanish kompilyatsiya vaqtida sodir bo'ladi.

> Bu tamoyilni amalga oshirishning bir ko'rinishi — **Strategy** pattern.

---

## 4. Patternlarning tasnifi

```
Oddiy ─────────────────────────────── Murakkab
  │                                       │
Idiomalar    Dizayn patternlar    Arxitektura patternlar
(1 til)      (ko'p tilda)         (butun tizim: MVC, CQRS)
```

**23 ta GoF Pattern — 3 guruh:**

```
┌──────────────┐   ┌──────────────┐   ┌───────────────────┐
│  Yaratuvchi  │   │ Tuzilmaviy   │   │   Xulq-atvoriy    │
│ (Creational) │   │ (Structural) │   │   (Behavioral)    │
│              │   │              │   │                   │
│ Factory      │   │ Adapter      │   │ Observer          │
│ Abstract F.  │   │ Bridge       │   │ Strategy          │
│ Builder      │   │ Composite    │   │ Command           │
│ Prototype    │   │ Decorator    │   │ Iterator          │
│ Singleton    │   │ Facade       │   │ State             │
│              │   │ Flyweight    │   │ Template Method   │
│              │   │ Proxy        │   │ Chain of Resp.    │
│              │   │              │   │ Mediator          │
│              │   │              │   │ Memento, Visitor  │
└──────────────┘   └──────────────┘   └───────────────────┘
```

**Yaratuvchi** — ob'ektlarni moslashuvchan yaratish, keraksiz bog'liqliksiz.

**Tuzilmaviy** — ob'ektlar orasida samarali munosabatlar o'rnatish.

**Xulq-atvoriy** — ob'ektlar orasidagi samarali muloqot va mas'uliyatni taqsimlash.

---

## 5. Tarix va nima uchun o'rganish kerak?

```
1977 → Kristofer Aleksandr — "Shablonlar tili" (arxitektura sohasi)
1994 → Gang of Four — 23 ta OOP pattern (GoF kitobi)
Bugun → 100+ pattern: Cloud, Concurrency, Enterprise...
```

**Patternlar beradi:**

| Foyda | Tushuntirish |
|-------|-------------|
| **Sinovdan o'tgan yechimlar** | G'ildirakni qayta ixtiro qilmaysiz |
| **Standart lug'at** | "Observer ishlatamiz" — darhol tushuniladi |
| **Kamroq xato** | Barcha yashirin muammolar allaqachon hal qilingan |

> **Ogohlantirish:** Pattern muammoga yechim, muammo esa patternni izlash uchun bahona emas. Oddiy muammo uchun murakkab pattern ishlatish — o'zi antipattern.

---

## 6. Xulosa

```
Yaxshi dizayn = Past Coupling + Yuqori Cohesion + 3 tamoyil

  Tamoyil 1: O'zgaruvchan narsani kapsulalash
  Tamoyil 2: Interfeyga dasturlash (konkretga emas)
  Tamoyil 3: Kompozitsiya > Meros

  Bu 3 tamoyil → 23 GoF patternning asosi

        ┌────────────────────────────────────────┐
        │           Yaxshi Kod                   │
        │  ✅ Moslashuvchan  ✅ Kengaytiruvchan  │
        │  ✅ Boshqaruvchan  ✅ Barqaror         │
        │  ✅ Qayta ishlatish mumkin             │
        └────────────────────────────────────────┘
```
