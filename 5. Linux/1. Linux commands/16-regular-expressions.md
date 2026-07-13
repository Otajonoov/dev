# 16. Regular expressions

> Manba: TLCL 19-bob · Muhit: Ubuntu 24.04, GNU grep · [← Oldingi: archiving-and-sync](15-archiving-and-sync.md) · [Kurs xaritasi](00-README.md) · [Keyingi: text-processing →](17-text-processing.md)

## Nima uchun kerak

Log dan "timeout yoki refused bilan tugagan connection xatolari"ni ajratish, konfigda IP formatini tekshirish, Go kodingizdagi `regexp` paketi, nginx `location ~` bloklari, `sed`/`awk`/vim qidiruvlari — hammasi bitta tilda gaplashadi: **regex**. Bu dars POSIX regex ni grep orqali o'rgatadi; keyin bu bilim 17-darsdagi sed dan tortib istalgan dasturlash tiligacha ko'chadi. Regexni bilgan va bilmagan backend developerning log bilan ishlash tezligi — o'nlab barobar farq.

## Nazariya

### Regex nima va dialektlar

Regular expression — matndagi **patternlarni** tasvirlash notatsiyasi. Shell wildcardlariga (03-dars) o'xshaydi, lekin ancha kuchli va **boshqa til**: `*` wildcardda "istalgan narsa", regexda "oldingi element 0+ marta"!

POSIX ikki dialekt belgilaydi:

