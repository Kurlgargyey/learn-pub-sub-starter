package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	amqp "github.com/rabbitmq/amqp091-go"
)

func client_repl(ch *amqp.Channel, exchange, key string, state *gamelogic.GameState) {
	for {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}
		switch input[0] {
		case "spawn":
			err := state.CommandSpawn(input)
			if err != nil {
				fmt.Printf("Error spawning: %s\n", err)
				continue
			}

		case "move":
			_, err := state.CommandMove(input)
			if err != nil {
				fmt.Printf("Error moving: %s\n", err)
				continue
			}

		case "status":
			state.CommandStatus()

		case "spam":
			fmt.Println("Spamming not allowed yet.")

		case "help":
			gamelogic.PrintClientHelp()
		case "quit":
			fmt.Println("Quitting...")
			return
		default:
			fmt.Println("Unknown command. Type 'help' for a list of commands.")
		}
	}
}