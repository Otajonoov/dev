# SSH boshqaruv, TFTP va FTP backup

Router va switchlarni masofadan boshqarish uchun Telnet o'rniga SSH ishlatiladi. Telnet username va passwordni ochiq matnda yuboradi. SSH esa sessiyani shifrlaydi.

TFTP va FTP konfiguratsiya backup/restore uchun ishlatiladi. TFTP oddiy, lekin autentifikatsiyasiz. FTP username/password bilan ishlaydi, lekin klassik FTP ham shifrlanmagan.

## SSH uchun asosiy talablar

- Hostname default bo'lmasligi kerak.
- Domain name kerak.
- RSA key yaratiladi.
- Local user kerak.
- VTY line SSH qabul qilishi kerak.

## SSH konfiguratsiya

```cisco
conf t
hostname R1
ip domain-name example.local
username admin privilege 15 secret StrongPass123
crypto key generate rsa modulus 2048
ip ssh version 2

line vty 0 4
 login local
 transport input ssh
 exec-timeout 10 0
end
```

Management source interface:

```cisco
conf t
ip ssh source-interface loopback0
end
```

VTYga ACL bilan faqat admin subnetga ruxsat berish:

```cisco
conf t
ip access-list standard MGMT-ONLY
 permit 192.168.100.0 0.0.0.255
 deny any

line vty 0 4
 access-class MGMT-ONLY in
end
```

## SSH tekshiruv

```cisco
show ip ssh
show ssh
show running-config | section line vty
show users
```

Clientdan:

```bash
ssh admin@192.168.100.1
```

## TFTP backup

Running-configni TFTP serverga saqlash:

```cisco
copy running-config tftp:
```

IOS odatda quyidagilarni so'raydi:

```text
Address or name of remote host []? 192.168.100.70
Destination filename [r1-confg]? R1-running-config.cfg
```

Startup-config backup:

```cisco
copy startup-config tftp:
```

TFTPdan konfiguratsiya qayta yuklash:

```cisco
copy tftp: running-config
```

Ehtiyot bo'ling: bu mavjud running-config ustiga merge qiladi, to'liq almashtirish emas.

## FTP backup

FTP uchun username/password sozlash:

```cisco
conf t
ip ftp username backupuser
ip ftp password BackupPass123
end
```

Backup:

```cisco
copy running-config ftp:
```

URL ko'rinishida ham ishlatish mumkin:

```cisco
copy running-config ftp://backupuser:BackupPass123@192.168.100.80/R1-running-config.cfg
```

## IOS image backup

Flashdagi fayllarni ko'rish:

```cisco
dir flash:
```

IOS imageni TFTPga ko'chirish:

```cisco
copy flash:c2900-universalk9-mz.SPA.157-3.M.bin tftp:
```

## Troubleshooting buyruqlari

```cisco
show ip interface brief
show ip route
ping 192.168.100.70
ping 192.168.100.70 source loopback0
show access-lists
show ip ssh
show users
dir flash:
```

TFTP ishlamasa, serverdagi firewall va TFTP root directoryni tekshiring. TFTP odatda UDP 69dan boshlaydi.

## Keng tarqalgan xatolar

- SSH uchun `ip domain-name` berilmaganligi.
- RSA key yaratilmaganligi.
- VTYda `transport input ssh` o'rniga Telnet ham ochiq qolishi.
- `login local` berilmaganligi yoki local user yo'qligi.
- TFTP serverga ping bor, lekin firewall UDP 69ni bloklayapti.
- `copy tftp: running-config` merge qilishini unutish.
- Backup fayl nomlarini tartibsiz qo'yish; qurilma nomi va sana yozish foydali.

## Xavfsizlik tavsiyalari

- `enable secret` ishlating, `enable password` emas.
- Local user uchun `secret` ishlating.
- SSH version 2ni majburiy qiling.
- VTYga ACL qo'ying.
- Keraksiz service va eski protokollarni o'chiring.

```cisco
conf t
no ip http server
no ip http secure-server
service password-encryption
end
```

`service password-encryption` kuchli himoya emas, faqat oddiy ko'rinishda password chiqishini kamaytiradi.

## Q&A

**Savol:** Telnet va SSH farqi nima?  
**Javob:** Telnet trafikni shifrlamaydi. SSH username, password va sessiya ma'lumotlarini shifrlaydi.

**Savol:** TFTP xavfsizmi?  
**Javob:** Yo'q, u autentifikatsiyasiz va oddiy. Uni management tarmoq ichida va ACL/firewall bilan cheklab ishlatish kerak.

**Savol:** Running-config va startup-config farqi nima?  
**Javob:** Running-config hozir RAMda ishlayotgan konfiguratsiya. Startup-config reloaddan keyin yuklanadigan NVRAMdagi konfiguratsiya.

