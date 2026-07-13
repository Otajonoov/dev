# 23. Advanced scripting

> Manba: TLCL 36-bob ("Exotica") · Muhit: Ubuntu 24.04, bash 5.2 · [← Oldingi: strings-numbers-arrays](22-strings-numbers-arrays.md) · [Kurs xaritasi](00-README.md) · [Keyingi: debugging-and-best-practices →](24-debugging-and-best-practices.md)

## Nima uchun kerak

Kitob bu bobni "Exotica" deb atagan, lekin backend uchun bular ekzotika emas — production zaruratlari: **trap** — script uzilganda temp fayllar/lock lar tozalanishi (Go dagi `defer` ning bash ekvivalenti!); **wait** — 5 serverga parallel deploy; **process substitution** — 20-darsdagi subshell muammosining professional yechimi. Bu dars scriptlaringizni "ishlaydi"dan "production-grade"ga ko'taradi.

## Nazariya

### Qayerda subshell paydo bo'ladi?

Subshell — scriptning child-process nusxasi: undagi o'zgarishlar (variable, `cd`) parentga **qaytmaydi**. Subshell paydo bo'ladigan joylar: `( ... )`, pipeline har qismi (20-darsdagi tuzoq!), `$(...)`, `&` bilan fon.

`{ ...; }` (group command) esa **joriy shellda** ishlaydi — gruppalash kerak-u subshell kerak bo'lmasa. Farqi jonli isbotlangan:

```console
$ x=tashqarida
$ ( x=ichkarida; echo "subshellda: $x" )
subshellda: ichkarida
$ echo "subshelldan keyin: $x"
subshelldan keyin: tashqarida        # o'zgarish yo'qoldi — subshell edi
$ ( cd /tmp && echo "ichkarida pwd: $(pwd)" )
ichkarida pwd: /tmp
$ pwd
/root/adv                            # parent joyida qoldi
```

Subshellning **foydali** tomoni ham shu: `( cd dir && tar xf ... )` — katalogni vaqtincha almashtirish, qaytishni o'ylamasdan (15-darsdagi tar-pipe misolida ishlatganmiz).

## Buyruqlar

### Group command: `{ }`

```console
$ { echo qator1; echo qator2; } > birga.txt
```

Bir nechta buyruq outputini **bitta** redirect ga yig'ish (uchtasini alohida `>>` bilan yozishdan toza). Sintaksis talabi: `{` dan keyin probel, oxirgi buyruqdan keyin `;`.

### Process substitution: `<(cmd)` va `>(cmd)`

Buyruq outputini **fayl sifatida** taqdim etadi:

```console
$ diff <(ls /usr/local/bin | head -3) <(ls /usr/sbin | head -3)
0a1,3
> accessdb
> add-shell
> addgroup
```

Ikkala `ls` natijasi "fayl"dek diff ga berildi — temp fayllarsiz! Va 20-darsdagi katta muammoning yechimi:

```console
$ while read -r line; do echo "PS orqali: $line"; done < <(echo "salom dunyo")
PS orqali: salom dunyo
```

`cmd | while` (subshell — o'zgarishlar yo'qoladi) o'rniga `while ... done < <(cmd)` — loop **joriy shellda**, variablelar saqlanadi. `mapfile -t arr < <(cmd)` ham shu oiladan (22-dars).

### `trap` — signal ushlash (bash "defer"i)

08-darsdagi signallar endi script ichida. Sintaksis: `trap 'buyruqlar' SIGNAL...`

**Cleanup pattern** — har production scriptning majburiy qismi (tekshirilgan):

```bash
#!/usr/bin/env bash
set -euo pipefail
tmpfile=$(mktemp)
cleanup() { rm -f "$tmpfile"; echo "cleanup: $tmpfile o'chirildi"; }
trap cleanup EXIT

echo "ishlayapman, tmp: $tmpfile"
false                                # kutilmagan xato!
echo "bu qator hech qachon chiqmaydi"
```

```console
$ ./trap-demo.sh
ishlayapman, tmp: /tmp/tmp.v4eGgqlASR
cleanup: /tmp/tmp.v4eGgqlASR o'chirildi     # yiqilishiga QARAMAY ishladi
```

