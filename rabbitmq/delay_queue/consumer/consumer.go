package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	// 建立链接
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	var (
		//exchange   = "x-delayed-message"
		queue = "delay_queue"
		//routingKey = "log_delay"
	)

	// 申请交换机
	//err = ch.ExchangeDeclare(
	//	exchange, // name
	//	exchange, // type
	//	true,     // durable
	//	false,    // auto-deleted
	//	false,    // internal
	//	false,    // no-wait
	//	amqp.Table{
	//		"x-delayed-type": "direct",
	//	})
	//if err != nil {
	//	failOnError(err, "交换机申请失败！")
	//	return
	//}

	// 声明一个常规的队列, 其实这个也没必要声明,因为 exchange 会默认绑定一个队列
	q, err := ch.QueueDeclare(
		queue, // name
		true,  // durable
		true,  // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	//err = ch.QueueBind(
	//	q.Name,     // queue name
	//	routingKey, // routing key
	//	exchange,   // exchange
	//	false,
	//	nil)
	//failOnError(err, "Failed to bind a queue")

	// 这里监听的是 test_logs
	msgs, err := ch.Consume(
		q.Name, // queue name
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("接收数据 [x] %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
