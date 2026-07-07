Функция L i s t e n создает объект n e t . L i s t e n e r , который прослушивает входя­ щие соединения на сетевом порту, в данном случае это TCP-порт l o c a l h o s t : 8000. Метод A c c e p t прослушивателя блокируется до тех пор, пока не будет сделан входя­ щий запрос на подключение, после чего возвращает объект n e t .C onn, представляю­ щий соединение. Функция h an d leC o n n обрабатывает одно полное клиентское соединение. Она в цикле выводит клиенту текущее временя, tim e .N o w (). Поскольку n e t.C o n n соот­ ветствует интерфейсу i o . W r i t e r , мы можем осуществлять вывод непосредственно в него. Цикл завершается, когда выполнение записи не удается, например потому, что клиент был отключен, и при этом h a n d le C o n n закрывает свою сторону соединения с помощью отложенного вызова C lo s e и переходит в состояние ожидания очередного запроса на подключение.

TCP — широко используемый сетевой протокол. Он лежит в основе HTTP, SSH и многих других протоколов, с которыми вы, вероятно, знакомы.

В этой концепции мы научимся писать TCP-сервер на языке Go, используя [`net`](https://pkg.go.dev/net)пакет.

## Пакет`net`​

Пакет Go [`net`](https://pkg.go.dev/net)обеспечивает доступ к сетевым примитивам.
Чтобы писать TCP-серверы на Go, вам необходимо знать следующие функции:

- [`net.Dial`](https://pkg.go.dev/net#Dial)
- [`net.Listen`](https://pkg.go.dev/net#Listen)
- [`net.Listener.Accept`](https://pkg.go.dev/net#TCPListener.Accept)
- [`net.Conn.Read`](https://pkg.go.dev/net#TCPConn.Read)
- [`net.Conn.Write`](https://pkg.go.dev/net#TCPConn.Write)

Начнем с рассмотрения [`net.Dial`](https://pkg.go.dev/net#Dial)и [`net.Listen`](https://pkg.go.dev/net#Listen).

[`net.Dial`](https://pkg.go.dev/net#Dial)используется для инициирования исходящих соединений.
Пример использования:
```go
// Connects to a TCP server on localhost:8080
conn, err := net.Dial("tcp", "localhost:8080") 
```

[`net.Listen`](https://pkg.go.dev/net#Listen)используется для создания серверов для приема входящих подключений.
Пример использования:
```go
// Starts a TCP server listening on localhost:8080
l, err := net.Listen("tcp", "localhost:8080") 
```


## Функция`net.Listen`​

Это интерфейс для [`net.Listen`](https://pkg.go.dev/net#Listen):

```go
func Listen(network string, address string) (Listener, error)
```

Чтобы создать TCP-сервер, вам нужно указать «tcp» в качестве `network`, а строку типа «localhost:8080» в качестве `address`:

```go
// Starts a TCP server listening on localhost:8080
l, err := net.Listen("tcp", "localhost:8080") 
```

## Интерфейс`net.Listener`​

[`net.Listener`](https://pkg.go.dev/net#Listener)— это интерфейс, возвращаемый из [`net.Listen`](https://pkg.go.dev/net#Listen).

```go
listener, err := net.Listen("tcp", "localhost:8080")
```

Вот функции внутри него:

```go
type Listener interface {
    // Accept waits for and returns the next connection to the listener.
    Accept() (Conn, error)

    // Close closes the listener.
    // Any blocked Accept operations will be unblocked and return errors.
    Close() error

    // Addr returns the listener's network address.
    Addr() Addr
}
```

После создания прослушивателя вы можете [`net.Listener.Accept()`](https://pkg.go.dev/net#TCPListener.Accept)дождаться подключения клиента.

Эта функция блокируется, если ни один клиент еще не подключился к серверу.

```go
// Block until we receive an incoming connection
conn, err := listener.Accept()
if err != nil {
    return err
}
```


## Интерфейс`net.Conn`​

[`net.Conn`](https://pkg.go.dev/net#Conn)— это интерфейс, возвращаемый из [`net.Listener.Accept()`](https://pkg.go.dev/net#TCPListener.Accept).

Важные функции этого интерфейса:

```go
type Conn interface {
    // Read reads data from the connection.
    Read(b []byte) (n int, err error)

    // Write writes data to the connection.
    Write(b []byte) (n int, err error)

    // Close closes the connection.
    // Any blocked Read or Write operations will be unblocked and return errors.
    Close() error
}
```

Вы можете использовать [`conn.Read()`](https://pkg.go.dev/net#TCPConn.Read)и [`conn.Write()`](https://pkg.go.dev/net#TCPConn.Write)для чтения и записи из соединения.

Чтобы прочитать данные из соединения, необходимо передать байтовый фрагмент в [`conn.Read`](https://pkg.go.dev/net#TCPConn.Read). Полученные данные будут сохранены в этом байтовом фрагменте. [`conn.Read`](https://pkg.go.dev/net#TCPConn.Read)возвращает количество прочитанных байтов:

```go
buf := make([]byte, 1024)
n, err := conn.Read(buf)
fmt.Printf("received %d bytes", n)
fmt.Printf("received the following data: %s", string(buf[:n]))
```

Чтобы записать данные в TCP-соединение, необходимо передать байтовый срез в [`conn.Write`](https://pkg.go.dev/net#TCPConn.Write). Функция возвращает количество записанных байтов:

```go
message := []byte("Hello, server!")
n, err := conn.Write(message)
fmt.Printf("sent %d bytes", n)
```


Теперь, когда вы знакомы с [net.Listen](https://pkg.go.dev/net#Listen) , [net.Listener](https://pkg.go.dev/net#Listener) и [net.Conn](https://pkg.go.dev/net#Conn) , давайте посмотрим, как объединить их для создания простого TCP-сервера, который отображает все получаемые им входные данные:

```go
package main

import (
    "fmt"
    "net"
)

func main() {
    // Listen for incoming connections
    listener, err := net.Listen("tcp", "localhost:8080")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    // Ensure we teardown the server when the program exits
    defer listener.Close()

    fmt.Println("Server is listening on port 8080")

    for {
        // Block until we receive an incoming connection
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Error:", err)
            continue
        }

        // Handle client connection
        handleClient(conn)
    }
}

func handleClient(conn net.Conn) {
    // Ensure we close the connection after we're done
    defer conn.Close()

    // Read data
    buf := make([]byte, 1024)
    n, err := conn.Read(buf)
    if err != nil {
        return
    }

    fmt.Println("Received data", buf[:n])

    // Write the same data back
    conn.Write(buf[:n])
}
```

## Краткое содержание

Теперь вы узнали, как использовать функции пакета [`net`](https://pkg.go.dev/net)для создания TCP-сервера.

Краткий обзор рассмотренных нами функций и интерфейсов:

- [`net.Listen`](https://pkg.go.dev/net#Listen): Возвращает [`net.Listener`](https://pkg.go.dev/net#Listener)экземпляр
- [`net.Listener.Accept`](https://pkg.go.dev/net#TCPListener.Accept): Блокирует до тех пор, пока клиент не подключится, возвращает [`net.Conn`](https://pkg.go.dev/net#Conn)экземпляр
- [`net.Conn.Read`](https://pkg.go.dev/net#TCPConn.Read): Считывает данные из соединения
- [`net.Conn.Write`](https://pkg.go.dev/net#TCPConn.Write): Записывает данные в соединение