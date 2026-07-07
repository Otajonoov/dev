`ORDER BY` natijalarni **katta-kichik** yoki **alfavit** tartibida saralaydi.

```sql
SELECT product, SUM(amount) AS total
FROM sales
GROUP BY product
ORDER BY total DESC;
```

|product|total|
|---|---|
|A|250|
|B|250|
