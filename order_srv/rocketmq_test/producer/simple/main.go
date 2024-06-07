package main

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

//普通消息
func main() {
	p, err := rocketmq.NewProducer(producer.WithNameServer([]string{"192.168.0.101:9876"}))
	if err != nil {
		panic("初始化失败")
	}
	err = p.Start()
	if err != nil {
		panic("开始失败")
	}

	//msg := &primitive.Message{
	//	Topic: "mytopic",
	//	Body:  []byte("this is my second message"),
	//}

	msg := primitive.NewMessage("mytopic", []byte("this is my test message 5"))
	res, err := p.SendSync(context.Background(), msg)
	if err != nil {
		panic("发送失败")
	} else {
		fmt.Printf("发送成功:%v", res.String())
	}
	err = p.Shutdown()
	if err != nil {
		panic("关闭失败")
	}
}
