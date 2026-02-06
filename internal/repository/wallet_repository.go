package repository

import (
	"e-wallet-go/internal/models"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WalletRepository interface {
	GetBalanceByUserID(userID string) (*models.Wallet, error)
	CreditBalance(userID string, amount float64, trxType, description string) error
	DebitBalance(userID string, amount float64, trxType, description string) error
	GetTransactionsByWalletID(walletID string, search string) ([]models.Transaction, error)
	GetAllTransactions(search string) ([]models.Transaction, error)
}

type walletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) WalletRepository {
	return &walletRepository{db: db}
}

func (r *walletRepository) GetBalanceByUserID(userID string) (*models.Wallet, error) {
	var wallet models.Wallet
	err := r.db.Where("user_id = ?", userID).First(&wallet).Error
	return &wallet, err
}

func (r *walletRepository) CreditBalance(userID string, amount float64, trxType, description string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var wallet models.Wallet

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ?", userID).First(&wallet).Error; err != nil {
			return err
		}

		wallet.Balance += amount
		if err := tx.Save(&wallet).Error; err != nil {
			return err
		}

		transaction := models.Transaction{
			WalletID:        wallet.ID,
			Amount:          amount,
			TransactionType: trxType,
			ReferenceID:     uuid.New().String(),
			Description:     description,
		}

		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *walletRepository) DebitBalance(userID string, amount float64, trxType, description string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var wallet models.Wallet

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ?", userID).First(&wallet).Error; err != nil {
			return err
		}

		if wallet.Balance < amount {
			return errors.New("insufficient balance")
		}

		wallet.Balance -= amount
		if err := tx.Save(&wallet).Error; err != nil {
			return err
		}

		transaction := models.Transaction{
			WalletID:        wallet.ID,
			Amount:          amount,
			TransactionType: trxType,
			ReferenceID:     uuid.New().String(),
			Description:     description,
		}

		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *walletRepository) GetTransactionsByWalletID(walletID string, search string) ([]models.Transaction, error) {
	var transactions []models.Transaction
	
	query := r.db.Where("wallet_id = ?", walletID)

	if search != "" {
		searchTerm := fmt.Sprintf("%%%s%%", search)
		query = query.Where(
			"description ILIKE ? OR transaction_type ILIKE ? OR reference_id ILIKE ? OR CAST(amount AS TEXT) ILIKE ?", 
			searchTerm, searchTerm, searchTerm, searchTerm,
		)
	}

	err := query.Preload("Wallet.User").Order("created_at desc").Find(&transactions).Error
	return transactions, err
}

func (r *walletRepository) GetAllTransactions(search string) ([]models.Transaction, error) {
	var transactions []models.Transaction

	query := r.db.Model(&models.Transaction{})

	if search != "" {
		searchTerm := fmt.Sprintf("%%%s%%", search)
		query = query.Where(
			"description ILIKE ? OR transaction_type ILIKE ? OR reference_id ILIKE ? OR CAST(amount AS TEXT) ILIKE ?", 
			searchTerm, searchTerm, searchTerm, searchTerm,
		)
	}

	err := query.Preload("Wallet.User").Order("created_at desc").Find(&transactions).Error
	return transactions, err
}