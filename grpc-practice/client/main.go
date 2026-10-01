package main

import (
	"context"
	"log"
	"time"

	pb "github.com/carloscfgos1980/gRPC-crash-course/grpc-practice/client/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	cert := "cert.pem"

	creds, err := credentials.NewClientTLSFromFile(cert, "")
	if err != nil {
		log.Fatalf("failed to load TLS credentials: %v", err)
	}
	addr := "localhost:50051"
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	addClient := pb.NewCalculatorClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	addReq := pb.AddRequest{
		A: 10,
		B: 20,
	}

	addRes, err := addClient.Add(ctx, &addReq)
	if err != nil {
		log.Fatalf("failed to call Add: %v", err)
	}
	log.Printf("Result: %d", addRes.Sum)

	// ------------------------------------------------------

	greeterClient := pb.NewGreeterClient(conn)

	greeterReq := pb.HelloRequest{
		Name: "Carlos",
	}

	greeterRes, err := greeterClient.Greet(ctx, &greeterReq)
	if err != nil {
		log.Fatalf("failed to call Greet: %v", err)
	}
	log.Printf("Message: %s", greeterRes.Message)
}
