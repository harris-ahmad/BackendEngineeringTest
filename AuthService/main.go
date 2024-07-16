package main

import (
	authpb "BackendEngineeringTest/AuthService/authproto"
	"math/rand"

	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"

	_ "github.com/lib/pq" // Import PostgreSQL driver
)

type AuthService struct {
	authpb.UnimplementedAuthServiceServer
}

var (
	db      *sql.DB
	channel *amqp.Channel
)

func publishOtpMessage(phoneNumber, otp string) {
	var err error
	mssgBody := fmt.Sprintf("OTP for phone number %s is %s", phoneNumber, otp)
	err = channel.Publish("", "otp", false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(mssgBody),
	})
	if err != nil {
		log.Fatalf("Failed to publish message: %v", err)
	}
}

func generateOTP() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

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

	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName)
	fmt.Println(connString)
	db, err = sql.Open("postgres", connString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
}

func (s *AuthService) SignupWithPhoneNumber(ctx context.Context, in *authpb.SignupWithPhoneNumberRequest) (*authpb.SignupWithPhoneNumberResponse, error) {
	// generate otp
	otp := generateOTP()
	var err error
	// save otpResponse to database
	_, err = db.Exec("INSERT INTO otps (phone_number, otp) VALUES ($1, $2)", in.PhoneNumber, otp)
	if err != nil {
		return nil, fmt.Errorf("FAILED TO SAVE OTP TO DATABASE: %v", err)
	}

	// send otpResponse to user
	publishOtpMessage(in.PhoneNumber, otp)
	return &authpb.SignupWithPhoneNumberResponse{Message: "Signup Successful, OTP sent"}, nil
}

func (s *AuthService) VerifyPhoneNumber(ctx context.Context, in *authpb.VerifyPhoneNumberRequest) (*authpb.VerifyPhoneNumberResponse, error) {
	var otp string
	err := db.QueryRow("SELECT otp FROM users WHERE phone_number = $1", in.PhoneNumber).Scan(&otp)
	if err != nil {
		return nil, fmt.Errorf("FAILED TO GET OTP FROM DATABASE: %v", err)
	}
	if in.VerificationCode == otp {
		_, err = db.Exec("UPDATE users SET is_verified = true WHERE phone_number = $1", in.PhoneNumber)
		if err != nil {
			return nil, fmt.Errorf("FAILED TO UPDATE USER TO VERIFIED: %v", err)
		}
		return &authpb.VerifyPhoneNumberResponse{Message: "Phone number verified"}, nil
	}
	return &authpb.VerifyPhoneNumberResponse{Message: "Invalid OTP"}, nil
}

func (s *AuthService) LoginWithPhoneNumber(ctx context.Context, in *authpb.LoginWithPhoneNumberRequest) (*authpb.LoginWithPhoneNumberResponse, error) {
	var verified bool
	err := db.QueryRow("SELECT is_verified FROM users WHERE phone_number = $1", in.PhoneNumber).Scan(&verified)
	if err != nil || !verified {
		return nil, fmt.Errorf("FAILED TO GET USER FROM DATABASE: %v", err)
	}
	return &authpb.LoginWithPhoneNumberResponse{Message: "Login Successful"}, nil
}

func (s *AuthService) GetProfile(ctx context.Context, in *authpb.GetProfileRequest) (*authpb.GetProfileResponse, error) {
	var profile authpb.Profile
	err := db.QueryRow("SELECT name, phone_number FROM users WHERE phone_number = $1", in.PhoneNumber).Scan(&profile.Name, &profile.PhoneNumber)
	if err != nil {
		return nil, fmt.Errorf("FAILED TO GET PROFILE FROM DATABASE: %v", err)
	}
	return &authpb.GetProfileResponse{Profile: &profile}, nil
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
}
