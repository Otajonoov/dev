# 00. Scope

## Nima Quramiz

CloneX quyidagilarni quradi:

- page-based storage engine
- heap table
- buffer manager
- WAL va recovery
- transaction manager
- MVCC visibility
- catalog
- SQL subset
- query executor
- B+Tree index
- PostgreSQL wire protocol subset

## Nima Qurmaymiz

Birinchi versiyada:

- full PostgreSQL compatibility yo'q
- full SQL standard yo'q
- cost-based optimizerning to'liq versiyasi yo'q
- parallel query yo'q
- replication yo'q
- extension system yo'q
- PL/pgSQL yo'q
- TOAST to'liq versiyasi yo'q

## Maqsadli SQL Subset

```sql
CREATE TABLE users (id INT, name TEXT);
INSERT INTO users VALUES (1, 'Ali');
SELECT * FROM users;
SELECT id, name FROM users WHERE id = 1;
UPDATE users SET name = 'Vali' WHERE id = 1;
DELETE FROM users WHERE id = 1;
CREATE INDEX users_id_idx ON users (id);
```

Keyin:

```sql
SELECT u.id, o.total
FROM users u
JOIN orders o ON o.user_id = u.id
WHERE o.total > 100
ORDER BY o.total DESC
LIMIT 10;
```

## Texnik Chegara

- Go standard library asosiy yo'l.
- Disk I/O uchun `os.File.ReadAt/WriteAt`.
- `mmap` birinchi versiyada yo'q.
- `cgo` yo'q.
- Network protocol keyingi bosqichda.
- Avval single-node, single-process.

## Bosh Prinsip

Storage engine ishlamaguncha SQL parserga o'tmang. SQL faqat storage, transaction va executor ustidagi interface.
