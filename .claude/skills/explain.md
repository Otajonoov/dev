---
name: explain
description: "Biror tushuncha yoki kodni batafsil o'zbek tilida tushuntiradi, diagrammalar bilan. Ishlatish: /explain <tushuncha yoki kod>"
user_invocable: true
---

# Explain skill

Foydalanuvchi tushuncha yoki kod beradi — sen uni chuqur tushuntirasan. `.claude/agents/teacher.md` dagi o'qitish prinsiplariga amal qil.

## Ish tartibi

1. Mavzuni aniqla
2. Oddiy tildan boshlab, bosqichma-bosqich chuqurlashtirib tushuntir (sodda → murakkab)
3. Har bir bosqichga diagramma va kod misoli qo'sh

## Tushuntirish strukturasi

### 1. Bir jumlada (5 yoshli bolaga tushuntirgandek)
### 2. Analogiya (real hayotdan, chegarasi bilan: "lekin farqi shundaki...")
### 3. Texnik tushuntirish (diagramma bilan)
Kod ortida ASLIDA nima sodir bo'ladi — memory, scheduler, pointer darajasida (notional machine).
### 4. Kod misoli (subgoal izohlar bilan: `// --- 1-qadam: ... ---`, output bilan)
### 5. 🤔 Predict savoli
"Agar shu joyni o'zgartirsak nima bo'ladi?" — javob `<details>` ichida yashirin.
### 6. Qachon/qayerda ishlatiladi (real use case)
### 7. ⚠️ Keng tarqalgan xatolar
Noto'g'ri tasavvur → nega noto'g'ri → to'g'risi qanday.
### 8. ✅ O'z-o'zini tekshir
2-3 ta savol, javoblar `<details>` ichida.

## Qoidalar

- O'zbek tilida, texnik atamalar ingliz tilida
- Kamida 2 ta Mermaid diagramma
- Kamida 2 ta Go kod misoli
- Paragraf 4 qatordan oshmasin
- Agar repo ichida tegishli kod bo'lsa — uni ham ko'rsatib tushuntir
