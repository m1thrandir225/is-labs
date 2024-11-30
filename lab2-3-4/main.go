package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"m1thrandir225/lab-2-3-4/api"
	"m1thrandir225/lab-2-3-4/auth"
	db "m1thrandir225/lab-2-3-4/db/sqlc"
	"m1thrandir225/lab-2-3-4/util"
)

type ginServerConfig struct {
	Config      util.Config
	Store       db.Store
	authManager *auth.AuthenticationManager
	tokenMaker  auth.TokenMaker
}

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal(err.Error())
	}

	database, err := sql.Open("sqlite3", "./sqlite.db")
	if err != nil {
		log.Fatal(err.Error())
	}

	defer database.Close()

	store := db.NewSQLiteStore(database)
	manager := auth.NewAuthenticationManager()
	tokenMaker := auth.NewJWTMaker([]byte(config.JWTKey))

	serverConfig := ginServerConfig{
		Config:      config,
		Store:       store,
		authManager: manager,
		tokenMaker:  tokenMaker,
	}

	runGinServer(serverConfig)

}

func runGinServer(serverConfig ginServerConfig) {
	server, err := api.NewServer(
		serverConfig.Store,
		serverConfig.authManager,
		serverConfig.tokenMaker,
		serverConfig.Config,
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = server.Start(serverConfig.Config.HTTPServerAddress)
}
