package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

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
	service := (os.Args[1])
	fmt.Println(service)

	// add expiration date if token exists already
	// to do: check if token is expired

	// generate token
	goodToken, err := jwt.GenerateToken("policy.json", service, "key.pem")

	if err != nil {
		log.Fatalf("Error on generating token: %v", err)
	}

	// Good Client
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", goodToken)

	start := time.Now()
	response, err := c.SayHello(ctx, &pb.HelloRequest{
		Phone:                     "91237123",
		StreetName:                "Marchstraße",
		StreetNumber:              32,
		ZipCode:                   13490,
		City:                      "Berlin",
		Country:                   "Germany",
		Email:                     "veryral@gmai.com",
		Name:                      "Mustermann",
		CreditCardNumber:          "123123478",
		CreditCardCvv:             234,
		CreditCardExpirationYear:  2023,
		CreditCardExpirationMonth: 12,
		Age:                       43})

	duration := time.Since(start).Microseconds()

	fmt.Println(duration)
	if err != nil {
		fmt.Println("Error on say hello:", err)
	}
	fmt.Println("-------------------------")
	fmt.Println("Message from server for goodclient:", response)

	// generate token
	//badToken, err := jwt.GenerateToken("client/policy.json", "service2", "client/key.pem")
	//
	//// Bad Client
	//ctx = metadata.AppendToOutgoingContext(context.Background(), "authorization", badToken)
	//
	//response, err = c.SayHello(ctx, &pb.HelloRequest{Name: "badclient"})
	//
	//if err != nil {
	//	log.Fatalf("Error on say hello: %v", err)
	//}
	//
	//log.Printf("Message from server for badclient:	%s", response)
	//log.Printf("-------------------------")
}
