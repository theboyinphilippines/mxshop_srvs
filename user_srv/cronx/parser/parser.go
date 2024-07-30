package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
)

func main() {
	// 自定义时间解析器，下面可以加入秒
	parser := cron.NewParser(
		cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)
	c := cron.New(cron.WithParser(parser))
	//等同于下面
	//c := cron.New(cron.WithSeconds())
	_, _ = c.AddFunc("*/1 * * * * *", func() {
		fmt.Println("every 1 second")
	})
	c.Start()
	select {}
}
