package main

import (
	"fmt"
	"go_rabbitmq/yanchi"
	"strconv"
	"time"
)

func main() {
	r := yanchi.NewRabbitMQDlx("yanchi_exchange_normal")
	for i := 0; i < 20; i++ {
		r.Publish(strconv.Itoa(i)+"r生产ttl=10s的消息", "qwer")
		//r.Publish(strconv.Itoa(i)+"r生产ttl=20s的消息", "asdf")
		r.PublishWithTTL(strconv.Itoa(i)+"r生产ttl=60s的消息", "zxcv", "60000") //主动设置延迟时间优化
		time.Sleep(time.Second * 1)
		fmt.Println(i)
	}
}
