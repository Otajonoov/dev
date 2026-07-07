# 5.8 Stage 6: Static Single Assignment (SSA) Generation

SSA - Static Single Assignment. Bu formda har bir value faqat bir marta assign qilinadi. Variable o'rniga immutable value'lar oqimi bilan ishlash compiler optimization uchun qulay.

![Illustration 132. SSA Generation bridges IR to machine code](images/132.png)

Go-to-SSA breakdown:

![Illustration 133. Go-to-SSA breakdown: clear and explicit](images/133.png)

## SSA value va block

SSA instruction odatda shunday tushuniladi:

- operation (`ADD`, `MUL`, `CMP`, `Load`, `Store`);
- type;
- arguments;
- auxiliary info;
- block.

![Illustration 134. Structure of an SSA value instruction](images/134.png)

Constant integer SSA value:

![Illustration 135. How a constant integer value is represented in SSA](images/135.png)

Register-based arithmetic:

![Illustration 136. Register-based arithmetic in action](images/136.png)

Memory load/store:

![Illustration 137. From memory to registers and back](images/137.png)

SSA memory token flow:

![Illustration 138. Memory token flow in SSA form](images/138.png)

## Branching va Phi

If/else SSA blocklarga bo'linadi:

![Illustration 139. Conditional branching from entry to return blocks](images/139.png)

`GOSSAFUNC` bilan SSA HTML visualization olish mumkin:

```bash
GOSSAFUNC=example go build
```

Kitobdagi visualization:

![Illustration 140. SSA visualization using GOSSAFUNC debug flag](images/140.png)

If/else logic SSA structure:

![Illustration 141. If-else logic converted to SSA structure](images/141.png)

Phi node control-flow birlashgan joyda value tanlaydi:

```go
if cond {
    y = a
} else {
    y = b
}
return y
```

SSA:

```text
y = Phi(a_from_then, b_from_else)
```

## SSA passes

SSA optimization pass pipeline:

![Illustration 142. Optimization begins the SSA pass pipeline](images/142.png)

Lowering generic SSA operationlarni target architecture operationlariga moslaydi:

![Illustration 143. Lowering adapts SSA to target architecture](images/143.png)

Layout phase block va instruction tartibini joylaydi:

![Illustration 144. Layout phase arranges blocks and instructions](images/144.png)

Register allocation SSA value'larni CPU registerlariga map qiladi:

![Illustration 145. SSA values mapped to CPU registers](images/145.png)

## Rewrite rules

SSA rewrite rules pattern matching bilan operationlarni soddalashtiradi:

![Illustration 146. Pattern-matching rules simplify expressions](images/146.png)

`x * 1` -> `x`:

![Illustration 147. Pattern match for multiplication by one](images/147.png)

`x * -1` -> `NEG x`:

![Illustration 148. Pattern match for multiplication by negative one](images/148.png)

Power-of-two multiplication:

![Illustration 149. Match Mul64 with power-of-two constant](images/149.png)

Shift bilan almashtirish:

![Illustration 150. Multiply replaced with left shift](images/150.png)

Architecture lowering:

![Illustration 151. Lowering adapts SSA to target architecture](images/151.png)

ARM64 constant add decision:

![Illustration 152. ARM64 constant add instruction decision](images/152.png)

## Conditional example va register allocation

Conditional panic/function example SSA:

![Illustration 153. SSA structure for conditional panic function](images/153.png)

Conditional variable assignment:

![Illustration 154. SSA branches for conditional variable assignment](images/154.png)

Liveness analysis:

![Illustration 155. Liveness analysis across SSA control flow](images/155.png)

Entry block:

![Illustration 156. The entry block b1](images/156.png)

ABI register mapping:

![Illustration 157. ABI forces v7 to R0, v8 to R1](images/157.png)

Register state flow:

![Illustration 158. v7 consumed; R0 now holds v19](images/158.png)

![Illustration 159. Register R2 assigned to value v19](images/159.png)

Block b4:

![Illustration 160. The block b4](images/160.png)

Phi:

![Illustration 161. Phi operation in b2](images/161.png)

Allocator Phi uchun bir xil register tanlashga harakat qiladi:

![Illustration 162. Allocator prefers same register for Phi](images/162.png)

Final register states:

![Illustration 163. v5 assigned to R2 in b4](images/163.png)

![Illustration 164. Register state: R2 = y, R0 = return value](images/164.png)

## Eslab qol

- SSA formda har value bir marta yaratiladi.
- Blocks control-flow'ni, values data-flow'ni ko'rsatadi.
- Phi node branchlardan kelgan value'ni birlashtiradi.
- Rewrite passes algebraic simplification va target-specific lowering qiladi.
- Register allocation SSA value'larni real CPU registerlariga joylaydi.
