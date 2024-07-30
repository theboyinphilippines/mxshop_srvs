package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"time"
)

// 测试下定时任务cron
func main() {
	//添加时区
	loc, _ := time.LoadLocation("Asia/Shanghai")
	c := cron.New(cron.WithLocation(loc))
	//添加任务
	_, _ = c.AddFunc("@every 1s", func() {
		fmt.Println("Every 1 second")
	})
	_, _ = c.AddFunc("*/1 * * * *", func() {
		fmt.Println("Every 1 minute")
	})
	c.Start()

	select {}
}
