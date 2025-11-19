package service

import (
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitPublisher struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
}

func NewRabbitPublisher(url, exchange string) (*RabbitPublisher, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, err
	}

	if err := ch.ExchangeDeclare(
		exchange, // name
		"direct", // kind
		true,     // durable
		false,    // autoDelete
		false,    // internal
		false,    // noWait
		nil,      // args
	); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, err
	}

	return &RabbitPublisher{conn: conn, channel: ch, exchange: exchange}, nil
}

func (p *RabbitPublisher) Publish(routingKey string, event interface{}) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.channel.Publish(
		p.exchange, // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate (deprecated in server but param still exists)
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (p *RabbitPublisher) Close() error {
	var firstErr error
	if p.channel != nil {
		if err := p.channel.Close(); err != nil {
			firstErr = err
		}
	}
	if p.conn != nil {
		if err := p.conn.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
