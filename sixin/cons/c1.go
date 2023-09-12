package main

import "go_rabbitmq/sixin"

func main() {
	r := sixin.NewRabbitMQDlx("exchange_1", "789")
	r.CosumeNormal()
}
