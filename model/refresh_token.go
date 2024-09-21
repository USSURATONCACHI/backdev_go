package model

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// ---- CreateToken
func (model *Model) RefreshToken(tokenString string, refreshToken uuid.UUID) (*JwtAndRefreshTokens, error) {
	is_valid, err := model.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if !is_valid {
		return nil, errors.New("cannot refresh invalid token")
	}

	
	keyFunc := model.getJwtKeyFunc()
	var claims Claims
	_, err = jwt.ParseWithClaims(tokenString, &claims, keyFunc)	
	if err != nil {
		return nil, errors.New("somehow failed to parse verified token")
	}

	tokenUuid, err := uuid.Parse(claims.ID)
	if err != nil {
		return nil, errors.New("token does not contain a valid uuid")
	}

	refreshTokenEntry, err := model.Database.Get_RefreshToken(tokenUuid)
	if err != nil || refreshTokenEntry == nil {
		return nil, errors.New("invalid JWT + Refresh tokens pair")
	}
	err = bcrypt.CompareHashAndPassword(refreshTokenEntry.RefreshBcryptHash, refreshToken[:])
	if err != nil {
		return nil, errors.New("invalid JWT + Refresh tokens pair")
	}

	err = model.Database.Remove_RefreshToken(tokenUuid)
	if err != nil {
		return nil, errors.New("failed to remove old refresh token: " + err.Error())
	}

	return model.CreateToken(claims.UserUuid)
}