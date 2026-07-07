# 05. Layer 2 Security

Layer 2 hujumlari ko'pincha LAN ichida sodir bo'ladi. Switch portlari ochiq qolsa, attacker soxta DHCP server, ARP spoofing, MAC flooding yoki trunk/VLAN hujumlari orqali tarmoqqa zarar yetkazishi mumkin.

## Port security

Port security access portda ruxsat etilgan MAC manzillarni cheklaydi.

Oddiy sozlash:

```cisco
conf t
interface f0/10
 switchport mode access
 switchport access vlan 10
 switchport port-security
 switchport port-security maximum 1
 switchport port-security mac-address sticky
 switchport port-security violation shutdown
end
```

Violation mode:

- **shutdown** - port err-disabled bo'ladi. Eng qat'iy va ko'p ishlatiladi.
- **restrict** - trafik bloklanadi, log/counter yoziladi.
- **protect** - trafik bloklanadi, lekin log kamroq.

Tekshirish:

```cisco
show port-security
show port-security interface f0/10
show interfaces status err-disabled
show mac address-table interface f0/10
```

Err-disabled portni tiklash:

```cisco
conf t
interface f0/10
 shutdown
 no shutdown
end
```

Avtomatik recovery:

```cisco
conf t
errdisable recovery cause psecure-violation
errdisable recovery interval 300
end
```

## DHCP snooping

DHCP snooping soxta DHCP serverdan himoya qiladi. Switch portlarni trusted va untrustedga ajratadi.

- **Trusted port** - haqiqiy DHCP server yoki uplink tomoni.
- **Untrusted port** - foydalanuvchi portlari. Bu portdan DHCP offer/ack kelmasligi kerak.

Sozlash:

```cisco
conf t
ip dhcp snooping
ip dhcp snooping vlan 10,20

interface g0/1
 description Uplink-to-DHCP-server
 ip dhcp snooping trust

interface range f0/1 - 24
 ip dhcp snooping limit rate 10
end
```

Tekshirish:

```cisco
show ip dhcp snooping
show ip dhcp snooping binding
show ip dhcp snooping statistics
```

Muhim: DHCP snooping binding table DAI uchun ham kerak bo'ladi.

## Dynamic ARP Inspection

DAI ARP spoofing hujumlarini kamaytiradi. Switch ARP paketlarini DHCP snooping binding table bilan solishtiradi.

Sozlash:

```cisco
conf t
ip dhcp snooping
ip dhcp snooping vlan 10
ip arp inspection vlan 10

interface g0/1
 description Uplink
 ip dhcp snooping trust
 ip arp inspection trust

interface range f0/1 - 24
 ip arp inspection limit rate 15
end
```

Tekshirish:

```cisco
show ip arp inspection
show ip arp inspection vlan 10
show ip arp inspection statistics
show ip dhcp snooping binding
```

Statik IP ishlatadigan hostlar uchun ARP ACL kerak bo'lishi mumkin:

```cisco
conf t
arp access-list STATIC_HOSTS
 permit ip host 192.168.10.50 mac host 0011.2233.4455
ip arp inspection filter STATIC_HOSTS vlan 10
end
```

## VLAN hoppingdan himoya

Access portlarda trunk negotiationni o'chiring:

```cisco
conf t
interface range f0/1 - 24
 switchport mode access
 switchport nonegotiate
 spanning-tree portfast
end
```

Trunk portlarda native VLANni foydalanuvchi VLANdan alohida qiling:

```cisco
conf t
vlan 999
 name NATIVE_UNUSED
interface g0/1
 switchport mode trunk
 switchport trunk native vlan 999
 switchport trunk allowed vlan 10,20,30
 switchport nonegotiate
end
```

Unused portlarni yopish:

```cisco
conf t
vlan 998
 name UNUSED_PORTS
interface range f0/20 - 24
 switchport mode access
 switchport access vlan 998
 shutdown
end
```

## STP security

### PortFast

End-device portlar uchun:

```cisco
conf t
interface range f0/1 - 24
 spanning-tree portfast
end
```

PortFastni switch-to-switch portda ishlatmang.

### BPDU Guard

PortFast portga BPDU kelsa portni err-disabled qiladi.

```cisco
conf t
spanning-tree portfast default
spanning-tree portfast bpduguard default
end
```

Yoki bitta interfeysda:

```cisco
conf t
interface f0/10
 spanning-tree bpduguard enable
end
```

### Root Guard

Kutilmagan switch root bridge bo'lib qolmasin:

```cisco
conf t
interface g0/2
 spanning-tree guard root
end
```

Tekshirish:

```cisco
show spanning-tree
show spanning-tree summary
show interfaces status err-disabled
```

## Common mistakes

- DHCP serverga boradigan uplinkni trusted qilmaslik.
- DHCP snoopingni global yoqib, VLANda yoqmaslik.
- DAI yoqilgan, lekin DHCP snooping binding table bo'sh.
- Statik IP hostlar uchun ARP ACL qo'shmaslik.
- Port security trunk portda ishlatishga urinish.
- PortFastni switchlar orasidagi linkda yoqish.
- Native VLANni default VLAN 1 holida qoldirish.

## Q&A

**Savol:** DHCP snooping nima uchun kerak?

**Javob:** Soxta DHCP server foydalanuvchilarga noto'g'ri gateway/DNS berishini to'xtatadi.

**Savol:** DAI DHCP snoopingsiz ishlaydimi?

**Javob:** Odatda DHCP snooping binding tablega tayanadi. Statik hostlar uchun ARP ACL bilan ishlatish mumkin.

**Savol:** Port security MAC spoofingni 100% to'xtatadimi?

**Javob:** Yo'q, lekin portdagi ruxsat etilgan MAClar sonini cheklab, oddiy MAC flooding va noto'g'ri ulanishlarni kamaytiradi.
