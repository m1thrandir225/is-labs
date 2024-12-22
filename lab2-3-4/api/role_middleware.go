package api

import (
	"errors"
	"log"
	db "m1thrandir225/lab-2-3-4/db/sqlc"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (server *Server) requireRole(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		payload, err := GetPayloadFromContext(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		}

		orgIDStr := ctx.Param("id")
		log.Println(orgIDStr)

		orgID, err := strconv.ParseInt(orgIDStr, 10, 64)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, errorResponse(errors.New("unauthorized")))
			return
		}

		user, err := server.store.GetUserByEmail(ctx, payload.Email)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, errorResponse(errors.New("unauthorized")))
			return
		}

		userRole, err := server.store.GetUserRole(ctx, db.GetUserRoleParams{
			UserID: user.ID,
			OrgID:  orgID,
		})
		if err != nil {
			ctx.AbortWithStatusJSON(403, errorResponse(errors.New("unauthorized")))
			return
		}

		for _, role := range roles {
			if userRole == role {
				ctx.Next()
				return
			}
		}

		ctx.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		ctx.Abort()
	}
}

func (server *Server) requireModerator() gin.HandlerFunc {
	return server.requireRole("moderator")
}

func (server *Server) requireAdmin() gin.HandlerFunc {
	return server.requireRole("admin", "moderator")
}
