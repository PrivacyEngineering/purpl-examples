package main

import (
	"context"
	"log"
	"net"

	"example.com/m/v2/pb"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/exp/rand"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type server struct {
	pb.UnimplementedPingPongServer
}

// CustomClaims is our custom metadata
type CustomClaims struct {
	Allowed   []string `json:"allowed"`   // Allowed fields
	Minimized []string `json:"minimized"` // To be minimized fields

	// everything alse not mentioned in either allowed or moinimzed will be reduced to -1 or nil

	jwt.StandardClaims // We'll use the standard claims
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Street: "Straße des 17 Juni ", Number: 135}, nil
}

func main() {
	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor),
	)
	pb.RegisterPingPongServer(s, &server{})
	if err := s.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func interceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	h, err := handler(ctx, req)
	if err != nil {
		return nil, err
	}

	// if r, ok := req.(*pb.HelloRequest); ok && strings.Contains(r.Name, "badclient") {
	// 	if hr, ok := h.(*pb.HelloReply); ok {
	// 		//hr.Number = reduceNumber(hr.Number)	// uncomment this line to reduce
	// 		hr.Number = noiseNumber(hr.Number) // uncomment this line to noise
	// 		// hr.Number = generalizeNumber(hr.Number) // uncomment this line to generalize
	// 	}
	// }

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if token := md.Get("authorization"); len(token) > 0 {
			tkn, err := jwt.ParseWithClaims(token[0], &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte("gpt-256-bit-secret"), nil
			})

			//log.Printf("Token: %s", token[0])

			if err != nil {
				return nil, err
			}

			if !tkn.Valid {
				return nil, jwt.NewValidationError("token is invalid", jwt.ValidationErrorMalformed)
			}

			claims, ok := tkn.Claims.(*CustomClaims)
			if !ok {
				return nil, jwt.NewValidationError("claims are not valid", jwt.ValidationErrorMalformed)
			}

			if len(claims.Allowed) > 0 && contains(claims.Allowed, "Number") {
				// pass
			} else if len(claims.Minimized) > 0 && contains(claims.Minimized, "Street") {
				// if hr, ok := h.(*pb.HelloReply); ok {
				// 	// DISCLAIMER: following functions are not implemented yet

				// 	// hr.Street = noiseStreet(hr.Street) // uncomment this line to noise
				// 	// hr.Street = generalizeStreet(hr.Street) // uncomment this line to generalize
				// }
			} else {
				// if hr, ok := h.(*pb.HelloReply); ok {
				// 	// DISCLAIMER: following functions are not implemented yet

				// 	// hr.Street = reduceStreet(hr.Street)	// uncomment this line to reduce
				// }
			}

			if len(claims.Allowed) > 0 && contains(claims.Allowed, "Number") {
				// pass
			} else if len(claims.Minimized) > 0 && contains(claims.Minimized, "Number") {
				if hr, ok := h.(*pb.HelloReply); ok {
					hr.Number = noiseNumber(hr.Number) // uncomment this line to noise
					// hr.Number = generalizeNumber(hr.Number) // uncomment this line to generalize
				}
			} else {
				if hr, ok := h.(*pb.HelloReply); ok {
					hr.Number = reduceNumber(hr.Number) // uncomment this line to reduce
				}
			}
		}
	}

	return h, nil
}

func reduceNumber(number int32) int32 {
	// receives a house number and returns -1 as "none".
	return -1
}

func noiseNumber(number int32) int32 {
	// receives a house number and returns noised version of it.
	return number - rand.Int31n(number) + rand.Int31n(number)
}

func generalizeNumber(number int32) int32 {
	// receives a house number and returns its range of 10's as the lower end of the interval.
	// e.g. 135 -> 131
	return number/10*10 + 1
}

// contains checks if a string is present in a slice
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
