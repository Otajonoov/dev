# 01. Security Concepts

Tarmoq xavfsizligi - tarmoq resurslarini ruxsatsiz kirish, buzish, o'g'irlash, xizmatni to'xtatish va noto'g'ri foydalanishdan himoya qilishdir. CCNA darajasida asosiy maqsad atamalarni tushunish, xavfni baholash va oddiy himoya choralarini qo'llashdir.

## CIA triad

Xavfsizlikning klassik uchligi:

- **Confidentiality** - ma'lumotni faqat ruxsat etilgan odam yoki tizim ko'radi. Misol: SSH, WPA2/WPA3, IPsec encryption.
- **Integrity** - ma'lumot yo'lda o'zgartirilmagan bo'lishi kerak. Misol: hashing, HMAC, routing authentication.
- **Availability** - xizmat kerak paytda ishlashi kerak. Misol: redundancy, DoS mitigation, monitoring.

## Threat, vulnerability, exploit, mitigation

| Atama | Ma'nosi | Oddiy misol |
|---|---|---|
| Threat | Xavf manbai yoki zarar yetkazuvchi holat | Internetdan kelgan attacker |
| Vulnerability | Zaiflik | Telnet yoqilgan, parol oddiy |
| Exploit | Zaiflikdan foydalanish usuli | Telnet parolini sniffing qilish |
| Mitigation | Xavfni kamaytirish | Telnet o'rniga SSH, ACL, kuchli parol |

Misol:

```text
Vulnerability: switch management VLAN hamma portlardan ochiq
Threat: ichki tarmoqdagi ruxsatsiz foydalanuvchi
Exploit: SSH brute-force yoki default parol bilan kirish
Mitigation: management ACL, AAA, kuchli parol, login block-for
```

## Attack surface

**Attack surface** - hujum qilinishi mumkin bo'lgan barcha kirish nuqtalari. Router yoki switch uchun misollar:

- VTY line orqali SSH/Telnet.
- Console port.
- SNMP.
- Web management.
- CDP/LLDP orqali oshkor bo'lgan ma'lumot.
- Ochiq access portlar.
- Trunk portlar va native VLAN.
- DHCP, ARP, STP kabi L2 protokollar.

Attack surface kamaytirish:

```cisco
conf t
no ip http server
no ip http secure-server
service password-encryption
no cdp run
end
```

Eslatma: `no cdp run` hamma interfeyslarda CDP ni o'chiradi. Agar faqat bitta interfeysda o'chirmoqchi bo'lsangiz:

```cisco
conf t
interface g0/1
 no cdp enable
end
```

## Common network attacks

### Reconnaissance

Attacker tarmoq haqida ma'lumot yig'adi: IP range, ochiq portlar, qurilma turi, OS versiya.

Mitigation:

- Keraksiz xizmatlarni o'chirish.
- Management access ni ACL bilan cheklash.
- Bannerda ortiqcha ma'lumot bermaslik.
- SNMP community string ni kuchli qilish yoki SNMPv3 ishlatish.

### Password attack

Parolni taxmin qilish, brute-force, dictionary attack.

Mitigation:

```cisco
conf t
security passwords min-length 10
login block-for 120 attempts 3 within 60
enable secret VeryStrongSecret123
username admin privilege 15 secret AnotherStrongSecret123
end
```

### Man-in-the-middle

Attacker ikki tomon orasiga kirib trafikni o'qishi yoki o'zgartirishi mumkin. LAN ichida ARP spoofing bunga misol bo'ladi.

Mitigation:

- Dynamic ARP Inspection.
- DHCP snooping.
- SSH va HTTPS.
- IPsec VPN.

### DoS/DDoS

Xizmatni ishlamay qoladigan darajada trafik yoki so'rov bilan band qilish.

Mitigation:

- Rate limiting va policing.
- ACL orqali keraksiz trafikni bloklash.
- Control Plane Policing.
- Monitoring va upstream provider bilan himoya.

## Security policy

Texnik sozlamalar siyosatsiz tartibsiz bo'lib qoladi. Oddiy security policy quyidagilarni belgilaydi:

- Kim qaysi qurilmaga kira oladi?
- Qaysi protokollar ruxsat etilgan?
- Parol uzunligi va almashtirish talabi qanday?
- Loglar qayerga yuboriladi?
- Guest Wi-Fi ichki tarmoqqa kira oladimi?
- Incident bo'lsa kim javob beradi?

## Defense in depth

**Defense in depth** - bitta himoyaga ishonib qolmaslik. Masalan, faqat parol emas:

```text
SSH + AAA + ACL + logging + least privilege + backup config
```

Shunda bitta qatlam xato bo'lsa, qolgan qatlamlar xavfni kamaytiradi.

## Cisco'da foydali tekshiruv komandalar

```cisco
show running-config
show version
show users
show line
show logging
show ip ssh
show access-lists
show interfaces status
show mac address-table
```

## Common mistakes

- Telnetni yoqib qo'yish va parolni ochiq matnda yuborish.
- `enable password` ishlatish, lekin `enable secret` ishlatmaslik.
- Bitta umumiy userni hamma administratorlar ishlatishi.
- Log va NTP sozlanmagani sabab incident vaqtini aniqlay olmaslik.
- Keraksiz xizmatlarni o'chirmaslik.
- Switch access portlarini ochiq qoldirish.

## Q&A

**Savol:** Xavfsizlikda eng birinchi nimani sozlash kerak?

**Javob:** Qurilma boshqaruv kirishini himoyalang: `enable secret`, local user, SSH, VTY ACL, Telnetni o'chirish.

**Savol:** Vulnerability va exploit farqi nima?

**Javob:** Vulnerability - zaif joy. Exploit - shu zaif joydan amalda foydalanish usuli.

**Savol:** Mitigation xavfni butunlay yo'q qiladimi?

**Javob:** Har doim emas. Ko'pincha xavfni qabul qilinadigan darajagacha kamaytiradi.
