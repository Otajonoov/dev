# DHCP Relay va `ip helper-address`

DHCP client odatda broadcast yuboradi. Broadcast paketlar routerdan o'tmaydi. Shuning uchun DHCP server boshqa subnet yoki VLANda turgan bo'lsa, router yoki Layer 3 switch DHCP relay vazifasini bajaradi.

Cisco IOSda DHCP relay uchun interfeysga `ip helper-address` beriladi. Bu buyruq client joylashgan gateway interfeysida yoziladi.

## Oddiy topologiya

```text
PC VLAN 10:        192.168.10.0/24
Gateway SVI:       192.168.10.1
DHCP Server:       192.168.100.10
Server VLAN:       192.168.100.0/24
```

PC DHCP Discover broadcast yuboradi. Router uni DHCP serverga unicast qilib uzatadi.

## Cisco konfiguratsiya

Layer 3 switch yoki routerda:

```cisco
conf t
interface vlan 10
 description USERS_VLAN10_GATEWAY
 ip address 192.168.10.1 255.255.255.0
 ip helper-address 192.168.100.10
 no shutdown
end
```

Router subinterface ishlatilsa:

```cisco
conf t
interface gigabitEthernet0/0.10
 encapsulation dot1Q 10
 ip address 192.168.10.1 255.255.255.0
 ip helper-address 192.168.100.10
end
```

## DHCP server tomoni

DHCP serverda VLAN 10 uchun pool bo'lishi kerak:

```text
Network:        192.168.10.0/24
Default gateway:192.168.10.1
DNS server:     8.8.8.8
```

Relay paketida `giaddr` maydoni gateway manzili bilan to'ldiriladi. Server shu `giaddr` bo'yicha qaysi subnetga IP berishni biladi.

## `ip helper-address` nimalarni uzatadi?

Default holatda Cisco `ip helper-address` faqat DHCP emas, bir nechta UDP broadcast xizmatlarini relay qiladi:

- UDP 67: BOOTP/DHCP server
- UDP 68: BOOTP/DHCP client
- UDP 69: TFTP
- UDP 53: DNS
- UDP 37: Time
- UDP 49: TACACS
- UDP 137/138: NetBIOS

CCNAda eng ko'p uchraydigani DHCP. Keraksiz UDP forwardni cheklash uchun `no ip forward-protocol udp <port>` ishlatilishi mumkin.

```cisco
conf t
no ip forward-protocol udp tftp
no ip forward-protocol udp domain
end
```

## Tekshiruv buyruqlari

```cisco
show running-config interface vlan 10
show ip interface vlan 10
show ip route 192.168.100.10
ping 192.168.100.10 source vlan 10
debug ip dhcp server packet
debug ip udp
```

`debug` buyruqlarini real tarmoqda ehtiyotkorlik bilan ishlating, chunki CPUga yuk berishi mumkin.

## Muammolar va yechimlar

| Muammo | Ehtimoliy sabab | Tekshiruv |
|---|---|---|
| Client IP olmayapti | `ip helper-address` noto'g'ri interfeysda | Gateway interfeys konfiguratsiyasi |
| Serverga paket yetmayapti | Routing yo'q | `show ip route`, `ping` |
| IP noto'g'ri subnetdan berilyapti | DHCP pool yoki gateway noto'g'ri | DHCP server pool sozlamalari |
| VLAN clientlari umuman ishlamayapti | SVI down yoki trunk muammosi | `show vlan brief`, `show interfaces trunk` |
| DHCP Offer qaytmayapti | ACL bloklayapti | `show access-lists` |

## Keng tarqalgan xatolar

- `ip helper-address`ni DHCP server joylashgan interfeysga yozish. To'g'ri joy: client gateway interfeysi.
- DHCP serverda kerakli subnet uchun pool yaratmaslik.
- DHCP serverdan client subnetga qaytish marshruti yo'qligi.
- VLAN trunkda kerakli VLAN o'tkazilmaganligi.

## Q&A

**Savol:** DHCP server har bir VLANda bo'lishi shartmi?  
**Javob:** Yo'q. Bitta markaziy DHCP server ko'p VLANlarga xizmat qilishi mumkin, lekin har bir client VLAN gatewayida `ip helper-address` kerak bo'ladi.

**Savol:** Ikkita DHCP server ko'rsatish mumkinmi?  
**Javob:** Ha.

```cisco
interface vlan 10
 ip helper-address 192.168.100.10
 ip helper-address 192.168.100.11
```

**Savol:** `ip helper-address` DHCP relaymi?  
**Javob:** Ha, Cisco IOSda DHCP relayni sozlashning eng ko'p ishlatiladigan usuli shu.

