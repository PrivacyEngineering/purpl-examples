package main

import (
	"context"
	"log"
	"net"

	"example.com/m/v2/pb"
	"golang.org/x/exp/rand"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type server struct {
	pb.UnimplementedPingPongServer
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
			log.Printf("Token: %s", token[0])
			if token[0] == "badtoken" {
				if hr, ok := h.(*pb.HelloReply); ok {
					//hr.Number = reduceNumber(hr.Number)	// uncomment this line to reduce
					hr.Number = noiseNumber(hr.Number) // uncomment this line to noise
					// hr.Number = generalizeNumber(hr.Number) // uncomment this line to generalize
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
