# 1. Kirish va maqsad

## 1.1. Nima uchun pure Go + `unsafe`?

Go tilida past darajadagi (low-level) ishlarni qilishning **uchta asosiy yo'li** bor:

| Yo'l | Tavsif | Murakkablik | Portativligi |
|------|--------|-------------|--------------|
| **Pure Go + `unsafe`** | Faqat Go tili va `unsafe` paketi | O'rta | Yuqori |
| **Cgo** | C kutubxonalarini chaqirish | Yuqori | Past |
| **Solod (Assembly)** | Plan9 assembly, `.s` fayllar | Juda yuqori | Past (CPU arch'ga bog'liq) |

### Nima uchun aynan **pure Go + `unsafe`**?

```mermaid
flowchart LR
    A[Maqsad: CS o'rganish] --> B{Yo'l tanlash}
    B -->|"Eng tezkor o'rganish"| C[Pure Go + unsafe]
    B -->|"Tashqi C kod"| D[Cgo]
    B -->|"CPU darajasi"| E[Solod/Assembly]

    C --> F[Go ekosistemida qoladi]
    C --> G[Hech qanday tashqi build kerak emas]
    C --> H[Goroutine va GC bilan integratsiya]
    C --> I[Cross-platform]

    style C fill:#90EE90
    style D fill:#FFE4B5
    style E fill:#FFB6C1
```

**Sabablari:**
1. **Go ekosistemasidan chiqmaysiz** — `go build`, `go test`, `go bench` hammasi ishlaydi
2. **Cross-platform** — Linux, macOS, Windows hammasida ishlaydi (kichik o'zgarishlar bilan)
3. **Goroutine va GC bilan integratsiya** — yozgan strukturangizni kodingizda ishlatish oson
4. **C bilmasangiz ham bo'ladi** — faqat Go bilish kifoya
5. **Asosiy fundamental bilim** — pointer, memory layout, alignment hammasini o'rganasiz
6. **Solod'ga kelajakda tayyorgarlik** — `unsafe`'da puxta bo'lsangiz, assembly osonroq ko'rinadi

## 1.2. Afzalliklari va kamchiliklari

```mermaid
mindmap
  root((Pure Go + unsafe))
    Afzalliklari
      Tez o'rganish
      Cross-platform
      Go GC bilan integratsiya
      Hech qanday tashqi tool yo'q
      Debugging oson (delve)
      Test va benchmark oson
    Kamchiliklari
      Maksimal tezlik emas
      GC ishtiroki bor
      Ba'zi optimizatsiya yo'q
      unsafe xavfli (pointer arithmetic)
      Go versiyasiga bog'liq (internal API)
```

| Afzallik | Kamchilik |
|----------|-----------|
| Tez o'rganasiz, sintaksis tanish | Maksimal tezlik (assembly darajasi) emas |
| Cross-platform (deyarli) | Go runtime ichidagi ba'zi narsalarga to'g'ridan-to'g'ri tegolmaysiz |
| GC bilan birga ishlaydi | `unsafe` xavfli — pointer xatosi `panic` chaqiradi |
| Test va benchmark stdlib bilan | Go versiyasi o'zgarsa, kod buzilishi mumkin |
| Pointer va memory tushunchasini chuqur o'rgatadi | C dasturchilari ko'pi hali ham Cgo'ni afzal ko'radi |

## 1.3. Kim uchun mos?

Bu yo'l mos keladi agar siz:
- Go'ni **o'rta darajada** bilsangiz (goroutine, channel, slice, map ishlata olsangiz)
- **CS fundamental** narsalarni o'rganmoqchi bo'lsangiz (pointer, memory, allocator)
- **Database internals**, **GC algorithms**, **lock-free data structures** ga qiziqsangiz
- **BadgerDB, CockroachDB, Redis** kabi loyihalarning ichki dunyosini ko'rmoqchi bo'lsangiz
- Production uchun **emas**, balki **bilim uchun** ishlamoqchi bo'lsangiz

Bu yo'l mos **kelmaydi** agar siz:
- Eng yuqori tezlik kerak bo'lsa (u holda assembly)
- C kutubxonalarini integratsiya qilish kerak bo'lsa (u holda Cgo)
- Faqat ishlaydigan kod kerak bo'lsa (u holda stdlib `slice`/`map` yetarli)

## 1.4. Umumiy yo'l (Roadmap diagrammasi)

```mermaid
flowchart TD
    Start([Boshlanish: Go o'rta daraja]) --> P1[Old shartlar: CS, OS, Memory]
    P1 --> L1["Bosqich 1: unsafe paketi"]
    L1 --> L2["Bosqich 2: reflect, runtime"]
    L2 --> L3["Bosqich 3: syscall, mmap"]

    L3 --> D1["Data Bosqich 1: Linked List, Stack, Queue"]
    D1 --> D2["Data Bosqich 2: Allocators (Bump, Pool, Slab, Buddy)"]
    D2 --> D3["Data Bosqich 3: Hash Maps"]
    D3 --> D4["Data Bosqich 4: Trees (BST, RB, B-Tree)"]
    D4 --> D5["Data Bosqich 5: Lock-free, LSM, Custom GC"]

    D5 --> Proj1[Loyiha: Mini Redis]
    Proj1 --> Proj2[Loyiha: Toy GC]
    Proj2 --> Proj3[Loyiha: In-memory DB]

    Proj3 --> End([Tamomlandi: CS bilim chuqur])

    style Start fill:#87CEEB
    style End fill:#90EE90
    style L1 fill:#FFE4B5
    style D1 fill:#FFE4B5
    style Proj1 fill:#DDA0DD
```

---

