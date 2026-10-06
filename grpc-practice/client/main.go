package main

import (
	"context"
	"log"
	"time"

	pb "github.com/carloscfgos1980/gRPC-crash-course/grpc-practice/client/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/metadata"
)

func main() {
	cert := "cert.pem"

	creds, err := credentials.NewClientTLSFromFile(cert, "")
	if err != nil {
		log.Fatalf("failed to load TLS credentials: %v", err)
	}
	addr := "localhost:50051"

	// conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	// conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(creds), grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)))
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	addClient := pb.NewCalculatorClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	addReq := &pb.AddRequest{
		A: 10,
		B: 20,
	}

	// addRes, err := addClient.Add(ctx, addReq)
	md := metadata.Pairs("authorization", "Bearer kjsfckafadcada")
	ctx = metadata.NewOutgoingContext(ctx, md)

	var resHeader metadata.MD
	var resTrailer metadata.MD

	addRes, err := addClient.Add(ctx, addReq, grpc.UseCompressor(gzip.Name), grpc.Header(&resHeader), grpc.Trailer(&resTrailer))
	if err != nil {
		log.Fatalf("failed to call Add: %v", err)
	}

	log.Printf("Result: %d", addRes.Sum)
	log.Printf("Response Header: %s", resHeader["key1"])
	log.Printf("Response Trailer: %s", resTrailer["key2"])

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
