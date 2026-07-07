> Process — bu qaysidir dasturiy taminotning komputer ichida bajarilayotgan jarayoni.

Bu nima degani masalan siz word dasturini ishga tushurdingiz bu komputerda word processni yaratadi va siz qilgan ishlar uning ichida amalga oshadi. **Process**'lar kamida bitta(_**main thread**_)yoki odatda bir necha **thread**'lardan tashkil topgan bo'ladi.

> Thread — bu process ichidagi element bo'lib dasturni haqiqatdan ishlashiga ya’ni siz buyurgan amallarni natijaga aylantirishga yordam beradi.

Processlar bir biri bilan xotira ulashmaydi va bu agar bir process ishdan chiqsa yoki xato ishlashni boshlasa boshqa processlarga tasir qilmaydi. Process thread natijalarini boshqaradi. Threadlar bir biri bilan xotira ulashadi va agarda birortasi noto'g'ri ishlashni boshlasa hammasiga tasir qiladi.

Biz bir vaqtda kompyuterda bir necha dasturlarni ishlatamiz bunda kompyuter bularni qanday boshqaradi?.Buning uchun bizga **scheduler** yordam beradi.Scheduler processdagi har bir **thread**ga ishlashi uchun keraklicha _**(teng emas)**_ vaqt berishga bir necha har xil algoritmlardan foydalanadi. Buning yordamida biz 4 corelik computerda ham har bir coreda 1 ta thread ishlata olishimizga qaramay ko`plar dasturlarni ishalata olamiz.Bu holat **context switching** deyiladi .Ya’ni 1 ta coredagi threadlar vaqt o`tishi bilan almashinadi bular ketma-ketlik shaklida bo`lmaydi muhumliligiga qarab o`sha thread oldinroq ishga tushuriladi.