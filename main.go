package main

import (
	"backdev_go/jwt"
	"encoding/base64"
	"io"

	// "crypto/sha512"
	// "github.com/BurntSushi/toml"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type AcessTokenPayload struct {
	UserUuid string     `json:"user_uuid"`
	UserName string     `json:"user_name"`
	RefreshToken string `json:"refresh_token"`
}

type AppConfig struct {
	Secret string
}

func main() {
	if (len(os.Args) != 2) {
		fmt.Println("Wrong amount of CLI arguments passed (must be 1)")
		os.Exit(1)
	}

	
	fsys := os.DirFS(".")
	filePath := os.Args[1]

	file, err := fsys.Open(filePath)
	if err != nil {
		fmt.Println("Error opening config file: ", err)
		os.Exit(2)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading config file: ", err)
		os.Exit(3)
	}
	

	var conf AppConfig
	_, err = toml.Decode(string(content), &conf)
	if err != nil {
		fmt.Println("Error parsing config file: ", err)
		os.Exit(4)
	}

	fmt.Println("Secret is: ", conf.Secret, "\n");

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

	jwt.JwtHello()
}
