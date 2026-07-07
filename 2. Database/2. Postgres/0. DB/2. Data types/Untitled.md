# 1. **Raqamli (Numeric) turlar**

|Type|Tavsif|Misol|
|---|---|---|
|`SMALLINT`|2 bayt, -32,768 dan 32,767 gacha|120|
|`INTEGER` (`INT`)|4 bayt, -2.1 mlrd dan 2.1 mlrd gacha|50000|
|`BIGINT`|8 bayt, juda katta sonlar uchun|9000000000|
|`SERIAL`|Avtoinkrement `INTEGER`|1, 2, 3...|
|`BIGSERIAL`|Avtoinkrement `BIGINT`|1, 2, 3...|
|`DECIMAL(p,s)`|Aniqlik bilan o‘nli sonlar|12.34|
|`NUMERIC(p,s)`|Xuddi `DECIMAL`, ko‘p ishlatiladi|12345.67|
|`REAL`|4 bayt float (floating point)|3.14|
|`DOUBLE PRECISION`|8 bayt float|3.1415926|

---

# 🔤 2. **Matnli (String/Text) turlar**

|Type|Tavsif|Misol|
|---|---|---|
|`CHAR(n)`|Fiks uzunlikdagi matn|'ABC'|
|`VARCHAR(n)`|Maks. uzunlikdagi matn|'PostgreSQL'|
|`TEXT`|Cheksiz uzunlikdagi matn|'Hello'|
|`NAME`|Ob'ekt nomlari uchun maxsus tur|'admin'|

> ✳ `TEXT` — PostgreSQL'da matn uchun eng ko‘p ishlatiladigan va tavsiya etilgan tur.

---

# 📅 3. **Sana va vaqt (Date/Time) turlari**

|Type|Tavsif|Misol|
|---|---|---|
|`DATE`|Faqat sana (yil-oy-kun)|'2025-06-23'|
|`TIME`|Faqat vaqt|'14:30:00'|
|`TIMESTAMP`|Sana + vaqt|'2025-06-23 14:30:00'|
|`TIMESTAMPTZ`|`TIMESTAMP` + vaqt zonasi|'2025-06-23 14:30:00+05'|
|`INTERVAL`|Vaqt farqlari (oraliqlar)|'2 days 3 hours'|

---

# ✅ 4. **Boolean turi**

|Type|Tavsif|Misol|
|---|---|---|
|`BOOLEAN`|`TRUE`, `FALSE`, `NULL` qabul qiladi|`true`|

---

# 📦 5. **Enumeratsiya (ENUM)**

Custom ma'lumot turi – qat'iy qiymatlar to‘plami:

`CREATE TYPE mood AS ENUM ('happy', 'sad', 'angry');`

> Juda foydali agar ma’lumot ma'lum to‘plam bilan cheklangan bo‘lsa.

---

# 🧩 6. **Massiv (ARRAY)**

|Type|Tavsif|Misol|
|---|---|---|
|`INTEGER[]`|Integer massiv|'{1,2,3}'|
|`TEXT[]`|Matn massiv|'{"a","b"}'|

`SELECT ARRAY[1,2,3] AS num_array;`

---

# 📑 7. **JSON va JSONB**

|Type|Tavsif|
|---|---|
|`JSON`|Yengil JSON, saqlaydi|
|`JSONB`|Binary JSON (tezroq va indekslanadi)|

`-- Misol SELECT '{"name": "Ali", "age": 25}'::jsonb;`

---

# 🧬 8. **UUID (Universally Unique ID)**

|Type|Tavsif|
|---|---|
|`UUID`|Unikal identifikator (128 bit)|

`SELECT gen_random_uuid();`

> Modullar orqali qo‘shiladi: `CREATE EXTENSION IF NOT EXISTS "pgcrypto";`

---

# 🧭 9. **Network/Internet turlari**

|Type|Tavsif|Misol|
|---|---|---|
|`INET`|IP-manzil (IPv4/IPv6)|'192.168.1.1'|
|`CIDR`|IP-blok (tarmoqlar uchun)|'192.168.0.0/16'|
|`MACADDR`|MAC manzil|'08:00:2b:01:02:03'|

---

# 📐 10. **Geometrik turlar**

|Type|Tavsif|
|---|---|
|`POINT`|X, Y koordinatalar|
|`LINE`|To‘g‘ri chiziq|
|`CIRCLE`|Aylana|

`SELECT '(1,2)'::POINT;`

---

# 🧪 11. **Money, Bit, Other**

|Type|Tavsif|
|---|---|
|`MONEY`|Pul birligi (ehtiyot bo‘lish kerak)|
|`BIT(n)`|Bit ketma-ketligi (0/1 lar)|
|`BYTEA`|Binary data (rasmlar, fayllar)|
|`TSVECTOR`|Full-text search uchun|

---

# 🧰 12. **Custom (User-defined) Types**

Siz PostgreSQL’da o‘z type'ingizni ham yaratishingiz mumkin:

`CREATE TYPE full_name AS (   first_name TEXT,   last_name TEXT );`

---

# 🧮 Qo‘shimcha: Ma’lumotlar turini aniqlash

`SELECT column_name, data_type FROM information_schema.columns WHERE table_name = 'your_table';`

---
