# 08. Struct, Interface, Generics

Bu bosqich type system compiler qismini emas, runtime modelini o'rganadi.

## Struct Layout

Maqsad: field offset, alignment, padding hisoblash.

```go
type Field struct {
    Name  string
    Size  uintptr
    Align uintptr
    Offset uintptr
}

type StructDesc struct {
    Fields []Field
    Size   uintptr
    Align  uintptr
}
```

API:

```go
func Layout(fields []Field) StructDesc
```

Test:

- `unsafe.Sizeof` bilan solishtirish
- `unsafe.Offsetof` bilan solishtirish
- empty struct
- nested struct
- bool/int64 padding

## Type Descriptor

Map, interface va GC simulator uchun type metadata kerak:

```go
type TypeDesc struct {
    Name      string
    Size      uintptr
    Align     uintptr
    HasPtr    bool
    Hash      func(ptr unsafe.Pointer, seed uint64) uint64
    Equal     func(a, b unsafe.Pointer) bool
    Copy      func(dst, src unsafe.Pointer)
    Zero      func(ptr unsafe.Pointer)
}
```

## Interface Model

Empty interface:

```go
type Eface struct {
    typ  *TypeDesc
    data unsafe.Pointer
}
```

Non-empty interface:

```go
type Itab struct {
    iface *InterfaceDesc
    typ   *TypeDesc
    fun   []uintptr
}

type Iface struct {
    tab  *Itab
    data unsafe.Pointer
}
```

Labda haqiqiy Go interface'ni buzib o'qishdan ko'ra o'z modelingizni yozing. Bu ancha xavfsiz va tushunarli.

## Generics Runtime Modeli

Compiler genericsni siz yozmaysiz. Lekin generic data structure runtimega nima kerakligini ko'rsatish mumkin:

```go
type Vector struct {
    typ  *TypeDesc
    data unsafe.Pointer
    len  int
    cap  int
}
```

Bu model `[]T` emas, runtime-style generic container:

- element size
- alignment
- copy
- zero
- hash/equal kerak bo'lsa type operation

## Testlar

- struct layout Go bilan mos
- interface dynamic type saqlaydi
- type descriptor orqali generic vector ishlaydi
- pointer-free va pointer-containing typelar farqlanadi