- **BRE** (Basic) — `grep` default: `^ $ . [ ] *` maxsus; `( ) { } ? + |` ni ishlatish uchun `\` kerak.
- **ERE** (Extended) — `grep -E`: `( ) { } ? + |` to'g'ridan-to'g'ri ishlaydi.

Amaliy tavsiya: **doim `grep -E` bilan ishlang** — sintaksis zamonaviy tillardagiga (Go `regexp`, JS, Python) yaqin, `\` changallari kam. (Perl/PCRE — undan ham boy dialekt: `\d`, lookahead; `grep -P` da ba'zi tizimlarda bor.)

### Pattern elementlari xaritasi

| Element | Ma'nosi |
|---------|---------|
| `.` | Istalgan **bitta** belgi |
| `^` / `$` | Qator boshi / oxiri (anchor) |
| `[abc]` / `[^abc]` | To'plamdan bittasi / to'plamdan boshqasi |
| `[a-z]`, `[0-9]` | Diapazon |
| `[[:alpha:]]` `[[:digit:]]` `[[:alnum:]]` `[[:space:]]` | POSIX klasslar (locale-xavfsiz) |
| `?` | Oldingi element 0 yoki 1 marta |
| `*` | 0+ marta |
| `+` | 1+ marta |
| `{n}` `{n,m}` `{n,}` | Aynan n / n..m / kamida n marta |
| `(abc)` | Guruhlash |
| `a\|b` | Yoki (alternation) |
| `\.` | Maxsus belgini literal qilish |

## Buyruqlar

### grep flaglari (regex bilan birga ishlaydiganlar)

| Flag | Nima qiladi |
|------|-------------|
| `-E` | Extended regex (ERE) |
| `-i` | Katta-kichik farqsiz |
| `-v` | Teskari (mos KELMAGAN qatorlar) |
| `-c` | Sanash |
| `-o` | Faqat topilmani chiqarish (qator emas) |
| `-n` | Qator raqami bilan |
| `-r` | Rekursiv (kataloq bo'ylab) |
| `-l` | Faqat fayl nomlari |
| `-A n` / `-B n` / `-C n` | Topilmadan keyin / oldin / atrofida n qator |
| `-w` | Butun so'z sifatida |
| `-F` | Regex EMAS, oddiy matn (tez va xavfsiz) |

Test fayllarimiz: `dirlist.txt` (`ls /usr/bin` natijasi), `phones.txt` (telefonlar), `app.log` (log qatorlari). Hammasi konteynerda verify qilingan.

### Literal va `.`

```console
$ grep zip dirlist.txt | head -3
bunzip2
bzip2
bzip2recover
$ grep ".zip" dirlist.txt | head -3
bunzip2
bzip2
bzip2recover
```

Nozik farq: `.zip` — "zip dan oldin yana bitta belgi bo'lsin" — shuning uchun `zipcloak` (zip bilan **boshlanadi**, oldida belgi yo'q) bu patternga tushmaydi.

### Anchorlar: `^` va `$`

```console
$ grep "^zip" dirlist.txt      # zip bilan BOSHLANADI
zip
zipcloak
zipdetails
$ grep "zip$" dirlist.txt      # zip bilan TUGAYDI
funzip
gunzip
gzip
$ grep "^zip$" dirlist.txt     # AYNAN zip
zip
```

Anchorsiz pattern qatorning istalgan joyida qidiriladi — aniqlik uchun deyarli har doim anchor qo'shiladi.

### Bracket to'plamlar

```console
$ grep "^[bg]zip" dirlist.txt      # b YOKI g bilan boshlanib zip
bzip2
bzip2recover
gzip
$ grep "[^b]zip" dirlist.txt       # zip, lekin oldida b BO'LMAGAN belgi
bunzip2
funzip
gunzip
gzip
streamzip
unzip
unzipsfx
```

Kutilmagan tafsilot: `bunzip2` ro'yxatda! Chunki undagi mosliq — `nzip` qismi (`n` — b emas). Inkor to'plami "b siz" degani emas, "shu pozitsiyada b dan boshqa **biror belgi bor**" degani — bo'sh joyga mos kelmaydi (shuning uchun qator boshidagi `zip` so'ziga ham tushmaydi).

Bracket ichida `^` birinchi bo'lsa — inkor; `-` chetda bo'lsa literal. Diapazon o'rniga **POSIX klasslarini** afzal ko'ring (`[A-Z]` locale ga qarab kutilmagan ishlashi mumkin): `[[:upper:]]`, `[[:digit:]]`.

### ERE: alternation, guruhlash, kvantifikatorlar

```console
$ grep -E "^(bz|gz|zip)" dirlist.txt | head -4
bzcat
bzcmp
bzdiff
bzegrep
```

Guruhsiz `^bz|gz|zip` — "bz bilan boshlanadi YOKI ichida gz bor YOKI ichida zip bor" bo'lib qolardi — qavslar aniqlik beradi.

Kvantifikatorlar bilan real misol — telefon formati `(nnn) nnn-nn-nn` ni tekshirish:

```console
$ cat phones.txt
(998) 901-23-45
99 890 123 45
(998) 555-11-22
998-90-123-45-67
(123) 456-78-90
$ grep -E "^\(?[0-9]{3}\)? [0-9]{3}-[0-9]{2}-[0-9]{2}$" phones.txt
(998) 901-23-45
(998) 555-11-22
(123) 456-78-90
```

Pattern o'qilishi: `\(?` — ochuvchi qavs ixtiyoriy (literal qavs uchun `\`), `[0-9]{3}` — 3 raqam, va h.k. Noto'g'ri formatlar filtrlandi.

### Log bilan ishlash — kundalik amaliyot

```console
$ grep -c ERROR app.log                 # nechta xato?
2
$ grep -iE "error" app.log | wc -l      # registrga qaramay (lowercase 'error' ham)
3
$ grep -E "ERROR|WARN" app.log
2026-07-10 10:00:15 ERROR db connection failed: timeout
2026-07-10 10:00:16 WARN  retrying in 5s
2026-07-10 10:01:02 ERROR db connection failed: refused
$ grep -oE "[0-9]+ms" app.log           # -o: faqat topilmalar (o'lchovlar!)
12ms
$ grep -n "ERROR" app.log               # qator raqami bilan (vim :N uchun)
2:2026-07-10 10:00:15 ERROR db connection failed: timeout
4:2026-07-10 10:01:02 ERROR db connection failed: refused
$ grep -A1 "retrying" app.log           # topilma + keyingi 1 qator (kontekst)
2026-07-10 10:00:16 WARN  retrying in 5s
2026-07-10 10:01:02 ERROR db connection failed: refused
```

### Regex boshqa joylarda ham

Xuddi shu til ishlaydi: `less` da `/pattern`, vim da `/` va `:%s`, `find -regex`, nginx `location ~ \.php$`, Go `regexp.MustCompile`. Bir marta o'rganib — hamma joyda o'qiysiz.

## Real-world scenariylar

**1. Multi-fayl xato tergovi.** Qaysi servislarda 5xx qaytgan:

```bash
grep -rlE "HTTP/1\.[01]\" 5[0-9]{2}" /var/log/nginx/
grep -c " 502 " /var/log/nginx/access.log
```

(`\.` — nuqtani literal qilish: `.` bo'lsa "1x1" ham mos kelardi!)

**2. IP manzillarni ajratib sanash** (17-darsda sort/uniq bilan davomi):

```bash
grep -oE "([0-9]{1,3}\.){3}[0-9]{1,3}" access.log | head
```

**3. Kod bazasida qidiruv.** TODO larni fayl:qator bilan; eski funksiya ishlatilgan joylar:

```bash
grep -rn "TODO|FIXME" --include="*.go" -E .
grep -rlw "OldAuthMiddleware" ./internal      # -w: aynan so'z (NewOldAuth... emas)
```

## Zamonaviy yondashuv

- **[ripgrep (rg)](https://github.com/BurntSushi/ripgrep)** — kod bazasida qidiruvning zamonaviy standarti: default rekursiv, `.gitignore` ni hurmat qiladi, binaryni o'tkazadi, 5-13x tez, default ERE-ga yaqin sintaksis (tekshirilgan):

```console
$ rg -c ERROR app.log
2
$ rg "conn.*(timeout|refused)" app.log
2026-07-10 10:00:15 ERROR db connection failed: timeout
2026-07-10 10:01:02 ERROR db connection failed: refused
```

Qoida: interaktiv kod qidiruv — `rg`; scriptlar va har-qanday-serverda ishlash — `grep` (POSIX kafolati).
- **PCRE (`grep -P`)**: `\d`, `\w`, lookahead `(?=...)` — GNU grep da bor (har platformada emas). Murakkab pattern kerak bo'lsa ko'pincha to'g'ri javob — awk/dasturlash tiliga o'tish.
- **Regex debugging**: [regex101.com](https://regex101.com) — patternni jonli tushuntirib beradi; murakkab regexni avval shu yerda quring.
- **`grep -F`** (fixed string) — pattern emas, oddiy matn qidirayotganda: tezroq va `.` `*` kabi belgilar bilan sürpriz yo'q. User inputini grep ga berayotganda **har doim** `-F` yoki qattiq quoting.

## Keng tarqalgan xatolar

1. **Shell wildcard bilan regexni aralashtirish.** `grep *.log fayl` — shell `*.log` ni fayl nomlariga ochib yuboradi! Regexda "istalgan narsa" — `.*`, va pattern **har doim quote ichida**: `grep ".*\.log" fayl`.

2. **`.` ni literal nuqta deb o'ylash.** `grep "1.5" versions.txt` — "1x5", "125" ham topiladi. Literal: `\.` yoki `-F`.

3. **BRE da `+`/`?` ishlamay "regex buzildi" deyish.** `grep "ab+"` BRE da literal "+" ni qidiradi. Yechim: `grep -E` (yoki BRE da `\+`).

4. **Anchorsiz validatsiya.** `grep -E "[0-9]{3}"` — "abc12345xyz" ham o'tadi (ichida 3 raqam bor-ku). Format tekshiruvida doim `^...$`.

5. **Alternation qamrovi.** `grep -E "^ERROR|WARN"` — bu "(^ERROR) yoki (WARN istalgan joyda)". To'g'ri: `^(ERROR|WARN)`.

6. **Greedy `.*` bilan ortiqcha qamrab olish.** `".*"` pattern `say "a" and "b"` da `"a" and "b"` ni butun oladi (eng uzun mos). POSIX grep da lazy `*?` yo'q — aniqroq pattern yozing: `"[^"]*"` (qo'shtirnoq ichida qo'shtirnoq bo'lmagan belgilar).

