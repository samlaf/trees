package benchmark

import (
	"crypto/md5"
	"testing"
	"trees"
	"trees/merkle"
	"trees/smt"
)

// Benchmarks
func BenchmarkImplementations(b *testing.B) {
	hashAlgo := md5.New()
	implementations := map[string]trees.Tree{
		"merkle": merkle.New(1<<10, hashAlgo),
		"smt":    smt.NewSMT64(hashAlgo),
	}

	for name, impl := range implementations {
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				impl.GetRoot()
			}
		})
	}
}
