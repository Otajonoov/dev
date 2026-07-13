# 14. Fayllarni qidirish

> Manba: TLCL 17-bob · Muhit: Ubuntu 24.04 · [← Oldingi: networking](13-networking.md) · [Kurs xaritasi](00-README.md) · [Keyingi: archiving-and-sync →](15-archiving-and-sync.md)

## Nima uchun kerak

"7 kundan eski loglarni top va o'chir", "1GB dan katta fayllarni ko'rsat", "shu configni qaysi katalogga qo'yganman?" — bularning hammasi `find` ishi. `find` — shunchaki qidiruv emas, **fayl tizimi bo'ylab so'rovlar tili** (SQL ga o'xshatsa bo'ladi: WHERE shartlari + amallar). U bilan `xargs` juftligi shell dagi eng kuchli avtomatlashtirish kombinatsiyalaridan biri. Log tozalash cron joblari, CI dagi artifact yig'ishlar — hammasi shu darsdan.

## Nazariya

### Ikki xil qidiruv strategiyasi

- **`locate`** — oldindan qurilgan **indeks bazasidan** qidiradi: chaqmoq tez, lekin baza `updatedb` bilan davriy yangilanadi (odatda kuniga bir marta cron da) — yangi fayllarni ko'rmasligi mumkin.
- **`find`** — fayl tizimini **jonli aylanib chiqadi**: har doim aktual, atributlar bo'yicha murakkab shartlar, natijalar ustida amallar; kattaroq daraxtda sekinroq.

Qoida: "nomini bilaman, qayerdaligini bilmayman" → locate; qolgan hamma narsa → find.

### find ning uch tarkibiy qismi

```
find  [qayerdan]  [testlar]  [amallar]
find  /var/log    -name "*.log" -mtime +7    -delete
```

- **Testlar** — filtrlar (tur, nom, hajm, vaqt, egasi...)
- **Operatorlar** — testlarni bog'lash (-and default, -or, -not, guruhlash `\( \)`)
- **Amallar** — topilganlar bilan nima qilish (-print default, -delete, -exec)

## Buyruqlar

### `locate`

```console
$ locate bin/zip
/usr/bin/zipdetails
$ locate zip | grep bin      # grep bilan aniqlashtirish
```

Baza eskirgan bo'lsa: `sudo updatedb`. Zamonaviy distributivlarda tez varianti — `plocate`.

### `find` — testlar

Test poligonimiz (verify uchun yaratilgan struktura):

```
lab/
├── src/        main.go, util.go, main_test.go
├── logs/       app.log, error.log, big.log (5MB)
├── backup/     old.tar (10 kun eski)
└── empty.txt, empty-dir/
```

**Tur bo'yicha** (`-type`): `f` fayl, `d` katalog, `l` symlink:

```console
$ find ~/lab -type d
/root/lab
/root/lab/backup
/root/lab/logs
/root/lab/empty-dir
/root/lab/src
```

**Nom bo'yicha** (`-name`, katta-kichik farqsiz: `-iname`):

```console
$ find ~/lab -type f -name "*.go"
/root/lab/src/main_test.go
/root/lab/src/main.go
/root/lab/src/util.go
```

Pattern **majburiy quote ichida** — aks holda shell o'zi expand qilib yuboradi (06-dars).

**Hajm bo'yicha** (`-size`): `+` katta, `-` kichik, belgisiz — aynan. Birliklar: `c` bayt, `k`, `M`, `G`:

```console
$ find ~/lab -type f -size +1M
/root/lab/logs/big.log
```

**Vaqt bo'yicha**: `-mtime n` — kontent o'zgarganiga `n*24` soat (`+7` — 7 kundan eski); `-mmin` — daqiqalarda; `-newer fayl` — fayldan yangi:

```console
$ find ~/lab -type f -mtime +7
/root/lab/backup/old.tar
```

