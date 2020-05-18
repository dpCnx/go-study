package main

import (
	"fmt"
	"github.com/dpCnx/go-study/demo/datastuct/binarysearchtree"
)

func main() {
	t := binarysearchtree.BinarySearchTree{}
	t.Add(1)
	t.Add(2)

	fmt.Println(t.GetSize())
}