`EXIT` — psevdo-signal: script **qanday tugashidan qat'i nazar** (normal, xato, Ctrl+C) ishlaydi. Aynan Go dagi `defer` semantikasi.

**Graceful shutdown** (tekshirilgan — TERM yuborilib):

```bash
trap 'echo "SIGTERM keldi — graceful shutdown"; exit 0' TERM
echo "PID: $$, kutyapman..."
sleep 30 & wait $!
```

```console
$ ./trap-sig.sh &  kill -TERM $!
PID: 7074, kutyapman...
SIGTERM keldi — graceful shutdown
```

Nozik nuqta: bash signal ni **joriy buyruq tugagach** qayta ishlaydi — shuning uchun `sleep 30 & wait` (wait signalda darhol uziladi), `sleep 30` emas. Docker konteynerdagi entrypoint scriptlarda (08-darsdagi PID 1 muammosi) shu pattern ishlatiladi.

### `wait` va parallel bajarish

`&` bilan fonga yuborilganlarni kutish (tekshirilgan — 3× tezlashuv):

```console
$ sleep 0.5 & sleep 0.5 & sleep 0.5 &
$ wait
$ # 3 parallel sleep: 506 ms (ketma-ket 1500 ms bo'lardi)
```

`wait PID` — bittasini kutib **exit codeini olish**:

```console
$ ( exit 3 ) & p=$!
$ wait $p; echo "child exit: $?"
child exit: 3
```

`$!` — oxirgi fon processning PID i. Pattern: PID larni massivga yig'ib, keyin har birini tekshirish (pastdagi scenariyda).

### Named pipes (FIFO)

Ikki **alohida** process orasida quvur (05-darsdagi `|` — bir buyruq qatori ichida edi; FIFO — fayl tizimida turadigan, istalgan ikki processni ulaydigan quvur):

```console
$ mkfifo mypipe
$ ls -l mypipe
prw-r--r-- mypipe                    # tur belgisi: p (pipe!)
$ echo "quvur orqali xabar" > mypipe &     # yozuvchi kutib turadi...
$ read -r msg < mypipe && echo "o'qildi: $msg"
o'qildi: quvur orqali xabar
```

Data diskka yozilmaydi — to'g'ridan-to'g'ri processdan processga. Ishlatiladigan joylar: katta dumplarni diskka tushirmasdan uzatish, servislarning oddiy IPC si.

### `mktemp` — xavfsiz temp fayllar

```console
$ t=$(mktemp) && echo "$t"
/tmp/tmp.J4Wm9XY3xk
$ d=$(mktemp -d)                     # katalog varianti
/tmp/tmp.r30NmXohEf
```

Qo'lda `/tmp/myscript.tmp` yozish — race condition va xavfsizlik teshigi (boshqa user oldindan symlink qo'yishi mumkin). Har doim mktemp + trap cleanup jufti.

## Real-world scenariylar

**1. Production script to'liq skeleti** (18-dars skeletiga trap qo'shildi):

```bash
#!/usr/bin/env bash
set -euo pipefail

workdir=$(mktemp -d)
cleanup() {
    rm -rf "$workdir"
    # lock fayl, fon processlar ham shu yerda tozalanadi
}
trap cleanup EXIT

main() {
    cd "$workdir"
    # ... asosiy ish ...
}
main "$@"
```

**2. Parallel deploy — xatolarni yo'qotmasdan:**

```bash
declare -A pids
for host in web1 web2 web3; do
    ssh "$host" "systemctl restart myapp" &
    pids[$host]=$!
done

failed=0
for host in "${!pids[@]}"; do
    if ! wait "${pids[$host]}"; then
        echo "XATO: $host yiqildi" >&2
        failed=1
    fi
done
exit "$failed"
```

**3. Ikki manba farqini tekshirish (process substitution):**

```bash
diff <(ssh prod1 "cat /etc/app.conf") <(ssh prod2 "cat /etc/app.conf") \
    && echo "configlar bir xil" || echo "FARQ BOR!"
diff <(sort local-users.txt) <(sort remote-users.txt)
```

