package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/carloscfgos1980/gRPC-crash-course/grpc-stream/server/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"
)

func main() {
	// conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(creds), grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)))
	// conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)))
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)))
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	ctx := context.Background()

	client := pb.NewCalculatorClient(conn)
	req := &pb.FibonacciRequest{N: 10}
	stream, err := client.GenerateFibonacci(ctx, req)
	if err != nil {
		log.Fatalln(err)
	}

	for {
		res, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				log.Println("Stream ended")
				break
			}
			log.Fatalln(err)
		}
		log.Println("Fibonacci number:", res.GetNumber())
	}

	//------------------------------------------------------

	stream2, err := client.SendNumbers(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	// Send numbers to the server using the stream

	for i := 0; i < 9; i++ {
		err := stream2.Send(&pb.NumberRequest{Number: int32(i)})
		if err != nil {
			log.Fatalln(err)
		}
		time.Sleep(time.Second)
	}

	res, err := stream2.CloseAndRecv()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Server resp after stream:", res.Sum)

	// -----------------------------------------------

	stream3, err := client.Chat(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	go func() {
		for {
			meg, err := stream3.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Fatalln(err)
			}
			log.Println("Server: ", meg.GetMessage())
		}
	}()

	messages := []string{"Hi", "How are you?", "Bye"}

	for _, m := range messages {
		err := stream3.Send(&pb.ChatMessage{Message: m})
		if err != nil {
			log.Fatalln(err)
		}
		time.Sleep(time.Second)
	}

	stream3.CloseSend()

	// -------------------------------------------------------------
}
