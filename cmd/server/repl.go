package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func server_repl(ch *amqp.Channel, exchange, key string) {
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

		case "help":
			gamelogic.PrintServerHelp()

		case "quit":
			fmt.Println("Quitting...")
			return

		default:
			fmt.Println("Invalid command. Type 'help' for a list of commands.")
		}
	}
}
