package util

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func HashPassword(base string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(base), bcrypt.DefaultCost)

	if err != nil {
		return "", fmt.Errorf("There was an error hashing your password: %s", err.Error())
	}

	return string(hashedPassword), nil
}
