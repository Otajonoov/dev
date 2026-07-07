# Routing Table

Routing table - router paketni qayerga yuborishni hal qilish uchun ishlatadigan jadval. Har bir route odatda destination network, mask/prefix, next-hop yoki exit interface, administrative distance va metric ma'lumotlarini saqlaydi.

## Routing table komponentlari

`show ip route` natijasida quyidagiga o'xshash yozuvlar ko'rinadi:

```cisco
R1# show ip route

Gateway of last resort is 192.168.12.2 to network 0.0.0.0

C    192.168.12.0/24 is directly connected, GigabitEthernet0/0
L    192.168.12.1/32 is directly connected, GigabitEthernet0/0
S*   0.0.0.0/0 [1/0] via 192.168.12.2
O    10.10.10.0/24 [110/20] via 192.168.12.2, 00:01:22, GigabitEthernet0/0
```

Muhim qismlar:

- `C`, `L`, `S`, `O` - route manbasi kodi.
- `10.10.10.0/24` - destination prefix.
- `[110/20]` - birinchi son administrative distance, ikkinchi son metric.
- `via 192.168.12.2` - next-hop router.
- `GigabitEthernet0/0` - paket chiqadigan interface.
- `S*` - static default route. Yulduzcha candidate default route ekanini bildiradi.

## Route kodlari

| Kod | Ma'nosi |
| --- | --- |
| `C` | Connected network |
| `L` | Local host route, router interfeys IP manzili |
| `S` | Static route |
| `S*` | Static default route |
| `O` | OSPF route |
| `D` | EIGRP route |
| `R` | RIP route |
| `B` | BGP route |

## Longest prefix match

Router avval administrative distance yoki metricga emas, eng aniq routega qaraydi. Bu "longest prefix match" deyiladi.

Misol:

```text
192.168.1.0/24     via R2
192.168.1.128/25   via R3
0.0.0.0/0          via ISP
```

Destination `192.168.1.150` bo'lsa, `/25` route tanlanadi, chunki u `/24` dan aniqroq. Destination `192.168.1.50` bo'lsa, `/24` route tanlanadi. Destination `8.8.8.8` bo'lsa, default route ishlaydi.

## Administrative distance

Administrative distance (AD) - route manbasiga ishonch darajasi. Son kichik bo'lsa, ishonch yuqori. AD faqat bir xil prefix uchun turli manbalardan route kelganda solishtiriladi.

| Route manbasi | Default AD |
| --- | ---: |
| Connected | 0 |
| Static | 1 |
| EIGRP summary | 5 |
| External BGP | 20 |
| EIGRP internal | 90 |
| OSPF | 110 |
| RIP | 120 |
| Floating static | Administrator bergan qiymat, masalan 200 |

Misol: routerda `10.1.1.0/24` static route ham, OSPF route ham bo'lsa, static route ishlaydi. Chunki static AD 1, OSPF AD 110.

## Metric

Metric - bir routing protokol ichida eng yaxshi yo'lni tanlash uchun ishlatiladi.

- OSPF metric: cost. Odatda bandwidth asosida hisoblanadi.
- RIP metric: hop count.
- EIGRP metric: bandwidth va delay kabi parametrlar.

AD route manbasini tanlaydi, metric esa shu manba ichida eng yaxshi yo'lni tanlaydi.

## Connected, local va host routes

Interfeysga IP berib `no shutdown` qilsangiz, router avtomatik `C` va `L` route yaratadi.

```cisco
interface GigabitEthernet0/0
 ip address 192.168.10.1 255.255.255.0
 no shutdown
```

Natija:

```cisco
C 192.168.10.0/24 is directly connected, GigabitEthernet0/0
L 192.168.10.1/32 is directly connected, GigabitEthernet0/0
```

`L` route host route hisoblanadi. IPv4 da host route maskasi `/32`, IPv6 da `/128`.

## Amaliy tekshiruv komandalar

```cisco
show ip route
show ip route connected
show ip route static
show ip route ospf
show ip route 10.10.10.5
show ip cef 10.10.10.5
show ip interface brief
```

`show ip route <ip>` aniq IP uchun qaysi route ishlashini ko'rsatadi. Troubleshootingda juda foydali.

## Common mistakes

- Longest prefix matchni unutish: default route borligi aniqroq route ishlamasligini anglatmaydi.
- AD va metricni aralashtirish: AD turli manbalar orasida, metric bitta protokol ichida ishlaydi.
- Interfeys `administratively down` bo'lsa connected route chiqmaydi.
- Static route next-hop manzili unreachable bo'lsa route jadvalda ko'rinmasligi yoki ishlamasligi mumkin.
- Return route yo'qligini unutish: ping ketishi mumkin, lekin javob qaytishi uchun teskari route ham kerak.

## Q&A

**Savol:** `L 192.168.1.1/32` nima uchun kerak?

**Javob:** Bu router interfeysining o'z IP manzili uchun local route. Router aynan shu IPga kelgan paketni o'ziga tegishli deb biladi.

**Savol:** Metrici past OSPF route static routedan ustun bo'ladimi?

**Javob:** Agar prefix bir xil bo'lsa, yo'q. Avval AD solishtiriladi. Static AD 1, OSPF AD 110.

**Savol:** Ikki route bir xil prefix, bir xil AD va bir xil metric bilan kelsa nima bo'ladi?

**Javob:** Router equal-cost load balancing qilishi mumkin. Cisco IOS odatda bir nechta teng yo'lni routing tablega qo'shadi.
