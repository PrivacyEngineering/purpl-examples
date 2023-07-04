package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"

	"example.com/m/v2/pb"
	"github.com/golang-jwt/jwt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type server struct {
	pb.UnimplementedPingPongServer
}

// CustomClaims is our custom metadata
type CustomClaims struct {
	Policy struct {
		Allowed     map[string]string `json:"allowed"`
		Generalized map[string]string `json:"generalized"`
		Noised      map[string]string `json:"noised"`
		Reduced     map[string]string `json:"reduced"`
	} `json:"policy"`

	jwt.StandardClaims
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

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if token := md.Get("authorization"); len(token) > 0 {
			tkn, err := jwt.ParseWithClaims(token[0], &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte("MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0nZJVCy1saJr1ZbdCuQcztKlb/VdnyjqtNXEabRop4VbieAhtp9HpYEXVELy32LOx7H+yzc2q5ehVh0QyfOa33arwLiu9yuZZk00F83TP6RpXxysEH92lD4knAlUBm73t+IGrwFaRiA5Yt2ZtKOUvhzp0yRgRJfeXpwBxCRzxhuMP5PI2yds6nfdYml2ksXf0pdTkphwAIebQ5A7Sj53qIxR09fHZ01pUgd6S+Bl4pg6J/EEgJyZrA/4I12nuIWSCAlLlb/t7b6ExL/EdeXvIUB7ytjHfwpyZ5qRqHAa6lcZVYNk5xx5CHo+nCc4ayngDotl20V5wmQSPYbVxD/vyQIDAQAB"), nil
			})

			// -------------------------
			// ! Validation not working !
			// -------------------------

			if err != nil {
				// return nil, err
			}

			if !tkn.Valid {
				// return nil, jwt.NewValidationError("token is invalid", jwt.ValidationErrorMalformed)
			}

			claims, ok := tkn.Claims.(*CustomClaims)
			if !ok {
				// return nil, jwt.NewValidationError("claims are not valid", jwt.ValidationErrorMalformed)
			}

			if claims.StandardClaims.VerifyIssuer("test", true) {
				// return nil, jwt.NewValidationError("issuer is invalid", jwt.ValidationErrorMalformed)
			}

			if claims.StandardClaims.VerifyExpiresAt(time.Now().Unix(), true) {
				// return nil, jwt.NewValidationError("token is expired", jwt.ValidationErrorExpired)
			}

			// -------------------------
			// ! Validation not working !
			// -------------------------

			// Invoke ProtoReflect
			reflectedMsg := h.(*pb.HelloReply).ProtoReflect()

			// Declare a slice to store field names
			var fieldNames []string

			reflectedMsg.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
				name, err := getLastPart(string(fd.FullName()))

				if err != nil {
					fmt.Println("Error:", err)
				} else {
					fieldNames = append(fieldNames, name)
					// fmt.Printf("Field: %s\tValue: %v\n", fd.FullName(), v)
				}

				return true
			})

			// Iterate over the fields of the message
			for _, field := range fieldNames {
				// Check if the field is in the allowed list
				if !contains(claims.Policy.Allowed, field) {
					// Check if the field is in one of the minimized lists
					if contains(claims.Policy.Generalized, field) {
						// Generalize the field
						switch reflectedMsg.Descriptor().Fields().ByName(protoreflect.Name(field)).Kind() {
						case protoreflect.Int32Kind:
							reflectedMsg.Set(reflectedMsg.Descriptor().Fields().ByName(protoreflect.Name(field)), protoreflect.ValueOf(generalizeInt(reflectedMsg.Get(reflectedMsg.Descriptor().Fields().ByName(protoreflect.Name(field))).Int())))
						case protoreflect.StringKind:
							reflectedMsg.Set(reflectedMsg.Descriptor().Fields().ByName(protoreflect.Name(field)), protoreflect.ValueOf(generalizeString(reflectedMsg.Get(reflectedMsg.Descriptor().Fields().ByName(protoreflect.Name(field))).String())))
						}
					} else if contains(claims.Policy.Noised, field) {
						// Noise the field
						switch reflectedMsg.Descriptor().Fields().ByName(protoreflect.Name(field)).Kind() {
						case protoreflect.Int32Kind:
							reflectedMsg.Set(reflectedMsg.Descriptor().Fields().ByName(protoreflect.Name(field)), protoreflect.ValueOf(noiseInt(reflectedMsg.Get(reflectedMsg.Descriptor().Fields().ByName(protoreflect.Name(field))).Int())))
						case protoreflect.StringKind:
							reflectedMsg.Set(reflectedMsg.Descriptor().Fields().ByName(protoreflect.Name(field)), protoreflect.ValueOf(noiseString(reflectedMsg.Get(reflectedMsg.Descriptor().Fields().ByName(protoreflect.Name(field))).String())))
						}
					} else if contains(claims.Policy.Reduced, field) {
						log.Printf("\nField: %v", field)
						// Reduce the field
						switch reflectedMsg.Descriptor().Fields().ByName(protoreflect.Name(field)).Kind() {
						case protoreflect.Int32Kind:
							reflectedMsg.Set(reflectedMsg.Descriptor().Fields().ByName(protoreflect.Name(field)), protoreflect.ValueOf(reduceInt(reflectedMsg.Get(reflectedMsg.Descriptor().Fields().ByName(protoreflect.Name(field))).Int())))
						case protoreflect.StringKind:
							reflectedMsg.Set(reflectedMsg.Descriptor().Fields().ByName(protoreflect.Name(field)), protoreflect.ValueOf(reduceString(reflectedMsg.Get(reflectedMsg.Descriptor().Fields().ByName(protoreflect.Name(field))).String())))
						}
					} else {
						//Suppress the field
						switch reflectedMsg.Descriptor().Fields().ByName(protoreflect.Name(field)).Kind() {
						case protoreflect.Int32Kind:
							reflectedMsg.Set(reflectedMsg.Descriptor().Fields().ByName(protoreflect.Name(field)), protoreflect.ValueOf(suppressInt(reflectedMsg.Get(reflectedMsg.Descriptor().Fields().ByName(protoreflect.Name(field))).Int())))
						case protoreflect.StringKind:
							reflectedMsg.Set(reflectedMsg.Descriptor().Fields().ByName(protoreflect.Name(field)), protoreflect.ValueOf(suppressString(reflectedMsg.Get(reflectedMsg.Descriptor().Fields().ByName(protoreflect.Name(field))).String())))
						}
					}
				}
			}
		}
	}

	return h, nil
}

