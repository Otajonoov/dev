Tranzaksiya - bu ma'lumotlar bazasi ustida bajariladigan operatsiyalar to'plami.

Tranzaksiyalar ma'lumotlar bazasining muvofiqligini (ziddiyatsizligini) ta'minlash vositalaridan biri hisoblanadi, jadvallarga qo'yiladigan yaxlitlik cheklovlari (constraints) bilan birga.

Tranzaksiya ikki natijaga ega bo'lishi mumkin: 
	1. Tranzaksiya bajarish jarayonida amalga oshirilgan ma'lumotlar o'zgarishlari muvaffaqiyatli ma'lumotlar bazasida fiksatsiya qilindi
	2. Tranzaksiya bekor qilinadi va uning doirasida bajarilgan barcha o'zgarishlar bekor qilinadi. Tranzaksiyani bekor qilish rollback (ortga qaytish) deb ataladi.

DBMS lar tranzaksiyalarni parallel, ya'ni bir vaqtda bajarish uchun maxsus mexanizmlarni taklif qiladi. Bunday mexanizmlar PostgreSQL da ham amalga oshirilgan.

PostgreSQL DBMS da tranzaksiyalarni amalga oshirish ko'p versiyali model (Multiversion Concurrency Control, MVCC) ga asoslangan. Ushbu model har bir SQL operator ma'lumotlarning snapshot (surati) deb ataladigan, ya'ni ma'lumotlar bazasining ma'lum bir vaqt momentidagi muvofiq holati (versiyasi) ni ko'rishini nazarda tutadi. Bunda parallel bajariladigan tranzaksiyalar, hatto ma'lumotlar bazasiga o'zgarishlar kiritayotganlar ham, ushbu suratdagi ma'lumotlarning muvofiqligini buzmaydi. Bunday natija PostgreSQL da parallel tranzaksiyalar bir xil jadval qatorlarini o'zgartirganda, ushbu qatorlarning alohida versiyalari yaratilishi va mos tranzaksiyalar uchun mavjud bo'lishi hisobiga erishiladi. Bu ma'lumotlar bazasi bilan ishlashni tezlashtiradi, lekin ko'proq disk maydoni va operativ xotira talab qiladi. Va MVCC ni qo'llashning yana bir muhim natijasi - o'qish operatsiyalari hech qachon yozish operatsiyalari tomonidan bloklanmaydi, yozish operatsiyalari esa hech qachon o'qish operatsiyalari tomonidan bloklanmaydi.

Ma'lumotlar bazalari nazariyasiga ko'ra, tranzaksiyalar quyidagi xususiyatlarga ega bo'lishi kerak. Ushbu to'rtta xususiyatni belgilash uchun ACID qisqartmasi ishlatiladi.

1. **Atomicity** - all or nothing
2. **Consistency**. Bu xususiyat tranzaksiyani muvaffaqiyatli bajarish natijasida ma'lumotlar bazasi bir muvofiq holatdan boshqa muvofiq holatga o'tkazilishini talab qiladi.
3. **Isolation**. Tranzaksiya bajarilishi vaqtida boshqa tranzaksiyalar imkon qadar minimal ta'sir ko'rsatishi kerak.
4. **Durability**. Tranzaksiyani muvaffaqiyatli fiksatsiya qilgandan keyin foydalanuvchi ma'lumotlarning ma'lumotlar bazasida ishonchli saqlanganligiga va keyinchalik tizimning mumkin bo'lgan nosozliklaridan qat'i nazar, undan olinishi mumkinligiga ishonch hosil qilishi kerak.

SQL standartida jami to'rtta daraja nazarda tutilgan. Har bir yuqori daraja oldingi darajaning barcha imkoniyatlarini o'z ichiga oladi. Default holda PostgreSQL Read Committed izolyatsiya darajasidan foydalanadi.

1. **Read Uncommitted**. Bu eng past izolyatsiya darajasi. SQL standartiga ko'ra bu darajada "iflos" (fiksatsiya qilinmagan) ma'lumotlarni o'qishga yo'l qo'yiladi. Biroq PostgreSQL da bu darajaga qo'yiladigan talablar standartdagidan qattiqroq: bu darajada "iflos" ma'lumotlarni o'qishga yo'l qo'yilmaydi.
2. **Read Committed**. "Iflos" (fiksatsiya qilinmagan) ma'lumotlarni o'qishga yo'l qo'yilmaydi. Shunday qilib, PostgreSQL da Read Uncommitted darajasi Read Committed darajasi bilan mos tushadi. Tranzaksiya faqat o'zi bajarish jarayonida amalga oshirgan fiksatsiya qilinmagan ma'lumotlar o'zgarishlarini ko'rishi mumkin.
3. **Repeatable Read**. "Iflos" (fiksatsiya qilinmagan) ma'lumotlarni o'qish va takrorlanmaydigan o'qishga yo'l qo'yilmaydi. PostgreSQL da bu darajada fantom o'qishga ham yo'l qo'yilmaydi. Shunday qilib, bu darajani amalga oshirish SQL standartida talab qilinganidan qattiqroqdir. Bu standartga zid emas.
4. **Serializable**. Yuqorida sanab o'tilgan fenomenlarning hech biriga, shu jumladan serializatsiya anomaliyasiga ham yo'l qo'yilmaydi.

---

