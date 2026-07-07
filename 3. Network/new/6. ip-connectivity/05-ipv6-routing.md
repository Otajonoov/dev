# IPv6 Routing

IPv6 routing g'oyasi IPv4 bilan o'xshash: router destination IPv6 manzilga qaraydi, routing table ichidan eng mos prefixni topadi va paketni next-hopga yuboradi. Farqi: IPv6 manzillar uzunroq, link-local manzillar muhimroq, default route `::/0`, host route esa `/128`.

## IPv6 routingni yoqish

Cisco routerda IPv6 routing global yoqilishi kerak:

```cisco
R1(config)# ipv6 unicast-routing
```

Bu komanda bo'lmasa, router IPv6 interfacelariga ega bo'lishi mumkin, lekin IPv6 paketlarni router sifatida forward qilmaydi.

## IPv6 interface manzili

```cisco
R1(config)# interface GigabitEthernet0/0
R1(config-if)# ipv6 address 2001:db8:12::1/64
R1(config-if)# no shutdown
```

Tekshirish:

```cisco
show ipv6 interface brief
show ipv6 interface GigabitEthernet0/0
```

Interfeysda global unicast manzil bilan birga avtomatik link-local manzil ham bo'ladi, odatda `FE80::/10` diapazonida.

## IPv6 routing table

```cisco
R1# show ipv6 route

C   2001:DB8:12::/64 [0/0]
     via GigabitEthernet0/0, directly connected
L   2001:DB8:12::1/128 [0/0]
     via GigabitEthernet0/0, receive
L   FE80::1/128 [0/0]
     via GigabitEthernet0/0, receive
S   2001:DB8:20::/64 [1/0]
     via FE80::2, GigabitEthernet0/0
```

Kodlar:

- `C` - connected IPv6 prefix.
- `L` - local `/128` route.
- `S` - static IPv6 route.
- `O` - OSPFv3 route.

## IPv6 static route

Global next-hop bilan:

```cisco
R1(config)# ipv6 route 2001:db8:20::/64 2001:db8:12::2
```

Exit interface bilan:

```cisco
R1(config)# ipv6 route 2001:db8:20::/64 GigabitEthernet0/0
```

Link-local next-hop bilan esa exit interface majburiy:

```cisco
R1(config)# ipv6 route 2001:db8:20::/64 GigabitEthernet0/0 fe80::2
```

Nega interface kerak? Chunki link-local manzil faqat bitta link ichida noyob. Turli interfacelarda bir xil `fe80::2` bo'lishi mumkin. Router qaysi linkdan chiqishni bilishi kerak.

## IPv6 default route

Default route `::/0`:

```cisco
R1(config)# ipv6 route ::/0 2001:db8:12::2
```

Link-local next-hop bilan:

```cisco
R1(config)# ipv6 route ::/0 GigabitEthernet0/0 fe80::2
```

## Floating IPv6 static route

IPv4 kabi IPv6 static routeda ham AD berish mumkin:

```cisco
R1(config)# ipv6 route 2001:db8:30::/64 2001:db8:99::2 200
```

Bu backup route sifatida ishlatiladi. Asosiy dinamik route yo'qolsa, AD 200 bo'lgan static route jadvalga kiradi.

## Longest prefix match IPv6da ham ishlaydi

```text
2001:db8::/32        via ISP
2001:db8:10::/48     via R2
2001:db8:10:5::/64   via R3
```

Destination `2001:db8:10:5::100` bo'lsa, `/64` tanlanadi. Destination `2001:db8:10:9::100` bo'lsa, `/48` tanlanadi.

## Verification

```cisco
show ipv6 route
show ipv6 route static
show ipv6 route 2001:db8:20::10
show ipv6 interface brief
show ipv6 neighbors
ping ipv6 2001:db8:20::10
traceroute ipv6 2001:db8:20::10
```

Next-hop link-local bo'lsa, neighbor table juda muhim:

```cisco
show ipv6 neighbors
```

## Troubleshooting tartibi

1. `ipv6 unicast-routing` borligini tekshiring.
2. Interfeys up/up va IPv6 manzil borligini tekshiring.
3. Destination prefix routing tableda borligini tekshiring.
4. Next-hop reachable ekanini ping qiling.
5. Link-local next-hop ishlatilsa, exit interface to'g'ri yozilganini tekshiring.
6. Teskari route borligini tekshiring.

## Common mistakes

- `ipv6 unicast-routing` komandasi unutiladi.
- Link-local next-hop uchun interface ko'rsatilmaydi.
- IPv6 default route `0.0.0.0/0` deb yoziladi. To'g'risi `::/0`.
- `/64` LAN prefixlar o'rniga noto'g'ri uzunlik ishlatiladi, SLAAC ishlamay qoladi.
- Pingda noto'g'ri format: ko'p IOS versiyalarida `ping ipv6 <address>` kerak.

## Q&A

**Savol:** IPv6 static route AD qiymati nechchi?

**Javob:** IPv4 static route kabi default AD 1.

**Savol:** IPv6da ARP bormi?

**Javob:** Yo'q. IPv6 ARP o'rniga NDP ishlatadi.

**Savol:** Link-local manzil routeda next-hop bo'la oladimi?

**Javob:** Ha, amaliyotda ko'p ishlatiladi. Lekin exit interface bilan birga yozish kerak.
