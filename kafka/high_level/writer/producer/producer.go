package main

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
)

type user struct {
	Name   string `json:"name"`
	Mobile string `json:"mobile"`
	Gender int    `json:"gender"`
	Age    int    `json:"age"`
}

func main() {
	// 创建一个writer 向topic-A发送消息
	w := &kafka.Writer{
		Addr:         kafka.TCP("localhost:9092"),
		Topic:        "new-topic",
		Balancer:     &kafka.LeastBytes{}, // 指定分区的balancer模式为最小字节分布
		RequiredAcks: kafka.RequireAll,    // ack模式
		Async:        true,                // 异步
	}

	//var userinfo user
	var userinfo = user{
		Name:   "john",
		Mobile: "13700031234",
		Gender: 1,
		Age:    20,
	}
	userinfoByte, _ := json.Marshal(userinfo) // 定义模型字段时，要打json tag
	err := w.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte("Key-A"),
			Value: []byte("你好吗"),
		},
		kafka.Message{
			Key:   []byte("Key-B"),
			Value: []byte("真的是"),
		},
		kafka.Message{
			Key:   []byte(userinfo.Mobile),
			Value: userinfoByte,
		},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := w.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}

}
