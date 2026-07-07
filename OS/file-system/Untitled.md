**File system** (fayl tizimi yoki FS) - bu kompyuter fayllarini qanday tashkil etish va ularga qanday kirish mumkinligini boshqaradigan tizim. **Local file system** operatsion tizimning bir qismi bo'lib, bir xil kompyuterda ishlaydigan dasturlarga xizmat ko'rsatadi.

**Distributed file system** esa tarmoq orqali bir-biriga ulangan kompyuterlar o'rtasida fayllarga kirish imkonini beruvchi protokoldir.

File system data storage xizmatini taqdim etadi va dasturlarga katta hajmdagi xotira (mass storage) bilan ishlash imkonini beradi. Agar file system bo'lmasa, dasturlar xotira bilan mos kelmaydigan usullarda ishlar edi, bu resurslar to'qnashuviga, ma'lumotlarning buzilishiga va yo'qolishiga olib kelardi.

File systemlar turli xil storage qurilmalari uchun ishlab chiqilgan: **hard disk drives (HDD)**, **solid-state drives (SSD)**, magnit lentalar va optik disklar.

## Tarix

1900-yillardan va kompyuterlar paydo bo'lishidan oldin "file system" va "filing system" atamalari qog'oz hujjatlarni tashkil etish, saqlash va topish usullarini tasvirlash uchun ishlatilgan. 1961-yilda bu atama kompyuterlashtirilgan fayl saqlashga qo'llanila boshlandi va 1964-yilga kelib umumiy foydalanishga kirdi.

## Arxitektura

Local file system arxitekturasini abstraction layerlar (qatlamlar) sifatida tasvirlash mumkin:

### 1. Logical file system layer (Mantiqiy fayl tizimi qatlami)

Bu qatlam application programming interface (API) orqali yuqori darajali kirish imkonini beradi. File operatsiyalari: open, close, read, write. Bu qatlam open file table entries va process uchun file descriptor'larni boshqaradi. File access, directory operations, security va protection'ni ta'minlaydi.

### 2. Virtual file system (VFS)

Bu ixtiyoriy qatlam bo'lib, bir vaqtning o'zida bir nechta physical file system'larni qo'llab-quvvatlaydi. Har bir physical file system - bu file system implementation.

### 3. Physical file system layer (Jismoniy fayl tizimi qatlami)

Bu qatlam storage device'ga past darajali kirish imkonini beradi. U data block'larni o'qiydi va yozadi, buffering va boshqa memory management'ni ta'minlaydi va block'larni storage media'da aniq joylarga joylashtiradi. Bu qatlam device driver'lar yoki channel I/O orqali storage device bilan ishlaydi.

## File System Atributlari

### File nomlari (Filenames)

**File name** - bu faylni dasturlar va ba'zi hollarda foydalanuvchilar uchun identifikatsiya qiladi. File name yagona (unique) bo'lishi kerak, shunda dastur aniq bir faylga murojaat qila oladi.

Aksariyat file systemlar file name uzunligini cheklaydi. Ba'zi file systemlar **case sensitive** (katta-kichik harflarni farqlaydi), boshqalari **case insensitive** (farqlamaydi). Masalan, "MYFILE" va "myfile" case insensitive tizimda bir xil fayl, case sensitive'da esa turli fayllar.

Zamonaviy file systemlar file name'da Unicode dan keng ko'lamli belgilarni qo'llab-quvvatlaydi.

### Directory'lar (Papkalar)

File systemlar odatda fayllarni **directory'lar** (yoki **folder'lar**) ga tashkil etishni qo'llab-quvvatlaydi. Bu fayllarni guruhlarga ajratadi.

Directory strukturasi **flat** (chiziqli) yoki **hierarchical** (ierarxik) bo'lishi mumkin - ya'ni directory ichida directory bo'lishi mumkin (subdirectory).

Birinchi bo'lib ixtiyoriy ierarxik directory'larni qo'llab-quvvatlovchi file system **Multics** operatsion tizimida ishlatilgan. Unix-like tizimlar, HFS+, FAT, NTFS ham ierarxik directory'larni qo'llab-quvvatlaydi.

