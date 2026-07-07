# Lab topologies

Bu fayl `new` katalogidagi mavzular uchun amaliy lab g'oyalarini beradi.

## Lab 1: VLAN va trunk

Topology:

```text
PC1 -- SW1 ===== SW2 -- PC2
       | trunk |
PC3 -- SW1     SW2 -- PC4
```

VLANlar:

```text
VLAN 10: PC1, PC2
VLAN 20: PC3, PC4
```

Vazifa:

- SW1 va SW2 da VLAN 10/20 yarat.
- PC portlarini access mode'ga o'tkaz.
- SW1-SW2 linkni trunk qil.
- Faqat VLAN 10 va 20 ni trunkdan o'tkaz.

Tekshir:

```text
show vlan brief
show interfaces trunk
show mac address-table
```

## Lab 2: Inter-VLAN routing

Topology:

```text
PC1 VLAN10 -- SW1 -- R1
PC2 VLAN20 -- SW1 -- R1
```

Router-on-a-stick:

```text
R1 G0/0.10 -> 192.168.10.1/24
R1 G0/0.20 -> 192.168.20.1/24
```

Vazifa:

- Router subinterface sozla.
- Encapsulation dot1q yoz.
- PC gatewaylarini to'g'ri ber.

Tekshir:

```text
show ip interface brief
show interfaces trunk
ping 192.168.20.10
```

## Lab 3: STP root bridge

Topology:

```text
      SW1
     /   \
   SW2---SW3
```

Vazifa:

- SW1 ni VLAN 10 uchun root bridge qil.
- SW2 ni secondary root qil.
- PortFast va BPDU Guardni access portlarda yoq.

Tekshir:

```text
show spanning-tree vlan 10
show spanning-tree summary
```

## Lab 4: EtherChannel LACP

Topology:

```text
SW1 == two links == SW2
```

Vazifa:

- Ikki physical linkni Port-channel 1 ga birlashtir.
- LACP active mode ishlat.
- Port-channelni trunk qil.

Tekshir:

```text
show etherchannel summary
show interfaces port-channel 1 trunk
```

## Lab 5: Static routing

Topology:

```text
LAN1 -- R1 -- R2 -- LAN2
```

Vazifa:

- R1 da LAN2 uchun static route yoz.
- R2 da LAN1 uchun static route yoz.
- Default route va floating static route bilan test qil.

Tekshir:

```text
show ip route
show ip route static
traceroute
```

## Lab 6: OSPF single area

Topology:

```text
R1 -- R2 -- R3
```

Vazifa:

- Barcha routerlarda OSPF process yoq.
- Area 0 ishlat.
- Router IDlarni qo'lda belgilab chiq.

Tekshir:

```text
show ip ospf neighbor
show ip ospf interface brief
show ip route ospf
```

## Lab 7: HSRP

Topology:

```text
        R1
       /  \
LAN -- SW  -- upstream
       \  /
        R2
```

Vazifa:

- R1 va R2 orasida virtual default gateway yarat.
- R1 active, R2 standby bo'lsin.
- R1 uplink tushsa, R2 active bo'lsin.

Tekshir:

```text
show standby brief
```

## Lab 8: DHCP relay

Topology:

```text
Client VLAN99 -- L3 Switch -- DHCP Server VLAN1
```

Vazifa:

- VLAN99 SVI da `ip helper-address` yoz.
- DHCP serverda VLAN99 uchun scope yarat.

Tekshir:

```text
show running-config interface vlan 99
debug ip dhcp server packet
```

## Lab 9: NAT/PAT

Topology:

```text
Inside LAN -- R1 -- Outside network
```

Vazifa:

- Inside/outside interfacelarni belgila.
- ACL bilan inside subnetni tanla.
- PAT overload sozla.

Tekshir:

```text
show ip nat translations
show ip nat statistics
```

## Lab 10: ACL

Vazifa:

- VLAN10 dan VLAN20 ga pingni blokla.
- VLAN10 dan web server 80/443 ga ruxsat ber.
- ACLni to'g'ri interface va directionga qo'y.

Tekshir:

```text
show access-lists
show ip interface
```

## Lab 11: Device management security

Vazifa:

- `enable secret` sozla.
- Local user yarat.
- VTY faqat SSH qabul qilsin.
- Telnetni o'chir.

Tekshir:

```text
show ip ssh
show running-config | section line vty
ssh admin@device-ip
```

## Lab 12: Automation mini lab

Vazifa:

- REST API endpointdan JSON output ol.
- JSON ichidan hostname yoki interface statusni ajrat.
- Ansible inventory yoz.

Minimal JSON misol:

```json
{
  "hostname": "R1",
  "interfaces": [
    {"name": "GigabitEthernet0/0", "status": "up"},
    {"name": "GigabitEthernet0/1", "status": "down"}
  ]
}
```