**Boshqa foydali testlar**: `-empty` (bo'sh fayl/katalog), `-user nom`, `-perm -o+w` (world-writable — xavfsizlik auditi!), `-maxdepth 2` (chuqurlikni cheklash — texnik jihatdan option, testlardan **oldin** yoziladi).

### Operatorlar

Default — AND (testlar ketma-ket = hammasi bajarilsin):

```console
$ find ~/lab -type f -name "*.log" -not -name "big*"
/root/lab/logs/error.log
/root/lab/logs/app.log
$ find ~/lab \( -name "*.go" -o -name "*.tar" \) | sort
/root/lab/backup/old.tar
/root/lab/src/main.go
/root/lab/src/main_test.go
/root/lab/src/util.go
```

`-o` — OR; qavslar shell dan himoya uchun `\( \)` shaklida. `!` = `-not`.

### Amallar

**`-exec`** — har topilgan fayl uchun buyruq; `{}` — fayl nomi o'rni, `\;` — yakun:

```console
$ find ~/lab -name "*.log" -exec ls -lh {} \;
... 0 /root/lab/logs/error.log
... 5.0M /root/lab/logs/big.log
```

**`\;` vs `+`** — muhim farq: `\;` har fayl uchun **alohida process** (1000 fayl = 1000 ta ls); `+` — fayllarni **bitta buyruqqa to'plab** beradi (xargs kabi, ancha tez):

```console
$ find ~/lab -name "*.go" -exec wc -l {} +
0 /root/lab/src/main_test.go
0 /root/lab/src/main.go
0 /root/lab/src/util.go
0 total
```

**`-ok`** — exec ning "so'rab ishlaydigan" varianti (har fayl uchun y/n). **`-delete`** — o'chirish:

```console
$ find ~/lab -name "*.tar" -print     # OLDIN nima o'chishini KO'RISH!
/root/lab/backup/old.tar
$ find ~/lab -name "*.tar" -delete
```

Temir qoida: `-delete` / `-exec rm` dan oldin xuddi shu shartlar bilan `-print` yoki `-ls` bilan dry-run.

### `xargs` — stdin dan buyruq yasash

`find` natijalarini boshqa buyruqqa argument qilib berish:

```console
$ find ~/lab -name "*.go" | xargs wc -l
...
0 total
```

**Probelli nomlar tuzog'i** — verify qilingan klassika:

```console
$ find ~/lab -name "*.txt" | xargs ls -l
ls: cannot access '/root/lab/xato': No such file or directory
ls: cannot access 'nom.txt': No such file or directory
```

`xato nom.txt` fayli probelda ikkiga bo'linib ketdi! Yechim — **null-separator juftligi**:

```console
$ find ~/lab -name "*.txt" -print0 | xargs -0 ls -l
... /root/lab/empty.txt
... '/root/lab/xato nom.txt'
```

`-print0` — nomlarni `\0` bilan ajratadi (fayl nomida bo'lishi mumkin bo'lmagan yagona belgi), `-0` — xargs shu formatda o'qiydi. **Fayllar ustida amal = har doim `-print0 | xargs -0`** (yoki `-exec ... +`).

Foydali xargs flaglari: `-I {}` (argument o'rnini belgilash: `xargs -I {} mv {} {}.bak`), `-P 4` (4 parallel process!), `-n 10` (buyruq boshiga 10 argument).

### `touch` va `stat` — yordamchilar

```bash
touch fayl              # bo'sh fayl yaratish YOKI mtime ni hozirga yangilash
touch -d "10 days ago" fayl    # vaqtni sun'iy o'rnatish (test uchun!)
```

```console
$ stat --format "%n | %s bayt | %U:%G | %a | %y" src/main.go
src/main.go | 0 bayt | root:root | 644 | 2026-07-10 10:21:43.196333000 +0000
```

`stat` — fayl metadata sining to'liq manbai (uch xil vaqt: atime o'qish, mtime kontent, ctime metadata).

## Real-world scenariylar

**1. Log tozalash cron jobi.** 30 kundan eski loglarni siqish, 90 kundan eskisini o'chirish:

```bash
find /var/log/myapp -name "*.log" -mtime +30 -not -name "*.gz" -exec gzip {} +
find /var/log/myapp -name "*.log.gz" -mtime +90 -delete
```

**2. Xavfsizlik auditi.** World-writable fayllar va egasiz fayllarni topish:

```bash
sudo find /etc /usr -type f -perm -o+w
sudo find / -xdev \( -nouser -o -nogroup \) 2>/dev/null
```

(`-xdev` — boshqa fayl tizimlarga (proc, docker mountlar) o'tmaslik.)

**3. Katta fayllar tergovi (12-dars davomi).** `du` katalogni ko'rsatdi, endi aniq fayllarni topamiz:

```bash
sudo find /var -xdev -type f -size +500M -exec ls -lh {} + 2>/dev/null
```

## Zamonaviy yondashuv

- **[fd](https://github.com/sharkdp/fd)** — find ning Rust dagi zamonaviy muqobili: `fd '\.go$'` (sodda sintaksis), default `.gitignore` ni hurmat qiladi, hidden fayllarni o'tkazib yuboradi, parallel (3-7x tez). Ubuntu paketi `fd-find`, buyruq nomi `fdfind` (tekshirilgan):

```console
$ fdfind "\.go$" ~/lab
/root/lab/src/main.go
/root/lab/src/main_test.go
/root/lab/src/util.go
```

Loyihada ishlashda fd qulay; scriptlar va serverlarda `find` universal (har joyda bor).
- **`plocate`** — locate ning tez zamonaviy implementatsiyasi (Ubuntu 22.04+ da default).
- **`fzf` bilan integratsiya**: `Ctrl+T` — fayl nomini interaktiv fuzzy-tanlash; `vim $(fzf)` patterni.
- **`-exec ... +` vs xargs**: zamonaviy find da `+` bo'lgani uchun oddiy holatlarda xargs shart emas; xargs `-P` (parallellik) va murakkab pipeline larda ustun.

## Keng tarqalgan xatolar

1. **`find . -name *.log` (quotesiz pattern).** Joriy katalogda .log fayl bo'lsa shell uni expand qilib yuboradi — natija tasodifiy. Har doim: `-name "*.log"`.

2. **`find | xargs` ni probelli nomlar bilan ishlatish.** Yuqorida isbotlandi — nomlar bo'linib ketadi. `-print0 | xargs -0` yoki `-exec {} +`.

3. **`-delete` ni sinamasdan ishga tushirish.** `-delete` testlardan OLDIN yozilsa (`find . -delete -name "*.tmp"`) — **hammasi** o'chadi (delete ham test kabi "true" qaytaradi va tartib muhim). Avval `-print` bilan dry-run, `-delete` esa har doim oxirida.

4. **Permission xatolar to'foni.** `find /` oddiy userda ekranni "Permission denied" ga to'ldiradi. Yechim: `2>/dev/null` (05-dars) yoki sudo.

5. **`-mtime 7` vs `-mtime +7` chalkashligi.** Belgisiz `7` — "aynan 7-8 kun oralig'ida", `+7` — "7 kundan eski", `-7` — "oxirgi 7 kunda". Log tozalashda kerakli deyarli har doim `+`.

6. **locate yangi faylni topolmasligi.** Baza kecha qurilgan. `sudo updatedb` yoki find ishlating.

## Amaliy mashqlar

Muhit: `docker run -it --rm ubuntu:24.04 bash`. Poligon yarating:

```bash
mkdir -p lab/{src,logs,backup} && cd lab
touch src/{main,util,main_test}.go logs/{app,error}.log
dd if=/dev/zero of=logs/big.log bs=1M count=5 2>/dev/null
touch -d "10 days ago" backup/old.tar
touch "xato nom.txt" empty.txt
```

**1.** `lab` da: faqat kataloglarni toping; faqat `.go` fayllarni toping; `_test.go` bo'lmagan `.go` fayllarni toping.

<details><summary>Yechim</summary>

```bash
find . -type d
find . -name "*.go"
find . -name "*.go" -not -name "*_test.go"
```
</details>

**2.** 1MB dan katta fayllarni hajmi bilan chiqaring.

<details><summary>Yechim</summary>

```console
$ find . -type f -size +1M -exec ls -lh {} +
... 5.0M ./logs/big.log
```
</details>

**3.** 7 kundan eski fayllarni toping. Nega `touch -d` bilan yaratganimiz chiqadi-yu boshqalar chiqmaydi?

<details><summary>Yechim</summary>

```bash
find . -type f -mtime +7      # ./backup/old.tar
```
`-mtime` faylning **modification time** metadata siga qaraydi — `touch -d "10 days ago"` aynan shuni orqaga surgan. Qolganlar hozir yaratildi.
</details>

**4.** Bo'sh fayl va kataloglarni toping, keyin **faqat bo'sh fayllarni** (kataloglarni emas) o'chiring — avval dry-run bilan.

<details><summary>Yechim</summary>

```bash
find . -empty                          # hammasi
find . -type f -empty -print           # dry-run
find . -type f -empty -delete
```
</details>

**5.** Probel muammosini o'zingiz ko'ring: `find . -name "*.txt" | xargs ls -l` va `-print0 | xargs -0` variantini solishtiring.

<details><summary>Yechim</summary>

Birinchisi `xato nom.txt` da ikkita xato beradi; ikkinchisi to'g'ri ishlaydi:
```bash
find . -name "*.txt" -print0 | xargs -0 ls -l
```
</details>

**6.** Barcha `.go` fayllarning nusxasini `.go.bak` sifatida yarating — bitta find/xargs qatori bilan.

<details><summary>Yechim</summary>

```bash
find . -name "*.go" -print0 | xargs -0 -I {} cp {} {}.bak
ls src/
```
`-I {}` — har argument uchun `{}` o'rniga qo'yish (bu rejimda har fayl alohida cp bilan).
</details>

**7.** (Qiyinroq) `find`ning o'zi bilan (xargs siz): logs katalogidagi `.log` fayllarni gzip qiling, lekin 1MB dan kichiklarini tegmang. Keyin natijani tekshiring.

<details><summary>Yechim</summary>

```bash
find logs -name "*.log" -size +1M -exec gzip -v {} +
ls -lh logs/     # big.log.gz paydo bo'ldi, kichiklar joyida
```
Shartlar zanjiri (nom + hajm) va `+` bilan samarali exec — production log-rotation logikasining yadrosi.
</details>

## Cheat sheet

| Buyruq | Nima qiladi | Eng ko'p ishlatiladigan variant |
|--------|-------------|--------------------------------|
| `locate` | Indeksdan tez qidiruv | `locate nom`, `sudo updatedb` |
| `find -name` | Nom bo'yicha | `find . -name "*.log"` (quote!) |
| `find -type` | Tur bo'yicha | `-type f` / `-type d` |
| `find -size` | Hajm | `-size +100M` |
| `find -mtime` | Yosh | `-mtime +7` (7 kundan eski) |
| `find -empty` | Bo'shlar | `-type f -empty` |
| Operatorlar | AND/OR/NOT | `\( -name a -o -name b \)`, `-not` |
| `-exec {} \;` | Har faylga buyruq | kichik ro'yxatlar |
| `-exec {} +` | To'plab buyruq | katta ro'yxatlar (tez) |
| `-delete` | O'chirish | oldin `-print` bilan dry-run! |
| `xargs -0` | Null-separated argumentlar | `find ... -print0 \| xargs -0 cmd` |
| `xargs -P` | Parallel | `xargs -0 -P4 gzip` |
| `stat` | Metadata | `stat fayl` |
| `fd` | Zamonaviy find | `fdfind pattern` (Ubuntu) |

## Qo'shimcha manbalar

- [GNU Findutils manual](https://www.gnu.org/software/findutils/manual/html_mono/find.html) — find ning to'liq rasmiy hujjati
- [fd — GitHub](https://github.com/sharkdp/fd) — zamonaviy muqobil
- [find, grep, xargs, newlines and null](https://v5.chriskrycho.com/journal/find-grep-xargs-newlines-null/) — -print0/-0 muammosi chuqur tahlili

---

[← Oldingi: 13 — networking](13-networking.md) · [Kurs xaritasi](00-README.md) · [Keyingi: 15 — archiving-and-sync →](15-archiving-and-sync.md)
