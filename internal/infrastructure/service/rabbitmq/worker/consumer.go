package worker

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/junicochandra/golang-api-service/internal/domain/repository"
	"github.com/junicochandra/golang-api-service/internal/infrastructure/service/rabbitmq"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/shopspring/decimal"
)

// TopUpMessage according to the payload sent from the usecase
type TopUpMessage struct {
	TransactionID string          `json:"transactionId"`
	AccountNumber string          `json:"accountNumber"`
	Amount        decimal.Decimal `json:"amount"`
	Currency      string          `json:"currency"`
	CreatedAt     time.Time       `json:"createdAt"`
}

type Consumer struct {
	rabbit          *rabbitmq.RabbitMQService
	transactionRepo repository.TransactionRepository
	accountRepo     repository.AccountRepository
	queueName       string
	logger          *log.Logger
}

func NewConsumer(r *rabbitmq.RabbitMQService, trx repository.TransactionRepository, acc repository.AccountRepository, queueName string, logger *log.Logger) *Consumer {
	return &Consumer{
		rabbit:          r,
		transactionRepo: trx,
		accountRepo:     acc,
		queueName:       queueName,
		logger:          logger,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	ch, err := c.rabbit.Channel()
	if err != nil {
		return err
	}

	// Close channel when context done
	go func() {
		<-ctx.Done()
		_ = ch.Cancel("", false)
		_ = ch.Close()
	}()

	// Set QoS
	if err := ch.Qos(1, 0, false); err != nil {
		return err
	}

	msgs, err := ch.Consume(
		c.queueName,
		"",    // consumer tag
		false, // autoAck
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	c.logger.Println("worker: waiting for messages...")

	for {
		select {
		case <-ctx.Done():
			c.logger.Println("worker: context done, stopping")
			return nil
		case d, ok := <-msgs:
			if !ok {
				c.logger.Println("worker: deliveries closed")
				return nil
			}
			if err := c.handleDelivery(d); err != nil {
				c.logger.Printf("worker: handleDelivery error: %v", err)
			}
		}
	}
}

func (c *Consumer) handleDelivery(d amqp.Delivery) error {
	var m TopUpMessage
	if err := json.Unmarshal(d.Body, &m); err != nil {
		c.logger.Printf("worker: invalid message: %v", err)
		_ = d.Reject(false) // send to DLX if configured
		return err
	}

	// Idempotency check using trx repo
	trx, err := c.transactionRepo.GetByTransactionID(m.TransactionID)
	if err != nil {
		c.logger.Printf("worker: db error GetByID: %v", err)
		_ = d.Nack(false, true) // requeue
		return err
	}
	if trx == nil {
		c.logger.Printf("worker: transaction not found: %s", m.TransactionID)
		_ = d.Reject(false)
		return errors.New("transaction not found")
	}
	if trx.Status == "completed" || trx.Status == "success" {
		_ = d.Ack(false)
		return nil
	}
	if trx.Status == "processing" {
		_ = d.Nack(false, true)
		return nil
	}

	// Set processing
	if err := c.transactionRepo.UpdateStatus(m.TransactionID, "processing"); err != nil {
		c.logger.Printf("worker: failed set processing: %v", err)
		_ = d.Nack(false, true)
		return err
	}

	// Get account & update balance atomically via repository method
	account, err := c.accountRepo.GetByAccountNumber(m.AccountNumber)
	if err != nil {
		c.logger.Printf("worker: get account error: %v", err)
		_ = c.transactionRepo.UpdateStatus(m.TransactionID, "failed_account_error")
		_ = d.Nack(false, true)
		return err
	}
	if account == nil {
		c.logger.Printf("worker: account not found: %s", m.AccountNumber)
		_ = c.transactionRepo.UpdateStatus(m.TransactionID, "failed_account_not_found")
		_ = d.Ack(false)
		return errors.New("account not found")
	}

	account.Balance = account.Balance.Add(m.Amount)
	account.UpdatedAt = time.Now()

	if err := c.accountRepo.UpdateBalanceTx(account); err != nil {
		c.logger.Printf("worker: UpdateBalanceTx error: %v", err)
		_ = c.transactionRepo.UpdateStatus(m.TransactionID, "failed_update_balance")
		_ = d.Nack(false, true)
		return err
	}

	if err := c.transactionRepo.UpdateStatus(m.TransactionID, "completed"); err != nil {
		c.logger.Printf("worker: warning: failed to mark trx completed: %v", err)
		_ = d.Ack(false)
		return err
	}

	_ = d.Ack(false)
	c.logger.Printf("worker: processed tx=%s acc=%s amount=%s", m.TransactionID, m.AccountNumber, m.Amount.String())
	return nil
}
