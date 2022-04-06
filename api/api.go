package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/adarsh1021/kv-serve/kvdb"
	"github.com/gorilla/mux"
)

func StartServer(port int) {
	log.Println("Starting server...")

	mux := mux.NewRouter()

	mux.HandleFunc("/", healthcheck)
	mux.HandleFunc("/db", handleListDbs).Methods("GET")
	mux.HandleFunc("/db/{dbId}", handleCreateDb).Methods("POST")
	mux.HandleFunc("/db/{dbId}/{key}", handleGetKey).Methods("GET")
	mux.HandleFunc("/db/{dbId}/{key}", handlePutKey).Methods("PUT", "POST")

	log.Println("Server ready...")

	http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}

func healthcheck(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "ok")
}

func handleCreateDb(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	dbId, ok := vars["dbId"]
	if !ok {
		fmt.Fprintf(w, "dbId missing")
	}

	_, err := kvdb.CreateOrOpen(dbId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "kvdb created")
}

func handleListDbs(w http.ResponseWriter, req *http.Request) {
	dbs, err := kvdb.ListDbs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(dbs)
}

func handlePutKey(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	dbId, ok := vars["dbId"]
	if !ok {
		fmt.Fprintf(w, "dbId missing")
	}

	key, ok := vars["key"]
	if !ok {
		fmt.Fprintf(w, "key missing")
	}

	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = kvdb.Put(dbId, []byte(key), bodyBytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// log.Println("Received PUT " + key + " at " + time.Now().UTC().Format(time.RFC3339Nano))
	fmt.Fprintf(w, "ok")
}

func handleGetKey(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	dbId, ok := vars["dbId"]
	if !ok {
		fmt.Fprintf(w, "dbId missing")
	}

	key, ok := vars["key"]
	if !ok {
		fmt.Fprintf(w, "key missing")
	}

	data, err := kvdb.Get(dbId, []byte(key))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, string(data))
}
