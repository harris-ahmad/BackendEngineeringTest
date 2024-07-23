package main

import (
	"context"
	"log"
	"time"

	authpb "github.com/harris-ahmad/BackendEngineeringTest/AuthService"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(":8080", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	client := authpb.NewAuthServiceClient(conn)

	// Sign up with phone number
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	signupResp, err := client.SignupWithPhoneNumber(ctx, &authpb.SignupWithPhoneNumberRequest{
		PhoneNumber: "+1234567890",
	})
	if err != nil {
		log.Fatalf("Could not sign up: %v", err)
	}
	log.Printf("Signup Response: %s", signupResp.Message)

	// Verify phone number (replace "123456" with the actual OTP received)
	verifyResp, err := client.VerifyPhoneNumber(ctx, &authpb.VerifyPhoneNumberRequest{
		PhoneNumber:      "+1234567890",
		VerificationCode: "123456",
	})
	if err != nil {
		log.Fatalf("Could not verify phone number: %v", err)
	}
	log.Printf("Verification Response: %s", verifyResp.Message)
}
