# 07. Permissions

> Manba: TLCL 9-bob · Muhit: Ubuntu 24.04, bash 5.2 · [← Oldingi: expansion-and-quoting](06-expansion-and-quoting.md) · [Kurs xaritasi](00-README.md) · [Keyingi: processes →](08-processes.md)

## Nima uchun kerak

"Permission denied" — backend developer eng ko'p ko'radigan xatolardan biri: Docker volume ga yozib bo'lmayapti, deploy user servisni o'qiy olmayapti, `.ssh/authorized_keys` ishlamayapti (permissions noto'g'ri bo'lsa SSH indamay rad etadi!). Bu xatolarni "chmod 777 urib ko'raman" bilan emas, modelni **tushunib** yechish kerak — chunki 777 production da xavfsizlik teshigi. Linux ko'p userli tizim sifatida tug'ilgan, permission modeli — uning skeleti.

## Nazariya

### User, group, world

Har bir fayl **bitta user** (owner) va **bitta group** ga tegishli; qolgan hamma — **world** (others). Har uchala daraja uchun alohida huquqlar beriladi. Kimligingiz (tekshirilgan):

```console
$ id
uid=1001(dev) gid=1001(dev) groups=1001(dev),27(sudo)
```

- **uid** — user raqami (nomlar odam uchun; kernel faqat raqamni biladi — Docker muammosining ildizi shu, pastda)
- **gid** — asosiy guruh; qo'shimcha guruhlar ham bo'ladi (masalan `sudo`, `docker`)
- root — har doim `uid=0`

Bu ma'lumot oddiy matn fayllarda: `/etc/passwd` (userlar), `/etc/group` (guruhlar), `/etc/shadow` (parol hashlari). Format (tekshirilgan):

```console
$ grep -E "^(root|dev)" /etc/passwd
root:x:0:0:root:/root:/bin/bash
dev:x:1001:1001::/home/dev:/bin/bash
```

(login : parol-plaeysxolderi : uid : gid : izoh : home : login shell)

`/etc/shadow` ni esa oddiy user o'qiy olmaydi — permission modeli o'zini himoya qilishda (tekshirilgan):

```console
$ file /etc/shadow
/etc/shadow: regular file, no read permission
$ less /etc/shadow
/etc/shadow: Permission denied
$ ls -l /etc/shadow
-rw-r----- 1 root shadow 599 Jul 10 09:48 /etc/shadow
```

### rwx — fayl va katalogda BOSHQA ma'noda

`ls -l` dagi birinchi 10 belgi: `-rw-r--r--` = tur (1) + owner (3) + group (3) + world (3).

Turlar: `-` fayl, `d` katalog, `l` symlink, `c` character device, `b` block device.

| Huquq | Faylda | Katalogda |
|-------|--------|-----------|
| `r` | Kontentni o'qish | Ichidagi **nomlar ro'yxatini** ko'rish |
| `w` | Kontentni o'zgartirish/qirqish (o'chirish EMAS!) | Ichida fayl **yaratish/o'chirish/rename** |
| `x` | Dastur sifatida ishga tushirish | Katalogga **kirish** (`cd`, ichidagi fayllarga yetish) |

Ikki kutilmagan fakt:
1. **Faylni o'chirish huquqi faylda emas — katalogda.** Faylga `w` yo'q bo'lsa ham, katalogga `w` bor bo'lsa o'chira olasiz.
2. **Katalogda `r` va `x` mustaqil.** Amaliy isbot (tekshirilgan):

```console
$ chmod 644 secret          # r bor, x yo'q
$ ls secret
inside.txt                   # nomlarni KO'RADI
$ cat secret/inside.txt
cat: secret/inside.txt: Permission denied    # lekin ichiga KIRA OLMAYDI
$ chmod 111 secret          # x bor, r yo'q
$ ls secret
ls: cannot open directory 'secret': Permission denied
$ cat secret/inside.txt
data                         # nomni bilsangiz — o'qiy olasiz!
```

Shuning uchun kataloglar uchun standart: `r` va `x` **birga** (755, 750, 700).

### Octal arifmetika

`r=4, w=2, x=1` — uch bit yig'indisi:

```
rwx = 4+2+1 = 7    rw- = 4+2 = 6    r-x = 4+1 = 5    r-- = 4
```

`chmod 640` = owner `rw-`, group `r--`, world `---`. Eng ko'p ishlatiladiganlar: `600` (maxfiy fayl), `644` (oddiy fayl), `700` (shaxsiy katalog/script), `755` (umumiy katalog/dastur).

## Buyruqlar

