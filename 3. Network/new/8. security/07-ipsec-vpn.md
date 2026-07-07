# 07. IPsec VPN

IPsec VPN IP trafikni autentifikatsiya, integrity va encryption orqali himoyalaydi. CCNA darajasida IPsec nima uchun kerakligi, asosiy komponentlari va site-to-site VPN ishlash tartibini tushunish muhim.

## VPN nima?

VPN - ishonchsiz tarmoq, masalan Internet orqali, xususiy tarmoqlarni xavfsiz bog'lash usuli.

Ko'p uchraydigan turlar:

- **Site-to-site VPN** - ikki filial tarmog'i orasida.
- **Remote access VPN** - foydalanuvchi laptop/telefonidan kompaniya tarmog'iga.

Misol:

```text
Branch LAN 192.168.10.0/24 -- R1 == Internet == R2 -- HQ LAN 10.10.10.0/24
```

## IPsec nimalarni beradi?

- **Confidentiality** - trafik shifrlanadi.
- **Integrity** - paket o'zgarmaganini tekshiradi.
- **Authentication** - tunnel peer haqiqiyligini tekshiradi.
- **Anti-replay** - eski paketlarni qayta yuborish hujumini kamaytiradi.

## IPsec komponentlari

### IKE

IKE (Internet Key Exchange) peerlar orasida xavfsiz kelishuv qiladi:

- Peer authentication.
- Encryption va hashing algoritmlarini kelishish.
- Key exchange.
- Security Association yaratish.

IKE versiyalari:

- **IKEv1** - eski, lekin hali uchraydi.
- **IKEv2** - yangi, soddaroq va yaxshiroq.

### ESP

ESP (Encapsulating Security Payload) IPsec trafikni shifrlash va integrity bilan himoyalashda ko'p ishlatiladi.

### AH

AH (Authentication Header) integrity/authentication beradi, lekin encryption bermaydi. Amaliy tarmoqlarda ESP ko'proq uchraydi.

## IKEv1 phase 1 va phase 2

IKEv1 tushunchasi:

- **Phase 1** - peerlar orasida IKE SA yaratiladi. Bu management tunnelga o'xshaydi.
- **Phase 2** - real data traffic uchun IPsec SA yaratiladi.

Phase 1 parametrlari:

- Encryption: AES, 3DES.
- Hash: SHA.
- Authentication: pre-shared key yoki certificate.
- DH group.
- Lifetime.

Phase 2 parametrlari:

- Transform set.
- Interesting traffic.
- PFS optional.
- Lifetime.

## Interesting traffic

IPsec qaysi trafikni tunnelga solishini ACL belgilaydi. Masalan, Branch LAN dan HQ LAN ga:

```cisco
ip access-list extended VPN_TRAFFIC
 permit ip 192.168.10.0 0.0.0.255 10.10.10.0 0.0.0.255
```

Bu ACL interface filter ACL emas. Bu IPsec uchun "qaysi trafik VPNga tushadi?" degan savolga javob beradi.

## Site-to-site IPsec misol

Quyidagi misol klassik crypto map uslubida. Real platforma va IOS versiyasiga qarab syntax farq qilishi mumkin.

R1:

```cisco
conf t
crypto isakmp policy 10
 encr aes 256
 hash sha
 authentication pre-share
 group 14
 lifetime 86400

crypto isakmp key SharedKey123 address 203.0.113.2

ip access-list extended VPN_TRAFFIC
 permit ip 192.168.10.0 0.0.0.255 10.10.10.0 0.0.0.255

crypto ipsec transform-set TS esp-aes 256 esp-sha-hmac
 mode tunnel

crypto map VPN-MAP 10 ipsec-isakmp
 set peer 203.0.113.2
 set transform-set TS
 match address VPN_TRAFFIC

interface g0/0
 description Internet
 crypto map VPN-MAP
end
```

R2 tomonda tarmoqlar teskari yoziladi:

```cisco
conf t
crypto isakmp policy 10
 encr aes 256
 hash sha
 authentication pre-share
 group 14
 lifetime 86400

crypto isakmp key SharedKey123 address 203.0.113.1

ip access-list extended VPN_TRAFFIC
 permit ip 10.10.10.0 0.0.0.255 192.168.10.0 0.0.0.255

crypto ipsec transform-set TS esp-aes 256 esp-sha-hmac
 mode tunnel

crypto map VPN-MAP 10 ipsec-isakmp
 set peer 203.0.113.1
 set transform-set TS
 match address VPN_TRAFFIC

interface g0/0
 description Internet
 crypto map VPN-MAP
end
```

## NAT exemption

VPN trafik NAT qilinmasligi kerak bo'lishi mumkin. Aks holda interesting traffic ACL mos kelmaydi.

Konseptual misol:

```cisco
ip access-list extended NAT_EXEMPT
 deny ip 192.168.10.0 0.0.0.255 10.10.10.0 0.0.0.255
 permit ip 192.168.10.0 0.0.0.255 any
```

Bu ACL NAT uchun ishlatiladi: HQ tomonga trafik NAT qilinmaydi, Internet tomonga NAT qilinadi.

## Troubleshooting commands

```cisco
show crypto isakmp sa
show crypto ipsec sa
show crypto map
show access-lists VPN_TRAFFIC
show running-config | section crypto
show ip route
ping 10.10.10.1 source 192.168.10.1
traceroute 10.10.10.1 source 192.168.10.1
```

IKEv2 uchun:

```cisco
show crypto ikev2 sa
show crypto ipsec sa
```

Debug:

```cisco
debug crypto isakmp
debug crypto ipsec
undebug all
```

Production tarmoqda debugni ehtiyotkorlik bilan ishlating.

## Muammo qayerda ekanini topish

1. Peer public IP ping bo'ladimi?
2. Crypto policy ikki tomonda mosmi?
3. Pre-shared key bir xilmi?
4. Interesting traffic ACL ikki tomonda mirror holatidami?
5. Route bor-mi?
6. NAT exemption to'g'rimi?
7. Firewall UDP 500, UDP 4500 va ESPni bloklamayaptimi?
8. `show crypto ipsec sa` counter oshyaptimi?

## Common mistakes

- Interesting traffic ACLda source/destination teskari yoki noto'g'ri.
- Ikki tomonda transform set mos emas.
- Pre-shared key bir xil emas.
- NAT exemption yo'q.
- Crypto map outside interfacega qo'yilmagan.
- Tunnelni test qilishda routerning outside IPsidan ping qilish; LAN source bilan test qilish kerak.
- Routing yo'q yoki return route noto'g'ri.

## Q&A

**Savol:** IPsec tunnel o'zi avtomatik ko'tariladimi?

**Javob:** Odatda interesting traffic chiqqanda tunnel negotiate bo'ladi. Shuning uchun LANdan LANga ping bilan test qiling.

**Savol:** ESP va IKE farqi nima?

**Javob:** IKE kelishuv va key exchange qiladi, ESP esa data trafikni himoyalaydi.

**Savol:** IPsec GRE bilan bir xilmi?

**Javob:** Yo'q. GRE tunnel yaratadi, lekin o'zi shifrlamaydi. IPsec GRE trafikni shifrlash uchun ishlatilishi mumkin.
