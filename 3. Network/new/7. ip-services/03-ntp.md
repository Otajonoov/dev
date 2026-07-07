# NTP: Network Time Protocol

NTP qurilmalar vaqtini sinxronlash uchun ishlatiladi. To'g'ri vaqt troubleshooting, syslog, sertifikat, SSH, SNMP va audit uchun juda muhim.

## NTP nima uchun kerak?

Routerda log chiqdi:

```text
Mar 1 00:00:12.123: %LINK-3-UPDOWN: Interface Gi0/1, changed state to up
```

Agar vaqt noto'g'ri bo'lsa, bu log qachon sodir bo'lganini aniqlash qiyin. NTP barcha qurilmalarni bir vaqt manbasiga bog'laydi.

## Asosiy konfiguratsiya

```cisco
conf t
clock timezone UZT 5 0
ntp server 192.168.100.20
end
```

Manba interfeysini belgilash:

```cisco
conf t
ntp source loopback0
end
```

Bu ayniqsa routing o'zgarishlarida barqarorlik beradi.

## NTP tekshiruv

```cisco
show clock detail
show ntp status
show ntp associations
show running-config | include ntp
```

`show ntp status`da `synchronized` ko'rinsa, qurilma NTP bilan vaqtni olgan.

## NTP autentifikatsiya

NTP autentifikatsiya noto'g'ri vaqt manbasidan himoya qiladi. CCNA darajasida asosiy g'oya: client va server bir xil key ishlatadi.

```cisco
conf t
ntp authenticate
ntp authentication-key 1 md5 CCNA_NTP_KEY
ntp trusted-key 1
ntp server 192.168.100.20 key 1
end
```

## Router NTP server sifatida

Lab muhitida bitta router vaqt manbasi sifatida ishlashi mumkin:

```cisco
conf t
clock set 10:30:00 21 May 2026
ntp master 5
end
```

`ntp master` real productionda ehtiyotkorlik bilan ishlatiladi. Agar tashqi ishonchli vaqt manbasi bo'lmasa, lab yoki kichik yopiq tarmoqda foydali.

## Syslog bilan bog'lash

Loglarda aniq vaqt chiqishi uchun:

```cisco
conf t
service timestamps log datetime msec localtime show-timezone
service timestamps debug datetime msec localtime show-timezone
end
```

## Muammolar va yechimlar

| Muammo | Sabab | Tekshiruv |
|---|---|---|
| NTP sync bo'lmayapti | NTP serverga reachability yo'q | `ping`, `show ip route` |
| Vaqt bir necha soat farq qiladi | Timezone noto'g'ri | `show clock detail` |
| `show ntp associations`da server bor, sync yo'q | Server ishonchsiz yoki stratum yuqori | Server holati |
| Auth ishlamayapti | Key ID yoki password mos emas | NTP configni solishtirish |
| ACL bloklayapti | UDP 123 yopiq | ACL/firewall tekshirish |

## Keng tarqalgan xatolar

- NTP server IPga ping borligini tekshirmasdan NTPni ayblash.
- Timezone va NTPni aralashtirish. NTP UTC vaqtni sinxronlaydi, timezone esa lokal ko'rsatishni belgilaydi.
- `ntp source` bergandan keyin server tomonida shu source IPga qaytish yo'li yo'qligi.
- Labda `clock set` qilingan, lekin reloaddan keyin vaqt yana xato bo'lishini unutish.

## Q&A

**Savol:** NTP qaysi portdan foydalanadi?  
**Javob:** UDP 123.

**Savol:** Switchlarda ham NTP kerakmi?  
**Javob:** Ha. Switch loglari ham troubleshootingda kerak bo'ladi.

**Savol:** NTP bo'lmasa tarmoq ishlamay qoladimi?  
**Javob:** Oddiy routing ishlashi mumkin, lekin log, sertifikat, autentifikatsiya va audit muammoli bo'ladi.

