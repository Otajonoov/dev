---
name: learn
description: "Berilgan matnni tarjima qiladi, ilmiy asoslangan pedagogika bilan tushuntiradi, diagrammalar bilan boyitadi va learning/ papkasiga saqlaydi. Ishlatish: /learn <kategoriya> <mavzu>"
user_invocable: true
---

# Learn skill

Foydalanuvchi matn beradi — sen uni o'rganish materiali sifatida qayta ishlaysan. Bu ishni **teacher agent uslubida** bajar: `.claude/agents/teacher.md` dagi 6 ilmiy prinsip va 7-qadam tushuntirish formulasiga to'liq amal qil.

## Ish tartibi

1. **Kategoriyani aniqla** — golang, database, network, algorithms, system-design, va h.k.
2. **Tartib raqamini aniqla** — `learning/{kategoriya}/` ichidagi mavjud fayllarni ko'r, keyingi raqamni qo'y
3. **Bog'liq mavzularni top** — `learning/` dagi mavjud fayllarni ko'zdan kechir; yangi mavzuni ular bilan bog'lash uchun (elaboration)
4. **Matnni qayta ishla:**
   - To'liq o'zbek tiliga tarjima qil (texnik atamalar ingliz tilida)
   - Har bir tushunchani 7-qadam formula bilan tushuntir: muammo → analogiya → sodda ta'rif → diagramma → worked example (subgoal izohlar bilan) → predict savoli → ko'p uchraydigan xatolar
   - Go kod misollarini qo'sh (o'zbekcha izohlar bilan, output bilan)
5. **Faylga saqlash** — `learning/{kategoriya}/XX-mavzu-nomi.md`
6. **README yangilash** — `learning/README.md` ga yangi mavzuni qo'sh

## Fayl strukturasi

```markdown
# Mavzu nomi

> Bir jumlada ta'rif

## Nega bu kerak?
Bu tushuncha bo'lmasa qanday muammo paydo bo'ladi — real vaziyat bilan.
Oldingi o'rganilgan mavzularga bog'lanish: "X mavzusini eslaysanmi? Bu unga o'xshaydi, faqat..."

## Asosiy tushunchalar

### Tushuncha 1
1. Analogiya (real hayotdan, chegarasi bilan)
2. Sodda ta'rif (bir jumla)
3. Mermaid diagramma
4. Go kod misoli — subgoal izohlar bilan (`// --- 1-qadam: ... ---`), output bilan
5. 🤔 Predict savoli — javob <details> ichida
6. ⚠️ Ko'p uchraydigan xato

## Xulosa
- Muhim fikrlar ro'yxati (5-7 ta)

## 🧠 Eslab qol
- Eng muhim 3-5 nuqta

## ✅ O'z-o'zini tekshir
3-5 ta savol ("nima bo'ladi agar...", "nega...", "farqi nima...").
Har javob <details><summary>💡 Javob</summary>...</details> ichida.

## 🛠 Amaliyot
1. Oson — mavjud kodni o'zgartirish (Modify)
2. O'rta — yarim tayyor skeleton'ni to'ldirish (`// TODO:` bilan)
3. Qiyin — noldan yozish (Make)
Har biriga hint <details> ichida.

## 🔁 Takrorlash
- Bog'liq mavzular: [[oldingi mavzu fayllari]]
- Jadval: ertaga → 3 kun → 1 hafta ("O'z-o'zini tekshir" ga qaytish)
```

## Qoidalar

- Hech narsa tashlab ketma — to'liq tarjima
- Texnik atama = ingliz tilida qoladi
- Har bir tushunchaga kamida 1 ta diagramma va 1 ta kod misoli
- Paragraf 4 qatordan oshmasin — matn devori taqiqlanadi
- Javoblar va hintlar HAR DOIM <details> ichida yashirin bo'lsin — foydalanuvchi avval o'zi o'ylab ko'rsin (retrieval practice)
