`LEFT JOIN` **asosiy (chap) jadvaldagi barcha qatorlarni** saqlab qoladi va ularga **o‘ng jadvaldagi mos keladigan qatorlarni bog‘laydi**. Agar mos keladigan qator topilmasa, `NULL` qiymatlar qaytariladi.

## **PostgreSQL dvijokining LEFT JOIN bajarish bosqichlari** 🛠

### **1️⃣ Chap va o‘ng jadvallarni skan qilish**

- PostgreSQL **avval har ikkala jadvalni o‘qish usulini** tanlaydi.
- **Skan usullari**:
    - **Sequential Scan** – Agar jadvalda indeks bo‘lmasa.
    - **Index Scan** – Agar `ON` shartidagi ustun indekslangan bo‘lsa.
    - **Bitmap Index Scan** – Katta hajmli ma’lumotlar bilan ishlashda indekslardan foydalanadi.

### **2️⃣ JOIN shartini bajarish (`ON` qismi)**

- Chap jadvaldagi **har bir qator** o‘ng jadvaldagi mos keladigan qatorlar bilan **taqqoslanadi**.
- Agar mos keladigan qator topilsa, **ikkala jadval qatori qo‘shiladi**.
- Agar mos keladigan qator **topilmasa, o‘ng jadval ustunlari `NULL` bo‘ladi**.

### **3️⃣ JOIN algoritmini tanlash**

PostgreSQL **ma’lumot hajmiga qarab** quyidagi JOIN algoritmlaridan birini ishlatadi:

1️⃣ **Nested Loop Join** (Kichik jadvallar uchun)

- Agar chap jadval kichik bo‘lsa va indeks mavjud bo‘lsa, samarali ishlaydi.
- Har bir chap qator uchun o‘ng jadvalni skan qiladi.
- **O(n × m) murakkablik**, ya’ni sekin ishlashi mumkin.

2️⃣ **Hash Join** (Katta jadvallar uchun)

- O‘ng jadval uchun **hash jadval** yaratiladi.
- Chap jadval skan qilinadi va har bir satr **hash jadvalga qarab moslashtiriladi**.
- **O(n + m) murakkablik**, bu katta jadvallar uchun samarali.

3️⃣ **Merge Join** (Saralangan jadvallar uchun)

- Agar ikkala jadval **saralangan bo‘lsa**, `Merge Join` ishlaydi.
- Har ikkala jadval **saralangan tartibda** o‘qilib, `ON` sharti bo‘yicha solishtiriladi.
- **O(n + m) murakkablik**, saralangan ma’lumotlar bilan eng tezkor usul.

### **4️⃣ Natijani shakllantirish**

- **Chap jadvaldagi barcha qatorlar natijaga kiritiladi**.
- **Mos keladigan o‘ng jadval qatorlari qo‘shiladi**.
- **Mos kelmaydigan qatorlarda o‘ng ustunlar `NULL` bo‘ladi**.

---

## **3. LEFT JOIN misollar bilan tushuntirish** 📌

**Misol jadvallar:**

📌 **table1 (chap jadval)**

|id|name|
|---|---|
|1|Alice|
|2|Bob|
|3|Charlie|

📌 **table2 (o‘ng jadval)**

|id|salary|
|---|---|
|1|5000|
|2|6000|

📌 **LEFT JOIN natijasi:**

```sql
SELECT t1.*, t2.salary
FROM table1 AS t1
LEFT JOIN table2 AS t2
ON t1.id = t2.id;
```

|id|name|salary|
|---|---|---|
|1|Alice|5000|
|2|Bob|6000|
|3|Charlie|NULL|

🛑 **Charlie uchun `NULL`**, chunki `table2` jadvalida `id = 3` yo‘q.

---

---