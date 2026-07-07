# Static Routing

Static route - administrator qo'lda yozadigan marshrut. Kichik tarmoqlar, stub networklar, default internet chiqishi va backup yo'llar uchun juda ko'p ishlatiladi.

## Oddiy static route

Topologiya:

```text
LAN-A 192.168.10.0/24 -- R1 -- 192.168.12.0/24 -- R2 -- LAN-B 192.168.20.0/24
```

R1 dan LAN-B ga route:

```cisco
R1(config)# ip route 192.168.20.0 255.255.255.0 192.168.12.2
```

R2 dan LAN-A ga route:

```cisco
R2(config)# ip route 192.168.10.0 255.255.255.0 192.168.12.1
```

Bu yerda `192.168.12.2` va `192.168.12.1` next-hop IP manzillar.

## Exit interface bilan static route

Point-to-point linklarda exit interface ko'rsatish ham mumkin:

```cisco
R1(config)# ip route 192.168.20.0 255.255.255.0 Serial0/0/0
```

Ethernet multi-access tarmoqlarda faqat exit interface ko'rsatish tavsiya etilmaydi, chunki router ARP jarayonida ko'proq ishlashi mumkin. Ethernetda next-hop IP yoki next-hop + interface ishlatish yaxshiroq:

```cisco
R1(config)# ip route 192.168.20.0 255.255.255.0 GigabitEthernet0/0 192.168.12.2
```

## Default route

Default route destination `0.0.0.0/0`. Router boshqa aniq route topmasa shu yo'lni ishlatadi.

```cisco
R1(config)# ip route 0.0.0.0 0.0.0.0 203.0.113.1
```

Internetga chiqadigan edge routerlarda ko'p ishlatiladi. Stub router uchun ham qulay: barcha noma'lum tarmoqlarni upstream routerga yuboradi.

## Host route

Host route bitta IP manzil uchun yoziladi. IPv4 host route maskasi `/32`.

```cisco
R1(config)# ip route 10.10.10.50 255.255.255.255 192.168.12.2
```

Bu route faqat `10.10.10.50` uchun ishlaydi. Longest prefix match sababli `/32` juda aniq route.

## Floating static route

Floating static route - backup static route. Buning AD qiymati asosiy route AD qiymatidan yuqori qilinadi.

Masalan, asosiy yo'l OSPF orqali kelsin (AD 110), backup statik route esa faqat OSPF route yo'qolganda ishlasin:

```cisco
R1(config)# ip route 10.20.30.0 255.255.255.0 192.168.99.2 200
```

Oxiridagi `200` - administrative distance. OSPF route bor paytda bu route ishlamaydi. OSPF yo'qolsa, backup route jadvalga tushadi.

## Recursive lookup

Next-hop IP ko'rsatilgan static routeda router avval destination routega qaraydi, keyin next-hop manzilga qanday yetishni alohida topadi. Bu recursive lookup deyiladi.

```cisco
S 10.10.10.0/24 [1/0] via 192.168.12.2
C 192.168.12.0/24 is directly connected, GigabitEthernet0/0
```

Router `10.10.10.0/24` uchun next-hop `192.168.12.2` ekanini biladi, keyin `192.168.12.2` connected tarmoqda ekanini topadi va paketni `GigabitEthernet0/0` dan chiqaradi.

## Verification

```cisco
show ip route static
show ip route 192.168.20.10
show running-config | include ^ip route
ping 192.168.20.10
traceroute 192.168.20.10
show arp
```

Muammo bo'lsa:

```cisco
show ip interface brief
show interfaces GigabitEthernet0/0
ping 192.168.12.2
```

Avval next-hop reachable ekanini tekshiring. Next-hopga ping bormasa, uzoq destinationga ham bormaydi.

## Static route o'chirish

Qanday yozilgan bo'lsa, shunday `no` bilan o'chiriladi:

```cisco
R1(config)# no ip route 192.168.20.0 255.255.255.0 192.168.12.2
```

## Common mistakes

- Faqat bir tomonga route yozish: trafik borishi va qaytishi uchun ikki tomonda ham route kerak.
- Wrong next-hop: next-hop odatda qo'shni routerning shu linkdagi IP manzili bo'lishi kerak.
- Default routega haddan tashqari ishonish: ichki tarmoqlar uchun aniq route kerak bo'lishi mumkin.
- Floating static AD qiymatini past qo'yish: backup route asosiy route o'rniga ishlay boshlaydi.
- Ethernetda faqat exit interface yozish: keraksiz ARP muammolariga olib kelishi mumkin.

## Q&A

**Savol:** Static route qachon jadvalga tushmaydi?

**Javob:** Next-hopga yetish uchun route bo'lmasa yoki exit interface down bo'lsa, static route ishlamasligi mumkin.

**Savol:** Default route ham static routemi?

**Javob:** Agar `ip route 0.0.0.0 0.0.0.0 ...` bilan yozilsa, ha, bu static default route.

**Savol:** Floating static route uchun qaysi AD tanlanadi?

**Javob:** Asosiy route AD qiymatidan yuqori bo'lishi kerak. Masalan, OSPF backupi uchun 111 yoki 200 ishlatish mumkin.
