package model

import (
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"backdev_go/db_io"
)

func (model *Model) createRawSignedJwtToken(tokenUuid uuid.UUID, userUuid uuid.UUID, userIp string) (string, error) {
	claims := Claims {
		UserUuid:  userUuid,
		UserName:  model.Syllables.HumanNameFromUuid(userUuid),
		UserIp:    userIp,
		TokenUuid: tokenUuid,

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID: tokenUuid.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	ss, err := token.SignedString(model.Secret[:])
	if err != nil {
		return "", errors.New("failed to sign a JWT token: " + err.Error())
	}

	return ss, nil
}

// Returns: (Success Result, Server error)
func (model *Model) CreateToken(userUuid uuid.UUID, userIp string) (*JwtAndRefreshTokens, error) {
	thisTokenUuid := uuid.New()
	refreshTokenUuid := uuid.New()

	refreshBcrypt, err := bcrypt.GenerateFromPassword(refreshTokenUuid[:], bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to create a bcrypt hash: " + err.Error())
	}

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
	jwtToken, err := model.createRawSignedJwtToken(thisTokenUuid, userUuid, userIp)
	if err != nil {
		return nil, err
	}

	// Return it
	result := JwtAndRefreshTokens {
		JwtToken: jwtToken,
		RefreshTokenBase64: base64.StdEncoding.EncodeToString(refreshTokenUuid[:]),
	}
	return &result, nil
}