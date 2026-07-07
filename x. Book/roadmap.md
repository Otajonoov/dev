Muhim qoida:

- Algorithms & Data Structures alohida "keyin" o'qilmaydi. Uni 1-bosqichdan boshlab har kuni parallel olib bor.
- Har bir katta mavzudan keyin kichik loyiha qil.
- Kitoblarni 100% yodlash emas, asosiy model va trade-offlarni tushunish muhim.

---

## 0. Parallel odat: Algorithms & Data Structures

Bu yo'nalishni butun roadmap davomida yonma-yon olib bor.

1. Grokking Algorithms
2. A Common-Sense Guide to Data Structures and Algorithms
3. The Algorithm Design Manual
4. Introduction to Algorithms, CLRS
5. Algorithm Design, Kleinberg & Tardos
6. The Art of Computer Programming, Knuth

Amaliy mashq:

- Har kuni 1-2 ta LeetCode / Codeforces / NeetCode masalasi.
- Har hafta kamida 1 ta mavzu: arrays, hash table, stack, queue, tree, graph, heap, trie, DP, greedy.
- MAANG interview uchun minimum: 300-500 ta sifatli masala.

---

## 1. Computer Architecture foundation

Kompyuter qanday ishlashini tushunish: bit, gate, CPU, memory, cache, assembly, performance.

1. Code: The Hidden Language of Computer Hardware and Software
2. The Elements of Computing Systems / Nand2Tetris
3. Computer Systems: A Programmer's Perspective, CS:APP
4. Digital Design and Computer Architecture: RISC-V Edition
5. Computer Organization and Design: The Hardware/Software Interface
6. Computer Architecture: A Quantitative Approach

Eslatma:

- `Code` va `Nand2Tetris` poydevor uchun.
- Asosiy og'irlik `CS:APP`da bo'ladi.
- Advanced architecture kitoblarini keyinroq, performance va systems chuqurligi kerak bo'lganda o'qi.

Amaliy loyiha:

- Tiny assembler yoki VM interpreter yoz.
- Cache/locality benchmarklar yoz.
- C/Go'da memory layout, stack/heap, syscall experimentlar qil.

---

## 2. Operating Systems

Process, thread, scheduling, virtual memory, file system, I/O, Linux API.

1. Operating Systems: Three Easy Pieces, OSTEP
2. Operating System Concepts
3. Modern Operating Systems
4. Operating Systems: Internals and Design Principles
5. The Linux Programming Interface
6. Operating Systems: Design and Implementation
7. Understanding the Linux Kernel

Eslatma:

- `OSTEP`ni to'liq va yaxshi tushunib o'qi.
- `The Linux Programming Interface`ni reference sifatida o'qi: file I/O, process, signal, thread, socket, epoll, IPC.
- Kernel internals kitoblarini senior-level chuqurlik uchun keyinroq qoldirish mumkin.

Amaliy loyiha:

- Mini shell yoz.
- Thread pool yoz.
- epoll asosida TCP echo/chat server yoz.
- Process monitor yoki log tailer yoz.

---

## 3. Go / Golang

Go'ni production backend darajasida o'rganish.

1. Head First Go
2. The Go Programming Language
3. Learning Go
4. Go in Practice
5. Let's Go!
6. Mastering Go
7. 100 Go Mistakes and How to Avoid Them
8. Efficient Go

Eslatma:

- Agar programming tajribang bor bo'lsa, `Head First Go`ni o'tkazib yuborish mumkin.
- `The Go Programming Language` va `Learning Go` asosiy core.
- `Let's Go!` real backend qurish uchun juda foydali.
- `100 Go Mistakes` va `Efficient Go` seni senior Go darajasiga olib chiqadi.

Amaliy loyiha:

- REST API + PostgreSQL + auth + migrations.
- Background worker + queue.
- Rate limiter.
- CLI tool.
- Profiling: pprof, benchmark, memory optimization.

---

## 4. Concurrency & Parallelism

