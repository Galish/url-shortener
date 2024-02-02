// Package auth provides functions for user authentication and authorization.
package auth

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

// A secretKey is used to sign the token.
const secretKey = "supersecretkey"

// JWTClaims represents data encoded into a token.
type JWTClaims struct {
	jwt.RegisteredClaims
	UserID string
}

// NewToken generates a token.
func NewToken(claims *JWTClaims) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken decodes the token string.
func ParseToken(tokenString string) (*JWTClaims, error) {
	var claims JWTClaims

	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}

			return []byte(secretKey), nil
		},
	)
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("token is not valid")
	}

	return &claims, nil
}
