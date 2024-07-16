package main

import (
	pb "BackendEngineeringTest/AuthService/proto"
	"log"
	"net"

	"google.golang.org/grpc"
)

type AuthService struct {
	pb.UnimplementedAuthServiceServer
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	rpcServer := grpc.NewServer()
	// Register the service
	pb.RegisterAuthServiceServer(rpcServer, &AuthService{})
	if err := rpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
