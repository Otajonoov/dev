
Я очень хочу помочь читателям, насколько это возможно, проявить свои лучшие качества в моделировании программного обеспечения с помощью самых эффективных из доступных средств предметно-ориентированного проектирования, или DDD (domain-driven design). Этот набор инструментов, представляющих собой совокупность шаблонов, впервые был проанализирован Эриком Эвансом (Eric Evans) в книге Domain-Driven Design: Tackling Complexity т the heart of Software’. Я бы хотел, чтобы принципы DDD освоил каждый разработчик. Если это значит, что я хочу внедрить принципы DDD в массы. пусть так и будет. DDD заслуживает этого. DDD — это инструментарий, которые люди имеют право использовать для создания самых сложных моделей программного обеспечения. Я написал эту книгу, чтобы сделать изучение и использование DDD максимально простым и доступным для самой широкой аудитории. Для людей. воспринимающих мир с помощью слуха, DDD открывает возможность обучения на основе общения в группе разработчиков модели, создающих ЕДИНЫЙ ЯЗЫК (UBIQUITOUS LANGUAGE) . Люди, воспринимающие реальность с помощью зрения и осязания, оценят визуальный и тактильный характер процесса использования инструментальных средств DDD для стратегического и тактического проектирования. Это особенно ярко проявляется при создании КАРТ KOHTEKCTOB (CONTEXT MAPS) и моделировании бизнес-процессов с помощью СОБЫТИЙНОГО ШТУРМА (EVENT STORMING). Таким образом, я полагаю, что DDD может удовлетворить потребности каждого, кто хочет учиться и достичь успеха с помощью разработки моделей.


Douglas Martin (Book Design kitobidan):
> "Dizayn kerakmi yoki yo'qmi degan savol ma'nosiz: **dizayn muqarrar**. Yaxshi dizaynga alternativa — yomon dizayn, dizaynsiz holat emas."

**Xulosa:** Biz modellashtirish bilan shug'ullanamiz — tan olsak ham, tan olmasak ham.

## Dizayn qimmat emasmi?

Agar siz puxta loyihalangan dasturiy ta'minot yaratish qimmat deb o'ylayotgan bo'lsangiz, **yomon dizayndan foydalanish yoki uni tuzatish qancha qimmatroq bo'lishini o'ylab ko'ring**.

DDD — bu **Ubiquitous Language**ni aniq belgilangan **Bounded Context** doirasida modellashtirish demakdir.



Несмотря на то что обсуждение архитектуры связано с технологией, мо-
дель предметной области должна быть свободной от технологии. С одной
стороны, именно поэтому транзакции управляются службами приложения,
а не моделью предметной области.

**Muhim tamoyil:** Domain model texnologiyadan xoli bo'lishi kerak.
**Misol:** Transaksiyalar Application Service tomonidan boshqariladi, Domain Model tomonidan emas.

### DDD bilan ishlaydigan arxitektura uslublari

**Ports and Adapters** — asosiy, lekin yagona emas. Quyidagilar ham qo'llaniladi:

1. **Event-Driven Architecture** va **Event Sourcing** (6-bobda)
2. **CQRS** (Command Query Responsibility Segregation)
3. **Reactive** va **Actor Model**
4. **REST** (Representational State Transfer)
5. **SOA** (Service-Oriented Architecture)
6. **Microservices** (Mikroservislar)
7. **Cloud Computing** (Bulutli hisoblash)

### Mikroservislar va Bounded Context

**Muhim farq:**

**1-fikr:** Mikroservis = Bounded Context

- Bu to'g'ri yondashuv
- Har bir Bounded Context alohida mikroservis bo'lishi mumkin

**2-fikr (ba'zilar):** Mikroservis Bounded Context'dan ancha kichik

- Faqat bitta kontseptsiyani modellaydi
- Misol: `Product` yoki `BacklogItem` — alohida mikroservis

**Haqiqat:** Agar `Product` va `BacklogItem` bir xil Bounded Context'da bo'lsa (masalan, Scrum konteksti), ular:

- **Lingvistik jihatdan** bir xil kontekst va semantik chegaralar ichida
- Faqat turli deployment modullari
- Context Map orqali o'zaro ta'sir ko'rsatishi mumkin
