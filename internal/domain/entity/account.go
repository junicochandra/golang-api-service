package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type Account struct {
	ID            uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        uint64          `gorm:"not null" json:"userId"`
	AccountNumber string          `gorm:"not null" json:"accountNumber"`
	Balance       decimal.Decimal `gorm:"type:decimal(18,2);not null;default:0.00" json:"balance"`
	Currency      string          `gorm:"size:10;not null;default:'IDR'" json:"currency"`
	UpdatedAt     time.Time       `json:"updatedAt"`
}
