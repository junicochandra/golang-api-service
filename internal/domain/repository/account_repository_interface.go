package repository

import "github.com/junicochandra/golang-api-service/internal/domain/entity"

type AccountRepository interface {
	GetByAccountNumber(accountNumber string) (*entity.Account, error)
	UpdateBalanceTx(account *entity.Account) error
}
