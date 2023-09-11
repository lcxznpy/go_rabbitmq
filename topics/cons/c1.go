package main

import "go_rabbitmq/topics"

func main() {
	rabbitMQTopics := topics.NewRabbitMQTopics("exchangeTopics", "#")
	rabbitMQTopics.ConsumeTopics()
}