### `chmod` — huquqlarni o'zgartirish

Ikki sintaksis. **Octal** (tekshirilgan):

```console
$ chmod 600 foo.txt && ls -l foo.txt
-rw------- 1 dev dev 0 Jul 10 09:49 foo.txt
$ chmod 644 foo.txt && ls -l foo.txt
-rw-r--r-- 1 dev dev 0 Jul 10 09:49 foo.txt
$ chmod 755 foo.txt && ls -l foo.txt
-rwxr-xr-x 1 dev dev 0 Jul 10 09:49 foo.txt
```

**Symbolic**: kim (`u`/`g`/`o`/`a`) + amal (`+`/`-`/`=`) + nima (`r`/`w`/`x`):

```console
$ chmod u-x,go-rx foo.txt && ls -l foo.txt
-rw------- 1 dev dev 0 Jul 10 09:49 foo.txt
$ chmod go+r foo.txt && ls -l foo.txt
-rw-r--r-- 1 dev dev 0 Jul 10 09:49 foo.txt
```

Qachon qaysi: aniq holat o'rnatish — octal (`chmod 644`); mavjudiga nisbatan tuzatish — symbolic (`chmod +x`). Rekursiv: `chmod -R` (ehtiyot: fayl va kataloglarga bir xil huquq berib yuboradi; farqlash uchun `chmod -R u+rwX` — katta `X` faqat kataloglar va allaqachon executable fayllarga `x` beradi).

Klassik amaliyot — script ni ishga tushirish huquqi (tekshirilgan):

```console
$ echo "echo salom dunyo" > hello.sh
$ ./hello.sh
-bash: ./hello.sh: Permission denied
$ chmod +x hello.sh && ./hello.sh
salom dunyo
```

### `umask` — yangi fayllarning default huquqi

Yangi fayl nazariy jihatdan 666 (fayl) / 777 (katalog) bilan tug'iladi; umask undan bitlarni **olib tashlaydi** (tekshirilgan):

```console
$ umask
0002
$ > foo.txt && ls -l foo.txt
-rw-rw-r-- 1 dev dev 0 Jul 10 09:49 foo.txt      # 666-002=664
$ umask 0000 && > bar.txt && ls -l bar.txt
-rw-rw-rw- 1 dev dev 0 Jul 10 09:49 bar.txt      # hech nima olinmadi
$ umask 0077 && > baz.txt && ls -l baz.txt
-rw------- 1 dev dev 0 Jul 10 09:49 baz.txt      # group/world dan hammasi olindi
```

Server best practice: umumiy serverlarda `umask 027` yoki `077` — yangi fayllar default yopiq bo'ladi.

### Maxsus bitlar: setuid, setgid, sticky

To'rtinchi octal raqam: `4000` setuid, `2000` setgid, `1000` sticky.

- **setuid** (`-rws------`): dastur **egasining huquqi bilan** ishlaydi. Klassik misol — `passwd`: oddiy user `/etc/shadow` ni o'zgartira olmaydi, lekin parolini almashtirishi kerak (tekshirilgan):

```console
$ ls -l /usr/bin/passwd /usr/bin/sudo
-rwsr-xr-x 1 root root  72056 May 30  2024 /usr/bin/passwd
-rwsr-xr-x 1 root root 335120 Mar  2 12:56 /usr/bin/sudo
```

`s` — x o'rnida. Xavfsizlik: setuid-root dasturlar soni minimal bo'lishi kerak; o'z scriptlaringizga hech qachon setuid bermang.

- **setgid katalogda** (`drwxrwsr-x`): ichida yaratilgan fayllar **katalog guruhini** meros oladi — jamoaviy kataloglar uchun (`chmod g+s dir`).
- **sticky katalogda** (`drwxrwxrwt`): hamma yoza oladi, lekin faqat **o'z faylini** o'chira oladi. `/tmp` shunday (tekshirilgan):

```console
$ ls -ld /tmp
drwxrwxrwt 1 root root 4096 Jul 10 09:35 /tmp
```

### `su` va `sudo` — boshqa identity

**`su -`** — boshqa userning (default: root) **to'liq login shellini** ochish; **uning** parolini so'raydi:

```bash
su -              # root shell (root paroli kerak)
su - postgres     # postgres user bo'lib
su - dev -c "whoami; pwd"     # bitta buyruq
```

```console
$ su - dev -c "whoami; pwd"
dev
/home/dev
```

(`-` muhim: login shell — userning environmenti va home katalogi yuklanadi.)

