# Syslog: loglarni markaziy yig'ish

Syslog qurilma hodisalarini yozib boradi: interface up/down, routing neighbor o'zgarishi, login urinishlari, konfiguratsiya o'zgarishi va boshqa voqealar. Markaziy syslog server troubleshooting va audit uchun juda foydali.

## Syslog xabar tuzilishi

Cisco log namunasi:

```text
May 21 10:12:33.456 UZT: %LINK-3-UPDOWN: Interface GigabitEthernet0/1, changed state to up
```

Bu yerda:

- `May 21 10:12:33.456 UZT`: vaqt.
- `%LINK`: facility yoki log manbai.
- `3`: severity.
- `UPDOWN`: hodisa turi.
- Qolgan qismi: xabar matni.

## Severity darajalari

| Raqam | Nomi | Ma'nosi |
|---|---|---|
| 0 | Emergency | Qurilma ishlamay qolgan |
| 1 | Alert | Zudlik bilan e'tibor kerak |
| 2 | Critical | Kritik xatolik |
| 3 | Error | Xatolik |
| 4 | Warning | Ogohlantirish |
| 5 | Notification | Muhim odatiy hodisa |
| 6 | Informational | Ma'lumot |
| 7 | Debug | Juda batafsil debug |

Kichik raqam muhimroq. Agar `logging trap 5` berilsa, serverga 0 dan 5 gacha bo'lgan loglar yuboriladi.

## Asosiy konfiguratsiya

```cisco
conf t
service timestamps log datetime msec localtime show-timezone
logging host 192.168.100.60
logging trap informational
logging source-interface loopback0
end
```

Console loglarni kamaytirish:

```cisco
conf t
no logging console
logging buffered 16384 informational
end
```

`no logging console` ayniqsa labda yoki productionda konsolni spamdan saqlaydi.

## VTY login loglari

Login urinishlarini ko'rish uchun:

```cisco
conf t
login on-failure log
login on-success log
end
```

## Tekshiruv buyruqlari

```cisco
show logging
show running-config | include logging
show clock detail
show ntp status
ping 192.168.100.60
```

Server tomonda UDP 514 ochiq bo'lishi kerak.

## Troubleshooting tartibi

1. Router syslog server IPga ping qila oladimi?
2. `logging host` to'g'ri yozilganmi?
3. `logging trap` juda cheklab qo'yilmaganmi?
4. Source interface IPga serverdan qaytish yo'li bormi?
5. Server UDP 514ni eshityaptimi?
6. Vaqt NTP bilan to'g'rimi?

## Keng tarqalgan xatolar

- NTP sozlanmaganligi sabab log vaqtlari noto'g'ri chiqishi.
- `logging trap errors` berib, informational hodisalarni kutish.
- Syslog server firewalli UDP 514ni bloklashi.
- `logging source-interface loopback0` berilgan, lekin serverga loopback IPdan qaytish route yo'qligi.
- Debug loglarni doimiy yoqib qo'yish.

## Q&A

**Savol:** Syslog qaysi portdan foydalanadi?  
**Javob:** Odatda UDP 514.

**Savol:** `logging buffered` va `logging host` farqi nima?  
**Javob:** `logging buffered` loglarni qurilma xotirasida saqlaydi. `logging host` loglarni tashqi syslog serverga yuboradi.

**Savol:** `debug` loglarini doimiy yoqish mumkinmi?  
**Javob:** Tavsiya qilinmaydi. Debug CPU va xotiraga yuk berishi mumkin.

