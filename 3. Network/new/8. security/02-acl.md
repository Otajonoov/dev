# 02. ACL

ACL (Access Control List) - router yoki multilayer switch orqali o'tayotgan trafikni ruxsat berish (`permit`) yoki rad etish (`deny`) uchun ishlatiladigan qoidalar ro'yxati. CCNA uchun ACL direction, placement, wildcard mask va troubleshooting juda muhim.

## ACL qanday ishlaydi?

ACL yuqoridan pastga qarab tekshiriladi. Birinchi mos kelgan qoida ishlaydi, keyingi qoidalar tekshirilmaydi.

Har bir ACL oxirida ko'rinmas qoida bor:

```text
deny ip any any
```

Shuning uchun kerakli trafikni albatta `permit` qilish kerak.

## Standard ACL

Standard ACL faqat **source IP** bo'yicha tekshiradi.

Raqam oralig'i:

- `1-99`
- `1300-1999`

Named standard ACL ham ishlatiladi.

Misol: `192.168.10.0/24` tarmog'iga ruxsat berish.

```cisco
conf t
access-list 10 permit 192.168.10.0 0.0.0.255
access-list 10 deny any
interface g0/0
 ip access-group 10 in
end
```

Named ACL:

```cisco
conf t
ip access-list standard MGMT_ONLY
 permit 192.168.10.0 0.0.0.255
 deny any
line vty 0 4
 access-class MGMT_ONLY in
end
```

`access-class` VTY line uchun ishlatiladi, `ip access-group` esa interface uchun.

## Extended ACL

Extended ACL source, destination, protocol va port bo'yicha tekshiradi.

Raqam oralig'i:

- `100-199`
- `2000-2699`

Misol: `192.168.10.0/24` foydalanuvchilariga web server `172.16.1.10` ga HTTP/HTTPS ruxsat, boshqa trafikni bloklash.

```cisco
conf t
ip access-list extended USERS_TO_WEB
 permit tcp 192.168.10.0 0.0.0.255 host 172.16.1.10 eq 80
 permit tcp 192.168.10.0 0.0.0.255 host 172.16.1.10 eq 443
 deny ip 192.168.10.0 0.0.0.255 any
 permit ip any any
interface g0/1
 ip access-group USERS_TO_WEB in
end
```

Oxirdagi `permit ip any any` boshqa tarmoqlarning trafikini tasodifan bloklab qo'ymaslik uchun qo'shilgan. Real tarmoqda security policy bo'yicha aniqroq qoida yoziladi.

## Wildcard mask

Wildcard mask subnet maskning teskarisiga o'xshaydi:

```text
/24 mask:      255.255.255.0
wildcard:        0.0.0.255

/30 mask:      255.255.255.252
wildcard:        0.0.0.3
```

Tez hisoblash:

```text
255.255.255.255
- subnet mask
= wildcard mask
```

Host uchun:

```cisco
host 192.168.1.10
```

Bu quyidagiga teng:

```cisco
192.168.1.10 0.0.0.0
```

Any uchun:

```cisco
any
```

Bu quyidagiga teng:

```cisco
0.0.0.0 255.255.255.255
```

## Direction: in yoki out

ACL interfacega ikki yo'nalishda qo'yiladi:

- **in** - trafik interfacega kirayotganda tekshiriladi.
- **out** - trafik interfacedan chiqayotganda tekshiriladi.

Misol:

```text
PC --- g0/1 Router g0/0 --- Server
```

Agar PCdan kelayotgan trafikni routerga kirishda to'xtatmoqchi bo'lsangiz:

```cisco
interface g0/1
 ip access-group ACL_NAME in
```

Agar server tomonga chiqayotgan trafikni tekshirmoqchi bo'lsangiz:

```cisco
interface g0/0
 ip access-group ACL_NAME out
```

## Placement: qayerga qo'yish kerak?

Umumiy qoida:

- **Standard ACL** - destinationga yaqin joylashtiriladi, chunki faqat source IP ko'radi.
- **Extended ACL** - sourcega yaqin joylashtiriladi, chunki aniq trafikni erta bloklay oladi.

Misol:

```text
Standard ACL: "192.168.10.0/24 hech qayerga bormasin" kabi qo'pol filtr.
Extended ACL: "192.168.10.0/24 faqat 172.16.1.10:443 ga borsin" kabi aniq filtr.
```

## ACL sequence number

Named ACL ichida qoidalarni tartib bilan boshqarish mumkin:

```cisco
conf t
ip access-list extended USERS_TO_WEB
 10 permit tcp 192.168.10.0 0.0.0.255 host 172.16.1.10 eq 443
 20 deny ip 192.168.10.0 0.0.0.255 any
end
```

Qoida qo'shish:

```cisco
conf t
ip access-list extended USERS_TO_WEB
 15 permit tcp 192.168.10.0 0.0.0.255 host 172.16.1.10 eq 80
end
```

Qoida o'chirish:

```cisco
conf t
ip access-list extended USERS_TO_WEB
 no 20
end
```

## VTY access uchun ACL

Management kirishni faqat admin tarmog'idan ruxsat berish:

```cisco
conf t
ip access-list standard ADMIN_NET
 permit 10.10.10.0 0.0.0.255
 deny any
line vty 0 4
 access-class ADMIN_NET in
 transport input ssh
end
```

## Troubleshooting commands

```cisco
show access-lists
show ip access-lists
show running-config | section access-list
show running-config interface g0/1
show ip interface g0/1
show logging
```

Counterlarni tozalash:

```cisco
clear access-list counters
```

Tekshirishda quyidagilarga qarang:

- ACL interfacega qo'yilganmi?
- Direction to'g'rimi?
- Source va destination joyi almashib ketmaganmi?
- Wildcard mask to'g'rimi?
- Kerakli `permit` qoida implicit deny dan oldin turibdimi?
- NAT yoki routing ACLdan oldin/keyin trafikni o'zgartiryaptimi?

## Common mistakes

- Extended ACLda source va destinationni almashtirib yozish.
- `deny ip any any` ni qo'yib, keyin kerakli `permit` qoidani pastga yozish.
- Standard ACLni sourcega yaqin qo'yib, source tarmoqning hamma trafikini keraksiz bloklash.
- VTY uchun `ip access-group` ishlatish; VTYda `access-class` kerak.
- Wildcard mask o'rniga subnet mask yozish.

## Q&A

**Savol:** ACL firewall o'rnini bosa oladimi?

**Javob:** Oddiy filtrlash uchun yordam beradi, lekin stateful inspection, application control va advanced threat protection uchun firewall kerak.

**Savol:** Bitta interfacega nechta ACL qo'yish mumkin?

**Javob:** Har bir protocol, har bir direction uchun odatda bittadan ACL. Masalan IPv4 inbound bitta, IPv4 outbound bitta.

**Savol:** Nega `show access-lists` counter nol?

**Javob:** Trafik ACLdan o'tmayotgan bo'lishi, ACL noto'g'ri interfeys/directionga qo'yilgan bo'lishi yoki qoida mos kelmayotgan bo'lishi mumkin.
