# REST constraints — 6 ta cheklov

## Muammo: "RESTful" so'zi hamma joyda, lekin ma'nosi noaniq

Ko'p jamoalar API'sini "RESTful" deb ataydi. Lekin so'rasang, kimdir faqat
JSON qaytarganini, kimdir `/getData` endpoint'i borligini nazarda tutadi.
"RESTful" so'zi ma'nosini yo'qotib qo'ygan.

Roy Fielding bu tushunchani mavhum qoldirmadi. U REST'ni **6 ta aniq cheklov
(constraint)** orqali ta'rifladi. Agar API bu cheklovlarga amal qilsa — u
haqiqatan RESTful. Amal qilmasa — shunchaki "JSON qaytaradigan HTTP API".

> Constraint (cheklov) — bu API'ga qo'yiladigan qoida. Har biri erkinlikni
> biroz cheklaydi, lekin evaziga soddalik, kengaytiriluvchanlik yoki
> ishonchlilik beradi. "Nima qila olmasliging" seni yaxshi dizaynga majburlaydi.

## Analogiya: yo'l harakati qoidalari

6 ta constraint'ni **yo'l harakati qoidalari**ga o'xshat.

- "O'ng tomonda yur", "svetoforga bo'ysun" — bular seni cheklaydi.
- Lekin aynan shu cheklovlar tufayli millionlab mashina xaos'siz harakatlanadi.
- Qoidasiz yo'l = tez, lekin halokatli. Qoidali yo'l = biroz sekin, lekin
  bashoratli va xavfsiz.

Xuddi shunday, REST constraint'lari erkinligingni cheklaydi, lekin evaziga
tizim **kengaytiriladigan, ishonchli va tushunarli** bo'ladi.

## Diagramma: 6 constraint bir qarashda

```mermaid
mindmap
  root((REST 6 constraint))
    Client-Server
      rollarni ajratish
      mustaqil rivojlanish
    Stateless
      server eslamaydi
      har so'rov mustaqil
    Cacheable
      javob keshlanadi
      Cache-Control ETag
    Uniform Interface
      resource identifikatsiya
      representation
      self-descriptive
      HATEOAS
    Layered System
      proxy load balancer
      qatlamlar bir-birini bilmaydi
    Code on Demand
      ixtiyoriy
      server kod yuboradi
```

Beshtasi **majburiy**, oxirgisi (**Code on Demand**) — ixtiyoriy. Endi
har birini alohida ko'ramiz.

---

## 1. Client-Server — rollarni ajratish

**G'oya:** client (frontend/mobil) va server (backend) qat'iy ajratilgan.

- **Client** — faqat interfeys bilan shug'ullanadi: ko'rsatish, so'rov yuborish.
- **Server** — ma'lumotni saqlash va qayta ishlash bilan shug'ullanadi.
- Ular orasidagi shartnoma (interface) o'zgarmagan holda qolsa, ikkalasi
  **mustaqil rivojlanadi**.

**Analogiya:** restoran. Mijoz (client) menyudan buyurtma qiladi, oshpaz
(server) tayyorlaydi. Oshxona qanday jihozlangani mijozni qiziqtirmaydi;
mijoz qanday kelgani oshpazni qiziqtirmaydi. Menyu (interface) ularni bog'laydi.

**Foydasi:** bir server — ko'p client (web, iOS, Android) ga xizmat qiladi.
Frontend'ni to'liq qayta yozsang ham, backend'ga tegmasa bo'ladi.

---

## 2. Stateless — server hech narsani eslamaydi

**G'oya:** server client haqida **hech qanday holat (state) saqlamaydi**.
Har bir so'rov o'zida kerakli barcha ma'lumotni olib kelishi shart.

**Analogiya:** har safar do'konga borganingda sotuvchi seni tanimaydi. Kim
ekaningni, nima olmoqchi ekaningni har gal qaytadan aytasan.

```
GET /api/orders
Authorization: Bearer eyJhbGciOiJIUzI1Ni...
Content-Type: application/json
```

E'tibor ber: har so'rovda **token** yuborilyapti. Server "bu foydalanuvchi
oldin login qilgan edi" deb eslamaydi — har so'rovda kimligini token orqali
qaytadan isbotlaysan.

**Notional machine — serverda aslida nima bo'ladi?**
Stateless server so'rovni oladi, tokendan foydalanuvchini aniqlaydi,
javobni tayyorlaydi va **so'rovga oid hamma narsani xotiradan tozalaydi**.
Ikkinchi so'rov kelganda, u birinchisi haqida hech narsa bilmaydi.

