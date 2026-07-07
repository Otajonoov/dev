![](../../assets/obsidian-images/Pasted%20image%2020250826140315.png)|![](../../assets/obsidian-images/Pasted%20image%20.png)
**OSI**
1. **Application layer**
2. **Presentation layer**
3. **Session layer**
4. **Transport layer**
5. **Network layer**
6. **Data Link layer**
7. **Physical layer**

**OSI modelidagi qo'shimcha qatlamlar:**

**Presentation layer:**
- Ma'lumotlarni siqish (compression)
- Shifrlash (encryption)
- Ma'lumotlar formatini tavsiflash
- Turli kompyuterlar o'rtasidagi format farqlarini hal qilish

**Session layer:**
- Ma'lumot almashishda chegaralash va sinxronizatsiya
- Seansni nazorat qilish
- Aloqa uzilganda tiklash imkoniyatlari

```
Application layer: Xabar (M)
↓
Transport layer: Header (Ht) + M = Segment
↓  
Network layer: Header (Hn) + Segment = Datagram
↓
Data Link layer: Header (Hl) + Datagram = Frame
```

Routerlar ham, data link switchlar ham packet switching bilan shug'ullanadi.
Oxirgi tizimlarga o'xshab, ushbu qurilmalarning apparat va dasturiy ta'minoti ham qatlamlar ko'rinishida tashkil etilgan. Ammo oxirgi tizimlardan farqli o'laroq, routerlar va switchlar protokol steklarining barcha qatlamlarini ishlatmaydi. Ular odatda pastki qatlamlarda ishlaydi.

Data link switch 1 va 2-qatlamlarni ishlatadi
Routerlarda esa 1 dan 3-qatlamgacha bo'lgan qatlamlar amalga oshirilgan. 
Bu shuni anglatadiki, masalan, Internetdagi routerlar IP protokoli (3-qatlam protokoli) bilan ishlay oladi, data link switchlar esa buni qila olmaydi. 
Ular IP-manzillarni taniy olmaydi, lekin Ethernet manzillari kabi 2-qatlam manzillari bilan ishlaydi. 

1. ==Application layer== xabari (M) ==transport== layerga uzatiladi. 
2. ==Transport== layer xabarni qabul qiladi va unga qo'shimcha ma'lumot ==(Header (Ht) + M = Segment)== qo'shadi, bu ma'lumot keyinchalik qabul qiluvchi tarafdagi transport layer tomonidan ishlatiladi. Application layer xabari transport layer header ma'lumoti bilan birgalikda transport layer ==segmentini== tashkil etadi. So'ngra transport layer segmentni network layerga uzatadi, 
3. ==Network layer== o'z navbatida o'z qatlami header ma'lumotini (==Header (Hn) + Segment = Datagram==) qo'shadi, masalan, manba va qabul qiluvchi oxirgi tizimlarning manzillarini, va shu tariqa network layer datagramini yaratadi. Keyin datagram ==data link== layerga uzatiladi
4. ==Data link layer== ham o'zining headerini ==(Header (Hl) + Datagram = Frame)== qo'shadi va data link layer frameni yaratadi.

Shunday qilib, biz ko'ramizki, har bir qatlamdagi packet ikki turdagi maydonni o'z ichiga oladi - header maydoni va data maydoni, bu maydonlar odatda undan yuqorida joylashgan qatlamdan kelgan packetni o'z ichiga oladi.

