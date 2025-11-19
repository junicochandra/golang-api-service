package worker

import (
	"encoding/json"
	"log"

	"github.com/junicochandra/golang-api-service/internal/domain/entity"
	repo "github.com/junicochandra/golang-api-service/internal/domain/repository"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TopUpEvent struct {
	EventId string `json:"eventId"`
	Type    string `json:"type"`
	Payload struct {
		TransactionId string `json:"transactionId"`
		AccountId     uint64 `json:"accountId"`
		Amount        int64  `json:"amount"`
	} `json:"payload"`
}

type Consumer struct {
	db              *gorm.DB
	conn            *amqp.Connection
	channel         *amqp.Channel
	transactionRepo repo.TransactionRepository
	accountRepo     repo.AccountRepository
}

func NewConsumer(db *gorm.DB, amqURL string, conn *amqp.Connection, channel *amqp.Channel, transactionRepo repo.TransactionRepository, accountRepo repo.AccountRepository) (*Consumer, error) {
	conn, err := amqp.Dial(amqURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	queueName := "topup_queue"
	exchange := "payments"
	ch.QueueDeclare(queueName, true, false, false, false, nil)
	ch.QueueBind(queueName, "topup", exchange, false, nil)

	return &Consumer{
		db:              db,
		conn:            conn,
		channel:         ch,
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
	}, nil
}

func (c *Consumer) Start(queueName string) error {
	msgs, err := c.channel.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}
	go func() {
		for d := range msgs {
			var ev TopUpEvent
			if err := json.Unmarshal(d.Body, &ev); err != nil {
				log.Println("invalid event:", err)
				d.Nack(false, false)
				continue
			}
			if err := c.handleTopUp(ev); err != nil {
				log.Println("handleTopUp error:", err)
				// requeue true (retry). In prod: use DLQ & retry limits.
				d.Nack(false, true)
			} else {
				d.Ack(false)
			}
		}
	}()
	return nil
}

func (c *Consumer) handleTopUp(ev TopUpEvent) error {
	// cek idempotency
	t, err := c.transactionRepo.GetByTransactionID(ev.Payload.TransactionId)
	if err != nil {
		return err
	}
	if t.Status == "completed" {
		return nil
	}

	// mulai transaction GORM
	return c.db.Transaction(func(tx *gorm.DB) error {
		// lock row FOR UPDATE menggunakan clause Locking
		var acc entity.Account
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", ev.Payload.AccountId).First(&acc).Error; err != nil {
			return err
		}

		newBalance := acc.Balance + ev.Payload.Amount
		if err := c.accountRepo.UpdateBalanceTx(tx, ev.Payload.AccountId, newBalance); err != nil {
			return err
		}

		// update transaction status via repository but repo uses the main db;
		// gunakan tx.Session(&gorm.Session{NewDB: true}) atau implement UpdateStatusTx optionally.
		// Simpler: update status using tx directly here:
		if err := tx.Model(&entity.Transaction{}).Where("transaction_id = ?", ev.Payload.TransactionId).Update("status", "completed").Error; err != nil {
			return err
		}

		return nil
	})
}
