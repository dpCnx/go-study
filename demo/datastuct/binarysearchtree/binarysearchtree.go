package binarysearchtree

type BinarySearchTree struct {
	BaseSearchTree
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

/*
	有问题
*/
func (t *BinarySearchTree) RemoveElement(e interface{}) {
	t.Remove(t.GetNode(e))
}

func (t *BinarySearchTree) GetNode(e interface{}) *Node {
	n := t.nodeRoot
	for n != nil {
		cmp := t.compare(e, n.Element)
		if cmp == 0 {
			return n
		} else if cmp > 0 {
			n = n.NodeRight
		} else {
			n = n.NodeLeft
		}
	}

	return nil
}

func (t *BinarySearchTree) Remove(n *Node) {
	if n == nil {
		return
	}

	t.size--

	if n.HastwoChild() {
		afterN := t.successor(n)

		n.Element = afterN.Element

		n = afterN
	}

	var replacement *Node

	if n.NodeLeft != nil {
		replacement = n.NodeLeft
	} else {
		replacement = n.NodeRight
	}

	if replacement != nil {
		replacement.NodeParent = n.NodeParent
		if n.NodeParent == nil {
			t.nodeRoot = replacement
		} else if n == n.NodeParent.NodeLeft {
			n.NodeParent.NodeLeft = replacement
		} else {
			n.NodeParent.NodeRight = replacement
		}
	} else if n.NodeParent == nil {
		t.nodeRoot = nil
	} else {
		if n == n.NodeParent.NodeLeft {
			n.NodeParent.NodeLeft = nil
		} else {
			n.NodeParent.NodeRight = nil
		}
	}
}

func (t *BinarySearchTree) compare(e1 interface{}, e2 interface{}) int {
	return e1.(Comparator).Compare(e1, e2)
}

func (t *BinarySearchTree) checkElementNotNil(e interface{}) {
	if e == nil {
		panic("element not nil")
	}
}