**Nega bu muhim?** Agar server holat saqlamasa, so'rovni **istalgan server
nusxasi** bajarishi mumkin. Bu **load balancing**ni osonlashtiradi va
tizimni gorizontal kengaytirishga (scale out) imkon beradi.

```mermaid
sequenceDiagram
    participant C as Client
    participant LB as Load Balancer
    participant S1 as Server 1
    participant S2 as Server 2
    C->>LB: GET /orders (token)
    LB->>S1: so'rovni yo'naltiradi
    S1-->>C: javob
    C->>LB: GET /profile (token)
    LB->>S2: boshqa serverga yo'naltiradi
    Note over S2: Server 2 birinchi so'rovni bilmaydi,<br/>lekin token bor -> baribir ishlaydi
    S2-->>C: javob
```

Diagrammada ko'rinadi: ikkinchi so'rov boshqa serverga tushdi, lekin token
har so'rovda borgani uchun muammo bo'lmadi. Bu **statelessness'ning kuchi**.

---

## 3. Cacheable — javoblarni keshlash

**G'oya:** har bir javob o'zining keshlanishi mumkin yoki mumkin emasligini
aniq bildirishi kerak.

**Analogiya:** kutubxonadan kitob nusxasini bir haftaga olib ketasan. Har
safar bir xil ma'lumot uchun kutubxonaga qaytish shart emas — nusxa qo'lingda.

Agar javob keshlanadigan bo'lsa, client uni ma'lum vaqt qayta ishlatadi va
serverga qayta murojaat qilmaydi. Bu:

- **Client uchun** — tezroq javob
- **Server uchun** — kam yuk
- **Tarmoq uchun** — kam trafik

Keshlash HTTP header'lari orqali boshqariladi:

```
Cache-Control: max-age=3600      # 1 soat kesh
Cache-Control: no-cache          # har safar tekshir
Cache-Control: no-store          # umuman keshlanma
ETag: "33a64df551425fcc55e4d42a148795d9f25f89d4"
```

**ETag** — resurs versiyasining "barmoq izi". Client keyingi so'rovda
`If-None-Match: "33a64..."` yuboradi; agar resurs o'zgarmagan bo'lsa, server
`304 Not Modified` qaytaradi va ma'lumotni qayta yubormaydi — trafik tejaladi.

---

## 4. Uniform Interface — yagona interfeys (eng muhimi)

Bu — REST'ni REST qiladigan **markaziy constraint**. U client va server
o'rtasidagi aloqani bir xil, izchil shaklda tashkil etadi. 4 qismdan iborat:

```mermaid
flowchart LR
    UI["Uniform Interface"] --> A["1. Resource identification<br/>/users/123"]
    UI --> B["2. Manipulation via<br/>representation (JSON)"]
    UI --> C["3. Self-descriptive<br/>messages"]
    UI --> D["4. HATEOAS<br/>hypermedia link'lar"]
```

**1. Resource identification** — har resurs noyob URI bilan aniqlanadi:

```
To'g'ri:    GET /users/123
Noto'g'ri:  GET /getUserById?id=123
```

**2. Manipulation through representations** — client resurs holatini uning
representation'i (masalan JSON) orqali o'zgartiradi. Sen jadval qatorini
to'g'ridan-to'g'ri emas, JSON obyektini yuborib o'zgartirasan.

**3. Self-descriptive messages** — har bir so'rov va javob o'zida yetarli
ma'lumot tutadi: qanday formatda, qanday qayta ishlash kerakligi. Buni
bir xil yondashuv ta'minlaydi:

- Naming convention (nomlash qoidasi)
- Ma'lumot formati (JSON/XML) — `Content-Type` header'ida
- HTTP metodlar va status kodlar

Maqsad: bitta API bilan tanishgan dasturchi boshqasini ham osongina tushunsin.

**4. HATEOAS** (Hypermedia As The Engine Of Application State) — client faqat
**bitta boshlang'ich URI**ni biladi, qolgan amallarni server yuborgan
**link'lar** orqali topadi:

```json
{
  "id": 5,
  "status": "pending",
  "_links": {
    "self": { "href": "/orders/5" },
    "cancel": { "href": "/orders/5/cancel" },
    "pay": { "href": "/orders/5/payment" }
  }
}
```

