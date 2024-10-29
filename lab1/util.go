package main

import (
	"golang.org/x/crypto/scrypt"
	"math/rand"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func randomString(length int) string {
	var sb strings.Builder

	k := len(alphabet)

	for i := 0; i < length; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func generateKey() (string, error) {
	salt := make([]byte, 32)

	referenceString := randomString(32)

	dk, err := scrypt.Key([]byte(referenceString), salt, 1<<15, 8, 1, 32)

	if err != nil {
		return "", err
	}

	return string(dk), nil
}
