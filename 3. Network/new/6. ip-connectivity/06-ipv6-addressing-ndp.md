# IPv6 Addressing, SLAAC va NDP

IPv6 addressing IPv4dan kengroq va boshqacharoq ko'rinadi, lekin asosiy mantiq oddiy: hostlar global unicast manzil bilan tarmoqda qatnashadi, link-local manzil bilan local link ichida gaplashadi, NDP esa ARP va ayrim discovery vazifalarini bajaradi.

## IPv6 manzil ko'rinishi

IPv6 manzil 128 bit. Hexadecimal bloklar bilan yoziladi:

```text
2001:0db8:0000:0000:0000:0000:0000:0001
```

Qisqartirish qoidalari:

- Blok boshidagi nollar tashlab yuboriladi.
- Ketma-ket `0000` bloklar bir marta `::` bilan qisqartiriladi.

Misol:

```text
2001:0db8:0000:0000:0000:0000:0000:0001
2001:db8::1
```

`::` bir manzilda faqat bir marta ishlatiladi. Aks holda manzilni qayta tiklash noaniq bo'lib qoladi.

## IPv6 manzil turlari

| Tur | Prefix | Ma'nosi |
| --- | --- | --- |
| Global unicast | Odatda `2000::/3` | Internet/routed tarmoqlarda ishlatiladi |
| Link-local | `fe80::/10` | Faqat bitta link ichida ishlaydi |
| Unique local | `fc00::/7`, ko'p amaliyotda `fd00::/8` | Private IPv6ga o'xshash |
| Multicast | `ff00::/8` | Bir guruh qurilmalarga yuborish |
| Loopback | `::1/128` | O'zini tekshirish |
| Unspecified | `::/128` | Hali manzil yo'q degani |

IPv6da broadcast yo'q. Broadcast o'rnida multicast ishlatiladi.

## Prefix uzunligi

LAN segmentlar uchun eng ko'p ishlatiladigan prefix `/64`.

```text
2001:db8:10:1::/64
```

SLAAC normal ishlashi uchun odatda `/64` kerak. Point-to-point linklarda ham ko'p tashkilotlar soddalik uchun `/64` ishlatadi, lekin ayrim dizaynlarda `/127` uchraydi.

## Link-local manzil

Har bir IPv6-enabled interface link-local manzilga ega bo'ladi. U `fe80::/10` ichida bo'ladi va faqat shu L2 segment ichida ishlaydi.

Cisco routerda qo'lda berish:

```cisco
R1(config)# interface GigabitEthernet0/0
R1(config-if)# ipv6 address fe80::1 link-local
```

Tekshirish:

```cisco
show ipv6 interface brief
show ipv6 interface GigabitEthernet0/0
```

Routing protokollar va static IPv6 route next-hoplari ko'pincha link-local manzillardan foydalanadi.

## Global unicast manzil berish

Qo'lda:

```cisco
R1(config)# interface GigabitEthernet0/0
R1(config-if)# ipv6 address 2001:db8:10:1::1/64
R1(config-if)# no shutdown
```

EUI-64 bilan:

```cisco
R1(config-if)# ipv6 address 2001:db8:10:1::/64 eui-64
```

EUI-64 MAC address asosida interface ID yaratadi. Zamonaviy hostlarda privacy extension sabab real MACdan doim foydalanilmasligi mumkin.

## SLAAC

SLAAC (Stateless Address Autoconfiguration) hostga DHCPsiz IPv6 manzil yaratish imkonini beradi. Router Router Advertisement (RA) yuboradi, host prefixni oladi va o'z interface ID qismini yaratadi.

Router interface:

```cisco
R1(config)# ipv6 unicast-routing
R1(config)# interface GigabitEthernet0/0
R1(config-if)# ipv6 address 2001:db8:10:1::1/64
R1(config-if)# no shutdown
```

Host RA orqali `2001:db8:10:1::/64` prefixni ko'radi va o'z manzilini yaratadi.

