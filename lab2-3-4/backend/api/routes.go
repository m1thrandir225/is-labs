package api

import "github.com/gin-gonic/gin"

func (server *Server) initializeRoutes() {
	router := gin.Default()

	v1 := router.Group("/api/v1")

	v1.POST("/login", server.Login)
	v1.POST("/register", server.Register)

	//Verify 2-fa route
	v1.POST("/verify-2fa", server.VerifyOTP)

	//protectedRoutes.GET("/home", server.Home)

	server.router = router
}
