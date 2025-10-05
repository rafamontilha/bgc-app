package main

import (
	"log"

	"bgc-app/internal/app"
	"bgc-app/internal/config"
	"bgc-app/internal/repository/postgres"
)

func main() {
	cfg := config.LoadConfig()

	db := postgres.MustConnect(cfg)
	defer db.Close()

	server := app.NewServer(cfg, db)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
