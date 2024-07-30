package merkle

import (
	"bytes"
	"crypto/md5"
	"testing"
)

type Entry struct {
	Key   uint64
	Value string
}

func TestImplementations(t *testing.T) {
	hasher := md5.New()

	tests := []struct {
		name     string
		height   uint64
		entries  []Entry
		wantRoot []byte
	}{
		{"empty tree of height 0 (only root)", 0, []Entry{}, []byte{212, 29, 140, 217, 143, 0, 178, 4, 233, 128, 9, 152, 236, 248, 66, 126}},
		{"empty tree of height 2", 2, []Entry{}, []byte{15, 136, 147, 220, 10, 175, 253, 145, 41, 35, 37, 133, 118, 164, 177, 71}},
		{"full tree of height 2", 2, []Entry{{0, "0"}, {1, "1"}, {2, "2"}, {3, "3"}}, []byte{95, 251, 160, 76, 58, 92, 236, 191, 59, 48, 197, 1, 164, 232, 200, 20}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mt := NewFromHeight(tt.height, hasher)
			for _, entry := range tt.entries {
				err := mt.Insert(entry.Key, entry.Value)
				if err != nil {
					t.Errorf("Insert() error = %v, want nil", err)
				}
			}
			gotRoot := mt.GetRoot()
			if !bytes.Equal(gotRoot, tt.wantRoot) {
				t.Errorf("GetRoot() = %v, want %v", gotRoot, tt.wantRoot)
			}
		})
	}
}