## Amaliy mashqlar

Muhit: `docker run -it --rm ubuntu:24.04 bash`. Tayyorlov:

```bash
ls /usr/bin > dirlist.txt
printf '2026-07-10 10:00:15 ERROR db timeout\n2026-07-10 10:00:16 WARN retry\n2026-07-10 10:01:02 ERROR db refused\nGET /api/users 200 12ms\nGET /health 200 1ms\n' > app.log
```

**1.** `dirlist.txt` da: (a) `gz` bilan boshlanadigan; (b) `2` bilan tugaydigan; (c) roppa-rosa 3 belgili buyruqlarni toping.

<details><summary>Yechim</summary>

```bash
grep "^gz" dirlist.txt
grep "2$" dirlist.txt
grep -E "^.{3}$" dirlist.txt | head
```
</details>

**2.** app.log dan faqat millisekund qiymatlarini (raqam+ms) ajratib oling.

<details><summary>Yechim</summary>

```console
$ grep -oE "[0-9]+ms" app.log
12ms
1ms
```
</details>

**3.** "ERROR yoki WARN bilan davom etadigan log qatorlari" — lekin sana bilan boshlanganlarigagina mos keladigan bitta pattern yozing.

<details><summary>Yechim</summary>

```bash
grep -E "^[0-9]{4}-[0-9]{2}-[0-9]{2} .* (ERROR|WARN)" app.log
# yoki qat'iyroq:
grep -E "^[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9:]{8} (ERROR|WARN)" app.log
```
</details>

