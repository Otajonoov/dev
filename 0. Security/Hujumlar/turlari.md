ARP Spoofing (security threat): zararli host noto'g'ri ARP reply yuborib, gateway ning MAC address i o'rniga o'z MAC address ini "joylab" oladi va MITM (man-in-the-middle) hujum o'tkazadi. Himoya: dynamic ARP inspection (DAI), arpwatch.

MAC flooding (security): attacker switch CAM table ni soxta MAC lar bilan to'ldiradi, switch hub kabi ishlay boshlaydi. port-security bilan oldini olish.

ARP spoofing: arpwatch yoki arp-scan bilan duplicate IP-MAC mapping ni topish.

- **Broadcast storm:** STP (Spanning Tree Protocol) yo'q yoki sinmagan — frame loop bo'lib qoladi va network "qulab" tushadi.

Ko'p uchraydigan xavflar:

Muammo	Izoh	Himoya
IP spoofing	Source IP soxtalashtiriladi	ingress/egress filtering, uRPF
Broadcast abuse	Broadcast orqali amplification	directed broadcastni bloklash
ARP spoofing	LAN ichida gateway MAC soxtalashtiriladi	DHCP snooping, dynamic ARP inspection
Fragmentation evasion	Fragmentlar bilan firewallni aldash	stateful firewall, reassembly
Private IP leak	RFC1918 address Internetga chiqib ketadi	edge filtering
Muhim qoida: