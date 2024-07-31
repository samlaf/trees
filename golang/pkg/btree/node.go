package btree

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// A node includes:
// 1. A fixed-size header, which contains (little-endian encoded):
//      1.1 The type of node (leaf or internal).
//      1.2 The number of keys.
// 2. A list of pointers to child nodes for internal nodes.
// 3. A list of offsets (internal pointers) to KVs, which can be used to binary search KVs.
// 4. A list of KV pairs.
// | type | nkeys |  nodePtrs  |   keyPtrs  | entries | unused |
// |  2B  |   2B  | nkeys * 8B | nkeys * 2B |     ...    |        |
//
// with the entries section containing key-values of arbitrary key/value lengths:
// | klen | vlen | key | val |
// |  2B  |  2B  | ... | ... |

type NodeType uint16 // 2 bytes
const (
	InternalNodeType NodeType = iota
	LeafNodeType
)

type BNode []byte // can be dumped to the disk
type Key []byte
type Value []byte
type NodePtr uint64
type Entry struct {
	Key   Key
	Value Value
}

func (e *Entry) Len() int {
	return entriesHeaderSizeBytes + len(e.Key) + len(e.Value)
}

func NewInternalBNode(nodeptrs []NodePtr) BNode {
	bNode := make(BNode, headerSizeBytes+lenPtrsSection(len(nodeptrs)))
	binary.LittleEndian.PutUint16(bNode[:2], uint16(InternalNodeType))
	binary.LittleEndian.PutUint16(bNode[2:4], uint16(len(nodeptrs)))
	for i, ptr := range nodeptrs {
		offset := headerSizeBytes + i*8
		binary.LittleEndian.PutUint64(bNode[offset:offset+8], uint64(ptr))
	}
	return bNode
}

func NewLeafBNode(entries []*Entry) BNode {
	headerAndPtrsSectionSize := headerSizeBytes + lenPtrsSection(len(entries))
	entryOffsets := make([]int, len(entries))
	entriesSize := 0
	prevEntryLen := 0
	for i, entry := range entries {
		if i == 0 {
			entryOffsets[i] = headerAndPtrsSectionSize
		} else {
			entryOffsets[i] = entryOffsets[i-1] + prevEntryLen
		}
		prevEntryLen = entry.Len()
		entriesSize += prevEntryLen
	}
	bNodeSize := headerAndPtrsSectionSize + entriesSize
	bNode := make(BNode, bNodeSize)

	binary.LittleEndian.PutUint16(bNode[:2], uint16(LeafNodeType))
	binary.LittleEndian.PutUint16(bNode[2:4], uint16(len(entries)))
	for i, entry := range entries {
		entryHeaderOffset := headerSizeBytes + len(entries)*nodePtrsSizeBytes + i*entryOffsetsSizeBytes
		binary.LittleEndian.PutUint16(bNode[entryHeaderOffset:entryHeaderOffset+keyLenSizeBytes], uint16(entryOffsets[i]))
		entryOffset := entryOffsets[i]
		binary.LittleEndian.PutUint16(bNode[entryOffset:entryOffset+keyLenSizeBytes], uint16(len(entry.Key)))
		binary.LittleEndian.PutUint16(bNode[entryOffset+keyLenSizeBytes:entryOffset+keyLenSizeBytes+valLenSizeBytes], uint16(len(entry.Value)))
		entryOffset += entriesHeaderSizeBytes
		copy(bNode[entryOffset:entryOffset+len(entry.Key)], entry.Key)
		copy(bNode[entryOffset+len(entry.Key):entryOffset+len(entry.Key)+len(entry.Value)], entry.Value)
	}
	return bNode
}

func (n BNode) Type() NodeType {
	return NodeType(binary.LittleEndian.Uint16(n[:2]))
}

func (n BNode) NumKeys() uint16 {
	return binary.LittleEndian.Uint16(n[2:4])
}

func (n BNode) NodePointer(i uint16) (uint64, error) {
	if i >= n.NumKeys() {
		return 0, fmt.Errorf("index %d out of range (numKeys=%d)", i, n.NumKeys())
	}
	offset := headerSizeBytes + i*nodePtrsSizeBytes
	return binary.LittleEndian.Uint64(n[offset : offset+8]), nil
}

// returns the offset of the i-th key-val (including header) in the key-values section
// make sure i < NumKeys() before calling this function
func (n BNode) entryOffsetUnchecked(i uint16) uint16 {
	keyPtrOffset := headerSizeBytes + n.NumKeys()*nodePtrsSizeBytes + i*entryOffsetsSizeBytes
	entryOffset := binary.LittleEndian.Uint16(n[keyPtrOffset : keyPtrOffset+entryOffsetsSizeBytes])
	return entryOffset
}

// returns the key of the i-th key-val
// make sure i < NumKeys() before calling this function
func (n BNode) KeyUnchecked(i uint16) []byte {
	entryOffset := n.entryOffsetUnchecked(i)
	keyLen := binary.LittleEndian.Uint16(n[entryOffset : entryOffset+keyLenSizeBytes])
	key := n[entryOffset+entriesHeaderSizeBytes : entryOffset+entriesHeaderSizeBytes+keyLen]
	return key
}

// returns the first kid node whose range intersects the key. (kid[i] <= key)
// TODO: binary search
func (n BNode) LookupLE(key []byte) uint16 {
	numKids := n.NumKeys()
	for i := uint16(0); i < numKids; i++ {
		if bytes.Compare(key, n.KeyUnchecked(i)) <= 0 {
			return i
		} else {
			continue
		}
	}
	return numKids
}

// ================== 4.2 NODE UPDATE =================
// add a new key to a leaf node
func leafInsert(new BNode, old BNode, idx uint16, key []byte, val []byte) {
}

func nodeAppendKV(new BNode, idx uint16, ptr uint64, key []byte, val []byte) {
}

// ================== 4.3 NODE SPLIT =================

// split a oversized node into 2 so that the 2nd node always fits on a page
func nodeSplit2(left BNode, right BNode, old BNode) {
	// code omitted...
}

// split a node if it's too big. the results are 1~3 nodes.
func nodeSplit3(old BNode) (uint16, [3]BNode) {
	// code omitted...
	return 0, [3]BNode{}
}
