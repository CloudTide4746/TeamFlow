package mq

import (
	"github.com/rabbitmq/amqp091-go"
)

// initMQ 初始化 RabbitMQ 连接
func initMQ() *amqp091.Connection {
	// 连接到 RabbitMQ 服务器
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	return conn
}

// NewChannel 创建 RabbitMQ 通道
func NewChannel(conn *amqp091.Connection) *amqp091.Channel {
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}

func NewQueue(ch *amqp091.Channel, name string) amqp091.Queue {
	queue, err := ch.QueueDeclare(
		name,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,
	)
	if err != nil {
		panic(err)
	}
	return queue
}

// CloseChannel 关闭 RabbitMQ 通道
func CloseChannel(ch *amqp091.Channel) {
	ch.Close()
}

// CloseConnection 关闭 RabbitMQ 连接
func CloseConnection(conn *amqp091.Connection) {
	conn.Close()
}

// NewMessage 创建 RabbitMQ 消息
func NewMessage(body []byte) amqp091.Publishing {
	return amqp091.Publishing{
		Body: body,
	}
}

// Publish 发布 RabbitMQ 消息
func Publish(ch *amqp091.Channel, exchange, routingKey string, msg amqp091.Publishing) {
	ch.Publish(exchange, routingKey, false, false, msg)
}

// Consume 消费 RabbitMQ 消息
func Consume(ch *amqp091.Channel, queue string, noAck bool) <-chan amqp091.Delivery {
	msgs, err := ch.Consume(queue, "", noAck, false, false, false, nil)
	if err != nil {
		panic(err)
	}
	return msgs
}

// Ack 确认 RabbitMQ 消息
func Ack(ch *amqp091.Channel, delivery amqp091.Delivery) {
	delivery.Ack(false)
}
