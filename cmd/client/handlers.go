package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishLog(ch *amqp.Channel, gs *gamelogic.GameState, msg string) error {
	err := pubsub.PublishGob(ch, routing.ExchangePerilTopic, routing.GameLogSlug+"."+gs.GetUsername(), msg)
	return err
}

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) pubsub.AckType {
	return func(ps routing.PlayingState) pubsub.AckType {
		defer fmt.Print("> ")
		gs.HandlePause(ps)
		return pubsub.Ack
	}
}

func handlerMove(ch *amqp.Channel, gs *gamelogic.GameState) func(gamelogic.ArmyMove) pubsub.AckType {
	return func(mov gamelogic.ArmyMove) pubsub.AckType {
		defer fmt.Print("> ")
		oc := gs.HandleMove(mov)
		switch oc {
		case gamelogic.MoveOutComeSafe:
			return pubsub.Ack
		case gamelogic.MoveOutcomeMakeWar:
			err := pubsub.PublishJSON(
				ch,
				"peril_topic",
				routing.WarRecognitionsPrefix+"."+gs.GetUsername(), gamelogic.RecognitionOfWar{
					Attacker: mov.Player,
					Defender: gs.GetPlayerSnap(),
				})
			if err != nil {
				fmt.Printf("Error publishing war recognition: %s\n", err)
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		default:
			return pubsub.NackDiscard
		}
	}
}

func handlerWar(ch *amqp.Channel, gs *gamelogic.GameState) func(gamelogic.RecognitionOfWar) pubsub.AckType {
	return func(row gamelogic.RecognitionOfWar) pubsub.AckType {
		defer fmt.Print("> ")
		oc, winner, loser := gs.HandleWar(row)
		logmsg := fmt.Sprintf("%s won a war against %s", winner, loser)
		switch oc {
		case gamelogic.WarOutcomeNotInvolved:
			return pubsub.NackRequeue
		case gamelogic.WarOutcomeNoUnits:
			return pubsub.NackDiscard
		case gamelogic.WarOutcomeYouWon:
			err := PublishLog(ch, gs, logmsg)
			if err != nil {
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		case gamelogic.WarOutcomeOpponentWon:
			err := PublishLog(ch, gs, logmsg)
			if err != nil {
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		case gamelogic.WarOutcomeDraw:
			logmsg = fmt.Sprintf("A war between %s and %s resulted in a draw", winner, loser)
			err := PublishLog(ch, gs, logmsg)
			if err != nil {
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		default:
			fmt.Printf("Error: Unknown war outcome: %v\n", oc)
			return pubsub.NackDiscard
		}
	}
}