## Zamonaviy yondashuv

- **`trap cleanup EXIT`** — zamonaviy bash style guide larning majburiy bandi; INT/TERM ni alohida ushlab EXIT ga yo'naltirish ham keng: `trap 'exit 130' INT` (130 = 128+2, signal konventsiyasi).
- **Parallellik darajalari**: bir nechta ish — `& + wait`; o'nlab bir xil ish — `xargs -P` (14-dars); yuzlab, progress bilan — GNU `parallel`. Bash job control bilan murakkab pool yasashga urinmang — bu chegara.
- **Named pipe o'rniga** zamonaviy stacklarda message queue (Redis, NATS) ishlatiladi; FIFO — bir mashina ichidagi yengil holatlar uchun qoladi.
- **`coproc`** — bash 4+ dagi ikki tomonlama pipe li fon process; kam ishlatiladi, bilib qo'yish kifoya.
- Process substitution `/dev/fd` mexanizmiga tayanadi — `#!/bin/sh` (dash) da **ishlamaydi**: bash-only scriptlarda ishlating (shebang to'g'ri bo'lsin).

## Keng tarqalgan xatolar

1. **Pipeline subshellida variable o'zgartirish** (20-darsdan qaytadi): `cmd | while read` — yo'qoladi. Endi to'liq arsenal bor: `< <(cmd)` yoki `mapfile`.

2. **trap siz temp fayllar.** Script yiqilsa `/tmp` axlatxonaga aylanadi, lock fayllar qolib keyingi ishga tushishni bloklaydi. mktemp yaratdingizmi — shu zahoti trap yozing.

3. **`wait` ni PID siz chaqirib exit codelarni yo'qotish.** Yalang `wait` — hammasini kutadi, lekin **xatolarni yutadi** (oxirgi kutilganniki qaytadi xolos). Har child muhim bo'lsa — PID larni saqlab alohida `wait $pid`.

4. **trap ichida murakkab logika.** Signal kontekstida uzun ish — yangi muammolar (reentrancy). Trap faqat funksiya chaqirsin: `trap cleanup EXIT` — logika funksiyada.

5. **FIFO ga yozuvchisiz o'quvchi ochish (yoki aksincha) — deadlock.** FIFO ikkala tomoni ochilguncha bloklaydi. Bir tomonni `&` fonda oching (yuqoridagi misoldagidek).

6. **`(cd dir; cmd)` deb yozib subshell ekanini unutish.** cd yiqilsa cmd **joriy** katalogda ishlaydi! To'g'ri: `(cd dir && cmd)` — 19-darsdagi qoidaning subshell varianti.

## Amaliy mashqlar

Muhit: `docker run -it --rm ubuntu:24.04 bash`

**1.** Group command bilan uchta buyruq (`date`, `uname -a`, `uptime`) outputini bitta `report.txt` ga yozing. Keyin xuddi shuni subshell bilan qiling va farqni ayting.

<details><summary>Yechim</summary>

```bash
{ date; uname -a; uptime; } > report.txt
( date; uname -a; uptime ) > report2.txt
```
Natija bir xil; farq — `{}` joriy shellda (tezroq), `()` subshellda. Redirect uchun `{}` yetadi.
</details>

**2.** Subshell izolyatsiyasini isbotlang: funksiyada `( cd /tmp; touch f )` bajarib, parent pwd o'zgarmasligini ko'rsating.

<details><summary>Yechim</summary>

```console
$ pwd; ( cd /tmp && touch subshell-isbotf ); pwd
/root
/root                    # joyimizda qoldik
$ ls /tmp/subshell-isbotf
/tmp/subshell-isbotf     # ish esa bajarilgan
```
</details>

**3.** Ikki katalog tarkibini temp faylsiz solishtiring: `/usr/bin` va `/usr/sbin` da bir xil nomli buyruqlar bormi?

<details><summary>Yechim</summary>

```bash
comm -12 <(ls /usr/bin | sort) <(ls /usr/sbin | sort) | head
```
Process substitution + comm (17-dars) — klassik juftlik.
</details>

**4.** `safe-work.sh`: mktemp bilan ish katalogi yaratib, trap EXIT bilan tozalasin. Scriptni o'rtasida `false` bilan yiqitib, katalog baribir o'chirilganini isbotlang.

<details><summary>Yechim</summary>

```bash
#!/usr/bin/env bash
set -euo pipefail
wd=$(mktemp -d)
trap 'rm -rf "$wd"; echo "tozalandi: $wd"' EXIT
echo "ishchi katalog: $wd"
touch "$wd/data"
false
```
Ishga tushirib: "tozalandi" chiqadi va `ls $wd` — yo'q.
</details>

**5.** SIGINT (Ctrl+C) ni ushlaydigan script: bosilganda "yakunlanmoqda..." deb 130 kodi bilan chiqsin. `kill -INT` bilan test qiling.

<details><summary>Yechim</summary>

```bash
#!/usr/bin/env bash
trap 'echo "yakunlanmoqda..."; exit 130' INT
echo "PID $$ — Ctrl+C ni kutyapman"
sleep 60 & wait $!
```
Test: fonda ishga tushirib `kill -INT <pid>`.
</details>

**6.** Parallel ping: 4 ta hostga (8.8.8.8, 1.1.1.1, yoq.local, 9.9.9.9) parallel `ping -c1` yuborib, har birining natijasini (OK/FAIL) alohida hisobot qiling.

<details><summary>Yechim</summary>

```bash
declare -A pids
for h in 8.8.8.8 1.1.1.1 yoq.local 9.9.9.9; do
    ping -c1 -W2 "$h" >/dev/null 2>&1 &
    pids[$h]=$!
done
for h in "${!pids[@]}"; do
    wait "${pids[$h]}" && echo "$h: OK" || echo "$h: FAIL"
done
```
</details>

**7.** (Qiyinroq) FIFO orqali "log kollektor": bitta terminalda FIFO dan o'qib har qatorga timestamp qo'shadigan o'quvchi; boshqa processlar FIFO ga yozadi.

<details><summary>Yechim</summary>

```bash
mkfifo /tmp/logpipe
# o'quvchi (fonda):
while read -r line; do echo "[$(date +%T)] $line"; done < /tmp/logpipe &
# yozuvchilar:
echo "birinchi xabar" > /tmp/logpipe
echo "ikkinchi xabar" > /tmp/logpipe
sleep 0.5; rm /tmp/logpipe
```
Har yozuvchi alohida process — FIFO ularni bitta oqimga yig'di. syslog ning soddalashtirilgan modeli.
</details>

## Cheat sheet

| Konstruksiya | Nima | Qachon |
|--------------|------|--------|
| `{ cmd1; cmd2; }` | Joriy shellda guruh | umumiy redirect |
| `( cmd1; cmd2 )` | Subshell | izolyatsiya (cd, var) |
| `<(cmd)` | Output → "fayl" | `diff <(a) <(b)`, `done < <(cmd)` |
| `trap fn EXIT` | Har qanday chiqishda | cleanup (defer!) |
| `trap 'exit 130' INT TERM` | Signal ushlash | graceful shutdown |
| `mktemp` / `mktemp -d` | Xavfsiz temp | + trap juftlikda |
| `cmd &` + `$!` | Fon + PID | parallellik |
| `wait` / `wait $pid` | Kutish / exit code olish | har child alohida |
| `mkfifo` | Named pipe | processlararo quvur |
| `$$` | Joriy shell PID | lock/log fayllarga |

## Qo'shimcha manbalar

- [Bash Reference — Compound Commands](https://www.gnu.org/software/bash/manual/html_node/Compound-Commands.html) — subshell/group rasmiy semantikasi
- [Greg's Wiki — SignalTrap](https://mywiki.wooledge.org/SignalTrap) — trap ning barcha nozikliklari
- [Greg's Wiki — ProcessManagement](https://mywiki.wooledge.org/ProcessManagement) — fon processlar va wait patternlari

---

[← Oldingi: 22 — strings-numbers-arrays](22-strings-numbers-arrays.md) · [Kurs xaritasi](00-README.md) · [Keyingi: 24 — debugging-and-best-practices →](24-debugging-and-best-practices.md)
