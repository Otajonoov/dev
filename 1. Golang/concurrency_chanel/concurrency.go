package main

import "fmt"

func f(n int) {
	for i := range 10 {
		println(n, ":", i)
	}
}

func mai() {
	for i := range 10 {
		go f(i)
	}
	var input string
	fmt.Scanln(&input)
}
