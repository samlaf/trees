package btree

// We follow the BTree implementation from https://build-your-own.org/database/04_btree_code_1

import (
	"bytes"
	"trees/internal/errors"
	"trees/pkg/btree/bnode"
	"trees/pkg/btree/constant"
	"trees/pkg/btree/pagemanager"
	"trees/pkg/btree/types"
)

type BTree struct {
	// pointer (a nonzero page number)
	rootPtr types.PagePtr
	// interface for managing on-disk pages
	pageManager pagemanager.PageManager
}

// replace a kid at idx with one or multiple kids
func nodeReplaceKidN(
	tree *BTree, new bnode.BNode, old bnode.BNode, idx uint16, kids ...bnode.BNode,
) {
	inc := uint16(len(kids))
	new.SetHeader(bnode.BNODE_NODE, old.NumKeys()+inc-1)
	new.CopyPtrsAndKVs(old, 0, 0, idx)
	for i, node := range kids {
		new.CopyPtrAndKV(idx+uint16(i), tree.pageManager.New(node), node.GetKey(0), nil)
		//                ^position      ^pointer                   ^key            ^val
	}
	new.CopyPtrsAndKVs(old, idx+inc, idx+1, old.NumKeys()-(idx+1))
}

// insert a KV into a node, the result might need to be split.
// the caller is responsible for deallocating the input node
// and splitting and allocating result nodes.
func (tree *BTree) insert(bNode bnode.BNode, key []byte, val []byte) bnode.BNode {
	errors.Assert(len(key) < constant.BTREE_MAX_KEY_SIZE, "key is too big")
	errors.Assert(len(val) < constant.BTREE_MAX_VAL_SIZE, "val is too big")
	// the result node.
	// it's allowed to be bigger than 1 page and will be split if so
	new := make(bnode.BNode, 2*constant.BTREE_PAGE_SIZE)

	// where to insert the key?
	idx := bNode.LookupLE(key)
	// act depending on the node type
	switch bNode.Type() {
	case bnode.BNODE_LEAF:
		// leaf, node.getKey(idx) <= key
		if bytes.Equal(key, bNode.GetKey(idx)) {
			// found the key, update it.
			bnode.LeafUpdate(new, bNode, idx, key, val)
		} else {
			// insert it after the position.
			bnode.LeafInsert(new, bNode, idx+1, key, val)
		}
	case bnode.BNODE_NODE:
		// internal node, insert it to a kid node.
		nodeInsert(tree, new, bNode, idx, key, val)
	default:
		panic("bad node!")
	}
	return new
}

// part of the treeInsert(): KV insertion to an internal node
func nodeInsert(
	tree *BTree, new bnode.BNode, node bnode.BNode, idx uint16,
	key []byte, val []byte,
) {
	kptr := node.GetPtr(idx)
	// recursive insertion to the kid node
	knode := tree.insert(tree.pageManager.Get(kptr), key, val)
	// split the result
	nsplit, split := knode.Split3()
	// deallocate the kid node
	tree.pageManager.Del(kptr)
	// update the kid links
	nodeReplaceKidN(tree, new, node, idx, split[:nsplit]...)
}

func (tree *BTree) Insert(key []byte, val []byte) {
	if tree.rootPtr == constant.NilPagePtr {
		// create the first node
		root := make(bnode.BNode, constant.BTREE_PAGE_SIZE)
		root.SetHeader(bnode.BNODE_LEAF, 2)
		// a dummy key, this makes the tree cover the whole key space.
		// thus a lookup can always find a containing node.
		root.CopyPtrAndKV(0, 0, nil, nil)
		root.CopyPtrAndKV(1, 0, key, val)
		tree.rootPtr = tree.pageManager.New(root)
		return
	}

	rootNode := tree.pageManager.Get(tree.rootPtr)
	defer tree.pageManager.Del(tree.rootPtr)
	node := tree.insert(rootNode, key, val)
	nsplit, split := node.Split3()
	if nsplit > 1 {
		// the root was split, add a new level.
		root := make(bnode.BNode, constant.BTREE_PAGE_SIZE)
		root.SetHeader(bnode.BNODE_NODE, nsplit)
		for i, knode := range split[:nsplit] {
			ptr, key := tree.pageManager.New(knode), knode.GetKey(0)
			root.CopyPtrAndKV(uint16(i), ptr, key, nil)
		}
		tree.rootPtr = tree.pageManager.New(root)
	} else {
		tree.rootPtr = tree.pageManager.New(split[0])
	}
}

