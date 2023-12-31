package sixin

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
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

// NewRabbitMQTopics 创建Topics模式下RabbitMQ实例
func NewRabbitMQDlx(exchangeName, routingKey string) *Rabbit {
	rabbitMQ := NewRabbitMQ("", exchangeName, routingKey) // 创建RabbitMQ实例
	var err error
	rabbitMQ.conn, err = amqp.Dial(rabbitMQ.MqUrl) // 获取connection
	rabbitMQ.failOnErr(err, "failed to connect rabbitmq!")
	rabbitMQ.channel, err = rabbitMQ.conn.Channel() // 获取channel
	rabbitMQ.failOnErr(err, "failed to open a channel")
	return rabbitMQ
}

// 生产者想普通交换机内发送消息
func (r Rabbit) Publish(msg string) {
	err := r.channel.ExchangeDeclare(
		r.Exchange, // 交换机名字
		"direct",   // 交换机类型，这里使用topic类型，即: Topics模式
		true,       // 是否持久化
		false,      // 是否自动删除
		false,      // true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,      // 是否阻塞处理
		nil,        // 额外的属性
	)
	r.failOnErr(err, "Failed to declare an exchange")
	//time.Sleep(time.Second * 2)
	err = r.channel.Publish(
		r.Exchange,
		r.Key, // Topics模式这里要指定key
		false, // 如果为true，根据自身exchange类型和routekey规则无法找到符合条件的队列会把消息返还给发送者
		false, // 如果为true，当exchange发送消息到队列后发现队列上没有消费者，则会把消息返还给发送者
		amqp.Publishing{
			Headers:      amqp.Table{},
			DeliveryMode: amqp.Persistent,
			Priority:     0,
			ContentType:  "text/plain",
			Body:         []byte(msg),
		},
	)
	if err != nil {
		log.Println(err)
	}
}

func (r Rabbit) CosumeNormal() {
	dlxExchangeName := "dlx_exchange" //死信交换机名称
	//1.声明正常交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange, // 交换机名字
		"direct",   // 交换机类型，这里使用topic类型，即: Topics模式
		true,       // 是否持久化
		false,      // 是否自动删除
		false,      // true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,      // 是否阻塞处理
		nil,        // 额外的属性
	)
	r.failOnErr(err, "Failed to declare an exchange")
	//2.设置队列参数
	argsQue := make(map[string]interface{})
	//添加死信队列交换机属性
	argsQue["x-dead-letter-exchange"] = dlxExchangeName
	//指定死信队列的路由key，不指定使用队列路由键
	//argsQue["x-dead-letter-routing-key"] = r.Key
	//添加过期时间
	//argsQue["x-message-ttl"] = 6000 //单位毫秒
	//队列最大长度
	//argsQue["x-max-length"] = 6

	//3.声明队列并添加参数
	queue, err := r.channel.QueueDeclare(
		"sixin_reject", // 随机生产队列名称
		true,           // 是否持久化
		false,          // 是否自动删除
		false,          // 是否具有排他性   ,排他性连接关闭自动删除队列
		false,          // 是否阻塞处理
		argsQue,        // 额外的属性
	)
	r.failOnErr(err, "Failed to declare a queue")

	//4.绑定队列和普通交换机
	err = r.channel.QueueBind(
		queue.Name, // 队列名
		r.Key,      // 路由参数，如果匹配消息发送的时候指定的路由参数，消息就投递到当前队列（在Topics模式下，这里的key要指定）
		r.Exchange, // 交换机名字，需要跟消息发送端定义的交换器保持一致
		false,      // 是否阻塞处理
		nil,        // 额外的属性
	)
	r.failOnErr(err, "QueueBind err:")

	// 5.消费消息
	msgs, err := r.channel.Consume(
		queue.Name, // 队列名称
		"",         // 用来区分多个消费者
		false,      // 是否自动应答
		false,      // 是否独有
		false,      // 设置为true，表示不能将同一个Connection中生产者发送的消息传递给这个Connection中的消费者
		false,      // 队列是否阻塞
		nil,        // 额外的属性
	)
	r.failOnErr(err, "Failed to Consume")

	fmt.Println(" [*] Waiting for messages.")
	//6.消费接收到的消息
	for {
		select {
		case message, ok := <-msgs:
			if !ok {
				continue
			}
			go func() {
				//处理消息
				time.Sleep(time.Second * 2)
				//确认接收到的消息
				if string(message.Body) == "5r生产的消息" {
					if err = message.Reject(false); err != nil {
						fmt.Println("d.reject err: ", err)
						return
					}
					fmt.Println("5r生产的消息 reject reject reject reject")
					return
				}
				if err = message.Ack(true); err != nil {
					//TODD: 获取到消息后，在过期时间内如果未进行确认，此消息就会流入到死信队列，此时进行消息确认就会报错
					fmt.Println("d.Ack err: ", err)
					return
				}
				fmt.Println("已确认", string(message.Body))
			}()
		case <-time.After(time.Second * 1):

		}
	}

}

func (r Rabbit) CosumeDLX() {
	//1.声明死信交换机
	err := r.channel.ExchangeDeclare(
		"dlx_exchange",      // 交换机名字
		amqp.ExchangeFanout, // 交换机类型，这里使用topic类型，即: Topics模式
		true,                // 是否持久化
		false,               // 是否自动删除
		false,               // true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,               // 是否阻塞处理
		nil,                 // 额外的属性
	)
	r.failOnErr(err, "Failed to declare an exchange")

	//2.声明死信队列
	queue, err := r.channel.QueueDeclare(
		"dlx_queue", // 随机生产队列名称
		true,        // 是否持久化
		false,       // 是否自动删除
		false,       // 是否具有排他性
		false,       // 是否阻塞处理
		nil,         // 额外的属性
	)
	r.failOnErr(err, "Failed to declare a queue")

	//3.死信队列和死信交换机绑定
	err = r.channel.QueueBind(
		queue.Name,     // 队列名
		"",             // 路由参数，如果匹配消息发送的时候指定的路由参数，消息就投递到当前队列（在Topics模式下，这里的key要指定）
		"dlx_exchange", // 交换机名字，需要跟消息发送端定义的交换器保持一致
		false,          // 是否阻塞处理
		nil,            // 额外的属性
	)
	r.failOnErr(err, "Failed to bind a DlxQueue")

	//4.推送消息
	msgs, err := r.channel.Consume(
		queue.Name, // 队列名称
		"",         // 用来区分多个消费者
		false,      // 是否自动应答
		false,      // 是否独有
		false,      // 设置为true，表示不能将同一个Connection中生产者发送的消息传递给这个Connection中的消费者
		false,      // 队列是否阻塞
		nil,        // 额外的属性
	)
	r.failOnErr(err, "Failed to Consume")

	fmt.Println(" [*] Waiting for messages.")
	//5.消费接收到的消息
	for {
		select {
		case message, ok := <-msgs:
			if !ok {
				continue
			}
			go func() {
				//处理消息
				time.Sleep(time.Second * 3)
				//确认接收到的消息
				if err = message.Ack(true); err != nil {
					fmt.Println("dlx d.Ack err: ", err)
					return
				}
				fmt.Println("已确认dlx", string(message.Body))
			}()
		case <-time.After(time.Second * 1):

		}
	}

}
