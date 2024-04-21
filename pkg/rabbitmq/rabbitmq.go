package rabbitmq

import (
	"crypto/tls"
	"github.com/streadway/amqp"
)

type IBroker struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

var broker *IBroker

func Init(dsn string) (*IBroker, error) {
	cfg := new(tls.Config)
	cfg.InsecureSkipVerify = true

	conn, err := amqp.DialTLS(dsn, cfg)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	broker = &IBroker{
		Connection: conn,
		Channel:    ch,
	}

	return broker, nil
}

func (b *IBroker) InitQueue(queue string) error {
	_, err := b.Channel.QueueDeclare(
		queue, // queue name
		true,  // durable
		false, // auto delete
		false, // exclusive
		false, // no wait
		nil,   // arguments
	)

	return err
}

func (b *IBroker) ConnQueue(queue string) (<-chan amqp.Delivery, error) {
	messages, err := b.Channel.Consume(
		queue, // queue name
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (b *IBroker) SendToQueue(queue string, message amqp.Publishing) error {
	err := b.Channel.Publish(
		"",
		queue,
		false,
		false,
		message,
	)

	return err
}

func Get() *IBroker {
	return broker
}
