package rabbitmq

import (
	"errors"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQService struct {
	url  string
	conn *amqp.Connection
}

func New(url string) (*RabbitMQService, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	return &RabbitMQService{
		url:  url,
		conn: conn,
	}, nil
}

// Channel returns a new channel from the connection
func (r *RabbitMQService) Channel() (*amqp.Channel, error) {
	if r.conn == nil {
		return nil, errors.New("rabbitmq: connection is nil")
	}
	return r.conn.Channel()
}

func (r *RabbitMQService) Publish(exchange, routingKey string, body []byte) error {
	ch, err := r.Channel()
	if err != nil {
		return err
	}
	// Close channel after publish
	defer ch.Close()

	// Make sure exchange exists (idempotent)
	if err := ch.ExchangeDeclare(
		exchange,
		"direct",
		true,  // durable
		false, // auto-deleted
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	return ch.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			Body:         body,
		},
	)
}

func (r *RabbitMQService) Close() {
	if r.conn != nil {
		_ = r.conn.Close()
	}
}
