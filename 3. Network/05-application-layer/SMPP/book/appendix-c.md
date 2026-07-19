# Appendix C. Standard TLV tag'lar

> Mustaqil lookup-hujjat: v3.4 Table 5-7 (§5.3) barcha SMPP-defined TLV'lari, value formatlari bilan.
> Kodda: `code/tlv/tlv.go` (Tag konstantalari) va `code/tlv/values.go` (tipli helper'lar).
> TLV mexanizmi: [3-bob](03-tlv.md). Length har doim FAQAT Value uzunligi.

## Tag bloklari

| Diapazon | Kim belgilaydi |
|---|---|
| 0x0001–0x00FF, 0x0200–0x05FF, 0x1200–0x13FF | SMPP-defined (Table 5-7) |
| **0x1400–0x3FFF** | **SMSC vendor-specific** |
| Qolgani | Reserved |

## To'liq jadval

"Texn." ustuni — qaysi tarmoq texnologiyasiga tegishli (Generic = hammasi). "Qayerda" — asosan uchraydigan PDU'lar (to'liq ruxsat jadvali spec §4'ning har PDU bo'limida).

| Tag | Nom | Value | Texn. | Qayerda / izoh |
|---|---|---|---|---|
| 0x0005 | dest_addr_subunit | 1 oktet: 0=unknown, 1=MS display, 2=mobile equipment, 3=smart card, 4=external unit | GSM | submit/data |
| 0x0006 | dest_network_type | 1 oktet: 1=GSM 2=TDMA 3=CDMA 4=PDC 5=PHS 6=iDEN 7=AMPS 8=Paging | Generic | submit/data |
| 0x0007 | dest_bearer_type | 1 oktet: 1=SMS 2=CSD 3=Packet 4=USSD 5=CDPD 6=DataTAC 7=FLEX 8=CellBroadcast | Generic | submit/data |
| 0x0008 | dest_telematics_id | 2 oktet | GSM | submit/data |
| 0x000D | source_addr_subunit | 1 oktet (0x0005 kabi) | GSM | deliver/data |
| 0x000E | source_network_type | 1 oktet (0x0006 kabi) | Generic | deliver/data |
| 0x000F | source_bearer_type | 1 oktet (0x0007 kabi) | Generic | deliver/data |
| 0x0010 | source_telematics_id | 1 oktet (⚠ dest'niki 2 — spec'ning o'z nosimmetriyasi) | GSM | deliver/data |
| 0x0017 | qos_time_to_live | 4 oktet Integer — TTL SEKUNDLARDA (data_sm'ning "validity"si) | Generic | data_sm |
| 0x0019 | payload_type | 1 oktet: 0=WDP, 1=WCMP | Generic | data_sm |
| 0x001D | additional_status_info_text | C-Octet, max 256 — ASCII izoh | Generic | data_sm_resp |
| 0x001E | receipted_message_id | C-Octet, max 65 — original message_id ([9-bob](09-dlr.md)) | Generic | deliver_sm (DLR) |
| 0x0030 | ms_msg_wait_facilities | 1 oktet: bit 7=indikator yoq/o'chir, bit 1-0: 00=voicemail 01=fax 10=email 11=other | GSM | submit |
| 0x0201 | privacy_indicator | 1 oktet: 0–3 | CDMA/TDMA | submit/deliver/data |
| 0x0202 | source_subaddress | 2–23 oktet | CDMA/TDMA | deliver/data |
| 0x0203 | dest_subaddress | 2–23 oktet | CDMA/TDMA | submit/data |
| 0x0204 | user_message_reference | 2 oktet — application korrelyatsiya raqami | Generic | hamma xabar PDU'lari |
| 0x0205 | user_response_code | 1 oktet | CDMA/TDMA | deliver |
| 0x020A | source_port | 2 oktet ([8-bob](08-concat.md): UDHI bilan birga EMAS) | Generic | data/submit |
| 0x020B | destination_port | 2 oktet | Generic | data/submit |
| 0x020C | sar_msg_ref_num | 2 oktet — barcha segmentlarda BIR XIL ([8-bob](08-concat.md)) | Generic | submit/data |
| 0x020D | language_indicator | 1 oktet | CDMA/TDMA | submit/deliver |
| 0x020E | sar_total_segments | 1 oktet, 1–255; ⚠ sar uchligi BIRGA yoki ignore | Generic | submit/data |
| 0x020F | sar_segment_seqnum | 1 oktet, 1'dan | Generic | submit/data |
| 0x0210 | sc_interface_version | 1 oktet (0x34) — [4-bob](04-bind-session.md): bind_resp'da | Generic | bind_resp |
| 0x0302 | callback_num_pres_ind | 1 oktet | TDMA | submit/deliver |
| 0x0303 | callback_num_atag | max 65 — display alfanumerik tag | TDMA | submit/deliver |
| 0x0304 | number_of_messages | 1 oktet, 0–99 | CDMA | submit |
| 0x0381 | callback_num | 4–19 oktet: [digit mode][TON][NPI][raqam...] — TAKRORLANISHI mumkin! | Hammasi | submit/deliver |
| 0x0420 | dpf_result | 1 oktet: 0=dpf o'rnatilmadi, 1=o'rnatildi ([10-bob](10-boshqa-operatsiyalar.md)) | Generic | data_sm_resp |
| 0x0421 | set_dpf | 1 oktet: 1=yetkaza olmasang dpf o'rnat | Generic | data_sm |
| 0x0422 | ms_availability_status | 1 oktet: 0=available, 1=denied, 2=unavailable | Generic | alert_notification |
| 0x0423 | network_error_code | **3 oktet struktura**: [type: 1=ANSI-136, 2=IS-95, 3=GSM] + [2 oktet kod] | Generic | deliver (DLR), data_sm_resp |
| 0x0424 | message_payload | var — matn/data; max "implementation specific" (64K — faqat TLV nazariy limiti) | Generic | submit/deliver/data |
| 0x0425 | delivery_failure_reason | 1 oktet: 0=dest unavailable, 1=dest invalid, 2=permanent net, 3=temp net | Generic | data_sm_resp (transaction) |
| 0x0426 | more_messages_to_send | 1 oktet: 0=oxirgisi, 1=yana bor (default) | GSM | submit/data |
| 0x0427 | message_state | 1 oktet: §5.2.28 qiymatlari (1=ENROUTE...8=REJECTED, [9-bob](09-dlr.md)) | Generic | deliver (DLR) |
| 0x0501 | ussd_service_op | 1 oktet: 0/1=PSSD/PSSR ind, 2/3=USSR/USSN req, 16–19=javoblar | GSM/USSD | data_sm |
| 0x1201 | display_time | 1 oktet: 0=temporary, 1=default, 2=invoke | CDMA/TDMA | submit/deliver |
| 0x1203 | sms_signal | 2 oktet | TDMA | submit |
| 0x1204 | ms_validity | 1 oktet: 0=store indefinitely, 1=power down, 2=SID based, 3=display only | CDMA/TDMA | submit |
| 0x130C | alert_on_message_delivery | **0 oktet — zero-length!** (codec'lar uchun klassik test case, [3-bob](03-tlv.md)) | CDMA | submit |
| 0x1380 | its_reply_type | 1 oktet | CDMA | deliver |
| 0x1383 | its_session_info | 2 oktet | CDMA | submit/deliver |

## Amaliy eslatmalar

1. **Notanish tag = ignore, lekin TASHLAMANG** (§3.3 forward compatibility) — bizning `tlv.Decode` saqlaydi, siyosat caller'da.
2. **Length'ni Tag+Length+Value jami deb hisoblash** — 3-bobdagi klassik xato: Length FAQAT Value.
3. **callback_num takrorlanishi mumkin** — shuning uchun TLV to'plami map emas, slice (`tlv.Find` birinchisini beradi, qolganlari uchun to'plamni aylanib chiqiladi).
4. Kundalik A2P ishida real uchraydigani o'ntacha: receipted_message_id, message_state, network_error_code (DLR); message_payload, sar_* (uzun xabar); sc_interface_version (bind); user_message_reference; qos_time_to_live, set_dpf/dpf_result (data_sm). Qolganlari — texnologiyaga xos ekzotika, lekin decoder hammasini tanishi kerak.
