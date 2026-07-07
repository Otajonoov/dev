package main

/*
#include <stdlib.h>
#include <string.h>
*/
import (
	"fmt"
	"strings"
	"unsafe"
)

type Array struct {
	data        unsafe.Pointer
	length      int
	elementSiza uintptr
}

func NewArray() {

}

func (a *Array) Get(index int) {}

func (a *Array) Set(index, value int) {}

func (a *Array) Len() int { return 0 }

func (a *Array) Cap() int {
	return a.Len()
}

func (a *Array) ToSlice() {
	s := "qwe qeqw eqw, eqw, e"
	str := strings.ReplaceAll(s, " ", "")

	fmt.Println(str)
}

func (a *Array) Copy()
