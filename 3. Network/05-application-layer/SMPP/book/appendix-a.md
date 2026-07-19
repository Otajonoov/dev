# Appendix A. Command ID jadvali

> Mustaqil lookup-hujjat: v3.4 Table 5-1 (§5.1.2.1) bo'yicha barcha command_id qiymatlari.
> Kodda: `code/pdu/command.go` (`Cmd*` konstantalari). Batafsil: PDU anatomiyasi — [2-bob](02-pdu-anatomiyasi.md), dispatcher — [10-bob](10-boshqa-operatsiyalar.md).

## Asosiy qoida

Response id = request id | 0x80000000 (bit 31). Request'lar 0x00000000–0x000001FF, response'lar 0x80000000–0x800001FF diapazonida (§3.2). Kodda: `id.IsResponse()`, `id.Resp()`.

## To'liq jadval

| PDU | Request ID | Response ID | Bob | Izoh |
|---|---|---|---|---|
| generic_nack | — | 0x80000000 | [11](11-error-handling.md) | FAQAT response ko'rinishda mavjud |
| bind_receiver | 0x00000001 | 0x80000001 | [4](04-bind-session.md) | |
| bind_transmitter | 0x00000002 | 0x80000002 | [4](04-bind-session.md) | |
| query_sm | 0x00000003 | 0x80000003 | [10](10-boshqa-operatsiyalar.md) | |
| submit_sm | 0x00000004 | 0x80000004 | [5](05-submit-deliver.md) | |
| deliver_sm | 0x00000005 | 0x80000005 | [5](05-submit-deliver.md), [9](09-dlr.md) | MO xabar HAM, DLR HAM |
| unbind | 0x00000006 | 0x80000006 | [4](04-bind-session.md) | |
| replace_sm | 0x00000007 | 0x80000007 | [10](10-boshqa-operatsiyalar.md) | |
| cancel_sm | 0x00000008 | 0x80000008 | [10](10-boshqa-operatsiyalar.md) | |
| bind_transceiver | 0x00000009 | 0x80000009 | [4](04-bind-session.md) | v3.4 yangiligi |
| outbind | 0x0000000B | — | [4](04-bind-session.md) | Resp YO'Q (javob — ESME'ning bind_receiver'i) |
| enquire_link | 0x00000015 | 0x80000015 | [4](04-bind-session.md) | |
| submit_multi | 0x00000021 | 0x80000021 | [10](10-boshqa-operatsiyalar.md) | |
| alert_notification | 0x00000102 | — | [10](10-boshqa-operatsiyalar.md) | Resp YO'Q |
| data_sm | 0x00000103 | 0x80000103 | [10](10-boshqa-operatsiyalar.md) | |

## Diapazonlar

| Diapazon | Ma'no |
|---|---|
| 0x0000000A, 0x0000000C–0x00000014, 0x00000016–0x00000020, 0x00000022–0x000000FF | Reserved |
| 0x00000100 | Reserved (SMPP 4.0'ning "SC-side" merosxo'ri) |
| 0x00000101 | Reserved |
| 0x00010000–0x000101FF | Reserved SMSC vendor uchun emas — SMPP kengaytmalari |
| **0x00010200–0x000102FF** | **SMSC vendor-specific** (+0x80010200–0x800102FF resp'lari) |
| Qolgan hamma qiymat | Reserved |

Notanish command_id kelganda: generic_nack + ESME_RINVCMDID (0x03) — [11-bob](11-error-handling.md); kodda `pdu.Decode` `ErrUnknownCommandID` qaytaradi.

## Eslab qolish uchun uch g'alatilik

1. **generic_nack'ning request formasi yo'q** — 0x00000000 reserved; shuning uchun u "javobi bo'lmagan yagona javob".
2. **outbind (0x0B) va alert_notification (0x102) resp kutmaydi** — session engine'da "javob yuborilmaydigan PDU'lar" istisno ro'yxati shu ikkitadan iborat.
3. **0x80000005 (deliver_sm_resp)** — 1-bob mashqidagi savol: bit 31 o'rnatilgani uchun bu RESPONSE, ya'ni uni ESME yuboradi, SMSC oladi; faqat BOUND_RX/BOUND_TRX sessiyalarda uchraydi.
