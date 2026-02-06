package services

import (
	"e-wallet-go/internal/models"
	"e-wallet-go/internal/repository"
)

type WalletService interface {
	GetUserBalance(userID string) (*models.Wallet, error)
	
	UserTopUp(userID string, amount float64) (*models.Wallet, error)
	UserWithdraw(userID string, amount float64) (*models.Wallet, error)

	AdminTopUpUser(targetUserID string, amount float64) (*models.Wallet, error)
	AdminDeductUser(targetUserID string, amount float64) (*models.Wallet, error)

	GetTransactionHistory(userID string, search string) ([]models.Transaction, error)
	GetAdminGlobalTransactions(search string) ([]models.Transaction, error)
}

type walletService struct {
	walletRepo repository.WalletRepository
}

func NewWalletService(walletRepo repository.WalletRepository) WalletService {
	return &walletService{walletRepo: walletRepo}
}

func (s *walletService) GetUserBalance(userID string) (*models.Wallet, error) {
	return s.walletRepo.GetBalanceByUserID(userID)
}

func (s *walletService) UserTopUp(userID string, amount float64) (*models.Wallet, error) {
	desc := "Top up by user"
	if err := s.walletRepo.CreditBalance(userID, amount, models.TRX_DEPOSIT, desc); err != nil {
		return nil, err
	}
	return s.walletRepo.GetBalanceByUserID(userID)
}

func (s *walletService) UserWithdraw(userID string, amount float64) (*models.Wallet, error) {
	desc := "Withdrawal by user"
	if err := s.walletRepo.DebitBalance(userID, amount, models.TRX_WITHDRAW, desc); err != nil {
		return nil, err
	}
	return s.walletRepo.GetBalanceByUserID(userID)
}

func (s *walletService) AdminTopUpUser(targetUserID string, amount float64) (*models.Wallet, error) {
	desc := "Bonus/Correction by Admin"
	if err := s.walletRepo.CreditBalance(targetUserID, amount, models.TRX_ADMIN_TOPUP, desc); err != nil {
		return nil, err
	}
	return s.walletRepo.GetBalanceByUserID(targetUserID)
}

func (s *walletService) AdminDeductUser(targetUserID string, amount float64) (*models.Wallet, error) {
	desc := "Penalty/Correction by Admin"
	if err := s.walletRepo.DebitBalance(targetUserID, amount, models.TRX_ADMIN_DEDUCTION, desc); err != nil {
		return nil, err
	}
	return s.walletRepo.GetBalanceByUserID(targetUserID)
}

func (s *walletService) GetTransactionHistory(userID string, search string) ([]models.Transaction, error) {
	wallet, err := s.walletRepo.GetBalanceByUserID(userID)
	if err != nil {
		return nil, err
	}
	return s.walletRepo.GetTransactionsByWalletID(wallet.ID.String(), search)
}

func (s *walletService) GetAdminGlobalTransactions(search string) ([]models.Transaction, error) {
	return s.walletRepo.GetAllTransactions(search)
}