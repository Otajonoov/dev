# 19. Branching: if, test va case

> Manba: TLCL 27 va 31-boblar · Muhit: Ubuntu 24.04, bash 5.2 · [← Oldingi: scripting-first-steps](18-scripting-first-steps.md) · [Kurs xaritasi](00-README.md) · [Keyingi: loops →](20-loops.md)

## Nima uchun kerak

"Backup muvaffaqiyatli bo'lsa — eskisini o'chir, bo'lmasa — alert yubor" — shartsiz script shunchaki buyruqlar ro'yxati. CI/CD ning butun mantig'i exit code larga qurilgan: `docker build` yiqilsa pipeline to'xtaydi. Bu darsda bash dagi shartlarning **uch dunyosi**ni ajratamiz: `[ ]` (klassik test), `[[ ]]` (zamonaviy bash), `(( ))` (arifmetika) — va qaysi birini qachon ishlatishni bir marta uzil-kesil hal qilamiz.

## Nazariya

### Hamma narsa exit code ga qurilgan

Har buyruq tugaganda **exit status** qaytaradi: `0` = muvaffaqiyat, `1-255` = xato (08-darsdagi 137 esingizdami?). Oxirgi buyruqniki — `$?` da (tekshirilgan):

```console
$ ls -d /usr/bin >/dev/null; echo $?
0
$ ls -d /bin/usr; echo $?
ls: cannot access '/bin/usr': ...
2
$ true; echo $?; false; echo $?
0
1
```

**`if` hech qanday "shart"ni bilmaydi — u shunchaki buyruq bajarib, exit code ga qaraydi:**

```bash
if buyruq; then
    ...
elif boshqa_buyruq; then
    ...
else
    ...
fi
```

