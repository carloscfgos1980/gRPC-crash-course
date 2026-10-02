package main

import (
	"io"
	"log"
	"net"
	"time"

	pb "github.com/carloscfgos1980/gRPC-crash-course/grpc-stream/server/proto/gen"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedCalculatorServer
}

func (s *server) GenerateFibonacci(req *pb.FibonacciRequest, stream pb.Calculator_GenerateFibonacciServer) error {
	n := req.N
	a, b := int32(0), int32(1)
	for i := int32(0); i < n; i++ {
		if err := stream.Send(&pb.FibonacciResponse{Number: a}); err != nil {
			return err
		}
		a, b = b, a+b
		// Simulate some processing delay
		time.Sleep(500 * time.Millisecond)
	}
	return nil
}

func (s *server) SendNumbers(stream pb.Calculator_SendNumbersServer) error {
	var sum int32
	req, err := stream.Recv()
	if err != nil {
		if err == io.EOF {
			return stream.SendAndClose(&pb.NumberResponse{Sum: sum})
		}
		return err
	}
	sum += req.GetNumber()
	return nil
}

func main() {
	port := ":50051"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterCalculatorServer(grpcServer, &server{})
	log.Printf("Server listening on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
