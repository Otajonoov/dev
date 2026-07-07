`FULL JOIN` yoki `FULL OUTER JOIN` – bu **chap (LEFT) va o‘ng (RIGHT) jadvallarning hamma qatorlarini chiqaradigan** JOIN turi. Agar biror satr **boshqa jadvalda topilmasa**, u holda `NULL` qo‘yiladi.

## **2. FULL JOIN qanday ishlaydi?**

🔹 `FULL JOIN` – bu **LEFT JOIN + RIGHT JOIN** kombinatsiyasi.

🔹 Chap jadvaldan ham, o‘ng jadvaldan ham barcha ma’lumotlar chiqadi.

🔹 Agar mos kelmagan qator bo‘lsa, `NULL` bilan to‘ldiriladi.

📌 **Sintaksis**:

```sql
sql
CopyEdit
SELECT *
FROM jadval1
FULL JOIN jadval2 ON jadval1.id = jadval2.id;

```

---

## **3. FULL JOIN misol bilan tushuntirish**

📌 **Jadvallar**:

Xodimlar (`employees`) va Maoshlar (`salaries`) jadvallarimiz bor.

**`employees` (Xodimlar) jadvali:**

|id|name|
|---|---|
|1|Ali|
|2|Bobur|
|3|Diyor|
|4|Elmurod|

**`salaries` (Maoshlar) jadvali:**

|id|amount|
|---|---|
|2|1500|
|3|2000|
|5|2500|

📌 **FULL JOIN so‘rovi**:

```sql
sql
CopyEdit
SELECT employees.id, employees.name, salaries.amount
FROM employees
FULL JOIN salaries ON employees.id = salaries.id;

```

**Natija (FULL JOIN):**

|id|name|amount|
|---|---|---|
|1|Ali|NULL|
|2|Bobur|1500|
|3|Diyor|2000|
|4|Elmurod|NULL|
|5|NULL|2500|

**Tushuntirish:**

- **Bobur va Diyor** – maosh bor, shuning uchun normal chiqadi. ✅
- **Ali va Elmurod** – maoshi yo‘q (`NULL` bo‘ldi). ⚠️
- **ID = 5 bo‘lgan maosh** bor, lekin xodim yo‘q (`NULL`). ⚠️

---

## **4. FULL JOIN qanday ishlatiladi?**

### **1️⃣ Barcha ma’lumotlarni olish**

Ba’zan **LEFT JOIN yoki RIGHT JOIN yetarli bo‘lmaydi**, chunki ikkala jadvaldagi barcha ma’lumotlarni olish kerak.

**Misol:** Barcha xodimlar va barcha maoshlarni chiqarish.

```sql
sql
CopyEdit
SELECT employees.id, employees.name, salaries.amount
FROM employees
FULL JOIN salaries ON employees.id = salaries.id;

```

📌 **Foyda**: Barcha ma’lumotlar chiqadi.

---

### **2️⃣ Mos kelmagan qatorlarni ajratib olish**

Ba’zan **faqat mos kelmagan qatorlarni topish kerak** (xodim bor, lekin maoshi yo‘q yoki maoshi bor, lekin xodim yo‘q).

📌 **Mos kelmagan qatorlarni topish**:

```sql
sql
CopyEdit
SELECT employees.id, employees.name, salaries.amount
FROM employees
FULL JOIN salaries ON employees.id = salaries.id
WHERE employees.id IS NULL OR salaries.id IS NULL;

```

**Natija:**

|id|name|amount|
|---|---|---|
|1|Ali|NULL|
|4|Elmurod|NULL|
|5|NULL|2500|

👆 **Foyda**:

- Xodimi bor, lekin maoshi yo‘q (Ali, Elmurod).
- Maoshi bor, lekin xodim yo‘q (`ID = 5`)

**FULL JOIN – bu**

- **LEFT JOIN + RIGHT JOIN birgalikda ishlashi**.
- **Chap va o‘ng jadvaldagi barcha qatorlarni chiqaradi**.
- **Mos kelmagan joylarda `NULL` bo‘ladi**.