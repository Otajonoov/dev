# REST API

REST API - ilovalar va tizimlar HTTP orqali ma'lumot almashishi uchun ishlatiladigan keng tarqalgan API uslubi. Network automationda REST API controller, cloud platforma yoki qurilma bilan dasturiy ishlashga yordam beradi.

## API nima?

API (Application Programming Interface) - bir tizim boshqa tizim bilan qanday gaplashishini belgilaydigan interfeys.

Oddiy misol:

```text
Automation script -> API -> Network controller -> Switches/Routers
```

Administrator CLI orqali "show interfaces" yozishi mumkin. Script esa API orqali shunga o'xshash ma'lumotni so'raydi.

## REST nima?

REST (Representational State Transfer) - resurslarga URL orqali murojaat qilish uslubi. Resurs deganda foydalanuvchi, qurilma, interface, VLAN, policy kabi obyekt tushuniladi.

Misol URL'lar:

```text
GET /api/devices
GET /api/devices/10
GET /api/vlans
POST /api/vlans
```

## Client va server

REST API'da odatda client request yuboradi, server response qaytaradi.

```text
Client:  GET /api/devices
Server:  200 OK + devices ro'yxati
```

Network misol:

```text
Python script -> REST API request -> Controller
Controller     -> REST API response -> Python script
```

## HTTP verbs va CRUD

CRUD - Create, Read, Update, Delete. REST API'da bu amallar HTTP verbs bilan bog'lanadi.

| CRUD | HTTP verb | Ma'nosi | Network misol |
|---|---|---|---|
| Create | POST | yangi resurs yaratish | yangi VLAN yaratish |
| Read | GET | ma'lumot o'qish | qurilmalar ro'yxatini olish |
| Update | PUT/PATCH | resursni o'zgartirish | interface description o'zgartirish |
| Delete | DELETE | resursni o'chirish | eski policy'ni o'chirish |

PUT odatda resursni to'liq almashtirish uchun, PATCH esa qisman o'zgartirish uchun ishlatiladi. Real tizimlarda vendor hujjatini tekshirish kerak.

## HTTP status kodlari

| Kod | Ma'nosi | Oddiy tushuncha |
|---|---|---|
| 200 OK | muvaffaqiyatli | so'rov bajarildi |
| 201 Created | yaratildi | yangi resurs yaratildi |
| 204 No Content | kontentsiz muvaffaqiyat | o'chirish yoki update bajarildi |
| 400 Bad Request | noto'g'ri so'rov | JSON yoki parametr xato |
| 401 Unauthorized | autentifikatsiya yo'q | token yoki login kerak |
| 403 Forbidden | ruxsat yo'q | login bor, huquq yetmaydi |
| 404 Not Found | topilmadi | URL yoki resurs noto'g'ri |
| 500 Internal Server Error | server xatosi | API tomonida muammo |

## Request tarkibi

REST API request odatda quyidagilardan iborat:

- method: GET, POST, PUT, PATCH, DELETE;
- URL: resurs manzili;
- headers: token, content type va boshqa metadata;
- body: POST/PUT/PATCH uchun yuboriladigan ma'lumot.

Misol:

```http
POST /api/vlans HTTP/1.1
Host: controller.example.local
Authorization: Bearer eyJhbGciOi...
Content-Type: application/json

{
  "id": 20,
  "name": "USERS"
}
```

## Response tarkibi

Response odatda status kodi, headers va body qaytaradi.

```http
HTTP/1.1 201 Created
Content-Type: application/json

{
  "id": 20,
  "name": "USERS",
  "status": "created"
}
```

## cURL bilan mini misollar

Qurilmalar ro'yxatini olish:

```bash
curl -X GET "https://controller.example.local/api/devices" \
  -H "Authorization: Bearer TOKEN"
```

Yangi VLAN yaratish:

```bash
curl -X POST "https://controller.example.local/api/vlans" \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"id":20,"name":"USERS"}'
```

VLAN nomini yangilash:

```bash
curl -X PATCH "https://controller.example.local/api/vlans/20" \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"CORP-USERS"}'
```

## Python bilan mini misol

```python
import requests

url = "https://controller.example.local/api/devices"
headers = {"Authorization": "Bearer TOKEN"}

response = requests.get(url, headers=headers, verify=False)

if response.status_code == 200:
    devices = response.json()
    for device in devices:
        print(device["hostname"], device["managementIp"])
else:
    print("API xato:", response.status_code)
```

Eslatma: `verify=False` laboratoriyada ishlatilishi mumkin, lekin production muhitda sertifikatlarni to'g'ri sozlash kerak.

## Authentication va authorization

Authentication - "siz kimsiz?" degan savolga javob.

Authorization - "sizga nima qilishga ruxsat bor?" degan savolga javob.

REST APIlarda keng uchraydigan usullar:

- username/password orqali token olish;
- API key;
- bearer token;
- OAuth token.

Header misol:

```http
Authorization: Bearer TOKEN
```

## Idempotency

Idempotent amal bir necha marta bajarilsa ham natija o'zgarmaydi.

Misollar:

- `GET` idempotent: ma'lumotni o'qiydi.
- `PUT` odatda idempotent: resursni belgilangan holatga keltiradi.
- `DELETE` odatda idempotent deb qaraladi: resurs o'chirilgan bo'lsa, qayta o'chirish yangi holat yaratmaydi.
- `POST` odatda idempotent emas: har safar yangi resurs yaratishi mumkin.

## Common mistakes

- **GET bilan o'zgartirish qilishga urinish.** GET faqat o'qish uchun bo'lishi kerak.
- **Content-Type qo'ymaslik.** JSON yuborilsa, `Content-Type: application/json` kerak bo'lishi mumkin.
- **Tokenni body ichiga qo'yish.** Ko'p APIlarda token headerda yuboriladi.
- **401 va 403ni chalkashtirish.** 401 - autentifikatsiya yo'q yoki noto'g'ri, 403 - ruxsat yetmaydi.
- **API hujjatini o'qimaslik.** URL, required fields va response formati vendor bo'yicha farq qiladi.

## Qisqa Q&A

**Savol:** REST API CLI'dan yaxshiroqmi?

**Javob:** Vazifaga bog'liq. CLI troubleshooting uchun qulay, REST API esa avtomatlashtirish va integratsiya uchun qulay.

**Savol:** Har bir REST API JSON ishlatadimi?

**Javob:** Ko'pchilik JSON ishlatadi, lekin XML yoki boshqa formatlar ham bo'lishi mumkin.

**Savol:** POST va PUT farqi nima?

**Javob:** POST ko'pincha yangi resurs yaratadi. PUT esa mavjud resursni belgilangan holatga yangilaydi yoki almashtiradi.
