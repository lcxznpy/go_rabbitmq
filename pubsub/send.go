package pubsub

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

// 连接信息amqp://用户名:密码@ip/Virtual Hosts
const MQURL = "amqp://dhxdl666:dhxdl666@127.0.0.1:5672//dhxdl666"

// Rabbit RabbitMQ结构体
type Rabbit struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	QueueName string // 队列名称
	Exchange  string // 交换机名称
	Key       string // bind Key 名称
	MqUrl     string // 连接信息
}

// NewRabbitMQ 创建Rabbit结构体实例
func NewRabbitMQ(queueName, exchange, key string) *Rabbit {
	return &Rabbit{
		QueueName: queueName,
		Exchange:  exchange,
		Key:       key,
		MqUrl:     MQURL,
	}
}

// Destroy 断开channel和connection
func (r Rabbit) Destroy() error {
	err := r.channel.Close()
	err = r.conn.Close()
	return err
}

// 错误处理函数
func (r Rabbit) failOnErr(err error, msg string) {
	if err != nil {
		log.Fatal(msg, err)
	}
}

// NewRabbitMQPubSub 创建Publish/Subscribe模式下RabbitMQ实例
func NewRabbitMQPubSub(exchangeName string) *Rabbit {
	rabbitMQ := NewRabbitMQ("", exchangeName, "")
	var err error
	rabbitMQ.conn, err = amqp.Dial(rabbitMQ.MqUrl) // 获取connection
	rabbitMQ.failOnErr(err, "failed to connect rabbitmq!")
	rabbitMQ.channel, err = rabbitMQ.conn.Channel() // 获取channel
	rabbitMQ.failOnErr(err, "failed to open a channel")
	return rabbitMQ
}

// PublishPub Publish/Subscribe模式 生产者
func (r Rabbit) PublishPub(msg string) {
	// 1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange, // 交换机名字
		"fanout",   // 交换机类型，这里使用fanout类型，即: 发布订阅模式
		true,       // 是否持久化
		false,      // 是否自动删除
		false,      // true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,      // 是否阻塞处理
		nil,        // 额外的属性
	)
	r.failOnErr(err, "Failed to declare an exchange")
	// 2.发送消息
	err = r.channel.Publish(
		r.Exchange,
		"",
		false, // 如果为true，根据自身exchange类型和routekey规则无法找到符合条件的队列会把消息返还给发送者
		false, // 如果为true，当exchange发送消息到队列后发现队列上没有消费者，则会把消息返还给发送者
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		},
	)
	if err != nil {
		log.Println(err)
	}
}

// ConsumeSub  Publish/Subscribe模式 消费者
func (r Rabbit) ConsumeSub() {
	// 1.试探性创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange, // 交换机名字
		"fanout",   // 交换机类型，这里使用fanout类型，即: 发布订阅模式
		true,       // 是否持久化
		false,      // 是否自动删除
		false,      // true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,      // 是否阻塞处理
		nil,        // 额外的属性
	)
	r.failOnErr(err, "Failed to declare an exchange")
	// 2.试探性创建队列，这里注意队列名称不要写
	queue, err := r.channel.QueueDeclare(
		"",    // 随机生产队列名称
		false, // 是否持久化
		false, // 是否自动删除
		true,  // 是否具有排他性
		false, // 是否阻塞处理
		nil,   // 额外的属性
	)
	r.failOnErr(err, "Failed to declare a queue")
	// 3.绑定队列到exchange中
	err = r.channel.QueueBind(
		queue.Name, // 队列名
		"",         // 路由参数，fanout类型交换机，自动忽略路由参数（在pub/sub模式下，这里的key要为空）
		r.Exchange, // 交换机名字，需要跟消息发送端定义的交换器保持一致
		false,      // 是否阻塞处理
		nil,        // 额外的属性
	)
	// 4.消费消息
	msgs, err := r.channel.Consume(
		queue.Name, // 队列名称
		"",         // 用来区分多个消费者
		true,       // 是否自动应答
		false,      // 是否独有
		false,      // 设置为true，表示不能将同一个Connection中生产者发送的消息传递给这个Connection中的消费者
		false,      // 队列是否阻塞
		nil,        // 额外的属性
	)
	r.failOnErr(err, "Failed to Consume")
	// 5.启用协程处理消息
	forever := make(chan bool) // 开个channel阻塞住，让开启的协程能一直跑着
	go func() {
		for delivery := range msgs {
			// 消息逻辑处理，可以自行设计逻辑
			fmt.Println("Received a message:", string(delivery.Body))
		}
	}()
	fmt.Println(" [*] Waiting for messages.")
	<-forever
}
