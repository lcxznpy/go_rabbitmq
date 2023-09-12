package main

import "go_rabbitmq/yanchi"

// 普通交换机和固定时间20s延迟队列设置，队列消费
func main() {
	r := yanchi.NewRabbitMQDlx("yanchi_exchange_normal")
	r.CosumeNormal2()
}