### Metadata

File content (ma'lumotlar)dan tashqari, file system quyidagi **metadata**'ni ham saqlaydi:

- **name** - fayl nomi
- **size** - o'lcham (block'lar soni yoki byte'lar)
- **created, last accessed, last modified** - vaqt belgilari
- **owner** - user va group
- **access permissions** - kirish ruxsatlari
- **file attributes** - read-only, executable va boshqalar
- **device type** - block, character, socket, subdirectory

File system metadata'ni faylning content'idan alohida saqlaydi.

Aksariyat file systemlar bir directory'dagi barcha fayllar nomlarini **directory table**'da saqlaydi. Ko'p file systemlar faqat ba'zi metadata'ni directory table'ga qo'yadi, qolganini esa butunlay alohida strukturada (masalan **inode**'da) saqlaydi.

Ba'zi file systemlar (NTFS, XFS, ext2/3/4, HFS+) **extended file attributes** yordamida qo'shimcha atributlarni qo'llab-quvvatlaydi.

## Storage Space Organization (Xotira bo'yicha tashkilot)

### Ajratish (Allocation)

Local file system xotira qaysi qismlarining qaysi faylga tegishli ekanligini va qaysi qismlar bo'sh ekanligini kuzatib boradi.

File system fayl yaratganda, u data uchun joy ajratadi. Ba'zi file systemlar boshlang'ich joy hajmini va keyin fayl o'sishi bilan qo'shimcha hajmlarni belgilashga ruxsat beradi yoki talab qiladi.

Faylni o'chirish uchun file system fayl joyini bo'sh deb belgilaydi - boshqa fayl uchun ishlatish mumkin.

### Granular ajratish va Slack space

