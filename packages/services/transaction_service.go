package services

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/AD-Tech-Projects/Ledger/packages/domain"
	"github.com/AD-Tech-Projects/Ledger/packages/infra/database/models"
	"github.com/AD-Tech-Projects/Ledger/packages/infra/database/repositories"
)

type TransactionService struct {
	transactionRepository *repositories.GormRepository[models.TransactionModel, domain.Transaction]
	walletRepository      *repositories.GormRepository[models.WalletModel, domain.Wallet]
}

func NewTransactionService(
	transactionRepository *repositories.GormRepository[models.TransactionModel, domain.Transaction],
	walletRepository *repositories.GormRepository[models.WalletModel, domain.Wallet],
) *TransactionService {
	return &TransactionService{
		transactionRepository: transactionRepository,
		walletRepository:      walletRepository,
	}
}

var (
	ErrInvalidTransaction  = errors.New("invalid transaction")
	ErrTransactionNotFound = errors.New("transaction not found")
)

func (s *TransactionService) Create(ctx context.Context, input CreateTransactionInput) (*domain.Transaction, error) {
	if err := validateCreateTransactionInput(&input); err != nil {
		return nil, err
	}

	transaction := &domain.Transaction{
		ID:               uuid.NewString(),
		Amount:           input.Amount,
		Type:             input.Type,
		Status:           domain.TRANSACTION_PENDING,
		Description:      input.Description,
		ReferenceID:      input.ReferenceID,
		ReversalOriginID: input.ReversalOriginID,
	}

	var transactionFailure error
	err := s.transactionRepository.Transaction(ctx, func(tx *gorm.DB) error {
		walletRepository := s.walletRepository.WithDB(tx)
		transactionRepository := s.transactionRepository.WithDB(tx)

		sourceWalletID, destinationWalletID := input.SourceWalletID, input.DestinationWalletID
		if transaction.Type == domain.TRANSACTION_REVERSAL {
			origin, err := lockReversalOrigin(ctx, transactionRepository, *input.ReversalOriginID)
			if err != nil {
				return err
			}
			transaction.ReversalOrigin = origin
			if transaction.Amount == 0 {
				transaction.Amount = origin.Amount
			}
			sourceWalletID, destinationWalletID = getReversalWalletIDs(origin)
		}

		wallets, err := lockWallets(ctx, walletRepository, sourceWalletID, destinationWalletID)
		if err != nil {
			return err
		}
		transaction.SourceWallet = walletByID(wallets, sourceWalletID)
		transaction.DestinationWallet = walletByID(wallets, destinationWalletID)
		transaction.SourceWalletID = sourceWalletID
		transaction.DestinationWalletID = destinationWalletID
		transaction.AssetID = transactionAssetID(transaction)

		// NOTE - a rejected transaction is recorded as FAILED in this same database transaction, which already holds the
		// wallet locks, and the wallets are not updated, so their balances stay unchanged
		transactionFailure = validateTransaction(transaction)
		if transactionFailure == nil {
			transactionFailure = applyTransaction(transaction)
		}
		if transactionFailure != nil {
			transaction.Status = domain.TRANSACTION_FAILED
			if err := transactionRepository.Create(ctx, transaction); err != nil {
				return fmt.Errorf("%w; failed to persist failed transaction: %v", transactionFailure, err)
			}
			return nil
		}

		for _, wallet := range wallets {
			if err := walletRepository.Update(ctx, wallet); err != nil {
				return fmt.Errorf("failed to update wallet %s: %w", wallet.ID, err)
			}
		}

		transaction.Status = domain.TRANSACTION_COMPLETED
		now := time.Now()
		transaction.CompletedAt = &now
		if err := transactionRepository.Create(ctx, transaction); err != nil {
			return fmt.Errorf("failed to create transaction: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	if transactionFailure != nil {
		return nil, transactionFailure
	}

	return transaction, nil
}

// lockReversalOrigin locks the original transaction row, so concurrent reversals of the same transaction are serialized,
// and rejects transactions that were already reversed.
func lockReversalOrigin(
	ctx context.Context,
	transactionRepository *repositories.GormRepository[models.TransactionModel, domain.Transaction],
	originID string,
) (*domain.Transaction, error) {
	origin, err := transactionRepository.GetByID(ctx, originID, repositories.ForUpdate[models.TransactionModel]())
	if err != nil {
		return nil, fmt.Errorf("failed to get reversal origin transaction: %w", err)
	}
	if origin == nil {
		return nil, fmt.Errorf("%w: reversal origin transaction %s not found", ErrInvalidTransaction, originID)
	}

	existingReversal, err := transactionRepository.FindOne(ctx, &domain.Transaction{
		Type:             domain.TRANSACTION_REVERSAL,
		Status:           domain.TRANSACTION_COMPLETED,
		ReversalOriginID: &originID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to check existing reversals: %w", err)
	}
	if existingReversal != nil {
		return nil, fmt.Errorf("%w: transaction %s was already reversed by %s", ErrInvalidTransaction, originID, existingReversal.ID)
	}

	return origin, nil
}

// lockWallets loads the wallets with SELECT ... FOR UPDATE, always in the same (sorted) order to avoid deadlocks
// between concurrent transactions that touch the same wallets.
func lockWallets(
	ctx context.Context,
	walletRepository *repositories.GormRepository[models.WalletModel, domain.Wallet],
	walletIDs ...*string,
) ([]*domain.Wallet, error) {
	ids := make([]string, 0, len(walletIDs))
	for _, id := range walletIDs {
		if id != nil {
			ids = append(ids, *id)
		}
	}
	sort.Strings(ids)

	wallets := make([]*domain.Wallet, 0, len(ids))
	for _, id := range ids {
		wallet, err := walletRepository.GetByID(ctx, id, repositories.ForUpdate[models.WalletModel]())
		if err != nil {
			return nil, fmt.Errorf("failed to get wallet %s: %w", id, err)
		}
		if wallet == nil {
			return nil, fmt.Errorf("%w: wallet %s", ErrWalletNotFound, id)
		}
		wallets = append(wallets, wallet)
	}

	return wallets, nil
}

func walletByID(wallets []*domain.Wallet, id *string) *domain.Wallet {
	if id == nil {
		return nil
	}
	for _, wallet := range wallets {
		if wallet.ID == *id {
			return wallet
		}
	}
	return nil
}

func transactionAssetID(transaction *domain.Transaction) string {
	for _, wallet := range []*domain.Wallet{transaction.SourceWallet, transaction.DestinationWallet} {
		if wallet != nil {
			return wallet.AssetID
		}
	}
	return ""
}

func (s *TransactionService) GetById(ctx context.Context, id string) (*domain.Transaction, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, fmt.Errorf("%w: transaction id must be a valid UUID", ErrInvalidTransaction)
	}

	transaction, err := s.transactionRepository.GetByID(ctx, id,
		repositories.Join[models.TransactionModel]("Asset"),
		repositories.Join[models.TransactionModel]("SourceWallet"),
		repositories.Join[models.TransactionModel]("SourceWallet.Owner"),
		repositories.Join[models.TransactionModel]("DestinationWallet"),
		repositories.Join[models.TransactionModel]("DestinationWallet.Owner"),
		repositories.Join[models.TransactionModel]("ReversalOrigin"))
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction by id: %w", err)
	}
	if transaction == nil {
		return nil, ErrTransactionNotFound
	}

	return transaction, nil
}

func (s *TransactionService) GetAll(ctx context.Context, filter *domain.Transaction) ([]*domain.Transaction, error) {
	transactions, err := s.transactionRepository.FindAll(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	return transactions, nil
}
