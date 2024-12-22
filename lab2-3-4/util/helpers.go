package util

import (
	"errors"
	"github.com/gin-gonic/gin"
	"m1thrandir225/lab-2-3-4/auth"
	"strings"
)

func GetPayloadFromContext(c *gin.Context) (*auth.Claims, error) {
	data, exists := c.Get("authorization_payload")
	if !exists {
		return nil, errors.New("authorization header not found")
	}
	payload := data.(*auth.Claims)
	return payload, nil
}

func IsDuplicateKeyError(err error) bool {
	// This will depend on your specific database and error handling
	// For SQLite, you might check for a unique constraint violation
	return strings.Contains(err.Error(), "UNIQUE constraint failed") ||
		strings.Contains(err.Error(), "unique constraint") ||
		strings.Contains(err.Error(), "duplicate key")
}
