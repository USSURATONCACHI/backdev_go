package model

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"backdev_go/db_io"
)

type Claims struct {
	UserUuid  uuid.UUID `json:"user_uuid"`
	UserName  string    `json:"user_name"`
	UserIp    string    `json:"user_ip"`
	TokenUuid uuid.UUID `json:"token_uuid"`

	jwt.RegisteredClaims
}

type Model struct {
	Secret [64]byte
	Database db_io.Database
	Syllables Syllables
	SmtpInfo SmtpInfo
}

type JwtAndRefreshTokens struct {
	JwtToken string
	RefreshTokenBase64 string
}
