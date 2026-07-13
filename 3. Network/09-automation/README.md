# 09 — Network Automation va Programmability

Bu modul tarmoqni **qo'lda CLI** orqali bittalab sozlashdan **kod, API va
controller** orqali ommaviy boshqarishga o'tishni tushuntiradi. Asosiy g'oya
oddiy: yuzlab qurilmani qo'lda boshqarish sekin, xatoli va nazoratsiz —
automation esa tez, takrorlanadigan va auditga qulay.

Modul CCNA darajasidagi automation va programmability mavzularini oladi va
ustiga pedagogika qatlami (analogiya, worked example, retrieval practice)
qo'shadi. Har darsda Mermaid diagrammalar, `curl`/RESTCONF/playbook namunalari
va zamonaviy (2025-2026) holat aks etadi.

## Nima o'rganiladi

- SDN va controller-based networking: control/data/management plane,
  northbound/southbound API, underlay/overlay/fabric.
- REST API va RESTCONF: CRUD, status kodlar, autentifikatsiya, idempotency.
- JSON, YAML va XML: data serialization network automation'da.
- Ansible va Terraform: configuration management, IaC, state, idempotency.
- AI, cloud va network management: AIOps, telemetry, closed-loop, AgenticOps.

## Darslar

1. [SDN va controller-based networking](01-sdn-controller-based.md)
2. [REST API va network automation](02-rest-api-va-network-automation.md)
3. [JSON va YAML](03-json-yaml.md)
4. [Ansible va Terraform](04-ansible-terraform.md)
5. [AI, cloud va network management](05-ai-va-cloud-network-management.md)

## O'qish tartibi

Darslar ketma-ket qurilgan — 1-darsdan boshlang. Controller (1) skript bilan
API orqali (2) gaplashadi, API body JSON/YAML (3) formatida bo'ladi, bu formatlar
Ansible/Terraform (4) da ishlatiladi, hammasining ustida esa AI/cloud boshqaruv
(5) turadi. Har darsni tugatgach `## ✅ O'z-o'zini tekshir` savollariga javob
bering va `## 🔁 Takrorlash` jadvaliga amal qiling.
