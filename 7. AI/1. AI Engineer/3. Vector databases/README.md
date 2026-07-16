# 3. Vector databases

2-bo'limda embedding'lar numpy massivda yashadi — bu o'rganish uchun to'g'ri, production uchun emas.
Bu bo'lim vektorlarni haqiqiy DB'ga ko'chiradi: siz allaqachon chuqur bilgan **Postgres** (pgvector)
asosiy qurol bo'ladi, **Qdrant** esa "qachon dedicated vector DB kerak" savoliga javob beradi.
Bo'lim tugagach: ANN indexlarni (IVFFlat, HNSW) ongli sozlash, recall o'lchash, hybrid search (RRF),
filtered search tuzoqlari va scale eskalatsiya zinapoyasi — hammasi qo'lda.

## Kim uchun

Middle backend dasturchi. 1-2 bo'limlar tugagan, Postgres chuqur tanish (EXPLAIN, B-tree/GIN,
planner, partitioning). Markaziy tayanch: **vector index — bu ham index**, faqat exact emas,
approximate: tezlik evaziga recall sotib olasiz va bu savdoni parametrlar bilan boshqarasiz.

## Darslar

| # | Dars | Nima o'rganasiz |
|---|---|---|
| 01 | [pgvector — Postgres ichida vector search](./01.%20pgvector%20—%20Postgres%20ichida%20vector%20search.md) | vector tipi, distance operatorlari (`<->` `<=>` `<#>`), exact kNN, nega B-tree yordam bermaydi |
| 02 | [Index turlari — IVFFlat va HNSW](./02.%20Index%20turlari%20—%20IVFFlat%20va%20HNSW.md) | Ikkala indexning mexanikasi, lists/probes, m/ef_construction/ef_search, recall o'lchash, halfvec |
| 03 | [Qdrant — dedicated vector DB qachon kerak](./03.%20Qdrant%20—%20dedicated%20vector%20DB%20qachon%20kerak.md) | pgvector chegarasi signallari, Qdrant collections/payload/`query_points`, qaror daraxti |
| 04 | [Hybrid search — full-text + vector, RRF](./04.%20Hybrid%20search%20—%20full-text%20+%20vector,%20RRF.md) | Term-based vs embedding-based, tsvector+GIN, RRF SQL, Qdrant prefetch+fusion |
| 05 | [Metadata filtering va scale](./05.%20Metadata%20filtering%20va%20scale.md) | Pre/post-filtering dilemmasi, iterative index scans, partitioning, quantization, RAM hisob-kitobi |
| 06 | [Bo'lim loyihasi — pgvector qidiruv servisi](./06.%20Bo'lim%20loyihasi%20—%20pgvector%20qidiruv%20servisi.md) | `semsearch` indexini pgvector'ga ko'chirish + FastAPI HTTP API + Docker compose |

## Asosiy qarorlar (2026 holati)

- **pgvector 0.8.5** — asosiy DB (Docker: `pgvector/pgvector:pg18-trixie`). 2026 konsensus:
  Postgres allaqachon stack'da bo'lsa va vektorlar ~5-10M dan oshmasa — pgvector; undan nari —
  avval pgvectorscale, keyin dedicated DB.
- **Qdrant** — dedicated vector DB vakili (kitob — Iusztin & Labonne — ham shu tanlagan):
  filterable HNSW, built-in quantization, gorizontal scale. Faqat joriy `query_points()` API.
- **Embeddings standarti o'zgarmagan:** `voyage-4` (1024 dim, normalizatsiyalangan → `<#>` eng arzon),
  `input_type` majburiy. Kitoblardagi FAISS misollari — kutubxona, DB emas: production'da
  CRUD/filter/backup qatlami kerak, shuning uchun bu bo'limda pgvector/Qdrant.

## Kod muhiti

```bash
python -m venv .venv && source .venv/bin/activate
pip install "psycopg[binary]" psycopg-pool pgvector voyageai python-dotenv fastapi uvicorn qdrant-client

# Postgres + pgvector
docker run -d --name pgv -e POSTGRES_PASSWORD=secret -p 5432:5432 pgvector/pgvector:pg18-trixie

# Qdrant (03-04 darslar uchun)
docker run -d --name qdrant -p 6333:6333 -v "$(pwd)/qdrant_storage:/qdrant/storage" qdrant/qdrant
```

`.env`: `VOYAGE_API_KEY=pa-...` va `DATABASE_URL=postgresql://postgres:secret@localhost:5432/postgres`.

## Bo'lim tugagach

`06. Bo'lim loyihasi`dagi qidiruv servisi Docker compose bilan ko'tarilib turishi kerak.
Portfolio zanjiri: `askops` + `quizgen` (1-bo'lim) → `semsearch` (2-bo'lim) →
**pgvector servisi (shu bo'lim)** → RAG savol-javob tizimi (4-bo'lim).
