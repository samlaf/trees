package crypto

import "hash"

func Hash(hasher hash.Hash, data []byte) []byte {
	hasher.Reset()
	hasher.Write(data)
	return hasher.Sum(nil)
}
