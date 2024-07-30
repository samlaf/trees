package benchmark

import (
	"fmt"
	"os"
	"testing"
	"trees/internal/projectpath"

	"github.com/syndtr/goleveldb/leveldb"
)

// Benchmarks
func BenchmarkLevelDBRandomWrite(b *testing.B) {
	tmpDirPath, err := os.MkdirTemp(projectpath.Root+"/.data/benchmark", "leveldb-random-write-*")
	requireNoErr(b, err)
	fmt.Println("tmpDirPath:", tmpDirPath)
	defer os.RemoveAll(tmpDirPath)
	db, err := leveldb.OpenFile(tmpDirPath, nil)
	requireNoErr(b, err)
	for i := 0; i < b.N; i++ {
		err = db.Put([]byte("key"), []byte("value"), nil)
		requireNoErr(b, err)
	}
}

func requireNoErr(b *testing.B, err error) {
	if err != nil {
		b.Fatal("unexpected error...", err)
	}
}
