package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func ValidateJWT(hash, password []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err != nil, err
}

func GenerateJWT(claims jwt.MapClaims) (string, error) {
	secret, err := GetEnv("JWT_SECRET")
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
