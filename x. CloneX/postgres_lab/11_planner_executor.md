# 11. Planner And Executor

Parser AST beradi. Planner ASTni execution planga aylantiradi. Executor planni bajaradi.

## Plan Nodes

```go
type Plan interface{}

type SeqScanPlan struct {
    Table string
    Filter Expr
}

type IndexScanPlan struct {
    Table string
    Index string
    Key Expr
}

type ProjectionPlan struct {
    Input Plan
    Columns []string
}

type InsertPlan struct {}
type UpdatePlan struct {}
type DeletePlan struct {}
```

## Executor Interface

```go
type Executor interface {
    Init(ctx *ExecContext) error
    Next() (Row, bool, error)
    Close() error
}
```

Executorlar:

- `SeqScan`
- `Filter`
- `Projection`
- `Insert`
- `Update`
- `Delete`
- `IndexScan`
- `NestedLoopJoin`

## Planner V1

Simple rule-based:

```text
SELECT WHERE indexed_col = const -> IndexScan
else -> SeqScan
```

## Executor Context

```go
type ExecContext struct {
    DB *DB
    Tx *tx.Tx
    Snapshot tx.Snapshot
}
```

## Testlar

- select all
- where filter
- insert plan
- update creates new tuple version
- delete marks tuple
- index scan and seq scan same result
