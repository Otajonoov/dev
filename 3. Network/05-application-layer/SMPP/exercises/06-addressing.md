# 6-bob mashqlari: Addressing

> Avval mashqlarni mustaqil bajaring, keyin yechimlarga qarang. Bob matni: [book/06-addressing.md](../book/06-addressing.md)

---

## Mashq 1. Besh manzilga TON/NPI tanlang

Quyidagi beshta manzil uchun TON/NPI juftligini tanlang va asoslang. Qaysilarida "yagona to'g'ri javob yo'q"ligini ham ayting:

1. `998901234567` (destination sifatida)
2. `Bank` (source sifatida)
3. `1234` (source sifatida — short code)
4. `+15551234567` (destination sifatida)
5. bo'sh source (SMSC default sender qo'ysin)

Bonus: har biri uchun `pdu` package'ida qaysi konstruktor/qiymat ishlatilishini yozing.

## Mashq 2. Nega alphanumeric sender'ga javob yozib bo'lmaydi?

Texnik zanjirni oching: abonent "Bank" sender'idan SMS oldi va "javob berish"ni bosdi. Telefon nima qilishga urinadi, qaysi bosqichda va nega muvaffaqiyatsiz bo'ladi? Javobingizda TP-OA/TON tushunchalari qatnashsin. Ikki tomonlama muloqot kerak bo'lgan xizmat (masalan "STOP" qabul qilish) qanday sender ishlatishi kerak?

## Mashq 3. address_range tahlili

bind_transceiver'da `address_range = "^9989[01]"` yuborildi.

1. Bu regex qanday manzillar to'plamini bildiradi? Kamida 2 mos va 2 mos EMAS misol keltiring.
2. Bu qiymat qaysi bind turlarida umuman ma'noga ega va nima UCHUN kerak bo'lishi ko'zda tutilgan?
3. Amalda bu sozlamadan nima kutish real — va MO routing aslida qayerda hal bo'ladi?

---

# Yechimlar

## Yechim 1

| # | Manzil | TON/NPI | Asos | Kod |
|---|---|---|---|---|
| 1 | `998901234567` | **1/1** | To'liq xalqaro E.164 format, country code bilan — eng portativ | `pdu.International("998901234567")` |
| 2 | `Bank` | **5/0** | Harflar bor → alphanumeric; 4 belgi ≤ 11, hammasi GSM7 | `pdu.Alphanumeric("Bank")` |
| 3 | `1234` | **3/0** — lekin YAGONA TO'G'RI JAVOB YO'Q | Short code TON'i standartlashmagan: 3/0 eng keng tarqalgan, ba'zi operatorlar 6/0 yoki 0/1 kutadi — hujjatdan tekshiriladi | `pdu.ShortCode("1234")` (doc'ida ogohlantirish bilan) |
| 4 | `+15551234567` | **1/1, lekin `+`SIZ: "15551234567"** | `+` TON=1 ma'nosining dublikati; matnda qoldirilsa RINVDSTADR xavfi | `pdu.International("+15551234567")` — konstruktor `+`ni o'zi olib tashlaydi |
| 5 | bo'sh | **0/0, addr=""** | NULL source — SMSC hisob default'ini qo'yadi (§5.2.8) | `pdu.NullSource()` |

"Yagona to'g'ri javob yo'q" bo'lganlari: #3 (short code — konventsiya) va qisman #5 (ba'zi operatorlar bo'sh source o'rniga aniq sender talab qiladi — hisob sozlamasiga bog'liq).

## Yechim 2

Zanjir: telefon SMS'ni ko'rsatganda sender sifatida TP-OA field'idagi manzilni oladi. "Bank" uchun TP-OA'da **alphanumeric** turi (TON=5 ekvivalenti) turibdi — bu GSM7 packing'dagi MATN, telefon raqamlash rejasidagi manzil EMAS. Abonent "javob berish"ni bosganda telefon yangi SMS'ning TP-DA (Destination Address) field'iga shu manzilni qo'yishga urinadi — lekin "Bank" so'zi hech qanday raqamlash rejasida yashamaydi: MSC/SMSC uni route qila olmaydi (kimga? qaysi tarmoqqa?). Ko'p telefonlar "javob berish" tugmasini bunday xabarda umuman o'chirib qo'yadi; yuborilgan taqdirda ham tarmoq rad etadi. Ya'ni muammo protokol taqiqida emas — **route qilinadigan manzilning yo'qligida**.

Ikki tomonlama xizmat uchun: **short code** (1234 — abonentlar unga yozadi, MO route operator tomonda sozlanadi) yoki oddiy **long number** (virtual MSISDN, TON=1). Ko'p mamlakatlarda "STOP" mexanizmi aynan short code orqali majburiy qilingan.

## Yechim 3

**1.** `^9989[01]` — "9989 bilan boshlanadi VA beshinchi belgisi 0 yoki 1". Mos: `998901234567` (5-belgi 0), `998911234567` (5-belgi 1). Mos emas: `998931234567` (5-belgi 3), `79261234567` (9989 bilan boshlanmaydi). E'tibor: regex'da `$` yo'q — uzunlik chegaralanmagan, faqat prefiks + bitta belgi sinfi.

**2.** Faqat **bind_receiver va bind_transceiver**'da (§5.2.7) — ya'ni MO/DLR oqimini QABUL QILADIGAN sessiyalarda. G'oya: bitta SMSC'ga bir nechta ESME ulangan bo'lsa, kelgan MO xabarning destination'iga (masalan qaysi short code'ga yozilgan) qarab qaysi ESME'ga route qilishni address_range hal qilsin. bind_transmitter'da ma'nosi yo'q — u qabul qilmaydi.

**3.** Real kutish: **hech narsa** — zamonaviy SMSC/aggregator'larning mutlaq ko'pchiligi address_range'ni ignore qiladi; ayrimlari notanish qiymatga bind'ni RAD ETADI (RBINDFAIL/RINVSYSID ko'rinishida). MO routing amalda **account (system_id) konfiguratsiyasida** hal bo'ladi: operator paneli/shartnomasida "shu short code'ning traffic'i shu accountga" deb sozlanadi. Shuning uchun default — bo'sh address_range; to'ldirish faqat operator hujjati aniq talab qilganda.
