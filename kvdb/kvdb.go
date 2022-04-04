package kvdb

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/syndtr/goleveldb/leveldb"
)

type KvDb struct {
	Name string `json:"name"`
	Ref  interface{}
}

type KvPair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

const DATA_PATH string = "/Users/adarsh/projects/kv-serve/data"

var db *leveldb.DB

func CreateOrOpen(dbId string) (*leveldb.DB, error) {
	dbPath := filepath.FromSlash(DATA_PATH + "/" + dbId)

	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	return db, err
}

func Open(dbId string) (*leveldb.DB, error) {
	storePath := filepath.FromSlash(DATA_PATH + "/" + dbId)

	db, err := leveldb.OpenFile(storePath, nil)
	if err != nil {
		log.Panic(err)
	}

	return db, err
}

func ListDbs() ([]KvDb, error) {
	files, err := ioutil.ReadDir(DATA_PATH)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var dbs []KvDb
	for _, file := range files {
		dbs = append(dbs, KvDb{Name: file.Name()})
	}

	return dbs, nil
}

func Put(dbId string, kv KvPair) error {
	db, _ := Open(dbId)
	defer db.Close()
	return db.Put([]byte(kv.Key), []byte(kv.Value), nil)
}

func Get(dbId string, key string) ([]byte, error) {
	db, _ := Open(dbId)
	defer db.Close()
	return db.Get([]byte(key), nil)
}
