# 17. Text processing

> Manba: TLCL 20-bob + 21-bob (printf) · Muhit: Ubuntu 24.04 · [← Oldingi: regular-expressions](16-regular-expressions.md) · [Kurs xaritasi](00-README.md) · [Keyingi: scripting-first-steps →](18-scripting-first-steps.md)

## Nima uchun kerak

"Access log dan eng ko'p so'ragan 10 ta IP", "ikki config farqi", "CSV dan 2 va 5-ustunlar", "deploy oldi/keyin diff" — bu vazifalar uchun Python script ochish — chumolini traktorda bosish. Unix text toollari (sort, cut, sed, awk...) pipeline da ulanib shu ishlarni **bir qatorda** hal qiladi. Bu — kursning eng "to'yimli" darsi: 16-darsdagi regex shu yerda samara beradi, natijasi esa har kungi log-analiz mahoratingiz.

## Nazariya

### Unix text modeli va tool tanlash

Ko'p Unix data — **qator = yozuv, maydon = ustun** (ajratuvchi: tab, `:`, probel). Shu model ustida har tool o'z rolida:

```mermaid
flowchart LR
    A[grep<br/>qatorlarni TANLAsh] --> B[cut / awk<br/>ustunlarni AJRATish]
    B --> C[sort + uniq<br/>TARTIB va HISOB]
    C --> D[sed<br/>matnni O'ZGARTIRish]
    D --> E[printf / awk<br/>FORMATlash]
```

Tanlash qoidasi: qator topish — `grep`; ustun/hisob-kitob — `awk`; oqimda almashtirish — `sed`; JSON — `jq`.

Test datamiz — `distros.txt` (tab-ajratilgan: nom, versiya, sana), barcha misollar verify qilingan.

## Buyruqlar

### `sort` — tartiblash

| Flag | Nima qiladi |
|------|-------------|
| `-n` | Raqam sifatida (aks holda "10" < "5"!) |
| `-r` | Teskari |
| `-k n` | n-maydon bo'yicha |
| `-t x` | Maydon ajratuvchisi |
| `-u` | Tartibla + dublikatlarni olib tashla |
| `-h` | Human-readable sonlar (2K, 1G) — `du -h` bilan |

```console
$ sort -k2 -n distros.txt | head -3        # versiya (2-maydon) bo'yicha raqamli
Fedora	5	03/20/2006
Ubuntu	6.10	10/26/2006
Ubuntu	7.10	10/18/2007
$ sort -t: -k7 /etc/passwd | head -2       # ':' ajratuvchi, 7-maydon (shell)
dev:x:1001:1001::/home/dev:/bin/bash
root:x:0:0:root:/root:/bin/bash
```

Klassik juftlik — disk yeguvchilar (12-darsdan tanish):

```console
$ du -s /usr/share/* | sort -nr | head -3
42436	/usr/share/vim
22088	/usr/share/perl
15180	/usr/share/doc
```

### `uniq` — dublikatlar (sort bilan juft)

`-c` — sanab chiqish (eng qimmat flag!), `-d` — faqat dublikatlar, `-u` — faqat yakkalar:

```console
$ cut -f1 distros.txt | sort | uniq -c | sort -nr
      4 Fedora
      3 Ubuntu
      2 SUSE
```

Mana u — **"top-N" patterni**: `... | sort | uniq -c | sort -nr | head`. Log analizning yarmi shu qolipda.

### `cut` — ustun kesish

`-f n` — maydon(lar), `-d x` — ajratuvchi (default: tab), `-c n-m` — belgi pozitsiyalari:

```console
$ cut -f1,2 distros.txt | head -2
SUSE	10.2
Fedora	10
$ cut -d: -f1,7 /etc/passwd | head -2
root:/bin/bash
daemon:/usr/sbin/nologin
```

`cut` sodda va tez, lekin ajratuvchi **bitta belgi** — "bir nechta probel" bilan ishlamaydi (u yerda `awk` yutadi).

### `paste` va `join` — birlashtirish

```console
$ paste names.txt dates.txt | head -2      # ustunlarni yonma-yon
SUSE	12/07/2006
Fedora	11/25/2008
$ join a.txt b.txt                          # umumiy kalit bo'yicha (SQL JOIN!)
1 apple qizil
2 banana sariq
```

`join` ikkala fayl **kalit bo'yicha sortlangan** bo'lishini talab qiladi.

### `comm` va `diff` — solishtirish

```console
$ comm f1.txt f2.txt        # 3 ustun: faqat-1chida / faqat-2chida / ikkalasida
a
		b
		c
		d
	e
$ comm -12 f1.txt f2.txt    # faqat umumiylar (1 va 2-ustunni bekitdik)
b
c
d
```

