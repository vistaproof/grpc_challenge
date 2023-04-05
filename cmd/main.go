package main

import (
	"fmt"
	"log"
	"net"

	"github.com/antstalepresh/grpc-challenge/types"

	"github.com/antstalepresh/grpc-challenge/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", types.ServerPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Set up the gRPC server
	srv := grpc.NewServer()
	types.RegisterGenericServiceServer(srv, &server.Server{})
	reflection.Register(srv)

	// Start the server
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
