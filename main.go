package main

import (
	logger "github.com/AD-Tech-Projects/Go-Libs/packages/logger"
	account "github.com/AD-Tech-Projects/Ledger/packages/domain/account"
	holder "github.com/AD-Tech-Projects/Ledger/packages/domain/holder"
	transaction "github.com/AD-Tech-Projects/Ledger/packages/domain/transaction"
)

func main() {
	logger.Debug("Starting the application...")

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

	transaction := transaction.MakeTransaction(
		"1",
		5000,
		transaction.DEPOSIT,
		&companyAccount,
		&personAccount,
	)
	transaction.Commit()

	logger.Info("Transaction completed successfully")
}
