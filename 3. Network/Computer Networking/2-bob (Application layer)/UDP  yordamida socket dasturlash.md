

Ushbu bo'limda biz UDP protokolidan foydalanuvchi oddiy kliyent-server dasturini yozamiz; keyingi bo'limda esa TCP'dan foydalanuvchi shunga o'xshash dasturni ko'rib chiqamiz.

2.1-bo'limdan eslaymizki, turli qurilmalarda ishlaydigan jarayonlar bir-birlari bilan socketlarga xabarlar yuborish orqali aloqa qilishadi. Biz har bir jarayonni uyga, jarayonning socketini esa eshikka o'xshatgan edik. Dastur eshikning bir tomonida - uy ichida ishlaydi; transport qatlamining protokoli esa eshikning narigi tomonida - tashqi dunyoda joylashgan. Dastur ishlab chiquvchisi dastur qatlami tomonidagi hamma narsani nazorat qiladi, ammo transport qatlami tomonida u kamroq ta'sir ko'rsata oladi.

Endi UDP socketlardan foydalanuvchi ikki jarayon o'rtasidagi o'zaro ta'sirni batafsil ko'rib chiqaylik. Jo'natuvchi jarayon o'zining socketi orqali ma'lumotlar paketini yuborishdan oldin, u avval unga manzil ma'lumotini biriktirib qo'yishi kerak. Paket jo'natuvchining socketi orqali o'tgach, Internet ushbu manzil ma'lumotidan foydalanib, paketni qabul qiluvchi jarayonning socketiga yo'naltiradi. Paket qabul qiluvchining socketiga yetgach, qabul qiluvchi jarayon uni socket orqali olib, paket mazmunini tekshiradi va tegishli harakatlarni amalga oshiradi.

Albatta, siz paketga biriktiladigan manzil ma'lumoti nimalardan iborat ekanligini so'rashingiz mumkin. Taxmin qilish qiyin emaski, uning bir qismi manzil hostining IP-manzilidhan iborat. Bu Internet'dagi routerlarga paketni tarmoq bo'ylab kerakli manzilga yo'naltirish imkonini beradi. Ammo hostda bir nechta tarmoq dasturlari jarayonlari ishlaishi mumkin va ularning har biri bir yoki bir nechta socketlardan foydalanishi mumkinligi sababli, manzil hostida ma'lum bir socketni aniqlash zarur. Socket yaratilganda unga port raqami deb ataladigan identifikator beriladi. Shuning uchun manzil ma'lumoti socketning port raqamini ham o'z ichiga oladi. Xulosa qilib aytganda, jo'natuvchi jarayon paketga qabul qiluvchi hostning IP-manzili va qabul qiluvchi socketning port raqamidan iborat manzil ma'lumotini biriktirib qo'yadi.

Bundan tashqari, tez orada ko'rib chiqadigan kabi, jo'natuvchi hostning IP-manzili va jo'natuvchi socketning port raqamidan iborat jo'natuvchi manzili ham paketga biriktiriladi. Ammo bu odatda UDP-dasturi kodida emas, balki jo'natuvchi hostning operatsion tizimi tomonidan avtomatik ravishda amalga oshiriladi.

UDP va TCP yordamida socket dasturlash usullarini namoyish etish uchun biz quyidagi oddiy kliyent-server dasturdan foydalanamiz:

