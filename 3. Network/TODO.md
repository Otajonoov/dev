1. HTTP/HTTPS
2. WebSocket
3. SSH
4. FTP/SFTP
5. SMPP
6. DNS
7. SSL (Secure Sockets Layer)

DNS, IP, TCP, UDP, HTTP/HTTPS
Reliability patterns
 пакет `net` in Go
 `Практика`: клиент и сервер на основе TCP и UDP
  `Практика`: Реализуйте чат на основе WebSocket с использованием [gorilla/websocket](https://github.com/gorilla/websocket)
  Теория: что такое балансировка нагрузки, подходы (round-robin, least connections).
  Реализуйте простую балансировку нагрузки на уровне приложения.
  Используйте HAProxy или NGINX для внешней балансировки
   Разбор популярных механизмов: Basic Auth, API Keys, Token-based Auth.
   1. net/http …
        2. Go-Auth для Basic и Digest Auth
        3. OAuth2 для работы с OAuth.
        4. **JWT**
        5. `Практика`
        - сервер с Basic Auth.
        - Используйте OAuth2 для авторизации через Google/Facebook API.
- `gRPC`
        
        1. Protocol Buffers (protobuf)
        2. gRPC-сервер и клиент
        3. Interceptors
        4. Реализуйте сервис с каждым из типов стриминга
        5. gRPC Gateway для предоставления REST поверх gRPC


![](../assets/obsidian-images/Pasted%20image%2020250826140229.png)|![](../assets/obsidian-images/Pasted%20image%20.png)

![](../assets/obsidian-images/Pasted%20image%2020250826140315.png)|![](../assets/obsidian-images/Pasted%20image%20.png)

![](../assets/obsidian-images/Pasted%20image%2020250810022815.png)

---
ccna tarmoq administratorligi


Tarmoq komponentlari:
1. Qurilma(device)
    1. ohirgi qurilma: pc, lapto
2. Uzatish muhiti
3. Xizmat(service)

-----------------------------

Model

1. OSI(open systems interconnection) // g'oya 1969 yillarda berilgan, amalda qo'llangan yili 1984

Layers:
    7) Application (ilova)
    6) Presentation (taqdimot) formatlar bilan ishlaydi, masalan pc da .word ustiga bosilsa uni word app ga yo'naltiradimi?
    5) Session (seans)
    4) Transport (transport) -> MTU 1500 byte
    3) Network (tarmoq) -> Ip addressing
    2) Data link (kanal) -> MAC address
    1) Phisical (fizikaviy) malumotlar 01 larga o'giriladi va 3 hil yo'l bilan jo'natiladi 
        1. Qabel bo'lsa elektr ko'rinishida
        2. Optika bo'lsa nur ko'rinishida
        3. wifi bo'lsa bez pravadnoy ko'rinishida

7,6,5 - tarmoqqa bog'liq emas, yani os ichida yoziladi
4,3,2,1 - tarmoqqa bog'liq

Layerlardagi malumotlar atamasi:

7,6,5 - data
4 - segment
3 - packet
2 - frame
1 - bits

![alt text](image.png)



2. TCP/IP protocols stack suit // ariginal nomi DOD bo'lgan