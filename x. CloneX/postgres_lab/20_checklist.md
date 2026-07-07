# 20. Checklist

## Storage

- [ ] Page size 8KB
- [ ] Page header
- [ ] Slot array
- [ ] Insert cell
- [ ] Delete cell
- [ ] Relation file
- [ ] ReadAt/WriteAt
- [ ] restart persistence

## Tuple/Heap

- [ ] TID
- [ ] Tuple header
- [ ] INT type
- [ ] BOOL type
- [ ] TEXT type
- [ ] Null bitmap
- [ ] Heap insert
- [ ] Heap get
- [ ] Seq scan

## Buffer

- [ ] Buffer frame
- [ ] Page table
- [ ] Pin/unpin
- [ ] Dirty flag
- [ ] Flush page
- [ ] CLOCK eviction
- [ ] WAL flush rule

## WAL/Recovery

- [ ] LSN
- [ ] WAL record header
- [ ] CRC
- [ ] Append
- [ ] Flush
- [ ] Begin/commit records
- [ ] Heap insert record
- [ ] Recovery replay
- [ ] Crash test

## Transactions/MVCC

- [ ] TxID
- [ ] Tx status
- [ ] Begin/commit/abort
- [ ] Snapshot
- [ ] Visibility
- [ ] Read Committed
- [ ] Repeatable Read
- [ ] Update creates new version

## Catalog/SQL

- [ ] Schema
- [ ] Table catalog
- [ ] Index catalog
- [ ] Lexer
- [ ] Parser
- [ ] AST
- [ ] CREATE TABLE
- [ ] INSERT
- [ ] SELECT WHERE
- [ ] UPDATE/DELETE

## Executor

- [ ] SeqScan
- [ ] Filter
- [ ] Projection
- [ ] Insert executor
- [ ] Update executor
- [ ] Delete executor
- [ ] IndexScan
- [ ] NestedLoopJoin
- [ ] Aggregate
- [ ] Sort/Limit

## Index

- [ ] B+Tree page format
- [ ] Leaf insert
- [ ] Internal insert
- [ ] Split leaf
- [ ] Split root
- [ ] Search
- [ ] Range scan
- [ ] IndexScan result matches SeqScan

## Concurrency/Maintenance

- [ ] Lock manager
- [ ] Row locks
- [ ] Deadlock detection
- [ ] Vacuum
- [ ] Free Space Map
- [ ] Visibility Map optional

## Wire Protocol

- [ ] StartupMessage
- [ ] AuthenticationOk
- [ ] Query
- [ ] RowDescription
- [ ] DataRow
- [ ] CommandComplete
- [ ] ReadyForQuery
- [ ] psql smoke test

## Verification

- [ ] unit tests
- [ ] fuzz tests
- [ ] crash tests
- [ ] race detector
- [ ] benchmarks
- [ ] pprof
- [ ] golden SQL tests
