package main

import (
	"os"

	logger "github.com/AD-Tech-Projects/Go-Libs/packages/logger"
	account "github.com/AD-Tech-Projects/Ledger/packages/domain/account"
	holder "github.com/AD-Tech-Projects/Ledger/packages/domain/holder"
	transaction "github.com/AD-Tech-Projects/Ledger/packages/domain/transaction"
	database "github.com/AD-Tech-Projects/Ledger/packages/infra/database"
	"github.com/joho/godotenv"
)

func main() {
	loadEnvs()

	connect := database.Connect()
	connect.AutoMigrate()

	testLedger()

	logger.Debug("Done")
}

func testLedger() {
	companyHolder := holder.NewCompanyAccountHolder("1", "12345678000195")
	personHolder := holder.NewIndividualAccountHolder("2", "123456789")

	companyAccount := account.NewAccount(
		"1",
		account.ACTIVE,
		account.CHECKING,
		account.USD,
		10000,
		&companyHolder,
	)
	personAccount := account.NewAccount(
		"2",
		account.ACTIVE,
		account.CHECKING,
		account.USD,
		0,
		&personHolder,
	)
	logger.Info("Company account:", companyAccount)
	logger.Info("Person account:", personAccount)

	transaction := transaction.MakeTransaction(
		"1",
		5000,
		transaction.DEPOSIT,
		&companyAccount,
		&personAccount,
	)
	logger.Info("Transaction created successfully", transaction)

	transaction.Commit()
	logger.Info("Transaction completed successfully", transaction)
	logger.Info("Company account:", companyAccount)
	logger.Info("Person account:", personAccount)
}

func loadEnvs() {
	err := godotenv.Load()
	if err != nil {
		logger.Error("Failed to load .env file, using system environment variables instead")
	}

	port := os.Getenv("APP_PORT")
	logger.Debug("Starting the application on port: " + port)
}
