package api

import (
	"github.com/gin-gonic/gin"
	"m1thrandir225/lab-2-3-4/auth"
	db "m1thrandir225/lab-2-3-4/db/sqlc"
	"m1thrandir225/lab-2-3-4/util"
)

type Server struct {
	store                 db.Store
	authenticationManager *auth.AuthenticationManager
	tokenMaker            auth.TokenMaker
	config                util.Config
	router                *gin.Engine
}

func NewServer(store db.Store, am *auth.AuthenticationManager, tm auth.TokenMaker, config util.Config) (*Server, error) {
	server := &Server{
		store:                 store,
		authenticationManager: am,
		tokenMaker:            tm,
		config:                config,
	}

	server.initializeRoutes()

	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
