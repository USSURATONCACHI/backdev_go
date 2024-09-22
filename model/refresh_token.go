package model

import (
	"backdev_go/db_io"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)



// Returns (Success result, Client error, Server error)
func (model *Model) RefreshToken(tokenString string, refreshTokenBase64 string, userIp string) (*JwtAndRefreshTokens, error, error) {
	// Process JWT token
	claims, clientError, serverError := model.refresh_CheckAndParseToken(tokenString)
	if clientError != nil || serverError != nil {
		return nil, clientError, serverError
	}

	// Process Refresh token
	clientError, serverError = model.refresh_ParseAndCheckRefreshToken(refreshTokenBase64, claims.TokenUuid)
	if clientError != nil || serverError != nil {
		return nil, clientError, serverError
	}

	// Check IP mismatch
	if claims.UserIp != userIp {
		clientError, serverError = model.refresh_ProcessIpMismatch(claims.UserIp, userIp, claims.UserName, claims.UserEmail)
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
	newTokenData := RawTokeninfo {
		UserUuid: claims.UserUuid,
		UserIp: userIp,
		UserEmail: claims.UserEmail,
	}
	success, clientError, serverError := model.CreateToken(newTokenData)
	return success, clientError, serverError
}


// ======
// =
// ==== Level 2 helpers
// =
// ======
func (model *Model) refresh_ParseAndCheckRefreshToken(refreshTokenBase64 string, jwtTokenUuid uuid.UUID) (error, error) {
	// Convert refresh token Base64 -> UUID
	refreshToken, clientError := ParseUuidFromBase64(refreshTokenBase64)
	if clientError != nil {
		return errors.New("incorrect refresh token"), nil
	}

	// Check DB for that token
	_, clientError, serverError := model.refresh_CheckDatabaseForEntry(jwtTokenUuid, refreshToken)
	if clientError != nil || serverError != nil {
		return clientError, serverError
	}

	return nil, nil
}

func (model *Model) refresh_CheckAndParseToken(tokenString string) (*Claims, error, error) {
	// Check token validity
	clientError, serverError := model.refresh_CheckTokenValid(tokenString)
	if clientError != nil {
		return nil, clientError, serverError
	}

	// Parse token claims
	var claims Claims
	_, serverError = jwt.ParseWithClaims(tokenString, &claims, model.getJwtKeyFunc())	
	if serverError != nil {
		return nil, nil, errors.New("failed to parse already verified token")
	}

	return &claims, nil, nil
}
// ======
// =
// ==== Level 2 helpers
// =
// ======


// ======
// =
// ==== Level 1 helpers
// =
// ======
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

func (model *Model) refresh_ProcessIpMismatch(oldIp string, newIp string, userName string, email string) (error, error) {
	fmt.Println("User IP changed, sending EMail warning")
	
	subject := "Backdev: You just changed your ip inside of session"
	body := "Dear " + userName + ", we detected an IP change in your session.\n" + 
	        "Previous IP: " + oldIp + " , current IP: " + newIp + " .\n"; 

	err := model.SmtpClient.SendEmail(subject, body, email);
	if err != nil {
		fmt.Println("SMTP sending error: ", err)
	}

	return nil, nil
}
// ======
// =
// ==== Level 1 helpers
// =
// ======