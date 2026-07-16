# 4. RAG

Zanjir shu yerda birlashadi: 1-bo'lim (LLM API), 2-bo'lim (embeddings) va 3-bo'lim (vector DB +
`vecsearch` retrieval servisi) ustiga endi **generation qatlami** quriladi — retrieval'dan kelgan
kontekst bilan Claude javob beradi, manba iqtiboslari (citations) bilan. Bo'lim tugagach:
RAG qachon kerakligini asoslash, retrieval sifatini o'lchash (golden set, recall@k, MRR),
reranking, query optimization (rewriting, multi-query, HyDE), advanced indexing (parent-document,
contextual retrieval) va javob sifatini baholash (faithfulness, relevance) — hammasi qo'lda.

## Kim uchun

Middle backend dasturchi. 1-3 bo'limlar tugagan deb hisoblanadi — ayniqsa 3-bo'lim loyihasi
(`vecsearch`): bu bo'limdagi SQL misollar o'sha sxemaga (chunks jadvali, HNSW, tsv, RRF) tayanadi.
Markaziy tamoyil: **retrieval sifati = tizim sifati** — javob yomon bo'lsa avval retrieval
o'lchanadi, model emas.

## Darslar

| # | Dars | Nima o'rganasiz |
|---|---|---|
| 01 | [RAG arxitekturasi — qachon kerak va birinchi pipeline](./01.%20RAG%20arxitekturasi%20—%20qachon%20kerak%20va%20birinchi%20pipeline.md) | RAG vs long context vs fine-tune (200K qoidasi), 3 modul, minimal end-to-end pipeline + citations |
| 02 | [Chunking strategiyalari — chuqurlashuv](./02.%20Chunking%20strategiyalari%20—%20chuqurlashuv.md) | 2026 benchmark (recursive 512 nega yutadi), chunk header, semantic chunking qachon arziydi |
| 03 | [Retrieval sifatini o'lchash — golden set, recall@k, MRR](./03.%20Retrieval%20sifatini%20o'lchash%20—%20golden%20set,%20recall@k,%20MRR.md) | Golden dataset qurish, recall@k/precision@k/MRR/NDCG qo'lda, diagnostika tartibi |
| 04 | [Reranking — cross-encoder bilan aniqlikni oshirish](./04.%20Reranking%20—%20cross-encoder%20bilan%20aniqlikni%20oshirish.md) | Bi- vs cross-encoder, `rerank-2.5`, ikki bosqichli funnel, "recall past bo'lsa rerank yordam bermaydi" |
| 05 | [Query optimization — rewriting, multi-query, HyDE](./05.%20Query%20optimization%20—%20rewriting,%20multi-query,%20HyDE.md) | Query-shape failures, suhbatda rewriting, multi-query + RRF, HyDE va gate pattern, self-query |
| 06 | [Advanced indexing — parent-document va contextual retrieval](./06.%20Advanced%20indexing%20—%20parent-document%20va%20contextual%20retrieval.md) | Small-to-big, auto-merging, Anthropic contextual retrieval (-67% failure) + prompt caching iqtisodi |
| 07 | [RAG javob sifati — faithfulness, relevance va citations](./07.%20RAG%20javob%20sifati%20—%20faithfulness,%20relevance%20va%20citations.md) | RAG triad, citations API chuqur, qo'lda mini-faithfulness/relevance, LLM-judge ogohlantirishlari |
| 08 | [Bo'lim loyihasi — docqa savol-javob tizimi](./08.%20Bo'lim%20loyihasi%20—%20docqa%20savol-javob%20tizimi.md) | `vecsearch` + generation: hybrid → rerank → Claude + citations, POST /ask, mini eval |

## Asosiy qarorlar (2026 holati)

- **Generation: `claude-opus-4-8`**, arzon LLM qadamlar (rewriting, multi-query, chunk-kontekst,
  claim ajratish) — `claude-haiku-4-5`. Retrieval standarti o'zgarmagan: `voyage-4` + `input_type`,
  reranking `rerank-2.5`.
- **Citations — API funksiyasi, prompt-hack emas:** retrieved chunk'lar `document` content block
  sifatida (`citations: {"enabled": true}`) yuboriladi, javob `cited_text`/`document_index` bilan
  keladi. Structured outputs bilan birga ishlamaydi (400).
- **Anthropic 200K qoidasi:** knowledge base ~200K tokendan kichik bo'lsa RAG shart emas —
  to'liq kontekst + prompt caching. RAG'ni ehtiyoj isbotlanganda quring.
- Kitob (Iusztin & Labonne) LangChain+OpenAI'da — kontseptlar kitobdan, kod bizning stack'da
  (raw Claude API + psycopg + pgvector).

## Kod muhiti

```bash
python -m venv .venv && source .venv/bin/activate
pip install anthropic voyageai "psycopg[binary]" psycopg-pool pgvector fastapi uvicorn python-dotenv numpy
```

`.env`: `ANTHROPIC_API_KEY`, `VOYAGE_API_KEY`, `DATABASE_URL` (3-bo'limdagi pgvector konteyner).
Postgres/`vecsearch` 3-bo'limdagi Docker compose bilan ko'tarilgan deb hisoblanadi.

## Bo'lim tugagach

`08. Bo'lim loyihasi`dagi `docqa` servisi ishlab turishi kerak: `POST /ask` o'zbekcha savolga
manba iqtibosli javob qaytaradi. Portfolio zanjiri: `askops` (1) → `semsearch` (2) →
`vecsearch` (3) → **`docqa` (shu bo'lim)** → eval harness (6-bo'lim) → production + Telegram bot (7-bo'lim).
