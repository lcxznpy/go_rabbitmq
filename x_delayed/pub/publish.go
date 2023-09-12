package main

import (
	"fmt"
	"go_rabbitmq/x_delayed"
	"strconv"
	"time"
)

func main() {
	r := x_delayed.NewRabbitMQXDelay("x_delayed_exchange", "x_delay_rkey")
	for i := 0; i < 20; i++ {
		r.PublishWithTTL(strconv.Itoa(i)+"r生产ttl=60s的消息", "60000")
		r.PublishWithTTL(strconv.Itoa(i)+"r生产ttl=5s的消息", "30000")
		time.Sleep(time.Second * 1)
		fmt.Println(i)
	}
}
