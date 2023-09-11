package main

import "go_rabbitmq/work"

func main() {
	rabbitMQWork := work.NewRabbitMQWork("test_durable")
	rabbitMQWork.ConsumeWork()
}
