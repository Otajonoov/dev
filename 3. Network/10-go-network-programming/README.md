# 10. Go tilida tarmoq dasturlash

Bu modul — butun "Network" kursining **amaliy cho'qqisi**. Avvalgi modullarda o'rgangan nazariyani (TCP/UDP, HTTP, WebSocket, gRPC, load balancing, xavfsizlik) endi **ishlaydigan Go kod**ga aylantiramiz. Har dars markazida to'liq, ishga tushadigan dastur turadi: import'laridan tortib kutilgan output'igacha.

## Nima o'rganiladi

- `net` package bilan xom TCP/UDP socket dasturlash (poydevor)
- Production darajasidagi HTTP server (timeout, middleware, graceful shutdown) va to'g'ri HTTP client
- gorilla/websocket bilan real vaqtli chat va hub pattern
- gRPC: protobuf, 4 xil streaming, interceptor, gateway, TLS
- O'z reverse proxy va load balancer'ingni noldan yozish (round-robin, least connections, health check)
- Concurrency idiomalari: "bir ulanish = bir goroutine", channel bilan aloqa, atomic bilan himoya

## Darslar ro'yxati

1. [01-net-package-asoslari.md](01-net-package-asoslari.md) — `net.Dial`, `net.Listen`, `Conn` interface, deadline/timeout, context bilan dial (butun modul poydevori)
2. [02-tcp-client-server.md](02-tcp-client-server.md) — TCP echo server, goroutine per connection, `bufio` bilan message framing, graceful shutdown, connection pool
3. [03-udp-client-server.md](03-udp-client-server.md) — `ListenUDP`, `ReadFromUDP`, UDP'da "ulanish yo'qligi", packet loss va retry
4. [04-http-server-va-client.md](04-http-server-va-client.md) — production HTTP server (mux, middleware, timeout'lar), `http.Client` (transport, connection reuse), graceful shutdown
5. [05-websocket-chat.md](05-websocket-chat.md) — gorilla/websocket bilan chat: hub pattern, read/write pump, ping/pong, broadcast
6. [06-grpc-go.md](06-grpc-go.md) — protobuf kod generatsiya, unary + 3 xil streaming, interceptor, gRPC-Gateway, TLS
7. [07-load-balancer-yasash.md](07-load-balancer-yasash.md) — application-level load balancer: round-robin, least connections, health check, `httputil.ReverseProxy`, NGINX/HAProxy taqqoslash

## O'qish tartibi

Darslarni **ketma-ket** o'qing — har biri avvalgisiga tayanadi. `net` package (1-dars) — hamma narsaning poydevori: TCP, UDP, HTTP, gRPC va load balancer oxir-oqibat shu `net.Listen`/`net.Dial` ustiga quriladi. 2-3 darslar xom socket bilan tanishtiradi, 4-6 darslar yuqori darajali protokollarni beradi, 7-dars esa barcha g'oyalarni bitta production tizimda birlashtiradi.

Har darsda: muammo/hook -> analogiya -> diagramma -> to'liq kod (bloklarga bo'lib) -> PRIMM predict savoli -> xulosa, eslab qol, o'z-o'zini tekshir, amaliyot (oson/o'rta/qiyin) va takrorlash bo'limlari bor. Amaliyot topshiriqlarini albatta bajaring — tarmoq dasturlash faqat kod yozib o'rganiladi.

## Talablar

- Go 1.22+ (yangi `ServeMux` pattern routing uchun)
- gorilla/websocket (5-dars) va google.golang.org/grpc + protoc (6-dars) uchun tashqi kutubxonalar
- Qolgan hamma narsa standart kutubxona bilan ishlaydi
