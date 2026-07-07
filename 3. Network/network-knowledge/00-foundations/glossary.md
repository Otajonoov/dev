# Network Glossariy — Atamalar lug'ati

> Network sohasidagi 80+ asosiy atama, alfavit tartibida. Har bir atama uchun: qisqa izoh + misol + cross-reference.

**Tezkor navigatsiya:** [A](#a) · [B](#b) · [C](#c) · [D](#d) · [E](#e) · [F](#f) · [G](#g) · [H](#h) · [I](#i) · [J](#j) · [L](#l) · [M](#m) · [N](#n) · [O](#o) · [P](#p) · [Q](#q) · [R](#r) · [S](#s) · [T](#t) · [U](#u) · [V](#v) · [W](#w)

---

## A

### ACK (Acknowledgment)
TCP'da qabul qilingan packet uchun tasdiqlash signali. Sender packet yuboradi, receiver javob sifatida ACK qaytaradi — bu ishonchli uzatishni ta'minlaydi.
- **Misol:** Three-way handshake oxirgi qadami: SYN → SYN-ACK → **ACK**
- **Layer:** L4 (Transport)
- **Cross-ref:** [../osi/04-transport.md](../osi/04-transport.md)

### ARP (Address Resolution Protocol)
IP address'ni MAC address'ga aylantiruvchi protocol. LAN ichida packet jo'natishdan oldin "192.168.1.1 ning MAC'i nima?" deb so'raydi.
- **Misol:** `arp -a` Linux komandasi ARP jadvalni ko'rsatadi
- **Layer:** L2/L3 oraliq
- **Cross-ref:** [../osi/02-data-link.md](../osi/02-data-link.md)

### Anycast
Bir IP address'ga bir nechta server javob beradi — eng yaqindagisi javob qaytaradi. CDN va DNS root server'larda ishlatiladi.
- **Misol:** Cloudflare 1.1.1.1 — dunyo bo'ylab 300+ joydan javob keladi
- **Cross-ref:** [../osi/04-transport.md](../osi/04-transport.md)

### API (Application Programming Interface)
Dasturlar bir-biri bilan gaplashish uchun belgilangan qoidalar to'plami. Network kontekstida socket API yoki REST API.
- **Misol:** `socket()`, `bind()`, `listen()` — POSIX socket API

---

## B

### Bandwidth
Network'da vaqt birligida o'tkazilishi mumkin bo'lgan maksimal ma'lumot hajmi. Bit/s yoki byte/s'da o'lchanadi.
- **Misol:** 100 Mbit/s home internet — sekundiga 100 million bit
- **Eslatma:** Bandwidth ≠ throughput (real tezlik kichikroq bo'ladi)

### BGP (Border Gateway Protocol)
Internet'ning **routing** protokollari ichida eng muhimi. ASN'lar (Autonomous System) o'rtasida route ma'lumotlarini almashadi.
- **Misol:** UzTelecom o'z ASN'i orqali Google ASN'iga qanday yetishni BGP orqali biladi
- **Layer:** L7 (TCP ustida ishlaydi, lekin routing layeriga oid)
- **Cross-ref:** [../deep-dives/routing-protocols.md](../deep-dives/routing-protocols.md)

### Broadcast
Bitta xabar — barcha qurilmalarga yetkaziladi. Faqat LAN ichida ishlaydi.
- **Misol:** ARP request — "Kim 192.168.1.1?" deb broadcast'da so'raladi
- **IPv4 broadcast address:** `255.255.255.255`
- **Eslatma:** IPv6'da broadcast yo'q — multicast ishlatiladi

---

## C

### CIDR (Classless Inter-Domain Routing)
IP address bilan subnet mask'ni qisqa formatda yozish: `192.168.1.0/24`. `/24` — birinchi 24 bit network qismi.
- **Misol:** `10.0.0.0/8` = 16 million IP, `192.168.1.0/24` = 256 IP
- **Cross-ref:** [../deep-dives/subnetting-cidr.md](../deep-dives/subnetting-cidr.md)

### Client-Server
Network arxitekturasi — bir taraf (client) so'roq qiladi, ikkinchi (server) javob beradi.
- **Misol:** Browser (client) → google.com (server)
- **Cross-ref:** [./what-is-network.md](./what-is-network.md)

### Congestion
Network'da haddan ortiq traffic — packet'lar router queue'da yig'ilib, drop bo'lishi mumkin.
- **TCP javob:** Slow start, congestion avoidance algoritmlari
- **Cross-ref:** [../deep-dives/tcp-handshake.md](../deep-dives/tcp-handshake.md)

---

## D

### DHCP (Dynamic Host Configuration Protocol)
Network'ga ulangan qurilmaga avtomatik IP address, gateway, DNS sozlashlarini beradi.
- **Misol:** Wi-Fi'ga ulanganingda telefoning DHCP server'dan IP oladi
- **Bosqichlar:** DISCOVER → OFFER → REQUEST → ACK (DORA)
- **Layer:** L7 (UDP port 67/68)
- **Cross-ref:** [../tcp-ip/01-network-access.md](../tcp-ip/01-network-access.md)

### DNS (Domain Name System)
Domain nomini (`google.com`) IP address'ga (`142.250.74.46`) aylantiradi. Internetning telefon kitobi.
- **Misol:** `nslookup google.com`, `dig +short cloudflare.com`
- **Layer:** L7 (UDP/TCP port 53)
- **Cross-ref:** [../deep-dives/dns-resolution.md](../deep-dives/dns-resolution.md)

### DNS over HTTPS (DoH)
DNS so'rovlarni HTTPS ichida shifrlangan holda yuborish. Privacy uchun.
- **Misol:** Cloudflare DoH endpoint: `https://1.1.1.1/dns-query`

### Datagram
UDP'da packet'ning nomi. Connection'siz, ishonchsiz uzatma birligi.
- **Cross-ref:** [../osi/04-transport.md](../osi/04-transport.md)

---

## E

### Encapsulation
Yuqori layer ma'lumotini quyi layer header bilan o'rab, packet yasash jarayoni.
- **Misol:** HTTP message → +TCP header → +IP header → +Ethernet header
- **Cross-ref:** [./osi-vs-tcpip.md](./osi-vs-tcpip.md)

### Ethernet
LAN'da eng keng tarqalgan L2 standarti. Cable orqali frame'lar uzatiladi.
- **Standart:** IEEE 802.3
- **Tezliklar:** 10 Mbit, 100 Mbit (Fast Ethernet), 1 Gbit, 10 Gbit, 100 Gbit
- **Layer:** L2

### EIGRP (Enhanced Interior Gateway Routing Protocol)
Cisco'ning ichki tarmoq routing protokoli. OSPF'ga muqobil.

---

## F

### Firewall
Network trafficsi orasidan kerakli'larni o'tkazib, kerakmas'larni bloklaydigan tizim.
- **Turlari:** Stateful (TCP state tracking), stateless (faqat packet header), L7 (deep packet inspection)
- **Misol:** Linux'da `iptables`, `nftables`, `ufw`

### Frame
L2 (Data Link) layeri'da PDU. Ethernet header + IP packet + FCS = frame.
- **Maksimal hajm:** 1500 byte payload (standart MTU)

### FTP (File Transfer Protocol)
Fayl uzatish uchun eski protokol. Port 21 (control), port 20 (data).
- **Eslatma:** Shifrlanmagan — bugun SFTP yoki HTTPS afzalroq
- **Layer:** L7

### FCS (Frame Check Sequence)
Ethernet frame'ning oxiridagi 4-byte checksum. Frame buzilganligini aniqlaydi.

---

## G

### Gateway
Bir network'dan boshqa network'ga o'tish nuqtasi. Odatda bu — sening router'ing.
- **Misol:** Sening uy router'ing default gateway = `192.168.1.1`
- **Komanda:** `ip route` Linux'da default gateway'ni ko'rsatadi

---

## H

### Handshake
Ulanish o'rnatish jarayoni — ikki taraf "salom" deyishadi.
- **TCP three-way handshake:** SYN → SYN-ACK → ACK
- **TLS handshake:** ClientHello → ServerHello → certificate → key exchange → finished
- **Cross-ref:** [../osi/04-transport.md](../osi/04-transport.md)

### HTTP (HyperText Transfer Protocol)
Web sahifalarni yuklash uchun protocol. Stateless, request-response model.
- **Versiyalar:** HTTP/1.0, HTTP/1.1, HTTP/2 (binary, multiplex), HTTP/3 (QUIC ustida)
- **2026 holati:** HTTP/3 — global traffic'ning 35% (Cloudflare ma'lumoti)
- **Layer:** L7 (TCP port 80, HTTP/3 — UDP port 443)

### HTTPS (HTTP Secure)
HTTP + TLS. Traffic shifrlanadi, server identifikatsiya qilinadi.
- **Layer:** L7 (TCP port 443)
- **Cross-ref:** [../deep-dives/tls-ssl.md](../deep-dives/tls-ssl.md)

### Hop
Packet bir router'dan boshqasiga o'tish — bitta hop. `traceroute` har hopni ko'rsatadi.
- **Misol:** Toshkent → Frankfurt: ~12 hop

### Hub
Eski L1 qurilma — bitta portga kelgan signalni boshqa hamma portlarga yuboradi. Bugun deyarli ishlatilmaydi (switch o'rniga keldi).

---

## I

### ICMP (Internet Control Message Protocol)
Network'da xato xabarlari va diagnostika uchun. `ping` va `traceroute` ICMP'ni ishlatadi.
- **Misol:** "Destination Unreachable", "Time Exceeded"
- **Layer:** L3

### IETF (Internet Engineering Task Force)
Internet standartlarini ishlab chiquvchi tashkilot. RFC'larni chiqaradi.
- **Misol:** RFC 793 (TCP), RFC 9114 (HTTP/3) — IETF tomonidan yozilgan

### IP (Internet Protocol)
Network layer'ning asosiy protokoli. IP address bilan packet'larni yetkazib berish vazifasi.
- **Versiyalar:** IPv4 (32 bit), IPv6 (128 bit)
- **Layer:** L3

### IPv4
32-bit address. Maksimal ~4.3 milliard address. Allaqachon tugagan (IANA 2011-yil).
- **Misol:** `192.168.1.1`, `8.8.8.8`
- **Format:** `A.B.C.D` har bir 0-255

### IPv6
128-bit address. 2^128 = ulkan miqdor. 2026-yilda **global adoption 50%+** (Google statistikasi).
- **Misol:** `2001:db8::1`, `fe80::1`
- **Cross-ref:** [../osi/03-network.md](../osi/03-network.md)

### IPsec (IP Security)
IP layer'da shifrlash va autentifikatsiya. VPN'larda asosiy texnologiya.
- **Misol:** Site-to-site VPN, IKEv2

### ISP (Internet Service Provider)
Internet xizmatini ko'rsatuvchi kompaniya. Tier 1/2/3 darajalari.
- **O'zbekiston misol:** UzTelecom, Beeline, Ucell, Uzonline

---

## J

### Jitter
Packet'larning kelishi orasidagi vaqt oralig'ining o'zgarishi (variance). Voice/video uchun muhim.
- **Misol:** Birinchi packet 20ms, ikkinchi 25ms, uchinchi 18ms — bu jitter
- **Cross-ref:** [../osi/04-transport.md](../osi/04-transport.md)

---

## L

### LAN (Local Area Network)
Lokal hududdagi network: uy, ofis, maktab. Odatda Ethernet va Wi-Fi.
- **Cross-ref:** [./what-is-network.md](./what-is-network.md)

### Latency
Packet'ning A nuqtadan B nuqtaga yetib borish uchun ketadigan vaqt.
- **Misol:** Toshkent → Frankfurt round-trip ~80ms
- **`ping` ko'rsatadi:** RTT (Round Trip Time)

### Load Balancer
Traffic'ni bir nechta server o'rtasida taqsimlovchi qurilma yoki dastur.
- **L4 LB:** TCP/UDP darajasida (NLB)
- **L7 LB:** HTTP darajasida (ALB, nginx, HAProxy)

---

## M

### MAC address (Media Access Control)
Tarmoq kartasi'ning fizik address'i. 48 bit, hex'da yoziladi.
- **Misol:** `aa:bb:cc:dd:ee:ff`
- **Layer:** L2
- **Komanda:** `ip link show`, `ifconfig`

### MAN (Metropolitan Area Network)
Shahar miqyosidagi tarmoq.

### MTU (Maximum Transmission Unit)
Bir frame'da yuborilishi mumkin bo'lgan maksimal byte miqdori.
- **Standart Ethernet:** 1500 byte
- **Jumbo frame:** 9000 byte (data center'larda)
- **Komanda:** `ip link show` — har interface MTU'sini ko'rsatadi

### Multicast
Bitta packet — ma'lum guruhga yuboriladi (broadcast'dan tor, unicast'dan keng).
- **Misol:** IPTV, video streaming
- **IPv4 range:** `224.0.0.0/4`

---

## N

### NAT (Network Address Translation)
Bir IP address'ni boshqasiga aylantirish. Ko'pchilikga IPv4 yetishmagani uchun ishlatiladi.
- **Misol:** Uydagi 192.168.1.X (private) → ISP'dan olingan public IP
- **Turlari:** SNAT, DNAT, PAT (Port Address Translation)
- **Cross-ref:** [../deep-dives/nat-and-firewall.md](../deep-dives/nat-and-firewall.md)

### NIC (Network Interface Card)
Tarmoq kartasi — kompyuterning network'ga ulanish qurilmasi.
- **Misol:** Ethernet kartasi, Wi-Fi adapter

---

## O

### OSI (Open Systems Interconnection)
ISO tomonidan ishlab chiqilgan 7-layerli network model (1984). Bugun konseptual ravishda ishlatiladi.
- **Cross-ref:** [./osi-vs-tcpip.md](./osi-vs-tcpip.md)

### OSPF (Open Shortest Path First)
IGP (Interior Gateway Protocol) — bir AS ichidagi routing. Dijkstra algoritmidan foydalanadi.
- **Layer:** L3
- **Cross-ref:** [../deep-dives/routing-protocols.md](../deep-dives/routing-protocols.md)

---

## P

### Packet
L3 (Network) layeri'da PDU. IP packet = IP header + TCP segment.
- **Cross-ref:** [./osi-vs-tcpip.md](./osi-vs-tcpip.md)

### PAN (Personal Area Network)
Shaxsiy yaqin tarmoq: Bluetooth, NFC.

### Payload
Header'siz toza ma'lumot qismi. "Konvert ichidagi xat".
- **Misol:** TCP segment'da: TCP header + **payload** (HTTP message)

### Ping
ICMP Echo Request/Reply yuborib, host yetishish va RTT'ni o'lchaydigan tool.
- **Komanda:** `ping google.com`

### Port
Transport layer'da application'ni identifikatsiya qiluvchi 16-bit raqam (0-65535).
- **Mashhur portlar:** 22 (SSH), 53 (DNS), 80 (HTTP), 443 (HTTPS), 3306 (MySQL), 5432 (PostgreSQL)
- **Komanda:** `ss -tlnp`, `netstat -an`

### Proxy
Foydalanuvchi va server o'rtasida vositachi server.
- **Forward proxy:** Foydalanuvchi tarafi (corporate proxy)
- **Reverse proxy:** Server tarafi (nginx, Cloudflare)

---

## Q

### QUIC
Google ishlab chiqgan, HTTP/3 asosida bo'lgan transport protokoli. UDP ustida ishlaydi, lekin TCP'ga o'xshash kafolatlar beradi.
- **Plus:** 0-RTT connection, multiplexing, head-of-line blocking yo'q
- **Layer:** L4 (UDP port 443)
- **Cross-ref:** [../deep-dives/http-evolution.md](../deep-dives/http-evolution.md)

---

## R

### RIP (Routing Information Protocol)
Eski va sodda routing protokol. Distance-vector. Bugun deyarli ishlatilmaydi (OSPF, EIGRP afzalroq).

### RFC (Request For Comments)
IETF'ning rasmiy hujjatlari. Internet standartlari shu RFC'larda yoziladi.
- **Misol:** RFC 793 (TCP), RFC 791 (IP), RFC 9110 (HTTP)
- **Saytda:** [datatracker.ietf.org](https://datatracker.ietf.org/)

### Router
L3 qurilma. IP packet'larni turli network'lar o'rtasida yo'naltiradi.
- **Komanda:** `traceroute` har bir router'ni ko'rsatadi

### Routing
Packet uchun manzilga eng yaxshi yo'lni tanlash jarayoni.
- **Cross-ref:** [../deep-dives/routing-protocols.md](../deep-dives/routing-protocols.md)

### RTT (Round Trip Time)
Packet borish + qaytish uchun ketadigan vaqt. `ping` ko'rsatadi.

---

## S

### Segment
L4 (TCP) PDU. UDP'da bu — datagram.

### SMTP (Simple Mail Transfer Protocol)
Email yuborish protokoli.
- **Layer:** L7 (TCP port 25, 587)

### Socket
OS API — application'ning network'dagi endpoint'i. `IP:port` jufti.
- **Misol Go'da:** `net.Listen("tcp", ":8080")` — TCP socket ochadi

### SSH (Secure Shell)
Shifrlangan remote login protokoli. Telnet'ning xavfsiz o'rinbosari.
- **Layer:** L7 (TCP port 22)
- **Komanda:** `ssh user@host`

### SSL (Secure Sockets Layer)
TLS'ning eski versiyasi (SSL 2.0, 3.0). Bugun ishlatilmaydi (deprecated). Lekin ko'pincha "SSL certificate" deb ataladi (aslida TLS).

### Subnet
Network'ning kichik bo'lagi. Mask bilan aniqlanadi.
- **Misol:** `192.168.1.0/24` — 256 IP'lik subnet
- **Cross-ref:** [../deep-dives/subnetting-cidr.md](../deep-dives/subnetting-cidr.md)

### Switch
L2 qurilma. MAC address'larga qarab frame'larni faqat kerakli portga yuboradi.
- **Farqi hub'dan:** Hub — broadcast, switch — selective

---

## T

### TCP (Transmission Control Protocol)
Ishonchli, tartibli, connection-oriented transport protokoli.
- **Three-way handshake:** SYN → SYN-ACK → ACK
- **Layer:** L4
- **Cross-ref:** [../deep-dives/tcp-handshake.md](../deep-dives/tcp-handshake.md)

### TLS (Transport Layer Security)
SSL'ning zamonaviy versiyasi. HTTPS, SMTPS, IMAPS uchun.
- **Versiyalar:** TLS 1.2 (2008), TLS 1.3 (2018, faster handshake)
- **Cross-ref:** [../deep-dives/tls-ssl.md](../deep-dives/tls-ssl.md)

### Throughput
Real uzatilayotgan ma'lumot tezligi (bandwidth'dan kichikroq).
- **Misol:** 100 Mbit bandwidth, lekin throughput 80 Mbit (overhead, retransmission)

### TTL (Time To Live)
Packet'ning maksimal hop soni. Har router 1 ga kamaytiradi. 0 bo'lsa — packet drop bo'ladi.
- **Standart qiymat:** Linux 64, Windows 128
- **`traceroute`:** TTL'ni manipulyatsiya qilib, har router'ni topadi

---

## U

### UDP (User Datagram Protocol)
Connectionless, unreliable, fast transport protokoli. TCP'dan oldindan handshake yo'q.
- **Foydalanish:** DNS, video streaming, gaming, VoIP, QUIC
- **Layer:** L4
- **Cross-ref:** [../osi/04-transport.md](../osi/04-transport.md)

### Unicast
Bitta sender → bitta receiver. Eng oddiy uzatma turi.

### URL (Uniform Resource Locator)
Internet'da resursning manzili. `https://example.com/path?query=1`
- **Qismlar:** scheme://user:pass@host:port/path?query#fragment

---

## V

### VLAN (Virtual LAN)
Bitta fizik switch'da bir nechta mantiqiy LAN yaratish.
- **Standart:** IEEE 802.1Q
- **Misol:** 10-VLAN — engineering, 20-VLAN — sales

### VPN (Virtual Private Network)
Public network ustidan shifrlangan tunnel. Sening traffiking VPN server orqali o'tadi.
- **Turlari:** IPsec, OpenVPN, WireGuard
- **Cross-ref:** [../deep-dives/tls-ssl.md](../deep-dives/tls-ssl.md)

### VoIP (Voice over IP)
Internet orqali ovoz uzatish. Telegram audio call, Zoom, Skype — VoIP.

---

## W

### WAN (Wide Area Network)
Keng hududdagi tarmoq: davlat, kontinent, dunyo.

### Wi-Fi
Simsiz LAN texnologiyasi. IEEE 802.11 standarti.
- **Versiyalar:** 802.11n (Wi-Fi 4), 802.11ac (Wi-Fi 5), 802.11ax (Wi-Fi 6/6E), 802.11be (Wi-Fi 7)
- **Layer:** L1 + L2

### Wireshark
Network packet'larni capture va analiz qilish uchun GUI tool. CLI versiyasi — `tshark`.
- **Komanda:** `tshark -i any -Y 'http'`

---

## Yodda saqlash

- Atamalarni **kontekstda** o'rganish — formula yoddan emas, qaerda ishlatilishi muhim
- Inglizcha qisqartmalar (TCP, IP, DNS, NAT, ARP) — **tarjima qilinmaydi**
- Har bir atama ortida bitta layer va bitta vazifa turibdi — **shu juftlikni eslab qol**
- RFC'lar — har bir protokol uchun "rasmiy kitob"

---

## Cross-references

- Network nima — umumiy tushuncha: [./what-is-network.md](./what-is-network.md)
- OSI vs TCP/IP modellar: [./osi-vs-tcpip.md](./osi-vs-tcpip.md)
- OSI layerlar: [../osi/](../osi/)
- TCP/IP layerlar: [../tcp-ip/](../tcp-ip/)
- Deep-dive mavzular: [../deep-dives/](../deep-dives/)

---

## Manbalar

- **Kitob:** Kurose & Ross, *Computer Networking: A Top-Down Approach*, 6-nashr (atamalar 1-bobdan boshlanadi)
- **RFC index:** [datatracker.ietf.org](https://datatracker.ietf.org/)
- **MDN Glossary:** [developer.mozilla.org/en-US/docs/Glossary](https://developer.mozilla.org/en-US/docs/Glossary)
- **Cloudflare Learning Center:** [cloudflare.com/learning/](https://www.cloudflare.com/learning/)
- **HTTP/3 statistikasi:** [HTTP/3 — Wikipedia](https://en.wikipedia.org/wiki/HTTP/3)
- **IPv6 deployment 2026:** [Statistics on the Adoption of IPv6 — Internet Society](https://www.internetsociety.org/deploy360/ipv6/statistics/)