Goroutine, channel, mutex, context, cancellation, memory model, lock-free thinking.

1. Concurrency in Go
2. Go 101
3. Mastering Concurrency in Go
4. 100 Go Mistakes and How to Avoid Them
5. Java Concurrency in Practice
6. C++ Concurrency in Action
7. Seven Concurrency Models in Seven Weeks
8. Patterns for Parallel Programming
9. Structured Parallel Programming
10. The Art of Multiprocessor Programming

Qo'shimcha PDF:

- Материалы для изучения Concurrency в Go.pdf

Eslatma:

- Avval Go concurrency'ni o'rgan.
- Keyin umumiy concurrency nazariyasiga o't: locks, atomics, memory model, race, deadlock, starvation.
- `The Art of Multiprocessor Programming` expert-level, shoshilmasdan o'qiladi.

Amaliy loyiha:

- Worker pool.
- Pipeline processing.
- Context cancellation bilan distributed task runner.
- Concurrent cache.
- Lock-free queue experiment.

---

## 5. Computer Networking

HTTP, TCP/IP, DNS, TLS, load balancing, latency, retransmission, congestion.

1. Computer Networking: A Top-Down Approach
2. Computer Networks
3. Data and Computer Communications
4. Computer Networks: A Systems Approach
5. Internetworking with TCP/IP, Volume One
6. TCP/IP Illustrated, Volume 1
7. High Performance Browser Networking
8. Routing TCP/IP, Volume 1

Eslatma:

- Backend uchun eng muhimlari: `Top-Down`, `TCP/IP Illustrated`, `High Performance Browser Networking`.
- Routing kitobi network engineer chuqurligi uchun, backend uchun to'liq majburiy emas.

Amaliy loyiha:

- HTTP server/client yoz.
- TCP proxy yoz.
- DNS resolver experiment.
- TLS handshake va HTTP/2 haqida notes yoz.
- Wireshark/tcpdump bilan packet trace tahlil qil.

---

## 6. Databases

SQL, relational model, indexing, transactions, query planning, storage engine, replication.

1. Learning SQL
2. SQL Queries for Mere Mortals
3. An Introduction to Database Systems, C. J. Date
4. Database Management Systems, Ramakrishnan & Gehrke
5. Database Systems: The Complete Book
6. Designing Data-Intensive Applications
7. Database Internals
8. Transaction Processing: Concepts and Techniques
9. Principles of Distributed Database Systems
10. Readings in Database Systems

Eslatma:

- SQL'ni yaxshi bilmasdan backend kuchsiz bo'ladi.
- `DDIA` va `Database Internals` distributed backend uchun eng muhim kitoblardan.
- `Transaction Processing` va `Red Book` research/expert bosqich.

Amaliy loyiha:

- PostgreSQL bilan real schema design.
- Index va EXPLAIN ANALYZE tahlili.
- Mini key-value store yoz.
- WAL yoki LSM-tree mini implementation.
- Replication/consistency trade-off notes.

---

## 7. Distributed Systems

Failure model, latency, replication, partitioning, consensus, logical clocks, fault tolerance.

1. Understanding Distributed Systems
2. Designing Data-Intensive Applications
3. Designing Distributed Systems
4. Distributed Systems, Tanenbaum & van Steen
5. Distributed Systems: Concepts and Design
6. Introduction to Reliable and Secure Distributed Programming
7. Specifying Systems: TLA+
8. Distributed Algorithms, Nancy Lynch

Eslatma:

- `Understanding Distributed Systems` developer-friendly start.
- `DDIA` bu bosqichning eng muhim kitobi.
- `Nancy Lynch` formal proof va MIT-level chuqurlik, uni keyinroq o'qi.

Amaliy loyiha:

- Distributed key-value store.
- Leader election experiment.
- Raft implementation yoki existing Raft library bilan service.
- Idempotency, retry, timeout, backoff, circuit breaker.
- Event-driven service + outbox pattern.

---

## 8. System Design

Interview va real architecture: scalability, reliability, caching, queues, sharding, observability.

