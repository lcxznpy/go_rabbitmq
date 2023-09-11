package main

import "go_rabbitmq/simple"

func main() {
	rabbitMQSimple := simple.NewRabbitMQSimple("test1")
	rabbitMQSimple.ConsumeSimple()
}
