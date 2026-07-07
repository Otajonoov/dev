# STP va Rapid PVST+

STP (Spanning Tree Protocol) Layer 2 looplarni oldini oladi. Switchlar orasida redundant linklar bo'lsa, Ethernet frame aylanib qolishi, broadcast storm yuzaga kelishi va MAC jadval beqarorlashishi mumkin. STP bitta loop-free logical topologiya yaratadi.

## Nega STP kerak?

Ethernetda TTL yo'q. Agar Layer 2 loop bo'lsa, broadcast frame tarmoqda qayta-qayta aylanishi mumkin. Natija:

- Broadcast storm.
- Switch CPU yuklanishi.
- MAC address table flapping.
- Foydalanuvchilarda paket yo'qolishi.
- Tarmoq "ishlayapti, lekin juda sekin" holati.

## Asosiy terminlar

- Root bridge - STP topologiyaning markazi.
- Bridge ID - priority + MAC address.
- Root port - root bridge tomonga eng yaxshi port.
- Designated port - segmentdagi forwarding port.
- Alternate port - backup, odatda blocking/discarding.
- BPDU - switchlar STP ma'lumot almashadigan frame.

Root bridge eng kichik Bridge IDga ega switch bo'ladi. Priority default 32768. Cisco PVST/Rapid PVSTda priority VLAN bo'yicha belgilanadi.

## Rapid PVST+

Rapid PVST+ - Cisco implementatsiyasi, har VLAN uchun alohida rapid spanning-tree ishlatadi. U klassik STPga qaraganda tezroq convergenceni beradi.

Yoqish:

```cisco
configure terminal
spanning-tree mode rapid-pvst
end
```

## Root bridge sozlash

Rootni tasodifga tashlamang. Distribution switchni root qiling.

SW1 primary:

```cisco
configure terminal
spanning-tree vlan 10,20,99 root primary
end
```

SW2 secondary:

```cisco
configure terminal
spanning-tree vlan 10,20,99 root secondary
end
```

Yoki priority qo'lda:

```cisco
spanning-tree vlan 10 priority 4096
spanning-tree vlan 20 priority 4096
```

Priority 4096 qadam bilan beriladi: 0, 4096, 8192, 12288 va hokazo.

## PortFast

PortFast endpoint portlari uchun. U portni tez forwardingga o'tkazadi. Switch-switch trunk portga PortFast bermang, agar aniq sabab bo'lmasa.

```cisco
interface fastEthernet0/1
 description PC-PORT
 switchport mode access
 switchport access vlan 10
 spanning-tree portfast
```

Global:

```cisco
spanning-tree portfast default
```

Bu odatda access portlarga ta'sir qiladi.

## BPDU Guard

BPDU Guard PortFast portga switch ulanib qolsa portni err-disable qiladi. Bu access layerda juda foydali himoya.

```cisco
interface fastEthernet0/1
 spanning-tree bpduguard enable
```

Global:

```cisco
spanning-tree portfast bpduguard default
```

Err-disable portni ko'rish:

```cisco
show interfaces status err-disabled
show errdisable recovery
```

Qo'lda tiklash:

```cisco
interface fastEthernet0/1
 shutdown
 no shutdown
```

## Tekshiruv komandalar

```cisco
show spanning-tree
show spanning-tree vlan 10
show spanning-tree root
show spanning-tree blockedports
show interfaces status
show mac address-table dynamic
```

`show spanning-tree vlan 10`da tekshiring:

- Root bridge qaysi switch?
- Local switch rootmi?
- Port role: Root, Desg, Altn.
- Port state: FWD, BLK yoki discarding.
- Cost qiymatlari.

## Troubleshooting

Muammo: Trunk bor, lekin bitta VLAN trafik o'tmayapti.

```cisco
show interfaces trunk
show spanning-tree vlan 20
show spanning-tree blockedports
```

Ehtimol STP VLAN 20 uchun portni blocking qilgan. Bu har doim xato emas; STP loopni oldini olayotgan bo'lishi mumkin.

Muammo: Root switch noto'g'ri joyda.

```cisco
show spanning-tree root
show spanning-tree vlan 10
```

Tuzatish:

```cisco
spanning-tree vlan 10 root primary
```

Muammo: Access port err-disabled.

```cisco
show interfaces status err-disabled
show logging | include BPDU|ERR
```

Agar BPDU Guard sabab bo'lsa, demak shu access portga switch yoki noto'g'ri qurilma ulangan bo'lishi mumkin.

## Common mistakes

- STPni o'chirib qo'yish.
- Root bridge tanlashni switchlarga qoldirish.
- PortFastni switch-switch portda ishlatish.
- BPDU Guardni yoqmaslik.
- Blocking portni avtomatik "muammo" deb o'ylash.
- VLANlar uchun root placementni tekshirmaslik.
- EtherChannel tuzmasdan parallel linklarni ulash va STP nima qilayotganini tushunmaslik.

## Qisqa Q&A

**Savol:** STP blocking port trafik o'tkazmaydimi?  
**Javob:** Oddiy data trafikni o'tkazmaydi, lekin BPDUlarni tinglaydi.

**Savol:** Rapid PVST+ har VLAN uchun alohida ishlaydimi?  
**Javob:** Ha, Cisco Rapid PVST+ har VLAN uchun alohida STP instance ishlatadi.

**Savol:** Root bridge doim eng kuchli switch bo'lishi kerakmi?  
**Javob:** Odatda distribution/core kabi markaziy switch root bo'lishi kerak. Maqsad trafik yo'lini boshqarish.
