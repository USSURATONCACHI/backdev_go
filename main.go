package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"time"

	"crypto/sha512"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gin-gonic/gin"
)



type AppConfig struct {
	SecretHashed [64]byte
}


func main() {
	// Read config
	writtenConfig, err := GetConfigFromCli()
	if (len(os.Args) != 2) {
		fmt.Println("Failed to parse config")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("Secret is: ", writtenConfig.Secret)

	config := AppConfig {
		SecretHashed: sha512.Sum512([]byte(writtenConfig.Secret)),
	}
	fmt.Println("Hashed secret base64 is: ", base64.StdEncoding.EncodeToString(config.SecretHashed[:]))

	// Run server
	router := gin.Default()

	router.GET("/auth", func(ctx *gin.Context) {
		Authorize(ctx, &config)
	})
	router.Run("localhost:8080")
}

type AuthorizeResponse struct {
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func Authorize(c *gin.Context, config *AppConfig) {
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
	ss, err := token.SignedString(config.SecretHashed[:])

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	response := AuthorizeResponse {
		AccessToken: ss,
		RefreshToken: "none",
	}

	c.IndentedJSON(http.StatusOK, response)
}