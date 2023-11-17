package auth

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

const secretKey = "supersecretkey"

type JWTClaims struct {
	jwt.RegisteredClaims
	UserID string
}

func NewToken(claims *JWTClaims) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	fmt.Println("token:", token)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	fmt.Println("tokenString:", tokenString)

	return tokenString, nil
}

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
