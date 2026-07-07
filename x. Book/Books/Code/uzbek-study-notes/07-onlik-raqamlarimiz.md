# 7-bob: O'nlik raqamlarimiz

## Bobning asosiy g'oyasi

Bu bobda muallif biz har kuni ishlatadigan 0 dan 9 gacha bo'lgan raqamlarni oddiy va o'zgarmas narsa deb emas, balki juda kuchli kod sifatida ko'rsatadi. Til kabi sonlarni yozish usuli ham kelishuvga tayangan belgilar tizimidir. Farqi shundaki, sonlar tildan ham mavhumroq ko'rinadi: "3" belgisi uchta olma, uchinchi kanal, yosh, hisob, retseptdagi miqdor yoki boshqa ko'plab ma'nolarni bildirishi mumkin.

Muallifning asosiy fikri shunday: bizning o'nlik sanoq tizimimiz tabiat qonuni emas, balki tarixiy va amaliy tanlovdir. Odamlar qadimdan barmoqlar bilan sanagani uchun 10 soni atrofida tizim qurish qulay bo'lgan. Ammo hisoblashning haqiqiy kuchi 10 sonining o'zida emas, balki pozitsion yozuv, nol va raqamlarning o'rniga qarab qiymat olishi kabi g'oyalardadir. Shu g'oya keyingi boblarda boshqa asosli sanoq tizimlarini, ayniqsa ikkilik tizimni tushunishga tayyorlaydi.

## Bosqichma-bosqich tushuntirish

### 1. Sonlar ham koddir

Avvalgi boblarda kod deganda harflar, nuqta-tirelar, elektr signallari va belgilar haqida gap ketgan edi. Bu bobda muallif shu fikrni sonlarga ko'chiradi. Biz ko'pincha "3" ni go'yo bevosita uchlikning o'zi deb qabul qilamiz. Aslida esa "3" - ma'lum bir miqdorni yozish uchun ishlatiladigan belgi.

Masalan, "uchta olma" miqdorini turli yo'llar bilan ifodalash mumkin:

- uchta olma rasmini chizish;
- uchta chiziq tortish: `|||`;
- "uch" deb yozish;
- `3` raqamini yozish;
- boshqa sanoq tizimida boshqa ko'rinish bilan yozish.

Demak, sonning ma'nosi va sonning yozilishi bir xil narsa emas. Yozuv - kod, ma'no esa miqdor.

### 2. 10 soni sehrli emas

Biz 10, 100, 1000, million, milliard kabi sonlarga alohida ahamiyat beramiz. Chunki o'nlik tizimda ular juda "dumaloq" ko'rinadi. Lekin muallif o'quvchini muhim savolga olib keladi: nega aynan 10?

Eng sodda javob - odam qo'llari. Ko'p jamiyatlarda odamlar sanashni barmoqlar orqali o'rgangan. Ikki qo'lda odatda o'nta barmoq bo'lgani uchun o'nlik tizim tabiiy va qulay ko'ringan. Agar odamzodda sakkizta yoki o'n ikkita barmoq bo'lganida, ehtimol bizning "dumaloq" sonlarimiz ham boshqacha bo'lardi.

Shu sababli 10 asosli tizimni mutlaq yagona yo'l deb emas, balki amaliy kelishuv deb ko'rish kerak. Bu fikr kompyuterlarni tushunishda juda muhim: kompyuterlar biz kabi 10 ta raqamga muhtoj emas, ular ikki holatli signallar bilan ham sonlarni ifodalay oladi.

### 3. Dastlabki sanash: rasm va chiziqlardan raqamlarga

Muallif sonlarning kelib chiqishini oddiy ehtiyojdan boshlaydi: odamlar narsalarni sanashi kerak bo'lgan. Qancha mol bor, nechta odam ishlayapti, savdoda qancha buyum almashildi, qarz qancha - bularning barchasi yozib qo'yishni talab qiladi.

Eng sodda usul - narsaning o'zini chizish. To'rtta o'rdakni hisobga olish uchun to'rtta o'rdak chizish mumkin. Ammo bu tezda noqulaylashadi. O'rdak ko'paysa, har birini chizish ortiqcha mehnatga aylanadi. Shunda fikr o'zgaradi: narsaning rasmini bir marta ko'rsatib, miqdorni alohida belgi bilan berish mumkin.

Bu o'zgarish juda muhim:

