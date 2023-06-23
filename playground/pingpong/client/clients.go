package main

import (
	"context"
	"log"
	"time"

	"example.com/m/v2/pb"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewPingPongClient(conn)

	// Good Client
	response, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: "goodclient"})
	if err != nil {
		log.Fatalf("Error on say hello: %v", err)
	}
	log.Printf("Message from server: ")
	log.Printf("Street: %s", response.GetStreet())
	log.Printf("Number: %d", response.GetNumber())
	time.Sleep(1 * time.Second)

	// Bad Client
	response, err = c.SayHello(context.Background(), &pb.HelloRequest{Name: "badclient"})
	if err != nil {
		log.Fatalf("Error on say hello: %v", err)
	}
	log.Printf("Message from server for badclient: ")
	log.Printf("Street: %s", response.GetStreet())
	log.Printf("Number: %d", response.GetNumber())
}
