package api

import (
	"github.com/gin-gonic/gin"
	db "m1thrandir225/lab-2-3-4/db/sqlc"
	"net/http"
)

type createRoleRequest struct {
	Name string `json:"name" binding:"required"`
}

func (server *Server) createRole(ctx *gin.Context) {
	var uriId UriId
	if err := ctx.ShouldBindUri(&uriId); err != nil {
		ctx.JSON(http.StatusBadRequest, "missing required parameter")
		return
	}

	var req createRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	role, err := server.store.CreateRole(ctx, db.CreateRoleParams{
		Name:  req.Name,
		OrgID: uriId.ID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, role)
}

type orgRoleUriIdReqest struct {
	OrgId  int64 `uri:"id" binding:"required"`
	RoleId int64 `uri:"role_id" binding:"required"`
}

func (server *Server) updateRole(ctx *gin.Context) {
	var uriId orgRoleUriIdReqest
	if err := ctx.ShouldBindUri(&uriId); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	var req createRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	role, err := server.store.UpdateRole(ctx, db.UpdateRoleParams{
		ID:   uriId.RoleId,
		Name: req.Name,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, role)
}

func (server *Server) listOrganizationRoles(ctx *gin.Context) {
	var uriId UriId
	if err := ctx.ShouldBindUri(&uriId); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	roles, err := server.store.ListOrganizationRoles(ctx, uriId.ID)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.JSON(200, roles)
}
