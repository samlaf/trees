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

func (t *BTree) Insert(key, val int) {
	newSplitNode := t.root.insert(key, val, t.order)
	if newSplitNode != nil {
		t.root = &Node{
			isLeaf:   false,
			keys:     []int{newSplitNode.keys[0]},
			children: []*Node{t.root, newSplitNode},
		}
	}
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

// We start with simple update in place instead of copy-on-write
// insert might returns a new node if the node was split
// it is the job of the caller (Btree.Insert) to handle the new node
func (n *Node) insert(key, val int, order int) *Node {
	if n.isLeaf {
		for i, k := range n.keys {
			if key == k {
				n.vals[i] = val
				return nil
			}
			if key < k {
				n.keys = append(n.keys, 0)
				n.vals = append(n.vals, 0)
				copy(n.keys[i+1:], n.keys[i:])
				copy(n.vals[i+1:], n.vals[i:])
				n.keys[i] = key
				n.vals[i] = val
				return n.maybeSplitLeaf(order)
			}
		}
		// if we haven't returned yet, key is greater than all keys,
		// so we insert at the end
		n.keys = append(n.keys, key)
		n.vals = append(n.vals, val)
	} else {
		for i, k := range n.keys {
			if key < k {
				newSplitNode := n.children[i].insert(key, val, order)
				if newSplitNode != nil {
					keyToInsert := newSplitNode.keys[0]
					n.children = append(n.children, nil)
					copy(n.children[i+2:], n.children[i+1:])
					n.children[i+1] = newSplitNode
					n.keys = append(n.keys, 0)
					copy(n.keys[i+1:], n.keys[i:])
					n.keys[i] = keyToInsert
					return n.maybeSplitInternal(order)
				}
			}
		}
		newSplitNode := n.children[len(n.children)-1].insert(key, val, order)
		if newSplitNode != nil {
			n.children = append(n.children, newSplitNode)
			n.keys = append(n.keys, newSplitNode.keys[0])
			return n.maybeSplitInternal(order)
		}
	}
	return nil
}

// returns a new node (rightmost) if the node was split
func (n *Node) maybeSplitLeaf(order int) *Node {
	if len(n.keys) <= order {
		return nil
	}
	mid := len(n.keys) / 2
	n.keys = n.keys[:mid]
	n.vals = n.vals[:mid]
	newNode := &Node{
		isLeaf: true,
		keys:   n.keys[mid:],
		vals:   n.vals[mid:],
	}
	return newNode
}

func (n *Node) maybeSplitInternal(order int) *Node {
	if len(n.keys) <= order {
		return nil
	}
	mid := len(n.keys) / 2
	newNode := &Node{
		isLeaf:   false,
		keys:     n.keys[mid+1:],
		children: n.children[mid+1:],
	}
	n.keys = n.keys[:mid]
	n.children = n.children[:mid+1]
	return newNode
}
