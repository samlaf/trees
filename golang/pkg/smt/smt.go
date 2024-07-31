package smt

import (
	"hash"
	"math/bits"
	"trees"
)

// SMT64 is a sparse Merkle tree
// we use uint64 as key, so max depth is 64
type SMT64 struct {
	hash hash.Hash
	//              1           \
	//           2     3        |                 nodeHashes
	//         4  5   6  7      |
	//           ........      /
	//   2^64 2^64+1    ...     2^65-2 2^65-1     leaveHashes (stored as keys 0,1,2...)
	//    |      |      ...     |      |
	//    0      1              2      3          leaves
	leaves      map[uint64]string
	leaveHashes map[uint64][]byte
	nodeHashes  map[uint64][]byte
	zeroHashes  *[257][]byte
}

var _ trees.Tree = &SMT64{}

func NewSMT64(hash hash.Hash) *SMT64 {
	return &SMT64{
		hash:        hash,
		leaves:      make(map[uint64]string),
		leaveHashes: make(map[uint64][]byte),
		nodeHashes:  make(map[uint64][]byte),
		zeroHashes:  makeZeroHashes(hash),
	}
}

func (smt *SMT64) Insert(key uint64, value string) error {
	smt.leaves[key] = value
	valueHash := smt.hashLeaf(value)
	smt.leaveHashes[key] = valueHash
	// update first layer of nodeHash, right above the leaves
	// special case b/c can't just /2 the leaf key b/c we're at max uint64 bits

	leftKey := key & ^uint64(1) // clear the last bit
	rightKey := leftKey + 1
	leftHash := smt.getLeafHashOrZeroHash(leftKey)
	rightHash := smt.getLeafHashOrZeroHash(rightKey)
	smt.hash.Reset()
	smt.hash.Write(leftHash)
	smt.hash.Write(rightHash)
	nodeHash := smt.hash.Sum(nil)
	key = key>>1 + 1<<63
	smt.nodeHashes[key] = nodeHash
	// update other nodeHashes all the way to the root
	for key > 1 {
		leftKey = key & ^uint64(1)
		rightKey = leftKey + 1
		leftHash = smt.getNodeHashOrZeroHash(leftKey)
		rightHash = smt.getNodeHashOrZeroHash(rightKey)
		smt.hash.Reset()
		smt.hash.Write(leftHash)
		smt.hash.Write(rightHash)
		nodeHash = smt.hash.Sum(nil)
		key = key >> 1
		smt.nodeHashes[key] = nodeHash
	}
	return nil
}

func (smt *SMT64) GetRoot() []byte {
	return smt.getNodeHashOrZeroHash(1)
}

func (smt *SMT64) getLeafHashOrZeroHash(key uint64) []byte {
	if hash, ok := smt.leaveHashes[key]; ok {
		return hash
	}
	return smt.zeroHashes[0]
}

func (smt *SMT64) getNodeHashOrZeroHash(key uint64) []byte {
	if hash, ok := smt.nodeHashes[key]; ok {
		return hash
	}
	return smt.zeroHashes[highestNonZeroBitPosition(key)]
}

func highestNonZeroBitPosition(num uint64) int {
	if num == 0 {
		return -1 // Return -1 if there is no non-zero bit
	}
	return bits.Len64(num) - 1
}

func (smt *SMT64) hashLeaf(value string) []byte {
	smt.hash.Reset()
	smt.hash.Write([]byte(value))
	return smt.hash.Sum(nil)
}

func makeZeroHashes(h hash.Hash) *[257][]byte {
	zeroHashes := new([257][]byte)
	hash := make([]byte, 32)
	h.Write(hash)
	hash = h.Sum(nil)
	zeroHashes[0] = hash
	for i := 1; i < 257; i++ {
		h.Reset()
		h.Write(hash) // left element
		h.Write(hash) // right element
		hash = h.Sum(nil)
		zeroHashes[i] = hash
	}
	return zeroHashes
}
