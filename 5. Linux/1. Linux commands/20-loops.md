# 20. Looplar: while, until, for

> Manba: TLCL 29 va 33-boblar · Muhit: Ubuntu 24.04, bash 5.2 · [← Oldingi: branching](19-branching.md) · [Kurs xaritasi](00-README.md) · [Keyingi: script-input →](21-script-input.md)

## Nima uchun kerak

"Har bir servisga health-check yubor", "DB tayyor bo'lguncha kutib tur", "shu 50 ta faylning har biriga konvertatsiya" — takrorlanadigan har ish loop. Deploy scriptidagi retry mantiq, migration kutish, batch processing — bularsiz production script yozilmaydi. Bu darsda bash looplarining barcha shakllari va ikkita mashhur tuzoq (pipe-subshell va bo'sh glob) real isbotlari bilan.

## Nazariya

### Loop ham exit code ustida ishlaydi

`if` kabi (19-dars), `while`/`until` ham shartni emas — **buyruq exit codeini** tekshiradi:

```bash
while buyruq; do ... done      # buyruq 0 qaytarar ekan davom
until buyruq; do ... done      # buyruq 0 qaytarGUNCHA davom (teskari while)
```

Shuning uchun `while true`, `until pg_isready` kabi shakllar tabiiy.

`for` esa boshqacha — **ro'yxat ustida** yuradi: so'zlar, glob natijalari, `{1..N}`, `$(cmd)` outputi.

## Buyruqlar

### `while` va `until`

```console
$ count=1
$ while [ "$count" -le 5 ]; do echo -n "$count "; count=$((count + 1)); done
1 2 3 4 5
$ count=1
$ until [ "$count" -gt 5 ]; do echo -n "$count "; count=$((count + 1)); done
1 2 3 4 5
```

`until` — shart "teskari" tabiiy o'qiladigan joylarda: `until tayyor; do kut; done`.

### `break` va `continue`

```console
$ n=0
$ while true; do
>     n=$((n+1))
>     [ "$n" -eq 3 ] && continue     # 3 ni tashlab ket
>     [ "$n" -gt 5 ] && break        # 5 dan keyin chiq
>     echo -n "$n "
> done
1 2 4 5
```

`while true` + `break` — "o'rtasidan chiqiladigan" looplarning standart shakli.

### Fayl o'qish: `while read`

Qatorma-qator o'qishning **to'g'ri usuli** (satrni maydonlarga bo'lib beradi):

```console
$ while read -r distro version; do
>     echo "Distro: $distro, versiya: $version"
> done < distros.txt
Distro: SUSE, versiya: 10.2
Distro: Fedora, versiya: 10
Distro: Ubuntu, versiya: 8.04
```

