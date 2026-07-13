# 02 - Network Layer va IP

Bu modul **Network layer** (OSI 3-qatlam) va **IP** protokolini chuqur o'rgatadi:
IPv4 header'dan tortib subnetting hisoblari, ARP, NAT va IPv6'gacha. Maqsad --
"IP address" degan tushunchani sehr bo'lishdan to'xtatib, uni binary darajada,
qadam-baqadam hisoblab tushunish.

## Nima o'rganiladi

- Network layer vazifasi, IPv4 header, TTL, fragmentation, MTU va PMTUD.
- IP address tuzilishi: binary, oktet, network/host qismlar.
- Subnetting, CIDR, VLSM va wildcard mask (ko'p ishlangan hisoblar bilan).
- Network/broadcast/host range topish (binary AND va block size usullari).
- Address turlari: public, private (RFC 1918), special va classful tarix.
- ARP jarayoni, ARP cache, default gateway va ARP spoofing himoyasi.
- NAT turlari (static/dynamic/PAT), CGNAT va NAT traversal.
- IPv6 addressing, SLAAC va NDP.

## Darslar ro'yxati (o'qish tartibi)

1. [01-network-layer-va-ipv4.md](01-network-layer-va-ipv4.md) -- Network layer, IPv4 header, TTL, fragmentation
2. [02-ip-addressing.md](02-ip-addressing.md) -- Binary, oktet, IP address tuzilishi
3. [03-subnetting-cidr-vlsm.md](03-subnetting-cidr-vlsm.md) -- Subnet mask, CIDR, VLSM, wildcard
4. [04-network-broadcast-host-range.md](04-network-broadcast-host-range.md) -- Network/broadcast/host range topish
5. [05-address-types-classful-classless.md](05-address-types-classful-classless.md) -- Public/private/special, classful tarix
6. [06-arp-va-default-gateway.md](06-arp-va-default-gateway.md) -- ARP, default gateway
7. [07-nat.md](07-nat.md) -- NAT (static/dynamic/PAT), CGNAT
8. [08-ipv6-addressing-ndp.md](08-ipv6-addressing-ndp.md) -- IPv6, SLAAC, NDP

## O'qish tartibi haqida

Darslarni **tartib bilan** o'qi -- har biri oldingisiga tayanadi. Ayniqsa
2-dars (binary) yaxshi o'zlashtirilmasa, 3-4 darslardagi subnetting hisoblari
qiyin bo'ladi. Har dars oxiridagi "O'z-o'zini tekshir" va "Amaliyot" bo'limlarini
albatta bajar -- subnetting faqat mashq bilan mustahkamlanadi.

Routing protokollari (OSPF, BGP) bu modulda emas -- ular alohida modulda
chuqur o'rganiladi. Bu modulda routing'ning faqat asosiy g'oyasi (longest
prefix match, default gateway) beriladi.
