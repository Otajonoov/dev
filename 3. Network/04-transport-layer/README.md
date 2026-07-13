# 04. Transport Layer (TCP, UDP)

Bu modul tarmoq stack'ining **transport qatlami** — backend dasturchi uchun eng muhim
qismini o'rgatadi. Network layer ma'lumotni **kompyuterga** yetkazsa, transport layer
uni to'g'ri **dasturga** (process'ga) yetkazadi va, kerak bo'lganda, ishonchsiz IP
ustida **ishonchli** aloqa quradi.

## Nima o'rganiladi

- Transport layer nima uchun kerak va nima uchun ikkita protokol (TCP va UDP) bor.
- **Port** va **socket** orqali multiplexing/demultiplexing qanday ishlaydi.
- **UDP** — sodda, tez, connectionless protokol; qachon tanlanadi; QUIC bilan bog'liqligi.
- **TCP** — reliable, connection-oriented protokol; sequence/ACK, retransmission, flow control.
- **Three-way handshake**, connection teardown, TIME_WAIT, half-open, SYN flood va SYN cookies.
- **Flow control** va **congestion control** (slow start, AIMD, CUBIC vs BBR).

## Darslar ro'yxati (o'qish tartibi)

1. [`01-transport-layer-vazifasi.md`](01-transport-layer-vazifasi.md) — transport
   layer roli; TCP va UDP nega ikkita.
2. [`02-multiplexing-demultiplexing.md`](02-multiplexing-demultiplexing.md) — port,
   socket, UDP 2-tuple va TCP 4-tuple.
3. [`03-udp.md`](03-udp.md) — UDP header, checksum, use cases, QUIC nima uchun UDP ustida.
4. [`04-tcp.md`](04-tcp.md) — TCP header, sequence/ACK, retransmission, SACK, flow control.
5. [`05-tcp-handshake-va-connection.md`](05-tcp-handshake-va-connection.md) —
   handshake, teardown, TIME_WAIT, SYN flood, SYN cookies.
6. [`06-flow-va-congestion-control.md`](06-flow-va-congestion-control.md) — rwnd/cwnd,
   slow start, AIMD, CUBIC vs BBR.

## Tavsiya etilgan tartib

Darslar bir-birining ustiga quriladi, shuning uchun **tartib bilan** o'qish tavsiya
etiladi: 01 → 02 avval tushunchaviy poydevor beradi; 03 (UDP) soddaroq, shu sabab
undan boshlash yaxshi; keyin 04 → 05 → 06 TCP'ni bosqichma-bosqich chuqurlashtiradi.
Har darsning oxiridagi **O'z-o'zini tekshir** va **Amaliyot** bo'limlarini o'tkazib
yubormang — bilim aynan shu yerda mustahkamlanadi.
