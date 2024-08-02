package pagemanager

import (
	"trees/internal/errors"
	"trees/pkg/btree/bnode"
	"trees/pkg/btree/constant"
	"unsafe"
)

// memory.Inmemory is meant to be used for testing the BTree implementations
type InMemory struct {
	ref   map[string]string      // the reference data
	pages map[uint64]bnode.BNode // in-memory pages
}

var _ PageManager = (*InMemory)(nil)

func NewInMemory() *InMemory {
	return &InMemory{
		ref:   map[string]string{},
		pages: map[uint64]bnode.BNode{},
	}
}

func (pm *InMemory) Get(ptr uint64) []byte {
	node, ok := pm.pages[ptr]
	errors.Assert(ok, "page not found")
	return node
}

func (pm *InMemory) New(node []byte) uint64 {
	errors.Assert(bnode.BNode(node).NumBytes() <= constant.BTREE_PAGE_SIZE, "node size exceeds page size")
	ptr := uint64(uintptr(unsafe.Pointer(&node[0])))
	errors.Assert(pm.pages[ptr] == nil, "page already exists")
	pm.pages[ptr] = node
	return ptr
}

func (pm *InMemory) Del(ptr uint64) {
	errors.Assert(pm.pages[ptr] != nil, "page not found")
	delete(pm.pages, ptr)
}