1. Kliyent klaviaturadan belgilar qatorini (ma'lumotlarni) o'qiydi va ma'lumotlarni serverga yuboradi.
2. Server ma'lumotlarni qabul qilib, belgilarni katta harflarga o'tkazadi.
3. Server o'zgartirilgan ma'lumotlarni kliyentga yuboradi.
4. Kliyent o'zgartirilgan ma'lumotlarni qabul qilib, qatorni o'z ekranida ko'rsatadi.

UDP protokoli orqali o'zaro ta'sir qiluvchi kliyent va server socketlarining asosiy harakatlarini ko'rsatuvchi diagramma berilgan.

Endi UDP protokolidan foydalanib ushbu oddiy dasturni amalga oshiruvchi kliyent-server dasturlar juftligiga qaraylik. Har bir dasturning batafsil qator-ba-qator tahlilini amalga oshiramiz. UDP-kliyentdan boshlaymiz, u serverga oddiy dastur qatlami xabarini yuboradi. Serverning kliyent xabarlarini qabul qilib, ularga javob bera olishi uchun u "tayyor" bo'lishi kerak - ya'ni kliyent o'z xabarini yuborishdan oldin jarayon sifatida ishga tushirilgan bo'lishi lozim.

Kliyent dasturi UDPClient.go, server dasturi esa UDPServer.go deb ataladi. Aslida "to'g'ri kod" xatolarni qayta ishlash kabi ko'plab qatorlarni o'z ichiga oladi, ammo biz asosiy nuqtalarni ta'kidlash uchun uni ataylab iloji boricha qisqa qildik. Bizning dasturimiz uchun biz serverning ixtiyoriy 12000 port raqamini tanladik.

**UDPClient.go**

Quyida dasturning kliyent qismi kodi keltirilgan:

go

```go
package main

import (
    "fmt"
    "net"
)

func main() {
    serverName := "hostname"
    serverPort := "12000"
    
    serverAddr, err := net.ResolveUDPAddr("udp", serverName+":"+serverPort)
    if err != nil {
        panic(err)
    }
    
    clientSocket, err := net.DialUDP("udp", nil, serverAddr)
    if err != nil {
        panic(err)
    }
    defer clientSocket.Close()
    
    var message string
    fmt.Print("Kichik harflarda gap kiriting: ")
    fmt.Scanln(&message)
    
    _, err = clientSocket.Write([]byte(message))
    if err != nil {
        panic(err)
    }
    
    buffer := make([]byte, 2048)
    n, serverAddr, err := clientSocket.ReadFromUDP(buffer)
    if err != nil {
        panic(err)
    }
    modifiedMessage := string(buffer[:n])
    
    fmt.Println(modifiedMessage)
}
```

Endi ushbu dastur kodining turli qatorlarini ko'rib chiqaylik.

go

```go
import (
    "fmt"
    "net"
)
```

`net` paketi barcha tarmoq aloqalarining asosi hisoblanadi. Ushbu qatorni kiritib, biz dasturimiz ichida socketlar yarata olamiz.

go

```go
serverName := "hostname"
serverPort := "12000"
```

Birinchi qatorda biz `serverName` o'zgaruvchisiga `"hostname"` qator qiymatini belgilaymiz. Bu yerda biz server IP-manzilini (masalan, "128.138.32.126") yoki server host nomini (masalan, "cis.poly.edu") o'z ichiga olgan qatorni qo'yishimiz kerak. Host nomi ishlatilganda u avtomatik ravishda IP-manzilga aylantiriladi. Ikkinchi qatorda biz `serverPort` string o'zgaruvchisining qiymatini "12000" ga o'rnatamiz.

go

```go
serverAddr, err := net.ResolveUDPAddr("udp", serverName+":"+serverPort)
clientSocket, err := net.DialUDP("udp", nil, serverAddr)
```

Ushbu qatorlarda biz kliyent socketini yaratamiz va uni `clientSocket` deb ataymiz. `ResolveUDPAddr` funksiyasi UDP manzilini yaratadi, `DialUDP` esa UDP ulanishini ochadi. Go'da biz kliyent socketining port raqamini yaratishda belgilamaymiz - operatsion tizimga buni bizning o'rnimizda qilishga ruxsat beramiz. Endi kliyent jarayonining "eshigi" yaratilgandi, biz xabarlar yaratib, ular orqali yuborishimiz mumkin.

go

```go
var message string
fmt.Print("Kichik harflarda gap kiriting: ")
fmt.Scanln(&message)
```

`fmt.Scanln()` Go'ning o'rnatilgan funksiyasi. U bajarilganda kliyent tomonidagi foydalanuvchiga "Kichik harflarda gap kiriting" so'zlari bilan taklifnoma taqdim etiladi. Shundan so'ng foydalanuvchi klaviaturadan qator kiritishi mumkin va u `message` o'zgaruvchisiga joylashtiriladi. Endi bizda socket va xabar bor, biz ushbu xabarni socket orqali manzil hostiga yuborishimiz mumkin.

go

```go
_, err = clientSocket.Write([]byte(message))
```

Ushbu qatorda `Write()` metodi yordamida xabarga manzil ma'lumoti qo'shiladi va butun natijalar paket jarayonning socketiga - `clientSocket`ga yuboriladi. (Ilgari aytilganidek, jo'natuvchi manzil ham paketga qo'shiladi, ammo bu kodni aniq yozishdan ko'ra, avtomatik ravishda amalga oshiriladi.) Kliyentdan serverga UDP socket orqali xabar yuborish shu yerda tugaydi. Ko'rinib turibdiki, bu juda oddiy! Paketni yuborgandan keyin kliyent serverdan ma'lumot olishni kutadi.

go

```go
buffer := make([]byte, 2048)
n, serverAddr, err := clientSocket.ReadFromUDP(buffer)
modifiedMessage := string(buffer[:n])
```

Ushbu qator yordamida Internet'dan kliyent socketiga kelayotgan paket ma'lumotlari `modifiedMessage` o'zgaruvchisiga, paketlarning manba manzili esa `serverAddr` o'zgaruvchisiga joylashtiriladi. Oxirgi o'zgaruvchi ham IP-manzilni, ham server port raqamini o'z ichiga oladi. Aslida UDPClient dasturiga bu ma'lumot kerak emas, chunki u boshidayoq server manzilini biladi, lekin shunga qaramay, ushbu qator dasturda mavjud. `ReadFromUDP()` metodi 2048 baytlik kirish buferini yaratadi.

go

```go
fmt.Println(modifiedMessage)
```

Ushbu qator o'zgartirilgan xabarni foydalanuvchi ekraniga chiqaradi. Bu asl qator bo'lib, undagi barcha belgilar katta harfga aylantirilgan.

go

```go
defer clientSocket.Close()
```

Bu yerda biz socketni yopamiz va jarayon tugaydi.

**UDPServer.go**

Endi dasturning server qismini ko'rib chiqaylik:

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
    
    serverAddr, err := net.ResolveUDPAddr("udp", ":"+serverPort)
    if err != nil {
        panic(err)
    }
    
    serverSocket, err := net.ListenUDP("udp", serverAddr)
    if err != nil {
        panic(err)
    }
    defer serverSocket.Close()
    
    fmt.Println("Server tayyor, qabul qilishga tayyormiz")
    
    for {
        buffer := make([]byte, 2048)
        n, clientAddr, err := serverSocket.ReadFromUDP(buffer)
        if err != nil {
            continue
        }
        message := string(buffer[:n])
        
        modifiedMessage := strings.ToUpper(message)
        
        _, err = serverSocket.WriteToUDP([]byte(modifiedMessage), clientAddr)
        if err != nil {
            continue
        }
    }
}
```

E'tibor bering, UDPServer dasturining boshlanishi UDPClient'ga juda o'xshash. U ham `net` paketini import qiladi, `serverPort` o'zgaruvchisining qiymatini 12000 ga o'rnatadi, UDP socket yaratadi. Birinchi sezilarli farq quyidagi qator:

go

```go
serverSocket, err := net.ListenUDP("udp", serverAddr)
```

Ushbu qator 12000 raqamli portni server socketi bilan bog'laydi (ya'ni tayinlaydi). Shunday qilib, UDPServer dasturida kodning qatorlari (dastur ishlab chiquvchisi tomonidan yozilgan) aniq ravishda port raqamini socketga tayinlaydi. Portni bog'lash shuni anglatadiki, endi kimdir bizning serverimizning 12000 portiga paket yuborsa, u shu socketga yo'naltiriladi. Keyin UDPServer dasturi cheksiz `for` tsikliga kiradi, bu kliyentlardan cheksiz miqdorda paketlarni qabul qilish va qayta ishlash imkonini beradi.

go

```go
buffer := make([]byte, 2048)
n, clientAddr, err := serverSocket.ReadFromUDP(buffer)
message := string(buffer[:n])
```

Keltirilgan qator biz UDPClient'da ko'rgan narsaga juda o'xshash. Server socketiga paket kelganda, paket ma'lumotlari `message` o'zgaruvchisiga, paketlar manbaasining ma'lumotlari esa `clientAddr` o'zgaruvchisiga joylashtiriladi.

`clientAddr` o'zgaruvchisi kliyentning IP-manzili va port raqamini o'z ichiga oladi. Bu ma'lumot dastur tomonidan ishlatiladi, chunki unda qaytish manzili uzatiladi va server endi o'z javobini qayerga yo'naltish kerakligini biladi.

go

```go
modifiedMessage := strings.ToUpper(message)
```

Ushbu qator bizning oddiy dasturimizning asosiy qismidir. Bu yerda biz kliyent tomonidan kiritilgan qatorni olib, `strings.ToUpper()` funksiyasidan foydalanib, uning belgilarini katta harflarga o'zgartiramiz.

go

```go
_, err = serverSocket.WriteToUDP([]byte(modifiedMessage), clientAddr)
```

Ushbu oxirgi kod qatori kliyent manzilini (IP-manzil va port raqami) o'zgartirilgan xabarga biriktiradi va natija paketini server socketiga yuboradi (ilgari aytilganidek, server manzili ham paketga biriktiriladi, ammo bu kodda emas, balki avtomatik ravishda amalga oshiriladi). Keyin Internet paketni ushbu kliyent manziliga yetkazadi.

Server paketni yuborgandan keyin cheksiz tsiklda qoladi va boshqa UDP-paketning kelishini kutadi (istalgan hostda ishga tushirilgan istalgan kliyentdan).

Dasturning qismlarini sinab ko'rish uchun shunchaki UDPClient.go'ni bir hostda, UDPServer.go'ni esa boshqa hostda ishga tushiring. Kliyent qismiga to'g'ri server nomi yoki IP-manzilni kiritishni unutmang.

Keyin siz UDPServer.go'ni server hostida ishga tushirasiz, bu serverda kliyent bilan aloqa qilishni kutadigan jarayon yaratadi. Shundan so'ng siz UDPClient.go'ni kliyentda ishga tushirasiz, shunda kliyent mashinasida jarayon yaratiladi. Nihoyat, siz shunchaki gapni kiritasiz va enter tugmasini bosasiz.

Siz o'zingizning UDP kliyent-server dasturingizni kliyent yoki server qismlarini o'zgartirish orqali ishlab chiqishingiz mumkin. Masalan, barcha belgilarni katta harfga o'tkazish o'rniga, server qatordagi "a" harflarining sonini hisoblab berishi yoki o'zgartirilgan xabarni olgandan keyin foydalanuvchi serverga boshqa gaplarni yuborishni davom ettirishi mumkin.