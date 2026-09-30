package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/carloscfgos1980/shop-gRPC/grpc-practice/server/proto/gen"
)

type server struct {
	pb.UnimplementedCalculatorServer
}

func (s *server) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	sum := req.A + req.B

	log.Println("Sum:", sum)
	return &pb.AddResponse{
		Sum: sum}, nil
}

func main() {
	port := ":50051"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	pb.RegisterCalculatorServer(grpcServer, &server{})

	log.Println("Server is listening on port", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
