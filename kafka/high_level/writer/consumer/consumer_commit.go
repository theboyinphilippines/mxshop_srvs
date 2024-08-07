package main

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
)

// kafka-go 也支持显式提交（手动提交）。当需要显式提交时不要调用 ReadMessage，而是调用 FetchMessage获取消息，然后调用 CommitMessages 显式提交
func main() {
	// 创建一个reader，指定GroupID，从 topic-A 消费消息
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092", "localhost:9093", "localhost:9094"},
		GroupID:  "my-consumer-group3", // 指定消费者组id
		Topic:    "topic-B",
		MaxBytes: 10e6, // 10MB
		//CommitInterval: time.Second,
	})

	ctx := context.Background()
	for {
		// 获取消息
		m, err := r.FetchMessage(ctx)
		if err != nil {
			break
		}
		// 处理消息
		fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
		// 显式提交
		if err := r.CommitMessages(ctx, m); err != nil {
			log.Fatal("failed to commit messages:", err)
		}
	}

}
