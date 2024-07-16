package main

import (
	pb "BackendEngineeringTest/AuthService/proto"
	"database/sql"
	"log"
	"net"

	"github.com/streadway/amqp"
	"google.golang.org/grpc"
)

type AuthService struct {
	pb.UnimplementedAuthServiceServer
}

var (
	db      *sql.DB
	channel *amqp.Channel
)

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
