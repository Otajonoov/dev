# 06 — API protokollari

Bu modul zamonaviy backend dasturchi kundalik ishlatadigan **API dizayni va
protokollari**ga bag'ishlangan: REST arxitekturasidan tortib autentifikatsiya,
real vaqtli WebSocket va yuqori unumli gRPC gacha.

Bu yerda **protokol nazariyasi va dizaynga** fokus qilinadi. HTTP asoslari va
TLS oldingi (application-layer) modulda ko'rilgan — bu yerda takrorlanmaydi.
Go bilan chuqur amaliyot esa keyingi (Go network programming) modulda bo'ladi;
bu modulda faqat kichik, tushuntiruvchi Go misollar bor.

## Nima o'rganiladi

- REST nima va nega u API dunyosining standartiga aylandi
- REST'ni haqiqatan RESTful qiladigan 6 ta constraint
- Toza, izchil URI (endpoint) dizayni qoidalari
- API'ni himoyalash: Basic Auth, API key, JWT, OAuth2/OIDC, PASETO
- WebSocket bilan real vaqtli, ikki tomonlama aloqa
- gRPC bilan tez, binary, kontraktga asoslangan servislararo aloqa

## Darslar

1. [REST nima](01-rest-nima.md) — tarix, Fielding, resource, Richardson modeli
2. [REST constraints](02-rest-constraints.md) — 6 ta arxitektura cheklovi
3. [REST resource naming](03-rest-resource-naming.md) — URI dizayni, 10 oltin qoida
4. [API autentifikatsiya](04-api-autentifikatsiya.md) — JWT, OAuth2, PKCE, PASETO
5. [WebSocket](05-websocket.md) — handshake, frame'lar, SSE bilan taqqoslash
6. [gRPC](06-grpc.md) — Protocol Buffers, HTTP/2, streaming, gRPC-Gateway

## O'qish tartibi

Darslar ketma-ket, yuqoridan pastga o'qilishi tavsiya etiladi. REST uchligi
(1-3) poydevor beradi: avval nima ekanini (1), keyin qoidalarini (2), so'ngra
amaliy dizaynini (3) o'rganasan. 4-dars har qanday API uchun xavfsizlik
qatlamini qo'shadi. 5 va 6-darslar REST'ga muqobil protokollarni ochadi:
WebSocket real vaqtli aloqa uchun, gRPC esa yuqori unumli servislararo aloqa uchun.

Har dars oxirida `O'z-o'zini tekshir`, `Amaliyot` va `Takrorlash` bo'limlari
bor — ularni tashlab ketma, bilim aynan eslab chiqarish orqali mustahkamlanadi.
