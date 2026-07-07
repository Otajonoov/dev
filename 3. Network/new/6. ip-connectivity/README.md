# 6. IP Connectivity

Bu bo'lim router qanday qaror qabul qilishi, statik va dinamik marshrutlash, birinchi hop redundancy va IPv6 asoslarini amaliy tarzda tushuntiradi. CCNA darajasida eng muhim savol: "Paket qaysi interfeysdan chiqadi va nima uchun?"

## Mavzular

1. [Routing table](01-routing-table.md)
2. [Static routing](02-static-routing.md)
3. [OSPFv2 single-area](03-ospfv2-single-area.md)
4. [FHRP: HSRP, VRRP, GLBP](04-fhrp.md)
5. [IPv6 routing](05-ipv6-routing.md)
6. [IPv6 addressing va NDP](06-ipv6-addressing-ndp.md)

## Asosiy g'oya

Router paketni qabul qilganda destination IP manzilga qaraydi, routing table ichidan eng mos route topadi va paketni keyingi hopga yuboradi. Eng mos route tanlashda odatda quyidagi tartib ishlaydi:

1. Longest prefix match: eng uzun maskali route yutadi.
2. Administrative distance: bir xil prefix uchun ishonch darajasi pastroq son bo'lsa yutadi.
3. Metric: bir xil routing protokol ichida eng yaxshi yo'l tanlanadi.
4. Equal-cost bo'lsa, router odatda load balancing qilishi mumkin.

## Ko'p ishlatiladigan komandalar

```cisco
show ip route
show ip route 192.168.10.50
show ip protocols
show ip interface brief
show interfaces
ping 8.8.8.8
traceroute 8.8.8.8
```

IPv6 uchun:

```cisco
show ipv6 route
show ipv6 interface brief
show ipv6 neighbors
ping ipv6 2001:db8::1
traceroute ipv6 2001:db8::1
```

## Route turlari qisqacha

| Tur | Misol | Ma'nosi |
| --- | --- | --- |
| Connected | `C 192.168.1.0/24 is directly connected` | Tarmoq router interfeysida bor |
| Local/Host | `L 192.168.1.1/32` | Router interfeysining o'z IP manzili |
| Static | `S 10.10.10.0/24 [1/0] via 192.168.1.2` | Administrator qo'lda yozgan route |
| Default | `S* 0.0.0.0/0 via ...` | Hech qaysi route mos kelmasa ishlatiladi |
| OSPF | `O 10.1.1.0/24 [110/2] via ...` | OSPF o'rgangan route |

## Troubleshooting tartibi

1. IP manzil va interfeys holatini tekshiring: `show ip interface brief`.
2. Routing table ichida destination route borligini tekshiring: `show ip route <destination>`.
3. Next-hop reachable ekanini tekshiring: `ping <next-hop>`.
4. Paket yo'lini ko'ring: `traceroute <destination>`.
5. Dinamik routing bo'lsa neighbor holatini tekshiring: `show ip ospf neighbor`.
6. Access-list, NAT yoki firewall paketni to'smayotganini alohida tekshiring.

## Tez Q&A

**Savol:** Routerda default route bo'lsa, hamma paketlar shundan chiqadimi?

**Javob:** Yo'q. Avval aniqroq route qidiriladi. Default route faqat boshqa hech bir route mos kelmaganda ishlaydi.

**Savol:** Static route OSPF routedan ustunmi?

**Javob:** Odatda ha, chunki static route AD 1, OSPF AD 110. Lekin longest prefix match birinchi turadi.

**Savol:** Connected route uchun alohida konfiguratsiya kerakmi?

**Javob:** Yo'q. Interfeysda IP bor va interface up/up bo'lsa, connected va local route avtomatik paydo bo'ladi.
