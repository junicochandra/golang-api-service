package repository

import (
	"errors"

	"github.com/junicochandra/golang-api-service/internal/domain/entity"
	accountRepo "github.com/junicochandra/golang-api-service/internal/domain/repository"
	"github.com/junicochandra/golang-api-service/internal/infrastructure/config/database"
	"gorm.io/gorm"
)

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) accountRepo.AccountRepository {
	return &accountRepository{db: database.DB}
}

func (repo *accountRepository) GetByAccountNumber(accountNumber string) (*entity.Account, error) {
	var account entity.Account
	if err := repo.db.Where("account_number = ?", accountNumber).First(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}

func (repo *accountRepository) UpdateBalanceTx(account *entity.Account) error {
	return repo.db.Save(account).Error
}
