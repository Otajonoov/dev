## Bosqich 4: Tree strukturalar

### 4.1. Binary Search Tree (BST)

```mermaid
flowchart TB
    R[8] --> L[3]
    R --> RR[10]
    L --> LL[1]
    L --> LR[6]
    LR --> LRL[4]
    LR --> LRR[7]
    RR --> RRR[14]
    RRR --> RRRL[13]
```

```go
type BSTNode[T constraints.Ordered] struct {
    val   T
    left  *BSTNode[T]
    right *BSTNode[T]
}

func (n *BSTNode[T]) Insert(v T) *BSTNode[T] {
    if n == nil {
        return &BSTNode[T]{val: v}
    }
    if v < n.val {
        n.left = n.left.Insert(v)
    } else if v > n.val {
        n.right = n.right.Insert(v)
    }
    return n
}

func (n *BSTNode[T]) Contains(v T) bool {
    if n == nil {
        return false
    }
    if v == n.val {
        return true
    }
    if v < n.val {
        return n.left.Contains(v)
    }
    return n.right.Contains(v)
}
```

### 4.2. AVL Tree (self-balancing)

Balance factor: `height(left) - height(right)`. |bf| <= 1 bo'lishi shart. Aks holda — rotation.

```mermaid
flowchart LR
    A[Insert qildik, balansga ziyon] --> B{Balance factor}
    B -->|> 1| LL[Left-Left rotation]
    B -->|< -1| RR[Right-Right rotation]
    B -->|0, ±1| OK[Balansda]
```

### 4.3. Red-Black Tree

5 ta qoidaga rioya qilgan self-balancing BST. Linux kernel `rb_tree`, Java `TreeMap`, Go map (eski versiyalarda).

### 4.4. B-Tree

Disk uchun mo'ljallangan. Bir node'da ko'p kalit. Database'lar (PostgreSQL, MySQL) ishlatadi.

```mermaid
flowchart TB
    R["[10, 20, 30]"] --> L1["[1, 5]"]
    R --> L2["[12, 15, 18]"]
    R --> L3["[22, 25]"]
    R --> L4["[35, 40]"]
```

### 4.5. B+ Tree

B-Tree + linked list barglarda. SQL `range scan` uchun. **bbolt** ishlatadi.

### 4.6. Trie / Radix Tree

String prefix qidirish. `etcd`, IP routing.

```mermaid
flowchart TB
    R[root] --> a[a]
    R --> b[b]
    a --> ap[ap]
    ap --> app[app]
    app --> appl[appl]
    appl --> apple[apple ●]
    appl --> apply[apply ●]
    b --> ba[ba]
    ba --> bal[bal]
    bal --> ball[ball ●]

    style apple fill:#90EE90
    style apply fill:#90EE90
    style ball fill:#90EE90
```

### 4.7. Skip List

Probabilistik balanced struktura. Redis (sorted set) ishlatadi.

```mermaid
flowchart LR
    subgraph L3["Level 3"]
        H3[H] --> N3a[10] --> N3b[40]
    end
    subgraph L2["Level 2"]
        H2[H] --> N2a[10] --> N2b[20] --> N2c[40]
    end
    subgraph L1["Level 1"]
        H1[H] --> N1a[5] --> N1b[10] --> N1c[15] --> N1d[20] --> N1e[30] --> N1f[40]
    end
```

