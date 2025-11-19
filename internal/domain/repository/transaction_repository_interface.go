package repository

import "github.com/junicochandra/golang-api-service/internal/domain/entity"

type TransactionRepository interface {
	Create(txn *entity.Transaction) error
	GetByTransactionID(transactionID string) (*entity.Transaction, error)
	UpdateStatus(transactionId string, status string) error
}
