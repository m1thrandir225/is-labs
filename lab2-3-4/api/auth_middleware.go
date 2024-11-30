package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"m1thrandir225/lab-2-3-4/auth"
	"net/http"
	"strings"
)

// Everything else that requires a signed 2fa access and refresh token
func TwoFactorMiddleware(tokenMaker auth.TokenMaker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("Authorization header missing")))
			ctx.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("Authorization header invalid")))
			ctx.Abort()
			return
		}

		claims, err := tokenMaker.ValidateToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("Bearer token invalid")))
			ctx.Abort()
			return
		}

		if !claims.TwoFAVerified {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"error":       "2FA verification required",
				"email":       claims.Email,
				"redirect_to": "/verify-2fa",
			})
			ctx.Abort()
			return
		}

		ctx.Set("user_email", claims.Email)
		ctx.Next()
	}
}

// Just verify if the token is valid (only the verify-2fa route will have this middleware)
func LimitedAccessMiddleware(tokenMaker auth.TokenMaker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("Authorization header missing")))
			ctx.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("Authorization header invalid")))
			ctx.Abort()
			return
		}

		_, err := tokenMaker.ValidateToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("Bearer token invalid")))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
