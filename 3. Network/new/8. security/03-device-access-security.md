# 03. Device Access Security

Router va switch xavfsizligining birinchi bosqichi - boshqaruv kirishini himoyalash. Agar attacker qurilmaning CLI muhitiga kira olsa, ACL, VLAN, routing va boshqa himoyalar ham o'zgartirilishi mumkin.

## Console, VTY va privileged mode

Cisco qurilmada asosiy kirish joylari:

- **Console line** - jismoniy console kabel orqali.
- **VTY line** - SSH yoki Telnet orqali masofadan.
- **Privileged EXEC mode** - `enable` orqali kiriladigan yuqori huquqli rejim.

Minimal himoya:

```cisco
conf t
enable secret Str0ngEnableSecret
service password-encryption
line console 0
 password ConsolePass123
 login
 exec-timeout 5 0
 logging synchronous
line vty 0 4
 password VtyPass123
 login
 exec-timeout 5 0
 transport input ssh
end
```

`service password-encryption` kuchli kriptografiya emas, lekin running-config ichida oddiy parollarni ochiq ko'rinishda qoldirmaydi. `enable secret` esa `enable password`dan yaxshiroq.

## Local user bilan login

Yaxshiroq yondashuv - har bir administrator uchun alohida user.

```cisco
conf t
username admin privilege 15 secret AdminSecret123
username netops privilege 5 secret NetopsSecret123
line vty 0 4
 login local
 transport input ssh
end
```

Privilege level:

- `0` - juda cheklangan.
- `1` - user EXEC.
- `15` - full admin.

CCNA uchun privilege tushunchasi yetarli; katta tarmoqlarda AAA authorization ishlatiladi.

## SSH sozlash

SSH uchun hostname, domain name va RSA key kerak.

```cisco
conf t
hostname R1
ip domain-name example.local
crypto key generate rsa modulus 2048
ip ssh version 2
ip ssh time-out 60
ip ssh authentication-retries 3
username admin privilege 15 secret AdminSecret123
line vty 0 4
 login local
 transport input ssh
end
```

Tekshirish:

```cisco
show ip ssh
show ssh
show users
```

Telnetni o'chirish:

```cisco
conf t
line vty 0 4
 transport input ssh
end
```

## Management ACL

SSHni faqat admin subnetdan ruxsat berish:

```cisco
conf t
ip access-list standard SSH_MGMT
 permit 10.10.10.0 0.0.0.255
 deny any log
line vty 0 4
 access-class SSH_MGMT in
end
```

Bu ACL qurilmaga SSH qilishni cheklaydi. U router orqali o'tayotgan oddiy trafikni filtrlamaydi.

## Login protection

Brute-force hujumlarni kamaytirish:

```cisco
conf t
security passwords min-length 10
login block-for 120 attempts 3 within 60
login delay 2
login on-failure log
login on-success log
end
```

Ma'nosi:

- `min-length 10` - parol kamida 10 belgi.
- `block-for 120 attempts 3 within 60` - 60 soniyada 3 marta xato bo'lsa, 120 soniya bloklash.
- `login on-failure log` - xato loginlarni logga yozish.

## Banner

Legal ogohlantirish banneri:

```cisco
conf t
banner motd #
Unauthorized access is prohibited.
All activities may be monitored.
#
end
```

Bannerda qurilma modeli, joylashuvi, ichki IP yoki admin ismini yozmang. Bu attacker uchun reconnaissance ma'lumot bo'lishi mumkin.

## Keraksiz xizmatlarni o'chirish

```cisco
conf t
no ip http server
no ip http secure-server
no service tcp-small-servers
no service udp-small-servers
no ip source-route
end
```

CDP/LLDP masalasi:

```cisco
conf t
interface g0/1
 no cdp enable
 no lldp transmit
 no lldp receive
end
```

CDP/LLDP ichki tarmoqda foydali, lekin untrusted portlarda ortiqcha ma'lumot chiqarishi mumkin.

## NTP va logging

Incident tahlili uchun vaqt to'g'ri bo'lishi kerak.

```cisco
conf t
clock timezone UZT 5
ntp server 10.10.10.5
logging host 10.10.10.20
logging trap informational
service timestamps log datetime msec
end
```

Tekshirish:

```cisco
show clock
show ntp status
show logging
```

## Configuration backup

Running configni startup configga saqlash:

```cisco
copy running-config startup-config
```

TFTPga backup:

```cisco
copy running-config tftp:
```

Config backup ham xavfsizlikning bir qismi. Qurilma buzilsa yoki noto'g'ri sozlama kiritilsa, tez tiklash mumkin.

## Troubleshooting commands

```cisco
show running-config | section line
show running-config | include username|enable secret|transport input
show ip ssh
show users
show line
show login
show logging
debug ip ssh
```

`debug` buyruqlarini production tarmoqda ehtiyotkorlik bilan ishlating.

## Common mistakes

- `transport input telnet ssh` qoldirish.
- RSA key yaratmasdan SSH ishlatmoqchi bo'lish.
- `login local` yozmasdan username yaratish.
- Domain name sozlamasdan `crypto key generate rsa` bajarish.
- VTY ACLda admin IPni unutib, o'zini qurilmadan bloklab qo'yish.
- `copy run start` qilmasdan qurilmani reboot qilish.

## Q&A

**Savol:** `enable password` va `enable secret` farqi nima?

**Javob:** `enable secret` kuchliroq hash bilan saqlanadi va `enable password`dan ustun turadi. Amalda `enable secret` ishlating.

**Savol:** Console portga ham parol kerakmi?

**Javob:** Ha. Jismoniy kirish bo'lsa, console orqali sozlamani o'zgartirish mumkin.

**Savol:** SSH ishlamayapti, birinchi nimani tekshiraman?

**Javob:** `show ip ssh`, hostname/domain, RSA key, `line vty` ichida `login local` va `transport input ssh`, shuningdek management ACLni tekshiring.
