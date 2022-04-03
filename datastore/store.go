package datastore

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/syndtr/goleveldb/leveldb"
)

type Store struct {
	Name string `json:"name"`
}

type KeyValuePair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

const DATA_PATH string = "/Users/adarsh/projects/kv-serve/data"

func CreateOrOpen(store Store) (*leveldb.DB, error) {
	storePath := filepath.FromSlash(DATA_PATH + "/" + store.Name)

	db, err := leveldb.OpenFile(storePath, nil)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	return db, err
}

func Open(store Store) (*leveldb.DB, error) {
	storePath := filepath.FromSlash(DATA_PATH + "/" + store.Name)

	db, err := leveldb.OpenFile(storePath, nil)
	if err != nil {
		log.Panic(err)
	}

	return db, err
}

func ListStores() ([]Store, error) {
	files, err := ioutil.ReadDir(DATA_PATH)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var stores []Store
	for _, file := range files {
		stores = append(stores, Store{Name: file.Name()})
	}

	return stores, nil
}

func Put(store Store, kv KeyValuePair) error {
	storedb, _ := Open(store)
	defer storedb.Close()
	return storedb.Put([]byte(kv.Key), []byte(kv.Value), nil)
}

func Get(store Store, key string) ([]byte, error) {
	storedb, _ := Open(store)
	defer storedb.Close()
	return storedb.Get([]byte(key), nil)
}
