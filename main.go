package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"crypto/sha512"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserUuid string     `json:"user_uuid"`
	UserName string     `json:"user_name"`
	RefreshToken string `json:"refresh_token"`

	jwt.RegisteredClaims
}


func main() {
	config, err := GetConfigFromCli()
	if (len(os.Args) != 2) {
		fmt.Println("Failed to parse config")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("Secret is: ", config.Secret)
	
	hashedSecret := sha512.Sum512([]byte(config.Secret))
	fmt.Println("Hashed secret base64 is: ", base64.StdEncoding.EncodeToString(hashedSecret[:]))

	claims := Claims {
		UserUuid: "aaaa",
		UserName: "John Pork",
		RefreshToken: "none",

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	ss, err := token.SignedString(hashedSecret[:])

	fmt.Println(ss, err)
}
