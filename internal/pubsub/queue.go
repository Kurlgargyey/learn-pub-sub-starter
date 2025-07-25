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

func DeclareAndBind(ch *amqp.Channel, exchange, queueName, key string, simpleQueueType int,
) (amqp.Queue, error) {
	// Declare the queue
	q, err := ch.QueueDeclare(
		queueName,
		simpleQueueType == Durable,   // durable
		simpleQueueType == Transient, // delete when unused
		simpleQueueType == Transient,
		false, // no-wait
		amqp.Table{
			"x-dead-letter-exchange": "peril_dlx",
		}, // arguments
	)
	if err != nil {
		return amqp.Queue{}, fmt.Errorf("failed to declare a queue: %w", err)
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
		return amqp.Queue{}, fmt.Errorf("failed to bind the queue: %w", err)
	}
	return q, nil
}
