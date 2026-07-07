---
name: learn
description: "Berilgan matnni tarjima qiladi, tushuntiradi, diagrammalar bilan boyitadi va learning/ papkasiga saqlaydi. Ishlatish: /learn <kategoriya> <mavzu>"
user_invocable: true
---

# Learn skill

Foydalanuvchi matn beradi — sen uni o'rganish materiali sifatida qayta ishlaysan.

## Ish tartibi

1. **Kategoriyani aniqla** — golang, database, network, algorithm, system-design, va h.k.
2. **Tartib raqamini aniqla** — `learning/{kategoriya}/` ichidagi mavjud fayllarni ko'r, keyingi raqamni qo'y
3. **Matnni qayta ishla:**
   - To'liq o'zbek tiliga tarjima qil
   - Har bir bo'limga Mermaid diagramma qo'sh
   - Go kod misollarini qo'sh (o'zbekcha izohlar bilan)
   - Analogiyalar va real hayotdan misollar qo'sh
4. **Faylga saqlash** — `learning/{kategoriya}/XX-mavzu-nomi.md`
5. **README yangilash** — `learning/README.md` ga yangi mavzuni qo'sh

## Fayl strukturasi

```markdown
# Mavzu nomi

> Bir jumlada ta'rif

## Kirish
Nima uchun bu mavzu muhim, qayerda ishlatiladi

## Asosiy tushunchalar

### Tushuncha 1
- Oddiy tilda tushuntirish
- Analogiya
- Mermaid diagramma
- Go kod misoli

## Xulosa
- Muhim fikrlar ro'yxati

## Eslab qol
- Eng muhim 3-5 ta nuqta

## Amaliyot
1. Topshiriq (oson)
2. Topshiriq (o'rta)
3. Topshiriq (qiyin)
```

## Qoidalar

- Hech narsa tashlab ketma — to'liq tarjima
- Texnik atama = o'zbekcha (inglizcha)
- Har bir tushunchaga kamida 1 ta diagramma
- Har bir tushunchaga kamida 1 ta kod misoli
