# 14. Diagrammalar to'plami

## 14.1. Umumiy yo'l xaritasi

```mermaid
flowchart TB
    Start([Boshlash]) --> Lev1[Daraja 1: Asoslar]
    Lev1 --> L11[Linked List]
    Lev1 --> L12[Stack/Queue]
    Lev1 --> L13[Vector]
    Lev1 --> L14[Bitset]

    L11 & L12 & L13 & L14 --> Lev2[Daraja 2: Allocator]

    Lev2 --> L21[Bump]
    Lev2 --> L22[Pool]
    Lev2 --> L23[Slab]
    Lev2 --> L24[Buddy]

    L21 & L22 & L23 & L24 --> Lev3[Daraja 3: Hash]

    Lev3 --> L31[Chained]
    Lev3 --> L32[Open Addressing]
    Lev3 --> L33[Concurrent]

    L31 & L32 & L33 --> Lev4[Daraja 4: Tree]

    Lev4 --> L41[BST/AVL]
    Lev4 --> L42[B-Tree/B+]
    Lev4 --> L43[Trie/Skip]

    L41 & L42 & L43 --> Lev5[Daraja 5: Advanced]

    Lev5 --> L51[Lock-free]
    Lev5 --> L52[LSM]
    Lev5 --> L53[GC]

    L51 & L52 & L53 --> End([Tamomladi])

    style Start fill:#87CEEB
    style End fill:#90EE90
```

## 14.2. Memory layout

```mermaid
flowchart TB
    subgraph Process["Process Address Space (64-bit)"]
        K[Kernel reserved]
        Stack[Stack ▼]
        Empty1[unmapped]
        Heap[Heap ▲]
        BSS[BSS]
        Data[Data]
        Text[Text/Code]
    end

    subgraph GoRuntime["Go Runtime Heaps"]
        GS[Goroutine Stacks 8KB+]
        MH[mheap arenas]
        MMAP[Direct mmap regions]
    end

    Heap --> GoRuntime

    style Stack fill:#90EE90
    style Heap fill:#FFE4B5
    style Text fill:#87CEEB
```

## 14.3. Slice strukturasi

```mermaid
flowchart LR
    subgraph SH["[]int slice header (24 bytes)"]
        D["Data: 0x4000"]
        L["Len: 3"]
        C["Cap: 5"]
    end

    subgraph BA["Backing array (heap)"]
        E0["[0]: 10"]
        E1["[1]: 20"]
        E2["[2]: 30"]
        E3["[3]: ?"]
        E4["[4]: ?"]
    end

    D --> E0
    L -.->|"end of len"| E2
    C -.->|"end of cap"| E4

    style SH fill:#FFE4B5
    style E0 fill:#90EE90
    style E1 fill:#90EE90
    style E2 fill:#90EE90
```

## 14.4. Map strukturasi (Go 1.18-1.23)

```mermaid
flowchart TB
    HMap[hmap struct]
    HMap -->|count| Cnt[N elements]
    HMap -->|B| Bcnt["log2(buckets)"]
    HMap -->|hash0| Seed[hash seed]
    HMap -->|buckets| BPtr[bucket array]

    BPtr --> B0[bucket 0]
    BPtr --> B1[bucket 1]
    BPtr --> Bn[bucket 2^B-1]

    B0 --> TH[tophash 0..7]
    B0 --> Keys[keys 0..7]
    B0 --> Vals[values 0..7]
    B0 --> OF[overflow ptr]

    OF --> OB[Overflow bucket]
    OB --> OF2[next overflow...]

    style HMap fill:#FFE4B5
    style B0 fill:#90EE90
```

## 14.5. Allocator hierarchy (Go runtime)

```mermaid
flowchart TB
    UC[User Code]
    UC -->|"new(), make()"| GC[gc malloc]

    GC -->|"size <= 32KB"| Small[Small allocation]
    GC -->|"size > 32KB"| Large[Large allocation]

    Small --> P_mcache["P.mcache (per-P, lock-free)"]
    P_mcache -->|"empty span"| MCentral[mcentral, per size class]
    MCentral -->|"empty"| MHeap[mheap, global]

    Large --> MHeap
    MHeap -->|"mmap()"| OS[OS Kernel]

    style P_mcache fill:#90EE90
    style MCentral fill:#FFE4B5
    style MHeap fill:#FFB6C1
    style OS fill:#DDA0DD
```

## 14.6. Lock-free CAS algoritm

```mermaid
sequenceDiagram
    participant T1 as Thread 1
    participant Mem as Shared memory
    participant T2 as Thread 2

    T1->>Mem: Read x = 10
    T2->>Mem: Read x = 10
    T1->>T1: compute new = 11
    T2->>T2: compute new = 12

    T1->>Mem: CAS(x, 10, 11)
    Mem-->>T1: success! x = 11

    T2->>Mem: CAS(x, 10, 12)
    Mem-->>T2: fail (x != 10)

    T2->>Mem: Read x = 11 (retry)
    T2->>T2: compute new = 13
    T2->>Mem: CAS(x, 11, 13)
    Mem-->>T2: success! x = 13
```

## 14.7. Garbage collector tricolor algoritm

```mermaid
flowchart LR
    subgraph Roots
        R1[Stack roots]
        R2[Global vars]
    end

    Roots -->|init: gray| Gray[Gray queue]
    Gray -->|process| Process{Process node}
    Process -->|mark refs gray| Gray
    Process -->|mark self black| Black[Black: live]

    White[White: untouched]
    White -->|"after sweep"| Free[Free memory]

    style Black fill:#333,color:#fff
    style Gray fill:#aaa
    style White fill:#fff,stroke:#333
    style Free fill:#90EE90
```

## 14.8. LSM Tree write path

```mermaid
sequenceDiagram
    participant App as Application
    participant WAL as WAL log
    participant MT as MemTable (RAM)
    participant L0 as L0 SSTable
    participant L1 as L1 SSTable

    App->>WAL: Write append (durability)
    App->>MT: Insert (key, value)

    Note over MT: MemTable to'lganda...
    MT->>L0: Flush as SSTable
    MT->>WAL: Truncate

    Note over L0,L1: Background compaction
    L0->>L1: Merge & compact
    L1->>L1: Sort, dedupe
```

---

