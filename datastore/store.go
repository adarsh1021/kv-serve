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

const DATA_PATH string = "/Users/adarsh/projects/kv-serve/data"

func CreateOrOpen(store Store) (*leveldb.DB, error) {
	storePath := filepath.FromSlash(DATA_PATH + "/" + store.Name)
	return leveldb.OpenFile(storePath, nil)
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
