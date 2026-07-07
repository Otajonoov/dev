# Wireless LAN va WLAN

WLAN (Wireless LAN) - foydalanuvchilarni kabelsiz tarmoqqa ulash usuli. CCNA darajasida asosiy tushunchalar: SSID, BSSID, access point, controller, radio band, channel, encryption, authentication va VLAN mapping.

## Asosiy terminlar

- AP (Access Point) - wireless mijozlarni wired tarmoqqa ulaydi.
- SSID - foydalanuvchi ko'radigan WLAN nomi.
- BSSID - AP radio interfeysining MAC manzili.
- ESSID - bir nechta AP bir xil SSID bilan ishlaganda umumiy nom.
- WLC (Wireless LAN Controller) - AP va WLAN siyosatlarini markaziy boshqaradi.
- Lightweight AP - controller orqali boshqariladigan AP.
- Autonomous AP - mustaqil sozlanadigan AP.
- CAPWAP - lightweight AP va WLC orasidagi tunnel/control protokol.

## 2.4 GHz va 5 GHz

2.4 GHz:

- Kengroq qamrov.
- Ko'proq interference.
- Amaliyotda faqat 1, 6, 11 non-overlapping channel sifatida ishlatiladi.

5 GHz:

- Ko'proq kanal.
- Kamroq interference.
- Yuqori throughput uchun yaxshiroq.
- Qamrov 2.4 GHzga qaraganda qisqaroq bo'lishi mumkin.

6 GHz Wi-Fi 6E/7 muhitlarda ishlatiladi, lekin CCNA asosiy e'tibor 2.4 va 5 GHzga qaratadi.

## SSID va VLAN mapping

Har SSID odatda ma'lum VLANga bog'lanadi:

```text
SSID CORP  -> VLAN 10 -> 192.168.10.0/24
SSID GUEST -> VLAN 40 -> 192.168.40.0/24
SSID VOICE -> VLAN 30 -> 192.168.30.0/24
```

AP switchga trunk port orqali ulansa, bir nechta SSID VLANlari APga yetib boradi.

Switch port misoli:

```cisco
interface gigabitEthernet0/10
 description AP-LOBBY
 switchport mode trunk
 switchport trunk native vlan 99
 switchport trunk allowed vlan 10,30,40,99
 spanning-tree portfast trunk
 no shutdown
```

Ba'zi APlar management uchun native VLANdan, client trafik uchun tagged VLANlardan foydalanadi. Ba'zi dizaynlarda AP management ham tagged bo'lishi mumkin. Muhimi: AP/WLC dizayni va switch trunk mos bo'lishi kerak.

## Xavfsizlik

Tavsiya qilingan variantlar:

- WPA2-Personal - kichik tarmoqlar uchun PSK.
- WPA2/WPA3-Enterprise - korporativ tarmoq uchun 802.1X/RADIUS.
- WPA3 - yangi va kuchliroq xavfsizlik.

Yomon yoki eski variantlar:

- Open SSID - faqat guest/captive portal muhitida ehtiyotkorlik bilan.
- WEP - ishlatmaslik kerak.
- WPA/TKIP - eski, tavsiya qilinmaydi.

Enterprise autentifikatsiya:

```text
Client -> AP/WLC -> RADIUS server
```

RADIUS foydalanuvchi loginini tekshiradi va kerak bo'lsa VLAN yoki policy qaytarishi mumkin.

## Roaming

Roaming - client bir APdan boshqasiga o'tishi. Yaxshi roaming uchun:

- Bir xil SSID.
- Mos security policy.
- Overlap qamrov yetarli.
- Kanal rejalashtirish to'g'ri.
- Juda baland transmit powerdan qochish.

Yomon dizayn belgilari:

- Client uzoq APga yopishib oladi.
- Pingda uzilishlar ko'p.
- Voice/video chaqiriqlarda uzilish.
- Bir joyda 2.4 GHz juda shovqinli.

## Controller asosiy oqimi

Lightweight AP odatda quyidagicha ishlaydi:

1. IP oladi, ko'pincha DHCP orqali.
2. WLCni topadi: DHCP option, DNS, broadcast yoki qo'lda sozlama.
3. CAPWAP tunnel o'rnatadi.
4. WLCdan konfiguratsiya oladi.
5. SSIDlarni e'lon qiladi.

## Cisco CLI misollar

Switchda AP portini tekshirish:

```cisco
show interfaces gigabitEthernet0/10 status
show interfaces gigabitEthernet0/10 switchport
show interfaces trunk
show power inline
show cdp neighbors detail
show lldp neighbors detail
```

Agar AP PoE bilan ishlasa:

```cisco
show power inline gigabitEthernet0/10
```

WLCda (platformaga qarab komanda farq qiladi) umumiy tekshiruv:

```cisco
show ap summary
show wlan summary
show client summary
show wireless client summary
show wireless profile policy summary
```

Packet Tracer yoki eski WLC lablarda:

```cisco
show client summary
show ap summary
```

## Troubleshooting

Muammo: AP yoqilmayapti.

- PoE bormi?
- Switch port shutdown emasmi?
- Kabel ishlaydimi?
- `show power inline` APga quvvat berilganini ko'rsatyaptimi?

Muammo: SSID ko'rinmayapti.

- WLAN enabledmi?
- AP WLCga join bo'lganmi?
- Radio admin upmi?
- SSID broadcast o'chirilmaganmi?
- Country/regulatory domain mosmi?

Muammo: Client ulanadi, lekin IP olmaydi.

- SSID VLAN mapping to'g'rimi?
- AP port trunk allowed VLANda client VLAN bormi?
- DHCP scope bormi?
- DHCP relay kerakmi?
- WLC interface/gateway to'g'rimi?

Muammo: Ulanadi, lekin internet yoki ichki tarmoq ishlamaydi.

- Client IP, mask, gateway to'g'rimi?
- VLAN routing ishlayaptimi?
- ACL yoki firewall bloklamayaptimi?
- Guest VLAN internetga ruxsatlanganmi?

## Common mistakes

- AP portini access qilib, bir nechta SSID VLAN kutish.
- AP management native VLANini switch trunk native VLANiga moslamaslik.
- Guest va corporate SSIDni bitta VLANga qo'yish.
- WPA2-Enterprise uchun RADIUS reachabilityni tekshirmaslik.
- 2.4 GHzda 1, 6, 11 o'rniga overlapping kanallar ishlatish.
- Transmit powerni hamma APda maksimal qilish.
- DHCP muammosini wireless muammo deb o'ylash.

## Qisqa Q&A

**Savol:** SSID VLANmi?  
**Javob:** Yo'q. SSID wireless tarmoq nomi, lekin u odatda bitta VLANga mapping qilinadi.

**Savol:** AP switchga access port orqali ulanishi mumkinmi?  
**Javob:** Ha, agar faqat bitta VLAN/SSID ishlatilsa. Bir nechta VLAN uchun trunk kerak bo'ladi.

**Savol:** Client ulanib IP olmasa, birinchi nima tekshiriladi?  
**Javob:** VLAN mapping, trunk allowed VLAN va DHCP.
