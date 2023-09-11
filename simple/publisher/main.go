package main

import (
	"fmt"
	"go_rabbitmq/simple"
)

func main() {
	rabbitMQsimple := simple.NewRabbitMQSimple("test1")
	rabbitMQsimple.PublishSimple("qwerqwer")
	fmt.Println("successssssss")
}
