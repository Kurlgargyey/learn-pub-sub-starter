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

	usr, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("Failed to get user: %s\n", err)
		return
	}

	if err != nil {
		log.Fatalf("Failed to declare and bind: %s\n", err)
		return
	}

	state := gamelogic.NewGameState(usr)
	_, err = pubsub.SubscribeJSON(conn, "peril_topic", "army_moves."+usr, "army_moves.*", pubsub.Transient, handlerMove(state))
	ch, err := pubsub.SubscribeJSON(conn, routing.ExchangePerilDirect, routing.PauseKey+"."+usr, routing.PauseKey, pubsub.Transient, handlerPause(state))
	if err != nil {
		log.Fatalf("Failed to subscribe to pause: %s\n", err)
		return
	}
	defer ch.Close()

	client_repl(ch, routing.ExchangePerilDirect, routing.PauseKey+"."+usr, state)
}
