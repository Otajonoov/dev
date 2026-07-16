
Redis — мощное хранилище данных в памяти, широко используется для кэширования данных

Redis — это база данных «ключ-значение» с открытым исходным кодом.

Использование: 
	1.  cache
	2. брокер сообщений


> #Cache. Поскольку Redis хранит данные в оперативной памяти, его можно использовать **для кэширования** с целью ускорения доступа к данным и снижения нагрузки на основную базу данных.
> 
   #Дополнительно: Redis поддерживает алгоритм **LRU (Least Recently Used)**, который автоматически удаляет устаревшие данные при нехватке памяти
   #Redis как брокер сообщений (Message Queue)
	Redis поддерживает **Pub/Sub** и **структуры данных List**, что позволяет использовать его для передачи сообщений между сервисами.

**Примеры:**

- Обмен сообщениями между микросервисами
- Реализация системы уведомлений в реальном времени
- Очереди задач (Producer/Consumer model)

> **Дополнительно:** Хотя Redis не является полноценным брокером сообщений, как Kafka или RabbitMQ, он отлично подходит для простых и легковесных очередей.

### String Commands
- `SET key value` – Set a string value.
- `GET key` – Get the value of a key.
- `DEL key` – Delete a key.
- `INCR key` – Increment the integer value.
- `DECR key` – Decrement the integer value.
- `APPEND key value` – Append value to key.
- `MGET key1 key2` – Get multiple keys.
- `MSET key1 value1 key2 value2` – Set multiple keys.

### Hash Commands
- `HSET key field value` – Set a field in a hash.
- `HGET key field` – Get a field value from a hash.
- `HDEL key field` – Delete a field from a hash.
- `HGETALL key` – Get all fields and values in a hash.
- `HINCRBY key field increment` – Increment a field by a value.

### List Commands
- `LPUSH key value` – Push value to the start of a list.
- `RPUSH key value` – Push value to the end of a list.
- `LPOP key` – Remove and return the first element.
- `RPOP key` – Remove and return the last element.
- `LRANGE key start stop` – Get a range of elements.
- `LLEN key` – Get the length of a list.

### Set Commands
- `SADD key member` – Add a member to a set.
- `SREM key member` – Remove a member from a set.
- `SISMEMBER key member` – Check if a value exists in a set.
- `SMEMBERS key` – Get all members in a set.

### Sorted Set Commands
- `ZADD key score member` – Add a member with a score.
- `ZREM key member` – Remove a member.
- `ZRANK key member` – Get the rank of a member.
- `ZRANGE key start stop [WITHSCORES]` – Get a range of members.

### Key Management
- `EXISTS key` – Check if a key exists.
- `EXPIRE key seconds` – Set a timeout for a key.
- `TTL key` – Get the time-to-live of a key.
- `PERSIST key` – Remove the expiration from a key.
- `RENAME key newkey` – Rename a key.

### Transactions and Scripting
- `MULTI` – Start a transaction.
- `EXEC` – Execute a transaction.
- `DISCARD` – Discard a transaction.
- `WATCH key` – Watch a key for changes.
- `UNWATCH` – Stop watching keys.

### Server Management
- `FLUSHALL` – Clear all data from all databases.
- `FLUSHDB` – Clear the current database.
- `INFO` – Get server information.
- `PING` – Check server connection.
- `SAVE` – Synchronously save data to disk.
- `BGSAVE` – Asynchronously save data to disk.

### Pub/Sub (Publish/Subscribe)
- `PUBLISH channel message` – Send a message to a channel.
- `SUBSCRIBE channel` – Subscribe to a channel.
- `UNSUBSCRIBE channel` – Unsubscribe from a channel.


