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
	goodToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhbGxvd2VkIjpbIlN0cmVldCIsIk51bWJlciJdfQ.hpcTmawTYz_FhbIOV3fDiQihD2CHqtRG0hYqmqxF3jE"                   // {"allowed": ["Street", "Number"]}
	badToken  = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhbGxvd2VkIjpbIlN0cmVldCJdLCJtaW5pbWl6ZWQiOlsiTnVtYmVyIl19.yPv5EEPaAuKb-QEBZ0zb42Esi3h9Qy6O6s7Dq3sx0HQ" // {"allowed": ["Street"], "minimized": ["Number"]}
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
