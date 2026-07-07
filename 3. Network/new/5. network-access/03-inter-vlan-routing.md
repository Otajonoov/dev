# Inter-VLAN Routing

VLANlar Layer 2 broadcast domenlarni ajratadi. VLAN 10 va VLAN 20 bir-biri bilan gaplashishi uchun Layer 3 qurilma kerak: router yoki Layer 3 switch. Shu jarayon inter-VLAN routing deyiladi.

## 3 ta keng tarqalgan usul

1. Har VLAN uchun alohida router interface - eski va kam ishlatiladi.
2. Router-on-a-stick - bitta router interface, bir nechta subinterface.
3. Layer 3 switch SVI - katta va tezkor campus tarmoqlarda keng ishlatiladi.

## Router-on-a-stick

Topologiya:

```text
PC1 -- SW1 ==trunk== R1
PC2 -- SW1

VLAN 10 -> 192.168.10.0/24
VLAN 20 -> 192.168.20.0/24
```

Switch:

```cisco
configure terminal
vlan 10
 name USERS
vlan 20
 name SALES

interface fastEthernet0/1
 description PC-VLAN10
 switchport mode access
 switchport access vlan 10

interface fastEthernet0/2
 description PC-VLAN20
 switchport mode access
 switchport access vlan 20

interface gigabitEthernet0/1
 description TRUNK-TO-R1
 switchport mode trunk
 switchport trunk allowed vlan 10,20
 no shutdown
end
```

Router:

```cisco
configure terminal
interface gigabitEthernet0/0
 description TO-SW1
 no shutdown

interface gigabitEthernet0/0.10
 description VLAN10-GATEWAY
 encapsulation dot1Q 10
 ip address 192.168.10.1 255.255.255.0

interface gigabitEthernet0/0.20
 description VLAN20-GATEWAY
 encapsulation dot1Q 20
 ip address 192.168.20.1 255.255.255.0
end
```

PC sozlamalari:

```text
PC1:
IP 192.168.10.10
Mask 255.255.255.0
Gateway 192.168.10.1

PC2:
IP 192.168.20.10
Mask 255.255.255.0
Gateway 192.168.20.1
```

## Native VLAN bilan subinterface

Agar native VLAN kerak bo'lsa:

```cisco
interface gigabitEthernet0/0.999
 encapsulation dot1Q 999 native
 ip address 192.168.99.1 255.255.255.0
```

Lekin CCNA lablarda ko'pincha data VLANlar taglangan subinterfacelar orqali ishlatiladi.

## Layer 3 switch orqali SVI

SVI (Switched Virtual Interface) - VLAN uchun virtual Layer 3 interface.

```cisco
configure terminal
ip routing

vlan 10
 name USERS
vlan 20
 name SALES

interface vlan 10
 description VLAN10-GATEWAY
 ip address 192.168.10.1 255.255.255.0
 no shutdown

interface vlan 20
 description VLAN20-GATEWAY
 ip address 192.168.20.1 255.255.255.0
 no shutdown
end
```

Muhim: L3 switchda `ip routing` yoqilmasa, SVIlar gateway kabi ishlashi cheklanishi mumkin.

Access portlar:

```cisco
interface fastEthernet0/1
 switchport mode access
 switchport access vlan 10

interface fastEthernet0/2
 switchport mode access
 switchport access vlan 20
```

## Routed port

Layer 3 switchda portni switchport emas, routed interface qilish mumkin:

```cisco
interface gigabitEthernet0/1
 no switchport
 ip address 10.0.0.1 255.255.255.252
 no shutdown
```

Bu odatda distribution/core yoki routerga ulanishda ishlatiladi.

## Tekshiruv komandalar

Router:

```cisco
show ip interface brief
show ip route
show running-config interface gigabitEthernet0/0.10
show arp
ping 192.168.20.10
```

Switch:

```cisco
show vlan brief
show interfaces trunk
show mac address-table
show ip interface brief
```

Layer 3 switch:

```cisco
show ip route
show ip interface brief
show running-config | include ip routing
show interfaces vlan 10
```

## Troubleshooting

Muammo: Bir VLAN ichida ping ishlaydi, boshqa VLANga ping ishlamaydi.

Tekshiring:

- PC default gateway to'g'rimi?
- Gateway interface `up/up` holatidami?
- Switch-router link trunkmi?
- Trunk allowed listda kerakli VLAN bormi?
- Router subinterface `encapsulation dot1Q` VLAN IDsi to'g'rimi?
- L3 switchda `ip routing` yoqilganmi?
- ACL trafikni bloklamayaptimi?

Router-on-a-stickda klassik xato:

```cisco
interface gigabitEthernet0/0
 shutdown
```

Subinterfacelar sozlangan bo'lsa ham, fizik interface shutdown bo'lsa hammasi ishlamaydi.

## Common mistakes

- PC gatewayi noto'g'ri VLAN gatewayiga berilgan.
- Router subinterface VLAN IDsi switch access VLAN bilan mos emas.
- Switch-router port access bo'lib qolgan.
- L3 switchda VLAN yaratilmagan, lekin SVI sozlangan.
- SVI `down/down`, chunki shu VLANda active port yo'q.
- `ip routing` unutib qoldirilgan.
- Bir subnet ikki VLANda ishlatilgan.

## Qisqa Q&A

**Savol:** Har VLANga gateway kerakmi?  
**Javob:** Agar VLANdan boshqa tarmoqqa chiqish kerak bo'lsa, ha.

**Savol:** Router-on-a-stickda routerga nechta kabel kerak?  
**Javob:** Odatda bitta trunk kabel yetadi.

**Savol:** SVI up bo'lishi uchun nima kerak?  
**Javob:** VLAN mavjud bo'lishi va shu VLANda kamida bitta active Layer 2 port yoki trunkda active VLAN bo'lishi kerak.
