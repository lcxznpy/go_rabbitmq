package main

import (
	"fmt"
	"go_rabbitmq/fabuqueren"
	"strconv"
	"time"
)

func main() {
	rabbitMQWork := fabuqueren.NewRabbitMQWork("test_durable")
	start := time.Now()
	rabbitMQWork.SetMode()

	for i := 0; i < 1000; i++ {
		//if i%2 == 0 {
		//	rabbitMQWork.PublishWork(strconv.Itoa(i) + " : hello,world!.")
		//} else {
		//	rabbitMQWork.PublishWork(strconv.Itoa(i) + " : hello,world!.....")
		//}
		rabbitMQWork.PublishWork(strconv.Itoa(i) + " : hello,world!.")

		fmt.Println(i)
	}

	cost := time.Since(start)
	fmt.Printf("cost=%s", cost)
}
