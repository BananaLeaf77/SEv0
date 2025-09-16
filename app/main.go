package main

import (
	"SEv0/config"
	api "SEv0/http"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file, Error is: ", err)
	}

	_, address, err := config.BootDB()
	if err != nil {
		log.Fatal("Could not connect to the database, Error is: ", err)
	}

	whatsmeow, ctx, err := config.InitWA(*address)
	if err != nil {
		log.Fatal("Could not initialize WhatsApp client, Error is: ", err)
	}

	api.InitApi(whatsmeow, ctx)
}
