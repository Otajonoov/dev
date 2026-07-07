package main

type Printer interface {
	PrintFile()
}

type user struct {
	name string
	Age int
}
