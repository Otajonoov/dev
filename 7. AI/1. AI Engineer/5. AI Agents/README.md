# 5. AI Agents

> Agent — sirli narsa emas: **while loop + tool calls**. 1-bo'limda o'zing yozgan tool use
> loop shu bo'limda to'laqonli agent'ga aylanadi: planning, memory, MCP, multi-agent
> patterns va production darajasidagi xavfsizlik bilan.

## Darslar

| # | Dars | Asosiy savol |
|---|------|--------------|
| 01 | [Agent nima — loop, environment, tools](01.%20Agent%20nima%20—%20loop,%20environment,%20tools.md) | Agent qachon kerak va loop qanday ishlaydi? |
| 02 | [Workflow patterns — chaining, routing, parallelization](02.%20Workflow%20patterns%20—%20chaining,%20routing,%20parallelization.md) | Hamma narsa agent emas — 5 workflow pattern qachon yetadi? |
| 03 | [Tool design — agent uchun yaxshi tool yozish](03.%20Tool%20design%20—%20agent%20uchun%20yaxshi%20tool%20yozish.md) | Nega tool sifati = agent sifati? Bash vs dedicated tool? |
| 04 | [Planning va reflection — ReAct, Reflexion, Tool Runner](04.%20Planning%20va%20reflection%20—%20ReAct,%20Reflexion,%20Tool%20Runner.md) | Reja va o'z-o'zini tekshirish agentni qanday kuchaytiradi? |
| 05 | [Agent memory va context engineering](05.%20Agent%20memory%20va%20context%20engineering.md) | Uzun sessiyada kontekst qanday boshqariladi? |
| 06 | [MCP — Model Context Protocol](06.%20MCP%20—%20Model%20Context%20Protocol.md) | Tool integratsiyasining standart protokoli — server yozish va ulash |
| 07 | [Multi-agent patterns va framework tanlovi](07.%20Multi-agent%20patterns%20va%20framework%20tanlovi.md) | Qachon bir nechta agent kerak? Framework yoki raw API? |
| 08 | [Agent xavfsizligi — sandbox, approval, audit](08.%20Agent%20xavfsizligi%20—%20sandbox,%20approval,%20audit.md) | Write action'li agentni production'ga qanday chiqariladi? |
| 09 | [Bo'lim loyihasi — repoagent (MCP bilan repo tahlilchisi)](09.%20Bo'lim%20loyihasi%20—%20repoagent%20%28MCP%20bilan%20repo%20tahlilchisi%29.md) | Hammasini birlashtirish: mustaqil task agent |

## Poydevor

- **1-bo'lim, 05-dars (Tool use)** — `while stop_reason == "tool_use"` loop'ni o'sha yerda
  yozgansiz; bu bo'lim shundan boshlanadi, takrorlamaydi.
- **1-bo'lim, 08-dars (Prompt injection)** — xavfsizlik darsi (08) shu poydevorga quriladi.
- **4-bo'lim (RAG)** — retriever agent'ning tool'i sifatida ham ishlaydi (03-darsda mashq bor).

## Bo'lim loyihasi: `repoagent`

Lokal git repo'ni tahlil qiluvchi mustaqil agent: MCP server (`repotools`) + raw API manual
loop klient. Planning, human-in-the-loop approval, audit log va notes memory bilan — bo'limdagi
barcha darslarning sintezi. Portfolio zanjiridan (askops → semsearch → vecsearch → docqa)
alohida turadi: bu "mustaqil task agent" namunasi.

## Texnik standart

- Kod: Python 3.12, raw Claude API (`anthropic` SDK), model `claude-opus-4-8`
  (arzon qadamlar: `claude-haiku-4-5`), API key `.env`da.
- MCP: rasmiy `mcp` SDK (`FastMCP`), transport stdio (lokal) / Streamable HTTP (remote).
- Framework'lar (LangGraph, CrewAI, smolagents) 07-darsda taqqoslanadi, lekin ishlaydigan
  kod framework'siz — Anthropic tavsiyasi: "start by using LLM APIs directly".
- API faktlarining haqiqat manbai: `../x. Manbalar/research-5-agents.md`.

## Keyingi bo'limga ko'prik

Agent failure modes (planning/tool/efficiency) shu bo'limda tanishtiriladi — ularni
tizimli O'LCHASH esa **6. Evaluation** bo'limining mavzusi.
