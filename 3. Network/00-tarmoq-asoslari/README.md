# 00. Tarmoq asoslari (Network Fundamentals)

Bu modul — butun **Network kursi**ning poydevori. Bu yerda "tarmoq nima",
"Internet qanday ishlaydi", "protokol nima" va "qatlamli modellar" (OSI,
TCP/IP) kabi eng asosiy tushunchalarni o'rganamiz. Bularsiz keyingi modullar
(IP, routing, transport, security va h.k.) tushunarsiz bo'lib qoladi.

## Nimani o'rganasan?

- Tarmoq va Internet nima, ular qanday farq qiladi
- Protokol nima va nega qatlamlarga bo'linadi (protocol layering)
- Internetga qanday ulanamiz (DSL, fiber, 5G, Starlink)
- Ma'lumot tarmoq ichida qanday harakatlanadi (packet switching)
- Internet qanday tuzilgan (ISP ierarxiyasi, IXP, CDN)
- Nega internet sekinlashadi (latency, packet loss, throughput)
- OSI va TCP/IP modellari, encapsulation jarayoni

## Darslar ro'yxati (o'qish tartibi)

1. [01-tarmoq-va-internet-nima](01-tarmoq-va-internet-nima.md) — tarmoq, Internet, LAN/WAN, client-server
2. [02-protokol-nima](02-protokol-nima.md) — protokol, layering, RFC, IETF
3. [03-access-networks](03-access-networks.md) — DSL, cable, FTTH, 5G, Starlink
4. [04-network-core-va-packet-switching](04-network-core-va-packet-switching.md) — packet vs circuit switching
5. [05-internet-tuzilishi-isp](05-internet-tuzilishi-isp.md) — Tier 1/2/3, IXP, CDN, peering
6. [06-latency-loss-throughput](06-latency-loss-throughput.md) — 4 delay, queuing, bottleneck
7. [07-osi-modeli](07-osi-modeli.md) — 7 qatlam, PDU nomlari
8. [08-tcpip-modeli-va-encapsulation](08-tcpip-modeli-va-encapsulation.md) — TCP/IP, encapsulation
9. [09-glossary](09-glossary.md) — atamalar lug'ati (ma'lumotnoma)

## Qanday o'qish kerak?

- Darslarni **tartib bilan** o'qi — har biri oldingisiga tayanadi.
- Har darsda `## ✅ O'z-o'zini tekshir` savollariga javobni **oldin o'zing** o'yla, keyin och.
- `## 🛠 Amaliyot` topshiriqlarini albatta bajar — passiv o'qish yetarli emas.
- `## 🔁 Takrorlash` jadvaliga amal qil: ertaga → 3 kun → 1 hafta.
- Notanish atama uchrasa — [09-glossary](09-glossary.md) dan qidir.

## Amaliy tayyorgarlik

Darslardagi mashqlar uchun quyidagi buyruqlar foydali bo'ladi:
`ping`, `traceroute` (Windows: `tracert`), `ipconfig`/`ip addr`, `curl -v`,
`nslookup`, `tcpdump` yoki Wireshark.

---

Keyingi modul: **01-network-access** (Ethernet, MAC, VLAN, switching).
