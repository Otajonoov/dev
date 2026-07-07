# 9. Automation and Programmability

Bu bo'lim tarmoqni avtomatlashtirish va dasturlashga kirish mavzularini CCNA darajasida tushuntiradi. Asosiy g'oya oddiy: tarmoq qurilmalarini faqat CLI orqali qo'lda sozlash o'rniga, controller, API, JSON, Ansible, Terraform va cloud boshqaruv vositalari orqali tezroq, takrorlanadigan va kamroq xatoli boshqarish mumkin.

## Nima uchun bu muhim?

Katta tarmoqda yuzlab switch, router, access point va firewall bo'lishi mumkin. Har bir qurilmaga alohida kirib konfiguratsiya qilish:

- ko'p vaqt oladi;
- inson xatosiga moyil;
- konfiguratsiyalar orasida farq paydo qiladi;
- audit va qayta tiklashni qiyinlashtiradi.

Automation va programmability esa konfiguratsiyani kod, shablon, API chaqiruvi yoki controller siyosati sifatida boshqarishga yordam beradi.

## Bo'limdagi mavzular

1. [SDN va controller-based networking](01-sdn-controller-based-networking.md)
2. [REST API](02-rest-api.md)
3. [JSON](03-json.md)
4. [Ansible va Terraform](04-ansible-terraform.md)
5. [AI, cloud va network management](05-ai-cloud-network-management.md)

## CCNA uchun eslab qolish kerak

- Traditional networking: har bir qurilma ko'pincha alohida CLI orqali boshqariladi.
- Controller-based networking: markaziy controller qurilmalarni siyosat va API orqali boshqaradi.
- Control plane: qaror qabul qiladi, masalan marshrut tanlash.
- Data plane: trafikni real vaqtda uzatadi.
- Management plane: qurilmani boshqarish, monitoring va konfiguratsiya qilish uchun ishlatiladi.
- Northbound API: controller va ilovalar orasidagi API.
- Southbound API: controller va tarmoq qurilmalari orasidagi API.
- REST API: HTTP asosidagi API uslubi.
- JSON: APIlarda juda ko'p ishlatiladigan ma'lumot formati.
- Ansible: ko'pincha konfiguratsiya va task automation uchun ishlatiladi.
- Terraform: infrastructure as code va resurslarni yaratish/o'zgartirish uchun ishlatiladi.

## Kichik misol

Qo'lda ishlash:

```text
Admin -> SSH -> Switch1
Admin -> SSH -> Switch2
Admin -> SSH -> Switch3
```

Avtomatlashtirilgan yondashuv:

```text
Admin -> Script/Ansible/Controller -> Switch1, Switch2, Switch3
```

Natija: bir xil sozlama bir nechta qurilmaga tez va nazoratli tarqatiladi.

## Qisqa Q&A

**Savol:** Automation CLI'ni butunlay almashtiradimi?

**Javob:** Yo'q. CLI hali ham troubleshooting va tekshiruv uchun muhim. Automation esa takroriy va ommaviy ishlarni soddalashtiradi.

**Savol:** CCNA darajasida dasturchi bo'lish shartmi?

**Javob:** Yo'q. Lekin API, JSON, HTTP verbs va automation vositalarining asosiy g'oyasini tushunish kerak.

**Savol:** Controller bo'lsa, router/switch qaror qabul qilmaydimi?

**Javob:** Vaziyatga bog'liq. SDN modelida ayrim qarorlar markaziy controllerga ko'chadi, lekin qurilmalar baribir data plane vazifasini bajaradi.
