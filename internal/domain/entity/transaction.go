package entity

import "time"

type Transaction struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	TransactionID string    `gorm:"size:64;uniqueIndex;not null" json:"transactionId"`
	TxType        string    `gorm:"size:32;not null" json:"txType"` // TOPUP, TRANSFER, PAYMENT
	AccountID     uint64    `json:"accountId"`
	Amount        int64     `gorm:"not null" json:"amount"`
	Status        string    `gorm:"size:16;not null;default:'pending'" json:"status"`
	Payload       string    `gorm:"type:json" json:"payload"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
