package auth

import "time"

type TokenMaker interface {
	GenerateToken(email string, twoFAVerified bool, timeDuration time.Duration) (string, error)
	ValidateToken(tokenStr string) (*Claims, error)
}