- rasm narsaning turini bildiradi;
- chiziqlar yoki raqamlar miqdorni bildiradi;
- miqdorni alohida kodlash hisoblashni osonlashtiradi.

Shu yo'l bilan insoniyat "narsani chizish"dan "miqdorni belgilash"ga o'tadi.

### 4. Rim raqamlari: ishlaydi, lekin chegarasi bor

Qadimgi sanoq tizimlari ichida rim raqamlari hozirgacha eng tanishlaridan biri. Ular soatlarda, yodgorliklarda, kitob sahifalarida, kinolar oxiridagi mualliflik yili yozuvlarida uchraydi.

Rim raqamlari asosiy belgilar to'plamiga tayanadi:

```text
I   V   X   L   C   D   M
1   5   10  50  100 500 1000
```

Bu tizimda son ko'proq belgilarni yig'ish orqali tuziladi. Masalan, 27 ni tasavvur qilish uchun 10 + 10 + 5 + 1 + 1 kabi qismlarni ketma-ket qo'yish mumkin. Bunday yozuv qo'shish va ayrishni ma'lum darajada qulay qilgan: belgilarni birlashtirib, keyin beshta `I` ni `V` ga, ikkita `V` ni `X` ga almashtirish kabi qoidalar ishlaydi.

Ammo muammo murakkab amallarda ko'rinadi. Rim raqamlari bilan ko'paytirish yoki bo'lish oson emas. Tizimda belgilar sonning ichki tuzilishini bizning hozirgi yozuvimiz kabi qulay ochib bermaydi. Shuning uchun u hisob-kitob uchun emas, ko'proq yozuv va belgilash uchun yashab qolgan.

### 5. Hind-arab raqamlari va katta burilish

Bugun ishlatadigan `0 1 2 3 4 5 6 7 8 9` raqamlari ko'pincha hind-arab raqamlari deb ataladi. Ular Hindistonda paydo bo'lgan, keyin arab matematiklari orqali Yevropaga tarqalgan. Muallif Al-Xorazmiyni ham eslatadi: uning algebra haqidagi ishlari Yevropaga katta ta'sir qilgan, "algoritm" so'zi ham uning nomi bilan bog'liq.

Bu tizimning kuchi raqamlarning shaklida emas. Kuch quyidagi uch fikrda:

- raqamning qiymati uning turgan joyiga bog'liq;
- o'n uchun alohida belgi shart emas;
- nol bo'sh o'rinni ko'rsatib, butun tizimni ushlab turadi.

Hind-arab yozuvi sonlarni qisqa, tartibli va hisoblashga qulay shaklda beradi. Aynan shu sababli u rim raqamlaridan ustun chiqdi.

### 6. Pozitsion yozuv: raqam joyiga qarab o'zgaradi

Pozitsion tizimda raqamning ma'nosi faqat o'z shakliga emas, balki qaysi ustunda turganiga ham bog'liq. `5` belgisi birlar xonasida 5 ni, o'nlar xonasida 50 ni, yuzlar xonasida 500 ni bildiradi.

Masalan, `4825` sonini bunday o'qish mumkin:

```text
4825 = 4 minglik + 8 yuzlik + 2 o'nlik + 5 birlik
     = 4 x 1000 + 8 x 100 + 2 x 10 + 5 x 1
```

Bu yerda har bir xona 10 ning darajasiga mos keladi:

```text
10^3  10^2  10^1  10^0
1000  100   10    1
  4     8    2    5
```

Shuning uchun `1`, `10`, `100`, `1000` yozuvlarida bir xil `1` belgisi turli qiymatlarni bildiradi. Qaysi qiymatni bildirishi uning o'ngdan nechanchi o'rinda turganiga bog'liq.

### 7. Nol: "hech narsa"dan ko'proq narsa

Nol oddiy "yo'q" belgisi emas. Bu bobda nol pozitsion yozuvning tayanchi sifatida ko'rsatiladi. U bo'sh xonani belgilaydi va raqamlarning qaysi o'rinda turganini aniq saqlaydi.

Masalan:

```text
25   = 2 o'nlik + 5 birlik
205  = 2 yuzlik + 0 o'nlik + 5 birlik
250  = 2 yuzlik + 5 o'nlik + 0 birlik
```

Uchala yozuvda ham `2` va `5` bor, lekin nol ularning joylashuvini o'zgartirib, butun qiymatni boshqa qiladi. Agar nol bo'lmaganida, `25`, `205` va `250` kabi sonlarni bir-biridan ajratish juda qiyin bo'lardi.

