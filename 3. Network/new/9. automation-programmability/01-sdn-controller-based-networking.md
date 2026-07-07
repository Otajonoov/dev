# SDN va Controller-Based Networking

SDN (Software-Defined Networking) - tarmoqni dasturiy boshqarish yondashuvi. Bunda tarmoqning mantiqiy boshqaruvi markaziy controller orqali amalga oshiriladi. Controller butun tarmoqni umumiy ko'rinishda ko'radi va siyosatlarni qurilmalarga tarqatadi.

## Traditional networking

An'anaviy tarmoqda har bir router yoki switch odatda mustaqil boshqariladi.

```text
Admin -> CLI/SSH -> Router1
Admin -> CLI/SSH -> Router2
Admin -> CLI/SSH -> Switch1
```

Masalan, VLAN qo'shish kerak bo'lsa, administrator har bir switchga kirib kerakli buyruqlarni yozadi.

Afzalliklari:

- sodda kichik tarmoqlar uchun tushunarli;
- CLI orqali aniq nazorat qilish mumkin;
- ko'p yillik amaliyot va troubleshooting usullari bor.

Kamchiliklari:

- katta tarmoqda sekin;
- inson xatosi ko'p uchraydi;
- konfiguratsiya farqlari paydo bo'lishi mumkin;
- umumiy tarmoq holatini ko'rish qiyinroq.

## Controller-based networking

Controller-based networkingda markaziy controller tarmoq qurilmalarini boshqaradi. Administrator controllerda siyosat yoki konfiguratsiya yaratadi, controller esa uni qurilmalarga yuboradi.

```text
Admin/App
   |
   | Northbound API
   v
Controller
   |
   | Southbound API
   v
Switches/Routers/APs
```

Misollar:

- Cisco DNA Center / Catalyst Center - enterprise network management va automation;
- Cisco SD-WAN Manager - SD-WAN boshqaruvi;
- wireless LAN controller - access pointlarni markaziy boshqaradi.

## Control plane, data plane, management plane

Tarmoq qurilmasidagi vazifalarni uchta plane orqali tushunish oson.

### Control plane

Control plane tarmoq qarorlarini qabul qiladi.

Misollar:

- routing table yaratish;
- OSPF/EIGRP/BGP kabi protokollar orqali marshrut o'rganish;
- STP orqali loop'ni oldini olish;
- qaysi yo'l eng yaxshi ekanini tanlash.

Oddiy misol:

```text
Router: "10.10.10.0/24 tarmog'iga borish uchun next-hop 192.168.1.2"
```

### Data plane

Data plane paketlarni real vaqtda uzatadi. U control plane qarorlaridan foydalanadi.

Misollar:

- frame'ni switch portidan chiqarish;
- IP paketni routing table asosida next-hop'ga yuborish;
- ACL yoki QoS asosida trafikni ruxsat berish, bloklash yoki belgilash.

Oddiy misol:

```text
Paket keldi -> forwarding table tekshirildi -> kerakli interfeysdan chiqdi
```

### Management plane

Management plane qurilmani boshqarish va monitoring qilish uchun ishlatiladi.

Misollar:

- SSH;
- HTTPS web UI;
- SNMP;
- syslog;
- NETCONF/RESTCONF;
- konfiguratsiya backup.

## SDN'da plane'lar qanday ajraladi?

Traditional routerda control plane va data plane ko'pincha bitta qurilma ichida ishlaydi. SDN modelida esa control plane qisman yoki to'liq controllerga ko'chishi mumkin.

```text
Traditional:
Router = Control plane + Data plane + Management plane

Controller-based:
Controller = markaziy control/management
Devices    = asosan data plane va local forwarding
```

Bu har doim ham "qurilma umuman o'ylamaydi" degani emas. Ko'p real tarmoqlarda controller siyosat beradi, qurilmalar esa local forwarding va ayrim local qarorlarni bajaradi.

## Northbound va southbound API

### Northbound API

Northbound API controller va yuqoridagi ilovalar orasida ishlaydi.

```text
Monitoring app / Automation script / Dashboard
                 |
                 v
              Controller
```

Masalan, dastur controllerdan "barcha switchlar holatini ber" deb so'rashi mumkin.

### Southbound API

Southbound API controller va tarmoq qurilmalari orasida ishlaydi.

```text
Controller
    |
    v
Switch / Router / AP
```

Misollar:

- NETCONF;
- RESTCONF;
- OpenFlow;
- gRPC/gNMI;
- vendor-specific protokollar.

CCNA uchun muhim farq:

```text
Northbound = controllerdan yuqoriga, ilovalar tomonga
Southbound = controllerdan pastga, qurilmalar tomonga
```

## Overlay, underlay va fabric

### Underlay

Underlay - fizik yoki asosiy IP tarmoq. U qurilmalar orasida haqiqiy bog'lanishni ta'minlaydi.

Misollar:

- switchlar orasidagi trunk linklar;
- routerlar orasidagi IP bog'lanishlar;
- OSPF yoki IS-IS bilan ishlaydigan asosiy transport tarmog'i.

### Overlay

Overlay - underlay ustida qurilgan mantiqiy tarmoq. U foydalanuvchi yoki servis trafikini mantiqiy ajratadi.

Misollar:

- VXLAN tunnel;
- SD-WAN overlay tunnel;
- VPN.

Oddiy tasvir:

```text
Overlay:  Tenant A virtual network
          Tenant B virtual network
              |
Underlay: physical/IP transport network
```

### Fabric

Fabric - controller tomonidan boshqariladigan, siyosatga asoslangan tarmoq arxitekturasi. Fabric ichida overlay va underlay birga ishlashi mumkin.

Masalan, campus fabric foydalanuvchilarni joylashuvdan qat'i nazar bir xil siyosat bilan boshqarishi mumkin.

## Amaliy ssenariy

Vazifa: yangi filialga 20 ta access point qo'shish.

Traditional yondashuv:

1. Har bir AP yoki switch alohida sozlanadi.
2. SSID, VLAN, security policy qo'lda kiritiladi.
3. Xato bo'lsa, qurilmalarni alohida tekshirish kerak.

Controller-based yondashuv:

1. Controllerda filial profili yaratiladi.
2. AP'lar controllerga ulanadi.
3. SSID, VLAN, security policy avtomatik tarqatiladi.
4. Dashboardda holat ko'rinadi.

## Common mistakes

- **Controller hamma narsani o'zi forward qiladi deb o'ylash.** Ko'pincha trafik qurilmalarning o'zida forward qilinadi.
- **Northbound va southboundni chalkashtirish.** Northbound ilovalarga, southbound qurilmalarga qaraydi.
- **Overlay underlay o'rnini bosadi deb o'ylash.** Overlay ishlashi uchun barqaror underlay kerak.
- **Management plane bilan control plane'ni bir xil deb bilish.** Management qurilmani boshqaradi, control plane esa forwarding qarorlarini tayyorlaydi.

## Qisqa Q&A

**Savol:** Controller ishdan chiqsa, butun tarmoq to'xtaydimi?

**Javob:** Dizaynga bog'liq. Ko'p tizimlarda mavjud forwarding davom etadi, lekin yangi siyosat, yangi qurilma onboarding yoki markaziy monitoring ta'sirlanishi mumkin.

**Savol:** SDN faqat data center uchunmi?

**Javob:** Yo'q. SDN campus, WAN, wireless va cloud tarmoqlarda ham uchraydi.

**Savol:** Fabric degani faqat bitta mahsulotmi?

**Javob:** Yo'q. Fabric umumiy arxitektura tushunchasi: ko'p qurilmalar yagona siyosat va controller orqali boshqariladi.
