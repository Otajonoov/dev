# 5-bob. How Go Code Turns Into Assembly

> **Bu material "The Anatomy of Go" kitobining 5-bobi asosida o'zbek tilida tayyorlangan mazmuniy tarjima va o'quv qo'llanma. Asosiy ma'no, compiler pipeline, kod misollari, buyruqlar va kitobdagi illustrationlar saqlangan; mavzular qo'shimcha diagrammalar bilan boyitilgan.**

## Bob nimani o'rgatadi?

Bu bob Go source code qanday qilib assembly, object file va final binary'ga aylanishini bosqichma-bosqich ko'rsatadi.

Katta pipeline:

```mermaid
flowchart LR
    A[Go source] --> B[Parsing AST]
    B --> C[Type checking]
    C --> D[IR construction]
    D --> E[Optimizations]
    E --> F[Walk lowering]
    F --> G[SSA generation]
    G --> H[Machine code]
    H --> I[Object archive]
    I --> J[Linking]
    J --> K[Executable binary]
```

## Mundarija

| Fayl | Mavzu | Qisqa tavsif |
|------|-------|--------------|
| [01_inspecting_build.md](01_inspecting_build.md) | Build jarayonini ko'rish | `go build -n`, compile command, importcfg, archive |
| [02_parsing_typechecking.md](02_parsing_typechecking.md) | Parsing va type checking | scanner, parser, syntax tree, scope, package declarations |
| [03_ir_export_data.md](03_ir_export_data.md) | IR va export data | IR construction, export data, relocations, package import information |
| [04_optimizations.md](04_optimizations.md) | Optimization | dead code, devirtualization, inlining, escape analysis |
| [05_walk_middle_end.md](05_walk_middle_end.md) | Walk phase | order, lowering, range/switch/string conversion simplification |
| [06_ssa.md](06_ssa.md) | SSA backend | SSA values, blocks, Phi, rewrite passes, lowering, register allocation |
| [07_codegen_linking.md](07_codegen_linking.md) | Machine code va linking | obj.Prog, assembler, relocations, archive, linker |
| [08_summary.md](08_summary.md) | Xulosa | Butun pipeline'ni bog'lash |
| [09_references.md](09_references.md) | Manbalar | Kitobda keltirilgan havolalar |

## Bobning katta savollari

1. `go build -n` nimani ko'rsatadi?
2. Source code token, syntax tree va IR'ga qanday aylanadi?
3. Type checker scope, declarations va package initialization order'ni qanday hal qiladi?
4. Export data nega Go separate compilation uchun juda muhim?
5. Devirtualization va inlining qachon ishlaydi?
6. Walk phase high-level Go construct'larini qanday soddalashtiradi?
7. SSA nima va nega compiler optimization uchun qulay?
8. Machine code generation va linker orasida relocation qanday ishlaydi?

Boshlash uchun [01_inspecting_build.md](01_inspecting_build.md) faylini oching.
