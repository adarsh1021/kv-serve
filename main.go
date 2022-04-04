package main

import (
	"github.com/adarsh1021/kv-serve/api"
	"github.com/adarsh1021/kv-serve/kvdb"
)

func main() {
	kvdb.SetupDb()
	api.StartServer()
	kvdb.OnExit()
}
