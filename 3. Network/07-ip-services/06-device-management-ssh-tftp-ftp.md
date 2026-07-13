# Device Management: SSH, TFTP va FTP (xavfsiz boshqaruv va backup)

## Muammo: Telnet parolni ko'chada ochiq tashiydi

Router'ni masofadan boshqarmoqchisan. Klassik usul - **Telnet**. Sen `admin`
va `StrongPass123` yozib kirasan. Lekin Telnet bu ma'lumotni **ochiq matnda**
(shifrsiz) yuboradi. Yo'lda kimdir tinglasa (masalan bir xil switch'ga ulangan
hujumchi), sening login va parolingni **to'g'ridan-to'g'ri o'qiydi**.

Ikkinchi muammo: qurilma konfiguratsiyasi qimmatli. Agar router yonib ketsa yoki
noto'g'ri buyruq bilan hammasi o'chib ketsa - sozlamalarni qayta yozish soatlab
vaqt oladi. Backup bo'lmasa, falokat.

Kerak: (1) **shifrlangan** masofaviy boshqaruv, (2) konfiguratsiyani muntazam
**backup** qilish usuli.

## Analogiya: pochta - ochiq otkritka va berk konvert

- **Telnet** - bu **ochiq otkritka**. Yo'ldagi har bir pochtachi matnni o'qiy
  oladi (login/parol ko'rinadi).
- **SSH** - bu **muhrlangan berk konvert**, ustiga qulf ilingan. Faqat qabul
  qiluvchi ocha oladi; yo'lda hech kim ichini ko'rmaydi.
- **TFTP/FTP backup** - bu muhim hujjatning **nusxasini seyfda** saqlash. Asl
  yo'qolsa, nusxadan tiklaysan.

**Analogiya chegarasi:** SSH faqat "berkitish" emas - u qarshi tomon haqiqatan
o'sha router ekanini ham tekshiradi (host key orqali), soxta routerga ulanmaslik
uchun.

## Sodda ta'rif

> **SSH** (Secure Shell) - qurilmani masofadan boshqarish uchun butun sessiyani
> (login, parol, buyruqlar) shifrlaydigan protocol. TCP porti **22**.

Telnet va SSH qisqa taqqoslash:

| | Telnet | SSH |
|---|---|---|
| Shifrlash | Yo'q (ochiq matn) | Bor (butun sessiya) |
| Port | TCP 23 | TCP 22 |
| Server tekshiruvi | Yo'q | Host key bilan bor |
| Ishlatish | Ishlatilmaydi | Standart |

## SSH qanday ishlaydi: kalit almashinuvi

SSH ulanganda ikki tomon avval **umumiy maxfiy kalitni** (shared secret)
kelishib oladi - lekin bu kalitni yo'lda ochiq yubormasdan. Bu sehr **key
exchange** (Diffie-Hellman kabi) orqali bo'ladi.

```mermaid
sequenceDiagram
    participant C as Client (sen)
    participant S as SSH Server (router)
    C->>S: 1. Ulanish, versiya kelishuvi
    S->>C: 2. Host key (server o'zini tanishtiradi)
    Note over C: Host key oldingisiga mos keladimi?
    C->>S: 3. Key exchange (DH) - umumiy sirni yasash
    S->>C: 3. Key exchange (DH)
    Note over C,S: Endi ikkalasida bir xil session key bor
    C->>S: 4. Login/parol (shifrlangan!)
    S->>C: 5. Sessiya ochildi (barcha buyruqlar shifrlangan)
```

Ikki muhim tushuncha:

- **Host key** - server o'zini tanishtiruvchi noyob kalit. Client uni eslab
  qoladi. Keyingi ulanishda host key o'zgargan bo'lsa, client ogohlantiradi
  ("balki soxta serverga ulanmoqdasan"). Bu **man-in-the-middle** hujumidan
  himoya.
- **Session key** - key exchange orqali yasalgan, faqat shu sessiyaga xos
  shifrlash kaliti. Butun suhbat shu bilan shifrlanadi.

## SSH uchun asosiy talablar

Router'da SSH ishlashi uchun 5 shart bajarilishi kerak:

```mermaid
flowchart TD
    A["1. Hostname default emas<br/>(Router emas, R1)"] --> B["2. ip domain-name bor"]
    B --> C["3. RSA/crypto key yaratilgan"]
    C --> D["4. Local user bor"]
    D --> E["5. VTY line SSH qabul qiladi"]
    E --> OK["SSH ishlaydi"]
```

Nega domain name kerak? Chunki RSA kalit nomi `hostname.domain` dan yasaladi
(masalan `R1.example.local`). Domain bo'lmasa, kalit yaratilmaydi.

## Worked example: SSH konfiguratsiya

```cisco
conf t
! --- 1-qadam: hostname va domain (kalit uchun asos) ---
hostname R1
ip domain-name example.local
! --- 2-qadam: local admin user (privilege 15 = to'liq huquq) ---
username admin privilege 15 secret StrongPass123
! --- 3-qadam: RSA kalit (2048 bit) va SSH v2 majburiy ---
crypto key generate rsa modulus 2048
ip ssh version 2
! --- 4-qadam: VTY faqat SSH qabul qilsin, Telnet yopiq ---
line vty 0 4
 login local
 transport input ssh
 exec-timeout 10 0
end
```

Muhim: `transport input ssh` - Telnet'ni **butunlay yopadi**. `exec-timeout 10 0`
- 10 daqiqa harakatsizlikda sessiyani uzadi (xavfsizlik).

VTY'ni faqat admin subnetiga cheklash (ACL):

```cisco
conf t
ip access-list standard MGMT-ONLY
 permit 192.168.100.0 0.0.0.255
 deny any
line vty 0 4
 access-class MGMT-ONLY in
end
```

### Tekshirish buyruqlari

```cisco
show ip ssh
show ssh
show running-config | section line vty
show users
```

Client'dan ulanish:

```bash
ssh admin@192.168.100.1
```

## Predict: nima bo'ladi?

> 🤔 **O'ylab ko'r:** Sen hostname, user, `ip ssh version 2` yozding va VTY'da
> `transport input ssh` qo'yding, lekin `crypto key generate rsa` va
> `ip domain-name` ni unutding. SSH ulanishga urinsang nima bo'ladi?

<details>
<summary>💡 Javobni ko'rish</summary>

SSH ishlamaydi. RSA kalit yo'q - shuning uchun shifrlash uchun asos yo'q. Kalit
esa `hostname.domain-name` dan yasaladi, `ip domain-name` bo'lmasa `crypto key
generate rsa` ham xato beradi. Client ulana olmaydi ("connection refused"). Va
Telnet ham yopiq (`transport input ssh`), demak umuman kira olmaysan - faqat
console orqali tuzatasan.
</details>

## TFTP backup: oddiy, lekin himoyasiz

**TFTP** (Trivial File Transfer Protocol) - konfiguratsiya faylini uzatishning
eng oddiy usuli. UDP 69 porti. Autentifikatsiyasiz - shuning uchun faqat ishonchli
management tarmoqda ishlatiladi.

Running-config'ni TFTP serverga saqlash:

```cisco
copy running-config tftp:
```

IOS quyidagilarni so'raydi:

```text
Address or name of remote host []? 192.168.100.70
Destination filename [r1-confg]? R1-2026-07-10.cfg
```

Startup-config backup va tiklash:

```cisco
copy startup-config tftp:
copy tftp: running-config
```

> ⚠️ `copy tftp: running-config` mavjud running-config **ustiga merge** qiladi
> (qo'shadi), to'liq almashtirmaydi. To'liq tiklash uchun boshqacha yondashuv
> kerak.

## FTP backup: autentifikatsiya bilan

**FTP** username/parol bilan ishlaydi (TFTP'dan bu jihati yaxshiroq), lekin
klassik FTP ham trafikni **shifrlamaydi**.

```cisco
conf t
ip ftp username backupuser
ip ftp password BackupPass123
end
```

Backup:

```cisco
copy running-config ftp:
```

URL ko'rinishida ham:

```cisco
copy running-config ftp://backupuser:BackupPass123@192.168.100.80/R1-config.cfg
```

## Running-config va startup-config: muhim farq

Bu farqni chalkashtirmaslik backup uchun hayotiy:

```mermaid
flowchart LR
    RAM["running-config<br/>(RAM - hozir ishlayotgan)"] -->|"copy run start"| NVRAM["startup-config<br/>(NVRAM - reload'da yuklanadi)"]
    NVRAM -->|"reload"| RAM
```

| | running-config | startup-config |
|---|---|---|
| Qayerda | RAM | NVRAM |
| Nima | Hozir ishlayotgan | Reload'da yuklanadigan |
| O'chganda | Reload'da yo'qoladi | Saqlanib qoladi |

`copy running-config startup-config` (yoki `wr`) - o'zgarishlarni saqlash. Buni
qilmasang, reload'dan keyin barcha sozlamalar yo'qoladi.

## IOS image backup

Falokatga tayyorgarlik uchun IOS image'ni ham saqlash mumkin:

```cisco
dir flash:
copy flash:c2900-universalk9-mz.SPA.157-3.M.bin tftp:
```

## Zamonaviy SSH best practices (2025-2026)

SSH asoslar bir xil qolgan, lekin xavfsizlik standartlari o'zgardi:

- **Ed25519 kalitlar** - zamonaviy tavsiya. 256-bitli Ed25519 kalit 4096-bitli
  RSA'day xavfsiz, lekin ancha kichik va tez. RSA hali hamma joyda qo'llab-quvvatlanadi,
  lekin yangi tizimlar uchun Ed25519 afzal.
- **Zaif algoritmlarni o'chirish** - eski cipher'lar (3DES, CBC), SHA-1 MAC va
  zaif key exchange'larni o'chirib, faqat zamonaviylarini qoldirish kerak.
- **Post-quantum tayyorlik** - OpenSSH 9.0+ hybrid key exchange
  (`sntrup761x25519`) qo'shdi, 2025-da OpenSSH 10.0 `mlkem768x25519` ni default
  yoqdi - kelajakdagi kvant kompyuterlariga qarshi.
- **Cisco qurilmalari uchun:** faqat SSHv2, kuchli RSA (2048-4096 bit) yoki EC
  kalit, VTY'da source-IP ACL, va zaif cipher'larni cheklash.
- **Kalitlarni muntazam almashtirish** (rotation) - kalit sizib chiqsa, zarar
  vaqtini cheklaydi.

```cisco
! Zamonaviy tavsiya: EC kalit va SSHv2 majburiy
conf t
crypto key generate ec keysize 384 label SSH-KEYS
ip ssh version 2
ip ssh dh min size 2048
end
```

TFTP/FTP shifrsizligini hisobga olib, zamonaviy muhitda backup uchun **SCP**
(SSH ustidan fayl ko'chirish) yoki **SFTP** ishlatish tavsiya etiladi - ular
shifrlangan.

## Xavfsizlik tavsiyalari (checklist)

- `enable secret` ishlat, `enable password` emas (secret hash'lanadi).
- Local user uchun `secret`, `password` emas.
- SSH version 2 majburiy.
- VTY'ga ACL qo'y (faqat management subnet).
- Keraksiz servislarni o'chir:

```cisco
conf t
no ip http server
no ip http secure-server
service password-encryption
end
```

> `service password-encryption` kuchli himoya emas - u parolni faqat oddiy
> ko'rinishda (yengil obfuscation) yashiradi, oson ochiladi. Haqiqiy himoya
> `secret` (kuchli hash) da.

## Ko'p uchraydigan xatolar

⚠️ **Xato 1: ip domain-name bermaslik.**
RSA kalit yaratilmaydi, SSH ishlamaydi. Domain kalit nomining bir qismi.

⚠️ **Xato 2: RSA key yaratmaslik.**
`crypto key generate rsa` bo'lmasa SSH shifrlay olmaydi.

⚠️ **Xato 3: Telnet'ni ochiq qoldirish.**
`transport input ssh telnet` yoki default `transport input all` qolsa, Telnet
hali ochiq - zaiflik. Faqat `transport input ssh`.

⚠️ **Xato 4: login local yoki user unutish.**
`login local` bor-u user yo'q bo'lsa (yoki aksincha) - kira olmaysan.

⚠️ **Xato 5: copy tftp merge ekanini unutish.**
`copy tftp: running-config` to'liq almashtirmaydi, mavjudga qo'shadi. Kutilmagan
natija berishi mumkin.

⚠️ **Xato 6: TFTP UDP 69 firewall'da yopiq.**
Serverga ping bor, lekin backup ishlamaydi - firewall UDP 69'ni bloklagan.

## Xulosa

- SSH masofaviy boshqaruvni shifrlaydi (TCP 22); Telnet (TCP 23) ochiq matn -
  ishlatilmaydi.
- SSH host key bilan serverni tekshiradi, key exchange bilan session key yasaydi.
- SSH 5 shart talab qiladi: hostname, domain-name, RSA key, local user, VTY SSH.
- TFTP (UDP 69) - oddiy, autentifikatsiyasiz; FTP - user/parol, lekin shifrsiz.
- `copy tftp: running-config` merge qiladi, almashtirmaydi.
- running-config (RAM, o'chuvchan) vs startup-config (NVRAM, saqlanadi).
- Zamonaviy: Ed25519 kalitlar, zaif cipher'larni o'chirish, SCP/SFTP backup.

## 🧠 Eslab qol

- SSH = TCP 22, Telnet = TCP 23 (ishlatilmaydi).
- SSH shartlari: hostname + domain + RSA key + user + VTY.
- Host key = serverni tekshirish (MITM'dan himoya).
- `copy run start` qilmasang, reload'da sozlama yo'qoladi.
- Zamonaviy kalit: Ed25519; zamonaviy backup: SCP/SFTP.

## ✅ O'z-o'zini tekshir (retrieval practice)

**1.** Nega SSH uchun `ip domain-name` shart?

<details>
<summary>Javob</summary>

RSA/crypto kalit nomi `hostname.domain-name` shaklida yasaladi (masalan
R1.example.local). Domain bo'lmasa `crypto key generate rsa` kalit yarata olmaydi,
SSH esa shu kalitsiz shifrlay olmaydi.
</details>

**2.** Host key nima uchun kerak va u qanday hujumdan himoya qiladi?

<details>
<summary>Javob</summary>

Host key server (router) o'zini tanishtiruvchi noyob kalit. Client uni eslab
qoladi; keyingi ulanishda o'zgargan bo'lsa ogohlantiradi. Bu man-in-the-middle
hujumidan himoya - soxta serverga (o'zini asl router qilib ko'rsatgan) ulanib
qolmaslik uchun.
</details>

**3.** `copy tftp: running-config` bilan to'liq eski konfiguratsiyani tiklab
bo'ladimi?

<details>
<summary>Javob</summary>

Yo'q. Bu buyruq mavjud running-config **ustiga merge** qiladi (qo'shadi),
to'liq almashtirmaydi. Ba'zi eski sozlamalar qolib, kutilmagan aralashma
bo'lishi mumkin. To'liq tiklash uchun boshqa yondashuv (masalan `configure
replace`) kerak.
</details>

**4.** VTY'da `transport input ssh telnet` qolsa qanday xavf bor?

<details>
<summary>Javob</summary>

Telnet hali ochiq. Kimdir Telnet orqali ulansa, login/parol ochiq matnda ketadi
va tinglagan hujumchi ko'radi. Faqat `transport input ssh` qo'yib, Telnet'ni
butunlay yopish kerak.
</details>

**5.** Zamonaviy muhitda TFTP/FTP o'rniga nima tavsiya etiladi va nega?

<details>
<summary>Javob</summary>

SCP yoki SFTP - ular SSH ustidan ishlaydi va fayl uzatishni shifrlaydi. TFTP
autentifikatsiyasiz va shifrsiz, klassik FTP esa parol so'rasa ham shifrlamaydi -
backup fayl (konfiguratsiya) yo'lda ochiq ketadi, bu maxfiy ma'lumotni fosh
qiladi.
</details>

## 🛠 Amaliyot

**1. Oson (Modify).** Yuqoridagi SSH konfiguratsiyasida `exec-timeout` ni 5
daqiqaga qisqartir va ikkinchi read-only user (`monitor`) qo'sh.

<details>
<summary>Hint</summary>

`exec-timeout 5 0`; `username monitor privilege 1 secret <parol>` (privilege 1
= cheklangan).
</details>

**2. O'rta (faded example).** SSH skeletini to'ldir:

```cisco
conf t
hostname CORE-SW
! TODO: domain-name yoz
! TODO: admin user yarat (privilege 15, secret bilan)
! TODO: RSA kalit yarat (2048 bit)
ip ssh version 2
line vty 0 4
 ! TODO: local user bilan login
 ! TODO: faqat SSH qabul qilsin
end
```

<details>
<summary>Hint</summary>

`ip domain-name example.local`; `username admin privilege 15 secret ...`;
`crypto key generate rsa modulus 2048`; `login local`; `transport input ssh`.
</details>

**3. Qiyin (Make).** Noldan to'liq xavfsiz management yoz: SSHv2 (EC yoki 2048
RSA), faqat 10.0.0.0/24 subnetdan VTY ruxsat, `enable secret`, keraksiz HTTP
servislar o'chirilgan, exec-timeout 10 daqiqa. Keyin running-config'ni TFTP
serverga (10.0.0.70) backup qilish buyrug'ini yoz.

<details>
<summary>Hint</summary>

Kerak: hostname, domain-name, crypto key, `ip ssh version 2`, `enable secret`,
`username`, VTY'da `login local` + `transport input ssh` + `access-class`,
standard ACL 10.0.0.0/24, `no ip http server`, `no ip http secure-server`.
Backup: `copy running-config tftp:` -> 10.0.0.70.
</details>

## 🔁 Takrorlash

- Bog'liq mavzular: NTP (SSH sertifikat va vaqt uchun); Syslog (SSH login
  urinishlari loglanadi - `login on-failure log`); ACL asoslari (VTY access-class
  uchun).
- Takrorlash jadvali:
  - **Ertaga:** SSH 5 shartini xotiradan ayt.
  - **3 kundan keyin:** host key va session key farqini tushuntir.
  - **1 haftadan keyin:** running-config vs startup-config va `copy tftp` merge
    xususiyatini ayt.
- **Feynman testi:** SSH'ni buyruq ishlatmasdan 3 jumlada tushuntir: nega Telnet
  o'rniga SSH, host key nima uchun kerak, va backup nega muhim?

## 📚 Manbalar

- SSH key best practices 2025 (Ed25519, rotation): https://www.brandonchecketts.com/archives/ssh-ed25519-key-best-practices-for-2025
- SSH and SCP in 2026 - hardening: https://dev.to/matheus_releaserun/ssh-and-scp-in-2026-configuration-security-hardening-and-advanced-tips-8nm
- Cisco IOS XE device hardening guide: https://www.thenetworkdna.com/2026/03/cisco-ios-xe-device-hardening-complete.html
- Comparing SSH keys (RSA, ECDSA, EdDSA): https://goteleport.com/blog/comparing-ssh-keys/
