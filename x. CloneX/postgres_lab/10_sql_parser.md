# 10. SQL Parser

Storage ishlagandan keyin SQL parser yozing.

## Lexer

Tokenlar:

- identifiers
- keywords
- integers
- strings
- operators
- punctuation

Keywords:

```text
CREATE TABLE INSERT INTO VALUES SELECT FROM WHERE UPDATE SET DELETE
CREATE INDEX ON JOIN ORDER BY LIMIT AND OR NULL TRUE FALSE
```

## AST

```go
type Statement interface{}

type CreateTableStmt struct {
    Name string
    Columns []ColumnDef
}

type InsertStmt struct {
    Table string
    Values []Expr
}

type SelectStmt struct {
    Columns []SelectItem
    From string
    Where Expr
    Limit *int
}
```

## Parser Tartibi

1. `CREATE TABLE`
2. `INSERT`
3. `SELECT * FROM table`
4. `WHERE`
5. `UPDATE`
6. `DELETE`
7. `CREATE INDEX`
8. `JOIN`
9. `ORDER BY`, `LIMIT`

## Expression Parser

Pratt parser ishlating:

- comparison: `=`, `!=`, `<`, `<=`, `>`, `>=`
- boolean: `AND`, `OR`
- arithmetic keyinroq

## Testlar

- golden SQL -> AST
- invalid syntax
- string escaping
- keyword vs identifier
- fuzz lexer/parser no panic
