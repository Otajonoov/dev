# 05. Application Layer — foydalanuvchi ko'radigan protokollar

Bu modul tarmoq stekining eng yuqori qatlamini — **Application Layer** ni
o'rgatadi. Bu yerda foydalanuvchi va dasturlar bevosita ishlatadigan protokollar
yashaydi: web (HTTP/HTTPS), nom yechish (DNS), fayl uzatish (FTP/SFTP), email
(SMTP/IMAP), markazsiz almashinuv (BitTorrent) va SMS (SMPP). Asosiy maqsad —
har protokolning *nima muammoni yechishini*, *qanday ishlashini* va *2025-2026
zamonaviy holatini* tushunish.

Har dars bir xil pedagogik tuzilishda: **muammo/hook -> real-hayot analogiya ->
sodda ta'rif -> Mermaid diagramma -> worked example (curl/dig/openssl misollari)
-> Xulosa, Eslab qol, O'z-o'zini tekshir, Amaliyot, Takrorlash**. Har dars
zamonaviy best practice bilan boyitilgan (HTTP/3, TLS 1.3 + post-quantum,
DoH/DoT/DoQ, SPF/DKIM/DMARC, SFTP, DHT/magnet, SMPP 3.4).

## Darslar ro'yxati (o'qish tartibi)

| # | Dars | Mavzu |
|---|------|-------|
| 1 | [01-application-layer-va-socketlar.md](01-application-layer-va-socketlar.md) | Application arxitekturalari (client-server vs P2P), socket, session/presentation |
| 2 | [02-dns.md](02-dns.md) | DNS ierarxiya, recursive/iterative, RR turlari, caching/TTL, DoH/DoT/DoQ |
| 3 | [03-http.md](03-http.md) | Request/response, metodlar, status kodlar, headerlar, cookie, keep-alive |
| 4 | [04-http-evolution.md](04-http-evolution.md) | HTTP/0.9 -> 1.1 -> 2 -> 3, multiplexing, QUIC, HOL blocking |
| 5 | [05-https-tls.md](05-https-tls.md) | TLS handshake, sertifikat/CA zanjiri, TLS 1.3, mTLS, post-quantum |
| 6 | [06-smtp-va-email.md](06-smtp-va-email.md) | SMTP, POP3/IMAP, email oqimi, SPF/DKIM/DMARC |
| 7 | [07-ftp-sftp.md](07-ftp-sftp.md) | FTP active/passive, FTPS vs SFTP farqi |
| 8 | [08-p2p-bittorrent.md](08-p2p-bittorrent.md) | P2P arxitektura, tracker, chunk, tit-for-tat, DHT/magnet |
| 9 | [09-smpp.md](09-smpp.md) | SMPP 3.4, PDU, bind rejimlari, enquire_link, data_coding, SMSC |

## Qanday o'qish kerak

1. Darslarni **tartib bilan** o'qi — har biri oldingisiga tayanadi (masalan HTTPS
   DNS va HTTP tushunchasini talab qiladi).
2. Har darsda avval **muammoni his qil**, keyin yechimni o'rgan.
3. Worked example dagi `curl`, `dig`, `openssl` buyruqlarini o'zing terminalda
   yozib sina.
4. Har dars oxiridagi **O'z-o'zini tekshir** savollariga javobni ochishdan oldin
   o'zing javob berishga urin (retrieval practice).
5. **Takrorlash jadvali** bo'yicha ertaga -> 3 kun -> 1 hafta oralig'ida qayt.

## Umumiy tekshiruv buyruqlari

```bash
dig +trace example.com                 # DNS to'liq zanjir
curl -v https://example.com            # HTTP/HTTPS request
curl -sI https://example.com | grep -i alt-svc   # HTTP/3 mavjudligi
openssl s_client -connect example.com:443 -tls1_3  # TLS handshake + cert
dig _dmarc.example.com TXT +short      # email himoyasi (DMARC)
sftp user@server                       # xavfsiz fayl uzatish
ss -tnp                                # ochiq socketlar
```

## Eslatma

REST, WebSocket, gRPC, JWT/OAuth kabi API protokollari alohida **06-api-protokollari**
modulida; Go bilan socket dasturlash **10-go-network-programming** modulida.
