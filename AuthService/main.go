package main

import (
	otppb "BackendEngineeringTest/AuthService/otpproto"
	authpb "BackendEngineeringTest/AuthService/proto"
	"crypto/x509"

	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type AuthService struct {
	authpb.UnimplementedAuthServiceServer
}

var (
	db        *sql.DB
	channel   *amqp.Channel
	otpClient otppb.OtpServiceClient
)

var (
	dbUser string = os.Getenv("DB_USER")
	dbPass string = os.Getenv("DB_PASS")
	dbName string = os.Getenv("DB_NAME")
	dbHost string = os.Getenv("DB_HOST")
	dbPort string = os.Getenv("DB_PORT")
)

func publishOtpMessage(phoneNumber, otp string) {
	var err error
	var mssgBody string
	mssgBody = fmt.Sprintf("OTP for phone number %s is %s", phoneNumber, otp)
	err = channel.Publish("", "otp", false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(mssgBody),
	})
	if err != nil {
		log.Fatalf("Failed to publish message: %v", err)
	}
}

func init() {
	var err error

	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}

	connString := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", dbUser, dbPass, dbName, dbHost, dbPort)
	db, err = sql.Open("postgres", connString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
}

func (s *AuthService) SignupWithPhoneNumber(ctx context.Context, in *authpb.SignupWithPhoneNumberRequest) (*authpb.SignupWithPhoneNumberResponse, error) {
	otpResponse, err := otpClient.GenerateOtp(ctx, &otppb.GenerateOtpRequest{
		PhoneNumber: in.PhoneNumber,
	})
	if err != nil {
		return nil, fmt.Errorf("FAILED TO GENERATE OTP: %v", err)
	}
	// save otpResponse to database
	_, err = db.Exec("INSERT INTO otps (phone_number, otp) VALUES ($1, $2)", in.PhoneNumber, otpResponse.Otp)
	if err != nil {
		return nil, fmt.Errorf("FAILED TO SAVE OTP TO DATABASE: %v", err)
	}

	// send otpResponse to user
	publishOtpMessage(in.PhoneNumber, otpResponse.Otp)
	return &authpb.SignupWithPhoneNumberResponse{VerificationCode: "Signup Successful, OTP sent"}, nil
}

func (s *AuthService) VerifyPhoneNumber(ctx context.Context, in *authpb.VerifyPhoneNumberRequest) (*authpb.VerifyPhoneNumberResponse, error) {
	return &authpb.VerifyPhoneNumberResponse{}, nil
}

func (s *AuthService) LoginWithPhoneNumber(ctx context.Context, in *authpb.LoginWithPhoneNumberRequest) (*authpb.LoginWithPhoneNumberResponse, error) {
	return &authpb.LoginWithPhoneNumberResponse{}, nil
}

func (s *AuthService) GetProfile(ctx context.Context, in *authpb.GetProfileRequest) (*authpb.GetProfileResponse, error) {
	return &authpb.GetProfileResponse{}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	rpcServer := grpc.NewServer()
	authpb.RegisterAuthServiceServer(rpcServer, &AuthService{})
	if err := rpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	// connect to OTP service
	certPool := x509.NewCertPool()
	conn, err := grpc.NewClient("localhost:8081", grpc.WithTransportCredentials(
		credentials.NewClientTLSFromCert(certPool, ""),
	))
	if err != nil {
		log.Fatalf("Failed to connect to OTP service: %v", err)
	}
	otpClient = otppb.NewOtpServiceClient(conn)
}
