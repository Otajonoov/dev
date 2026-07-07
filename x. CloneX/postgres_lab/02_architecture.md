# 02. Architecture

## Layerlar

```text
client/wire
sql parser
planner
executor
access methods
storage
buffer manager
wal/recovery
transaction/mvcc
```

## Core Paketlar

```text
internal/storage   page, relation file, tuple
internal/buffer    buffer pool
internal/wal       wal records, log manager
internal/tx        transactions, snapshots, mvcc
internal/catalog   tables, columns, indexes
internal/parser    lexer, parser, AST
internal/planner   logical/physical plans
internal/executor  plan execution
internal/index     btree
internal/wire      PostgreSQL protocol subset
```

## Data Flow

INSERT:

```text
SQL -> AST -> InsertPlan -> Executor
    -> Tx Begin
    -> Heap Insert
    -> WAL InsertRecord
    -> Buffer dirty
    -> Commit
```

SELECT:

```text
SQL -> AST -> SelectPlan
    -> SeqScan/IndexScan
    -> MVCC visibility
    -> Filter
    -> Projection
    -> Result rows
```

Recovery:

```text
Open DB
  read checkpoint
  replay WAL from checkpoint
  rebuild dirty pages
  rollback incomplete tx
```

## Birinchi Versiya Chegarasi

V0:

- bitta process
- bitta database
- simple catalog
- single writer
- no network
- no index

V1:

- SQL subset
- WAL
- MVCC Read Committed
- B+Tree

V2:

- wire protocol
- locks
- vacuum
- planner improvements
