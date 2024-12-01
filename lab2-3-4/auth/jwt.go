package auth

import (
	"github.com/golang-jwt/jwt"
	"time"
)

type JWTMaker struct {
	Key []byte
}

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func NewJWTMaker(key []byte) *JWTMaker {
	return &JWTMaker{Key: key}
}

func (maker *JWTMaker) GenerateToken(email string, timeDuration time.Duration) (string, error) {
	expirationTime := time.Now().Add(timeDuration)

	jwtClaims := jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
		IssuedAt:  time.Now().Unix(),
		Issuer:    "sebastijanzindl",
	}

	claims := &Claims{
		Email:          email,
		StandardClaims: jwtClaims,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(maker.Key)
}

func (maker *JWTMaker) ValidateToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return maker.Key, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrTokenInvalid
	}

	return claims, nil
}
