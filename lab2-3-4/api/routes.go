package api

import (
	"github.com/gin-gonic/gin"
)

func (server *Server) initializeRoutes() {
	router := gin.Default()

	v1 := router.Group("/api/v1")

	v1.POST("/login", server.login)
	v1.POST("/register", server.register)
	v1.POST("/verify-2fa", server.verifyOTP)

	// Protected routes
	authRoutes := v1.Group("/").Use(authMiddleware(server.tokenMaker))

	// Resource management (Admin/Moderator only)
	authRoutes.POST("/organizations/:org_id/resources", server.requireAdmin(), server.createResource)

	// Resource access (all authenticated users)
	authRoutes.GET("/organizations/:org_id/resources", server.listOrganizationResources)
	authRoutes.GET("/organizations/:org_id/resources/:resource_id", server.getResource)

	// Access request workflow
	authRoutes.POST("/access-requests", server.requestAccess)                       // Request access
	authRoutes.GET("/access-requests/me", server.listUserAccessRequests)            // List own requests
	authRoutes.GET("/access-requests/check/:resource_id", server.checkActiveAccess) // Check active access

	// Access request management (Admin/Moderator only)
	authRoutes.GET("/organizations/:org_id/access-requests", server.requireAdmin(), server.listPendingAccessRequests)
	authRoutes.PUT("/access-requests/:request_id", server.requireAdmin(), server.evaluateAccessRequest)

	server.router = router
}
