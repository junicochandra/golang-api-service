package messaging

type BrokerPort interface {
	Connect() error
	Close() error
	Publish(queue string, body []byte) error
	ConsumeOne(queue string) ([]byte, error)
}
