# FHRP: HSRP, VRRP, GLBP

FHRP (First Hop Redundancy Protocol) - hostlar ishlatadigan default gateway uchun redundancy beradi. Oddiy LANda PC default gateway sifatida bitta routerning IP manzilini biladi. Agar shu router ishdan chiqsa, PC boshqa routerni avtomatik tanlamaydi. FHRP virtual gateway yaratib, bu muammoni hal qiladi.

## FHRP qanday ishlaydi

Ikki yoki undan ko'p router bitta virtual IP va virtual MAC manzilni baham ko'radi. Hostlar default gateway sifatida virtual IPni yozadi.

```text
PC default gateway: 192.168.10.254

R1 real IP: 192.168.10.1
R2 real IP: 192.168.10.2
Virtual IP: 192.168.10.254
```

R1 active/master bo'lsa, virtual IPga kelgan trafikni R1 qabul qiladi. R1 ishlamay qolsa, R2 virtual gateway rolini oladi.

## HSRP

HSRP - Cisco proprietary FHRP. Active va Standby router tushunchalarini ishlatadi.

Asosiy konfiguratsiya:

```cisco
R1(config)# interface GigabitEthernet0/1
R1(config-if)# ip address 192.168.10.1 255.255.255.0
R1(config-if)# standby 10 ip 192.168.10.254
R1(config-if)# standby 10 priority 110
R1(config-if)# standby 10 preempt
```

```cisco
R2(config)# interface GigabitEthernet0/1
R2(config-if)# ip address 192.168.10.2 255.255.255.0
R2(config-if)# standby 10 ip 192.168.10.254
R2(config-if)# standby 10 priority 100
R2(config-if)# standby 10 preempt
```

Muhim joylar:

- `10` - HSRP group number.
- Priority yuqori bo'lgan router active bo'ladi.
- Default priority `100`.
- `preempt` bo'lmasa, yuqori priorityli router qaytib kelganda active rolini avtomatik qaytarib olmaydi.

Tekshirish:

```cisco
show standby brief
show standby
```

## Interface tracking

Router LAN tomonda active bo'lishi mumkin, lekin uning WAN uplinki uzilgan bo'lsa, trafik noto'g'ri yo'lga ketadi. Tracking priorityni kamaytiradi.

```cisco
R1(config)# interface GigabitEthernet0/1
R1(config-if)# standby 10 track GigabitEthernet0/0 30
```

GigabitEthernet0/0 down bo'lsa, R1 prioritysi 30 ga kamayadi. R2 prioritysi yuqoriroq bo'lib qolsa, active rolini oladi.

## VRRP

VRRP - standart protokol. Cisco bo'lmagan qurilmalar bilan ishlashda foydali. Master va Backup tushunchalarini ishlatadi.

```cisco
R1(config)# interface GigabitEthernet0/1
R1(config-if)# ip address 192.168.10.1 255.255.255.0
R1(config-if)# vrrp 10 ip 192.168.10.254
R1(config-if)# vrrp 10 priority 110
R1(config-if)# vrrp 10 preempt
```

```cisco
R2(config)# interface GigabitEthernet0/1
R2(config-if)# ip address 192.168.10.2 255.255.255.0
R2(config-if)# vrrp 10 ip 192.168.10.254
R2(config-if)# vrrp 10 priority 100
R2(config-if)# vrrp 10 preempt
```

Tekshirish:

```cisco
show vrrp brief
show vrrp
```

## GLBP

GLBP - Cisco proprietary. Faqat redundancy emas, balki load balancing ham beradi. Bitta virtual IP bo'ladi, lekin bir nechta virtual MAC ishlatiladi. GLBP ichida AVG (Active Virtual Gateway) va AVF (Active Virtual Forwarder) rollari bor.

```cisco
R1(config)# interface GigabitEthernet0/1
R1(config-if)# ip address 192.168.10.1 255.255.255.0
R1(config-if)# glbp 10 ip 192.168.10.254
R1(config-if)# glbp 10 priority 110
R1(config-if)# glbp 10 preempt
```

```cisco
R2(config)# interface GigabitEthernet0/1
R2(config-if)# ip address 192.168.10.2 255.255.255.0
R2(config-if)# glbp 10 ip 192.168.10.254
R2(config-if)# glbp 10 priority 100
R2(config-if)# glbp 10 preempt
```

Tekshirish:

```cisco
show glbp brief
show glbp
```

## HSRP, VRRP, GLBP taqqoslash

| Protokol | Turi | Rollar | Load balancing |
| --- | --- | --- | --- |
| HSRP | Cisco proprietary | Active/Standby | Yo'q, odatda bitta active |
| VRRP | Standard | Master/Backup | Yo'q, odatda bitta master |
| GLBP | Cisco proprietary | AVG/AVF | Ha |

## Troubleshooting

```cisco
show standby brief
show vrrp brief
show glbp brief
show ip interface brief
show arp
show mac address-table
ping 192.168.10.254
traceroute 8.8.8.8
```

Host tomonda:

```text
Default gateway virtual IPga tengmi?
Host ARP table virtual MACni ko'ryaptimi?
Failoverdan keyin gateway ping bo'ladimi?
```

## Common mistakes

- Host default gatewayga real router IP yozish. FHRP ishlashi uchun virtual IP yozilishi kerak.
- HSRP/VRRP group number yoki virtual IP mos kelmasligi.
- `preempt` qo'yilmagani sabab active rol qaytmasligi.
- Uplink tracking yo'qligi: LAN gateway tirik, lekin internet yo'li o'lik.
- VLAN/SVI noto'g'ri: ikkala router bir xil L2 broadcast domainda bo'lishi kerak.

## Q&A

**Savol:** HSRP virtual IP router real IP manzili bo'lishi mumkinmi?

**Javob:** HSRPda odatda virtual IP alohida bo'ladi. VRRPda esa real IP egasi master bo'lishi kabi holatlar uchrashi mumkin, lekin amaliyotda alohida virtual IP ishlatish tushunarliroq.

**Savol:** FHRP routing protokol o'rnini bosadimi?

**Javob:** Yo'q. FHRP hostlarning birinchi gateway muammosini hal qiladi. Routerlar orasidagi marshrutlash uchun static route, OSPF va boshqa routing kerak bo'lishi mumkin.

**Savol:** GLBP nima uchun kerak?

**Javob:** Bir nechta router orqali host trafiklarini taqsimlash uchun. HSRP/VRRPda odatda bitta router active/master bo'ladi, GLBP esa load balancing bera oladi.
