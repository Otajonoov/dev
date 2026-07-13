# 09. Glossary — Tarmoq atamalari lug'ati

> Tarmoq asoslari bo'yicha 90+ asosiy atama, alifbo tartibida. Har atama uchun:
> qisqa izoh + misol. Bu dars emas — tezkor ma'lumotnoma. Atama uchratganda
> shu yerdan qidir.

**Tezkor navigatsiya:** [A](#a) · [B](#b) · [C](#c) · [D](#d) · [E](#e) ·
[F](#f) · [G](#g) · [H](#h) · [I](#i) · [J](#j) · [L](#l) · [M](#m) · [N](#n) ·
[O](#o) · [P](#p) · [Q](#q) · [R](#r) · [S](#s) · [T](#t) · [U](#u) · [V](#v) · [W](#w)

---

## A

### ACK (Acknowledgment)
TCP'da qabul qilingan ma'lumot uchun tasdiqlash signali. Ishonchli uzatishni ta'minlaydi.
- **Misol:** three-way handshake: SYN → SYN-ACK → **ACK**. Layer: 4.

### Access network
Foydalanuvchi qurilmasini ISP'ning birinchi router'iga ulaydigan fizik tarmoq.
- **Misol:** DSL, fiber (FTTH), 5G. Batafsil: [03-access-networks](03-access-networks.md).

### API (Application Programming Interface)
Dasturlar bir-biri bilan gaplashish uchun belgilangan qoidalar to'plami.
- **Misol:** POSIX socket API: `socket()`, `bind()`, `listen()`.

### ARP (Address Resolution Protocol)
IP address'ni MAC address'ga aylantiruvchi protokol. LAN ichida ishlaydi.
- **Misol:** `arp -a` ARP jadvalni ko'rsatadi. Layer: 2/3 oraliq.

### ARPANET
Internetning bobosi (1969). Birinchi packet-switched tarmoq, 4 universitet ulangan.

### AS (Autonomous System)
Mustaqil boshqariladigan tarmoq birligi. Har ISP o'z AS raqamiga ega.
- **Misol:** BGP AS'lar orasida route almashadi.

### Anycast
Bir IP address'ga bir nechta server javob beradi — eng yaqindagisi. CDN va DNS'da.
- **Misol:** Cloudflare 1.1.1.1 — dunyo bo'ylab ko'p joydan javob.

---

## B

### Backbone
Internetning yuqori tezlikdagi asosiy magistral tarmog'i (Tier 1 ISP'lar).

### Bandwidth
Vaqt birligida o'tkazilishi mumkin bo'lgan **maksimal nazariy** ma'lumot hajmi.
- **Misol:** 100 Mbit/s. Eslatma: bandwidth ≠ throughput.

### BGP (Border Gateway Protocol)
Internet'ning eng muhim routing protokoli. AS'lar orasida route ma'lumotini almashadi.

### Bit
Ma'lumotning eng kichik birligi (0 yoki 1). Physical qatlamdagi PDU.

### Broadcast
Bitta xabar — LAN'dagi barcha qurilmaga yetkaziladi.
- **IPv4 broadcast:** `255.255.255.255`. IPv6'da broadcast yo'q (multicast bor).

### Bottleneck
Yo'lning eng sekin bo'g'ini — u butun throughput'ni cheklaydi.

---

## C

### Cable (HFC)
Kabel televidenie liniyalari orqali internet. Optik + koaksial. Standart: DOCSIS.
- **Misol:** DOCSIS 4.0 — 10 Gbit/s down. Baham ko'riladi.

### CDN (Content Delivery Network)
Kontent nusxalarini foydalanuvchilarga yaqin serverlarga tarqatuvchi tarmoq.
- **Misol:** Cloudflare, Akamai. Latencyni kamaytiradi.

### CIDR (Classless Inter-Domain Routing)
IP + subnet mask'ni qisqa yozish: `192.168.1.0/24`.

### Circuit switching
Yo'lni oldindan band qilib uzatish (eski telefon). Packet switching'ning raqibi.
- **Batafsil:** [04-network-core-va-packet-switching](04-network-core-va-packet-switching.md).

### Client-Server
Bir taraf so'raydi (client), ikkinchi javob beradi (server).
- **Misol:** Browser → google.com.

### Congestion
Tarmoqda haddan ortiq trafik — paketlar navbatda yig'ilib, drop bo'lishi mumkin.

---

## D

### Datagram
UDP'da paketning nomi. Connection'siz, ishonchsiz uzatma birligi.

### Decapsulation
Qabul qiluvchi tarafda har qatlam o'z header'ini olib tashlashi (encapsulation teskarisi).

### DHCP (Dynamic Host Configuration Protocol)
Qurilmaga avtomatik IP, gateway, DNS beradi.
- **Bosqichlar:** DISCOVER → OFFER → REQUEST → ACK (DORA). Port 67/68.

### DNS (Domain Name System)
Domen nomini (`google.com`) IP'ga (`142.250.74.46`) aylantiradi. Internetning telefon kitobi.
- **Misol:** `nslookup google.com`. Port 53.

### DoD model
TCP/IP modelining original nomi (US Department of Defense).

### DSL (Digital Subscriber Line)
Eski telefon (mis) simlari orqali internet. 12-24 Mbit/s. Sekin, arzon.

---

## E

### Encapsulation
Har qatlam yuqoridan kelgan ma'lumotni o'z header'i bilan o'rashi (matryoshka).
- **Batafsil:** [08-tcpip-modeli-va-encapsulation](08-tcpip-modeli-va-encapsulation.md).

### End-to-End Delay
Manbadan manzilgacha to'liq kechikish (barcha hop'lar yig'indisi).

### Ethernet
LAN'da eng keng tarqalgan Layer 2 standarti. Kabel orqali frame uzatadi.
- **Standart:** IEEE 802.3. Tezliklar: 100 Mbit - 100 Gbit.

### Extranet
Tashkilot + ishonchli sheriklar foydalanadigan tarmoq (B2B portal).

---

## F

### FCS (Frame Check Sequence)
Ethernet frame oxiridagi 4-byte checksum. Frame buzilganini aniqlaydi.

### FDM (Frequency-Division Multiplexing)
Kanalni chastota bo'yicha bo'lish (radiostansiyalar kabi). Circuit switching'da.

### Firewall
Kerakli trafikni o'tkazib, kerakmasini bloklaydigan tizim.
- **Misol:** Linux `iptables`, `nftables`, `ufw`.

### Frame
Layer 2 (Data Link) PDU. Ethernet header + IP packet + FCS.
- **Maksimal:** 1500 byte payload (standart MTU).

### FTTH (Fiber To The Home)
Optik tola to'g'ridan-to'g'ri uyga. Eng tez, eng past latency (1-4 ms).
- **Standartlar:** GPON, XGS-PON.

### FTP (File Transfer Protocol)
Fayl uzatish uchun eski protokol. Port 21. Shifrlanmagan (SFTP afzalroq).

### FWA (Fixed Wireless Access)
Uyga o'rnatilgan 5G/radio modem orqali internet. Sim tortmasdan.

---

## G

### Gateway
Bir tarmoqdan boshqasiga o'tish nuqtasi. Odatda uy router'ing.
- **Misol:** default gateway = `192.168.1.1`.

### GPON / XGS-PON
FTTH fiber standartlari. GPON ~900 Mbit/s, XGS-PON 2.5-10 Gbit/s simmetrik.

---

## H

### Handshake
Ulanish o'rnatish jarayoni — ikki taraf "salom" deyishadi.
- **TCP:** SYN → SYN-ACK → ACK. **TLS:** ClientHello → ServerHello → ...

### Header
Qatlam qo'shadigan xizmat ma'lumoti (manzil, port, tartib raqami kabi).

### Hop
Paketning bir router'dan boshqasiga o'tishi. `traceroute` har hopni ko'rsatadi.
- **Misol:** Toshkent → Frankfurt: ~12 hop.

### HTTP (HyperText Transfer Protocol)
Web sahifalarni yuklash protokoli. Request-response, stateless.
- **Versiyalar:** HTTP/1.1, HTTP/2, HTTP/3 (QUIC ustida). Port 80.

### HTTPS
HTTP + TLS. Trafik shifrlanadi, server tekshiriladi. Port 443.

### Hub
Eski Layer 1 qurilma — kelgan signalni hamma portga yuboradi. Switch o'rniga keldi.

---

## I

### ICMP (Internet Control Message Protocol)
Xato xabarlari va diagnostika. `ping` va `traceroute` ICMP ishlatadi. Layer: 3.

### IETF (Internet Engineering Task Force)
Internet standartlarini (RFC) ishlab chiquvchi tashkilot.

### IP (Internet Protocol)
Network/Internet qatlamining asosiy protokoli. IP address bilan paketlarni yetkazadi.
- **Versiyalar:** IPv4 (32 bit), IPv6 (128 bit). Layer: 3.

### IPv4
32-bit address, ~4.3 milliard. Allaqachon tugagan.
- **Misol:** `192.168.1.1`, `8.8.8.8`.

### IPv6
128-bit address, ulkan miqdor. 2026-yilda global adoption **50%+**.
- **Misol:** `2001:db8::1`.

### ISP (Internet Service Provider)
Internet xizmatini ko'rsatuvchi kompaniya. Tier 1/2/3 darajalari.
- **Batafsil:** [05-internet-tuzilishi-isp](05-internet-tuzilishi-isp.md).

### IXP (Internet Exchange Point)
ISP'lar uchrashib bepul (peering) ma'lumot almashadigan bino.
- **Misol:** DE-CIX (Frankfurt), TAS-IX (Toshkent).

### Intranet
Faqat tashkilot ichidagi xodimlar uchun tarmoq (HR portali).

---

## J

### Jitter
Paketlar kelishi orasidagi vaqt oralig'ining o'zgaruvchanligi. Video/ovoz uchun muhim.
- **Misol:** 20 ms, 25 ms, 18 ms — bu jitter.

---

## L

### LAN (Local Area Network)
Lokal tarmoq: uy, ofis, maktab. Ethernet va Wi-Fi.

### Last mile (oxirgi mil)
Foydalanuvchi uyigacha bo'lgan oxirgi ulanish qismi — eng qimmat va murakkab.

### Latency
Paketning A dan B ga yetib borish vaqti (kechikish).
- **Misol:** Toshkent → Frankfurt RTT ~80 ms. Batafsil: [06-latency-loss-throughput](06-latency-loss-throughput.md).

### Layer (qatlam)
Tarmoq vazifalarini bo'lish usuli. Har qatlam bitta ish qiladi.

### LEO (Low Earth Orbit)
Past orbitadagi sun'iy yo'ldoshlar (Starlink). GEO'dan ancha past latency (~30 ms).

### Load Balancer
Trafikni bir nechta server orasida taqsimlovchi qurilma/dastur.
- **L4 LB:** IP/port; **L7 LB:** HTTP mazmuni.

---

## M

### MAC address (Media Access Control)
Tarmoq kartasining fizik address'i (48 bit, `aa:bb:cc:dd:ee:ff`). O'zgarmas. Layer: 2.

### MAN (Metropolitan Area Network)
Shahar miqyosidagi tarmoq (optik tola).

### MTU (Maximum Transmission Unit)
Bir frame'da yuborilishi mumkin maksimal byte.
- **Standart Ethernet:** 1500 byte. Jumbo: 9000 byte.

### Multicast
Bitta paket — ma'lum guruhga (broadcast'dan tor, unicast'dan keng).
- **IPv4 range:** `224.0.0.0/4`. Misol: IPTV.

### Multihoming
Bir tarmoq bir nechta ISP'ga ulanadi (zaxira uchun).

---

## N

### NAT (Network Address Translation)
Bir IP'ni boshqasiga aylantirish. IPv4 yetishmasligi tufayli.
- **Misol:** uydagi `192.168.1.X` (private) → ISP public IP.

### Network core
Internetning "ichi", o'zaro bog'langan routerlar to'ri.

### NIC (Network Interface Card)
Tarmoq kartasi — kompyuterning tarmoqqa ulanish qurilmasi.

---

## O

### OSI (Open Systems Interconnection)
ISO'ning 7-qatlamli konseptual modeli (1984). O'quv tili.
- **Batafsil:** [07-osi-modeli](07-osi-modeli.md).

### OSPF (Open Shortest Path First)
Bir AS ichidagi routing protokoli (IGP). Dijkstra algoritmi. Layer: 3.

---

## P

### Packet
Layer 3 (Network) PDU. IP header + TCP segment.

### Packet loss
Router buferi to'lganda paketning tashlab yuborilishi (yo'qolishi).

### Packet switching
Xabarni paketlarga bo'lib mustaqil yo'naltirish. Internet shunday ishlaydi.

### PAN (Personal Area Network)
Shaxsiy yaqin tarmoq: Bluetooth, NFC (1-10 m).

### Payload
Header'siz toza ma'lumot. "Konvert ichidagi xat".

### PDU (Protocol Data Unit)
Ma'lumotning har qatlamdagi nomi: Data → Segment → Packet → Frame → Bits.

### Peering
Teng ISP'lar orasida bepul ma'lumot almashish kelishuvi.

### Ping
ICMP Echo yuborib, host yetishishi va RTT'ni o'lchaydigan tool.
- **Komanda:** `ping google.com`.

### PoP (Point of Presence)
ISP'ning har shahardagi router/server "ofisi".

### Port
Transport qatlamda dasturni identifikatsiya qiluvchi 16-bit raqam (0-65535).
- **Mashhur:** 22 (SSH), 53 (DNS), 80 (HTTP), 443 (HTTPS).

### Propagation delay
Bitning fizik masofani bosib o'tish vaqti (`d/s`).

### Protocol
Qurilmalar o'rtasidagi ma'lumot almashish qoidalari.
- **Batafsil:** [02-protokol-nima](02-protokol-nima.md).

### Proxy
Foydalanuvchi va server o'rtasida vositachi server.

---

## Q

### QoS (Quality of Service)
Muhim trafikka (video call) ustuvorlik berish mexanizmi. Queuing delayni kamaytiradi.

### QUIC
UDP ustida ishlaydigan yangi transport protokol (HTTP/3 poydevori). TCP'ga o'xshash kafolat.
- **Plus:** 0-RTT, multiplexing, head-of-line blocking yo'q.

### Queuing delay
Paketning router buferida navbatda kutish vaqti. Eng o'zgaruvchan delay.

---

## R

### RFC (Request For Comments)
IETF'ning rasmiy protokol hujjatlari. 9000+ mavjud.
- **Misol:** RFC 791 (IP), RFC 9293 (TCP), RFC 9110 (HTTP).

### Router
Layer 3 qurilma. IP paketlarni turli tarmoqlar orasida yo'naltiradi.

### Routing
Paket uchun manzilga eng yaxshi yo'lni tanlash jarayoni.

### Routing table
Router'ning "qaysi IP qaysi yo'nalishga" degan jadvali.

### RTT (Round Trip Time)
Paket borish + qaytish vaqti. `ping` ko'rsatadi.

---

## S

### Segment
Layer 4 (TCP) PDU. UDP'da bu — datagram.

### Session
Ikki taraf o'rtasidagi seans (dialog) — OSI 5-qatlam.

### SMTP (Simple Mail Transfer Protocol)
Email yuborish protokoli. Port 25, 587. Layer: 7.

### Socket
OS API — dasturning tarmoqdagi endpoint'i (`IP:port`).
- **Go misoli:** `net.Listen("tcp", ":8080")`.

### SSH (Secure Shell)
Shifrlangan remote login protokoli. Port 22. Telnet o'rinbosari.

### Starlink
LEO sun'iy yo'ldosh interneti. ~128-200 Mbit/s, ~25-50 ms latency. Uzoq joylar uchun.

### Store-and-forward
Router paketni to'liq oladi, keyin uzatadi.

### Subnet
Tarmoqning kichik bo'lagi, mask bilan aniqlanadi.
- **Misol:** `192.168.1.0/24` = 256 IP.

### Switch
Layer 2 qurilma. MAC address'ga qarab frame'larni faqat kerakli portga yuboradi.
- **Farqi hub'dan:** hub broadcast, switch selective.

---

## T

### TAS-IX
O'zbekistonning asosiy IXP'si (Toshkent). Mahalliy trafik shu yerda almashadi.

### TCP (Transmission Control Protocol)
Ishonchli, tartibli, connection-oriented transport protokoli. Layer: 4.
- **Handshake:** SYN → SYN-ACK → ACK.

### TCP/IP model
Internetning haqiqiy 4 qatlamli modeli.
- **Batafsil:** [08-tcpip-modeli-va-encapsulation](08-tcpip-modeli-va-encapsulation.md).

### TDM (Time-Division Multiplexing)
Kanalni vaqt oralig'i bo'yicha bo'lish. Circuit switching'da.

### Throughput
Real uzatilayotgan ma'lumot tezligi (bandwidth'dan kichik yoki teng).
- **Qoida:** `min{barcha liniyalar}` (bottleneck).

