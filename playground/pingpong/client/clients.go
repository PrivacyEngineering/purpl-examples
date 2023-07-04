package main

import (
	"context"
	"log"

	"example.com/m/v2/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	// -----------------------------
	// {
	// 	"policy": {
	// 	  "allowed": {
	// 		"name": "string",
	// 		"sex": "string"
	// 	  },
	// 	  "generalized": {
	// 		"phoneNumber": "string"
	// 	  },
	// 	  "noised": {
	// 		"age": "int"
	// 	  },
	// 	  "reduced": {
	// 		"street": "string"
	// 	  }
	// 	},
	// 	"exp": 1688843806,
	// 	"iss": "test"
	//   }
	goodToken = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJwb2xpY3kiOnsiYWxsb3dlZCI6eyJuYW1lIjoic3RyaW5nIiwic2V4Ijoic3RyaW5nIn0sImdlbmVyYWxpemVkIjp7InBob25lTnVtYmVyIjoic3RyaW5nIn0sIm5vaXNlZCI6eyJhZ2UiOiJpbnQifSwicmVkdWNlZCI6eyJzdHJlZXQiOiJzdHJpbmcifX0sImV4cCI6MTY4ODg0MzgwNiwiaXNzIjoidGVzdCJ9.zrAKl4lH5_hPjB77HKDudNqn3Zq3eqmUA0ELteLvjVOeAKUI6gtpMRy-uvmVQtiI1XUSeZNwUDEdfSS2--thIlPiYIL2UQxIa5k1GPi3zWOk05akLj6aLngJxEB7fq18eQ_6ZzQMC6UA_kOrZMLjLzMSxZrU50_AALl0KP-M3New9Bck0SC56-IU4PTWQh2oQdD9Bh3LXE3T-bzZPRkUJ6op_fDHZ1-kLJN8JqzzcQNg7Fw__fxHHInSbQ_QJ3-rmNoNjyDWOdDB-CWxNpYnbS6QgVdJjwoG0AzMRtMtEiZKeKwbhATIic0UC42oaBU9tN6f_g4x-6IDzboos4IcAw"

	// -----------------------------
	// payload for badTokenTwo:
	// {
	// 	"policy": {
	// 	"allowed": {},
	// 	"generalized": {
	// 		"age": "int",
	// 		"name": "string",
	// 		"phoneNumber": "string",
	// 		"sex": "string",
	// 		"street": "string"
	// 	 },
	// 	 "noised": {},
	//   "reduced": {}
	// 	},
	// 	"exp": 1688483421,
	// 	"iss": "test"
	//   }
	badToken = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJwb2xpY3kiOnsiYWxsb3dlZCI6e30sImdlbmVyYWxpemVkIjp7ImFnZSI6ImludCIsIm5hbWUiOiJzdHJpbmciLCJwaG9uZU51bWJlciI6InN0cmluZyIsInNleCI6InN0cmluZyIsInN0cmVldCI6InN0cmluZyJ9LCJub2lzZWQiOnt9LCJyZWR1Y2VkIjp7fX0sImV4cCI6MTY4ODQ4MzQyMSwiaXNzIjoidGVzdCJ9.v04mp6rZuPNjrv015KVRU9ijFiUrH109AzPRI8VWOxSq5Ene9HTOPGLYdaY8D66LSC1La2EEpnxJnwKr0i8t2e3JgcgJta0ex2hMqFDoc3YsaIMkSrnMaHB6xkH2eh35agrC-2LRMJ7nMhAnhnkJci9R7oCySI9wg5eN2NJVonOmGnWlQ88QL-hDUFattIF1dxfHZK7CiwmwB1tkPnxClebrHr6ADnP5Vv3uZKBhtGy8Q3kxPNta01XcztTFGSFLDSnNswidV4uWi7opbm_d7JohLn6jQ-2vcjWrQhhRZzJjZ_uaynRVYEYtJ_tlk3nhQMYN_EWCarp9romPYxYXyw"
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
	log.Printf("-------------------------")
	log.Printf("Message from server for goodclient:	%s", response)

	// Bad Client
	ctx = metadata.AppendToOutgoingContext(context.Background(), "authorization", badToken)

	response, err = c.SayHello(ctx, &pb.HelloRequest{Name: "badclient"})

	if err != nil {
		log.Fatalf("Error on say hello: %v", err)
	}

	log.Printf("Message from server for badclient:	%s", response)
	log.Printf("-------------------------")
}
