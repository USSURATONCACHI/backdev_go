package model

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"backdev_go/db_io"
)

type Claims struct {
	UserUuid uuid.UUID  `json:"user_uuid"`
	UserName string     `json:"user_name"`
	RefreshToken string `json:"refresh_token"`

	jwt.RegisteredClaims
}

type Model struct {
	Secret [64]byte
	Database db_io.Database

	StartSyllables []string
	MiddleSyllables []string
	FinalSyllables []string
}