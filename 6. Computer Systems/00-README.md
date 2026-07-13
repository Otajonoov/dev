# Computer Systems — CS:APP asosidagi kurs

> Manba: **"Computer Systems: A Programmer's Perspective"** (2-nashr) — R. Bryant, D. O'Hallaron.
> Auditoriya: Go backend developer (3 yil tajriba). Til: o'zbekcha, terminlar English.
> Verify muhiti: Docker `csapp` konteyner (Ubuntu 24.04, `--platform linux/amd64`), gcc.

<!--
JARAYON ESLATMALARI (yangi sessiya uchun):
- TASK.md — to'liq spetsifikatsiya (template, qoidalar, forbidden ro'yxati). Har dars TASK.md dagi template'ga mos yozilishi SHART.
- PDF matni: scratchpad/csapp_text.txt (===== PAGE N ===== markerlar, pypdf bilan extract qilingan, 1080 sahifa).
  Yo'q bo'lsa qayta extract: scratchpad/extract_pdf.py. TOC: scratchpad/csapp_toc.txt.
  PDF sahifasi = kitob sahifasi + 35 (masalan 3.7 Procedures kitobda 219, PDF da 254).
- Verify konteyner: docker run -d --name csapp --platform linux/amd64 ubuntu:24.04 sleep infinity
  docker exec csapp bash -c "apt-get update && apt-get install -y build-essential gdb binutils gcc make valgrind bsdmainutils file golang-go strace"
  Heredoc: docker exec -i csapp bash <<'EOF' ... EOF
