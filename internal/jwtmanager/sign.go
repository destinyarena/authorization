package jwtmanager

import (
	"github.com/dgrijalva/jwt-go"
)

// Sign encapsulates jwt claims into a JWT token string
func (jm *jwtManager) Sign(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jm.SigningMethod, claims)

	tokenString, err := token.SignedString([]byte(jm.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, err
}
