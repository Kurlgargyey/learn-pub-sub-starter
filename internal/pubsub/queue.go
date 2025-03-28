package pubsub

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	// durable means the queue will survive a broker restart
	Durable = iota
	// transient means the queue will be deleted when the last consumer unsubscribes
	Transient
)

func DeclareAndBind(conn *amqp.Connection, exchange, queueName, key string, simpleQueueType int,
) (*amqp.Channel, amqp.Queue, error) {
	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("failed to open a channel: %w", err)
	}
	// Declare the queue
	q, err := ch.QueueDeclare(
		queueName,
		simpleQueueType == Durable,   // durable
		simpleQueueType == Transient, // delete when unused
		simpleQueueType == Transient,
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("failed to declare a queue: %w", err)
	}
	// Bind the queue to the exchange
	err = ch.QueueBind(
		q.Name,
		key,
		exchange,
		false,
		nil,
	)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("failed to bind the queue: %w", err)
	}
	return ch, q, nil
}
