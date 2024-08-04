package pagemanager

import (
	"trees/pkg/btree/constant"
	"trees/pkg/btree/types"
)

// OnDisk pagemanager manages pages on disk, using mmap
// TODO: figure out the proper interface/interaction between kvstore, btree, and ondisk page manager
type OnDisk struct {
	Path string // file name
	// internals
	fd   int
	mmap struct {
		total  int      // mmap size, can be larger than the file size
		chunks [][]byte // multiple mmaps, can be non-continuous
	}
	page struct {
		flushed uint64   // database size in number of pages
		temp    [][]byte // newly allocated pages
	}
}

var _ PageManager = (*OnDisk)(nil)

func (od *OnDisk) Get(ptr types.PagePtr) []byte {
	start := types.PagePtr(0)
	for _, chunk := range od.mmap.chunks {
		end := start + types.PagePtr(len(chunk))/constant.BTREE_PAGE_SIZE
		if ptr < end {
			offset := constant.BTREE_PAGE_SIZE * (ptr - start)
			return chunk[offset : offset+constant.BTREE_PAGE_SIZE]
		}
		start = end
	}
	panic("bad ptr")
}

func (od *OnDisk) New(node []byte) types.PagePtr {
	ptr := od.page.flushed + uint64(len(od.page.temp)) // just append
	od.page.temp = append(od.page.temp, node)
	return types.PagePtr(ptr)
}

func (od *OnDisk) Del(ptr types.PagePtr) {
	// do nothing
}
