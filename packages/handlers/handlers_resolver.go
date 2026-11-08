package handlers

import (
	"github.com/AD-Tech-Projects/Ledger/packages/infra/database/models"
	"github.com/AD-Tech-Projects/Ledger/packages/infra/database/repositories"
	"github.com/AD-Tech-Projects/Ledger/packages/services"
	"gorm.io/gorm"
)

type HandlerResolver struct {
	Wallet      *WalletHandler
	Transaction *TransactionHandler
}

func NewHandlerResolver(db *gorm.DB) *HandlerResolver {
	// * Repositories
	assetRepository := repositories.NewRepository(db, models.AssetToDomain, models.AssetFromDomain)
	ownerRepository := repositories.NewRepository(db, models.OwnerToDomain, models.OwnerFromDomain)
	walletRepository := repositories.NewRepository(db, models.WalletToDomain, models.WalletFromDomain)
	transactionRepository := repositories.NewRepository(db, models.TransactionToDomain, models.TransactionFromDomain)

	// ? Services
	assetService := services.NewAssetService(assetRepository)
	ownerService := services.NewOwnerService(ownerRepository)
	walletService := services.NewWalletService(walletRepository)
	transactionService := services.NewTransactionService(transactionRepository, walletRepository)

	// ! Handlers
	walletHandler := NewWalletHandler(ownerService, assetService, walletService)
	transactionHandler := NewTransactionHandler(transactionService)

	return &HandlerResolver{
		Wallet:      walletHandler,
		Transaction: transactionHandler,
	}
}
