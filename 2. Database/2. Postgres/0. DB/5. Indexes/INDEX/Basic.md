
Что такое индексы? 
	• Легализованные костыли для ускорения (pk) 

Какова расплата? 
	• Замедление записи в таблицы (80% к 20%) 
	• Дополнительные объемы дискового пространства 
	• Усложненное техническое обслуживание (bloat) 
	

Все ли индексы полезны? 
	• Нет, нет и нет.


---
Indeks - bu jadval bilan bog'langan maxsus ma'lumot tuzilmasi bo'lib, undagi ma'lumotlar asosida yaratiladi. Indekslarni yaratishning asosiy maqsadi - ma'lumotlar bazasining ishlash samaradorligini oshirish.

Jadvallardagi qatorlar tartibsiz holda saqlanadi. SELECT, UPDATE va DELETE operatsiyalarini bajarish paytida DBMS kerakli qatorlarni topishi kerak. Ushbu qidiruvni tezlashtirish uchun indeks yaratiladi. Aslida u quyidagicha tashkil etilgan: jadvalning ma'lum bir qatoridagi ma'lumotlar asosida, ushbu qatorga mos keladigan indeks elementi (yozuv) qiymati shakllantiriladi. Indeks elementi va jadval qatori o'rtasidagi moslikni saqlab qolish uchun har bir elementga qatorga pointer (ko'rsatkich) joylashtiriladi. Indeks tartiblangan struktura hisoblanadi. Undagi elementlar (yozuvlar) saralangan holda saqlanadi, bu esa indeksda ma'lumotlarni qidirishni sezilarli darajada tezlashtiradi. Indeksda kerakli yozuv topilgandan so'ng, DBMS to'g'ridan-to'g'ri havola orqali jadvalning mos qatoriga o'tadi. Indeks yozuvlari jadval qatorlarining bir yoki bir nechta maydonlari qiymatlari asosida shakllantirilishi mumkin. Ushbu maydonlarning qiymatlari turli usullar bilan birlashtirilib va o'zgartirilishi mumkin. Bularning barchasi ma'lumotlar bazasi ishlab chiquvchisi tomonidan indeks yaratish paytida belgilanadi.

Jadvaldagi ma'lum qatorlarni qidirishda DBMS ning planner (rejalashtiruvchi) deb ataladigan maxsus quyi tizimi, ushbu jadval uchun WHERE shartida ko'rsatilgan ustunlar asosida yaratilgan indeks borligini tekshiradi. Agar bunday indeks mavjud bo'lsa, planner uning ushbu aniq holatda ishlatilishi maqsadga muvofiqligini baholaydi. Agar uning ishlatilishi maqsadga muvofiq bo'lsa, avval indeksda kerakli qiymatlar qidiriladi, so'ngra agar bunday qiymatlar topilsa, indeks yozuvlarida saqlangan pointer'lar yordamida jadvalga murojaat qilinadi. Shunday qilib, jadvaldagi qatorlarni to'liq ko'rib chiqish tartiblangan indeksda qidiruv va to'g'ridan-to'g'ri pointer (havola) orqali jadval qatoriga o'tish bilan almashtirilishi mumkin.

Ikki jadvalni birlashtiruvchi JOIN da ishtirok etayotgan ustun bo'yicha yaratilgan indeks jadvallardan yozuvlarni SELECT qilish jarayonini tezlashtirishga yordam berishi mumkin. Saralangan tartibda yozuvlarni SELECT qilishda ham indeks yordam berishi mumkin, agar saralash indeks yaratilgan ustunlar bo'yicha bajariladigan bo'lsa.

Agar SQL query da ORDER BY bo'lsa, indeks tanlangan qatorlarni saralash bosqichidan qochish imkonini berishi mumkin. Biroq, agar SQL query jadvalning katta qismini ko'rib chiqsa, tanlangan qatorlarni aniq saralash indeksdan foydalanishdan tezroq bo'lishi mumkin. Ma'lumotlarga kirishni tezlashtirish maqsadida indekslar yaratishda, yaratilayotgan indeks ishlatilishi kerak bo'lgan tipik query larda tanlanadigan jadval qatorlarining taxminiy ulushi (selektivlik) ni hisobga olish kerak. Agar bu ulush katta bo'lsa (ya'ni selektivlik past), u holda indeksning mavjudligi kutilgan ta'sir bermasligi mumkin. Indekslar jadvaldan faqat qatorlarning kichik ulushi tanlanganida, ya'ni yuqori selektivlikda foydaliroqdir.

ORDER BY ni LIMIT n bilan birgalikda ishlatish holatida aniq saralash (indeks bo'lmasa) birinchi n ta qatorni aniqlash uchun jadvalning barcha qatorlarini qayta ishlashni talab qiladi. Ammo agar ORDER BY bajariladigan ustunlar bo'yicha indeks mavjud bo'lsa, bu birinchi n ta qatorni qolgan qatorlarni umuman skanerlashsiz to'g'ridan-to'g'ri olish mumkin.

Indekslar yaratishda nafaqat indekslangan ustundagi qiymatlarning o'sish tartibi, balki kamayish tartibi ham ishlatilishi mumkin. Default holda tartib o'suvchi, bunda indekslangan ustunlarda mavjud bo'lishi mumkin bo'lgan NULL qiymatlari oxirgi o'rinda turadi. Indeks yaratishda ASC (o'sish tartibi), DESC (kamayish tartibi), NULLS FIRST (bu qiymatlar birinchi o'rinda) va NULLS LAST (bu qiymatlar oxirgi o'rinda) kalit so'zlari yordamida default xatti-harakatni o'zgartirish mumkin. Masalan:
```sql
CREATE INDEX indeks-nomi
ON jadval-nomi ( ustun-nomi NULLS FIRST, ... );

CREATE INDEX indeks-nomi
ON jadval-nomi ( ustun-nomi DESC NULLS LAST, ... );
```


PostgreSQL juda qiziqarli indeks turini qo'llab-quvvatlaydi - qisman indekslar. Bunday indeks jadvalning barcha qatorlari uchun emas, balki faqat ularning qism to'plami uchun shakllantiriladi. Bunga indeks predikati deb ataladigan shartli ifodadan foydalanish orqali erishiladi. Predikat WHERE bandi yordamida kiritiladi.