**`sudo`** — zamonaviy standart: bitta buyruqni root (yoki boshqa user) sifatida bajarish, **o'z parolingiz** bilan. Kim nimaga haqli — `/etc/sudoers` da (tahrirlash faqat `visudo` orqali!):

```console
$ sudo whoami
root
$ sudo -l          # menga nima ruxsat etilgan?
User dev may run the following commands on 1af0fe06007c:
    (ALL : ALL) ALL
```

`su` vs `sudo` farqi: sudo — granulyar (aniq buyruqlargacha cheklash mumkin), audit log yozadi, root parolini tarqatishni talab qilmaydi. Production da: root bilan to'g'ridan-to'g'ri kirish o'chiriladi, adminlar sudo ishlatadi.

### `chown` / `chgrp` — egasini almashtirish

Faqat root (yoki sudo) qila oladi (tekshirilgan):

```console
$ sudo chown dev report.txt && ls -l report.txt
-rw-r--r-- 1 dev root 0 Jul 10 09:49 report.txt      # faqat user
$ sudo chown dev:dev report.txt && ls -l report.txt
-rw-r--r-- 1 dev dev 0 Jul 10 09:49 report.txt       # user:group
$ sudo chown root: report.txt && ls -l report.txt
-rw-r--r-- 1 root root 0 Jul 10 09:49 report.txt     # "root:" = user + uning guruhi
$ sudo chgrp dev report.txt && ls -l report.txt
-rw-r--r-- 1 root dev 0 Jul 10 09:49 report.txt      # faqat group
```

Oddiy user o'z faylini ham boshqaga "bera olmaydi":

```console
$ chown root ~/foo.txt
chown: changing ownership of '/home/dev/foo.txt': Operation not permitted
```

Rekursiv: `sudo chown -R appuser:appgroup /srv/myapp` — deploy kataloglarini sozlashda standart.

## Real-world scenariylar

**1. Docker volume "Permission denied".** Konteyner `uid=999` (masalan postgres image) bilan ishlaydi, host katalogi esa sizning `uid=1000` ingizga tegishli. Kernel **faqat raqamlarni** taqqoslaydi — nomlar mos kelishi ahamiyatsiz. Yechimlar (yaxshidan yomonga):

```bash
# 1) konteyner uid iga moslash:
sudo chown -R 999:999 ./pgdata
# 2) konteynerni o'z uid ingiz bilan ishga tushirish:
docker run --user "$(id -u):$(id -g)" -v "$PWD/data:/data" myapp
# 3) chmod 777 ./pgdata  ← HECH QACHON production da
```

**2. SSH kalitlari ishlamayapti.** `ssh` juda talabchan: `~/.ssh` — `700`, `authorized_keys`/kalitlar — `600` bo'lishi shart, aks holda **indamay** parol so'rayveradi:

```bash
chmod 700 ~/.ssh && chmod 600 ~/.ssh/authorized_keys ~/.ssh/id_*
```

**3. Deploy user minimal huquq bilan.** Go binary deploy qilinadigan server:

```bash
sudo useradd -r -s /usr/sbin/nologin appuser     # login qilolmaydigan servis user
sudo chown -R appuser:appuser /srv/myapp
sudo chmod 750 /srv/myapp                        # world uchun yopiq
sudo chmod 640 /srv/myapp/config.yaml            # config: owner rw, group r
```

CI/CD da esa sudoers ga aniq buyruqgina ochiladi: `deploy ALL=(ALL) NOPASSWD: /usr/bin/systemctl restart myapp` — butun root emas.

## Zamonaviy yondashuv

- **`chmod 777` — anti-pattern.** "Ishlab ketdi" degani "to'g'ri" degani emas: istalgan user (shu jumladan buzilgan servis) faylni o'zgartira oladi. Web-upload katalogiga 777 berish — attacker scriptini yozdirib olishning klassik yo'li. Standart: kataloglar 755, fayllar 644, maxfiylar 600, undan keyin **aniq ehtiyojga qarab** toraytirish/kengaytirish.
- **ACL (`setfacl`/`getfacl`)** — user/group/world uchligi yetmasa: bitta faylga bir nechta user/guruhga alohida huquq. `setfacl -m u:jenkins:r /var/log/app.log`. Klassik modeldan murakkabroq, kerak bo'lgandagina.
- **Capabilities** — "root yoki hech nima" o'rniga: `setcap cap_net_bind_service=+ep ./server` — Go binary 80-portni **root bo'lmasdan** eshitadi. setuid-root ning xavfsiz muqobili.
- **Konteynerlarda rootless yo'nalishi**: image da `USER app` directivasi, rootless Docker/Podman. "Konteyner ichida root" — hostdagi uid 0 bilan bir xil raqam ekanini unutmang.
- `stat -c '%a %U:%G %n' fayl` — huquqni octal ko'rish (`ls -l` dan aniqroq script uchun).

