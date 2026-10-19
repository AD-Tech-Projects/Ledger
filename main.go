package main

import (
	"os"

	logger "github.com/AD-Tech-Projects/Go-Libs/packages/logger"
	database "github.com/AD-Tech-Projects/Ledger/packages/infra/database"
	"github.com/joho/godotenv"
)

func main() {
	loadEnvs()

	connect := database.Connect()
	connect.AutoMigrate()

	logger.Debug("Done")
}

func loadEnvs() {
	err := godotenv.Load()
	if err != nil {
		logger.Error("Failed to load .env file, using system environment variables instead")
	}

	port := os.Getenv("APP_PORT")
	logger.Info("Starting the application on port: " + port)
}
