package kvdb

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/bluele/gcache"
	"github.com/syndtr/goleveldb/leveldb"
)

var DATA_DIR string

var gc gcache.Cache

func SetupDb(dataDir string, maxDbCacheEntries int) {
	log.Println("Starting db...")
	DATA_DIR = dataDir
	gc = gcache.New(maxDbCacheEntries).
		LRU().
		EvictedFunc(func(key, value interface{}) {
			// Close the db reference on eviction
			db := value.(*leveldb.DB)
			db.Close()
		}).
		Build()
	warmUpCache()
	log.Println("Ready")
}

func warmUpCache() {
	dbIds, _ := ListDbs()
	for _, dbId := range dbIds {
		dbPath := filepath.FromSlash(DATA_DIR + "/" + dbId)
		db, err := leveldb.OpenFile(dbPath, nil)
		if err != nil {
			log.Panic(err)
		}
		gc.Set(dbId, db)
	}
}

func OnExitCleanup() {
	log.Println("Cleaning up...")
	// Cleanup cache items
	cache := gc.GetALL(false)
	for _, i := range cache {
		db := i.(*leveldb.DB)
		db.Close()
	}
}

func open(dbId string) (*leveldb.DB, error) {
	dbPath := filepath.FromSlash(DATA_DIR + "/" + dbId)
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
	files, err := ioutil.ReadDir(DATA_DIR)
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
