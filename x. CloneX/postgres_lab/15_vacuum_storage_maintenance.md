# 15. Vacuum And Storage Maintenance

MVCC eski tuple versionlarni qoldiradi. Vacuum ularni tozalaydi.

## Dead Tuple

Tuple dead agar:

```text
xmax committed
and xmax < oldest active snapshot
```

## Vacuum V1

Heap pagelarni scan:

```text
dead tuple -> slot dead/free
free space update
```

## Free Space Map

Insert uchun qaysi pageda joy borligini topish:

```go
type FSM struct {
    free map[PageID]uint16
}
```

Keyin page-backed FSM.

## Visibility Map

Page hamma tuplelari visible bo'lsa:

```text
index-only scan uchun foydali
```

## Bloat Metrics

Stats:

- live tuples
- dead tuples
- free bytes
- pages
- bloat ratio

## Testlar

- delete + vacuum frees space
- active snapshot borida vacuum tuple o'chirmaydi
- old snapshot tugagach vacuum o'chiradi
- free space map insert page tanlaydi
