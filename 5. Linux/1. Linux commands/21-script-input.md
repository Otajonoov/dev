# 21. Script input: read va positional parametrlar

> Manba: TLCL 28 va 32-boblar · Muhit: Ubuntu 24.04, bash 5.2 · [← Oldingi: loops](20-loops.md) · [Kurs xaritasi](00-README.md) · [Keyingi: strings-numbers-arrays →](22-strings-numbers-arrays.md)

## Nima uchun kerak

Hozirgacha scriptlarimiz "kar" edi — tashqaridan ma'lumot olmasdi. Haqiqiy tool esa argument oladi: `deploy.sh --env prod --version v2.4.1`. Bu darsda scriptni **CLI dasturga** aylantiramiz: positional parametrlar, `shift`, flag parsing (o'sha `-v --output` lar qanday ishlashi), interaktiv `read` va input validatsiyasi. Go dagi `os.Args` va `flag` paketining shell tomondagi ekvivalenti.

## Nazariya

### Ikki xil input kanali

1. **Argumentlar** — chaqiruv paytida: `./script arg1 arg2` → `$1`, `$2` (Go: `os.Args`). Avtomatlashtirish uchun asosiy usul.
2. **stdin** — ish paytida: `read` bilan (interaktiv savollar) yoki pipe/redirect dan.

Qoida: script **avtomatlashtiriladigan** bo'lsin — hamma narsani argumentdan oling, `read` ni faqat haqiqatan interaktiv stsenariylar uchun saqlang (CI da savol so'raydigan script — osilib qoladigan script).

## Buyruqlar

### Positional parametrlar (tekshirilgan)

```bash
#!/usr/bin/env bash
echo "script nomi (\$0): $0"
echo "1-argument (\$1): ${1:-berilmagan}"
echo "argumentlar soni (\$#): $#"
echo "hammasi (\$@): $@"
```

```console
$ ./posit.sh olma anor "uzum shingili"
script nomi ($0): ./posit.sh
1-argument ($1): olma
argumentlar soni ($#): 3
hammasi ($@): olma anor uzum shingili
$ ./posit.sh
1-argument ($1): berilmagan
argumentlar soni ($#): 0
```

