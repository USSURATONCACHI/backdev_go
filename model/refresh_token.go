package model

import (
	"backdev_go/db_io"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (model *Model) refresh_CheckTokenValid(tokenString string) (error, error) {
	is_valid, err := model.ValidateToken(tokenString)
	if err != nil {
		return err, nil
	}
	if !is_valid {
		return errors.New("token is invalid"), nil
	}
	return nil, nil
}
func (model *Model) refresh_CheckDatabaseForEntry(tokenUuid uuid.UUID, refreshToken uuid.UUID) (*db_io.RefreshToken, error, error) {
	refreshTokenEntry, serverError := model.Database.Get_RefreshToken(tokenUuid)
	if serverError != nil {
		return nil, nil, serverError
	}
	if refreshTokenEntry == nil {
		return nil, errors.New("invalid JWT + Refresh tokens pair"), nil
	}
	clientError := bcrypt.CompareHashAndPassword(refreshTokenEntry.RefreshBcryptHash, refreshToken[:])
	if clientError != nil {
		return nil, errors.New("invalid JWT + Refresh tokens pair"), nil
	}
	return refreshTokenEntry, nil, nil
}
func (model *Model) refresh_ProcessIpMismatch(oldIp string, newIp string) (error, error) {
	fmt.Println("User IP changed, sending EMail warning")
	return nil, nil
}

// Returns (Success result, Client error, Server error)
func (model *Model) RefreshToken(tokenString string, refreshTokenBase64 string, userIp string) (*JwtAndRefreshTokens, error, error) {
	// Check token validity
	clientError, serverError := model.refresh_CheckTokenValid(tokenString)
	if clientError != nil {
		return nil, clientError, serverError
	}

	// Convert refresh token Base64 -> UUID
	refreshTokenBytes, clientError := base64.StdEncoding.DecodeString(refreshTokenBase64)
	if clientError != nil || len(refreshTokenBytes) != 16 {
		return nil, errors.New("incorrect refresh token"), serverError
	}
	var refreshToken uuid.UUID
	copy(refreshToken[:], refreshTokenBytes)

	// Parse token claims
	var claims Claims
	_, serverError = jwt.ParseWithClaims(tokenString, &claims, model.getJwtKeyFunc())	
	if serverError != nil {
		return nil, nil, errors.New("failed to parse already verified token")
	}

	// Check DB for that token
	_, clientError, serverError = model.refresh_CheckDatabaseForEntry(claims.TokenUuid, refreshToken)
	if clientError != nil || serverError != nil {
		return nil, clientError, serverError
	}

	// Check IP mismatch
	if claims.UserIp != userIp {
		clientError, serverError = model.refresh_ProcessIpMismatch(claims.UserIp, userIp)
		if clientError != nil || serverError != nil {
			return nil, clientError, serverError
		}
	}

	// Delete used token
	serverError = model.Database.Remove_RefreshToken(claims.TokenUuid)
	if serverError != nil {
		return nil, nil, errors.New("failed to remove old refresh token: " + serverError.Error())
	}

	// Success
	success, serverError := model.CreateToken(claims.UserUuid, userIp)
	return success, nil, serverError
}