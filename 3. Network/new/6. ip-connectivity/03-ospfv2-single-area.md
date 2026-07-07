# OSPFv2 Single-Area

OSPFv2 - IPv4 uchun link-state routing protokoli. Single-area OSPF deganda barcha routerlar bir area ichida ishlaydi, odatda `area 0`. CCNA darajasida OSPF neighbor adjacency, router-id, wildcard mask, DR/BDR va routing tabledagi OSPF routelarni tushunish muhim.

## OSPF asoslari

OSPF routerlar bir-biri bilan Hello paket almashadi, neighbor bo'ladi, keyin link-state ma'lumotlarini almashadi. Har bir router LSDB (link-state database) tuzadi va SPF algoritmi orqali eng yaxshi yo'llarni hisoblaydi.

OSPF route routing tableda `O` kodi bilan ko'rinadi:

```cisco
O 10.10.20.0/24 [110/2] via 192.168.12.2, 00:00:32, GigabitEthernet0/0
```

Bu yerda AD `110`, metric esa OSPF cost.

## Router ID

Router-id - OSPF ichida routerni tanitadigan 32-bit qiymat. IP manzilga o'xshaydi, lekin interface IP bo'lishi shart emas.

Tanlash tartibi:

1. `router-id` komandasi bilan qo'lda berilgan qiymat.
2. Eng katta up/up loopback IP.
3. Eng katta up/up physical interface IP.

Tavsiya: har doim qo'lda bering.

```cisco
R1(config)# router ospf 1
R1(config-router)# router-id 1.1.1.1
```

Router-id o'zgartirilgandan keyin OSPF processni qayta ishga tushirish kerak bo'lishi mumkin:

```cisco
R1# clear ip ospf process
```

## Network komandasi va wildcard mask

OSPFda `network` komandasi qaysi interfeyslar OSPFga kirishini aniqlaydi.

```cisco
R1(config)# router ospf 1
R1(config-router)# network 192.168.12.0 0.0.0.255 area 0
R1(config-router)# network 10.1.1.0 0.0.0.255 area 0
```

Wildcard mask subnet maskning teskarisi:

| Subnet mask | Wildcard mask |
| --- | --- |
| 255.255.255.0 | 0.0.0.255 |
| 255.255.255.252 | 0.0.0.3 |
| 255.255.0.0 | 0.0.255.255 |
| 255.255.255.255 | 0.0.0.0 |

Bitta interfeysni aniq qo'shish:

```cisco
R1(config-router)# network 192.168.12.1 0.0.0.0 area 0
```

## Interface ostida OSPF yoqish

Zamonaviy va aniqroq usul:

```cisco
R1(config)# interface GigabitEthernet0/0
R1(config-if)# ip ospf 1 area 0
```

Bu usul wildcard mask xatosini kamaytiradi.

## Neighbor adjacency

Neighbor bo'lish uchun asosiy shartlar:

- Ikki router bir L2 segmentda bo'lishi kerak.
- Area ID bir xil bo'lishi kerak.
- Hello/dead timer mos bo'lishi kerak.
- Subnet va mask mos bo'lishi kerak.
- Authentication sozlangan bo'lsa, parollar mos bo'lishi kerak.
- Stub/NSSA kabi area flaglari mos bo'lishi kerak.
- Router-id takrorlanmasligi kerak.

Tekshirish:

```cisco
show ip ospf neighbor
show ip ospf interface brief
show ip ospf interface GigabitEthernet0/0
show ip protocols
```

Neighbor holatlari:

| Holat | Ma'nosi |
| --- | --- |
| Down | Hello ko'rinmayapti |
| Init | Hello keldi, lekin o'z routerimiz neighbor listda yo'q |
| 2-Way | Ikki tomon bir-birini ko'rdi |
| ExStart/Exchange | DB almashishga tayyorgarlik |
| Loading | LSA ma'lumotlari olinmoqda |
| Full | Adjacency to'liq |

Broadcast tarmoqlarda barcha neighborlar bilan `Full` bo'lish shart emas. DROTHER routerlar bir-biri bilan `2-Way` holatda qolishi normal.

## DR va BDR

Ethernet kabi broadcast tarmoqlarda OSPF DR (Designated Router) va BDR (Backup Designated Router) saylaydi. Maqsad: LSA almashuvini kamaytirish.

Saylovda:

1. Eng katta OSPF priority yutadi.
2. Priority teng bo'lsa, eng katta router-id yutadi.
3. Priority `0` bo'lsa, router DR/BDR bo'la olmaydi.

```cisco
interface GigabitEthernet0/0
 ip ospf priority 100
```

DR/BDR saylovi preemptive emas. Ya'ni keyinroq kuchliroq router qo'shilsa, mavjud DR avtomatik almashtirilmaydi.

## Passive interface

LAN tomonda OSPF route e'lon qilish kerak, lekin user segmentida neighbor kerak emas. Bunday holatda passive interface ishlatiladi.

```cisco
R1(config)# router ospf 1
R1(config-router)# passive-interface GigabitEthernet0/1
```

Hamma interfacelarni passive qilib, faqat router-router linkni ochish:

```cisco
R1(config-router)# passive-interface default
R1(config-router)# no passive-interface GigabitEthernet0/0
```

## Default route e'lon qilish

Edge routerda default route bor bo'lsa, OSPF orqali tarqatish:

```cisco
R1(config)# ip route 0.0.0.0 0.0.0.0 203.0.113.1
R1(config)# router ospf 1
R1(config-router)# default-information originate
```

## Troubleshooting

```cisco
show ip ospf neighbor
show ip ospf interface brief
show ip route ospf
show ip protocols
show running-config | section router ospf
ping 224.0.0.5
```

Kerak bo'lsa vaqtincha debug:

```cisco
debug ip ospf hello
undebug all
```

Production tarmoqda debugni ehtiyotkorlik bilan ishlating.

## Common mistakes

- Wildcard maskni subnet mask bilan adashtirish.
- Area bir tomonda `0`, boshqa tomonda `1` bo'lib qolishi.
- Router-id takrorlanishi.
- Passive interface sabab Hello chiqmasligi.
- Broadcast segmentda `2-Way` holatni xato deb o'ylash.
- DR/BDR saylovini darhol o'zgaradi deb kutish.

## Q&A

**Savol:** OSPF process ID bir xil bo'lishi shartmi?

**Javob:** Yo'q. `router ospf 1` va `router ospf 10` qo'shni bo'lishi mumkin. Muhimi area, timer, subnet va boshqa OSPF parametrlar mos bo'lishi.

**Savol:** OSPF metric nimaga asoslanadi?

**Javob:** Costga. Cisco IOSda cost odatda bandwidth asosida hisoblanadi, lekin interface ostida qo'lda ham berish mumkin: `ip ospf cost 10`.

**Savol:** Neighbor `2-Way` bo'lsa har doim muammomi?

**Javob:** Yo'q. Broadcast networkda DROTHER-DROTHER orasida `2-Way` normal holat.
