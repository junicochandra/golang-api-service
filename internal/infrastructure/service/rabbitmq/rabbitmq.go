package rabbitmq

import (
	"context"
	"errors"
	"time"

	"github.com/junicochandra/golang-api-service/internal/app/messaging"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	url  string
	conn *amqp.Connection
	ch   *amqp.Channel
}

func New(url string) messaging.BrokerPort {
	return &RabbitMQ{url: url}
}

func (r *RabbitMQ) Connect() error {
	if r.conn != nil && !r.conn.IsClosed() {
		return nil
	}
	conn, err := amqp.Dial(r.url)
	if err != nil {
		return err
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return err
	}
	r.conn = conn
	r.ch = ch
	return nil
}

func (r *RabbitMQ) Close() error {
	var err error
	if r.ch != nil {
		if e := r.ch.Close(); e != nil && err == nil {
			err = e
		}
		r.ch = nil
	}
	if r.conn != nil {
		if e := r.conn.Close(); e != nil && err == nil {
			err = e
		}
		r.conn = nil
	}
	return err
}

func (r *RabbitMQ) ensureChannel() error {
	if r.ch == nil {
		return errors.New("channel is nil; call Connect() first")
	}
	return nil
}

func (r *RabbitMQ) Publish(queue string, body []byte) error {
	if err := r.ensureChannel(); err != nil {
		return err
	}

	_, err := r.ch.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return err
	}
	return r.ch.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)
}

func (r *RabbitMQ) ConsumeOne(queue string) ([]byte, error) {
	if err := r.ensureChannel(); err != nil {
		return nil, err
	}

	_, err := r.ch.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	msgs, err := r.ch.Consume(queue, "", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	select {
	case m := <-msgs:
		return m.Body, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
