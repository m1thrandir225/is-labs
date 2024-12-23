package api

import (
	"github.com/gin-gonic/gin"
)

type UriId struct {
	ID int64 `uri:"id" binding:"required"`
}

func (server *Server) initializeRoutes() {
	router := gin.Default()

	v1 := router.Group("/api/v1")

	v1.POST("/login", server.login)
	v1.POST("/register", server.register)
	v1.POST("/verify-2fa", server.verifyOTP)

	// Protected routes
	authRoutes := v1.Group("/").Use(authMiddleware(server.tokenMaker))

	// Organization management
	authRoutes.POST("/organizations", server.createOrganization)
	authRoutes.GET("/organizations", server.listUserOrganizations)
	authRoutes.GET("/organizations/:id", server.getOrganization)

	// Organization user management (Admin/Moderator only)
	authRoutes.POST("/organizations/:id/users", server.requireAdmin(), server.addUserToOrganization)
	authRoutes.DELETE("/organizations/:id/users/:user_id", server.requireAdmin(), server.removeUserFromOrganization)
	authRoutes.PUT("/organizations/:id/users/:user_id/role", server.requireModerator(), server.updateUserRole)

	// Role management (Admin/Moderator only)
	authRoutes.GET("/organizations/:id/roles", server.requireAdmin(), server.listOrganizationRoles)
	authRoutes.POST("/organizations/:id/roles", server.requireModerator(), server.createRole)
	authRoutes.PUT("/organizations/:id/roles/:role_id", server.requireModerator(), server.updateRole)

	// Resource management
	authRoutes.POST("/organizations/:id/resources", server.requireAdmin(), server.createResource)
	authRoutes.GET("/organizations/:id/resources", server.listOrganizationResources)
	authRoutes.GET("/organizations/:id/resources/:resource_id", server.getResource)
	authRoutes.PUT("/organizations/:id/resources/:resource_id", server.requireAdmin(), server.updateResource)
	authRoutes.DELETE("/organizations/:id/resources/:resource_id", server.requireAdmin(), server.deleteResource)

	// Resource permissions (Admin/Moderator only)
	authRoutes.POST("/organizations/:id/resources/:resource_id/permissions", server.requireAdmin(), server.setResourcePermissions)
	authRoutes.GET("/organizations/:id/resources/:resource_id/permissions/:role_id", server.requireAdmin(), server.getResourcePermissions)

	// Just-in-time access management
	authRoutes.POST("/resources/:id/access", server.requestAccess) // Request 15-min access
	authRoutes.GET("/resources/:id/access", server.checkAccess)    // Check current access status
	authRoutes.GET("/access/active", server.listActiveAccess)      // List all active access grants

	// User profile and settings
	authRoutes.GET("/me", server.getCurrentUser)
	authRoutes.PUT("/me", server.updateUser)
	authRoutes.PUT("/me/password", server.updatePassword)

	server.router = router
}
