# 01. Network Access (Physical + Data Link, L2 switching)

Bu modul tarmoqning eng quyi ikki qatlamini — **Physical (L1)** va **Data Link
(L2)** ni CCNA darajasida o'rgatadi. Bu yerda foydalanuvchi qurilmalari tarmoqqa
ulanadi: kabel va signaldan tortib, VLAN, trunk, STP, EtherChannel va simsiz
tarmoqgacha. Asosiy maqsad — switch portlari orqali foydalanuvchilarni ulash,
broadcast domenlarni ajratish, redundant linklarni xavfsiz ishlatish va nosozliklarni
tez topish.

Har dars bir xil pedagogik tuzilishda: **muammo/hook -> real-hayot analogiya ->
sodda ta'rif -> Mermaid diagramma -> worked example (Cisco CLI) -> Xulosa, Eslab qol,
O'z-o'zini tekshir, Amaliyot, Takrorlash**. Mavzular 2025–2026 professional best
practice bilan boyitilgan (800G/1.6T kabellash, Wi-Fi 7/8, WPA3, VLAN hardening).

## Darslar ro'yxati (o'qish tartibi)

| # | Dars | Mavzu |
|---|------|-------|
| 1 | [01-physical-layer.md](01-physical-layer.md) | Kabel turlari, signal, encoding, bandwidth, hub |
| 2 | [02-data-link-ethernet-mac.md](02-data-link-ethernet-mac.md) | Ethernet frame, MAC, switch MAC table, collision/broadcast domain, ARP |
| 3 | [03-vlan.md](03-vlan.md) | VLAN, access port, voice VLAN, segmentatsiya |
| 4 | [04-trunk-8021q.md](04-trunk-8021q.md) | Trunk, 802.1Q tag, native VLAN, DTP xavfsizligi |
| 5 | [05-inter-vlan-routing.md](05-inter-vlan-routing.md) | Router-on-a-stick, L3 switch SVI |
| 6 | [06-stp.md](06-stp.md) | STP, Rapid PVST+, root bridge, PortFast, BPDU/Root/Loop Guard |
| 7 | [07-etherchannel-lacp.md](07-etherchannel-lacp.md) | EtherChannel, LACP, load balancing |
| 8 | [08-cdp-lldp.md](08-cdp-lldp.md) | CDP, LLDP, LLDP-MED, discovery xavfsizligi |
| 9 | [09-wireless-wlan.md](09-wireless-wlan.md) | WLAN, AP, WLC, SSID/VLAN, Wi-Fi 6E/7, WPA3 |

## Qanday o'qish kerak

1. Darslarni **tartib bilan** o'qi — har biri oldingisiga tayanadi (masalan STP
   VLAN va broadcast domain tushunchasini talab qiladi).
2. Har darsda avval **muammoni his qil**, keyin yechimni o'rgan.
3. Worked example dagi Cisco CLI ni Packet Tracer yoki lab da **o'zing yozib** sina.
4. Har dars oxiridagi **O'z-o'zini tekshir** savollariga javobni ochishdan oldin
   o'zing javob berishga urin (retrieval practice).
5. **Takrorlash jadvali** bo'yicha ertaga -> 3 kun -> 1 hafta oralig'ida qayt.

## Umumiy tekshiruv buyruqlari

```cisco
show vlan brief
show interfaces status
show interfaces trunk
show spanning-tree
show etherchannel summary
show cdp neighbors detail
show lldp neighbors detail
show mac address-table
show ip interface brief
```
