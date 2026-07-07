# 13. Keyingi qadamlar

## 13.1. Solod (Assembly) yo'liga o'tish

`unsafe` puxta bo'lganda, Plan9 assembly o'rganish oson:

```asm
TEXT ·Add(SB), NOSPLIT, $0-24
    MOVQ a+0(FP), AX
    MOVQ b+8(FP), BX
    ADDQ BX, AX
    MOVQ AX, ret+16(FP)
    RET
```

**Resurslar:**
- [Go assembly](https://go.dev/doc/asm)
- [9p.io/sys/doc/asm.html](http://9p.io/sys/doc/asm.html)
- `runtime/asm_amd64.s` — Go runtime asm

## 13.2. Real loyihalarda qo'llash

**Source o'qish dasturi:**

| Hafta | Loyiha | Qaysi qism |
|-------|--------|-----------|
| 1 | bbolt | `db.go`, `tx.go` (B+tree) |
| 2 | BadgerDB | `db.go`, `levels.go` (LSM) |
| 3 | ristretto | `cache.go`, `policy.go` (TinyLFU) |
| 4 | CockroachDB | `pkg/storage/` (RocksDB integration) |
| 5 | go runtime | `runtime/malloc.go`, `runtime/mgc.go` |

## 13.3. Open Source contribution

1. **Go GitHub'da issue topish** — "good first issue", "help wanted" label
2. **golang/go runtime'iga PR** — kichik fix bilan boshlash
3. **BadgerDB / bbolt PR** — DB internal'ga
4. **O'z paketingni e'lon qilish** — pkg.go.dev
5. **Blog yozish** — o'rganganingizni boshqalarga ulashish

## 13.4. Mahalliy Open Source

O'zbekistondagi Go community:
- Go Uzbekistan Telegram guruhi
- Local meetup'lar
- O'zingiz biror mini-loyiha yaratib, GitHub'ga qo'yish

---

