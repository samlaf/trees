package btree

import (
	"trees/pkg/btree/pagemanager"
)

type C struct {
	tree BTree
	ref  map[string]string // the reference data
}

func newC() *C {
	pageManager := pagemanager.NewInMemory()
	tree := BTree{
		pageManager: pageManager,
	}
	return &C{
		tree: tree,
		ref:  map[string]string{},
	}
}
