package api

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	db "m1thrandir225/lab-2-3-4/db/sqlc"
	"net/http"
	"time"
)

const (
	accessDuration = 15 * time.Minute
)

type requestAccessRequest struct {
	Reason string `json:"reason" binding:"required"`
}

func (server *Server) requestAccess(ctx *gin.Context) {
	var uriId UriId
	if err := ctx.ShouldBindUri(&uriId); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req requestAccessRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	payload, err := GetPayloadFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("user not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Get resource to check organization
	resource, err := server.store.GetResource(ctx, uriId.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("resource not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Get user's role in the organization
	userOrg, err := server.store.GetUserOrganization(ctx, db.GetUserOrganizationParams{
		UserID: user.ID,
		OrgID:  resource.OrgID,
	})
	log.Println(userOrg.RoleID)
	if err != nil {
		ctx.JSON(http.StatusForbidden, errorResponse(errors.New("not a member of organization")))
		return
	}

	// Check role permissions for the resource
	permissions, err := server.store.GetRolePermissions(ctx, db.GetRolePermissionsParams{
		RoleID:     userOrg.RoleID,
		ResourceID: uriId.ID,
	})
	log.Println(permissions)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("no permissions defined")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if !permissions.CanRead {
		ctx.JSON(http.StatusForbidden, errorResponse(errors.New("insufficient permissions")))
		return
	}

	accessRequest, err := server.store.CreateAccessRequest(ctx, db.CreateAccessRequestParams{
		UserID:     user.ID,
		ResourceID: uriId.ID,
		Status:     "approved",
		Reason:     req.Reason,
		ExpiresAt:  time.Now().Add(accessDuration),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"access_id":  accessRequest.ID,
		"status":     "approved",
		"expires_at": accessRequest.ExpiresAt,
		"duration":   "15 minutes",
	})
}

func (server *Server) checkAccess(ctx *gin.Context) {
	var uriId UriId
	if err := ctx.ShouldBindUri(&uriId); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := GetPayloadFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	access, err := server.store.GetActiveAccessRequest(ctx, db.GetActiveAccessRequestParams{
		UserID:     user.ID,
		ResourceID: uriId.ID,
		ExpiresAt:  time.Now(),
	})

	if err != nil {
		remainingTime := time.Duration(0)
		ctx.JSON(http.StatusOK, gin.H{
			"has_access":     false,
			"message":        "No active access. Please request new access.",
			"time_remaining": remainingTime.String(),
		})
		return
	}

	remainingTime := time.Until(access.ExpiresAt).Round(time.Second)
	ctx.JSON(http.StatusOK, gin.H{
		"has_access":     true,
		"expires_at":     access.ExpiresAt,
		"time_remaining": remainingTime.String(),
	})
}

func (server *Server) listActiveAccess(ctx *gin.Context) {
	payload, err := GetPayloadFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	activeAccess, err := server.store.ListActiveUserAccess(ctx, db.ListActiveUserAccessParams{
		UserID:    user.ID,
		ExpiresAt: time.Now(),
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Add remaining time for each access
	var response []gin.H
	for _, access := range activeAccess {
		remainingTime := time.Until(access.ExpiresAt).Round(time.Second)
		response = append(response, gin.H{
			"access_id":      access.ID,
			"resource_id":    access.ResourceID,
			"resource_name":  access.ResourceName,
			"expires_at":     access.ExpiresAt,
			"time_remaining": remainingTime.String(),
		})
	}

	ctx.JSON(http.StatusOK, response)
}
