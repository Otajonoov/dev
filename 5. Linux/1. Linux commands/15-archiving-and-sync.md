# 15. Arxivlash va sinxronizatsiya

> Manba: TLCL 18-bob · Muhit: Ubuntu 24.04 · [← Oldingi: finding-files](14-finding-files.md) · [Kurs xaritasi](00-README.md) · [Keyingi: regular-expressions →](16-regular-expressions.md)

## Nima uchun kerak

Deploy artifactlari, DB backuplar, log arxivlari, `docker save` — backend hayoti arxivlar bilan to'la. `tar` sintaksisini har safar googlashdan charchagan bo'lsangiz — bu dars uni miyaga o'rnatadi. `rsync` esa bundan ham qimmatroq skill: faqat **farqni** ko'chiradi, uzilsa davom ettiradi — 100GB ni har safar boshidan scp qilish o'rniga. Bonus: 2026 da gzip emas, `zstd` ishlatish kerakligini benchmarkda ko'rasiz.

## Nazariya

### Arxivlash ≠ siqish

Unix falsafasida bu ikki **alohida** ish:

- **Arxivlash** (`tar`) — ko'p fayl/katalogni bitta faylga yig'ish (metadata: yo'llar, permissions, timestamps saqlanadi);
- **Siqish** (`gzip`, `zstd`, `xz`) — bitta faylni kichraytirish.

`archive.tar.gz` = avval tar, keyin gzip. Windows dagi `zip` ikkalasini bitta formatda qiladi — shuning uchun u ekotizimlar orasida almashinuvda ishlatiladi, Unix ichida esa tar+kompressor standarti.

Siqish **lossless** (ma'lumot to'liq qaytadi — matn, kod, DB uchun) va **lossy** (JPEG, MP3 — taxminiy qaytadi) bo'ladi; biz faqat lossless bilan ishlaymiz. Allaqachon siqilgan datani (JPEG, video, siqilgan arxiv) qayta siqish deyarli hech narsa bermaydi.

### Kompressorlar oilasi (2026 holati)

| Tool | Kengaytma | Xarakteri |
|------|-----------|-----------|
| `gzip` | .gz | Universal klassika — har joyda bor |
| `bzip2` | .bz2 | Eskirgan oraliq variant |
| `xz` | .xz | Eng zich, lekin juda sekin |
| `zstd` | .zst | **Zamonaviy g'olib**: gzip -9 darajasini ~9x tez beradi |

## Buyruqlar

### `gzip` / `gunzip`

Faylni **o'rnida** almashtiradi (tekshirilgan — ~4.7x kichraydi):

```console
$ ls -l /etc > foo.txt && ls -l foo.txt
... 6791 ... foo.txt
$ gzip foo.txt && ls -l foo.*
... 1455 ... foo.txt.gz
$ gunzip foo.txt.gz
```

Foydali flaglar: `-k` (originalni saqlab qolish), `-9`/`-1` (zichroq/tezroq), `-v`. Siqilganini ochmasdan o'qish — `zcat` / `zless` (tekshirilgan):

```console
$ zcat foo.txt.gz | head -2
total 620
drwxr-xr-x 3 root root    4096 Jul 10 10:17 X11
```

Log arxivlarini o'qishda har kuni kerak: `zcat access.log.2.gz | grep 500`. Bonus: `zgrep 500 access.log.*.gz` — to'g'ridan-to'g'ri.

### `tar` — arxivator

Rejimlar (birinchi harf): `c` create, `x` extract, `t` list, `r` append. `f arxiv` — fayl nomi (deyarli har doim kerak, **oxirida** turadi).

```console
$ tar cf playground.tar playground      # yaratish
$ tar tf playground.tar | head -4       # ko'rish (t = list)
playground/
playground/dir-2/
playground/dir-2/file-a
playground/dir-2/file-b
$ mkdir restore && cd restore
$ tar xf ../playground.tar              # ochish (JORIY katalogga)
```

Siqish bilan birga — bitta harf qo'shiladi: `z` gzip, `j` bzip2, `J` xz, `--zstd`:

```console
$ tar czf playground.tgz playground     # yaratish + gzip (10240 → 228 bayt)
$ tar xzf playground.tgz                # ochish (zamonaviy tar z ni o'zi sezadi)
```

Yodlash formulasi: **c**reate/e**x**tract + **z**ip + **f**ile → `czf` / `xzf`. Ko'rish uchun `v` (verbose) qo'shiladi: `tar xzvf`.

Arxivdan bitta faylni chiqarish (tekshirilgan):

```bash
tar xf playground.tar --wildcards "playground/dir-2/file-a"
```

Muhim semantika: tar **relative yo'llarni** saqlaydi (GNU tar boshidagi `/` ni avtomatik olib tashlaydi) — ochilganda joriy katalogga "playground/..." bo'lib tushadi. Shuning uchun: arxivni **doim `tar tf` bilan ko'rib** keyin oching — ba'zi arxivlar katalogsiz, to'g'ridan-to'g'ri fayllarni sochib yuboradi ("tarbomb").

### `zip` / `unzip`

Windows dunyosi bilan almashinuv uchun:

```console
$ zip -rq playground.zip playground     # -r rekursiv, -q quiet
$ unzip -l playground.zip | head -4     # ochmasdan ko'rish
Archive:  playground.zip
  Length      Date    Time    Name
---------  ---------- -----   ----
        0  2026-07-10 10:24   playground/
$ unzip playground.zip                  # ochish
```

### `zstd` — zamonaviy siqish (benchmark bilan)

25MB test fayl (20MB nol + 5MB random) — tekshirilgan natija:

```console
$ time gzip -k data.bin
real	0m0.196s
$ time zstd -q data.bin
real	0m0.016s
$ ls -l data.bin*
26214400 data.bin
 5266224 data.bin.gz
 5243677 data.bin.zst
```

**Bir xil siqish darajasi, 12x tezlik.** Katta backuplar uchun: `tar --zstd -cf backup.tar.zst dir/` yoki `tar -I 'zstd -T0' -cf ...` (`-T0` — barcha yadrolarda parallel). Qoida: yangi ish — zstd; eski tizimlar bilan moslik kerak — gzip; har bayt qimmat va vaqt ko'p — xz.

### `rsync` — aqlli sinxronizatsiya

rsync **delta-transfer** ishlatadi: manba va maqsadni solishtirib, faqat farqni ko'chiradi. Birinchi ishga tushirish — to'liq nusxa, keyingilari — sekundlar (tekshirilgan):

```console
$ rsync -av playground/ pg-copy/        # birinchi marta: hammasi
$ rsync -av playground/ pg-copy/        # ikkinchi marta: hech nima ko'chmadi
sent 213 bytes ...
$ touch playground/dir-1/yangi.txt
$ rsync -av playground/ pg-copy/        # faqat yangi fayl
sending incremental file list
dir-1/yangi.txt
```

Asosiy flaglar: `-a` (archive: rekursiv + permissions + timestamps + symlinks — deyarli har doim kerak), `-v`, `-z` (kanalda siqish — tarmoq orqali foydali), `--delete` (manbada yo'q fayllarni maqsaddan ham o'chirish — **ko'zgu** rejim), `-n`/`--dry-run`:

```console
$ rsync -av --delete playground/ pg-copy/
deleting dir-1/yangi.txt
$ rsync -avn --delete playground/ pg-copy/     # -n: faqat REJANI ko'rsatadi
```

`--delete` bilan ishlashdan oldin **har doim `-n` bilan dry-run** — yo'nalishni adashtirsangiz ma'lumot o'chadi.

**Trailing slash semantikasi** — rsync da (cp dan farqli!) muhim (tekshirilgan):

```console
$ rsync -a playground d1/      # slashsiz:  d1/playground/... (katalog O'ZI)
$ rsync -a playground/ d2/     # slash bilan: d2/dir-1... (KONTENTI)
```

**SSH orqali remote sinxronizatsiya** (tekshirilgan):

```console
$ rsync -az playground/ localhost:/tmp/pg-remote/
$ ssh localhost "ls /tmp/pg-remote"
dir-1  dir-2  dir-3 ...
```

Zamonaviy rsync default ssh ishlatadi (`-e ssh` shart emas). Uzilgan katta ko'chirishni davom ettirish: `rsync -az --partial --progress`.

## Real-world scenariylar

**1. Deploy artifact.** CI da build natijasini yig'ish va serverda ochish:

```bash
tar --zstd -cf app-v2.4.1.tar.zst -C build/ .     # -C: build/ ichidan (yo'llar toza)
# serverda:
mkdir -p /srv/app/releases/v2.4.1
tar --zstd -xf app-v2.4.1.tar.zst -C /srv/app/releases/v2.4.1
```

`-C` flagi — "shu katalogga o'tib ishla": arxiv ichida ortiqcha yo'l qatlamlari qolmaydi.

**2. DB backupni siqib remote ga oqizish** (diskka oraliq fayl yozmasdan — 13-darsdagi ssh-pipe):

```bash
pg_dump mydb | zstd -T0 | ssh backup-host "cat > /backups/mydb-$(date +%F).sql.zst"
```

**3. Katta katalogni serverlar orasida ko'chirish.** scp o'rniga rsync — uzilsa davom etadi, ikkinchi urinish faqat qolganini oladi:

```bash
rsync -az --partial --progress /data/uploads/ newserver:/data/uploads/
# cron da har kuni ko'zgu-backup:
rsync -az --delete /srv/app/ backup-host:/mirrors/app/
```

## Zamonaviy yondashuv

- **zstd hamma joyda**: docker image layerlari, btrfs, .deb paketlar (Ubuntu 21.10+), kernel modullar — industriya gzip dan ko'chib bo'ldi. `zstd --long -19 -T0` — arxivlash uchun, `zstd -3` (default) — kundalik.
- **Backup falsafasi**: tar+cron — boshlang'ich daraja; jiddiy backup uchun deduplikatsiyali toollar: **restic**, **borgbackup** (snapshot, shifrlash, retention policy). "3-2-1 qoidasi": 3 nusxa, 2 xil vosita, 1 tasi boshqa joyda.
- **`tar` GNU vs BSD**: macOS dagi tar (BSD) flaglari biroz farq qiladi; Linux scriptlarini Mac da ishlatganda ehtiyot bo'ling.
- Katta fayllarni tez ko'chirishda siqishni o'ylab yoqing: LAN ichida (1-10Gbit) siqish CPU ga tiqilib **sekinlashtirishi** mumkin; internet orqali — deyarli har doim foyda (`-z` yoki zstd pipe).

## Keng tarqalgan xatolar

1. **`tar xf` ni ko'rmasdan ishga tushirish — "tarbomb".** Arxiv ichi katalogsiz bo'lsa, joriy katalogingizga yuzlab fayl sochiladi. Oldin `tar tf arxiv | head`; yoki har doim yangi katalogga: `mkdir out && tar xf a.tar -C out`.

2. **rsync da trailing slash ni adashtirish.** `rsync -a src dst` vs `rsync -a src/ dst` — natija bitta daraja farq qiladi; `--delete` bilan birga bu ma'lumot yo'qotishga aylanadi. Qoida: ikkalasiga ham slash qo'ying (`src/ dst/`) va dry-run bilan tekshiring.

3. **`--delete` ni dry-runsiz.** Yo'nalish adashsa (bo'sh katalogdan to'liq katalogga emas, teskarisiga) — to'liq katalog "sinxronlanib" bo'shaydi. `-n` — bir soniya, ma'lumot — qaytmas.

4. **Siqilgan faylni qayta siqishga urinish.** `gzip access.log.gz` — foyda nol, vaqt behuda. JPEG/MP4/zip larni tar ga olayotganda ham `z` siz qilsangiz tezroq bo'ladi.

5. **tar da absolute path bilan arxiv yaratish.** `tar cf b.tar /etc/nginx` — ochilganda `etc/nginx/...` bo'lib joriy katalogga tushadi (GNU tar `/` ni olib tashlaydi), lekin eski/boshqa tar realizatsiyalarda to'g'ridan-to'g'ri `/etc` ni ustidan yozishga urinishi mumkin. Toza usul: `tar cf b.tar -C /etc nginx`.

6. **scp bilan katta katalogni ko'chirib, uzilishdan keyin boshidan boshlash.** rsync `--partial` bilan davom ettiradi — katta hajmlarda scp o'rniga rsync odat bo'lsin.

## Amaliy mashqlar

Muhit: `docker run -it --rm ubuntu:24.04 bash` (`apt update && apt install -y zip unzip rsync zstd`)

**1.** `/etc` ro'yxatidan fayl yasab (`ls -lR /etc > big.txt` — kattaroq bo'lsin), uni gzip/zstd/xz bilan siqib hajm va tezlikni solishtiring.

<details><summary>Yechim</summary>

```bash
ls -lR /etc > big.txt 2>/dev/null
time gzip -k big.txt && time zstd -qk big.txt && time xz -k big.txt
ls -l big.txt*
```
Kutilma: zstd — eng tez, xz — eng kichik, gzip — o'rtada.
</details>

**2.** `project/{src,docs}` strukturasini yarating, ichiga fayllar qo'ying va `.tar.gz` arxiv yasang. Arxivni **ochmasdan** ichini ko'ring.

<details><summary>Yechim</summary>

```bash
mkdir -p project/{src,docs} && touch project/src/a.go project/docs/readme.md
tar czf project.tgz project
tar tzf project.tgz
```
</details>

**3.** Arxivdan faqat `readme.md` ni chiqaring (butun arxivni ochmasdan).

<details><summary>Yechim</summary>

```bash
tar xzf project.tgz --wildcards "*/readme.md"
find project -name readme.md
```
</details>

**4.** Xuddi shu strukturani zip qiling va unzip -l bilan tar tf farqini ko'ring. Qaysi holatda zip tanlaysiz?

<details><summary>Yechim</summary>

```bash
zip -r project.zip project && unzip -l project.zip
```
Zip: Windows userlarga yuborish, alohida fayllarni tez olish (indeksli format). Unix ichki ishlari uchun tar+zstd.
</details>

**5.** rsync bilan `project/` ni `backup/` ga ko'chiring. Fayl qo'shib qayta sinxronlang — faqat yangi fayl ketganini ko'rsating. Keyin faylni o'chirib `--delete` bilan ko'zguni tekshiring.

<details><summary>Yechim</summary>

```bash
rsync -av project/ backup/
touch project/src/new.go
rsync -av project/ backup/          # faqat new.go
rm project/src/new.go
rsync -avn --delete project/ backup/   # dry-run: "deleting src/new.go"
rsync -av --delete project/ backup/
```
</details>

**6.** Trailing slash tajribasi: `rsync -a project d1/` va `rsync -a project/ d2/` natijalarini `ls` bilan solishtirib, qoidani o'z so'zingiz bilan yozing.

<details><summary>Yechim</summary>

```console
$ ls d1
project          # slashsiz: katalog O'ZI ko'chdi
$ ls d2
docs  src        # slash bilan: KONTENTI ko'chdi
```
Qoida: "manba/ = ichidagilar; manba = o'zi bilan".
</details>

**7.** (Qiyinroq) "Oqim ichida arxiv": `project` ni diskka oraliq arxiv yozmasdan boshqa katalogga tar-pipe orqali ko'chiring (permissions saqlangan holda).

<details><summary>Yechim</summary>

```bash
mkdir -p /tmp/dest
tar cf - project | (cd /tmp/dest && tar xf -)
find /tmp/dest | head
```
`-` = stdout/stdin. Bu pattern ssh bilan ham ishlaydi: `tar cf - dir | ssh host "cd /dst && tar xf -"` — rsync bo'lmagan joyda klassik yechim.
</details>

## Cheat sheet

| Buyruq | Nima qiladi | Eng ko'p ishlatiladigan variant |
|--------|-------------|--------------------------------|
| `tar czf` | Arxiv yaratish (+gzip) | `tar czf a.tgz dir/` |
| `tar xzf` | Ochish | `tar xzf a.tgz -C maqsad/` |
| `tar tzf` | Ichini ko'rish | ochishdan OLDIN |
| `tar --zstd -cf` | Zamonaviy siqish bilan | `tar --zstd -cf a.tar.zst dir/` |
| `gzip`/`gunzip` | Siqish/ochish | `gzip -k9 f`, `zcat f.gz` |
| `zstd`/`unzstd` | Tez siqish | `zstd -T0 f` |
| `zip -r`/`unzip` | Windows almashinuvi | `unzip -l` (ko'rish) |
| `rsync -av` | Sinxronizatsiya | `rsync -az src/ host:dst/` |
| `rsync --delete` | Ko'zgu rejimi | oldin `-n` (dry-run)! |
| `rsync --partial --progress` | Katta fayllar | uzilsa davom etadi |
| `zgrep`/`zcat` | Siqilgan loglarni o'qish | `zgrep ERROR app.log.3.gz` |

## Qo'shimcha manbalar

- [GNU tar manual](https://www.gnu.org/software/tar/manual/) — rasmiy hujjat
- [rsync man page](https://download.samba.org/pub/rsync/rsync.1) — flaglar to'liq ro'yxati
- [GZIP vs BZIP2 vs XZ vs ZSTD benchmark](https://changethisfile.com/blog/gzip-bzip2-xz-zstd) — kompressorlar taqqoslamasi

---

[← Oldingi: 14 — finding-files](14-finding-files.md) · [Kurs xaritasi](00-README.md) · [Keyingi: 16 — regular-expressions →](16-regular-expressions.md)
