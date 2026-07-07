# 5.5 IR Construction, Export Data va Relocations

Type checking tugagach, compiler syntax tree'dan pastroq darajadagi Intermediate Representation (IR) quradi. IR optimization va keyingi lowering bosqichlari uchun qulayroq format.

![Illustration 92. Compiler builds IR from syntax and types](images/92.png)

Syntax tree source code tuzilishini ifodalaydi; IR esa compiler uchun ishlov berish osonroq, normalized representation.

## Export data nima?

Go separate compilation qiladi: package A compile bo'lishi uchun package B source code'ining hammasi shart emas. A package B'ni import qilsa, compiler B'ning export data'sini o'qiydi.

Export data quyidagilarni saqlaydi:

- exported type'lar;
- function signature'lari;
- constants;
- method set'lar;
- generic type/function metadata;
- package path va object metadata.

Source -> exportable bitstream:

![Illustration 93. From syntax to exportable bitstream](images/93.png)

String deduplication ham bor: `"foo"` bir marta saqlanadi, ko'p joyda reloc orqali refer qilinadi:

![Illustration 94. "foo" stored once, referenced multiple times](images/94.png)

Package element compact layout:

![Illustration 95. Compact layout for package element](images/95.png)

String dedup reloc index orqali:

![Illustration 96. String deduplication via reloc index](images/96.png)

`RelocString` type va index reference:

![Illustration 97. RelocString type and index reference](images/97.png)

Encoder stringlarni relocation orqali map qiladi:

![Illustration 98. foo encoder maps strings via relocations](images/98.png)

Encoded header relocation indexlarga qaraydi:

![Illustration 99. Encoded header points to relocation indices](images/99.png)

Final encoded object data:

![Illustration 100. Final encoded object data for 'foo' package element](images/100.png)

## Import qiluvchi package uchun export data

`main` package `foo`ni import qilsa, `main` compiler'i `foo` archive ichidagi export data'ni o'qiydi. Shunda `foo.Foo` signature, constants va type'lar ma'lum bo'ladi.

`main` uchun relocation value'lar:

![Illustration 102. main requires three relocation values](images/102.png)

Final encoded object data:

![Illustration 103. Final encoded object data for 'main' package](images/103.png)

Constant object relocation:

![Illustration 104. Structure of a constant's object relocation](images/104.png)

Position, type va value relocation orqali encoded bo'ladi:

![Illustration 106. Position encoded for constant a](images/106.png)

![Illustration 107. Type encoded via relocation reference](images/107.png)

![Illustration 109. Constant type encoding complete](images/109.png)

![Illustration 110. Final step: encoding the constant value](images/110.png)

![Illustration 111. Constant value encoded as int64 type](images/111.png)

Name relocation:

![Illustration 113. Name relocation for constant a](images/113.png)

Barcha relocation types:

![Illustration 114. All four relocation types for constant a](images/114.png)

## IR example: scoped if va range

Compiler high-level construct'ni IR componentlarga ajratadi. Scoped if error check:

![Illustration 115. Breakdown of scoped if error check](images/115.png)

Range loop:

![Illustration 117. Range loop broken into IR components](images/117.png)

IR construction tugagach optimization bosqichi boshlanadi:

![Illustration 118. Optimizations run after IR construction completes](images/118.png)

## Eslab qol

- IR syntax tree'dan keyingi compiler-friendly representation.
- Export data Go separate compilation imkonini beradi.
- Import qiluvchi package dependency source'ini to'liq o'qimay, export data orqali type/signature ma'lumot oladi.
- Relocation export data ichida takrorlanuvchi string/type/position/value reference'larini compact saqlashga yordam beradi.
- IR keyingi optimization va lowering bosqichlari uchun asos.
