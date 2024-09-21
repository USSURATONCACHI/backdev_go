package model

import (
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

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




func (model *Model) CreateToken(userUuid uuid.UUID) (*jwt.Token, error) {
	thisTokenUuid := uuid.New()
	refreshTokenUuid := uuid.New()

	refreshBcrypt, err := bcrypt.GenerateFromPassword(refreshTokenUuid[:], bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to create a bcrypt hash: " + err.Error())
	}
	refreshBcryptBase64 := base64.StdEncoding.EncodeToString(refreshBcrypt)

	// Add refresh token
	err = model.Database.Add_RefreshToken(
		db_io.RefreshToken{
			JwtTokenUuid: thisTokenUuid,
			RefreshBcryptHash: refreshBcrypt,
		},
	)
	if err != nil {
		return nil, errors.New("failed to create a refresh token: " + err.Error())
	}
	
	// Generate JWT token
	claims := Claims {
		UserUuid: userUuid,
		UserName: model.GenerateNameFromUuid(userUuid),
		RefreshToken: refreshBcryptBase64,

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID: thisTokenUuid.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	return token, nil
}

func (model *Model) CreateTokenString(userUuid uuid.UUID) (string, error) {
	token, err := model.CreateToken(userUuid)
	if err != nil {
		return "", err
	}

	ss, err := token.SignedString(model.Secret[:])
	if err != nil {
		return "", err
	}

	return ss, nil
}



func (model *Model) ValidateToken(tokenSigned string) (bool, error) {
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS512.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Method.Alg())
		}

		return model.Secret[:], nil
	}

	token, err := jwt.Parse(tokenSigned, keyFunc)
	if err != nil {
		return false, err
	}

	claims, ok := token.Claims.(Claims)
	if !ok {
		return false, fmt.Errorf("failed to parse claims of token")
	}

	if time.Now().Before(claims.NotBefore.Time) {
		return false, fmt.Errorf("token not yet valid")
	}
	if time.Now().After(claims.ExpiresAt.Time) {
		return false, fmt.Errorf("token expired")
	}

	return true, nil
}

