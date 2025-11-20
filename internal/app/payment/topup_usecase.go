package payment

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/junicochandra/golang-api-service/internal/app/payment/dto"
	"github.com/junicochandra/golang-api-service/internal/domain/entity"
	"github.com/junicochandra/golang-api-service/internal/domain/repository"
	"github.com/shopspring/decimal"
)

var (
	ErrNotFound = errors.New("User not found")
)

type topUpUseCase struct {
	accountRepo     repository.AccountRepository
	transactionRepo repository.TransactionRepository
}

func NewTopUpUseCase(accountRepo repository.AccountRepository, transactionRepo repository.TransactionRepository) TopUpUseCase {
	return &topUpUseCase{
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
	}
}

func (u *topUpUseCase) CreateTopUp(req *dto.TopUpRequest) (*dto.TopUpResponse, error) {
	// Validate request
	if req == nil {
		return nil, fmt.Errorf("request is nil")
	}
	if req.AccountNumber == "" {
		return nil, fmt.Errorf("account number is required")
	}

	amountDecimal := decimal.NewFromFloat(float64(req.Amount))
	if amountDecimal.Cmp(decimal.Zero) <= 0 {
		return nil, fmt.Errorf("amount must be greater than 0")
	}

	// Check account existence
	account, err := u.accountRepo.GetByAccountNumber(req.AccountNumber)
	if err != nil {
		return nil, err
	}

	if account == nil {
		return nil, ErrNotFound
	}

	// Create Transaction
	txID := uuid.New().String()
	txn := &entity.Transaction{
		TransactionID:     txID,
		Type:              "topup",
		SenderAccountID:   req.AccountNumber,
		ReceiverAccountID: req.AccountNumber,
		Amount:            amountDecimal,
		Status:            "pending",
	}

	if err := u.transactionRepo.Create(txn); err != nil {
		return nil, err
	}

	// Update Balance
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
