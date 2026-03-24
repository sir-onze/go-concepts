package main

import (
	"context"
	"log"
	"net"
	"sync"

	pb "mini-platform/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedStoreServiceServer
	mu    sync.Mutex
	store map[string]string
}

func (s *server) Store(ctx context.Context, req *pb.StoreRequest) (*pb.StoreResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.store[req.Id] = req.Data

	return &pb.StoreResponse{Success: true}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterStoreServiceServer(grpcServer, &server{
		store: make(map[string]string),
	})

	log.Printf("store-service listening on :50051")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
