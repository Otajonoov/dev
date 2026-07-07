# 12. B+Tree Index

Postgres B-Tree index access methodini soddalashtirib yozing.

## Key -> TID

```go
type IndexEntry struct {
    Key Value
    TID storage.TID
}
```

## Page Types

```go
const (
    BTreeLeaf = 1
    BTreeInternal = 2
)
```

Leaf page:

```text
[keys + tids][right sibling]
```

Internal page:

```text
[separator keys + child page ids]
```

## API

```go
type BTree struct {}

func (b *BTree) Insert(key Value, tid TID) error
func (b *BTree) Search(key Value) ([]TID, error)
func (b *BTree) Delete(key Value, tid TID) error
func (b *BTree) Range(lo, hi Value) Iterator
```

## Bosqichlar

1. In-memory B+Tree.
2. Page-backed leaf nodes.
3. Page-backed internal nodes.
4. Split leaf.
5. Split root.
6. Delete without merge.
7. Range scan.

## MVCC

Index entry tuple versionga ishora qiladi. Visibility heap tuple header orqali tekshiriladi.

## Testlar

- insert/search
- duplicate keys
- leaf split
- internal split
- range scan sorted
- restartdan keyin index ishlaydi
- index scan result seq scan bilan mos
