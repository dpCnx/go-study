package main

import (
	"fmt"
	"unsafe"
)

/*type slice struct {
	array unsafe.Pointer // 元素指针
	len   int // 长度
	cap   int // 容量
}*/

func main() {

	s := make([]int, 9, 20)
	len := *(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&s)) + uintptr(8)))
	fmt.Println(len)
	//Len: &s => pointer => uintptr => pointer => *int => int

	Cap := *(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&s)) + uintptr(16)))
	fmt.Println(Cap)
	//Cap: &s => pointer => uintptr => pointer => *int => int

	a := A{
		name: "d",
		age:  18,
	}

	aname := *(*string)(unsafe.Pointer(&a))
	fmt.Println(aname)

	aage := *(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&a)) + unsafe.Offsetof(a.age)))
	fmt.Println(aage)

	aage1 := (*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&a)) + unsafe.Offsetof(a.age)))
	*aage1 = 28
	fmt.Println(a.age)
}

type A struct {
	name string
	age  int
}
