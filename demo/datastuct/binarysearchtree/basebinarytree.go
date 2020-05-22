package binarysearchtree

import (
	"container/list"
	"fmt"
)

type BaseSearchTree struct {
	nodeRoot  *Node
	size      int
	VCallBack func(v interface{}) bool
	IsStop    bool //是否停止遍历
}

func (t *BaseSearchTree) GetSize() int {
	return t.size
}

func (t *BaseSearchTree) IsEmpty() bool {
	return t.size == 0
}

func (t *BaseSearchTree) Clear() {
	t.nodeRoot = nil
	t.size = 0
}

//前序遍历

func (t *BaseSearchTree) PreOrderPrint() {
	t.PreOrderPrintLooper(t.nodeRoot)
}

func (t *BaseSearchTree) PreOrderPrintLooper(n *Node) {

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

func (t *BaseSearchTree) MidOrderPrint() {
	t.MidOrderPrintLooper(t.nodeRoot)
}

func (t *BaseSearchTree) MidOrderPrintLooper(n *Node) {

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

func (t *BaseSearchTree) AfterOrderPrint() {
	t.AfterOrderPrintLooper(t.nodeRoot)
}

func (t *BaseSearchTree) AfterOrderPrintLooper(n *Node) {

	if n == nil {
		return
	}

	t.AfterOrderPrintLooper(n.NodeLeft)
	t.AfterOrderPrintLooper(n.NodeRight)
	fmt.Print(n.Element)
	fmt.Print(" ")
}

//层序遍历
func (t *BaseSearchTree) LevelOrderPrint() {
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

//获取高度

func (t *BaseSearchTree) Height() int {

	if t.nodeRoot == nil {
		return 0
	}

	// 树的高度
	var height = 0
	// 存储着每一层的元素数量
	var levelSize = 1

	l := list.New()
	l.PushBack(t.nodeRoot)

	for l.Len() != 0 {

		n := l.Remove(l.Front()).(*Node)
		levelSize--

		if n.NodeLeft != nil {
			l.PushBack(n.NodeLeft)
		}
		if n.NodeRight != nil {
			l.PushBack(n.NodeRight)
		}

		if levelSize == 0 {
			levelSize = l.Len()
			height++
		}
	}

	return height
}

//判断是不是完全二叉树
func (t *BaseSearchTree) IsComplete() bool {

	if t.nodeRoot == nil {
		return false
	}

	l := list.New()
	l.PushBack(t.nodeRoot)

	leaf := false

	for l.Len() != 0 {

		n := l.Remove(l.Front()).(*Node)
		if leaf && !n.IsLeaf() {
			return false
		}

		if n.NodeLeft != nil {
			l.PushBack(n.NodeLeft)
		} else if n.NodeRight != nil {
			return false
		}

		if n.NodeRight != nil {
			l.PushBack(n.NodeRight)
		} else {
			leaf = true
		}

	}

	return true

}

//反转二叉树，前序遍历

func (t *BaseSearchTree) FPreOrderPrint() {
	t.FPreOrderPrintLooper(t.nodeRoot)
}

func (t *BaseSearchTree) FPreOrderPrintLooper(n *Node) {

	if n == nil {
		return
	}

	n.NodeLeft, n.NodeRight = n.NodeRight, n.NodeLeft

	t.PreOrderPrintLooper(n.NodeLeft)
	t.PreOrderPrintLooper(n.NodeRight)
}

//前驱节点
func (t *BaseSearchTree) predecessor(node *Node) *Node {
	if node == nil {
		return nil
	}

	n := node.NodeLeft

	if n != nil {
		for n.NodeRight != nil {
			n = n.NodeRight
		}

		return n
	}

	for node.NodeParent != nil && node == node.NodeParent.NodeLeft {
		node = node.NodeParent
	}

	return node.NodeParent
}

//后驱节点
func (t *BaseSearchTree) successor(node *Node) *Node {
	if node == nil {
		return nil
	}

	n := node.NodeRight

	if n != nil {
		for n.NodeLeft != nil {
			n = n.NodeLeft
		}

		return n
	}

	for node.NodeParent != nil && node == node.NodeParent.NodeRight {
		node = node.NodeParent
	}

	return node.NodeParent
}





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

func (n *Node) IsLeaf() bool {
	return n.NodeRight == nil && n.NodeLeft == nil
}

func (n *Node) HastwoChild() bool {
	return n.NodeRight != nil && n.NodeLeft != nil
}