`diff` — o'zgarishlar tili (git ning asosi!). Zamonaviy standart format — unified (`-u`):

```console
$ diff -u f1.txt f2.txt
--- f1.txt	2026-07-10 10:33:45 +0000
+++ f2.txt	2026-07-10 10:33:45 +0000
@@ -1,4 +1,4 @@
-a
 b
 c
 d
+e
```

O'qish: `-` olib tashlangan, `+` qo'shilgan, `@@ -1,4 +1,4 @@` — qaysi qatorlar oralig'i. `patch` — diff ni qo'llash:

```console
$ diff -u f1.txt f2.txt > patchfile
$ patch f1-patched.txt < patchfile
patching file f1-patched.txt
```

`git diff`/`git apply` — xuddi shu mexanizmning evolyutsiyasi.

### `tr` — belgi almashtirish

```console
$ echo "kichik harflar" | tr a-z A-Z
KICHIK HARFLAR
$ echo "salom	dunyo" | tr "\t" ","        # tab → vergul (TSV→CSV)
salom,dunyo
$ echo "aabbccdd" | tr -d "b"               # -d: o'chirish
aaccdd
$ echo "aaabbbccc" | tr -s "abc"            # -s: ketma-ketlarni siqish
abc
```

`tr` — **belgi** darajasida (1:1); so'z/pattern almashtirish — sed ishi. Klassik: `tr -d "\r"` — Windows CRLF dan tozalash.

### `sed` — stream editor

Eng ko'p ishlatiladigan buyruq — almashtirish `s/pattern/almashtirma/flaglar`:

```console
$ echo "front" | sed "s/front/back/"
back
$ echo "aaa bbb aaa" | sed "s/aaa/ccc/g"    # g — qatordagi HAMMAsi (aks holda faqat 1chisi)
ccc bbb ccc
```

