package btree_serde

// We follow the BTree implementation from https://build-your-own.org/database/04_btree_code_1

type BTree struct {
	// pointer (a nonzero page number)
	root uint64
	// callbacks for managing on-disk pages
	get func(uint64) []byte // dereference a pointer
	new func([]byte) uint64 // allocate a new page
	del func(uint64)        // deallocate a page
}

// ================== 4.2 NODE UPDATE =================

// replace a link with one or multiple links
func (t *BTree) nodeReplaceKidN(new BNode, old BNode, idx uint16, kids ...BNode) {
}

// ================== 4.4 INSERTION =================
// insert a KV into a node, the result might be split.
// the caller is responsible for deallocating the input node
// and splitting and allocating result nodes.
func (t *BTree) treeInsert(node BNode, key []byte, val []byte) BNode {
	return BNode{}
}

// part of the treeInsert(): KV insertion to an internal node
func (t *BTree) nodeInsert(new BNode, node BNode, idx uint16, key []byte, val []byte) {
}
