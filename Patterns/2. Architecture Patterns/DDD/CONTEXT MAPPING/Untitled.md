
СВЯЗЫВАНИЕ КОНТЕКСТОВ (CONTEXT MAPPING).

### Context Mapping nima?

**Context Mapping** — bu turli Bounded Context'lar o'rtasidagi integratsiya va munosabatlarni tavsiflaydi. Ikki Bounded Context o'rtasidagi chiziq Context Map'ni yaratadi va ular o'rtasida dinamik o'zaro ta'sir borligini ko'rsatadi.

Ikki turli Bounded Context'da ikki xil Ubiquitous Language ishlatilgani uchun, bu chiziqni ikki til o'rtasidagi tarjima sifatida tushunish mumkin.

### Context Mapping turlari

**1. Partnership (Hamkorlik)**

- Ikki jamoa o'z maqsadlarini muvofiqlashtirish uchun hamkorlik quradilar
- Yoki ikkalasi ham muvaffaqiyatga erishadi, yoki ikkalasi ham mag'lub bo'ladi
- Jamoalar tez-tez uchrashib, kalendar rejalarini sinxronlashtiradilar
- Bu munosabatlarni uzoq vaqt davomida saqlab qolish qiyin bo'lishi mumkin

**2. Shared Kernel (Umumiy yadro)**

- Ikki yoki undan ortiq jamoa kichik, lekin umumiy modeldan foydalanadi
- Jamoalar qaysi model elementlarini birgalikda ishlatishlarini kelishib olishlari kerak
- Ochiq muloqot va doimiy kelishuv talab qiladi

**3. Customer-Supplier (Mijoz-Yetkazib beruvchi)**

- Supplier yuqori oqim konteksti (upstream - U)
- Customer pastki oqim konteksti (downstream - D)
- Supplier bu munosabatlarda hukmronlik qiladi
- Supplier nima va qachon berishni belgilaydi

**4. Conformist (Moslashuvchi)**

- Yuqori oqim jamoasi pastki oqim ehtiyojlarini qondirish uchun rag'batlantirmaydi
- Pastki oqim jamoasi Ubiquitous Language tarjimasini saqlay olmaydi
- Shuning uchun yuqori oqim modeliga moslashadi
- Masalan, Amazon.com bilan integratsiya qilmoqchi bo'lgan kompaniyalar ko'pincha Conformist bo'ladi

**5. Anticorruption Layer (Himoya qatlami)**

- Pastki oqim jamoasi o'z Ubiquitous Language va yuqori oqim tili o'rtasida tarjima qatlamini yaratadi
- Bu qatlam pastki oqim modelini yuqori oqimdan izolyatsiya qiladi
- Har doim imkon qadar Anticorruption Layer yaratishga harakat qilish kerak

**6. Open Host Service (Ochiq xost xizmati)**

- Bounded Context'ga kirish uchun protokol yoki interfeys taqdim etadi
- Protokol ochiq, shuning uchun integratsiya qilish nisbatan oson
- API xizmatlari yaxshi hujjatlashtirilgan va foydalanish uchun qulay

**7. Published Language (Umumiy til)**

- Yaxshi hujjatlashtirilgan ma'lumot almashish tili
- Istalgan miqdordagi Bounded Context'lar uchun oson foydalanish va tarjima imkonini beradi
- XML Schema, JSON Schema, Protobuf yoki Avro yordamida aniqlanishi mumkin

**8. Separate Ways (Alohida yo'llar)**

- Bir yoki bir nechta Bounded Context bilan integratsiya katta foyda keltirmaydi
- O'z maxsus yechimingizni o'z Bounded Context'ingizda ishlab chiqasiz

**9. Big Ball of Mud (Iflos loy to'pi)**

- Bunday tizim yaratishdan vabo kasalligidek qochish kerak
- Agar bunday tizim bilan ishlashga majbur bo'lsangiz, Anticorruption Layer yarating
- Bu tilni hech qachon gapirmang!

---
