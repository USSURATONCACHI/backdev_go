package model

import (
	"encoding/base64"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"backdev_go/db_io"
)

type RawTokeninfo struct {
	UserUuid uuid.UUID
	UserIp string
	UserEmail string
}

func (model *Model) createRawSignedJwtToken(tokenUuid uuid.UUID, info RawTokeninfo) (string, error) {
	claims := Claims {
		UserUuid:  info.UserUuid,
		UserName:  model.Syllables.HumanNameFromUuid(info.UserUuid),
		UserIp:    info.UserIp,
		TokenUuid: tokenUuid,
		UserEmail: info.UserEmail,

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

// Returns: (Success Result, Client error, Server error)
func (model *Model) CreateToken(info RawTokeninfo) (*JwtAndRefreshTokens, error, error) {
	// Check validity of data
	err := model.create_CheckTokenValidity(info)
	if err != nil {
		return nil, err, nil
	}

	// Generate UUID pair
	thisTokenUuid := uuid.New()
	refreshTokenUuid := uuid.New()

	refreshBcrypt, err := bcrypt.GenerateFromPassword(refreshTokenUuid[:], bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, errors.New("failed to create a bcrypt hash: " + err.Error())
	}

	// Add refresh token
	err = model.Database.Add_RefreshToken(
		db_io.RefreshToken{
			JwtTokenUuid: thisTokenUuid,
			RefreshBcryptHash: refreshBcrypt,
		},
	)
	if err != nil {
		return nil, nil, errors.New("failed to create a refresh token: " + err.Error())
	}
	
	// Generate JWT token
	jwtToken, err := model.createRawSignedJwtToken(thisTokenUuid, info)
	if err != nil {
		return nil, nil, err
	}

	// Return it
	result := JwtAndRefreshTokens {
		JwtToken: jwtToken,
		RefreshTokenBase64: base64.StdEncoding.EncodeToString(refreshTokenUuid[:]),
	}
	return &result, nil, nil
}

func (model *Model) create_CheckTokenValidity(info RawTokeninfo) error {
	if len(info.UserEmail) == 0 {
		return errors.New("email must be specified")
	}
	if len(info.UserIp) == 0 {
		return errors.New("IP must be specified")
	}

	if strings.Contains(info.UserEmail, "\n") {
		return errors.New("email cannot have newlines")
	}
	if strings.Contains(info.UserEmail, " ") || strings.Contains(info.UserEmail, "\t") {
		return errors.New("do not use whitespaces in email addresses")
	}

	if strings.Contains(info.UserIp, "\n") {
		return errors.New("IP cannot have newlines")
	}

	return nil
}