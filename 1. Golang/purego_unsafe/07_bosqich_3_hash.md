## Bosqich 3: Hash strukturalar

### 3.1. Hash Map (Chained — separate chaining)

```mermaid
flowchart LR
    subgraph Buckets
        B0["Bucket 0"]
        B1["Bucket 1"]
        B2["Bucket 2"]
        B3["Bucket 3"]
    end

    B0 --> N1["k1, v1"]
    N1 --> N2["k5, v5"]
    N2 --> Nil1[nil]

    B1 --> Nil2[nil]

    B2 --> N3["k2, v2"]
    N3 --> Nil3[nil]

    B3 --> N4["k3, v3"]
    N4 --> N5["k7, v7"]
    N5 --> Nil4[nil]
```

```go
type Bucket[K comparable, V any] struct {
    key  K
    val  V
    next *Bucket[K, V]
}

type HashMap[K comparable, V any] struct {
    buckets []*Bucket[K, V]
    size    int
    hash    func(K) uint64
}

func (m *HashMap[K, V]) Get(k K) (V, bool) {
    h := m.hash(k) % uint64(len(m.buckets))
    for b := m.buckets[h]; b != nil; b = b.next {
        if b.key == k {
            return b.val, true
        }
    }
    var zero V
    return zero, false
}

func (m *HashMap[K, V]) Put(k K, v V) {
    h := m.hash(k) % uint64(len(m.buckets))
    for b := m.buckets[h]; b != nil; b = b.next {
        if b.key == k {
            b.val = v
            return
        }
    }
    m.buckets[h] = &Bucket[K, V]{key: k, val: v, next: m.buckets[h]}
    m.size++
    if m.size > len(m.buckets)*8 {
        m.resize()
    }
}
```

### 3.2. Hash Map (Open Addressing — Linear Probing)

Go map'ning eski versiyalari shunday ishlagan. Endi Swiss Tables.

```mermaid
flowchart LR
    subgraph Slots
        S0["k1"]
        S1["k5 (probed)"]
        S2["empty"]
        S3["k2"]
        S4["empty"]
    end

    H["hash(k5) -> 0, lekin band, +1 -> 1"] -.-> S1
```

### 3.3. Robin Hood Hashing

"Adolatli" probe distance: kim ko'p uzoq probe qilgan bo'lsa, o'rnini saqlaydi.

### 3.4. Cuckoo Hashing

Ikki hash funksiya, guaranteed O(1) lookup.

### 3.5. Concurrent Hash Map (Sharded)

```mermaid
flowchart TB
    Map[ConcurrentMap] --> S0[Shard 0 + Mutex]
    Map --> S1[Shard 1 + Mutex]
    Map --> S2[Shard 2 + Mutex]
    Map --> SN[Shard N + Mutex]

    style Map fill:#FFE4B5
```

```go
const N = 32

type ConcurrentMap[V any] struct {
    shards [N]struct {
        mu sync.RWMutex
        m  map[string]V
    }
}

func (c *ConcurrentMap[V]) shard(k string) *struct {
    mu sync.RWMutex
    m  map[string]V
} {
    h := fnv.New64a()
    h.Write([]byte(k))
    return &c.shards[h.Sum64()%N]
}

func (c *ConcurrentMap[V]) Get(k string) (V, bool) {
    s := c.shard(k)
    s.mu.RLock()
    defer s.mu.RUnlock()
    v, ok := s.m[k]
    return v, ok
}
```

### 3.6. Bloom Filter

Probabilistic data structure: "balki bor" yoki "aniq yo'q".

```mermaid
flowchart LR
    K[Key x] --> H1["hash1(x) = 3"]
    K --> H2["hash2(x) = 7"]
    K --> H3["hash3(x) = 12"]

    H1 --> B[Bit array]
    H2 --> B
    H3 --> B

    B --> Set["Set bits 3, 7, 12 = 1"]
```

```go
type BloomFilter struct {
    bits  *BitSet
    k     int // hash funksiya soni
    seeds []uint64
}

func (b *BloomFilter) Add(s string) {
    for _, seed := range b.seeds {
        h := hashWith(seed, s) % uint64(b.bits.size)
        b.bits.Set(int(h))
    }
}

func (b *BloomFilter) MayContain(s string) bool {
    for _, seed := range b.seeds {
        h := hashWith(seed, s) % uint64(b.bits.size)
        if !b.bits.Get(int(h)) {
            return false // aniq yo'q
        }
    }
    return true // balki bor
}
```