Client "bu buyurtmani bekor qilish mumkinmi?" degan savolga javobni javobning
o'zidan oladi — `cancel` link'i bormi yoki yo'qmi.

---

## 5. Layered System — qatlamli tizim

**G'oya:** arxitektura bir necha qatlamdan iborat bo'lishi mumkin. Har bir
komponent faqat o'ziga bevosita ulangan qatlam bilan gaplashadi, uning
orqasida nima borligini bilmaydi.

```mermaid
flowchart LR
    C[Client] --> P[Proxy] --> LB[Load Balancer] --> GW[API Gateway] --> S[API Server] --> DB[(Database)]
```

Client "men to'g'ridan-to'g'ri serverga gaplashyapmanmi yoki proxy orqalimi"
degan savolni bilishi shart emas. U shunchaki API Gateway bilan gaplashadi.

**Foydasi:**

- **Xavfsizlik** — firewall, WAF qatlamlarini qo'shish
- **Kengaytirish** — load balancer bilan yukni taqsimlash
- **Kesh qatlami** — CDN yoki reverse proxy'da keshlash

Statelessness bilan birgalikda bu constraint tizimni cheksiz gorizontal
kengaytirishga imkon beradi.

---

## 6. Code on Demand (ixtiyoriy) — server kod yuboradi

Yagona **ixtiyoriy** constraint. Server client'ga bajariladigan **kod**
(script yoki applet) yuborishi mumkin. Client uni bajaradi va o'z
funksionalligini kengaytiradi.

**Klassik misol:** web-sahifa server'dan JavaScript yuklab oladi va brauzerda
ishga tushiradi. Client oldindan hamma logikaga ega bo'lishi shart emas.

Nega ixtiyoriy? Chunki u client'ni murakkablashtiradi va xavfsizlik
xatarlarini oshiradi. Ko'p API'lar undan foydalanmaydi.

---

## Umumiy jadval

| # | Constraint | Bir jumlada | Majburiy? |
| --- | --- | --- | --- |
| 1 | Client-Server | Rollarni ajrat | Ha |
| 2 | Stateless | Server eslamaydi | Ha |
| 3 | Cacheable | Javob keshlanadi | Ha |
| 4 | Uniform Interface | Bir xil interfeys | Ha |
| 5 | Layered System | Qatlamlarga ajrat | Ha |
| 6 | Code on Demand | Server kod yuborishi mumkin | Yo'q |

---

### 🤔 O'ylab ko'r

Agar server har bir foydalanuvchi uchun login holatini **xotirada (session)**
saqlasa, qaysi constraint buziladi? Va bu load balancing'ga qanday ta'sir qiladi?

<details>
<summary>💡 Javobni ko'rish</summary>

**Stateless** constraint buziladi. Agar login holati Server 1 xotirasida
saqlangan bo'lsa, load balancer keyingi so'rovni Server 2'ga yuborsa,
Server 2 foydalanuvchini tanimaydi — "session lost" xatosi chiqadi.

Yechim: holatni serverda emas, **token**da (client tomonida) saqlash. Shunda
istalgan server so'rovni bajarib beradi. Aynan shu sabab REST API'lar odatda
session emas, JWT/token ishlatadi.

</details>

---

## ⚠️ Ko'p uchraydigan xatolar

**1-xato: "JSON qaytarsa REST bo'ldi" deb o'ylash.**
JSON — bu formatgina. Agar `/api` bitta endpoint bo'lsa va hamma narsa POST
bo'lsa, bu Level 0, REST emas. Constraint'lar formatdan muhimroq.

**2-xato: server session'ida foydalanuvchi holatini saqlash.**
Bu stateless'ni buzadi va gorizontal kengaytirishni buzadi. Holatni token
yoki client'da saqla.

**3-xato: HATEOAS majburiy deb o'ylash.**
HATEOAS — Uniform Interface'ning bir qismi bo'lsa-da, amalda ko'p API'lar
uni qo'llamaydi (Level 2). Bu jamoaning ongli tanlovi bo'lishi mumkin.

---

## Xulosa

- REST 6 ta constraint bilan ta'riflanadi; 5 tasi majburiy, biri ixtiyoriy.
- **Client-Server** — rollarni ajratadi, mustaqil rivojlanishga yo'l ochadi.
- **Stateless** — server holat saqlamaydi, har so'rov mustaqil; scale osonlashadi.
- **Cacheable** — javoblar keshlanadi; `Cache-Control`, `ETag` boshqaradi.
- **Uniform Interface** — eng markaziy; resource + representation +
  self-descriptive + HATEOAS.
