# 8. Security

Bu bo'lim CCNA darajasida tarmoq xavfsizligini tushunish va amalda sozlash uchun kerak bo'ladigan asosiy mavzularni jamlaydi. Maqsad - nazariya bilan birga Cisco IOS buyruqlari, tekshirish komandalarini va ko'p uchraydigan xatolarni ko'rsatish.

## Bo'lim xaritasi

1. [Security concepts](./01-security-concepts.md) - threat, vulnerability, exploit, mitigation, CIA triad, attack surface.
2. [ACL](./02-acl.md) - standard va extended ACL, direction, placement, wildcard mask, troubleshooting.
3. [Device access security](./03-device-access-security.md) - parollar, SSH, line security, banner, login protection.
4. [AAA, RADIUS, TACACS+](./04-aaa-radius-tacacs.md) - authentication, authorization, accounting va server protokollari.
5. [Layer 2 security](./05-l2-security.md) - DHCP snooping, Dynamic ARP Inspection, port security, STP himoyalari.
6. [Wireless security](./06-wireless-security.md) - WPA2/WPA3, PSK, Enterprise, guest Wi-Fi, WLC tekshiruvlari.
7. [IPsec VPN](./07-ipsec-vpn.md) - IPsec, IKE, ESP, site-to-site VPN umumiy ishlash tartibi.

## Xavfsizlikni o'rganishda asosiy fikr

Tarmoq xavfsizligi bitta buyruq bilan hal bo'ladigan narsa emas. U quyidagi qatlamlarda quriladi:

- **Physical security** - qurilmaga jismoniy kirishni cheklash.
- **Management security** - SSH, AAA, kuchli parol, log va monitoring.
- **Control plane security** - routing protokollari, STP, DHCP, ARP kabi xizmatlarni himoyalash.
- **Data plane security** - trafikni ACL, segmentation, VPN, firewall orqali nazorat qilish.
- **Endpoint va wireless security** - mijoz qurilmalar, Wi-Fi autentifikatsiyasi va shifrlash.

## Tez eslab qolish

- **Threat** - zarar yetkazishi mumkin bo'lgan xavf manbai.
- **Vulnerability** - tizimdagi zaif joy.
- **Exploit** - zaiflikdan foydalanish usuli.
- **Mitigation** - xavfni kamaytirish chorasi.
- **ACL** - trafikni permit yoki deny qilish uchun qoidalar ro'yxati.
- **AAA** - kim kirdi, nima qila oladi, nima qildi degan savollarga javob beradi.
- **DHCP snooping** - soxta DHCP serverlardan himoya qiladi.
- **DAI** - ARP spoofing hujumlarini kamaytiradi.
- **Port security** - switch portida ruxsat etilgan MAC manzillarni nazorat qiladi.
- **WPA2/WPA3** - Wi-Fi trafikini himoyalash standartlari.
- **IPsec** - IP trafikni autentifikatsiya va shifrlash orqali himoyalaydi.

## Amaliy yondashuv

Har bir xavfsizlik sozlamasidan keyin quyidagilarni tekshiring:

```text
show running-config
show ip interface brief
show interfaces status
show logging
```

Sozlama ishlamayotganda savolni shu tartibda bering:

1. Trafik qayerdan qayerga ketmoqda?
2. Qaysi interfeysda va qaysi yo'nalishda filtr bor?
3. Default holat permitmi yoki denymi?
4. Qurilma vaqt, DNS, gateway va route bo'yicha to'g'ri ishlayaptimi?
5. Log yoki counter nimani ko'rsatyapti?

## Kichik Q&A

**Savol:** CCNA uchun firewall chuqur o'rganilishi kerakmi?

**Javob:** Asosiy tushuncha kerak, lekin CCNA ko'proq ACL, device hardening, L2 security, wireless security va VPN asoslariga urg'u beradi.

**Savol:** Xavfsizlik sozlamalarini qayerdan boshlash kerak?

**Javob:** Avval management access: kuchli parollar, SSH, local user yoki AAA, keraksiz xizmatlarni o'chirish. Keyin ACL, L2 security va monitoring.

**Savol:** Eng ko'p uchraydigan xato nima?

**Javob:** ACL direction yoki placement noto'g'ri tanlanadi. Natijada kerakli trafik bloklanadi yoki zararli trafik umuman filtrdan o'tmaydi.
