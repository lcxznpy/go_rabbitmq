package main

import (
	"fmt"
	"go_rabbitmq/sixin"
	"strconv"
	"time"
)

func main() {
	r := sixin.NewRabbitMQDlx("exchange_1", "123")
	for i := 0; i < 20; i++ {
		r.Publish(strconv.Itoa(i) + "r生产的消息")
		time.Sleep(time.Second * 1)
		fmt.Println(i)
	}
}
