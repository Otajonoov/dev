# 08. Security (Tarmoq xavfsizligi)

Bu modul tarmoq xavfsizligini CCNA darajasida — nazariya bilan birga real
Cisco IOS konfiguratsiyasi, tekshirish buyruqlari va ko'p uchraydigan
xatolar orqali o'rgatadi. Har dars: muammo → analogiya → diagramma →
worked example (Cisco CLI) → o'z-o'zini tekshir → amaliyot.

## Nima o'rganiladi

- Xavfsizlik atamalari (CIA triad, threat/vulnerability/exploit/mitigation)
  va zamonaviy hujum turlari (malware, phishing, DDoS, MITM).
- Trafikni nazorat qilish: firewall (stateless/stateful/NGFW/WAF) va ACL.
- Qurilma boshqaruvini himoyalash: SSH, enable secret, banner, login himoyasi.
- Markazlashgan autentifikatsiya: AAA, RADIUS, TACACS+.
- Layer 2 himoyalari: port security, DHCP snooping, DAI, VLAN hopping.
- Simsiz xavfsizlik: WEP → WPA3 evolyutsiyasi, PSK vs Enterprise.
- Shifrlangan tunnellar: IPsec VPN va zamonaviy WireGuard.

## Darslar

1. [Security concepts va hujumlar](./01-security-concepts-va-hujumlar.md) —
   CIA triad, threat/vulnerability/exploit/mitigation, attack surface, hujum turlari.
2. [Firewall](./02-firewall.md) — stateless vs stateful, NGFW, WAF, zero trust.
3. [ACL](./03-acl.md) — standard vs extended, wildcard mask, direction, placement.
4. [Device access security](./04-device-access-security.md) — SSH, enable secret,
   banner, login protection, management ACL.
5. [AAA, RADIUS, TACACS+](./05-aaa-radius-tacacs.md) — authentication,
   authorization, accounting va server protokollari.
6. [Layer 2 security](./06-l2-security.md) — port security, DHCP snooping,
   DAI, VLAN hopping, STP himoyalari.
7. [Wireless security](./07-wireless-security.md) — WPA2/WPA3, SAE, PSK vs
   Enterprise, guest Wi-Fi, threatlar.
8. [VPN va IPsec](./08-vpn-ipsec.md) — site-to-site vs remote access, IKE,
   ESP/AH, tunnel mode, WireGuard.

## O'qish tartibi

Darslar ketma-ket qurilgan: avval **til va tahdid** (1), keyin **trafik
nazorati** (2-3), so'ng **qurilma va login** (4-5), keyin **L2 va wireless**
(6-7), oxirida **VPN** (8). Har dars oldingisiga tayanadi (masalan AAA
RADIUS'i wireless Enterprise'da qaytadi), shuning uchun tartib bilan o'qish
tavsiya etiladi.

> Xavfsizlik oltin qoidasi: u qatlam-qatlam quriladi (**defense in depth**).
> Bitta himoya teshilsa, keyingisi xavfni ushlab qoladi.
