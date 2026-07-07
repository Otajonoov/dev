# SNMP: Monitoring uchun asoslar

SNMP (Simple Network Management Protocol) router, switch, firewall va serverlardan holat ma'lumotlarini olish uchun ishlatiladi. Monitoring tizimlari SNMP orqali interface traffic, CPU, memory, uptime va boshqa ko'rsatkichlarni o'qiydi.

## SNMP komponentlari

- Managed device: router yoki switch.
- NMS (Network Management System): monitoring server, masalan Zabbix, PRTG, LibreNMS.
- Agent: qurilma ichidagi SNMP xizmati.
- MIB: o'qiladigan obyektlar bazasi.
- OID: aniq bir obyekt identifikatori.

## SNMP versiyalari

| Versiya | Xususiyat | CCNA uchun eslatma |
|---|---|---|
| SNMPv1 | Eski, community string | Tavsiya qilinmaydi |
| SNMPv2c | Oddiy, tez, community string | Ko'p lablarda uchraydi |
| SNMPv3 | Auth va encryption qo'llaydi | Eng xavfsiz variant |

## SNMPv2c konfiguratsiya

Read-only community:

```cisco
conf t
access-list 50 permit 192.168.100.50
snmp-server community CCNA-RO RO 50
snmp-server location Tashkent-DC-Rack1
snmp-server contact netops@example.local
end
```

Bu yerda faqat `192.168.100.50` monitoring serveri SNMP o'qiy oladi.

Trap yuborish:

```cisco
conf t
snmp-server host 192.168.100.50 version 2c CCNA-RO
snmp-server enable traps
end
```

## SNMPv3 konfiguratsiya

SNMPv3 username, authentication va privacy ishlatadi.

```cisco
conf t
snmp-server group NMS-GROUP v3 priv
snmp-server user nmsuser NMS-GROUP v3 auth sha AuthPass123 priv aes 128 PrivPass123
snmp-server host 192.168.100.50 version 3 priv nmsuser
end
```

Umumiy ma'no:

- `auth`: foydalanuvchini tekshiradi.
- `priv`: trafikni shifrlaydi.
- `sha`: authentication algoritmi.
- `aes`: encryption algoritmi.

## Tekshiruv buyruqlari

```cisco
show snmp
show snmp user
show snmp group
show running-config | include snmp
show access-lists
```

NMS serverdan test:

```bash
snmpwalk -v2c -c CCNA-RO 192.168.1.1 sysName
snmpwalk -v3 -l authPriv -u nmsuser -a SHA -A AuthPass123 -x AES -X PrivPass123 192.168.1.1 sysName
```

## SNMP va xavfsizlik

- Community stringni password kabi ko'ring.
- `public` va `private` communitylardan foydalanmang.
- SNMPni management VLAN yoki ACL bilan cheklang.
- Imkon bo'lsa SNMPv3 ishlating.
- Internet tomondan SNMP ochiq qolmasin.

## Muammolar va yechimlar

| Muammo | Sabab | Tekshiruv |
|---|---|---|
| NMS qurilmani o'qiy olmayapti | ACL monitoring serverni bloklagan | `show access-lists` |
| SNMPv2c ishlamayapti | Community string noto'g'ri | Config va NMS sozlamasi |
| SNMPv3 ishlamayapti | Auth/priv parol yoki algoritm mos emas | `show snmp user` |
| Trap kelmayapti | `snmp-server host` yoki traps yo'q | `show run | include snmp-server` |
| OID topilmayapti | MIB/NMS moslashmagan | NMS MIB sozlamalari |

## Keng tarqalgan xatolar

- SNMP communityni ACLsiz ochib qo'yish.
- SNMPv3da `authPriv`, `authNoPriv`, `noAuthNoPriv` darajalarini aralashtirish.
- NMS IP o'zgarganda ACLni yangilamaslik.
- Trapni sozlab, lekin `snmp-server enable traps`ni unutish.

## Q&A

**Savol:** SNMP polling va trap farqi nima?  
**Javob:** Pollingda NMS qurilmadan ma'lumot so'raydi. Trapda qurilma hodisa bo'lsa NMSga o'zi xabar yuboradi.

**Savol:** SNMPv2c xavfsizmi?  
**Javob:** Cheklangan. Community string ochiq matn sifatida ketadi. Management tarmoq va ACL bilan cheklash zarur.

**Savol:** CCNAda SNMPv3ni chuqur bilish kerakmi?  
**Javob:** Asosiy tushuncha, auth/priv farqi va umumiy konfiguratsiya mantiqi yetarli.

