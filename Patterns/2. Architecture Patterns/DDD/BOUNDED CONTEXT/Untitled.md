
ОГРАНИЧЕННЫЙ КОНТЕКСТ

### Bounded Context

Bu — **semantik kontekstual chegara**. Bounded Context ichida dasturiy ta'minotning har bir komponenti:

- Aniq ma'noga ega
- Aniq vazifalarni bajaradi
- Shu kontekstga xos
- Semantik jihatdan asoslangan

**Muhim:** Dasturiy ta'minotni yaratishni boshlaganingizda, Bounded Context nazariy bo'ladi (**problem space**). Model chuqurroq va aniqroq bo'lgani sari, Bounded Context tezda **solution space**ga o'tadi va dasturiy model loyiha manba kodida aks etadi.

### Problem space va Solution space

**Problem space** — yuqori darajadagi strategik tahlil va loyihalash bosqichlari amalga oshiriladigan soha. Bu yerda:

- Loyihaga ta'sir qiluvchi asosiy omillar muhokama qilinadi
- Muhim maqsadlar va xavflar ajratiladi
- Context Map ishlatiladi
- Oddiy diagrammalar qo'llaniladi

**Solution space** — yechim amalda amalga oshiriladigan soha. Bu yerda:

- **Core Domain** (Asosiy domen) sifatida aniqlangan yechim amalga oshiriladi
- Bounded Context manba kodiga aylanadi (asosiy va test kod)
- Boshqa Bounded Context'lar bilan integratsiya kodi yoziladi

### Core Domain

**Core Domain** — tashkilotning asosiy strategik tashabbusi sifatida ishlab chiqilayotgan Bounded Context.

Bu sizning tashkilotingiz uchun **eng muhim dasturiy model**, chunki:

- Tashkilotga raqobatdagi ustunlik berishi kerak
- Biznesingizning kamida asosiy muammolarini hal qilishi kerak
- Eng yaxshi resurslar shu yerga ajratilishi kerak

## Bir jamoa — bir Bounded Context qoidasi

**Qat'iy qoida:**

- Har bir Bounded Context ustida **faqat bitta jamoa** ishlashi kerak
- Har bir Bounded Context **alohida manba kodi repository**ga ega bo'lishi kerak
- Bitta jamoa bir nechta Bounded Context ustida ishlashi mumkin
- Lekin bir Bounded Context ustida bir nechta jamoa ishlashi **mumkin emas**

**Afzalliklari:**

- Boshqa jamoa sizning kodingizga o'zgartirish kiritishi mumkin emas
- Sizning jamoangiz o'z kodi va ma'lumotlar bazasiga ega
- Rasmiy interfeys orqali Bounded Context'dan foydalaniladi
- Bu DDD dan foydalanishning ustunliklaridan biri

**Muhim:** Manba kodi va ma'lumotlar bazasi sxemasini har bir Bounded Context uchun aniq ajrating. Qabul testlari va test modullarini asosiy manba kodi bilan birga saqlang.