### Tier 1 / 2 / 3
ISP darajalari. Tier 1 — global backbone (hech kimga to'lamaydi).

### TLS (Transport Layer Security)
SSL'ning zamonaviy versiyasi. HTTPS uchun shifrlash.
- **Versiyalar:** TLS 1.2, TLS 1.3.

### Topology
Qurilmalar fizik/mantiqiy joylashuvi: bus, star, ring, mesh, tree.

### Transit
Kichik ISP katta ISP'ga Internetga chiqish uchun to'laydigan kelishuv.

### Transmission delay
Paketni liniyaga "itarish" vaqti (`L/R`).

### TTL (Time To Live)
Paketning maksimal hop soni. Har router 1 ga kamaytiradi. 0 bo'lsa drop.
- **Standart:** Linux 64, Windows 128.

---

## U

### UDP (User Datagram Protocol)
Connectionless, tez, ishonchsiz transport protokoli. Handshake yo'q. Layer: 4.
- **Foydalanish:** DNS, video, gaming, VoIP, QUIC.

### Unicast
Bitta sender → bitta receiver. Eng oddiy uzatma turi.

### URL (Uniform Resource Locator)
Internet resursining manzili.
- **Format:** `scheme://host:port/path?query#fragment`.

---

## V

### VLAN (Virtual LAN)
Bitta fizik switch'da bir nechta mantiqiy LAN. Standart: IEEE 802.1Q.

### VPN (Virtual Private Network)
Public tarmoq ustidan shifrlangan tunnel.
- **Turlari:** IPsec, OpenVPN, WireGuard.

### VoIP (Voice over IP)
Internet orqali ovoz uzatish (Telegram call, Zoom).

---

## W

### WAN (Wide Area Network)
Keng hududdagi tarmoq: davlat, kontinent, dunyo.

### Wi-Fi
Simsiz LAN texnologiyasi. IEEE 802.11.
- **Versiyalar:** Wi-Fi 4/5/6/6E/7. Layer: 1+2.

### Wireshark
Paketlarni ushlab tahlil qiluvchi GUI tool. CLI: `tshark`.

### WWW (World Wide Web)
Internet ustida ishlaydigan bitta xizmat (HTTP orqali web sahifalar). Internet ≠ WWW.

---

## Yodda saqlash

- Atamalarni **kontekstda** o'rgan — qayerda ishlatilishi muhim.
- Inglizcha qisqartmalar (TCP, IP, DNS, NAT, ARP) **tarjima qilinmaydi**.
- Har atama ortida bitta **qatlam** va bitta **vazifa** turibdi — shu juftlikni eslab qol.
- RFC'lar — har protokol uchun "rasmiy kitob".

---

## 📚 Manbalar

- Kurose & Ross, *Computer Networking: A Top-Down Approach*, 6-nashr (1-bob atamalari)
- [MDN Web Docs Glossary — Mozilla](https://developer.mozilla.org/en-US/docs/Glossary)
- [Cloudflare Learning Center](https://www.cloudflare.com/learning/)
- [RFC index — IETF Datatracker](https://datatracker.ietf.org/)
- [IPv6 Adoption 2026 — Internet Society Pulse](https://pulse.internetsociety.org/en/blog/2026/04/18-years-later-ipv6-reaches-majority/)