- MUHIM (GDB): `csapp` (QEMU emulyatsiya) ichida gdb JONLI debug qila olmaydi — ptrace yo'q
  ("Cannot PTRACE_GETREGS"). Statik `gdb -batch -ex 'disassemble f'` ishlaydi, `run`/`break` YO'Q.
  Yechim — jonli x86-64 debug uchun ALOHIDA native arm64 konteyner + cross-compiler + qemu gdbserver:
    docker run -d --name csapp_arm ubuntu:24.04 sleep infinity
    docker exec csapp_arm bash -c "apt-get update && apt-get install -y gcc-x86-64-linux-gnu qemu-user gdb-multiarch"
    # kompilyatsiya: x86_64-linux-gnu-gcc -Og -g -static -o t t.c
    # server:       qemu-x86_64 -g 1234 ./t &
    # klient:       gdb-multiarch -batch -q -x cmds.gdb ./t
    #   cmds.gdb: set architecture i386:x86-64 / target remote localhost:1234 / break f / continue / bt ...
  Eslatma: QEMU ostida stack address'lar 0x2aaa... ko'rinishida (native Linux'da 0x7fff...) — darsda halol aytilsin.
- Pipeline (TASK.md Phase 2): kitob bo'limini o'qish -> 2-3 web search -> HAR kod konteynerda verify -> sintez (teacher agent yoki o'zi) -> README checklist yangilash -> qisqa report.
- PERFORMANCE o'lchovlari (12-18 darslar, pipeline/cache): QEMU pipeline'ni va cache'ni emulyatsiya QILMAYDI —
  sorted-vs-unsorted, ILP, cache miss effektlari `csapp` (QEMU) da KO'RINMAYDI. Yechim: o'lchovlarni `csapp_arm`
  NATIVE arm64 host'da o'tkaz (haqiqiy apparat pipeline/cache). Pipeline/cache hodisalari arxitekturadan mustaqil —
  arm64 o'lchovi x86-64 g'oyasini to'g'ri ko'rsatadi. Darsda "o'lchov Apple Silicon/arm64 native, effekt universal" deb halol ayt.
  csapp_arm da native gcc bor. Kompilyator dead-code elimination'dan saqlanish: volatile sink + natijani ishlatish.
- Subagent tekshiruvi: kirill harflar (grep -P '[а-яА-ЯёЁ]'), U+02BC apostrof, faqat shu papka ichidagi linklar, verify qilinmagan output.
- MUHIM SABOQ (09-dars): subagent brief uzilib qolsa, assembly listinglarni O'ZI TO'QIB chiqaradi (add3/pick misollari
  noto'g'ri chiqqan edi). Shuning uchun subagent yozgan HAR assembly blokini `gcc -Og -S` bilan qayta kompilyatsiya
  qilib solishtir. Brief to'liq va uzilmagan bo'lsa muammо kamroq, lekin tekshiruv baribir majburiy.
- Linux kursi linklari: ../5. Linux/1. Linux commands/ (05-redirection-and-pipelines.md, 08-processes.md).
- Bitta iteratsiya = bitta dars. Parallel yozish taqiqlangan.
-->

## Holat

- [x] Phase 0 — Setup (PDF extract, Docker verify muhiti, README skeleti)
- [x] Phase 1 — Reja foydalanuvchi tomonidan TASDIQLANGAN (2026-07-12, 34 dars)
- [ ] Phase 2 — Darslar (quyidagi checklist)
- [ ] Phase 3 — Final assembly (cross-linklar, sifat tekshiruvlari)

## Kurs xaritasi

| # | Dars | Bob | Holat |
|---|------|-----|-------|
| 01 | [Tour of Computer Systems](01-tour-of-computer-systems.md) | 1 | ✅ |
| 02 | [Information Storage: bits, bytes, hex](02-information-storage.md) | 2.1 | ✅ |
| 03 | [Integer Representation](03-integer-representation.md) | 2.2 | ✅ |
| 04 | [Integer Arithmetic](04-integer-arithmetic.md) | 2.3 | ✅ |
| 05 | [Floating Point (IEEE 754)](05-floating-point.md) | 2.4 | ✅ |
| 06 | [Machine-Level Basics: registers, encodings](06-machine-level-basics.md) | 3.1–3.3, 3.13 | ✅ |
| 07 | [Data Movement va Arithmetic](07-data-movement-arithmetic.md) | 3.4–3.5 | ✅ |
| 08 | [Control Flow](08-machine-control-flow.md) | 3.6 | ✅ |
| 09 | [Procedures va Stack](09-procedures-stack.md) | 3.7, 3.13 | ✅ |
| 10 | [Arrays, Structs, Pointers](10-arrays-structs-pointers.md) | 3.8–3.10 | ✅ |
| 11 | [GDB va Buffer Overflow](11-gdb-buffer-overflow.md) | 3.11–3.12 | ✅ |
| 12 | [CPU Pipeline asoslari (qisqartirilgan)](12-cpu-pipeline.md) | 4.1, 4.4–4.5 | ✅ |
| 13 | [Kompilyator optimizatsiyasi chegaralari](13-compiler-optimization.md) | 5.1–5.6 | ✅ |
| 14 | [Zamonaviy CPU va Instruction-Level Parallelism](14-instruction-level-parallelism.md) | 5.7–5.10 | ✅ |
| 15 | [Profiling va Bottleneck'lar](15-profiling-bottlenecks.md) | 5.11–5.14 | ✅ |
| 16 | [Storage texnologiyalari va Locality](16-storage-locality.md) | 6.1–6.3 | ✅ |
| 17 | [Cache Memories](17-cache-memories.md) | 6.4 | ✅ |
| 18 | [Cache-Friendly Kod](18-cache-friendly-code.md) | 6.5–6.6 | ✅ |
| 19 | [Static Linking va ELF](19-static-linking.md) | 7.1–7.9 | ✅ |
| 20 | [Dynamic Linking va PIC](20-dynamic-linking.md) | 7.10–7.13 | ✅ |
| 21 | [Exceptions va Processlar](21-exceptions-processes.md) | 8.1–8.3 | ✅ |
| 22 | [Process Control: fork, exec, wait](22-process-control.md) | 8.4 | ✅ |
| 23 | [Signals](23-signals.md) | 8.5–8.7 | ✅ |
| 24 | [Virtual Memory tushunchalari](24-virtual-memory.md) | 9.1–9.6 | ✅ |
| 25 | [Linux Memory System va mmap](25-linux-memory-mmap.md) | 9.7–9.8 | ✅ |
| 26 | [Dynamic Memory Allocation (malloc ichi)](26-dynamic-memory-allocation.md) | 9.9 | ✅ |
| 27 | [Garbage Collection va Memory Bug'lar](27-garbage-collection.md) | 9.10–9.11 | ✅ |
| 28 | [Unix I/O: fayl deskriptorlari](28-unix-io.md) | 10.1–10.4 | ✅ |
| 29 | [Fayl metadata, sharing, redirection](29-file-metadata-sharing.md) | 10.5–10.9 | ✅ |
| 30 | Sockets Interface | 11.1–11.4 | ⬜ |
| 31 | Web Server ichidan | 11.5–11.6 | ⬜ |
| 32 | Concurrency modellari: process, epoll, thread | 12.1–12.3 | ⬜ |
| 33 | Shared Variables va Semaphore'lar | 12.4–12.5 | ⬜ |
| 34 | Parallelism va Concurrency muammolari | 12.6–12.7 | ⬜ |

## Qanday o'rganish kerak

(Phase 3 da to'ldiriladi.)
