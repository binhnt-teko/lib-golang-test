package main

import (
	"log"

	"github.com/blcvn/lib-golang-test/storage/lvdbsnap/store"
	"github.com/gin-gonic/gin"
	"github.com/syndtr/goleveldb/leveldb"
)

func main() {
	db, err := leveldb.OpenFile("/tmp/foo.db", nil)
	if err != nil {
		log.Fatal("Yikes!")
	}
	defer db.Close()

	router := gin.Default()

	kvStore := store.New(db)

	router.POST("/data", kvStore.HandleSetValue)
	router.GET("/data", kvStore.HandleGetValue)
	router.GET("/all", kvStore.HandleGetAll)
	router.GET("/all-snapshot", kvStore.HandleGetAllSnapshot)
	router.POST("/snapshot", kvStore.HandleSnapshot)

	router.Run(":8080")
}
