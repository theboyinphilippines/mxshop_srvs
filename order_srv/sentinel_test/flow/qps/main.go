package main

import (
	"fmt"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/flow"
	"log"
)

func main() {
	err := sentinel.InitDefault()
	if err != nil {
		log.Fatalf("Unexpected error: %+v", err)
	}

	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "some-test",
			TokenCalculateStrategy: flow.Direct, //qps
			ControlBehavior:        flow.Reject, //直接拒绝
			Threshold:              500,
			StatIntervalInMs:       1000, //1秒最大请求并发500个
		},
		{
			Resource:               "some-test2",
			TokenCalculateStrategy: flow.Direct,     //限流qps
			ControlBehavior:        flow.Throttling, //匀速通过
			Threshold:              10,
			StatIntervalInMs:       1000, //1秒匀速通过10个请求
		},
	})
	if err != nil {
		log.Fatalf("Unexpected error: %+v", err)
		return
	}

	for i := 0; i < 10; i++ {
		e, b := sentinel.Entry("some-test", sentinel.WithTrafficType(base.Inbound))
		if b != nil {
			fmt.Println("限流了")
		} else {
			fmt.Println("通过了")
			e.Exit()
		}
	}
}
