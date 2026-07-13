package main

import (
	"container/list"
	"fmt"
)

func main() {
	myList := list.New()

	// Добавляем каждый элемент в конец
	myList.PushBack(1)
	myList.PushBack(2)
	myList.PushBack(3)

	// Пробегаемся по списку и печатаем не пустые элементы
	// Мы не можем пробежаться привычным способом как с массивами,
	// поэтому придется использовать метод Front()
	// которая вернет первый элемент и затем с помощью Next
	// получать следующий элемент пока он не будет равен nil, что означает конец списка
	for element := myList.Front(); element != nil; element = element.Next() {
		fmt.Println(element.Value)
	}
}
