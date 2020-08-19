package jwtmanager

import (
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

// Decrypt extracts and validates the claims from the JWT string
func (jm *jwtManager) Decrypt(tokenString string, claims jwt.Claims) error {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jm.SigningMethod {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(jm.Secret), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("Invalid token")
	}

	return nil
}
