package handlers

import (
	"github.com/AD-Tech-Projects/Ledger/packages/infra/database/models"
	"github.com/AD-Tech-Projects/Ledger/packages/infra/database/repositories"
	"github.com/AD-Tech-Projects/Ledger/packages/services"
	"gorm.io/gorm"
)

type HandlerResolver struct {
	Account     *AccountHandler
	Transaction *TransactionHandler
}

func NewHandlerResolver(db *gorm.DB) *HandlerResolver {
	holderRepository := repositories.NewRepository(db, models.AccountHolderToDomain, models.AccountHolderFromDomain)
	accountRepository := repositories.NewRepository(db, models.AccountToDomain, models.AccountFromDomain)
	transactionRepository := repositories.NewRepository(db, models.TransactionToDomain, models.TransactionFromDomain)

	holderService := services.NewAccountHolderService(holderRepository)
	accountService := services.NewAccountService(accountRepository)
	transactionService := services.NewTransactionService(transactionRepository)

	accountHandler := NewAccountHandler(holderService, accountService)
	transactionHandler := NewTransactionHandler(transactionService)

	return &HandlerResolver{
		Account:     accountHandler,
		Transaction: transactionHandler,
	}
}
