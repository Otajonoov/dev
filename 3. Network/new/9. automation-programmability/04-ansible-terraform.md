# Ansible va Terraform

Ansible va Terraform automation dunyosida ko'p uchraydigan ikki vosita. Ikkalasi ham infratuzilmani avtomatlashtirishga yordam beradi, lekin ishlash falsafasi va kuchli tomonlari farq qiladi.

## Ansible nima?

Ansible - konfiguratsiya boshqaruvi va task automation vositasi. U playbook deb ataladigan YAML fayllar orqali vazifalarni bajaradi.

Network misollar:

- switchda VLAN yaratish;
- router konfiguratsiyasini backup qilish;
- interface description o'zgartirish;
- bir nechta qurilmada bir xil buyruqni bajarish;
- konfiguratsiyani shablon asosida tarqatish.

Ansible odatda agentless ishlaydi, ya'ni boshqarilayotgan qurilmaga alohida agent o'rnatish talab qilinmaydi. Network qurilmalarda SSH, NETCONF, RESTCONF yoki vendor modullari ishlatilishi mumkin.

## Ansible asosiy tushunchalari

| Tushuncha | Izoh |
|---|---|
| Inventory | qaysi qurilmalar boshqarilishini ko'rsatadi |
| Playbook | bajariladigan vazifalar ro'yxati |
| Task | bitta amal |
| Module | Ansible bajaradigan tayyor funksiya |
| Variable | qayta ishlatiladigan qiymat |
| Template | Jinja2 yordamida konfiguratsiya matni yaratish |

## Inventory misoli

```ini
[switches]
sw1 ansible_host=192.168.10.11
sw2 ansible_host=192.168.10.12

[routers]
r1 ansible_host=192.168.10.1
```

Bu inventoryda `switches` va `routers` guruhlari bor.

## Playbook mini misoli

Quyidagi misol tushuncha berish uchun soddalashtirilgan:

```yaml
---
- name: Show version on switches
  hosts: switches
  gather_facts: no

  tasks:
    - name: Run show version
      ansible.netcommon.cli_command:
        command: show version
      register: output

    - name: Print result
      debug:
        var: output.stdout
```

Bu playbook `switches` guruhidagi qurilmalarda `show version` buyrug'ini bajaradi.

## Ansible qachon foydali?

- bir xil konfiguratsiyani ko'p qurilmaga berish kerak bo'lsa;
- backup va audit kerak bo'lsa;
- qo'lda bajariladigan tasklar takrorlansa;
- konfiguratsiya shablonlarini ishlatish kerak bo'lsa.

## Terraform nima?

Terraform - infrastructure as code (IaC) vositasi. U infratuzilma resurslarini deklarativ konfiguratsiya orqali yaratish, o'zgartirish va o'chirishga yordam beradi.

Network va cloud misollar:

- cloud VPC/VNet yaratish;
- subnetlar yaratish;
- security group/firewall rule sozlash;
- load balancer resurslarini yaratish;
- SD-WAN yoki controller resurslarini provider orqali boshqarish.

## Deklarativ yondashuv

Terraformda siz "qanday bajarish kerak"dan ko'ra "qanday holat bo'lishi kerak"ni yozasiz.

Masalan:

```hcl
resource "example_vlan" "users" {
  id   = 10
  name = "USERS"
}
```

Ma'nosi: `USERS` nomli VLAN 10 mavjud bo'lishi kerak.

Terraform provider real platforma bilan gaplashadi. Provider AWS, Azure, Google Cloud, VMware, network controller yoki boshqa tizim uchun bo'lishi mumkin.

## Terraform asosiy tushunchalari

| Tushuncha | Izoh |
|---|---|
| Provider | qaysi platforma bilan ishlashni belgilaydi |
| Resource | yaratiladigan yoki boshqariladigan obyekt |
| Data source | mavjud obyekt haqida ma'lumot o'qish |
| State | Terraform biladigan joriy holat |
| Plan | bajarilishidan oldin o'zgarishlar ro'yxati |
| Apply | rejalangan o'zgarishlarni bajarish |
| Destroy | resurslarni o'chirish |

## Terraform workflow

```text
terraform init
terraform plan
terraform apply
```

Izoh:

- `init` provider va loyihani tayyorlaydi;
- `plan` nima o'zgarishini ko'rsatadi;
- `apply` o'zgarishlarni bajaradi.

## Cloud subnet mini misoli

Quyidagi kod vendor-neutral tushuncha berish uchun soddalashtirilgan:

```hcl
resource "cloud_network" "campus_lab" {
  name = "campus-lab"
}

resource "cloud_subnet" "users" {
  name       = "users-subnet"
  network_id = cloud_network.campus_lab.id
  cidr       = "10.10.10.0/24"
}
```

Bu yerda maqsad: network va unga tegishli subnet bo'lishi kerak.

## Ansible va Terraform farqi

| Savol | Ansible | Terraform |
|---|---|---|
| Asosiy vazifa | konfiguratsiya va task automation | resurslarni IaC orqali boshqarish |
| Til/uslub | YAML playbook | HCL konfiguratsiya |
| Yondashuv | ko'proq procedural + declarative modullar | deklarativ |
| State | odatda markaziy state shart emas | state juda muhim |
| Kuchli tomoni | qurilma konfiguratsiyasi, tasklar | cloud/infrastructure resurslari |

Oddiy eslab qolish:

```text
Ansible  = configure and automate tasks
Terraform = provision and manage infrastructure resources
```

## Idempotency nima uchun muhim?

Automationda idempotency - bir xil automation bir necha marta ishga tushsa ham natija keraksiz o'zgarmasligi.

Masalan, VLAN 10 allaqachon bor bo'lsa:

- yaxshi automation: "VLAN 10 bor, o'zgartirish shart emas";
- yomon automation: har safar yangi xato yoki duplicate yaratadi.

## Common mistakes

- **Ansible va Terraformni bir xil deb bilish.** Ular bir-birini to'ldiradi, lekin vazifasi bir xil emas.
- **Terraform state faylini e'tiborsiz qoldirish.** State noto'g'ri boshqarilsa, real infratuzilma bilan kod orasida farq paydo bo'ladi.
- **Playbookni test qilmasdan productionga qo'llash.** Avval lab yoki kichik guruhda sinash kerak.
- **Credentiallarni kod ichida ochiq yozish.** Token, parol va API key maxfiy saqlanishi kerak.
- **Manual o'zgarishlar va IaCni aralashtirish.** Koddan tashqarida o'zgartirish drift keltirib chiqaradi.

## Qisqa Q&A

**Savol:** Ansible ishlatish uchun Python bilish shartmi?

**Javob:** Yo'q, asosiy playbooklar YAML bilan yoziladi. Lekin Python tushunchasi murakkab modullar va integratsiyalar uchun foydali.

**Savol:** Terraform router CLI buyruqlarini bajaradimi?

**Javob:** Odatda yo'q. Terraform resurs holatini boshqarishga mo'ljallangan. CLI tasklar uchun Ansible ko'proq mos.

**Savol:** Ikkalasini bir loyihada ishlatish mumkinmi?

**Javob:** Ha. Masalan, Terraform cloud subnetlarni yaratadi, Ansible esa qurilmalar konfiguratsiyasini sozlaydi.
