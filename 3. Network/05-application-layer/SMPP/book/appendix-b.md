# Appendix B. command_status to'liq jadvali

> Mustaqil lookup-hujjat: v3.4 Table 5-2 (§5.1.3) barcha kodlari + reaksiya tasnifi.
> Kodda: `code/pdu/status.go` (nomlar) va `code/client/retry.go` (`Classify`).
> Tasnif ustuni spec'da YO'Q — industriya konsensusi ([11-bob](11-error-handling.md)).
> Tasnif kaliti: **T** = Transient (backoff bilan retry), **P** = Permanent (retry ma'nosiz),
> **S** = Session-level (rebind/credential yo'li), **U** = Unknown (operator hujjati + cheklangan retry).

## v3.4 kodlari

| Hex | Dec | Nom | Ma'no | Tasnif |
|---|---|---|---|---|
| 0x00 | 0 | ESME_ROK | Xato yo'q | — |
| 0x01 | 1 | ESME_RINVMSGLEN | Message length invalid (sm_length buzuq) | P |
| 0x02 | 2 | ESME_RINVCMDLEN | Command length invalid (framing!) | P* |
| 0x03 | 3 | ESME_RINVCMDID | Notanish command_id | P* |
| 0x04 | 4 | ESME_RINVBNDSTS | Bu state'da bu PDU mumkin emas (Table 2-1) | S |
| 0x05 | 5 | ESME_RALYBND | Allaqachon bound | S |
| 0x06 | 6 | ESME_RINVPRTFLG | priority_flag buzuq (⚠ RINVPASWD EMAS!) | P |
| 0x07 | 7 | ESME_RINVREGDLVFLG | registered_delivery buzuq | P |
| 0x08 | 8 | ESME_RSYSERR | SMSC ichki xatosi | T |
| 0x0A | 10 | ESME_RINVSRCADR | Source manzil buzuq (⚠ RMSGQFUL EMAS!) | P |
| 0x0B | 11 | ESME_RINVDSTADR | Destination manzil buzuq | P |
| 0x0C | 12 | ESME_RINVMSGID | message_id topilmadi (query/cancel/replace) | P |
| 0x0D | 13 | ESME_RBINDFAIL | Bind muvaffaqiyatsiz (umumiy) | S |
| 0x0E | 14 | ESME_RINVPASWD | Parol xato | S |
| 0x0F | 15 | ESME_RINVSYSID | system_id notanish | S |
| 0x11 | 17 | ESME_RCANCELFAIL | cancel_sm bo'lmadi | P |
| 0x13 | 19 | ESME_RREPLACEFAIL | replace_sm bo'lmadi | P |
| 0x14 | 20 | ESME_RMSGQFUL | Navbat to'la (⚠ dec 20 — 0x20 EMAS!) | T |
| 0x15 | 21 | ESME_RINVSERTYP | service_type buzuq | P |
| 0x33 | 51 | ESME_RINVNUMDESTS | number_of_dests buzuq (submit_multi) | P |
| 0x34 | 52 | ESME_RINVDLNAME | Distribution List topilmadi | P |
| 0x40 | 64 | ESME_RINVDESTFLAG | dest_flag 1/2 emas (submit_multi) | P |
| 0x42 | 66 | ESME_RINVSUBREP | "submit with replace" buzuq | P |
| 0x43 | 67 | ESME_RINVESMCLASS | esm_class buzuq | P |
| 0x44 | 68 | ESME_RCNTSUBDL | DL'ga submit qilib bo'lmaydi | P |
| 0x45 | 69 | ESME_RSUBMITFAIL | submit_sm/submit_multi bo'lmadi (sabab noma'lum) | U |
| 0x48 | 72 | ESME_RINVSRCTON | Source TON buzuq | P |
| 0x49 | 73 | ESME_RINVSRCNPI | Source NPI buzuq | P |
| 0x50 | 80 | ESME_RINVDSTTON | Dest TON buzuq | P |
| 0x51 | 81 | ESME_RINVDSTNPI | Dest NPI buzuq | P |
| 0x53 | 83 | ESME_RINVSYSTYP | system_type buzuq (⚠ RINVDCS EMAS — u v5.0!) | P |
| 0x54 | 84 | ESME_RINVREPFLAG | replace_if_present buzuq | P |
| 0x55 | 85 | ESME_RINVNUMMSGS | Xabar soni buzuq | P |
| 0x58 | 88 | ESME_RTHROTTLED | Rate limit oshildi — SEKINLA | T |
| 0x61 | 97 | ESME_RINVSCHED | schedule_delivery_time buzuq | P |
| 0x62 | 98 | ESME_RINVEXPIRY | validity_period buzuq/limitdan katta | P |
| 0x63 | 99 | ESME_RINVDFTMSGID | sm_default_msg_id buzuq/topilmadi | P |
| 0x64 | 100 | ESME_RX_T_APPN | ESME (receiver): VAQTINCHA xato — qayta urin | T |
| 0x65 | 101 | ESME_RX_P_APPN | ESME (receiver): DOIMIY xato — urinma | P |
| 0x66 | 102 | ESME_RX_R_APPN | ESME (receiver): rad etildi | P |
| 0x67 | 103 | ESME_RQUERYFAIL | query_sm bo'lmadi | P |
| 0xC0 | 192 | ESME_RINVOPTPARSTREAM | TLV tail buzuq (parse bo'lmadi) | P |
| 0xC1 | 193 | ESME_ROPTPARNOTALLWD | Bu TLV bu PDU'da mumkin emas | P |
| 0xC2 | 194 | ESME_RINVPARLEN | TLV length buzuq | P |
| 0xC3 | 195 | ESME_RMISSINGOPTPARAM | Kutilgan TLV yo'q | P |
| 0xC4 | 196 | ESME_RINVOPTPARAMVAL | TLV value buzuq | P |
| 0xFE | 254 | ESME_RDELIVERYFAILURE | Yetkazish xatosi (data_sm_resp, transaction mode) | U |
| 0xFF | 255 | ESME_RUNKNOWNERR | Noma'lum xato | U |
| 0x400–0x4FF | 1024–1279 | (vendor) | SMSC vendor-specific — operator hujjatiga qarang | U |

\* 0x02/0x03 odatda **generic_nack ichida** keladi ([11-bob](11-error-handling.md)): tasnifi "permanent", lekin haqiqiy reaksiya — framing shubhasi bo'lsa reconnect.

Jadvalda YO'Q qiymatlar (0x09, 0x10, 0x12, 0x16–0x32, 0x35–0x3F, ...): **reserved** — vendor kodi EMAS (vendor faqat 0x400–0x4FF). Bunday kod kelsa: buggy SMSC yoki versiya aralashuvi; `String()` uni `unknown(0x...)` deb ko'rsatadi.

## v5.0 qo'shimchalari (v3.4'da MAVJUD EMAS)

v5.0 (Table 4-2) 0x100 blokidan yangi kodlar qo'shgan — v3.4 sessiyasida bularni YUBORISH xato, lekin buggy/aralash stack'lardan KELIB QOLISHI mumkin (bizda `unknown` bo'lib ko'rinadi):

| Hex | Nom | Ma'no |
|---|---|---|
| 0x100 | ESME_RSERTYPUNAUTH | service_type uchun ruxsat yo'q |
| 0x101 | ESME_RPROHIBITED | Operatsiya taqiqlangan |
| 0x102 | ESME_RSERTYPUNAVAIL | service_type hozir ishlamayapti |
| 0x103 | ESME_RSERTYPDENIED | service_type rad etilgan |
| **0x104** | **ESME_RINVDCS** | **data_coding buzuq — eski darslar v3.4'niki deb adashtirgan kod!** |
| 0x105 | ESME_RINVSRCADDRSUBUNIT | source_addr_subunit buzuq |
| 0x106 | ESME_RINVDSTADDRSUBUNIT | dest_addr_subunit buzuq |
| 0x107–0x112 | ESME_RINVBCAST* / ESME_RBCAST* oilasi | broadcast_sm operatsiyalari xatolari (v5.0 yangi PDU'lari) |

v3.4'da data_coding xatosi maxsus kodga EGA EMAS — RSYSERR/RSUBMITFAIL/vendor kod bilan (yoki jim qabul qilinib buzuq matn bilan) namoyon bo'ladi ([7-bob](07-encoding.md), [11-bob](11-error-handling.md)).

## Ishlatish shpargalkasi

1. Log'da har doim NOM yozing (`pdu.CommandStatus(x).String()`) — hex/dec chalkashligi yo'qoladi.
2. Reaksiya: `client.Classify` → T=queue+backoff (max-age!), P=failed+log, S=sessiya tuzat (bind'ga ham backoff), U=operator jadvali.
3. DLR `err:` bilan ADASHTIRMANG — u boshqa fazo ([9-bob](09-dlr.md)).
