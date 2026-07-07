# 04. AAA, RADIUS, TACACS+

AAA - Authentication, Authorization, Accounting. Katta tarmoqlarda har bir router yoki switchda alohida user-parol saqlash noqulay va xavfli. AAA markazlashgan login, huquq va audit beradi.

## AAA nimani anglatadi?

- **Authentication** - foydalanuvchi kim? Login/parol, token, sertifikat.
- **Authorization** - foydalanuvchi nima qila oladi? Masalan, faqat `show` buyruqlari yoki full admin.
- **Accounting** - foydalanuvchi nima qildi? Login va buyruqlar logi.

Oddiy misol:

```text
Admin SSH qiladi -> R1 AAA serverdan tekshiradi -> ruxsat bo'lsa CLI ochiladi -> bajarilgan buyruqlar accountingga yoziladi
```

## RADIUS va TACACS+ farqi

| Xususiyat | RADIUS | TACACS+ |
|---|---|---|
| Transport | UDP | TCP |
| Port | UDP 1812/1813 | TCP 49 |
| Encryption | Asosan parol qismi | Butun payload |
| AAA ajratilishi | Authentication va authorization ko'proq birga | Authentication, authorization, accounting alohida |
| Ko'p ishlatiladi | Network access, Wi-Fi, VPN | Network device administration |

CCNA eslab qolish:

- Wi-Fi 802.1X va VPN uchun RADIUS ko'p uchraydi.
- Router/switch admin access uchun TACACS+ qulay.

## AAA yoqish

AAA ishlashi uchun global rejimda yoqiladi:

```cisco
conf t
aaa new-model
end
```

Ehtiyot bo'ling: `aaa new-model` mavjud login usullariga ta'sir qilishi mumkin. Masofadan ishlayotgan bo'lsangiz, local fallback tayyor bo'lsin.

## Local fallback bilan TACACS+

```cisco
conf t
aaa new-model

username localadmin privilege 15 secret LocalAdminSecret123

tacacs server TAC1
 address ipv4 10.10.10.30
 key TacacsSharedKey123

aaa group server tacacs+ TAC-GROUP
 server name TAC1

aaa authentication login default group TAC-GROUP local
aaa authorization exec default group TAC-GROUP local
aaa accounting exec default start-stop group TAC-GROUP

line vty 0 4
 login authentication default
 transport input ssh
end
```

`group TAC-GROUP local` ma'nosi: avval TACACS+ serverdan tekshiradi, server ishlamasa local userga tushadi.

## RADIUS misoli

```cisco
conf t
aaa new-model

username localadmin privilege 15 secret LocalAdminSecret123

radius server RAD1
 address ipv4 10.10.10.40 auth-port 1812 acct-port 1813
 key RadiusSharedKey123

aaa group server radius RAD-GROUP
 server name RAD1

aaa authentication login default group RAD-GROUP local
aaa authorization exec default group RAD-GROUP local
aaa accounting exec default start-stop group RAD-GROUP

line vty 0 4
 login authentication default
 transport input ssh
end
```

## Method list

AAA method list - qaysi login uchun qaysi tekshiruv usuli ishlatilishini belgilaydi.

Default method list:

```cisco
aaa authentication login default group TAC-GROUP local
```

Named method list:

```cisco
aaa authentication login SSH_LOGIN group TAC-GROUP local

line vty 0 4
 login authentication SSH_LOGIN
```

Console uchun alohida qilish mumkin:

```cisco
aaa authentication login CONSOLE_LOGIN local

line console 0
 login authentication CONSOLE_LOGIN
```

## Authorization

Exec shell ochishga ruxsat:

```cisco
aaa authorization exec default group TAC-GROUP local
```

Buyruqlar bo'yicha authorization:

```cisco
aaa authorization commands 15 default group TAC-GROUP local
```

Bu production muhitda foydali, lekin noto'g'ri sozlansa adminni kerakli komandadan ham mahrum qilishi mumkin.

## Accounting

Login sessiyalarini yozish:

```cisco
aaa accounting exec default start-stop group TAC-GROUP
```

Privilege 15 buyruqlarini yozish:

```cisco
aaa accounting commands 15 default start-stop group TAC-GROUP
```

## Test va troubleshooting

Serverga reachability:

```cisco
ping 10.10.10.30
traceroute 10.10.10.30
```

AAA tekshirish:

```cisco
show running-config | section aaa
show aaa servers
show tacacs
show radius statistics
test aaa group tacacs+ admin AdminPassword legacy
test aaa group radius admin AdminPassword legacy
```

Debug:

```cisco
debug aaa authentication
debug aaa authorization
debug tacacs
debug radius
```

Debugni tugatish:

```cisco
undebug all
```

## Common mistakes

- `aaa new-model`dan oldin local admin yaratmaslik.
- `local` fallback qo'shmaslik.
- Shared key router va serverda bir xil emas.
- RADIUS/TACACS+ portlari firewall tomonidan bloklangan.
- Qurilma source IP manzili serverda ruxsat etilmagan.
- NTP noto'g'ri bo'lgani sabab log va accounting vaqti chalkash.
- Method list yaratilgan, lekin `line vty`ga biriktirilmagan.

## Q&A

**Savol:** AAA server ishlamasa qurilmaga kira olamanmi?

**Javob:** Agar method listda `local` fallback bo'lsa va local user mavjud bo'lsa, ha.

**Savol:** RADIUS yoki TACACS+ qaysi biri yaxshi?

**Javob:** Device administration uchun ko'pincha TACACS+ yaxshi, network access va Wi-Fi uchun RADIUS ko'p ishlatiladi.

**Savol:** Accounting nima uchun kerak?

**Javob:** Kim qachon kirgani va nima qilganini bilish uchun. Bu troubleshooting va audit uchun juda muhim.
