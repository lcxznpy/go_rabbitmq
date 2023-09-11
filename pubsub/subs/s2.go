package main

import (
	"go_rabbitmq/pubsub"
)

func main() {
	rabbitMQPubSub := pubsub.NewRabbitMQPubSub("qwerqwer")
	rabbitMQPubSub.ConsumeSub()
}
