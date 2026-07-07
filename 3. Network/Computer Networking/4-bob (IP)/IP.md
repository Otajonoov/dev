
## Internet Protocol (IP) haqida qisqacha

Network layer protokoli **IP (Internet Protocol)** quyidagi xususiyatlarga ega:

- Hostlar o'rtasida logical connection ta'minlaydi
- **Unreliable service** - ishonchsiz xizmat
- Segment'larni yetkazishga kafolat bermaydi
- Tartibni saqlashga kafolat bermaydi
- Har bir host o'zining **IP address**'iga ega

---

## Network Layer nima?

Network layer - bu internetdagi barcha host va routerlarni bog'laydigan muhim qism. 
Transport layer faqat ma'lumotni jo'natuvchi va qabul qiluvchi o'rtasida uzatsa, network layer bu jarayonni qanday amalga oshirishni ta'minlaydi.

**Forwarding va Routing o'rtasidagi farq:**

- **Forwarding** - bitta router ichida paketni kirish portidan chiqish portiga yo'naltirish
- **Routing** - butun network bo'ylab paketlar uchun optimal yo'lni aniqlash

Bu xuddi pochta tizimiga o'xshaydi: forwarding - bu pochta bo'limida xatni to'g'ri qutiga tashlash, routing esa xat uchun butun shahar bo'ylab eng yaxshi marshrutni tanlash.

## Network Layer Arxitekturasi

Oddiy network strukturasi:

```
Host X1 → Router M1 → Router M2 → Host X2
```

Bu jarayonda:

1. X1 host transport layerdan segmentlarni oladi
2. Ularni datagramlarga (network layer paketlari) o'rash
3. Eng yaqin routerga (M1) jo'natish
4. Routerlar paketlarni maqsadgacha yo'naltirish
5. X2 hostda datagramlardan segmentlarni chiqarib, transport layerga uzatish

## Protocol Stack tuzilishi

**Host tizimlarda:**

- Application Layer
- Transport Layer
- Network Layer
- Data Link Layer
- Physical Layer

**Routerlarda:**

- Network Layer
- Data Link Layer
- Physical Layer

Routerlar faqat paketlarni yo'naltirish bilan shug'ullanadiganligi sababli, ularga yuqori layerlar kerak emas.

## Ikki Asosiy Model

1. **Datagram modeli** - har bir paket mustaqil ravishda yo'naltiriladi
2. **Virtual circuit modeli** - oldindan belgilangan yo'l bo'yicha uzatish

## Amaliy Ahamiyati

Network layer internetning "asosiy mexanizmi" hisoblanadi. U:

- IP addresslash tizimini ta'minlaydi
- Paketlarning yo'lini aniqlaydi
- Network traffic boshqaruvi
- Xatoliklarni qayta ishlash (ICMP orqali)
- IPv4 dan IPv6 ga o'tishni ta'minlaydi

Network layer bo'lmasa, internet faqat bir-biriga to'g'ridan-to'g'ri ulangan kompyuterlar to'plami bo'lib qolar edi. Aynan bu layer tufayli biz istalgan nuqtadan dunyoning istalgan nuqtasiga ma'lumot jo'nata olamiz.

---

## Forwarding va Routing o'rtasidagi muhim farq

### 1. Forwarding (Yo'naltirish)

- **Bir router ichidagi** jarayon
- Paket router kirish interfaceiga kelganda uni to'g'ri chiqish interfaceiga yo'naltirish
- Masalan: X1 dan kelgan paket M1 routeriga kelsa, uni keyingi routerga yo'naltirish

**Real hayotiy misol:** Bu xuddi chorrahada yo'l ko'rsatuvchi politsiyachiga o'xshaydi - u har bir mashinani to'g'ri yo'nalishga yo'naltiradi.

### 2. Routing (Marshrutlash)

- **Butun network bo'yicha** jarayon
- Jo'natuvchidan qabul qiluvchigacha eng yaxshi yo'lni aniqlash
- Routing algoritmlari bu yo'llarni hisoblaydi

**Real hayotiy misol:** Bu sayohat rejasini tuzishga o'xshaydi - xaritani o'rganib, maqsadgacha eng yaxshi marshrutni tanlash.

## Forwarding Table qanday ishlaydi

Har bir routerda **forwarding table** mavjud:

```
Header qiymati | Chiqish interfacesi
0100           | 3
0101           | 2  
0111           | 2
1001           | 1
```

**Jarayon:**

1. Paket keladi (masalan, header qiymat: 0111)
2. Router forwarding tabledan qidiradi
3. Interface 2 ga yo'naltiradi
4. Paket keyingi routerga jo'natiladi

## Routing Algorithm va Forwarding Table bog'liqlik

**Routing algorithm** forwarding tableni to'ldiradi:

- **Markazlashgan** - bir joyda hisoblanib, barchaga tarqatiladi
- **Taqsimlangan** - har bir router o'z algoritmini ishga tushiradi

