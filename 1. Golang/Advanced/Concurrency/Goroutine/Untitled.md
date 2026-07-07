**`Parallelism`** bir nechta ishlarni bir vaqtda mustaqil bajarish.

**`Concurrency`** vazifalar bir vaqtning o‘zida bajarilganday ko‘rinadi, lekin ular asosan o‘zaro navbat bilan bajariladi.

`Data Race` - bir nechta gorutina bir vaqtning o’zida bir o’zgaruvchiga murojat qlish natijasida yuzaga keladi.

![image.png](https://prod-files-secure.s3.us-west-2.amazonaws.com/5f55dc4a-9892-4a92-b9aa-0193e96a8298/c09d1eb2-f274-4c61-9383-9938a0a66337/image.png)

`Data Race` ni oldini olishning bir nechta usullari mavjud:

![image.png](https://prod-files-secure.s3.us-west-2.amazonaws.com/5f55dc4a-9892-4a92-b9aa-0193e96a8298/84e29aa3-ef72-4ec6-8eb8-0d01d7436ae5/image.png)

![image.png](https://prod-files-secure.s3.us-west-2.amazonaws.com/5f55dc4a-9892-4a92-b9aa-0193e96a8298/b4df8c29-34aa-4562-a1d9-6808ba2fef9f/image.png)

`Race Condition` `Deadlock` `Livelock` `Starvation`

`Context Cancellations Not Propagated`

`Unbuffered Channel Blocking`

Savol: WARNING: DATA RACE

1. Race flagi bilan dasturni run qilganimizda hatoni qanday topadi?

memory fence, или memory barrier

1. `Synchronous` – `t**hread**` ga bitta vazifa tayinlanadi va bajarish jarayoni boshlanadi. Vazifa bajarilishi tugallangach, boshqa vazifa bilan shug‘ullanish imkoniyati paydo bo‘ladi. Ushbu modelda vazifani to‘xtatib, o‘rtada boshqa vazifani bajarish imkoni yo‘q.

Keling, ushbu modelning bir oqimlik va ko‘p oqimlik senariylarda qanday ishlashini muhokama qilamiz.

Synchronous Single-Threaded Synchronous Multi-Threaded Asynchronous Single-Threaded Asynchronous Multi-Threaded

`Synchronous Single-Thread` – agar bizda bajarilishi kerak bo‘lgan bir nechta vazifamiz bo‘lsa va tizim bizga barcha vazifalar bilan ishlay oladigan bitta **thread** ni bersa, **thread** vazifalarni ketma-ketlikda birma-bir oladi va jarayon quyidagicha ko‘rinadi:

![image.png](https://prod-files-secure.s3.us-west-2.amazonaws.com/5f55dc4a-9892-4a92-b9aa-0193e96a8298/ead068dd-2047-453c-b93f-1c91ef6aed16/image.png)

Bu yerda bizda bitta **thread** (Thread 1) va bajarilishi kerak bo‘lgan 4 ta vazifa bor.

**Thread** har bir vazifani ketma-ket bajarishni boshlaydi va ularning barchasini tugatadi.

`Synchronous Multi-thread` – bir nechta **`thread`** lardan foydalanadi.

Bizda **thread pool** va bir nechta vazifalar mavjud. Demak, **thread**lar quyidagicha ishlashi mumkin:

![image.png](https://prod-files-secure.s3.us-west-2.amazonaws.com/5f55dc4a-9892-4a92-b9aa-0193e96a8298/aea31630-2505-4e1b-9b21-1ba0ebb2745c/image.png)

Bu yerda bizda 4 ta **`thread`** va bajarilishi kerak bo‘lgan 4 ta vazifa bor, har bir **`thread`** o‘z vazifasi bilan ishlashni boshlaydi, bo‘shagan **`thread`** boshqa vazifani oladi.

1. **`Asynchronous` -** sinxron dasturlash modelidan farqli o‘laroq, bunda **`thread`** bir marta vazifani bajarishni boshlagandan so‘ng, uni to‘xtatib, hozirgi holatini saqlab qo‘yishi va shu bilan bir vaqtda boshqa vazifani bajarishni boshlashi mumkin.

`Asynchronous Single-Threaded`

![image.png](https://prod-files-secure.s3.us-west-2.amazonaws.com/5f55dc4a-9892-4a92-b9aa-0193e96a8298/9ecf1ab2-a977-4c36-a184-ba3503040c75/image.png)

Bu yerda biz ko‘rishimiz mumkinki, bitta **`thread`** barcha vazifalarni bajarishga javobgar.

Agar tizimda ko‘p **`thread`** ishlash imkoniyati bo‘lsa, unda barcha **`thread` lar** quyida ko‘rsatilgandek asinxron modelda ishlashi mumkin.

`Asynchronous Multi-Threaded`

**Go'da** schedule **qanday ishlaydi?**

- **G** — gorutina
- **M** — OS oqimi (_M_ — _machine_, ya’ni mashina degani)
- **P** — CPU yadro (_P_ — _processor_, ya’ni protsessor degani)

Har bir OC `thread` operatsion tizim `schedule’i` tomonidan CPU yadrosiga (_processor_) biriktiriladi. Keyin har bir `gorutina` OC `thread` ida ishga tushiriladi.

Gorutina OC `thread`ga nisbatan soddaroq tuzilishga ega. U quyidagi ishlardan birini amalga oshirishi mumkin:

- **Bajarilmoqda (`executing`)** — OC `thread` ga `gorutina` bajarilishga biriktiriladi va undagi ko‘rsatmalar bajariladi.
- **Bajarilishga tayyor (`runnable`)** — gorutina bajarilish holatiga o‘tishni kutmoqda.
- **Kutilmoqda (`waiting`)** — gorutina to‘xtatilgan

Go’da ikki turdagi navbat mavjudi: har bir `*processor*` uchun bitta **`lokal` navbat** va barcha `*processor`* larda ishlashga mo‘ljallangan **`global` navbat**.

![image.png](https://prod-files-secure.s3.us-west-2.amazonaws.com/5f55dc4a-9892-4a92-b9aa-0193e96a8298/26e9679b-6925-44c3-a936-52acf1be3899/image.png)

Quyida `rschedule` ning ishlash tartibi:

![image.png](https://prod-files-secure.s3.us-west-2.amazonaws.com/5f55dc4a-9892-4a92-b9aa-0193e96a8298/200db1ed-381e-46fb-b93c-f8c035d6b130/image.png)

1. Har 61-chi iteratsiyada Go `schedule` `global` navbatdan gorutina mavjudligini tekshiradi.
2. Agar bo'lmasa, u o'zining lokal navbatini tekshiradi.
3. Agar global va lokal navbatlar bo'sh bo'lsa, `schedule` boshqa lokal navbatlardan gorutinani olishi mumkin. Bu `schedule` da **"ishni o'g'irlash"**(`work stealing`) deb ataladi.

При реализации паттерна Worker Pool мы увидели, что оптимальное количество горутин в пуле зависит от типа рабочей нагрузки. Если рабочая нагрузка, вы- полняемая рабочими процессами, является I/O-bound, то это значение зависит от внешней системы. И наоборот, если рабочая нагрузка будет типа CPU-bound, то оптимальное количество горутин близко к количеству доступных потоков. Знание типа рабочей нагрузки (I/O- или CPU-bound) очень важно при разра- ботке конкурентных приложений.

В чем смысл вызова функции cancel как функции defer? Внутри себя context. WithTimeout создает горутину, которая будет храниться в памяти в течение 4 секунд или до тех пор, пока не будет вызвана cancel. Следовательно, вызов cancel в каче- стве функции defer означает, что при выходе из родительской функции контекст будет отменен, а созданная горутина остановлена. Это мера предосторожности, чтобы при возвращении мы не оставили в памяти сохраненные объекты.