package trees

// Tree is the interface for any cryptographic tree
type Tree interface {
	Insert(key uint64, value string) error
	GetRoot() []byte
}
