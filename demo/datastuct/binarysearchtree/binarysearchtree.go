package binarysearchtree

type BinarySearchTree struct {
	nodeRoot   *Node
	size       int
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
