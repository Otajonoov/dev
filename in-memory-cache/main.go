package main

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	data        map[string]interface{}
	expiration  map[string]time.Time
	mutex       sync.RWMutex
	defaultTTL  time.Duration
	cleanupTick time.Duration
}

func NewCache(defaultTTL, cleanupTick time.Duration) *Cache {
	cache := &Cache{
		data:        make(map[string]interface{}),
		expiration:  make(map[string]time.Time),
		defaultTTL:  defaultTTL,
		cleanupTick: cleanupTick,
	}

	go cache.startCleanup()

	return cache
}

func (c *Cache) startCleanup() {
	ticker := time.NewTicker(c.cleanupTick)
	for {
		select {
		case <-ticker.C:
			c.cleanup()
		}
	}
}

func (c *Cache) cleanup() {
	currentTime := time.Now()

	c.mutex.Lock()
	defer c.mutex.Unlock()

	for key, expirationTime := range c.expiration {
		if currentTime.After(expirationTime) {
			delete(c.data, key)
			delete(c.expiration, key)
		}
	}

}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data[key] = value
	c.expiration[key] = time.Now().Add(ttl)
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	value, ok := c.data[key]
	return value, ok
}

func (c *Cache) Delete(key string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.data, key)
	delete(c.expiration, key)

	return nil
}

func main() {
	cache := NewCache(5*time.Second, 1*time.Second)

	cache.Set("key1", "value1", 2*time.Second)
	cache.Set("key2", "value2", 4*time.Second)

	value, ok := cache.Get("key1")
	fmt.Println(value, ok)

	value, ok = cache.Get("key2")
	fmt.Println(value, ok)

	time.Sleep(5 * time.Second)

	value1, ok1 := cache.Get("key1")
	if ok1 {
		fmt.Println("Value 1:", value1)
	} else {
		fmt.Println("Value 1 expired")
	}

	value2, ok2 := cache.Get("key2")
	if ok2 {
		fmt.Println("Value 2:", value2)
	} else {
		fmt.Println("Value 2 expired")
	}

}
