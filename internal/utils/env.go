package utils

import (
	"errors"
	"os"
)

func GetEnv(key string) (string, error) {
	if value, ok := os.LookupEnv(key); ok {
		return value, nil
	}
	return "", errors.New("Cannot find environment variable " + key)
}
