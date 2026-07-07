# Network Interview Savol-Javoblar

50+ ta interview savoli — junior, mid, senior darajalarga ajratilgan. Har bir javob real misollar va Linux komandalar bilan.

## Tarkib

- [Junior darajasi](#junior-darajasi-0-2-yil) — 18 savol
- [Mid darajasi](#mid-darajasi-2-5-yil) — 18 savol
- [Senior darajasi](#senior-darajasi-5-yil) — 17 savol

---

## Junior darajasi (0-2 yil)

### Q1. OSI model 7 ta layer ni sanab bering. Har biri nima qiladi?

**Tag:** OSI, basics

**Javob:**

1. **Physical (L1)** — bit larni elektr/optika/radio signaliga aylantirish. Misol: Ethernet kabel, fiber, Wi-Fi radio.
2. **Data Link (L2)** — bitta link ichida frame uzatish, MAC address. Misol: Ethernet, Wi-Fi (802.11), ARP.
3. **Network (L3)** — host'lar o'rtasida packet routing, IP address. Misol: IPv4, IPv6, ICMP.
4. **Transport (L4)** — process-process aloqasi, port. Misol: TCP, UDP.
5. **Session (L5)** — session boshqarish, dialog. Misol: NetBIOS, RPC.
6. **Presentation (L6)** — encoding, encryption, compression. Misol: TLS, JPEG.
7. **Application (L7)** — foydalanuvchi protokollari. Misol: HTTP, DNS, SSH.

Eslab qolish: **A**ll **P**eople **S**eem **T**o **N**eed **D**ata **P**rocessing (7→1).

[OSI README](../osi/README.md)

---

### Q2. TCP va UDP farqi nima?

**Tag:** Transport, TCP, UDP

**Javob:**

| Mezon | TCP | UDP |
|-------|-----|-----|
| Connection | Bor (3-way handshake) | Yo'q |
| Reliability | Ishonchli (ACK, retransmit) | Ishonchsiz |
| Order | Ordered (sequence number) | Tartibsiz |
| Speed | Sekinroq (overhead) | Tezroq |
| Header | 20-60 byte | 8 byte |
| Use case | HTTP, SSH, email | DNS, video, gaming, VoIP |

TCP — ishonchli, lekin overhead bor. UDP — tez, lekin yo'qotish/tartibni dastur o'zi hal qiladi.

[Transport layer](../osi/04-transport.md)

---

### Q3. DNS qanday ishlaydi?

**Tag:** Application, DNS

**Javob:**

DNS — domain name (`google.com`) → IP address (`142.250.190.46`) o'tkaziladi. Bosqichlar:

1. Browser local cache'ni tekshiradi
2. OS resolver (yoki `/etc/hosts`) tekshiradi
3. Configured DNS server'ga so'rov yuboriladi (port 53, UDP)
4. DNS server o'z cache'ini tekshiradi, yo'q bo'lsa **recursive** qidiradi:
   - Root servers (`.`) → TLD (`.com`) → Authoritative (`google.com`)
5. IP qaytariladi, browser shu IP ga TCP ulanishni boshlaydi

```bash
dig google.com
dig +trace google.com   # to'liq resolution chain
```

[DNS deep-dive](../deep-dives/dns-resolution.md)

---

### Q4. MAC address va IP address farqi nima?

**Tag:** Data Link, Network

**Javob:**

- **MAC address (L2):** 48 bit (6 byte), `aa:bb:cc:dd:ee:ff`. NIC vendor tomonidan beriladi (OUI), o'zgarmaydi (deyarli). Faqat **bitta network segment ichida** ishlaydi.
- **IP address (L3):** 32 bit IPv4 yoki 128 bit IPv6. Network admin tomonidan beriladi (yoki DHCP). **Internet bo'ylab** routerlar orqali yo'naltiriladi.

ARP — IP'dan MAC topadi (bir segment ichida).

```bash
ip link show       # MAC
ip addr show       # IP
ip neigh           # ARP table (IP↔MAC)
```

---

### Q5. ping va traceroute farqi?

**Tag:** Network, ICMP

**Javob:**

- **ping** — destination'ga ICMP Echo Request yuboradi va Echo Reply'ni kutadi. RTT (round-trip time) va packet loss'ni ko'rsatadi.
- **traceroute** (Linux'da `tracepath` ham) — har TTL=1, 2, 3, ... bilan UDP/ICMP packet yuboradi. Har router TTL=0 ga yetganida ICMP Time Exceeded qaytaradi. Shu yo'l bilan butun yo'l ko'rinadi.

```bash
ping -c 4 8.8.8.8
traceroute -n google.com
mtr google.com    # ping + traceroute hybrid
```

---

### Q6. HTTP status code 5 ta gruh

**Tag:** Application, HTTP

**Javob:**

- **1xx (Informational):** 100 Continue, 101 Switching Protocols
- **2xx (Success):** 200 OK, 201 Created, 204 No Content
- **3xx (Redirect):** 301 Moved Permanently, 302 Found, 304 Not Modified
- **4xx (Client error):** 400 Bad Request, 401 Unauthorized, 403 Forbidden, 404 Not Found, 429 Too Many Requests
- **5xx (Server error):** 500 Internal Server Error, 502 Bad Gateway, 503 Service Unavailable, 504 Gateway Timeout

[HTTP evolution](../deep-dives/http-evolution.md)

---

### Q7. Subnet mask nima uchun kerak?

**Tag:** Network, addressing

**Javob:**

Subnet mask IP address'ni **network qism** va **host qism**ga ajratadi. Misol: `192.168.1.10/24`:
- Mask `255.255.255.0` — birinchi 24 bit network (`192.168.1`), oxirgi 8 bit host (`.10`).
- Network: `192.168.1.0`, Broadcast: `192.168.1.255`, Host range: `.1`–`.254` (254 host).

Subnet mask routerga "bu IP qaysi network'ga tegishli?" deb aytadi va to'g'ri yo'naltiradi.

[Subnetting deep-dive](../deep-dives/subnetting-cidr.md)

---

### Q8. Private vs Public IP

**Tag:** Network, addressing

**Javob:**

**Private IP** (RFC 1918) — Internet'da routelash mumkin emas, faqat lokal network'da:
- `10.0.0.0/8` (16M host)
- `172.16.0.0/12` (1M host)
- `192.168.0.0/16` (65K host)

**Public IP** — Internet'da unique, ICANN/RIR tomonidan ajratiladi.

NAT — bir nechta private IP'ni bitta public IP orqali Internet'ga chiqaradi.

[NAT deep-dive](../deep-dives/nat-and-firewall.md)

---

### Q9. Switch va Router farqi

**Tag:** Data Link, Network

**Javob:**

- **Switch (L2):** MAC address bo'yicha frame yo'naltiradi, **bitta network segment ichida**. MAC address table yuritadi.
- **Router (L3):** IP address bo'yicha packet yo'naltiradi, **turli network'lar orasida**. Routing table yuritadi.

Modern qurilmalar L3 switch — switch'ning tezligi + router'ning routing imkoniyati bilan.

---

### Q10. ARP nima qiladi?

**Tag:** Data Link

**Javob:**

ARP (Address Resolution Protocol) — IP address'dan MAC address topadi (bitta segment ichida).

1. Host: "192.168.1.1 kimning IP'si?" — ARP request **broadcast** qilinadi
2. Egasi: "U menga tegishli, MAC'im aa:bb:cc:..." — ARP reply **unicast**
3. Cache'ga saqlanadi (`ip neigh`)

ARP spoofing hujum — yolg'on reply yuborib, traffikni yo'naltirish (MitM).

---

### Q11. DHCP qanday ishlaydi (DORA)?

**Tag:** Network Access

**Javob:**

**DORA** — DHCP'ning 4 bosqichi:

1. **Discover** — Client broadcast: "DHCP server bormi?" (UDP 67/68)
2. **Offer** — Server: "Sizga 192.168.1.10 berishim mumkin"
3. **Request** — Client: "Yarayman, shu IP'ni olaman"
4. **Acknowledge** — Server: "OK, lease sizniki, 24 soat"

Lease tugashidan oldin renew bo'ladi.

```bash
sudo dhclient -v eth0   # manual DHCP
journalctl -u systemd-networkd
```

---

### Q12. URL, URI, URN farqi

**Tag:** Application

**Javob:**

- **URI (Uniform Resource Identifier)** — umumiy nom, **Identifier**.
- **URL (Locator)** — qayerda joylashganligi, ya'ni URI + scheme + host (`https://example.com/page`).
- **URN (Name)** — nomi (`urn:isbn:0451450523`).

URL ⊂ URI, URN ⊂ URI.

---

### Q13. HTTP vs HTTPS

**Tag:** Application, security

**Javob:**

- **HTTP** — plain text, port 80. Hech qanday encryption yo'q — sniff qilish oson.
- **HTTPS** — TLS layer ostida HTTP, port 443. Confidentiality + integrity + authentication (server certificate).

HTTPS handshake (TLS 1.3) qo'shimcha 1 RTT (yoki 0 RTT TFO bilan).

```bash
openssl s_client -connect example.com:443 -tls1_3
```

[TLS deep-dive](../deep-dives/tls-ssl.md)

---

### Q14. localhost va 127.0.0.1

**Tag:** Network

**Javob:**

- **127.0.0.1** — IPv4 loopback address. Aslida butun `127.0.0.0/8` block loopback uchun.
- **localhost** — DNS nomi, `/etc/hosts`'da `127.0.0.1` ga o'rnatilgan (yoki IPv6 `::1`).

Loopback packet'lar real interface'ga chiqmaydi — kernel ichida qaytadi (`lo` interface).

---

### Q15. NAT nima qiladi?

**Tag:** Network, NAT

**Javob:**

NAT (Network Address Translation) — private IP'larni public IP'ga (yoki teskari) o'tkazadi.

- **SNAT** (Source NAT) — outbound. Private → Public (uy router → Internet).
- **DNAT** (Destination NAT) — inbound. Public → Private (port forwarding).
- **PAT** (Port Address Translation) — bitta public IP, bir nechta private host. Port'lar bilan ajratiladi.

```bash
iptables -t nat -L -n -v
nft list table ip nat
```

[NAT deep-dive](../deep-dives/nat-and-firewall.md)

---

### Q16. Hub, Switch, Router farqi

**Tag:** Hardware

**Javob:**

- **Hub (L1):** Hamma portga signal yuboradi (broadcast). Bitta collision domain. Eskirgan.
- **Switch (L2):** MAC bo'yicha to'g'ri portga frame yuboradi. Har port — alohida collision domain. Bitta broadcast domain (VLAN bo'lmasa).
- **Router (L3):** Network'lar orasida IP packet yo'naltiradi. Har interface — alohida network.

---

### Q17. IPv4 va IPv6 farqi

**Tag:** Network

**Javob:**

| Mezon | IPv4 | IPv6 |
|-------|------|------|
| Bit | 32 | 128 |
| Format | `192.168.1.1` | `2001:db8::1` |
| Address space | 4.3 billion | 3.4 × 10³⁸ |
| Header | 20-60 byte | 40 byte fixed |
| Fragmentation | Router + host | Faqat host (PMTUD) |
| NAT | Kerak | Kerak emas (yetarli) |
| ARP | Bor | Yo'q (ICMPv6 NDP) |

2026'da global IPv6 adoption ~50%.

[Network layer](../osi/03-network.md)

---

### Q18. Encapsulation/Decapsulation nima

**Tag:** OSI, fundamentals

**Javob:**

**Encapsulation** — har layer yuqoridan kelgan data'ga **header** (va ba'zan trailer) qo'shadi:
- L7 data → L4 + TCP header → segment
- → L3 + IP header → packet
- → L2 + Ethernet header/trailer → frame
- → L1 bits

**Decapsulation** — qabul tarafda teskari: har layer o'z header'ini o'qib, olib tashlaydi va yuqori layerga yuboradi.

[OSI vs TCP/IP](../00-foundations/osi-vs-tcpip.md)

---

## Mid darajasi (2-5 yil)

### Q19. TCP three-way handshake batafsil

**Tag:** TCP, transport

**Javob:**

```
Client                          Server
  |                               |
  |---- SYN seq=x ---------------->|   (SYN_SENT → SYN_RCVD)
  |<--- SYN+ACK seq=y, ack=x+1 ---|
  |---- ACK ack=y+1 -------------->|   (ESTABLISHED)
```

1. **SYN:** Client server'ga "ulanmoqchiman" — Initial Sequence Number (ISN) `x`.
2. **SYN-ACK:** Server "OK, men ham" — o'z ISN `y`, va `x+1` ni ACK qiladi.
3. **ACK:** Client `y+1` ni ACK qiladi.

ISN tasodifiy tanlanadi (RFC 6528) — security va eski connection'lar bilan adashmaslik uchun.

```bash
tcpdump -tttt -n -i any -c 6 'port 80'
```

[TCP handshake deep-dive](../deep-dives/tcp-handshake.md)

---

### Q20. TLS 1.3 vs TLS 1.2

**Tag:** TLS, security

**Javob:**

- **TLS 1.2 (2008):** 2 RTT handshake. Ko'p cipher suite (yomon va yaxshi aralash). Static RSA key exchange xavfli (forward secrecy yo'q).
- **TLS 1.3 (2018, RFC 8446):** **1 RTT** handshake (0-RTT mumkin). Faqat AEAD cipher (AES-GCM, ChaCha20-Poly1305). Forward secrecy majburiy. Eski cipherlar olib tashlangan (RC4, MD5, SHA-1, CBC).

TLS 1.3'da `Encrypted Extensions` — server hello'dan keyin SNI ham encrypt qilinadi (ECH bilan to'liq private).

[TLS deep-dive](../deep-dives/tls-ssl.md)

---

### Q21. HTTP/2 va HTTP/3 farqi

**Tag:** HTTP, transport

**Javob:**

- **HTTP/2 (2015):** Binary framing, multiplexing (bitta TCP connection ustida ko'p stream), HPACK header compression, server push. Lekin TCP head-of-line blocking — bitta packet yo'qolsa, BARCHA stream to'xtaydi.
- **HTTP/3 (2022):** UDP/QUIC ustida. Stream'lar mustaqil — bittasi yo'qolsa, boshqalari davom etadi. Connection migration (Wi-Fi → 4G — connection saqlanadi). 0-RTT renegotiation.

2026'da HTTP/3 adoption ~22-35%.

[HTTP evolution](../deep-dives/http-evolution.md)

---

### Q22. TIME_WAIT state — nima va nima uchun 2 MSL?

**Tag:** TCP

**Javob:**

TCP connection yopilganda, FIN yuborgan tomonda **TIME_WAIT** state'i 2*MSL (Maximum Segment Lifetime, odatda 30s × 2 = 60s) davom etadi. Sabablari:

1. **Eski packet'lar** network'da hali yurishi mumkin — yangi connection xato'na ulashlirib qo'ymasin.
2. **Final ACK yo'qolsa** — peer FIN'ni qaytaradi va biz yana ACK yuborishimiz kerak.

High RPS server'da TIME_WAIT exhaustion bo'ladi — `net.ipv4.tcp_tw_reuse=1` (client side), `SO_REUSEPORT` server side.

```bash
ss -tan | grep TIME-WAIT | wc -l
sysctl net.ipv4.tcp_tw_reuse
```

---

### Q23. "google.com browserga kirganda nima sodir bo'ladi?"

**Tag:** End-to-end

**Javob:**

1. **DNS:** `gethostbyname("google.com")` — local cache → `/etc/hosts` → `nsswitch.conf` → DNS server (UDP 53). Recursive resolution: root → `.com` → `ns1.google.com`. IP keladi (mas. `142.250.190.46`).
2. **ARP:** Default gateway MAC'ni topish (bitta segment ichida) — agar cache'da yo'q bo'lsa.
3. **TCP handshake:** Client → SYN → Server → SYN-ACK → ACK (port 443).
4. **TLS handshake:** ClientHello → ServerHello + cert → key exchange → Finished. (TLS 1.3'da 1 RTT)
5. **HTTP request:** `GET / HTTP/2` — headers (Host, User-Agent, Cookie...).
6. **HTTP response:** Status, headers, body (HTML).
7. **Browser parse:** HTML → DOM → CSS, JS, image — har biri yangi request (HTTP/2 multiplexing yoki HTTP/3 QUIC).
8. **Render** + scripts execute.

Yopilish: TLS close_notify → TCP FIN handshake → TIME_WAIT.

---

### Q24. NAT turlari (Cone vs Symmetric)

**Tag:** NAT

**Javob:**

NAT mapping behaviour 4 ta:

1. **Full Cone:** Bitta internal `(IP, port)` → bitta external `(IP, port)`. Har kim shu external'ga yuborgan packet ichkariga keladi.
2. **Restricted Cone:** External'dan kelgan packet faqat oldin internal client yuborgan IP'dan qabul qilinadi.
3. **Port Restricted Cone:** IP + port'gacha cheklov.
4. **Symmetric:** Har destination uchun yangi external port ajratiladi. P2P uchun eng qiyin (STUN ishlamaydi, TURN kerak).

[NAT deep-dive](../deep-dives/nat-and-firewall.md)

---

### Q25. ARP poisoning hujumi

**Tag:** Security, ARP

**Javob:**

Hujumchi LAN'da yolg'on ARP reply yuboradi: "Gateway IP — mening MAC'im". Boshqa hostlar cache'ni yangilaydi va traffic hujumchi orqali o'tadi (MitM).

Himoya:
- **Static ARP** muhim host'lar uchun
- **Dynamic ARP Inspection (DAI)** managed switch'da
- **arpwatch** — ARP table o'zgarishini monitor qilish
- **802.1X** authentication

```bash
arpwatch -i eth0
ip neigh show           # current ARP cache
```

---

### Q26. BGP nima va nima uchun Internetning poydevori?

**Tag:** Routing, BGP

**Javob:**

BGP (Border Gateway Protocol, RFC 4271) — **AS (Autonomous System)** lar orasida routing axboroti almashish protokoli. Internet — 100K+ AS'dan iborat. Har AS o'z prefix'larini e'lon qiladi va boshqa AS'lardan o'rganadi.

- **eBGP** — turli AS'lar orasida (TCP 179)
- **iBGP** — bitta AS ichida

Path Vector — har route qanday AS'lardan o'tganini ko'radi (`AS_PATH`). Loop avoidance — o'z AS'ni AS_PATH'da ko'rsa, drop qiladi.

Real incident: 2008 Pakistan-YouTube hijack, 2018 MyEtherWallet, 2024 Cloudflare 1.1.1.1.

[Routing protocols](../deep-dives/routing-protocols.md)

---

### Q27. MTU va Path MTU Discovery

**Tag:** Network

**Javob:**

**MTU** — interface'dan o'ta oladigan maksimal frame hajmi (Ethernet odatda 1500 byte).

**PMTUD** — yo'l bo'ylab eng kichik MTU'ni topish:
1. Source DF (Don't Fragment) bit qo'yib packet yuboradi
2. Yo'lda MTU yetmasa, router ICMP "Fragmentation Needed" qaytaradi
3. Source kichikroq MTU'ga o'tadi

Muammo: ICMP block bo'lsa — **black hole**, packet'lar yo'qoladi sababini bilmasdan.

```bash
ping -M do -s 1472 google.com   # 1472 + 28 ICMP/IP = 1500
tracepath google.com            # PMTUD ko'rsatadi
```

---

### Q28. /24 va /23 farqi (host soni)

**Tag:** Subnetting

**Javob:**

- `/24` — 32-24=8 bit host, 2⁸=256 IP, **254 usable** (network + broadcast minus). Misol: `192.168.1.0/24`.
- `/23` — 9 bit host, 512 IP, **510 usable**. 2 ta /24 ni birlashtirgan.

Formula: usable = 2^(32-prefix) − 2.

```bash
ipcalc 192.168.1.0/23
sipcalc 10.0.0.0/16
```

[Subnetting deep-dive](../deep-dives/subnetting-cidr.md)

---

### Q29. TCP congestion control (Tahoe/Reno/Cubic/BBR)

**Tag:** TCP

**Javob:**

- **Tahoe (1988):** Slow start → congestion avoidance. Yo'qotish → cwnd=1, slow start qaytadan.
- **Reno:** Fast retransmit + fast recovery. 3 dup ACK → cwnd/2, davom etadi (slow start'siz).
- **Cubic (Linux default):** cwnd time'ning kub funksiyasi sifatida o'sadi. High BDP (Bandwidth-Delay Product) link'da yaxshi.
- **BBR (Google, 2016):** Loss'ga emas, **bandwidth + RTT estimation**'ga asoslanadi. Bufferbloat'ni oldini oladi.

```bash
sysctl net.ipv4.tcp_congestion_control
sudo sysctl -w net.ipv4.tcp_congestion_control=bbr
```

---

### Q30. CIDR vs Classful

**Tag:** Subnetting

**Javob:**

- **Classful (legacy):** A (`/8`), B (`/16`), C (`/24`) — fixed mask. IP space israfi (B class — 65K host, lekin ko'pchilik 5K kerak).
- **CIDR (RFC 4632, 1993):** Variable Length — `/13`, `/27`, `/30` — har qanday. VLSM (Variable Length Subnet Masking) bilan bir prefix ichida turli subnet hajmi.

CIDR Internet'ning IPv4 omon qolishini ta'minladi.

---

### Q31. Sticky session vs Round-robin LB

**Tag:** Load balancing

**Javob:**

- **Round-robin:** Har request keyingi backend'ga. Stateless app'lar uchun.
- **Sticky session (session affinity):** Bitta client'ning hamma so'rovi bitta backend'ga. Cookie (`JSESSIONID`) yoki source IP hash bilan.

Stateless arxitektura afzal — backend'da session storage shared (Redis) bo'lsa, sticky kerak emas. Sticky — failover muammoli (backend tushsa, session yo'qoladi).

---

### Q32. CORS qanday ishlaydi

**Tag:** HTTP, security

**Javob:**

Same-Origin Policy — JS faqat o'z domain'idan resource ololadi. CORS — bu cheklovni ochish mexanizmi.

1. **Simple request** (GET, POST, HEAD; standart header): browser `Origin` header qo'shadi. Server `Access-Control-Allow-Origin: *` (yoki specific) qaytarsa — OK.
2. **Preflight (OPTIONS):** Murakkab so'rovdan oldin browser `OPTIONS` yuboradi. Server `Access-Control-Allow-Methods`, `Allow-Headers` bilan javob.

Credentials (cookie): `Access-Control-Allow-Credentials: true` + Origin'da `*` ishlatish mumkin emas.

---

### Q33. WebSocket vs Long polling

**Tag:** Application

**Javob:**

- **Long polling:** Client HTTP so'rov yuboradi, server javobni ushlab turadi (pending) data kelguncha. Keyin yana so'rov. Overhead bor (har so'rovga TLS, header).
- **WebSocket:** HTTP Upgrade orqali full-duplex TCP connection. Bitta connection — ikki tarafga real-time data. Header overhead minimal.

WebSocket — chat, real-time dashboard, gaming. Long polling — eskirib qolgan, lekin firewall'lar WebSocket bloklasa hali ham ishlatiladi.

```
GET /chat HTTP/1.1
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Key: ...
```

---

### Q34. JWT vs Session cookie

**Tag:** Auth

**Javob:**

- **Session cookie:** Server session ID generate qiladi, storage'da (Redis) saqlaydi. Cookie ID'ni client'ga yuboradi. Har request'da server lookup qiladi. Server stateful.
- **JWT:** Server signed token (header.payload.signature) yaratadi, client cookie/localStorage'da saqlaydi. Har request'da signature verify — DB lookup yo'q. Stateless.

JWT trade-off: revocation qiyin (token expire bo'lguncha haqiqiy), katta hajm (cookie 4KB limit'iga yaqin).

---

### Q35. DNS caching va negative cache

**Tag:** DNS

**Javob:**

DNS response'larida **TTL** field bor — qancha vaqt cache'da saqlash mumkin. Resolver TTL davomida o'sha javobni qaytaradi (DNS server'ni qiynamasdan).

**Negative cache (RFC 2308):** NXDOMAIN javobni ham cache qilish — `SOA`'ning `MINIMUM` field TTL'ini belgilaydi. Yo'q domain'larni qayta-qayta so'ramasdan.

```bash
dig example.com         # answer + ttl
dig nonexistent12345.com  # NXDOMAIN + soa ttl
```

[DNS deep-dive](../deep-dives/dns-resolution.md)

---

### Q36. SYN flood hujumi

**Tag:** Security, TCP

**Javob:**

Hujumchi server'ga ko'p SYN yuboradi (spoofed source IP), server SYN-ACK qaytaradi va SYN_RCVD state'ida resource band qiladi. ACK kelmasa, queue to'lib qoladi — yangi connection qabul qilinmaydi.

Himoya:
- **SYN cookies (Bernstein, 1996):** Server SYN_RCVD state saqlamaydi — cookie SEQ ichida encode qilinadi.
- `net.ipv4.tcp_syncookies=1`
- Rate limit (iptables `--limit`)
- SYN proxy (HAProxy, F5)

```bash
sysctl net.ipv4.tcp_syncookies
```

---

## Senior darajasi (5+ yil)

### Q37. TCP TIME_WAIT exhaustion — high RPS optimallashtirish

**Tag:** TCP, performance

**Javob:**

Symptom: `cannot assign requested address`, `EADDRNOTAVAIL`. Har TIME_WAIT 60s davomida (IP, port) tuple'ni band qiladi. Outgoing connection'lar uchun `net.ipv4.ip_local_port_range` (default 32768-60999 = 28K port). 1000 RPS × 60s = 60K ulanish — limit oshib ketadi.

Yechimlar:
1. **`net.ipv4.tcp_tw_reuse=1`** — outbound TIME_WAIT'larni reuse qilish (timestamp tekshiriladi).
2. **`net.ipv4.ip_local_port_range = 1024 65535`** — port range kengaytirish.
3. **Connection pool / keep-alive** — yangi connection'lar yarata olmaslik.
4. **HTTP/2 multiplexing** — bitta connection ustida ko'p request.
5. **`SO_LINGER` 0** — connection RST bilan yopish (TIME_WAIT yo'q, lekin reliability yo'qoladi).

`tcp_tw_recycle` — RHEL 7+da olib tashlangan, NAT bilan muammo qiladi (timestamp).

---

### Q38. 10 Gbps NIC TCP throughput maxim

**Tag:** Performance, kernel

**Javob:**

10 Gbps tomon real throughput uchun:

1. **TCP buffer'larni oshirish:**
   ```
   net.core.rmem_max = 134217728
   net.core.wmem_max = 134217728
   net.ipv4.tcp_rmem = 4096 87380 134217728
   ```
2. **TSO/GSO/GRO** (TCP/Generic Segmentation Offload) — kernel emas, NIC segmentatsiya qilsin: `ethtool -K eth0 tso on gso on gro on`.
3. **RSS (Receive Side Scaling)** — bir nechta queue, har CPU'ga taqsimlash. `ethtool -L eth0 combined N`.
4. **CPU affinity / IRQ pinning** — interrupt'larni NUMA local CPU'ga.
5. **MTU 9000 (jumbo frames)** — header overhead kamayadi.
6. **BBR congestion control:** `sysctl -w net.ipv4.tcp_congestion_control=bbr`.
7. **Window scaling** majburiy: `net.ipv4.tcp_window_scaling=1`.
8. **NUMA-aware** — NIC qaysi NUMA node'da, app shu node'da.

Test: `iperf3 -P 16 -t 60` (16 parallel stream).

---

### Q39. Symmetric NAT orqasidagi 2 host P2P (STUN/TURN/ICE)

**Tag:** NAT, P2P

**Javob:**

Symmetric NAT'da har destination uchun yangi external port — STUN ishlamaydi.

**ICE (Interactive Connectivity Establishment, RFC 8445):**

1. **STUN** — har host o'z external (IP, port)'ni topadi.
2. **Hole punching** — ikkala host bir-biriga packet yuboradi (NAT'ning timer'ini ochadi).
3. Symmetric NAT bilan ikkala host bir-birining yangi mapping'ini bilmaydi → fail.
4. **TURN (Traversal Using Relays around NAT)** — relay server media'ni o'rtadan o'tkazadi. Bandwidth qimmat lekin universal.

WebRTC ICE candidate'lar: host (LAN), srflx (server reflexive — STUN), relay (TURN).

---

### Q40. DNS cache poisoning + DNSSEC

**Tag:** Security, DNS

**Javob:**

**Cache poisoning (Kaminsky, 2008):** Hujumchi resolver'ga yolg'on response yuboradi, resolver cache'iga kirib qoladi. Asl javob keyinroq kelsa, cache'da yolg'on saqlanadi.

DNS UDP, transaction ID 16 bit — brute force mumkin. Source port randomization (RFC 5452) — qo'shimcha 16 bit entropy.

**DNSSEC** (RFC 4033):
- **RRSIG** — har record uchun signature
- **DNSKEY** — public key
- **DS** — delegation signer (parent zone'da)
- Chain of trust: root → TLD → zone

Resolver signature verify qiladi — yolg'on response signature'siz topiladi. Lekin DNSSEC adoption ~10-25%, key rollover qiyin.

[DNS deep-dive](../deep-dives/dns-resolution.md)

---

### Q41. BGP route hijack + RPKI

**Tag:** BGP, security

**Javob:**

**Hijack:** Hujumchi AS o'ziga tegishli bo'lmagan prefix'ni e'lon qiladi. Internet routing shu yo'lga o'tadi (more specific prefix'ga yoki shorter AS_PATH'ga). Misol: 2008 Pakistan AS17557 YouTube prefix'ini hijack qildi → global outage.

**RPKI (Resource Public Key Infrastructure, RFC 6480):**
- IP block egasi — RIR (ARIN, RIPE) — public key signed **ROA (Route Origin Authorization)** chiqaradi: "Bu prefix faqat AS X tomonidan e'lon qilinishi mumkin".
- Router validator (`rpki-client`, `routinator`) ROA'larni yig'adi va BGP route'larni `valid` / `invalid` / `notfound` deb belgilaydi.
- `invalid` route'larni drop qilish.

2026'da ~50% prefix RPKI bilan himoyalangan. Cloudflare, Google, AT&T `invalid` drop qiladi.

---

### Q42. QUIC nima uchun UDP ustida emas, TCP ustida emas?

**Tag:** QUIC, transport

**Javob:**

TCP'da:
- **Head-of-line blocking** — bitta packet yo'qolsa, BARCHA yuqori-layer stream to'xtaydi.
- **Kernel'da implementatsiya** — yangi feature deploy qilish 5-10 yil (OS update kerak).
- **Middlebox ossification** — TCP option'larni proxy/firewall blocklaydi.

QUIC (UDP ustida user-space):
- Stream'lar mustaqil — bittasi yo'qolsa, boshqalari davom etadi.
- TLS 1.3 majburiy integratsiyalashgan (handshake = TLS handshake).
- Connection migration — Wi-Fi → 4G connection saqlanadi (connection ID, IP/port emas).
- Library update (Chrome, Firefox) tez deploy.
- 0-RTT — qaytalangan ulanish data bilan ketadi.

Trade-off: UDP throttling ba'zi network'larda, CPU cost (user-space + crypto).

---

### Q43. Service mesh (Istio/Linkerd) sidecar proxy

**Tag:** Microservices, networking

**Javob:**

Sidecar proxy (Envoy in Istio, linkerd2-proxy in Linkerd) — har Pod yonida ishlaydigan L7 proxy. App TCP traffic'i avval sidecar'dan o'tadi.

Imkoniyatlar:
- **mTLS** — service'lar orasida automatic certificate rotation.
- **Routing** — canary, A/B testing (header'ga qarab traffic split).
- **Circuit breaking, retry, timeout** — resilience pattern.
- **Observability** — traffic metrika, distributed tracing.

Topology: iptables `REDIRECT` — Pod traffic'i sidecar'ga, sidecar destination Pod sidecar'ga.

Cost: latency (~1-2 ms), resource (~50-100 MB RAM × Pod count). eBPF (Cilium) — sidecar-less, kernel'da implement.

---

### Q44. Kubernetes Pod-to-Pod networking

**Tag:** Kubernetes

**Javob:**

K8s networking 4 ta talab:
1. Har Pod o'z IP — `Pod IP`.
2. Pod ↔ Pod NAT'siz (cluster ichida).
3. Pod ↔ Node NAT'siz.
4. Service IP — virtual.

**CNI plugins:**
- **Flannel** — VXLAN encapsulation.
- **Calico** — BGP, har Node BGP speaker. Encapsulation'siz (yoki IPIP).
- **Cilium** — eBPF, kernel'da L3-L7 routing.

**kube-proxy modes:**
- **iptables** (default) — Service IP → Pod IP DNAT, har rule O(N).
- **IPVS** — kernel-level LB, hash table O(1).
- **eBPF (Cilium kube-proxy replacement)** — fastest.

DNS: CoreDNS Pod, har service `<name>.<ns>.svc.cluster.local`.

---

### Q45. Anycast routing va CDN

**Tag:** Routing

**Javob:**

**Anycast** — bir IP address dunyoning bir nechta joyida e'lon qilingan. BGP eng yaqin (AS_PATH, MED, IGP cost) instance'ga yo'naltiradi.

Misol: `1.1.1.1` (Cloudflare) — 320+ shahar, har biri shu IP'ni e'lon qiladi. Foydalanuvchi eng yaqin PoP'ga ulanadi.

CDN'da:
- Statik content cache — eng yaqin edge.
- DDoS himoya — hujum global'ga taqsimlanadi (Cloudflare 2024 19.6 Tbps absorb).
- DNS — root server'lar anycast (13 logical, 1000+ instance).

Trade-off: TCP connection mid-session'da PoP almashishi mumkin (BGP routing o'zgarsa) → connection break. Connection-aware routing kerak.

---

### Q46. L4 vs L7 load balancer

**Tag:** Load balancing

**Javob:**

- **L4 (TCP/UDP):** IP+port asosida, payload o'qimaydi. Tez (kernel-space, eBPF/IPVS), connection-level (sticky). Misol: HAProxy TCP mode, AWS NLB, F5 LTM.
- **L7 (HTTP):** Header, URL, cookie asosida. Smart routing (path → service A, path → service B). Slower (decrypt TLS, parse HTTP). Misol: nginx, HAProxy HTTP, Envoy, AWS ALB.

Trade-off:
- L4 — TLS passthrough mumkin (backend decrypts), faster, host'larni TCP/UDP arbitrary protokol uchun.
- L7 — Smart routing, observability, WAF, rate-limit per-route.

Real arxitektura: L4 LB (DDoS layer) → L7 LB (routing) → backend.

---

### Q47. TCP BBR vs Cubic

**Tag:** TCP, congestion

**Javob:**

**Cubic** (loss-based):
- Yo'qotish → cwnd kamaytirish.
- Buffer to'lguncha cwnd o'sadi.
- **Bufferbloat** — router buffer to'lganda RTT katta o'sadi, latency yomon.

**BBR** (model-based, Google 2016):
- **Bandwidth × RTT** estimate qiladi (BDP — Bandwidth-Delay Product).
- cwnd = BDP — buffer to'lmaydi.
- Loss bo'lsa ham (random loss) cwnd kamaytirmaydi.
- Throughput Cubic'dan 2-25× ko'p (lossy yoki long-RTT link'da).

Trade-off: BBR Cubic flow'lar bilan bo'lsa — adolatsiz (BBR ko'proq oladi). v2/v3 (CCS) o'rtacha.

```bash
sysctl -w net.ipv4.tcp_congestion_control=bbr
sysctl -w net.core.default_qdisc=fq   # fair queueing — BBR uchun majburiy
```

---

### Q48. ECN (Explicit Congestion Notification)

**Tag:** Network

**Javob:**

ECN (RFC 3168) — router congestion'ni packet drop o'rniga **bit'da signal** qiladi:
- IP header'da 2 bit `ECN` field
- TCP header'da `ECE` (ECN-Echo), `CWR` (Congestion Window Reduced) flag

Jarayon:
1. SYN'da client: "ECN-capable" bilan negotiate.
2. Router buffer to'la bo'lganda packet drop o'rniga **CE (Congestion Experienced)** bit qo'yadi.
3. Receiver `ECE` flag bilan ACK'da xabar beradi.
4. Sender cwnd kamaytiradi (loss bo'lganday).

Afzalligi: zero packet loss, latency past. Lekin middlebox'lar ECN bit'larni clear qiladi yoki connection'ni RST qiladi (ossification).

L4S (Low Latency Low Loss Scalable, RFC 9330) — modern variant, datacenter va Internet'da.

---

### Q49. eBPF networking

**Tag:** Linux, performance

**Javob:**

**eBPF** — kernel'da safe sandboxed programs ishga tushirish. Networking'da:

- **XDP (eXpress Data Path):** NIC driver'da packet — eng tez (DPDK level). DDoS mitigation, LB.
- **TC (Traffic Control)** — qdisc'da, kontainer ingress/egress.
- **Socket filtering** — `setsockopt(SO_ATTACH_BPF)`.
- **kprobes/uprobes** — observability (bpftrace, BCC).

Real ishlatilishi:
- **Cilium** — K8s CNI, kube-proxy replacement.
- **Katran (Facebook)** — XDP-based L4 LB, per-server 10M PPS.
- **Cloudflare DDoS** — XDP'da.
- **Calico** — eBPF dataplane.

Trade-off: kernel 4.18+, debugging qiyin (verifier reject), maps lifecycle.

---

### Q50. mTLS qanday ishlaydi

**Tag:** Security, TLS

**Javob:**

Mutual TLS — server **va** client ham certificate bilan authenticate.

TLS 1.3 handshake'da:
1. ClientHello → ServerHello + cert + **CertificateRequest**
2. Server cert verify → Client cert + signature → server verify
3. Finished

Server `ClientCAs` bilan trust store belgilab beradi. Client cert chain `Issuer` shu CA'ga borishi kerak.

Use cases:
- Service-to-service auth (Istio, Linkerd) — sertifikatlar avtomatik rotation.
- Bank-to-bank API (PSD2)
- IoT device auth.

Go misol:
```go
caCert, _ := os.ReadFile("ca.crt")
caPool := x509.NewCertPool()
caPool.AppendCertsFromPEM(caCert)
tls.Config{
    ClientAuth: tls.RequireAndVerifyClientCert,
    ClientCAs:  caPool,
}
```

[TLS deep-dive](../deep-dives/tls-ssl.md)

---

### Q51. DDoS himoyasi

**Tag:** Security

**Javob:**

DDoS turlari:
1. **Volumetric** (UDP flood, amplification) — bandwidth tugatish. DNS amp, NTP amp, memcached (50000× amplification).
2. **Protocol** — SYN flood, Smurf, ping of death — server resource.
3. **Application** — Slowloris (sekin HTTP), HTTP flood, Layer 7 attack — app server.

Himoya:
- **Anycast network (Cloudflare, AWS Shield)** — global'ga taqsim.
- **Rate limit** — per-IP, per-ASN.
- **SYN cookies** — protocol attack.
- **WAF** — application layer.
- **BGP RTBH (Remotely Triggered Black Hole)** — hujum prefix'ini ISP'da drop.
- **eBPF/XDP** — line-rate filter.

Recent: 2024 Cloudflare 4.2 Tbps record, Github 1.35 Tbps memcached amp.

---

### Q52. Zero-trust networking

**Tag:** Security, architecture

**Javob:**

Eski model (perimeter-based): "Network ichida bo'lsang, ishonchli". Hujumchi periphery yorilsa — butun ichkari ochiq.

**Zero-trust:** Hech kimga ishonma, har request'da verify.

Komponentlar:
- **Identity-based access** — IP emas, identity (user, service).
- **mTLS** har komponentlar orasida.
- **Just-in-time access** — short-lived credential.
- **Micro-segmentation** — har workload alohida policy (Cilium, Calico NetworkPolicy).
- **Continuous verification** — har request'da policy check.

Implementations: BeyondCorp (Google), Tailscale, Cloudflare Zero Trust, AWS Verified Access.

---

### Q53. WireGuard vs IPsec

**Tag:** VPN

**Javob:**

| Mezon | IPsec | WireGuard |
|-------|-------|-----------|
| Year | 1995 | 2018 |
| Code | 100K+ qator | ~4000 qator |
| Crypto | AES, SHA, 3DES, ko'p | ChaCha20-Poly1305, Curve25519 (faqat) |
| Key exchange | IKEv2 (murakkab) | Noise protocol |
| Setup | Qiyin | Oson (`wg0.conf`) |
| Performance | Yaxshi | A'lo (kernel module) |
| NAT traversal | NAT-T | UDP, oson |

WireGuard kernel'da (Linux 5.6+). Tailscale — WireGuard ustida managed mesh VPN.

```bash
wg genkey | tee privatekey | wg pubkey > publickey
ip link add wg0 type wireguard
wg setconf wg0 wg0.conf
ip addr add 10.0.0.1/24 dev wg0
```

---

## Yakunda

Bu fayllar **interview tayyorgarligi** uchun. Real ishda har savol bo'yicha layer fayli yoki deep-dive faylini chuqurroq o'qing:

- [OSI layers](../osi/README.md)
- [TCP/IP layers](../tcp-ip/README.md)
- [Deep-dives](../deep-dives/)
- [Glossary](../00-foundations/glossary.md)
- [Troubleshooting cases](./troubleshooting-cases.md)

**Tip:** Interview'da javob berishdan oldin "Eslab qolaman, OSI L4 — transport layer..." deb sekin tartiblang. Hech vaqt taxminga asoslanib javob bermang — bilmasangiz "to'g'ri eslamayman, lekin shu kontekstda mantiqan..." degan yaxshi.
