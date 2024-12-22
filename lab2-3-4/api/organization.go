package api

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	db "m1thrandir225/lab-2-3-4/db/sqlc"
	"net/http"
)

type createOrganizationRequest struct {
	Name string `json:"name" binding:"required"`
}

type updateUserRoleRequest struct {
	RoleID int64 `json:"role_id" binding:"required"`
}

type organizationIdRequest struct {
	OrganizationId int64 `uri:"org_id" binding:"required"`
}

type addUserToOrgRequest struct {
	UserID int64 `json:"user_id" binding:"required"`
	RoleID int64 `json:"role_id" binding:"required"`
}

func (server *Server) createOrganization(ctx *gin.Context) {
	var req createOrganizationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	org, err := server.store.CreateOrganization(ctx, req.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Create default roles for the organization
	err = server.store.CreateInitialRoles(ctx, org.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Add creator as moderator
	payload, err := GetPayloadFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	user, err := server.store.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	modRoleId, err := server.store.GetModeratorRole(ctx, org.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = server.store.AddUserToOrganization(ctx, db.AddUserToOrganizationParams{
		UserID: user.ID,
		OrgID:  org.ID,
		RoleID: modRoleId,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, org)
}

func (server *Server) getOrganization(ctx *gin.Context) {
	var requestData UriId

	if err := ctx.ShouldBindUri(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	org, err := server.store.GetOrganization(ctx, requestData.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, org)
}

func (server *Server) deleteOrganization(ctx *gin.Context) {
	var requestData organizationIdRequest

	if err := ctx.ShouldBindUri(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteOrganization(ctx, requestData.OrganizationId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (server *Server) addUserToOrganization(ctx *gin.Context) {
	var uriID UriId
	if err := ctx.ShouldBindUri(&uriID); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req addUserToOrgRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Only admin can add users
	_, err := server.store.AddUserToOrganization(ctx, db.AddUserToOrganizationParams{
		UserID: req.UserID,
		OrgID:  uriID.ID,
		RoleID: req.RoleID,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Status(http.StatusCreated)
}

type UserOrgIdsRequest struct {
	UserID int64 `uri:"user_id" binding:"required"`
	OrgId  int64 `uri:"id" binding:"required"`
}

func (server *Server) removeUserFromOrganization(ctx *gin.Context) {
	var uriID UserOrgIdsRequest
	if err := ctx.ShouldBindUri(&uriID); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.RemoveUserFromOrganization(ctx, db.RemoveUserFromOrganizationParams{
		UserID: uriID.UserID,
		OrgID:  uriID.OrgId,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (server *Server) updateUserRole(ctx *gin.Context) {
	var uriID UserOrgIdsRequest
	if err := ctx.ShouldBindUri(&uriID); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req updateUserRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	err := server.store.UpdateUserOrganizationRole(ctx, db.UpdateUserOrganizationRoleParams{
		UserID: uriID.UserID,
		OrgID:  uriID.OrgId,
		RoleID: req.RoleID,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.Status(http.StatusOK)
}

func (server *Server) listUserOrganizations(ctx *gin.Context) {
	payload, err := GetPayloadFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	orgs, err := server.store.ListUserOrganizations(ctx, user.ID)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.JSON(200, orgs)
}
