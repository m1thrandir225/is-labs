package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (server *Server) Home(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Hello World",
	})
}
