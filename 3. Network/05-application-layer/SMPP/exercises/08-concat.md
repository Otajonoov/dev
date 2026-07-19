# 8-bob mashqlari: Concatenation

> Avval mashqlarni mustaqil bajaring, keyin yechimlarga qarang. Bob matni: [book/08-concat.md](../book/08-concat.md)

---

## Mashq 1. UDH'ni o'qish

deliver_sm keldi: esm_class=0x40, short_message boshlanishi:

```
05 00 03 5A 03 02 53 61 6C 6F 6D ...
```

1. UDH'ni bayt-ma-bayt oching: bu qaysi xabarning nechanchi segmenti?
2. Xabar to'liq yig'ilishi uchun telefon yana nimalarni kutyapti (qiymatlari bilan)?
3. Matn qaysi oktetdan boshlanadi va bu segmentda maksimal necha GSM7 belgi bo'lishi mumkin?
4. Xuddi shu UDH 16-bit reference bilan qanday ko'rinardi?

## Mashq 2. 300 belgilik UCS2 hisobi

300 belgilik kirill matn (surrogate'siz) UDH 8-bit usulida yuboriladi.

1. Nechta segmentga bo'linadi? Har segmentda nechta belgi? Hisob bilan.
2. Har segmentning short_message'i necha OKTET bo'ladi?
3. Jami nechta submit_sm, message_id va (registered_delivery=1 bo'lsa) DLR kutiladi?
4. `coding.CountSegments` va `coding.Split` bilan tekshiring.

## Mashq 3. Unutilgan sar_total_segments

ESME uzun xabarni sar usulida yubordi, lekin bug tufayli faqat `sar_msg_ref_num` va `sar_segment_seqnum` TLV'lari ketdi — `sar_total_segments` yo'q.

1. Spec bo'yicha SMSC bu PDU'ga nima qiladi? Aniq §'ni keltiring.
2. Natijada abonent nimani ko'radi?
3. Bu xulq (xato qaytarmaslik!) qaysi umumiy SMPP tamoyilidan kelib chiqadi va yuboruvchi bug'ni qanday sezadi?
4. `coding.SarFromTLVs` bu holatda nima qaytaradi — kod bilan ko'rsating.

---

# Yechimlar

## Yechim 1

**1.** `05` UDHL=5; `00` IEI=concat 8-bit; `03` IEDL=3; `5A` ref=0x5A (90); `03` total=3; `02` seq=2. Demak: **ref=0x5A xabarning 3 segmentidan 2-chisi**.

**2.** Xuddi shu yuboruvchidan ref=0x5A, total=3 bo'lgan **seq=1 va seq=3** segmentlarini. Uchalasi kelgunicha ko'p telefonlar xabarni ko'rsatmay turadi (ba'zilari qismlarni alohida ko'rsatib, keyin birlashtiradi).

**3.** Matn 6-oktetdan (offset 6: `53 61 6C 6F 6D` = "Salom..."). Sig'im: UDH 6 oktet oldi → (140−6)×8/7 = **153 belgi** (GSM7, extension'siz).

**4.** `06 08 04 00 5A 03 02` — UDHL=6, IEI=0x08, IEDL=4, ref 2 oktetga kengaydi (0x005A), sig'im esa 152 ga tushardi.

## Yechim 2

**1.** UCS2 segment sig'imi (UDH bilan): 67 belgi. 300 > 70 → concat: ceil(300/67) = **5 segment**: 67+67+67+67+32.

**2.** Dastlabki to'rttasi: 6 (UDH) + 67×2 = **140 oktet** (havo limitining o'zi — tasodif emas!); oxirgisi: 6 + 32×2 = **70 oktet**.

**3.** **5 ta submit_sm** (har biri o'z sequence_number'i bilan!), 5 ta message_id, 5 ta DLR. "Xabar yetkazildi" = beshala DLR ham DELIVRD bo'lganda (9-bob korrelyatsiyasi).

**4.**

```go
dc, n, _ := coding.CountSegments(matn300)   // DCUCS2, 5
segs, _ := coding.Split(matn300, dc, coding.MethodUDH8, ref)
// len(segs) == 5; len(segs[0].Data) == 140; len(segs[4].Data) == 70
```

## Yechim 3

**1.** §5.3.2.22 (sar_msg_ref_num ta'rifi): "When present, the PDU must also contain the sar_total_segments and sar_segment_seqnum parameters. **Otherwise this parameter shall be ignored.**" Uchlik to'liq emas → SMSC sar TLV'larni butunlay E'TIBORSIZ qoldiradi va PDU'larni ODDIY, bog'lanmagan xabarlar sifatida qabul qiladi. Xato qaytarilmaydi!

**2.** Abonent uzun xabar o'rniga **N ta alohida, "kesilgan" SMS** oladi — tartibi ham aralash bo'lishi mumkin.

**3.** Bu 3-bobdagi forward-compatibility falsafasining davomi: TLV darajasidagi nomukammallik — jim ignore, xato emas. Yuboruvchi buni submit_sm_resp'dan SEZA OLMAYDI (hammasi status=0!) — faqat abonent shikoyati yoki test telefonidagi tekshiruv orqali. Saboq: sar usulini ishlatishdan oldin jonli qurilmada sinash shart (7-bobdagi test-matritsaga qo'shimcha band), va TLV'larni qo'lda emas, `SarInfo.TLVs()` kabi UCHLIKNI YAXLIT beradigan helper bilan qo'shish — bittasini unutish strukturaviy imkonsiz bo'ladi.

**4.**

```go
tlvs := []tlv.TLV{
	tlv.U16(tlv.SarMsgRefNum, 0x5A),
	tlv.U8(tlv.SarSegmentSeqnum, 1),
	// sar_total_segments YO'Q
}
info, found, err := coding.SarFromTLVs(tlvs)
// found == false, err == nil — spec bo'yicha "uchlik yo'q" holati;
// info bo'sh. Qabul qiluvchi kod bu PDU'ni oddiy xabar sifatida davom ettiradi.
```
