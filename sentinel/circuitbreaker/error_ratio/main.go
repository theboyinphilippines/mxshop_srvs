package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/logging"
	"github.com/alibaba/sentinel-golang/util"
)

type stateChangeTestListener struct {
}

func (s *stateChangeTestListener) OnTransformToClosed(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	fmt.Printf("rule.steategy: %+v, From %s to Closed, time: %d\n", rule.Strategy, prev.String(), util.CurrentTimeMillis())
}

func (s *stateChangeTestListener) OnTransformToOpen(prev circuitbreaker.State, rule circuitbreaker.Rule, snapshot interface{}) {
	fmt.Printf("rule.steategy: %+v, From %s to Open, snapshot: %.2f, time: %d\n", rule.Strategy, prev.String(), snapshot, util.CurrentTimeMillis())
}

func (s *stateChangeTestListener) OnTransformToHalfOpen(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	fmt.Printf("rule.steategy: %+v, From %s to Half-Open, time: %d\n", rule.Strategy, prev.String(), util.CurrentTimeMillis())
}

// 熔断器-错误统计
func main() {
	total := 0
	totalBlock := 0
	totalError := 0
	totalPass := 0
	conf := config.NewDefaultConfig()
	// for testing, logging output to console
	conf.Sentinel.Log.Logger = logging.NewConsoleLogger()
	err := sentinel.InitWithConfig(conf)
	if err != nil {
		log.Fatal(err)
	}
	ch := make(chan struct{})
	// Register a state change listener so that we could observer the state change of the internal circuit breaker.
	circuitbreaker.RegisterStateChangeListeners(&stateChangeTestListener{})

	_, err = circuitbreaker.LoadRules([]*circuitbreaker.Rule{
		// Statistic time span=5s, recoveryTimeout=3s
		{
			Resource:         "abc",
			Strategy:         circuitbreaker.ErrorRatio,
			RetryTimeoutMs:   3000, //3s之后尝试恢复
			MinRequestAmount: 10,   //静默数，10个请求以内全部通过
			StatIntervalMs:   5000, //5s统计一次
			Threshold:        0.4,  //错误数比率不超过40%
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	logging.Info("[CircuitBreaker ErrorCount] Sentinel Go circuit breaking demo is running. You may see the pass/block metric in the metric log.")
	go func() {
		for {
			total++
			e, b := sentinel.Entry("abc")
			if b != nil {
				// g1 blocked
				//rand.Uint64()%20)，20以内随机取一个数
				totalBlock++
				fmt.Println("携程熔断了")
				time.Sleep(time.Duration(rand.Uint64()%20) * time.Millisecond)
			} else {
				totalPass++
				if rand.Uint64()%20 > 9 {
					// Record current invocation as error.
					totalError++
					sentinel.TraceError(e, errors.New("biz error"))
				}
				// g1 passed
				time.Sleep(time.Duration(rand.Uint64()%20+10) * time.Millisecond)
				e.Exit()
			}
		}
	}()
	go func() {
		for {
			total++
			e, b := sentinel.Entry("abc")
			if b != nil {
				// g2 blocked
				totalBlock++
				time.Sleep(time.Duration(rand.Uint64()%20) * time.Millisecond)
			} else {
				// g2 passed
				totalPass++
				time.Sleep(time.Duration(rand.Uint64()%30) * time.Millisecond)
				e.Exit()
			}
		}
	}()

	go func() {
		for {
			time.Sleep(time.Second)
			fmt.Printf("totalError ratio: %v\n", float64(totalError)/float64(total))
		}
	}()
	<-ch
}
