package main

import (
	"github.com/segmentio/kafka-go"
	"log"
	"net"
	"strconv"
)

// 创建topic
// to create topics when auto.create.topics.enable='false'
func main() {
	// 指定要创建的topic名称
	topic := "new-topic"

	// 连接至任意kafka节点（连接到非leader节点）
	conn, err := kafka.Dial("tcp", "localhost:9092")
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	// 获取当前控制节点信息
	controller, err := conn.Controller()
	log.Printf("controller：%v", controller)
	if err != nil {
		panic(err.Error())
	}
	var controllerConn *kafka.Conn
	// 连接至leader节点（通过非leader节点连接到leader节点）
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		panic(err.Error())
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}

	// 创建topic
	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		panic(err.Error())
	}
}
