# QoS: Classification, Marking, Queuing, Policing va Shaping

QoS (Quality of Service) tarmoqda muhim trafikni ustuvor qilish uchun ishlatiladi. Masalan, voice trafik kechikishga sezgir, oddiy file download esa biroz kutishi mumkin.

QoS bandwidth yaratmaydi. U mavjud bandwidthni tartibli taqsimlaydi.

## QoS jarayoni

1. Classification: trafikni aniqlash.
2. Marking: trafikga belgi qo'yish.
3. Queuing: navbatda ustuvorlik berish.
4. Policing: limitdan oshgan trafikni drop yoki remark qilish.
5. Shaping: trafikni silliq qilib, belgilangan tezlikda yuborish.

## Classification

Trafikni ACL, protocol, DSCP yoki interface bo'yicha ajratish mumkin.

```cisco
conf t
access-list 101 permit udp any any range 16384 32767

class-map match-any VOICE-RTP
 match access-group 101
end
```

Yoki DSCP bo'yicha:

```cisco
class-map match-any VOICE
 match dscp ef
```

## Marking

Voice uchun keng tarqalgan DSCP qiymat: EF (Expedited Forwarding), decimal 46.

```cisco
conf t
class-map match-any VOICE-RTP
 match access-group 101

policy-map MARK-VOICE
 class VOICE-RTP
  set dscp ef
 class class-default
  set dscp default

interface gigabitEthernet0/0
 service-policy input MARK-VOICE
end
```

## Queuing

LLQ (Low Latency Queueing) voice kabi kechikishga sezgir trafik uchun priority queue beradi.

```cisco
conf t
class-map match-any VOICE
 match dscp ef

class-map match-any BUSINESS
 match dscp af31

policy-map WAN-OUT
 class VOICE
  priority percent 20
 class BUSINESS
  bandwidth percent 30
 class class-default
  fair-queue

interface serial0/0/0
 service-policy output WAN-OUT
end
```

`priority` juda katta berilsa, boshqa trafik och qolishi mumkin. Shuning uchun voice trafikni real ehtiyojga mos limitlash kerak.

## Policing

Policing trafikni belgilangan limitdan oshsa darhol drop yoki remark qiladi. Bu ko'proq kiruvchi trafikni cheklashda ishlatiladi.

```cisco
conf t
policy-map POLICE-GUEST
 class class-default
  police 10000000 conform-action transmit exceed-action drop

interface gigabitEthernet0/1
 service-policy input POLICE-GUEST
end
```

Bu misolda guest trafik 10 Mbps bilan cheklanadi.

## Shaping

Shaping trafikni navbatga qo'yib, belgilangan tezlikda chiqaradi. Bu provider link tezligi physical interface tezligidan past bo'lsa foydali.

```cisco
conf t
policy-map SHAPE-WAN
 class class-default
  shape average 50000000

interface gigabitEthernet0/1
 service-policy output SHAPE-WAN
end
```

Bu misolda chiqish trafik 50 Mbps atrofida shakllantiriladi.

## Policing va shaping farqi

| Xususiyat | Policing | Shaping |
|---|---|---|
| Harakat | Oshgan trafikni drop/remark qiladi | Oshgan trafikni navbatga qo'yadi |
| Yo'nalish | Ko'pincha input | Ko'pincha output |
| Natija | Keskinroq | Silliqroq |
| Qo'llanish | Cheklash | Provider tezligiga moslash |

## Tekshiruv buyruqlari

```cisco
show policy-map
show policy-map interface
show class-map
show running-config | section policy-map
show interfaces
```

`show policy-map interface` eng foydali buyruqlardan biri, chunki match, drop va queue statistikalarini ko'rsatadi.

## Keng tarqalgan xatolar

- QoSni noto'g'ri yo'nalishda qo'llash: input kerak joyda output yoki aksincha.
- ACL classification trafikni match qilmayotgani.
- DSCP marking tarmoqning keyingi qismida trust qilinmasligi.
- Voice priority queuega juda katta foiz berish.
- Shapingni physical interface tezligiga emas, real provider tezligiga moslamaslik.

## Q&A

**Savol:** QoS internetni tezlashtiradimi?  
**Javob:** Yo'q. QoS bandwidthni oshirmaydi, faqat muhim trafikni tartiblaydi.

**Savol:** Voice uchun qaysi DSCP ishlatiladi?  
**Javob:** Odatda EF, ya'ni DSCP 46.

**Savol:** Marking qayerda qilinadi?  
**Javob:** Eng yaxshi amaliyot bo'yicha trafik tarmoqqa kirgan joyda, ya'ni access edge yoki ishonchli IP phone/switch nuqtasida.

