package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"m1thrandir225/lab-2-3-4/auth"
)

func GetPayloadFromContext(c *gin.Context) (*auth.Claims, error) {
	data, exists := c.Get(authorizationPayloadKey)
	if !exists {
		return nil, errors.New("authorization header not found")
	}
	payload := data.(*auth.Claims)
	return payload, nil
}
