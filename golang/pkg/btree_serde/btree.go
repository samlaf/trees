package btree_serde

import "fmt"

// Contrary to pkg/btree, pkg/btree_serde provides a btree implementation where on-disk nodes
// are deserialized into golang structs and serialized back into bytes.
// This makes them less efficient (potentially?) but easier to work with.

// According to Knuth's definition, a B-tree of order m is a tree which satisfies the following properties:
// 1. Every node has at most m children.
// 2. Every node, except for the root and the leaves, has at least ⌈m/2⌉ children.
// 3. The root node has at least two children unless it is a leaf.
// 4. All leaves appear on the same level.
// 5. A non-leaf node with k children contains k−1 keys.
type BTree struct {
	root  *Node
	order int
}

func NewBTree(order int) *BTree {
	return &BTree{order: order, root: NewLeafNode()}
}

func (t *BTree) Lookup(key int) (int, bool) {
	return t.root.lookup(key)
}

// Node can either be a leaf (isLeaf = true) or an internal node (isLeaf = false)
type Node struct {
	isLeaf bool
	// internal node should have len(children)-1 keys
	keys []int
	// only leaves have vals
	vals []int
	// only internal nodes have children
	children []*Node
}

func NewLeafNode() *Node {
	return &Node{isLeaf: true}
}

func NewInternalNode(keys []int, children []*Node) (*Node, error) {
	if len(keys) != len(children)-1 {
		return nil, fmt.Errorf("number of keys (%d) + 1 != number of children (%d)", len(keys), len(children))
	}
	for i, k := range keys {
		childKeys := children[i].keys
		if len(childKeys) == 0 || childKeys[len(childKeys)-1] >= k {
			return nil, fmt.Errorf("child %d has invalid keys", i)
		}
	}
	lastChildKeys := children[len(children)-1].keys
	if len(lastChildKeys) == 0 || lastChildKeys[0] < keys[len(keys)-1] {
		return nil, fmt.Errorf("last child has invalid keys")
	}
	return &Node{
		isLeaf:   false,
		keys:     keys,
		children: children,
	}, nil
}

func (n *Node) lookup(key int) (int, bool) {
	if n.isLeaf {
		for i, k := range n.keys {
			if k == key {
				return n.vals[i], true
			}
		}
		return 0, false
	}
	// internal node
	// child0 | key0 | child1 | key1 | ... | childn-1 | keyn-1 | childn
	// childi has all keys in the range [ keyi,  keyi+1 )
	for i, k := range n.keys {
		if key < k {
			return n.children[i].lookup(key)
		}
	}
	return n.children[len(n.children)-1].lookup(key)
}
