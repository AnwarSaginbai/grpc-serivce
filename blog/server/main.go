package main

import (
	pb "github.com/AnwarSaginbai/grpc-service/blog/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

const addr = ":8080"

type Server struct {
	pb.BlogServiceServer
}

func main() {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen addr: %v", err)
	}
	log.Printf("listening the addr %s", addr)

	repo, err := setupDB()
	if err != nil {
		log.Fatalf("failed to setup db: %v", err)
	}
	_ = repo
	srv := grpc.NewServer()
	pb.RegisterBlogServiceServer(srv, &Server{})
	if err := srv.Serve(listener); err != nil {
		log.Fatalf("failed to serve the server: %v", err)
	}
}
