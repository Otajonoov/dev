# 14-bob mashqlari: Mock SMSC

> Avval mashqlarni mustaqil bajaring, keyin yechimlarga qarang. Bob matni: [book/14-mock-smsc.md](../book/14-mock-smsc.md)

---

## Mashq 1. "Bir system_id — bir bind" cheklovi

Ko'p operatorlar bitta system_id'ga bir vaqtda bitta (yoki N ta) sessiya beradi. Serverga `OneBindPerSystemID bool` konfiguratsiyasini qo'shing: yoqilganda ikkinchi bind urinishiga `RALYBND` qaytsin.

1. Qaysi strukturada tekshiruv qilasiz va qaysi lock ostida? (Race'ga e'tibor: ikki bind BIR VAQTDA kelsa?)
2. Kod yozing (server.go'ga patch) va test yozing: birinchi client bound, ikkinchi Dial RALYBND bilan yiqiladi.
3. Nozik savol: birinchi sessiya UZILGANDAN keyin qancha vaqt o'tib yangi bind qabul qilinishi kerak? Real operatorlarda bu qanday muammo tug'diradi (reconnect + RALYBND juftligi)?

## Mashq 2. BOUND_TX sessiyaga deliver_sm?

SMSC BOUND_TX sessiyaga deliver_sm yuborishga "urinsa" nima bo'lishi kerak?

1. Spec bo'yicha javob: Table 2-1 deliver_sm'ni qaysi state'larda ruxsat etadi va bu holat umuman yuzaga kelishi mumkinmi (kim aybdor bo'ladi)?
2. Bizning serverda bu strukturaviy qanday oldini olingan? (deliverableConn'ga qarang.)
3. Teskari tomonni test qiling: TX-only bound client'ga DLR yo'naltirilmasligi va log yozilishini tekshiruvchi test yozing.

## Mashq 3. Duplicate-DLR quirk'i

Real SMSC'lar deliver_sm_resp kechiksa DLR'ni QAYTA yuboradi — client duplicate'ga tayyor bo'lishi kerak (9-bob).

1. `DuplicateDLR bool` quirk'ini qo'shing: har DLR ikki marta (kichik oraliq bilan) yuborilsin.
2. Client tomonda "idempotent DLR handler" qanday ko'rinadi? `dlr.Table`ning qaysi xususiyati bu yerda ishlaydi va qaysi tartibda Forget qilish kerak?
3. Testni yozing: quirk yoqiq, client BITTA xabar holatini ikki marta emas, BIR marta yakunlashini tekshiring.

---

# Yechimlar

## Yechim 1

**1.** `Server.conns` allaqachon bor — bind paytida shu ro'yxatdan system_id bo'yicha bound sessiya qidiriladi. Lock: `Server.mu` (conns'ni himoya qiladigan mutex) — TEKSHIRUV VA O'RNATISH bitta lock ostida bo'lishi shart, aks holda ikki parallel bind ikkalasi ham "yo'q ekan" deb o'tib ketadi (check-then-act race).

**2.** Patch (g'oyasi):

```go
// serverConn.handleBind ichida, auth'dan keyin:
if c.srv.cfg.OneBindPerSystemID {
	c.srv.mu.Lock()
	busy := false
	for other := range c.srv.conns {
		if other == c {
			continue
		}
		st, id := other.status()
		if id == b.SystemID && st != session.Open && st != session.Closed {
			busy = true
			break
		}
	}
	c.srv.mu.Unlock()
	if busy {
		frame, _ := pdu.BindResp{Mode: h.ID.Resp(), Status: uint32(pdu.StatusRAlyBnd)}.Encode(h.Sequence)
		c.write(frame)
		return true
	}
}
```

(Qat'iy atomiklik uchun busy-tekshiruv va setBound'ni bitta lock ostiga olish yaxshiroq — mashqning "race" savoli aynan shu.) Test: birinchi `client.Dial` muvaffaqiyatli; ikkinchisi xato va `strings.Contains(err.Error(), "ESME_RALYBND")`.

**3.** Darhol — eski TCP sessiya YOPILGANI server tomonda aniqlangan zahoti. Muammo shunda: client tomonda uzilish sezilib reconnect boshlanadi, server tomonda esa o'lik sessiya hali "tirik" ko'rinadi (half-open! inactivity timer otguncha). Natija: reconnect'ning birinchi urinishlari RALYBND oladi. Davo: server tomonda tezroq o'lik-aniqlash (enquire_link/inactivity), client tomonda RALYBND'ni "biroz kutib qayta urinish" deb talqin qilish (session-level, lekin vaqtinchalik tabiatli — 11-bob tasnifining chekka holati).

## Yechim 2

**1.** Table 2-1: deliver_sm faqat **BOUND_RX va BOUND_TRX**da. BOUND_TX'ga deliver_sm kelishi — SERVER aybi (client "qabul qilaman" deb va'da bermagan). Client tomoni bunga RINVBNDSTS'li deliver_sm_resp qaytarishi mumkin edi — lekin yaxshi server bu holatga umuman yo'l qo'ymaydi.

**2.** `deliverableConn` FAQAT `BoundRX || BoundTRX` sessiyalarni qaytaradi — DLR/MO'ni TX sessiyaga yuborish kod yo'li mavjud emas; mos sessiya bo'lmasa log + tashlash.

**3.** Test g'oyasi: `Mode: pdu.CmdBindTransmitter` bilan client Dial; registered submit; so'ng qisqa kutish — `dc.wait` O'RNIGA teskari tekshiruv: kanalga hech narsa KELMASLIGI (`select` + timeout) va (Logf'ni yig'ib) "DLR yo'naltirib bo'lmadi" logi yozilgani.

## Yechim 3

**1.** `sendDLR` oxirida:

```go
if s.cfg.Quirks.DuplicateDLR {
	time.AfterFunc(20*time.Millisecond, func() {
		if err := c.write(frame); err != nil {
			s.logf("smsc: duplicate DLR yozish: %v", err)
		}
	})
}
```

(Ehtiyot: frame'ni qayta ishlatish xavfsiz — u immutable []byte; lekin YANGI seq bilan yuborish realroq — real SMSC retry'da yangi seq ishlatadi.)

**2.** Birinchi final DLR: `Resolve` topadi → biznes holat yakunlanadi → `Forget`. Ikkinchi (duplicate) DLR: `Resolve` endi TOPMAYDI (Forget o'chirgan) → "notanish DLR" yo'li → jim log. Ya'ni idempotentlik kaliti — **Forget'ni faqat holat COMMIT bo'lgandan keyin** chaqirish tartibi: aks holda (Forget oldin, commit keyin, commit yiqildi) haqiqiy DLR'ni ham yo'qotib qo'yasiz.

**3.** Test: quirk yoqiq server, client OnDeliver ichida hisoblagich: `Resolve` muvaffaqiyatlari soni. 1 submit → 2 deliver_sm keladi (kanalda ikkita), lekin resolvedCount == 1 va holat-yakunlash funksiyasi bir marta chaqilgan. (Diqqat: deliver_sm'lar soni 2 ekanini ham tekshiring — quirk haqiqatan ishlayotganiga ishonch.)
