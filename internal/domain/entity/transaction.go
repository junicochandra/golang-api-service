package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type Transaction struct {
	ID                int64           `json:"id" db:"id"`
	TransactionID     string          `json:"transactionId" db:"transaction_id"`
	Type              string          `json:"type" db:"type"`                             // transfer | topup | payment
	SenderAccountID   string          `json:"senderAccountId" db:"sender_account_id"`     // nullable
	ReceiverAccountID string          `json:"receiverAccountId" db:"receiver_account_id"` // nullable
	Amount            decimal.Decimal `json:"amount" db:"amount"`
	Status            string          `json:"status" db:"status"`           // pending | processing | success | failed
	Reference         *string         `json:"reference" db:"reference"`     // nullable
	Description       *string         `json:"description" db:"description"` // nullable
	Payload           *string         `json:"payload" db:"payload"`         // nullable (text)
	CreatedAt         time.Time       `json:"createdAt" db:"created_at"`
	UpdatedAt         time.Time       `json:"updatedAt" db:"updated_at"`
}
