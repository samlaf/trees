package binary

type Node struct {
	Value int
	Left  *Node
	Right *Node
}

func NewLeafNode(value int) *Node {
	return &Node{Value: value}
}

// Insert inserts a new value into the tree. It returns true if the value was
// inserted, false if the value already exists in the tree.
func (n *Node) Insert(value int) bool {
	if value < n.Value {
		if n.Left == nil {
			n.Left = NewLeafNode(value)
			return true
		} else {
			return n.Left.Insert(value)
		}
	} else if value > n.Value {
		if n.Right == nil {
			n.Right = NewLeafNode(value)
			return true
		} else {
			return n.Right.Insert(value)
		}
	} else {
		// value already exists. we don't do anything
		return false
	}
}

type Tree struct {
	Root        *Node
	NumElements int
}

func NewTree() *Tree {
	return &Tree{
		Root:        nil,
		NumElements: 0,
	}
}

func (t *Tree) Insert(value int) {
	if t.Root == nil {
		t.Root = &Node{Value: value}
		t.NumElements++
	} else {
		if t.Root.Insert(value) {
			t.NumElements++
		}
	}
}

func (t )
