package main

import (
	"context"

	pb "github.com/harris-ahmad/BackendEngineeringTest/OtpService/proto"

	"github.com/streadway/amqp"
)

type OtpService struct {
	pb.UnimplementedOtpServiceServer
}

var channel *amqp.Channel

func (s *OtpService) GenerateOTP(ctx context.Context, in *pb.GenerateOtpRequest) (*pb.GenerateOtpResponse, error) {
	return &pb.GenerateOtpResponse{}, nil
}
