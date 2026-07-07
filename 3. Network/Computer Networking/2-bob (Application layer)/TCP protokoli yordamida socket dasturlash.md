
UDP'dan farqli o'laroq, TCP ulanish o'rnatuvchi protokoldir. Bu shuni anglatadiki, kliyent va server bir-birlariga ma'lumot yuborishdan oldin, avval qo'l berib ko'rishadi va TCP-ulanishni o'rnatadilar. Ushbu ulanishning bir tomoni kliyent socketi bilan, ikkinchi tomoni esa server socketi bilan bog'langan. TCP-ulanishi o'rnatilganda u bilan kliyent socketining manzili (IP-manzil va port raqami) va server socketining manzili (IP-manzil va port raqami) bog'lanadi. Bir tomon o'rnatilgan TCP-ulanish yordamida ikkinchi tomonga ma'lumot yuborishga harakat qilganda, u shunchaki ularni o'z socketiga uzatadi. Bu UDP'dan farq qiladi, bunda server avval paketga manzil ma'lumotini biriktirib, keyin paketni socketga yuborishi kerak.

Endi TCP orqali ishlashda kliyent va server dasturlari qanday o'zaro ta'sir qilishini batafsil ko'rib chiqaylik. Kliyentning vazifasi server bilan aloqa o'rnatishni boshlashdir. Kliyentning dastlabki aloqasiga javob berish uchun server "tayyor" bo'lishi kerak. Tayyorlik ikki narsani anglatadi: birinchidan, UDP protokoli holatidagi kabi, TCP-server kliyent o'z xabarini yuborishga harakat qilishidan oldin jarayon sifatida ishga tushirilgan bo'lishi kerak; ikkinchidan, server dasturida maxsus "eshik", ya'ni ixtiyoriy hostda ishlaydigan kliyent jarayonidan dastlabki aloqani qabul qiladigan maxsus socket bo'lishi kerak. Uy va eshik o'xshashligimizdan foydalanib, biz ba'zan kliyentning dastlabki aloqasini "kirish eshigiga taqillatish" deb ataymiz.

Server jarayoni ishga tushganda, kliyent jarayoni server bilan TCP-ulanishni boshlashi mumkin. Bu kliyent dasturida TCP-socket yaratish orqali amalga oshiriladi. Bunda kliyent serverning kirish socketining manzilini, ya'ni server hostining IP-manzilini va socket port raqamini ko'rsatadi. O'z socketini yaratgandan keyin kliyent uch karra qo'l berib ko'rishishni boshlaydi va server bilan TCP-ulanishni o'rnatadi. Ushbu qo'l berib ko'rishish kliyent va server dasturlari uchun butunlay ko'rinmas va transport qatlamida sodir bo'ladi.

Uch karra qo'l berib ko'rishish davomida kliyent jarayoni server jarayonining kirish eshigiga taqillatadi. Server "taqillatishni eshitganda", u yangi eshik ochadi, aniqroq aytganda, shu aniq taqillatayotgan kliyent uchun mo'ljallangan yangi socket yaratadi. Quyidagi misalimizda kirish eshigi - bu biz `serverSocket` deb atagan TCP-socket; ulanish o'rnatuvchi kliyent uchun yaratilgan yangi socket esa `connectionSocket` deb ataladi.

