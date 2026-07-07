Протокол – это ==набор правил, который определяет, как обмениваются данными устройства и программы==. В информатике это может быть протокол передачи данных, который используется для связывания компьютеров в сети.

| Protocol number | Protocol |
|---:|---|
| 1 | ICMP |
| 2 | IGMP |
| 6 | TCP |
| 17 | UDP |
| 47 | GRE |
| 50 | ESP |
| 51 | AH |
| 89 | OSPF |


Ko'p uchraydigan xavflar:

| Muammo | Izoh | Himoya |
|---|---|---|
| IP spoofing | Source IP soxtalashtiriladi | ingress/egress filtering, uRPF |
| Broadcast abuse | Broadcast orqali amplification | directed broadcastni bloklash |
| ARP spoofing | LAN ichida gateway MAC soxtalashtiriladi | DHCP snooping, dynamic ARP inspection |
| Fragmentation evasion | Fragmentlar bilan firewallni aldash | stateful firewall, reassembly |
| Private IP leak | RFC1918 address Internetga chiqib ketadi | edge filtering |

ARP spoofing nega IPv4 bilan bog'liq?
IPv4 lokal tarmoqda MAC address topish uchun ARP'ga tayanadi. ARP o'zi ishonchli authentication qilmaydi.

Shuning uchun attacker:

Gateway IP menman, MAC mana shu.
deb yolg'on ARP reply yuborishi mumkin. Bu man-in-the-middle yoki traffic hijack xavfiga olib keladi.

