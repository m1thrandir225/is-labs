package auth

import "time"

type TokenMaker interface {
	GenerateToken(email string, timeDuration time.Duration) (string, error)
	ValidateToken(tokenStr string) (*Claims, error)
}
