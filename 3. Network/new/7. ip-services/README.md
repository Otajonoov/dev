# 7. IP Services

Bu bo'lim CCNA darajasida tarmoqdagi amaliy IP xizmatlarini tushuntiradi. Maqsad: buyruqlarni yodlash emas, balki qaysi xizmat nima uchun kerakligini, qayerda sozlanishini va nosozlik paytida qaysi `show` yoki `debug` buyrug'i yordam berishini tushunish.

## Mavzular

| Fayl | Mavzu | Qisqa mazmun |
|---|---|---|
| [01-dhcp-relay.md](01-dhcp-relay.md) | DHCP Relay | Boshqa subnetdagi DHCP serverga `ip helper-address` orqali so'rov uzatish |
| [02-nat-cisco.md](02-nat-cisco.md) | NAT | Static NAT, Dynamic NAT, PAT, port forwarding va tekshiruv |
| [03-ntp.md](03-ntp.md) | NTP | Router/switch vaqtini sinxronlash, autentifikatsiya asoslari |
| [04-snmp.md](04-snmp.md) | SNMP | SNMP v2c va v3, monitoring, trap, xavfsizlik |
| [05-syslog.md](05-syslog.md) | Syslog | Log severity darajalari, remote syslog serverga yuborish |
| [06-qos.md](06-qos.md) | QoS | Classification, marking, queuing, policing, shaping |
| [07-ssh-tftp-ftp.md](07-ssh-tftp-ftp.md) | SSH/TFTP/FTP | Xavfsiz boshqaruv, konfiguratsiya backup/restore |

## CCNA uchun asosiy g'oya

IP services odatda trafikni "uzatish"dan ko'ra tarmoqni boshqarish, kuzatish va barqaror ishlatish uchun kerak:

- DHCP Relay: client va DHCP server boshqa VLAN/subnetda bo'lsa.
- NAT/PAT: private IP manzillarni internetga chiqarish yoki ichki serverni tashqaridan ochish.
- NTP: loglar, sertifikatlar, autentifikatsiya va troubleshooting uchun to'g'ri vaqt.
- SNMP: qurilma holatini monitoring qilish.
- Syslog: loglarni markaziy serverda yig'ish.
- QoS: muhim trafikni, masalan voice yoki video trafikni, ustuvor qilish.
- SSH/TFTP/FTP: qurilmani xavfsiz boshqarish va konfiguratsiyalarni saqlash.

## Tezkor tekshiruv buyruqlari

```cisco
show ip interface brief
show running-config
show ip route
show access-lists
show logging
show ntp status
show snmp
show ip nat translations
show policy-map interface
```

## Amaliy maslahatlar

- Har bir xizmatda avval reachability tekshiring: `ping`, `traceroute`, `show ip route`.
- ACL ishlatilsa, NAT yoki QoS muammosi aslida ACL match bo'lmayotganidan kelib chiqishi mumkin.
- Vaqt xato bo'lsa, syslog va troubleshooting chalkashadi; NTPni erta sozlash foydali.
- Management uchun Telnet o'rniga SSH ishlating.
- Backupdan oldin `show running-config` va `show startup-config` farqini biling.

