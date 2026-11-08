package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/AD-Tech-Projects/Ledger/packages/domain"
	"github.com/AD-Tech-Projects/Ledger/packages/infra/database/models"
	"github.com/AD-Tech-Projects/Ledger/packages/infra/database/repositories"
)

type WalletService struct {
	walletRepository *repositories.GormRepository[models.WalletModel, domain.Wallet]
}

func NewWalletService(
	walletRepository *repositories.GormRepository[models.WalletModel, domain.Wallet],
) *WalletService {
	return &WalletService{
		walletRepository: walletRepository,
	}
}

var (
	ErrInvalidWallet  = errors.New("invalid wallet")
	ErrWalletNotFound = errors.New("wallet not found")
)

func (s *WalletService) Create(ctx context.Context, wallet *domain.Wallet) (*domain.Wallet, error) {
	if wallet == nil {
		return nil, fmt.Errorf("%w: wallet is required", ErrInvalidWallet)
	}

	if err := validateWallet(wallet); err != nil {
		return nil, err
	}

	if wallet.ID == "" {
		wallet.ID = uuid.New().String()
	}

	wallet.FreeBalance = 0
	wallet.ReservedBalance = 0
	wallet.HoldBalance = 0

	if err := s.walletRepository.Create(ctx, wallet); err != nil {
		return nil, fmt.Errorf("failed to create wallet: %w", err)
	}

	return wallet, nil
}

func (s *WalletService) GetById(ctx context.Context, id string) (*domain.Wallet, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, fmt.Errorf("%w: wallet id must be a valid UUID", ErrInvalidWallet)
	}

	wallet, err := s.walletRepository.GetByID(ctx, id,
		repositories.Join[models.WalletModel]("Asset"),
		repositories.Join[models.WalletModel]("Owner"),
		repositories.Join[models.WalletModel]("Account"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet by id: %w", err)
	}
	if wallet == nil {
		return nil, ErrWalletNotFound
	}

	return wallet, nil
}

func (s *WalletService) GetAll(ctx context.Context, filter *domain.Wallet) ([]*domain.Wallet, error) {
	wallets, err := s.walletRepository.FindAll(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallets: %w", err)
	}

	return wallets, nil
}

func (s *WalletService) Deactivate(ctx context.Context, id string) error {
	wallet, err := s.GetById(ctx, id)
	if err != nil {
		return err
	}

	wallet.Status = domain.WALLET_INACTIVE
	if err := s.walletRepository.Update(ctx, wallet); err != nil {
		return fmt.Errorf("failed to deactivate wallet: %w", err)
	}

	return nil
}
