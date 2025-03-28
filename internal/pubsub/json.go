package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
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
    conn *amqp.Connection,
    exchange,
    queueName,
    key string,
    simpleQueueType int, // an enum to represent "durable" or "transient"
    handler func(T),
) (*amqp.Channel, error) {
	ch, _, err := DeclareAndBind(conn, exchange, queueName, key, simpleQueueType)

	if err != nil {
		return nil, fmt.Errorf("failed to declare and bind queue: %w", err)
	}
	cons, err := ch.Consume(
		queueName,
		"",false,false,false,false,nil,
		)
	if err != nil {
		return nil, fmt.Errorf("failed to consume messages: %w", err)
	}
	go func() {
		for msg := range cons {
			var val T
			if err := json.Unmarshal(msg.Body, &val); err != nil {
				fmt.Printf("failed to unmarshal JSON: %s\n", err)
				continue
			}
			handler(val)
			if err := msg.Ack(false); err != nil {
				fmt.Printf("failed to ack message: %s\n", err)
			}
		}
	}()
	return ch, nil
	}