package main

import (
	"fmt"
	"go_rabbitmq/topics"
	"strconv"
	"time"
)

func main() {
	rabbitMQTopics1 := topics.NewRabbitMQTopics("exchangeTopics", "dhxdl.1011.top")
	rabbitMQTopics2 := topics.NewRabbitMQTopics("exchangeTopics", "dhxdl.xyz.top")
	for i := 0; i < 100; i++ {
		rabbitMQTopics1.PublishTopics(strconv.Itoa(i) + "rabbitMQTopics1生产的消息")
		rabbitMQTopics2.PublishTopics(strconv.Itoa(i) + "rabbitMQTopics2生产的消息")
		time.Sleep(time.Second * 1)
		fmt.Println(i)
	}
}
