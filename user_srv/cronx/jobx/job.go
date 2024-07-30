package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"time"
)

type GreetingJob struct {
	Name string
}

func (g GreetingJob) Run() {
	fmt.Println("Hello ", g.Name)
}

// 除了直接将无参函数作为回调外，cron还支持Job接口
func main() {
	//添加时区
	loc, _ := time.LoadLocation("Asia/Shanghai")
	c := cron.New(cron.WithLocation(loc))
	//添加任务
	_, _ = c.AddJob("@every 1s", GreetingJob{Name: "nihao"})
	_, _ = c.AddJob("@every 3.5s", GreetingJob{Name: "haha"})
	c.Start()
	select {}

}