1. System Design Interview: An Insider's Guide
2. Web Scalability for Startup Engineers
3. Designing Distributed Systems
4. Designing Data-Intensive Applications
5. Building Microservices
6. Release It!
7. Site Reliability Engineering
8. Database Internals
9. Distributed Systems
10. Distributed Algorithms

Eslatma:

- `System Design Interview` framework beradi.
- `Web Scalability` amaliy web scaling beradi.
- `Release It!` production failure thinking beradi.
- Real system design uchun `DDIA`, `SRE`, `Database Internals` juda muhim.

Amaliy mashq:

- URL shortener design.
- Chat system design.
- News feed design.
- Payment system design.
- Rate limiter design.
- Distributed job scheduler design.
- Metrics/logs/traces bilan observability design.

---

## 9. DevOps, Observability, Security

Production Linux, CI/CD, Kubernetes, SRE, monitoring, tracing, incident response, secure systems.

1. UNIX and Linux System Administration Handbook
2. The Phoenix Project
3. The DevOps Handbook
4. Accelerate
5. Continuous Delivery
6. Kubernetes in Action
7. Site Reliability Engineering
8. Observability Engineering
9. Distributed Systems Observability
10. Web Application Security
11. Real-World Cryptography
12. Security Engineering

Eslatma:

- Backend engineer deploy, logs, metrics, tracing, rollback, incident, security basics'ni bilishi kerak.
- `SRE` va `Observability Engineering` production mindset beradi.
- Security'ni oxirga tashlab qo'yma: auth, secrets, TLS, OWASP, threat modeling kerak.

Amaliy loyiha:

- Service'ni Docker bilan package qil.
- Kubernetes'ga deploy qil.
- Prometheus/Grafana metrics qo'sh.
- OpenTelemetry tracing qo'sh.
- CI/CD pipeline qur.
- Load test va incident postmortem yoz.

---

## Tavsiya qilingan asosiy ketma-ketlik

Bu eng amaliy va kuchli tartib:

1. Code
2. Nand2Tetris
3. CS:APP
4. OSTEP
5. The Linux Programming Interface
6. The Go Programming Language
7. Learning Go
8. Let's Go!
9. Concurrency in Go
10. 100 Go Mistakes and How to Avoid Them
11. Efficient Go
12. Grokking Algorithms
13. The Algorithm Design Manual
14. CLRS
15. Computer Networking: A Top-Down Approach
16. TCP/IP Illustrated, Volume 1
17. High Performance Browser Networking
18. Learning SQL
19. SQL Queries for Mere Mortals
20. Database Systems: The Complete Book
21. Designing Data-Intensive Applications
22. Database Internals
23. Understanding Distributed Systems
24. Distributed Systems, Tanenbaum & van Steen
25. System Design Interview
26. Web Scalability for Startup Engineers
27. Designing Distributed Systems
28. Building Microservices
29. Release It!
30. Site Reliability Engineering
31. Observability Engineering
32. Kubernetes in Action
33. Web Application Security
34. Real-World Cryptography
35. Security Engineering

---

## Minimal core path

Agar vaqt cheklangan bo'lsa, avval mana shu kitoblarni tugat:

1. CS:APP
2. OSTEP
3. The Go Programming Language
4. Learning Go
5. Concurrency in Go
6. 100 Go Mistakes and How to Avoid Them
7. Grokking Algorithms
8. The Algorithm Design Manual
9. Computer Networking: A Top-Down Approach
10. Learning SQL
11. Database Systems: The Complete Book
12. Designing Data-Intensive Applications
13. Database Internals
14. Understanding Distributed Systems
15. System Design Interview
16. Release It!
17. Site Reliability Engineering
18. Observability Engineering

---

## Yakuniy qoida

Kuchli software engineer bo'lish formulasi:

Kitob + kod + masala + design + production tajriba.

Faqat o'qish yetmaydi. Har bir kitobdan keyin kichik, lekin ishlaydigan loyiha qil.
