package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
)

// install and run with cmd
// go get github.com/redis/go-redis/v9
// brew install redis
// brew services start redis
// redis-cli ping

var ctx = context.Background()

func main() {
	// client 1
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "1234",
		DB:       0,
	})

	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Error connecting to Redis: ", err)
	}
	fmt.Println("Connected to Redis:", pong)

	// client 2
	//opt, err := redis.ParseURL("redis://:1234@localhost:6379/0") // "redis://<user>:<pass>@localhost:6379/<db>"
	//if err != nil {
	//	panic(err)
	//}
	//client2 := redis.NewClient(opt)

	// with key: value
	//_ = client.Set(ctx, "key", "value", 0).Err()
	//
	//val, _ := client.Get(ctx, "key").Result()
	//fmt.Println("Key:", val)
	//
	// with map
	//user := map[string]interface{}{
	//	"name": "Quvonchbek",
	//	"age":  25,
	//}
	//
	//_ = client.HSet(ctx, "quvonchbek", user).Err()
	//
	//userInfo, _ := client.HGet(ctx, "quvonchbek", "name").Result()
	//fmt.Println("Name:", userInfo)
	//
	//// if you want to get all values
	//userAllInfo, _ := client.HGetAll(ctx, "quvonchbek").Result()
	//fmt.Println("Name:", userAllInfo)

	//// with struct
	//u := User{Name: "Quvonchbek", Age: 25}
	//
	//_ = client.Set(ctx, "quvonchbek", u, 0).Err()
	//
	//userInfo, _ = client.Get(ctx, "quvonchbek").Result()
	//fmt.Println("Name:", userInfo)

	// with list
	// LPush insert values at the head of the list
	//_ = client.LPush(ctx, "list", "1", "2", "3").Err()
	//
	//// RPush insert values at the back/ tail of the list
	//_ = client.RPush(ctx, "list1", "4", "5", "6").Err()
	//
	//// LPop remove the first element in the list
	//list, _ := client.LPop(ctx, "list").Result()
	//fmt.Println("Popped Task: ", list)
	//
	//// RPopr remove the last element in the list
	//list, _ = client.RPop(ctx, "list").Result()
	//fmt.Println("Popped Task: ", list)

	// Set a key with expiration time
	//_ = client.Set(ctx, "key", "value", 3*time.Second).Err()
	//
	//// Get the value of a key
	//val, _ := client.Get(ctx, "key").Result()
	//fmt.Println("Key:", val)
	//
	//// get key expiration time
	//time.Sleep(4 * time.Second)
	//expiration, _ := client.TTL(ctx, "key").Result()
	//fmt.Println("Key expiration time:", expiration)

	// publish/subscribe messaging
	//_ = client.Publish(ctx, "mychannel", "Hello, World!").Err()
	//
	//// subscribe to a channel
	//pubsub := client.Subscribe(ctx, "mychannel")
	//channel := pubsub.Channel()
	//
	//fmt.Println("Subscribed to channel: mychannel")
	//
	//for msg := range channel {
	//	fmt.Printf("Received message: %s\n", msg.Payload)
	//}
}

type User struct {
	Name string
	Age  int
}
