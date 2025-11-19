package dto

import "github.com/shopspring/decimal"

type TopUpRequest struct {
	AccountNumber string `json:"accountNumber"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
}

type TopUpResponse struct {
	TransactionID uint64          `json:"transactionId"`
	AccountNumber string          `json:"accountNumber"`
	Amount        decimal.Decimal `json:"amount"`
	BalanceBefore decimal.Decimal `json:"balanceBefore"`
	BalanceAfter  decimal.Decimal `json:"balanceAfter"`
	Currency      string          `json:"currency"`
	Status        string          `json:"status"`
	Message       string          `json:"message"`
}
