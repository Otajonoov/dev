# 17. Observability And Testing

Database testlari faqat unit test emas. Crash, fuzz, concurrency, golden output kerak.

## Unit Tests

```bash
go test ./...
```

Har package mustaqil test qilinadi.

## Golden Tests

SQL input -> expected rows:

```text
testdata/sql/select_basic.sql
testdata/sql/select_basic.golden
```

## Fuzz Tests

- page insert/delete no panic
- tuple encode/decode
- SQL lexer/parser
- WAL record decode

## Crash Tests

Processni ataylab to'xtatish:

```text
insert -> WAL written -> crash before data flush -> recovery
insert -> data dirty -> no commit -> crash -> not visible
```

Test harness:

- child process starts DB
- commands bajaradi
- process kill
- DB reopen
- invariants tekshiriladi

## Benchmarks

```bash
go test -bench=. -benchmem ./...
```

Benchmarklar:

- page insert
- heap insert
- seq scan
- buffer hit/miss
- WAL append
- B+Tree insert/search
- SQL select

## Metrics

Ichki stats:

- buffer hit ratio
- dirty pages
- WAL bytes written
- transactions committed/aborted
- dead tuples
- index pages

## Tools

- `go test -race`
- `pprof`
- `go tool trace`
- `benchstat`
- `go test -run TestCrash`
