package kvdb

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/bluele/gcache"
	"github.com/syndtr/goleveldb/leveldb"
)

// type KvDb struct {
// 	Name string `json:"name"`
// 	Id   string `json:"id"`
// 	Ref  *leveldb.DB
// }

type KvPair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

const DATA_PATH string = "/Users/adarsh/projects/kv-serve/data"

var gc gcache.Cache = gcache.New(100).LRU().Build()

func SetupDb() {

}

func CreateOrOpen(dbId string) (*leveldb.DB, error) {
	var db *leveldb.DB
	i, err := gc.Get(dbId)
	if err != nil {
		dbPath := filepath.FromSlash(DATA_PATH + "/" + dbId)
		db, err = leveldb.OpenFile(dbPath, nil)
		gc.Set(dbId, db)
	} else {
		db = i.(*leveldb.DB)
	}

	if err != nil {
		log.Panic(err)
	}

	return db, err
}

func ListDbs() ([]string, error) {
	files, err := ioutil.ReadDir(DATA_PATH)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var dbs []string
	for _, file := range files {
		dbs = append(dbs, file.Name())
	}

	return dbs, nil
}

func Put(dbId string, kv KvPair) error {
	// var i interface{}
	i, _ := gc.Get(dbId)
	db := i.(*leveldb.DB)

	return db.Put([]byte(kv.Key), []byte(kv.Value), nil)
}

func Get(dbId string, key string) ([]byte, error) {
	// db, _ := Open(dbId)
	// defer db.Close()
	i, _ := gc.Get(dbId)
	db := i.(*leveldb.DB)

	return db.Get([]byte(key), nil)
}
