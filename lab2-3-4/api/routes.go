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
	authRoutes.POST("/organizations", server.createOrganization)   // okay
	authRoutes.GET("/organizations", server.listUserOrganizations) //okay
	authRoutes.GET("/organizations/:id", server.getOrganization)   // okay

	// Organization user management (Admin/Moderator only)
	authRoutes.POST("/organizations/:id/users", server.requireAdmin(), server.addUserToOrganization)                 // okay
	authRoutes.DELETE("/organizations/:id/users/:user_id", server.requireAdmin(), server.removeUserFromOrganization) //okay
	authRoutes.PUT("/organizations/:id/users/:user_id/role", server.requireAdmin(), server.updateUserRole)           //okay

	// Role management (Admin/Moderator only)
	authRoutes.GET("/organizations/:id/roles", server.requireAdmin(), server.listOrganizationRoles) //okay
	authRoutes.POST("/organizations/:id/roles", server.requireAdmin(), server.createRole)           //okays
	authRoutes.PUT("/organizations/:id/roles/:role_id", server.requireAdmin(), server.updateRole)   //okay

	// Resource management
	authRoutes.POST("/organizations/:id/resources", server.requireAdmin(), server.createResource)                //okay
	authRoutes.GET("/organizations/:id/resources", server.listOrganizationResources)                             //okay
	authRoutes.GET("/organizations/:id/resources/:resource_id", server.getResource)                              //okay
	authRoutes.PUT("/organizations/:id/resources/:resource_id", server.requireAdmin(), server.updateResource)    //okay
	authRoutes.DELETE("/organizations/:id/resources/:resource_id", server.requireAdmin(), server.deleteResource) //okay

	// Resource permissions (Admin/Moderator only)
	authRoutes.POST("/organizations/:id/resources/:resource_id/permissions", server.requireAdmin(), server.setResourcePermissions)         //okay
	authRoutes.GET("/organizations/:id/resources/:resource_id/permissions/:role_id", server.requireAdmin(), server.getResourcePermissions) //okay

	// Just-in-time access management
	authRoutes.POST("/resources/:id/access", server.requestAccess) //okay // Request 15-min access
	authRoutes.GET("/resources/:id/access", server.checkAccess)    //okay    // Check current access status
	authRoutes.GET("/access/active", server.listActiveAccess)      //okay      // List all active access grants

	// User profile and settings
	authRoutes.GET("/me", server.getCurrentUser) //okay
	authRoutes.PUT("/me", server.updateUser)     //okay
	authRoutes.PUT("/me/password", server.updatePassword)

	server.router = router
}
