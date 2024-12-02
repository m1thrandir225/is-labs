package main

import (
	"database/sql"
	"log"
	"m1thrandir225/lab-2-3-4/api"
	"m1thrandir225/lab-2-3-4/auth"
	db "m1thrandir225/lab-2-3-4/db/sqlc"
	"m1thrandir225/lab-2-3-4/mail"
	"m1thrandir225/lab-2-3-4/util"

	_ "github.com/mattn/go-sqlite3"
)

type ginServerConfig struct {
	Config      util.Config
	Store       db.Store
	otpService  *auth.OTPService
	tokenMaker  auth.TokenMaker
	mailService mail.MailService
}

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal(err.Error())
	}

	database, err := sql.Open("sqlite3", config.DBSource)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer database.Close()

	store := db.NewSQLiteStore(database)
	otpService := auth.NewOTPService(store)

	mailService := mail.NewResendMail(
		config.SMTPHost,
		config.SMTPPort,
		config.SMTPUsername,
		config.SMTPPassword,
	)
	tokenMaker := auth.NewJWTMaker([]byte(config.JWTKey))

	serverConfig := ginServerConfig{
		Config:      config,
		Store:       store,
		otpService:  otpService,
		tokenMaker:  tokenMaker,
		mailService: mailService,
	}

	runGinServer(serverConfig)

}

func runGinServer(serverConfig ginServerConfig) {
	server, err := api.NewServer(
		serverConfig.Store,
		serverConfig.otpService,
		serverConfig.tokenMaker,
		serverConfig.Config,
		serverConfig.mailService,
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = server.Start(serverConfig.Config.HTTPServerAddress)
}
