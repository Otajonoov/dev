# 6. Evaluation

> Testsiz merge qilmaysiz, monitoringsiz deploy qilmaysiz — eval'siz LLM feature ham
> chiqarilmaydi. 4-bo'limda recall@k va mini-faithfulness'ni qo'lda yozgansiz; bu bo'lim
> o'sha poydevorni to'laqonli kasbiy ko'nikmaga aylantiradi: golden dataset, ishonchli
> LLM-as-judge, CI'dagi regression testing va model tanlash qarori.

## Darslar

| # | Dars | Asosiy savol |
|---|------|--------------|
| 01 | [Nega eval — evaluation-driven development](01.%20Nega%20eval%20—%20evaluation-driven%20development.md) | Mezonlar kodni yozishdan OLDIN qanday aniqlanadi va biznesga qanday bog'lanadi? |
| 02 | [Golden dataset qurish](02.%20Golden%20dataset%20qurish.md) | 10 savollik demo set qanday qilib yashovchi, versiyalangan datasetga aylanadi? |
| 03 | [LLM-as-judge — chuqur](03.%20LLM-as-judge%20—%20chuqur.md) | Judge qanday yoziladi, qayerda buziladi (bias'lar) va qanday kalibrlanadi? |
| 04 | [Offline va online eval — regression testing va CI](04.%20Offline%20va%20online%20eval%20—%20regression%20testing%20va%20CI.md) | Har o'zgarishda regression qanday ushlanadi va production'da eval qanday davom etadi? |
| 05 | [Model selection va benchmark'lar](05.%20Model%20selection%20va%20benchmark'lar.md) | Public benchmark'ga qachon ishonish mumkin va model qarori qanday qabul qilinadi? |
| 06 | [Bo'lim loyihasi — evalharness](06.%20Bo'lim%20loyihasi%20—%20evalharness.md) | Hammasini birlashtirish: docqa uchun to'liq eval harness |

## Poydevor

- **4-bo'lim, 03-dars** — golden set kirish, recall@k/precision@k/MRR/NDCG qo'lda
  implementatsiya. Bu bo'lim metrikalarni QAYTA TUSHUNTIRMAYDI — ishlatadi.
- **4-bo'lim, 07-dars** — RAG triad, mini-faithfulness judge, answer relevance.
  03-dars (judge) shu poydevorga quriladi.
- **4-bo'lim, 08-dars (docqa)** — bo'lim loyihasi aynan shu tizimni baholaydi;
  `docqa` ko'tarilgan bo'lishi kerak (docker compose up).
- **1-bo'lim, 04-dars (structured output)** — judge chiqishi `messages.parse` + Pydantic
  bilan olinadi; `quizgen` tajribasi sintetik savol generatsiyasida qayta ishlaydi.

## Bo'lim loyihasi: `evalharness`

`docqa` uchun to'liq eval harness — portfolio zanjirining 5-bo'g'ini
(askops → semsearch → vecsearch → docqa → **evalharness**): versiyalangan golden dataset
(JSONL, o'zbekcha, slice'lar bilan), retrieval eval (recall@5/MRR), LLM-judge generation
eval (faithfulness + relevance, kalibrlash qadami bilan), regression testing
(baseline + threshold + exit code — CI-ready), Batches API rejimi (50% arzon) va
markdown hisobot.

## Texnik standart

- Kod: Python 3.12, raw Claude API (`anthropic` SDK); default judge `claude-haiku-4-5`,
  kalibrlash va murakkab hukmlar `claude-opus-4-8`; API key `.env`da.
- Judge chiqishi: `client.messages.parse(..., output_format=PydanticModel)` —
  JSON'ni regex bilan qazish yo'q.
- Eval run transporti: Batches API (50% chegirma) — katta setlar uchun standart.
- Token/narx hisobi: `client.messages.count_tokens` (tiktoken EMAS).
- Eval framework'lari (Ragas, DeepEval, promptfoo, LangSmith, Braintrust, TruLens)
  landshaft sifatida o'rganiladi — ishlaydigan kod framework'siz: kontseptni nolda
  qurgan odam istalgan framework'ni bir kunda o'zlashtiradi.
- API faktlarining haqiqat manbai: `../x. Manbalar/research-6-evaluation.md`.

## Keyingi bo'limga ko'prik

Eval "nimani chiqarayotganimiz sifatlimi?" savoliga javob beradi. **7. Production**
bo'limi keyingi savolni oladi: shu tizimni real traffic ostida qanday tez, arzon va
kuzatiladigan qilib ishlatamiz (serving, caching, observability, guardrails).
