package main

import "go_rabbitmq/routing"

func main() {
	rabbitMQRouting1 := routing.NewRabbitMQRouting("exchange2", "dhxdl2")
	rabbitMQRouting1.ConsumeRouting()
}
