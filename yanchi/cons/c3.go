package main

import "go_rabbitmq/yanchi"

// 死信交换机和死信队列设置，死信队列消费
func main() {
	r := yanchi.NewRabbitMQDlx("yanchi_exchange_normal")
	r.CosumeDLX()
}
