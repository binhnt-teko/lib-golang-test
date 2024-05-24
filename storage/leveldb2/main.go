package main

import (
	"fmt"
	"log"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

func Recovery(filename string) {
	db, err := leveldb.RecoverFile(filename, &opt.Options{
		Strict: opt.NoStrict,
	})
	if err != nil {
		log.Fatal("RecoverFile failed")
	}
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Printf("Close db error: %s \n", err.Error())

		}
	}()

}
func ReadFile(fileName string) {
	optA := &opt.Options{
		ReadOnly: true,
		// Strict:   opt.NoStrict,
	}
	db, err := leveldb.OpenFile(fileName, optA)
	if err != nil {
		log.Fatal("Yikes!")
	}
	defer db.Close()
	rang := util.BytesPrefix([]byte(""))
	iter := db.NewIterator(rang, &opt.ReadOptions{
		DontFillCache: false,
		Strict:        opt.NoStrict,
	})

	fmt.Println("Start iterator")
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		fmt.Printf("key: %s | value: %s\n", key, value)
	}
	iter.Release()
	err = iter.Error()
	if err != nil {
		fmt.Printf("Error iterator: %s", err.Error())
		return
	}
	fmt.Println("End iterator")
}

func main() {
	fileName := "undone_txs.db"
	// Recovery(fileName)
	ReadFile(fileName)
}
