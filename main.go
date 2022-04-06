package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/adarsh1021/kv-serve/api"
	"github.com/adarsh1021/kv-serve/kvdb"
)

func main() {
	var dataDir string
	const (
		dataDirDefault = "/data"
		dataDirUsage   = "Path to data directory"
	)
	flag.StringVar(&dataDir, "data-dir", dataDirDefault, dataDirUsage)
	flag.StringVar(&dataDir, "d", dataDirDefault, dataDirUsage+" (shorthand)")

	var maxDbCacheEntries int
	const (
		maxDbCacheEntriesDefault = 100
		maxDbCacheEntriesUsage   = "Maximun number of databases that can be open at a time"
	)
	flag.IntVar(&maxDbCacheEntries, "max-db-cache-entries", maxDbCacheEntriesDefault, maxDbCacheEntriesUsage)
	flag.IntVar(&maxDbCacheEntries, "c", maxDbCacheEntriesDefault, maxDbCacheEntriesUsage+" (shorthand)")

	var port int
	const (
		portDefault = 9090
		portUsage   = "Server port"
	)
	flag.IntVar(&port, "port", portDefault, portUsage)
	flag.IntVar(&port, "p", portDefault, portUsage+" (shorthand)")

	flag.Parse()

	log.Println("Data dir: " + dataDir)
	log.Println("Max entries in db cache: " + fmt.Sprint(maxDbCacheEntries))

	kvdb.SetupDb(dataDir, maxDbCacheEntries)
	api.StartServer(port)
	kvdb.OnExitCleanup()
}
