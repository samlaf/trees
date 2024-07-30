package main

import (
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

func main() {
	db, err := leveldb.OpenFile("leveldb", &opt.Options{ErrorIfMissing: true})
	panicOnError(err)
	defer db.Close()

	err = db.Put([]byte("key"), []byte("value"), nil)
	panicOnError(err)
	data, err := db.Get([]byte("key"), nil)
	panicOnError(err)
	fmt.Println(string(data))
	// err = db.Delete([]byte("key"), nil)
	// panicOnError(err)

}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
