# O'qish roadmap

Bu roadmap `new` katalogidagi mavzularni tartibli o'rganish uchun.

## 1-bosqich: Layer 2 ni mustahkamlash

O'qish:

1. `5. network-access/01-vlan.md`
2. `5. network-access/02-trunk-8021q.md`
3. `5. network-access/03-inter-vlan-routing.md`
4. `5. network-access/04-stp-rapid-pvst.md`

Maqsad:

- VLAN nima uchun broadcast domain ajratishini tushunish.
- Access port va trunk port farqini bilish.
- Native VLAN va allowed VLAN xatolarini topish.
- STP loopdan qanday himoya qilishini tushunish.

Lab:

```text
SW1 --- SW2
 |       |
PC1     PC2

VLAN 10: PC1
VLAN 20: PC2
Trunk: SW1-SW2
```

Tekshir:

```text
show vlan brief
show interfaces trunk
show spanning-tree
```

## 2-bosqich: L3 routing asoslari

O'qish:

1. `6. ip-connectivity/01-routing-table.md`
2. `6. ip-connectivity/02-static-routing.md`
3. `6. ip-connectivity/03-ospfv2-single-area.md`

Maqsad:

- Routing table satrlarini o'qiy olish.
- Longest prefix matchni qo'llash.
- Static, default, host va floating static route farqini tushunish.
- OSPF neighbor bo'lmasa sababini topish.

Tekshir:

```text
show ip route
show ip protocols
show ip ospf neighbor
show ip ospf interface brief
```

## 3-bosqich: Redundancy va IPv6

O'qish:

1. `6. ip-connectivity/04-fhrp.md`
2. `6. ip-connectivity/05-ipv6-routing.md`
3. `6. ip-connectivity/06-ipv6-addressing-ndp.md`

Maqsad:

- HSRP/VRRP/GLBP nima uchun kerakligini tushunish.
- IPv6 link-local, global unicast, multicastni farqlash.
- NDP ARP o'rnini qanday bosishini tushunish.
- IPv6 static/default route yozish.

## 4-bosqich: IP Services

O'qish:

1. `7. ip-services/01-dhcp-relay.md`
2. `7. ip-services/02-nat-cisco.md`
3. `7. ip-services/03-ntp.md`
4. `7. ip-services/04-snmp.md`
5. `7. ip-services/05-syslog.md`
6. `7. ip-services/06-qos.md`
7. `7. ip-services/07-ssh-tftp-ftp.md`

Maqsad:

- DHCP broadcast nega routerdan o'tmasligini tushunish.
- `ip helper-address` vazifasini bilish.
- PAT va static NATni sozlash.
- Syslog severity va SNMP monitoringni o'qiy olish.
- QoS nima uchun kerakligini tushunish.

## 5-bosqich: Security

O'qish:

1. `8. security/01-security-concepts.md`
2. `8. security/02-acl.md`
3. `8. security/03-device-access-security.md`
4. `8. security/04-aaa-radius-tacacs.md`
5. `8. security/05-l2-security.md`
6. `8. security/06-wireless-security.md`
7. `8. security/07-ipsec-vpn.md`

Maqsad:

- Threat, vulnerability, exploit, mitigation farqini bilish.
- ACL direction va placementni to'g'ri tanlash.
- Device managementni SSH-only qilish.
- Port security, DHCP snooping, DAI ishlashini tushunish.

## 6-bosqich: Automation

O'qish:

1. `9. automation-programmability/01-sdn-controller-based-networking.md`
2. `9. automation-programmability/02-rest-api.md`
3. `9. automation-programmability/03-json.md`
4. `9. automation-programmability/04-ansible-terraform.md`
5. `9. automation-programmability/05-ai-cloud-network-management.md`

Maqsad:

- Traditional vs controller-based networking farqini bilish.
- Northbound va southbound API rolini tushunish.
- REST, CRUD va HTTP verbsni bilish.
- JSONni o'qiy olish.
- Ansible/Terraform nima uchun ishlatilishini tushunish.

## Har bo'limdan keyingi test

O'zingga 3 savol ber:

1. Bu mavzuni 5 yoshli bolaga qanday tushuntiraman?
2. Cisco qurilmada buni qaysi command bilan tekshiraman?
3. Bu mavzuda eng ko'p uchraydigan xato nimadan iborat?
