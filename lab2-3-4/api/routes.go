package api

import (
	"github.com/gin-gonic/gin"
	"m1thrandir225/lab-2-3-4/dto"
)

func (server *Server) initializeRoutes() {
	router := gin.Default()

	v1 := router.Group("/api/v1")
	protectedRoutes := v1.Group("/")
	protectedRoutes.Use(authMiddleware(server.tokenMaker))

	adminRoutes := protectedRoutes.Group("/admin")
	adminRoutes.Use(roleMiddleware(dto.RoleAdmin, server.store))

	moderatorRoutes := protectedRoutes.Group("/moderator")
	moderatorRoutes.Use(roleMiddleware(dto.RoleModerator, server.store))

	v1.POST("/login", server.Login)
	v1.POST("/register", server.Register)
	v1.POST("/refresh-token", server.RefreshToken)

	//Verify 2-fa route
	v1.POST("/verify-2fa", server.VerifyOTP)

	protectedRoutes.GET("/home", server.Home)
	adminRoutes.POST("/get-role", server.GetUserRole)
	moderatorRoutes.POST("/update-role", server.UpdateUserRole)

	server.router = router
}
