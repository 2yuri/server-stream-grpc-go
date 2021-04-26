package main

import (
	"fmt"

	"log"
	"net"

	"github.com/hyperyuri/server-stream-grpc-go/pb"

	"google.golang.org/grpc"
)

func main() {
	port := ":7000"
	fmt.Printf("server running at port %v\n", port)

	grpcServer := grpc.NewServer()

	pb.RegisterTestServiceServer(grpcServer, &TestServiceServer{})

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("cannot start server", err)
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
