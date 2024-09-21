package model

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ---- ValidateToken
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

