package payment

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/junicochandra/golang-api-service/internal/app/payment/dto"
	"github.com/junicochandra/golang-api-service/internal/domain/entity"
	"github.com/junicochandra/golang-api-service/internal/domain/repository"
	"github.com/junicochandra/golang-api-service/internal/infrastructure/service/rabbitmq"
	"github.com/shopspring/decimal"
)

var (
	ErrNotFound = errors.New("User not found")
)

type TopUpMessage struct {
	TransactionID string          `json:"transactionId"`
	AccountNumber string          `json:"accountNumber"`
	Amount        decimal.Decimal `json:"amount"`
	Currency      string          `json:"currency"`
	CreatedAt     time.Time       `json:"createdAt"`
}

func (m *TopUpMessage) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

type topUpUseCase struct {
	accountRepo     repository.AccountRepository
	transactionRepo repository.TransactionRepository
	rabbitSvc       *rabbitmq.RabbitMQService
}

func NewTopUpUseCase(accountRepo repository.AccountRepository, transactionRepo repository.TransactionRepository, rabbitSvc *rabbitmq.RabbitMQService) TopUpUseCase {
	return &topUpUseCase{
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
		rabbitSvc:       rabbitSvc,
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

	// Create Transaction (pending)
	txID := uuid.New().String()
	txn := &entity.Transaction{
		TransactionID:     txID,
		Type:              "topup",
		SenderAccountID:   req.AccountNumber,
		ReceiverAccountID: req.AccountNumber,
		Amount:            amountDecimal,
		Status:            "pending",
		CreatedAt:         time.Now(),
	}

	if err := u.transactionRepo.Create(txn); err != nil {
		return nil, err
	}

	// Prepare message and publish
	msg := &TopUpMessage{
		TransactionID: txID,
		AccountNumber: req.AccountNumber,
		Amount:        amountDecimal,
		Currency:      "IDR",
		CreatedAt:     time.Now(),
	}
	body, err := msg.Marshal()
	if err != nil {
		_ = u.transactionRepo.UpdateStatus(txID, "failed_marshal")
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}

	if u.rabbitSvc == nil {
		_ = u.transactionRepo.UpdateStatus(txID, "failed_no_rabbit")
		return nil, fmt.Errorf("rabbitmq service not initialized")
	}

	if err := u.rabbitSvc.Publish("topup.exchange", "topup.created", body); err != nil {
		_ = u.transactionRepo.UpdateStatus(txID, "failed_publish")
		return nil, fmt.Errorf("failed to publish topup message: %w", err)
	}

	// Success: return pending response (balance not yet updated)
	return &dto.TopUpResponse{
		AccountNumber: req.AccountNumber,
		Amount:        amountDecimal,
		BalanceBefore: account.Balance,
		BalanceAfter:  account.Balance,
		Currency:      "IDR",
		Status:        "pending",
	}, nil
}
