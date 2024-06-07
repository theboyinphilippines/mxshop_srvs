package main

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"time"
)

//事务消息
type orderListener struct {
}

//业务
func (o *orderListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	return primitive.UnknowState
}

//回查
func (o *orderListener) CheckLocalTransaction(msgExt *primitive.MessageExt) primitive.LocalTransactionState {
	return primitive.CommitMessageState
}

func main() {
	p, err := rocketmq.NewTransactionProducer(&orderListener{}, producer.WithNameServer([]string{"192.168.0.101:9876"}))
	if err != nil {
		panic("初始化失败")
	}
	err = p.Start()
	if err != nil {
		panic("开始失败")
	}
	msg := primitive.NewMessage("transtopic", []byte("this is my transaction message"))
	res, err := p.SendMessageInTransaction(context.Background(), msg)
	if err != nil {
		panic("发送失败")
	} else {
		fmt.Printf("发送成功:%v", res.String())
	}
	//这里需要回查，不能立马Shutdown
	time.Sleep(time.Hour)
	err = p.Shutdown()
	if err != nil {
		panic("关闭失败")
	}
}
