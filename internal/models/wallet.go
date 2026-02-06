package models

import "github.com/google/uuid"

type Wallet struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID       uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	WalletNumber string    `gorm:"unique;not null" json:"wallet_number"`
	Balance      float64   `gorm:"type:decimal(15,2);default:0" json:"balance"`
	User         *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (Wallet) TableName() string {
	return "wallets"
}
