
ПОДОБЛАСТИ (SUBDOMAIN)

### Subdomain nima?

Oddiy qilib aytganda, **Subdomain** — bu sizning domain'ingizning bir qismi. Subdomain'ni ma'lum bir domain'ning alohida mantiqiy modeli sifatida tasavvur qilishingiz mumkin. Ko'pgina domain'lar odatda juda katta va murakkab bo'ladi, shuning uchun biz faqat muayyan loyihada foydalanish kerak bo'lgan Subdomain'lar bilan qiziqamiz.

Subdomain'ni talqin qilishning yana bir usuli — uni biznesingizning **Core Domain**'i uchun yechim topishga imkon beruvchi aniq belgilangan kompetensiya sohasi sifatida ko'rishdir. Bu ma'lum Subdomain'da bitta yoki bir nechta **domain ekspertlari** mavjudligini anglatadi.

### Subdomain'larning turlari

Loyihada uchta asosiy Subdomain turi mavjud:

**1. Core Domain (Asosiy yadro)**

- Bu strategik investitsiya sohasi
- Aniq belgilangan **Bounded Context**'da **Ubiquitous Language** yaratishga katta resurslar sarflanadi
- Tashkilot uchun eng yuqori ustuvorlikka ega
- Sizga raqobatdosh ustunlik beradi
- Bu yerga eng ko'p investitsiya kiritilishi kerak

**2. Supporting Subdomain (Yordamchi Subdomain)**

- Maxsus modellashtirish talab qiladi, ammo tayyor yechimlar mavjud emas
- Core Domain qadar katta investitsiya talab qilmaydi
- Tashqi tashkilotlarga topshirilishi mumkin
- Core Domain'siz muvaffaqiyatli ishlay olmaydi

**3. Generic Subdomain (Universal Subdomain)**

- Tayyor sotib olinishi mumkin
- Tashqi tashkilotlarga yoki alohida bo'linmaga topshiriladi
- Elite dasturchilar bu yerda ishlamamasligi kerak
- Katta investitsiya talab qilmaydi
- Generic Subdomain'ni Core Domain deb adashmaslik kerak

Legacy tizimlar ko'pincha **Big Ball of Mud** (Iflos loy to'pi) deb ataladigan narsani tashkil qiladi. Bunday tizimlarda bir nechta mantiqiy model aralashib ketgan va ularni ajratish juda qiyin. Har bir mantiqiy modelni Subdomain sifatida ko'rib chiqish mumkin.

**Muhim qoida:** DDD yondashuvida Bounded Context va Subdomain o'rtasida bir-birga mos kelish bo'lishi kerak. Ya'ni, bitta Bounded Context'da bitta Subdomain modeli bo'lishi kerak.