`[ ... ]` — ham buyruq (aslida `test` ning sinonimi, `/usr/bin/[` faylini 05-darsda ko'rgan edik!). Shuning uchun `if grep -q ERROR log`, `if ping -c1 host` to'g'ridan-to'g'ri ishlaydi — testga o'rash shart emas.

### Uch sintaksis — qachon qaysi biri

| Sintaksis | Nima | Qachon |
|-----------|------|--------|
| `[ ... ]` / `test` | POSIX klassika | sh-portativlik kerak bo'lsa |
| `[[ ... ]]` | Bash kengaytmasi: pattern (`==`), regex (`=~`), xavfsizroq quoting | **bash scriptda default tanlov** |
| `(( ... ))` | Arifmetik kontekst: `<` `>` `==` odatiy ma'noda | sonlar bilan ishlaganda |

## Buyruqlar

### Fayl testlari (eng ko'p ishlatiladiganlar)

| Test | Ma'nosi |
|------|---------|
| `-e` | Mavjud |
| `-f` | Oddiy fayl |
| `-d` | Katalog |
| `-r` / `-w` / `-x` | O'qish / yozish / bajarish mumkin |
| `-s` | Mavjud va bo'sh emas |
| `-L` | Symlink |
| `f1 -nt f2` | f1 f2 dan yangi (newer than) |

Tekshirilgan:

```console
$ touch mavjud.txt
$ [ -f mavjud.txt ] && echo "-f: oddiy fayl"
-f: oddiy fayl
$ [ -s mavjud.txt ] || echo "-s: fayl bo'sh"
-s: fayl bo'sh
$ [ ! -x mavjud.txt ] && echo "! -x: executable emas"
! -x: executable emas
```

(`!` — inkor; `&&`/`||` bilan qisqa if — pastda.)

### String testlari

| Test | Ma'nosi |
|------|---------|
| `-n "$s"` | Bo'sh emas |
| `-z "$s"` | Bo'sh (zero length) |
| `"$s1" = "$s2"` | Teng (`==` ham bo'ladi `[[ ]]` da) |
| `"$s1" != "$s2"` | Teng emas |

```console
$ s="salom"
$ [ -n "$s" ] && echo "-n: bo'sh emas"
-n: bo'sh emas
$ [ -z "$bosh_var" ] && echo "-z: bo'sh/mavjud emas"
-z: bo'sh/mavjud emas
```

**Quote — hayotiy zarurat** (tekshirilgan xavf):

```console
$ unset bosh
$ [ $bosh = "x" ]
bash: [: =: unary operator expected      # variable g'oyib bo'lib sintaksis buzildi!
$ [ "$bosh" = "x" ] || echo "quote bilan: shunchaki false"
quote bilan: shunchaki false
```

### Integer testlari

`[ ]` ichida: `-eq -ne -lt -le -gt -ge` (equal, not-equal, less-than...):

```console
$ n=7
$ [ "$n" -gt 5 ] && [ "$n" -lt 10 ] && echo "5 < n < 10"
5 < n < 10
```

`(( ))` da — odatiy matematik belgilar (o'qish osonroq):

```console
$ ((n > 5)) && echo "(( )) arifmetik: n>5"
(( )) arifmetik: n>5
```

### `[[ ]]` — zamonaviy test: pattern va regex

```console
$ v="v2.4.1"
$ [[ "$v" == v2.* ]] && echo "pattern match: v2.x"
pattern match: v2.x
$ [[ "$v" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]] && echo "semver OK"
semver OK
```

`==` o'ng tomonida quotesiz pattern — glob sifatida (03-dars); `=~` — ERE regex (16-dars to'g'ridan-to'g'ri scriptda!). Bonus: `[[ ]]` ichida word splitting yo'q — unquoted variable ham sinmaydi (baribir quote qiling — odat buzilmasin).

### `elif` zanjiri — amaliy misol

```bash
check_temp() {
    local t=$1
    if [ "$t" -ge 80 ]; then
        echo "KRITIK"
    elif [ "$t" -ge 60 ]; then
        echo "ogohlantirish"
    else
        echo "normal"
    fi
}
```

```console
$ check_temp 95; check_temp 65; check_temp 30
KRITIK
ogohlantirish
normal
```

### `&&` va `||` — bir qatorlik branching

05-darsdagi control operatorlar shartli mantiqning qisqa shakli:

```console
$ mkdir -p tmpdir && echo "yaratildi"        # oldingisi OK bo'lsa
yaratildi
$ [ -d yoq-katalog ] || echo "yo'q ekan"     # oldingisi yiqilsa
yo'q ekan
```

Idiomatik patternlar:

```bash
command -v jq >/dev/null || { echo "jq kerak" >&2; exit 1; }
[ -f .env ] && source .env
cd /srv/app || exit 1          # cd yiqilsa davom etmaslik — MUHIM!
```

### `case` — ko'p variantli tanlov

O'nlab `elif` o'rniga — pattern bo'yicha tanlash:

```bash
case "$1" in
    *.go)            echo "Go fayli" ;;
    *.md|*.txt)      echo "hujjat" ;;        # | — "yoki"
    *.tar.gz|*.tgz)  echo "arxiv" ;;
    .*)              echo "hidden" ;;
    *)               echo "nomalum" ;;       # default (har doim oxirida)
esac
```

```console
$ check_ext main.go; check_ext README.md; check_ext app.tgz; check_ext .bashrc; check_ext data.bin
Go fayli
hujjat
arxiv
hidden
nomalum
```

Patternlar — pathname expansion qoidalari (`*`, `?`, `[...]`, `[[:class:]]`). Birinchi mos kelgan g'olib — `;;` dan keyin chiqadi. Bash 4+ da `;;&` — moslikdan keyin **davom etish** (tekshirilgan):

```console
$ c="a"
$ case "$c" in
>   [[:alpha:]]) echo "harf" ;;&
>   [aeiou])     echo "unli" ;;
> esac
harf
unli
```

## Real-world scenariylar

**1. Deploy scriptdagi himoya zanjiri:**

```bash
#!/usr/bin/env bash
set -euo pipefail

[ -f "./app" ] || { echo "binary topilmadi — avval build qiling" >&2; exit 1; }

if [[ "$(uname -m)" != "x86_64" ]]; then
    echo "OGOHLANTIRISH: kutilmagan arxitektura" >&2
fi

if curl -sf --max-time 5 localhost:8080/healthz >/dev/null; then
    echo "eski versiya ishlayapti — graceful restart"
    systemctl reload myapp
else
    systemctl start myapp
fi
```

**2. Environment bo'yicha tarmoqlanish (case):**

```bash
case "${DEPLOY_ENV:-dev}" in
    prod)        replicas=5; log_level=warn ;;
    staging)     replicas=2; log_level=info ;;
    dev|local)   replicas=1; log_level=debug ;;
    *)           echo "nomalum muhit: $DEPLOY_ENV" >&2; exit 1 ;;
esac
```

**3. Healthcheck + retry mantiq asosi** (loop 20-darsda to'ldiriladi):

```bash
if ! pg_isready -h "$DB_HOST" -q; then
    echo "DB javob bermayapti" >&2
    exit 1
fi
```

## Zamonaviy yondashuv

- **Default tanlovlar**: bash script → `[[ ]]` (string/fayl/pattern) va `(( ))` (sonlar); `#!/bin/sh` portativ script → faqat `[ ]`. Aralashtirib yurmang — bitta scriptda bitta uslub.
- **`grep -q`, `curl -sf` kabi "jim" rejimlar** if bilan juftlikda: bizni output emas, exit code qiziqtiradi.
- **ShellCheck** bu darsning xatolarini avtomatik topadi: SC2086 (quote), SC2166 (`-a/-o` o'rniga `&&`/`||`), SC2181 (`$?` ni if ga o'rash o'rniga to'g'ridan-to'g'ri `if cmd`).
- **`if [ $? -eq 0 ]` anti-pattern**: buyruqni to'g'ridan-to'g'ri if ga qo'ying: `if cmd; then`. `$?` faqat exit code ni saqlab keyinroq ishlatish kerak bo'lganda.
- Eski `[ -a ]`/`[ -o ]` (and/or ichkarida) — deprecated: alohida `[ ] && [ ]` yoki `[[ ... && ... ]]`.

## Keng tarqalgan xatolar

1. **`[$x = 5]` — probellarsiz.** `[` — buyruq, argumentlari probel bilan ajratilishi shart: `[ "$x" = 5 ]`. Xuddi shunday `if[` ham ishlamaydi.

2. **Testda variable ni quotesiz qoldirish.** Yuqorida ko'rsatildi: bo'sh variable `unary operator expected` beradi. `[ ]` ichida **har doim** `"$var"`. (`[[ ]]` kechiradi, lekin odatni buzmang.)

3. **String bilan sonni adashtirish: `[ "$n" = 5 ]` vs `[ "$n" -eq 5 ]`.** `=` — string taqqoslash: `[ "05" = "5" ]` — **false**! Sonlar uchun `-eq` yoki `(( ))`.

4. **`[ ]` ichida `<` `>` ishlatish.** `[ 2 > 10 ]` — bu **redirect**! (`10` nomli fayl yaratiladi va test true qaytadi.) Sonlar: `-lt`/`-gt` yoki `((2 > 10))`.

5. **`cd` xatosini tekshirmaslik.** `cd $dir; rm -rf *` — cd yiqilsa (katalog yo'q), rm **joriy katalogda** ishlaydi. Har doim: `cd "$dir" || exit 1`.

6. **case da default (`*)`) ni unutish.** Kutilmagan qiymat indamay o'tib ketadi. Har case da oxirgi variant — `*)` bilan xato/log.

## Amaliy mashqlar

Muhit: `docker run -it --rm ubuntu:24.04 bash`

**1.** Bir qatorlik iflar: fayl mavjud bo'lsa "bor", bo'lmasa "yo'q" — `&&`/`||` bilan; keyin to'liq if/else bilan.

<details><summary>Yechim</summary>

```bash
[ -f /etc/passwd ] && echo bor || echo yo'q
# to'liq (va aniqroq — && || zanjirida yon effekt xavfi bor):
if [ -f /etc/passwd ]; then echo bor; else echo "yo'q"; fi
```
</details>

**2.** `check_file` funksiyasi: argumentdagi yo'l uchun "katalog / oddiy fayl / symlink / mavjud emas" deb aytsin.

<details><summary>Yechim</summary>

```bash
check_file() {
    if   [ -L "$1" ]; then echo "symlink"       # MUHIM: -L tekshiruvi -f/-d dan OLDIN
    elif [ -d "$1" ]; then echo "katalog"
    elif [ -f "$1" ]; then echo "oddiy fayl"
    else echo "mavjud emas"; fi
}
check_file /etc; check_file /etc/passwd; check_file /bin; check_file /yoq
```
(`-L` avval — chunki symlink `-f`/`-d` testlaridan ham o'tadi, target orqali.)
</details>

**3.** Foydalanuvchi kiritgan qiymat ("$1") 1-100 oralig'idagi son ekanini tekshiring: regex bilan son ekani, `(( ))` bilan oraliq.

<details><summary>Yechim</summary>

```bash
n="$1"
if [[ "$n" =~ ^[0-9]+$ ]] && ((n >= 1 && n <= 100)); then
    echo "OK: $n"
else
    echo "1-100 oralig'ida son kiriting" >&2; exit 1
fi
```
Regex avval — aks holda `((n >= 1))` matnda sintaksis xato beradi.
</details>

**4.** `[ "05" = "5" ]` va `[ "05" -eq "5" ]` natijalarini tekshirib farqni tushuntiring.

<details><summary>Yechim</summary>

```console
$ [ "05" = "5" ]; echo $?
1        # string sifatida har xil
$ [ "05" -eq "5" ]; echo $?
0        # son sifatida teng
```
Port raqamlari, ID lar bilan ishlaganda bu farq real buglar manbai.
</details>

**5.** case bilan mini-router: `$1` ga qarab `start|stop|restart|status` amallarini echo qilsin, noma'lum buyruqda usage ko'rsatib exit 1 qilsin. `st*` kabi pattern ishlatmang — aniq variantlar.

<details><summary>Yechim</summary>

```bash
case "${1:-}" in
    start)    echo "ishga tushirilmoqda..." ;;
    stop)     echo "to'xtatilmoqda..." ;;
    restart)  echo "qayta..." ;;
    status)   echo "holat: OK" ;;
    *)        echo "Usage: $0 {start|stop|restart|status}" >&2; exit 1 ;;
esac
```
Bu — klassik init-script qolipi; systemd dan oldingi barcha servislar shunday boshqarilardi.
</details>

**6.** Internetni tekshiruvchi qism yozing: `ping -c1 8.8.8.8` jim ishlasa "tarmoq bor", bo'lmasa "yo'q" — `if` ga buyruqni **to'g'ridan-to'g'ri** qo'yib ($? siz).

<details><summary>Yechim</summary>

```bash
if ping -c1 -W2 8.8.8.8 >/dev/null 2>&1; then
    echo "tarmoq bor"
else
    echo "tarmoq yo'q"
fi
```
</details>

**7.** (Qiyinroq) `backup.sh` skeleti: (a) argument berilmagan bo'lsa usage bilan chiqsin; (b) manba katalog mavjudligini tekshirsin; (c) maqsad faylı allaqachon mavjud bo'lsa `.old` ga rename qilsin; (d) tar yaratib, exit codeni tekshirib xabar bersin.

<details><summary>Yechim</summary>

```bash
#!/usr/bin/env bash
set -euo pipefail

src="${1:-}"
[ -n "$src" ] || { echo "Usage: $0 <katalog>" >&2; exit 1; }
[ -d "$src" ] || { echo "XATO: '$src' katalog emas" >&2; exit 1; }

dest="/tmp/$(basename "$src")-backup.tgz"
[ -f "$dest" ] && mv "$dest" "$dest.old"

if tar czf "$dest" -C "$(dirname "$src")" "$(basename "$src")"; then
    echo "OK: $dest ($(du -h "$dest" | cut -f1))"
else
    echo "XATO: arxivlash yiqildi" >&2; exit 1
fi
```
</details>

## Cheat sheet

| Konstruksiya | Nima | Misol |
|--------------|------|-------|
| `$?` | Oxirgi exit code | 0 = OK |
| `if cmd; then...fi` | Buyruq exit codeiga qarab | `if grep -q x f; then` |
| `[ -f/-d/-e/-s "$f" ]` | Fayl testlari | quote majburiy! |
| `[ -n/-z "$s" ]` | String bo'sh(mas)ligi | — |
| `[ "$a" = "$b" ]` | String tenglik | sonlar uchun EMAS |
| `[ "$n" -eq/-lt/-gt N ]` | Integer | yoki `(( ))` |
| `[[ "$s" == pat* ]]` | Glob pattern | bash only |
| `[[ "$s" =~ regex ]]` | ERE regex | bash only |
| `((n > 5))` | Arifmetika | odatiy belgilar |
| `cmd && ok \|\| fail` | Qisqa shart | yon effektga ehtiyot |
| `case "$x" in pat) ...;; esac` | Ko'p tarmoq | oxirida `*)` |

## Qo'shimcha manbalar

- [Bash Reference — Conditional Constructs](https://www.gnu.org/software/bash/manual/html_node/Conditional-Constructs.html) — rasmiy hujjat
- [BashGuide: Tests and Conditionals](https://mywiki.wooledge.org/BashGuide/TestsAndConditionals) — chuqur va aniq qo'llanma
- [ShellCheck wiki](https://www.shellcheck.net/wiki/) — shartlardagi klassik xatolar katalogi

---

[← Oldingi: 18 — scripting-first-steps](18-scripting-first-steps.md) · [Kurs xaritasi](00-README.md) · [Keyingi: 20 — loops →](20-loops.md)
