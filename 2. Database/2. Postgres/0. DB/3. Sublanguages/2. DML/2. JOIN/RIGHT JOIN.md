`RIGHT JOIN` **o‘ng (RIGHT) jadvaldagi barcha qatorlarni saqlaydi** va **chap (LEFT) jadvaldagi mos keladigan qatorlarni bog‘laydi**. Agar chap jadvalda mos keladigan qator topilmasa, `NULL` qiymatlar qaytariladi.

## **1. RIGHT JOIN qanday ishlaydi?**

### **SQL sintaksisi**:

```sql
SELECT t1.*, t2.*
FROM table1 AS t1
RIGHT JOIN table2 AS t2
ON t1.id = t2.id;
```

📌 **Natija quyidagicha bo‘ladi**:

- **O‘ng jadvaldagi barcha qatorlar saqlanadi**.
- **Agar `ON` shartiga mos keladigan chap jadval qatori bo‘lsa, bog‘lanadi**.
- **Agar mos keladigan chap jadval qatori bo‘lmasa, chap jadval ustunlari `NULL` bo‘ladi**.

---

## **2. PostgreSQL dvijokining RIGHT JOIN bajarish bosqichlari** 🛠

### **1️⃣ Chap va o‘ng jadvallarni skan qilish**

PostgreSQL **har ikkala jadvalni qanday o‘qish kerakligini** aniqlaydi.

- **Skan usullari**:
    - **Sequential Scan** – Agar jadvalda indeks bo‘lmasa.
    - **Index Scan** – Agar `ON` shartidagi ustun indekslangan bo‘lsa.
    - **Bitmap Index Scan** – Indekslangan ma’lumotlarni samarali ishlatish uchun.

### **2️⃣ JOIN shartini bajarish (`ON` qismi)**

- O‘ng jadvaldagi **har bir qator** chap jadvaldagi mos keladigan qatorlar bilan **taqqoslanadi**.
- Agar mos keladigan qator **bo‘lsa**, ikkala jadval qatori bog‘lanadi.
- Agar mos keladigan qator **bo‘lmasa, chap jadval ustunlari `NULL` bo‘ladi**.

### **3️⃣ JOIN algoritmini tanlash**

PostgreSQL **ma’lumot hajmiga qarab** quyidagi JOIN algoritmlaridan birini ishlatadi:

1️⃣ **Nested Loop Join** (Kichik jadvallar uchun)

- Agar chap jadval kichik bo‘lsa va indeks mavjud bo‘lsa, samarali ishlaydi.
- Har bir o‘ng qator uchun chap jadvalni skan qiladi.
- **O(n × m) murakkablik**, ya’ni sekin ishlashi mumkin.

2️⃣ **Hash Join** (Katta jadvallar uchun)

- Chap jadval uchun **hash jadval** yaratiladi.
- O‘ng jadval skan qilinadi va har bir satr **hash jadvalga qarab moslashtiriladi**.
- **O(n + m) murakkablik**, bu katta jadvallar uchun samarali.

3️⃣ **Merge Join** (Saralangan jadvallar uchun)

- Agar ikkala jadval **saralangan bo‘lsa**, `Merge Join` ishlaydi.
- Har ikkala jadval **saralangan tartibda** o‘qilib, `ON` sharti bo‘yicha solishtiriladi.
- **O(n + m) murakkablik**, saralangan ma’lumotlar bilan eng tezkor usul.

### **4️⃣ Natijani shakllantirish**

- **O‘ng jadvaldagi barcha qatorlar natijaga kiritiladi**.
- **Mos keladigan chap jadval qatorlari qo‘shiladi**.
- **Mos kelmaydigan qatorlarda chap ustunlar `NULL` bo‘ladi**.

---

## **3. RIGHT JOIN misollar bilan tushuntirish** 📌

📌 **table1 (chap jadval)**

|id|name|
|---|---|
|1|Alice|
|2|Bob|

📌 **table2 (o‘ng jadval)**

|id|salary|
|---|---|
|1|5000|
|2|6000|
|3|7000|

📌 **RIGHT JOIN natijasi:**

```sql
SELECT t1.*, t2.salary
FROM table1 AS t1
RIGHT JOIN table2 AS t2
ON t1.id = t2.id;
```

|id|name|salary|
|---|---|---|
|1|Alice|5000|
|2|Bob|6000|
|NULL|NULL|7000|

🛑 **id = 3 uchun `NULL`**, chunki `table1` jadvalida `id = 3` yo‘q