Ajratuvchini almashtirish mumkin (yo'llardagi `/` lar bilan qulay): `sed "s|/var/log|/tmp|"`. 

**Addressing** — qaysi qatorlarga qo'llash:

```console
$ sed -n "1,3p" distros.txt      # faqat 1-3 qatorlar (-n: default chiqarishni o'chir)
SUSE	10.2	12/07/2006
Fedora	10	11/25/2008
SUSE	11.0	06/19/2008
$ sed -n "/SUSE/p" distros.txt   # regexga mos qatorlar (grep ekvivalenti)
$ sed "/^Fedora/d" distros.txt   # Fedora qatorlarini O'CHIRish
```

Guruhlar bilan qayta joylash — sana formatini ISO ga (tekshirilgan):

```console
$ sed -E "s|([0-9]{2})/([0-9]{2})/([0-9]{4})|\3-\1-\2|" distros.txt | head -2
SUSE	10.2	2006-12-07
Fedora	10	2008-11-25
```

`-E` — ERE (16-dars!), `\1..\3` — qavsdagi guruhlarga havola. Faylni joyida tahrirlash: `sed -i` (**oldin backupsiz sinab ko'ring** yoki `-i.bak`):

```console
$ sed -i "s|/2006|/06|g" d2.txt
```

vim dagi `:%s/...` bilan bir xil sintaksis ekanini payqadingizmi? (10-dars) — bilim ko'chdi.

### `printf` — formatlangan chiqish (21-bobdan)

`echo` dan farqi — aniq format nazorati (C dagi printf, Go `fmt.Printf` bilan bir oila):

```console
$ printf "%s uchun %d ta pod, load: %.2f\n" myapp 3 0.6789
myapp uchun 3 ta pod, load: 0.68
$ printf "%-10s|%5d|\n" qator 42        # -10: chapga tekislab 10 joy; 5: o'ngga
qator     |   42|
$ printf "%05d\n" 42                     # nol bilan to'ldirish
00042
$ printf "%x\n" 255                      # hex
ff
```

Scriptlarda tabulyar hisobot chiqarish — printf; oddiy xabar — echo.

### `awk` — ustunlar tili (zamonaviy qo'shimcha)

Kitob awk ni chetlab o'tadi, lekin backend uchun bu majburiy minimum. Model: har qator uchun `pattern { action }`; `$1..$n` — maydonlar, `$0` — butun qator, `NF` — maydonlar soni, `-F` — ajratuvchi:

```console
$ awk "{print \$1}" distros.txt | sort -u
Fedora SUSE Ubuntu
$ awk -F: '$3 >= 1000 {print $1, $3}' /etc/passwd     # shart + tanlash
ubuntu 1000
dev 1001
$ awk '{sum += $2} END {print "jami:", sum}' distros.txt   # agregatsiya!
jami: 74.44
```

`cut` bir nechta probelda ojiz, `awk` esa default istalgan whitespace ni ajratuvchi deb oladi — `ps aux`, `ls -l` outputlari bilan ishlashda awk yagona to'g'ri tanlov.

### `jq` — JSON uchun (zamonaviy qo'shimcha)

Structured loglar davri asbobi (tekshirilgan):

```console
$ echo '{"level":"error","msg":"db down","ms":42}' | jq -r '.msg'
db down
$ cat logs.json | jq -r 'select(.level=="error") | .ms'
42
```

sed/awk qator modeliga qurilgan — JSON ning ichki strukturasi uchun jq yagona adekvat tanlov.

## Real-world scenariylar

**1. Access log dan top-10 IP:**

```bash
awk '{print $1}' access.log | sort | uniq -c | sort -nr | head -10
```

**2. Deploy oldi/keyin config diff.** Har deployda config snapshot olib solishtirish:

```bash
diff -u /etc/nginx/nginx.conf.bak-2026-07-01 /etc/nginx/nginx.conf
# yoki ikki serverni: diff <(ssh s1 cat /etc/app.conf) <(ssh s2 cat /etc/app.conf)
```

**3. Sekin endpointlar hisoboti** (nginx log oxirgi ustunda request_time deb faraz qilaylik):

```bash
awk '{sum[$7]+=$NF; cnt[$7]++} END {for (u in sum) printf "%-40s %8.3f %6d\n", u, sum[u]/cnt[u], cnt[u]}' access.log | sort -k2 -nr | head
```

(URL bo'yicha o'rtacha vaqt va soni — awk assotsiativ massivlari bilan bir qatorda.)

## Zamonaviy yondashuv

- **JSON loglar** → `jq`; **CSV katta hajmda** → csvkit / [Miller (mlr)](https://miller.readthedocs.io/) / DuckDB (`duckdb -c "select ... from 'file.csv'"`) — SQL bilan CSV.
- **`sd`** — sed ning soddalashtirilgan Rust muqobili (`sd 'eski' 'yangi' fayl`); sed bilimi baribir kerak (universal).
- **`diff` ranglarda**: `diff --color=auto -u`, yoki `git diff --no-index f1 f2` — git o'rnatilgan har joyda chiroyli diff.
- **`column -t`** — ustunlarni tekislab chiqarish: `mount | column -t`.
- Kitobdagi `aspell` (imlo tekshiruv) va 21-bobdagi `nl/fold/fmt/pr/groff` — bugungi backend ishida deyarli uchramaydi; printf dan boshqasini bilish shart emas.

## Keng tarqalgan xatolar

1. **`sort` da `-n` ni unutish.** `sort` matn sifatida: `10 < 5` (chunki "1"<"5"). Raqamlarda doim `-n` (yoki `du -h` bilan `-h`).

2. **`uniq` ni sortsiz ishlatish.** Faqat qo'shni dublikatlarni ko'radi (05-darsda aytilgan, ammo bu xato shu qadar keng tarqalganki — yana). `sort | uniq -c` — ayrilmas juftlik.

3. **`sed -i` ni birinchi urinishda ishlatish.** Regex xato bo'lsa fayl buzildi. Avval `-i` siz ekranga, tekshirib keyin `-i` (yoki `-i.bak`).

4. **`cut` ga ko'p probelli output berish.** `ps aux | cut -d" " -f2` — probellar soni o'zgaruvchan, natija tasodifiy. To'g'ri: `awk '{print $2}'`.

5. **`echo` bilan format chiqarishga urinish.** `echo` portativ emas (`-e` bash/dash da har xil), floatlarni formatlamaydi. Scriptda aniq format kerak bo'lsa — `printf`.

6. **JSON ni grep/sed bilan parse qilish.** Bir qatorli JSON da ishlagandek ko'rinadi, keyin ko'p qatorli/escaped data kelib sinadi. JSON → faqat `jq`.

## Amaliy mashqlar

Muhit: `docker run -it --rm ubuntu:24.04 bash` (`apt update && apt install -y jq gawk`). Tayyorlov:

```bash
printf "SUSE\t10.2\t12/07/2006\nFedora\t10\t11/25/2008\nSUSE\t11.0\t06/19/2008\nUbuntu\t8.04\t04/24/2008\nFedora\t8\t11/08/2007\nUbuntu\t6.10\t10/26/2006\nFedora\t9\t05/13/2008\nUbuntu\t7.10\t10/18/2007\nFedora\t5\t03/20/2006\n" > distros.txt
```

**1.** distros.txt ni: (a) nom bo'yicha; (b) versiya bo'yicha raqamli; (c) sana yili bo'yicha sortlang (maslahat: `-k` da maydon ichidagi pozitsiya `3.7` shaklida bo'ladi).

<details><summary>Yechim</summary>

```bash
sort distros.txt
sort -k2 -n distros.txt
sort -t$'\t' -k3.7 distros.txt      # 3-maydonning 7-belgisidan (yil)
```
</details>

**2.** Har bir distributiv nomi necha marta uchraydi — kamayish tartibida.

<details><summary>Yechim</summary>

```console
$ cut -f1 distros.txt | sort | uniq -c | sort -nr
      4 Fedora
      3 Ubuntu
      2 SUSE
```
</details>

**3.** `/etc/passwd` dan: shelli `/bin/bash` bo'lgan userlarning faqat nomlarini chiqaring — bir marta cut bilan, bir marta awk bilan.

<details><summary>Yechim</summary>

```bash
grep ":/bin/bash$" /etc/passwd | cut -d: -f1
awk -F: '$7 == "/bin/bash" {print $1}' /etc/passwd
```
</details>

**4.** sed bilan distros.txt dagi sanalarni `MM/DD/YYYY` dan `YYYY-MM-DD` ga o'tkazing (fayl o'zgarmasin, ekranga).

<details><summary>Yechim</summary>

```bash
sed -E "s|([0-9]{2})/([0-9]{2})/([0-9]{4})|\3-\1-\2|" distros.txt
```
</details>

**5.** distros.txt nusxasini olib, undan SUSE qatorlarini o'chiring va Fedora ni RedHat ga almashtiring — bitta sed chaqiruvida, faylning ichida (`-i`).

<details><summary>Yechim</summary>

```bash
cp distros.txt d.txt
sed -i -e "/^SUSE/d" -e "s/^Fedora/RedHat/" d.txt
cat d.txt
```
`-e` — bir nechta buyruqni ulash (yoki `sed "cmd1; cmd2"`).
</details>

**6.** printf bilan jadval: uch qatorli "servis | port | holat" hisobotini ustunlari tekis chiqaring (masalan api/8080/OK, db/5432/OK, cache/6379/DOWN).

<details><summary>Yechim</summary>

```console
$ printf "%-8s %6s %-6s\n" SERVIS PORT HOLAT api 8080 OK db 5432 OK cache 6379 DOWN
SERVIS     PORT HOLAT
api        8080 OK
db         5432 OK
cache      6379 DOWN
```
printf format stringni argumentlar tugaguncha qayta ishlatadi — bir chaqiruvda ko'p qator!
</details>

**7.** (Qiyinroq) awk bilan mini-hisobot: distros.txt dan har distributivning eng katta versiyasini chiqaring.

<details><summary>Yechim</summary>

```console
$ awk '{if ($2 > max[$1]) max[$1] = $2} END {for (d in max) print d, max[d]}' distros.txt
SUSE 11.0
Ubuntu 8.04
Fedora 9
```
(Raqamli taqqoslash uchun `$2+0` ishlatish mumkin; Fedora "10" vs "9" string sifatida — awk da `+0` bilan: `if ($2+0 > max[$1]+0)`.)
</details>

## Cheat sheet

| Buyruq | Nima qiladi | Eng ko'p ishlatiladigan variant |
|--------|-------------|--------------------------------|
| `sort` | Tartiblash | `sort -nr`, `sort -t: -k3 -n`, `sort -u` |
| `uniq` | Dublikatlar | `sort \| uniq -c \| sort -nr` (top-N!) |
| `cut` | Ustun kesish | `cut -d: -f1,7` |
| `paste` | Yonma-yon ulash | `paste f1 f2` |
| `join` | Kalit bo'yicha JOIN | sortlangan fayllarda |
| `comm` | To'plamlarni solishtirish | `comm -12` (kesishma) |
| `diff` | Farqlar | `diff -u old new` |
| `patch` | Diff ni qo'llash | `patch < fayl.diff` |
| `tr` | Belgi almashtirish | `tr a-z A-Z`, `tr -d "\r"` |
| `sed` | Oqim tahriri | `sed -E "s/old/new/g"`, `sed -n "/re/p"`, `-i.bak` |
| `printf` | Format | `printf "%-10s %5.2f\n" ...` |
| `awk` | Ustun + hisob | `awk '{print $2}'`, `'{s+=$1} END {print s}'` |
| `jq` | JSON | `jq -r '.field'`, `select(...)` |

## Qo'shimcha manbalar

- [sed & awk (learnbyexample)](https://learnbyexample.github.io/) — GNU sed va awk bo'yicha bepul chuqur kitoblar
- [jq manual](https://jqlang.github.io/jq/manual/) — rasmiy jq hujjati
- [The AWK Programming Language (Aho, Kernighan, Weinberger)](https://awk.dev/) — awk mualliflarining kitobi, 2-nashr

---

[← Oldingi: 16 — regular-expressions](16-regular-expressions.md) · [Kurs xaritasi](00-README.md) · [Keyingi: 18 — scripting-first-steps →](18-scripting-first-steps.md)
