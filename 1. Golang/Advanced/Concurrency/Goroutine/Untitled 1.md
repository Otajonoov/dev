**==Goroutine==** – bu Go tilining runtime tomonidan boshqariladigan yengil thread.

Go tilida quyidagi atamalar qo'llaniladi:

- **G (Goroutine)** – Gorutina
- **M (Machine)** – Mashina
- 
 Каждая Машина работает в отдельном потоке и способна выполнять только одну Горутину в момент времени. Планировщик операционной системы, в которой работает программа, переключает Машины. Число работающих Машин ограничено переменной среды `GOMAXPROCS` или функцией `runtime.GOMAXPROCS(n int)`. По умолчанию она равна количеству ядер процессора компьютера, на которой было запущено приложение.

Har bir **Mashina (M)** alohida thread da ishlaydi va bir vaqtning o'zida faqat bitta **gorutina** ni bajarishi mumkin. Dastur ishlayotgan operatsion tizimning **scheduleri** **Mashinalar** ni almashtiradi. Ishchi **Mashinalar** soni `GOMAXPROCS` o‘zgaruvchisi yoki `runtime.GOMAXPROCS(n int)` funksiyasi bilan cheklangan. Odatiy holatda bu kompyuter protsessorining yadrolari soniga teng bo‘ladi.**

Funksiyani **gorutina** sifatida ishga tushirish uchun `go func()` yozish kifoya, bu yerda `func()` – siz ishga tushirmoqchi bo‘lgan funksiya.

### Qisqacha tushuntirish:

Gorutinalar – Go tilida parallel va samarali ishlash uchun ishlatiladigan "yengil" threadlar. Ular asosiy dasturga ta'sir qilmasdan bir-biridan mustaqil ishlashi mumkin. `go` kalit so‘zi orqali funksiyani gorutina sifatida ishga tushirish mumkin.