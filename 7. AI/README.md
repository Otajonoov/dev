# AI yo'li

Go backend fonidan AI sohasiga o'tish uchun umumiy papka. Ikkita kasb yo'lini o'z ichiga oladi — ikkalasi ham bitta poydevorga tayanadi, lekin maqsadi va muddati farq qiladi.

## Ikki yo'l

```
7. AI/
├── 1. AI Engineer/    ← tezkor yo'l (2–4 oy): tayyor LLM'larni mahsulotga ulash
└── 2. ML Engineer/    ← chuqur yo'l (1–2 yil): model qurish, o'rgatish, deploy qilish
```

| | AI Engineer | ML Engineer |
| --- | --- | --- |
| **Maqsad** | LLM API, RAG, agent'lar bilan mahsulot qurish | Model qurish va production'da yuritish |
| **Talab** | Python + backend tajriba yetadi | + Matematika, klassik ML, Deep learning |
| **Muddat** | 2–4 oy | 1–2 yil |
| **Bozor** | O'zbekistonda hozir eng ko'p so'raladigan kombinatsiya (backend + AI) | Chuqurlashish va katta kompaniyalar uchun |

## Tavsiya etilgan tartib

1. **Python** — `2. ML Engineer/1. Python ecosystem/1. Python` (Basics + Advanced tayyor)
2. **1. AI Engineer** — to'liq yo'l, portfolio loyihalar bilan → ish topish
3. **2. ML Engineer** — ishlab yurgan holda chuqurlashish (Matematika → ML → Deep learning → ...)

Ya'ni AI Engineer — kirish nuqtasi, ML Engineer — uzoq muddatli chuqurlik. Ikkala yo'l bir-birini inkor qilmaydi: AI Engineer'da o'rganilgan RAG/embeddings amaliyoti ML Engineer'dagi `5. LLM internals` (Transformers, Fine-tuning) nazariyasi bilan to'ldiriladi.

## Umumiy poydevor

Python ecosystem, Database (Postgres, Redis), Algorithm — `2. ML Engineer` ichida turadi (1, 9, 10), lekin ikkala yo'lga birdek xizmat qiladi. Repo root'dagi Linux, DevOps, System Design, Network ham shu yo'lning infratuzilma qismini yopadi.

## Til va format

- Tushuntirishlar o'zbek tilida, texnik atamalar ingliz tilida
- Diagrammalar — Mermaid formatda