Local file system odatda xotira joyini **granular** (dona-dona) tarzda ajratadi, ko'pincha bir nechta physical unit'lar (masalan byte'lar).

Masalan, 1980-yillarning boshida Apple DOS'da 140 KB floppy disk'da 256-byte sector'lar ishlatilgan.

Granular ajratish natijasida har bir fayl uchun (granular hajmga ko'paytma bo'lgan noyob hajmdagi fayllardan tashqari) **slack space** (ishlatilmagan joy) qoladi. 512-byte allocation uchun o'rtacha ishlatilmagan joy 256 byte. 64 KB cluster'lar uchun o'rtacha 32 KB.

Odatda allocation unit size xotira sozlanganda belgilanadi. Nisbatan kichik hajm tanlash ortiqcha overhead'ga olib keladi. Katta hajm esa ortiqcha ishlatilmagan joyga olib keladi.

### Fragmentation (Parchalanish)

File system fayllarni yaratganda, o'zgartirganda va o'chirganda, underlying storage representation **fragmented** (parchalangan) bo'lishi mumkin. Fayllar va fayllar orasidagi bo'sh joy ketma-ket bo'lmagan allocation block'larni egallaydi.

Agar faylni saqlash uchun zarur bo'lgan joy ketma-ket block'larda ajratilmasa, fayl fragmented bo'ladi. Fayllar o'chirilganda bo'sh joy ham fragmented bo'ladi.

Fragmentation oxirgi foydalanuvchi uchun ko'rinmaydi va tizim to'g'ri ishlashda davom etadi. Lekin bu **hard disk drive**'lar kabi ketma-ket block'lar bilan yaxshi ishlaydigan ba'zi storage hardware'larda performance'ni pasaytirishi mumkin. **Solid-state drive**'lar fragmentation'dan ta'sirlanmaydi.

## Access Control (Kirish Nazorati)

File system ko'pincha o'zi boshqaradigan data'ga **access control**'ni qo'llab-quvvatlaydi.

Access control'ning maqsadi ko'pincha ba'zi userlarni ba'zi fayllarni o'qish yoki o'zgartirishdan to'xtatishdir.

Access control shuningdek, dasturlar tomonidan kirishni cheklashi mumkin, bu ma'lumotlarning nazorat ostida o'zgartirilishini ta'minlaydi. Misollar: faylning metadata'sida yoki boshqa joyda saqlangan parollar, **permission bit**'lar, **access control list**'lar yoki **capability**'lar shaklida file permission'lar.

File data'ni shifrlash usullari ba'zan file system'ga kiritilgan. Bu juda samarali, chunki file system utility'lari data'ni samarali boshqarish uchun encryption seed'ni bilishiga hojat yo'q.

### Storage Quota

Ba'zi operatsion tizimlar system administrator'ga **disk quota**'ni yoqish imkonini beradi, bu foydalanuvchining storage joy foydalanishini cheklaydi.

## Data Integrity (Ma'lumotlar Yaxlitligi)

File system odatda saqlangan ma'lumotlarning ham oddiy operatsiyalarda, ham istisno holatlarda **consistent** (izchil) bo'lib qolishini ta'minlaydi:

- Dastur faylga kirishni tugatganini (faylni yopishni) xabar qilmaslik
- Dastur to'satdan to'xtashi (crash)
- Media nosozligi
- Remote tizimlarga ulanishni yo'qotish
- Operatsion tizim nosozligi
- Tizimni qayta yuklash (soft reboot)
- Elektr uzilishi (hard reboot)

Istisno holatlarda tiklanish metadata'ni, directory entry'larni yangilashni va bufferlangan lekin storage media'ga yozilmagan data bilan ishlashni o'z ichiga olishi mumkin.

## Data Access (Ma'lumotlarga Kirish)

### Byte Stream Access

Ko'p file systemlar data'ga **stream of bytes** (byte oqimi) sifatida kiradi. Odatda file data'sini o'qish uchun dastur memory buffer'ni taqdim etadi va file system media'dan data'ni oladi, keyin buffer'ga yozadi. Write operatsiyasi dastur byte'lar buffer'ini taqdim etishi, file system esa o'qib, media'ga saqlashini o'z ichiga oladi.

### Record Access

Ba'zi file systemlar yoki file system ustidagi layer'lar dasturga **record** (yozuv) ni aniqlash imkonini beradi, shunda dastur data'ni strukturaviy tarzda o'qishi va yozishi mumkin - tartibsiz byte ketma-ketligi emas.

## File System Types (Turlari)

### Disk File Systems

Disk file system disk storage media'ning data'ga random kirish qobiliyatidan foydalanadi. Qo'shimcha fikrlar: dastlab so'ralganidan keyin data'ga kirish tezligi va keyingi data ham so'ralishi mumkinligi haqidagi taxmin.

Misollar: FAT (FAT12, FAT16, FAT32), exFAT, NTFS, ReFS, HFS+, HPFS, APFS, UFS, ext2/3/4, XFS, btrfs, ZFS, ReiserFS.

### Optical Discs

**ISO 9660** va **Universal Disk Format (UDF)** - Compact Disc'lar, DVD'lar va Blu-ray disc'lar uchun ikkita keng tarqalgan format.

### Flash File Systems

**Flash file system** flash memory device'larning maxsus qobiliyatlari, performance'i va cheklovlarini hisobga oladi. Ko'pincha disk file system flash memory device'ni underlying storage media sifatida ishlatishi mumkin, lekin flash device uchun maxsus mo'ljallangan file system'dan foydalanish ancha yaxshi.

### Tape File Systems

**Tape file system** - lentaga fayllarni saqlash uchun mo'ljallangan file system va tape format. Magnit lentalar sequential storage media bo'lib, random data access vaqti disklarga qaraganda ancha uzoqroq.

Disk file system'da odatda **master file directory** va ishlatilgan va bo'sh data region'larning xaritasi mavjud. Har qanday fayl qo'shish, o'zgartirish yoki o'chirish directory va ishlatilgan/bo'sh xaritalarni yangilashni talab qiladi.

Lenta esa linear motion'ni (chiziqli harakat) talab qiladi. Lentaning bir uchidan ikkinchi uchiga o'tish bir necha soniyadan bir necha daqiqagacha vaqt olishi mumkin.

### Database File Systems

File management uchun yana bir konsepsiya - **database-based file system**. Ierarxik boshqaruvdan tashqari yoki qo'shimcha ravishda, fayllar o'z xususiyatlari bilan identifikatsiya qilinadi: file turi, mavzu, muallif yoki shunga o'xshash boy metadata.

### Network File Systems

**Network file system** - remote file access protokoli uchun client sifatida ishlaydigan file system. Local interfacelardan foydalanuvchi dasturlar remote tarmoqqa ulangan kompyuterlardagi ierarxik directory'lar va fayllarni shaffof tarzda yaratishi, boshqarishi va ularga kirishi mumkin.

Misollar: NFS, AFS, SMB protokollari uchun clientlar, FTP va WebDAV uchun file-system-like clientlar.

### Shared Disk File Systems

**Shared disk file system** - bir nechta mashinalar (odatda serverlar) bir xil external disk subsystem'ga (odatda **storage area network**) kirish huquqiga ega. File system shu subsystem'ga kirishni arbitration qiladi, write collision'larni oldini oladi.

Misollar: GFS2 (Red Hat), GPFS/Spectrum Scale (IBM), CXFS (SGI), StorNext (Quantum).

### Flat File Systems

**Flat file system**'da subdirectory'lar yo'q - barcha fayllar uchun directory entry'lar bitta directory'da saqlanadi.

Floppy disk media birinchi paydo bo'lganda, bu turdagi file system nisbatan kichik data hajmi tufayli yetarli edi. CP/M mashinalari flat file system'ni qo'llab-quvvatlagan.

Oddiy bo'lsa-da, flat file systemlar fayllar soni ko'paygan sari noqulay bo'lib qoladi va data'ni tegishli fayllar guruhlariga tashkil etishni qiyinlashtiradi.

Yaqinda flat file system oilasiga Amazon S3 qo'shildi - bu remote storage service. Faqat **bucket**'lar (cheksiz o'lchamdagi disk kabi) va **object**'lar (fayl konsepsiyasiga o'xshash) mavjud.

## Operating System Implementations

### Unix va Unix-like OS'lar

Unix-like operatsion tizimlar **virtual file system** yaratadi, bu barcha ulangan storage device'lardagi barcha fayllarni bitta ierarxiyada ko'rsatadi. Ya'ni bitta **root directory** mavjud va tizimdagi barcha fayllar uning ostida joylashgan.

Unix-like tizimlar har bir device'ga device name beradi, lekin bu faylga kirish usuli emas. Boshqa device'dagi faylarga kirish uchun, operatsion tizimga directory tree'da qayerda paydo bo'lishi kerakligini aytish kerak. Bu jarayon **mounting** deb ataladi.

Masalan, CD-ROM'dagi faylarga kirish uchun operatsion tizimga "Bu CD-ROM'dan file system'ni ol va `/media` directory ostida ko'rsat" deyish kerak. `/media` directory ko'p Unix sistemlarda mavjud bo'lib, maxsus removable media (CD, DVD, USB, floppy) uchun mount point sifatida mo'ljallangan.

### Linux

Linux ko'plab file systemlarni qo'llab-quvvatlaydi, lekin system disk uchun keng tarqalgan tanlovlar: **ext*** oilasi (ext2, ext3, ext4), **XFS**, **JFS** va **btrfs**.

Raw flash uchun (flash translation layer yoki MTD'siz): **UBIFS**, **JFFS2**, **YAFFS**.

**SquashFS** - keng tarqalgan compressed read-only file system.

### macOS

macOS **Apple File System (APFS)**'dan foydalanadi, bu 2017-yilda classic Mac OS'dan meros bo'lib qolgan **HFS Plus (HFS+)** ni almashtirdi.

HFS Plus - metadata-rich va case-preserving lekin odatda case-insensitive file system. macOS'ning Unix ildizlari tufayli HFS Plus'ga Unix permission'lar qo'shilgan. Keyingi versiyalarda **journaling** qo'shildi.

File name'lar 255 gacha belgi bo'lishi mumkin. HFS Plus file name'larni saqlash uchun Unicode'dan foydalanadi.

macOS FAT file systemlarni (16 va 32) o'qish va yozishni qo'llab-quvvatlaydi. NTFS'ni o'qishi mumkin, lekin yozish uchun qo'shimcha sozlama yoki third-party software kerak.

### Microsoft Windows

Windows **FAT**, **NTFS**, **exFAT**, **Live File System** va **ReFS** file systemlaridan foydalanadi.

Windows user darajasida bir diskni yoki partition'ni boshqasidan ajratish uchun **drive letter** abstraction'dan foydalanadi. Masalan, `C:\WINDOWS` yo'li C harfi bilan ifodalangan partition'dagi WINDOWS directory'ni bildiradi.

#### FAT

FAT file system'lar oilasi deyarli barcha personal kompyuter operatsion tizimlari tomonidan qo'llab-quvvatlanadi: barcha Windows va MS-DOS/PC DOS versiyalari, OS/2, DR-DOS.

Yillar davomida file system FAT12'dan FAT16 va FAT32'ga kengaytirilgan. Subdirectory'lar, codepage qo'llab-quvvatlash, extended attribute'lar va long filename'lar qo'shilgan.

FAT12 va FAT16 root directory'dagi entry'lar sonida va FAT-formatted disk'lar yoki partition'lar maksimal hajmida cheklovlarga ega edi.

FAT32 bu cheklovlarni hal qiladi (4 GB fayl hajmi cheklovi bundan mustasno), lekin NTFS bilan solishtirganda cheklangan.

FAT12, FAT16 va FAT32 file name uchun 8 belgi va extension uchun 3 belgi (**8.3 filename limit**) chekloviga ega.

#### NTFS

NTFS Windows NT operatsion tizimi bilan 1993-yilda taqdim etildi va **ACL-based permission control**'ni qo'llab-quvvatladi. Boshqa xususiyatlar: hard link'lar, multiple file stream'lar, attribute indexing, quota tracking, sparse file'lar, encryption, compression, **reparse point**'lar.

#### exFAT

exFAT file system overhead jihatidan NTFS ustidan ma'lum afzalliklarga ega. exFAT FAT file systemlar (FAT12, FAT16, FAT32) bilan backward compatible emas.

macOS va Windows'da to'liq qo'llab-quvvatlanadigan va 4 GB dan katta fayllarni saqlashi mumkin bo'lgan yagona file system - bu exFAT.

## Design Limitations (Dizayn Cheklovlari)

File systemlar saqlanishi mumkin bo'lgan data hajmini cheklaydi - bu odatda file system loyihalashtirilgan paytdagi storage device'larning tipik hajmi va yaqin kelajakda kutilgan hajm bilan bog'liq.

Storage hajmlari deyarli eksponensial tezlikda oshgani uchun (Moore qonuni), yangi storage device'lar ko'pincha mavjud file system cheklovlarini taqdimotdan keyin bir necha yil ichida oshib ketadi. Bu tobora ortib borayotgan hajmga ega yangi file systemlarni talab qiladi.

1980-yillarning boshlaridagi 50 KB dan 512 KB gacha xotiraga ega bo'lgan uy kompyuterlari file system'lari yuzlab gigabyte hajmga ega zamonaviy storage tizimlar uchun oqilona tanlov bo'lmaydi. Xuddi shunday, zamonaviy file systemlar o'sha eski tizimlar uchun oqilona tanlov bo'lmaydi.

## Xulosa

File system - bu zamonaviy computing'ning fundamental qismi bo'lib, data'ni tashkil etish, saqlash va unга kirish uchun tizimli yondashuv taqdim etadi. 1960-yillardan buyon file systemlar oddiy flat strukturalardan murakkab distributed va object-based tizimlarga rivojlandi. Har bir file system turi - disk, tape, flash, network - o'zining maxsus use case'lari va optimizatsiyalari bilan mo'ljallangan.

---

**Eslatma:** Ushbu tarjimada texnik atamalar asl ko'rinishida qoldirilgan (masalan: file system, directory, metadata, buffer, cache, etc.) chunki ular IT sohasida standart terminologiya hisoblanadi va tarjima qilish ularn

i tushunarliligini kamaytirishi mumkin.