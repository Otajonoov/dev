![](../../assets/obsidian-images/Pasted%20image%2020250810022815.png)

Har bir layer yuqorida joylashgan layer'ga **xizmat ko'rsatish modelini** taklif etadi. Bu xizmat ikki yo'l bilan amalga oshiriladi:

- **1**: o'z ichida ma'lum amallarni bajarish orqali
- **2**: pastda joylashgan darajaning xizmatlaridan foydalanish orqali

Masalan, n-darajaning xizmatlari bir tarmoq chetidan ikkinchi chetiga ishonchli message yetkazib berishni o'z ichiga olishi mumkin. Bu (n-1) darajaning ishonchsiz message uzatish xizmatiga n-darajaning yo'qolgan messagelarni aniqlash va qayta yuborish funksiyalarini qo'shish orqali amalga oshiriladi.

##### 1. Physical #Physical_Layer

Data link layer qo'shni network node'lari o'rtasida frame'larni uzatish bilan shug'ullansa, physical layer bu node'lar o'rtasida frame'ning **individual bitlarini** uzatish uchun mo'ljallangan

Physical layer protokollari ishlatiladigan communication linkga va shu linkning haqiqiy transmission muhitiga bog'liq:

- Twisted-pair copper wire
- Single-mode fiber optic
- Boshqa muhitlar

Masalan, Ethernet ko'plab physical layer protokollarini qo'llab-quvvatlaydi:

- Bir protokol twisted copper pair uchun
- Boshqasi coaxial cable uchun
- Uchinchisi fiber optic uchun

Har bir holatda bitlar communication link bo'yicha turli usullarda transmit qilinadi.

##### 2. Data link #DataLink_Layer

Data link layerdagi ma'lumot birligini **==frame==** deb ataymiz.

Protokollar:
- **Ethernet**
- **Wi-Fi**

Network layer ==datagramni== sourcdan destinationgacha routerlar chaini orqali uzatishni ta'minlaydi. Packetni routedagi bir nodedan (host yoki router) keyingisiga o'tkazish uchun network layer data link layer servicelaridan foydalanadi.

Har bir node'da network layer ==datagramni== pastga data link layeriga uzatadi, u keyingi node'gacha yetkazib beradi. Keyin data link layer datagramni yuqoriga network layeriga uzatadi.

Data link layer tomonidan taklif etiladigan servicelar ma'lum communication linkda ishlatiladigan data link layer protocoliga bog'liq:

- Ba'zi data link layer protokollari communication link bo'yicha transmitting node'dan receiving node'gacha reliable delivery'ni ta'minlaydi
- Bu reliability TCP tomonidan taklif etiladigan reliability'dan farq qiladi

##### 3. Network #Network_Layer 

Network layerdagi ma'lumot birligi **==datagram==** deb ataladi.

Network layeri **==datagrammalar==** deb nomlanuvchi data qismlarini bir network hostidan boshqasiga uzatish uchun javobgar.

Transport layer protokollari (TCP va UDP) ==transport== layer ==segmentini== va destination addressini network layeriga uzatadi. Network layer o'z navbatida bu ==segmentni== receiver hostning transport layeriga yetkazib berish serviceni ta'minlaydi.

Network layeri **IP protokolini** o'z ichiga oladi:

- Datagram fieldlarini belgilaydi
- End systemlar va routerlar bu fieldlar bilan qanday amal qilishini belgilaydi
- Butun Internet uchun yagona protokol hisoblanadi

Network layerda **routing protokollari** ham mavjud:

- Source host va destination host o'rtasidagi datagrammalarning routelarini belgilaydi
- Internet - bu network of networks bo'lgani uchun, har bir network o'z routing protocolidan foydalanishi mumkin
##### 4. Transport #Transport_Layer 

Transport layeri application layer "message" larini end applicationlar o'rtasida uzatishni amalga oshiradi. 
Transport layerdagi ma'lumot birligini **==segment==** deb ataymiz.

==TCP== protokoli applicationlarga **connection-oriented** servicelarni taklif qiladi:

- Application layer messagelarning receiverlarga reliable yetkazib berilishi
- Flow control (ya'ni flow tezligini tartibga solish)
- Uzun messagelarni qisqa segmentlarga bo'lish
- Congestion control mexanizmi

==UDP== protokoli applicationlarga **connectionless** servicelarni taklif qiladi:

- Transmissionning reliability kafolatlanmaydi
- Flow control yo'q
- Congestion control yo'q
##### 5. Application #Application_Layer

- **HTTP** - web documentlarni so'rash va uzatishni ta'minlaydi
- **SMTP** - email messagelarini uzatish
- **FTP** - file almashish uchun
- **DNS** - insonlar tushunadigan nomlarni 32-bitli network addresslariga aylantirish
- Application layerdagi ma'lumot birligini **message** deb ataymiz.






#### Layerlar O'rtasidagi Ma'lumot Oqimi 

Protocol layerlarining ishlashi **data encapsulation** printsipiga asoslangan:

1. **Application layer** - message hosil qiladi
2. **Transport layer** - messageni segmentga o'raydi
3. **Network layer** - segmentni datagramga o'raydi
4. **Data link layer** - datagramni framega o'raydi
5. **Physical layer** - frameni bitlar sequencega aylantiradi

Receiver tomonda bu jarayon teskari tartibda amalga oshiriladi.