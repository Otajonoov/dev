# 21. Resources

## PostgreSQL Source Reading

O'qish tartibi:

1. `src/include/storage/bufpage.h` - page layout
2. `src/include/access/htup_details.h` - heap tuple header
3. `src/backend/access/heap/heapam.c` - heap access method
4. `src/backend/storage/buffer/bufmgr.c` - buffer manager
5. `src/backend/access/transam/xlog.c` - WAL
6. `src/backend/access/transam/xact.c` - transactions
7. `src/backend/utils/time/snapmgr.c` - snapshots
8. `src/backend/access/nbtree/` - B-Tree index
9. `src/backend/parser/` - parser
10. `src/backend/executor/` - executor
11. `src/backend/tcop/postgres.c` - query loop
12. `src/backend/libpq/` - wire protocol

## Books

- Database Internals - Alex Petrov
- Designing Data-Intensive Applications - Martin Kleppmann
- Architecture of a Database System - Hellerstein, Stonebraker, Hamilton
- PostgreSQL 14 Internals
- Readings in Database Systems

## Docs

- PostgreSQL documentation: storage, MVCC, WAL, indexes
- PostgreSQL Frontend/Backend Protocol
- SQL standard basics

## Go Tools

- `go test`
- `go test -race`
- `go test -fuzz`
- `go test -bench`
- `pprof`
- `go tool trace`
- `benchstat`

## Codebases To Study

- PostgreSQL
- SQLite
- CockroachDB
- etcd bbolt
- BadgerDB
- Pebble

## CloneX Reading Rule

Har Postgres source file o'qilganda:

1. faqat kerakli struct/algorithmni yozib oling
2. Go'da mini versiyasini yozing
3. test bilan behaviorni mustahkamlang
4. to'liq Postgres murakkabligini darhol ko'chirmang
