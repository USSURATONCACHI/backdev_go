package jwt

import (
	"encoding/base64"
	"crypto/sha512"
	"encoding/json"
	"errors"
)

type JwtHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

func CreateJwtSha512(payload string, secret string) (string, error) {
	header := GetSha512Header()
	headerJson, err := json.Marshal(header)
	if err != nil {
		return "", errors.New("Failed to convert header to JSON: " + err.Error())
	}

	stringToSign := string(headerJson) + "." + payload + "." + secret
	shasum := sha512.Sum512([]byte(stringToSign))

	headerBase64 := base64.StdEncoding.EncodeToString(headerJson)
	payloadBase64 := base64.StdEncoding.EncodeToString([]byte(payload))
	signatureBase64 := base64.StdEncoding.EncodeToString(shasum[:])

	jwtToken := headerBase64 + "." + payloadBase64 + "." + signatureBase64
	return jwtToken, nil
}

func GetSha512Header() JwtHeader {
	return JwtHeader {
		Alg: "SHA512",
		Typ: "JWT",
	}
}