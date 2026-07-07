# 09. Catalog And Schema

Catalog database metadata saqlaydi: table, column, index, type.

## Schema

```go
type TypeID uint16

const (
    TypeInt TypeID = iota
    TypeBigInt
    TypeBool
    TypeText
)

type Column struct {
    Name string
    Type TypeID
    NotNull bool
}

type Schema struct {
    Columns []Column
}
```

## Catalog Tables

Minimal metadata:

```go
type TableDesc struct {
    OID     uint32
    Name    string
    RelPath string
    Schema  Schema
}

type IndexDesc struct {
    OID       uint32
    Name      string
    TableOID  uint32
    Columns   []int
    IndexPath string
}
```

## Storage

Avval JSON yoki binary file ishlatish mumkin:

```text
catalog/
  tables.meta
  indexes.meta
```

Keyin catalogni o'z heap tablelaringiz ichida saqlang.

## API

```go
func (c *Catalog) CreateTable(name string, schema Schema) (TableDesc, error)
func (c *Catalog) GetTable(name string) (TableDesc, error)
func (c *Catalog) CreateIndex(name, table string, cols []string) (IndexDesc, error)
func (c *Catalog) ListTables() []TableDesc
```

## Testlar

- create table
- duplicate table error
- lookup by name
- restartdan keyin catalog tiklanadi
- schema validation
