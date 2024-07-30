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

	fmt.Printf("writing %d items to db\n", b.N)
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := fmt.Sprintf("value-%d", i)
		err = db.Put([]byte(key), []byte(value), nil)
		requireNoErr(b, err)
	}

	iter := db.NewIterator(nil, nil)
	fmt.Printf("Reading all items from db\n")
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		fmt.Println("key:", string(key), "value:", string(value))
	}
}

func requireNoErr(b *testing.B, err error) {
	if err != nil {
		b.Fatal("unexpected error...", err)
	}
}
