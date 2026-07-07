# 18. Milestones

## Milestone 1: Page Store

Natija:

- 8KB page
- slotted page
- relation file
- disk write/read

Demo:

```bash
go test ./internal/storage
```

## Milestone 2: Heap Table

Natija:

- tuple encode/decode
- insert/get/seqscan
- restart persistence

## Milestone 3: Buffer Manager

Natija:

- buffer pool
- dirty flush
- CLOCK eviction

## Milestone 4: WAL + Recovery

Natija:

- WAL append
- commit flush
- crash recovery

## Milestone 5: MVCC

Natija:

- begin/commit/abort
- snapshots
- visibility
- update creates new version

## Milestone 6: SQL Subset

Natija:

- CREATE TABLE
- INSERT
- SELECT WHERE
- UPDATE
- DELETE

## Milestone 7: B+Tree

Natija:

- CREATE INDEX
- IndexScan
- range scan

## Milestone 8: Joins/Aggregates

Natija:

- nested loop join
- count/sum
- group by
- order by

## Milestone 9: Wire Protocol

Natija:

```bash
psql -h localhost -p 55432
```

Minimal simple query ishlaydi.

## Milestone 10: Maintenance

Natija:

- vacuum
- free space map
- dead tuple stats
