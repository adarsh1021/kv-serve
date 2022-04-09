package api

import (
	context "context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	"github.com/adarsh1021/kv-serve/kvdb"
	"github.com/gorilla/mux"
	grpc "google.golang.org/grpc"
)

type server struct {
	UnimplementedKvServiceServer
}

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

func StartGrpcServer(port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	RegisterKvServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
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

	fmt.Fprintf(w, dbId+" created")
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

func (s *server) CreteDb(ctx context.Context, in *CreateDbRequest) (*CreateDbResponse, error) {
	return &CreateDbResponse{}, nil
}

func (s *server) GetKey(ctx context.Context, in *GetKeyRequest) (*GetKeyResponse, error) {
	return &GetKeyResponse{Key: "Hello", Value: []byte("000000")}, nil
}
