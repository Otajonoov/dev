# EtherChannel va LACP

EtherChannel - bir nechta fizik linkni bitta mantiqiy link sifatida ishlatish. Switchlar orasida ikki yoki undan ortiq kabel bo'lsa, STP ularning ayrimlarini bloklab qo'yishi mumkin. EtherChannel esa ularni bitta Port-channel qilib, bandwidth va redundancy beradi.

## Nima beradi?

- Bir nechta linkdan parallel foydalanish.
- STP bitta logical linkni ko'radi.
- Bitta link uzilsa, kanal qolgan linklar bilan ishlayveradi.
- Konfiguratsiya Port-channel interfeysida markazlashadi.

Muhim: EtherChannel bitta flow uchun bandwidthni "qo'shib" yubormaydi. Load-balancing odatda source/destination MAC, IP yoki port asosida flowlarni linklarga taqsimlaydi.

## LACP nima?

LACP (Link Aggregation Control Protocol) - IEEE 802.3ad standart protokoli. U EtherChannelni kelishib hosil qiladi.

LACP rejimlari:

| Mode | Ma'nosi |
|---|---|
| `active` | LACPni faol boshlaydi |
| `passive` | LACPga javob beradi |
| `on` | Protokolsiz majburiy EtherChannel |

Ishlaydigan kombinatsiyalar:

```text
active + active  -> ishlaydi
active + passive -> ishlaydi
passive + passive -> ishlamaydi
on + on -> ishlaydi, lekin LACP yo'q
```

CCNA uchun odatda `active` ishlatish eng aniq va xavfsiz.

## Trunk EtherChannel sozlash

SW1:

```cisco
configure terminal
interface range gigabitEthernet0/1 - 2
 description TO-SW2-EC1
 switchport mode trunk
 switchport trunk allowed vlan 10,20,99
 channel-group 1 mode active
 no shutdown

interface port-channel1
 description LACP-TO-SW2
 switchport mode trunk
 switchport trunk allowed vlan 10,20,99
 no shutdown
end
```

SW2:

```cisco
configure terminal
interface range gigabitEthernet0/1 - 2
 description TO-SW1-EC1
 switchport mode trunk
 switchport trunk allowed vlan 10,20,99
 channel-group 1 mode active
 no shutdown

interface port-channel1
 description LACP-TO-SW1
 switchport mode trunk
 switchport trunk allowed vlan 10,20,99
 no shutdown
end
```

## Access EtherChannel

Ba'zida serverga ikki link access VLANda ulanadi:

```cisco
interface range gigabitEthernet0/3 - 4
 description SERVER-LACP
 switchport mode access
 switchport access vlan 50
 channel-group 2 mode active

interface port-channel2
 description SERVER-BOND
 switchport mode access
 switchport access vlan 50
```

Server tomonda ham LACP bond/team sozlangan bo'lishi kerak.

## Moslik talablari

EtherChannel a'zolarida odatda quyidagilar bir xil bo'lishi kerak:

- Speed.
- Duplex.
- Access yoki trunk mode.
- Access VLAN.
- Native VLAN.
- Allowed VLAN list.
- STP parametrlar.

Amaliy qoida: kerakli Layer 2 sozlamalarni Port-channel interfeysida bering, a'zo portlarda esa kanalga qo'shishdan oldingi asosiy moslikni tekshiring.

## Tekshiruv komandalar

```cisco
show etherchannel summary
show etherchannel detail
show interfaces port-channel1
show interfaces trunk
show spanning-tree interface port-channel1
show running-config interface port-channel1
```

`show etherchannel summary`dagi belgilar:

```text
P - bundled in port-channel
I - stand-alone
S - Layer2
R - Layer3
U - in use
D - down
```

Yaxshi holatga misol:

```text
Group  Port-channel  Protocol  Ports
1      Po1(SU)       LACP      Gi0/1(P) Gi0/2(P)
```

## Troubleshooting

Muammo: Port-channel ishlamayapti.

```cisco
show etherchannel summary
show etherchannel detail
show interfaces status
show running-config interface gigabitEthernet0/1
show running-config interface gigabitEthernet0/2
show running-config interface port-channel1
```

Tekshiring:

- Ikkala tomonda channel-group raqami mahalliy, lekin mode mosmi?
- `active/passive` kombinatsiyasi ishlaydimi?
- A'zo portlar speed/duplex bo'yicha mosmi?
- Biri access, biri trunk emasmi?
- Allowed VLAN listlar bir xilmi?
- Portlar shutdown emasmi?

Muammo: Bir port `I` holatda.

Bu port channelga bundling bo'lmaganini bildiradi. Sabab ko'pincha konfiguratsiya mismatch.

Tuzatishdan oldin konfiguratsiyani solishtiring, keyin portni qayta qo'shing:

```cisco
interface gigabitEthernet0/2
 no channel-group 1
 channel-group 1 mode active
```

## Common mistakes

- Bir tomonda `active`, ikkinchi tomonda `on`.
- `passive + passive` qilib qo'yish.
- Port-channelda trunk allowed VLANni berib, fizik portlarda boshqacha qoldirish.
- A'zo portlardan birida native VLAN farq qilishi.
- EtherChannel tuzilmasdan parallel kabellarni ulab, STP bloklashini unutish.
- Port-channel konfiguratsiyasini emas, faqat fizik portlarni tekshirish.

## Qisqa Q&A

**Savol:** EtherChannel STPni o'chiradimi?  
**Javob:** Yo'q. STP Port-channelni bitta logical link sifatida ko'radi.

**Savol:** 2 ta 1 Gbps link doim bitta file transferga 2 Gbps beradimi?  
**Javob:** Odatda yo'q. Load-balancing flowlarga qarab taqsimlanadi.

**Savol:** LACP va PAgP farqi nima?  
**Javob:** LACP ochiq standart, PAgP Cisco proprietary. CCNAda LACP ko'proq tavsiya qilinadi.