**4.** Semantik versiya (`v1.2.3` formati) ni tekshiradigan pattern yozing va sinang: `v1.2.3` o'tsin, `v1.2`, `1.2.3`, `v1.2.3-rc` o'tmasin.

<details><summary>Yechim</summary>

```console
$ printf "v1.2.3\nv1.2\n1.2.3\nv1.2.3-rc\n" | grep -E "^v[0-9]+\.[0-9]+\.[0-9]+$"
v1.2.3
```
`^...$` anchorlar va `\.` literal nuqtalar — ikkala tuzoqdan himoya.
</details>

**5.** `dirlist.txt` dagi tarkibida raqam BO'LMAGAN buyruqlar sonini toping.

<details><summary>Yechim</summary>

```bash
grep -vE "[0-9]" dirlist.txt | wc -l
# yoki POSIX klass bilan: grep -v "[[:digit:]]"
```
</details>

**6.** grep kontekst flaglari bilan: app.log da "WARN" topilmasining oldingi va keyingi qatorini birga chiqaring.

<details><summary>Yechim</summary>

```bash
grep -C1 WARN app.log
```
`-C1` = `-B1 -A1`. Stack trace tergovlarida `-A20` juda foydali.
</details>

**7.** (Qiyinroq) `/etc/passwd` dan login shelli `nologin` yoki `false` bo'lgan userlar sonini bitta grep bilan toping (format: oxirgi maydon).

<details><summary>Yechim</summary>

```console
$ grep -cE "(nologin|false)$" /etc/passwd
```
Anchor `$` — aynan qator oxirida (oxirgi maydon shell); alternation guruh ichida.
</details>

## Cheat sheet

| Pattern / flag | Ma'nosi | Misol |
|----------------|---------|-------|
| `^` `$` | Qator boshi/oxiri | `^ERROR`, `\.go$` |
| `.` / `\.` | Istalgan belgi / literal nuqta | `1\.2\.3` |
| `[abc]` `[^abc]` | To'plam / inkor | `^[bg]zip` |
| `[[:digit:]]` | POSIX klass | `[[:upper:]]` |
| `?` `*` `+` | 0-1 / 0+ / 1+ | `https?` |
| `{3}` `{1,3}` | Takror soni | `[0-9]{1,3}` |
| `(a\|b)` | Guruh + yoki | `^(GET\|POST)` |
| `grep -E` | ERE (tavsiya) | `grep -E "a\|b"` |
| `grep -o` | Faqat topilma | `-oE "[0-9]+ms"` |
| `grep -v` | Teskari | `-v DEBUG` |
| `grep -rn` | Rekursiv + raqam | kod qidiruv |
| `grep -A/-B/-C` | Kontekst | `-C2 panic` |
| `grep -F` | Regexsiz matn | user input uchun |
| `rg` | Zamonaviy grep | `rg pattern` (rekursiv default) |

## Qo'shimcha manbalar

- [regex101.com](https://regex101.com) — patternlarni jonli qurish va tushunish
- [learnbyexample: GNU grep & ripgrep](https://learnbyexample.github.io/learn_gnugrep_ripgrep/) — bepul, chuqur amaliy kitob
- [RegexOne](https://regexone.com/) — interaktiv regex darslik

---

[← Oldingi: 15 — archiving-and-sync](15-archiving-and-sync.md) · [Kurs xaritasi](00-README.md) · [Keyingi: 17 — text-processing →](17-text-processing.md)