---
Qani buni so`rov berish misolida ko`ramiz , masalan , siz **_google.com_** deb yozdingiz va enter tugmasini bosdingiz .Bu **GET** requestga o`tkaziladi va [_https://google.com_](https://google.com/) va bir nech ta headerlar bilan **_application layer_**da paydo bo`ladi .Application layerda bu bizga kelgan ma’lumotlarni Nodejs orqali qayta ishlash joyi (bizning http dasturlar) .Xo`sh davom etamiz agarda siz **https** protocoldan foydalanayotgan bo`lsangiz _encryption_ **layer** 6 (**_presentation layer_**)da bo`ladi agarda ishlatilmagan bo`lsa layer 5 (**session**) ga shunday berib yuboriladi . Biz **http** protocol ishlatmoqdamiz va bu **_server o`rtasida bog`lanishga asoslangan_** shuning uchun buyerda tepadan kelgan ma’lumotga _tag sifatida session id qo`shiladi_ ( bu bizga bir necha so`rovlar kelganda bu qaysi sessionga egaligini bilib olish uchun kerak) va layer 4 **_Transport layer_**ga ma’lumotlar uzatiladi . U malumotlarni shunchaki bir to`da bitlar deb hisoblaydi.Agarda u malum hajmdan ko`p bo`lsa buni qismlarga bo`ladi va uni **_segment_** deb ataydi va har bir segmentga **source port** va **destination port** va yana bir necha narsalar qo`shadi .Bu yerda yana **sequence number (tartib raqam)** ham bo`ladi buni yordamida qabul qilgan server ma’lumotlar tartibini bilib oladi , yo`qolgan segmentlarni qayta so`ray olishi mumkin .Bu malumotlarni hammasini **_network layer_**ga yuboradi (u portlar haqida hech narsani bilmaydi ) bu qismda segmentlar bir necha qismga bo`linadi va ular **_packetlar_** deyiladi network layer har bir packetga **source ip adress** and **destination ip address** qo`shib hech qanday xatoliklarga tekshirmagan holda **data link layer**ga ma’lumotlarni o`tkazib yuboradi. Bu layer ham uni bir necha qismlarga bo`ladi va ularni **frame** deb ataydi va u har bir framega **_source_** **_MAC address_** va **_target MAC address_** qo`shadi(_Agarda target uchun mac addressni bilmasa local networklar uchun ARP protocol orqali aniqlab oladi tashqaridagi komputerlar uchun oldingi maqolamda aytganimdek routerning mac addressini beradi va router uni forward qilib yuboradi_ ).Bu yerda biz mac addresslarnin bilishimiz kerak.Bundan so`ng u haqiqiy hayotga yani **_physical layer_**ga yani kabellarga 0 va 1 lar bo`lib o`tadi . Bizda malumotlar bor ammo uni qayerga uzatishni bilmaymiz chunki u shunchaki elektr toki va unda yo`nalish mavjud emas. Bu mahalliy tarmoqda juda xavfli bo`lishi mumkin chunki bu ma’lumot tarmoqqa ulangan hammaga uzatiladi va u encryption qilinmagan bo`lsa hamma malumotlarinigiz kundek ravshan bo`ladi (Shuning uchun ovqatlanish va jamoat joylaridagi wifilarga ulanishda ehtiyot bo`ling) . Bu har xil ishlaydi agarda wifi shunchaki switch sifatida ishlatilingan bo`lsa u hamma ulanganlarga uzatadi agarda u internetga ulangan bo`lsa(_so`rov tashqariga uzatilgan bo`lsa_) aqlli holda faqatgina kerakli manzilga uzatilinadi. Local tarmoqda kompyuterlar **_frame_**larni qabul qilib , **_mac address_** orqali tekshirishadi agarda u _bunga tegishli bo`lmasa shunchaki shu yo`nalish kelayotgan packetlarni etibor bermaydi_ . Ma’lumotni uzatishda davom etamiz. U siz xohlagan routergacha yetib bordi va **_data link layer_**ga ko`tariladi tarmoqda **_mac address_** orqali kerakli manzil aniqlanib uzatiladi. Agarda bu so`rov mac address orqali bitmasa kelayotgan so`rov bo`lsa **_network layer_**gacha ko`tariladi agarda **NAT** ishlatilingan bo`lsa **_transport layer_**gacha chiqib portni ham aniqlaydi va map qilib kerakli kompyuterga uzatadi .Va u kerakli kompyuterda layerlar orqali tepaga ko`tarilib boradi 0 va 1 lar **frame**ga aylanadi **framelar packet**ga , va **packetlar segment**ga aylanib oxirida layer 7 da bizga ko`rinadigan ma’lumotlarga aylanadi .**_Bu yerda portlar transport layer va application layerni bog`lab beradi .


---














![](../../assets/obsidian-images/Pasted%20image%2020250827113038.png)|![](../../assets/obsidian-images/Pasted%20image%20.png)