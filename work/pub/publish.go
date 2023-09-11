package main

import (
	"fmt"
	"go_rabbitmq/work"
	"strconv"
	"time"
)

func main() {
	rabbitMQWork := work.NewRabbitMQWork("test_durable")
	for i := 0; i < 20; i++ {
		//if i%2 == 0 {
		//	rabbitMQWork.PublishWork(strconv.Itoa(i) + " : hello,world!.")
		//} else {
		//	rabbitMQWork.PublishWork(strconv.Itoa(i) + " : hello,world!.....")
		//}
		rabbitMQWork.PublishWork(strconv.Itoa(i) + " : hello,world!.")
		time.Sleep(time.Second * 1)
		fmt.Println(i)
	}
}
