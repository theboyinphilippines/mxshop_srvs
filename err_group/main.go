package main

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"time"
)

func main() {

	eg, ctx := errgroup.WithContext(context.Background())
	// 5秒后抛出错误，ctx会cancel
	eg.Go(func() error {
		fmt.Println("doing task 1")
		time.Sleep(5 * time.Second)
		return errors.New("taks 1 has error")
	})
	// task 1， 5秒后抛出错误，会cancel，这里要接收到cancel
	eg.Go(func() error {
		for {
			select {
			case <-time.After(1 * time.Second):
				//1秒后执行task 2
				fmt.Println("doing task 2")
			case <-ctx.Done():
				//ctx.Done() 接收cancel的通道
				return ctx.Err()
			}
		}
	})

	eg.Go(func() error {
		for {
			select {
			case <-time.After(1 * time.Second):
				//1秒后执行task 2
				fmt.Println("doing task 3")
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	})

	err := eg.Wait()
	if err != nil {
		fmt.Println("task falied")
	} else {
		fmt.Println("task success")
	}

}
