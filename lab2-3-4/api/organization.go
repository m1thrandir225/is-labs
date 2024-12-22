package api

import (
	"github.com/gin-gonic/gin"
	"m1thrandir225/lab-2-3-4/auth"
	db "m1thrandir225/lab-2-3-4/db/sqlc"
	"strconv"
)

type createOrganizationRequest struct {
	Name string `json:"name" binding:"required"`
}

type addUserToOrgRequest struct {
	UserID int64 `json:"user_id" binding:"required"`
	RoleID int64 `json:"role_id" binding:"required"`
}

func (server *Server) createOrganization(ctx *gin.Context) {
	var req createOrganizationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	org, err := server.store.CreateOrganization(ctx, req.Name)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	// Create default roles for the organization
	err = server.store.CreateInitialRoles(ctx, org.ID)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	// Add creator as moderator
	payload := ctx.MustGet(authorizationPayloadKey).(*auth.Claims)

	modRoleId, err := server.store.GetModeratorRole(ctx, org.ID)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
	}

	_, err = server.store.AddUserToOrganization(ctx, db.AddUserToOrganizationParams{
		UserID: user.ID,
		OrgID:  org.ID,
		RoleID: modRoleId,
	})

	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.JSON(201, org)
}

func (server *Server) addUserToOrganization(ctx *gin.Context) {
	var req addUserToOrgRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	orgIDStr := ctx.Param("org_id")
	orgID, err := strconv.ParseInt(orgIDStr, 10, 64)

	// Only admin can add users
	_, err = server.store.AddUserToOrganization(ctx, db.AddUserToOrganizationParams{
		UserID: req.UserID,
		OrgID:  orgID,
		RoleID: req.RoleID,
	})

	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.Status(201)
}

func (server *Server) removeUserFromOrganization(ctx *gin.Context) {
	userIDStr := ctx.Param("user_id")
	orgIDStr := ctx.Param("org_id")

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		ctx.JSON(400, errorResponse(err))
	}

	orgID, err := strconv.ParseInt(orgIDStr, 10, 64)
	if err != nil {
		ctx.JSON(400, errorResponse(err))
	}

	err = server.store.RemoveUserFromOrganization(ctx, db.RemoveUserFromOrganizationParams{
		UserID: userID,
		OrgID:  orgID,
	})

	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.Status(204)
}
