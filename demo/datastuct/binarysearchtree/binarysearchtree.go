package binarysearchtree

import (
	"container/list"
	"fmt"
)

type BinarySearchTree struct {
	nodeRoot  *Node
	size      int
	VCallBack CB
	IsStop    bool //是否停止遍历
}

type CB func(v interface{}) bool

type Node struct {
	Element    interface{}
	NodeRight  *Node
	NodeLeft   *Node
	NodeParent *Node
}

func CreateNode(e interface{}, parent *Node) *Node {
	return &Node{
		Element:    e,
		NodeParent: parent,
	}
}

func (t *BinarySearchTree) GetSize() int {
	return t.size
}

func (t *BinarySearchTree) IsEmpty() bool {
	return t.size == 0
}

func (t *BinarySearchTree) Clear() {
	t.nodeRoot = nil
	t.size = 0
}

func (t *BinarySearchTree) Add(e interface{}) {
	t.checkElementNotNil(e)

	//添加第一个节点
	if t.nodeRoot == nil {
		t.nodeRoot = CreateNode(e, nil)
		t.size++
		return
	}

	parent := t.nodeRoot
	node := t.nodeRoot
	cmd := 0

	for node != nil {
		parent = node
		cmd = t.compare(e, node.Element)
		if cmd > 0 {
			node = node.NodeRight
		} else if cmd < 0 {
			node = node.NodeLeft
		} else { //相等
			node.Element = e
			return
		}
	}

	newNode := CreateNode(e, parent)
	if cmd > 0 {
		parent.NodeRight = newNode
	} else {
		parent.NodeLeft = newNode
	}
	t.size++

}

func (t *BinarySearchTree) compare(e1 interface{}, e2 interface{}) int {
	return e1.(Comparator).Compare(e1, e2)
}

func (t *BinarySearchTree) checkElementNotNil(e interface{}) {
	if e == nil {
		panic("element not nil")
	}
}

//前序遍历

func (t *BinarySearchTree) PreOrderPrint() {
	t.PreOrderPrintLooper(t.nodeRoot)
}

func (t *BinarySearchTree) PreOrderPrintLooper(n *Node) {

	if n == nil || t.IsStop {
		return
	}

	if t.VCallBack != nil {
		t.IsStop = t.VCallBack(n.Element)
	}
	//fmt.Println(n.Element)
	t.PreOrderPrintLooper(n.NodeLeft)
	t.PreOrderPrintLooper(n.NodeRight)
}

//中序遍历

func (t *BinarySearchTree) MidOrderPrint() {
	t.MidOrderPrintLooper(t.nodeRoot)
}

func (t *BinarySearchTree) MidOrderPrintLooper(n *Node) {

	if n == nil || t.IsStop {
		return
	}

	t.MidOrderPrintLooper(n.NodeLeft)

	if t.IsStop {
		return
	}

	if t.VCallBack != nil {
		t.IsStop = t.VCallBack(n.Element)
	}
	//fmt.Println(n.Element)
	t.MidOrderPrintLooper(n.NodeRight)
}

//后序遍历

func (t *BinarySearchTree) AfterOrderPrint() {
	t.AfterOrderPrintLooper(t.nodeRoot)
}

func (t *BinarySearchTree) AfterOrderPrintLooper(n *Node) {

	if n == nil {
		return
	}

	t.AfterOrderPrintLooper(n.NodeLeft)
	t.AfterOrderPrintLooper(n.NodeRight)
	fmt.Println(n.Element)
}

//层序遍历
func (t *BinarySearchTree) LevelOrderPrint() {
	if t.nodeRoot == nil {
		return
	}

	l := list.New()
	l.PushBack(t.nodeRoot)

	for l.Len() != 0 {
		n := l.Remove(l.Front()).(*Node)
		if t.VCallBack != nil {
			_ = t.VCallBack(n.Element)
		}
		//fmt.Println(n.Element)
		if n.NodeLeft != nil {
			l.PushBack(n.NodeLeft)
		}
		if n.NodeRight != nil {
			l.PushBack(n.NodeRight)
		}
	}

}
