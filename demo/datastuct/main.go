package main

import (
	"github.com/dpCnx/go-study/demo/datastuct/binarysearchtree"
	"log"
)

func main() {
	var a binarysearchtree.Myint = 9
	var b binarysearchtree.Myint = 3
	var c binarysearchtree.Myint = 7
	var d binarysearchtree.Myint = 5
	var e binarysearchtree.Myint = 8
	var f binarysearchtree.Myint = 15
	var g binarysearchtree.Myint = 13
	var h binarysearchtree.Myint = 14
	var i binarysearchtree.Myint = 4

	t := binarysearchtree.BinarySearchTree{}
	t.Add(a)
	t.Add(b)
	t.Add(c)
	t.Add(d)
	t.Add(e)
	t.Add(f)
	t.Add(g)
	t.Add(h)
	t.Add(i)

	y := 0
	t.VCallBack = func(v interface{}) bool {
		if y == 2 {
			return true
		}

		log.Println(v)
		y++
		return false
	}

	t.PreOrderPrint()
}
