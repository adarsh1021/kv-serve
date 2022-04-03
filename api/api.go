package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adarsh1021/kv-serve/datastore"
)

func StartServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", healthcheck)
	mux.HandleFunc("/store", handleStore) // Create database
	// mux.HandleFunc("/:db_name")       // GET, PUT/POST, DELETE

	http.ListenAndServe(":8080", mux)
}

func healthcheck(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "ok")
}

func handleStore(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		stores, err := datastore.ListStores()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(stores)

	case http.MethodPost:
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
}
