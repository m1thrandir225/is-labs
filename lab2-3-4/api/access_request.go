package api

import (
	"errors"
	"m1thrandir225/lab-2-3-4/auth"
	db "m1thrandir225/lab-2-3-4/db/sqlc"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type createAccessRequestRequest struct {
	ResourceID int64  `json:"resource_id" binding:"required"`
	Reason     string `json:"reason" binding:"required"`
	Duration   string `json:"duration" binding:"required"` // e.g., "24h", "48h"
}

type accessRequestResponse struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	ResourceID int64     `json:"resource_id"`
	Status     string    `json:"status"`
	Reason     string    `json:"reason"`
	ExpiresAt  time.Time `json:"expires_at"`
	CreatedAt  time.Time `json:"created_at"`
}

func (server *Server) requestAccess(ctx *gin.Context) {
	var req createAccessRequestRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	payload := ctx.MustGet(authorizationPayloadKey).(*auth.Claims)

	// Parse duration
	duration, err := time.ParseDuration(req.Duration)
	if err != nil {
		ctx.JSON(400, errorResponse(errors.New("invalid duration format")))
		return
	}

	user, err := server.store.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}
	// Create access request
	accessRequest, err := server.store.CreateAccessRequest(ctx, db.CreateAccessRequestParams{
		UserID:     user.ID,
		ResourceID: req.ResourceID,
		Status:     "pending",
		Reason:     req.Reason,
		ExpiresAt:  time.Now().Add(duration),
	})
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.JSON(201, accessRequestResponse{
		ID:         accessRequest.ID,
		UserID:     accessRequest.UserID,
		ResourceID: accessRequest.ResourceID,
		Status:     accessRequest.Status,
		Reason:     accessRequest.Reason,
		ExpiresAt:  accessRequest.ExpiresAt,
		CreatedAt:  accessRequest.CreatedAt,
	})
}

// Evaluate and update access request (Admin/Moderator only)
func (server *Server) evaluateAccessRequest(ctx *gin.Context) {
	type evaluateRequest struct {
		Status string `json:"status" binding:"required,oneof=approved rejected"`
	}

	var req evaluateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	requestIDStr := ctx.Param("request_id")

	requestID, err := strconv.ParseInt(requestIDStr, 10, 64)
	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	err = server.store.UpdateAccessRequestStatus(ctx, db.UpdateAccessRequestStatusParams{
		ID:     requestID,
		Status: req.Status,
	})
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.Status(200)
}

// List pending access requests (Admin/Moderator only)
func (server *Server) listPendingAccessRequests(ctx *gin.Context) {
	orgIDStr := ctx.Param("org_id")

	orgID, err := strconv.ParseInt(orgIDStr, 10, 64)
	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	requests, err := server.store.ListPendingAccessRequests(ctx, orgID)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.JSON(200, requests)
}

// List user's access requests
func (server *Server) listUserAccessRequests(ctx *gin.Context) {
	payload := ctx.MustGet(authorizationPayloadKey).(*auth.Claims)

	user, err := server.store.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	requests, err := server.store.ListUserAccessRequests(ctx, user.ID)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.JSON(200, requests)
}

// Check active access
func (server *Server) checkActiveAccess(ctx *gin.Context) {
	payload := ctx.MustGet(authorizationPayloadKey).(*auth.Claims)
	resourceIDStr := ctx.Param("resource_id")

	resourceID, err := strconv.ParseInt(resourceIDStr, 10, 64)
	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	access, err := server.store.GetActiveAccessRequest(ctx, db.GetActiveAccessRequestParams{
		UserID:     user.ID,
		ResourceID: resourceID,
		ExpiresAt:  time.Now(),
	})
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			ctx.JSON(404, gin.H{"has_access": false})
			return
		}
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.JSON(200, gin.H{
		"has_access": true,
		"expires_at": access.ExpiresAt,
	})
}
