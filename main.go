package main

import (
	"encoding/base64"
	"fmt"
	"os"

	"crypto/sha512"

	"backdev_go/model"
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
	mdl := model.Model {
		Secret: sha512.Sum512([]byte(writtenConfig.Secret)),

		StartSyllables: writtenConfig.StartSyllables,
		MiddleSyllables: writtenConfig.MiddleSyllables,
		FinalSyllables: writtenConfig.FinalSyllables,
	}
	fmt.Println("Your base64 of secret: ", base64.StdEncoding.EncodeToString(mdl.Secret[:]))
	

	// Run server
	server := CreateServer(mdl)
	server.Run("localhost:8080")
}
