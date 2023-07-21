package main

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"io/ioutil"
	"log"

	"github.com/Siar-Akbayin/jwt-go-auth"
)

func main() {
	tokenString, err := grpc_go_auth.GenerateToken("policy.json", "trackingService-maximal", "purpose1", "key.pem", 2)
	if err != nil {
		log.Println(err)
	}
	// Read the RSA public key from file
	keyData, err := ioutil.ReadFile("public.pem") // Replace with the path to your RSA public key file
	if err != nil {
		log.Fatalf("Error reading public key: %v", err)
	}

	// Parse the RSA public key
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(keyData)
	if err != nil {
		log.Fatalf("Error parsing public key: %v", err)
	}

	// Verify the token signature using the RSA public key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		log.Fatalf("Error verifying token: %v", err)
	}

	if token.Valid {
		fmt.Println("Token signature is valid")
	} else {
		fmt.Println("Token signature is invalid")
	}

	fmt.Println("Here is your token: ", tokenString)
}
