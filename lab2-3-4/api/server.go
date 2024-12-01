package api

import (
	"github.com/gin-gonic/gin"
	"m1thrandir225/lab-2-3-4/auth"
	db "m1thrandir225/lab-2-3-4/db/sqlc"
	"m1thrandir225/lab-2-3-4/mail"
	"m1thrandir225/lab-2-3-4/util"
)

type Server struct {
	store       db.Store
	otpService  *auth.OTPService
	tokenMaker  auth.TokenMaker
	config      util.Config
	router      *gin.Engine
	mailService mail.MailService
}

func NewServer(
	store db.Store,
	otpService *auth.OTPService,
	tm auth.TokenMaker,
	config util.Config,
	mailService mail.MailService,
) (*Server, error) {
	server := &Server{
		store:       store,
		otpService:  otpService,
		tokenMaker:  tm,
		config:      config,
		mailService: mailService,
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
