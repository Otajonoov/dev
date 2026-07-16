package main

import (
	"fmt"
	lru "github.com/hashicorp/golang-lru"
)

var cache *lru.Cache

func init() {
	cache, _ = lru.NewWithEvict(2,
		func(key interface{}, value interface{}) {
			fmt.Printf("Evicted: key=%v value=%v\n", key, value)
		})
}

func man() {
	cache.Add(1, "a")         // добавит ключ 1
	cache.Add(2, "b")         // добавит ключ 2; теперь кеш заполнен полностью
	fmt.Println(cache.Get(1)) // "a true"; теперь 1 -- наименее давно
	// использовавшееся значение
	cache.Add(3, "c")         // добавит ключ 3, ключ 2 будет выселен
	fmt.Println(cache.Get(2)) // "<nil> false" (не найдено)
}
