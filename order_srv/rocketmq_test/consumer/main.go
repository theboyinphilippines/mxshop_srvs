package main

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"time"
)

func main() {
	c, err := rocketmq.NewPushConsumer(consumer.WithNameServer([]string{"192.168.0.101:9876"}), consumer.WithGroupName("test"))
	if err != nil {
		panic(any(err))
	}
	err = c.Subscribe("mytopic", consumer.MessageSelector{}, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for k, _ := range msgs {
			fmt.Printf("消费的消息是：%v", msgs[k])
		}
		//for _, msg := range msgs {
		//	fmt.Printf("消费的消息是：%v \n", msg )
		//}
		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		panic(any(err))
	}
	err = c.Start()
	if err != nil {
		panic(any(err))
	}
	time.Sleep(time.Hour)
	err = c.Shutdown()
}
