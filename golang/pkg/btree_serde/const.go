package btree_serde

import "trees/internal/errors"

// Node:
// | type | nkeys |  nodePtrs  |   keyOffsets  | key-values | unused |
// |  2B  |   2B  | nkeys * 8B |   nkeys * 2B  |     ...    |        |
//
// Key-Val:
// | klen | vlen | key | val |
// |  2B  |  2B  | ... | ... |
const (
	// header section
	nodeTypeSizeBytes = 2
	numKeysSizeBytes  = 2
	headerSizeBytes   = nodeTypeSizeBytes + numKeysSizeBytes

	// pointers section
	nodePtrsSizeBytes     = 8
	entryOffsetsSizeBytes = 2

	// entries section
	keyLenSizeBytes        = 2
	valLenSizeBytes        = 2
	entriesHeaderSizeBytes = keyLenSizeBytes + valLenSizeBytes

	maxKeySizeBytes = 1000
	maxValSizeBytes = 3000
	pageSizeBytes   = 4096 // 4KB
)

func init() {
	// we make sure that the Sizes set above allow for a node to fit in a page
	node1MaxSize := headerSizeBytes + 1*nodePtrsSizeBytes + 1*entryOffsetsSizeBytes + 1*maxKeySizeBytes + 1*maxValSizeBytes
	errors.Assert(node1MaxSize <= pageSizeBytes, "node1MaxSize <= PageSize")
}

func lenPtrsSection(nkeys int) int {
	return nkeys * (nodePtrsSizeBytes + entryOffsetsSizeBytes)
}
