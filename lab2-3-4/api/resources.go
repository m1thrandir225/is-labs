package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	db "m1thrandir225/lab-2-3-4/db/sqlc"
	"strconv"
	"time"
)

type createResourceRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type resourceResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	OrgID     int64     `json:"org_id"`
	CreatedAt time.Time `json:"created_at"`
}

// Create resource (Admin/Moderator only)
func (server *Server) createResource(ctx *gin.Context) {
	var req createResourceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	orgIDStr := ctx.Param("org_id")
	orgID, err := strconv.ParseInt(orgIDStr, 10, 64)
	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	resource, err := server.store.CreateResource(ctx, db.CreateResourceParams{
		Name:  req.Name,
		OrgID: orgID,
	})
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.JSON(201, resourceResponse{
		ID:        resource.ID,
		Name:      resource.Name,
		OrgID:     resource.OrgID,
		CreatedAt: resource.CreatedAt,
	})
}

// Get single resource
func (server *Server) getResource(ctx *gin.Context) {
	resourceIDStr := ctx.Param("resource_id")
	resourceID, err := strconv.ParseInt(resourceIDStr, 10, 64)
	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	resource, err := server.store.GetResource(ctx, resourceID)
	if err != nil {
		ctx.JSON(404, errorResponse(errors.New("resource not found")))
		return
	}

	ctx.JSON(200, resourceResponse{
		ID:        resource.ID,
		Name:      resource.Name,
		OrgID:     resource.OrgID,
		CreatedAt: resource.CreatedAt,
	})
}

// List organization resources
func (server *Server) listOrganizationResources(ctx *gin.Context) {
	orgIDStr := ctx.Param("org_id")
	orgID, err := strconv.ParseInt(orgIDStr, 10, 64)
	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	resources, err := server.store.ListOrganizationResources(ctx, orgID)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
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

	ctx.JSON(200, response)
}
