package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// 通过死信队列和ttl实现延迟队列
func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	var (
		exchange   = "x-delayed-message"
		queue      = "delay_queue"
		routingKey = "log_delay"
		body       string
	)
	// 申请交换机
	err = ch.ExchangeDeclare(
		exchange,
		exchange, //交换机类型 x-delayed-message
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-delayed-type": "direct",
		})

	if err != nil {
		failOnError(err, "交换机申请失败！")
		return
	}
	if err = ch.QueueBind(queue, routingKey, exchange, false, nil); err != nil {
		failOnError(err, "绑定交换机失败！")
		return
	}

	body = "==========10000=================" + time.Now().Local().Format("2006-01-02 15:04:05")
	// 将消息发送到延时队列上
	err = ch.Publish(
		exchange,   // exchange 这里为空则不选择 exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
			Headers: map[string]interface{}{
				"x-delay": "10000", // 消息从交换机过期时间,毫秒（x-dead-message插件提供）
			},
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s", body)

	body = "==========20000=================" + time.Now().Local().Format("2006-01-02 15:04:05")
	// 将消息发送到延时队列上
	err = ch.Publish(
		exchange,   // exchange 这里为空则不选择 exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
			Headers: map[string]interface{}{
				"x-delay": "20000", // 消息从交换机过期时间,毫秒（x-dead-message插件提供）
			},
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s", body)

	body = "==========5000=================" + time.Now().Local().Format("2006-01-02 15:04:05")
	// 将消息发送到延时队列上
	err = ch.Publish(
		exchange,   // exchange 这里为空则不选择 exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
			Headers: map[string]interface{}{
				"x-delay": "5000", // 消息从交换机过期时间,毫秒（x-dead-message插件提供）
			},
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s", body)
}
