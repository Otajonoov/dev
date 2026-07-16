# 2. Embeddings

Semantic search'ning poydevori. Bu bo'lim tugagach siz matnni vektorga aylantirish, ma'no bo'yicha
qidirish, embedding model tanlash va matnni to'g'ri tayyorlashni bilasiz — 3-bo'lim (vector databases)
va 4-bo'lim (RAG) aynan shu poydevor ustiga quriladi.

## Kim uchun

Middle backend dasturchi. 1-bo'lim (LLM API va Prompt Engineering) tugagan deb hisoblanadi.
Markaziy analogiya backend dunyosidan: **hash = aniq moslik, embedding = semantik moslik**.

## Darslar

| # | Dars | Nima o'rganasiz |
|---|---|---|
| 01 | [Embedding nima — vektor va semantic similarity](./01.%20Embedding%20nima%20—%20vektor%20va%20semantic%20similarity.md) | Exact/lexical/semantic moslik, embedding ta'rifi, birinchi `vo.embed()`, provider-agnostik pattern |
| 02 | [Embedding modellari va tanlash](./02.%20Embedding%20modellari%20va%20tanlash.md) | 2026 model manzarasi, MTEB va uning tuzoqlari, o'z datada eval, Matryoshka, quantization, narx |
| 03 | [Similarity metrics — cosine, dot product, euclidean](./03.%20Similarity%20metrics%20—%20cosine,%20dot%20product,%20euclidean.md) | Uch metrika, normalizatsiya teoremasi, threshold kalibrlash, similarity matritsa, dedup |
| 04 | [Matn tayyorlash va chunking](./04.%20Matn%20tayyorlash%20va%20chunking.md) | Chunking strategiyalari, overlap, tozalash, markdown-aware chunker, contextual retrieval |
| 05 | [Bo'lim loyihasi — semantic search CLI](./05.%20Bo'lim%20loyihasi%20—%20semantic%20search%20CLI.md) | `semsearch` — lokal fayllar ustida semantic search, incremental index, provider-agnostik |

## Muhim fakt: Anthropic embedding API bermaydi

Kurs Claude asosida, lekin **Anthropic o'z embedding modelini taklif qilmaydi** — rasmiy hujjat
embeddings uchun **Voyage AI** ni tavsiya qiladi. Shu sababli bu bo'limda:

- Asosiy provider — **Voyage AI** (`voyage-4`, har model uchun 200M token bepul kvota)
- Kod provider-agnostik pattern bilan yoziladi (provider almashtirish — bitta funksiya)
- API'siz mashq uchun lokal muqobil: `sentence-transformers` (`BAAI/bge-m3`)

Berryman kitobida (2024) OpenAI `text-embedding-3-small` + FAISS ishlatilgan — kontseptlar
o'sha-o'sha, lekin darslar 2026 ekotizimiga moslangan. FAISS/pgvector 3-bo'limda; bu bo'limda
similarity hisoblash ataylab numpy'da — "qora quti"siz.

## Kod muhiti

```bash
python -m venv .venv && source .venv/bin/activate
pip install voyageai python-dotenv numpy
echo 'VOYAGE_API_KEY=pa-...' > .env
# lokal muqobil uchun (ixtiyoriy, ~2GB model yuklaydi):
pip install sentence-transformers
```

Asosiy model — `voyage-4` ($0.06/1M token, 200M bepul). Arzon qadam: `voyage-4-lite`.
LLM kerak bo'lgan joyda 1-bo'lim standarti: `claude-opus-4-8` / `claude-haiku-4-5`.

## Bo'lim tugagach

`05. Bo'lim loyihasi` dagi `semsearch` CLI ishlab turishi kerak. Portfolio zanjiri davom etadi:
`askops` + `quizgen` (1-bo'lim) → **`semsearch` (shu bo'lim)** → pgvector servisi (3-bo'lim) →
RAG savol-javob tizimi (4-bo'lim).
