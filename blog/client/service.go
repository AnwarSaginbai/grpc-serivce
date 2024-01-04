package main

import (
	"context"
	"fmt"
	pb "github.com/AnwarSaginbai/grpc-service/blog/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"log"
)

func doCreateBlog(c pb.BlogServiceClient) string {
	log.Println("[CREATE] invoked in client-side")
	blog := &pb.Blog{
		AuthorId: "Clement",
		Title:    "My First Time",
		Content:  "Hello from GRPC service",
	}
	res, err := c.CreateBlog(context.Background(), blog)
	if err != nil {
		log.Fatal()
	}
	log.Printf("[INSERTED] blog id: %s", res.Id)
	return res.Id
}

func doGetByIdBlog(c pb.BlogServiceClient, id string) *pb.Blog {
	log.Println("[GET] invoked in client-side")
	IdBlog := &pb.BlogId{
		Id: id,
	}
	res, err := c.ReadBlog(context.Background(), IdBlog)
	if err != nil {
		fmt.Errorf("failed to find %v", err)
	}
	log.Printf("[READ] blog was find: %v", res)
	return res
}

func doUpdateBlog(c pb.BlogServiceClient, id string) {
	log.Println("[UPDATE] invoked in client-side")
	data := &pb.Blog{
		Id:       id,
		AuthorId: "Not Anwar",
		Title:    "changed",
		Content:  "changed",
	}
	_, err := c.UpdateBlog(context.Background(), data)
	if err != nil {
		fmt.Errorf("failed to update %v", err)
	}
	log.Println("[UPDATED] blog")
}

func doGetAllBlogs(c pb.BlogServiceClient) {
	log.Println("[GET ALL] invoked in client-side")
	stream, err := c.ListsBlog(context.Background(), &emptypb.Empty{})
	if err != nil {
		fmt.Errorf("failed to list all blogs %v", err)
	}
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal()
		}
		log.Println(msg)
	}
}

func doDeleteBlog(c pb.BlogServiceClient, id string) {
	log.Println("[DELETE] invoked in client-side")
	IdBlog := &pb.BlogId{
		Id: id,
	}
	_, err := c.DeleteBlog(context.Background(), IdBlog)
	if err != nil {
		fmt.Errorf("failed to find %v", err)
	}
	log.Println("[DELETED]")
}
