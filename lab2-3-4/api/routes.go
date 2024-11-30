package api

import "github.com/gin-gonic/gin"

func (server *Server) initializeRoutes() {
	router := gin.Default()

	v1 := router.Group("/api/v1")
	protectedRoutes := v1.Group("/").Use(TwoFactorMiddleware(server.tokenMaker))

	v1.POST("/login", server.Login)
	v1.POST("/register", server.Register)

	//Verify 2-fa route
	v1.POST("/verify-2fa", server.VerifyOTP).Use(LimitedAccessMiddleware(server.tokenMaker))

	protectedRoutes.GET("/home", server.Home)

	server.router = router
}
