
 `GROUP BY`
- GROUP BY agregatsiya qilish uchun kerak bo’lgan ustunni aniqlaydi. Ustunlar bo’yicha guruhlab, har bir guruh uchun natijalar chiqariladi.
- GROUP BY bir xil qiymatlarga ega bo‘lgan satrlarni _bitta guruh_ qiladi. Har bir guruh uchun siz **AGGREGATE FUNCTION** (masalan: `SUM`, `AVG`, `COUNT`, `MIN`, `MAX`) ishlatishingiz mumkin.

```sql
-- Har bir product_id bo‘yicha sotilgan umumiy miqdor 
SELECT product_id, SUM(quantity)  FROM orders  GROUP BY product_id;
```

izda quyidagi `sales` jadvali bor:

|id|product|amount|
|---|---|---|
|1|A|100|
|2|B|200|
|3|A|150|
|4|B|50|

#### Siz har bir mahsulotdan qancha sotilganini bilmoqchisiz:

```sql
SELECT product, SUM(amount) FROM sales GROUP BY product;
```

###  Natija:

| product | sum |
| ------- | --- |
| A       | 250 |
| B       | 250 |

GROUP BY product — `A` larni alohida, `B` larni alohida guruhlab, har biri uchun `SUM(amount)` ishlatadi.