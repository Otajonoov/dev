# 06. Wireless Security

Wireless tarmoqlarda xavfsizlik ayniqsa muhim, chunki signal devordan ham tashqariga chiqishi mumkin. Attacker kabelga ulanmasdan ham SSIDni ko'rishi, autentifikatsiyani sinashi yoki noto'g'ri sozlangan guest tarmoqdan foydalanishi mumkin.

## Asosiy tushunchalar

- **SSID** - Wi-Fi tarmoq nomi.
- **BSSID** - access point radio MAC manzili.
- **WPA/WPA2/WPA3** - Wi-Fi xavfsizlik standartlari.
- **PSK** - umumiy parol, uy va kichik ofislarda ko'p ishlatiladi.
- **Enterprise** - har bir user alohida autentifikatsiya qilinadi, odatda 802.1X va RADIUS bilan.
- **Open network** - parolsiz tarmoq, trafik himoyasiz bo'lishi mumkin.

## WPA2 va WPA3

### WPA2-Personal

WPA2-Personal PSK ishlatadi. Hamma foydalanuvchilar bitta parolni biladi.

Afzallik:

- Sodda.
- Kichik tarmoqlar uchun qulay.

Kamchilik:

- Parol tarqalsa, hamma qurilmada almashtirish kerak.
- Zaif parol dictionary attackga moyil.

Tavsiya:

```text
Kamida 14-16 belgi, oddiy so'z emas, kompaniya nomi yoki telefon raqam emas.
```

### WPA2-Enterprise

802.1X va RADIUS orqali har bir user yoki qurilma alohida tekshiriladi.

```text
Client -> AP/WLC -> RADIUS server -> authentication result
```

Afzallik:

- Har bir user alohida.
- User ketganda faqat uning accounti o'chiriladi.
- Accounting va policy qo'llash osonroq.

### WPA3

WPA3 yangi xavfsizlik imkoniyatlarini beradi:

- **SAE** - PSK o'rniga kuchliroq handshake.
- Offline dictionary attackga yaxshiroq qarshilik.
- WPA3-Enterprise kuchliroq security variantlarini qo'llashi mumkin.

Real hayotda WPA2/WPA3 mixed mode uchraydi, chunki eski qurilmalar WPA3ni qo'llamasligi mumkin.

## Guest Wi-Fi

Guest tarmoq ichki corporate LANdan ajratilgan bo'lishi kerak.

Yaxshi amaliyot:

- Guest SSID alohida VLANda.
- Guest VLANdan internal serverlarga ACL bilan deny.
- Internetga NAT orqali chiqadi.
- Captive portal yoki vaqtinchalik access ishlatilishi mumkin.
- Guest tarmoqdan management IPga kirish taqiqlanadi.

ACL misoli:

```cisco
conf t
ip access-list extended GUEST_FILTER
 deny ip 192.168.50.0 0.0.0.255 10.0.0.0 0.255.255.255
 deny ip 192.168.50.0 0.0.0.255 172.16.0.0 0.15.255.255
 deny ip 192.168.50.0 0.0.0.255 192.168.0.0 0.0.255.255
 permit ip 192.168.50.0 0.0.0.255 any
interface vlan 50
 ip access-group GUEST_FILTER in
end
```

Bu RFC1918 private tarmoqlariga kirishni bloklab, qolgan trafikni ruxsat qiladi.

## Cisco WLC tekshiruvlari

WLC platformasiga qarab buyruqlar farq qilishi mumkin, lekin CCNA darajasida quyidagi tekshiruvlar foydali:

```cisco
show wlan summary
show client summary
show ap summary
show client detail <client-mac>
show advanced 802.11a summary
show advanced 802.11b summary
```

Catalyst 9800 WLC IOS-XE uslubida:

```cisco
show wireless summary
show wireless client summary
show wireless client mac-address aaaa.bbbb.cccc detail
show wlan summary
show ap summary
```

## AP va switch port

Lightweight AP odatda access yoki trunk portga ulanadi. Agar AP bir nechta SSID/VLAN tashisa, trunk kerak bo'lishi mumkin.

Misol:

```cisco
conf t
interface g0/10
 description AP-01
 switchport mode trunk
 switchport trunk native vlan 10
 switchport trunk allowed vlan 10,20,50
 spanning-tree portfast trunk
end
```

Oddiy bitta VLAN AP:

```cisco
conf t
interface g0/10
 description AP-01
 switchport mode access
 switchport access vlan 20
 spanning-tree portfast
end
```

## Wireless threatlar

### Rogue AP

Ruxsatsiz access point tarmoqqa ulanadi.

Mitigation:

- Switch port security.
- WLC rogue detection.
- 802.1X wired access.
- Unused portlarni shutdown qilish.

### Evil twin

Attacker haqiqiy SSIDga o'xshash soxta SSID yaratadi.

Mitigation:

- WPA2/WPA3-Enterprise.
- Sertifikat tekshiruvi.
- Foydalanuvchilarga noma'lum sertifikatni qabul qilmaslikni o'rgatish.

### Weak PSK

Oddiy parol dictionary attack bilan topilishi mumkin.

Mitigation:

- Kuchli PSK.
- PSKni davriy almashtirish.
- Enterprise authenticationga o'tish.

## Common mistakes

- Guest Wi-Fi ni internal VLANga ulash.
- WPA2-Personal parolini hamma joyda bir xil ishlatish.
- Eski WPA yoki WEPni yoqib qo'yish.
- SSIDni yashirishni xavfsizlik deb o'ylash.
- RADIUS sertifikat tekshiruvini e'tiborsiz qoldirish.
- AP switch portida kerakli VLANlarni trunk allowed listga qo'shmaslik.

## Q&A

**Savol:** SSIDni hide qilish xavfsizlik beradimi?

**Javob:** Juda kam. SSID baribir wireless trafikdan aniqlanishi mumkin. Asosiy himoya WPA2/WPA3 va to'g'ri autentifikatsiya.

**Savol:** WPA3 har doim ishlatilsinmi?

**Javob:** Agar clientlar qo'llasa, ha. Eski clientlar bo'lsa mixed mode kerak bo'lishi mumkin.

**Savol:** Guest tarmoqda parol bo'lishi yetarlimi?

**Javob:** Yo'q. Guest VLAN internal tarmoqlardan ACL/firewall bilan ajratilishi kerak.