## Keng tarqalgan xatolar

1. **Har muammoga `chmod 777`.** Ishlaydi, chunki hamma narsani ochib tashlaydi — muammoni emas, himoyani yo'q qiladi. To'g'ri yo'l: `ls -l` + `id` bilan **kim kirolmayotganini** aniqlab, minimal huquq berish.

2. **`chmod -R 644 dir/` — kataloglar "sinishi".** Kataloglardan `x` olib tashlanadi — ichiga kirib bo'lmay qoladi. To'g'ri: `chmod -R u=rwX,go=rX dir/` (katta X) yoki fayl/katalogni alohida: `find dir -type f -exec chmod 644 {} +`.

3. **"Faylga w yo'q — o'chirib bo'lmaydi" deb o'ylash.** O'chirish — katalog operatsiyasi. Read-only fayl ham katalogga `w` bo'lsa o'chadi (`rm` faqat qo'shimcha so'raydi).

4. **`sudo echo x > /etc/conf` ishlamasligi.** Redirect sizning shellingizda ochiladi (05-darsda ko'rdik). To'g'ri: `echo x | sudo tee /etc/conf`.

5. **`/etc/sudoers` ni to'g'ridan-to'g'ri tahrirlash.** Sintaksis xatosi = sudo butunlay ishlamay qoladi = tizimga root kira olmaysiz. Faqat `sudo visudo` (saqlashdan oldin tekshiradi) yoki `/etc/sudoers.d/` + `visudo -cf` bilan.

6. **Docker da uid/gid raqam ekanini unutish.** Host da `dev(1000)`, konteynerda `node(1000)` — kernel uchun bu **bitta identity**. Fayllar "begona user niki" ko'rinsa — `ls -ln` bilan raqamlarni ko'ring.

## Amaliy mashqlar

Muhit: `docker run -it --rm ubuntu:24.04 bash`, ichida test user yarating: `useradd -m -s /bin/bash talaba`

**1.** `id` va `ls -ln ~` yordamida: sizning uid ingiz nechchi? `talaba` niki-chi? `/etc/passwd` dan tekshiring.

<details><summary>Yechim</summary>

```console
# id
uid=0(root) gid=0(root) groups=0(root)
# grep talaba /etc/passwd
talaba:x:1001:1001::/home/talaba:/bin/bash
```
Konteynerda birinchi yaratilgan user odatda 1000 yoki 1001 oladi (Ubuntu 24.04 image da `ubuntu` useri 1000 ni band qilgan bo'lishi mumkin).
</details>

**2.** `maxfiy.txt` yarating va shunday sozlangki: faqat owner o'qiy va yoza olsin. Ikki usulda (octal va symbolic) qiling, `talaba` user bilan o'qishga urinib isbotlang.

<details><summary>Yechim</summary>

```console
# echo "sir" > maxfiy.txt && chmod 600 maxfiy.txt        # octal
# chmod u=rw,go= maxfiy.txt                              # symbolic (xuddi shu)
# su - talaba -c "cat /root/maxfiy.txt"
cat: /root/maxfiy.txt: Permission denied
```
(Aslida talaba /root ga ham kira olmaydi — katalog `x` huquqi ham yo'q. Ikki qavat himoya.)
</details>

**3.** Katalogda `r` va `x` farqini o'zingiz ko'rsating: `demo/` katalogiga faqat `r` bering — nima ishlaydi? Faqat `x` bering — nima ishlaydi?

<details><summary>Yechim</summary>

```console
# mkdir demo && echo ichki > demo/f.txt && chown -R talaba demo
# chmod 444 demo && su - talaba -c "ls /root/demo; cat /root/demo/f.txt"
f.txt
cat: /root/demo/f.txt: Permission denied
# chmod 111 demo && su - talaba -c "ls /root/demo; cat /root/demo/f.txt"
ls: cannot open directory '/root/demo': Permission denied
ichki
```
`r` — ro'yxat, `x` — kirish. (Bu mashq uchun demo /root emas, /tmp ostida bo'lgani qulayroq.)
</details>

**4.** umask ni `077` qilib fayl va katalog yarating. Huquqlarini tushuntiring. Nega katalog `700` bo'ldi, fayl esa `600`?

<details><summary>Yechim</summary>

```console
# umask 077 && touch f && mkdir d && ls -ld f d
drwx------ 2 root root 4096 ... d
-rw------- 1 root root    0 ... f
```
Boshlang'ich baza fayl uchun 666 (x berilmaydi!), katalog uchun 777. 666-077→600, 777-077→700. Shuning uchun yangi fayllar hech qachon avtomatik executable bo'lmaydi.
</details>

**5.** `passwd` dasturi qanday qilib oddiy userga `/etc/shadow` ni o'zgartirish imkonini beradi? Buni `ls -l` output i bilan isbotlang va bit nomini ayting.

<details><summary>Yechim</summary>

```console
# ls -l /usr/bin/passwd
-rwsr-xr-x 1 root root 72056 ... /usr/bin/passwd
```
Owner triadasidagi `s` — **setuid bit**: dastur ishga tushganda effective uid = fayl egasi (root) bo'ladi. Dastur o'zi tekshiradi: user faqat O'Z parolini almashtira oladi.
</details>

**6.** `talaba` uchun sudo ni faqat **bitta buyruqqa** oching: `apt update`. Tekshiring: `apt update` ishlaydi, `apt install` esa yo'q.

<details><summary>Yechim</summary>

```console
# apt-get install -y sudo
# echo "talaba ALL=(root) NOPASSWD: /usr/bin/apt update" > /etc/sudoers.d/talaba
# visudo -cf /etc/sudoers.d/talaba        # sintaksisni TEKSHIRISH shart!
/etc/sudoers.d/talaba: parsed OK
# su - talaba -c "sudo apt update" | tail -1     # ishlaydi
# su - talaba -c "sudo apt install htop"
Sorry, user talaba is not allowed to execute '/usr/bin/apt install htop' as root...
```
Granulyar sudo — CI/CD deploy userlari uchun standart pattern.
</details>

**7.** (Qiyinroq) Jamoaviy katalog qurng: `shared/` — `devs` guruhidagi hamma yoza olsin, yangi fayllar avtomatik `devs` guruhiga tegishli bo'lsin, va hech kim boshqaning faylini o'chira olmasin.

<details><summary>Yechim</summary>

```console
# groupadd devs && mkdir /srv/shared
# chown root:devs /srv/shared
# chmod 2770 /srv/shared          # 2=setgid + rwxrwx---
# chmod +t /srv/shared            # sticky qo'shish (birga: chmod 3770)
# ls -ld /srv/shared
drwxrws--T 2 root devs 4096 ... /srv/shared
```
`s` (group o'rnida) — setgid: fayllar guruhni meros oladi; `T` — sticky: faqat egasi o'chiradi. `/tmp` (drwxrwxrwt) xuddi shu g'oyaning world varianti.
</details>

## Cheat sheet

| Buyruq | Nima qiladi | Eng ko'p ishlatiladigan variant |
|--------|-------------|--------------------------------|
| `id` | Kimligim (uid/gid/guruhlar) | `id`, `id username` |
| `chmod` | Huquq o'zgartirish | `chmod 600 f`, `chmod +x f`, `chmod -R u=rwX,go=rX d` |
| `umask` | Default huquq maskasi | `umask 027` (server), `umask 077` (paranoid) |
| `su` | Boshqa user shelli | `su - user`, `su - -c 'cmd'` |
| `sudo` | Bitta buyruq root sifatida | `sudo cmd`, `sudo -l`, `sudo visudo` |
| `chown` | Egasini almashtirish | `sudo chown -R user:group dir/` |
| `chgrp` | Guruhni almashtirish | `sudo chgrp devs f` |
| `stat` | Octal ko'rish | `stat -c '%a %U:%G' f` |
| Qiymatlar | — | 600 maxfiy, 644 fayl, 700/755 katalog, `.ssh`=700, kalit=600 |
| Maxsus | — | 4000 setuid, 2000 setgid, 1000 sticky |

## Qo'shimcha manbalar

- [Linux File Permissions — Arch Wiki](https://wiki.archlinux.org/title/File_permissions_and_attributes) — chuqur va aniq spravochnik
- [Why chmod 777 is dangerous](https://www.baeldung.com/linux/why-is-chmod-r-777-destructive) — 777 anti-patterni tahlili
- [Docker: fix volume permission denied](https://oneuptime.com/blog/post/2026-01-24-fix-permission-denied-docker-volumes/view) — uid/gid mismatch yechimlari

---

[← Oldingi: 06 — expansion-and-quoting](06-expansion-and-quoting.md) · [Kurs xaritasi](00-README.md) · [Keyingi: 08 — processes →](08-processes.md)
