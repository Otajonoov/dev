# Network Interview Savol-Javoblar

Backend / DevOps / SRE intervyusiga tayyorgarlik uchun 90+ savol, mavzu bo'yicha guruhlangan. Har savolning javobi `<details>` ichida yashiringan — avval o'zing javob berishga urin, keyin och.

## Mavzular

1. Asoslar (OSI, encapsulation, qurilmalar)
2. L2 / Network Access (ARP, DHCP, VLAN)
3. IP addressing va subnetting
4. Routing (static, OSPF, BGP)
5. Transport (TCP, UDP)
6. Application (HTTP, DNS, TLS)
7. Security va advanced
8. Backend end-to-end ("URL yozganda nima bo'ladi?")

> **Intervyu tip:** Javob berishdan oldin "Bu L4 — transport layer..." deb sekin tartibla. Bilmasang, taxmin qilma: "aniq eslamayman, lekin bu kontekstda mantiqan..." degan yaxshiroq.

---

## 1. Asoslar

### Q1. OSI model 7 ta layer ni sanab bering. Har biri nima qiladi?

<details>
<summary>Javobni ko'rish</summary>

1. **Physical (L1)** — bitlarni signalga aylantirish (kabel, fiber, radio).
2. **Data Link (L2)** — bitta link ichida frame, MAC address (Ethernet, Wi-Fi, ARP).
3. **Network (L3)** — hostlar orasida packet routing, IP address (IPv4, IPv6, ICMP).
4. **Transport (L4)** — process-process aloqa, port (TCP, UDP).
5. **Session (L5)** — session/dialog boshqarish (RPC, NetBIOS).
6. **Presentation (L6)** — encoding, encryption, compression (TLS, JPEG).
7. **Application (L7)** — foydalanuvchi protokollari (HTTP, DNS, SSH).

Eslab qolish: **A**ll **P**eople **S**eem **T**o **N**eed **D**ata **P**rocessing (L7 -> L1).
</details>

### Q2. Encapsulation / Decapsulation nima?

<details>
<summary>Javobni ko'rish</summary>