Nolning yana bir kuchi shundaki, u hisoblash qoidalarini soddalashtiradi. U bo'sh o'rinni ko'rsatgani uchun katta sonlar ustida qo'shish, ko'paytirish va bo'lishni ustunlar bo'yicha bajarish mumkin bo'ladi.

### 8. O'nlik kasrlar ham shu qoidaga bo'ysunadi

Muallif pozitsion tizim faqat vergulgacha bo'lgan butun sonlarda emas, verguldan keyingi kasrlarda ham ishlashini ko'rsatadi. Verguldan chap tomonda 10 ning musbat darajalari bor: birlik, o'nlik, yuzlik, minglik. Verguldan o'ng tomonda esa 10 ning manfiy darajalari bor: o'ndan bir, yuzdan bir, mingdan bir.

Masalan, `42705.684` sonini shunday o'qish mumkin:

```text
4 x 10000
2 x 1000
7 x 100
0 x 10
5 x 1
6 x 0.1
8 x 0.01
4 x 0.001
```

Darajalar bilan esa umumiy naqsh yanada ravshanlashadi:

```text
10^4  10^3  10^2  10^1  10^0 . 10^-1  10^-2  10^-3
  4     2     7     0     5       6       8       4
```

Demak, pozitsion tizimda vergul faqat chegarani bildiradi: chapda darajalar 0 ga qarab kamayadi, o'ngda esa manfiy darajalarga o'tadi.

### 9. Nega arifmetika jadval bilan ishlaydi?

O'nlik tizimning go'zalligi shundaki, katta sonlar ustidagi amallar kichik amallarga bo'linadi. `300 + 400 = 700`, chunki aslida bir xonada `3 + 4 = 7` qoidasi ishlayapti. `3000 + 4000` ham shu fikrga tayanadi: faqat xona qiymati boshqa.

Qo'shish jadvali va ko'paytirish jadvali shuning uchun kerak bo'lgan. Biz barcha katta sonlarni yodlamaymiz. Biz bir xonali raqamlar bilan ishlashni o'rganamiz, keyin shu kichik qoidalarni har bir ustunga qo'llaymiz.

Masalan:

```text
  386
+ 247
-----
  633
```

Bu yerda har bir ustunda bir xonali qo'shish ishlaydi:

- birliklarda: 6 + 7 = 13, 3 yoziladi, 1 keyingi xonaga o'tadi;
- o'nliklarda: 8 + 4 + 1 = 13, 3 yoziladi, 1 keyingi xonaga o'tadi;
- yuzliklarda: 3 + 2 + 1 = 6.

Ko'paytirishda ham shunga o'xshash: katta amal bir xonali ko'paytirishlar va xona bo'yicha siljitishlardan tuziladi.

### 10. Eng muhim xulosa: pozitsion tizim boshqa asoslarga ham o'tadi

Bob oxirida muallif keyingi katta g'oyaga eshik ochadi. Pozitsion yozuv faqat 10 asosida ishlaydigan usul emas. Agar sanoq tizimi 8 asosli bo'lsa, unda ham xonalar bo'ladi, faqat ular 10 ning darajalari emas, 8 ning darajalari bo'ladi.

O'nlik tizimda xonalar:

```text
... 1000  100  10  1
```

Sakkizlik tizimda xonalar:

```text
... 512   64   8   1
```

Bu fikr nihoyatda muhim, chunki kompyuterlar keyinroq 2 asosli tizim bilan ishlaydi. Ya'ni asos o'zgarsa ham pozitsion yozuv g'oyasi qoladi: har bir xona asosning navbatdagi darajasini bildiradi.

## Original vizual diagramma

Quyidagi diagramma bobdagi asosiy tushunchani yangi shaklda ko'rsatadi: bitta raqam belgisi joyiga qarab boshqa qiymatga aylanadi, nol esa bo'sh joyni ushlab turadi.

```text
Pozitsion yozuvning skeleti

            chapga yurish: qiymat asosga ko'payadi
        <---------------------------------------------

   [10^4]   [10^3]   [10^2]   [10^1]   [10^0] . [10^-1] [10^-2]
   10000     1000      100       10        1       0.1     0.01
      |        |        |        |        |        |        |
      4        2        7        0        5   .    6        8

            --------------------------------------------->
            o'ngga yurish: qiymat asosga bo'linadi

Nolning roli:
42705.68 sonida 0 o'nliklar xonasini egallab turadi.
U "bu yerda o'nlik yo'q, lekin keyingi raqamlar joyidan siljimasin" deydi.
```

