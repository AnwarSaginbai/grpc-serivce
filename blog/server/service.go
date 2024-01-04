package main

import (
	"context"
	"fmt"
	pb "github.com/AnwarSaginbai/grpc-service/blog/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
)

func (s *Server) CreateBlog(ctx context.Context, in *pb.Blog) (*pb.BlogId, error) {
	log.Printf("[CREATE] function was invoked with %v\n", in)
	data := BlogItem{
		Author:  in.AuthorId,
		Title:   in.Title,
		Content: in.Content,
	}
	res, err := collection.InsertOne(ctx, data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("internal error: %v\n", err))
	}
	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("cannot convert to object id"))
	}

	return &pb.BlogId{Id: oid.Hex()}, nil

}

func (s *Server) ReadBlog(ctx context.Context, in *pb.BlogId) (*pb.Blog, error) {
	log.Printf("[GET] function was invoked with %v\n", in)
	oid, err := primitive.ObjectIDFromHex(in.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("failed to convert id: %v\n", in.Id))
	}
	filter := bson.M{"_id": oid}
	data := &BlogItem{}

	err = collection.FindOne(ctx, filter).Decode(data)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "not find")
	}
	return documentToBlog(data), nil
}

func (s *Server) UpdateBlog(ctx context.Context, in *pb.Blog) (*empty.Empty, error) {
	log.Printf("[PUT] function was invoked with %v\n", in)
	oid, err := primitive.ObjectIDFromHex(in.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("invalid argument: %v", in.Id))
	}
	data := &BlogItem{
		Author:  in.AuthorId,
		Title:   in.Title,
		Content: in.Content,
	}
	res, err := collection.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": data})
	if err != nil {
		return nil, status.Error(codes.Internal, "could not update")
	}
	if res.MatchedCount == 0 {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("not found"))
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteBlog(ctx context.Context, in *pb.BlogId) (*empty.Empty, error) {
	log.Printf("[DELETE] function was invoked with %v\n", in)
	id, err := primitive.ObjectIDFromHex(in.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "not found")
	}
	filter := bson.M{"_id": id}
	res, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error")
	}
	if res.DeletedCount == 0 {
		return nil, status.Errorf(codes.NotFound, "not found")
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) ListsBlog(in *empty.Empty, stream pb.BlogService_ListsBlogServer) error {
	log.Printf("[GET ALL] function was invoked with %v\n", in)

	cursor, err := collection.Find(context.Background(), primitive.D{{}})
	if err != nil {
		return status.Errorf(codes.Internal, "internal error")
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		data := &BlogItem{}
		err = cursor.Decode(data)
		if err != nil {
			return status.Errorf(codes.Internal, "internal error")
		}
		stream.Send(documentToBlog(data))
	}
	if err = cursor.Err(); err != nil {
		status.Errorf(codes.Internal, "internal error")
	}
	return nil
}
