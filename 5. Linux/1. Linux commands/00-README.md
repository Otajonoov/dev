# Linux Basic — "The Linux Command Line" asosidagi kurs

> Manba: William E. Shotts, "The Linux Command Line" (2-nashr) + har bir mavzu bo'yicha zamonaviy web manbalar.
> Auditoriya: 3+ yillik tajribali backend developer (Go, Docker, PostgreSQL, CI/CD).
> Barcha buyruqlar **Ubuntu 24.04 (bash 5.2)** muhitida real ishga tushirilib tekshirilgan.

## Kurs xaritasi

### Blok A — Shell fundamentlari

| # | Dars | Tavsif |
|---|------|--------|
| 01 | [shell-and-terminal](01-shell-and-terminal.md) | Shell nima, terminal/tty, readline editing va history tricks |
| 02 | [navigation-and-filesystem](02-navigation-and-filesystem.md) | FHS daraxti, `cd`/`ls`/`file`/`less`, symlink va hardlink |
| 03 | [file-operations](03-file-operations.md) | Wildcards, `mkdir`/`cp`/`mv`/`rm`/`ln` |
| 04 | [commands-and-documentation](04-commands-and-documentation.md) | `type`/`which`/`man`/`apropos`/`alias` — mustaqil o'rganish skili |
| 05 | [redirection-and-pipelines](05-redirection-and-pipelines.md) | stdin/stdout/stderr, `>` `>>` `2>` `\|`, filters, `tee` |
| 06 | [expansion-and-quoting](06-expansion-and-quoting.md) | Brace/parameter/command substitution, quoting qoidalari |
| 07 | [permissions](07-permissions.md) | rwx, `chmod`/`umask`/`chown`, `su`/`sudo`, SUID/SGID |
| 08 | [processes](08-processes.md) | `ps`/`top`, jobs, signals, `kill` |

### Blok B — Environment va konfiguratsiya

| # | Dars | Tavsif |
|---|------|--------|
| 09 | [environment](09-environment.md) | Env variables, `.bashrc`/`.profile` zanjiri, PS1 prompt |
| 10 | [vim-basics](10-vim-basics.md) | Modes, editing, search/replace |

### Blok C — Kundalik asboblar

| # | Dars | Tavsif |
|---|------|--------|
| 11 | [package-management](11-package-management.md) | apt/dpkg, dnf/rpm, repository va dependency mexanikasi |
| 12 | [storage-and-filesystems](12-storage-and-filesystems.md) | `mount`, `df`/`du`, `lsblk`, fdisk/mkfs/fsck |
| 13 | [networking](13-networking.md) | `ping`/`ip`/`ss`, `ssh`/`scp`, `wget`/`curl` |
| 14 | [finding-files](14-finding-files.md) | `find`, `locate`, `xargs` |
| 15 | [archiving-and-sync](15-archiving-and-sync.md) | `tar`/`gzip`/`zip`, `rsync` |
| 16 | [regular-expressions](16-regular-expressions.md) | BRE/ERE, `grep` chuqur |
| 17 | [text-processing](17-text-processing.md) | `sort`/`uniq`/`cut`/`join`/`diff`/`tr`/`sed` + `printf` |

### Blok D — Shell scripting

| # | Dars | Tavsif |
|---|------|--------|
| 18 | [scripting-first-steps](18-scripting-first-steps.md) | Shebang, PATH, variables, here docs, functions |
| 19 | [branching](19-branching.md) | `if`, exit codes, `test`/`[[ ]]`/`(( ))`, `case` |
| 20 | [loops](20-loops.md) | `while`/`until`/`for`, fayl o'qish looplari |
| 21 | [script-input](21-script-input.md) | `read`, IFS, positional params, `shift`, options parsing |
| 22 | [strings-numbers-arrays](22-strings-numbers-arrays.md) | Parameter expansion tricks, arithmetic, arrays |
| 23 | [advanced-scripting](23-advanced-scripting.md) | Subshells, process substitution, `trap`, `wait`, named pipes |
| 24 | [debugging-and-best-practices](24-debugging-and-best-practices.md) | `set -x`, defensive scripting, shellcheck, `set -euo pipefail` |

## Progress

- [x] 01 — shell-and-terminal
- [x] 02 — navigation-and-filesystem
- [x] 03 — file-operations
- [x] 04 — commands-and-documentation
- [x] 05 — redirection-and-pipelines
- [x] 06 — expansion-and-quoting
- [x] 07 — permissions
- [x] 08 — processes
- [x] 09 — environment
- [x] 10 — vim-basics
- [x] 11 — package-management
- [x] 12 — storage-and-filesystems
- [x] 13 — networking
- [x] 14 — finding-files
- [x] 15 — archiving-and-sync
- [x] 16 — regular-expressions
- [x] 17 — text-processing
- [x] 18 — scripting-first-steps
- [x] 19 — branching
- [x] 20 — loops
- [x] 21 — script-input
- [x] 22 — strings-numbers-arrays
- [x] 23 — advanced-scripting
- [x] 24 — debugging-and-best-practices

## Qanday o'rganish kerak

1. **Tartib bilan boring** — har bir dars oldingilariga tayanadi (learning path shunday qurilgan).
2. **Terminalda takrorlang** — har bir code block copy-paste qilib ishlatiladigan holatda. O'qish ≠ o'rganish; buyruqni o'zingiz terib ko'ring.
3. **Docker da mashq qiling** — toza muhit uchun: `docker run -it --rm ubuntu:24.04 bash` (kursdagi barcha misollar shu muhitda tekshirilgan).
4. **Amaliy mashqlarni yeching** — yechimga qaramasdan avval o'zingiz urinib ko'ring, har darsda 5-7 ta mashq bor.
5. **Cheat sheet ni alohida saqlang** — har dars oxiridagi jadval kundalik ishda tez qarash uchun.

## Kurs holati

**Kurs to'liq yakunlangan: 24/24 dars.** Har dars TLCL (2-nashr) tegishli boblari + har mavzu bo'yicha 1-3 web-qidiruv (best practices, keng tarqalgan xatolar, zamonaviy muqobillar) sintezidan yozilgan. Har bir code block va uning chiqishi Ubuntu 24.04 (bash 5.2) konteynerida real ishga tushirilib tekshirilgan; dnf/rpm misollari Fedora 41 da, mount/mkfs misollari privileged konteynerda verify qilingan. Kitobda bor, lekin eskirgani uchun kursga kirmagan mavzular: printing (22-bob), groff/pr (21-bob), C kompilyatsiya (23-bob).
