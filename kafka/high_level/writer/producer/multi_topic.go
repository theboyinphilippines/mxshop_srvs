package main

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
)

// 分别设置每条消息的message.topic，可以实现将消息发送至多个topic
func main() {
	w := &kafka.Writer{
		Addr: kafka.TCP("localhost:9092"),
		// 注意: 当此处不设置Topic时,后续的每条消息都需要指定Topic
		Balancer: &kafka.LeastBytes{},
	}

	err := w.WriteMessages(context.Background(),
		// 注意: 每条消息都需要指定一个 Topic, 否则就会报错
		kafka.Message{
			Topic: "topic-A",
			Key:   []byte("Key-A"),
			Value: []byte("Hello World!"),
		},
		kafka.Message{
			Topic: "topic-B",
			Key:   []byte("Key-B"),
			Value: []byte("One!"),
		},
		kafka.Message{
			Topic: "topic-C",
			Key:   []byte("Key-C"),
			Value: []byte("Two!"),
		},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := w.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}

}
