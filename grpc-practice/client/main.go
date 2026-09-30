package main

import (
	"context"
	"log"
	"time"

	pb "github.com/carloscfgos1980/shop-gRPC/grpc-practice/client/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	addr := "localhost:50051"
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewCalculatorClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := pb.AddRequest{
		A: 10,
		B: 20,
	}

	res, err := client.Add(ctx, &req)
	if err != nil {
		log.Fatalf("failed to call Add: %v", err)
	}
	log.Printf("Result: %d", res.Sum)
}
