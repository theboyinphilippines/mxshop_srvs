package main

import (
	"context"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"log"
)

// 直接使用第三方日志库，例如下面示例代码中使用了zap日志库。
func main() {
	// 创建一个writer 向zap-topic发送消息
	w := &kafka.Writer{
		Addr:                   kafka.TCP("localhost:9092"),
		Topic:                  "zap-topic",
		Balancer:               &kafka.LeastBytes{}, // 指定分区的balancer模式为最小字节分布
		RequiredAcks:           kafka.RequireAll,    // ack模式
		Async:                  true,                // 异步
		AllowAutoTopicCreation: true,
		Logger:                 kafka.LoggerFunc(zap.NewExample().Sugar().Infof),
		ErrorLogger:            kafka.LoggerFunc(zap.NewExample().Sugar().Errorf),
	}

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
			Key:   []byte("Key-C"),
			Value: []byte("嘻嘻哈哈的"),
		},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := w.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}

}
