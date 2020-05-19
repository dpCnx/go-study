package binarysearchtree

type Myint int

func (m Myint) Compare(e interface{}, e2 interface{}) int {
	return int(e.(Myint) - e2.(Myint))
}
