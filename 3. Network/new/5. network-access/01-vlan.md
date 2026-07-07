# VLAN

VLAN (Virtual LAN) - bitta fizik switch ichida alohida Layer 2 broadcast domenlar yaratish usuli. Masalan, bitta switchda buxgalteriya, savdo va menejment kompyuterlari ulangan bo'lsa, ularni VLAN 10, VLAN 20, VLAN 99 qilib ajratish mumkin.

## VLAN nima qiladi?

VLAN quyidagilarni beradi:

- Broadcast trafikni chegaralaydi.
- Foydalanuvchilarni bo'lim yoki vazifa bo'yicha ajratadi.
- Xavfsizlikni yaxshilaydi, chunki turli VLANlar bevosita Layer 2 orqali gaplashmaydi.
- IP subnetlarni tartibli joylashtirishga yordam beradi.

Muhim: VLANning o'zi routing qilmaydi. VLAN 10 dagi PC VLAN 20 dagi PC bilan gaplashishi uchun router yoki Layer 3 switch kerak.

## Access port

Access port odatda bitta endpoint uchun ishlatiladi: PC, printer, kamera, IP phone. Access port faqat bitta data VLANga tegishli bo'ladi.

```cisco
configure terminal
vlan 10
 name USERS
vlan 20
 name SALES
vlan 99
 name MGMT

interface fastEthernet0/1
 description PC-1
 switchport mode access
 switchport access vlan 10
 no shutdown

interface fastEthernet0/2
 description PC-2
 switchport mode access
 switchport access vlan 20
 no shutdown
end
```

## Voice VLAN

IP phone va PC bir portga ulanadigan holatda voice VLAN ishlatiladi. Telefon voice VLANdan, telefon ortidagi PC esa access VLANdan foydalanadi.

```cisco
interface fastEthernet0/10
 description IP-PHONE-AND-PC
 switchport mode access
 switchport access vlan 10
 switchport voice vlan 30
 spanning-tree portfast
 no shutdown
```

## VLAN database va tekshiruv

```cisco
show vlan brief
show interfaces status
show running-config interface fastEthernet0/1
show mac address-table vlan 10
```

`show vlan brief` natijasida VLANlar va ularga tegishli access portlar ko'rinadi. Trunk portlar bu jadvalda access port kabi chiqmasligi mumkin.

## Default VLAN

Cisco switchlarda odatda barcha portlar boshlang'ich holatda VLAN 1 ga tegishli. VLAN 1 ni o'chirib bo'lmaydi, lekin foydalanuvchi trafikini VLAN 1 da qoldirmaslik tavsiya qilinadi.

Yaxshi amaliyot:

```cisco
vlan 999
 name UNUSED

interface range fastEthernet0/3 - 24
 description UNUSED-PORT
 switchport mode access
 switchport access vlan 999
 shutdown
```

## VLAN va IP subnet

Odatda har bir VLANga alohida IP subnet beriladi:

```text
VLAN 10 USERS  -> 192.168.10.0/24
VLAN 20 SALES  -> 192.168.20.0/24
VLAN 99 MGMT   -> 192.168.99.0/24
```

PC konfiguratsiyasi:

```text
PC in VLAN 10
IP: 192.168.10.10
Mask: 255.255.255.0
Gateway: 192.168.10.1
```

Gateway odatda router subinterface yoki L3 switch SVI manzili bo'ladi.

## Troubleshooting

Muammo: PC bir VLAN ichidagi boshqa PCga ping qilmayapti.

Tekshirish:

```cisco
show vlan brief
show interfaces fastEthernet0/1 switchport
show interfaces fastEthernet0/1 status
show mac address-table interface fastEthernet0/1
```

Endpoint tomonda:

```text
IP manzil to'g'rimi?
Subnet mask to'g'rimi?
PC kabeli to'g'ri portgami?
Port shutdown emasmi?
```

Muammo: Port VLANga qo'yilgan, lekin VLAN yo'q.

```cisco
show vlan brief
configure terminal
vlan 10
 name USERS
end
```

Ba'zi switchlarda port mavjud bo'lmagan VLANga qo'yilsa, VLAN avtomatik yaratilmaydi yoki noto'g'ri holat yuzaga keladi. Shuning uchun avval VLAN yaratish yaxshi.

## Common mistakes

- VLAN 1 ni hamma narsa uchun ishlatish.
- Portga `switchport access vlan 10` berib, `switchport mode access` ni unutish.
- VLAN nomini yaratib, portni boshqa VLANga qo'yish.
- PC gatewayini boshqa VLANdagi IPga sozlash.
- Trunk portni `show vlan brief`da access port sifatida izlash.
- VLAN bor, lekin uplink trunkda bu VLANga ruxsat berilmagan.

## Qisqa Q&A

**Savol:** VLAN 10 dagi PC VLAN 20 dagi PCga ping qilishi kerakmi?  
**Javob:** Faqat inter-VLAN routing sozlangan bo'lsa.

**Savol:** VLAN nomi trafikga ta'sir qiladimi?  
**Javob:** Yo'q. Nom faqat administratorga qulaylik beradi.

**Savol:** Bitta access port bir nechta data VLANda bo'la oladimi?  
**Javob:** Yo'q. Access port bitta data VLANda bo'ladi. Voice VLAN alohida maxsus holat.
