package main

import (
	"log"
	"time"

	"github.com/streadway/amqp"
)

func main() {
	// 连接到 RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	// 创建通道并设置为手动确认模式
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	if err = ch.Confirm(false, true); err != nil {
		log.Fatalf("Could not enable publisher confirms: %v", err)
	}

	// 声明交换机和队列
	exchange := "my_exchange"
	queue := "my_queue"
	if err = ch.ExchangeDeclare(
		exchange, "direct", true, false, false, false, nil,
	); err != nil {
		log.Fatalf("Failed to declare an exchange: %v", err)
	}

	if _, err = ch.QueueDeclare(
		queue, true, false, false, false, nil,
	); err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	if err = ch.QueueBind(
		queue, "my_routing_key", exchange, false, nil,
	); err != nil {
		log.Fatalf("Failed to bind a queue: %v", err)
	}

	// 发布消息
	body := "Hello World!"
	err = ch.Publish(
		exchange, "my_routing_key",
		false, false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
	}

	// 消费消息
	msgs, err := ch.Consume(
		queue, "", true, false, false, false, nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	// 处理消息
	for d := range msgs {
		log.Printf("Received a message: %s", d.Body)
		// 模拟消息处理过程
		time.Sleep(1 * time.Second)

		// 手动确认消息
		if err := d.Ack(false); err != nil {
			log.Printf("Failed to ack message: %v", err)
		}
	}
}