**Encapsulation** — har layer yuqoridan kelgan data ga header (va ba'zan trailer) qo'shadi:

- L7 data -> +TCP header -> **segment**
- -> +IP header -> **packet**
- -> +Ethernet header/trailer -> **frame**
- -> L1 bits

**Decapsulation** — qabul tarafda teskari: har layer o'z header ini o'qib olib tashlaydi va yuqoriga uzatadi.
</details>

### Q3. MAC address va IP address farqi?

<details>
<summary>Javobni ko'rish</summary>

- **MAC (L2):** 48 bit, `aa:bb:cc:dd:ee:ff`, NIC ga vendor beradi (OUI), faqat **bitta segment ichida** ishlaydi.
- **IP (L3):** 32 bit (IPv4) yoki 128 bit (IPv6), admin/DHCP beradi, **internet bo'ylab** routerlar orqali yo'naltiriladi.

Har hop da MAC o'zgaradi, IP odatda o'zgarmaydi (NAT bo'lmasa). ARP IP dan MAC ni topadi.
</details>

### Q4. Hub, Switch, Router farqi?

<details>
<summary>Javobni ko'rish</summary>

- **Hub (L1):** hamma portga signal (broadcast), bitta collision domain — eskirgan.
- **Switch (L2):** MAC bo'yicha to'g'ri portga frame, har port alohida collision domain, bitta broadcast domain (VLAN yo'q bo'lsa).
- **Router (L3):** network'lar orasida IP packet, har interfeys alohida network.

Modern qurilma — **L3 switch** (switch tezligi + routing).
</details>

### Q5. ping va traceroute farqi?

<details>
<summary>Javobni ko'rish</summary>

- **ping** — ICMP Echo Request/Reply, RTT va packet loss ko'rsatadi.
- **traceroute** — TTL=1,2,3... bilan paket yuboradi. Har router TTL=0 da ICMP Time Exceeded qaytaradi, shu bilan butun yo'l ko'rinadi.

`mtr` — ping + traceroute gibrid, davomiy monitoring.
</details>

### Q6. localhost va 127.0.0.1 farqi?

<details>
<summary>Javobni ko'rish</summary>

- **127.0.0.1** — IPv4 loopback address (butun `127.0.0.0/8` block loopback).
- **localhost** — DNS/`hosts` nomi, `127.0.0.1` (yoki IPv6 `::1`) ga o'rnatilgan.

Loopback paket fizik interfeysga chiqmaydi — kernel ichida `lo` orqali qaytadi.
</details>

---

## 2. L2 / Network Access

### Q7. ARP nima qiladi?

<details>
<summary>Javobni ko'rish</summary>

ARP (Address Resolution Protocol) — IP dan MAC ni topadi (bitta segment ichida):

1. Host broadcast: "192.168.1.1 kimning IP si? MAC ni ayt".
2. Egasi unicast: "Menman, MAC: aa:bb:cc...".
3. Natija cache ga saqlanadi (`ip neigh`).

IPv6 da ARP yo'q — o'rniga ICMPv6 NDP (Neighbor Discovery).
</details>

### Q8. ARP poisoning (spoofing) hujumi va himoya?

<details>
<summary>Javobni ko'rish</summary>

Hujumchi LAN da yolg'on ARP reply yuboradi: "Gateway IP — mening MAC im". Hostlar cache ni yangilaydi, trafik hujumchi orqali o'tadi (MitM).

Himoya: static ARP muhim hostlar uchun; **Dynamic ARP Inspection (DAI)** managed switch da; `arpwatch`; 802.1X authentication; DHCP snooping.
</details>

### Q9. DHCP qanday ishlaydi (DORA)?

<details>
<summary>Javobni ko'rish</summary>

**DORA** — 4 bosqich (UDP 67/68):

1. **Discover** — client broadcast: "DHCP server bormi?"
2. **Offer** — server: "192.168.1.10 berishim mumkin".
3. **Request** — client: "Shu IP ni olaman".
4. **Acknowledge** — server: "OK, lease 24 soat".

DISCOVER broadcast, chunki client hali o'z IP sini ham, server manzilini ham bilmaydi. Cross-VLAN uchun `ip helper-address` (DHCP relay) kerak.
</details>

### Q10. VLAN nima va nima uchun kerak?

<details>
<summary>Javobni ko'rish</summary>

VLAN — bitta jismoniy switch ustidagi mantiqiy broadcast domain. Bir port VLAN 10, boshqasi VLAN 20 bo'lsa, ular L2 da bir-birini ko'rmaydi.

Foyda: segmentatsiya (bo'limlarni ajratish), xavfsizlik, broadcast domain kichraytirish. VLAN lar orasida gaplashish uchun L3 (inter-VLAN routing / router-on-a-stick / SVI) kerak.

802.1Q tag — frame ga 4 baytlik VLAN ID qo'shadi (trunk link da).
</details>

---

## 3. IP addressing va subnetting

### Q11. Subnet mask nima uchun kerak?

<details>
<summary>Javobni ko'rish</summary>

IP addressning o'zi qaysi qismi network, qaysi qismi host ekanini aytmaydi. Mask IP ni **network + host** ga ajratadi. Maskdagi `1` bitlar network, `0` bitlar host.

`192.168.1.10/24`: birinchi 24 bit network, oxirgi 8 bit host. Router "bu IP qaysi network ga tegishli?" degan savolga mask orqali javob topadi.
</details>

### Q12. Nima uchun usable host `-2` qilinadi? /31 va /32 chi?

<details>
<summary>Javobni ko'rish</summary>

Har subnetda 2 ta address hostga berilmaydi: **network address** (subnet nomi) va **broadcast address** (hammaga yuborish). Shuning uchun usable = `2^host_bits - 2`.

- **/31** — point-to-point link (2 router) da broadcast kerak emas, shuning uchun 2 ta address ham ishlatiladi (RFC 3021).
- **/32** — bitta aniq host (loopback, host route, firewall qoidasi).
</details>

### Q13. Private vs Public IP. NAT muammoni qanday yengillashtirdi?

<details>
<summary>Javobni ko'rish</summary>

**Private (RFC 1918):** internetda route qilinmaydi — `10.0.0.0/8`, `172.16.0.0/12`, `192.168.0.0/16`.
**Public:** internetda unique, RIR ajratadi.

IPv4 address soni cheklangan. NAT bitta public IP orqasida ko'p private host ishlashiga imkon berdi — router portlar orqali qaysi ichki host qaysi ulanish ochganini eslab turadi (PAT).
</details>

### Q14. Classful vs Classless (CIDR)?

<details>
<summary>Javobni ko'rish</summary>

- **Classful:** birinchi octet ga qarab fixed mask — A `/8`, B `/16`, C `/24`. IP space isrofi (B class 65K host).
- **CIDR (RFC 4632):** variable length — `/13`, `/27`, `/30`. VLSM bilan bir prefix ichida turli subnet hajmi.

"192.168.x.x doim /24" — **xato**, chunki 192.168 faqat private range ni bildiradi, mask ni emas. Subnet chegarasini doim mask/CIDR belgilaydi.
</details>

### Q15. /24 va /23 farqi (host soni)?

<details>
<summary>Javobni ko'rish</summary>

- `/24` — 8 host bit, 256 IP, **254 usable**.
- `/23` — 9 host bit, 512 IP, **510 usable** (2 ta /24 birlashtirilgan).

Formula: usable = `2^(32-prefix) - 2`.
</details>

### Q16. IPv4 vs IPv6?

<details>
<summary>Javobni ko'rish</summary>

| Mezon | IPv4 | IPv6 |
|-------|------|------|
| Bit | 32 | 128 |
| Format | `192.168.1.1` | `2001:db8::1` |
| Header | 20-60 bayt | 40 bayt fixed |
| Fragmentation | router + host | faqat host (PMTUD) |
| NAT | kerak | kerak emas |
| ARP | bor | yo'q (ICMPv6 NDP) |

SLAAC — IPv6 da host o'ziga avtomatik address beradi (Router Advertisement dan prefix oladi).
</details>

### Q17. Longest Prefix Match nima? 0.0.0.0/0 chi?

<details>
<summary>Javobni ko'rish</summary>

Router destination IP uchun routing table dan **eng aniq (eng uzun prefix)** route ni tanlaydi. Masalan `10.10.20.55` uchun `/24`, `/16`, `/8`, `/0` mavjud bo'lsa — `/24` tanlanadi.

`0.0.0.0/0` — **default route** (network bit yo'q, hamma destination mos). Faqat aniqroq route bo'lmaganda ishlatiladi.
</details>

### Q18. 169.254.x.x ko'rsam nima o'ylashim kerak?

<details>
<summary>Javobni ko'rish</summary>

Ko'pincha **DHCP ishlamaganini**. Host DHCP dan address ololmasa, o'ziga link-local/APIPA address beradi (`169.254.0.0/16`). Lokal segmentda cheklangan aloqa bo'lishi mumkin, lekin gateway va internet ishlamaydi.

Tekshirish: DHCP relay (`ip helper`), switch port VLAN, DHCP server scope.
</details>

### Q19. TTL nima uchun kerak? Fragmentation nega muammo?

<details>
<summary>Javobni ko'rish</summary>

**TTL** — paket routing loop da cheksiz aylanmasligi uchun. Har router TTL ni 1 kamaytiradi, 0 bo'lsa drop qiladi. Traceroute shu mexanizmdan foydalanadi.

**Fragmentation** muammosi: performance pasayadi, firewall tekshiruvi murakkablashadi, ICMP block bo'lsa PMTUD buziladi (black hole). Yechim — paket hajmini path MTU ga moslash (MSS clamping).
</details>

---

## 4. Routing

### Q20. Static va dynamic routing farqi? Administrative Distance?

<details>
<summary>Javobni ko'rish</summary>

- **Static** — admin qo'lda yozadi. Oddiy, bashoratli, lekin qo'lda scale bo'lmaydi.
- **Dynamic (OSPF, BGP)** — routerlar bir-biridan o'rganadi, o'zgarishga avtomatik moslashadi (convergence).

**Administrative Distance (AD)** — route manbaining ishonchliligi (kichik = ishonchli): Connected 0, Static 1, eBGP 20, OSPF 110, RIP 120. Bir destinationga ikki manbadan route kelsa, AD kichigi tanlanadi.
</details>

### Q21. OSPF neighbor qanday FULL bo'ladi?

<details>
<summary>Javobni ko'rish</summary>

OSPF routerlar Hello paket bilan qo'shni topadi, keyin state machine: Down -> Init -> 2-Way -> ExStart -> Exchange -> Loading -> **Full**. Adjacency uchun area ID, subnet, hello/dead timer, authentication mos kelishi shart.

DR/BDR — broadcast segmentda saylanadi (LSA flooding ni kamaytirish uchun). Router ID — eng katta loopback IP yoki qo'lda `router-id`.
</details>

### Q22. BGP nima va nima uchun internetning poydevori?

<details>
<summary>Javobni ko'rish</summary>

BGP (RFC 4271) — **AS (Autonomous System)** lar orasida routing axboroti almashadi (TCP 179). Internet 100K+ AS dan iborat. Har AS o'z prefixlarini e'lon qiladi.

- **eBGP** — turli AS orasida; **iBGP** — bitta AS ichida.
- **Path Vector** — `AS_PATH` orqali loop avoidance (o'z AS ni ko'rsa drop).

Real incident: 2008 Pakistan-YouTube hijack, 2018 MyEtherWallet.
</details>

### Q23. BGP route hijack va RPKI?

<details>
<summary>Javobni ko'rish</summary>

**Hijack** — hujumchi AS o'ziga tegishli bo'lmagan prefix ni e'lon qiladi (more-specific yoki qisqaroq AS_PATH). Trafik shu yo'lga o'tadi.

**RPKI (RFC 6480):** IP block egasi signed **ROA (Route Origin Authorization)** chiqaradi — "bu prefix faqat AS X tomonidan e'lon qilinsin". Router validator route larni valid/invalid/notfound deb belgilaydi, `invalid` ni drop qiladi.
</details>

### Q24. Anycast routing va CDN?

<details>
<summary>Javobni ko'rish</summary>

**Anycast** — bir IP dunyoning ko'p joyida e'lon qilinadi. BGP eng yaqin instance ga yo'naltiradi. Masalan `1.1.1.1` (Cloudflare) — 320+ shahar.

CDN da: statik content eng yaqin edge dan; DDoS himoya global taqsimlanadi; DNS root serverlar anycast (13 logical, 1000+ instance).

Trade-off: BGP routing o'zgarsa TCP session mid-session PoP almashib break bo'lishi mumkin.
</details>

### Q25. MTU va Path MTU Discovery?

<details>
<summary>Javobni ko'rish</summary>

**MTU** — interfeysdan o'ta oladigan maksimal frame (Ethernet 1500 bayt).

**PMTUD:** source DF (Don't Fragment) bilan paket yuboradi; yo'lda MTU yetmasa router ICMP "Fragmentation Needed" qaytaradi; source kichikroq MTU ga o'tadi.

Muammo: ICMP block bo'lsa — **black hole** (paketlar sababsiz yo'qoladi). Test: `ping -M do -s 1472 host`, `tracepath host`.
</details>

---

## 5. Transport (TCP, UDP)

### Q26. Nima uchun transport layer kerak, network layer host'larni ulasa?

<details>
<summary>Javobni ko'rish</summary>

Router faqat network layer da ishlaydi — u transport segmenti ichiga qaramaydi. Network layer **host-to-host** yetkazadi, lekin bir hostda ko'p process (browser, ssh, DB) bor.

Transport layer **process-to-process** aloqani ta'minlaydi — **port** raqami orqali multiplexing/demultiplexing qiladi. Qaysi kelgan segment qaysi process ga borishini port belgilaydi.
</details>

### Q27. TCP va UDP farqi?

<details>
<summary>Javobni ko'rish</summary>

| Mezon | TCP | UDP |
|-------|-----|-----|
| Connection | 3-way handshake | yo'q |
| Reliability | ishonchli (ACK, retransmit) | ishonchsiz |
| Order | tartibli (sequence number) | tartibsiz |
| Speed | sekinroq (overhead) | tez |
| Header | 20-60 bayt | 8 bayt |
| Use case | HTTP, SSH, email | DNS, video, gaming, VoIP |

UDP — tez, lekin yo'qotish/tartibni app o'zi hal qiladi.
</details>

### Q28. TCP three-way handshake batafsil. Nima uchun handshake kerak?

<details>
<summary>Javobni ko'rish</summary>

```
Client                          Server
  |---- SYN seq=x --------------->|  (SYN_SENT -> SYN_RCVD)
  |<--- SYN+ACK seq=y, ack=x+1 --|
  |---- ACK ack=y+1 ------------->|  (ESTABLISHED)
```

1. **SYN:** client ISN `x` yuboradi.
2. **SYN-ACK:** server ISN `y` va `x+1` ni ACK qiladi.
3. **ACK:** client `y+1` ni ACK qiladi.

**Nima uchun:** ikkala tomon boshlang'ich sequence number ga kelishib olishi kerak (paketlarni ajratish/tartiblash uchun). ISN tasodifiy (RFC 6528) — hujum va eski ulanish bilan adashmaslik uchun.
</details>

### Q29. TCP half-open va SYN flood hujumi?

<details>
<summary>Javobni ko'rish</summary>

**Half-open:** client SYN yuboradi, server SYN-ACK qaytaradi va SYN_RCVD holatida ACK kutadi. ACK kelmasa, server resource ushlab turadi (memory).

**SYN flood:** hujumchi spoofed source IP bilan ko'p SYN yuboradi, server SYN_RCVD queue si to'lib qoladi -> legitim ulanishlar rad etiladi.

Himoya: **SYN cookies** (server state saqlamaydi, cookie SEQ ichida encode qilinadi), `tcp_syncookies=1`, rate limit, SYN proxy.
</details>

### Q30. TCP slow start nima?

<details>
<summary>Javobni ko'rish</summary>

Handshake dan keyin TCP kichik congestion window (cwnd) bilan boshlaydi va har ACK da eksponensial oshiradi, toki packet loss yoki threshold ga yetguncha. Bu tarmoqni birdan katta yuk bilan qulatib yubormaslik uchun.

Trade-off: handshake dan darhol katta POST yuborsang, slow start tufayli boshida sekin. Oddiy GET larda sezilmaydi. Yechim — keep-alive connection ushlab turish.
</details>

### Q31. TIME_WAIT state nima va nega 2*MSL?

<details>
<summary>Javobni ko'rish</summary>

Connection yopilganda FIN yuborgan tomonda TIME_WAIT 2*MSL (odatda 60s) davom etadi. Sabablar:

1. Network da hali yurgan eski paketlar yangi ulanishga adashib kirmasin.
2. Final ACK yo'qolsa, peer FIN ni qaytaradi va biz yana ACK yuboramiz.

High RPS server da TIME_WAIT exhaustion bo'ladi -> `net.ipv4.tcp_tw_reuse=1` (client side), connection pooling, port range kengaytirish. `tcp_tw_recycle` kernel 4.12+ da olib tashlangan (NAT bilan buziladi).
</details>

### Q32. TCP congestion control — Cubic vs BBR?

<details>
<summary>Javobni ko'rish</summary>

**Cubic (Linux default, loss-based):** loss bo'lsa cwnd kamaytiradi, buffer to'lguncha o'sadi. Bufferbloat — router buffer to'lganda RTT juda o'sadi.

**BBR (Google 2016, model-based):** bandwidth x RTT (BDP) estimate qiladi, cwnd = BDP (buffer to'lmaydi). Random loss da cwnd kamaytirmaydi. Lossy/long-RTT link da throughput 2-25x ko'p.

```bash
sysctl -w net.ipv4.tcp_congestion_control=bbr
sysctl -w net.core.default_qdisc=fq
```
</details>

### Q33. ECN (Explicit Congestion Notification)?

<details>
<summary>Javobni ko'rish</summary>

ECN (RFC 3168) — router congestion ni packet drop o'rniga **bit da signal** qiladi. IP header da 2 bit ECN, TCP da ECE/CWR flag.

Router buffer to'la bo'lganda paket drop qilmasdan CE (Congestion Experienced) bit qo'yadi -> receiver ECE bilan xabar beradi -> sender cwnd kamaytiradi. Afzallik: zero packet loss, past latency. Kamchilik: ba'zi middlebox ECN bit ni tozalaydi (ossification).
</details>

---

## 6. Application (HTTP, DNS, TLS)

### Q34. DNS qanday ishlaydi? Caching va negative cache?

<details>
<summary>Javobni ko'rish</summary>

Domain -> IP. Bosqichlar: browser cache -> OS resolver / `hosts` -> configured DNS (UDP 53) -> recursive qidiruv: root `.` -> TLD `.com` -> authoritative.

**Caching:** har javobda **TTL** — qancha vaqt saqlash. Resolver TTL davomida shu javobni qaytaradi.
**Negative cache (RFC 2308):** NXDOMAIN ham cache qilinadi (SOA MINIMUM TTL) — yo'q domain ni qayta-qayta so'ramaslik uchun.
</details>

### Q35. HTTP status code 5 ta guruh?

<details>
<summary>Javobni ko'rish</summary>

- **1xx** Informational — 100 Continue, 101 Switching Protocols.
- **2xx** Success — 200 OK, 201 Created, 204 No Content.
- **3xx** Redirect — 301 Moved, 302 Found, 304 Not Modified.
- **4xx** Client error — 400, 401, 403, 404, 429 Too Many Requests.
- **5xx** Server error — 500, 502 Bad Gateway, 503 Unavailable, 504 Gateway Timeout.

502 vs 504: 502 — upstream noto'g'ri javob/uzildi; 504 — upstream vaqtida javob bermadi.
</details>

### Q36. HTTP vs HTTPS. TLS 1.3 vs TLS 1.2?

<details>
<summary>Javobni ko'rish</summary>

**HTTP** — plain text, port 80, sniff qilish oson. **HTTPS** — TLS ostida, port 443, confidentiality + integrity + authentication.

- **TLS 1.2 (2008):** 2 RTT handshake, ko'p cipher (yomon+yaxshi), static RSA (forward secrecy yo'q).
- **TLS 1.3 (2018):** **1 RTT** (0-RTT mumkin), faqat AEAD cipher, forward secrecy majburiy, eski cipherlar olib tashlangan.
</details>

### Q37. HTTP/2 va HTTP/3 farqi?

<details>
<summary>Javobni ko'rish</summary>

- **HTTP/2 (2015):** binary framing, multiplexing (bitta TCP ustida ko'p stream), HPACK compression. Muammo: **TCP head-of-line blocking** — bitta paket yo'qolsa BARCHA stream to'xtaydi.
- **HTTP/3 (2022):** UDP/QUIC ustida. Streamlar mustaqil — biri yo'qolsa boshqasi davom etadi. Connection migration (Wi-Fi -> 4G saqlanadi). 0-RTT.
</details>

### Q38. QUIC nima uchun TCP emas, UDP ustida?

<details>
<summary>Javobni ko'rish</summary>

TCP muammolari: head-of-line blocking; kernel da implement (yangi feature 5-10 yil); middlebox ossification.

QUIC (UDP ustida, user-space): streamlar mustaqil; TLS 1.3 integratsiyalashgan; connection migration; library update tez deploy; 0-RTT.

Trade-off: ba'zi network da UDP throttling, CPU cost (user-space crypto).
</details>

### Q39. WebSocket vs Long polling?

<details>
<summary>Javobni ko'rish</summary>

- **Long polling:** client HTTP so'rov yuboradi, server data kelguncha ushlab turadi, keyin yana so'rov. Overhead bor.
- **WebSocket:** HTTP Upgrade orqali full-duplex TCP. Bitta ulanish — ikki tarafga real-time data, minimal overhead.

WebSocket — chat, real-time dashboard, gaming.
</details>

### Q40. CORS qanday ishlaydi?

<details>
<summary>Javobni ko'rish</summary>

Same-Origin Policy — JS faqat o'z origin idan resource oladi. CORS bu cheklovni ochadi.

1. **Simple request:** browser `Origin` header qo'shadi, server `Access-Control-Allow-Origin` bilan javob.
2. **Preflight (OPTIONS):** murakkab so'rovdan oldin browser OPTIONS yuboradi, server `Allow-Methods/Headers` bilan javob.

Credentials bilan `Allow-Origin: *` ishlatib bo'lmaydi — aniq origin kerak.
</details>

### Q41. JWT vs Session cookie?

<details>
<summary>Javobni ko'rish</summary>

- **Session cookie:** server session ID generate qiladi, storage (Redis) da saqlaydi, har requestda lookup. Server **stateful**.
- **JWT:** signed token (header.payload.signature), client saqlaydi, har requestda signature verify — DB lookup yo'q. **Stateless**.

JWT trade-off: revocation qiyin (expire gacha valid), katta hajm.
</details>

### Q42. URL, URI, URN farqi?

<details>
<summary>Javobni ko'rish</summary>

- **URI** — umumiy identifier.
- **URL** — Locator, qayerda joylashgan (`https://example.com/page`).
- **URN** — Name (`urn:isbn:0451450523`).

URL va URN — URI ning kichik to'plamlari.
</details>

---

## 7. Security va advanced

### Q43. mTLS qanday ishlaydi?

<details>
<summary>Javobni ko'rish</summary>

Mutual TLS — server **va** client ham certificate bilan authenticate. TLS handshake da server `CertificateRequest` yuboradi, client o'z cert + signature ni beradi, server `ClientCAs` trust store bo'yicha tekshiradi.

Use case: service-to-service auth (Istio, Linkerd — auto rotation), bank API (PSD2), IoT device auth.
</details>

### Q44. DDoS turlari va himoya?

<details>
<summary>Javobni ko'rish</summary>

Turlari: **Volumetric** (UDP flood, DNS/NTP/memcached amplification — bandwidth); **Protocol** (SYN flood, Smurf — resource); **Application** (Slowloris, HTTP flood — app server).

Himoya: anycast network (Cloudflare, AWS Shield); rate limit (per-IP, per-ASN); SYN cookies; WAF; BGP RTBH blackhole; eBPF/XDP line-rate filter.
</details>

### Q45. Zero-trust networking?

<details>
<summary>Javobni ko'rish</summary>

Eski (perimeter): "network ichida bo'lsang ishonchli" — periphery yorilsa hammasi ochiq.

**Zero-trust:** hech kimga ishonma, har requestda verify. Komponentlar: identity-based access (IP emas); mTLS; just-in-time credential; micro-segmentation (NetworkPolicy); continuous verification.

Implementatsiya: BeyondCorp (Google), Tailscale, Cloudflare Zero Trust.
</details>

### Q46. WireGuard vs IPsec?

<details>
<summary>Javobni ko'rish</summary>

| Mezon | IPsec | WireGuard |
|-------|-------|-----------|
| Kod | 100K+ qator | ~4000 qator |
| Crypto | ko'p (AES, 3DES...) | ChaCha20-Poly1305 (faqat) |
| Key exchange | IKEv2 (murakkab) | Noise protocol |
| Setup | qiyin | oson (`wg0.conf`) |
| Performance | yaxshi | a'lo (kernel module) |

Tailscale — WireGuard ustida managed mesh VPN.
</details>

### Q47. L4 vs L7 load balancer?

<details>
<summary>Javobni ko'rish</summary>

- **L4 (TCP/UDP):** IP+port asosida, payload o'qimaydi. Tez (kernel-space, IPVS/eBPF), TLS passthrough. Misol: AWS NLB, HAProxy TCP.
- **L7 (HTTP):** header/URL/cookie asosida smart routing (path -> service). Sekinroq (TLS decrypt, HTTP parse), lekin WAF, rate-limit per-route. Misol: nginx, Envoy, AWS ALB.

Real arxitektura: L4 (DDoS) -> L7 (routing) -> backend.
</details>

### Q48. Kubernetes Pod-to-Pod networking?

<details>
<summary>Javobni ko'rish</summary>

4 ta talab: har Pod o'z IP; Pod<->Pod NAT'siz; Pod<->Node NAT'siz; Service IP virtual.

**CNI:** Flannel (VXLAN), Calico (BGP), Cilium (eBPF). **kube-proxy:** iptables (default, O(N)), IPVS (hash O(1)), eBPF (fastest). DNS: CoreDNS, `<name>.<ns>.svc.cluster.local`.
</details>

### Q49. Service mesh sidecar proxy (Istio/Linkerd)?

<details>
<summary>Javobni ko'rish</summary>

Har Pod yonida L7 proxy (Envoy). App trafigi avval sidecar dan o'tadi. Imkoniyat: mTLS auto rotation; canary/A-B routing; circuit breaking/retry/timeout; observability.

Cost: latency ~1-2 ms, resource ~50-100 MB RAM/Pod. eBPF (Cilium) — sidecar-less alternativa.
</details>

### Q50. Symmetric NAT orqasidagi P2P (STUN/TURN/ICE)?

<details>
<summary>Javobni ko'rish</summary>

Symmetric NAT da har destination uchun yangi external port — STUN ishlamaydi.

**ICE (RFC 8445):** STUN (o'z external IP:port ni topish) -> hole punching -> symmetric bo'lsa fail -> **TURN** relay media ni o'rtadan o'tkazadi (bandwidth qimmat, lekin universal). WebRTC candidate: host, srflx (STUN), relay (TURN).
</details>

### Q51. DNS cache poisoning va DNSSEC?

<details>
<summary>Javobni ko'rish</summary>

**Poisoning (Kaminsky 2008):** hujumchi resolver ga yolg'on response yuboradi (16 bit transaction ID brute force). Himoya: source port randomization (RFC 5452).

**DNSSEC (RFC 4033):** RRSIG (signature), DNSKEY (public key), DS (delegation). Chain of trust: root -> TLD -> zone. Resolver signature verify qiladi. Adoption past (~10-25%), key rollover qiyin.
</details>

### Q52. eBPF networking?

<details>
<summary>Javobni ko'rish</summary>

eBPF — kernel da safe sandboxed dastur. Networking da: **XDP** (NIC driver da, eng tez, DDoS/LB); **TC** (container ingress/egress); socket filtering; kprobes (observability).

Real: Cilium (K8s CNI), Katran (Facebook L4 LB, 10M PPS), Cloudflare DDoS. Kamchilik: kernel 4.18+, debugging qiyin (verifier).
</details>

---

## 8. Backend end-to-end

### Q53. "google.com ni browserga yozganda nima bo'ladi?"

<details>
<summary>Javobni ko'rish</summary>

1. **DNS:** browser/OS cache -> `hosts` -> DNS server (UDP 53). Recursive: root -> `.com` -> authoritative. IP keladi.
2. **ARP:** default gateway MAC ni topish (cache da yo'q bo'lsa).
3. **TCP handshake:** SYN -> SYN-ACK -> ACK (port 443).
4. **TLS handshake:** ClientHello -> ServerHello + cert -> key exchange -> Finished (TLS 1.3 da 1 RTT).
5. **HTTP request:** `GET / HTTP/2` + headers (Host, Cookie...).
6. **HTTP response:** status + headers + body (HTML).
7. **Browser parse:** HTML -> DOM, CSS/JS/image — har biri yangi request (HTTP/2 multiplexing).
8. **Render** + JS execute.

Yopilish: TLS close_notify -> TCP FIN -> TIME_WAIT.
</details>

### Q54. Connection pooling nima uchun kerak? (backend)

<details>
<summary>Javobni ko'rish</summary>

Har HTTP request uchun yangi TCP ulanish — bu har safar 3-way handshake + TLS negotiation + slow start. 50 ta resurs bo'lgan sahifa = 50 ta ulanish sozlash.

**Connection pool** oldindan ochilgan ulanishlarni qayta ishlatadi. Pool to'la bo'lsa, eski ulanish tartibli yopiladi va TIME_WAIT ga o'tadi. Keep-alive: fewer yangi ulanish -> fewer port band -> fewer TIME_WAIT.

Go, Java HttpClient, .NET HttpClient buni default qiladi — lekin har requestda yangi client yaratsang, pool ishlamaydi (Case 12/14 ga qara).
</details>

### Q55. HTTP keep-alive nima va nima uchun latency ni kamaytiradi?

<details>
<summary>Javobni ko'rish</summary>

Keep-alive (persistent connection) — bitta TCP ulanishni ko'p HTTP request uchun qayta ishlatadi (`Connection: keep-alive`). Har request uchun handshake takrorlanmaydi -> latency kamayadi, CPU/port tejaladi.

HTTP/1.1 da default yoqilgan. Reverse proxy (nginx -> upstream) da alohida sozlash kerak (`keepalive 64;`), aks holda upstream ga har requestda yangi ulanish (Case 9 ni ko'r).
</details>

### Q56. Sticky session vs Round-robin LB?

<details>
<summary>Javobni ko'rish</summary>

- **Round-robin:** har request keyingi backend ga. Stateless app uchun.
- **Sticky (session affinity):** bir client ning hamma so'rovi bir backend ga (cookie yoki source IP hash).

Stateless arxitektura afzal — session shared storage (Redis) da bo'lsa, sticky kerak emas. Sticky da backend tushsa session yo'qoladi (failover muammoli).
</details>

---

## Yakunda

Bu savollar **intervyu tayyorgarligi** uchun tayanch. Real chuqurlik uchun har mavzuni tegishli modulda o'qi:

- Transport (TCP/UDP) — 04-transport-layer moduli;
- IP addressing/subnetting — 02-network-layer-ip moduli;
- Routing (OSPF/BGP) — 03-routing moduli;
- HTTP/DNS/TLS — 05-application-layer va 06-api-protokollari modullari;
- Security — 08-security moduli;
- Amaliy troubleshooting — [02-troubleshooting-cases.md](./02-troubleshooting-cases.md).

> **Yakuniy tip:** Intervyuda diagramma chizib tushuntir (OSI stack, TCP handshake, DNS resolution). Vizual javob "bilaman"dan "chuqur tushunaman"ga farqni ko'rsatadi.

## 📚 Manbalar

- Manba savollar: kursning eski `network-knowledge` va `Computer Networking` konspektlaridan singdirilgan (git tarixida saqlangan)
- [InterviewBit: 70+ Networking Interview Questions](https://www.interviewbit.com/networking-interview-questions/)
- [GeeksforGeeks: Top 50 TCP/IP Interview Questions](https://www.geeksforgeeks.org/blogs/top-50-tcp-ip-interview-questions-and-answers/)
- [Apache HttpComponents: Connections in TIME_WAIT State](https://cwiki.apache.org/confluence/display/HTTPCOMPONENTS/FrequentlyAskedConnectionManagementQuestions)
- [Microsoft DevBlogs: The Art of HTTP Connection Pooling](https://devblogs.microsoft.com/premier-developer/the-art-of-http-connection-pooling-how-to-optimize-your-connections-for-peak-performance/)
- [MDN: HTTP Keep-Alive header](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Keep-Alive)
