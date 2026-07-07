
`INNER JOIN` ikkala jadvalda **faqat mos keladigan** qatorlarni qaytaradi.

## **Misol jadvallar**

### **1. `products` jadvali**

|product_id|product_name|category_id|
|---|---|---|
|33|Geitost|4|
|34|Sasquatch Ale|1|
|35|Steeleye Stout|1|
|36|Inlagd Sill|8|

### **2. `categories` jadvali**

|category_id|category_name|
|---|---|
|1|Beverages|
|2|Condiments|
|3|Confections|
|4|Dairy Products|

---

## **INNER JOIN qo‘llash**

```sql
SELECT p.product_id, p.product_name, c.category_name
FROM products p
INNER JOIN categories c
ON p.category_id = c.category_id;
```

### **INNER JOIN qanday ishlaydi?**

1. **Chap jadval (`products`) dan `category_id` ni oladi.**
2. **O‘ng jadval (`categories`) dan `category_id` ni oladi.**
3. **category_id lar mos keladigan qatorlarni tanlaydi.**
4. **Faqat mos kelgan qatorlarni natijaga kiritadi.**

---

## **Natija jadvali**

|product_id|product_name|category_name|
|---|---|---|
|33|Geitost|Dairy Products|
|34|Sasquatch Ale|Beverages|
|35|Steeleye Stout|Beverages|

**Nega `Inlagd Sill` natijaga kirmadi?**

- `category_id = 8`, lekin `categories` jadvalida `8` yo‘q.
- `INNER JOIN` faqat mos keladigan qiymatlarni oladi, shuning uchun **o‘sha qator natijada yo‘q**.

---

## **INNER JOIN qanday optimallashtiriladi?**

### **1. Indeks qo‘shish**

Agar `category_id` ustunlari ustida indeks bo‘lsa, PostgreSQL **Index Scan** yordamida tezroq natija topa oladi:

```sql
CREATE INDEX idx_category_id ON products(category_id);
CREATE INDEX idx_category_id2 ON categories(category_id);
```

### **2. EXPLAIN yordamida tekshirish**

Query qanday ishlashini tushunish uchun PostgreSQL-ning `EXPLAIN ANALYZE` komandasi ishlatiladi:

```sql
EXPLAIN ANALYZE
SELECT p.product_id, p.product_name, c.category_name
FROM products p
INNER JOIN categories c
ON p.category_id = c.category_id;
```

Bu natijada PostgreSQL qaysi **JOIN algoritmini** ishlatganini ko‘rsatadi (`Nested Loop`, `Hash Join`, yoki `Merge Join`).

---

## **1. JOIN jarayoni qanday ishlaydi?**

PostgreSQL `INNER JOIN` bajarayotganda quyidagi **asosiy bosqichlarni** bajaradi:

1️⃣ **Har bir jadvalni skan qilish**

- Har ikkala jadvaldan **kerakli qatorlarni topish uchun** `Sequential Scan`, `Index Scan` yoki `Bitmap Index Scan` ishlatiladi.

2️⃣ **JOIN shartini bajarish**

- JOIN shartiga mos keladigan qatorlar bir-biriga **bog‘lanadi**.
- Natija jadvali hosil bo‘ladi.

3️⃣ **Natijadan faqat kerakli ustunlarni olish**

- `SELECT` dagi ustunlar **ajratib olinadi** va natijaga qaytariladi.

## **INNER JOIN jarayoni SQL dvijokida qanday ishlaydi?**

1. **Har bir jadval skanerlash (**`Sequential Scan`, `Index Scan`, `Bitmap Scan`**)**
    
    - SQL dvijoki avval har ikkala jadvaldan ma’lumotlarni olish usulini aniqlaydi.
    - Agar **indeks mavjud bo‘lsa**, **Index Scan** ishlaydi. Agar `ON` shartida ishlatiladigan ustun indeksi mavjud bo‘lsa, PostgreSQL **indeksdan foydalanadi**.
    - Agar indeks bo‘lmasa, **Sequential Scan** ishlaydi. PostgreSQL butun jadvallarni **boshlanishidan oxirigacha** o‘qiydi.
2. Keyin **JOIN shartini tekshiradi** (`ON` qismi).
    
    - Dvijok har bir satrni boshqa jadvaldagi mos keladigan satr bilan solishtiradi.
    - Agar `ON` sharti bajarilsa, satrlar natijaga kiritiladi.

2.1 Ushbu bosqichda **uchta JOIN algoritmidan biri tanlanadi.**

- PostgreSQL ma’lumotlar hajmiga qarab `INNER JOIN`ni amalga oshirish uchun **uchta algoritmdan** birini tanlaydi:
    - **Nested Loop Join** – Agar bitta jadval kichik bo‘lsa.
    - **Hash Join** – Katta hajmli jadvallar uchun.
    - **Merge Join** – Har ikkala jadval **saralangan** bo‘lsa, samarali ishlaydi.