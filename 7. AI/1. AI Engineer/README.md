# AI Engineer yo'li

Tayyor LLM'larni (Claude, OpenAI, Gemini) real mahsulotlarga integratsiya qilish kasbi. Model qurilmaydi — API orqali ishlatiladi, shuning uchun backend tajriba + Python bilan 2–4 oyda ishga yaroqli darajaga chiqish mumkin. O'zbekiston bozorida (bank, telekom, e-commerce) aynan shu kombinatsiya so'raladi.

## Boshlashdan oldin

- **Python** — `../2. ML Engineer/1. Python ecosystem/1. Python` (Basics + Advanced tayyor)
- **Postgres** — `../2. ML Engineer/9. Database` (pgvector shu bilimga tayanadi)

## Struktura va tartib

```
1. LLM API va Prompt Engineering/  ← boshlanish nuqtasi
   API bilan ishlash (Claude, OpenAI), streaming, tool use,
   structured output, tizimli prompt yozish
   → loyiha: CLI chatbot + savol generatsiya skripti

2. Embeddings/
   Matnni vektorga aylantirish, semantic similarity,
   embedding modellari va ularni tanlash
   → loyiha: semantic search CLI

3. Vector databases/
   pgvector, Qdrant — indekslash (IVFFlat, HNSW), hybrid search
   (Postgres bilimiga bog'lanadi)
   → loyiha: pgvector qidiruv servisi

4. RAG/
   Retrieval-Augmented Generation — chunking, retrieval,
   reranking, advanced RAG. Ish e'lonlarida №1 ko'nikma
   → loyiha: hujjat savol-javob tizimi

5. AI Agents/
   Tool use, agent loop, MCP, multi-agent patterns,
   agent xavfsizligi. 2025–2026 ning eng tez o'sayotgan sohasi
   → loyiha: mustaqil task agent

6. Evaluation/
   Golden dataset, LLM-as-judge, offline/online eval,
   regression testing — eval'siz LLM feature deploy qilinmaydi
   → loyiha: RAG tizimi uchun eval harness

7. Production/
   Serving (FastAPI + SSE), semantic caching, cost/latency
   optimization, observability, guardrails — backend
   tajribaning raqobat ustunligi aynan shu yerda
   → loyiha: production RAG + Telegram bot (yakuniy portfolio)
```

Alohida "Amaliyot" papkasi yo'q — har mavzuning nazariyasi va amaliyoti o'z papkasi ichida yonma-yon turadi, har bo'lim o'zining loyihasi bilan tugaydi (loyihalar zanjiri portfolio hosil qiladi). Kursni qurish rejasi va yakuniy struktura — `TASK.md` da.

## Chegara: bu yerda nima YO'Q

Model ichki tuzilishi (Transformers arxitekturasi) va fine-tuning — `../2. ML Engineer/5. LLM internals` da. AI Engineer uchun ular shart emas, lekin chuqurlashishda keyingi qadam.

## Til va format

- Tushuntirishlar o'zbek tilida, texnik atamalar ingliz tilida
- Diagrammalar — Mermaid formatda
