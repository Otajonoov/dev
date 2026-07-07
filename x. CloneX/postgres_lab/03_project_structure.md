# 03. Project Structure

Tavsiya qilingan tuzilma:

```text
CloneX/
  cmd/clonex/
    main.go
  internal/
    storage/
      page.go
      slotted_page.go
      tuple.go
      relation.go
      tid.go
    buffer/
      manager.go
      frame.go
      replacer_clock.go
    wal/
      manager.go
      record.go
      lsn.go
      recovery.go
    tx/
      tx.go
      xid.go
      snapshot.go
      mvcc.go
    catalog/
      catalog.go
      schema.go
      types.go
    parser/
      lexer.go
      parser.go
      ast.go
    planner/
      plan.go
      planner.go
    executor/
      executor.go
      seqscan.go
      insert.go
      update.go
      delete.go
    index/
      btree.go
      page.go
    wire/
      server.go
      protocol.go
  testdata/
  roadmap/
```

## Package Qoidalari

- `storage` SQL bilmasin.
- `buffer` tuple formatni bilmasin.
- `wal` SQL bilmasin.
- `executor` storage va tx bilan ishlaydi.
- `parser` storage package import qilmasin.
- `catalog` schema metadata manbasi bo'lsin.

## Public API

Ichki API:

```go
type DB struct {
    Catalog *catalog.Catalog
    Buffer  *buffer.Manager
    WAL     *wal.Manager
    Tx      *tx.Manager
}

func Open(path string) (*DB, error)
func (db *DB) Close() error
```

SQLsiz test API:

```go
func (db *DB) CreateTable(name string, schema catalog.Schema) (*storage.Relation, error)
func (db *DB) Insert(table string, values []types.Value) (storage.TID, error)
func (db *DB) SeqScan(table string) ([]Row, error)
```