## Muhim tushunchalar

- Raqam: sonni yozish uchun ishlatiladigan belgi.
- Son: miqdor yoki tartib haqidagi mavhum tushuncha; u turli belgilar bilan ifodalanishi mumkin.
- Sanoq tizimi: sonlarni yozish va o'qish qoidalari to'plami.
- Asos: pozitsion tizimda har bir xona necha marta kattalashishini belgilaydigan son; o'nlik tizimda asos 10.
- O'nlik tizim: 0 dan 9 gacha bo'lgan o'nta raqamga tayangan, asosi 10 bo'lgan pozitsion tizim.
- Pozitsion yozuv: raqam qiymati uning son ichidagi o'rniga bog'liq bo'lgan yozuv usuli.
- Xona qiymati: birlik, o'nlik, yuzlik, minglik kabi raqam turgan o'rinning qiymati.
- Daraja: asosning necha marta o'ziga ko'paytirilishini bildiradi; masalan, 10^3 = 1000.
- Nol: bo'sh xonani belgilovchi va pozitsion yozuvni aniq saqlovchi raqam.
- Rim raqamlari: I, V, X, L, C, D, M belgilariga tayangan qadimgi sanoq yozuvi.
- Hind-arab raqamlari: hozir ishlatiladigan 0-9 raqamlari va pozitsion yozuvga asoslangan tizim.
- O'nlik kasr: verguldan keyin 10 ning manfiy darajalari orqali ifodalanadigan son qismi.

## Kichik misol

`3048` sonini olaylik. Unda `3`, `0`, `4`, `8` raqamlari bor. Bu oddiy ketma-ketlik emas, har bir raqamning o'rni bor.

```text
3048 = 3 x 1000 + 0 x 100 + 4 x 10 + 8 x 1
     = 3000 + 0 + 40 + 8
     = 3048
```

Bu yerda nol juda muhim. Agar nolni olib tashlasak, `348` hosil bo'ladi. Bu esa boshqa son:

```text
348 = 3 x 100 + 4 x 10 + 8 x 1
```

Demak, `3048` dagi `0` "hech narsa qo'shmayapti"dek ko'rinadi, lekin aslida `3` ning minglik xonasida qolishini ta'minlayapti. Nol qiymat qo'shmaydi, lekin joyni saqlaydi.

## O'zini tekshirish savollari

1. Nega muallif sonlarni ham kod deb qarashni taklif qiladi?
2. `3` belgisi bilan "uchta narsa" tushunchasi orasidagi farq nimada?
3. Nima uchun 10 asosli tizim tarixan qulay bo'lgan, lekin majburiy yagona tizim emas?
4. Rim raqamlari qo'shish uchun ma'lum darajada ishlasa ham, nega murakkab arifmetikada noqulay?
5. Hind-arab raqamlarining eng katta ustunligi raqam shakllaridami yoki yozuv qoidalaridami?
6. Pozitsion yozuvda `7` raqami qachon 7, qachon 70, qachon 700 bo'ladi?
7. Nol `25`, `205` va `250` sonlarini qanday farqlantiradi?
8. Verguldan keyingi raqamlar nima uchun 10 ning manfiy darajalari bilan bog'lanadi?
9. Katta sonlarni qo'shish va ko'paytirish nega bir xonali amallarga bo'linadi?
10. Pozitsion yozuv g'oyasi nima uchun keyingi boblarda boshqa sanoq tizimlarini tushunishga yordam beradi?

## Qisqa xulosa

7-bob bizga juda tanish bo'lgan o'nlik raqamlarni begona ko'z bilan qayta ko'rsatadi. Sonlarni yozish usuli tabiiy haqiqat emas, balki koddir. O'nlik asos barmoqlar sababli qulay bo'lgan tarixiy tanlov, lekin hisoblashning haqiqiy kuchi pozitsion yozuv va noldadir. Raqam joyiga qarab qiymat oladi, nol esa bo'sh joyni ushlab turadi. Shu sababli katta sonlar, kasrlar va arifmetik amallar tartibli va sodda qoidalarga bo'ysunadi. Bobning eng muhim tayyorgarligi shuki: agar pozitsion tizim 10 asosida ishlasa, u boshqa asoslarda ham ishlashi mumkin.
