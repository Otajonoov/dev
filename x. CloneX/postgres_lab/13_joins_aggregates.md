# 13. Joins, Aggregates, Sort

Storage va simple executor ishlagandan keyin query capability kengayadi.

## Join

Birinchi join:

```text
Nested Loop Join
```

Plan:

```go
type NestedLoopJoinPlan struct {
    Left Plan
    Right Plan
    Cond Expr
}
```

Keyin:

- Hash Join
- Merge Join

## Aggregates

Boshlang'ich aggregate:

```sql
SELECT count(*) FROM users;
SELECT user_id, count(*) FROM orders GROUP BY user_id;
```

Aggregate functions:

- `count`
- `sum`
- `min`
- `max`

## Sort

Avval in-memory sort:

```sql
ORDER BY col ASC|DESC
```

Keyin external sort:

```text
run generation -> merge
```

## Limit

`LIMIT` executor pipeline ichida erta to'xtashi kerak.

## Testlar

- nested loop join
- hash join result nested loop bilan mos
- count/sum
- group by
- order by
- limit
