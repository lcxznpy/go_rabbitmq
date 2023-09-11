package main

import (
	"fmt"
	"go_rabbitmq/pubsub"
	"strconv"
	"time"
)

func main() {
	rabbitMQPubSub := pubsub.NewRabbitMQPubSub("qwerqwer")
	for i := 0; i < 100; i++ {
		rabbitMQPubSub.PublishPub("订阅模式生产的第" + strconv.Itoa(i) + "条数据")
		fmt.Println("订阅模式生产第" + strconv.Itoa(i) + "条数据")
		time.Sleep(time.Second * 1)
	}
}
