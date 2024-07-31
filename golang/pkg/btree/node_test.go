package btree

import "testing"

func TestNode(t *testing.T) {
	emptyLeafNode := NewLeafBNode(nil)
	leafNode := NewLeafBNode([]*Entry{
		{Key: []byte("a"), Value: []byte("1")},
		{Key: []byte("b"), Value: []byte("2")},
	})
	t.Run("leaf node constructor", func(t *testing.T) {
		wantLeafNode := []byte{
			// header
			1, 0,
			2, 0,
			// nodePtrs
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			// entryOffsets
			24, 0,
			30, 0,
			// entries
			1, 0, 1, 0, 97, 49,
			1, 0, 1, 0, 98, 50,
		}
		if string(leafNode) != string(wantLeafNode) {
			t.Errorf("expected leaf node to be \n%v, got \n%v", wantLeafNode, leafNode)
		}
	})

	t.Run("empty leaf node lookup", func(t *testing.T) {
		idx := emptyLeafNode.LookupLE([]byte("a"))
		if idx != 0 {
			t.Errorf("expected not to find anything")
		}
	})
	t.Run("numKeys", func(t *testing.T) {
		if emptyLeafNode.NumKeys() != 0 {
			t.Errorf("expected 0 keys, got %d", emptyLeafNode.NumKeys())
		}
		if leafNode.NumKeys() != 2 {
			t.Errorf("expected 2 keys, got %d", leafNode.NumKeys())
		}
	})
	t.Run("entryOffsetUnchecked", func(t *testing.T) {
		gotOffset := leafNode.entryOffsetUnchecked(0)
		wantOffset := headerSizeBytes + 2*(nodePtrsSizeBytes+entryOffsetsSizeBytes) + 0
		if gotOffset != headerSizeBytes+2*(nodePtrsSizeBytes+entryOffsetsSizeBytes) {
			t.Errorf("expected offset %d, got %d", wantOffset, gotOffset)
		}
		gotOffset = leafNode.entryOffsetUnchecked(1)
		wantOffset = headerSizeBytes + 2*(nodePtrsSizeBytes+entryOffsetsSizeBytes) + 6
		if gotOffset != uint16(wantOffset) {
			t.Errorf("expected offset %d, got %d", wantOffset, gotOffset)
		}
	})
	t.Run("get key", func(t *testing.T) {
		key := leafNode.KeyUnchecked(0)
		if string(key) != "a" {
			t.Errorf("expected key a, got %s", key)
		}
		key = leafNode.KeyUnchecked(1)
		if string(key) != "b" {
			t.Errorf("expected key b, got %s", key)
		}
	})
	t.Run("leaf node lookup", func(t *testing.T) {
		idx := leafNode.LookupLE([]byte("a"))
		if idx != 0 {
			t.Errorf("expected to find a at idx 0, found %d", idx)
		}
		idx = leafNode.LookupLE([]byte("b"))
		if idx != 1 {
			t.Errorf("expected to find b at idx 1, found %d", idx)
		}
		idx = leafNode.LookupLE([]byte("c"))
		if idx != 2 {
			t.Errorf("expected to find c at idx 2, found %d", idx)
		}
	})

}
