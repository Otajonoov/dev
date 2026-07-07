### **`Для чего используется Redis Pub/Sub`**

- **Обмен сообщениями в реальном времени** :
    - Чат-приложения для мгновенной доставки сообщений между пользователями или группами.
    - Push-уведомления для мобильных и веб-приложений.
- **Межведомственная связь** :
    - Упрощение коммуникации в архитектуре микросервисов, где сервисы обмениваются легковесными событиями.

### **`Преимущества Redis Pub/Sub`**

- **Низкая задержка**

    : сообщения доставляются за миллисекунды.
    
- **Простой API**
    
    : прост в реализации как для издателей, так и для подписчиков.
    
- **Масштабируемость** : Redis Cluster позволяет масштабировать данные на несколько узлов для повышения пропускной способности.
    

### **`Ограничения Redis Pub/Sub`**

- **Отсутствие сохранения**
    
    : сообщения недолговечны; они теряются, если на момент публикации ни один подписчик не подключен.
    
- **Ограниченная масштабируемость для больших подписок**
    
    : по мере роста числа каналов или подписчиков производительность может снижаться.
    
- **Отсутствие подтверждений или гарантий доставки**
    
    : не предоставляются гарантии доставки, как очереди сообщений (например, Kafka или RabbitMQ).
    

```go
package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

var ctx = context.Background()

// Message represents a chat message
type Message struct {
	Action  string `json:"action"`
	User    string `json:"user,omitempty"`
	Message string `json:"message,omitempty"`
}

type Server struct {
	clients   map[*websocket.Conn]string // Map of active WebSocket connections
	redisPub  *redis.Client             // Redis client for publishing
	redisSub  *redis.Client             // Redis client for subscribing
	mu        sync.Mutex                // Mutex for synchronizing access to clients
	upgrader  websocket.Upgrader        // WebSocket upgrader
	subscribe chan *websocket.Conn      // Channel to handle new subscriptions
}

func NewServer() *Server {
	return &Server{
		clients: make(map[*websocket.Conn]string),
		redisPub: redis.NewClient(&redis.Options{
			Addr: os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		}),
		redisSub: redis.NewClient(&redis.Options{
			Addr: os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		}),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		subscribe: make(chan *websocket.Conn),
	}
}

// Handle WebSocket connections
func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	var user string

	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			s.removeClient(conn)
			return
		}

		var msg Message
		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			log.Println("Invalid message format:", err)
			continue
		}

		switch msg.Action {
		case "join":
			user = msg.User
			s.addClient(conn, user)
			s.publishMessage(Message{Action: "control", User: user, Message: "joined the chat room"})
		case "message":
			msg.User = user
			s.publishMessage(msg)
		}
	}
}

// Add a new client
func (s *Server) addClient(conn *websocket.Conn, user string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clients[conn] = user
}

// Remove a client
func (s *Server) removeClient(conn *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if user, ok := s.clients[conn]; ok {
		delete(s.clients, conn)
		s.publishMessage(Message{Action: "control", User: user, Message: "left the chat room"})
	}
}

// Publish a message to the Redis channel
func (s *Server) publishMessage(msg Message) {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Println("Error marshaling message:", err)
		return
	}
	if err := s.redisPub.Publish(ctx, "chat", string(msgBytes)).Err(); err != nil {
		log.Println("Redis publish error:", err)
	}
}

// Listen to Redis channel and broadcast to WebSocket clients
func (s *Server) listenAndBroadcast() {
	sub := s.redisSub.Subscribe(ctx, "chat")
	ch := sub.Channel()

	for msg := range ch {
		var message Message
		if err := json.Unmarshal([]byte(msg.Payload), &message); err != nil {
			log.Println("Error unmarshaling message:", err)
			continue
		}

		s.mu.Lock()
		for conn := range s.clients {
			if err := conn.WriteJSON(message); err != nil {
				log.Println("WebSocket write error:", err)
			}
		}
		s.mu.Unlock()
	}
}

func main() {
	server := NewServer()

	go server.listenAndBroadcast()

	http.HandleFunc("/ws", server.handleWebSocket)

	log.Println("Chat server running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("HTTP server error:", err)
	}
}

```