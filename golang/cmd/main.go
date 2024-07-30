package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/celestiaorg/nmt"
	"github.com/celestiaorg/nmt/namespace"
)

func main() {
	// the tree will use this namespace size (number of bytes)
	nidSize := 1
	// the leaves that will be pushed
	data := [][]byte{
		append(namespace.ID{0}, []byte("leaf_0")...),
		append(namespace.ID{0}, []byte("leaf_1")...),
		append(namespace.ID{1}, []byte("leaf_2")...),
		append(namespace.ID{1}, []byte("leaf_3")...)}
	// Init a tree with the namespace size as well as
	// the underlying hash function:
	tree := nmt.New(sha256.New(), nmt.NamespaceIDSize(nidSize))
	for _, d := range data {
		if err := tree.Push(d); err != nil {
			panic(fmt.Sprintf("unexpected error: %v", err))
		}
	}
	// compute the root
	root, err := tree.Root()
	if err != nil {
		panic(fmt.Sprintf("unexpected error: %v", err))
	}
	// the root's min/max namespace is the min and max namespace of all leaves:
	minNS := nmt.MinNamespace(root, tree.NamespaceSize())
	maxNS := nmt.MaxNamespace(root, tree.NamespaceSize())
	if bytes.Equal(minNS, namespace.ID{0}) {
		fmt.Printf("Min namespace: %x\n", minNS)
	}
	if bytes.Equal(maxNS, namespace.ID{1}) {
		fmt.Printf("Max namespace: %x\n", maxNS)
	}

	// compute proof for namespace 0:
	proof, err := tree.ProveNamespace(namespace.ID{0})
	if err != nil {
		panic("unexpected error")
	}

	// verify proof using the root and the leaves of namespace 0:
	leafs := [][]byte{
		append(namespace.ID{0}, []byte("leaf_0")...),
		append(namespace.ID{0}, []byte("leaf_1")...),
	}

	if proof.VerifyNamespace(sha256.New(), namespace.ID{0}, leafs, root) {
		fmt.Printf("Successfully verified namespace: %x %x\n", namespace.ID{0}, ID{0})
	}

	if proof.VerifyNamespace(sha256.New(), namespace.ID{2}, leafs, root) {
		panic(fmt.Sprintf("Proof for namespace %x, passed for namespace: %x\n", namespace.ID{0}, namespace.ID{2}))
	}
}

type ID []byte

func (nid ID) String() string {
	return hex.EncodeToString(nid)
}
