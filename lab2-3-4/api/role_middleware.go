package api

import (
	"errors"
	"m1thrandir225/lab-2-3-4/auth"
	db "m1thrandir225/lab-2-3-4/db/sqlc"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (server *Server) requireRole(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		payload := ctx.MustGet(authorizationPayloadKey).(*auth.Claims)

		orgIDStr := ctx.Param("org_id")
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

		ctx.AbortWithStatusJSON(403, errorResponse(errors.New("insufficient permissions")))
	}
}

func (server *Server) requireModerator() gin.HandlerFunc {
	return server.requireRole("moderator")
}

func (server *Server) requireAdmin() gin.HandlerFunc {
	return server.requireRole("admin", "moderator")
}
