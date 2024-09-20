package main

import (
	"encoding/base64"
	"fmt"
)

type AcessTokenPayload struct {
	UserUuid string     `json:"user_uuid"`
	UserName string     `json:"user_name"`
	RefreshToken string `json:"refresh_token"`
}

func main() {
	HelloWorld()

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

}
