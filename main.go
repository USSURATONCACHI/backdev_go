package main

import (
	"encoding/base64"
	"fmt"
	"os"

	"crypto/sha512"
)




func main() {
	// Read config
	writtenConfig, err := GetConfigFromCli()
	if (len(os.Args) != 2) {
		fmt.Println("Failed to parse config")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Create model
	model := Model {
		secret: sha512.Sum512([]byte(writtenConfig.Secret)),
	}
	fmt.Println("Your base64 of secret: ", base64.StdEncoding.EncodeToString(model.secret[:]))
	

	// Run server
	server := CreateServer(model)
	server.Run("localhost:8080")
}
