package utils

import (
	"github.com/golang-jwt/jwt/v5"
)

func ValidateJWT(hash string) (*jwt.Token, error) {
	token, err := jwt.Parse(hash, func(token *jwt.Token) (any, error) {
		secret, err := getSecret()
		return []byte(secret), err
	})
	return token, err
}

func getSecret() (string, error) {
	secret, err := GetEnv("JWT_SECRET")
	return secret, err
}

func GenerateJWT(claims jwt.MapClaims) (string, error) {
	secret, err := getSecret()
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return signedString, nil
}
