package merkle

import (
	"fmt"
	"hash"
	"trees"
	"trees/internal/crypto"
	"trees/internal/math"
)

type MerkleTree struct {
	// number of leaves
	n      uint64
	leaves []string
	hasher hash.Hash
}

var _ trees.Tree = (*MerkleTree)(nil)

func New(n uint64, hasher hash.Hash) *MerkleTree {
	n = math.NextPowerOf2(n)
	return &MerkleTree{
		n:      n,
		leaves: make([]string, n),
		hasher: hasher,
	}
}

func NewFromHeight(height uint64, hasher hash.Hash) *MerkleTree {
	n := 1 << uint64(height)
	return New(uint64(n), hasher)
}

func (mt *MerkleTree) Insert(key uint64, value string) error {
	if key > mt.n {
		return fmt.Errorf("key %d is greater than max number of leaves %d", key, mt.n)
	}
	mt.leaves[key] = value
	return nil
}

func (mt *MerkleTree) GetRoot() []byte {
	return mt.getRootRecursive(mt.leaves)
}

func (mt *MerkleTree) getRootRecursive(elements []string) []byte {
	if len(elements) == 1 {
		return crypto.Hash(mt.hasher, []byte(elements[0]))
	}
	mid := len(elements) / 2
	leftHash := mt.getRootRecursive(elements[:mid])
	rightHash := mt.getRootRecursive(elements[mid:])
	return crypto.Hash(mt.hasher, append(leftHash, rightHash...))
}