RA flaglar:

- `A` flag: SLAAC orqali address yaratish mumkin.
- `M` flag: managed address, DHCPv6 orqali address oling.
- `O` flag: other config, masalan DNSni DHCPv6dan oling.

Cisco IOSda DHCPv6 bilan bog'liq sozlamalar:

```cisco
interface GigabitEthernet0/0
 ipv6 nd managed-config-flag
 ipv6 nd other-config-flag
```

## NDP nima qiladi?

NDP (Neighbor Discovery Protocol) ICMPv6 asosida ishlaydi va IPv6da bir nechta vazifani bajaradi:

- Neighbor Solicitation (NS): "Bu IPv6 manzil kimda?"
- Neighbor Advertisement (NA): "Bu manzil menda, MAC mana."
- Router Solicitation (RS): host routerlardan RA so'raydi.
- Router Advertisement (RA): router prefix, default gateway va flaglarni e'lon qiladi.
- DAD (Duplicate Address Detection): manzil takror emasligini tekshiradi.

IPv4dagi ARP vazifasini NDP bajaradi, lekin NDP bundan kengroq.

## Solicited-node multicast

IPv6 ARP kabi broadcast yubormaydi. Host biror IPv6 manzilning MACini topmoqchi bo'lsa, solicited-node multicast manzilga NS yuboradi.

Misol:

```text
IPv6: 2001:db8::1234
Solicited-node multicast: ff02::1:ff00:1234
```

Bu broadcastdan ko'ra samaraliroq, chunki hamma hostlar emas, faqat tegishli multicast guruhini tinglayotgan hostlar javob beradi.

## Neighbor table

```cisco
show ipv6 neighbors
```

Natija:

```cisco
IPv6 Address                              Age Link-layer Addr State Interface
FE80::2                                    12 aabb.cc00.0200  REACH Gi0/0
2001:DB8:10:1::10                          1 aabb.cc00.1010  STALE Gi0/0
```

State misollari:

- `REACH` - neighbor yaqinda reachable bo'lgan.
- `STALE` - ma'lumot bor, lekin yangilanishi mumkin.
- `DELAY/PROBE` - reachability tekshirilmoqda.

## Verification

```cisco
show ipv6 interface brief
show ipv6 interface GigabitEthernet0/0
show ipv6 neighbors
show ipv6 route connected
ping ipv6 fe80::2
ping ipv6 2001:db8:10:1::10
```

Link-local pingda ba'zi IOS versiyalar interface so'raydi, chunki link-local manzil faqat link ichida noyob:

```cisco
ping ipv6 fe80::2
```

Keyin chiqish interface sifatida `GigabitEthernet0/0` tanlanadi.

## Common mistakes

- `::` qisqartirishni bir manzilda ikki marta ishlatish.
- SLAAC uchun `/64` o'rniga boshqa prefix ishlatish.
- Link-local manzilni boshqa subnetga route qilinadi deb o'ylash.
- IPv6da broadcast bor deb hisoblash.
- ICMPv6ni to'liq bloklash: NDP, PMTUD va RA ishlamay qolishi mumkin.
- Duplicate address muammosini ko'rmaslik: DAD xatolari interface loglarida chiqadi.

## Q&A

**Savol:** Host default gateway sifatida qaysi manzilni ishlatadi?

**Javob:** Ko'pincha routerning link-local manzilini. RA orqali host default router haqida biladi.

**Savol:** IPv6da NAT shartmi?

**Javob:** Odatda yo'q. IPv6 katta manzil maydoni sabab end-to-end routingni qo'llab-quvvatlaydi. Security uchun firewall ishlatiladi, NAT shart emas.

**Savol:** NDP ARPning aynan o'zi-mi?

**Javob:** Yo'q. NDP ARPga o'xshash neighbor MAC topish vazifasini bajaradi, lekin router discovery, SLAAC va DAD kabi qo'shimcha vazifalari ham bor.
