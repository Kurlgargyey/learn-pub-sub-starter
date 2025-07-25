package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AckType int

const (
	Ack AckType = iota
	NackRequeue
	NackDiscard
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	body, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	ch.PublishWithContext(context.Background(),
		exchange,
		key,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	return nil
}

func SubscribeJSON[T any](
	ch *amqp.Channel,
	exchange,
	queueName,
	key string,
	simpleQueueType int, // an enum to represent "durable" or "transient"
	handler func(T) AckType,
) error {
	_, err := DeclareAndBind(ch, exchange, queueName, key, simpleQueueType)

	if err != nil {
		return fmt.Errorf("failed to declare and bind queue: %w", err)
	}
	cons, err := ch.Consume(
		queueName,
		"", false, false, false, false, nil,
	)
	if err != nil {
		return fmt.Errorf("failed to consume messages: %w", err)
	}
	go func() {
		for msg := range cons {
			var val T
			if err := json.Unmarshal(msg.Body, &val); err != nil {
				fmt.Printf("failed to unmarshal JSON: %s\n", err)
				continue
			}
			ack := handler(val)
			switch ack {
			case NackRequeue:
				err = msg.Nack(false, true)
			case NackDiscard:
				err = msg.Nack(false, false)
			default:
				err = msg.Ack(false)
			}
			if err != nil {
				fmt.Printf("failed to ack message: %s\n", err)
			}
		}
	}()
	return nil
}
