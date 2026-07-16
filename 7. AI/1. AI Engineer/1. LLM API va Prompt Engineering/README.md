# 1. LLM API va Prompt Engineering

AI Engineer yo'lining boshlanish nuqtasi. Bu bo'lim tugagach siz LLM'ni **production darajasida**
integratsiya qila olasiz: streaming, tool use, structured output, xato/retry/rate limit boshqaruvi,
narx nazorati va prompt injection himoyasi.

## Kim uchun

Middle backend dasturchi. JSON, REST, retry, exponential backoff, connection pool kabi tushunchalar
ma'lum deb qabul qilinadi — analogiyalar shu dunyodan olinadi. Python bilish talab qilinadi.

## Darslar

| # | Dars | Nima o'rganasiz |
|---|---|---|
| 01 | [LLM dasturchi ko'zi bilan](./01.%20LLM%20dasturchi%20ko'zi%20bilan.md) | Token, context window, autoregressiya, sampling, hallucination va truth bias |
| 02 | [Claude API — birinchi so'rov](./02.%20Claude%20API%20—%20birinchi%20so'rov.md) | SDK, messages, model tanlash va narx, `stop_reason`, xato turlari, retry, rate limit |
| 03 | [Streaming](./03.%20Streaming.md) | SSE event'lari, TTFT, cancel, timeout — nega streaming ko'pincha majburiy |
| 04 | [Structured output](./04.%20Structured%20output.md) | Constrained decoding, Pydantic bilan `messages.parse`, strict tool, tuzoqlar |
| 05 | [Tool use](./05.%20Tool%20use.md) | Agent loop, `tool_result`, parallel tool'lar, iteratsiya limiti, xavfli amallar |
| 06 | [Prompt engineering asoslari](./06.%20Prompt%20engineering%20asoslari.md) | Prompt anatomiyasi, XML teglar, few-shot va uning bias'lari, prompt registry |
| 07 | [Ilg'or prompt texnikalari](./07.%20Ilg'or%20prompt%20texnikalari.md) | CoT va reasoning-model paradoksi, self-consistency, decomposition, Reflexion |
| 08 | [Prompt injection va himoya](./08.%20Prompt%20injection%20va%20himoya.md) | Indirect injection, OWASP qatlamlari, dual-LLM, violation vs false refusal rate |
| 09 | [Bo'lim loyihasi](./09.%20Bo'lim%20loyihasi.md) | `askops` CLI chatbot + `quizgen` savol generatori — portfolio uchun |

## Muhim ogohlantirish: kitoblar 2024, API 2026

Bu bo'lim ikkita kuchli kitobga tayanadi (Berryman & Ziegler — "Prompt Engineering for LLMs";
Chip Huyen — "AI Engineering"), lekin ularda markaziy o'rin egallagan bir nechta amaliyot
Claude'ning joriy modellarida **umuman ishlamaydi**:

| Kitobda | 2026, Claude Opus 4.8 / Sonnet 5 |
|---|---|
| `temperature`, `top_p`, `top_k` | 400 xato — parametr olib tashlangan |
| Assistant prefill ("inception trick") | 400 xato — o'rniga structured output |
| `budget_tokens` bilan extended thinking | 400 xato — o'rniga `adaptive` + `effort` |
| `logprobs` bilan calibration | Anthropic logprob'larni ochmaydi |

Darslarda bu farqlar ochiq belgilangan. Tushunchalar baribir kerak — ular OpenAI-compatible
API'larda va lokal modellarda ishlaydi, va intervyuda so'raladi.

## Kod muhiti

```bash
python -m venv .venv && source .venv/bin/activate
pip install anthropic python-dotenv pydantic
echo 'ANTHROPIC_API_KEY=sk-ant-...' > .env
```

Asosiy model — `claude-opus-4-8`. Arzon/tez qadamlarda `claude-haiku-4-5`.
Har misol mustaqil ishga tushadigan fayl.

## Bo'lim tugagach

`09. Bo'lim loyihasi` dagi ikkala loyiha ishlab turishi kerak. Ular keyingi bo'limlarda
kengaytiriladi: semantic search, pgvector, RAG, eval harness va production Telegram bot —
bularning hammasi birgalikda portfolio hosil qiladi.
