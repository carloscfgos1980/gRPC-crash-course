package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/carloscfgos1980/gRPC-crash-course/grpc-practice/server/proto/gen"

	_ "google.golang.org/grpc/encoding/gzip"
)

type server struct {
	pb.UnimplementedCalculatorServer
	pb.UnimplementedGreeterServer
}

func (s *server) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		log.Println("Metadata:", md)
	}
	log.Printf("metadata: %v", md)
	authHeader, ok := md["authorization"]
	if ok {
		log.Println("Authorization Header:", authHeader)
	}
	log.Println("Authorization Header:", authHeader)

	// Set response with headers
	responseHeader := metadata.Pairs("key1", "value1")
	err := grpc.SendHeader(ctx, responseHeader)
	if err != nil {
		log.Printf("failed to send header: %v", err)
	}

	sum := req.A + req.B

	log.Println("Sum:", sum)
	return &pb.AddResponse{
		Sum: sum}, nil
}

func (s *server) Greet(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	message := fmt.Sprintf("Hello, %s, nice to meet you", req.Name)
	log.Println("Greet:", message)
	return &pb.HelloResponse{
		Message: message}, nil
}

func main() {
	cert := "cert.pem"
	key := "key.pem"
	creds, err := credentials.NewServerTLSFromFile(cert, key)
	if err != nil {
		log.Fatalf("failed to load TLS credentials: %v", err)
	}
	port := ":50051"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer(grpc.Creds(creds))

	pb.RegisterCalculatorServer(grpcServer, &server{})
	pb.RegisterGreeterServer(grpcServer, &server{})

	log.Println("Server is listening on port", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
