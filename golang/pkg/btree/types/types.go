package types

// PagePtr is a pointer to a page and/or node in the BTree
// It is defined for clarity and self-documentation (for eg of the pagemanager interface),
// but is always expected to be an uint64, as it is also used in the in-memory node implementation.
// See for eg bnode/bnode.go GetPtr/SetPtr which explicitly convert PagePtr to uint64.
type PagePtr uint64
