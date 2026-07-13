# 18. Scripting: birinchi qadamlar

> Manba: TLCL 24-26 boblar · Muhit: Ubuntu 24.04, bash 5.2 · [← Oldingi: text-processing](17-text-processing.md) · [Kurs xaritasi](00-README.md) · [Keyingi: branching →](19-branching.md)

## Nima uchun kerak

17 dars davomida yig'ilgan buyruqlar arsenali endi **dasturlarga** aylanadi. Har kuni terminalda takrorlanayotgan 5 ta buyruq — bu yozilmagan script. Deploy, backup, log tozalash, environment tayyorlash — Go dasturchi sifatida siz baribir bash script yozasiz (Makefile, CI stepplari, Dockerfile RUN lari — hammasi shell). Farq faqat: **tushunib** yozasizmi yoki stackoverflow dan yamab. Blok D (18-24 darslar) — tushunib yozish.

## Nazariya

### Script nima?

Script — buyruqlar ketma-ketligi yozilgan oddiy matn fayli. Shell uni o'qib, xuddi siz qo'lda tergandek bajaradi. Uch shart:

1. **Yozish** — istalgan matn muharriri (syntax highlighting bilan — vim, VS Code).
2. **Executable qilish** — `chmod 755` (hamma uchun) yoki `700` (faqat o'zingiz). 07-dars bilimlari ishga tushdi.
3. **PATH ga joylash** — yoki `./script` deb aniq yo'l bilan chaqirish.

### Shebang — `#!`

Fayl birinchi qatori `#!/bin/bash` — kernel ga "bu faylni qaysi interpretator bilan bajarish"ni aytadi. `#!` (shebang) siz fayl bajarilsa, joriy shell taxmin qiladi — har doim yozing. Zamonaviy tavsiya: `#!/usr/bin/env bash` — bash ni PATH dan topadi (macOS/nix kabi bash boshqa joyda turgan tizimlar uchun portativroq).

`#` bilan boshlangan qolgan hamma narsa — komment (qator oxirida ham bo'ladi).

### Script qayerda tursin?

- `~/bin` — shaxsiy scriptlar (PATH ga qo'shiladi: 09-dars)
- `/usr/local/bin` — tizimdagi barcha userlar uchun (FHS bo'yicha "qo'lda o'rnatilgan" joy — 02-dars)
- Loyiha ichida — `./scripts/` katalogi (repoda versiyalanadi)

## Buyruqlar

### Birinchi script — to'liq sikl (verify qilingan)

```bash
#!/bin/bash
# Bu bizning birinchi scriptimiz.
echo "Hello World!"
```

```console
$ ls -l hello_world
-rw-r--r-- ...                        # x yo'q hali
$ ./hello_world
bash: ./hello_world: Permission denied
$ chmod 755 hello_world
$ ./hello_world
Hello World!
$ cp hello_world ~/bin/ && cd /tmp
$ hello_world                          # endi istalgan joydan!
Hello World!
```

Nega `./` kerak edi? Xavfsizlik: shell joriy katalogdan **qidirmaydi** (PATH da `.` yo'q) — aks holda birov `/tmp` ga `ls` nomli zararli fayl tashlab qo'yishi mumkin edi.

### Variables va constants

```bash
title="Sahifa sarlavhasi"        # = atrofida PROBEL YO'Q (04-darsdagi alias kabi)
TITLE="Tizim hisoboti"           # konvensiya: konstantalar KATTA, o'zgaruvchilar kichik
CURRENT_TIME="$(date +"%x %r %Z")"    # command substitution natijasi
```

Ishlatish — `$title` yoki aniqroq `${title}` (qo'shni matndan ajratish kerak bo'lsa: `${title}_v2`). **Har doim quote ichida**: `"$title"` (06-dars — bu yerda ham amal qiladi).

### Here document — ko'p qatorli matn

`cat << belgi` dan `belgi` gacha hamma narsa stdin ga boradi; ichida `$var` va `$(...)` ishlaydi:

```bash
#!/bin/bash
TITLE="Tizim hisoboti: $HOSTNAME"
CURRENT_TIME="$(date +"%x %r %Z")"
TIMESTAMP="$USER tomonidan $CURRENT_TIME da yaratildi"

cat << _EOF_
<html>
  <head><title>$TITLE</title></head>
  <body>
    <h1>$TITLE</h1>
    <p>$TIMESTAMP</p>
  </body>
</html>
_EOF_
```

Tekshirilgan natija:

```console
$ ./sys_info
<html>
  <head><title>Tizim hisoboti: 1af0fe06007c</title></head>
  <body>
    <h1>Tizim hisoboti: 1af0fe06007c</h1>
    <p>root tomonidan 07/10/26 10:37:22 AM UTC da yaratildi</p>
  </body>
</html>
```

Ikki variatsiya:
- `<< "_EOF_"` (limit string **quotelangan**) — ichidagi `$` **expand bo'lmaydi** (tekshirilgan) — config shablonlarida literal `$` kerak bo'lsa;
- `<<-_EOF_` (defis bilan) — boshidagi **tab** larni olib tashlaydi (indentatsiya qilingan heredoc).

Amaliyot: heredoc — scriptdan config fayl generatsiya qilish, ko'p qatorli usage xabari, SQL blokini psql ga berish uchun standart usul.

### Shell functions va `local`

Kod takrorlanishini yig'ish — funksiya. Ikki sintaksis (birinchisi keng tarqalgan):

```bash
report_disk() {
    local title="Disk holati"     # local — FAQAT funksiya ichida yashaydi
    echo "=== $title ==="
    df -h / | tail -1
}
```

Tekshirilgan (local isboti bilan):

```console
$ ./funcs.sh
=== Disk holati ===
overlay         453G   13G  417G   4% /
=== Uptime ===
 10:37:22 up  1:21,  0 user,  load average: 0.00, 0.05, 0.17
title tashqarida: [yoq]        # local o'z ishini qildi
```

Qoidalar:
- Funksiya **chaqirilishidan oldin** ta'riflangan bo'lishi kerak (script boshida funksiyalar, oxirida chaqiruvlar);
- Funksiya ichida `local` siz e'lon qilingan variable — **global** (bug manbai №1!);
- Funksiya — bu 04-darsdagi savolga javob: parametr oladigan "alias" (parametrlar 21-darsda: `$1`, `$2`).

### Top-down design va stub lar

Katta vazifani yozish usuli: avval yuqori darajani skeletini yozing, funksiyalarni **stub** (vaqtinchalik plug) qilib:

```bash
report_disk() { echo "report_disk: KEYINROQ"; }    # stub
report_uptime() { echo "report_uptime: KEYINROQ"; }

report_disk
report_uptime
```

Script **har doim ishlaydigan holatda** rivojlanadi — har qadamda test qilib borasiz (kitobning eng qimmatli maslahati: "keep the script running").

### Chiroyli kod

```bash
# Uzun buyruqlarni davom ettirish — \ va indentatsiya:
find playground \
    \( -type f -not -perm 0600 \) \
    -exec chmod 0600 '{}' ';'

ls --all --directory       # scriptda uzun flaglar — o'qiluvchanroq (-ad o'rniga)
```

Script — bir marta yoziladi, ko'p o'qiladi: kelajakdagi o'zingiz uchun yozing.

## Real-world scenariylar

**1. Health-check hisoboti.** Har servisda takrorlanadigan tekshiruvlarni bitta scriptga:

```bash
#!/usr/bin/env bash
report_service() {
    local name="$1"
    echo "--- $name ---"
    systemctl is-active "$name" || true
}
echo "Hisobot: $(hostname), $(date +%F)"
report_service nginx
report_service postgresql
df -h / | tail -1
```

**2. Config generatsiya (heredoc bilan).** CI da environmentga qarab config yasash:

```bash
cat > app.env << _EOF_
APP_ENV=$DEPLOY_ENV
DATABASE_URL=$DATABASE_URL
BUILD_SHA=$(git rev-parse --short HEAD)
_EOF_
```

**3. Takrorlanadigan buyruqlar to'plami → `~/bin`.** Har kuni ishlatadigan zanjiringizni scriptga aylantiring (masalan `logs-errors`: `ssh prod "tail -500 /var/log/app.log" | grep -E "ERROR|panic"`). Bir hafta ichida `~/bin` sizning shaxsiy tool to'plamingizga aylanadi.

## Zamonaviy yondashuv

- **Zamonaviy script skeleti** (24-darsda to'liq ochamiz, hozirdan shu qolipda yozing):

```bash
#!/usr/bin/env bash
set -euo pipefail

log() { echo "[$(date +%T)] $*" >&2; }

main() {
    log "boshlandi"
    # ...
}

main "$@"
```

`main` patterni: butun fayl avval o'qiladi, keyin bajariladi (yarim yuklangan scriptning ishga tushib ketishidan himoya); logika tartibli.
- **UPPER_CASE faqat export/konstanta** uchun, oddiy variablelar — `lower_case` (env variable bilan to'qnashuvdan qochish).
- **ShellCheck** ni birinchi kundan ishlating (VS Code extension) — xatolarni yozayotganingizda ko'rsatadi.
- Qachon bash, qachon Go/Python? Empirik qoida: **~100 qatordan yoki assotsiativ data strukturalardan oshsa** — bash dan chiqing. Bash — glue (yelimlash) tili, application tili emas.

## Keng tarqalgan xatolar

1. **`var = "qiymat"` (probel bilan).** Bash buni "var buyrug'ini = va qiymat argumentlari bilan chaqir" deb tushunadi. To'g'ri: `var="qiymat"` — probelsiz.

2. **Shebang ni unutish yoki 2-qatorga yozish.** `#!` faqat faylning **eng birinchi** ikki bayti sifatida ishlaydi. Undan oldin bo'sh qator ham bo'lmasin.

3. **`chmod +x` ni unutib "Permission denied" ga hayron bo'lish.** Yuqoridagi verify sessiyasida ataylab ko'rsatildi — bu har bir yangi script bilan bo'ladigan ritual.

4. **Funksiyani ta'rifidan oldin chaqirish.** Bash interpretator — funksiya chaqiruv paytida mavjud bo'lishi kerak. `main` patterni bu muammoni tizimli hal qiladi.

5. **`local` ni unutish.** Funksiya ichidagi `i`, `result` kabi variablelar global bo'lib, boshqa funksiyanikini ustidan yozadi. Funksiya ichidagi HAR variable — `local`.

6. **Heredoc limit stringidan keyin probel/tab.** `_EOF_` qatorida ko'rinmas probel bo'lsa — "here-document delimited by end-of-file" warning va script oxirigacha hammasi heredoc ichida. Limit string qatori toza bo'lsin.

## Amaliy mashqlar

Muhit: `docker run -it --rm ubuntu:24.04 bash`

**1.** `hello` scripti yarating: ismingiz bilan salomlashsin va bugungi sanani chiqarsin. To'liq sikl: yozish → chmod → `./` bilan → `~/bin` orqali PATH dan ishga tushirish.

<details><summary>Yechim</summary>

```bash
mkdir -p ~/bin && export PATH="$HOME/bin:$PATH"
cat > ~/bin/hello <<"EOF"
#!/usr/bin/env bash
echo "Salom, $USER! Bugun: $(date +%F)"
EOF
chmod 755 ~/bin/hello
hello
```
</details>

**2.** Probel xatosini o'zingiz ko'ring: `x = 5` deb yozilgan script nima deydi? Nima uchun?

<details><summary>Yechim</summary>

```console
$ bash -c "x = 5"
bash: line 1: x: command not found
```
Bash `x` ni buyruq deb qidirdi. Assignment sintaksisi qat'iy: `x=5`.
</details>

**3.** Heredoc bilan `motd` scripti: hostname, uptime va disk holatini chiroyli "banner" ichida chiqarsin.

<details><summary>Yechim</summary>

```bash
#!/usr/bin/env bash
cat << _EOF_
=========================================
  Server : $(hostname)
  Uptime : $(uptime -p 2>/dev/null || uptime)
  Disk   : $(df -h / | tail -1 | awk '{print $5 " band"}')
=========================================
_EOF_
```
</details>

**4.** Quoted heredoc farqini isbotlang: bitta scriptda ikkala variantni chiqarib, `$HOME` qaysi birida expand bo'lishini ko'rsating.

<details><summary>Yechim</summary>

```bash
cat << EOF
expand: $HOME
EOF
cat << "EOF"
literal: $HOME
EOF
```
Natija: birinchisi `/root`, ikkinchisi `$HOME` matnini chiqaradi.
</details>

**5.** `sys_report` scripti: uch funksiya (`report_os`, `report_mem`, `report_disk`) yozing, har birida `local title` ishlating, oxirida uchchalasini chaqiring.

<details><summary>Yechim</summary>

```bash
#!/usr/bin/env bash
report_os()   { local title="OS";   echo "== $title =="; grep PRETTY /etc/os-release; }
report_mem()  { local title="RAM";  echo "== $title =="; free -h | grep Mem; }
report_disk() { local title="Disk"; echo "== $title =="; df -h / | tail -1; }
report_os
report_mem
report_disk
```
</details>

**6.** local siz nima bo'lishini ko'rsating: ikkala funksiyada ham `count` variablesini local siz ishlatib, bir-birini buzishini demonstratsiya qiling. Keyin `local` bilan tuzating.

<details><summary>Yechim</summary>

```bash
f1() { count=1; }
f2() { count=99; }
f1; f2
echo "$count"      # 99 — f2 f1 nikini yozib yubordi (global!)
# local bilan:
g1() { local count=1; echo "g1: $count"; }
g2() { local count=99; echo "g2: $count"; }
g1; g2; echo "tashqarida: ${count:-yo'q emas, 99 qoldi}"
```
</details>

**7.** (Qiyinroq) Top-down amaliyoti: "backup_home" scriptini avval 3 ta stub funksiya bilan yozing (`check_space`, `create_archive`, `verify_archive`), har qadamda ishga tushirib, keyin stublarni birma-bir real kod bilan almashtiring.

<details><summary>Yechim</summary>

```bash
#!/usr/bin/env bash
check_space()    { echo "[stub] joy tekshirildi"; }
create_archive() { echo "[stub] arxiv yaratildi"; }
verify_archive() { echo "[stub] tekshirildi"; }

main() { check_space; create_archive; verify_archive; }
main
```
Keyin birma-bir: `check_space` → `df -h /tmp | tail -1`; `create_archive` → `tar czf /tmp/home-$(date +%F).tgz ~/bin`; `verify_archive` → `tar tzf /tmp/home-*.tgz > /dev/null && echo OK`. Har almashtirishdan keyin ishga tushiring — script hech qachon "singan" holatda qolmaydi.
</details>

## Cheat sheet

| Element | Sintaksis | Eslatma |
|---------|-----------|---------|
| Shebang | `#!/usr/bin/env bash` | faylning 1-qatori |
| Komment | `# matn` | qator oxirida ham |
| Ishga tushirish huquqi | `chmod 755 script` | 700 — shaxsiy |
| Variable | `nom="qiymat"` | probelsiz `=`, ishlatishda `"$nom"` |
| Konstanta | `NOM="qiymat"` | konvensiya: KATTA harf |
| Command subst | `natija="$(cmd)"` | backticks emas |
| Heredoc | `cat << _EOF_ ... _EOF_` | `<< "_EOF_"` — expand o'chadi |
| Funksiya | `nom() { ...; }` | ta'rif chaqiruvdan oldin |
| Local | `local var="x"` | funksiya ichida HAR DOIM |
| Skelet | `set -euo pipefail` + `main "$@"` | zamonaviy standart |
| Joylashuv | `~/bin`, `/usr/local/bin` | PATH dagi kataloglar |

## Qo'shimcha manbalar

- [Bash Reference Manual — Shell Functions](https://www.gnu.org/software/bash/manual/html_node/Shell-Functions.html) — rasmiy hujjat
- [Google Shell Style Guide](https://google.github.io/styleguide/shellguide.html) — sanoat standarti uslub qo'llanmasi
- [Modern Bash Best Practices (Mechanical Rock)](https://www.mechanicalrock.io/blog/modern-bash) — zamonaviy skelet va patternlar

---

[← Oldingi: 17 — text-processing](17-text-processing.md) · [Kurs xaritasi](00-README.md) · [Keyingi: 19 — branching →](19-branching.md)
