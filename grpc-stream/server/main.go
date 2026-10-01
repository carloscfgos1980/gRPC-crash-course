package main

import (
	pb "github.com/carloscfgos1980/gRPC-crash-course/grpc-stream/server/proto/gen"
)

type server struct {
	pb.UnimplementedCalculatorServer
}

func main() {}
