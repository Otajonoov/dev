# CDP va LLDP

CDP va LLDP qo'shni qurilmalarni aniqlash uchun ishlatiladi. Ular switch, router, IP phone yoki access point qaysi portga ulanganini tez bilishga yordam beradi.

## CDP nima?

CDP (Cisco Discovery Protocol) - Cisco proprietary protokoli. Cisco qurilmalarda odatda yoqilgan bo'ladi. CDP Layer 2 da ishlaydi va IP sozlanmagan bo'lsa ham qo'shni Cisco qurilma haqida ma'lumot berishi mumkin.

CDP ko'rsatishi mumkin:

- Qo'shni qurilma nomi.
- Local interface.
- Qo'shni interface.
- Platforma.
- Capability: router, switch, phone.
- Management IP.
- IOS versiya.

## LLDP nima?

LLDP (Link Layer Discovery Protocol) - IEEE 802.1AB ochiq standart. Turli vendor qurilmalari orasida foydali: Cisco, Juniper, HP/Aruba, Linux server, IP phone va boshqalar.

Cisco IOSda LLDP ba'zida default o'chirilgan bo'lishi mumkin:

```cisco
configure terminal
lldp run
end
```

## CDP komandalar

```cisco
show cdp
show cdp neighbors
show cdp neighbors detail
show cdp interface
show cdp traffic
```

Interface bo'yicha:

```cisco
show cdp neighbors gigabitEthernet0/1 detail
```

CDPni global o'chirish:

```cisco
no cdp run
```

Faqat bitta portda o'chirish:

```cisco
interface gigabitEthernet0/10
 no cdp enable
```

## LLDP komandalar

```cisco
show lldp
show lldp neighbors
show lldp neighbors detail
show lldp interface
show lldp traffic
```

LLDPni portda sozlash:

```cisco
interface gigabitEthernet0/1
 lldp transmit
 lldp receive
```

Portda o'chirish:

```cisco
interface gigabitEthernet0/10
 no lldp transmit
 no lldp receive
```

## Amaliy ishlatish

Topologiyani tez tekshirish:

```cisco
show cdp neighbors
show lldp neighbors
```

Misol:

```text
Device ID        Local Intrfce     Capability  Platform  Port ID
SW2              Gig 0/1           S           C2960     Gig 0/1
R1               Gig 0/2           R           ISR4321   Gig 0/0
```

Bu natijadan bilamiz:

- SW1 Gi0/1 porti SW2 Gi0/1 ga ulangan.
- SW1 Gi0/2 porti R1 Gi0/0 ga ulangan.

## IP phone va voice VLAN

CDP/LLDP-MED IP phonega voice VLAN haqida ma'lumot berishda ishlatilishi mumkin. Cisco IP phone odatda CDP orqali voice VLANni biladi.

```cisco
interface fastEthernet0/10
 description IP-PHONE
 switchport mode access
 switchport access vlan 10
 switchport voice vlan 30
 spanning-tree portfast
```

Phone CDP/LLDP orqali VLAN 30 voice VLAN ekanini bilishi mumkin.

## Xavfsizlik

CDP/LLDP foydali, lekin ma'lumot oshkor qiladi. Masalan, qurilma nomi, platforma, IOS versiya va management IP ko'rinishi mumkin.

Tavsiya:

- Ichki trunk/uplinklarda yoqilgan bo'lishi mumkin.
- User-facing ishonchsiz portlarda o'chirish mumkin.
- Internet, ISP, guest yoki third-party portlarda o'chirish yaxshi amaliyot.

```cisco
interface gigabitEthernet0/24
 description ISP-HANDOFF
 no cdp enable
 no lldp transmit
 no lldp receive
```

## Troubleshooting

Muammo: Qo'shni qurilma ko'rinmayapti.

```cisco
show cdp
show cdp interface gigabitEthernet0/1
show lldp
show lldp interface gigabitEthernet0/1
show interfaces gigabitEthernet0/1 status
show interfaces gigabitEthernet0/1
```

Tekshiring:

- Interface up/upmi?
- CDP yoki LLDP global yoqilganmi?
- Portda protokol o'chirilmaganmi?
- Qo'shni qurilma CDP/LLDP qo'llaydimi?
- Firewall yoki L2 filtering yo'qmi?

Muammo: Native VLAN mismatch xabari.

CDP ba'zan trunkdagi native VLAN mismatchni logda ko'rsatadi:

```text
%CDP-4-NATIVE_VLAN_MISMATCH
```

Tekshirish:

```cisco
show interfaces trunk
show running-config interface gigabitEthernet0/1
```

## Common mistakes

- CDP/LLDP natijasini routing jadvali deb o'ylash.
- Qo'shni IP ko'rinadi, demak ping ishlashi kerak deb o'ylash.
- Ishonchsiz portlarda CDP/LLDPni yoqib qo'yish.
- LLDP global yoqilmaganini unutish.
- Packet Tracerda hamma LLDP funksiyalari real IOS kabi ishlamasligini hisobga olmaslik.

## Qisqa Q&A

**Savol:** CDP IP kerak bo'lmasdan ishlaydimi?  
**Javob:** Ha, u Layer 2 discovery protokoli.

**Savol:** LLDP Cisco qurilmalarida ishlaydimi?  
**Javob:** Ha, lekin ba'zi IOSlarda `lldp run` bilan yoqish kerak.

**Savol:** CDP xavflimi?  
**Javob:** Ichki tarmoqda foydali, lekin ishonchsiz portlarda qurilma ma'lumotlarini oshkor qilishi mumkin.
