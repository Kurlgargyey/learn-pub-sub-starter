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

	usr, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("Failed to get user: %s\n", err)
		return
	}

	state := gamelogic.NewGameState(usr)
	err = pubsub.SubscribeJSON(ch, routing.ExchangePerilTopic, "army_moves."+usr, "army_moves.*", pubsub.Transient, handlerMove(ch, state))
	if err != nil {
		log.Fatalf("Failed to subscribe to moves: %s\n", err)
		return
	}
	err = pubsub.SubscribeJSON(ch, routing.ExchangePerilDirect, routing.PauseKey+"."+usr, routing.PauseKey, pubsub.Transient, handlerPause(state))
	if err != nil {
		log.Fatalf("Failed to subscribe to pause: %s\n", err)
		return
	}

	err = pubsub.SubscribeJSON(ch, routing.ExchangePerilTopic, "war", routing.WarRecognitionsPrefix+".*", pubsub.Durable, handlerWar(ch, state))
	if err != nil {
		log.Fatalf("Failed to subscribe to war recognitions: %s\n", err)
		return
	}

	defer ch.Close()

	client_repl(ch, state)
}
