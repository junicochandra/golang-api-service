package repository

import (
	"errors"

	"github.com/junicochandra/golang-api-service/internal/domain/entity"
	transactionRepo "github.com/junicochandra/golang-api-service/internal/domain/repository"
	"github.com/junicochandra/golang-api-service/internal/infrastructure/config/database"
	"gorm.io/gorm"
)

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) transactionRepo.TransactionRepository {
	return &transactionRepository{db: database.DB}
}

func (repo *transactionRepository) Create(txn *entity.Transaction) error {
	return repo.db.Create(txn).Error
}

func (repo *transactionRepository) GetByTransactionID(transactionId string) (*entity.Transaction, error) {
	var txn entity.Transaction
	if err := repo.db.Where("transaction_id = ?", transactionId).First(&txn).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &txn, nil
}

func (repo *transactionRepository) UpdateStatus(transactionId string, status string) error {
	return repo.db.Model(&entity.Transaction{}).Where("transaction_id = ?", transactionId).Update("status", status).Error
}
