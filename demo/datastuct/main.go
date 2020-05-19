package main

import (
	"fmt"
	"github.com/dpCnx/go-study/demo/datastuct/binarysearchtree"
)

func main() {
	var a binarysearchtree.Myint = 1
	var b binarysearchtree.Myint = 2

	t := binarysearchtree.BinarySearchTree{}
	t.Add(a)
	t.Add(b)

	fmt.Println(t.GetSize())

}
