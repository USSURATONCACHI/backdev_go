package main

import (
	"backdev_go/jwt"
	"encoding/base64"

	// "crypto/sha512"
	// "github.com/BurntSushi/toml"
	"fmt"
	"os"

	
)

type AcessTokenPayload struct {
	UserUuid string     `json:"user_uuid"`
	UserName string     `json:"user_name"`
	RefreshToken string `json:"refresh_token"`
}


func main() {
	config, err := GetConfigFromCli()
	if (len(os.Args) != 2) {
		fmt.Println("Failed to parse config")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("Secret is: ", config.Secret, "\n");



	msg := "I am sigma!"
	
	fmt.Printf("Base string: %s\n", msg)
	encoded := base64.StdEncoding.EncodeToString([]byte(msg))

	fmt.Printf("Encoded: %s\n", encoded)

	decoded, err := base64.StdEncoding.DecodeString(encoded)

	if (err != nil) {
		fmt.Printf("Decode failed: %s\n", err)
		return
	}

	fmt.Printf("Decoded back: %s\n", decoded)

	jwt.JwtHello()
}
