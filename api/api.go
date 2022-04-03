package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adarsh1021/kv-serve/datastore"
	"github.com/gorilla/mux"
)

func StartServer() {
	mux := mux.NewRouter()

	mux.HandleFunc("/", healthcheck)
	mux.HandleFunc("/store", handleCreateStore).Methods("POST")
	mux.HandleFunc("/store", handleListStores).Methods("GET")
	mux.HandleFunc("/store/{store}", handleStorePut).Methods("PUT", "POST") // GET, PUT/POST, DELETE
	mux.HandleFunc("/store/{store}/{key}", handleStoreGet).Methods("GET")

	http.ListenAndServe(":8080", mux)
}

func healthcheck(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "ok")
}

func handleCreateStore(w http.ResponseWriter, req *http.Request) {
	stores, err := datastore.ListStores()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(stores)
}

func handleListStores(w http.ResponseWriter, req *http.Request) {
	var newStore datastore.Store
	err := json.NewDecoder(req.Body).Decode(&newStore)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = datastore.CreateOrOpen(newStore)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Key value store created")
}

func handleStorePut(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	storeName, ok := vars["store"]
	if !ok {
		fmt.Fprintf(w, "storeName missing")
	}

	var kv datastore.KeyValuePair
	err := json.NewDecoder(req.Body).Decode(&kv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = datastore.Put(datastore.Store{Name: storeName}, kv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "ok")
}

func handleStoreGet(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	storeName, _ := vars["store"]
	key := vars["key"]

	data, err := datastore.Get(datastore.Store{Name: storeName}, key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, string(data))
}
