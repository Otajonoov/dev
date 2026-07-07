# CCNA new: to'ldirilgan network knowledge xaritasi

Bu katalog `network-knowledge` ichida kamroq yoritilgan CCNA mavzularini alohida, amaliy va Cisco CLI bilan tushuntirish uchun ochildi.

Maqsad:

- Network Access mavzularini chuqurlashtirish: VLAN, trunk, STP, EtherChannel, CDP/LLDP, WLAN.
- IP Connectivity mavzularini amaliy qilish: routing table, static routing, OSPFv2, FHRP, IPv6.
- IP Services mavzularini to'ldirish: DHCP relay, NAT, NTP, SNMP, Syslog, QoS, SSH/TFTP/FTP.
- Security Fundamentals bo'limini alohida qilish: ACL, device hardening, AAA, L2 security, wireless security, IPsec.
- Automation/Programmability bo'limini qo'shish: SDN, REST API, JSON, Ansible/Terraform, AI/cloud management.

## Tuzilma

```text
5. network-access/
6. ip-connectivity/
7. ip-services/
8. security/
9. automation-programmability/
```

## O'qish tartibi

1. [Network Access](5.%20network-access/README.md)
   - Avval VLAN va trunkni tushun.
   - Keyin inter-VLAN routing, STP, EtherChannel.
   - Oxirida CDP/LLDP va wireless.

2. [IP Connectivity](6.%20ip-connectivity/README.md)
   - Avval routing table va longest prefix match.
   - Keyin static/default/floating route.
   - Keyin OSPFv2 va FHRP.
   - Oxirida IPv6 addressing, NDP va IPv6 routing.

3. [IP Services](7.%20ip-services/README.md)
   - DHCP relay, NAT/PAT, NTP, SNMP, Syslog.
   - QoS va management protokollarini keyin o'qi.

4. [Security](8.%20security/README.md)
   - Security terminology.
   - ACL.
   - Device access security.
   - AAA.
   - L2 va wireless security.
   - IPsec overview.

5. [Automation and Programmability](9.%20automation-programmability/README.md)
   - SDN va controller-based networking.
   - REST API, JSON.
   - Ansible/Terraform.
   - AI, ML va cloud network management.

## CCNA uslubidagi mental model

Har mavzuni quyidagi 5 savol bilan o'rgan:

1. Bu mavzu qaysi layerda ishlaydi?
2. Qaysi muammoni hal qiladi?
3. Cisco qurilmada qanday sozlanadi?
4. Qanday tekshiriladi?
5. Eng ko'p uchraydigan xato nima?

## Amaliy lab tavsiyasi

Packet Tracer, GNS3, EVE-NG yoki haqiqiy Cisco qurilmada quyidagi tartibda lab qil:

1. 2 ta switch + 2 ta VLAN.
2. Trunk + router-on-a-stick.
3. STP root bridge tanlash.
4. EtherChannel LACP.
5. Static route + default route.
6. Single-area OSPF.
7. HSRP gateway redundancy.
8. DHCP relay.
9. PAT orqali Internet simulation.
10. Standard va extended ACL.
11. SSH-only device management.
12. REST API va JSON mini test.

## Tezkor checklist

```text
L2: VLAN, trunk, STP, EtherChannel, CDP/LLDP, WLAN
L3: routing table, static route, OSPF, FHRP, IPv6
Services: DHCP relay, NAT, NTP, SNMP, Syslog, QoS, SSH, TFTP
Security: ACL, AAA, port security, DHCP snooping, DAI, WPA, IPsec
Automation: SDN, controller, REST, JSON, Ansible, Terraform, AI/cloud
```

## Manba yo'nalishi

Bu katalog Cisco CCNA 200-301 v1.1 topic bloklari bilan mos yurish uchun tuzildi:

- Network Fundamentals
- Network Access
- IP Connectivity
- IP Services
- Security Fundamentals
- Automation and Programmability

Lekin bu fayllar exam dump emas. Ular mavzuni tushunish, lab qilish va troubleshooting fikrlashini shakllantirish uchun yozilgan.
