package main

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
)

// 对于从单个主题-分区（topic-partition）消费消息
// 使用reader消费消息
func main() {
	// 创建Reader
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:9092"},
		Topic:     "my-topic",
		Partition: 0,
		MaxBytes:  10e6, // 10MB
	})
	r.SetOffset(0) // 设置Offset

	// 接收消息
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	}

	// 程序退出前关闭Reader
	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}
