package rabbitmq

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type TopologyConfig struct {
	Exchange   string
	ExchangeTy string
	Queue      string
	RoutingKey string
	// Optional: DLX string
	DLX string
}

func DeclareTopology(r *RabbitMQService, cfg TopologyConfig) error {
	ch, err := r.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// Declare exchange
	if err := ch.ExchangeDeclare(
		cfg.Exchange,
		cfg.ExchangeTy,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	// Declare queue (optionally set DLX)
	table := amqp.Table{}
	if cfg.DLX != "" {
		table["x-dead-letter-exchange"] = cfg.DLX
	}

	_, err = ch.QueueDeclare(
		cfg.Queue,
		true,
		false,
		false,
		false,
		table,
	)
	if err != nil {
		return err
	}

	// Bind queue to exchange
	if err := ch.QueueBind(
		cfg.Queue,
		cfg.RoutingKey,
		cfg.Exchange,
		false,
		nil,
	); err != nil {
		return err
	}

	log.Printf("rabbitmq: declared exchange=%s queue=%s routingKey=%s", cfg.Exchange, cfg.Queue, cfg.RoutingKey)
	return nil
}
