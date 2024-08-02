package btree_serde

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBTreeLookup(t *testing.T) {
	// create a bunch of nodes to use in trees below
	leafNode1 := NewLeafNode()
	leafNode1.keys = []int{1}
	leafNode1.vals = []int{1}
	leafNode1To4 := NewLeafNode()
	leafNode1To4.keys = []int{1, 2, 3, 4}
	leafNode1To4.vals = []int{1, 2, 3, 4}
	leafNode5To8 := NewLeafNode()
	leafNode5To8.keys = []int{5, 6, 7, 8}
	leafNode5To8.vals = []int{5, 6, 7, 8}
	leafNode9To12 := NewLeafNode()
	leafNode9To12.keys = []int{9, 10, 11, 12}
	leafNode9To12.vals = []int{9, 10, 11, 12}
	leafNode13To16 := NewLeafNode()
	leafNode13To16.keys = []int{13, 14, 15, 16}
	leafNode13To16.vals = []int{13, 14, 15, 16}

	internalNode1To8, err := NewInternalNode([]int{5}, []*Node{leafNode1To4, leafNode5To8})
	require.NoError(t, err)
	internalNode9To16, err := NewInternalNode([]int{13}, []*Node{leafNode9To12, leafNode13To16})
	require.NoError(t, err)
	internalNode1To16, err := NewInternalNode([]int{9}, []*Node{internalNode1To8, internalNode9To16})
	require.NoError(t, err)

	// we use trees of order 4 for testing (2-3-4 trees)
	EmptyTree := NewBTree(4)
	Tree1 := NewBTree(4)
	Tree1.root = leafNode1
	Tree1To4 := NewBTree(4)
	Tree1To4.root = leafNode1To4
	Tree1To8 := NewBTree(4)
	Tree1To8.root = internalNode1To8
	Tree1To16 := NewBTree(4)
	Tree1To16.root = internalNode1To16

	t.Run("lookup", func(t *testing.T) {
		tests := []struct {
			name      string
			tree      *BTree
			key       int
			wantVal   int
			wantFound bool
		}{
			{"empty", EmptyTree, 1, 0, false},
			{"1", Tree1, 1, 1, true},
			{"1To4", Tree1To4, 1, 1, true},
			{"1To4", Tree1To4, 4, 4, true},
			{"1To4", Tree1To4, 5, 0, false},
			{"1To8", Tree1To8, 1, 1, true},
			{"1To8", Tree1To8, 4, 4, true},
			{"1To8", Tree1To8, 5, 5, true},
			{"1To8", Tree1To8, 8, 8, true},
			{"1To8", Tree1To8, 9, 0, false},
			{"1To16", Tree1To16, 1, 1, true},
			{"1To16", Tree1To16, 4, 4, true},
			{"1To16", Tree1To16, 5, 5, true},
			{"1To16", Tree1To16, 8, 8, true},
			{"1To16", Tree1To16, 9, 9, true},
			{"1To16", Tree1To16, 12, 12, true},
			{"1To16", Tree1To16, 13, 13, true},
			{"1To16", Tree1To16, 16, 16, true},
			{"1To16", Tree1To16, 17, 0, false},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				gotVal, gotFound := tt.tree.Lookup(tt.key)
				require.Equal(t, tt.wantVal, gotVal)
				require.Equal(t, tt.wantFound, gotFound)
			})
		}
	})
}

func TestBTreeInsert(t *testing.T) {
	// create a bunch of nodes to use in trees below
	leafNode1To4 := NewLeafNode()
	leafNode1To4.keys = []int{1, 2, 3, 4}
	leafNode1To4.vals = []int{1, 2, 3, 4}

	// we use trees of order 4 for testing (2-3-4 trees)
	EmptyTree := NewBTree(4)
	Tree1To4 := NewBTree(4)
	Tree1To4.root = leafNode1To4

	t.Run("insert", func(t *testing.T) {
		tests := []struct {
			name          string
			initialTree   *BTree
			keys          []int
			vals          []int
			wantFinalTree *BTree
		}{
			{"empty", EmptyTree, []int{1, 2, 3, 4}, []int{1, 2, 3, 4}, Tree1To4},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				require.Equal(t, len(tt.keys), len(tt.vals))
				for i, key := range tt.keys {
					tt.initialTree.Insert(key, tt.vals[i])
				}
				if !reflect.DeepEqual(tt.initialTree, tt.wantFinalTree) {
					t.Errorf("got %+v, want %+v", tt.initialTree.root, tt.wantFinalTree.root)
				}
			})
		}
	})
}
