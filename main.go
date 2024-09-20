package main

import (
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
	fmt.Println("Secret is: ", config.Secret)

	fmt.Println("Hello, world!")
}
