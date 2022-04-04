package kvdb

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/bluele/gcache"
	"github.com/syndtr/goleveldb/leveldb"
)

const DATA_PATH string = "/Users/adarsh/projects/kv-serve/data"

var gc gcache.Cache

func SetupDb() {
	gc = gcache.New(100).LRU().Build()
	warmUpCache()
}

func warmUpCache() {
	dbIds, _ := ListDbs()
	for _, dbId := range dbIds {
		dbPath := filepath.FromSlash(DATA_PATH + "/" + dbId)
		db, err := leveldb.OpenFile(dbPath, nil)
		if err != nil {
			log.Panic(err)
		}
		gc.Set(dbId, db)
	}
}

func OnExit() {
	// Cleanup cache items
	cache := gc.GetALL(false)
	for _, i := range cache {
		db := i.(*leveldb.DB)
		db.Close()
	}
}

func open(dbId string) (*leveldb.DB, error) {
	dbPath := filepath.FromSlash(DATA_PATH + "/" + dbId)
	return leveldb.OpenFile(dbPath, nil)
}

func CreateOrOpen(dbId string) (*leveldb.DB, error) {
	var db *leveldb.DB
	var err error = nil

	i, getErr := gc.Get(dbId)
	if getErr != nil {
		db, err = open(dbId)
		if err != nil {
			log.Panic(err)
			return nil, err
		}

		gc.Set(dbId, db)
	} else {
		db = i.(*leveldb.DB)
	}

	return db, err
}

func ListDbs() ([]string, error) {
	files, err := ioutil.ReadDir(DATA_PATH)
	if err != nil {
		log.Panic(err)
		return nil, err
	}

	var dbs []string
	for _, file := range files {
		dbs = append(dbs, file.Name())
	}

	return dbs, nil
}

func Put(dbId string, key []byte, value []byte) error {
	db, err := CreateOrOpen(dbId)
	if err != nil {
		log.Panic(err)
		return err
	}

	return db.Put(key, value, nil)
}

func Get(dbId string, key []byte) ([]byte, error) {
	db, err := CreateOrOpen(dbId)
	if err != nil {
		log.Panic(err)
		return nil, err
	}

	return db.Get(key, nil)
}
