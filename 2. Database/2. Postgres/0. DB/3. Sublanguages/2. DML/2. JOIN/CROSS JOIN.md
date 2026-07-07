### **1. CROSS JOIN nima?**

`CROSS JOIN` – **ikkita jadvaldagi har bir qatorni bir-biriga qo‘shish** uchun ishlatiladi.

Agar **birinchi jadvalda `N` ta qator** va **ikkinchi jadvalda `M` ta qator** bo‘lsa, natijada **`N × M` ta qator** hosil bo‘ladi.

### **2. CROSS JOIN qanday ishlaydi?**

- `CROSS JOIN` ikkita jadval o‘rtasida hech qanday bog‘liqlik yoki `ON` sharti talab qilmaydi.
- Barcha qatorlar har biri bilan kombinatsiya qilinadi.

📌 **Oddiy misol**

```sql
SELECT employees.name, projects.project_name
FROM employees
CROSS JOIN projects;
```

👆 **Natija:**

- Har bir xodim **barcha loyihalar bilan bog‘lanadi**.
- Agar `employees` jadvalida 3 **ta xodim** va `projects` jadvalida 2 **ta loyiha** bo‘lsa, natijada 3 **× 2 = 6 ta qator** hosil bo‘ladi.

---

### **3. CROSS JOIN natijasi qanday shakllanadi?**

### **Misol**:

Ikkita **jadvalimiz bor**:

📌 **Xodimlar (`employees`)**

|id|name|
|---|---|
|1|Ali|
|2|Bobur|
|3|Diyor|

📌 **Loyihalar (`projects`)**

|id|project_name|
|---|---|
|1|Web App|
|2|Mobile App|

✅ **CROSS JOIN natijasi:**

|[employees.name](http://employees.name)|projects.project_name|
|---|---|
|Ali|Web App|
|Ali|Mobile App|
|Bobur|Web App|
|Bobur|Mobile App|
|Diyor|Web App|
|Diyor|Mobile App|

👆 Ko‘rinib turibdiki, har bir xodim **barcha loyihalar bilan bog‘langan**.