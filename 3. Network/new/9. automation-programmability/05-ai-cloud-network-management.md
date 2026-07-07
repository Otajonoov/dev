# AI, Cloud va Network Management

Zamonaviy tarmoqlarda monitoring, automation, cloud boshqaruv va AI/ML asosidagi tahlil tobora ko'proq ishlatiladi. CCNA darajasida maqsad - bu texnologiyalar nima qilishini, qayerda foydali ekanini va qanday ehtiyot choralar kerakligini tushunish.

## Network management nima?

Network management - tarmoqni kuzatish, sozlash, muammolarni aniqlash va hisobot qilish jarayoni.

Asosiy vazifalar:

- qurilmalar holatini ko'rish;
- interface up/down holatini kuzatish;
- bandwidth va packet loss monitoring;
- konfiguratsiya backup;
- log va alertlarni tahlil qilish;
- security va compliance tekshirish.

## Traditional management

An'anaviy yondashuvda administrator ko'pincha alohida vositalardan foydalanadi:

```text
SSH/CLI     -> konfiguratsiya va troubleshooting
SNMP        -> monitoring
Syslog      -> log yig'ish
NetFlow     -> trafik tahlili
Spreadsheet -> inventar ro'yxati
```

Bu ishlaydi, lekin katta tarmoqda ma'lumotlar bo'linib ketadi.

## Cloud-based network management

Cloud-based managementda boshqaruv platformasi cloud orqali ishlaydi. Qurilmalar internet yoki private ulanish orqali cloud controllerga ulanadi.

```text
Admin browser
     |
     v
Cloud management platform
     |
     v
Routers / Switches / APs / Firewalls
```

Afzalliklari:

- markaziy dashboard;
- remote boshqaruv;
- tez deployment;
- avtomatik update va monitoring;
- ko'p filiallarni bitta joydan boshqarish.

Ehtiyot bo'lish kerak:

- internet bog'lanishiga qaramlik;
- identity va access control;
- data privacy;
- change control;
- vendor lock-in.

## On-premises va cloud management farqi

| Savol | On-premises management | Cloud management |
|---|---|---|
| Platforma qayerda? | kompaniya data centerida | cloud provider yoki vendor cloudida |
| Yangilash | ko'pincha admin bajaradi | ko'pincha vendor boshqaradi |
| Remote filiallar | VPN yoki maxsus ulanish kerak bo'lishi mumkin | internet orqali osonroq |
| Nazorat | ko'proq local nazorat | servis modeliga bog'liq |
| Skalalash | resurs talab qiladi | odatda tezroq |

## AI va ML asoslari

AI (Artificial Intelligence) - tizimning "aqlli" vazifalarni bajarishga yordam beruvchi texnologiyalar umumiy nomi.

ML (Machine Learning) - AI ichidagi yo'nalish bo'lib, tizim ma'lumotlardan pattern o'rganadi.

Network managementda ML quyidagilar uchun ishlatilishi mumkin:

- anomal trafikni aniqlash;
- odatiy bandwidth patternlarini o'rganish;
- qurilma nosozligini oldindan taxmin qilish;
- alertlarni ustuvorlik bo'yicha saralash;
- root cause analysis uchun tavsiya berish.

## Generative AI nima?

Generative AI matn, kod, rasm yoki boshqa kontent yaratishi mumkin bo'lgan AI turidir. Network sohasida u yordamchi sifatida ishlatilishi mumkin.

Misollar:

- log xabarlarini sodda tilda tushuntirish;
- konfiguratsiya shabloni yaratishda yordam berish;
- troubleshooting checklist taklif qilish;
- API yoki Ansible playbook namunasi yozish;
- hujjat tayyorlash.

Muhim: generative AI javobi har doim tekshirilishi kerak. U noto'g'ri buyruq, mos kelmaydigan vendor syntax yoki xavfli o'zgarish taklif qilishi mumkin.

## AIOps tushunchasi

AIOps - IT operations jarayonlarida AI/MLdan foydalanish. Network monitoringda AIOps ko'p alertlarni tahlil qilib, muhimlarini ajratishga yordam beradi.

