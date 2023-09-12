package main

import "go_rabbitmq/yanchi"

// 设置可变的ttl优化固定时间的延迟队列
func main() {
	r := yanchi.NewRabbitMQDlx("yanchi_exchange_normal")
	r.CosumeNormal3()
}
