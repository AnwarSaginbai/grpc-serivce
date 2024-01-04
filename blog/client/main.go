package main

import (
	pb "github.com/AnwarSaginbai/grpc-service/blog/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

const addr = ":8080"

func main() {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewBlogServiceClient(conn)

	//doCreateBlog(client)
	//doGetByIdBlog(client, "6596ea901c1b1df8b0e056d8")
	//doUpdateBlog(client, "6596ea901c1b1df8b0e056d8")
	//doGetAllBlogs(client)
	doDeleteBlog(client, "6596f90e93b61ae7f7f36b13")
}
