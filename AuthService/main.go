package main

import (
	pb "BackendEngineeringTest/AuthService/proto"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
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

func init() {
	var err error

	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}

	var (
		dbUser string = os.Getenv("DB_USER")
		dbPass string = os.Getenv("DB_PASS")
		dbName string = os.Getenv("DB_NAME")
		dbHost string = os.Getenv("DB_HOST")
		dbPort string = os.Getenv("DB_PORT")
	)

	connString := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", dbUser, dbPass, dbName, dbHost, dbPort)
	db, err = sql.Open("postgres", connString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	rpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(rpcServer, &AuthService{})
	if err := rpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
