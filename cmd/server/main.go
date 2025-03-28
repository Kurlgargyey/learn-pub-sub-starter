package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	const connectionString = "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(connectionString)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s\n", err)
		return
	}
	defer conn.Close()
	fmt.Println("Connected to RabbitMQ")
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s\n", err)
		return
	}
	defer ch.Close()
	fmt.Println("Channel opened")
	gamelogic.PrintServerHelp()

	exchange := routing.ExchangePerilDirect
	key := routing.PauseKey

	for {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}
		switch input[0] {
		case "pause":
			val := routing.PlayingState{
					IsPaused: true,
				}

			pubsub.PublishJSON(ch, exchange, key, val)
			fmt.Printf("Published message to exchange %s with key %s: %+v\n", exchange, key, val)

		case "resume":
			val := routing.PlayingState{
					IsPaused: false,
				}

			pubsub.PublishJSON(ch, exchange, key, val)
			fmt.Printf("Published message to exchange %s with key %s: %+v\n", exchange, key, val)

		case "quit":
			fmt.Println("Quitting...")
			return

		default: fmt.Println("Invalid command. Type 'help' for a list of commands.")
		}
	}
}
