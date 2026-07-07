# 5.2-5.4 Parsing, Syntax Tree va Type Checking

Go package compile bo'lishi bir nechta bosqichdan o'tadi. Dastlab source text tokenlarga ajraladi, keyin syntax tree quriladi, so'ng type checker Go qoidalarini tekshiradi.

## 5.2 Package compile overview

Package compile qilishning soddalashtirilgan ko'rinishi:

```mermaid
flowchart LR
    A[Source files] --> B[Scanner]
    B --> C[Tokens]
    C --> D[Parser]
    D --> E[Syntax tree]
    E --> F[Type checker]
    F --> G[Typed syntax]
    G --> H[IR]
```

## 5.3 Stage 1: parsing va syntax tree

Scanner source code'ni tokenlarga ajratadi. Masalan:

```go
if err != nil {
    return err
}
```

Tokenlar:

- `if`
- identifier `err`
- operator `!=`
- identifier/predeclared `nil`
- `{`
- `return`
- identifier `err`
- `}`

Kitobdagi token ko'rinishi:

![Illustration 77. Tokens corresponding to the if block's content and closing.](images/77.png)

Parser tokenlardan syntax tree quradi:

![Illustration 78. Visual syntax tree of a Go program](images/78.png)

`if` statement strukturasi:

![Illustration 79. Structure of an if statement in Go](images/79.png)

Syntax tree hali "bu type to'g'rimi?", "variable mavjudmi?", "assignment mosmi?" kabi savollarga to'liq javob bermaydi. U faqat source code grammatik structure'ni ushlaydi.

## 5.4 Stage 2: type checking va scope resolution

Type checker syntax tree ustidan yurib, Go language qoidalarini tekshiradi:

![Illustration 80. Type checking stage in Go compiler](images/80.png)

U quyidagilarni bajaradi:

- package name consistency;
- import resolution;
- top-level declarations yig'ish;
- scope qurish;
- type declaration, const, var va function signature tekshirish;
- function body type checking;
- package-level initialization order aniqlash.

Package-level declarations:

![Illustration 81. Package-level declarations: const, type, and func](images/81.png)

Constant declaration tree:

![Illustration 82. Constant declaration tree for global scope](images/82.png)

Package scope object'larni saqlaydi:

![Illustration 83. Package scope with declared objects](images/83.png)

Compiler type-checking state'larni rang/kod kabi kuzatadi:

![Illustration 84. Type-checking states tracked using color codes](images/84.png)

## Scope va dependency resolution

Type checker variable yoki type nomini ko'rganda current scope'dan boshlaydi, topilmasa parent scope'ga ko'tariladi. Bu 2-bobdagi scope modelining compiler pipeline ichidagi davomidir.

Type dependency'lar ham hal qilinadi:

![Illustration 85. Type checker resolving variable dependencies](images/85.png)

Circular reference bo'lsa ham compiler ma'lum darajada resolution'ni davom ettirishi va xatoni aniqroq chiqarishi mumkin:

![Illustration 86. Type resolution proceeds despite circular references](images/86.png)

Global type/signature tekshiruvlardan keyin function body'lar to'liq tekshiriladi:

![Illustration 87. From global types to full function checks](images/87.png)

## Package-level initialization order

Global variable'lar bir-biriga bog'liq bo'lishi mumkin:

```go
var a = f()
var b = a + 1
var c = b + 1
```

Compiler dependency graph tuzadi:

![Illustration 88. Resolving package-level initialization dependencies](images/88.png)

Graph:

![Illustration 89. Dependency graph of global object initialization](images/89.png)

`f()` downstream dependency bilan almashtirilishi mumkin:

![Illustration 90. f() replaced with its downstream links](images/90.png)

Agar dependency cycle qolsa, initialization order mumkin emas:

![Illustration 91. Cycle detected if dependencies remain](images/91.png)

## Eslab qol

- Scanner source text'ni tokenlarga ajratadi.
- Parser tokenlardan syntax tree quradi.
- Type checker Go qoidalari, scope va type compatibility'ni tekshiradi.
- Package-level declarations avval yig'iladi, keyin body'lar tekshiriladi.
- Global initialization dependency graph orqali tartiblanadi.
