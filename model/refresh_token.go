package model

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Returns (Success result, Client error, Server error)
func (model *Model) RefreshToken(tokenString string, refreshToken uuid.UUID, userIp string) (*JwtAndRefreshTokens, error, error) {
	// Check that token is valid
	is_valid, err := model.ValidateToken(tokenString)
	if err != nil {
		return nil, err, nil
	}
	if !is_valid {
		return nil, errors.New("token is invalid"), nil
	}

	// Parse token claims
	keyFunc := model.getJwtKeyFunc()
	var claims Claims
	_, err = jwt.ParseWithClaims(tokenString, &claims, keyFunc)	
	if err != nil {
		return nil, nil, errors.New("somehow failed to parse verified token")
	}

	tokenUuid, err := uuid.Parse(claims.ID)
	if err != nil {
		return nil, nil, errors.New("token does not contain a valid uuid")
	}

	if claims.UserIp != userIp {
		fmt.Println("User IP changed, sending EMail warning")
		panic("Not implemented yet")
	}

	// Check DB for that token
	refreshTokenEntry, err := model.Database.Get_RefreshToken(tokenUuid)
	if err != nil {
		return nil, nil, err
	}
	if refreshTokenEntry == nil {
		return nil, errors.New("invalid JWT + Refresh tokens pair"), nil
	}
	err = bcrypt.CompareHashAndPassword(refreshTokenEntry.RefreshBcryptHash, refreshToken[:])
	if err != nil {
		return nil, errors.New("invalid JWT + Refresh tokens pair"), nil
	}

	// Delete used token
	err = model.Database.Remove_RefreshToken(tokenUuid)
	if err != nil {
		return nil, nil, errors.New("failed to remove old refresh token: " + err.Error())
	}

	// Success
	success, serverError := model.CreateToken(claims.UserUuid, userIp)
	return success, nil, serverError
}