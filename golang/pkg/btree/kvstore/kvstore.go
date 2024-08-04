package kvstore

import (
	"encoding/binary"
	"fmt"
	"os"
	"path"
	"syscall"
	"trees/pkg/btree"
	"trees/pkg/btree/constant"
)

type KV struct {
	Path string // file name
	// internals
	fd   int
	tree btree.BTree
	// we use mmap to READ the underlying db file
	// we use multiple mmap "chunks" (different mmap calls each) instead of
	// a single one that we grow with mremap because mremap can potentially change
	// the address of the memory region, which would invalidate the pointers
	// and hinder concurrent readers
	mmap struct {
		// mmap size, can be larger than the file size
		// sum of all mmap chunks sizes
		totalSizeBytes int
		chunks         [][]byte // multiple mmaps, can be non-continuous
	}
	// we use pages to WRITE to the underlying db file
	// remember that the btree is implemented using copy-on-write, so any modification
	// results in new nodes/pages being created: everytime btree calls 'new',
	// to create a new node, it gets appended to temp and eventually flushed to disk
	pages struct {
		flushed uint64   // database size in number of pages
		temp    [][]byte // newly allocated pages
	}
}

// KV is a wrapper around a BTree which first updates the in-memory BTree structs,
// and then writes the changes to the file.
func (db *KV) Get(key []byte) ([]byte, bool) {
	return db.tree.Get(key)
}
func (db *KV) Set(key []byte, val []byte) error {
	db.tree.Insert(key, val)
	return updateFile(db)
}
func (db *KV) Del(key []byte) (bool, error) {
	deleted := db.tree.Delete(key)
	return deleted, updateFile(db)
}

func updateFile(db *KV) error {
	// 1. Write new nodes.
	if err := writePages(db); err != nil {
		return err
	}
	// 2. `fsync` to enforce the order between 1 and 3.
	if err := syscall.Fsync(db.fd); err != nil {
		return err
	}
	// 3. Update the root pointer atomically.
	if err := writeMetaPage(db); err != nil {
		return err
	}
	// 4. `fsync` to make everything persistent.
	return syscall.Fsync(db.fd)
}

func createFileSync(file string) (int, error) {
	// obtain the directory fd
	dirflags := os.O_RDONLY | syscall.O_DIRECTORY
	dirfd, err := syscall.Open(path.Dir(file), dirflags, 0o644)
	if err != nil {
		return -1, fmt.Errorf("open directory: %w", err)
	}
	defer syscall.Close(dirfd)

	// open or create the file
	flags := os.O_RDWR | os.O_CREATE
	err = syscall.Chdir(path.Dir(file)) // change the working directory
	if err != nil {
		return -1, fmt.Errorf("chdir: %w", err)
	}
	fd, err := syscall.Open(path.Base(file), flags, 0o644)
	if err != nil {
		return -1, fmt.Errorf("open file: %w", err)
	}
	// fsync the directory
	if err = syscall.Fsync(dirfd); err != nil {
		_ = syscall.Close(fd) // may leave an empty file
		return -1, fmt.Errorf("fsync directory: %w", err)
	}
	return fd, nil
}

// maybeCreateNewMmapChunk extends the mmap region by creating a new chunk, if needed.
// the new chunk is at least 64 MB, and otherwise a power of 2 of the current total size.
func maybeCreateNewMmapChunk(db *KV, size int) error {
	if size <= db.mmap.totalSizeBytes {
		return nil // enough range
	}
	newMmapLen := max(db.mmap.totalSizeBytes, 64<<20)
	for db.mmap.totalSizeBytes+newMmapLen < size {
		newMmapLen *= 2 // still not enough?
	}
	chunk, err := syscall.Mmap(
		db.fd, int64(db.mmap.totalSizeBytes), newMmapLen,
		syscall.PROT_READ, syscall.MAP_SHARED, // read-only
	)
	if err != nil {
		return fmt.Errorf("mmap: %w", err)
	}
	db.mmap.totalSizeBytes += newMmapLen
	db.mmap.chunks = append(db.mmap.chunks, chunk)
	return nil
}

func (db *KV) pageAppend(node []byte) uint64 {
	ptr := db.pages.flushed + uint64(len(db.pages.temp)) // just append
	db.pages.temp = append(db.pages.temp, node)
	return ptr
}

func writePages(db *KV) error {
	// extend the mmap if needed
	size := (int(db.pages.flushed) + len(db.pages.temp)) * constant.BTREE_PAGE_SIZE
	if err := maybeCreateNewMmapChunk(db, size); err != nil {
		return err
	}
	// write data pages to the file
	offset := int64(db.pages.flushed * constant.BTREE_PAGE_SIZE)
	for _, tempPage := range db.pages.temp {
		n, err := syscall.Pwrite(db.fd, tempPage, offset)
		if err != nil {
			return err
		}
		if n != len(tempPage) {
			return fmt.Errorf("incomplete write: wrote %d bytes instead of %d", n, len(tempPage))
		}
		offset += int64(n)
	}

	// discard in-memory data
	db.pages.flushed += uint64(len(db.pages.temp))
	db.pages.temp = db.pages.temp[:0]
	return nil
}

// META PAGE STUFF
const DB_SIG = "BuildYourOwnDB06" // not compatible between chapters

// | sig | root_ptr | page_used |
// | 16B |    8B    |     8B    |
func serializeMeta(db *KV) []byte {
	var data [32]byte
	copy(data[:16], []byte(DB_SIG))
	binary.LittleEndian.PutUint64(data[16:], uint64(db.tree.RootPtr))
	binary.LittleEndian.PutUint64(data[24:], db.pages.flushed)
	return data[:]
}

func loadMeta(db *KV, data []byte)

func readRoot(db *KV, fileSize int64) error {
	if fileSize == 0 { // empty file
		db.pages.flushed = 1 // the meta page is initialized on the 1st write
		return nil
	}
	// read the page
	data := db.mmap.chunks[0]
	loadMeta(db, data)
	// verify the page
	// ...
	return nil
}

// 3. Update the meta page. it must be atomic.
func writeMetaPage(db *KV) error {
	if _, err := syscall.Pwrite(db.fd, serializeMeta(db), 0); err != nil {
		return fmt.Errorf("write meta page: %w", err)
	}
	return nil
}
