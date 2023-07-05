package main

import (
	"context"
	"log"
	"net"

	"example.com/m/v2/pb"

	"google.golang.org/grpc"

	// contribution:
	"github.com/louisloechel/purposelimiter"
)

type server struct {
	pb.UnimplementedPingPongServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Name: "Ken Guru", PhoneNumber: "+0123456789", Street: "Straße des 17 Juni", Age: 35, Sex: "male"}, nil
}

func main() {
	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(purposelimiter.UnaryServerInterceptor()),
	)

	pb.RegisterPingPongServer(s, &server{})
	if err := s.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
