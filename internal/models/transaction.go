package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	TRX_DEPOSIT         = "DEPOSIT"
	TRX_WITHDRAW        = "WITHDRAW"
	TRX_ADMIN_TOPUP     = "ADMIN_TOPUP"
	TRX_ADMIN_DEDUCTION = "ADMIN_DEDUCTION"
)

type Transaction struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	WalletID        uuid.UUID `gorm:"type:uuid;not null" json:"wallet_id"`
	Amount          float64   `gorm:"type:decimal(15,2);not null" json:"amount"`
	TransactionType string    `gorm:"not null" json:"transaction_type"`
	ReferenceID     string    `gorm:"unique;not null" json:"reference_id"`
	Description     string    `json:"description"`
	CreatedAt       time.Time `json:"created_at"`
	Wallet          *Wallet   `gorm:"foreignKey:WalletID" json:"wallet,omitempty"`
}

func (Transaction) TableName() string {
	return "transactions"
}
