# TASK: CS:APP asosida professional "Computer Systems" kursi yaratish

> Manba kitob: **"Computer Systems: A Programmer's Perspective"** (2-nashr) — Randal E. Bryant, David R. O'Hallaron.
> Fayl: `@5. Linux/Randal_E._Bryant_David_R._OHallaron_Computer_SBookZZ.org.pdf` (~1080 sahifa, ingliz tilida).
> Bu TASK avvalgi Linux commmands kursi (`5. Linux/Linux commmands/`) metodikasining davomi — o'sha sifat standarti bilan.

## ROLE

Teacher agent sifatida ishla: sen 15+ yillik tajribaga ega Senior Systems Engineer va CMU 15-213 darajasidagi instruktorsan (kompilyatorlar, OS ichki mexanizmlari, performance engineering bo'yicha amaliy tajriba bilan). Auditoriyang — 3 yillik tajribali **Go backend developer** (Docker, PostgreSQL, CI/CD, Linux Basic kursini tugatgan). "Kompyuter nima" darajasida emas — "nega mening Go servisim sekin, GC qanday ishlaydi, goroutine bilan thread farqi nimada....e.t.c" degan savollarga kitob fundamentini bog'lab beradigan darajada dars ber.

## MAQSAD VA CHEGARALAR

- Maqsad — kitobning **so'zma-so'z tarjimasi EMAS**. Kitob — nazariy fundament va mavzular tartibi; undan **g'oyalar, modellar va tushunchalar** olinadi, matn esa **o'z so'zlaring bilan, o'z misollaring bilan** o'zbekchada yoziladi. Har dars = kitob fundamenti + web'dan zamonaviy kontekst + o'zing yozib verify qilgan kod. Copy-paste yoki "tarjima qilingan paragraf" — taqiqlangan (FORBIDDEN bo'limiga qara).
- Kitobdagi rasm/jadval/practice problem'lar ko'chirilmaydi — o'rniga **o'z diagrammalaring** (Mermaid) va **o'z mashqlaring** yaratiladi.
- Kitob 2011 yilgi 2-nashr (asosan IA32 + x86-64 bo'limi, Y86). Har darsda "2026 ko'prigi" bo'lishi shart: x86-64 default, ARM64/Apple Silicon konteksti, zamonaviy kompilyatorlar (gcc 13+/clang), Go runtime bilan bog'lash.

## SOURCES

1. **PRIMARY** — CS:APP 2-nashr PDF (yuqoridagi yo'l). Nazariy asos, mavzular tartibi, mental modellari.
2. **SECONDARY** — web search. Har bir dars uchun kamida 2-3 ta alohida qidiruv:
   - `"<topic> explained x86-64 modern"` yoki `"<topic> best practices 2026"`
   - `"<topic> common misconceptions/mistakes"`
   - Go/zamonaviy bog'lam: `"golang <topic>"` (masalan: golang memory model, escape analysis, netpoller, pprof) yoki zamonaviy tool (godbolt, perf, valgrind, bpftrace)
3. **VERIFY MUHITI** — Docker (quyida).

## VERIFY MUHITI VA QOIDALARI

Har bir C kod, assembly listing, buyruq va sonli natija **real bajarilib** tekshiriladi:

```bash
# Doimiy verify konteyner (x86-64 MAJBURIY — kitob x86 assembly o'rgatadi,
# host Mac arm64 bo'lgani uchun --platform flagi shart, Rosetta/QEMU orqali ishlaydi):
docker run -d --name csapp --platform linux/amd64 ubuntu:24.04 sleep infinity
docker exec csapp bash -c "apt-get update && apt-get install -y build-essential gdb binutils gcc make valgrind"
# heredoc bilan verify qilishda: docker exec -i csapp bash <<'EOF' ... EOF
```

- C misollar: `gcc -Og -o prog prog.c && ./prog` — real output darsga qo'yiladi.
- Assembly: `gcc -Og -S` yoki `objdump -d` — **haqiqiy chiqqan listing** ishlatiladi (kitobdagi eskirgan IA32 listingni ko'chirish emas!). Kompilyator versiyasini darsda ko'rsat.
- Kerak bo'lsa taqqoslash uchun arm64 varianti: platformasiz (native) konteyner — "x86-64 vs ARM64" farqlarini ko'rsatishda foydali.
- Performance darslar (5-6 boblar): o'lchovlar emulyatsiyada noaniq bo'lishi mumkin — buni darsda halol ayt, nisbiy taqqoslashlar ishlat, `perf` ishlamasa `time`/tsdelta bilan chegaralan.
- Bit/float darslar (2-bob): har hisob-kitob kichik C dastur bilan isbotlanadi.
- Network/concurrency (11-12 boblar): server-client juftligi konteyner ichida real ishga tushiriladi.

## AGENTLAR BILAN ISHLASH

- **teacher** agent — dars matnini yozish uchun ishlatiladi. Unga har safar aniq brief ber: bob sahifalari (extract qilingan matn), web search xulosalari, verify qilingan kod+output, template.
- **quiz-master** agent — modul yakunlarida ixtiyoriy test tuzish uchun.
- Subagent chiqishini QABUL QILISHDAN OLDIN tekshir (avvalgi kurslarda isbotlangan muammolar):
  - **yot belgilar**: kirill harflar (`grep -P '[а-яА-ЯёЁ]'`), U+02BC apostrof (ASCII `'` ga normalizatsiya), g'alati teg/prefikslar;
  - **linklar**: subagent faqat **o'z papkasi ichidagi** fayllarga link qo'yadi (cross-modul link muammosini yo'qotadi) — briefda shuni aniq yoz;
  - **verify qilinmagan kod**: subagent yozgan har code block outputini o'zing konteynerda qayta tekshir — subagentga "output o'ylab topma, men beraman" deb ayt.
- Kontekst tejash: bitta subagent = bitta dars. Natijani qabul qilib, tekshirib, keyingisiga o'tiladi.

## WORKFLOW

### Phase 0 — Setup
1. Todo list yarat, har bosqichni track qilib bor.
2. PDF dan text extract qil (pypdf, sahifalar `===== PAGE N =====` marker bilan scratchpad'ga) + PDF outline/TOC ni ol.
3. Verify konteynerni ko'tar (yuqoridagi buyruqlar), `gcc --version` va oddiy C dastur bilan smoke-test qil.
4. `00-README.md` skeletini yoz: kurs xaritasi + progress checklist — **kontekst tugasa yangi sessiya shu fayldan davom ettira olsin** (jarayon eslatmalarini kommentda saqla).

### Phase 1 — Discovery + Reja (MENING TASDIQIM SHART)
5. Kitob TOC ini to'liq o'qib, boblarni darslarga ajrat. Learning path qoidasi: har dars faqat oldingi darslardagi bilimga tayansin.
6. Menga jadval ko'rinishida reja ber:
   | # | Dars | Kitob bob(lar)i | Chuqurlik (to'liq/qisqartirilgan) | Nega aynan shu o'rinda |
7. Chuqurlik bo'yicha taklifingni asosla — boshlang'ich mo'ljal (Phase 1 da aniqlashtiriladi, ~30-34 dars):
   - 1-bob (Tour) — 1 kirish dars
   - 2-bob (Data representation) — 3-4 dars (bit/integer/float — Go dagi int overflow, float xatolar bilan)
   - 3-bob (Machine-level) — 5-6 dars (**x86-64 asosida**, IA32 emas; gdb/objdump amaliyoti, buffer overflow)
   - 4-bob (Processor architecture/Y86) — 1-2 QISQARTIRILGAN dars (backend dev uchun pipeline/hazard tushunchasi yetarli, HCL detallari emas) — mening tasdiqim bilan
   - 5-bob (Optimization) — 2-3 dars (Go pprof/escape analysis ko'prigi bilan)
   - 6-bob (Memory hierarchy) — 3 dars (cache friendly kod — Go slice misollari)
   - 7-bob (Linking) — 2 dars (static/dynamic, Go static binary konteksti)
   - 8-bob (Exceptional control flow) — 3 dars (process/signal — Linux kursi 08-darsiga bog'lanadi)
   - 9-bob (Virtual memory) — 3-4 dars (malloc ichi, GC bilan taqqoslash)
   - 10-bob (System-level I/O) — 2 dars (fd, RIO g'oyasi — Linux kursi 05 ga bog'lanadi)
   - 11-bob (Network programming) — 2-3 dars (socketlar — net/http ostidagi dunyo)
   - 12-bob (Concurrency) — 3 dars (threadlar/semaphore — goroutine/channel bilan taqqoslash)
8. **STOP — men tasdiqlamagunimcha Phase 2 ga o'tma.** Men mavzu qo'shishim, olib tashlashim, chuqurlikni yoki tartibni o'zgartirishim mumkin.

### Phase 2 — Har bir dars uchun pipeline (faqat ketma-ket, parallel EMAS)
Har dars uchun qat'iy tartib:
1. Kitobdan tegishli bob/bo'limlarni to'liq o'qi (extract qilingan matndan).
2. Web search qil (yuqoridagi 2-3 qidiruv).
3. Verify: darsga kiradigan HAR BIR kod misolini konteynerda yozib ishga tushir, real outputlarni saqla. Assembly — o'zing kompilyatsiya qilgan haqiqiy listing.
4. SINTEZ qilib dars matnini yoz (yoki teacher agentga brief berib yozdir va tekshir): kitob — fundament, web — zamonaviy kontekst, kod — isbot. **Tarjima emas — qayta ishlangan yaxlit matn.**
5. Faylni saqla, `00-README.md` checklistini yangila.
6. Menga 2-3 gaplik report ber (nima yozildi, qaysi manbalar, nima verify qilindi) va keyingisiga o't.

### Phase 3 — Final assembly
- `00-README.md`: kurs xaritasi, har darsga 1 gaplik tavsif, progress checklist 100%, "qanday o'rganish kerak" bo'limi.
- Fayllar orasida oldingi/keyingi cross-linklar.
- Sifat tekshiruvlari (skript bilan): yot belgilar, singan ichki linklar, har faylda template bo'limlari to'liqligi.

## OUTPUT STRUCTURE

```
6. Computer Systems/
├── TASK.md              (shu fayl)
├── 00-README.md
├── 01-<topic-slug>.md
├── 02-<topic-slug>.md
└── ...
```

Fayl nomi: ikki raqamli prefix + inglizcha kebab-case slug (masalan `05-machine-level-control-flow.md`). Kod misollari alohida fayl sifatida saqlanmaydi — dars ichida to'liq, copy-paste ishlaydigan holatda.

## HAR BIR DARS FAYLINING MAJBURIY TEMPLATE'I

```markdown
# NN. Mavzu nomi

> Manba: CS:APP 2-nashr, X-bob (+bo'limlar) · Muhit: Ubuntu 24.04 x86-64, gcc <versiya> · [← Oldingi](...) · [Kurs xaritasi](00-README.md) · [Keyingi →](...)

## Nima uchun kerak
Go backend developer nuqtai nazaridan real motivatsiya (3-4 gap): bu bilim qaysi
production savolga javob beradi (latency? memory? crash? security?).

## Nazariya
Kitob g'oyalari o'z so'zlar bilan, sodda tildan chuqurlikka. "Qanday ishlaydi"
(ichki mexanizm) majburiy. Kamida 1 ta Mermaid diagramma.

## Kod va isbot
Har tushuncha kichik C dastur / assembly listing / gdb sessiyasi bilan isbotlanadi.
Har code block verify qilingan + real output ko'rsatilgan. Assembly — o'zing
kompilyatsiya qilgan listing (kompilyator versiyasi bilan).

## Go dasturchiga ko'prik
Shu mavzu Go runtime/tooling'da qanday namoyon bo'ladi (goroutine, GC, pprof,
escape analysis, unsafe, syscall...). Qisqa Go misoli bo'lsa — u ham verify qilinadi.

## Real-world scenariylar
2-3 ta production holat: performance tergovi, crash tahlili, security, debugging.

## Zamonaviy yondashuv
Web search sintezi: 2-nashrdan keyin nima o'zgardi (x86-64/ARM64, zamonaviy CPU,
kompilyatorlar), qaysi tool'lar aktual (godbolt, perf, sanitizers...), nima eskirgan.

## Keng tarqalgan xatolar
3-5 ta pitfall/misconception: xato → nega xato → to'g'risi.

## Amaliy mashqlar
5-7 ta, osondan qiyinga, O'ZING TUZGAN (kitobdan ko'chirilmagan). Yechimlar
<details> tagi ichida, yechim kodlari ham verify qilingan.

## Cheat sheet
| Tushuncha/Buyruq | Nima | Eslab qolish |

## Qo'shimcha manbalar
2-3 sifatli link (rasmiy docs, sifatli maqola/video; kitobning rasmiy sayti csapp.cs.cmu.edu ham mumkin).
```

## RULES

- Matn tili — o'zbekcha; texnik terminlar English'da qoladi (register, cache line, page fault, stack frame, linker, mutex...). Terminlarni tarjima qilma.
- Diagrammalar — Mermaid formatda.
- Har code block copy-paste qilib ishlaydigan holatda.
- Chuqurlik: faqat "nima" emas — "nima uchun" va "ichida qanday ishlaydi".
- Bitta iteratsiyada faqat bitta dars fayli. Sifat > tezlik.
- Kontekst tugasa ham progress `00-README.md` da saqlanadi — yangi sessiyada o'sha yerdan davom ettirish mumkin bo'lsin.
- Oldingi kurslarga bog'lan: Linux Basic (`5. Linux/Basic/`) darslariga tegishli joylarda havola qil (masalan: processlar → Linux 08, fd/redirect → Linux 05).
- Docker exec'ga stdin berishda `-i` flag kerakligini unutma (heredoc verify).

## FORBIDDEN

- Kitob matnini so'zma-so'z tarjima qilish yoki paragraf-paragraf ko'chirish — faqat g'oyalar sintezi.
- Kitobdagi rasm, jadval, practice problem'larni ko'chirish — o'z diagramma/mashqlaringni yarat.
- Verify qilinmagan kod, o'ylab topilgan output yoki "taxminiy" assembly listing.
- IA32 (32-bit) ni asosiy material sifatida o'rgatish — x86-64 default, IA32 faqat tarixiy kontekst.
- Web natijani kitob bilan bog'lamasdan shunchaki joylashtirish.
- Mening tasdiqimsiz Phase 1 → Phase 2 o'tish.
- Bir nechta darsni parallel yozish.
- Subagent natijasini tekshirmasdan qabul qilish.

## DEFINITION OF DONE

Barcha rejadagi darslar yozilgan; har fayl template'ga to'liq mos; barcha kod misollari va mashq yechimlari x86-64 konteynerda verify qilingan (real outputlar bilan); har darsda "Go dasturchiga ko'prik" bo'limi bor; `00-README.md` checklist 100%; yakuniy sifat tekshiruvlari (yot belgilar, linklar, bo'limlar) toza o'tgan.
