package main

import (
	"fmt"
	"go_rabbitmq/routing"
	"strconv"
	"time"
)

func main() {
	rabbitMQRouting1 := routing.NewRabbitMQRouting("exchange2", "dhxdl1")
	rabbitMQRouting2 := routing.NewRabbitMQRouting("exchange2", "dhxdl2")
	for i := 0; i < 20; i++ {
		rabbitMQRouting1.PublishRouting(strconv.Itoa(i) + " : Hello dhxdl1!")
		rabbitMQRouting2.PublishRouting(strconv.Itoa(i) + " : Hello dhxdl2!")
		time.Sleep(time.Second * 1)
		fmt.Println(i)
	}
}
