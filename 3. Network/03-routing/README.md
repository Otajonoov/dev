# 03 - Routing

Bu modul **routing** ni -- ya'ni "paket qaysi yo'ldan ketadi va nima uchun" degan
markaziy savolni -- asosdan chuqurgacha o'rgatadi. Router qanday qaror qabul
qilishidan (routing table, longest prefix match) boshlab, static va dynamic
routing (OSPF, BGP), diagnostika (ICMP, ping, traceroute), gateway redundancy
(FHRP) va IPv6 routing gacha.

Maqsad: routing ni "sehr" bo'lishdan to'xtatib, har bir qarorni -- qaysi route,
nega, qanday tekshirish -- Cisco CLI misollari va real hodisalar bilan tushunish.

## Nima o'rganiladi

- Routing table ni o'qish: manba kodlari, AD, metric, longest prefix match, FIB/CEF.
- Static routing: oddiy, default, host va floating static route.
- Dynamic routing protokollari: IGP vs EGP, distance-vector vs link-state vs path-vector.
- OSPF: neighbor adjacency, DR/BDR, LSA, area, konfiguratsiya va troubleshooting.
- BGP: eBGP/iBGP, path attributes, Internet routing, hijack/route leak, RPKI.
- ICMP: message turlari, ping, traceroute TTL tryuki, PMTUD.
- FHRP: HSRP, VRRP, GLBP, virtual IP/MAC, preempt, tracking.
- IPv6 routing va dual-stack: static/default route, link-local next-hop.

## Darslar ro'yxati (o'qish tartibi)

1. [01-routing-table-va-longest-prefix.md](01-routing-table-va-longest-prefix.md) -- Routing table, AD, metric, longest prefix match
2. [02-static-routing.md](02-static-routing.md) -- Static, default, host, floating static route
3. [03-routing-protocols-overview.md](03-routing-protocols-overview.md) -- IGP/EGP, DV/LS/PV, protokollar taqqoslash
4. [04-ospf.md](04-ospf.md) -- OSPF: adjacency, DR/BDR, LSA, area, konfiguratsiya
5. [05-bgp.md](05-bgp.md) -- BGP: eBGP/iBGP, path attributes, hijack, RPKI
6. [06-icmp-ping-traceroute.md](06-icmp-ping-traceroute.md) -- ICMP, ping, traceroute, PMTUD
7. [07-fhrp.md](07-fhrp.md) -- FHRP: HSRP, VRRP, GLBP
8. [08-ipv6-routing.md](08-ipv6-routing.md) -- IPv6 routing va dual-stack

## O'qish tartibi haqida

Darslarni **tartib bilan** o'qi -- har biri oldingisiga tayanadi. Ayniqsa 1-dars
(routing table, longest prefix match, AD/metric) yaxshi o'zlashtirilmasa, keyingi
darslar qiyin bo'ladi. 3-dars (protokollar umumiy ko'rinishi) OSPF va BGP ga
kirish darvozasi.

IP addressing, subnetting, ARP, NAT va IPv6 addressing/NDP bu modulda emas -- ular
avvalgi modulda o'rganilgan va bu yerda faqat kerak joyda eslatiladi. Bu modul
routing ga -- yo'l tanlashga -- e'tibor qaratadi.

Har dars oxiridagi "O'z-o'zini tekshir" va "Amaliyot" bo'limlarini albatta bajar.
Routing amaliy fan -- Packet Tracer, GNS3 yoki FRRouting da sinab ko'rgan bilangina
mustahkamlanadi.
