package main

import (
	"context"
	"log"

	"example.com/m/v2/pb"
	jwt "github.com/Siar-Akbayin/jwt-go-auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewPingPongClient(conn)

	// add expiration date if token exists already
	// to do: check if token is expired

	// generate token
	goodToken, err := jwt.GenerateToken("client/policy.json", "service1", "client/key.pem")

	if err != nil {
		log.Fatalf("Error on generating token: %v", err)
	}

	// Good Client
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", goodToken)

	response, err := c.SayHello(ctx, &pb.HelloRequest{Name: "goodclient"})

	if err != nil {
		log.Fatalf("Error on say hello: %v", err)
	}
	log.Printf("-------------------------")
	log.Printf("Message from server for goodclient:	%s", response)

	// generate token
	badToken, err := jwt.GenerateToken("client/policy.json", "service2", "client/key.pem")

	// Bad Client
	ctx = metadata.AppendToOutgoingContext(context.Background(), "authorization", badToken)

	response, err = c.SayHello(ctx, &pb.HelloRequest{Name: "badclient"})

	if err != nil {
		log.Fatalf("Error on say hello: %v", err)
	}

	log.Printf("Message from server for badclient:	%s", response)
	log.Printf("-------------------------")
}
