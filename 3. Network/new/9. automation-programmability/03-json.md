# JSON

JSON (JavaScript Object Notation) - ma'lumotni matn ko'rinishida ifodalash formati. Network automationda JSON REST API request va response'larida juda ko'p ishlatiladi.

## JSON nima uchun kerak?

API orqali ma'lumot yuborishda ikkala tomon ham tushunadigan format kerak. JSON:

- o'qish oson;
- ko'p dasturlash tillarida qo'llab-quvvatlanadi;
- obyekt, ro'yxat, matn, son va boolean qiymatlarni ifodalaydi;
- REST APIlarda keng tarqalgan.

Misol:

```json
{
  "hostname": "SW1",
  "managementIp": "192.168.10.11",
  "model": "Catalyst 9300",
  "online": true
}
```

## JSON syntax asoslari

JSON key-value juftliklardan tuziladi.

```json
{
  "key": "value"
}
```

Qoidalar:

- obyekt `{}` bilan yoziladi;
- array `[]` bilan yoziladi;
- key nomlari qo'shtirnoq ichida bo'lishi kerak;
- string qiymatlar qo'shtirnoq ichida bo'ladi;
- har bir key-value orasida `:` ishlatiladi;
- juftliklar `,` bilan ajratiladi;
- oxirgi elementdan keyin vergul qo'yilmaydi.

To'g'ri:

```json
{
  "vlan": 10,
  "name": "USERS"
}
```

Noto'g'ri:

```json
{
  vlan: 10,
  "name": "USERS",
}
```

Xatolar:

- `vlan` key'i qo'shtirnoqsiz yozilgan;
- oxirgi elementdan keyin ortiqcha vergul bor.

## JSON data types

| Type | Misol | Izoh |
|---|---|---|
| string | `"SW1"` | matn |
| number | `10`, `3.14` | son |
| boolean | `true`, `false` | rost/yolg'on |
| null | `null` | qiymat yo'q |
| object | `{ "id": 1 }` | key-value tuzilma |
| array | `[1, 2, 3]` | ro'yxat |

## Object

Object bir nechta key-value qiymatlarni saqlaydi.

```json
{
  "interface": "GigabitEthernet1/0/1",
  "description": "User access port",
  "enabled": true
}
```

## Array

Array ro'yxatni ifodalaydi.

```json
[
  "SW1",
  "SW2",
  "R1"
]
```

Object ichida array:

```json
{
  "hostname": "SW1",
  "vlans": [10, 20, 30]
}
```

Array ichida objectlar:

```json
[
  {
    "id": 10,
    "name": "USERS"
  },
  {
    "id": 20,
    "name": "VOICE"
  }
]
```

## Network API response misoli

```json
{
  "devices": [
    {
      "hostname": "SW1",
      "managementIp": "192.168.10.11",
      "role": "access",
      "online": true
    },
    {
      "hostname": "R1",
      "managementIp": "192.168.10.1",
      "role": "router",
      "online": true
    }
  ],
  "count": 2
}
```

Bu yerda:

- `devices` - array;
- array ichida device objectlari bor;
- `count` - number;
- `online` - boolean.

## REST API request body misoli

Yangi VLAN yaratish uchun JSON body:

```json
{
  "id": 30,
  "name": "GUEST",
  "status": "active"
}
```

Interface sozlash uchun:

```json
{
  "interface": "GigabitEthernet1/0/10",
  "mode": "access",
  "accessVlan": 30,
  "description": "Guest printer"
}
```

## JSON va Python

Python dictionary ko'rinishi JSONga o'xshaydi, lekin aynan bir xil emas.

Python:

```python
device = {
    "hostname": "SW1",
    "online": True,
    "location": None
}
```

JSON:

```json
{
  "hostname": "SW1",
  "online": true,
  "location": null
}
```

Farqlar:

- Python: `True`, `False`, `None`
- JSON: `true`, `false`, `null`

Python'da JSONni o'qish:

```python
import json

text = '{"hostname": "SW1", "online": true}'
data = json.loads(text)

print(data["hostname"])
print(data["online"])
```

Python obyektini JSONga aylantirish:

```python
import json

data = {
    "hostname": "SW1",
    "vlans": [10, 20, 30]
}

json_text = json.dumps(data, indent=2)
print(json_text)
```

## YAML bilan qisqa taqqoslash

Ansible ko'pincha YAML ishlatadi, APIlar esa ko'pincha JSON ishlatadi.

JSON:

```json
{
  "vlan": 10,
  "name": "USERS"
}
```

YAML:

```yaml
vlan: 10
name: USERS
```

JSON qat'iyroq syntaxga ega. YAML o'qish uchun qulayroq, lekin indentation xatolariga sezgir.

## Common mistakes

- **Single quote ishlatish.** JSONda string uchun qo'shtirnoq `"` ishlatiladi, `'` emas.
- **Oxirgi vergul qo'yish.** JSONda trailing comma ruxsat etilmaydi.
- **Boolean qiymatni katta harf bilan yozish.** JSONda `true` va `false`, Python'dagi `True` va `False` emas.
- **Key nomini qo'shtirnoqsiz yozish.** JSONda key ham string bo'lishi kerak.
- **IP addressni son deb yozish.** IP address string bo'lishi kerak: `"192.168.1.1"`.

## Qisqa Q&A

**Savol:** JSON dasturlash tilimi?

**Javob:** Yo'q. JSON ma'lumot formati.

**Savol:** JSON fayl konfiguratsiya sifatida ishlatiladimi?

**Javob:** Ha, ayrim tizimlarda konfiguratsiya yoki API payload sifatida ishlatiladi.

**Savol:** JSONda comment yozish mumkinmi?

**Javob:** Standart JSONda comment yo'q. Comment kerak bo'lsa, hujjat yoki alohida maydon ishlatiladi.
