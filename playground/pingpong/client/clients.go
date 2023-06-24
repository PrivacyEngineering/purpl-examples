package main

import (
	"context"
	"log"
	"time"

	"example.com/m/v2/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	goodToken = "goodtoken"
	badToken  = "badtoken"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewPingPongClient(conn)

	// Good Client
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", goodToken)
	response, err := c.SayHello(ctx, &pb.HelloRequest{Name: "goodclient"})
	if err != nil {
		log.Fatalf("Error on say hello: %v", err)
	}
	log.Printf("Message from server: ")
	log.Printf("Street: %s", response.GetStreet())
	log.Printf("Number: %d", response.GetNumber())
	time.Sleep(1 * time.Second)

	// Bad Client
	ctx = metadata.AppendToOutgoingContext(context.Background(), "authorization", badToken)
	response, err = c.SayHello(ctx, &pb.HelloRequest{Name: "badclient"})
	if err != nil {
		log.Fatalf("Error on say hello: %v", err)
	}
	log.Printf("Message from server for badclient: ")
	log.Printf("Street: %s", response.GetStreet())
	log.Printf("Number: %d", response.GetNumber())
}
