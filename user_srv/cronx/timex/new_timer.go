package timex

import (
	"fmt"
	"time"
)

//使用time标准库来做定时任务

func Timer() {
	t := time.NewTicker(3 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			job()
			t.Reset(5 * time.Second)
		}
	}
}

func job() {
	fmt.Println("执行定时具体任务....")
}