// ------ minimzation functions ------

// Suppression functions
func suppressInt(number int64) int32 {
	// receives an integer (e.g., house number) and returns -1 as "none".
	return -1
}
func suppressString(text string) string {
	// receives a string (e.g., street name) and cuts it off after the 5th character.
	return ""
}

// Noising functions
func noiseInt(number int64) int64 {
	// receives a house number and returns noised version of it.
	// rand.Int31 returns a non-negative pseudo-random 31-bit integer as an int32 from the default Source.
	return number - rand.Int63n(number) + rand.Int63n(number)
}
func noiseString(string) string {
	// receives a string and returns noised version of it.
	return ""
}

// Generalization functions
func generalizeInt(number int64) int64 {
	// receives an integer (e.g., house number) and returns its range of 10's as the lower end of the interval.
	// e.g. 135 -> 131
	return number/10*10 + 1
}
func generalizeString(text string) string {
	// receives a string (e.g., street name) and returns the first character.
	return text[0:1]
}

// Reduction functions
func reduceInt(number int64) int64 {
	return number / 10
}

func reduceString(text string) string {
	// receives a string (e.g., street name) and returns the first 4 characters.
	return text[0:3]
}

// ------ utiliy functions ------

// contains checks if a field is present in a map
func contains(m map[string]string, key string) bool {
	_, ok := m[key]
	return ok
}

// getLastPart returns the last part of a string separated by dots
// e.g., main.HelloReply.name --> name
func getLastPart(s string) (string, error) {
	parts := strings.Split(s, ".")
	if len(parts) < 1 {
		return "", errors.New("input string is empty")
	}
	return parts[len(parts)-1], nil
}
