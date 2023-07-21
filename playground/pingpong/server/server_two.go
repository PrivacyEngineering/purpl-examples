package main

import (
	"context"
	purposelimiter "github.com/louisloechel/purpl"
	"log"
	"net"

	"example.com/m/v2/pb"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedPingPongServer
}

// Send a message containing name, phone number, street, age and sex
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{
		Phone:                     in.Phone,
		StreetName:                in.StreetName,
		StreetNumber:              in.StreetNumber,
		ZipCode:                   in.ZipCode,
		City:                      in.City,
		Country:                   in.Country,
		Email:                     in.Email,
		Name:                      in.Name,
		CreditCardNumber:          in.CreditCardNumber,
		CreditCardCvv:             in.CreditCardCvv,
		CreditCardExpirationYear:  in.CreditCardExpirationYear,
		CreditCardExpirationMonth: in.CreditCardExpirationMonth,
		Age:                       in.Age}, nil
}

func main() {
	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// path to public key
	keyPath := "server/public.pem"

	s := grpc.NewServer(
		grpc.UnaryInterceptor(purposelimiter.UnaryServerInterceptor(keyPath)), // <--- don't forget to pass public key to interceptor
	)

	pb.RegisterPingPongServer(s, &server{})
	if err := s.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