TCP-socketlarni birinchi marta ko'rib chiqayotganlar ba'zan kirish socketi (barcha kliyentlar uchun aloqa o'rnatish boshlang'ich nuqtasi bo'lib, server bilan bog'lanishni kutayotganlar uchun) va serverning yangi yaratilgan ulanish socketi (har bir kliyent bilan bog'lanish uchun ketma-ket yaratiladi) tushunchalarini aralashtirib yuborishadi.

Dasturlar nuqtai nazaridan, kliyent socketi va server ulanish socketi to'g'ridan-to'g'ri bog'langan. Ko'rsatilganidek, kliyent jarayoni o'z socketiga ixtiyoriy miqdorda bayt yuborishi mumkin va TCP protokoli server jarayoni (ulanish socketi orqali) yuborilgan har bir baytni ma'lum tartibda olishini kafolatlaydi. Shunday qilib, TCP kliyent va server jarayonlari o'rtasida ishonchli yetkazib berish xizmatini taqdim etadi. Odamlar eshikdan ikki tomonlama o'tishi mumkin bo'lgani kabi, kliyent jarayoni ham o'z socketi orqali nafaqat ma'lumot yuborishi, balki qabul qilishi ham mumkin. Xuddi shunday, server jarayoni ham o'z ulanish socketi orqali nafaqat qabul qiladi, balki ma'lumot yuboradi ham.

Biz TCP yordamida socket dasturlash usullarini namoyish etish uchun xuddi shu oddiy kliyent-server dasturdan foydalanamiz: kliyent serverga bitta qator ma'lumot yuboradi, server qator belgilarini katta harflarga o'tkazadi va qatorni qayta kliyentga yuboradi. TCP protokolining transport xizmati orqali socket bilan bog'liq kliyent va serverning asosiy o'zaro ta'sirini ko'rsatuvchi diagramma keltirilgan.

**TCPClient.go**

Quyida dasturning kliyent qismi kodi keltirilgan:

go

```go
package main

import (
    "fmt"
    "net"
)

func main() {
    serverName := "servername"
    serverPort := "12000"
    
    clientSocket, err := net.Dial("tcp", serverName+":"+serverPort)
    if err != nil {
        panic(err)
    }
    defer clientSocket.Close()
    
    var sentence string
    fmt.Print("Kichik harflarda gap kiriting: ")
    fmt.Scanln(&sentence)
    
    _, err = clientSocket.Write([]byte(sentence))
    if err != nil {
        panic(err)
    }
    
    buffer := make([]byte, 1024)
    n, err := clientSocket.Read(buffer)
    if err != nil {
        panic(err)
    }
    modifiedSentence := string(buffer[:n])
    
    fmt.Println("Serverdan:", modifiedSentence)
}
```

Endi UDP protokoli bilan amalga oshirilgandan sezilarli darajada farq qiladigan ba'zi kod qatorlarini ko'rib chiqaylik. Birinchi farqli qator - kliyent socketini yaratish.

go

```go
clientSocket, err := net.Dial("tcp", serverName+":"+serverPort)
```

Ushbu qatorda `clientSocket` deb ataladigan kliyent socketi yaratiladi. Birinchi parametr tarmoq turini ko'rsatadi ("tcp"), ikkinchi parametr esa server manzilini belgilaydi. E'tibor bering, biz kliyent socketining port raqamini yaratishda belgilamaymiz, balki buni bizning o'rnimizda operatsion tizimga qilishga ruxsat beramiz.

UDPClient'da ko'rganlarimizdan juda farq qiladigan keyingi kod qatori:

go

```go
clientSocket, err := net.Dial("tcp", serverName+":"+serverPort)
```

Eslab qoling, kliyent va server TCP-socket yordamida bir-birlariga ma'lumot yuborishdan oldin, ular o'rtasida TCP-ulanishi o'rnatilgan bo'lishi kerak. Yuqoridagi qator aynan kliyent va server o'rtasidagi ulanishni boshlaydi. `Dial()` funksiyasining parametri ulanishning server qismining manzilidir. Ushbu kod qatori bajarilgandan keyin uch karra qo'l berib ko'rishish amalga oshiriladi va kliyent va server o'rtasida TCP-ulanishi o'rnatiladi.

go

```go
var sentence string
fmt.Print("Kichik harflarda gap kiriting: ")
fmt.Scanln(&sentence)
```

UDPClient dasturidagi kabi, bu yerda ham dastur foydalanuvchi tomonidan kiritilgan gapni oladi. `sentence` string o'zgaruvchisiga foydalanuvchi enter tugmasini bosmagunicha kiritgan barcha belgilar yoziladi. Keyingi qator ham UDPClient dasturida ko'rganlarimizdan sezilarli darajada farq qiladi:

go

```go
_, err = clientSocket.Write([]byte(sentence))
```

Ushbu qator `sentence` string o'zgaruvchisini kliyent socketi orqali TCP-ulanishga yuboradi. E'tibor bering, dastur aniq ravishda paket yaratmaydi va UDP-socketlar holatidagi kabi unga manzil ma'lumotini biriktirmaydi.

Buning o'rniga u shunchaki `sentence` qatori ma'lumotlarini TCP-ulanishga tashlaydi. Shundan keyin kliyent serverdan ma'lumot olishni kutadi.

go

```go
buffer := make([]byte, 1024)
n, err := clientSocket.Read(buffer)
modifiedSentence := string(buffer[:n])
```

Serverdan kelayotgan belgilar `modifiedSentence` string o'zgaruvchisiga joylashtiriladi. Katta harfli belgilardan iborat qatorni chop etgandan keyin biz kliyent socketini yopamiz:

go

```go
defer clientSocket.Close()
```

Oxirgi qator socketni yopadi va demak, kliyent va server o'rtasidagi TCP-ulanishni yopadi.

Aslida u kliyentdan serverga TCP-xabar yuborilishiga olib keladi (3.5-bo'limga qarang).

**TCPServer.go**

Endi dasturning server qismini ko'raylik:

go

```go
package main

import (
    "fmt"
    "net"
    "strings"
)

func main() {
    serverPort := "12000"
    
    serverSocket, err := net.Listen("tcp", ":"+serverPort)
    if err != nil {
        panic(err)
    }
    defer serverSocket.Close()
    
    fmt.Println("Server tayyor, qabul qilishga tayyormiz")
    
    for {
        connectionSocket, err := serverSocket.Accept()
        if err != nil {
            continue
        }
        
        buffer := make([]byte, 1024)
        n, err := connectionSocket.Read(buffer)
        if err != nil {
            connectionSocket.Close()
            continue
        }
        sentence := string(buffer[:n])
        
        capitalizedSentence := strings.ToUpper(sentence)
        
        _, err = connectionSocket.Write([]byte(capitalizedSentence))
        if err != nil {
            connectionSocket.Close()
            continue
        }
        
        connectionSocket.Close()
    }
}
```

Yana UDPServer va TCPClient dasturlarida keltirilganlardan sezilarli darajada farq qiladigan qatorlarni ko'rib chiqamiz. TCPClient dasturidagi kabi, server ham quyidagi qator yordamida TCP-socket yaratadi:

go

```go
serverSocket, err := net.Listen("tcp", ":"+serverPort)
```

Keyin biz server port raqamini (`serverPort` o'zgaruvchisi) socketimiz bilan bog'laymiz:

go

```go
serverSocket, err := net.Listen("tcp", ":"+serverPort)
```

Ammo TCP holatida `serverSocket` o'zgaruvchisi bizning kirish socketimiz bo'ladi. Kirish "eshigini" tayyorlagandan keyin biz unga taqillatuvchi kliyentlarni kutamiz:

go

```go
serverSocket, err := net.Listen("tcp", ":"+serverPort)
```

`Listen()` funksiyasi server TCP-ulanish so'rovlarini "tinglaydi" degan ma'noni anglatadi.

go

```go
connectionSocket, err := serverSocket.Accept()
```

Kliyent eshikka taqillatganda, dastur server socketi uchun `Accept()` metodini ishga tushiradi va u serverda yangi socket - shu aniq kliyent uchun mo'ljallangan ulanish socketini yaratadi. Keyin kliyent va socket qo'l berib ko'rishishni yakunlaydi va shu bilan `clientSocket` va `connectionSocket` o'rtasida TCP-ulanishni yaratadi, ulanish o'rnatilgandan keyin kliyent va server bir-birlari bilan ma'lumot almashishlari mumkin, bunda bir tomondan ikkinchi tomonga ma'lumotlar kafolat bilan va ularning yetkazilish tartibi kafolatlanib yetkaziladi.

go

```go
connectionSocket.Close()
```

O'zgartirilgan qatorni kliyentga yuborgandan keyin biz ulanish socketini yopamiz. Ammo server socketi ochiq qolganligi sababli, boshqa har qanday kliyent eshikka taqillatib, serverga o'zgartirish uchun yangi qator yuborishi mumkin.

Shu bilan biz TCP yordamida socket dasturlash muhokamani yakunlaymiz. Sizga ikki dasturni ikki xil hostda ishga tushirish, shuningdek ularni o'z xohishingiz bo'yicha biroz o'zgartirib ko'rish taklif etiladi. UDP uchun dasturlar juftligini TCP uchun dasturlar jufti bilan solishtiring, ular nimada farq qilishini ko'ring. Sizga kitob oxirida tasvirlangan ko'plab socket dasturlash topshiriqlarini bajarish ham taqdim etiladi. Oxir-oqibat, biz umid qilamizki, ushbu dasturlarni, shuningdek murakkabroq socket dasturlari muvaffaqiyatli yozgandan keyin, siz o'zingizning mashhur bo'ladigan, sizni boy va mashhur qiladigan dasturingizni ishlab chiqasiz va, ehtimol, ushbu kitob mualliflarini eslab qolasiz!