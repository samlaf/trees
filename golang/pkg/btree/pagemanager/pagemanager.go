package pagemanager

type PageManager interface {
	// Get a page by its number
	Get(uint64) []byte
	// Allocate a new page and return its number
	New([]byte) uint64
	// Deallocate a page by its number
	Del(uint64)
}
