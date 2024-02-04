package main

import (
	"kakeibo-api/db"
	"kakeibo-api/server"
	"kakeibo-api/utils"
	"log"
	"os"
)

// @title Kakeibo API Service
// @version 1.0
// @BasePath /v1/
func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	if os.Getenv("TIMEZONE") != "" {
		utils.SetTimezone(os.Getenv("TIMEZONE"))
	}

	db.Init()
	log.Println("Database connected.")

	log.Println("Starting HTTP Server...")
	server.Start()
	log.Println("Stopping HTTP Server.")
}