// remove a key from a leaf node
func leafDelete(new bnode.BNode, old bnode.BNode, idx uint16)

// merge 2 nodes into 1
func nodeMerge(new bnode.BNode, left bnode.BNode, right bnode.BNode)

// replace 2 adjacent links with 1
func nodeReplace2Kid(
	new bnode.BNode, old bnode.BNode, idx uint16, ptr types.PagePtr, key []byte,
)

// should the updated kid be merged with a sibling?
func shouldMerge(
	tree *BTree, node bnode.BNode,
	idx uint16, updated bnode.BNode,
) (int, bnode.BNode) {
	if updated.NumBytes() > constant.BTREE_PAGE_SIZE/4 {
		return 0, bnode.BNode{}
	}

	if idx > 0 {
		sibling := bnode.BNode(tree.pageManager.Get(node.GetPtr(idx - 1)))
		merged := sibling.NumBytes() + updated.NumBytes() - constant.HEADER_SIZE
		if merged <= constant.BTREE_PAGE_SIZE {
			return -1, sibling // left
		}
	}
	if idx+1 < node.NumKeys() {
		sibling := bnode.BNode(tree.pageManager.Get(node.GetPtr(idx + 1)))
		merged := sibling.NumBytes() + updated.NumBytes() - constant.HEADER_SIZE
		if merged <= constant.BTREE_PAGE_SIZE {
			return +1, sibling // right
		}
	}
	return 0, bnode.BNode{}
}

// delete a key from the tree
func treeDelete(tree *BTree, node bnode.BNode, key []byte) bnode.BNode

// delete a key from an internal node; part of the treeDelete()
func nodeDelete(tree *BTree, node bnode.BNode, idx uint16, key []byte) bnode.BNode {
	// recurse into the kid
	kptr := node.GetPtr(idx)
	updated := treeDelete(tree, tree.pageManager.Get(kptr), key)
	if len(updated) == 0 {
		return bnode.BNode{} // not found
	}
	tree.pageManager.Del(kptr)

	new := bnode.BNode(make([]byte, constant.BTREE_PAGE_SIZE))
	// check for merging
	mergeDir, sibling := shouldMerge(tree, node, idx, updated)
	switch {
	case mergeDir < 0: // left
		merged := bnode.BNode(make([]byte, constant.BTREE_PAGE_SIZE))
		nodeMerge(merged, sibling, updated)
		tree.pageManager.Del(node.GetPtr(idx - 1))
		nodeReplace2Kid(new, node, idx-1, tree.pageManager.New(merged), merged.GetKey(0))
	case mergeDir > 0: // right
		merged := bnode.BNode(make([]byte, constant.BTREE_PAGE_SIZE))
		nodeMerge(merged, updated, sibling)
		tree.pageManager.Del(node.GetPtr(idx + 1))
		nodeReplace2Kid(new, node, idx, tree.pageManager.New(merged), merged.GetKey(0))
	case mergeDir == 0 && updated.NumKeys() == 0:
		errors.Assert(node.NumKeys() == 1 && idx == 0, "1 empty child but no sibling")
		new.SetHeader(bnode.BNODE_NODE, 0) // the parent becomes empty too
	case mergeDir == 0 && updated.NumKeys() > 0: // no merge
		nodeReplaceKidN(tree, new, node, idx, updated)
	}
	return new
}

func (tree *BTree) Delete(key []byte) {
	if tree.rootPtr == 0 {
		return
	}
	node := treeDelete(tree, tree.pageManager.Get(tree.rootPtr), key)
	if len(node) == 0 {
		tree.pageManager.Del(tree.rootPtr)
		tree.rootPtr = 0
	} else {
		tree.rootPtr = tree.pageManager.New(node)
	}
}
