package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"m1thrandir225/lab-2-3-4/auth"
	"strings"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker auth.TokenMaker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(authorizationHeaderKey)

		if len(authHeader) == 0 {
			err := errors.New("authorization Header not provided")
			ctx.AbortWithStatusJSON(401, errorResponse(err))
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			err := errors.New("invalid header format")
			ctx.AbortWithStatusJSON(401, errorResponse(err))
			return
		}

		authorizationType := fields[0]
		if strings.ToLower(authorizationType) != authorizationTypeBearer {
			err := errors.New("invalid authorization type")
			ctx.AbortWithStatusJSON(401, errorResponse(err))
			return
		}

		authToken := fields[1]
		payload, err := tokenMaker.ValidateToken(authToken)
		if err != nil {
			ctx.AbortWithStatusJSON(401, errorResponse(err))
			return
		}
		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}
