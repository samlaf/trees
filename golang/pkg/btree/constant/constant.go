package constant

import (
	"trees/internal/errors"
	"trees/pkg/btree/types"
)

// sizes are all in bytes
const HEADER_SIZE = 4

const BTREE_PAGE_SIZE = 4096
const BTREE_MAX_KEY_SIZE = 1000
const BTREE_MAX_VAL_SIZE = 3000

func init() {
	// we want to make sure we can fit a node into a page
	node1max := HEADER_SIZE + 8 + 2 + 4 + BTREE_MAX_KEY_SIZE + BTREE_MAX_VAL_SIZE
	errors.Assert(node1max <= BTREE_PAGE_SIZE, "node1max <= BTREE_PAGE_SIZE")
}

const NilPagePtr types.PagePtr = 0
