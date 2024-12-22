package api

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	db "m1thrandir225/lab-2-3-4/db/sqlc"
	"net/http"
	"time"
)

type createResourceRequest struct {
	Name string `json:"name" binding:"required"`
}

type resourceResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	OrgID     int64     `json:"org_id"`
	CreatedAt time.Time `json:"created_at"`
}

type resourcePermissionRequest struct {
	RoleID    int64 `json:"role_id" binding:"required"`
	CanRead   bool  `json:"can_read"`
	CanWrite  bool  `json:"can_write"`
	CanDelete bool  `json:"can_delete"`
}

type getResourcePermissionRequest struct {
	RoleID int64 `json:"resource_id" binding:"required"`
}

type updateResourceRequest struct {
	Name string `json:"name"`
}

// Create resource (Admin/Moderator only)
func (server *Server) createResource(ctx *gin.Context) {
	var uriId UriId
	if err := ctx.ShouldBindUri(&uriId); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var req createResourceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	resource, err := server.store.CreateResource(ctx, db.CreateResourceParams{
		Name:  req.Name,
		OrgID: uriId.ID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, resourceResponse{
		ID:        resource.ID,
		Name:      resource.Name,
		OrgID:     resource.OrgID,
		CreatedAt: resource.CreatedAt,
	})
}

type orgIdResourceIdRequest struct {
	OrgId      int64 `uri:"id" binding:"required"`
	ResourceID int64 `uri:"resource_id" binding:"required"`
}

type orgIdResourceIdRoleIdRequest struct {
	OrgId      int64 `uri:"id" binding:"required"`
	ResourceID int64 `uri:"resource_id" binding:"required"`
	RoleID     int64 `uri:"role_id" binding:"required"`
}

func (server *Server) getResource(ctx *gin.Context) {
	var uriId orgIdResourceIdRequest
	if err := ctx.ShouldBindUri(&uriId); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := GetPayloadFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	access, err := server.store.GetActiveAccessRequest(ctx, db.GetActiveAccessRequestParams{
		UserID:     user.ID,
		ResourceID: uriId.ResourceID,
		ExpiresAt:  time.Now(),
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusForbidden, errorResponse(errors.New("access expired, please request a new access request")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if time.Now().After(access.ExpiresAt) {
		ctx.JSON(http.StatusForbidden, errorResponse(errors.New("access expired, please request a new access request")))
		return
	}

	resource, err := server.store.GetResource(ctx, uriId.ResourceID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("resource not found")))
		return
	}

	ctx.JSON(http.StatusOK, resourceResponse{
		ID:        resource.ID,
		Name:      resource.Name,
		OrgID:     resource.OrgID,
		CreatedAt: resource.CreatedAt,
	})
}

// List organization resources
func (server *Server) listOrganizationResources(ctx *gin.Context) {
	var uriId UriId
	if err := ctx.ShouldBindUri(&uriId); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resources, err := server.store.ListOrganizationResources(ctx, uriId.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var response []resourceResponse
	for _, resource := range resources {
		response = append(response, resourceResponse{
			ID:        resource.ID,
			Name:      resource.Name,
			OrgID:     resource.OrgID,
			CreatedAt: resource.CreatedAt,
		})
	}

	ctx.JSON(http.StatusOK, response)
}

func (server *Server) setResourcePermissions(ctx *gin.Context) {
	var uriId orgIdResourceIdRequest
	if err := ctx.ShouldBindUri(&uriId); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req resourcePermissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	permissions, err := server.store.CreateRolePermission(ctx, db.CreateRolePermissionParams{
		RoleID:     req.RoleID,
		ResourceID: uriId.ResourceID,
		CanRead:    req.CanRead,
		CanWrite:   req.CanWrite,
		CanDelete:  req.CanDelete,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return

		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, permissions)
}

func (server *Server) getResourcePermissions(ctx *gin.Context) {
	var uriId orgIdResourceIdRoleIdRequest
	if err := ctx.ShouldBindUri(&uriId); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	permissions, err := server.store.GetRolePermissions(ctx, db.GetRolePermissionsParams{
		RoleID:     uriId.RoleID,
		ResourceID: uriId.ResourceID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return

		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, permissions)
}

func (server *Server) updateResource(ctx *gin.Context) {
	var uriId orgIdResourceIdRequest
	if err := ctx.ShouldBindUri(&uriId); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req updateResourceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	resource, err := server.store.UpdateResource(ctx, db.UpdateResourceParams{
		ID:   uriId.ResourceID,
		Name: req.Name,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, resource)
}

func (server *Server) deleteResource(ctx *gin.Context) {
	var uriId orgIdResourceIdRequest
	if err := ctx.ShouldBindUri(&uriId); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	err := server.store.DeleteResource(ctx, uriId.ResourceID)
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
