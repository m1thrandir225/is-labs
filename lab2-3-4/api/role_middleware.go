package api

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"m1thrandir225/lab-2-3-4/auth"
	db "m1thrandir225/lab-2-3-4/db/sqlc"
	"m1thrandir225/lab-2-3-4/dto"
	"net/http"
)

func roleMiddleware(role dto.Role, store db.Store) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		payload, ok := ctx.Get(authorizationPayloadKey)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New("unauthorized")))
		}

		userPayload, ok := payload.(*auth.Claims)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New("unauthorized")))
		}

		user, err := store.GetUserByEmail(ctx, userPayload.Email)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				ctx.AbortWithStatusJSON(http.StatusNotFound, errorResponse(errors.New("the user was not found")))
			}
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(errors.New("there was an internal server error")))
		}

		if user.Role != role {
			ctx.AbortWithStatusJSON(http.StatusForbidden, errorResponse(errors.New("you don't have permission to access this resource")))
		}
		ctx.Next()
	}
}
