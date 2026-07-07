### Soket Nima? 
Soket PROTOKOL EMAS, balki INTERFEYS!

**Soket** - bu tarmoq orqali ma'lumot almashish uchun dasturiy interfeys. Uni **eshik** metaforasi bilan tushuntirish mumkin:


### Soket Qanday Ishlaydi?

**1. Soket yaratilishi:**

```
┌─────────────────────────┐
│   APPLICATION LAYER     │ ← Sizning ilovangiz (HTTP, FTP va h.k.)
│   (Ilova qatlami)       │
├─────────────────────────┤
│      🔌 SOKET API       │ ← Bu yerda soket! (Interfeys)
├─────────────────────────┤
│   TRANSPORT LAYER       │ ← TCP/UDP protokollari
│   (Transport qatlami)   │
├─────────────────────────┤
│   NETWORK LAYER         │ ← IP protokoli
│   (Tarmoq qatlami)      │
├─────────────────────────┤
│   DATA LINK LAYER       │
├─────────────────────────┤
│   PHYSICAL LAYER        │
└─────────────────────────┘
```


**2. Ma'lumot uzatish jarayoni:**

**Mijoz tomonida:**

1. Ilova xabarni soketga yozadi
2. Soket xabarni TCP ga topshiradi
3. TCP xabarni paketlarga bo'lib, IP ga beradi
4. IP paketlarni tarmoq orqali yuboradi

**Server tomonida:**

1. Tarmoqdan paketlar keladi
2. IP paketlarni TCP ga beradi
3. TCP paketlarni yig'ib, asl xabarni tiklaydi
4. Soket orqali ilova xabarni oladi

### Soket Nimalar Ustiga Qurilgan?

**1. Operating System (OS) darajasida:**

- Soket OS tomonidan boshqariladigan resurs
- File descriptor (fayl tavsiflovchi) kabi ishlaydi
- OS soketlar uchun bufer ajratadi


### Soket Qanday Ishlaydi - Batafsil:

**1. Ilova soket yaratadi:**

python

```python
sock = socket.socket(AF_INET, SOCK_STREAM)
# Bu OS dan "menga TCP ulanish kerak" deyapti
```

**2. OS nima qiladi:**

```
Ilova: "Soket yarat"
    ↓
OS: "OK, mana file descriptor: 5"
    ↓
OS ichida: - TCP uchun bufer ajratadi
          - Port tayinlaydi
          - IP stack bilan bog'laydi
```

**3. Ma'lumot yuborilganda:**

```
Ilova: sock.send("Salom")
    ↓
Soket API: OS ga system call
    ↓
OS: TCP protokoliga topshiradi
    ↓
TCP: Segmentlarga bo'ladi, header qo'shadi
    ↓
IP: Paketlarga o'raydi
    ↓
Tarmoq: Yuboradi
```


### Soket (Socket)

Soket - bu jarayon va tarmoq o'rtasidagi dasturiy interfeys. Jarayonlar xabarlarni soket orqali yuboradi va qabul qiladi. Soketni uyning eshigiga o'xshatish mumkin - xabar shu "eshik" orqali chiqib ketadi va boshqa xostdagi "eshik"ka yetib boradi.

Механизм cookie,
определенный в документе RFC 6265562, позволяет веб-сайтам отслежи-
вать состояние пользовательского соединения.

**Cookie** - foydalanuvchi holatini kuzatish mexanizmi, 4 ta komponentdan iborat:

- Server javobidagi Set-cookie sarlavhasi
- Mijoz so'rovidagi Cookie sarlavhasi
- Brauzerdagi cookie-fayl
- Serverdagi ma'lumotlar bazasi


**Cookie vazifasi** - foydalanuvchini identifikatsiya qilish va sessiya holatini saqlash
**Proksi-server (Web cache)** - tez-tez so'raladigan ob'ektlarning lokal nusxalarini saqlaydi
**CDN (Content Distribution Networks)** - geografik taqsimlangan proksi-serverlar tarmog'**i**

![](../../../assets/obsidian-images/Pasted%20image%2020250819111638.png)