package payment

import (
	"errors"
	"fmt"
	"time"

	"github.com/junicochandra/golang-api-service/internal/app/payment/dto"
	"github.com/junicochandra/golang-api-service/internal/domain/repository"
	"github.com/shopspring/decimal"
)

var (
	ErrNotFound = errors.New("User not found")
)

type topUpUseCase struct {
	accountRepo repository.AccountRepository
}

func NewTopUpUseCase(accountRepo repository.AccountRepository) TopUpUseCase {
	return &topUpUseCase{accountRepo: accountRepo}
}

func (u *topUpUseCase) CreateTopUp(req *dto.TopUpRequest) (*dto.TopUpResponse, error) {
	account, err := u.accountRepo.GetByAccountNumber(req.AccountNumber)
	if err != nil {
		return nil, err
	}

	if account == nil {
		return nil, ErrNotFound
	}

	var balanceBefore = account.Balance
	account.Balance = account.Balance.Add(decimal.NewFromInt(req.Amount))
	account.UpdatedAt = time.Now()

	if err := u.accountRepo.UpdateBalanceTx(account); err != nil {
		return nil, fmt.Errorf("failed to update balance: %w", err)
	}

	return &dto.TopUpResponse{
		AccountNumber: req.AccountNumber,
		Amount:        decimal.NewFromInt(req.Amount),
		BalanceBefore: balanceBefore,
		BalanceAfter:  account.Balance,
		Currency:      "IDR",
		Status:        "pending",
	}, nil
}

// type Publisher interface {
// 	Publish(routingKey string, event interface{}) error
// }

// type TopUpUseCase struct {
// 	accountRepo     repo.AccountRepository
// 	transactionRepo repo.TransactionRepository
// 	publisher       Publisher
// 	db              *gorm.DB
// }

// func NewTopUpUseCase(accountRepo repo.AccountRepository, transactionRepo repo.TransactionRepository, publisher Publisher, db *gorm.DB) *TopUpUseCase {
// 	return &TopUpUseCase{
// 		accountRepo:     accountRepo,
// 		transactionRepo: transactionRepo,
// 		publisher:       publisher,
// 		db:              db,
// 	}
// }

// func (u *TopUpUseCase) CreateTopUp(accountID uint64, amount int64) (string, error) {
// 	txID := uuid.New().String()
// 	txn := &entity.Transaction{
// 		TransactionID: txID,
// 		TxType:        "TOPUP",
// 		AccountID:     accountID,
// 		Amount:        amount,
// 		Status:        "pending",
// 	}
// 	if err := u.transactionRepo.Create(txn); err != nil {
// 		return "", err
// 	}

// 	event := map[string]interface{}{
// 		"eventId": uuid.New().String(),
// 		"type":    "TopUpRequested",
// 		"payload": map[string]interface{}{
// 			"transactionId": txID,
// 			"accountId":     accountID,
// 			"amount":        amount,
// 		},
// 	}

// 	if err := u.publisher.Publish("topup", event); err != nil {
// 		// if publish fails, mark transaction as failed
// 		_ = u.transactionRepo.UpdateStatus(txID, "failed")
// 		return "", err
// 	}

// 	return txID, nil
// }
