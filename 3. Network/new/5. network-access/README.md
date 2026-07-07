# 5. Network Access

Bu bo'lim CCNA darajasida access layer mavzularini tushuntiradi: VLAN, trunk, inter-VLAN routing, STP/Rapid PVST+, EtherChannel, CDP/LLDP va WLAN. Asosiy maqsad - switch portlari orqali foydalanuvchilarni tarmoqqa ulash, broadcast domenlarni ajratish, redundant linklarni xavfsiz ishlatish va nosozliklarni tez topish.

## Bo'lim xaritasi

| Fayl | Mavzu | Nima uchun kerak |
|---|---|---|
| `01-vlan.md` | VLAN | Bitta switch ichida mantiqiy tarmoqlar yaratish |
| `02-trunk-8021q.md` | 802.1Q trunk | VLANlarni switchlar orasida tashish |
| `03-inter-vlan-routing.md` | Inter-VLAN routing | Turli VLANlar orasida Layer 3 aloqa |
| `04-stp-rapid-pvst.md` | STP va Rapid PVST+ | Looplarni oldini olish va redundant linklarni boshqarish |
| `05-etherchannel-lacp.md` | EtherChannel/LACP | Bir nechta linkni bitta mantiqiy kanal qilish |
| `06-cdp-lldp.md` | CDP/LLDP | Qo'shni qurilmalarni aniqlash |
| `07-wireless-wlan.md` | WLAN | Simsiz tarmoq asoslari va Cisco WLAN tushunchalari |

## Network access nima?

Access layer - foydalanuvchi qurilmalari tarmoqqa ulanadigan joy. Bu yerda odatda quyidagilar qilinadi:

- PC, printer, IP phone, access point portlari sozlanadi.
- VLAN orqali trafik ajratiladi.
- Trunk orqali VLANlar boshqa switch yoki routerga olib o'tiladi.
- STP yordamida Layer 2 looplardan himoyalanadi.
- EtherChannel bilan uplink sig'imi va barqarorligi oshiriladi.
- CDP/LLDP yordamida topologiya tekshiriladi.
- WLAN orqali mobil qurilmalar uchun SSID va xavfsizlik sozlanadi.

## Tez eslab qolish

- VLAN - Layer 2 broadcast domen.
- Trunk - bir nechta VLAN trafikini bitta linkdan o'tkazadi.
- 802.1Q - VLAN tag qo'yish standarti.
- Native VLAN - trunkda taglanmaydigan VLAN.
- Inter-VLAN routing - VLANlar orasida router yoki L3 switch orqali aloqa.
- STP - loopni oldini oladi.
- Rapid PVST+ - Cisco per-VLAN rapid STP varianti.
- EtherChannel - bir nechta fizik linkdan bitta logical link.
- LACP - EtherChannel uchun ochiq standart protokol.
- CDP - Cisco discovery protokoli.
- LLDP - vendor-neutral discovery protokoli.

## Minimal lab topologiya

```text
PC1 -- SW1 ==trunk== SW2 -- PC2
          \
           ==trunk== R1 yoki L3-SW
```

Misol VLANlar:

```text
VLAN 10  USERS
VLAN 20  SALES
VLAN 30  VOICE
VLAN 99  MGMT
VLAN 999 BLACKHOLE/native
```

## Umumiy tekshiruv komandalar

```cisco
show vlan brief
show interfaces status
show interfaces trunk
show spanning-tree
show etherchannel summary
show cdp neighbors detail
show lldp neighbors detail
show ip interface brief
show running-config interface gigabitEthernet0/1
```

## Amaliy maslahatlar

- Har bir access portga aniq VLAN bering.
- Ishlatilmaydigan portlarni shutdown qiling va "parking" VLANga qo'ying.
- Trunklarda allowed VLAN ro'yxatini cheklang.
- Native VLANni foydalanuvchi VLANi qilmaslik yaxshi amaliyot.
- STP root switchni tasodifga tashlab qo'ymang, qo'lda belgilang.
- EtherChannel a'zolarida speed, duplex, trunk/access mode va allowed VLANlar bir xil bo'lishi kerak.
- CDP/LLDPni ichki tarmoqda foydali, lekin tashqi yoki ishonchsiz portlarda xavfli deb biling.

## Ko'p uchraydigan xatolar

- VLAN yaratilmagan, lekin port shu VLANga qo'yilgan.
- Trunkning bir tomoni access, ikkinchi tomoni trunk bo'lib qolgan.
- Native VLAN ikki tomonda mos emas.
- Inter-VLAN routingda gateway IP noto'g'ri berilgan.
- STP root kutilmagan switchga o'tib ketgan.
- EtherChannel portlari alohida-alohida sozlanib, `channel-group` parametrlari mos kelmagan.
- WLANda SSID ishlaydi, lekin VLAN mapping yoki default gateway xato.

## Qisqa Q&A

**Savol:** VLAN routerga o'xshab trafikni yo'naltiradimi?  
**Javob:** Yo'q. VLAN Layer 2 ajratish beradi. VLANlar orasida aloqa uchun Layer 3 routing kerak.

**Savol:** Trunk har doim switch-switch orasida bo'ladimi?  
**Javob:** Ko'pincha shunday, lekin switch-router, switch-firewall, switch-access point orasida ham trunk bo'lishi mumkin.

**Savol:** STP kerakmi, agar loop yo'q bo'lsa?  
**Javob:** Ha, access tarmoqlarda xato kabel ulanishi tez-tez bo'ladi. STP himoya qatlamidir.

**Savol:** CCNA uchun qaysi komandalar eng muhim?  
**Javob:** `show vlan brief`, `show interfaces trunk`, `show spanning-tree`, `show etherchannel summary`, `show ip interface brief`.
