# SMPP v3.4 — Go'da noldan, production darajasida

SMS gateway dasturchisi uchun **kitob darajasidagi o'quv material + bob-ma-bob o'sib boradigan Go implementatsiyasi**. Uslub modeli — Thorsten Ball, "Writing an Interpreter in Go": har bob = nazariya + shu bobda yoziladigan, testlangan kod.

**Maqsad:** materialni boshidan oxirigacha o'qigan Go dasturchi SMPP v3.4'ni tashqi manbasiz, to'liq va to'g'ri qayta yoza oladi. Har texnik fakt spec reference bilan; har PDU qo'lda tekshirilgan hex dump bilan; har hex golden test bilan qotirilgan.

## Mundarija

### Boblar

| # | Bob | Mavzu | Kod milestone |
|---|---|---|---|
| 1 | [SMS ekotizimi va SMPP'ning o'rni](book/01-sms-ekotizimi.md) | SMSC/ESME, MO/MT, versiyalar, muqobillar | `pdu/command.go` |
| 2 | [PDU anatomiyasi](book/02-pdu-anatomiyasi.md) | 16-baytlik header, data type'lar, TCP framing | `pdu/header,frame,types` |
| 3 | [TLV — optional parameter'lar](book/03-tlv.md) | Tag/Length/Value, forward compatibility | `tlv/` |
| 4 | [Bind va session lifecycle](book/04-bind-session.md) | TX/RX/TRX, state machine, timer'lar | `pdu/bind,simple`, `session/state` |
| 5 | [submit_sm va deliver_sm](book/05-submit-deliver.md) | 17 field, esm_class/registered_delivery bitlari, vaqt formati | `pdu/submit_sm,deliver_sm,esm,time` |
| 6 | [Addressing: TON/NPI](book/06-addressing.md) | Jadvallar, alphanumeric sender, real kombinatsiyalar | `pdu/address.go` |
| 7 | [Text encoding](book/07-encoding.md) | GSM7/UCS2, limitlar matematikasi, o'zbek matni | `coding/gsm7,ucs2,detect` |
| 8 | [Concatenation](book/08-concat.md) | UDH / sar_* / message_payload, splitter | `coding/udh,split` |
| 9 | [Delivery receipt (DLR)](book/09-dlr.md) | Tolerant parser, hex/dec korrelyatsiya | `dlr/` |
| 10 | [Qolgan operatsiyalar](book/10-boshqa-operatsiyalar.md) | data_sm, query, cancel, replace, multi, alert + dispatcher | `pdu/*_sm.go`, `pdu/decode.go` |
| 11 | [Error handling](book/11-error-handling.md) | command_status, generic_nack, retry tasnifi | `pdu/status.go`, `client/retry.go` |
| 12 | [Session engine](book/12-session-engine.md) | Goroutine'lar, pending window, flow control | `session/` to'liq |
| 13 | [ESME client API](book/13-client-api.md) | Uch kutubxona taqqosi, reconnect, SubmitLong | `client/` |
| 14 | [Mock SMSC](book/14-mock-smsc.md) | Server state machine, DLR generator, quirk rejimlari | `smsc/server.go` |
| 15 | [Testing va tooling](book/15-testing.md) | Fuzzing (real bug hikoyasi!), benchmark, Wireshark, CI | fuzz/bench, `ci.sh` |
| 16 | [Production](book/16-production.md) | TLS, monitoring, PII masking, quirk katalogi, v5.0 | `client/tls,metrics,logmask`, e2e |

### Ilovalar (mustaqil lookup-hujjatlar)

- [A — Command ID jadvali](book/appendix-a.md)
- [B — command_status to'liq jadvali (+tasnif, v5.0)](book/appendix-b.md)
- [C — Standard TLV tag'lar](book/appendix-c.md)
- [D — GSM 03.38 alphabet](book/appendix-d.md)
- [E — SMPP timer'lari](book/appendix-e.md)
- [F — Glossary (60+ termin)](book/appendix-f.md)
- [G — Tezkor shpargalka: DLR/TON/NPI/DCS](book/appendix-g.md)

## Papkalar

```
SMPP/
├── book/               # kitob: 16 bob + 7 ilova
├── exercises/          # har bob uchun mashqlar VA yechimlar (01–16)
├── code/               # Go module (stdlib-only), bob-ma-bob qurilgan
│   ├── pdu/            # header, framing, barcha 15 PDU codec, dispatcher, status
│   ├── tlv/            # optional parameter'lar
│   ├── coding/         # GSM7/UCS2, normalize (o'zbek matni), UDH, splitter
│   ├── dlr/            # tolerant DLR parser + hex/dec korrelyatsiya
│   ├── session/        # session engine: window, sequencer, timer'lar
│   ├── client/         # ESME client: reconnect, rate limit, TLS, metrics, masking
│   ├── smsc/           # mock SMSC (quirk rejimlari bilan)
│   ├── examples/       # localsmsc, prometheus adapter, e2e (yakuniy demo)
│   └── ci.sh           # fmt + vet + build + test -race + qisqa fuzz
└── resources/          # original spec PDF'lar (haqiqat manbai) + linklar
```

## Tezkor start

```bash
cd code

./ci.sh                       # to'liq tekshiruv zanjiri
go test ./examples/e2e -v     # yakuniy demo: bind → o'zbek/kirill matn →
                              # concat → DLR (hex/dec quirk bilan korrelyatsiya)
go run ./examples/localsmsc   # jonli mock SMSC (quirk flag'lari bilan)
```

## Haqiqat manbalari

1. [resources/SMPP_v3_4_Issue1_2.pdf](resources/SMPP_v3_4_Issue1_2.pdf) — original spec, yakuniy hakam
2. [resources/SMPP_v5.pdf](resources/SMPP_v5.pdf) — v5.0 (farqlar uchun)
3. 3GPP TS 23.038 / 23.040 — encoding va UDH normativi
4. [resources/links.md](resources/links.md) — izohlangan tashqi linklar

Kitob uslubi: spec fakti va industriya konventsiyasi HAR DOIM ajratiladi; har bobda kamida bitta "⚠ Amaliyotda" bloki — real operator nozikliklari.
