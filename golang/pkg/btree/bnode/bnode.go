package bnode

import (
	"bytes"
	"encoding/binary"
	"trees/internal/errors"
	"trees/pkg/btree/constant"
)

type BNode []byte // can be dumped to the disk

const (
	BNODE_NODE = 1 // internal nodes without values
	BNODE_LEAF = 2 // leaf nodes with values
)

func (node BNode) Type() uint16 {
	return binary.LittleEndian.Uint16(node[0:2])
}
func (node BNode) NumKeys() uint16 {
	return binary.LittleEndian.Uint16(node[2:4])
}
func (node BNode) SetHeader(btype uint16, nkeys uint16) {
	binary.LittleEndian.PutUint16(node[0:2], btype)
	binary.LittleEndian.PutUint16(node[2:4], nkeys)
}

// pointers
func (node BNode) GetPtr(idx uint16) uint64 {
	errors.Assert(idx < node.NumKeys(), "idx < node.nkeys()")
	pos := constant.HEADER_SIZE + 8*idx
	return binary.LittleEndian.Uint64(node[pos:])
}
func (node BNode) SetPtr(idx uint16, val uint64) {
	errors.Assert(idx < node.NumKeys(), "idx < node.nkeys()")
	pos := constant.HEADER_SIZE + 8*idx
	binary.LittleEndian.PutUint64(node[pos:], val)
}

// offset list
func offsetPos(node BNode, idx uint16) uint16 {
	errors.Assert(1 <= idx && idx <= node.NumKeys(), "1 <= idx && idx <= node.nkeys()")
	return constant.HEADER_SIZE + 8*node.NumKeys() + 2*(idx-1)
}
func (node BNode) GetOffset(idx uint16) uint16 {
	if idx == 0 {
		return 0
	}
	return binary.LittleEndian.Uint16(node[offsetPos(node, idx):])
}
func (node BNode) SetOffset(idx uint16, offset uint16) {
	if idx == 0 {
		return
	}
	binary.LittleEndian.PutUint16(node[offsetPos(node, idx):], offset)
}

// key-values
func (node BNode) kvPos(idx uint16) uint16 {
	errors.Assert(idx <= node.NumKeys(), "idx <= node.nkeys()")
	return constant.HEADER_SIZE + 8*node.NumKeys() + 2*node.NumKeys() + node.GetOffset(idx)
}
func (node BNode) GetKey(idx uint16) []byte {
	errors.Assert(idx < node.NumKeys(), "idx < node.nkeys()")
	pos := node.kvPos(idx)
	klen := binary.LittleEndian.Uint16(node[pos:])
	return node[pos+4 : pos+4+klen]
}
func (node BNode) GetVal(idx uint16) []byte {
	errors.Assert(idx < node.NumKeys(), "idx < node.nkeys()")
	pos := node.kvPos(idx)
	klen := binary.LittleEndian.Uint16(node[pos:])
	vlen := binary.LittleEndian.Uint16(node[pos+2:])
	return node[pos+4+klen : pos+4+klen+vlen]
}

// node size in bytes
func (node BNode) NumBytes() uint16 {
	return node.kvPos(node.NumKeys())
}

// returns the first kid node whose range intersects the key. (kid[i] <= key)
// TODO: binary search
func (node BNode) LookupLE(key []byte) uint16 {
	nkeys := node.NumKeys()
	found := uint16(0)
	// the first key is a copy from the parent node,
	// thus it's always less than or equal to the key.
	for i := uint16(1); i < nkeys; i++ {
		cmp := bytes.Compare(node.GetKey(i), key)
		if cmp <= 0 {
			found = i
		}
		if cmp >= 0 {
			break
		}
	}
	return found
}

// LeafUpdate is similar to leafInsert; it updates an existing key instead of inserting a duplicate key.
func LeafUpdate(
	new BNode, old BNode, idx uint16,
	key []byte, val []byte,
) {
	new.SetHeader(BNODE_LEAF, old.NumKeys())
	copy(new, old)
	kvPos := new.kvPos(idx)
	binary.LittleEndian.PutUint16(new[kvPos:], uint16(len(key)))
	binary.LittleEndian.PutUint16(new[kvPos+2:], uint16(len(val)))
	copy(new[kvPos+4:], key)
	copy(new[kvPos+4+uint16(len(key)):], val)
}

// add a new key to a leaf node
func LeafInsert(
	new BNode, old BNode, idx uint16, key []byte, val []byte,
) {
	new.SetHeader(BNODE_LEAF, old.NumKeys()+1) // setup the header
	new.CopyPtrsAndKVs(old, 0, 0, idx)
	new.CopyPtrAndKV(idx, 0, key, val)
	new.CopyPtrsAndKVs(old, idx+1, idx, old.NumKeys()-idx)
}

// copy a ptr and KV into node at index idx, overwriting any existing data
// ptr can be set to 0 when copying into a leaf node, and val can be set to nil when writing to a node.
// TODO: given that there are 1 more ptr than key, what if we ever need to store a new ptr at the last position (that doesnt have an index?)?
func (n BNode) CopyPtrAndKV(idx uint16, ptr uint64, key []byte, val []byte) {
	// ptrs
	n.SetPtr(idx, ptr)
	// KVs
	pos := n.kvPos(idx)
	binary.LittleEndian.PutUint16(n[pos+0:], uint16(len(key)))
	binary.LittleEndian.PutUint16(n[pos+2:], uint16(len(val)))
	copy(n[pos+4:], key)
	copy(n[pos+4+uint16(len(key)):], val)
	// the offset of the next key
	n.SetOffset(idx+1, n.GetOffset(idx)+4+uint16((len(key)+len(val))))
}

// CopyPtrsAndKVs copies n ptrs and KVs from src BNode starting at index srcIdx, to dst BNode starting at index dstIdx
// It assumes that the destination BNode has enough space to hold the KVs, and will overwrite any existing KVs.
func (dst BNode) CopyPtrsAndKVs(
	src BNode, dstIdx uint16, srcIdx uint16, n uint16,
) {
	for i := uint16(0); i < n; i++ {
		dst.CopyPtrAndKV(dstIdx+i, src.GetPtr(srcIdx+i), src.GetKey(srcIdx+i), src.GetVal(srcIdx+i))
	}
}

// split a oversized node into 2 so that the 2nd node always fits on a page
func nodeSplit2(left BNode, right BNode, old BNode) {
	// code omitted...
}

// split a node if it's too big. the results are 1~3 nodes.
func (old BNode) Split3() (uint16, [3]BNode) {
	if old.NumBytes() <= constant.BTREE_PAGE_SIZE {
		old = old[:constant.BTREE_PAGE_SIZE]
		return 1, [3]BNode{old} // not split
	}
	left := BNode(make([]byte, 2*constant.BTREE_PAGE_SIZE)) // might be split later
	right := BNode(make([]byte, constant.BTREE_PAGE_SIZE))
	nodeSplit2(left, right, old)
	if left.NumBytes() <= constant.BTREE_PAGE_SIZE {
		left = left[:constant.BTREE_PAGE_SIZE]
		return 2, [3]BNode{left, right} // 2 nodes
	}
	leftleft := BNode(make([]byte, constant.BTREE_PAGE_SIZE))
	middle := BNode(make([]byte, constant.BTREE_PAGE_SIZE))
	nodeSplit2(leftleft, middle, left)
	errors.Assert(leftleft.NumBytes() <= constant.BTREE_PAGE_SIZE, "leftleft.nbytes() <= BTREE_PAGE_SIZE")
	return 3, [3]BNode{leftleft, middle, right} // 3 nodes
}
