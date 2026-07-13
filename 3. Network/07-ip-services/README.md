# 07 - IP Services

Bu modul tarmoqning kundalik ishlashini ta'minlaydigan amaliy **IP xizmatlarini**
o'rgatadi: qurilmalarga IP tarqatish, vaqtni sinxronlash, holatni kuzatish,
loglarni yig'ish, muhim trafikni ustuvor qilish va qurilmalarni xavfsiz boshqarish.

Maqsad - buyruqlarni yodlash emas, balki **qaysi xizmat nima uchun kerakligini**,
qayerda sozlanishini va nosozlikda qaysi `show`/`debug` buyrug'i yordam berishini
tushunish. Har dars CCNA uslubidagi Cisco CLI misollari bilan boradi va ustiga
zamonaviy (2025-2026) best practice qatlamini qo'shadi.

## Nima o'rganiladi

- Qurilmalarga IP sozlamalarni avtomatik berish (DHCP DORA, relay).
- Butun tarmoq soatini bitta ishonchli manbaga moslash (NTP, stratum).
- Qurilmalar holatini markazdan kuzatish (SNMP, MIB/OID, trap).
- Loglarni markaziy serverga yig'ish (Syslog, severity 0-7).
- Cheklangan bandwidth'ni muhim trafikka ustuvorlik berib taqsimlash (QoS).
- Qurilmani shifrlangan boshqarish va konfiguratsiyani backup qilish (SSH/TFTP/FTP).

## Darslar

| # | Dars | Qisqa mazmun |
|---|---|---|
| 1 | [01-dhcp.md](01-dhcp.md) | DHCP DORA jarayoni, lease, `ip helper-address` relay, DHCP snooping |
| 2 | [02-ntp.md](02-ntp.md) | Vaqt sinxronizatsiyasi, stratum, konfiguratsiya, NTS xavfsizlik |
| 3 | [03-snmp.md](03-snmp.md) | SNMP versiyalari, MIB/OID, polling vs trap, telemetry/gNMI |
| 4 | [04-syslog.md](04-syslog.md) | Severity 0-7, facility, markaziy logging, zamonaviy pipeline |
| 5 | [05-qos.md](05-qos.md) | Classification, DSCP marking, queuing, policing vs shaping |
| 6 | [06-device-management-ssh-tftp-ftp.md](06-device-management-ssh-tftp-ftp.md) | SSH key exchange/host key, TFTP/FTP backup, best practices |

## O'qish tartibi

Darslar ketma-ket o'qish uchun tuzilgan va bir-biriga bog'lanadi:

1. **DHCP** - qurilma tarmoqqa ulanganda birinchi bo'ladigan jarayon.
2. **NTP** - keyingi barcha xizmatlar (log, sertifikat, audit) uchun to'g'ri
   vaqt asos.
3. **SNMP** va **Syslog** - monitoring va logging birga ishlaydi; ikkalasi
   NTP vaqtiga tayanadi.
4. **QoS** - trafik boshqaruvi, mustaqil o'qilishi ham mumkin.
5. **Device Management (SSH/TFTP/FTP)** - qurilmani xavfsiz boshqarish va backup;
   SSH sertifikat/vaqt uchun NTP bilan bog'liq.

## Amaliy maslahatlar

- Har xizmatda avval **reachability** tekshir: `ping`, `traceroute`, `show ip route`.
- ACL ishlatilsa, DHCP relay, SNMP yoki QoS muammosi aslida ACL match
  bo'lmayotganidan bo'lishi mumkin.
- Vaqt xato bo'lsa syslog va troubleshooting chalkashadi - **NTP'ni erta sozla**.
- Boshqaruv uchun Telnet o'rniga har doim **SSH**.
- Backupdan oldin `show running-config` va `show startup-config` farqini bil.

## Tezkor tekshiruv buyruqlari

```cisco
show ip interface brief
show ip route
show ntp status
show snmp
show logging
show policy-map interface
show ip ssh
```