Oddiy misol:

```text
100 ta alert:
- 80 tasi bitta uplink muammosidan kelgan
- 15 tasi dependency sababli paydo bo'lgan
- 5 tasi alohida tekshirilishi kerak

AIOps: asosiy ehtimoliy sabab - Core-SW1 uplink packet loss
```

## Telemetry

Telemetry - qurilmalardan holat va performance ma'lumotlarini yig'ish jarayoni. Traditional SNMP pollingdan farqli ravishda streaming telemetry real vaqtga yaqinroq ma'lumot berishi mumkin.

Misollar:

- CPU va memory;
- interface counters;
- packet drop;
- latency;
- routing neighbor holati;
- wireless client tajribasi.

```text
Network devices -> telemetry stream -> collector/controller -> dashboard/AI analysis
```

## Closed-loop automation

Closed-loop automationda tizim muammoni aniqlaydi, qaror qiladi va avtomatik tuzatish amalini bajaradi.

Soddalashtirilgan oqim:

```text
Detect -> Analyze -> Decide -> Act -> Verify
```

Misol:

1. Controller branch routerda packet loss ko'payganini aniqlaydi.
2. Telemetry va loglarni tahlil qiladi.
3. Backup link yaxshi ekanini ko'radi.
4. Traffic policy'ni backup link tomonga o'zgartiradi.
5. Natijani tekshiradi.

CCNA darajasida eslash kerak: closed-loop automation kuchli, lekin noto'g'ri siyosat bo'lsa, muammoni kattalashtirishi mumkin. Shuning uchun approval, rollback va audit muhim.

## Security va governance

Automation, cloud va AI ishlatilganda xavfsizlik yanada muhim bo'ladi.

Yaxshi amaliyotlar:

- least privilege: foydalanuvchiga faqat kerakli ruxsat berish;
- MFA ishlatish;
- API tokenlarni maxfiy saqlash;
- audit loglarni yoqish;
- change approval jarayonini saqlash;
- rollback rejasini tayyorlash;
- konfiguratsiyani backup qilish.

## Amaliy ssenariy

Muammo: filialdagi foydalanuvchilar Wi-Fi sekinligidan shikoyat qilyapti.

Traditional tekshiruv:

1. AP holatini CLI yoki controllerdan tekshirish.
2. Interferensiya, signal, client soni va uplinkni ko'rish.
3. Loglarni o'qish.
4. Qo'lda xulosa chiqarish.

AI/cloud yordamidagi tekshiruv:

1. Cloud dashboard client experience score ko'rsatadi.
2. ML odatiy holatdan chetga chiqishni aniqlaydi.
3. Tizim "2.4 GHz bandda interference yuqori" degan tavsiya beradi.
4. Admin tavsiyani tekshiradi va kanal/power siyosatini o'zgartiradi.

## Common mistakes

- **AI hamma muammoni avtomatik hal qiladi deb o'ylash.** AI yordamchi, yakuniy javobgar odatda administrator.
- **Cloud dashboard borligi backup shart emas degani emas.** Backup va export hali ham kerak.
- **Alert soniga qarab muammo og'irligini baholash.** Bitta asosiy muammo yuzlab alert yaratishi mumkin.
- **Telemetryni faqat monitoring deb ko'rish.** Telemetry automation va analytics uchun ham asos bo'ladi.
- **API tokenlarni oddiy faylda saqlash.** Maxfiylik va access control talab qilinadi.

## Qisqa Q&A

**Savol:** Cloud management bo'lsa, qurilmalar internetsiz ishlamay qoladimi?

**Javob:** Ko'p qurilmalar local forwardingni davom ettiradi, lekin markaziy monitoring, yangi policy va remote boshqaruv ta'sirlanishi mumkin.

**Savol:** AI tavsiyasini darhol productionda bajarish kerakmi?

**Javob:** Yo'q. Tavsiyani tekshirish, riskni baholash va kerak bo'lsa approval olish kerak.

**Savol:** Telemetry SNMPni almashtiradimi?

**Javob:** Har doim emas. Ko'p tarmoqlarda SNMP va telemetry birga ishlatiladi.
