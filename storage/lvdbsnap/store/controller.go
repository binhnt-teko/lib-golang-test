package store

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/syndtr/goleveldb/leveldb"
)

type KV struct {
	Key   string
	Value string
}
type KVStore struct {
	db       *leveldb.DB
	snapshot *leveldb.Snapshot
}

func New(db *leveldb.DB) *KVStore {
	return &KVStore{
		db: db,
	}
}

func (s *KVStore) HandleSetValue(c *gin.Context) {
	var kv KV
	if err := c.BindJSON(&kv); err != nil {
		c.JSON(202, map[string]interface{}{
			"message": fmt.Sprintf("BindJSON error: %s ", err.Error()),
		})
		return
	}

	if err := s.db.Put([]byte(kv.Key), []byte(kv.Value), nil); err != nil {
		c.JSON(202, map[string]interface{}{
			"message": fmt.Sprintf("Put error : %s ", err.Error()),
		})
		return
	}
	c.JSON(200, map[string]interface{}{
		"message": "success",
	})
}
func (s *KVStore) HandleGetValue(c *gin.Context) {
	key := c.Request.URL.Query().Get("key")
	if key == "" {
		c.JSON(202, map[string]interface{}{
			"message": "NOT FIND KEY",
		})
		return
	}
	value, err := s.db.Get([]byte(key), nil)
	if err != nil {
		c.JSON(202, map[string]interface{}{
			"message": fmt.Sprintf("Error: %s ", err.Error()),
		})
		return
	}
	kv := &KV{
		Key:   key,
		Value: string(value),
	}
	c.JSON(200, kv)
	return
}
func (s *KVStore) HandleGetAll(c *gin.Context) {
	list := []*KV{}
	iter := s.db.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		list = append(list, &KV{
			Key:   string(key),
			Value: string(value),
		})
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		c.JSON(202, map[string]interface{}{
			"message": fmt.Sprintf("Iterator Error: %s ", err.Error()),
		})
		return
	}
	c.JSON(200, list)
	return

}
func (s *KVStore) HandleGetAllSnapshot(c *gin.Context) {
	list := []*KV{}
	if s.snapshot == nil {
		c.JSON(202, map[string]interface{}{
			"message": fmt.Sprintf("Must take snapshot before"),
		})
		return
	}
	iter := s.snapshot.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		list = append(list, &KV{
			Key:   string(key),
			Value: string(value),
		})
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		c.JSON(202, map[string]interface{}{
			"message": fmt.Sprintf("Iterator Error: %s ", err.Error()),
		})
		return
	}
	c.JSON(200, list)
	return

}
func (s *KVStore) HandleSnapshot(c *gin.Context) {
	snapshot, err := s.db.GetSnapshot()
	if err != nil {
		c.JSON(202, map[string]interface{}{
			"message": fmt.Sprintf("GetSnapshot error: %s ", err.Error()),
		})
	}
	s.snapshot = snapshot
	c.JSON(200, "Snapshot ok")
	return
}
