# Cisco NAT: Static NAT, Dynamic NAT, PAT va Port Forwarding

NAT (Network Address Translation) IP manzilni boshqa IP manzilga o'zgartiradi. Eng ko'p holatda private IP manzillar internetga chiqish uchun public IPga tarjima qilinadi.

## NAT atamalari

- Inside local: ichki qurilmaning haqiqiy private IP manzili. Masalan, `192.168.1.10`.
- Inside global: tashqaridan ko'rinadigan public IP manzil. Masalan, `203.0.113.10`.
- Outside local/global: tashqi host manzili. CCNAda asosan inside local va inside global muhim.

## NAT interfeys rollari

NAT ishlashi uchun router interfeyslariga rol beriladi:

```cisco
interface gigabitEthernet0/0
 description LAN
 ip address 192.168.1.1 255.255.255.0
 ip nat inside

interface gigabitEthernet0/1
 description ISP
 ip address 203.0.113.2 255.255.255.252
 ip nat outside
```

## PAT (NAT Overload)

PAT ko'p private hostlarni bitta public IP orqali internetga chiqaradi. Port raqamlari yordamida sessiyalar ajratiladi.

```cisco
conf t
access-list 10 permit 192.168.1.0 0.0.0.255

interface gigabitEthernet0/0
 ip nat inside

interface gigabitEthernet0/1
 ip nat outside

ip nat inside source list 10 interface gigabitEthernet0/1 overload
end
```

Default route ham kerak bo'ladi:

```cisco
ip route 0.0.0.0 0.0.0.0 203.0.113.1
```

## Dynamic NAT pool

Dynamic NAT ichki hostlarni public IP pooldan vaqtincha manzilga bog'laydi. PATdan farqli ravishda bitta public IP bir vaqtda bitta inside hostga biriktiriladi.

```cisco
conf t
access-list 20 permit 192.168.10.0 0.0.0.255
ip nat pool PUBLIC_POOL 203.0.113.10 203.0.113.20 netmask 255.255.255.0
ip nat inside source list 20 pool PUBLIC_POOL
end
```

## Static NAT

Static NAT bitta private IPni doimiy bitta public IPga bog'laydi. Serverlarni tashqaridan ko'rsatishda ishlatiladi.

```cisco
conf t
ip nat inside source static 192.168.1.50 203.0.113.50
end
```

## Port forwarding

Port forwarding bitta public IPdagi ma'lum portni ichki serverga uzatadi. Masalan, public `203.0.113.2:8080` ichki web server `192.168.1.50:80`ga:

```cisco
conf t
ip nat inside source static tcp 192.168.1.50 80 203.0.113.2 8080
end
```

SSH port forwarding:

```cisco
ip nat inside source static tcp 192.168.1.60 22 203.0.113.2 2222
```

Tashqaridan ulanish:

```bash
ssh -p 2222 admin@203.0.113.2
```

## Tekshiruv buyruqlari

```cisco
show ip nat translations
show ip nat statistics
show running-config | include ip nat
show access-lists
show ip route
clear ip nat translation *
debug ip nat
```

`show ip nat translations` namunasi:

```text
Pro  Inside global      Inside local       Outside local      Outside global
tcp  203.0.113.2:49152  192.168.1.10:49152 198.51.100.5:443  198.51.100.5:443
```

## Troubleshooting tartibi

1. LAN host gatewayga ping qila oladimi?
2. Router ISP next-hopga ping qila oladimi?
3. Default route bormi?
4. `ip nat inside` va `ip nat outside` to'g'ri interfeyslardami?
5. ACL ichki subnetni match qilyaptimi?
6. NAT translation yaralyaptimi?
7. ISP yoki upstream ACL trafikni bloklamayaptimi?

## Keng tarqalgan xatolar

- Inside va outside rollarini almashtirib yuborish.
- NAT ACLda wildcard maskni noto'g'ri yozish.
- Default route yo'qligi.
- Static NAT uchun public IP upstreamda routerga yo'naltirilmaganligi.
- Port forwardingda ichki server firewalli portni yopib qo'yganligi.

## Q&A

**Savol:** PAT va NAT bir xilmi?  
**Javob:** PAT NATning bir turi. PATda ko'p ichki hostlar bitta yoki kam sonli public IPni portlar orqali bo'lishadi.

**Savol:** NAT routingni almashtiradimi?  
**Javob:** Yo'q. NAT manzilni o'zgartiradi, lekin paket qayerga borishini routing hal qiladi.

**Savol:** NAT xavfsizlik devorimi?  
**Javob:** To'liq firewall emas. U ichki manzillarni yashiradi, lekin security policy uchun ACL yoki firewall kerak.

