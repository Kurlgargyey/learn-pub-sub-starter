package main

import (
	"fmt"
	"log"
	os "os"
	"os/signal"

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

	for {
		input := gamelogic.GetInput()

	}

	exchange := routing.ExchangePerilDirect
	key := routing.PauseKey
	val := routing.PlayingState{
		IsPaused: true,
	}

	pubsub.PublishJSON(ch, exchange, key, val)
	fmt.Printf("Published message to exchange %s with key %s: %+v\n", exchange, key, val)
	// wait for the message to be sent

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("Shutting down...")
}
