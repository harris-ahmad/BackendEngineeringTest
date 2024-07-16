package main

import (
	pb "BackendEngineeringTest/AuthService/proto"
	"context"
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

func (s *AuthService) SignupWithPhoneNumber(ctx context.Context, in *pb.SignupWithPhoneNumberRequest) (*pb.SignupWithPhoneNumberResponse, error) {
	return &pb.SignupWithPhoneNumberResponse{}, nil
}

func (s *AuthService) VerifyPhoneNumber(ctx context.Context, in *pb.VerifyPhoneNumberRequest) (*pb.VerifyPhoneNumberResponse, error) {
	return &pb.VerifyPhoneNumberResponse{}, nil
}

func (s *AuthService) LoginWithPhoneNumber(ctx context.Context, in *pb.LoginWithPhoneNumberRequest) (*pb.LoginWithPhoneNumberResponse, error) {
	return &pb.LoginWithPhoneNumberResponse{}, nil
}

func (s *AuthService) GetProfile(ctx context.Context, in *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	return &pb.GetProfileResponse{}, nil
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