- `$0` — script nomi, `$1`..`$9` (undan yuqori: `${10}`), `$#` — soni
- `${1:-default}` — 1-argument bo'sh bo'lsa default (22-darsda bu oila to'liq)
- Funksiyalar ichida ham xuddi shu `$1`, `$#` — **funksiya argumentlari** (18-darsdagi savolga javob)

### `"$@"` vs `"$*"` — muhim farq (tekshirilgan)

```console
$ ./at-vs-star.sh bir ikki "uch to'rt"
--- "$@" (TO'G'RI):
[bir]
[ikki]
[uch to'rt]
--- "$*" (bitta string):
[bir ikki uch to'rt]
```

`"$@"` — har argument alohida, quoting saqlanadi (**99% holatda kerakligi shu**); `"$*"` — hammasi bitta stringga yelimlanadi (faqat log xabari kabi joylarda). Argumentlarni boshqa buyruqqa uzatish: har doim `"$@"` — 18-darsdagi `main "$@"` shundan.

### `shift` — argumentlarni "yeyish"

Har `shift` da `$2`→`$1`, `$3`→`$2`... va `$#` kamayadi (tekshirilgan):

```console
$ ./shift.sh a b c
qoldi 3: birinchisi=a
qoldi 2: birinchisi=b
qoldi 1: birinchisi=c
```

### Option parsing — CLI tool qolipi

`while` + `case` + `shift` uchligi (to'liq tekshirilgan):

```bash
#!/usr/bin/env bash
verbose=0; out=""
while [ $# -gt 0 ]; do
    case "$1" in
        -v|--verbose) verbose=1 ;;
        -o|--output)  shift; out="$1" ;;      # qiymatli flag: yana bir shift
        -h|--help)    echo "Usage: $0 [-v] [-o FILE] args..."; exit 0 ;;
        -*)           echo "noma'lum flag: $1" >&2; exit 1 ;;
        *)            break ;;                # flag emas — pozitsion argumentlar boshlandi
    esac
    shift
done
echo "verbose=$verbose out=$out qolgan argumentlar: $*"
```

```console
$ ./opts.sh -v -o natija.txt fayl1 fayl2
verbose=1 out=natija.txt qolgan argumentlar: fayl1 fayl2
$ ./opts.sh --help
Usage: ./opts.sh [-v] [-o FILE] args...
$ ./opts.sh -x
noma'lum flag: -x
```

Bu qolip — deyarli barcha professional bash toollarning skeleti. (Bash builtin `getopts` ham bor — faqat qisqa flaglar uchun; uzun flaglar kerak bo'lgani uchun amalda ko'pincha yuqoridagi qo'lda parsing ishlatiladi.)

### `read` — stdin dan o'qish

```console
$ echo "42" | { read -r javob; echo "o'qildi: $javob"; }
o'qildi: 42
$ printf "olma anor uzum\n" | { read -r m1 m2; echo "m1=$m1 m2=[$m2]"; }
m1=olma m2=[anor uzum]
```

Muhim semantika: variablelar soni yetmasa, **oxirgisi qolgan hammasini oladi**. Foydali flaglar:

| Flag | Nima |
|------|------|
| `-r` | Backslash literal (HAR DOIM) |
| `-p "matn"` | Prompt ko'rsatish |
| `-s` | Yashirin kiritish (parollar) |
| `-t N` | N soniya timeout |
| `-n N` | N belgi bilan cheklash |

```console
$ echo "Alisher" | ./readdemo.sh      # read -rp "Ism kiriting > " ism
Salom, Alisher!
```

### IFS bilan maydonlash

`IFS` — read ning ajratuvchisi. Bitta buyruq uchun vaqtincha o'zgartirish (20-darsdagi mashqdan tanish):

```console
$ IFS=: read -r user pw uid gid <<< "root:x:0:0"
$ echo "user=$user uid=$uid"
user=root uid=0
```

(`<<<` — here string: bitta qatorni stdin ga beradi.)

### Validatsiya — ishonmang, tekshiring

```bash
validate() {
    local n="$1"
    [[ "$n" =~ ^[0-9]+$ ]] || { echo "'$n' son emas" >&2; return 1; }
    echo "OK: $n"
}
```

```console
$ validate 42; validate abc
OK: 42
'abc' son emas
```

Qoida: tashqaridan kelgan **har qanday** qiymat (argument, read, env) ishlatilishidan oldin tekshiriladi — ayniqsa u keyin `rm`, `ssh`, SQL ga ketadigan bo'lsa.

## Real-world scenariylar

**1. Professional deploy script interfeysi:**

```bash
#!/usr/bin/env bash
set -euo pipefail

usage() {
    cat <<EOF
Usage: $0 --env <dev|staging|prod> --version <vX.Y.Z> [--dry-run]
EOF
    exit 1
}

env=""; version=""; dry_run=0
while [ $# -gt 0 ]; do
    case "$1" in
        --env)     shift; env="$1" ;;
        --version) shift; version="$1" ;;
        --dry-run) dry_run=1 ;;
        -h|--help) usage ;;
        *) echo "noma'lum: $1" >&2; usage ;;
    esac
    shift
done

[[ "$env" =~ ^(dev|staging|prod)$ ]] || { echo "yaroqsiz env: '$env'" >&2; usage; }
[[ "$version" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]] || { echo "yaroqsiz versiya" >&2; usage; }
```

**2. Xavfli amaldan oldin tasdiqlash:**

```bash
read -rp "PRODUCTION bazasi o'chiriladi. Davom etish uchun 'ha' yozing: " javob
[ "$javob" = "ha" ] || { echo "bekor qilindi"; exit 1; }
```

**3. Parolni xavfsiz so'rash:**

```bash
read -rsp "DB parol: " db_pass; echo
export PGPASSWORD="$db_pass"
```

(`-s` — terminalda ko'rinmaydi; argument sifatida berilgan parol esa `ps` da hammaga ko'rinadi — shuning uchun read!)

## Zamonaviy yondashuv

- **Har toolga `-h/--help`** — hatto 20 qatorlik shaxsiy scriptga ham. 3 oydan keyingi o'zingiz minnatdor bo'lasiz. `usage()` funksiya + heredoc — standart shakl.
- **12-factor uyg'unligi**: konfiguratsiya prioriteti — flag > env variable > default: `env="${DEPLOY_ENV:-dev}"` bilan boshlab, flag bo'lsa ustidan yozish.
- **CI muhitida `read` bloklanadi** — script interaktivlikni tekshirsin: `[ -t 0 ]` (stdin terminal ekanini bilish) yoki `--yes` flagi bilan tasdiqlarni chetlab o'tish imkonini bering.
- Katta CLI toollar uchun bash o'rniga Go (`cobra`) — subcommand, autocompletion, tiplar. Bash parsing — ~5 flaggacha ideal.

## Keng tarqalgan xatolar

1. **`$@` ni quotesiz ishlatish.** Probelli argumentlar bo'linib ketadi — huddi `"$@"` dagi farq jadvalida ko'rsatilganidek. Har doim `"$@"`.

2. **`$1` mavjudligini tekshirmaslik.** `set -u` bilan script "unbound variable" bilan yiqiladi, usiz — jimgina bo'sh string bilan davom etib g'alati xato beradi. Standart: `src="${1:?Usage: $0 <src>}"` yoki aniq `[ $# -ge 1 ] ||usage`.

3. **Qiymatli flagda ikkinchi shiftni unutish.** `-o) out="$1"` (shiftsiz) — out ga `-o` ning o'zi tushadi. Case ichida `shift; out="$1"` yoki `out="$2"; shift 2` — bitta uslub tanlab doim shunga rioya qiling.

4. **`read` ga `-r` siz odatlanish.** `C:\new\table` kabi input `C:new able` bo'lib qoladi. Muscle memory: `read -r`.

5. **Parolni argument qilib olish.** `./script --password secret123` — `ps aux` da, shell history da ochiq qoladi. Parol/token: `read -s`, env variable yoki fayl orqali.

6. **`$*` bilan argumentlarni uzatish.** `child_script $*` — probellar buziladi. To'g'ri: `child_script "$@"`.

## Amaliy mashqlar

Muhit: `docker run -it --rm ubuntu:24.04 bash`

**1.** `args-info` scripti: nechta argument kelganini, birinchisini va hammasini ko'rsatsin; argumentsiz chaqirilsa usage chiqarib exit 1 qilsin.

<details><summary>Yechim</summary>

```bash
#!/usr/bin/env bash
[ $# -ge 1 ] || { echo "Usage: $0 <arg>..." >&2; exit 1; }
echo "soni: $#, birinchisi: $1, hammasi: $*"
```
</details>

**2.** `"$@"` va `"$*"` farqini isbotlovchi script yozing va `a "b c" d` bilan chaqiring.

<details><summary>Yechim</summary>

```bash
#!/usr/bin/env bash
echo '"$@" bilan:'; for x in "$@"; do echo " [$x]"; done
echo '"$*" bilan:'; for x in "$*"; do echo " [$x]"; done
```
`"$@"` — 3 element ([a], [b c], [d]); `"$*"` — 1 element.
</details>

**3.** `sum` scripti: berilgan barcha sonlarni qo'shsin (shift bilan), son bo'lmagan argumentda xato bersin.

<details><summary>Yechim</summary>

```bash
#!/usr/bin/env bash
total=0
while [ $# -gt 0 ]; do
    [[ "$1" =~ ^-?[0-9]+$ ]] || { echo "'$1' son emas" >&2; exit 1; }
    total=$((total + $1))
    shift
done
echo "$total"
```
</details>

**4.** `greet` scripti: `-n ISM` flagi bilan ism, `-l` flagi bilan lotincha "Salve!" rejimi; flaglar istalgan tartibda kelsin.

<details><summary>Yechim</summary>

```bash
#!/usr/bin/env bash
name="dunyo"; latin=0
while [ $# -gt 0 ]; do
    case "$1" in
        -n) shift; name="$1" ;;
        -l) latin=1 ;;
        *)  echo "nomalum: $1" >&2; exit 1 ;;
    esac
    shift
done
if [ "$latin" -eq 1 ]; then echo "Salve, $name!"; else echo "Salom, $name!"; fi
```
</details>

**5.** read + validatsiya: userdan 1-10 oralig'ida son so'rang; noto'g'ri kiritsa qayta so'rasin (3 urinishdan keyin taslim bo'lsin).

<details><summary>Yechim</summary>

```bash
#!/usr/bin/env bash
for attempt in 1 2 3; do
    read -rp "1-10 oralig'ida son > " n
    if [[ "$n" =~ ^[0-9]+$ ]] && ((n >= 1 && n <= 10)); then
        echo "Rahmat: $n"; exit 0
    fi
    echo "noto'g'ri, yana urinib ko'ring"
done
echo "3 urinish tugadi" >&2; exit 1
```
</details>

**6.** `/etc/passwd` ning istalgan qatorini `IFS=: read` bilan maydonlarga ajratib, "USER (uid=N) home=H shell=S" formatida chiqaring.

<details><summary>Yechim</summary>

```bash
IFS=: read -r user _ uid _ _ home shell <<< "$(grep '^root:' /etc/passwd)"
echo "$user (uid=$uid) home=$home shell=$shell"
# natija: root (uid=0) home=/root shell=/bin/bash
```
</details>

**7.** (Qiyinroq) `confirm-rm` scripti: argumentlardagi fayllarni o'chirishdan oldin har biri uchun `[y/N]` so'rasin; `--yes` flagi bilan so'ramasdan o'chirsin (CI rejimi).

<details><summary>Yechim</summary>

```bash
#!/usr/bin/env bash
yes=0
[ "${1:-}" = "--yes" ] && { yes=1; shift; }
for f in "$@"; do
    [ -e "$f" ] || { echo "yo'q: $f" >&2; continue; }
    if [ "$yes" -eq 1 ]; then
        rm "$f" && echo "o'chirildi: $f"
    else
        read -rp "o'chirilsinmi '$f'? [y/N] " a
        [ "$a" = "y" ] && rm "$f" && echo "o'chirildi: $f"
    fi
done
```
</details>

## Cheat sheet

| Element | Nima | Eslatma |
|---------|------|---------|
| `$0` `$1`..`$9` `${10}` | Script nomi va argumentlar | funksiyada ham ishlaydi |
| `$#` | Argumentlar soni | `[ $# -ge 1 ]` tekshiruv |
| `"$@"` | Har argument alohida | uzatishda HAR DOIM |
| `"$*"` | Bitta string | faqat log/xabar uchun |
| `shift` / `shift 2` | Argumentlarni surish | parsing sikllari |
| `${1:-def}` / `${1:?xato}` | Default / majburiy | 22-darsda to'liq |
| `while+case+shift` | Flag parsing qolipi | `-*)` va `*)` ni unutmang |
| `read -r var` | Stdin dan | `-r` majburiy |
| `read -rp "..." v` | Prompt bilan | interaktiv |
| `read -rs` | Yashirin (parol) | `ps` da ko'rinmaydi |
| `read -t 5` | Timeout | CI osilmasin |
| `IFS=: read -r a b` | Maydonlash | passwd/CSV |
| `<<<"str"` | Here string | bitta qator stdin |

## Qo'shimcha manbalar

- [Bash Reference — Special Parameters](https://www.gnu.org/software/bash/manual/html_node/Special-Parameters.html) — `$@`, `$*`, `$#` rasmiy semantikasi
- [BashFAQ/035 — How can I handle command-line options?](https://mywiki.wooledge.org/BashFAQ/035) — parsing usullari taqqoslamasi
- [Google Shell Style Guide — Flags](https://google.github.io/styleguide/shellguide.html) — sanoat uslubi

---

[← Oldingi: 20 — loops](20-loops.md) · [Kurs xaritasi](00-README.md) · [Keyingi: 22 — strings-numbers-arrays →](22-strings-numbers-arrays.md)