- **Layered System** — proxy, load balancer, gateway qatlamlariga ajratadi.
- **Code on Demand** — ixtiyoriy; server client'ga kod yuborishi mumkin.

## 🧠 Eslab qol

- 6 constraint bor, 5 majburiy, 1 (Code on Demand) ixtiyoriy.
- Stateless = server eslamaydi = scale oson.
- Uniform Interface — REST'ning yuragi.
- JSON qaytarish REST'ni anglatmaydi.

## ✅ O'z-o'zini tekshir (retrieval practice)

**1. Nega statelessness load balancing'ni osonlashtiradi?**

<details>
<summary>Javob</summary>

Server holat saqlamagani uchun har bir so'rovni **istalgan server nusxasi**
bajara oladi. Load balancer so'rovni bo'sh serverga yuborsa bo'ladi — holat
tokenda kelgani uchun. Bu gorizontal scale (scale out) imkonini beradi.

</details>

**2. ETag va `304 Not Modified` qanday trafik tejaydi?**

<details>
<summary>Javob</summary>

Client `If-None-Match` bilan eski ETag'ni yuboradi. Agar resurs o'zgarmagan
bo'lsa, server `304 Not Modified` qaytaradi va **ma'lumotni qayta yubormaydi** —
faqat "o'zgarmadi" degan qisqa javob keladi. Katta JSON qayta uzatilmaydi.

</details>

**3. Layered System constraint qanday xavfsizlik afzalligini beradi?**

<details>
<summary>Javob</summary>

Client to'g'ridan-to'g'ri serverga emas, oraliq qatlamlar (proxy, gateway,
firewall/WAF) orqali gaplashadi. Bu hujumlarni oldindan filtrlash, DDoS'dan
himoya va serverni tashqi dunyodan yashirish imkonini beradi.

</details>

**4. Uniform Interface'ning 4 qismini ayt.**

<details>
<summary>Javob</summary>

1) Resource identification (URI), 2) Manipulation through representations
(JSON/XML), 3) Self-descriptive messages, 4) HATEOAS (hypermedia link'lar).

</details>

## 🛠 Amaliyot

**1. Oson (Modify).** Yuqoridagi HATEOAS JSON misolini o'zgartir: buyurtma
`status` `"paid"` bo'lsa, `pay` link'ini olib tashla, o'rniga `refund` link'i
qo'sh. Nima uchun bu HATEOAS'ning kuchi?

**2. O'rta (faded example).** Quyidagi javobga to'g'ri kesh header'larini qo'sh:

```
GET /products/42
# TODO: bu javob 10 daqiqa keshlansin        -> ______________
# TODO: resurs versiyasi "abc123" bo'lsin     -> ______________
```

<details>
<summary>Hint</summary>

`Cache-Control: max-age=600` va `ETag: "abc123"`.

</details>

**3. Qiyin (Make).** Bitta oddiy "TODO list" API dizayn qil, unda ikkita
constraint buzilgan bo'lsin (masalan stateless va uniform interface).
Keyin har bir buzilishni tuzat va nima uchun tuzatilganini yoz.

<details>
<summary>Hint</summary>

Buzilish: `GET /getTodos` va serverda session. Tuzatish: `GET /todos` va
har so'rovda `Authorization` token.

</details>

## 🔁 Takrorlash

- Oldingi dars: [REST nima](01-rest-nima.md). Keyingi dars:
  [Resource naming](03-rest-resource-naming.md) — Uniform Interface'ning
  amaliy davomi.
- Takrorlash jadvali: **ertaga** 6 constraint'ni xotiradan yoz (5+1) →
  **3 kundan keyin** stateless sequence diagrammasini qayta chiz →
  **1 haftadan keyin** "TODO API" amaliyotini takrorla.
- **Feynman testi:** "6 constraint" ni yo'l harakati qoidalari analogiyasi
  orqali 3 jumlada tushuntir. Nega cheklov erkinlikdan yaxshiroq?

## 📚 Manbalar

- Roy Fielding dissertatsiyasi, 5-bob (REST) —
  https://ics.uci.edu/~fielding/pubs/dissertation/rest_arch_style.htm
- REST Architectural Constraints — https://restfulapi.net/rest-architectural-constraints/
- Richardson Maturity Model — https://restfulapi.net/richardson-maturity-model/
