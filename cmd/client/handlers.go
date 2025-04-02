package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) pubsub.AckType {
	return func(ps routing.PlayingState) pubsub.AckType {
		defer fmt.Print("> ")
		gs.HandlePause(ps)
		return pubsub.Ack
	}
}

func handlerMove(ch *amqp.Channel, gs *gamelogic.GameState) func(gamelogic.ArmyMove) pubsub.AckType{
	return func(mov gamelogic.ArmyMove) pubsub.AckType {
		defer fmt.Print("> ")
		oc := gs.HandleMove(mov)
		switch oc {
		case gamelogic.MoveOutComeSafe:
			return pubsub.Ack
		case gamelogic.MoveOutcomeMakeWar:
			pubsub.PublishJSON(
				ch,
				"peril_topic",
				routing.WarRecognitionsPrefix+"."+gs.GetUsername(),mov)
			return pubsub.NackRequeue
		default:
			return pubsub.NackDiscard
		}
	}
}