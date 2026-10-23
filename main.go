package main

import (
	"os"

	logger "github.com/AD-Tech-Projects/Go-Libs/packages/logger"
	api "github.com/AD-Tech-Projects/Ledger/packages/api"
	handlers "github.com/AD-Tech-Projects/Ledger/packages/handlers"
	database "github.com/AD-Tech-Projects/Ledger/packages/infra/database"
	"github.com/joho/godotenv"
)

func main() {
	loadEnvs()

	db := database.Connect()

	resolver := handlers.NewHandlerResolver(db)
	router := api.NewRouter(resolver)

	api.Listen(os.Getenv("APP_PORT"), router)

	logger.Debug("Done")
}

func loadEnvs() {
	err := godotenv.Load()
	if err != nil {
		logger.Error("Failed to load .env file, using system environment variables instead")
	}

}
