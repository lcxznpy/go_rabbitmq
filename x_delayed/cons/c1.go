package main

import "go_rabbitmq/x_delayed"

func main() {
	r := x_delayed.NewRabbitMQXDelay("x_delayed_exchange", "x_delay_rkey")
	r.ConsumeXDelay()
}
