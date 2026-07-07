# 16. PostgreSQL Wire Protocol

Storage, SQL va executor ishlagandan keyin `psql` bilan ulanish uchun protocol subset yozing.

## Protocol V3 Minimal

Qo'llab-quvvatlanadigan flow:

```text
StartupMessage
AuthenticationOk
ParameterStatus
ReadyForQuery
Query
RowDescription
DataRow
CommandComplete
ReadyForQuery
Terminate
```

## Server

```go
type Server struct {
    DB *DB
}

func (s *Server) Listen(addr string) error
```

## Query Message

Avval simple query protocol:

```text
Q + length + query string
```

Extended protocol keyin:

- Parse
- Bind
- Execute
- Sync

## Authentication

V0:

```text
trust auth
```

Keyin:

- cleartext password
- MD5/SASL optional

## Testlar

- raw TCP client
- startup handshake
- simple select returns rows
- insert command complete
- `psql -h localhost -p 55432` bilan smoke test
