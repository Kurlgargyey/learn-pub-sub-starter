package pubsub

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishGob[T any](ch *amqp.Channel, exchange, key string, val T) error {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(val)
	if err != nil {
		return fmt.Errorf("failed to encode Gob: %w", err)
	}
	ch.PublishWithContext(context.Background(),
		exchange,
		key,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/gob",
			Body:        buffer.Bytes(),
		})
	return nil
}

func SubscribeGob[T any](
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
			decoder := gob.NewDecoder(bytes.NewReader(msg.Body))
			if err := decoder.Decode(&val); err != nil {
				fmt.Printf("failed to decode Gob: %s\n", err)
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
