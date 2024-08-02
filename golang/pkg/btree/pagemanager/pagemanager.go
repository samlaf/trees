package pagemanager

import "trees/pkg/btree/types"

type PageManager interface {
	// Get a page by its number
	Get(types.PagePtr) []byte
	// Allocate a new page and return its number
	New([]byte) types.PagePtr
	// Deallocate a page by its number
	Del(types.PagePtr)
}
