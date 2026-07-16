# Amaliyot — 100 kunlik challenge

Asosiy fayl: **[100-kunlik-challenge.md](100-kunlik-challenge.md)** — 16 blok, 100 kun, ~110 LeetCode misol.

## Qoidalar

1. **Har kun bitta katak.** Kun yakunida `- [ ]` ni `- [x]` ga o'zgartir — shu tracking.
2. **Avval nazariya, keyin misol.** Har blok boshidagi 📖 havoladagi faylni o'qimasdan misolga o'tma.
3. **30 daqiqa qoidasi.** Misol ustida 30 daqiqa o'ylab yo'l topilmasa — hint yoki yechimni o'rgan, lekin keyin **yechimga qaramay o'zing qayta yoz**. Bu mag'lubiyat emas, o'rganish usuli.
4. **Yechimdan keyin 2 savol:** time/space complexity qancha? Boshqa qanday yechim bor edi?
5. **Kun o'tkazib yuborilsa** — kechagi kunni bugun qil, jadvalni surish mumkin. 100 kun ketma-ketligi emas, 100 ta bajarilgan kun muhim.
6. **Takrorlash kunlari (🔁) muqaddas.** Ular esdan chiqarishga qarshi eng kuchli qurol (spaced repetition) — o'tkazib yuborma.

## Yechimlarni saqlash

Yechimlarni `yechimlar/` papkasida mavzu bo'yicha saqla:

```
yechimlar/
├── two-pointers/
│   ├── 344-reverse-string.go
│   └── 15-3sum.go
├── linked-list/
├── graph/
└── ...
```

Fayl nomi: `<masala-raqami>-<nomi>.go`. Har faylning boshiga bitta izoh: yondashuv va complexity.

## Kunlik misollar soni qanday belgilangan?

- **Easy** misollar — kuniga 2 tagacha (ular bir texnikaning ikki ko'rinishi)
- **Medium** misollar — kuniga 1 ta (chuqur o'ylashni talab qiladi)
- Katta bloklarga (Two Pointers — 12, Binary Search — 10, Tree — 10) ko'proq kun ajratilgan
- Har 2-3 blokda 🔁 takrorlash kuni bor

## Progress hisobi

Nechta kun bajarilganini tez ko'rish uchun:

```bash
grep -c "\- \[x\]" 100-kunlik-challenge.md
```
