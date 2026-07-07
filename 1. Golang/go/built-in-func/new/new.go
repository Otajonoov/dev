package main

/*
#include <stdlib.h>
#include <string.h>
*/
import "C"
import (
	"fmt"
)

func main() {
	// C'da hotira ajratish (malloc)
	size := C.size_t(100)
	ptr := C.malloc(size)

	// Muhim: hotira ishlatilganidan keyin ozod qilish
	defer C.free(ptr)

	// Pointer bilan ishlash
	fmt.Printf("Ajratilgan hotira manzili: %p\n", ptr)
}
