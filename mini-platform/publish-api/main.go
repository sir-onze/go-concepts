package main

import (
	"encoding/json"
	"log"
	"net/http"

	"mini-platform/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PublishRequest struct {
	ID   string `json:"id"`
	Data string `json:"data"`
}

func publishHandler(w http.ResponseWriter, r *http.Request) {
	var req PublishRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	// create the grpc connection
	conn, err := grpc.Dial(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		http.Error(w, "could not connect to store", http.StatusInternalServerError)
		return
	}
	defer conn.Close()
	// create the grpc client
	client := proto.NewStoreServiceClient(conn)
	_, err = client.Store(r.Context(), &proto.StoreRequest{Id: req.ID, Data: req.Data})
	if err != nil {
		http.Error(w, "could not store data", http.StatusInternalServerError)
		return
	}
	// call the Store method
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("published"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/publish", publishHandler)
	log.Println("publish-api listening on :8080")
	http.ListenAndServe(":8080", mux)

}
