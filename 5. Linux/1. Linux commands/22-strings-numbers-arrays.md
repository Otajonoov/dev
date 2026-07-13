# 22. Stringlar, sonlar va massivlar

> Manba: TLCL 34 va 35-boblar · Muhit: Ubuntu 24.04, bash 5.2 · [← Oldingi: script-input](21-script-input.md) · [Kurs xaritasi](00-README.md) · [Keyingi: advanced-scripting →](23-advanced-scripting.md)

## Nima uchun kerak

`basename`, `dirname`, `sed`, `cut` uchun har safar alohida process ochmasdan — bash o'zi string bilan ishlay oladi: fayl nomidan kengaytmani kesish, default qiymatlar, almashtirish — hammasi **parameter expansion** ichida. Massivlar esa "server ro'yxati ustida ishlash" kabi vazifalarni stringga qo'shib-ajratish azobisiz hal qiladi (va 06-darsdagi SC2086 muammosining to'g'ri javobi ham shu). Bu dars — bash ning "standart kutubxonasi".

## Nazariya

### `${...}` — shunchaki qavs emas

06-darsda `$VAR` ni ko'rgan edik. To'liq shakl `${VAR}` ichida esa butun mini-til yashaydi: default qiymatlar, kesish, almashtirish, registr. Bularning hammasi **bitta processda** (subprocess yo'q — sed/awk chaqirishdan tezroq) ishlaydi.

Pattern kesishda yodlash qoidasi (klaviaturadagi joylashuv bo'yicha):
- `#` — **chapdan** (boshidan) kesadi ($ belgisidan chapda turadi)
- `%` — **o'ngdan** (oxiridan) kesadi
- bittasi — eng qisqa moslik, ikkitasi (`##`, `%%`) — eng uzun moslik

## Buyruqlar

### Default qiymatlar oilasi (hammasi tekshirilgan)

| Sintaksis | Ma'nosi |
|-----------|---------|
| `${var:-def}` | var bo'sh bo'lsa def **qaytar** (var o'zgarmaydi) |
| `${var:=def}` | var bo'sh bo'lsa def **o'rnat va qaytar** |
| `${var:?xabar}` | var bo'sh bo'lsa xabar bilan **xato + exit** |
| `${var:+qiymat}` | var **mavjud bo'lsa** qiymat qaytar |

```console
$ unset foo
$ echo "[${foo:-standart}]"
[standart]
$ echo "[${foo:=ornatildi}]" && echo "endi foo: $foo"
[ornatildi]
endi foo: ornatildi
$ echo "${bosh:?bu variable majburiy}"
bash: bosh: bu variable majburiy       # exit=1 — script to'xtaydi
```

`:?` — 03-darsda va'da qilingan himoya: `rm -rf "${APP_DIR:?APP_DIR o'rnatilmagan}"` — endi bo'sh variable falokat keltirmaydi.

### Uzunlik, substring

```console
$ s="production-server-01"
$ echo "${#s}"          # uzunlik
20
$ echo "${s:0:10}"      # offset:length
production
$ echo "${s: -2}"       # oxiridan 2 (probel MAJBURIY: :- bilan adashmasin)
01
```

### Prefix/suffix kesish — eng ko'p ishlatiladigani

```console
$ f="/srv/app/releases/app-v2.4.1.tar.gz"
$ echo "${f##*/}"       # eng uzun */ ni boshidan kesish = basename
app-v2.4.1.tar.gz
$ echo "${f%/*}"        # oxirgi /* ni kesish = dirname
/srv/app/releases
$ echo "${f%.tar.gz}"   # suffixni kesish
/srv/app/releases/app-v2.4.1
$ v="app-v2.4.1.tar.gz"
$ echo "${v#*-}"        # boshidan birinchi '-' gacha
v2.4.1.tar.gz
$ echo "${v%.*}"        # oxiridan bitta kengaytma
app-v2.4.1.tar
```

20-darsdagi `${f%.log}` sirlari ochildi. `basename`/`dirname` buyruqlaridan tezroq (subprocess yo'q) va loop ichida sezilarli.

### Almashtirish va registr

```console
$ path="/usr/local/bin:/usr/bin:/usr/local/sbin"
$ echo "${path/local/LOCAL}"      # birinchi moslik
/usr/LOCAL/bin:/usr/bin:/usr/local/sbin
$ echo "${path//local/LOCAL}"     # HAMMA moslik (//)
/usr/LOCAL/bin:/usr/bin:/usr/LOCAL/sbin
$ name="alisher"
$ echo "${name^} ${name^^} ${name,,}"    # Birinchi / HAMMA katta / hamma kichik
Alisher ALISHER alisher
```

### Arifmetika chuqurroq

06-darsdagi `$(( ))` ning to'liq imkoniyatlari (tekshirilgan):

```console
$ echo $((0xff)) $((2#11111111)) $((010))     # hex, binary, octal!
255 255 8
$ echo "daraja: $((2**10)), qoldiq: $((17 % 5))"
daraja: 1024, qoldiq: 2
$ a=5; ((a += 3)); ((a++)); echo "$a"
9
$ echo $(( a > 5 ? 1 : 0 ))                    # ternary
1
```

Diqqat: `$((010))` = **8** — bosh nol octal degani! Sana/vaqt bilan ishlashda (`08`, `09` soatlar) klassik bug manbai. Davosi: `$((10#$soat))` — majburiy o'nlik.

Kasrli hisob — bash da yo'q; `bc` (arbitrary precision kalkulyator):

```console
$ echo "scale=4; 10/3" | bc
3.3333
$ echo "scale=2; (1500*0.12)/12" | bc     # oylik foiz misoli
15.00
```

(Yoki `awk 'BEGIN{printf "%.4f\n", 10/3}'` — bc o'rnatilmagan joyda.)

### Massivlar

```console
$ days=(dush sesh chor pay jum shan yak)
$ echo "birinchi: ${days[0]}, hammasi: ${days[@]}"
birinchi: dush, hammasi: dush sesh chor pay jum shan yak
$ echo "soni: ${#days[@]}, indekslar: ${!days[@]}"
soni: 7, indekslar: 0 1 2 3 4 5 6
$ days+=(bayram)                    # qo'shish
$ echo "${days[@]:1:3}"             # kesim (slice)
sesh chor pay
$ unset "days[7]"                   # element o'chirish
```

**Eng muhim qoida** — loop va uzatishda `"${arr[@]}"` (quote bilan!): har element probeli bo'lsa ham butun qoladi (tekshirilgan):

```console
$ files=("app config.yaml" "main.go" "run.sh")
$ for f in "${files[@]}"; do echo "[$f]"; done
[app config.yaml]
[main.go]
[run.sh]
```

06-darsdagi ShellCheck maslahati amalda: **bir nechta argumentni variable da saqlash kerakmi — string emas, massiv**: `opts=(-v --color=auto); ls "${opts[@]}"`.

### Assotsiativ massivlar (bash 4+)

Kalit-qiymat (Go dagi map):

```console
$ declare -A ports
$ ports[api]=8080; ports[db]=5432; ports[cache]=6379
$ for svc in "${!ports[@]}"; do echo "$svc -> ${ports[$svc]}"; done
db -> 5432
api -> 8080
cache -> 6379
```

(`declare -A` majburiy; iteratsiya tartibi kafolatlanmagan — xuddi Go map kabi!)

### `mapfile` — fayl/buyruqdan massivga

20-darsdagi subshell muammosining eng toza yechimi:

```console
$ mapfile -t users < <(cut -d: -f1 /etc/passwd | head -3)
$ echo "massivda ${#users[@]} ta: ${users[@]}"
massivda 3 ta: root daemon bin
```

(`< <(...)` — process substitution, 23-darsda; `-t` — newline larni olib tashlash.)

## Real-world scenariylar

**1. Deploy artifakt nomini parchalash:**

```bash
artifact="app-v2.4.1-linux-amd64.tar.gz"
base="${artifact%.tar.gz}"          # app-v2.4.1-linux-amd64
version="${base#app-}"; version="${version%%-*}"   # v2.4.1
echo "versiya: $version"
```

**2. Batch rename — kengaytma almashtirish:**

```bash
for f in *.jpeg; do
    mv "$f" "${f%.jpeg}.jpg"
done
```

**3. Server guruhlari (assotsiativ massiv bilan):**

```bash
declare -A group=( [web]="web1 web2" [db]="db1" [cache]="redis1 redis2" )
for host in ${group[$1]:?guruh topilmadi}; do
    ssh "$host" uptime
done
```

## Zamonaviy yondashuv

- **Parameter expansion vs sed/awk**: bitta variable ustida — expansion (tez, subprocess yo'q); oqim/fayl ustida — sed/awk. `basename`/`dirname` o'rniga `##*/` va `%/*` — hot loop larda odat qiling.
- **Massivlar bash 3 da cheklangan (macOS default!)** — assotsiativ massivlar bash 4+. Portativlik kerak bo'lsa tekshiring: `((BASH_VERSINFO[0] >= 4))`.
- **`printf -v var ...`** — natijani echo qilmasdan variablega formatlash: `printf -v ts '%(%F_%T)T' -1` (joriy vaqt!).
- Data strukturalari murakkablashsa (nested, JSON) — bu bash chegarasi: `jq` yoki Go/Python. Assotsiativ massiv ichida massiv — bash da yo'q, va emulyatsiya qilishga urinish — kod hidining o'zi.

## Keng tarqalgan xatolar

1. **`$arr` deb butun massivni olmoqchi bo'lish.** `$arr` == `${arr[0]}` — faqat birinchi element! Butun massiv: `"${arr[@]}"`.

2. **`${arr[@]}` ni quotesiz ishlatish.** Probelli elementlar bo'linadi — massivning butun ma'nosi yo'qoladi. Har doim `"${arr[@]}"`.

3. **`#` va `%` ni adashtirish.** `${f#...}` boshidan, `${f%...}` oxiridan. Eslash: `#` — komment **boshida** turadi; `%` — foiz sonning **oxirida**.

4. **Bosh nolli sonlar octal bo'lishi.** `$((09))` — "value too great for base" xatosi! Sana parsing da: `$((10#$oy))`.

5. **`${s:-2}` va `${s: -2}` farqi.** Probelsiz `:-` — default qiymat operatori; oxirdan substring uchun probel shart: `${s: -2}` (yoki `${s:(-2)}`).

6. **declare -A siz assotsiativ massiv.** `ports[api]=8080` declare siz — bash indeksni arifmetik kontekstda hisoblab (`api`=0), oddiy massivning 0-elementiga yozadi. Kalitlar jimgina ustma-ust tushadi. `declare -A` — majburiy.

## Amaliy mashqlar

Muhit: `docker run -it --rm ubuntu:24.04 bash`

**1.** `url="https://api.example.com:8443/v2/users"` dan protokol, host:port va path ni faqat parameter expansion bilan ajrating.

<details><summary>Yechim</summary>

```bash
url="https://api.example.com:8443/v2/users"
proto="${url%%://*}"                 # https
rest="${url#*://}"                   # api.example.com:8443/v2/users
hostport="${rest%%/*}"               # api.example.com:8443
path="/${rest#*/}"                   # /v2/users
echo "$proto | $hostport | $path"
```
</details>

**2.** Default qiymatlar bilan config o'qish: `PORT`, `HOST`, `LOG_LEVEL` env variablelarini defaultlar (8080, 0.0.0.0, info) bilan chiqaring; `DATABASE_URL` esa majburiy bo'lsin.

<details><summary>Yechim</summary>

```bash
port="${PORT:-8080}"
host="${HOST:-0.0.0.0}"
log_level="${LOG_LEVEL:-info}"
db="${DATABASE_URL:?DATABASE_URL majburiy}"
echo "$host:$port ($log_level)"
```
</details>

**3.** Fayl nomlari massivi yasab (`logs=(app.log err.log "access log.txt")`), har birining kengaytmasiz nomini chiqaring — probelli nom buzilmasin.

<details><summary>Yechim</summary>

```bash
logs=(app.log err.log "access log.txt")
for f in "${logs[@]}"; do
    echo "${f%.*}"
done
# app / err / access log
```
</details>

**4.** `${#}` operatorlari bilan: `/etc/passwd` dagi eng uzun user nomini toping (faqat bash, awk siz).

<details><summary>Yechim</summary>

```bash
longest=""
while IFS=: read -r name _; do
    [ "${#name}" -gt "${#longest}" ] && longest="$name"
done < /etc/passwd
echo "$longest (${#longest} belgi)"
```
</details>

**5.** Octal tuzoqni ko'rsating: `soat="09"` bilan `$((soat + 1))` nima beradi? To'g'ri varianti?

<details><summary>Yechim</summary>

```console
$ soat="09"; echo $((soat + 1))
bash: 09: value too great for base (error token is "09")
$ echo $((10#$soat + 1))
10
```
</details>

**6.** Assotsiativ massiv bilan mini-DNS: 3 ta host→IP jufti saqlang; argumentda kelgan host uchun IP qaytaring, topilmasa xato.

<details><summary>Yechim</summary>

```bash
#!/usr/bin/env bash
declare -A dns=( [web1]=10.0.0.11 [db1]=10.0.0.21 [cache1]=10.0.0.31 )
host="${1:?Usage: $0 <host>}"
ip="${dns[$host]:-}"
[ -n "$ip" ] && echo "$ip" || { echo "topilmadi: $host" >&2; exit 1; }
```
</details>

**7.** (Qiyinroq) "Top-5 katta fayl" hisobotini massivlar bilan: `du /usr/share/* | sort -nr | head -5` natijasini mapfile ga oling va tartib raqami bilan formatlab chiqaring.

<details><summary>Yechim</summary>

```bash
mapfile -t top < <(du -s /usr/share/* 2>/dev/null | sort -nr | head -5)
for i in "${!top[@]}"; do
    printf "%d) %s\n" "$((i+1))" "${top[$i]}"
done
```
`"${!top[@]}"` — indekslar ro'yxati; printf bilan raqamlash.
</details>

## Cheat sheet

| Sintaksis | Nima | Misol |
|-----------|------|-------|
| `${v:-d}` / `${v:=d}` | Default (o'qish / o'rnatish) | `"${PORT:-8080}"` |
| `${v:?msg}` | Majburiy variable | `"${DIR:?kerak}"` |
| `${#v}` | Uzunlik | — |
| `${v:o:l}` | Substring | `${s:0:10}`, `${s: -2}` |
| `${v#p}` / `${v##p}` | Boshidan kesish (qisqa/uzun) | `${f##*/}` = basename |
| `${v%p}` / `${v%%p}` | Oxiridan kesish | `${f%.*}` — kengaytmasiz |
| `${v/a/b}` / `${v//a/b}` | Almashtirish (1 / hamma) | — |
| `${v^^}` / `${v,,}` | KATTA / kichik | — |
| `$((expr))` | Arifmetika | `**` daraja, `10#$n` baza |
| `bc` | Kasrli hisob | `echo "scale=2; 10/3" \| bc` |
| `arr=(a b c)` | Massiv | `"${arr[@]}"` — hamma (quote!) |
| `${#arr[@]}` / `${!arr[@]}` | Soni / indekslar | — |
| `arr+=(x)` | Qo'shish | — |
| `declare -A m` | Map | `${m[key]}`, `"${!m[@]}"` |
| `mapfile -t a < <(cmd)` | Buyruq → massiv | subshellsiz |

## Qo'shimcha manbalar

- [Bash Reference — Shell Parameter Expansion](https://www.gnu.org/software/bash/manual/html_node/Shell-Parameter-Expansion.html) — to'liq rasmiy ro'yxat
- [BashGuide — Arrays](https://mywiki.wooledge.org/BashGuide/Arrays) — massivlar bo'yicha eng yaxshi qo'llanma
- [Bash Hackers — Parameter Expansion cheat sheet](https://bash-hackers.gabe565.com/syntax/pe/) — tez spravka

---

[← Oldingi: 21 — script-input](21-script-input.md) · [Kurs xaritasi](00-README.md) · [Keyingi: 23 — advanced-scripting →](23-advanced-scripting.md)