`-r` — backslash larni literal qoldirish (har doim qo'ying). E'tibor: fayl **loopga redirect** qilindi (`done < fayl`).

**Pipe-subshell tuzog'i** — nega `cat fayl | while` yomon (jonli isbot):

```console
$ last=""
$ cat distros.txt | while read -r d v; do last="$d"; done
$ echo "pipe dan keyin last: [$last]"
pipe dan keyin last: []                    # O'ZGARISH YO'QOLDI!
$ while read -r d v; do last="$d"; done < distros.txt
$ echo "redirect bilan last: [$last]"
redirect bilan last: [Ubuntu]              # ishladi
```

Sabab: pipe ning har tomoni **alohida subshell** (05-dars — parallel processlar!) — loop ichidagi variable o'zgarishlari parent ga qaytmaydi. Yechim: `done < fayl` yoki `done < <(cmd)` (process substitution — 23-dars).

### `for` — klassik shakl

```console
$ for i in A B C D; do echo -n "$i "; done
A B C D
$ for i in {1..5}; do echo -n "$i "; done
1 2 3 4 5
$ for f in *.log; do echo "fayl: $f"; done
fayl: app.log
fayl: err.log
$ for u in $(cut -d: -f1 /etc/passwd | head -3); do echo -n "$u "; done
root daemon bin
```

Ro'yxat manbalari: literal so'zlar, brace expansion, **glob** (fayllar ustida eng xavfsiz usul — 02-darsdagi "ls ni parse qilmang" muammosining javobi!), command substitution (faqat probel/newline muammosi yo'q data uchun).

**Bo'sh glob tuzog'i va `nullglob`** (jonli isbot):

```console
$ for f in *.yoq; do echo "topildi: $f"; done
topildi: *.yoq                    # mos fayl yo'q — literal PATTERN keldi!
$ shopt -s nullglob
$ for f in *.yoq; do echo "topildi: $f"; done
                                  # endi loop umuman ishlamadi — to'g'ri
```

06-darsdagi qoida esingizdami: mos kelmagan glob o'z holicha qoladi. Fayllar ustidagi har loopda: `shopt -s nullglob` yoki ichida `[ -e "$f" ] || continue`.

### C-uslubdagi `for`

Hisoblagichli looplar uchun (Go dagi `for i := 0; i < 5; i++` bilan aynan):

```console
$ for ((i = 0; i < 5; i++)); do echo -n "$i "; done
0 1 2 3 4
```

06-darsdagi `{1..$N}` ishlamasligi muammosining yechimi ham shu: `for ((i=1; i<=N; i++))`.

## Real-world scenariylar

**1. DB tayyor bo'lishini kutish** (docker-compose, CI da klassika):

```bash
until pg_isready -h "$DB_HOST" -q; do
    echo "DB kutilmoqda..."
    sleep 2
done
echo "DB tayyor!"
```

**2. Retry with limit** (verify qilingan pattern):

```bash
attempt=1; max=4
until curl -sf "$URL/healthz" >/dev/null; do
    if [ "$attempt" -ge "$max" ]; then
        echo "barcha $max urinish yiqildi" >&2
        exit 1
    fi
    echo "urinish $attempt yiqildi, qayta..."
    attempt=$((attempt+1))
    sleep $((attempt * 2))        # ortib boruvchi kutish (backoff)
done
```

**3. Ko'p serverga buyruq:**

```bash
for host in app1 app2 app3; do
    echo "--- $host ---"
    ssh "$host" "df -h / | tail -1" || echo "XATO: $host javob bermadi"
done
```

**4. Batch qayta ishlash (safe glob bilan):**

```bash
shopt -s nullglob
for f in /var/log/app/*.log; do
    gzip "$f"
done
```

## Zamonaviy yondashuv

- **Loop o'rniga to'g'ri tool**: qatorlarni sanash/agregatsiya — `awk` (17-dars) loopdan yuz baravar tez; fayllar ustida amal — `find -exec` (14-dars). Bash loop — orkestratsiya uchun (buyruqlarni ketma-ket boshqarish), data processing uchun emas.
- **Parallellik**: ketma-ket `for host in ...; ssh` o'rniga — `xargs -P` (14-dars) yoki [GNU parallel](https://www.gnu.org/software/parallel/): `parallel -j4 gzip ::: *.log`.
- **`mapfile`** (bash 4+): fayl qatorlarini massivga bir amalda: `mapfile -t lines < fayl` — subshell tuzog'isiz (22-darsda massivlar bilan).
- **`seq` vs `{1..N}` vs C-for**: `{1..N}` faqat literal son; dinamik chegara — C-uslub `for ((...))` (bash) yoki `seq` (POSIX).
- ShellCheck bu darsning ham qo'riqchisi: SC2044 (`for f in $(find ...)` anti-patterni), SC2013 (`for l in $(cat ...)`).

## Keng tarqalgan xatolar

1. **`for line in $(cat fayl)` bilan qator o'qish.** Word splitting qatorlarni **so'zlarga** bo'lib yuboradi ("bir qator" ≠ "bir element"). Qatorlar uchun: `while read -r line; do ... done < fayl`.

2. **`cat fayl | while read` da variable yo'qolishi.** Yuqorida isbotlandi — pipe = subshell. `done < fayl` ishlating.

3. **`read` da `-r` ni tashlab ketish.** Backslashli data (Windows yo'llari, regex) buziladi. `read -r` — refleks bo'lsin.

4. **Bo'sh glob bilan loop.** `for f in *.tmp; do rm "$f"; done` — fayl bo'lmasa `rm '*.tmp'` bajariladi (xato, yaxshiyamki). `nullglob` yoki mavjudlik tekshiruvi.

5. **Loop ichida quotesiz `$f`.** Probelli fayl nomida buziladi (06-dars qaytadi va qaytadi): `gzip "$f"`.

6. **Cheksiz loopda `sleep` siz polling.** `while ! check; do :; done` — CPU 100%. Har pollingda `sleep`.

## Amaliy mashqlar

Muhit: `docker run -it --rm ubuntu:24.04 bash`

**1.** 1 dan 20 gacha sonlardan faqat 3 ga bo'linadiganlarini chiqaring — while bilan bir marta, C-uslub for bilan bir marta.

<details><summary>Yechim</summary>

```bash
i=1; while [ $i -le 20 ]; do ((i % 3 == 0)) && echo -n "$i "; i=$((i+1)); done; echo
for ((i=1; i<=20; i++)); do ((i % 3 == 0)) && echo -n "$i "; done; echo
```
</details>

**2.** `/etc/passwd` ni `while read` bilan o'qib, har user uchun "NOM (uid: N)" chiqaring (IFS maslahati: `IFS=:`).

<details><summary>Yechim</summary>

```bash
while IFS=: read -r name _ uid _; do
    echo "$name (uid: $uid)"
done < /etc/passwd | head -5
```
`IFS=:` faqat read uchun ajratuvchini almashtiradi; `_` — keraksiz maydonlarni yutish (21-darsda batafsil).
</details>

**3.** Pipe-subshell tuzog'ini o'zingiz isbotlang: fayl qatorlarini sanaydigan counter ni pipe bilan va redirect bilan yozib, farqni ko'rsating.

<details><summary>Yechim</summary>

```bash
cnt=0; cat /etc/passwd | while read -r _; do cnt=$((cnt+1)); done; echo "pipe: $cnt"
cnt=0; while read -r _; do cnt=$((cnt+1)); done < /etc/passwd; echo "redirect: $cnt"
```
Birinchisi 0 (subshellda qoldi), ikkinchisi haqiqiy son.
</details>

**4.** `backoff.sh`: mavjud bo'lmagan URL ga curl bilan 3 marta urinib (2s, 4s, 6s kutishlar bilan), oxirida xato bilan chiqsin.

<details><summary>Yechim</summary>

```bash
#!/usr/bin/env bash
for attempt in 1 2 3; do
    if curl -sf --max-time 3 "http://yoq.example.local/health"; then
        echo "OK"; exit 0
    fi
    echo "urinish $attempt yiqildi"
    sleep $((attempt * 2))
done
echo "3 urinish ham yiqildi" >&2; exit 1
```
</details>

**5.** Joriy katalogdagi har `.log` fayl uchun `nomi-YYYY-MM-DD.log` nusxa yarating; fayl yo'q bo'lsa loop umuman ishlamasin.

<details><summary>Yechim</summary>

```bash
shopt -s nullglob
today=$(date +%F)
for f in *.log; do
    cp "$f" "${f%.log}-$today.log"
done
```
`${f%.log}` — suffix kesish (22-darsda parameter expansion to'liq).
</details>

**6.** Ichma-ich loop: 3 "server" × 3 "servis" jadval chiqaring (`srv1/api OK` ko'rinishida).

<details><summary>Yechim</summary>

```bash
for host in srv1 srv2 srv3; do
    for svc in api worker cache; do
        echo "$host/$svc OK"
    done
done
```
</details>

**7.** (Qiyinroq) "Log kuzatuvchi": `until` bilan fayl paydo bo'lishini kutib (har soniyada tekshirib, max 10s), paydo bo'lgach oxirgi 3 qatorini chiqaring. Boshqa terminaldan faylni yaratib test qiling.

<details><summary>Yechim</summary>

```bash
#!/usr/bin/env bash
f="/tmp/kutilayotgan.log"
waited=0
until [ -f "$f" ]; do
    [ "$waited" -ge 10 ] && { echo "10s da paydo bo'lmadi" >&2; exit 1; }
    sleep 1; waited=$((waited+1))
done
tail -3 "$f"
```
Test: `(sleep 3; printf "a\nb\nc\nd\n" > /tmp/kutilayotgan.log) & ./watcher.sh`
</details>

## Cheat sheet

| Konstruksiya | Nima | Misol |
|--------------|------|-------|
| `while cmd; do...done` | 0 qaytarar ekan | `while [ $i -le 5 ]` |
| `until cmd; do...done` | 0 qaytarguncha | `until pg_isready` |
| `for x in ro'yxat` | Ro'yxat ustida | `for f in *.log` |
| `for ((i=0;i<n;i++))` | Hisoblagich | C/Go uslubi |
| `break` / `continue` | Chiqish / o'tkazish | `while true` + break |
| `while read -r ...` | Qatorma-qator | `done < fayl` (pipe EMAS!) |
| `IFS=: read -r a b` | Maydonlarga bo'lish | passwd kabi fayllar |
| `shopt -s nullglob` | Bo'sh glob = bo'sh ro'yxat | fayl looplarida |
| `mapfile -t arr` | Fayl → massiv | subshellsiz |
| Retry qolipi | `until cmd; do [ limit ] && break; sleep N; done` | deploy/healthcheck |

## Qo'shimcha manbalar

- [Bash Reference — Looping Constructs](https://www.gnu.org/software/bash/manual/html_node/Looping-Constructs.html) — rasmiy hujjat
- [BashFAQ/001: How can I read a file line-by-line?](https://mywiki.wooledge.org/BashFAQ/001) — while read ning barcha nozikliklari
- [GNU Parallel tutorial](https://www.gnu.org/software/parallel/parallel_tutorial.html) — looplarni parallellashtirish

---

[← Oldingi: 19 — branching](19-branching.md) · [Kurs xaritasi](00-README.md) · [Keyingi: 21 — script-input →](21-script-input.md)