## Packet Switch vs Router

**Packet Switch** - umumiy nom (packet ni bir interfacedan boshqasiga yo'naltiruvchi qurilma)

- **Switch** (Layer 2) - Data Link Layer header asosida ishlaydi
- **Router** (Layer 3) - Network Layer header asosida ishlaydi

## Connection Setup (uchinchi funksiya)

Ba'zi network arxitekturalarda (ATM, Frame Relay, MPLS) **connection setup** ham kerak:

- Har bir router yo'lda handshake qilishi kerak
- Ma'lumot uzatish boshlanishidan oldin connection holatini o'rnatish
- TCP dagi handshake ga o'xshash jarayon

## Amaliy Ahamiyat

**Forwarding** - tez va oddiy (mikrosekundlarda) 
**Routing** - murakkab va vaqt talab etadi (soniyalar/daqiqalar)

Bu ikki funksiya birgalikda internetning ishlashini ta'minlaydi. 
Forwarding operativ tezlikda, Routing esa strategik rejalashtirish darajasida ishlaydi.

---

# Network Layer Service Modellari: Internetdan ATM gacha

## Network Layer Service nima berishi mumkin?

Transport layer paket jo'natganda, network layer turli xil **xizmat kafolatlari** bera oladi. Bu xuddi pochta xizmatida turli darajadagi yetkazib berish xizmatlari mavjud bo'lgani kabi.

## Asosiy Service turlari

### 1. Yagona Paketlar uchun xizmatlar

- **Guaranteed delivery** - paket albatta yetib borishini kafolatlash
- **Guaranteed delivery with bounded delay** - ma'lum vaqt ichida (masalan, 100ms) yetib borishini kafolatlash

### 2. Paket oqimi uchun xizmatlar

- **In-order packet delivery** - paketlar jo'natilgan tartibda yetib borishi
- **Guaranteed minimal bandwidth** - kafolatlangan minimal tezlik (masalan, 1 Mbps)
- **Guaranteed maximum jitter** - paketlar orasidagi vaqt intervali barqaror bo'lishi
- **Security services** - shifrlash va autentifikatsiya

## Uchta Asosiy Service Model

### 1. Internet: Best-effort Service

```
Kafolatlar: Hech qanday kafolat YO'Q
- Paket yo'qolishi mumkin
- Tartib buzilishi mumkin  
- Vaqt kechikishi mumkin
- Xavfsizlik kafolati yo'q
```

**Real hayotiy misol:** Bu oddiy pochta xizmatiga o'xshaydi - xat jo'natasiz, lekin yetib borishini kafolatlamaysiz.

### 2. ATM CBR (Constant Bit Rate)

```
Kafolatlar:
✓ Doimiy tezlik kafolati
✓ Paket yo'qolmasligi
✓ Tartib saqlanishi  
✓ Vaqt munosabatlari barqaror
✓ Overload bo'lmasligi
```

**Qo'llanish:** Audio va video uchun ideal - telefon kompaniyalari uchun yaratilgan.

### 3. ATM ABR (Available Bit Rate)

```
Kafolatlar:
✓ Minimal tezlik kafolati (MCR)
✓ Tartib saqlanishi
✓ Feedback mexanizmi
△ Paket yo'qolishi mumkin (lekin kamroq)
△ Vaqt munosabatlari yo'q
```

## Service Model Comparison jadval

|Arxitektura|Bandwidth|Yo'qolmaslik|Tartib|Vaqt|Feedback|
|---|---|---|---|---|---|
|**Internet**|Yo'q|Yo'q|Buziladi|Yo'q|Yo'q|
|**ATM CBR**|Kafolat|Ha|Saqlanadi|Ha|Kerak emas|
|**ATM ABR**|Minimal|Yo'q|Saqlanadi|Yo'q|Ha|

## Qaysi model yaxshiroq?

**Internet Best-effort** afzalliklari:

- Oddiy va tez
- Flexible - har qanday trafik turi
- Masshtablanadi
- Transport layer kompensatsiya qiladi (TCP reliability beradi)

**ATM modellari** afzalliklari:

- Real-time traffic uchun yaxshi
- Quality of Service (QoS) kafolatlari
- Professional telecom standard

## Amaliy Natija

**Internet oddiy, ATM murakkab** - lekin internet g'alaba qildi chunki:

- Oddiyligi masshtablanishni osonlashtirdi
- TCP/UDP yuqori layerlarda kerakli kafolatlarni berdi
- Arzon va flexible ekanligini isbotladi

Bugungi kunda ham internet asosan **best-effort** model bilan ishlaydi, lekin QoS texnologiyalar orqali ba'zi kafolatlar qo'shilmoqda.

---
# Virtual Circuit va Datagram Networks: Ikki Xil Yondashuv

## Virtual Circuit Networks (Masalan: ATM, Frame Relay)

### Virtual Circuit tuzilishi

Virtual circuit uchta asosiy elementdan tashkil topadi:

1. **Route** - jo'natuvchidan qabul qiluvchigacha yo'l (routerlar ketma-ketligi)
2. **VC raqamlari** - har bir link uchun alohida raqam
3. **Forwarding table** - barcha routerlarda maxsus yozuvlar

### VC raqamlari qanday ishlaydi?

**Misol:** A dan B ga paket jo'natish

```
A → M1 → M2 → B
   12   22   32
```

Paket har bir routerda yangi VC raqam oladi:

- A dan chiqganda: VC=12
- M1 dan chiqganda: VC=22
- M2 dan chiqganda: VC=32

### Forwarding Table misoli (M1 router):

|Kirish Interface|Kirish VC|Chiqish Interface|Chiqish VC|
|---|---|---|---|
|1|12|2|22|
|2|63|1|18|
|3|7|2|17|

### VC ning uch bosqichi:

1. **VC Setup** - yo'l o'rnatish, resurs ajratish
2. **Data Transfer** - ma'lumot uzatish
3. **VC Teardown** - connection yopish

**Real hayotiy misol:** Bu telefon qo'ng'irog'iga o'xshaydi - avval raqam terish, so'ngra gaplashish, keyin telefon qo'yish.

## Datagram Networks (Internet)

### Asosiy prinsip

Har bir paket **mustaqil** yo'naltiriladi:

- Connection setup yo'q
- Har bir paket destination address olib yuradi
- Routerlar forwarding table asosida qaror qabul qiladi

### Forwarding Table misoli:

32-bitli IP address uchun **prefix matching**:

|Prefix|Interface|
|---|---|
|11001000 00010111 00010|0|
|11001000 00010111 00011000|1|
|11001000 00010111 00011|2|
|boshqalar|3|

### Longest Prefix Matching

Agar bir nechta prefix mos kelsa, **eng uzun** prefix tanlanadi:

**Misol:** `11001000 00010111 00011000 10101010`

- 21-bit prefix (3-qator): mos keladi ✓
- 24-bit prefix (2-qator): mos keladi ✓✓
- **24-bit tanlandi** (uzunroq)

## Virtual Circuit vs Datagram Taqqoslash

|Xususiyat|Virtual Circuit|Datagram|
|---|---|---|
|**Setup**|Ha, kerak|Yo'q|
|**State info**|Routerlarda saqlanadi|Minimal|
|**Addressing**|VC raqam|Destination address|
|**Path**|Fixed (bir xil yo'l)|Variable (har xil yo'l)|
|**Packet order**|Kafolatlanadi|Buzilishi mumkin|
|**Failure recovery**|Murakkab|Oddiy|
|**Resource management**|Yaxshi|Limited|

## Qaysi yondashuv yaxshiroq?

### Virtual Circuit afzalliklari:

- QoS kafolatlari
- Resource reservation
- Predictable performance
- Paket tartibi saqlanadi

### Datagram afzalliklari:

- **Oddiylik** - connectionless
- **Flexibility** - har bir paket mustaqil
- **Fault tolerance** - bir router ishlamasa, boshqa yo'l topadi
- **Scalability** - kamroq state ma'lumot

## Amaliy Natija

**Internet datagram modelni tanladi** chunki:

- Internetning asosiy maqsadi **survivability** edi (urush sharoitida ham ishlashi)
- Oddiy va arzon
- Har xil network texnologiyalarini birlashtirish oson
- End-to-end principi: reliability Transport layerda (TCP) hal qilinadi

**Virtual Circuit** maxsus telecom sohalarda hali ham ishlatiladi (MPLS, ATM) chunki QoS kafolatlari muhim

---

# Virtual Circuit va Datagram Networklarning Kelib Chiqishi

## Virtual Circuit: Telefon Dunyosidan Meros

### Tarixi Kelib Chiqishi

Virtual circuit tushunchasi **telefon tarmog'idan** kelib chiqqan:

- Telefon tarmog'i **physical connection** (fizik ulanish) printsipiga asoslangan
- Eski diskli telefonlar (rotary phones) juda oddiy qurilmalar edi
- Network barcha murakkablikni o'ziga olishi kerak edi

**Real hayotiy misol:** Bu eski telefonlarga o'xshaydi - siz faqat raqam terasiz, qolgan hamma ishni telefon stantsiyasi bajaradi.

### Virtual Circuit murakkabligi:

```
Telefon Network                Virtual Circuit Network
- Connection setup     →       - VC setup  
- Circuit switching    →       - State maintenance
- "Stupid" endpoints   →       - Router complexity
- Centralized control  →       - Signaling protocols
```

## Datagram: Internet Innovatsiyasi

### Internet Arxitektorlarining Fikrlash Tarzi

Internet yaratuvchilari **boshqacha yondashuvni** tanladi:

- Kompyuterlar aqlli qurilmalar (telefonlardan farqli)
- Network layerni **oddiy** qilish mumkin
- Murakkablikni **end systemlarga** o'tkazish mumkin

### End-to-End Prinsipi

Internet model:

```
Simple Network Core + Intelligent Endpoints = Scalable Internet

Oddiy Network + Aqlli End Systemlar = Masshtablanadigan Internet
```

## Har Bir Yondashuvning Afzalliklari

### Virtual Circuit (Telefon Model):

**Afzalliklari:**

- QoS kafolatlari
- Centralized control
- Predictable performance

**Kamchiliklari:**

- Murakkab infrastructure
- Har xil texnologiyalarni birlashtirish qiyin
- Yangi service qo'shish murakkab

### Datagram (Internet Model):

**Afzalliklari:**

- **Network simplicity** - soda infrastruktura
- **Heterogeneous network support** - har xil texnologiyalar (Satellite, Ethernet, fiber, radio)
- **Innovation at edges** - yangi servicelar osonlik bilan qo'shiladi
- **Scalability** - million/milliard qurilmalar bilan ishlaydi

**Misol:** World Wide Web qanchalik tez paydo bo'ldi:

- HTTP protokol yaratildi
- Server ishga tushirildi
- Butun dunyo foydalana boshladi
- Hech qanday network infrastructure o'zgarishi kerak bo'lmadi

## Internet Model Success Factors

### 1. Minimal Network Layer Requirements

```
Best-effort service = Minimal guarantees
↓
Easy interconnection of different technologies
```

### 2. Intelligence at End Systems

```
Email servers    →    Application Layer
TCP reliability  →    Transport Layer  
DNS service     →    Application Layer
HTTP/Web        →    Application Layer
```

### 3. Service Innovation

Yangi service yaratish uchun:

- Faqat server dasturi yozish
- Protocol belgilash
- Network infrastructure o'zgarmaydi

## Amaliy Natija va Ta'sirlar

### Internetning Muvaffaqiyat Sabablari:

1. **Simplicity** - oddiy network core
2. **Flexibility** - har xil texnologiyalar
3. **Innovation** - edge da innovatsiya
4. **Scalability** - global darajada ishlaydi

### Bugungi Holat:

- Internet **datagram model** asosida global network bo'ldi
- Virtual circuit faqat maxsus sohalarda (telecom, enterprise MPLS)
- **Hybrid approaches** - MPLS, SD-WAN (ikki modelni birlashtiradi)

Bu internet arxitekturasining eng muhim qarorlari biri bo'ldi - **simple network, smart endpoints** prinsipi internetning global miqyosdagi muvaffaqiyatiga olib keldi.

---

# Router Ichidagi Arxitektura

## Routerning To'rt Asosiy Komponenti

### 1. Input Portlar

**Funksiyalari:**

- **Physical layer** - kiruvchi signalni qabul qilish
- **Data link layer** - frame processing
- **Lookup va forwarding** - eng muhim funksiya

**Lookup jarayoni:**

```
Packet keladi → Header tekshiriladi → Forwarding table → Output port aniqlanadi
```

### 2. Switching Fabric (Kommutatsiya matritsasi)

- Input portlarni output portlar bilan bog'laydi
- Routerning **ichki network**i
- Paketlarni bir portdan boshqasiga o'tkazadi

### 3. Output Portlar

- Switching fabricdan paket oladi
- **Buffering** - navbat saqlash
- **Physical layer** - tashqariga jo'natish
- Ko'pincha input port bilan birgalikda joylashgan (ikki tomonlama link uchun)

### 4. Routing Processor

- **Routing protokollarini** ishlatadi
- **Forwarding table**ni yaratadi va yangilaydi
- Network managament funksiyalari
- **Software** orqali ishlaydi

## Ikki Xil Layer: Hardware vs Software

### Data Plane (Hardware - nanosekund)

```
Input Port → Switching Fabric → Output Port
     ↓              ↓              ↓
  Lookup        Switching      Buffering
(Hardware)    (Hardware)     (Hardware)
```

**Real misol:** 10 Gbps tezlik, 64-byte paket = 51.2 nanosekund vaqt!

### Control Plane (Software - millisekund/sekund)

```
Routing Processor:
- Routing protokollari
- Network management
- Table updating
```

## Nima uchun Hardware kerak?

**Hisob-kitob:**

- **10 Gbps** tezlik
- **64 byte** paket o'lchami
- Har bir paket uchun **51.2 nanosekund** vaqt
- Bir necha port bir vaqtda ishlasa, N marta tezroq kerak

Software bu tezlikda ishlay olmaydi!

## Real Hayotiy Analogiya: Chorrahali Perimeter

### Router Komponentlari = Chorhaha Elementlari

| Router Component      | Chorhaha Analogi                    |
| --------------------- | ----------------------------------- |
| **Input Port**        | Kirish yo'li + shlagbaum + operator |
| **Lookup Table**      | Operatorda turgan yo'l xaritasi     |
| **Switching Fabric**  | Aylana chorraha                     |
| **Output Port**       | Chiqish yo'li                       |
| **Routing Processor** | Traffic boshqaruv markazi           |

### Bottleneck Muammolari:

1. **Operator sekin ishlasa** = Input port lookup sekin
2. **Chorrahada tiqilinch** = Switching fabric to'lib ketishi
3. **Bir yo'lga ko'p mashina** = Output port buffer to'lishi
4. **Emergency vehicles** = Priority traffic management

## Arxitektura Turlari

### Hardware Implementation:

- **ASIC** (Application-Specific Integrated Circuit) - custom hardware
- **NPU** (Network Processing Unit) - Intel, Broadcom chiplar

### Performance Talablari:

- **Throughput:** Gbps/Tbps darajasida
- **Latency:** Microsecond/nanosecund
- **Packet rate:** Million packets per second

## Zamonaviy Router Challenges

1. **Speed** - 100Gbps, 400Gbps, 800Gbps
2. **Scale** - million routes, billion connections
3. **Features** - QoS, security, analytics
4. **Power** - energy efficiency
5. **Cost** - affordable solutions

Router arxitekturasi bu **speed, scale va features** orasidagi muvozanatni topish san'ati hisoblanadi. Hardware acceleration va parallel processing orqali zamonaviy routerlar terabit darajasida ma'lumot qayta ishlay oladi.

---

Matn asosida **Router (маршрутизатор)ning Input Port ishlash jarayoni** haqida eng muhim ma'lumotlarni qisqacha bayon qilaman:

## Input Port ning asosiy vazifasi

Input port routerning eng muhim qismi bo'lib, kiruvchi paketlarni qayta ishlash uchun mas'uldir. Bu jarayon bir necha bosqichdan iborat:

### 1. Ma'lumotlarni qayta ishlash bosqichlari

Input port quyidagi tartibda ishlaydi:

- **Physical layer** va **Data Link layer** protokollari orqali liniya bilan bog'lanish
- Paketni parsing qilish va formatini tekshirish
- Forwarding table orqali kerakli output portni topish
- Paketni switching matrix orqali yuborish uchun navbatga qo'yish

### 2. Forwarding table bilan ishlash

Eng muhim jarayon - bu **lookup** operatsiyasi. Router forwarding table dan foydalanib, paketni qaysi output portga yuborishni aniqlaydi. Bu table:

- Routing processor tomonidan yaratiladi va yangilanadi
- Har bir input portda shadow copy saqlanadi
- PCI bus orqali interfeys kartalariga uzatiladi
- **Longest prefix matching** algoritmidan foydalanadi

### 3. Tezlik talablari

Gigabit tezlikdagi liniyalarda lookup operatsiyasi **nanosekund**larda bajarilishi kerak. Buning uchun:

- Oddiy linear search yetarli emas
- Maxsus algoritmlar kerak (Gupta va Ruiz-Sanchez tadqiqotlariga binoan)
- **DRAM** yoki **SRAM** memory turlaridan foydalaniladi
- **TCAM (Ternary Content Address Memory)** tez lookup uchun ishlatiladi

### 4. Amaliy misol

Cisco 8500 routerlarda har bir input portda 64 KB associative memory mavjud bo'lib, u har qanday IP address uchun bir xil vaqtda table dan qiymat olish imkonini beradi.

### 5. Match-Action prinsipi

Input port "**match plus action**" prinsipiga asoslangan:

- **Match**: IP addressni forwarding table da topish
- **Action**: paketni kerakli output portga yuborish

Bu printsip faqat routerlarda emas, balki switch, firewall va NAT qurilmalarida ham qo'llaniladi.

### Xulosa

Input port routerning "miyasi" hisoblanib, paketlarni to'g'ri yo'nalishga yuborish uchun javobgar. Uning samarali ishlashi butun tarmoqning performance ini belgilaydi.

---

Matn asosida **Router switching (коммутация)** jarayoni haqida eng muhim ma'lumotlarni bayon qilaman:

## Switching Matrix ning roli

Switching matrix routerning asosiy qismi bo'lib, paketlarni input portdan output portga o'tkazish vazifasini bajaradi. Bu jarayon uchun **uchta asosiy usul** mavjud:

### 1. Memory orqali switching (коммутация через память)

**Eski model:**

- Eng oddiy va birinchi router modellari
- CPU (routing processor) nazorati ostida ishlaydi
- Input port paket olganida, CPU ga interrupt signal beradi
- Paket CPU memory ga ko'chiriladi
- CPU destination addressni topib, output portga yuboradi

**Cheklanishi:** Agar memory throughput N paket/sekund bo'lsa, umumiy tezlik N/2 dan kam bo'ladi, chunki har bir paket ikki marta ko'chiriladi.

**Zamonaviy model:**

- Input interface board o'zi lookup va memory operatsiyalarini bajaradi
- Shared memory ga to'g'ridan-to'g'ri yoziladi
- Cisco Catalyst 8500 bu usuldan foydalanadi

### 2. Bus orqali switching (коммутация через шину)

**Ishlash prinsipi:**

- Input port paketga **internal label** (header) qo'shadi
- Bu label qaysi output portga borishni ko'rsatadi
- Paket shared bus orqali yuboriladi
- Barcha output portlar paketni oladi, lekin faqat kerakli port saqlaydi
- Label o'chiriladi

**Cheklanishi:** Bir vaqtda faqat bitta paket bus orqali o'ta oladi - bu "transport aylanasi"da bitta mashinaning o'tishiga o'xshaydi.

**Amaliy misol:** Cisco 5600 32 Gbit/s tezlikda bus switching ishlatadi.

### 3. Interconnection network orqali switching

**Crossbar switch:**

- **2N bus** dan iborat (N input va N output portlar uchun)
- Har bir vertical bus horizontal bus bilan kesishadi
- Controller kesishish nuqtalarini ochish/yopish orqali boshqaradi

**Afzalliklari:**

- **Parallel processing** imkoniyati
- Bir vaqtda bir necha paket uzatilishi mumkin
- Masalan: A portdan Y portga va B portdan X portga bir vaqtda yuborish

**Cheklanishi:** Agar ikki xil input portdan bitta output portga paket kelsa, ular ketma-ket uzatiladi.

**Murakkab network:** Ko'p bosqichli switching orqali bir necha paketni bir vaqtda bitta output portga yuborish mumkin.

**Amaliy misol:** Cisco 12000 seriyasi interconnection network ishlatadi.

### Xulosa

Switching usuli routerning **performance va cost**ini belgilaydi:

- **Memory switching** - oddiy, lekin sekin
- **Bus switching** - o'rtacha, kichik tarmoqlar uchun yetarli
- **Interconnection network** - eng tez, lekin murakkab va qimmat

Zamonaviy yuqori tezlikdagi routerlarda asosan interconnection network ishlatiladi.

---
Matn asosida **Queue formation (очереди формирования)** jarayoni haqida eng muhim ma'lumotlarni bayon qilaman:

## Queuing ning kelib chiqishi

Routerlarda queue (navbat) - bu **traffic intensity** va **switching matrix** tezligi o'rtasidagi nomutanosiblik natijasida hosil bo'ladi. Xuddi yo'l chorrahasida avtomobillar navbat hosil qilgani kabi, paketlar ham input va output portlarda kutishi mumkin.

## Output Port Queuing

### Asosiy printsip

Agar switching matrix tezligi **R_matrix = N × R_connection** bo'lsa:

- Input portlarda deyarli queue hosil bo'lmaydi
- Lekin output portda **bottleneck** paydo bo'ladi

### Misol bilan tushuntirish

- N ta input portdan bitta output portga paket kelsa
- Output port faqat **bitta paket/time unit** jo'nata oladi
- Natijada N-1 ta paket queue da kutadi
- Keyingi time unit da yana N ta paket kelishi mumkin
- Bu jarayon davom etsa, **packet loss** yuz beradi

### Packet Scheduling algoritmlari

**FCFS (First-Come-First-Served):**

- Eng oddiy usul - kelgan tartibda jo'natish

**WFQ (Weighted Fair Queuing):**

- Har bir connection dan adolatli taqsimlash
- Quality of Service kafolatlari uchun muhim

## Buffer sizini hisoblash

### Klassik formula

```
B = RTT × C
```

- B = Buffer size
- RTT = Round Trip Time
- C = Channel capacity

**Misol:** 10 Gbit/s kanal, 250ms RTT uchun: B = 2.5 Gbit buffer kerak

### Zamonaviy formula (ko'p TCP flow lar uchun)

```
B = (RTT × C) / √N
```

- N = TCP flow lar soni
- Katta N da buffer size sezilarli kamayadi

## Active Queue Management (AQM)

### RED (Random Early Detection) algoritmi

**Ishlash prinsipi:**

- **min_th** va **max_th** threshold lar belgilanadi
- Queue length < min_th: barcha paketlar qabul qilinadi
- Queue length > max_th: yangi paketlar drop qilinadi
- min_th < Queue length < max_th: paketlar ehtimoliy drop qilinadi

**Afzalligi:** Buffer to'lishdan oldin congestion signali berish

## Input Port Queuing

### Head-of-Line (HOL) Blocking muammosi

**Misol orqali:**

- Ikki input portdan bitta output portga paket kelsa
- Birinchi paket jo'natiladi
- Ikkinchi paket kutadi
- Ikkinchi paketning ortidagi paket boshqa output portga ketishi kerak bo'lsa ham kutishga majbur

**Oqibati:**

- Input arrival rate 58% ga yetganda unlimited queue growth
- Packet loss ehtimoli oshadi

### Yechimlar

- **Virtual output queuing** - har bir input portda har output port uchun alohida queue
- More sophisticated scheduling algorithms
- Non-blocking switch fabrics

## Amaliy tavsiyalar

**Buffer management:**

- **Tail drop** - yangi paketlarni tashlab yuborish
- **Selective drop** - ma'lum paketlarni tanlab tashlab yuborish
- **Packet marking** - paket headeriga congestion belgilari qo'yish

**Performance optimization:**

- Queue length monitoring
- Dynamic threshold adjustment
- Load balancing across multiple paths

### Xulosa

Queue management - routerning eng muhim qismi bo'lib, network performance va packet loss ni to'g'ridan-to'g'ri belgilaydi. Zamonaviy routerlarda murakkab algoritm va hardware optimizationlar qo'llaniladi.

---

Matn asosida **Routing Control Level** va **IP Protocol** haqida eng muhim ma'lumotlarni bayon qilaman:

## Routing Control Level (маршрутизация boshqaruv darajasi)

### An'anaviy arxitektura

- Barcha boshqaruv funktsiyalari **routing processor** da joylashgan
- Network bo'ylab **decentralized** boshqaruv
- Har bir router o'z algoritmlarini ishlatadi
- Control message lar orqali o'zaro aloqa

### Zamonaviy yondashuv

Ba'zi tadqiqotchilar yangi arxitektura taklif qilmoqdalar:

- **Hardware forwarding** va **software control** ning ajratilishi
- Ba'zi funktsiyalar routerdan tashqarida (masalan, centralized serverda)
- **API** orqali aniq o'zaro aloqa qoidalari
- **Centralized route calculation** distributed hisoblash o'rniga

## IP Protocol tuzilishi

Internet **Network layer** uch komponentdan iborat:

### 1. IP Protocol

- Addressing qoidalari
- Datagram formati
- Packet processing qoidalari

### 2. Routing Protocol

- Route selection (RIP, OSPF, BGP)
- Forwarding table yaratish

### 3. ICMP Protocol

- Error reporting
- Network diagnostics

## IPv4 Datagram formati

### Asosiy headerlar (jami 32 bit = 4 byte)

**Version (4 bit)** - IP protocol versiyasi **Header Length (4 bit)** - Header uzunligi (odatda 20 byte) **Type of Service (8 bit)** - QoS talablari (past latency, yuqori throughput) **Total Length (16 bit)** - Butun datagram uzunligi (maksimal 65,535 byte)

**Identification (16 bit)** - Fragmentation uchun **Flags (3 bit)** - Fragmentation boshqaruvi  
**Fragment Offset (13 bit)** - Fragment joylashuvi

**Time to Live (8 bit)** - Har routerda 1 ga kamayadi, 0 bo'lsa delete qilinadi **Protocol (8 bit)** - Upper layer protocol (6=TCP, 17=UDP) **Header Checksum (16 bit)** - Header xatoliklarini aniqlash

**Source IP Address (32 bit)** **Destination IP Address (32 bit)** **Options (variable)** - Ixtiyoriy qo'shimcha ma'lumotlar **Data** - Transport layer segment (TCP/UDP)

## IP Fragmentation jarayoni

### Sababi

- Har bir **Data Link layer** ning o'z **MTU** (Maximum Transmission Unit) si bor
- Ethernet: 1500 bytes, ba'zi WAN linklar: 576 bytes
- Katta datagram kichik MTU orqali o'ta olmaydi

### Fragmentation algoritmi

**Misol:** 4000 byte datagram, MTU=1500

```
Original: 4000 bytes (20 header + 3980 data)

Fragment 1: 1500 bytes (20 header + 1480 data)
- ID=777, Offset=0, Flag=1

Fragment 2: 1500 bytes (20 header + 1480 data)  
- ID=777, Offset=185, Flag=1

Fragment 3: 1040 bytes (20 header + 1020 data)
- ID=777, Offset=370, Flag=0 (last fragment)
```

### Reassembly (qayta yig'ish)

- Faqat **destination host** da bajariladi
- Routerlarda emas (network core simple bo'lishi uchun)
- Agar biron fragment yo'qolsa, butun datagram delete qilinadi
- TCP layer keyin missing data ni qayta so'raydi

## Fragmentation ning salbiy tomonlari

### Xavfsizlik muammolari

- **Jolt2 attack**: nol offset bo'lmagan kichik fragmentlar
- **Overlapping fragments**: noto'g'ri offset lar bilan
- Host crash ga olib kelishi mumkin

### Performance issues

- Router va host murakkabligi oshadi
- DoS attack ga mo'yallik

## IPv6 ning afzalliklari

IPv6 da fragmentation **router level** da yo'q qilingan:

- Faqat source host fragmentga bo'lishi mumkin
- Network core yanada sodda
- Security risklari kamaygan
- Packet processing tezligi oshgan

### Xulosa

IP protocol Internet ning asosi bo'lib, addressing, routing va error handling ni ta'minlaydi. Fragmentation IPv4 ning muhim xususiyati, lekin IPv6 da optimized qilingan. Zamonaviy network arxitekturasida control plane ni data plane dan ajratish tendentsiyasi kuchaymoqda.

---

Matn asosida **IPv4 Addressing** haqida eng muhim ma'lumotlarni bayon qilaman:

## IP Address tuzilishi

### Asosiy xususiyatlar

- Har bir **interface** (host yoki router) ning o'ziga xos IP addressi bor
- **32 bit** uzunlikda (4 byte)
- Jami **2^32 ≈ 4 milliard** address mavjud
- **Dotted decimal** formatda yoziladi: 193.32.216.9

### Binary ko'rinishi

```
193.32.216.9 = 11000001 00100000 11011000 00001001
```

## Subnet tushunchasi

### Subnet nima?

**Subnet** - routersiz bir-biriga ulangan interfacеlar to'plami

### Subnet aniqlash qoidasi:

1. Har bir interfacеni host yoki routerdan ajrating
2. Paydo bo'lgan **isolated network** lar - subnet
3. Har bir subnet **subnet address** ga ega

### Misol:

- Subnet: **223.1.1.0/24**
- **/24** - subnet mask (24 ta leftmost bit subnet uchun)
- Host addresslar: 223.1.1.1, 223.1.1.2, 223.1.1.3
- Router interface: 223.1.1.4

## CIDR (Classless Interdomain Routing)

### Format: **a.b.c.d/x**

- **x** - network prefix uzunligi
- Network part: birinchi x bit
- Host part: qolgan (32-x) bit

### Afzalliklari

- **Route aggregation** (маршрутизация umumlashtirish)
- Routing table hajmini kamaytirish
- Hierarchical addressing

### Address aggregation misoli:

```
Provider: 200.23.16.0/20
├── Org 0: 200.23.16.0/23
├── Org 1: 200.23.18.0/23  
├── Org 2: 200.23.20.0/23
└── Org 7: 200.23.30.0/23
```

## Address olish jarayoni

### Organizatsiya uchun

1. **ISP** dan address block olish
2. **ICANN** orqali global boshqaruv
3. **Regional registry** lar:
    - ARIN (Shimoliy Amerika)
    - RIPE (Evropa, Yaqin Sharq)
    - APNIC (Osiyo-Tinch okeani)
    - LACNIC (Lotin Amerika)

### Host uchun - DHCP Protocol

## DHCP (Dynamic Host Configuration Protocol)

### 4 bosqichli jarayon:

**1. DHCP DISCOVER**

- Client: broadcast message (0.0.0.0 → 255.255.255.255)
- UDP port 67 ga yuboriladi

**2. DHCP OFFER**

- Server: IP address taklifi
- Lease time bilan birga

**3. DHCP REQUEST**

- Client: tanlangan serverni tasdiqlaydi

**4. DHCP ACK**

- Server: final tasdiqlash

### DHCP ning afzalliklari

- **Plug-and-play** functionality
- Automatic configuration
- Mobile device lar uchun ideal
- Address pool management

## NAT (Network Address Translation)

### Muammo

- Private network lar bir xil address space ishlatadi
- Public Internet da unique address kerak

### Yechim - NAT

- **Private address space**: 10.0.0.0/8, 192.168.0.0/16
- NAT router **port number** ham o'zgartiradi
- **Translation table** orqali boshqariladi

### NAT ishlash misoli:

```
Internal: (10.0.0.1, 3345) → External: (138.76.29.7, 5001)
```

### NAT Translation Table:

|Private Side|Public Side|
|---|---|
|10.0.0.1:3345|138.76.29.7:5001|

## NAT muammolari

### P2P Application lar

- **Incoming connection** qabul qila olmaydi
- **Server role** o'ynay olmaydi

### Yechimlar:

1. **Connection reversal**
2. **UPnP Protocol**
3. **Application relay**

## UPnP (Universal Plug and Play)

### Vazifasi

- NAT ni automatic configure qilish
- Port mapping yaratish
- External host larga access berish

### Misol:

```
Internal: BitTorrent (10.0.0.1:3345)
UPnP mapping: (138.76.29.7:5001) → (10.0.0.1:3345)
```

## Tarixiy context

### Class-based addressing (eski)

- **Class A**: /8 (16 million host)
- **Class B**: /16 (65K host)
- **Class C**: /24 (254 host)

**Muammo**: Address space ning noto'g'ri ishlatilishi

### CIDR (zamonaviy)

- Har qanday prefix length (/1 dan /30 gacha)
- Optimal address allocation
- Better route aggregation

### Xulosa

IPv4 addressing - Internet ning asosi. CIDR, DHCP va NAT texnologiyalari address scarcity muammosini hal qilishda muhim rol o'ynaydi. IPv6 ga o'tish davom etmoqda, lekin IPv4 hali ham keng qo'llaniladi.

---

